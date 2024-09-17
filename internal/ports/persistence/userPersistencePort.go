package persistence

import (
	"user-auth-hexagonal-architecture/internal/domain"
)

// UserPersistencePort is a secondary (driven) port to decouple the core layer from the persistence layer
type UserPersistencePort interface {
	SaveUser(username string, hashedPassword string) error
	FindUser(username string) (domain.User, error)
	IsUsernameAvailable(username string) (bool, error)
}
