package mesh_platform

import (
	"context"

	core_types "github.com/solo-io/service-mesh-hub/pkg/api/core.zephyr.solo.io/v1alpha1/types"
	zephyr_discovery "github.com/solo-io/service-mesh-hub/pkg/api/discovery.zephyr.solo.io/v1alpha1"
	zephyr_discovery_controller "github.com/solo-io/service-mesh-hub/pkg/api/discovery.zephyr.solo.io/v1alpha1/controller"
	k8s_apps "github.com/solo-io/service-mesh-hub/pkg/api/kubernetes/apps/v1"
	k8s_apps_controller "github.com/solo-io/service-mesh-hub/pkg/api/kubernetes/apps/v1/controller"
	k8s_core "github.com/solo-io/service-mesh-hub/pkg/api/kubernetes/core/v1"
	k8s_core_controller "github.com/solo-io/service-mesh-hub/pkg/api/kubernetes/core/v1/controller"
	"github.com/solo-io/service-mesh-hub/pkg/env"
	mc_manager "github.com/solo-io/service-mesh-hub/services/common/mesh-platform/k8s"
	meshservice_discovery "github.com/solo-io/service-mesh-hub/services/mesh-discovery/pkg/discovery/mesh-service/k8s"
	meshworkload_discovery "github.com/solo-io/service-mesh-hub/services/mesh-discovery/pkg/discovery/mesh-workload/k8s"
	"github.com/solo-io/service-mesh-hub/services/mesh-discovery/pkg/discovery/mesh/k8s"
	"github.com/solo-io/service-mesh-hub/services/mesh-discovery/pkg/wire"
)

type MeshWorkloadScannerFactoryImplementations map[core_types.MeshType]meshworkload_discovery.MeshWorkloadScannerFactory

// this is the main entrypoint for discovery
// when a cluster is registered, we handle that event and spin up new resource controllers for that cluster
func NewDiscoveryClusterHandler(
	localManager mc_manager.AsyncManager,
	meshScanners []k8s.MeshScanner,
	meshWorkloadScannerFactories MeshWorkloadScannerFactoryImplementations,
	discoveryContext wire.DiscoveryContext,
) (mc_manager.AsyncManagerHandler, error) {

	// these clients operate against the local cluster, so we use the local manager's client
	localClient := localManager.Manager().GetClient()
	localMeshServiceClient := discoveryContext.ClientFactories.MeshServiceClientFactory(localClient)
	localMeshWorkloadClient := discoveryContext.ClientFactories.MeshWorkloadClientFactory(localClient)
	localMeshClient := discoveryContext.ClientFactories.MeshClientFactory(localClient)

	localMeshWorkloadEventWatcher := discoveryContext.EventWatcherFactories.MeshWorkloadEventWatcherFactory.Build(localManager, "mesh-workload-apps_controller")

	localMeshController := discoveryContext.EventWatcherFactories.MeshControllerFactory.Build(localManager, "mesh-controller")

	// we don't store the local manager on the struct to avoid mistakenly conflating the local manager with the remote manager
	handler := &discoveryClusterHandler{
		localMeshClient:               localMeshClient,
		meshScanners:                  meshScanners,
		localMeshWorkloadClient:       localMeshWorkloadClient,
		localManager:                  localManager,
		meshWorkloadScannerFactories:  meshWorkloadScannerFactories,
		discoveryContext:              discoveryContext,
		localMeshServiceClient:        localMeshServiceClient,
		localMeshWorkloadEventWatcher: localMeshWorkloadEventWatcher,
		localMeshEventWatcher:         localMeshController,
	}

	return handler, nil
}

type discoveryClusterHandler struct {
	localManager     mc_manager.AsyncManager
	discoveryContext wire.DiscoveryContext

	// clients that operate against the local cluster
	localMeshClient         zephyr_discovery.MeshClient
	localMeshWorkloadClient zephyr_discovery.MeshWorkloadClient
	localMeshServiceClient  zephyr_discovery.MeshServiceClient

	// controllers that operate against the local cluster
	localMeshWorkloadEventWatcher zephyr_discovery_controller.MeshWorkloadEventWatcher
	localMeshEventWatcher         zephyr_discovery_controller.MeshEventWatcher

	// scanners
	meshScanners                 []k8s.MeshScanner
	meshWorkloadScannerFactories MeshWorkloadScannerFactoryImplementations
}

