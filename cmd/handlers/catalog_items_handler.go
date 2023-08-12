package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rhpds/zerotouch-api/cmd/models"
)

// TODO: should contain a reference to the model to retrieve data from K8s API
// e.g.: CatalogItems
type CatalogItemsHandler struct {
	catalogItemRepo *models.CatalogItemRepo
	rcController    *models.ResourceClaimsController
}

// Make sure we conform to the StrictServer interface
var _ StrictServerInterface = (*CatalogItemsHandler)(nil)

func NewCatalogItemsHandler(catalogItemRepo *models.CatalogItemRepo, rcController *models.ResourceClaimsController) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItemRepo: catalogItemRepo,
		rcController:    rcController,
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

	//	h.cache.Resync()

	// List all ResourceClaims
	// claims := h.cache.List()
	// for _, v := range claims {
	// 	fmt.Printf("%+v\n\n", v.(*v1.ResourceClaim))
	// }

	// Create a ResourceClaim
	rc := models.ResourceClaimParameters{
		Name:         "test-auto-3.babylon-empty-config.prod",
		Namespace:    "user-kmalgich-redhat-com",
		ProviderName: "tests.babylon-empty-config.prod",
		Purpose:      "Testing",
		Start:        time.Now(),
		Stop:         time.Now().Add(1 * time.Hour),
	}

	err := h.rcController.CreateResourceClaim(rc)
	if err != nil {
		fmt.Printf("Error creating ResourceClaim: %s\n", err.Error())
	}

	return Health200JSONResponse{
		Status: &status,
	}, nil
}
