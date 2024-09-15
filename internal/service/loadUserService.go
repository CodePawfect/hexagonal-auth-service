package service

import "user-auth-hexagonal-architecture/internal/ports/persistence"

type LoadUserService struct {
	userPersistence persistence.UserPersistencePort
}

func NewLoadUserService(userPersistence persistence.UserPersistencePort) *LoadUserService {
	return &LoadUserService{userPersistence}
}

func (lu *LoadUserService) LoadUser(username string, password string) (bool, error) {
	//TODO implement
	return true, nil
}
