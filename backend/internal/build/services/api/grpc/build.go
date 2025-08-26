//go:build wireinject
// +build wireinject

package grpcapi

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/config"
	authmgr "github.com/dinnerdonebetter/backend/internal/domain/auth/managers"
	identitymgr "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	grocerylistpreparation "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/grocerylistpreparation"
	mealplanningmgr "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	recipeanalysis "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	oauthmgr "github.com/dinnerdonebetter/backend/internal/domain/oauth/manager"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/qrcodes"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/platform/server/grpc"
	auditrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	authrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auth"
	dataprivacysrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/dataprivacy"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	mealplanningrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	notificationsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/notifications"
	oauthrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
	settingsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/settings"
	webhooksrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"
	auditsvc "github.com/dinnerdonebetter/backend/internal/services/audit/grpc"
	authsvc "github.com/dinnerdonebetter/backend/internal/services/auth/grpc"
	dataprivacysvc "github.com/dinnerdonebetter/backend/internal/services/dataprivacy/grpc"
	identitysvc "github.com/dinnerdonebetter/backend/internal/services/identity/grpc"
	internalopssvc "github.com/dinnerdonebetter/backend/internal/services/internalops/grpc"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc"
	mealplanfinalizer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/services/notifications/grpc"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/services/oauth/grpc"
	settingssvc "github.com/dinnerdonebetter/backend/internal/services/settings/grpc"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/services/webhooks/grpc"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (*GRPCService, error) {
	wire.Build(ConfigProviders,
		// core
		metricscfg.Providers,
		loggingcfg.ProvidersLogConfig,
		tracingcfg.ProvidersTracingConfig,
		msgconfig.MessageQueueProviders,
		authentication.AuthProviders,
		sessions.Providers,
		observability.Providers,
		random.ProvidersRandom,
		postgres.Providers,
		grpc.ProvidersGRPC,
		qrcodes.Providers,
		// repos
		auditrepo.Providers,
		authrepo.Providers,
		dataprivacysrepo.Providers,
		identityrepo.Providers,
		notificationsrepo.Providers,
		settingsrepo.Providers,
		webhooksrepo.Providers,
		oauthrepo.Providers,
		mealplanningrepo.Providers,
		// services
		auditsvc.Providers,
		authsvc.Providers,
		dataprivacysvc.Providers,
		identitysvc.Providers,
		internalopssvc.Providers,
		notificationssvc.Providers,
		settingssvc.Providers,
		webhookssvc.Providers,
		oauthsvc.Providers,
		mealplanningsvc.Providers,
		// manager
		identitymgr.Providers,
		oauthmgr.Providers,
		mealplanningmgr.ProvidersManagers,
		authmgr.Providers,
		// workers
		mealplanfinalizer.ProvidersMealPlanFinalizer,
		mealplangrocerylistinitializer.ProvidersMealPlanGroceryListInitializer,
		mealplantaskcreator.ProvidersMealPlanTaskCreator,
		// misc
		recipeanalysis.ProvidersRecipeAnalysis,
		grocerylistpreparation.ProvidersGroceryListPreparation,
		ProvideUserTextSearcher,
		BuildUnaryServerInterceptors,
		BuildStreamServerInterceptors,
		BuildRegistrationFuncs,
		NewGRPCService,
	)

	return nil, nil
}
