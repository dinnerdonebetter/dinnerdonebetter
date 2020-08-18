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
	scs "github.com/alexedwards/scs/v2"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	// PostgresProviderKey is the string we use to refer to postgres
	PostgresProviderKey = "postgres"
)

// ProvideDatabaseConnection provides a database implementation dependent on the configuration.
func (cfg *ServerConfig) ProvideDatabaseConnection(logger logging.Logger) (*sql.DB, error) {
	switch cfg.Database.Provider {
	case PostgresProviderKey:
		return postgres.ProvidePostgresDB(logger, cfg.Database.ConnectionDetails)
	default:
		return nil, fmt.Errorf("invalid database type selected: %q", cfg.Database.Provider)
	}
}

// ProvideDatabaseClient provides a database implementation dependent on the configuration.
func (cfg *ServerConfig) ProvideDatabaseClient(ctx context.Context, logger logging.Logger, rawDB *sql.DB) (database.DataManager, error) {
	if rawDB == nil {
		return nil, errors.New("nil DB connection provided")
	}

	debug := cfg.Database.Debug || cfg.Meta.Debug

	ocsql.RegisterAllViews()
	ocsql.RecordStats(rawDB, cfg.Metrics.DBMetricsCollectionInterval)

	var dbc database.DataManager
	switch cfg.Database.Provider {
	case PostgresProviderKey:
		dbc = postgres.ProvidePostgres(debug, rawDB, logger)
	default:
		return nil, fmt.Errorf("invalid database type selected: %q", cfg.Database.Provider)
	}

	return dbclient.ProvideDatabaseClient(ctx, rawDB, dbc, debug, logger)
}

// ProvideSessionManager provides a session manager based on some settings.
// There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,
// but it obviously doesn't belong in the database package, or maybe it does
func ProvideSessionManager(authConf AuthSettings, dbConf DatabaseSettings, db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()

	switch dbConf.Provider {
	case PostgresProviderKey:
		sessionManager.Store = postgresstore.New(db)
	}

	sessionManager.Lifetime = authConf.CookieLifetime
	// elaborate further here later if you so choose

	return sessionManager
}
