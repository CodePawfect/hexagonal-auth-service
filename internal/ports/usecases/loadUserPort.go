package usecases

// LoadUserPort is a primary (driving) port to decouple the core layer from the adapter layer
type LoadUserPort interface {
	LoadUser(username string, password string) (string, error)
}
