package handlers

import (
	"context"
	"net/http"

	"github.com/rhpds/zerotouch-api/cmd/models"
)

type CatalogItemsHandler struct {
	catalogItemsController *models.CatalogItemsController
	rcController           *models.ResourceClaimsController
}

// Make sure we conform to the StrictServer interface
var _ StrictServerInterface = (*CatalogItemsHandler)(nil)

func NewCatalogItemsHandler(catalogItemsControler *models.CatalogItemsController, rcController *models.ResourceClaimsController) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItemsController: catalogItemsControler,
		rcController:           rcController,
	}
}

// TODO: add pagination
func (h *CatalogItemsHandler) ListCatalogItems(ctx context.Context, request ListCatalogItemsRequestObject) (ListCatalogItemsResponseObject, error) {
	catalogItemList := h.catalogItemsController.ListAll()

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
	catalogItem, ok, err := h.catalogItemsController.GetByName(request.Name)
	if err != nil {
		return GetCatalogItemdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	if !ok {
		return GetCatalogItemdefaultJSONResponse{
			StatusCode: http.StatusNotFound,
			Body: Error{
				Code:    http.StatusNotFound,
				Message: "Not Found",
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

	// // Create a ResourceClaim
	// rc := models.ResourceClaimParameters{
	// 	Name:         "test-auto-1.babylon-empty-config.prod",
	// 	Namespace:    "user-kmalgich-redhat-com",
	// 	ProviderName: "tests.babylon-empty-config.prod",
	// 	Purpose:      "Testing",
	// 	Start:        time.Now(),
	// 	Stop:         time.Now().Add(1 * time.Hour),
	// 	End:          time.Now().Add(24 * time.Hour),
	// }

	// err := h.rcController.CreateResourceClaim(rc)
	// if err != nil {
	// 	fmt.Printf("Error creating ResourceClaim: %s\n", err.Error())
	// }

	// // Get a ResourceClaim status
	// claimStatus, ok, err := h.rcController.GetResourceClaimStatus("user-kmalgich-redhat-com", "test-auto-1.babylon-empty-config.prod")
	// if err != nil {
	// 	fmt.Printf("Error getting ResourceClaim status: %s\n", err.Error())
	// }

	// if !ok {
	// 	fmt.Printf("ResourceClaim not found\n")
	// }

	// fmt.Printf("Claim status: %+v\n", claimStatus)

	// err := h.rcController.DeleteResourceClaim("user-kmalgich-redhat-com", "test-auto-1.babylon-empty-config.prod")
	// if err != nil {
	// 	fmt.Printf("Error deleting ResourceClaim: %s\n", err.Error())
	// }

	return Health200JSONResponse{
		Status: &status,
	}, nil
}
