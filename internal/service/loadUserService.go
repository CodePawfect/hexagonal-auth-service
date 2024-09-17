// Package service implements the application's business logic and use cases.
// It acts as an intermediary between the adapter layer and the domain layer,
package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
	"user-auth-hexagonal-architecture/internal/ports/persistence"
)

// LoadUserService handles the business logic for user authentication.
// It implements the LoadUserPort interface from the usecases package.
type LoadUserService struct {
	userPersistence persistence.UserPersistencePort
}

// NewLoadUserService creates a new instance of LoadUserService.
//
// Parameters:
//   - userPersistence: An implementation of UserPersistencePort for retrieving user data
//
// Returns:
//   - *LoadUserService: A pointer to the newly created LoadUserService
func NewLoadUserService(userPersistence persistence.UserPersistencePort) *LoadUserService {
	return &LoadUserService{userPersistence}
}

// LoadUser authenticates a user and generates a JWT token upon successful authentication.
//
// This method performs the following steps:
// 1. Retrieves the user from the persistence layer using the provided username.
// 2. Compares the provided password with the stored (hashed) password.
// 3. If authentication is successful, generates a JWT token with user claims.
//
// Parameters:
//   - username: A string representing the username of the user to authenticate.
//   - password: A string representing the password to verify.
//
// Returns:
//   - string: A signed JWT token string if authentication is successful.
//   - error: An error in the following cases:
//   - If the user is not found in the persistence layer.
//   - If the provided password doesn't match the stored password.
//   - If there's an error during password comparison.
//   - If there's an error while creating or signing the JWT token.
//
// The JWT token includes the following claims:
//   - username: The authenticated user's username.
//   - role: The user's role.
//   - exp: The expiration time of the token (set to 24 hours from creation).
//
// Note:
//   - This method uses bcrypt for password comparison.
//   - The JWT signing key is hardcoded for demonstration purposes.
//     In a production environment, this should be securely managed.
//   - Error messages for authentication failures are intentionally vague
//     to prevent information leakage.
func (lu *LoadUserService) LoadUser(username string, password string) (string, error) {
	user, err := lu.userPersistence.FindUser(username)
	if err != nil {
		return "", fmt.Errorf("error finding user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", fmt.Errorf("invalid username or password")
		}
		return "", fmt.Errorf("error comparing passwords: %w", err)
	}

	var jwtKey = []byte("my_secret_key") // This is only for demo purposes
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	signedString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error while creating jwt: %w", err)
	}

	return signedString, nil
}
