package v1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const (
	catalog_items = "catalogitems"
)

// -----------------------------------------------------------------------------
// Client Sets
// -----------------------------------------------------------------------------
type BabylonResourcesInterface interface {
	CatalogItems(namespace string) CatalogItemsInterface
}

type BabylonResourcesClient struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*BabylonResourcesClient, error) {
	AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.UnversionedRESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &BabylonResourcesClient{restClient: client}, nil
}

func (c *BabylonResourcesClient) CatalogItems(namespace string, ctx context.Context) CatalogItemsInterface {
	return &catalogItemsClient{
		restClient: c.restClient,
		ns:         namespace,
		ctx:        ctx,
	}
}

// -----------------------------------------------------------------------------
// CatalogItems
// -----------------------------------------------------------------------------
type CatalogItemsInterface interface {
	List(opts metav1.ListOptions) (*CatalogItemList, error)
	Get(name string, options metav1.GetOptions) (*CatalogItem, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type catalogItemsClient struct {
	restClient rest.Interface
	ns         string
	ctx        context.Context
}

func (c *catalogItemsClient) List(opts metav1.ListOptions) (*CatalogItemList, error) {
	result := CatalogItemList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(catalog_items).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *catalogItemsClient) Get(name string, opts metav1.GetOptions) (*CatalogItem, error) {
	result := CatalogItem{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(catalog_items).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *catalogItemsClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(catalog_items).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(c.ctx)
}
