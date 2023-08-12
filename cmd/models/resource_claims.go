package models

import (
	"context"
	"fmt"

	"github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1/clientsets/poolboy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// TODO: Remove this function
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

func GetStore() cache.Store {
	config, err := clientcmd.BuildConfigFromFlags("", "/home/kmalgich/.kube/config")
	if err != nil {
		panic(err)
	}

	poolboyClientSet, err := poolboy.NewForConfig(config, context.Background())
	if err != nil {
		panic(err)
	}

	store := poolboy.WatchResourceResources(poolboyClientSet, "")

	return store
}
