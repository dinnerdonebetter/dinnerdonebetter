package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	recipestepproductsservice "github.com/prixfixeco/api_server/internal/services/recipestepproducts"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	"github.com/prixfixeco/api_server/internal/encoding"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/messagequeue/redis"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	metricscfg "github.com/prixfixeco/api_server/internal/observability/metrics/config"
	"github.com/prixfixeco/api_server/internal/observability/metrics/prometheus"
	"github.com/prixfixeco/api_server/internal/observability/tracing/cloudtrace"
	tracingcfg "github.com/prixfixeco/api_server/internal/observability/tracing/config"
	"github.com/prixfixeco/api_server/internal/observability/tracing/jaeger"
	"github.com/prixfixeco/api_server/internal/routing"
	"github.com/prixfixeco/api_server/internal/server"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	householdinvitationsservice "github.com/prixfixeco/api_server/internal/services/householdinvitations"
	householdsservice "github.com/prixfixeco/api_server/internal/services/households"
	mealplanoptionsservice "github.com/prixfixeco/api_server/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/api_server/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/api_server/internal/services/mealplans"
	mealsservice "github.com/prixfixeco/api_server/internal/services/meals"
	recipesservice "github.com/prixfixeco/api_server/internal/services/recipes"
	recipestepingredientsservice "github.com/prixfixeco/api_server/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/api_server/internal/services/recipestepinstruments"
	recipestepsservice "github.com/prixfixeco/api_server/internal/services/recipesteps"
	usersservice "github.com/prixfixeco/api_server/internal/services/users"
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
	defaultPort              = 8000
	defaultCookieDomain      = ".prixfixe.local"
	debugCookieSecret        = "HEREISA32CHARSECRETWHICHISMADEUP"
	debugCookieSigningKey    = "DIFFERENT32CHARSECRETTHATIMADEUP"
	devPostgresDBConnDetails = "postgres://dbuser:hunter2@pgdatabase:5432/prixfixe?sslmode=disable"
	defaultCookieName        = authservice.DefaultCookieName

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	// message provider topics
	dataChangesTopicName = "data_changes"

	pasetoSecretSize      = 32
	maxAttempts           = 50
	defaultPASETOLifetime = 1 * time.Minute

	contentTypeJSON    = "application/json"
	workerQueueAddress = "worker_queue:6379"

	pasteoIssuer = "prixfixe_service"
)

