package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	auth0fga "github.com/auth0-lab/fga-go-sdk"
	"github.com/joho/godotenv"
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
	storeID, ok := os.LookupEnv("FGA_STORE_ID")
	if !ok {
		log.Fatal("'FGA_STORE_ID' environment variable must be set")
	}

	clientID, ok := os.LookupEnv("FGA_CLIENT_ID")
	if !ok {
		log.Fatal("'FGA_CLIENT_ID' environment variable must be set")
	}

	clientSecret, ok := os.LookupEnv("FGA_CLIENT_SECRET")
	if !ok {
		log.Fatal("'FGA_CLIENT_SECRET' environment variable must be set")
	}

	// Initialise the Auth0 FGA middleware
	fgaConfig, err := auth0fga.NewConfiguration(auth0fga.UserConfiguration{
		StoreId:      storeID,
		ClientId:     clientID,
		ClientSecret: clientSecret,
		Environment:  "us",
	})

	if err != nil {
		log.Fatal("'Oops, %v", err)
	}

	log.Printf("-------- AUTH0 FGA --------")
	log.Printf("  - Store ID: %v", storeID)
	log.Printf("  - Cliend ID: %v", clientID)
	log.Printf("  - Client Secret: %v%v%v", clientSecret[0:4], strings.Repeat("*", len(clientSecret)-8), clientSecret[len(clientSecret)-4:])
	log.Printf("---------------------------")

	// Create a new data store
	var documents []data.Document
	store := data.Store{
		Documents: documents,
	}

	// Create a new routerservice
	service := router.Service{
		Store:     &store,
		FGAClient: auth0fga.NewAPIClient(fgaConfig).Auth0FgaApi,
	}

	// Create router and routes
	r := router.CreateRouter(service)

	// Start the server
	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

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
