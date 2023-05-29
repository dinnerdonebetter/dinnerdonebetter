//go:build wireinject
// +build wireinject

package build

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/observability"
	logcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/server/rpc"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.InstanceConfig,
) (*rpc.Server, error) {
	wire.Build(
		config.ServiceConfigProviders,
		database.DBProviders,
		tracingcfg.Providers,
		observability.Providers,
		logcfg.Providers,
		postgres.Providers,
		rpc.Providers,
	)

	return nil, nil
}
