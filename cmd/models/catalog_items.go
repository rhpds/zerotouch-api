package models

import (
	"context"
	"errors"

	babylon "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1/clientsets"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	babylonClientSet, err := babylon.NewForConfig(config)
	if err != nil {
		return 0, err
	}

	catalogItems, err := babylonClientSet.CatalogItems("", context.Background()).List(metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	for item := range catalogItems.Items {

		var info CatalogItemInfo
		info.Name = catalogItems.Items[item].ObjectMeta.Name

		annotations := catalogItems.Items[item].ObjectMeta.Annotations
		info.DisplayName = annotations["babylon.gpte.redhat.com/displayName"]
		info.Description = annotations["babylon.gpte.redhat.com/description"]
		info.DescriptionFormat = annotations["babylon.gpte.redhat.com/descriptionFormat"]

		labels := catalogItems.Items[item].ObjectMeta.Labels
		info.Provider = labels["babylon.gpte.redhat.com/Provider"]
		info.Id = labels["gpte.redhat.com/asset-uuid"]

		c.catalogItems[info.Name] = info
	}

	return len(c.catalogItems), nil
}
