package handlers

import (
	"context"
	"fmt"
	"net/http"

	v1 "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1"
	"github.com/rhpds/zerotouch-api/cmd/models"
	"k8s.io/client-go/tools/cache"
)

// TODO: should contain a reference to the model to retrieve data from K8s API
// e.g.: CatalogItems
type CatalogItemsHandler struct {
	catalogItemRepo *models.CatalogItemRepo

	// TODO: testing
	cache cache.Store
}

// Make sure we conform to the StrictServer interface
var _ StrictServerInterface = (*CatalogItemsHandler)(nil)

func NewCatalogItemsHandler(catalogItemRepo *models.CatalogItemRepo) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItemRepo: catalogItemRepo,
		cache:           models.GetStore(), //TODO: testing
	}
}

func (h *CatalogItemsHandler) ListCatalogItems(ctx context.Context, request ListCatalogItemsRequestObject) (ListCatalogItemsResponseObject, error) {
	catalogItemList, err := h.catalogItemRepo.ListAll()
	if err != nil {
		return ListCatalogItemsdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	items := make([]CatalogItem, 0, len(catalogItemList))
	for _, v := range catalogItemList {
		items = append(items, CatalogItem{
			Name:              v.Name,
			DisplayName:       v.DisplayName,
			Description:       v.Description,
			DescriptionFormat: v.DescriptionFormat,
			Id:                v.Id,
			Provider:          v.Provider,
		})

	}

	return ListCatalogItems200JSONResponse(items), nil
}

func (h *CatalogItemsHandler) GetCatalogItem(ctx context.Context, request GetCatalogItemRequestObject) (GetCatalogItemResponseObject, error) {
	catalogItem, err := h.catalogItemRepo.GetByName(request.Name)
	if err != nil {
		return GetCatalogItemdefaultJSONResponse{
			StatusCode: http.StatusNotFound,
			Body: Error{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}}, nil
	}

	return GetCatalogItem200JSONResponse(CatalogItem{
		Name:              catalogItem.Name,
		DisplayName:       catalogItem.DisplayName,
		Description:       catalogItem.Description,
		DescriptionFormat: catalogItem.DescriptionFormat,
		Id:                catalogItem.Id,
		Provider:          catalogItem.Provider,
	}), nil
}

func (h *CatalogItemsHandler) Health(ctx context.Context, request HealthRequestObject) (HealthResponseObject, error) {
	status := OK

	h.cache.Resync()

	// keys := h.cache.ListKeys()
	// for _, v := range keys {

	// 	item, ok, err := h.cache.GetByKey(v)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	if !ok {
	// 		fmt.Println("not found")
	// 	}

	// 	fmt.Printf("%+v\n", item.(*v1.ResourceClaim))
	// }

	claims := h.cache.List()
	for _, v := range claims {
		fmt.Printf("%+v\n\n", v.(*v1.ResourceClaim))
	}

	return Health200JSONResponse{
		Status: &status,
	}, nil
}
