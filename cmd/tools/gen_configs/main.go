package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	analyticsconfig "github.com/prixfixeco/backend/internal/analytics/config"
	"github.com/prixfixeco/backend/internal/config"
	dbconfig "github.com/prixfixeco/backend/internal/database/config"
	emailconfig "github.com/prixfixeco/backend/internal/email/config"
	"github.com/prixfixeco/backend/internal/email/sendgrid"
	"github.com/prixfixeco/backend/internal/encoding"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/messagequeue/redis"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	logcfg "github.com/prixfixeco/backend/internal/observability/logging/config"
	metricscfg "github.com/prixfixeco/backend/internal/observability/metrics/config"
	"github.com/prixfixeco/backend/internal/observability/metrics/prometheus"
	"github.com/prixfixeco/backend/internal/observability/tracing/cloudtrace"
	tracingcfg "github.com/prixfixeco/backend/internal/observability/tracing/config"
	"github.com/prixfixeco/backend/internal/observability/tracing/jaeger"
	"github.com/prixfixeco/backend/internal/routing"
	"github.com/prixfixeco/backend/internal/server"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	householdinvitationsservice "github.com/prixfixeco/backend/internal/services/householdinvitations"
	householdsservice "github.com/prixfixeco/backend/internal/services/households"
	mealplaneventsservice "github.com/prixfixeco/backend/internal/services/mealplanevents"
	"github.com/prixfixeco/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/prixfixeco/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/backend/internal/services/mealplans"
	"github.com/prixfixeco/backend/internal/services/mealplantasks"
	mealsservice "github.com/prixfixeco/backend/internal/services/meals"
	recipepreptasksservice "github.com/prixfixeco/backend/internal/services/recipepreptasks"
	recipesservice "github.com/prixfixeco/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/prixfixeco/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/prixfixeco/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/prixfixeco/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/prixfixeco/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/prixfixeco/backend/internal/services/recipestepvessels"
	usersservice "github.com/prixfixeco/backend/internal/services/users"
	validingredientmeasurementunitsservice "github.com/prixfixeco/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/prixfixeco/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/prixfixeco/backend/internal/services/validingredients"
	validingredientstateingredientsservice "github.com/prixfixeco/backend/internal/services/validingredientstateingredients"
	validingredientstatesservice "github.com/prixfixeco/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/prixfixeco/backend/internal/services/validinstruments"
	validmeasurementconversionsservice "github.com/prixfixeco/backend/internal/services/validmeasurementconversions"
	validmeasurementunitsservice "github.com/prixfixeco/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/prixfixeco/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/prixfixeco/backend/internal/services/validpreparations"
	webhooksservice "github.com/prixfixeco/backend/internal/services/webhooks"
	websocketsservice "github.com/prixfixeco/backend/internal/services/websockets"
	"github.com/prixfixeco/backend/internal/storage"
	"github.com/prixfixeco/backend/internal/uploads"
)

const (
	defaultPort         = 8000
	defaultCookieDomain = ".prixfixe.local"
	/* #nosec G101 */
	debugCookieSecret = "HEREISA32CHARSECRETWHICHISMADEUP"
	/* #nosec G101 */
	debugCookieSigningKey    = "DIFFERENT32CHARSECRETTHATIMADEUP"
	devPostgresDBConnDetails = "postgres://dbuser:hunter2@pgdatabase:5432/prixfixe?sslmode=disable"
	defaultCookieName        = authservice.DefaultCookieName

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	// message provider topics.
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
		Provider:               routing.ChiProvider,
		EnableCORSForLocalhost: true,
		SilenceRouteLogging:    false,
	}

	devRoutingConfig = routing.Config{
		Provider:               routing.ChiProvider,
		EnableCORSForLocalhost: true,
		SilenceRouteLogging:    true,
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
)

func saveConfig(ctx context.Context, outputPath string, cfg *config.InstanceConfig, indent, validate bool) error {
	/* #nosec G301 */
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o777); err != nil {
		// okay, who gives a shit?
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

	/* #nosec G306 */
	return os.WriteFile(outputPath, output, 0o644)
}

type configFunc func(ctx context.Context, filePath string) error

