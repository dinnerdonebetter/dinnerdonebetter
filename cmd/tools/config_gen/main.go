package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	"github.com/prixfixeco/api_server/internal/encoding"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	metricscfg "github.com/prixfixeco/api_server/internal/observability/metrics/config"
	"github.com/prixfixeco/api_server/internal/observability/metrics/prometheus"
	tracingcfg "github.com/prixfixeco/api_server/internal/observability/tracing/config"
	jaeger "github.com/prixfixeco/api_server/internal/observability/tracing/jaeger"
	"github.com/prixfixeco/api_server/internal/observability/tracing/xray"
	"github.com/prixfixeco/api_server/internal/search"
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
	defaultPort              = 8000
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

	devEnvLogConfig = logcfg.Config{
		Level:    logging.InfoLevel,
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
		SecureOnly: true,
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
	"environments/dev/config_files/service-config.json":               devEnvironmentServerConfig,
	"environments/dev/config_files/worker-config.json":                devEnvironmentWorkerConfig,
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
		Logging: devEnvLogConfig,
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
			Metrics: metricscfg.Config{
				Provider: metricscfg.ProviderPrometheus,
				Prometheus: &prometheus.Config{
					RuntimeMetricsCollectionInterval: 5 * time.Second,
				},
				//Cloudwatch: &cloudwatch.Config{
				//	CollectorEndpoint:                "0.0.0.0:4317",
				//	MetricsCollectionInterval:        5 * time.Second,
				//	RuntimeMetricsCollectionInterval: 5 * time.Second,
				//},
			},
			Tracing: tracingcfg.Config{
				Provider: tracingcfg.ProviderXRay,
				XRay: &xray.Config{
					CollectorEndpoint:         "0.0.0.0:4317",
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
		Search: search.Config{
			Provider: search.ElasticsearchProvider,
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
				MinimumUsernameLength: 4,
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

func devEnvironmentWorkerConfig(ctx context.Context, filePath string) error {
	cfg := buildDevEnvironmentServerConfig()

	cfg.Observability.Tracing.Provider = ""

	return saveConfig(ctx, filePath, cfg, false, false)
}

func localDevelopmentConfig(ctx context.Context, filePath string) error {
	cfg := &config.InstanceConfig{
		Logging: localLogConfig,
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
			Metrics: metricscfg.Config{
				Provider: "prometheus",
				Prometheus: &prometheus.Config{
					RuntimeMetricsCollectionInterval: time.Second,
				},
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
					Issuer:       pasteoIssuer,
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
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstruments: validinstrumentsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			ValidIngredients: validingredientsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			ValidPreparations: validpreparationsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			Meals: mealsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			Recipes: recipesservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			RecipeSteps: recipestepsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			MealPlans: mealplansservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
		},
	}

	return saveConfig(ctx, filePath, cfg, true, true)
}

func buildIntegrationTestsConfig() *config.InstanceConfig {
	return &config.InstanceConfig{
		Logging: localLogConfig,
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
			Metrics: metricscfg.Config{
				Provider: "",
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
				MinimumUsernameLength: 4,
				MinimumPasswordLength: 8,
			},
			Webhooks: webhooksservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			Websockets: websocketsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstruments: validinstrumentsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			ValidIngredients: validingredientsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			ValidPreparations: validpreparationsservice.Config{
				SearchIndexPath:      localElasticsearchLocation,
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			Meals: mealsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			Recipes: recipesservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			RecipeSteps: recipestepsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			MealPlans: mealplansservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				PreWritesTopicName:   preWritesTopicName,
				PreUpdatesTopicName:  preUpdatesTopicName,
				PreArchivesTopicName: preArchivesTopicName,
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
