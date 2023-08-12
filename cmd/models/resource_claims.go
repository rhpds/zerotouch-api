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

// TODO:
// What is spec.lifespan?
type ResourceClaimParameters struct {
	Name         string
	Namespace    string
	ProviderName string
	Purpose      string
	Start        time.Time
	Stop         time.Time
}

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
	store := poolboy.WatchResourceResources(poolboyClientSet, "")

	return &ResourceClaimsController{
		clientSet: poolboyClientSet,
		store:     store,
	}, nil
}

func (c *ResourceClaimsController) CreateResourceClaim(parameters ResourceClaimParameters) error {
	rc := &v1.ResourceClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: parameters.Name,
		},
		Spec: v1.ResourceClaimSpec{
			Provider: v1.ResourceClaimProvider{
				Name: parameters.ProviderName,
				ParameterValues: v1.ResourceClaimParameterValues{
					Purpose:        parameters.Purpose,
					StartTimeStamp: parameters.Start.Format(time.RFC3339Nano), //TODO: Check if this is the correct format
					EndTimeStamp:   parameters.Stop.Format(time.RFC3339Nano),
				},
			},
			Lifespan: v1.ResourceClaimLifespan{
				End: "2023-08-14T00:00:00Z",
			},
		},
	}

	_, err := c.clientSet.ResourceClaims(parameters.Namespace).Create(rc)
	if err != nil {
		return err
	}

	return nil

}

//
// TODO: Bellow are draft functions, they should be removed
//
//

func PrintResourceClaims() (int, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "/home/kmalgich/.kube/config")
	if err != nil {
		return 0, err
	}

	poolboyClientSet, err := poolboy.NewForConfig(config, context.Background())
	if err != nil {
		return 0, err
	}

	resourceClaims, err := poolboyClientSet.ResourceClaims("").List(metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	for _, resourceClaim := range resourceClaims.Items {
		fmt.Println(resourceClaim.Name)
	}

	return len(resourceClaims.Items), nil

	//----
	// store := poolboy.WatchResourceResources(poolboyClientSet, "")

	// store.Resync()

	// keys := store.ListKeys()

	// for key := range keys {
	// 	fmt.Println(key)
	// }

	// return len(keys), nil
}
