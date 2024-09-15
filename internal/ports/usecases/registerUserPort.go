package usecases

// RegisterUserPort is a primary (driving) port to decouple the core layer from the adapter layer
type RegisterUserPort interface {
	RegisterUser(username string, password string) error
}
