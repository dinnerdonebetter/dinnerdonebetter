package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/prixfixeco/api_server/internal/authentication"
	"github.com/prixfixeco/api_server/internal/config"
	"github.com/prixfixeco/api_server/internal/database"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	"github.com/prixfixeco/api_server/internal/encoding"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/search"
	"github.com/prixfixeco/api_server/internal/secrets"
	"github.com/prixfixeco/api_server/internal/server"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	householdsservice "github.com/prixfixeco/api_server/internal/services/households"
	mealplanoptionsservice "github.com/prixfixeco/api_server/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/api_server/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/api_server/internal/services/mealplans"
	recipesservice "github.com/prixfixeco/api_server/internal/services/recipes"
	recipestepingredientsservice "github.com/prixfixeco/api_server/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/api_server/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/prixfixeco/api_server/internal/services/recipestepproducts"
	recipestepsservice "github.com/prixfixeco/api_server/internal/services/recipesteps"
	validingredientpreparationsservice "github.com/prixfixeco/api_server/internal/services/validingredientpreparations"
	validingredientsservice "github.com/prixfixeco/api_server/internal/services/validingredients"
	validinstrumentsservice "github.com/prixfixeco/api_server/internal/services/validinstruments"
	validpreparationsservice "github.com/prixfixeco/api_server/internal/services/validpreparations"
	webhooksservice "github.com/prixfixeco/api_server/internal/services/webhooks"
	websocketsservice "github.com/prixfixeco/api_server/internal/services/websockets"
	"github.com/prixfixeco/api_server/internal/storage"
	"github.com/prixfixeco/api_server/internal/uploads"
)

const (
	defaultPort              = 8888
	defaultCookieDomain      = ".prixfixe.local"
	debugCookieSecret        = "HEREISA32CHARSECRETWHICHISMADEUP"
	debugCookieSigningKey    = "DIFFERENT32CHARSECRETTHATIMADEUP"
	devPostgresDBConnDetails = "postgres://dbuser:hunter2@pgdatabase:5432/prixfixe?sslmode=disable"
	defaultCookieName        = authservice.DefaultCookieName

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	// database providers.
	postgres = "postgres"

	// test user stuff.
	defaultPassword = "password"

	localElasticsearchLocation = "http://elasticsearch:9200"

	// message provider topics
	preWritesTopicName   = "pre_writes"
	preUpdatesTopicName  = "pre_updates"
	preArchivesTopicName = "pre_archives"

	pasetoSecretSize      = 32
	maxAttempts           = 50
	defaultPASETOLifetime = 1 * time.Minute

	contentTypeJSON    = "application/json"
	workerQueueAddress = "worker_queue:6379"
)

var (
	examplePASETOKey = generatePASETOKey()

	noopTracingConfig = tracing.Config{
		SpanCollectionProbability: 0,
	}

	localServer = server.Config{
		Debug:           true,
		HTTPPort:        defaultPort,
		StartupDeadline: time.Minute,
		// HTTPSCertificateFile:    "/etc/certs/cert.pem",
		// HTTPSCertificateKeyFile: "/etc/certs/key.pem",
	}

	localCookies = authservice.CookieConfig{
		Name:       defaultCookieName,
		Domain:     defaultCookieDomain,
		HashKey:    debugCookieSecret,
		SigningKey: debugCookieSigningKey,
		Lifetime:   authservice.DefaultCookieLifetime,
		SecureOnly: true,
	}

	localTracingConfig = tracing.Config{
		Provider:                  "jaeger",
		SpanCollectionProbability: 1,
		Jaeger: &tracing.JaegerConfig{
			CollectorEndpoint: "http://localhost:14268/api/traces",
			ServiceName:       "prixfixe_service",
		},
	}
)

func initializeLocalSecretManager(ctx context.Context) secrets.SecretManager {
	logger := logging.NewNoopLogger()

	cfg := &secrets.Config{
		Provider: secrets.ProviderLocal,
		Key:      "SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU=",
	}

	k, err := secrets.ProvideSecretKeeper(ctx, cfg)
	if err != nil {
		panic(err)
	}

	sm, err := secrets.ProvideSecretManager(logger, k)
	if err != nil {
		panic(err)
	}

	return sm
}

func encryptAndSaveConfig(ctx context.Context, outputPath string, cfg *config.InstanceConfig) error {
	sm := initializeLocalSecretManager(ctx)
	output, err := sm.Encrypt(ctx, cfg)
	if err != nil {
		return fmt.Errorf("encrypting config: %v", err)
	}

	if err = os.MkdirAll(filepath.Dir(outputPath), 0777); err != nil {
		// that's okay
		_ = err
	}

	return os.WriteFile(outputPath, []byte(output), 0644)
}

type configFunc func(ctx context.Context, filePath string) error

