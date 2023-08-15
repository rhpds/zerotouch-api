package handlers

import (
	"context"
	"fmt"
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

func (h *CatalogItemsHandler) CreateProvision(ctx context.Context, request CreateProvisionRequestObject) (CreateProvisionResponseObject, error) {
	// TODO: Get namespace from environment variable
	// TODO: Get default time values from CI
	rc := models.ResourceClaimParameters{
		Name:         request.Body.Name,
		Namespace:    "user-kmalgich-redhat-com",
		ProviderName: request.Body.ProviderName,
		Purpose:      request.Body.Purpose,
		Start:        request.Body.Start,
		Stop:         request.Body.Stop,
	}

	rcInfo, err := h.rcController.CreateResourceClaim(rc)
	if err != nil {
		return CreateProvisiondefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("Can't start provision %s: %s", request.Body.Name, err.Error()),
			}}, nil
	}

	return CreateProvision201JSONResponse{
		Body: ProvisionInfo{
			Name:      rcInfo.Name,
			UID:       rcInfo.UID,
			CreatedAt: &rcInfo.CreatedAt,
		},
		Headers: CreateProvision201ResponseHeaders{
			Location: fmt.Sprintf("/provision/%s", rcInfo.Name),
		},
	}, nil
}

func (h *CatalogItemsHandler) DeleteProvision(ctx context.Context, request DeleteProvisionRequestObject) (DeleteProvisionResponseObject, error) {
	err := h.rcController.DeleteResourceClaim("user-kmalgich-redhat-com", request.Name)
	if err != nil {
		return DeleteProvisiondefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	return DeleteProvision204Response{}, nil
}

func (h *CatalogItemsHandler) GetProvisionStatus(ctx context.Context, request GetProvisionStatusRequestObject) (GetProvisionStatusResponseObject, error) {
	claimStatus, ok, err := h.rcController.GetResourceClaimStatus("user-kmalgich-redhat-com", request.Name)
	if err != nil {
		return GetProvisionStatusdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	if !ok {
		return GetProvisionStatusdefaultJSONResponse{
			StatusCode: http.StatusNotFound,
			Body: Error{
				Code:    http.StatusNotFound,
				Message: "Not Found",
			}}, nil
	}

	return GetProvisionStatus200JSONResponse(ProvisionStatus{
		State:          claimStatus.State,
		GUID:           claimStatus.GUID,
		RandomString:   claimStatus.RandomString,
		RuntimeDefault: claimStatus.RuntimeDefault,
		RuntimeMaximum: claimStatus.RuntimeMaximum,
	}), nil

}

func (h *CatalogItemsHandler) Health(ctx context.Context, request HealthRequestObject) (HealthResponseObject, error) {
	status := OK

	return Health200JSONResponse{
		Status: &status,
	}, nil
}
