package identifiers

import (
	"github.com/segmentio/ksuid"
)

// New produces a new string ID.
func New() string {
	return ksuid.New().String()
}
