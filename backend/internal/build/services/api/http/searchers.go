package api

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"

	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
	textsearchcfg "github.com/verygoodsoftwarenotvirus/platform/search/text/config"
)

// ProvideTextSearchConfig provides a pointer to the text search config for dependency injection.
func ProvideTextSearchConfig(cfg *config.APIServiceConfig) *textsearchcfg.Config {
	return &cfg.TextSearch
}

// ProvideUserTextSearcher provides a user text searcher for the identity manager.
func ProvideUserTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (identityindexing.UserTextSearcher, error) {
	return textsearchcfg.ProvideIndex[identityindexing.UserSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		identityindexing.IndexTypeUsers,
	)
}
