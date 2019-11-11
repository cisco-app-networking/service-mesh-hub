// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"fmt"

	"github.com/solo-io/go-utils/hashutils"
	"go.uber.org/zap"
)

type NetworkBridgeSnapshot struct {
	MeshBridges   MeshBridgeList
	Meshes        MeshList
	MeshIngresses MeshIngressList
}

func (s NetworkBridgeSnapshot) Clone() NetworkBridgeSnapshot {
	return NetworkBridgeSnapshot{
		MeshBridges:   s.MeshBridges.Clone(),
		Meshes:        s.Meshes.Clone(),
		MeshIngresses: s.MeshIngresses.Clone(),
	}
}

func (s NetworkBridgeSnapshot) Hash() uint64 {
	return hashutils.HashAll(
		s.hashMeshBridges(),
		s.hashMeshes(),
		s.hashMeshIngresses(),
	)
}

func (s NetworkBridgeSnapshot) hashMeshBridges() uint64 {
	return hashutils.HashAll(s.MeshBridges.AsInterfaces()...)
}

func (s NetworkBridgeSnapshot) hashMeshes() uint64 {
	return hashutils.HashAll(s.Meshes.AsInterfaces()...)
}

func (s NetworkBridgeSnapshot) hashMeshIngresses() uint64 {
	return hashutils.HashAll(s.MeshIngresses.AsInterfaces()...)
}

func (s NetworkBridgeSnapshot) HashFields() []zap.Field {
	var fields []zap.Field
	fields = append(fields, zap.Uint64("meshBridges", s.hashMeshBridges()))
	fields = append(fields, zap.Uint64("meshes", s.hashMeshes()))
	fields = append(fields, zap.Uint64("meshIngresses", s.hashMeshIngresses()))

	return append(fields, zap.Uint64("snapshotHash", s.Hash()))
}

type NetworkBridgeSnapshotStringer struct {
	Version       uint64
	MeshBridges   []string
	Meshes        []string
	MeshIngresses []string
}

func (ss NetworkBridgeSnapshotStringer) String() string {
	s := fmt.Sprintf("NetworkBridgeSnapshot %v\n", ss.Version)

	s += fmt.Sprintf("  MeshBridges %v\n", len(ss.MeshBridges))
	for _, name := range ss.MeshBridges {
		s += fmt.Sprintf("    %v\n", name)
	}

	s += fmt.Sprintf("  Meshes %v\n", len(ss.Meshes))
	for _, name := range ss.Meshes {
		s += fmt.Sprintf("    %v\n", name)
	}

	s += fmt.Sprintf("  MeshIngresses %v\n", len(ss.MeshIngresses))
	for _, name := range ss.MeshIngresses {
		s += fmt.Sprintf("    %v\n", name)
	}

	return s
}

func (s NetworkBridgeSnapshot) Stringer() NetworkBridgeSnapshotStringer {
	return NetworkBridgeSnapshotStringer{
		Version:       s.Hash(),
		MeshBridges:   s.MeshBridges.NamespacesDotNames(),
		Meshes:        s.Meshes.NamespacesDotNames(),
		MeshIngresses: s.MeshIngresses.NamespacesDotNames(),
	}
}
