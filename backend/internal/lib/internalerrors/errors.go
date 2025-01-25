package internalerrors

import (
	"errors"
	"fmt"
)

var ErrCircuitBroken = errors.New("service circuit broken")

// NilConfigError returns a nil config error.
func NilConfigError(name string) error {
	return fmt.Errorf("nil config provided for %s", name)
}