var files = map[string]configFunc{
	"environments/local/service.config":                                   localDevelopmentConfig,
	"environments/testing/config_files/frontend-tests.config":             frontendTestsConfig,
	"environments/testing/config_files/integration-tests-postgres.config": buildIntegrationTestForDBImplementation(postgres, devPostgresDBConnDetails),
}

func mustHashPass(password string) string {
	hashed, err := authentication.ProvideArgon2Authenticator(logging.NewNoopLogger()).
		HashPassword(context.Background(), password)

	if err != nil {
		panic(err)
	}

	return hashed
}

func generatePASETOKey() []byte {
	b := make([]byte, pasetoSecretSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return b
}

func localDevelopmentConfig(ctx context.Context, filePath string) error {
	cfg := &config.InstanceConfig{
		Meta: config.MetaSettings{
			Debug:   true,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Events: msgconfig.Config{
			Provider: msgconfig.ProviderRedis,
			RedisConfig: msgconfig.RedisConfig{
				QueueAddress: workerQueueAddress,
			},
		},
		Server: localServer,
		Database: dbconfig.Config{
			Debug:             true,
			RunMigrations:     true,
			MaxPingAttempts:   maxAttempts,
			Provider:          postgres,
			ConnectionDetails: devPostgresDBConnDetails,
		},
		Observability: observability.Config{
			Metrics: metrics.Config{
				Provider:                         "prometheus",
				RouteToken:                       "",
				RuntimeMetricsCollectionInterval: time.Second,
			},
			Tracing: localTracingConfig,
		},
		Uploads: uploads.Config{
			Debug: true,
			Storage: storage.Config{
				UploadFilenameKey: "avatar",
				Provider:          "filesystem",
				BucketName:        "avatars",
				AzureConfig:       nil,
				GCSConfig:         nil,
				S3Config:          nil,
				FilesystemConfig: &storage.FilesystemConfig{
					RootDirectory: "/avatars",
				},
			},
		},
		Search: search.Config{
			Provider: search.ElasticsearchProvider,
			Address:  localElasticsearchLocation,
		},
		Services: config.ServicesConfigurations{
			Households: householdsservice.Config{
				PreWritesTopicName: preWritesTopicName,
			},
			Auth: authservice.Config{
				PASETO: authservice.PASETOConfig{
					Issuer:       "prixfixe_service",
					Lifetime:     defaultPASETOLifetime,
					LocalModeKey: examplePASETOKey,
				},
				Cookies:               localCookies,
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 4,
				MinimumPasswordLength: 8,
			},
			Webhooks: webhooksservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			Websockets: websocketsservice.Config{
				Logging: logging.Config{
					Name:     "webhook",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidInstruments: validinstrumentsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "valid_instruments",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidIngredients: validingredientsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "valid_ingredients",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidPreparations: validpreparationsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "valid_preparations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "valid_ingredient_preparations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			Recipes: recipesservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipes",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeSteps: recipestepsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipe_steps",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipe_step_instruments",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipe_step_ingredients",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipe_step_products",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			MealPlans: mealplansservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "meal_plans",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "meal_plan_options",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "meal_plan_option_votes",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
		},
	}

	return encryptAndSaveConfig(ctx, filePath, cfg)
}

func frontendTestsConfig(ctx context.Context, filePath string) error {
	cfg := &config.InstanceConfig{
		Meta: config.MetaSettings{
			Debug:   false,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Events: msgconfig.Config{
			Provider: msgconfig.ProviderRedis,
			RedisConfig: msgconfig.RedisConfig{
				QueueAddress: workerQueueAddress,
			},
		},
		Server: localServer,
		Database: dbconfig.Config{
			Debug:             true,
			RunMigrations:     true,
			Provider:          postgres,
			ConnectionDetails: devPostgresDBConnDetails,
			MaxPingAttempts:   maxAttempts,
		},
		Observability: observability.Config{
			Metrics: metrics.Config{
				Provider:                         "prometheus",
				RouteToken:                       "",
				RuntimeMetricsCollectionInterval: time.Second,
			},
			Tracing: noopTracingConfig,
		},
		Uploads: uploads.Config{
			Debug: true,
			Storage: storage.Config{
				UploadFilenameKey: "avatar",
				Provider:          "memory",
				BucketName:        "avatars",
			},
		},
		Search: search.Config{
			Provider: search.ElasticsearchProvider,
			Address:  localElasticsearchLocation,
		},
		Services: config.ServicesConfigurations{
			Households: householdsservice.Config{
				PreWritesTopicName: preWritesTopicName,
			},
			Auth: authservice.Config{
				PASETO: authservice.PASETOConfig{
					Issuer:       "prixfixe_service",
					Lifetime:     defaultPASETOLifetime,
					LocalModeKey: examplePASETOKey,
				},
				Cookies:               localCookies,
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 4,
				MinimumPasswordLength: 8,
			},
			Webhooks: webhooksservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			Websockets: websocketsservice.Config{
				Logging: logging.Config{
					Name:     "webhook",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidInstruments: validinstrumentsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "valid_instruments",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidIngredients: validingredientsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "valid_ingredients",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidPreparations: validpreparationsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "valid_preparations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "valid_ingredient_preparations",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			Recipes: recipesservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipes",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeSteps: recipestepsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipe_steps",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipe_step_instruments",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipe_step_ingredients",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipe_step_products",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			MealPlans: mealplansservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "meal_plans",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "meal_plan_options",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "meal_plan_option_votes",
					Level:    logging.InfoLevel,
					Provider: logging.ProviderZerolog,
				},
			},
		},
	}

	return encryptAndSaveConfig(ctx, filePath, cfg)
}

func buildIntegrationTestForDBImplementation(dbVendor, dbDetails string) configFunc {
	return func(ctx context.Context, filePath string) error {
		startupDeadline := time.Minute

		cfg := &config.InstanceConfig{
			Meta: config.MetaSettings{
				Debug:   false,
				RunMode: testingEnv,
			},
			Events: msgconfig.Config{
				Provider: msgconfig.ProviderRedis,
				RedisConfig: msgconfig.RedisConfig{
					QueueAddress: workerQueueAddress,
				},
			},
			Encoding: encoding.Config{
				ContentType: contentTypeJSON,
			},
			Server: server.Config{
				Debug:           false,
				HTTPPort:        defaultPort,
				StartupDeadline: startupDeadline,
			},
			Database: dbconfig.Config{
				Debug:             false,
				RunMigrations:     true,
				Provider:          dbVendor,
				MaxPingAttempts:   maxAttempts,
				ConnectionDetails: database.ConnectionDetails(dbDetails),
			},
			Observability: observability.Config{
				Metrics: metrics.Config{
					Provider:                         "",
					RouteToken:                       "",
					RuntimeMetricsCollectionInterval: time.Second,
				},
				Tracing: localTracingConfig,
			},
			Uploads: uploads.Config{
				Debug: false,
				Storage: storage.Config{
					Provider:    "memory",
					BucketName:  "avatars",
					AzureConfig: nil,
					GCSConfig:   nil,
					S3Config:    nil,
				},
			},
			Search: search.Config{
				Provider: search.ElasticsearchProvider,
				Address:  localElasticsearchLocation,
			},
			Services: config.ServicesConfigurations{
				Households: householdsservice.Config{
					PreWritesTopicName: preWritesTopicName,
				},
				Auth: authservice.Config{
					PASETO: authservice.PASETOConfig{
						Issuer:       "prixfixe_service",
						Lifetime:     defaultPASETOLifetime,
						LocalModeKey: examplePASETOKey,
					},
					Cookies: authservice.CookieConfig{
						Name:       defaultCookieName,
						Domain:     defaultCookieDomain,
						HashKey:    debugCookieSecret,
						SigningKey: debugCookieSigningKey,
						Lifetime:   authservice.DefaultCookieLifetime,
						SecureOnly: false,
					},
					Debug:                 false,
					EnableUserSignup:      true,
					MinimumUsernameLength: 4,
					MinimumPasswordLength: 8,
				},
				Webhooks: webhooksservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
				},
				Websockets: websocketsservice.Config{
					Logging: logging.Config{
						Name:     "webhook",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidInstruments: validinstrumentsservice.Config{
					SearchIndexPath:      localElasticsearchLocation,
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "valid_instruments",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidIngredients: validingredientsservice.Config{
					SearchIndexPath:      localElasticsearchLocation,
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "valid_ingredients",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidPreparations: validpreparationsservice.Config{
					SearchIndexPath:      localElasticsearchLocation,
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "valid_preparations",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				ValidIngredientPreparations: validingredientpreparationsservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "valid_ingredient_preparations",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				Recipes: recipesservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "recipes",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeSteps: recipestepsservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "recipe_steps",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeStepInstruments: recipestepinstrumentsservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "recipe_step_instruments",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeStepIngredients: recipestepingredientsservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "recipe_step_ingredients",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				RecipeStepProducts: recipestepproductsservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "recipe_step_products",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				MealPlans: mealplansservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "meal_plans",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				MealPlanOptions: mealplanoptionsservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "meal_plan_options",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
				MealPlanOptionVotes: mealplanoptionvotesservice.Config{
					PreWritesTopicName:   preWritesTopicName,
					PreUpdatesTopicName:  preUpdatesTopicName,
					PreArchivesTopicName: preArchivesTopicName,
					Logging: logging.Config{
						Name:     "meal_plan_option_votes",
						Level:    logging.InfoLevel,
						Provider: logging.ProviderZerolog,
					},
				},
			},
		}

		return encryptAndSaveConfig(ctx, filePath, cfg)
	}
}

func main() {
	ctx := context.Background()

	for filePath, fun := range files {
		if err := fun(ctx, filePath); err != nil {
			log.Fatalf("error rendering %s: %v", filePath, err)
		}
	}
}
