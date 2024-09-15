package main

import (
	"log"
	"net/http"
	"user-auth-hexagonal-architecture/adapters/persistence/user"
	"user-auth-hexagonal-architecture/adapters/web/api"
	"user-auth-hexagonal-architecture/internal/service"
)

func main() {
	userPersistence, err := persistence.NewUserPersistenceAdapter("mongodb://user:password@localhost:27017", "demo")
	if err != nil {
		log.Fatalf("Failed to create user persistence adapter: %v", err)
	}

	registerUserService := service.NewRegisterUserService(userPersistence)
	loadUserService := service.NewLoadUserService(userPersistence)
	userApiAdapter := api.NewUserApiAdapter(registerUserService, loadUserService)
	userApiAdapter.InitUserRoutes()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
