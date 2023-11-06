// Package handlers provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// Defines values for StatusStatus.
const (
	ERROR StatusStatus = "ERROR"
	OK    StatusStatus = "OK"
)

// CatalogItem defines model for CatalogItem.
type CatalogItem struct {
	Description       string `json:"Description"`
	DescriptionFormat string `json:"DescriptionFormat"`
	DisplayName       string `json:"DisplayName"`
	Name              string `json:"Name"`
	Provider          string `json:"Provider"`
	Id                string `json:"id"`
}

// CatalogItemRating defines model for CatalogItemRating.
type CatalogItemRating struct {
	RatingScore  string `json:"RatingScore"`
	TotalRatings string `json:"TotalRatings"`
}

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// ProvisionInfo defines model for ProvisionInfo.
type ProvisionInfo struct {
	CreatedAt time.Time `json:"CreatedAt"`
	Name      string    `json:"Name"`
	UID       string    `json:"UID"`
}

// ProvisionParams defines model for ProvisionParams.
type ProvisionParams struct {
	Name         string     `json:"Name"`
	ProviderName string     `json:"ProviderName"`
	Purpose      string     `json:"Purpose"`
	Start        time.Time  `json:"Start"`
	Stop         *time.Time `json:"Stop,omitempty"`
}

// ProvisionStatus defines model for ProvisionStatus.
type ProvisionStatus struct {
	GUID                string  `json:"GUID"`
	LabUserInterfaceUrl *string `json:"labUserInterfaceUrl,omitempty"`
	LifespanEnd         string  `json:"lifespanEnd"`
	RuntimeDefault      string  `json:"runtimeDefault"`
	RuntimeMaximum      string  `json:"runtimeMaximum"`
	State               string  `json:"state"`
}

// RatingDetails defines model for RatingDetails.
type RatingDetails struct {
	Comment       *string `json:"Comment,omitempty"`
	Email         string  `json:"Email"`
	ProvisionName string  `json:"ProvisionName"`
	Rating        int     `json:"Rating"`
	Useful        *string `json:"Useful,omitempty"`
}

// Status defines model for Status.
type Status struct {
	Message *string       `json:"message,omitempty"`
	Status  *StatusStatus `json:"status,omitempty"`
}

// StatusStatus defines model for Status.Status.
type StatusStatus string

// CreateServiceRequestParams defines parameters for CreateServiceRequest.
type CreateServiceRequestParams struct {
	// XGrecaptchaToken Google Recaptcha Token
	XGrecaptchaToken *string `json:"X-Grecaptcha-Token,omitempty"`
}

// CreateRatingJSONRequestBody defines body for CreateRating for application/json ContentType.
type CreateRatingJSONRequestBody = RatingDetails

// CreateServiceRequestJSONRequestBody defines body for CreateServiceRequest for application/json ContentType.
type CreateServiceRequestJSONRequestBody = ProvisionParams

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List all Catalog Items
	// (GET /catalogItems)
	ListCatalogItems(w http.ResponseWriter, r *http.Request)
	// Get a Catalog Item
	// (GET /catalogItems/{name})
	GetCatalogItem(w http.ResponseWriter, r *http.Request, name string)
	// Health check
	// (GET /health)
	Health(w http.ResponseWriter, r *http.Request)
	// Create a new rating for a provisioned item
	// (POST /ratings)
	CreateRating(w http.ResponseWriter, r *http.Request)
	// Get Catalog Item ratings
	// (GET /ratings/{name})
	GetRating(w http.ResponseWriter, r *http.Request, name string)
	// Create a service request for a new provision
	// (POST /serviceRequest)
	CreateServiceRequest(w http.ResponseWriter, r *http.Request, params CreateServiceRequestParams)
	// Service request to destroy provision
	// (DELETE /serviceRequest/{name})
	DeleteServiceRequest(w http.ResponseWriter, r *http.Request, name string)
	// Get status of service request for provision
	// (GET /serviceRequest/{name})
	GetServiceRequestStatus(w http.ResponseWriter, r *http.Request, name string)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// List all Catalog Items
