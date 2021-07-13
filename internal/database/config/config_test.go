package config

import (
	"context"
	"database/sql"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"

	"github.com/stretchr/testify/assert"
)

const (
	invalidProvider = "blah"
)

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			Provider:          PostgresProvider,
			ConnectionDetails: "example_connection_string",
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}

func TestConfig_ProvideDatabaseConnection(T *testing.T) {
	T.Parallel()

	T.Run("standard for postgres", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := &Config{
			Provider:          PostgresProvider,
			ConnectionDetails: "example_connection_string",
		}

		db, err := ProvideDatabaseConnection(logger, cfg)
		assert.NotNil(t, db)
		assert.NoError(t, err)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := &Config{
			Provider:          invalidProvider,
			ConnectionDetails: "example_connection_string",
		}

		db, err := ProvideDatabaseConnection(logger, cfg)
		assert.Nil(t, db)
		assert.Error(t, err)
	})
}

func TestConfig_ProvideDatabasePlaceholderFormat(T *testing.T) {
	T.Parallel()

	T.Run("standard for postgres", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider:          PostgresProvider,
			ConnectionDetails: "example_connection_string",
		}

		pf, err := cfg.ProvideDatabasePlaceholderFormat()
		assert.NotNil(t, pf)
		assert.NoError(t, err)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider:          invalidProvider,
			ConnectionDetails: "example_connection_string",
		}

		pf, err := cfg.ProvideDatabasePlaceholderFormat()
		assert.Nil(t, pf)
		assert.Error(t, err)
	})
}

func TestConfig_ProvideJSONPluckQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard for postgres", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider:          PostgresProvider,
			ConnectionDetails: "example_connection_string",
		}

		assert.NotEmpty(t, cfg.ProvideJSONPluckQuery())
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider:          invalidProvider,
			ConnectionDetails: "example_connection_string",
		}

		assert.Empty(t, cfg.ProvideJSONPluckQuery())
	})
}

func TestConfig_ProvideCurrentUnixTimestampQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard for postgres", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider:          PostgresProvider,
			ConnectionDetails: "example_connection_string",
		}

		assert.NotEmpty(t, cfg.ProvideCurrentUnixTimestampQuery())
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider:          invalidProvider,
			ConnectionDetails: "example_connection_string",
		}

		assert.Empty(t, cfg.ProvideCurrentUnixTimestampQuery())
	})
}

func TestProvideSessionManager(T *testing.T) {
	T.Parallel()

	T.Run("with nil database", func(t *testing.T) {
		t.Parallel()

		cookieConfig := authservice.CookieConfig{}
		cfg := Config{
			Provider:          PostgresProvider,
			ConnectionDetails: "example_connection_string",
		}

		sessionManager, err := ProvideSessionManager(cookieConfig, cfg, nil)
		assert.Nil(t, sessionManager)
		assert.Error(t, err)
	})

	T.Run("standard for postgres", func(t *testing.T) {
		t.Parallel()

		cookieConfig := authservice.CookieConfig{}
		cfg := Config{
			Provider:          PostgresProvider,
			ConnectionDetails: "example_connection_string",
		}

		sessionManager, err := ProvideSessionManager(cookieConfig, cfg, &sql.DB{})
		assert.NotNil(t, sessionManager)
		assert.NoError(t, err)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		cookieConfig := authservice.CookieConfig{}
		cfg := Config{
			Provider:          invalidProvider,
			ConnectionDetails: "example_connection_string",
		}

		sessionManager, err := ProvideSessionManager(cookieConfig, cfg, &sql.DB{})
		assert.Nil(t, sessionManager)
		assert.Error(t, err)
	})
}
