package authentication

type User interface {
	GetID() string
	GetUsername() string
	GetEmail() string
	GetFirstName() string
	GetLastName() string
}
