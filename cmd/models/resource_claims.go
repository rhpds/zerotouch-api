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
	clientSet      *poolboy.PoolboyResourcesClient
	store          cache.Store
	namespace      string
	OnStatusUpdate func(details ResourceClaimDetails)
}

type ResourceClaimParameters struct {
	Name         string
	ProviderName string
	Purpose      string
	Start        string
	Stop         string
	Lifespan     *string
}

type ResourceClaimDetails struct {
	Name           string
	Provider       string
	GUID           string
	LabURL         string
	RuntimeDefault string
	RuntimeMaximum string
	State          string
	LifespanStart  string
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
	rcController := ResourceClaimsController{
		namespace: namespace,
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	rcController.clientSet, err = poolboy.NewForConfig(config, ctx)
	if err != nil {
		return nil, err
	}

	rcEventHandlers := poolboy.ResourceClaimEvents{
		UpdateEvent: rcController.OnResourceClaimUpdate,
	}

	rcController.store = poolboy.WatchResourceClaims(
		rcController.clientSet,
		namespace,
		rcEventHandlers,
	)

	return &rcController, nil
}

func (c *ResourceClaimsController) CreateResourceClaim(
	parameters ResourceClaimParameters,
) (ResourceClaim, error) {
	var lifespan *v1.ResourceClaimSpecLifespan = nil
	if parameters.Lifespan != nil {
		lifespan = &v1.ResourceClaimSpecLifespan{
			End: *parameters.Lifespan,
		}
	}

	rc := &v1.ResourceClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: parameters.Name,
		},
		Spec: v1.ResourceClaimSpec{
			Provider: v1.ResourceClaimProvider{
				Name: parameters.ProviderName,
				ParameterValues: v1.ResourceClaimParameterValues{
					Purpose:        parameters.Purpose,
					StartTimeStamp: parameters.Start,
					StopTimeStamp:  parameters.Stop,
				},
			},
			Lifespan: lifespan,
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

func (c *ResourceClaimsController) GetResourceClaimDetails(
	name string,
) (*ResourceClaimDetails, bool, error) {
	item, ok, err := c.store.GetByKey(fmt.Sprintf("%s/%s", c.namespace, name))
	if err != nil || !ok {
		return nil, ok, err
	}

	rc := item.(*v1.ResourceClaim)

	// Status is not available yet
	if (v1.ResourceClaimStatusSummary{}) == rc.Status.Summary {
		return nil, ok, nil
	}

	return &ResourceClaimDetails{
		Name:           rc.Name,
		Provider:       rc.Spec.Provider.Name,
		GUID:           rc.Status.Summary.ProvisionData.GUID,
		LabURL:         rc.Status.Summary.ProvisionData.LabURL,
		RuntimeDefault: rc.Status.Summary.RuntimeDefault,
		RuntimeMaximum: rc.Status.Summary.RuntimeMaximum,
		State:          rc.Status.Summary.State,
		LifespanEnd:    rc.Status.Lifespan.End,
	}, ok, nil
}

func (c *ResourceClaimsController) GetUID(name string) (*string, error) {
	item, ok, err := c.store.GetByKey(fmt.Sprintf("%s/%s", c.namespace, name))
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("%s not found", name)
	}

	rc := item.(*v1.ResourceClaim)

	return (*string)(&rc.ObjectMeta.UID), nil
}

// -----------------------------------------------------------------------------
// Event Handlers
// -----------------------------------------------------------------------------
func (c *ResourceClaimsController) OnResourceClaimUpdate(oldRC, newRC *v1.ResourceClaim) {
	if oldRC.Status.Summary.State != newRC.Status.Summary.State {
		if c.OnStatusUpdate != nil {
			c.OnStatusUpdate(ResourceClaimDetails{
				Name:           newRC.Name,
				Provider:       newRC.Spec.Provider.Name,
				GUID:           newRC.Status.Summary.ProvisionData.GUID,
				LabURL:         newRC.Status.Summary.ProvisionData.LabURL,
				RuntimeDefault: newRC.Status.Summary.RuntimeDefault,
				RuntimeMaximum: newRC.Status.Summary.RuntimeMaximum,
				State:          newRC.Status.Summary.State,
				LifespanStart:  newRC.Status.Lifespan.Start,
				LifespanEnd:    newRC.Status.Lifespan.End,
			})
		}
	}
}
