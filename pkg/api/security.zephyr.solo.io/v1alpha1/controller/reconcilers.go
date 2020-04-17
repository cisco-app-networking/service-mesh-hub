// Definitions for the Kubernetes Controllers
package controller

import (
	"context"

	security_zephyr_solo_io_v1alpha1 "github.com/solo-io/service-mesh-hub/pkg/api/security.zephyr.solo.io/v1alpha1"

	"github.com/pkg/errors"
	"github.com/solo-io/skv2/pkg/ezkube"
	"github.com/solo-io/skv2/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// Reconcile Upsert events for the VirtualMeshCertificateSigningRequest Resource.
// implemented by the user
type VirtualMeshCertificateSigningRequestReconciler interface {
	ReconcileVirtualMeshCertificateSigningRequest(obj *security_zephyr_solo_io_v1alpha1.VirtualMeshCertificateSigningRequest) (reconcile.Result, error)
}

// Reconcile deletion events for the VirtualMeshCertificateSigningRequest Resource.
// Deletion receives a reconcile.Request as we cannot guarantee the last state of the object
// before being deleted.
// implemented by the user
type VirtualMeshCertificateSigningRequestDeletionReconciler interface {
	ReconcileVirtualMeshCertificateSigningRequestDeletion(req reconcile.Request)
}

type VirtualMeshCertificateSigningRequestReconcilerFuncs struct {
	OnReconcileVirtualMeshCertificateSigningRequest         func(obj *security_zephyr_solo_io_v1alpha1.VirtualMeshCertificateSigningRequest) (reconcile.Result, error)
	OnReconcileVirtualMeshCertificateSigningRequestDeletion func(req reconcile.Request)
}

func (f *VirtualMeshCertificateSigningRequestReconcilerFuncs) ReconcileVirtualMeshCertificateSigningRequest(obj *security_zephyr_solo_io_v1alpha1.VirtualMeshCertificateSigningRequest) (reconcile.Result, error) {
	if f.OnReconcileVirtualMeshCertificateSigningRequest == nil {
		return reconcile.Result{}, nil
	}
	return f.OnReconcileVirtualMeshCertificateSigningRequest(obj)
}

func (f *VirtualMeshCertificateSigningRequestReconcilerFuncs) ReconcileVirtualMeshCertificateSigningRequestDeletion(req reconcile.Request) {
	if f.OnReconcileVirtualMeshCertificateSigningRequestDeletion == nil {
		return
	}
	f.OnReconcileVirtualMeshCertificateSigningRequestDeletion(req)
}

// Reconcile and finalize the VirtualMeshCertificateSigningRequest Resource
// implemented by the user
type VirtualMeshCertificateSigningRequestFinalizer interface {
	VirtualMeshCertificateSigningRequestReconciler

	// name of the finalizer used by this handler.
	// finalizer names should be unique for a single task
	VirtualMeshCertificateSigningRequestFinalizerName() string

	// finalize the object before it is deleted.
	// Watchers created with a finalizing handler will a
	FinalizeVirtualMeshCertificateSigningRequest(obj *security_zephyr_solo_io_v1alpha1.VirtualMeshCertificateSigningRequest) error
}

type VirtualMeshCertificateSigningRequestReconcileLoop interface {
	RunVirtualMeshCertificateSigningRequestReconciler(ctx context.Context, rec VirtualMeshCertificateSigningRequestReconciler, predicates ...predicate.Predicate) error
}

type virtualMeshCertificateSigningRequestReconcileLoop struct {
	loop reconcile.Loop
}

func NewVirtualMeshCertificateSigningRequestReconcileLoop(name string, mgr manager.Manager) VirtualMeshCertificateSigningRequestReconcileLoop {
	return &virtualMeshCertificateSigningRequestReconcileLoop{
		loop: reconcile.NewLoop(name, mgr, &security_zephyr_solo_io_v1alpha1.VirtualMeshCertificateSigningRequest{}),
	}
}

func (c *virtualMeshCertificateSigningRequestReconcileLoop) RunVirtualMeshCertificateSigningRequestReconciler(ctx context.Context, reconciler VirtualMeshCertificateSigningRequestReconciler, predicates ...predicate.Predicate) error {
	genericReconciler := genericVirtualMeshCertificateSigningRequestReconciler{
		reconciler: reconciler,
	}

	var reconcilerWrapper reconcile.Reconciler
	if finalizingReconciler, ok := reconciler.(VirtualMeshCertificateSigningRequestFinalizer); ok {
		reconcilerWrapper = genericVirtualMeshCertificateSigningRequestFinalizer{
			genericVirtualMeshCertificateSigningRequestReconciler: genericReconciler,
			finalizingReconciler: finalizingReconciler,
		}
	} else {
		reconcilerWrapper = genericReconciler
	}
	return c.loop.RunReconciler(ctx, reconcilerWrapper, predicates...)
}

// genericVirtualMeshCertificateSigningRequestHandler implements a generic reconcile.Reconciler
type genericVirtualMeshCertificateSigningRequestReconciler struct {
	reconciler VirtualMeshCertificateSigningRequestReconciler
}

func (r genericVirtualMeshCertificateSigningRequestReconciler) Reconcile(object ezkube.Object) (reconcile.Result, error) {
	obj, ok := object.(*security_zephyr_solo_io_v1alpha1.VirtualMeshCertificateSigningRequest)
	if !ok {
		return reconcile.Result{}, errors.Errorf("internal error: VirtualMeshCertificateSigningRequest handler received event for %T", object)
	}
	return r.reconciler.ReconcileVirtualMeshCertificateSigningRequest(obj)
}

func (r genericVirtualMeshCertificateSigningRequestReconciler) ReconcileDeletion(request reconcile.Request) {
	if deletionReconciler, ok := r.reconciler.(VirtualMeshCertificateSigningRequestDeletionReconciler); ok {
		deletionReconciler.ReconcileVirtualMeshCertificateSigningRequestDeletion(request)
	}
}

// genericVirtualMeshCertificateSigningRequestFinalizer implements a generic reconcile.FinalizingReconciler
type genericVirtualMeshCertificateSigningRequestFinalizer struct {
	genericVirtualMeshCertificateSigningRequestReconciler
	finalizingReconciler VirtualMeshCertificateSigningRequestFinalizer
}

func (r genericVirtualMeshCertificateSigningRequestFinalizer) FinalizerName() string {
	return r.finalizingReconciler.VirtualMeshCertificateSigningRequestFinalizerName()
}

func (r genericVirtualMeshCertificateSigningRequestFinalizer) Finalize(object ezkube.Object) error {
	obj, ok := object.(*security_zephyr_solo_io_v1alpha1.VirtualMeshCertificateSigningRequest)
	if !ok {
		return errors.Errorf("internal error: VirtualMeshCertificateSigningRequest handler received event for %T", object)
	}
	return r.finalizingReconciler.FinalizeVirtualMeshCertificateSigningRequest(obj)
}
