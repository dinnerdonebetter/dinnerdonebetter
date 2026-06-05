package grpc

import (
	"os"
	"testing"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authorization"
)

// TestMain registers all domain permissions into the platform role sets before any test in this
// package runs, mirroring what authorization.RegisterCoreDomainPermissions() does at application startup.
func TestMain(m *testing.M) {
	authorization.RegisterCoreDomainPermissions()
	os.Exit(m.Run())
}
