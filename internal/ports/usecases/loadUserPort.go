package usecases

type LoadUserPort interface {
	LoadUser(username string, password string) (bool, error)
}
