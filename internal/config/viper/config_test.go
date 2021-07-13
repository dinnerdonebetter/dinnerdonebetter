package viper

import (
	"testing"
	"time"

	"gitlab.com/prixfixe/prixfixe/internal/capitalism"
	"gitlab.com/prixfixe/prixfixe/internal/config"
	"gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/server"
	"gitlab.com/prixfixe/prixfixe/internal/services/audit"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/services/invitations"
	"gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	"gitlab.com/prixfixe/prixfixe/internal/services/recipestepingredients"
	"gitlab.com/prixfixe/prixfixe/internal/services/recipestepproducts"
	"gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	"gitlab.com/prixfixe/prixfixe/internal/services/reports"
	"gitlab.com/prixfixe/prixfixe/internal/services/validingredientpreparations"
	"gitlab.com/prixfixe/prixfixe/internal/services/validingredients"
	"gitlab.com/prixfixe/prixfixe/internal/services/validinstruments"
	"gitlab.com/prixfixe/prixfixe/internal/services/validpreparationinstruments"
	"gitlab.com/prixfixe/prixfixe/internal/services/validpreparations"
	"gitlab.com/prixfixe/prixfixe/internal/storage"
	"gitlab.com/prixfixe/prixfixe/internal/uploads"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestBuildViperConfig(t *testing.T) {
	t.Parallel()

	actual := BuildViperConfig()
	assert.NotNil(t, actual)
}

func TestFromConfig(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleConfig := &config.InstanceConfig{
			Server: server.Config{
				HTTPPort:        1234,
				Debug:           false,
				StartupDeadline: time.Minute,
			},
			Meta: config.MetaSettings{
				RunMode: config.DevelopmentRunMode,
			},
			Encoding: encoding.Config{
				ContentType: "application/json",
			},
			Capitalism: capitalism.Config{
				Enabled:  false,
				Provider: capitalism.StripeProvider,
				Stripe: &capitalism.StripeConfig{
					APIKey:        "whatever",
					SuccessURL:    "whatever",
					CancelURL:     "whatever",
					WebhookSecret: "whatever",
				},
			},
			Observability: observability.Config{
				Metrics: metrics.Config{
					Provider:                         "",
					RouteToken:                       "",
					RuntimeMetricsCollectionInterval: 2 * time.Second,
				},
				Tracing: tracing.Config{
					Jaeger: &tracing.JaegerConfig{
						CollectorEndpoint: "things",
						ServiceName:       "stuff",
					},
					Provider:                  "blah",
					SpanCollectionProbability: 0,
				},
			},
			Uploads: uploads.Config{
				Storage: storage.Config{
					FilesystemConfig: &storage.FilesystemConfig{RootDirectory: "/blah"},
					AzureConfig: &storage.AzureConfig{
						BucketName: "blahs",
						Retrying:   &storage.AzureRetryConfig{},
					},
					GCSConfig: &storage.GCSConfig{
						ServiceAccountKeyFilepath: "/blah/blah",
						BucketName:                "blah",
						Scopes:                    nil,
					},
					S3Config:          &storage.S3Config{BucketName: "blahs"},
					BucketName:        "blahs",
					UploadFilenameKey: "blahs",
					Provider:          "blahs",
				},
				Debug: false,
			},
			Services: config.ServicesConfigurations{
				AuditLog: audit.Config{
					Enabled: true,
				},
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
				ValidInstruments: validinstruments.Config{
					SearchIndexPath: "/valid_instruments_index_path",
					Logging: logging.Config{
						Name:     "valid_instruments",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidPreparations: validpreparations.Config{
					SearchIndexPath: "/valid_preparations_index_path",
					Logging: logging.Config{
						Name:     "valid_preparations",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidIngredients: validingredients.Config{
					SearchIndexPath: "/valid_ingredients_index_path",
					Logging: logging.Config{
						Name:     "valid_ingredients",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidIngredientPreparations: validingredientpreparations.Config{
					Logging: logging.Config{
						Name:     "valid_ingredient_preparations",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidPreparationInstruments: validpreparationinstruments.Config{
					Logging: logging.Config{
						Name:     "valid_preparation_instruments",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				Recipes: recipes.Config{
					Logging: logging.Config{
						Name:     "recipes",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeSteps: recipesteps.Config{
					Logging: logging.Config{
						Name:     "recipe_steps",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeStepIngredients: recipestepingredients.Config{
					Logging: logging.Config{
						Name:     "recipe_step_ingredients",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeStepProducts: recipestepproducts.Config{
					Logging: logging.Config{
						Name:     "recipe_step_products",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				Invitations: invitations.Config{
					Logging: logging.Config{
						Name:     "invitations",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				Reports: reports.Config{
					Logging: logging.Config{
						Name:     "reports",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
			},
			Database: dbconfig.Config{
				Provider:                  "postgres",
				MetricsCollectionInterval: 2 * time.Second,
				Debug:                     true,
				RunMigrations:             true,
				ConnectionDetails:         database.ConnectionDetails("postgres://username:passwords@host/table"),
				CreateTestUser: &types.TestUserCreationConfig{
					Username:       "username",
					Password:       "password",
					HashedPassword: "hashashashashash",
					IsServiceAdmin: false,
				}},
		}

		actual, err := FromConfig(exampleConfig)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		actual, err := FromConfig(nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid config", func(t *testing.T) {
		t.Parallel()

		exampleConfig := &config.InstanceConfig{}

		actual, err := FromConfig(exampleConfig)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
