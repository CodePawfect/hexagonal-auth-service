package api

import (
	"encoding/json"
	"log"
	"net/http"
	"user-auth-hexagonal-architecture/internal/ports/usecases"
)

type UserAdapter struct {
	registerUserPort usecases.RegisterUserPort
	loadUserPort     usecases.LoadUserPort
}

type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserApiAdapter(registerUserPort usecases.RegisterUserPort, loadUserPort usecases.LoadUserPort) *UserAdapter {
	return &UserAdapter{registerUserPort, loadUserPort}
}

func (ua *UserAdapter) InitUserRoutes() {
	http.HandleFunc("POST /user/register", ua.handleUserRegister)
}

func (ua *UserAdapter) handleUserRegister(w http.ResponseWriter, r *http.Request) {
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
