package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/database"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	householdinvitationsservice "github.com/prixfixeco/api_server/internal/services/householdinvitations"
	mealsservice "github.com/prixfixeco/api_server/internal/services/meals"

	"github.com/prixfixeco/api_server/internal/config"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	"github.com/prixfixeco/api_server/internal/encoding"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/search"
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

	localEmailConfig = emailconfig.Config{
		Provider: "",
		APIToken: "",
	}

	localCustomerDataPlatformConfig = customerdataconfig.Config{
		Provider: "",
		APIToken: "",
	}
)

func saveConfig(ctx context.Context, outputPath string, cfg *config.InstanceConfig, indent, validate bool) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0777); err != nil {
		// that's okay
		_ = err
	}

	if validate {
		if err := cfg.ValidateWithContext(ctx); err != nil {
			return err
		}
	}

	var (
		output []byte
		err    error
	)

	if indent {
		output, err = json.MarshalIndent(cfg, "", "\t")
	} else {
		output, err = json.Marshal(cfg)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, output, 0644)
}

type configFunc func(ctx context.Context, filePath string) error

var files = map[string]configFunc{
	"environments/dev/config_files/service-config.json":                     devEnvironmentConfig,
	"environments/local/config_files/service-config.json":                   localDevelopmentConfig,
	"environments/testing/config_files/integration-tests-config.json":       integrationTestConfig,
	"environments/testing/config_files/integration-tests-local-config.json": localIntegrationTestConfig,
}

