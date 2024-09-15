// Package service implements the application's business logic and use cases.
// It acts as an intermediary between the adapter layer and the domain layer,
package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"user-auth-hexagonal-architecture/internal/ports/persistence"
)

// RegisterUserService handles the business logic for user registration.
// It implements the RegisterUserPort interface from the usecases package.
type RegisterUserService struct {
	userPersistence persistence.UserPersistencePort
}

// NewRegisterUserService creates a new instance of RegisterUserService.
//
// Parameters:
//   - userPersistence: An implementation of UserPersistencePort for storing user data
//
// Returns:
//   - *RegisterUserService: A pointer to the newly created RegisterUserService
func NewRegisterUserService(userPersistence persistence.UserPersistencePort) *RegisterUserService {
	return &RegisterUserService{userPersistence}
}

// RegisterUser handles the registration of a new user.
//
// This method performs the following steps:
// 1. Hashes the provided password using bcrypt
// 2. Saves the user's username and hashed password using the persistence layer
//
// Parameters:
//   - username: The username for the new user
//   - password: The plain text password for the new user
//
// Returns:
//   - error: An error if registration fails, nil otherwise
//
// Possible errors:
//   - If password hashing fails
//   - If saving the user to the persistence layer fails
//
// Note: This method uses bcrypt's DefaultCost for password hashing.
func (lu *RegisterUserService) RegisterUser(username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return lu.userPersistence.SaveUser(username, string(hashedPassword))
}