var (
	examplePASETOKey = generatePASETOKey()

	localRoutingConfig = routing.Config{
		Provider:            routing.ChiProvider,
		SilenceRouteLogging: false,
	}

	devRoutingConfig = routing.Config{
		Provider:            routing.ChiProvider,
		SilenceRouteLogging: true,
	}

	devEnvLogConfig = logcfg.Config{
		Level:    logging.DebugLevel,
		Provider: logcfg.ProviderZerolog,
	}

	localLogConfig = logcfg.Config{
		Level:    logging.DebugLevel,
		Provider: logcfg.ProviderZerolog,
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
		BlockKey:   debugCookieSigningKey,
		Lifetime:   authservice.DefaultCookieLifetime,
		SecureOnly: false,
	}

	localMetricsConfig = metricscfg.Config{
		Provider: "prometheus",
		Prometheus: &prometheus.Config{
			RuntimeMetricsCollectionInterval: time.Second,
		},
	}

	localTracingConfig = tracingcfg.Config{
		Provider: "jaeger",
		Jaeger: &jaeger.Config{
			SpanCollectionProbability: 1,
			CollectorEndpoint:         "http://tracing-server:14268/api/traces",
			ServiceName:               "prixfixe_service",
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
		if err := cfg.ValidateWithContext(ctx, true); err != nil {
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
	"environments/dev/config_files/service-config.json":               devEnvironmentServerConfig,
	"environments/local/config_files/service-config.json":             localDevelopmentConfig,
	"environments/testing/config_files/integration-tests-config.json": integrationTestConfig,
}

func generatePASETOKey() []byte {
	b := make([]byte, pasetoSecretSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return b
}

func buildDevEnvironmentServerConfig() *config.InstanceConfig {
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
		Routing: devRoutingConfig,
		Meta: config.MetaSettings{
			Debug:   true,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Events: msgconfig.Config{
			Consumers: msgconfig.ProviderConfig{
				Provider: msgconfig.ProviderRedis,
			},
			Publishers: msgconfig.ProviderConfig{
				Provider: msgconfig.ProviderPubSub,
			},
		},
		Email:        emailConfig,
		CustomerData: customerDataPlatformConfig,
		Server: server.Config{
			Debug:           true,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: dbconfig.Config{
			Debug:           true,
			RunMigrations:   true,
			MaxPingAttempts: maxAttempts,
		},
		Observability: observability.Config{
			Logging: devEnvLogConfig,
			Metrics: metricscfg.Config{},
			Tracing: tracingcfg.Config{
				Provider: tracingcfg.ProviderCloudTrace,
				CloudTrace: &cloudtrace.Config{
					ProjectID:                 "prixfixe-dev",
					ServiceName:               "prixfixe_api",
					SpanCollectionProbability: 1,
				},
			},
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
		Services: config.ServicesConfigurations{
			Auth: authservice.Config{
				PASETO: authservice.PASETOConfig{
					Issuer:   pasteoIssuer,
					Lifetime: defaultPASETOLifetime,
				},
				Cookies:               cookieConfig,
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
			},
		},
	}

	return cfg
}

func devEnvironmentServerConfig(ctx context.Context, filePath string) error {
	cfg := buildDevEnvironmentServerConfig()

	return saveConfig(ctx, filePath, cfg, false, false)
}

func localDevelopmentConfig(ctx context.Context, filePath string) error {
	cfg := &config.InstanceConfig{
		Routing: localRoutingConfig,
		Meta: config.MetaSettings{
			Debug:   true,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Events: msgconfig.Config{
			Consumers: msgconfig.ProviderConfig{
				Provider: msgconfig.ProviderRedis,
				RedisConfig: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
			Publishers: msgconfig.ProviderConfig{
				Provider: msgconfig.ProviderRedis,
				RedisConfig: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
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
			Logging: localLogConfig,
			Metrics: localMetricsConfig,
			Tracing: localTracingConfig,
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
		Services: config.ServicesConfigurations{
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Households: householdsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			HouseholdInvitations: householdinvitationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Auth: authservice.Config{
				PASETO: authservice.PASETOConfig{
					Issuer:       pasteoIssuer,
					Lifetime:     defaultPASETOLifetime,
					LocalModeKey: examplePASETOKey,
				},
				Cookies:               localCookies,
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
				DataChangesTopicName:  dataChangesTopicName,
			},
			Webhooks: webhooksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Websockets: websocketsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstruments: validinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredients: validingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparations: validpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Meals: mealsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Recipes: recipesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeSteps: recipestepsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlans: mealplansservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
		},
	}

	return saveConfig(ctx, filePath, cfg, true, true)
}

func buildIntegrationTestsConfig() *config.InstanceConfig {
	return &config.InstanceConfig{
		Routing: localRoutingConfig,
		Meta: config.MetaSettings{
			Debug:   false,
			RunMode: testingEnv,
		},
		Events: msgconfig.Config{
			Consumers: msgconfig.ProviderConfig{
				Provider: msgconfig.ProviderRedis,
				RedisConfig: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
			Publishers: msgconfig.ProviderConfig{
				Provider: msgconfig.ProviderRedis,
				RedisConfig: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
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
			Debug:             true,
			RunMigrations:     true,
			MaxPingAttempts:   maxAttempts,
			ConnectionDetails: devPostgresDBConnDetails,
		},
		Observability: observability.Config{
			Logging: localLogConfig,
			Metrics: localMetricsConfig,
			Tracing: localTracingConfig,
		},
		Uploads: uploads.Config{
			Debug: false,
			Storage: storage.Config{
				Provider:   "memory",
				BucketName: "avatars",
				S3Config:   nil,
			},
		},
		Services: config.ServicesConfigurations{
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Households: householdsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			HouseholdInvitations: householdinvitationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Auth: authservice.Config{
				PASETO: authservice.PASETOConfig{
					Issuer:       pasteoIssuer,
					Lifetime:     defaultPASETOLifetime,
					LocalModeKey: examplePASETOKey,
				},
				Cookies: authservice.CookieConfig{
					Name:       defaultCookieName,
					Domain:     defaultCookieDomain,
					HashKey:    debugCookieSecret,
					BlockKey:   debugCookieSigningKey,
					Lifetime:   authservice.DefaultCookieLifetime,
					SecureOnly: false,
				},
				Debug:                 false,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
				DataChangesTopicName:  dataChangesTopicName,
			},
			Webhooks: webhooksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Websockets: websocketsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstruments: validinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredients: validingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparations: validpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Meals: mealsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Recipes: recipesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeSteps: recipestepsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlans: mealplansservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
		},
	}
}

func integrationTestConfig(ctx context.Context, filePath string) error {
	cfg := buildIntegrationTestsConfig()

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
