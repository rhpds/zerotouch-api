package models

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	v1 "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1"
	"github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1/clientsets/poolboy"
)

type ResourceClaimsController struct {
	clientSet *poolboy.PoolboyResourcesClient
	store     cache.Store
	namespace string
}

type ResourceClaimParameters struct {
	Name         string
	ProviderName string
	Purpose      string
	Start        time.Time
	Stop         time.Time
}

type ResourceClaimStatus struct {
	GUID           string
	LabURL         string
	RuntimeDefault string
	RuntimeMaximum string
	State          string
	LifespanEnd    string
}

type ResourceClaim struct {
	Name      string
	UID       string
	CreatedAt time.Time
}

func NewResourceClaimsController(
	kubeconfigPath string,
	namespace string,
	ctx context.Context,
) (*ResourceClaimsController, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	poolboyClientSet, err := poolboy.NewForConfig(config, ctx)
	if err != nil {
		return nil, err
	}

	store := poolboy.WatchResources(poolboyClientSet, namespace)

	return &ResourceClaimsController{
		clientSet: poolboyClientSet,
		store:     store,
		namespace: namespace,
	}, nil
}

func (c *ResourceClaimsController) CreateResourceClaim(
	parameters ResourceClaimParameters,
) (ResourceClaim, error) {
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
					StopTimeStamp:  parameters.Stop.UTC().Format(time.RFC3339),
				},
			},
		},
	}

	ret, err := c.clientSet.ResourceClaims(c.namespace).Create(rc)
	if err != nil {
		return ResourceClaim{}, err
	}

	return ResourceClaim{
		Name:      ret.Name,
		UID:       string(ret.ObjectMeta.UID),
		CreatedAt: ret.ObjectMeta.CreationTimestamp.Time.UTC(),
	}, nil
}

func (c *ResourceClaimsController) DeleteResourceClaim(name string) error {
	return c.clientSet.ResourceClaims(c.namespace).Delete(name, &metav1.DeleteOptions{})
}

func (c *ResourceClaimsController) GetResourceClaimStatus(
	name string,
) (*ResourceClaimStatus, bool, error) {
	item, ok, err := c.store.GetByKey(fmt.Sprintf("%s/%s", c.namespace, name))
	if err != nil || !ok {
		return nil, ok, err
	}

	rc := item.(*v1.ResourceClaim)

	// Status is not available yet
	if (v1.ResourceClaimStatusSummary{}) == rc.Status.Summary {
		return nil, ok, nil
	}

	return &ResourceClaimStatus{
		GUID:           rc.Status.Summary.ProvisionData.GUID,
		LabURL:         rc.Status.Summary.ProvisionData.LabURL,
		RuntimeDefault: rc.Status.Summary.RuntimeDefault,
		RuntimeMaximum: rc.Status.Summary.RuntimeMaximum,
		State:          rc.Status.Summary.State,
		LifespanEnd:    rc.Status.Lifespan.End,
	}, ok, nil
}
