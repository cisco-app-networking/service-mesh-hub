// Code generated by solo-kit. DO NOT EDIT.

package v1alpha1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type ClusterRbacConfigWatcher interface {
	// watch cluster-scoped ClusterRbacConfigs
	Watch(opts clients.WatchOpts) (<-chan ClusterRbacConfigList, <-chan error, error)
}

type ClusterRbacConfigClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(name string, opts clients.ReadOpts) (*ClusterRbacConfig, error)
	Write(resource *ClusterRbacConfig, opts clients.WriteOpts) (*ClusterRbacConfig, error)
	Delete(name string, opts clients.DeleteOpts) error
	List(opts clients.ListOpts) (ClusterRbacConfigList, error)
	ClusterRbacConfigWatcher
}

type clusterRbacConfigClient struct {
	rc clients.ResourceClient
}

func NewClusterRbacConfigClient(rcFactory factory.ResourceClientFactory) (ClusterRbacConfigClient, error) {
	return NewClusterRbacConfigClientWithToken(rcFactory, "")
}

func NewClusterRbacConfigClientWithToken(rcFactory factory.ResourceClientFactory, token string) (ClusterRbacConfigClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &ClusterRbacConfig{},
		Token:        token,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base ClusterRbacConfig resource client")
	}
	return NewClusterRbacConfigClientWithBase(rc), nil
}

func NewClusterRbacConfigClientWithBase(rc clients.ResourceClient) ClusterRbacConfigClient {
	return &clusterRbacConfigClient{
		rc: rc,
	}
}

func (client *clusterRbacConfigClient) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *clusterRbacConfigClient) Register() error {
	return client.rc.Register()
}

func (client *clusterRbacConfigClient) Read(name string, opts clients.ReadOpts) (*ClusterRbacConfig, error) {
	opts = opts.WithDefaults()

	resource, err := client.rc.Read("", name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*ClusterRbacConfig), nil
}

func (client *clusterRbacConfigClient) Write(clusterRbacConfig *ClusterRbacConfig, opts clients.WriteOpts) (*ClusterRbacConfig, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(clusterRbacConfig, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*ClusterRbacConfig), nil
}

func (client *clusterRbacConfigClient) Delete(name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()

	return client.rc.Delete("", name, opts)
}

func (client *clusterRbacConfigClient) List(opts clients.ListOpts) (ClusterRbacConfigList, error) {
	opts = opts.WithDefaults()

	resourceList, err := client.rc.List("", opts)
	if err != nil {
		return nil, err
	}
	return convertToClusterRbacConfig(resourceList), nil
}

func (client *clusterRbacConfigClient) Watch(opts clients.WatchOpts) (<-chan ClusterRbacConfigList, <-chan error, error) {
	opts = opts.WithDefaults()

	resourcesChan, errs, initErr := client.rc.Watch("", opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	clusterRbacConfigsChan := make(chan ClusterRbacConfigList)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				clusterRbacConfigsChan <- convertToClusterRbacConfig(resourceList)
			case <-opts.Ctx.Done():
				close(clusterRbacConfigsChan)
				return
			}
		}
	}()
	return clusterRbacConfigsChan, errs, nil
}

func convertToClusterRbacConfig(resources resources.ResourceList) ClusterRbacConfigList {
	var clusterRbacConfigList ClusterRbacConfigList
	for _, resource := range resources {
		clusterRbacConfigList = append(clusterRbacConfigList, resource.(*ClusterRbacConfig))
	}
	return clusterRbacConfigList
}
