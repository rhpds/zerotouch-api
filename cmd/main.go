package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/rhpds/zerotouch-api/cmd/handlers"
	"github.com/rhpds/zerotouch-api/cmd/log"
	"github.com/rhpds/zerotouch-api/cmd/models"
)

func main() {
	log.InitLoggers(true)

	// Loading the .env file. Environemnt variables will overwrite variables
	// defined at the .env file
	err := godotenv.Load()
	if err != nil {
		log.Logger.Info(err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//------------------
	// OpenAPI validation
	//------------------
	swagger, err := handlers.GetSwagger()
	if err != nil {
		log.Err.Fatal("Error loading swagger spec", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "api_key", "Authorization", "X-Grecaptcha-Token"},
		// ExposedHeaders:   []string{"Link"},
		// AllowCredentials: false,
		MaxAge: 300, // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/", mainRouter(swagger))
	r.Mount("/swagger.json", swaggerSpecRouter(swagger))

	log.Logger.Info("Starting server on port " + port)
	log.Err.Fatal(http.ListenAndServe(":"+port, r))
}

func mainRouter(swagger *openapi3.T) http.Handler {
	//
	// Get k8s configuration to authenticate server
	//
	kuebeconfig := os.Getenv("KUBECONFIG")
	if kuebeconfig == "" {
		log.Logger.Info("KUBECONFIG not set, using in-cluster config")
	} else {
		log.Logger.Info("Using KUBECONFIG: " + kuebeconfig)
	}

	//
	// Create controllers for CatalogItems and resourceClaims
	//
	rcNamespace := os.Getenv("RESOURCECLAIM_NAMESPACE")
	if rcNamespace == "" {
		log.Err.Fatal("RESOURCECLAIM_NAMESPACE not set, exiting")
	}
	log.Logger.Info("Using RESOURCECLAIM_NAMESPACE: " + rcNamespace)

	catalogItemsController, err := models.NewCatalogItemsController(
		kuebeconfig,
		context.Background(),
		"",
	)
	if err != nil {
		log.Err.Fatal("Error creating catalog items controller", err)
	}

	resourceClaimsController, err := models.NewResourceClaimsController(
		kuebeconfig,
		rcNamespace,
		context.Background(),
	)
	if err != nil {
		log.Err.Fatal("Error creating resource claims controller", err)
	}

	//
	// Get configurations for the Google Recaptcha Enterprise
	//

	recaptchaConfig := handlers.RecaptchaConfig{
		ProjectID:        "rhdp-ui",
		RecapthcaSiteKey: "6LfGWQcoAAAAABsPO_pHXwQK0HAaRXn5FL2eWN-O",
		Threshold:        0.5,
		Disabled:         false,
	}

	recaptchaConfig.Disabled, _ = strconv.ParseBool(os.Getenv("RECAPTCHA_DISABLED"))
	if recaptchaConfig.Disabled {
		log.Logger.Warn(
			"Recaptcha validation disabled, make sure that RECAPTCHA_DISABLED environment variable is not set for production use.",
		)
	} else {

		recaptchaConfig.AuthKey = os.Getenv("RECAPTCHA_AUTH_KEY")
		if recaptchaConfig.AuthKey == "" {
			log.Err.Fatal("RECAPTCHA_AUTH_KEY not set, exiting")
		}

		if len(os.Getenv("RECAPTCHA_THRESHOLD")) > 0 {
			threshold, err := strconv.ParseFloat(os.Getenv("RECAPTCHA_THRESHOLD"), 64)
			if err != nil {
				log.Logger.Error("Incorrect RECAPTCHA_THRESHOLD value, using default")
			} else {
				recaptchaConfig.Threshold = threshold
			}
		}

		log.Logger.Info("Google Recaptcha config",
			"recaptchaConfig.ProjectID", recaptchaConfig.ProjectID,
			"recaptchaConfig.RecapthcaSiteKey", recaptchaConfig.RecapthcaSiteKey,
			"recaptchaConfig.Threshold", fmt.Sprintf("%f", recaptchaConfig.Threshold),
			"recaptchaConfig.Debug", fmt.Sprintf("%v", recaptchaConfig.Disabled),
		)
	}
	//
	// Create Handler
	//
	catalogItemsHandler := handlers.NewCatalogItemsHandler(
		catalogItemsController,
		resourceClaimsController,
		&recaptchaConfig,
	)

	strictHandler := handlers.NewStrictHandler(catalogItemsHandler, nil)
	r := chi.NewRouter()

	// Use validation middleware to validate requests
	r.Use(middleware.OapiRequestValidator(swagger))

	// Register handlers
	handlers.HandlerFromMux(strictHandler, r)

	return r
}

func swaggerSpecRouter(swagger *openapi3.T) http.Handler {
	data, _ := swagger.MarshalJSON()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	})
	return r
}
