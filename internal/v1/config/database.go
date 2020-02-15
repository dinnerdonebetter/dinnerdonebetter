package config

import (
	"context"
	"errors"
	"fmt"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	dbclient "gitlab.com/prixfixe/prixfixe/database/v1/client"
	"gitlab.com/prixfixe/prixfixe/database/v1/queriers/postgres"

	"contrib.go.opencensus.io/integrations/ocsql"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	postgresProviderKey = "postgres"
)

// ProvideDatabase provides a database implementation dependent on the configuration
func (cfg *ServerConfig) ProvideDatabase(ctx context.Context, logger logging.Logger) (database.Database, error) {
	var (
		debug             = cfg.Database.Debug || cfg.Meta.Debug
		connectionDetails = cfg.Database.ConnectionDetails
	)

	switch cfg.Database.Provider {
	case postgresProviderKey:
		rawDB, err := postgres.ProvidePostgresDB(logger, connectionDetails)
		if err != nil {
			return nil, fmt.Errorf("establish postgres database connection: %w", err)
		}
		ocsql.RegisterAllViews()
		ocsql.RecordStats(rawDB, cfg.Metrics.DBMetricsCollectionInterval)

		pgdb := postgres.ProvidePostgres(debug, rawDB, logger)

		return dbclient.ProvideDatabaseClient(ctx, rawDB, pgdb, debug, logger)
	default:
		return nil, errors.New("invalid database type selected")
	}
}
