package poolboy

import (
	"context"
	"time"

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
	resource_claim = "resourceclaims"
	group_name     = "poolboy.gpte.redhat.com"
	group_version  = "v1"
)

var (
	schemaGroupVersion = schema.GroupVersion{
		Group:   group_name,
		Version: group_version,
	}

	schemeBuilder = runtime.NewSchemeBuilder(func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(schemaGroupVersion,
			&v1.ResourceClaim{},
			&v1.ResourceClaimList{},
		)

		metav1.AddToGroupVersion(scheme, schemaGroupVersion)
		return nil
	})
)

// -----------------------------------------------------------------------------
// Client Sets
// -----------------------------------------------------------------------------
type PoolboyResourcesInterface interface {
	ResourceClaims(namespace string) ResourceClaimsInterface
}

type PoolboyResourcesClient struct {
	restClient rest.Interface
	ctx        context.Context
}

func NewForConfig(c *rest.Config, ctx context.Context) (*PoolboyResourcesClient, error) {
	schemeBuilder.AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &schemaGroupVersion
	config.APIPath = "/apis"

	// It seems that serializer bellow affects watch events
	//config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.UnversionedRESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &PoolboyResourcesClient{
		restClient: client,
		ctx:        ctx,
	}, nil
}

func (c *PoolboyResourcesClient) ResourceClaims(namespace string) ResourceClaimsInterface {
	return &resourceClaimsClient{
		restClient: c.restClient,
		ns:         namespace,
		ctx:        c.ctx,
	}
}

// -----------------------------------------------------------------------------
// ResourceClaims
// -----------------------------------------------------------------------------
type ResourceClaimsInterface interface {
	List(opts metav1.ListOptions) (*v1.ResourceClaimList, error)
	Get(name string, options metav1.GetOptions) (*v1.ResourceClaim, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Create(*v1.ResourceClaim) (*v1.ResourceClaim, error)
	Update(*v1.ResourceClaim) (*v1.ResourceClaim, error)
	Delete(name string, options *metav1.DeleteOptions) error
}

type resourceClaimsClient struct {
	restClient rest.Interface
	ns         string
	ctx        context.Context
}

func (c *resourceClaimsClient) List(opts metav1.ListOptions) (*v1.ResourceClaimList, error) {
	result := v1.ResourceClaimList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resource_claim).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *resourceClaimsClient) Get(name string, opts metav1.GetOptions) (*v1.ResourceClaim, error) {
	result := v1.ResourceClaim{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resource_claim).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *resourceClaimsClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resource_claim).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(c.ctx)
}

func (c *resourceClaimsClient) Create(item *v1.ResourceClaim) (*v1.ResourceClaim, error) {
	result := v1.ResourceClaim{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(resource_claim).
		Body(item).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *resourceClaimsClient) Update(item *v1.ResourceClaim) (*v1.ResourceClaim, error) {
	result := v1.ResourceClaim{}
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Resource(resource_claim).
		Name(item.Name).
		Body(item).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *resourceClaimsClient) Delete(name string, opts *metav1.DeleteOptions) error {
	return c.restClient.
		Delete().
		Namespace(c.ns).
		Resource(resource_claim).
		Name(name).
		Body(opts).
		Do(c.ctx).
		Error()
}

func WatchResourceResources(clientSet PoolboyResourcesInterface, namespace string) cache.Store {
	resourceClaimStore, resourceClaimController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (runtime.Object, error) {
				return clientSet.ResourceClaims(namespace).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.ResourceClaims(namespace).Watch(lo)
			},
		},
		&v1.ResourceClaim{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go resourceClaimController.Run(wait.NeverStop)
	return resourceClaimStore
}
