package authorization

import (
	"os"
	"testing"
)

// TestMain registers the meal-planning domain permissions into the platform role sets before
// any test in this package runs, mirroring what mealplanningregistration.RegisterForGRPCAPI
// does at application startup.
func TestMain(m *testing.M) {
	RegisterServiceAdminPermissions(MealPlanningServiceAdminPermissions...)
	RegisterServiceDataAdminPermissions(MealPlanningServiceDataAdminPermissions...)
	RegisterAccountAdminPermissions(MealPlanningAccountAdminPermissions...)
	RegisterAccountMemberPermissions(MealPlanningAccountMemberPermissions...)

	os.Exit(m.Run())
}
