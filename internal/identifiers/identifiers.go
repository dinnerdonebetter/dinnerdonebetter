package identifiers

import (
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
)

func newID(useXID bool) string {
	if useXID {
		return xid.New().String()
	}
	return ksuid.New().String()
}

// New produces a new string ID.
func New() string {
	return newID(false)
}

func parseID(x string, useXID bool) error {
	var err error

	if useXID {
		_, err = xid.FromString(x)
		return err
	}

	_, err = ksuid.Parse(x)
	return err
}

// Validate validates a string ID.
func Validate(x string) error {
	return parseID(x, false)
}
