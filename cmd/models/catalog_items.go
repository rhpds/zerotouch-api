package models

import (
	"context"
	"errors"

	v1 "github.com/rhpds/zerotouch-api/cmd/kube/api/types"
	"k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	error_not_found = "not found"
)

type CatalogItemInfo struct {
	Name              string
	DisplayName       string
	Description       string
	DescriptionFormat string
	Id                string
	Provider          string
}

type CatalogItemRepo struct {
	catalogItems map[string]CatalogItemInfo
}

func NewCatalogItemRepo() *CatalogItemRepo {
	return &CatalogItemRepo{
		catalogItems: make(map[string]CatalogItemInfo),
	}
}

// TODO: add pagination
func (c *CatalogItemRepo) ListAll() ([]CatalogItemInfo, error) {
	r := make([]CatalogItemInfo, 0, len(c.catalogItems))
	for _, v := range c.catalogItems {
		r = append(r, v)
	}

	return r, nil
}

func (c *CatalogItemRepo) GetByName(name string) (CatalogItemInfo, error) {
	item, ok := c.catalogItems[name]
	if !ok {
		return CatalogItemInfo{}, errors.New(error_not_found)
	}

	return item, nil
}

func (c *CatalogItemRepo) Refresh(kubeconfig string) (int, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		return 0, err
	}

	v1.AddToScheme(scheme.Scheme)

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1.GroupName, Version: v1.GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	restClent, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		return 0, err
	}

	result := v1.CatalogItemList{}
	//TODO: Remove hardcoded resource
	err = restClent.
		Get().
		Resource("catalogitems").
		Do(context.Background()).
		Into(&result)

	if err != nil {
		return 0, err
	}

	for item := range result.Items {

		var info CatalogItemInfo
		info.Name = result.Items[item].ObjectMeta.Name

		annotations := result.Items[item].ObjectMeta.Annotations
		info.DisplayName = annotations["babylon.gpte.redhat.com/displayName"]
		info.Description = annotations["babylon.gpte.redhat.com/description"]
		info.DescriptionFormat = annotations["babylon.gpte.redhat.com/descriptionFormat"]

		labels := result.Items[item].ObjectMeta.Labels
		info.Provider = labels["babylon.gpte.redhat.com/Provider"]
		info.Id = labels["gpte.redhat.com/asset-uuid"]

		c.catalogItems[info.Name] = info
	}

	return len(c.catalogItems), nil
}
