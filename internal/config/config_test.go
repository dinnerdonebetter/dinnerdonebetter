package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	server "gitlab.com/prixfixe/prixfixe/internal/server"
	auditservice "gitlab.com/prixfixe/prixfixe/internal/services/audit"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerConfig_EncodeToFile(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &InstanceConfig{
			Server: server.Config{
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
			Observability: observability.Config{
				Metrics: metrics.Config{
					Provider:                         "",
					RouteToken:                       "",
					RuntimeMetricsCollectionInterval: 2 * time.Second,
				},
			},
			Services: ServicesConfigurations{
				AuditLog: auditservice.Config{
					Enabled: true,
				},
				Auth: authservice.Config{
					Cookies: authservice.CookieConfig{
						Name:     "prixfixe_cookie",
						Domain:   "https://verygoodsoftwarenotvirus.ru",
						Lifetime: time.Second,
					},
					MinimumUsernameLength: 4,
					MinimumPasswordLength: 8,
					EnableUserSignup:      true,
				},
				ValidInstruments: validinstrumentsservice.Config{
					SearchIndexPath: "/valid_instruments_index_path",
				},
				ValidPreparations: validpreparationsservice.Config{
					SearchIndexPath: "/valid_preparations_index_path",
				},
				ValidIngredients: validingredientsservice.Config{
					SearchIndexPath: "/valid_ingredients_index_path",
				},
			},
			Database: dbconfig.Config{
				Provider:                  "postgres",
				MetricsCollectionInterval: 2 * time.Second,
				Debug:                     true,
				RunMigrations:             true,
				ConnectionDetails:         database.ConnectionDetails("postgres://username:passwords@host/table"),
			},
		}

		f, err := ioutil.TempFile("", "")
		require.NoError(t, err)

		assert.NoError(t, cfg.EncodeToFile(f.Name(), json.Marshal))
	})

	T.Run("with error marshaling", func(t *testing.T) {
		t.Parallel()

		cfg := &InstanceConfig{}

		f, err := ioutil.TempFile("", "")
		require.NoError(t, err)

		assert.Error(t, cfg.EncodeToFile(f.Name(), func(interface{}) ([]byte, error) {
			return nil, errors.New("blah")
		}))
	})
}

func TestServerConfig_ProvideDatabaseClient(T *testing.T) {
	T.Parallel()

	T.Run("supported providers", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		for _, provider := range []string{"postgres"} {
			cfg := &InstanceConfig{
				Database: dbconfig.Config{
					Provider: provider,
				},
			}

			x, err := ProvideDatabaseClient(ctx, logger, &sql.DB{}, cfg)
			assert.NotNil(t, x)
			assert.NoError(t, err)
		}
	})

	T.Run("with nil *sql.DB", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		cfg := &InstanceConfig{}

		x, err := ProvideDatabaseClient(ctx, logger, nil, cfg)
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := &InstanceConfig{
			Database: dbconfig.Config{
				Provider: "provider",
			},
		}

		x, err := ProvideDatabaseClient(ctx, logger, &sql.DB{}, cfg)
		assert.Nil(t, x)
		assert.Error(t, err)
	})
}
