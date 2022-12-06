package identifiers

import (
	"github.com/rs/xid"
)

// New produces a new string ID.
func New() string {
	return xid.New().String()
}

func parseID(x string) error {
	var err error

	_, err = xid.FromString(x)
	return err
}

// Validate validates a string ID.
func Validate(x string) error {
	return parseID(x)
}
