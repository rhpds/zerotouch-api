package models

import (
	"context"
	"strings"

	v1 "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1"
	babylon "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1/clientsets/babylon"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// const (
// 	error_not_found = "not found"
// )

type CatalogItemsController struct {
	clientSet *babylon.BabylonResourcesClient
	store     cache.Store
}

type CatalogItemInfo struct {
	Name              string
	NameSpace         string
	DisplayName       string
	Description       string
	DescriptionFormat string
	Id                string
	Provider          string
}

func NewCatalogItemsController(kubeconfigPath string, ctx context.Context, namespace string) (*CatalogItemsController, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	babylonClientSet, err := babylon.NewForConfig(config, ctx)
	if err != nil {
		return nil, err
	}

	store := babylon.WatchResources(babylonClientSet, namespace)

	return &CatalogItemsController{
		clientSet: babylonClientSet,
		store:     store,
	}, nil
}

func (c *CatalogItemsController) ListAll() []CatalogItemInfo {
	items := c.store.List()
	r := make([]CatalogItemInfo, 0, len(items))
	for _, item := range items {
		catalogItem := item.(*v1.CatalogItem)
		anotations := catalogItem.ObjectMeta.Annotations
		labels := catalogItem.ObjectMeta.Labels

		r = append(r, CatalogItemInfo{
			Name:              catalogItem.ObjectMeta.Name,
			NameSpace:         catalogItem.ObjectMeta.Namespace,
			DisplayName:       anotations["babylon.gpte.redhat.com/displayName"],
			Description:       anotations["babylon.gpte.redhat.com/description"],
			DescriptionFormat: anotations["babylon.gpte.redhat.com/descriptionFormat"],
			Id:                labels["gpte.redhat.com/asset-uuid"],
			Provider:          labels["babylon.gpte.redhat.com/Provider"],
		})
	}

	return r
}

// Key is a string that uniquely identifies a CatalogItem in a store
// and this is a string "namespace/name". We need to find the key
// because currently all we have is the name. and we need to extract namespace.
// TODO: check if we can update UI implementation to provide us namespace.
func (c *CatalogItemsController) findKey(name string) string {
	keys := c.store.ListKeys()
	for _, key := range keys {
		if strings.Contains(strings.ToLower(key), strings.ToLower(name)) {
			return key
		}
	}
	return ""
}

func (c *CatalogItemsController) GetByName(name string) (CatalogItemInfo, bool, error) {
	key := c.findKey(name)

	item, ok, err := c.store.GetByKey(key)
	if err != nil || !ok {
		return CatalogItemInfo{}, false, err
	}

	return CatalogItemInfo{
		Name:              item.(*v1.CatalogItem).ObjectMeta.Name,
		NameSpace:         item.(*v1.CatalogItem).ObjectMeta.Namespace,
		DisplayName:       item.(*v1.CatalogItem).ObjectMeta.Annotations["babylon.gpte.redhat.com/displayName"],
		Description:       item.(*v1.CatalogItem).ObjectMeta.Annotations["babylon.gpte.redhat.com/description"],
		DescriptionFormat: item.(*v1.CatalogItem).ObjectMeta.Annotations["babylon.gpte.redhat.com/descriptionFormat"],
		Id:                item.(*v1.CatalogItem).ObjectMeta.Labels["gpte.redhat.com/asset-uuid"],
		Provider:          item.(*v1.CatalogItem).ObjectMeta.Labels["babylon.gpte.redhat.com/Provider"],
	}, true, nil
}