// (GET /catalogItems)
func (_ Unimplemented) ListCatalogItems(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get a Catalog Item
// (GET /catalogItems/{name})
func (_ Unimplemented) GetCatalogItem(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Health check
// (GET /health)
func (_ Unimplemented) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create a new rating for a provisioned item
// (POST /ratings)
func (_ Unimplemented) CreateRating(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get Catalog Item ratings
// (GET /ratings/{name})
func (_ Unimplemented) GetRating(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create a service request for a new provision
// (POST /serviceRequest)
func (_ Unimplemented) CreateServiceRequest(w http.ResponseWriter, r *http.Request, params CreateServiceRequestParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Service request to destroy provision
// (DELETE /serviceRequest/{name})
func (_ Unimplemented) DeleteServiceRequest(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get status of service request for provision
// (GET /serviceRequest/{name})
func (_ Unimplemented) GetServiceRequestStatus(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// ListCatalogItems operation middleware
func (siw *ServerInterfaceWrapper) ListCatalogItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListCatalogItems(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetCatalogItem operation middleware
func (siw *ServerInterfaceWrapper) GetCatalogItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, chi.URLParam(r, "name"), &name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetCatalogItem(w, r, name)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Health operation middleware
func (siw *ServerInterfaceWrapper) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Health(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateRating operation middleware
func (siw *ServerInterfaceWrapper) CreateRating(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateRating(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetRating operation middleware
func (siw *ServerInterfaceWrapper) GetRating(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, chi.URLParam(r, "name"), &name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRating(w, r, name)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateServiceRequest operation middleware
func (siw *ServerInterfaceWrapper) CreateServiceRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params CreateServiceRequestParams

	headers := r.Header

	// ------------- Optional header parameter "X-Grecaptcha-Token" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Grecaptcha-Token")]; found {
		var XGrecaptchaToken string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "X-Grecaptcha-Token", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Grecaptcha-Token", runtime.ParamLocationHeader, valueList[0], &XGrecaptchaToken)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X-Grecaptcha-Token", Err: err})
			return
		}

		params.XGrecaptchaToken = &XGrecaptchaToken

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateServiceRequest(w, r, params)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteServiceRequest operation middleware
func (siw *ServerInterfaceWrapper) DeleteServiceRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, chi.URLParam(r, "name"), &name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteServiceRequest(w, r, name)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetServiceRequestStatus operation middleware
func (siw *ServerInterfaceWrapper) GetServiceRequestStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, chi.URLParam(r, "name"), &name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetServiceRequestStatus(w, r, name)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/catalogItems", wrapper.ListCatalogItems)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/catalogItems/{name}", wrapper.GetCatalogItem)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/health", wrapper.Health)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/ratings", wrapper.CreateRating)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/ratings/{name}", wrapper.GetRating)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/serviceRequest", wrapper.CreateServiceRequest)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/serviceRequest/{name}", wrapper.DeleteServiceRequest)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/serviceRequest/{name}", wrapper.GetServiceRequestStatus)
	})

	return r
}

type ListCatalogItemsRequestObject struct {
}

type ListCatalogItemsResponseObject interface {
	VisitListCatalogItemsResponse(w http.ResponseWriter) error
}

type ListCatalogItems200JSONResponse []CatalogItem

func (response ListCatalogItems200JSONResponse) VisitListCatalogItemsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetCatalogItemRequestObject struct {
	Name string `json:"name"`
}

type GetCatalogItemResponseObject interface {
	VisitGetCatalogItemResponse(w http.ResponseWriter) error
}

type GetCatalogItem200JSONResponse CatalogItem

func (response GetCatalogItem200JSONResponse) VisitGetCatalogItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetCatalogItem404Response struct {
}

func (response GetCatalogItem404Response) VisitGetCatalogItemResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type GetCatalogItem500JSONResponse Error

func (response GetCatalogItem500JSONResponse) VisitGetCatalogItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type HealthRequestObject struct {
}

type HealthResponseObject interface {
	VisitHealthResponse(w http.ResponseWriter) error
}

type Health200JSONResponse Status

func (response Health200JSONResponse) VisitHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type CreateRatingRequestObject struct {
	Body *CreateRatingJSONRequestBody
}

type CreateRatingResponseObject interface {
	VisitCreateRatingResponse(w http.ResponseWriter) error
}

type CreateRating201Response struct {
}

func (response CreateRating201Response) VisitCreateRatingResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type CreateRating500JSONResponse Error

func (response CreateRating500JSONResponse) VisitCreateRatingResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetRatingRequestObject struct {
	Name string `json:"name"`
}

type GetRatingResponseObject interface {
	VisitGetRatingResponse(w http.ResponseWriter) error
}

type GetRating200JSONResponse CatalogItemRating

func (response GetRating200JSONResponse) VisitGetRatingResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetRating500JSONResponse Error

func (response GetRating500JSONResponse) VisitGetRatingResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type CreateServiceRequestRequestObject struct {
	Params CreateServiceRequestParams
	Body   *CreateServiceRequestJSONRequestBody
}

type CreateServiceRequestResponseObject interface {
	VisitCreateServiceRequestResponse(w http.ResponseWriter) error
}

type CreateServiceRequest201ResponseHeaders struct {
	Location string
}

type CreateServiceRequest201JSONResponse struct {
	Body    ProvisionInfo
	Headers CreateServiceRequest201ResponseHeaders
}

func (response CreateServiceRequest201JSONResponse) VisitCreateServiceRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprint(response.Headers.Location))
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response.Body)
}

