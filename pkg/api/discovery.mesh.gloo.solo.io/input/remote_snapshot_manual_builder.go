// Code generated by skv2. DO NOT EDIT.

/*
	Utility for manually building input snapshots. Used primarily in tests.
*/
package input

import (
	appmesh_k8s_aws_v1beta2 "github.com/aws/aws-app-mesh-controller-for-k8s/apis/appmesh/v1beta2"
	appmesh_k8s_aws_v1beta2_sets "github.com/solo-io/external-apis/pkg/api/appmesh/appmesh.k8s.aws/v1beta2/sets"

	v1_sets "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/sets"
	v1 "k8s.io/api/core/v1"

	apps_v1_sets "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1/sets"
	apps_v1 "k8s.io/api/apps/v1"
)

type InputDiscoveryInputSnapshotManualBuilder struct {
	name string

	meshes appmesh_k8s_aws_v1beta2_sets.MeshSet

	configMaps v1_sets.ConfigMapSet
	services   v1_sets.ServiceSet
	pods       v1_sets.PodSet
	endpoints  v1_sets.EndpointsSet
	nodes      v1_sets.NodeSet

	deployments  apps_v1_sets.DeploymentSet
	replicaSets  apps_v1_sets.ReplicaSetSet
	daemonSets   apps_v1_sets.DaemonSetSet
	statefulSets apps_v1_sets.StatefulSetSet
}

func NewInputDiscoveryInputSnapshotManualBuilder(name string) *InputDiscoveryInputSnapshotManualBuilder {
	return &InputDiscoveryInputSnapshotManualBuilder{
		name: name,

		meshes: appmesh_k8s_aws_v1beta2_sets.NewMeshSet(),

		configMaps: v1_sets.NewConfigMapSet(),
		services:   v1_sets.NewServiceSet(),
		pods:       v1_sets.NewPodSet(),
		endpoints:  v1_sets.NewEndpointsSet(),
		nodes:      v1_sets.NewNodeSet(),

		deployments:  apps_v1_sets.NewDeploymentSet(),
		replicaSets:  apps_v1_sets.NewReplicaSetSet(),
		daemonSets:   apps_v1_sets.NewDaemonSetSet(),
		statefulSets: apps_v1_sets.NewStatefulSetSet(),
	}
}

func (i *InputDiscoveryInputSnapshotManualBuilder) Build() DiscoveryInputSnapshot {
	return NewDiscoveryInputSnapshot(
		i.name,

		i.meshes,

		i.configMaps,
		i.services,
		i.pods,
		i.endpoints,
		i.nodes,

		i.deployments,
		i.replicaSets,
		i.daemonSets,
		i.statefulSets,
	)
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddMeshes(meshes []*appmesh_k8s_aws_v1beta2.Mesh) *InputDiscoveryInputSnapshotManualBuilder {
	i.meshes.Insert(meshes...)
	return i
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddConfigMaps(configMaps []*v1.ConfigMap) *InputDiscoveryInputSnapshotManualBuilder {
	i.configMaps.Insert(configMaps...)
	return i
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddServices(services []*v1.Service) *InputDiscoveryInputSnapshotManualBuilder {
	i.services.Insert(services...)
	return i
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddPods(pods []*v1.Pod) *InputDiscoveryInputSnapshotManualBuilder {
	i.pods.Insert(pods...)
	return i
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddEndpoints(endpoints []*v1.Endpoints) *InputDiscoveryInputSnapshotManualBuilder {
	i.endpoints.Insert(endpoints...)
	return i
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddNodes(nodes []*v1.Node) *InputDiscoveryInputSnapshotManualBuilder {
	i.nodes.Insert(nodes...)
	return i
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddDeployments(deployments []*apps_v1.Deployment) *InputDiscoveryInputSnapshotManualBuilder {
	i.deployments.Insert(deployments...)
	return i
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddReplicaSets(replicaSets []*apps_v1.ReplicaSet) *InputDiscoveryInputSnapshotManualBuilder {
	i.replicaSets.Insert(replicaSets...)
	return i
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddDaemonSets(daemonSets []*apps_v1.DaemonSet) *InputDiscoveryInputSnapshotManualBuilder {
	i.daemonSets.Insert(daemonSets...)
	return i
}
func (i *InputDiscoveryInputSnapshotManualBuilder) AddStatefulSets(statefulSets []*apps_v1.StatefulSet) *InputDiscoveryInputSnapshotManualBuilder {
	i.statefulSets.Insert(statefulSets...)
	return i
}