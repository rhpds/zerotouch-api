package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rhpds/zerotouch-api/cmd/log"
	"github.com/rhpds/zerotouch-api/cmd/models"
	"github.com/rhpds/zerotouch-api/cmd/recaptcha"
)

type RecaptchaConfig struct {
	ProjectID        string
	AuthKey          string
	RecapthcaSiteKey string
	Threshold        float64
	Disabled         bool
}

type CatalogItemsHandler struct {
	catalogItemsController *models.CatalogItemsController
	rcController           *models.ResourceClaimsController
	recaptchaConfig        *RecaptchaConfig
}

func NewCatalogItemsHandler(
	catalogItemsController *models.CatalogItemsController,
	rcController *models.ResourceClaimsController,
	recaptchaConfig *RecaptchaConfig,
) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItemsController: catalogItemsController,
		rcController:           rcController,
		recaptchaConfig:        recaptchaConfig,
	}
}

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
		log.Logger.Error("can't retrieve CatalogItem", "name", request.Name, "error", err.Error())

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

	var token string
	if request.Params.XGrecaptchaToken != nil {
		token = *request.Params.XGrecaptchaToken // reCAPTCHA token is not provided
	}

	if !h.recaptchaConfig.Disabled &&
		!h.verifyRecaptchaToken(token, "login") {
		return CreateProvision401JSONResponse(Error{
			Code:    http.StatusUnauthorized,
			Message: "reCAPTCHA Token verification failed",
		}), nil
	}

	rcInfo, err := h.rcController.CreateResourceClaim(rc)
	if err != nil {
		log.Logger.Error(
			"can't create provision",
			"provision name", request.Body.Name,
			"error", err.Error(),
		)

		return CreateProvision500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	log.Logger.Info("provision created", "provision name", request.Body.Name)

	// TODO: Provide lifespan.end with the response
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
	var token string
	if request.Params.XGrecaptchaToken != nil {
		token = *request.Params.XGrecaptchaToken // reCAPTCHA Token is not provided
	}

	if !h.recaptchaConfig.Disabled &&
		!h.verifyRecaptchaToken(token, "login") {
		return DeleteProvision401JSONResponse(Error{
			Code:    http.StatusUnauthorized,
			Message: "reCAPTCHA Token verification failed",
		}), nil
	}

	err := h.rcController.DeleteResourceClaim(request.Name)
	if err != nil {
		log.Logger.Error(
			"can't delete provision",
			"provision", request.Name,
			"error", err.Error(),
		)

		return DeleteProvision500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	log.Logger.Info("provision deleted", "provision name", request.Name)

	return DeleteProvision204Response{}, nil
}

// TODO: Provide lifespan end
func (h *CatalogItemsHandler) GetProvisionStatus(
	ctx context.Context,
	request GetProvisionStatusRequestObject,
) (GetProvisionStatusResponseObject, error) {
	claimStatus, ok, err := h.rcController.GetResourceClaimStatus(request.Name)
	if err != nil {
		log.Logger.Error(
			"can't retrieve provision status",
			"provision", request.Name,
			"error", err.Error(),
		)

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
		LabUserInterfaceUrl: &claimStatus.LabURL,
		RuntimeDefault:      claimStatus.RuntimeDefault,
		RuntimeMaximum:      claimStatus.RuntimeMaximum,
		LifespanEnd:         claimStatus.LifespanEnd,
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

// Helpers
func (h *CatalogItemsHandler) verifyRecaptchaToken(
	token string,
	action string,
) bool {
	assessmentParams := recaptcha.AssessmentParams{
		ProjectID:        h.recaptchaConfig.ProjectID,
		AuthKey:          h.recaptchaConfig.AuthKey,
		RecapthcaSiteKey: h.recaptchaConfig.RecapthcaSiteKey,
	}

	assessment, err := recaptcha.CreateAssessment(token, action, assessmentParams)
	if err != nil {
		log.Logger.Error("can't crate grecaptcha assessment", "error", err.Error())
		return false
	}

	if !assessment.IsTokenValid() {
		log.Logger.Debug("invalid token", "reason", assessment.GetInvalidReason())
		return false
	}

	if !assessment.IsActionValid() {
		log.Logger.Debug(
			"invalid token action",
			"expected: "+assessment.GetExpectedAction(),
			"actual: "+assessment.GetAction(),
		)
		return false
	}

	if !assessment.IsScoreValid(h.recaptchaConfig.Threshold) {
		log.Logger.Debug("token score is low", "score reason", assessment.GetScoreReasons())
		return false
	}

	return true
}
