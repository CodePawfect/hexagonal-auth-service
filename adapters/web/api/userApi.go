// Package api provides HTTP handlers for domain-related operations in a hexagonal architecture.
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"user-auth-hexagonal-architecture/internal/ports/usecases"
)

// UserApi handles HTTP requests for user operations.
// It acts as an adapter between the HTTP layer and the application's use cases.
type UserApi struct {
	registerUserPort usecases.RegisterUserPort
	loadUserPort     usecases.LoadUserPort
}

// userRequest represents the expected JSON structure for user registration requests.
type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewUserApiAdapter creates a new UserApi with the given use case ports.
//
// Parameters:
//   - registerUserPort: Port for user registration use case
//   - loadUserPort: Port for user loading use case
//
// Returns:
//   - *UserApi: A pointer to the newly created UserApi
func NewUserApiAdapter(registerUserPort usecases.RegisterUserPort, loadUserPort usecases.LoadUserPort) *UserApi {
	return &UserApi{registerUserPort, loadUserPort}
}

// InitUserRoutes sets up the HTTP routes for user-related operations.
//
// This method registers the necessary HTTP handlers with the given ServeMux.
func (ua *UserApi) InitUserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /user/register", ua.handleUserRegister)
	//TODO: add route for loading a user
}

// handleUserRegister handles HTTP POST requests for user registration.
//
// It decodes the JSON request body, calls the RegisterUser use case,
// and responds with appropriate HTTP status codes.
//
// The function expects a JSON body with "username" and "password" fields.
// On success, it responds with HTTP 201 Created.
// On failure, it responds with either 400 Bad Request for invalid JSON
// or 500 Internal Server Error for registration failures.
//
// Parameters:
//   - w: HTTP ResponseWriter to write the response
//   - r: HTTP Request containing the registration data
//
// Note: This method logs errors but does not return them to the caller.
func (ua *UserApi) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	var userRequest userRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err = ua.registerUserPort.RegisterUser(userRequest.Username, userRequest.Password)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Registering new user failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (ua *UserApi) handleLoadUser(w http.ResponseWriter, r *http.Request) {
	//TODO: implement handler for api route for loading a user
}
