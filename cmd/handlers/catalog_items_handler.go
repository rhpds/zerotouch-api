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

func NewCatalogItemsHandler(
	catalogItemsController *models.CatalogItemsController,
	rcController *models.ResourceClaimsController,
) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItemsController: catalogItemsController,
		rcController:           rcController,
	}
}

// TODO: add pagination
func (h *CatalogItemsHandler) ListCatalogItems(
	ctx context.Context,
	request ListCatalogItemsRequestObject,
) (ListCatalogItemsResponseObject, error) {
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

func (h *CatalogItemsHandler) GetCatalogItem(
	ctx context.Context,
	request GetCatalogItemRequestObject,
) (GetCatalogItemResponseObject, error) {
	catalogItem, ok, err := h.catalogItemsController.GetByName(request.Name)
	if err != nil {
		return GetCatalogItem500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	if !ok {
		return GetCatalogItem404Response{}, nil
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

func (h *CatalogItemsHandler) CreateProvision(
	ctx context.Context,
	request CreateProvisionRequestObject,
) (CreateProvisionResponseObject, error) {
	rc := models.ResourceClaimParameters{
		Name:         request.Body.Name,
		ProviderName: request.Body.ProviderName,
		Purpose:      request.Body.Purpose,
		Start:        request.Body.Start,
		Stop:         request.Body.Stop,
	}

	rcInfo, err := h.rcController.CreateResourceClaim(rc)
	if err != nil {
		return CreateProvision500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	return CreateProvision201JSONResponse{
		Body: ProvisionInfo{
			Name:      rcInfo.Name,
			UID:       rcInfo.UID,
			CreatedAt: rcInfo.CreatedAt,
		},
		Headers: CreateProvision201ResponseHeaders{
			Location: fmt.Sprintf("/provision/%s", rcInfo.Name),
		},
	}, nil
}

func (h *CatalogItemsHandler) DeleteProvision(
	ctx context.Context,
	request DeleteProvisionRequestObject,
) (DeleteProvisionResponseObject, error) {
	err := h.rcController.DeleteResourceClaim(request.Name)
	if err != nil {
		return DeleteProvision500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	return DeleteProvision204Response{}, nil
}

func (h *CatalogItemsHandler) GetProvisionStatus(
	ctx context.Context,
	request GetProvisionStatusRequestObject,
) (GetProvisionStatusResponseObject, error) {
	claimStatus, ok, err := h.rcController.GetResourceClaimStatus(request.Name)
	if err != nil {
		return GetProvisionStatus500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	if !ok {
		return GetProvisionStatus404Response{}, nil
	}

	if claimStatus == nil {
		return GetProvisionStatus202Response{}, nil
	}

	return GetProvisionStatus200JSONResponse(ProvisionStatus{
		State:               claimStatus.State,
		GUID:                claimStatus.GUID,
		LabUserInterfaceUrl: &claimStatus.ShowroomURL,
		RandomString:        claimStatus.RandomString,
		RuntimeDefault:      claimStatus.RuntimeDefault,
		RuntimeMaximum:      claimStatus.RuntimeMaximum,
	}), nil
}

func (h *CatalogItemsHandler) Health(
	ctx context.Context,
	request HealthRequestObject,
) (HealthResponseObject, error) {
	status := OK

	return Health200JSONResponse{
		Status: &status,
	}, nil
}
