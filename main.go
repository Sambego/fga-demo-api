package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	openfgaclient "github.com/openfga/go-sdk/client"
	"github.com/sambego/fga-demo-api/data"
	"github.com/sambego/fga-demo-api/router"
)

func main() {
	// Load the .env file
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Check for environment variables
	apiUrl, ok := os.LookupEnv("FGA_API_URL")
	if !ok {
		log.Fatal("'FGA_STORE_ID' environment variable must be set")
	}

	// Initialise the Auth0 FGA middleware
	fgaClient, err := openfgaclient.NewSdkClient(&openfgaclient.ClientConfiguration{
		ApiUrl: apiUrl,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("FGA API URL: %v", apiUrl)

	// Create a new FGA Store
	body := openfgaclient.ClientCreateStoreRequest{Name: "FGA Demo"}
	fgaStore, err := fgaClient.CreateStore(context.Background()).Body(body).Execute()

	// Set the new store as the store ID for the Open FGA Client
	fgaClient.SetStoreId(fgaStore.Id)
	log.Printf("NEW FGA STORE ID: %v", fgaStore.Id)

	if err != nil {
		log.Fatal(err)
	}

	// Read the fga-model.json file, containing our model
	jsonData, err := os.ReadFile("./fga-model.json")

	if err != nil {
		log.Fatal(err)
	}

	// Convert to a ClientWriteAuthorizationModelRequest
	var fgaModel openfgaclient.ClientWriteAuthorizationModelRequest
	err = json.Unmarshal(jsonData, &fgaModel)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new FGA Model with the contents of our fga-model.json file
	fgaModelData, err := fgaClient.WriteAuthorizationModel(context.Background()).Body(fgaModel).Execute()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("NEW FGA MODEL ID: %v", fgaModelData.AuthorizationModelId)

	// Create a new data store
	var documents []data.Document
	store := data.Store{
		Documents: documents,
	}

	// Create a new routerservice
	service := router.Service{
		Store:     &store,
		FGAClient: fgaClient,
	}

	// Create router and routes
	r := router.CreateRouter(service)

	// Port
	var port = ":4000"

	// Start the server
	srv := &http.Server{
		Addr:         port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Printf("SERVER URL: http://localhost%s", port)

	// Error handling
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
