// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sync"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	skstats "github.com/solo-io/solo-kit/pkg/stats"

	"github.com/solo-io/go-utils/errutils"
)

var (
	// Deprecated. See mRbacResourcesIn
	mRbacSnapshotIn = stats.Int64("rbac.zephyr.solo.io/emitter/snap_in", "Deprecated. Use rbac.zephyr.solo.io/emitter/resources_in. The number of snapshots in", "1")

	// metrics for emitter
	mRbacResourcesIn    = stats.Int64("rbac.zephyr.solo.io/emitter/resources_in", "The number of resource lists received on open watch channels", "1")
	mRbacSnapshotOut    = stats.Int64("rbac.zephyr.solo.io/emitter/snap_out", "The number of snapshots out", "1")
	mRbacSnapshotMissed = stats.Int64("rbac.zephyr.solo.io/emitter/snap_missed", "The number of snapshots missed", "1")

	// views for emitter
	// deprecated: see rbacResourcesInView
	rbacsnapshotInView = &view.View{
		Name:        "rbac.zephyr.solo.io/emitter/snap_in",
		Measure:     mRbacSnapshotIn,
		Description: "Deprecated. Use rbac.zephyr.solo.io/emitter/resources_in. The number of snapshots updates coming in.",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}

	rbacResourcesInView = &view.View{
		Name:        "rbac.zephyr.solo.io/emitter/resources_in",
		Measure:     mRbacResourcesIn,
		Description: "The number of resource lists received on open watch channels",
		Aggregation: view.Count(),
		TagKeys: []tag.Key{
			skstats.NamespaceKey,
			skstats.ResourceKey,
		},
	}
	rbacsnapshotOutView = &view.View{
		Name:        "rbac.zephyr.solo.io/emitter/snap_out",
		Measure:     mRbacSnapshotOut,
		Description: "The number of snapshots updates going out",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
	rbacsnapshotMissedView = &view.View{
		Name:        "rbac.zephyr.solo.io/emitter/snap_missed",
		Measure:     mRbacSnapshotMissed,
		Description: "The number of snapshots updates going missed. this can happen in heavy load. missed snapshot will be re-tried after a second.",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	view.Register(
		rbacsnapshotInView,
		rbacsnapshotOutView,
		rbacsnapshotMissedView,
		rbacResourcesInView,
	)
}

type RbacSnapshotEmitter interface {
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *RbacSnapshot, <-chan error, error)
}

type RbacEmitter interface {
	RbacSnapshotEmitter
	Register() error
	Mesh() MeshClient
}

func NewRbacEmitter(meshClient MeshClient) RbacEmitter {
	return NewRbacEmitterWithEmit(meshClient, make(chan struct{}))
}

func NewRbacEmitterWithEmit(meshClient MeshClient, emit <-chan struct{}) RbacEmitter {
	return &rbacEmitter{
		mesh:      meshClient,
		forceEmit: emit,
	}
}

type rbacEmitter struct {
	forceEmit <-chan struct{}
	mesh      MeshClient
}

func (c *rbacEmitter) Register() error {
	if err := c.mesh.Register(); err != nil {
		return err
	}
	return nil
}

func (c *rbacEmitter) Mesh() MeshClient {
	return c.mesh
}

func (c *rbacEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *RbacSnapshot, <-chan error, error) {

	if len(watchNamespaces) == 0 {
		watchNamespaces = []string{""}
	}

	for _, ns := range watchNamespaces {
		if ns == "" && len(watchNamespaces) > 1 {
			return nil, nil, errors.Errorf("the \"\" namespace is used to watch all namespaces. Snapshots can either be tracked for " +
				"specific namespaces or \"\" AllNamespaces, but not both.")
		}
	}

	errs := make(chan error)
	var done sync.WaitGroup
	ctx := opts.Ctx
	/* Create channel for Mesh */
	type meshListWithNamespace struct {
		list      MeshList
		namespace string
	}
	meshChan := make(chan meshListWithNamespace)

	var initialMeshList MeshList

	currentSnapshot := RbacSnapshot{}

	for _, namespace := range watchNamespaces {
		/* Setup namespaced watch for Mesh */
		{
			meshes, err := c.mesh.List(namespace, clients.ListOpts{Ctx: opts.Ctx, Selector: opts.Selector})
			if err != nil {
				return nil, nil, errors.Wrapf(err, "initial Mesh list")
			}
			initialMeshList = append(initialMeshList, meshes...)
		}
		meshNamespacesChan, meshErrs, err := c.mesh.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Mesh watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, meshErrs, namespace+"-meshes")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-ctx.Done():
					return
				case meshList := <-meshNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case meshChan <- meshListWithNamespace{list: meshList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}
	/* Initialize snapshot for Meshes */
	currentSnapshot.Meshes = initialMeshList.Sort()

	snapshots := make(chan *RbacSnapshot)
	go func() {
		// sent initial snapshot to kick off the watch
		initialSnapshot := currentSnapshot.Clone()
		snapshots <- &initialSnapshot

		timer := time.NewTicker(time.Second * 1)
		previousHash := currentSnapshot.Hash()
		sync := func() {
			currentHash := currentSnapshot.Hash()
			if previousHash == currentHash {
				return
			}

			sentSnapshot := currentSnapshot.Clone()
			select {
			case snapshots <- &sentSnapshot:
				stats.Record(ctx, mRbacSnapshotOut.M(1))
				previousHash = currentHash
			default:
				stats.Record(ctx, mRbacSnapshotMissed.M(1))
			}
		}
		meshesByNamespace := make(map[string]MeshList)

		for {
			record := func() { stats.Record(ctx, mRbacSnapshotIn.M(1)) }

			select {
			case <-timer.C:
				sync()
			case <-ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case meshNamespacedList := <-meshChan:
				record()

				namespace := meshNamespacedList.namespace

				skstats.IncrementResourceCount(
					ctx,
					namespace,
					"mesh",
					mRbacResourcesIn,
				)

				// merge lists by namespace
				meshesByNamespace[namespace] = meshNamespacedList.list
				var meshList MeshList
				for _, meshes := range meshesByNamespace {
					meshList = append(meshList, meshes...)
				}
				currentSnapshot.Meshes = meshList.Sort()
			}
		}
	}()
	return snapshots, errs, nil
}
