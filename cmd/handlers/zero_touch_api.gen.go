// Package handlers provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
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

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
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
	Name         string    `json:"Name"`
	ProviderName string    `json:"ProviderName"`
	Purpose      string    `json:"Purpose"`
	Start        time.Time `json:"Start"`
	Stop         time.Time `json:"Stop"`
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

// Status defines model for Status.
type Status struct {
	Message *string       `json:"message,omitempty"`
	Status  *StatusStatus `json:"status,omitempty"`
}

// StatusStatus defines model for Status.Status.
type StatusStatus string

// CreateProvisionParams defines parameters for CreateProvision.
type CreateProvisionParams struct {
	// XGrecaptchaToken Google Recaptcha Token
	XGrecaptchaToken *string `json:"X-Grecaptcha-Token,omitempty"`
}

// CreateProvisionJSONRequestBody defines body for CreateProvision for application/json ContentType.
type CreateProvisionJSONRequestBody = ProvisionParams

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
	// Create new provision
	// (POST /provision)
	CreateProvision(w http.ResponseWriter, r *http.Request, params CreateProvisionParams)
	// Destroy provision
	// (DELETE /provision/{name})
	DeleteProvision(w http.ResponseWriter, r *http.Request, name string)
	// Get provision status
	// (GET /provision/{name})
	GetProvisionStatus(w http.ResponseWriter, r *http.Request, name string)
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

// Create new provision
// (POST /provision)
func (_ Unimplemented) CreateProvision(w http.ResponseWriter, r *http.Request, params CreateProvisionParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Destroy provision
// (DELETE /provision/{name})
func (_ Unimplemented) DeleteProvision(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get provision status
// (GET /provision/{name})
func (_ Unimplemented) GetProvisionStatus(w http.ResponseWriter, r *http.Request, name string) {
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

// CreateProvision operation middleware
func (siw *ServerInterfaceWrapper) CreateProvision(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params CreateProvisionParams

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
		siw.Handler.CreateProvision(w, r, params)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteProvision operation middleware
func (siw *ServerInterfaceWrapper) DeleteProvision(w http.ResponseWriter, r *http.Request) {
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
		siw.Handler.DeleteProvision(w, r, name)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetProvisionStatus operation middleware
func (siw *ServerInterfaceWrapper) GetProvisionStatus(w http.ResponseWriter, r *http.Request) {
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
		siw.Handler.GetProvisionStatus(w, r, name)
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
		r.Post(options.BaseURL+"/provision", wrapper.CreateProvision)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/provision/{name}", wrapper.DeleteProvision)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/provision/{name}", wrapper.GetProvisionStatus)
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

type CreateProvisionRequestObject struct {
	Params CreateProvisionParams
	Body   *CreateProvisionJSONRequestBody
}

type CreateProvisionResponseObject interface {
	VisitCreateProvisionResponse(w http.ResponseWriter) error
}

type CreateProvision201ResponseHeaders struct {
	Location string
}

type CreateProvision201JSONResponse struct {
	Body    ProvisionInfo
	Headers CreateProvision201ResponseHeaders
}

func (response CreateProvision201JSONResponse) VisitCreateProvisionResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprint(response.Headers.Location))
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response.Body)
}

type CreateProvision401JSONResponse Error

func (response CreateProvision401JSONResponse) VisitCreateProvisionResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type CreateProvision500JSONResponse Error

func (response CreateProvision500JSONResponse) VisitCreateProvisionResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type DeleteProvisionRequestObject struct {
	Name string `json:"name"`
}

type DeleteProvisionResponseObject interface {
	VisitDeleteProvisionResponse(w http.ResponseWriter) error
}

type DeleteProvision204Response struct {
}

func (response DeleteProvision204Response) VisitDeleteProvisionResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type DeleteProvision500JSONResponse Error

func (response DeleteProvision500JSONResponse) VisitDeleteProvisionResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetProvisionStatusRequestObject struct {
	Name string `json:"name"`
}

type GetProvisionStatusResponseObject interface {
	VisitGetProvisionStatusResponse(w http.ResponseWriter) error
}

type GetProvisionStatus200JSONResponse ProvisionStatus

func (response GetProvisionStatus200JSONResponse) VisitGetProvisionStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetProvisionStatus202Response struct {
}

func (response GetProvisionStatus202Response) VisitGetProvisionStatusResponse(w http.ResponseWriter) error {
	w.WriteHeader(202)
	return nil
}

type GetProvisionStatus404Response struct {
}

func (response GetProvisionStatus404Response) VisitGetProvisionStatusResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type GetProvisionStatus500JSONResponse Error

func (response GetProvisionStatus500JSONResponse) VisitGetProvisionStatusResponse(w http.ResponseWriter) error {
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
	// Create new provision
	// (POST /provision)
	CreateProvision(ctx context.Context, request CreateProvisionRequestObject) (CreateProvisionResponseObject, error)
	// Destroy provision
	// (DELETE /provision/{name})
	DeleteProvision(ctx context.Context, request DeleteProvisionRequestObject) (DeleteProvisionResponseObject, error)
	// Get provision status
	// (GET /provision/{name})
	GetProvisionStatus(ctx context.Context, request GetProvisionStatusRequestObject) (GetProvisionStatusResponseObject, error)
}

type StrictHandlerFunc = runtime.StrictHttpHandlerFunc
type StrictMiddlewareFunc = runtime.StrictHttpMiddlewareFunc

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
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
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
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
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
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// CreateProvision operation middleware
func (sh *strictHandler) CreateProvision(w http.ResponseWriter, r *http.Request, params CreateProvisionParams) {
	var request CreateProvisionRequestObject

	request.Params = params

	var body CreateProvisionJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateProvision(ctx, request.(CreateProvisionRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateProvision")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateProvisionResponseObject); ok {
		if err := validResponse.VisitCreateProvisionResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// DeleteProvision operation middleware
func (sh *strictHandler) DeleteProvision(w http.ResponseWriter, r *http.Request, name string) {
	var request DeleteProvisionRequestObject

	request.Name = name

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteProvision(ctx, request.(DeleteProvisionRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteProvision")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(DeleteProvisionResponseObject); ok {
		if err := validResponse.VisitDeleteProvisionResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// GetProvisionStatus operation middleware
func (sh *strictHandler) GetProvisionStatus(w http.ResponseWriter, r *http.Request, name string) {
	var request GetProvisionStatusRequestObject

	request.Name = name

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetProvisionStatus(ctx, request.(GetProvisionStatusRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetProvisionStatus")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetProvisionStatusResponseObject); ok {
		if err := validResponse.VisitGetProvisionStatusResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"JoYbKhufx21Fc7O92/4VAAD//+nuWF0JFwAA",
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
