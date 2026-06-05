package authorization

import (
	"os"
	"testing"
)

// TestMain registers all domain permissions into the platform role sets before any test in this
// package runs, mirroring what authorization.RegisterCoreDomainPermissions() and
// mealplanningregistration.RegisterForGRPCAPI() do at application startup.
func TestMain(m *testing.M) {
	RegisterCoreDomainPermissions()

	RegisterServiceAdminPermissions(MealPlanningServiceAdminPermissions...)
	RegisterServiceDataAdminPermissions(MealPlanningServiceDataAdminPermissions...)
	RegisterAccountAdminPermissions(MealPlanningAccountAdminPermissions...)
	RegisterAccountMemberPermissions(MealPlanningAccountMemberPermissions...)

	os.Exit(m.Run())
}
