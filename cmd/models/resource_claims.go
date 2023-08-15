package models

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1"
	"github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1/clientsets/poolboy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type ResourceClaimsController struct {
	clientSet *poolboy.PoolboyResourcesClient
	store     cache.Store
}

type ResourceClaimParameters struct {
	Name         string
	Namespace    string
	ProviderName string
	Purpose      string
	Start        time.Time
	Stop         time.Time
	// End          time.Time
}

type ResourceClaimStatus struct {
	GUID           string
	RandomString   string
	RuntimeDefault string
	RuntimeMaximum string
	State          string
}

type ResourceClaim struct {
	Name      string
	UID       string
	CreatedAt time.Time
}

// TODO: Add namespace parameter
func NewResourceClaimsController(kubeconfigPath string, ctx context.Context) (*ResourceClaimsController, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	poolboyClientSet, err := poolboy.NewForConfig(config, ctx)
	if err != nil {
		return nil, err
	}

	// Watch for resource claims in the all namespaces (last parameter)
	// and store them in cache
	store := poolboy.WatchResources(poolboyClientSet, "")

	return &ResourceClaimsController{
		clientSet: poolboyClientSet,
		store:     store,
	}, nil
}

func (c *ResourceClaimsController) CreateResourceClaim(parameters ResourceClaimParameters) (ResourceClaim, error) {
	rc := &v1.ResourceClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: parameters.Name,
		},
		Spec: v1.ResourceClaimSpec{
			Provider: v1.ResourceClaimProvider{
				Name: parameters.ProviderName,
				ParameterValues: v1.ResourceClaimParameterValues{
					Purpose:        parameters.Purpose,
					StartTimeStamp: parameters.Start.UTC().Format(time.RFC3339),
					EndTimeStamp:   parameters.Stop.UTC().Format(time.RFC3339),
				},
			},
			// Lifespan: v1.ResourceClaimLifespan{
			// 	End: parameters.End.UTC().Format(time.RFC3339),
			// },
		},
	}

	ret, err := c.clientSet.ResourceClaims(parameters.Namespace).Create(rc)
	if err != nil {
		return ResourceClaim{}, err
	}

	return ResourceClaim{
		Name:      ret.Name,
		UID:       string(ret.ObjectMeta.UID),
		CreatedAt: ret.ObjectMeta.CreationTimestamp.Time.UTC(),
	}, nil
}

func (c *ResourceClaimsController) DeleteResourceClaim(namespace string, name string) error {
	return c.clientSet.ResourceClaims(namespace).Delete(name, &metav1.DeleteOptions{})
}

func (c *ResourceClaimsController) GetResourceClaimStatus(namespace string, name string) (*ResourceClaimStatus, bool, error) {

	item, ok, err := c.store.GetByKey(fmt.Sprintf("%s/%s", namespace, name))
	if err != nil || !ok {
		return nil, ok, err
	}

	rc := item.(*v1.ResourceClaim)

	return &ResourceClaimStatus{
		GUID:           rc.Status.Summary.ProvisionData.GUID,
		RandomString:   rc.Status.Summary.ProvisionData.RandomString,
		RuntimeDefault: rc.Status.Summary.RuntimeDefault,
		RuntimeMaximum: rc.Status.Summary.RuntimeMaximum,
		State:          rc.Status.Summary.State,
	}, ok, nil
}
