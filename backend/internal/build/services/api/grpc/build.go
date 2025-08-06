//go:build wireinject
// +build wireinject

package grpcapi

import (
	"context"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/config"
	identitymgr "github.com/dinnerdonebetter/backend/internal/domain/identity/managers"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	auditrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	dataprivacysrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/dataprivacy"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	notificationsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/notifications"
	settingsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/settings"
	webhooksrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"
	auditsvc "github.com/dinnerdonebetter/backend/internal/services/audit/grpc"
	dataprivacysvc "github.com/dinnerdonebetter/backend/internal/services/dataprivacy/grpc"
	identitysvc "github.com/dinnerdonebetter/backend/internal/services/identity/grpc"
	internalopssvc "github.com/dinnerdonebetter/backend/internal/services/internalops/grpc"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/services/notifications/grpc"
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
		// repos
		auditrepo.Providers,
		dataprivacysrepo.Providers,
		identityrepo.Providers,
		notificationsrepo.Providers,
		settingsrepo.Providers,
		webhooksrepo.Providers,
		// services
		auditsvc.Providers,
		dataprivacysvc.Providers,
		identitysvc.Providers,
		internalopssvc.Providers,
		notificationssvc.Providers,
		settingssvc.Providers,
		webhookssvc.Providers,
		// managers
		identitymgr.Providers,
		// misc
		ProvideUserTextSearcher,
		// BuildUnaryServerInterceptors,
		// BuildStreamServerInterceptors,
		// BuildRegistrationFuncs,
		NewGRPCService,
	)

	return nil, nil
}
