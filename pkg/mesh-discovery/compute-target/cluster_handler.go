package compute_target

import (
	"context"

	k8s_apps "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1"
	k8s_apps_controller "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1/controller"
	k8s_core "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	k8s_core_controller "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/controller"
	smh_discovery "github.com/solo-io/service-mesh-hub/pkg/api/discovery.smh.solo.io/v1alpha1"
	smh_discovery_controller "github.com/solo-io/service-mesh-hub/pkg/api/discovery.smh.solo.io/v1alpha1/controller"
	mc_manager "github.com/solo-io/service-mesh-hub/pkg/common/compute-target/k8s"
	k8s_tenancy "github.com/solo-io/service-mesh-hub/pkg/mesh-discovery/discovery/cluster-tenancy/k8s"
	meshworkload_discovery "github.com/solo-io/service-mesh-hub/pkg/mesh-discovery/discovery/mesh-workload/k8s"
	"github.com/solo-io/service-mesh-hub/pkg/mesh-discovery/discovery/mesh/k8s"
	"github.com/solo-io/service-mesh-hub/pkg/mesh-discovery/wire"
)

// this is the main entrypoint for discovery
// when a cluster is registered, we handle that event and spin up new resource controllers for that cluster
func NewDiscoveryClusterHandler(
	localManager mc_manager.AsyncManager,
	meshScanners []k8s.MeshScanner,
	clusterTenancyScannerFactories []k8s_tenancy.ClusterTenancyScannerFactory,
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
		localMeshClient:                localMeshClient,
		meshScanners:                   meshScanners,
		localMeshWorkloadClient:        localMeshWorkloadClient,
		localManager:                   localManager,
		discoveryContext:               discoveryContext,
		localMeshServiceClient:         localMeshServiceClient,
		localMeshWorkloadEventWatcher:  localMeshWorkloadEventWatcher,
		localMeshEventWatcher:          localMeshController,
		clusterTenancyScannerFactories: clusterTenancyScannerFactories,
	}

	return handler, nil
}

type discoveryClusterHandler struct {
	localManager     mc_manager.AsyncManager
	discoveryContext wire.DiscoveryContext

	// clients that operate against the local cluster
	localMeshClient         smh_discovery.MeshClient
	localMeshWorkloadClient smh_discovery.MeshWorkloadClient
	localMeshServiceClient  smh_discovery.MeshServiceClient

	// controllers that operate against the local cluster
	localMeshWorkloadEventWatcher smh_discovery_controller.MeshWorkloadEventWatcher
	localMeshEventWatcher         smh_discovery_controller.MeshEventWatcher

	// scanners
	meshScanners                   []k8s.MeshScanner
	clusterTenancyScannerFactories []k8s_tenancy.ClusterTenancyScannerFactory
}

type clusterDependentDeps struct {
	deploymentEventWatcher k8s_apps_controller.DeploymentEventWatcher
	podEventWatcher        k8s_core_controller.PodEventWatcher
	meshWorkloadScanners   meshworkload_discovery.MeshWorkloadScanners
	clusterTenancyScanners []k8s_tenancy.ClusterTenancyRegistrar
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

	clusterTenancyFinder := k8s_tenancy.NewClusterTenancyFinder(
		clusterName,
		initializedDeps.clusterTenancyScanners,
		initializedDeps.podClient,
		m.localMeshClient,
	)

	if err = meshFinder.StartDiscovery(initializedDeps.deploymentEventWatcher); err != nil {
		return err
	}

	if err = clusterTenancyFinder.StartRegistration(ctx, initializedDeps.podEventWatcher, m.localMeshEventWatcher); err != nil {
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

	var clusterTenancyScanners []k8s_tenancy.ClusterTenancyRegistrar
	//for _, tenancyScannerFactory := range m.clusterTenancyScannerFactories {
	//	clusterTenancyScanners = append(clusterTenancyScanners, tenancyScannerFactory(m.localMeshClient))
	//}

	return &clusterDependentDeps{
		deploymentEventWatcher: deploymentEventWatcher,
		podEventWatcher:        podEventWatcher,
		serviceEventWatcher:    serviceEventWatcher,
		serviceClient:          serviceClient,
		podClient:              podClient,
		deploymentClient:       deploymentClient,
		clusterTenancyScanners: clusterTenancyScanners,
	}, nil
}
