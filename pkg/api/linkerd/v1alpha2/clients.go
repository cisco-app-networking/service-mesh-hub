package v1alpha2

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/linkerd/linkerd2/controller/gen/apis/serviceprofile/v1alpha2"
)

// clienset for the /v1alpha2 APIs
type Clientset interface {
	// clienset for the v1alpha2/v1alpha2 APIs
	ServiceProfiles() ServiceProfileClient
}

type clientSet struct {
	client client.Client
}

func NewClientsetFromConfig(cfg *rest.Config) (*clientSet, error) {
	scheme := scheme.Scheme
	if err := AddToScheme(scheme); err != nil {
		return nil, err
	}
	client, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	return NewClientset(client), nil
}

func NewClientset(client client.Client) *clientSet {
	return &clientSet{client: client}
}

// clienset for the v1alpha2/v1alpha2 APIs
func (c *clientSet) ServiceProfiles() ServiceProfileClient {
	return NewServiceProfileClient(c.client)
}

// Reader knows how to read and list ServiceProfiles.
type ServiceProfileReader interface {
	// Get retrieves a ServiceProfile for the given object key
	GetServiceProfile(ctx context.Context, key client.ObjectKey) (*ServiceProfile, error)

	// List retrieves list of ServiceProfiles for a given namespace and list options.
	ListServiceProfile(ctx context.Context, opts ...client.ListOption) (*ServiceProfileList, error)
}

// Writer knows how to create, delete, and update ServiceProfiles.
type ServiceProfileWriter interface {
	// Create saves the ServiceProfile object.
	CreateServiceProfile(ctx context.Context, obj *ServiceProfile, opts ...client.CreateOption) error

	// Delete deletes the ServiceProfile object.
	DeleteServiceProfile(ctx context.Context, key client.ObjectKey, opts ...client.DeleteOption) error

	// Update updates the given ServiceProfile object.
	UpdateServiceProfile(ctx context.Context, obj *ServiceProfile, opts ...client.UpdateOption) error

	// If the ServiceProfile object exists, update its spec. Otherwise, create the ServiceProfile object.
	UpsertServiceProfileSpec(ctx context.Context, obj *ServiceProfile, opts ...client.UpdateOption) error

	// Patch patches the given ServiceProfile object.
	PatchServiceProfile(ctx context.Context, obj *ServiceProfile, patch client.Patch, opts ...client.PatchOption) error

	// DeleteAllOf deletes all ServiceProfile objects matching the given options.
	DeleteAllOfServiceProfile(ctx context.Context, opts ...client.DeleteAllOfOption) error
}

// StatusWriter knows how to update status subresource of a ServiceProfile object.
type ServiceProfileStatusWriter interface {
	// Update updates the fields corresponding to the status subresource for the
	// given ServiceProfile object.
	UpdateServiceProfileStatus(ctx context.Context, obj *ServiceProfile, opts ...client.UpdateOption) error

	// Patch patches the given ServiceProfile object's subresource.
	PatchServiceProfileStatus(ctx context.Context, obj *ServiceProfile, patch client.Patch, opts ...client.PatchOption) error
}

// Client knows how to perform CRUD operations on ServiceProfiles.
type ServiceProfileClient interface {
	ServiceProfileReader
	ServiceProfileWriter
	ServiceProfileStatusWriter
}

type serviceProfileClient struct {
	client client.Client
}

func NewServiceProfileClient(client client.Client) *serviceProfileClient {
	return &serviceProfileClient{client: client}
}

func (c *serviceProfileClient) GetServiceProfile(ctx context.Context, key client.ObjectKey) (*ServiceProfile, error) {
	obj := &ServiceProfile{}
	if err := c.client.Get(ctx, key, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (c *serviceProfileClient) ListServiceProfile(ctx context.Context, opts ...client.ListOption) (*ServiceProfileList, error) {
	list := &ServiceProfileList{}
	if err := c.client.List(ctx, list, opts...); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *serviceProfileClient) CreateServiceProfile(ctx context.Context, obj *ServiceProfile, opts ...client.CreateOption) error {
	return c.client.Create(ctx, obj, opts...)
}

func (c *serviceProfileClient) DeleteServiceProfile(ctx context.Context, key client.ObjectKey, opts ...client.DeleteOption) error {
	obj := &ServiceProfile{}
	obj.SetName(key.Name)
	obj.SetNamespace(key.Namespace)
	return c.client.Delete(ctx, obj, opts...)
}

func (c *serviceProfileClient) UpdateServiceProfile(ctx context.Context, obj *ServiceProfile, opts ...client.UpdateOption) error {
	return c.client.Update(ctx, obj, opts...)
}

func (c *serviceProfileClient) UpsertServiceProfileSpec(ctx context.Context, obj *ServiceProfile, opts ...client.UpdateOption) error {
	existing, err := c.GetServiceProfile(ctx, client.ObjectKey{Name: obj.GetName(), Namespace: obj.GetNamespace()})
	if err != nil {
		if errors.IsNotFound(err) {
			return c.CreateServiceProfile(ctx, obj)
		}
		return err
	}
	existing.Spec = obj.Spec
	return c.client.Update(ctx, existing, opts...)
}

func (c *serviceProfileClient) PatchServiceProfile(ctx context.Context, obj *ServiceProfile, patch client.Patch, opts ...client.PatchOption) error {
	return c.client.Patch(ctx, obj, patch, opts...)
}

func (c *serviceProfileClient) DeleteAllOfServiceProfile(ctx context.Context, opts ...client.DeleteAllOfOption) error {
	obj := &ServiceProfile{}
	return c.client.DeleteAllOf(ctx, obj, opts...)
}

func (c *serviceProfileClient) UpdateServiceProfileStatus(ctx context.Context, obj *ServiceProfile, opts ...client.UpdateOption) error {
	return c.client.Status().Update(ctx, obj, opts...)
}

func (c *serviceProfileClient) PatchServiceProfileStatus(ctx context.Context, obj *ServiceProfile, patch client.Patch, opts ...client.PatchOption) error {
	return c.client.Status().Patch(ctx, obj, patch, opts...)
}
