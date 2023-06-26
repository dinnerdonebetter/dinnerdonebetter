package cryptography

import (
	"errors"
)

var (
	ErrIncorrectKeyLength = errors.New("secret is not the right length")
)
