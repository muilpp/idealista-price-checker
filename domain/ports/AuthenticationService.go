package ports

type AuthenticationService interface {
	GetToken() string
}
