package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
	"user-auth-hexagonal-architecture/adapters/persistence/user"
	"user-auth-hexagonal-architecture/adapters/web/api"
	"user-auth-hexagonal-architecture/internal/service"
)

func main() {
	// dependency injection brings ports and adapters together
	mongoClient := createMongoClient()
	userPersistence, err := persistence.NewUserPersistenceMongoAdapter(mongoClient, "demo")
	if err != nil {
		log.Fatalf("Failed to create user persistence adapter: %v", err)
	}

	registerUserService := service.NewRegisterUserService(userPersistence)
	loadUserService := service.NewLoadUserService(userPersistence)
	userApiAdapter := api.NewUserApiAdapter(registerUserService, loadUserService)

	mux := http.NewServeMux()
	userApiAdapter.InitUserRoutes(mux)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// createMongoClient creates a new MongoDB client and returns it.
func createMongoClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://user:password@localhost:27017"))
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %w", err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %w", err)
	}

	return mongoClient
}
