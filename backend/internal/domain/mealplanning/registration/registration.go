// Package registration provides consolidated DI registration functions for the mealplanning domain.
// Domain: mealplanning — remove this package when swapping the domain.
package registration

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/build/services/api/grpc/domainreg"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	domaindataprivacy "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/dataprivacy"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/grocerylistpreparation"
	mealplanningmgr "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	mealplanningprivacy "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/privacy"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	mealplanningsvcpb "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	mealplanningrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/auth/grpc/interceptors"
	mealplanningcfg "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/config"
	mealplanningsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/grpc"
	eatingindexing "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/indexing"
	mealplanfinalizer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"

	"github.com/primandproper/platform/observability/logging"
	"github.com/primandproper/platform/observability/tracing"

	"github.com/samber/do/v2"
	"google.golang.org/grpc"
)

func registerRepository(i do.Injector) {
	mealplanningrepo.RegisterMealPlanningRepository(i)
}

func registerDataPrivacyCollector(i do.Injector) {
	do.Provide[[]domaindataprivacy.UserDataCollector](i, func(i do.Injector) ([]domaindataprivacy.UserDataCollector, error) {
		return []domaindataprivacy.UserDataCollector{
			mealplanningprivacy.NewCollector(
				do.MustInvoke[mealplanning.Repository](i),
				do.MustInvoke[logging.Logger](i),
				do.MustInvoke[tracing.TracerProvider](i),
			),
		}, nil
	})
}

// RegisterForGRPCAPI registers all mealplanning components needed by the gRPC API server.
func RegisterForGRPCAPI(i do.Injector) {
	// Register meal-planning permissions into the platform role sets.
	authorization.RegisterServiceAdminPermissions(authorization.MealPlanningServiceAdminPermissions...)
	authorization.RegisterServiceDataAdminPermissions(authorization.MealPlanningServiceDataAdminPermissions...)
	authorization.RegisterAccountAdminPermissions(authorization.MealPlanningAccountAdminPermissions...)
	authorization.RegisterAccountMemberPermissions(authorization.MealPlanningAccountMemberPermissions...)

	registerRepository(i)
	registerDataPrivacyCollector(i)
	mealplanningmgr.RegisterManagers(i)
	mealplanningsvc.RegisterMealPlanningService(i)
	mealplanfinalizer.RegisterMealPlanFinalizer(i)
	mealplangrocerylistinitializer.RegisterMealPlanGroceryListInitializer(i)
	mealplantaskcreator.RegisterMealPlanTaskCreator(i)
	recipeanalysis.RegisterRecipeAnalyzer(i)
	grocerylistpreparation.RegisterGroceryListCreator(i)

	// Provide the mealplanning service config (previously wired in grpcapi/config.go).
	do.Provide[*mealplanningcfg.Config](i, func(i do.Injector) (*mealplanningcfg.Config, error) {
		svc := do.MustInvoke[*config.ServicesConfig](i)
		return &svc.MealPlanning, nil
	})

	// Provide the gRPC registration func so extras.go can append it without knowing the type.
	do.Provide[domainreg.ExtraRegistrationFuncs](i, func(i do.Injector) (domainreg.ExtraRegistrationFuncs, error) {
		svc := do.MustInvoke[mealplanningsvcpb.MealPlanningServiceServer](i)
		return domainreg.ExtraRegistrationFuncs{
			func(s *grpc.Server) {
				mealplanningsvcpb.RegisterMealPlanningServiceServer(s, svc)
			},
		}, nil
	})

	// Provide method permissions so extras.go can merge them without knowing the type.
	do.Provide[domainreg.ExtraMethodPermissions](i, func(i do.Injector) (domainreg.ExtraMethodPermissions, error) {
		perms := do.MustInvoke[mealplanningsvc.MealPlanningMethodPermissions](i)
		return domainreg.ExtraMethodPermissions(interceptors.MethodPermissionsMap(perms)), nil
	})
}

// RegisterForDataChangeHandler registers mealplanning components needed by the async message handler.
func RegisterForDataChangeHandler(i do.Injector) {
	registerRepository(i)
	registerDataPrivacyCollector(i)
	eatingindexing.RegisterMealPlanningSearchers(i)
	eatingindexing.RegisterMealPlanningDataIndexer(i)
}

// RegisterForSearchIndexScheduler registers mealplanning components needed by the search index scheduler.
func RegisterForSearchIndexScheduler(i do.Injector) {
	registerRepository(i)
	eatingindexing.RegisterMealPlanningDataIndexer(i)
}

// RegisterForMCPServer registers mealplanning components needed by the MCP server.
func RegisterForMCPServer(i do.Injector) {
	registerRepository(i)
}
