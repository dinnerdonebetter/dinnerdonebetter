package config

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	database "github.com/prixfixeco/api_server/internal/database"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	"github.com/prixfixeco/api_server/internal/encoding"
	observability "github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	server "github.com/prixfixeco/api_server/internal/server"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	validingredientsservice "github.com/prixfixeco/api_server/internal/services/validingredients"
	validinstrumentsservice "github.com/prixfixeco/api_server/internal/services/validinstruments"
	validpreparationsservice "github.com/prixfixeco/api_server/internal/services/validpreparations"
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
				Auth: authservice.Config{
					Cookies: authservice.CookieConfig{
						Name:     "prixfixecookie",
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
				ValidIngredients: validingredientsservice.Config{
					SearchIndexPath: "/valid_ingredients_index_path",
				},
				ValidPreparations: validpreparationsservice.Config{
					SearchIndexPath: "/valid_preparations_index_path",
				},
			},
			Database: dbconfig.Config{
				Provider:          "postgres",
				Debug:             true,
				RunMigrations:     true,
				ConnectionDetails: database.ConnectionDetails("postgres://username:passwords@host/table"),
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

			x, err := ProvideDatabaseClient(ctx, logger, cfg)
			assert.NotNil(t, x)
			assert.NoError(t, err)
		}
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

		x, err := ProvideDatabaseClient(ctx, logger, cfg)
		assert.Nil(t, x)
		assert.Error(t, err)
	})
}
