package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	postgres "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding/postgres"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/Masterminds/squirrel"
	postgresstore "github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// PostgresProvider is the string used to refer to postgres.
	PostgresProvider = "postgres"

	// DefaultMetricsCollectionInterval is the default amount of time we wait between database metrics queries.
	DefaultMetricsCollectionInterval = 2 * time.Second
)

var (
	errInvalidDatabase = errors.New("invalid database")
	errNilDBProvided   = errors.New("invalid DB connection provided")
)

type (
	// Config represents our database configuration.
	Config struct {
		CreateTestUser            *types.TestUserCreationConfig `json:"createTestUser" mapstructure:"create_test_user" toml:"create_test_user,omitempty"`
		Provider                  string                        `json:"provider" mapstructure:"provider" toml:"provider,omitempty"`
		ConnectionDetails         database.ConnectionDetails    `json:"connectionDetails" mapstructure:"connection_details" toml:"connection_details,omitempty"`
		MetricsCollectionInterval time.Duration                 `json:"metricsCollectionInterval" mapstructure:"metrics_collection_interval" toml:"metrics_collection_interval,omitempty"`
		Debug                     bool                          `json:"debug" mapstructure:"debug" toml:"debug,omitempty"`
		RunMigrations             bool                          `json:"runMigrations" mapstructure:"run_migrations" toml:"run_migrations,omitempty"`
		MaxPingAttempts           uint8                         `json:"maxPingAttempts" mapstructure:"max_ping_attempts" toml:"max_ping_attempts,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates an DatabaseSettings struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.ConnectionDetails, validation.Required),
		validation.Field(&cfg.Provider, validation.In(PostgresProvider)),
		validation.Field(&cfg.CreateTestUser, validation.When(cfg.CreateTestUser != nil, validation.Required).Else(validation.Nil)),
	)
}

// ProvideDatabaseConnection provides a database implementation dependent on the configuration.
func ProvideDatabaseConnection(logger logging.Logger, cfg *Config) (*sql.DB, error) {
	switch cfg.Provider {
	case PostgresProvider:
		return postgres.ProvidePostgresDB(logger, cfg.ConnectionDetails)
	default:
		return nil, fmt.Errorf("%w: %q", errInvalidDatabase, cfg.Provider)
	}
}

// ProvideDatabasePlaceholderFormat provides .
func (cfg *Config) ProvideDatabasePlaceholderFormat() (squirrel.PlaceholderFormat, error) {
	switch cfg.Provider {
	case PostgresProvider:
		return squirrel.Dollar, nil
	default:
		return nil, fmt.Errorf("%w: %q", errInvalidDatabase, cfg.Provider)
	}
}

// ProvideJSONPluckQuery provides a query for extracting a value out of a JSON dictionary for a given database.
func (cfg *Config) ProvideJSONPluckQuery() string {
	switch cfg.Provider {
	case PostgresProvider:
		return `%s.%s->'%s'`
	default:
		return ""
	}
}

// ProvideCurrentUnixTimestampQuery provides a database implementation dependent on the configuration.
func (cfg *Config) ProvideCurrentUnixTimestampQuery() string {
	switch cfg.Provider {
	case PostgresProvider:
		return `extract(epoch FROM NOW())`
	default:
		return ""
	}
}

// ProvideSessionManager provides a session manager based on some settings.
// There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,
// but it obviously doesn't belong in the database package, or maybe it does.
func ProvideSessionManager(cookieConfig authservice.CookieConfig, dbConf Config, db *sql.DB) (*scs.SessionManager, error) {
	sessionManager := scs.New()

	if db == nil {
		return nil, errNilDBProvided
	}

	switch dbConf.Provider {
	case PostgresProvider:
		sessionManager.Store = postgresstore.New(db)
	default:
		return nil, fmt.Errorf("%w: %q", errInvalidDatabase, dbConf.Provider)
	}

	sessionManager.Lifetime = cookieConfig.Lifetime
	sessionManager.Cookie.Name = cookieConfig.Name
	sessionManager.Cookie.Domain = cookieConfig.Domain
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = cookieConfig.SecureOnly

	return sessionManager, nil
}
