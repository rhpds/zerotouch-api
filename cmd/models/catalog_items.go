package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	v1 "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1"
	babylon "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1/clientsets/babylon"
)

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
	AssetUUID         string
	Provider          string
	DefaultLifespan   string
}

func NewCatalogItemsController(
	kubeconfigPath string,
	ctx context.Context,
	namespace string,
) (*CatalogItemsController, error) {
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
		annotations := catalogItem.ObjectMeta.Annotations
		labels := catalogItem.ObjectMeta.Labels

		r = append(r, CatalogItemInfo{
			Name:              catalogItem.ObjectMeta.Name,
			NameSpace:         catalogItem.ObjectMeta.Namespace,
			DisplayName:       annotations["babylon.gpte.redhat.com/displayName"],
			Description:       annotations["babylon.gpte.redhat.com/description"],
			DescriptionFormat: annotations["babylon.gpte.redhat.com/descriptionFormat"],
			AssetUUID:         labels["gpte.redhat.com/asset-uuid"],
			Provider:          labels["babylon.gpte.redhat.com/Provider"],
		})
	}

	return r
}

// Key is a string that uniquely identifies a CatalogItem in a store
// and this is a string in "namespace/name" format. We need to "build" the key
// because currently all we have is the name, and we need to extract namespace
// from the keys array.
func (c *CatalogItemsController) findKey(name string) string {
	keys := c.store.ListKeys()
	for _, key := range keys {
		strSlice := strings.Split(key, "/")
		providerName := strSlice[len(strSlice)-1]
		if strings.Compare(strings.ToLower(providerName), strings.ToLower(name)) == 0 {
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
		AssetUUID:         item.(*v1.CatalogItem).ObjectMeta.Labels["gpte.redhat.com/asset-uuid"],
		Provider:          item.(*v1.CatalogItem).ObjectMeta.Labels["babylon.gpte.redhat.com/Provider"],
		DefaultLifespan:   item.(*v1.CatalogItem).Spec.Lifespan.Default,
	}, true, nil
}

func (ci *CatalogItemInfo) GetDefaultLifespan() (time.Duration, error) {
	var duration time.Duration
	lifespan := ci.DefaultLifespan

	duration, err := time.ParseDuration(lifespan)
	if err == nil {
		return duration, nil
	}

	value, err := strconv.ParseInt(lifespan[:len(lifespan)-1], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("unknown value \"%s\" in duration \"%s\"", lifespan[:len(lifespan)-1], lifespan)
	}

	unit := lifespan[len(lifespan)-1]
	if unit != 'd' {
		return duration, fmt.Errorf("unknown unit \"%c\" in duration \"%s\"", unit, lifespan)
	}

	return time.ParseDuration(fmt.Sprintf("%dh", 24*value))
}
