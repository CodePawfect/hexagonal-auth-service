package persistence

// UserPersistencePort is a secondary (driven) port to decouple the core layer from the persistence layer
type UserPersistencePort interface {
	SaveUser(username string, hashedPassword string) error
	LoadUser(username string, password string) (string, error)
	IsUsernameAvailable(username string) (bool, error)
}