var files = map[string]configFunc{
	"environments/dev/config_files/service-config.json":                              devEnvironmentServerConfig,
	"environments/local/config_files/service-config.json":                            localDevelopmentServiceConfig,
	"environments/local/config_files/queue-loader-config.json":                       localDevelopmentWorkerConfig,
	"environments/local/config_files/meal-plan-finalizer-config.json":                localDevelopmentWorkerConfig,
	"environments/local/config_files/meal-plan-task-creator-config.json":             localDevelopmentWorkerConfig,
	"environments/local/config_files/meal-plan-grocery-list-initializer-config.json": localDevelopmentWorkerConfig,
	"environments/testing/config_files/integration-tests-config.json":                integrationTestConfig,
}

func generatePASETOKey() []byte {
	b := make([]byte, pasetoSecretSize)

	return b
}

func buildDevEnvironmentServerConfig() *config.InstanceConfig {
	cookieConfig := authservice.CookieConfig{
		Name:       defaultCookieName,
		Domain:     ".prixfixe.dev",
		Lifetime:   (24 * time.Hour) * 30,
		SecureOnly: true,
	}

	emailConfig := emailconfig.Config{
		Provider: emailconfig.ProviderSendgrid,
		Sendgrid: &sendgrid.Config{
			WebAppURL: "https://www.prixfixe.dev",
		},
	}

	analyticsConfig := analyticsconfig.Config{
		Provider: analyticsconfig.ProviderSegment,
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
			Consumers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
			},
			Publishers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderPubSub,
			},
		},
		Email:     emailConfig,
		Analytics: analyticsConfig,
		Server: server.Config{
			Debug:           true,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: dbconfig.Config{
			Debug:           true,
			LogQueries:      true,
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
			Users: usersservice.Config{
				PublicMediaURLPrefix: "https://media.prixfixe.dev/avatars",
				Uploads: uploads.Config{
					Debug: true,
					Storage: storage.Config{
						UploadFilenameKey: "avatar",
						Provider:          storage.GCPCloudStorageProvider,
						BucketName:        "media.prixfixe.dev",
						BucketPrefix:      "avatars/",
						GCPConfig: &storage.GCPConfig{
							BucketName: "media.prixfixe.dev",
						},
					},
				},
			},
			Recipes: recipesservice.Config{
				// note, this should effectively be "https://media.prixfixe.dev" + bucket prefix
				PublicMediaURLPrefix: "https://media.prixfixe.dev/recipe_media",
				Uploads: uploads.Config{
					Debug: true,
					Storage: storage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          storage.GCPCloudStorageProvider,
						BucketName:        "media.prixfixe.dev",
						BucketPrefix:      "recipe_media/",
						GCPConfig: &storage.GCPConfig{
							BucketName: "media.prixfixe.dev",
						},
					},
				},
			},
			RecipeSteps: recipestepsservice.Config{
				// note, this should effectively be "https://media.prixfixe.dev" + bucket prefix
				PublicMediaURLPrefix: "https://media.prixfixe.dev/recipe_media",
				Uploads: uploads.Config{
					Debug: true,
					Storage: storage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          storage.GCPCloudStorageProvider,
						BucketName:        "media.prixfixe.dev",
						BucketPrefix:      "recipe_media/",
						GCPConfig: &storage.GCPConfig{
							BucketName: "media.prixfixe.dev",
						},
					},
				},
			},
		},
	}

	return cfg
}

func devEnvironmentServerConfig(ctx context.Context, filePath string) error {
	cfg := buildDevEnvironmentServerConfig()

	return saveConfig(ctx, filePath, cfg, false, false)
}

