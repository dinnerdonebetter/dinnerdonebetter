package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	dbclient "gitlab.com/prixfixe/prixfixe/database/v1/client"
	"gitlab.com/prixfixe/prixfixe/database/v1/queriers/postgres"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	postgresProviderKey = "postgres"
)

// ProvideDatabase provides a database implementation dependent on the configuration.
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

		return dbclient.ProvideDatabaseClient(ctx, logger, rawDB, pgdb, debug, cfg.Database.CreateDummyUser)
	default:
		return nil, errors.New("invalid database type selected")
	}
}

// ProvideDatabaseConnection provides a database implementation dependent on the configuration.
func ProvideDatabaseConnection(logger logging.Logger, dbSettings DatabaseSettings) (*sql.DB, error) {
	return postgres.ProvidePostgresDB(logger, dbSettings.ConnectionDetails)
}

// ProvideSessionManager provides a session manager based on some settings.
// There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,
// but it obviously doesn't belong in the database package, or maybe it does
func ProvideSessionManager(authConf AuthSettings, db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()

	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = authConf.CookieLifetime

	return sessionManager
}
