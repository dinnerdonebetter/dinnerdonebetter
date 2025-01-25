package authentication

type User interface {
	GetID() string
	GetEmail() string
	GetFirstName() string
	GetLastName() string
}