type clusterDependentDeps struct {
	deploymentEventWatcher k8s_apps_controller.DeploymentEventWatcher
	podEventWatcher        k8s_core_controller.PodEventWatcher
	meshWorkloadScanners   meshworkload_discovery.MeshWorkloadScannerImplementations
	serviceEventWatcher    k8s_core_controller.ServiceEventWatcher
	serviceClient          k8s_core.ServiceClient
	podClient              k8s_core.PodClient
	deploymentClient       k8s_apps.DeploymentClient
}

func (m *discoveryClusterHandler) ClusterAdded(ctx context.Context, mgr mc_manager.AsyncManager, clusterName string) error {
	initializedDeps, err := m.initializeClusterDependentDeps(mgr, clusterName)
	if err != nil {
		return err
	}
	meshFinder := k8s.NewMeshFinder(
		ctx,
		clusterName,
		m.meshScanners,
		m.localMeshClient,
		mgr.Manager().GetClient(),
		initializedDeps.deploymentClient,
	)

	meshWorkloadFinder := meshworkload_discovery.NewMeshWorkloadFinder(
		ctx,
		clusterName,
		m.localMeshClient,
		m.localMeshWorkloadClient,
		initializedDeps.meshWorkloadScanners,
		initializedDeps.podClient,
	)

	meshServiceFinder := meshservice_discovery.NewMeshServiceFinder(
		ctx,
		clusterName,
		env.GetWriteNamespace(),
		initializedDeps.serviceClient,
		m.localMeshServiceClient,
		m.localMeshWorkloadClient,
		m.localMeshClient,
	)

	err = meshFinder.StartDiscovery(initializedDeps.deploymentEventWatcher)
	if err != nil {
		return err
	}

	err = meshWorkloadFinder.StartDiscovery(initializedDeps.podEventWatcher, m.localMeshEventWatcher)
	if err != nil {
		return err
	}

	err = meshServiceFinder.StartDiscovery(initializedDeps.serviceEventWatcher, m.localMeshWorkloadEventWatcher)
	if err != nil {
		return err
	}

	return nil
}

func (m *discoveryClusterHandler) ClusterRemoved(cluster string) error {
	// TODO: Not deleting any entities for now
	return nil
}

func (m *discoveryClusterHandler) initializeClusterDependentDeps(mgr mc_manager.AsyncManager, clusterName string) (*clusterDependentDeps, error) {
	deploymentEventWatcher := m.discoveryContext.EventWatcherFactories.DeploymentEventWatcherFactory.Build(mgr, clusterName)
	podEventWatcher := m.discoveryContext.EventWatcherFactories.PodEventWatcherFactory.Build(mgr, clusterName)
	serviceEventWatcher := m.discoveryContext.EventWatcherFactories.ServiceEventWatcherFactory.Build(mgr, clusterName)

	remoteClient := mgr.Manager().GetClient()

	serviceClient := m.discoveryContext.ClientFactories.ServiceClientFactory(remoteClient)
	podClient := m.discoveryContext.ClientFactories.PodClientFactory(remoteClient)
	deploymentClient := m.discoveryContext.ClientFactories.DeploymentClientFactory(remoteClient)
	replicaSetClient := m.discoveryContext.ClientFactories.ReplicaSetClientFactory(remoteClient)

	meshWorkloadScanners := make(meshworkload_discovery.MeshWorkloadScannerImplementations)
	for meshType, scannerFactory := range m.meshWorkloadScannerFactories {
		ownerFetcher := m.discoveryContext.ClientFactories.OwnerFetcherClientFactory(
			deploymentClient,
			replicaSetClient,
		)

		meshWorkloadScanners[meshType] = scannerFactory(ownerFetcher, m.localMeshClient)
	}

	return &clusterDependentDeps{
		deploymentEventWatcher: deploymentEventWatcher,
		podEventWatcher:        podEventWatcher,
		meshWorkloadScanners:   meshWorkloadScanners,
		serviceEventWatcher:    serviceEventWatcher,
		serviceClient:          serviceClient,
		podClient:              podClient,
		deploymentClient:       deploymentClient,
	}, nil
}