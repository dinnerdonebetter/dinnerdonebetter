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
	return newID(true)
}