func generatePASETOKey() []byte {
	b := make([]byte, pasetoSecretSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return b
}

func loggingConfigWithName(name string) logging.Config {
	return logging.Config{
		Name:     name,
		Level:    logging.InfoLevel,
		Provider: logging.ProviderZerolog,
	}
}

func devEnvironmentConfig(ctx context.Context, filePath string) error {
	cookieConfig := authservice.CookieConfig{
		Name:       defaultCookieName,
		Domain:     ".prixfixe.dev",
		Lifetime:   authservice.DefaultCookieLifetime,
		SecureOnly: true,
	}

	emailConfig := emailconfig.Config{
		Provider: emailconfig.ProviderSendgrid,
		APIToken: "",
	}

	customerDataPlatformConfig := customerdataconfig.Config{
		Provider: customerdataconfig.ProviderSegment,
		APIToken: "",
	}

	cfg := &config.InstanceConfig{
		Meta: config.MetaSettings{
			Debug:   true,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Events: msgconfig.Config{
			Provider: msgconfig.ProviderSQS,
		},
		Email:        emailConfig,
		CustomerData: customerDataPlatformConfig,
		Server:       localServer,
		Database: dbconfig.Config{
			Debug:           true,
			RunMigrations:   true,
			MaxPingAttempts: maxAttempts,
		},
		Observability: observability.Config{
			Metrics: metrics.Config{
				Provider:                         "prometheus",
				RuntimeMetricsCollectionInterval: time.Second,
			},
			Tracing: noopTracingConfig,
		},
		Uploads: uploads.Config{
			Debug: true,
			Storage: storage.Config{
				UploadFilenameKey: "avatar",
				Provider:          "filesystem",
				BucketName:        "avatars.prixfixe.dev",
				S3Config:          nil,
				FilesystemConfig: &storage.FilesystemConfig{
					RootDirectory: "/avatars",
				},
			},
		},
		Search: search.Config{
			Provider: search.ElasticsearchProvider,
		},
		Services: config.ServicesConfigurations{
			Households: householdsservice.Config{
				Logging: loggingConfigWithName("households"),
			},
			HouseholdInvitations: householdinvitationsservice.Config{
				Logging: loggingConfigWithName("household_invitations"),
			},
			Auth: authservice.Config{
				Logging: loggingConfigWithName("authentication"),
				PASETO: authservice.PASETOConfig{
					Issuer:   "prixfixe_service",
					Lifetime: defaultPASETOLifetime,
				},
				Cookies:               cookieConfig,
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 4,
				MinimumPasswordLength: 8,
			},
			Webhooks: webhooksservice.Config{
				Logging: loggingConfigWithName("webhooks"),
			},
			Websockets: websocketsservice.Config{
				Logging: loggingConfigWithName("websocket"),
			},
			ValidInstruments: validinstrumentsservice.Config{
				Logging: loggingConfigWithName("valid_instruments"),
			},
			ValidIngredients: validingredientsservice.Config{
				Logging: loggingConfigWithName("valid_ingredients"),
			},
			ValidPreparations: validpreparationsservice.Config{
				Logging: loggingConfigWithName("valid_preparations"),
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				Logging: loggingConfigWithName("valid_ingredient_preparations"),
			},
			Meals: mealsservice.Config{
				Logging: loggingConfigWithName("meals"),
			},
			Recipes: recipesservice.Config{
				Logging: loggingConfigWithName("recipes"),
			},
			RecipeSteps: recipestepsservice.Config{
				Logging: loggingConfigWithName("recipe_steps"),
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				Logging: loggingConfigWithName("recipe_step_instruments"),
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				Logging: loggingConfigWithName("recipe_step_ingredients"),
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				Logging: loggingConfigWithName("recipe_step_products"),
			},
			MealPlans: mealplansservice.Config{
				Logging: loggingConfigWithName("meal_plans"),
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				Logging: loggingConfigWithName("meal_plan_options"),
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				Logging: loggingConfigWithName("meal_plan_option_votes"),
			},
		},
	}

	return saveConfig(ctx, filePath, cfg, true, false)
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
		Email:        localEmailConfig,
		CustomerData: localCustomerDataPlatformConfig,
		Server:       localServer,
		Database: dbconfig.Config{
			Debug:             true,
			RunMigrations:     true,
			MaxPingAttempts:   maxAttempts,
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
				BucketName:        "avatars.prixfixe.dev",
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
			HouseholdInvitations: householdinvitationsservice.Config{
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
			Meals: mealsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipes",
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

	return saveConfig(ctx, filePath, cfg, true, true)
}

func buildIntegrationTestsConfig() *config.InstanceConfig {
	return &config.InstanceConfig{
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
		Email:        localEmailConfig,
		CustomerData: localCustomerDataPlatformConfig,
		Server: server.Config{
			Debug:           false,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: dbconfig.Config{
			Debug:             false,
			RunMigrations:     true,
			MaxPingAttempts:   maxAttempts,
			ConnectionDetails: devPostgresDBConnDetails,
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
			HouseholdInvitations: householdinvitationsservice.Config{
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
			Meals: mealsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
				Logging: logging.Config{
					Name:     "recipes",
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
}

func integrationTestConfig(ctx context.Context, filePath string) error {
	cfg := buildIntegrationTestsConfig()

	return saveConfig(ctx, filePath, cfg, true, true)
}

func localIntegrationTestConfig(ctx context.Context, filePath string) error {
	cfg := buildIntegrationTestsConfig()

	cfg.Events.RedisConfig.QueueAddress = msgconfig.MessageQueueAddress(strings.ReplaceAll(workerQueueAddress, "worker_queue", "localhost"))
	cfg.Search.Address = search.IndexPath(strings.ReplaceAll(localElasticsearchLocation, "elasticsearch", "localhost"))
	cfg.Database.ConnectionDetails = database.ConnectionDetails(strings.ReplaceAll(devPostgresDBConnDetails, "pgdatabase", "localhost"))

	cfg.Services.ValidInstruments.SearchIndexPath = strings.ReplaceAll(localElasticsearchLocation, "elasticsearch", "localhost")
	cfg.Services.ValidIngredients.SearchIndexPath = strings.ReplaceAll(localElasticsearchLocation, "elasticsearch", "localhost")
	cfg.Services.ValidPreparations.SearchIndexPath = strings.ReplaceAll(localElasticsearchLocation, "elasticsearch", "localhost")

	return saveConfig(ctx, filePath, cfg, true, true)
}

func main() {
	ctx := context.Background()

	for filePath, fun := range files {
		if err := fun(ctx, filePath); err != nil {
			log.Fatalf("error rendering %s: %v", filePath, err)
		}
	}
}