type CreateServiceRequest401JSONResponse Error

func (response CreateServiceRequest401JSONResponse) VisitCreateServiceRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type CreateServiceRequest500JSONResponse Error

func (response CreateServiceRequest500JSONResponse) VisitCreateServiceRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type DeleteServiceRequestRequestObject struct {
	Name string `json:"name"`
}

type DeleteServiceRequestResponseObject interface {
	VisitDeleteServiceRequestResponse(w http.ResponseWriter) error
}

type DeleteServiceRequest204Response struct {
}

func (response DeleteServiceRequest204Response) VisitDeleteServiceRequestResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type DeleteServiceRequest500JSONResponse Error

func (response DeleteServiceRequest500JSONResponse) VisitDeleteServiceRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetServiceRequestStatusRequestObject struct {
	Name string `json:"name"`
}

type GetServiceRequestStatusResponseObject interface {
	VisitGetServiceRequestStatusResponse(w http.ResponseWriter) error
}

type GetServiceRequestStatus200JSONResponse ProvisionStatus

func (response GetServiceRequestStatus200JSONResponse) VisitGetServiceRequestStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetServiceRequestStatus202Response struct {
}

func (response GetServiceRequestStatus202Response) VisitGetServiceRequestStatusResponse(w http.ResponseWriter) error {
	w.WriteHeader(202)
	return nil
}

type GetServiceRequestStatus404Response struct {
}

func (response GetServiceRequestStatus404Response) VisitGetServiceRequestStatusResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type GetServiceRequestStatus500JSONResponse Error

