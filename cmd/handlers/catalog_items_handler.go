package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rhpds/zerotouch-api/cmd/log"
	"github.com/rhpds/zerotouch-api/cmd/models"
	"github.com/rhpds/zerotouch-api/cmd/ratings"
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
	ratingsClient          *ratings.RatingClient
}

func NewCatalogItemsHandler(
	catalogItemsController *models.CatalogItemsController,
	rcController *models.ResourceClaimsController,
	recaptchaConfig *RecaptchaConfig,
	ratingsClient *ratings.RatingClient,
) *CatalogItemsHandler {
	return &CatalogItemsHandler{
		catalogItemsController: catalogItemsController,
		rcController:           rcController,
		recaptchaConfig:        recaptchaConfig,
		ratingsClient:          ratingsClient,
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
			Id:                v.AssetUUID,
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
		Id:                catalogItem.AssetUUID,
		Provider:          catalogItem.Provider,
	}), nil
}

func (h *CatalogItemsHandler) CreateServiceRequest(
	ctx context.Context,
	request CreateServiceRequestRequestObject,
) (CreateServiceRequestResponseObject, error) {
	var stop string
	if request.Body.Stop != nil {
		stopTimeStamp := *request.Body.Stop
		stop = stopTimeStamp.UTC().Format(time.RFC3339)
	}

	catalogItem, found, err := h.catalogItemsController.GetByName(request.Body.ProviderName)
	if err != nil {
		log.Logger.Error(
			"can't create provision",
			"provision name", request.Body.Name,
			"error", err.Error(),
		)

		return CreateServiceRequest500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	if !found {
		log.Logger.Error(
			"can't create provision",
			"provision name", request.Body.Name,
			"provider not found", request.Body.ProviderName,
		)

		return CreateServiceRequest500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Provider %s not found", request.Body.ProviderName),
		}), nil
	}

	rc := models.ResourceClaimParameters{
		Name:         request.Body.Name,
		ProviderName: request.Body.ProviderName,
		Purpose:      request.Body.Purpose,
		Start:        request.Body.Start.UTC().Format(time.RFC3339),
		Stop:         stop,
	}

	lifespanDuration, err := catalogItem.GetDefaultLifespan()
	if err == nil {
		lifespanEnd := request.Body.Start
		lifespanEnd = lifespanEnd.Add(lifespanDuration)
		lifespan := lifespanEnd.Format(time.RFC3339)

		rc.Lifespan = &lifespan
	}

	var token string
	if request.Params.XGrecaptchaToken != nil {
		token = *request.Params.XGrecaptchaToken // reCAPTCHA token is not provided
	}

	if !h.recaptchaConfig.Disabled &&
		!h.verifyRecaptchaToken(token, "login") {
		return CreateServiceRequest401JSONResponse(Error{
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

		return CreateServiceRequest500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	log.Logger.Info("provision created", "provision name", request.Body.Name)

	// TODO: Provide lifespan.end with the response
	return CreateServiceRequest201JSONResponse{
		Body: ProvisionInfo{
			Name:      rcInfo.Name,
			UID:       rcInfo.UID,
			CreatedAt: rcInfo.CreatedAt,
		},
		Headers: CreateServiceRequest201ResponseHeaders{
			Location: fmt.Sprintf("/provision/%s", rcInfo.Name),
		},
	}, nil
}

func (h *CatalogItemsHandler) DeleteServiceRequest(
	ctx context.Context,
	request DeleteServiceRequestRequestObject,
) (DeleteServiceRequestResponseObject, error) {
	err := h.rcController.DeleteResourceClaim(request.Name)
	if err != nil {
		log.Logger.Error(
			"can't delete provision",
			"provision", request.Name,
			"error", err.Error(),
		)

		return DeleteServiceRequest500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	log.Logger.Info("provision deleted", "provision name", request.Name)

	return DeleteServiceRequest204Response{}, nil
}

// TODO: Provide lifespan end
func (h *CatalogItemsHandler) GetServiceRequestStatus(
	ctx context.Context,
	request GetServiceRequestStatusRequestObject,
) (GetServiceRequestStatusResponseObject, error) {
	claimStatus, ok, err := h.rcController.GetResourceClaimStatus(request.Name)
	if err != nil {
		log.Logger.Error(
			"can't retrieve provision status",
			"provision", request.Name,
			"error", err.Error(),
		)

		return GetServiceRequestStatus500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	if !ok {
		return GetServiceRequestStatus404Response{}, nil
	}

	if claimStatus == nil {
		return GetServiceRequestStatus202Response{}, nil
	}

	return GetServiceRequestStatus200JSONResponse(ProvisionStatus{
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

func (h *CatalogItemsHandler) GetRating(
	ctx context.Context,
	request GetRatingRequestObject,
) (GetRatingResponseObject, error) {

	catalogItem, found, err := h.catalogItemsController.GetByName(request.Name)
	if err != nil {
		log.Logger.Error(
			"can't create rating",
			"provision name", request.Name,
			"error", err.Error(),
		)

		return GetRating500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	if !found {
		log.Logger.Error(
			"can't create rating",
			"provision name", request.Name,
			"provider not found", request.Name,
		)

		return GetRating500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Provider %s not found", request.Name),
		}), nil
	}


	rating, err := h.ratingsClient.GetRatings(catalogItem.AssetUUID)
	if err != nil {
		return GetRating500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	return GetRating200JSONResponse{
		RatingScore:  fmt.Sprintf("%.1f", rating.RatingScore),
		TotalRatings: strconv.Itoa(rating.TotalRatings),
	}, nil
}

func (h *CatalogItemsHandler) CreateRating(
	ctx context.Context,
	request CreateRatingRequestObject,
) (CreateRatingResponseObject, error) {
	rating := ratings.NewRating{
		Email: request.Body.Email,
	}

	if request.Body.Rating != nil {
		rating.Rating = *request.Body.Rating
	}

	if request.Body.Comment != nil {
		rating.Comment = *request.Body.Comment
	}

	if request.Body.Useful != nil {
		rating.Useful = *request.Body.Useful
	}

	uid, err := h.rcController.GetUID(request.Body.ProvisionName)
	if err != nil {
		return CreateRating500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	rating.RequestID = *uid

	_, err = h.ratingsClient.SetRating(rating)
	if err != nil {
		return CreateRating500JSONResponse(Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}), nil
	}

	return CreateRating201Response{}, nil
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
