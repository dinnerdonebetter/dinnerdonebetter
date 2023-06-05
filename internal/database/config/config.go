package config

import (
	"context"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"

	"github.com/alexedwards/scs/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config represents our database configuration.
	Config struct {
		_ struct{}

		ConnectionDetails database.ConnectionDetails `json:"connectionDetails" mapstructure:"connection_details" toml:"connection_details,omitempty"`
		Debug             bool                       `json:"debug"             mapstructure:"debug"              toml:"debug,omitempty"`
		LogQueries        bool                       `json:"logQueries"        mapstructure:"log_queries"        toml:"log_queries,omitempty"`
		RunMigrations     bool                       `json:"runMigrations"     mapstructure:"run_migrations"     toml:"run_migrations,omitempty"`
		MaxPingAttempts   uint8                      `json:"maxPingAttempts"   mapstructure:"max_ping_attempts"  toml:"max_ping_attempts,omitempty"`
		PingWaitPeriod    time.Duration              `json:"pingWaitPeriod"    mapstructure:"ping_wait_period"   toml:"ping_wait_period,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates an DatabaseSettings struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.ConnectionDetails, validation.Required),
	)
}

// ProvideSessionManager provides a session manager based on some settings.
// There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,
// but it obviously doesn't belong in the database package, or maybe it does.
func ProvideSessionManager(cookieConfig authservice.CookieConfig, dm database.DataManager) (*scs.SessionManager, error) {
	sessionManager := scs.New()

	sessionManager.Lifetime = cookieConfig.Lifetime
	sessionManager.Cookie.Name = cookieConfig.Name
	sessionManager.Cookie.Domain = cookieConfig.Domain
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = cookieConfig.SecureOnly

	sessionManager.Store = dm.ProvideSessionStore()

	return sessionManager, nil
}