func (response GetServiceRequestStatus500JSONResponse) VisitGetServiceRequestStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// List all Catalog Items
	// (GET /catalogItems)
	ListCatalogItems(ctx context.Context, request ListCatalogItemsRequestObject) (ListCatalogItemsResponseObject, error)
	// Get a Catalog Item
	// (GET /catalogItems/{name})
	GetCatalogItem(ctx context.Context, request GetCatalogItemRequestObject) (GetCatalogItemResponseObject, error)
	// Health check
	// (GET /health)
	Health(ctx context.Context, request HealthRequestObject) (HealthResponseObject, error)
	// Create a new rating for a provisioned item
	// (POST /ratings)
	CreateRating(ctx context.Context, request CreateRatingRequestObject) (CreateRatingResponseObject, error)
	// Get Catalog Item ratings
	// (GET /ratings/{name})
	GetRating(ctx context.Context, request GetRatingRequestObject) (GetRatingResponseObject, error)
	// Create a service request for a new provision
	// (POST /serviceRequest)
	CreateServiceRequest(ctx context.Context, request CreateServiceRequestRequestObject) (CreateServiceRequestResponseObject, error)
	// Service request to destroy provision
	// (DELETE /serviceRequest/{name})
	DeleteServiceRequest(ctx context.Context, request DeleteServiceRequestRequestObject) (DeleteServiceRequestResponseObject, error)
	// Get status of service request for provision
	// (GET /serviceRequest/{name})
	GetServiceRequestStatus(ctx context.Context, request GetServiceRequestStatusRequestObject) (GetServiceRequestStatusResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHttpHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHttpMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// ListCatalogItems operation middleware
func (sh *strictHandler) ListCatalogItems(w http.ResponseWriter, r *http.Request) {
	var request ListCatalogItemsRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ListCatalogItems(ctx, request.(ListCatalogItemsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListCatalogItems")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListCatalogItemsResponseObject); ok {
		if err := validResponse.VisitListCatalogItemsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetCatalogItem operation middleware
func (sh *strictHandler) GetCatalogItem(w http.ResponseWriter, r *http.Request, name string) {
	var request GetCatalogItemRequestObject

	request.Name = name

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetCatalogItem(ctx, request.(GetCatalogItemRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetCatalogItem")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetCatalogItemResponseObject); ok {
		if err := validResponse.VisitGetCatalogItemResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Health operation middleware
func (sh *strictHandler) Health(w http.ResponseWriter, r *http.Request) {
	var request HealthRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Health(ctx, request.(HealthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Health")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(HealthResponseObject); ok {
		if err := validResponse.VisitHealthResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// CreateRating operation middleware
func (sh *strictHandler) CreateRating(w http.ResponseWriter, r *http.Request) {
	var request CreateRatingRequestObject

	var body CreateRatingJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateRating(ctx, request.(CreateRatingRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateRating")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateRatingResponseObject); ok {
		if err := validResponse.VisitCreateRatingResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetRating operation middleware
func (sh *strictHandler) GetRating(w http.ResponseWriter, r *http.Request, name string) {
	var request GetRatingRequestObject

	request.Name = name

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetRating(ctx, request.(GetRatingRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetRating")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetRatingResponseObject); ok {
		if err := validResponse.VisitGetRatingResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// CreateServiceRequest operation middleware
func (sh *strictHandler) CreateServiceRequest(w http.ResponseWriter, r *http.Request, params CreateServiceRequestParams) {
	var request CreateServiceRequestRequestObject

	request.Params = params

	var body CreateServiceRequestJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateServiceRequest(ctx, request.(CreateServiceRequestRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateServiceRequest")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateServiceRequestResponseObject); ok {
		if err := validResponse.VisitCreateServiceRequestResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// DeleteServiceRequest operation middleware
func (sh *strictHandler) DeleteServiceRequest(w http.ResponseWriter, r *http.Request, name string) {
	var request DeleteServiceRequestRequestObject

	request.Name = name

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteServiceRequest(ctx, request.(DeleteServiceRequestRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteServiceRequest")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(DeleteServiceRequestResponseObject); ok {
		if err := validResponse.VisitDeleteServiceRequestResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetServiceRequestStatus operation middleware
func (sh *strictHandler) GetServiceRequestStatus(w http.ResponseWriter, r *http.Request, name string) {
	var request GetServiceRequestStatusRequestObject

	request.Name = name

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetServiceRequestStatus(ctx, request.(GetServiceRequestStatusRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetServiceRequestStatus")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetServiceRequestStatusResponseObject); ok {
		if err := validResponse.VisitGetServiceRequestStatusResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"AYvqHQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
