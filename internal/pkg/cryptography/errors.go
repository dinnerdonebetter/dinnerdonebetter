package cryptography

import (
	"errors"
)

var (
	ErrIncorrectSecretLength = errors.New("secret is not the right length")
)
