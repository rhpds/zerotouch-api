// Package handlers provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
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
	Rating        *int    `json:"Rating,omitempty"`
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
	"HQAA",
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