func buildDevConfig() *config.InstanceConfig {
	return &config.InstanceConfig{
		Routing: localRoutingConfig,
		Meta: config.MetaSettings{
			Debug:   true,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Events: msgconfig.Config{
			Consumers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				RedisConfig: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
			Publishers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				RedisConfig: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
		},
		Server: localServer,
		Database: dbconfig.Config{
			Debug:             true,
			RunMigrations:     true,
			LogQueries:        true,
			MaxPingAttempts:   maxAttempts,
			ConnectionDetails: devPostgresDBConnDetails,
		},
		Observability: observability.Config{
			Logging: localLogConfig,
			Metrics: localMetricsConfig,
			Tracing: localTracingConfig,
		},
		Services: config.ServicesConfigurations{
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				Uploads: uploads.Config{
					Debug: true,
					Storage: storage.Config{
						UploadFilenameKey: "avatar",
						Provider:          storage.FilesystemProvider,
						BucketName:        "avatars",
						FilesystemConfig: &storage.FilesystemConfig{
							RootDirectory: "/uploads",
						},
					},
				},
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
			ValidIngredientStates: validingredientstatesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidMeasurementUnits: validmeasurementunitsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientStateIngredients: validingredientstateingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparationInstruments: validpreparationinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstrumentMeasurementUnits: validingredientmeasurementunitsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Meals: mealsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Recipes: recipesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				PublicMediaURLPrefix: "https://example.website.lol",
				Uploads: uploads.Config{
					Debug: true,
					Storage: storage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          storage.FilesystemProvider,
						BucketName:        "recipe_media",
						FilesystemConfig: &storage.FilesystemConfig{
							RootDirectory: "/uploads",
						},
					},
				},
			},
			RecipeSteps: recipestepsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				PublicMediaURLPrefix: "https://example.website.lol",
				Uploads: uploads.Config{
					Debug: true,
					Storage: storage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          storage.FilesystemProvider,
						BucketName:        "recipe_media",
						FilesystemConfig: &storage.FilesystemConfig{
							RootDirectory: "/uploads",
						},
					},
				},
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepVessels: recipestepvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepCompletionConditions: recipestepcompletionconditionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlans: mealplansservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanEvents: mealplaneventsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanTasks: mealplantasks.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipePrepTasks: recipepreptasksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanGroceryListItems: mealplangrocerylistitems.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidMeasurementConversions: validmeasurementconversionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
		},
	}
}

func localDevelopmentServiceConfig(ctx context.Context, filePath string) error {
	cfg := buildDevConfig()

	return saveConfig(ctx, filePath, cfg, true, true)
}

func localDevelopmentWorkerConfig(ctx context.Context, filePath string) error {
	cfg := buildDevConfig()

	cfg.Database.LogQueries = false
	cfg.Database.RunMigrations = false

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
			Consumers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				RedisConfig: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
			Publishers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				RedisConfig: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Server: server.Config{
			Debug:           false,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: dbconfig.Config{
			Debug:             true,
			RunMigrations:     true,
			LogQueries:        true,
			MaxPingAttempts:   maxAttempts,
			ConnectionDetails: devPostgresDBConnDetails,
		},
		Observability: observability.Config{
			Logging: logcfg.Config{
				Level:    logging.InfoLevel,
				Provider: logcfg.ProviderZerolog,
			},
			Metrics: localMetricsConfig,
			Tracing: localTracingConfig,
		},
		Services: config.ServicesConfigurations{
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				Uploads: uploads.Config{
					Debug: false,
					Storage: storage.Config{
						Provider:   "memory",
						BucketName: "avatars",
						S3Config:   nil,
					},
				},
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
			ValidIngredientStates: validingredientstatesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidMeasurementUnits: validmeasurementunitsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientStateIngredients: validingredientstateingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparationInstruments: validpreparationinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstrumentMeasurementUnits: validingredientmeasurementunitsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Meals: mealsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Recipes: recipesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				PublicMediaURLPrefix: "https://media.example.website/lol",
				Uploads: uploads.Config{
					Debug: false,
					Storage: storage.Config{
						Provider:   "memory",
						BucketName: "recipes",
						S3Config:   nil,
					},
				},
			},
			RecipeSteps: recipestepsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				PublicMediaURLPrefix: "https://media.example.website/lol",
				Uploads: uploads.Config{
					Debug: false,
					Storage: storage.Config{
						Provider:   "memory",
						BucketName: "recipes",
						S3Config:   nil,
					},
				},
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepVessels: recipestepvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepCompletionConditions: recipestepcompletionconditionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlans: mealplansservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanEvents: mealplaneventsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanTasks: mealplantasks.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipePrepTasks: recipepreptasksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanGroceryListItems: mealplangrocerylistitems.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidMeasurementConversions: validmeasurementconversionsservice.Config{
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
