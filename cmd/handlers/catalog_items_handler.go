package handlers

import (
	"context"
	"net/http"

	"github.com/rhpds/zerotouch-api/cmd/models"
)

// TODO: should contain a reference to the model to retrieve data from K8s API
// e.g.: CatalogItems
type CatalogItemsHandler struct {
	catalogItemRepo *models.CatalogItemRepo
}

// Make sure we conform to the StrictServer interface
var _ StrictServerInterface = (*CatalogItemsHandler)(nil)

func NewCatalogItemsHandler(catalogItemRepo *models.CatalogItemRepo) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItemRepo: catalogItemRepo,
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
		})

	}

	return ListCatalogItems200JSONResponse(items), nil
}

func (h *CatalogItemsHandler) GetCatalogItem(ctx context.Context, request GetCatalogItemRequestObject) (GetCatalogItemResponseObject, error) {
	catalogItem, err := h.catalogItemRepo.GetByName(request.Name)
	if err != nil {
		return GetCatalogItemdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: "Not implemented",
			}}, nil
	}

	return GetCatalogItem200JSONResponse(CatalogItem{
		Name:              catalogItem.Name,
		DisplayName:       catalogItem.DisplayName,
		Description:       catalogItem.Description,
		DescriptionFormat: catalogItem.DescriptionFormat,
		Id:                catalogItem.Id,
	}), nil
}

func (h *CatalogItemsHandler) Health(ctx context.Context, request HealthRequestObject) (HealthResponseObject, error) {
	status := OK

	return Health200JSONResponse{
		Status: &status,
	}, nil
}
