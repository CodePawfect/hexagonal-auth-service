package usecases

type RegisterUserPort interface {
	RegisterUser(username string, password string) error
}
