package babylon

import (
	"context"

	v1 "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const (
	catalog_items = "catalogitems"
	group_name    = "babylon.gpte.redhat.com"
	group_version = "v1"
)

var (
	schemaGroupVersion = schema.GroupVersion{
		Group:   group_name,
		Version: group_version,
	}

	schemeBuilder = runtime.NewSchemeBuilder(func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(schemaGroupVersion,
			&v1.CatalogItem{},
			&v1.CatalogItemList{},
		)

		metav1.AddToGroupVersion(scheme, schemaGroupVersion)
		return nil
	})
)

// -----------------------------------------------------------------------------
// Client Sets
// -----------------------------------------------------------------------------
type BabylonResourcesInterface interface {
	CatalogItems(namespace string) CatalogItemsInterface
}

type BabylonResourcesClient struct {
	restClient rest.Interface
	ctx        context.Context
}

func NewForConfig(c *rest.Config, ctx context.Context) (*BabylonResourcesClient, error) {
	schemeBuilder.AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &schemaGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.UnversionedRESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &BabylonResourcesClient{
		restClient: client,
		ctx:        ctx,
	}, nil
}

func (c *BabylonResourcesClient) CatalogItems(namespace string) CatalogItemsInterface {
	return &catalogItemsClient{
		restClient: c.restClient,
		ns:         namespace,
		ctx:        c.ctx,
	}
}

// -----------------------------------------------------------------------------
// CatalogItems
// -----------------------------------------------------------------------------
type CatalogItemsInterface interface {
	List(opts metav1.ListOptions) (*v1.CatalogItemList, error)
	Get(name string, options metav1.GetOptions) (*v1.CatalogItem, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type catalogItemsClient struct {
	restClient rest.Interface
	ns         string
	ctx        context.Context
}

func (c *catalogItemsClient) List(opts metav1.ListOptions) (*v1.CatalogItemList, error) {
	result := v1.CatalogItemList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(catalog_items).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *catalogItemsClient) Get(name string, opts metav1.GetOptions) (*v1.CatalogItem, error) {
	result := v1.CatalogItem{}
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

func WatchResources(client BabylonResourcesInterface, namespace string) cache.Store {
	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return client.CatalogItems(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.CatalogItems(namespace).Watch(options)
			},
		},
		&v1.CatalogItem{},
		0,
		cache.ResourceEventHandlerFuncs{},
	)

	// TODO: Provide chan object to stop the controller
	go controller.Run(wait.NeverStop)

	return store
}
