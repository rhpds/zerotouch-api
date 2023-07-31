package main

import (
	"net/http"
	"os"
	"strconv"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rhpds/zerotouch-api/cmd/handlers"
	"github.com/rhpds/zerotouch-api/cmd/log"
	"github.com/rhpds/zerotouch-api/cmd/models"
)

func main() {
	log.InitLoggers(true)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	kuebeconfig := os.Getenv("KUBECONFIG")
	if kuebeconfig == "" {
		log.Logger.Info("KUBECONFIG not set, using in-cluster config")
	} else {
		log.Logger.Info("Using KUBECONFIG: " + kuebeconfig)
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

	// Create an instance of the API handler that satisfies the generated interface
	catalogItemRepo := models.NewCatalogItemRepo()
	itemsLoaded, err := catalogItemRepo.Refresh(kuebeconfig)
	if err != nil {
		log.Err.Fatal("Error loading catalog items", err)
	}

	log.Logger.Info("Loaded " + strconv.Itoa(itemsLoaded) + " catalog items")

	catalogItemsHandler := handlers.NewCatalogItemsHandler(catalogItemRepo) // TODO: create new model controller and attache it to the api handler

	strictHandler := handlers.NewStrictHandler(catalogItemsHandler, nil)
	r := chi.NewRouter()

	// Use validation middleware to validate requests
	r.Use(middleware.OapiRequestValidator(swagger))

	// Register handlers
	handlers.HandlerFromMux(strictHandler, r)

	log.Logger.Info("Starting server on port " + port)
	log.Err.Fatal(http.ListenAndServe(":"+port, r))
}
