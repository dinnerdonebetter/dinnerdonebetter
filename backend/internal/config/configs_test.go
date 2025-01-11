package config

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config/envvars"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerConfig_EncodeToFile(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &APIServiceConfig{
			Server: http.Config{
				HTTPPort:        1234,
				Debug:           false,
				StartupDeadline: time.Minute,
			},
			Meta: MetaSettings{
				RunMode: DevelopmentRunMode,
			},
			Encoding: encoding.Config{
				ContentType: "application/json",
			},
			Observability: observability.Config{},
			Services: ServicesConfig{
				Auth: authservice.Config{
					MinimumUsernameLength: 4,
					MinimumPasswordLength: 8,
					EnableUserSignup:      true,
				},
			},
			Database: databasecfg.Config{
				Debug:         true,
				RunMigrations: true,
				ConnectionDetails: databasecfg.ConnectionDetails{
					Username:   "username",
					Password:   "password",
					Database:   "table",
					Host:       "host",
					DisableSSL: true,
				},
			},
		}

		f, err := os.CreateTemp("", "")
		require.NoError(t, err)

		assert.NoError(t, cfg.EncodeToFile(f.Name(), json.Marshal))
	})

	T.Run("with error marshaling", func(t *testing.T) {
		t.Parallel()

		cfg := &APIServiceConfig{}

		f, err := os.CreateTemp("", "")
		require.NoError(t, err)

		assert.Error(t, cfg.EncodeToFile(f.Name(), func(any) ([]byte, error) {
			return nil, errors.New("blah")
		}))
	})
}

//nolint:paralleltest // because we set env vars for this, we can't
func TestFetchForApplication(T *testing.T) {
	T.Run("standard", func(t *testing.T) {
		ctx := context.Background()

		cfg := &APIServiceConfig{
			Database: databasecfg.Config{
				Debug: true,
			},
		}
		cfgBytes, err := json.Marshal(cfg)
		require.NoError(t, err)

		configFilepath := t.TempDir() + "/config.json"
		require.NoError(t, os.WriteFile(configFilepath, cfgBytes, 0o0644))

		t.Setenv(ConfigurationFilePathEnvVarKey, configFilepath)

		actual, err := FetchForApplication(ctx, GetAPIServiceConfigFromGoogleCloudRunEnvironment)
		assert.NoError(t, err)
		assert.NotNil(t, actual)

		assert.Equal(t, actual.Database.Debug, true)
	})

	// prior TODOs count here too
	T.Run("overrides meta", func(t *testing.T) {
		ctx := context.Background()

		cfg := &APIServiceConfig{
			Database: databasecfg.Config{
				Debug: true,
			},
		}
		cfgBytes, err := json.Marshal(cfg)
		require.NoError(t, err)

		configFilepath := t.TempDir() + "/config.json"
		require.NoError(t, os.WriteFile(configFilepath, cfgBytes, 0o0644))

		t.Setenv(ConfigurationFilePathEnvVarKey, configFilepath)
		t.Setenv(envvars.MetaDebugEnvVarKey, "false")

		actual, err := FetchForApplication(ctx, GetAPIServiceConfigFromGoogleCloudRunEnvironment)
		assert.NoError(t, err)
		assert.NotNil(t, actual)

		assert.Equal(t, actual.Meta.Debug, false)
	})
}
