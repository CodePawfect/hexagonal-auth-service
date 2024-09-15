package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"user-auth-hexagonal-architecture/internal/ports/persistence"
)

type RegisterUserService struct {
	userPersistence persistence.UserPersistencePort
}

func NewRegisterUserService(userPersistence persistence.UserPersistencePort) *RegisterUserService {
	return &RegisterUserService{userPersistence}
}

func (lu *RegisterUserService) RegisterUser(username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return lu.userPersistence.SaveUser(username, string(hashedPassword))
}
