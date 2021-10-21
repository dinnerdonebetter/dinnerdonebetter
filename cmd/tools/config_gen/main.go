package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/prixfixe/prixfixe/internal/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/config"
	database "gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	msgconfig "gitlab.com/prixfixe/prixfixe/internal/messagequeue/config"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	"gitlab.com/prixfixe/prixfixe/internal/secrets"
	server "gitlab.com/prixfixe/prixfixe/internal/server"
	accountsservice "gitlab.com/prixfixe/prixfixe/internal/services/accounts"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	frontendservice "gitlab.com/prixfixe/prixfixe/internal/services/frontend"
	mealplanoptionsservice "gitlab.com/prixfixe/prixfixe/internal/services/mealplanoptions"
	mealplanoptionvotesservice "gitlab.com/prixfixe/prixfixe/internal/services/mealplanoptionvotes"
	mealplansservice "gitlab.com/prixfixe/prixfixe/internal/services/mealplans"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepingredients"
	recipestepinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepinstruments"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/internal/services/webhooks"
	websocketsservice "gitlab.com/prixfixe/prixfixe/internal/services/websockets"
	storage "gitlab.com/prixfixe/prixfixe/internal/storage"
	uploads "gitlab.com/prixfixe/prixfixe/internal/uploads"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	defaultPort              = 8888
	defaultCookieDomain      = "localhost"
	debugCookieSecret        = "HEREISA32CHARSECRETWHICHISMADEUP"
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
	}

	localCookies = authservice.CookieConfig{
		Name:       defaultCookieName,
		Domain:     defaultCookieDomain,
		HashKey:    debugCookieSecret,
		SigningKey: debugCookieSecret,
		Lifetime:   authservice.DefaultCookieLifetime,
		SecureOnly: false,
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
	}

	return os.WriteFile(outputPath, []byte(output), 0644)
}

type configFunc func(ctx context.Context, filePath string) error

var files = map[string]configFunc{
	"environments/local/service.config":                                   localDevelopmentConfig,
	"environments/testing/config_files/frontend-tests.config":             frontendTestsConfig,
	"environments/testing/config_files/integration-tests-postgres.config": buildIntegrationTestForDBImplementation(postgres, devPostgresDBConnDetails),
}

func buildLocalFrontendServiceConfig() frontendservice.Config {
	return frontendservice.Config{}
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
			CreateTestUser: &types.TestUserCreationConfig{
				Username:       "username",
				Password:       defaultPassword,
				HashedPassword: mustHashPass(defaultPassword),
				IsServiceAdmin: true,
			},
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
			Accounts: accountsservice.Config{
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
			Frontend: buildLocalFrontendServiceConfig(),
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
			Accounts: accountsservice.Config{
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
			Frontend: buildLocalFrontendServiceConfig(),
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
				CreateTestUser: &types.TestUserCreationConfig{
					Username:       "exampleUser",
					Password:       "integration-tests-are-cool",
					HashedPassword: mustHashPass("integration-tests-are-cool"),
					IsServiceAdmin: true,
				},
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
				Accounts: accountsservice.Config{
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
						SigningKey: debugCookieSecret,
						Lifetime:   authservice.DefaultCookieLifetime,
						SecureOnly: false,
					},
					Debug:                 false,
					EnableUserSignup:      true,
					MinimumUsernameLength: 4,
					MinimumPasswordLength: 8,
				},
				Frontend: buildLocalFrontendServiceConfig(),
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
