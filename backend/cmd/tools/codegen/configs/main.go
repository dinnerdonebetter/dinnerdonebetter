package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	logcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics/otelgrpc"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing/oteltrace"
	"github.com/dinnerdonebetter/backend/internal/pkg/testutils"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/search/text/algolia"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/dataprivacy"
	householdinstrumentownershipsservice "github.com/dinnerdonebetter/backend/internal/services/householdinstrumentownerships"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	"github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	"github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/oauth2clients"
	recipepreptasksservice "github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	reciperatingsservice "github.com/dinnerdonebetter/backend/internal/services/reciperatings"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientgroupsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	validingredientstateingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	validingredientstatesservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	validmeasurementconversionsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunitconversions"
	validmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	validpreparationvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationvessels"
	validvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validvessels"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
	workersservice "github.com/dinnerdonebetter/backend/internal/services/workers"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
)

const (
	defaultPort = 8000
	/* #nosec G101 */
	debugCookieHashKey = "HEREISA32CHARSECRETWHICHISMADEUP"
	/* #nosec G101 */
	devPostgresDBConnDetails = "postgres://dbuser:hunter2@pgdatabase:5432/dinner-done-better?sslmode=disable"

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	// message provider topics.
	dataChangesTopicName              = "data_changes"
	outboundEmailsTopicName           = "outbound_emails"
	searchIndexRequestsTopicName      = "search_index_requests"
	userDataAggregationTopicName      = "user_data_aggregation_requests"
	webhookExecutionRequestsTopicName = "webhook_execution_requests"

	maxAttempts = 50

	contentTypeJSON               = "application/json"
	workerQueueAddress            = "worker_queue:6379"
	localOAuth2TokenEncryptionKey = debugCookieHashKey
)

func saveConfig(ctx context.Context, outputPath string, cfg *config.APIServiceConfig, indent, validate bool) error {
	/* #nosec G301 */
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o750); err != nil {
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

func buildLocalDevConfig() *config.APIServiceConfig {
	return &config.APIServiceConfig{
		Routing: routingcfg.Config{
			Provider:               routingcfg.ChiProvider,
			EnableCORSForLocalhost: true,
			SilenceRouteLogging:    false,
		},
		Queues: config.QueuesConfig{
			DataChangesTopicName:              dataChangesTopicName,
			OutboundEmailsTopicName:           outboundEmailsTopicName,
			SearchIndexRequestsTopicName:      searchIndexRequestsTopicName,
			UserDataAggregationTopicName:      userDataAggregationTopicName,
			WebhookExecutionRequestsTopicName: webhookExecutionRequestsTopicName,
		},
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
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
			Publishers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
		},
		Search: searchcfg.Config{
			Algolia:  &algolia.Config{},
			Provider: searchcfg.AlgoliaProvider,
		},
		Server: http.Config{
			Debug:           true,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: dbconfig.Config{
			OAuth2TokenEncryptionKey: localOAuth2TokenEncryptionKey,
			Debug:                    true,
			RunMigrations:            true,
			LogQueries:               true,
			MaxPingAttempts:          maxAttempts,
			PingWaitPeriod:           time.Second,
			ConnectionDetails:        devPostgresDBConnDetails,
		},
		Observability: observability.Config{
			Logging: logcfg.Config{
				Level:          logging.DebugLevel,
				Provider:       logcfg.ProviderSlog,
				OutputFilepath: "/var/log/dinnerdonebetter/api-service.log",
			},
			Metrics: metricscfg.Config{
				Otel: &otelgrpc.Config{
					BaseName:           "ddb.api-svc",
					CollectorEndpoint:  "otel-collector:4318",
					CollectionInterval: time.Second,
				},
				Provider: metricscfg.ProviderOtel,
			},
			Tracing: tracingcfg.Config{
				Provider: tracingcfg.ProviderOtel,
				Otel: &oteltrace.Config{
					SpanCollectionProbability: 1,
					CollectorEndpoint:         "http://tracing-server:14268/api/traces",
					ServiceName:               "dinner_done_better_service",
				},
			},
		},
		Services: config.ServicesConfig{
			AuditLogEntries: auditlogentriesservice.Config{},
			Auth: authservice.Config{
				OAuth2: authservice.OAuth2Config{
					Domain:               "http://localhost:9000",
					AccessTokenLifespan:  time.Hour,
					RefreshTokenLifespan: time.Hour,
					Debug:                false,
				},
				SSO: authservice.SSOConfigs{
					Google: authservice.GoogleSSOConfig{
						CallbackURL: "https://app.dinnerdonebetter.dev/auth/google/callback",
					},
				},
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
				DataChangesTopicName:  dataChangesTopicName,
				JWTAudience:           "localhost",
				JWTSigningKey:         base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
				JWTLifetime:           5 * time.Minute,
			},
			DataPrivacy: dataprivacyservice.Config{
				Uploads: uploads.Config{
					Storage: objectstorage.Config{
						FilesystemConfig: &objectstorage.FilesystemConfig{RootDirectory: "/tmp"},
						BucketName:       "userdata",
						Provider:         objectstorage.FilesystemProvider,
					},
					Debug: false,
				},
				DataChangesTopicName:         dataChangesTopicName,
				UserDataAggregationTopicName: userDataAggregationTopicName,
			},
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "avatar",
						Provider:          objectstorage.FilesystemProvider,
						BucketName:        "avatars",
						FilesystemConfig: &objectstorage.FilesystemConfig{
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
			Webhooks: webhooksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstruments: validinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidVessels: validvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparationVessels: validpreparationvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredients: validingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientGroups: validingredientgroupsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparations: validpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserIngredientPreferences: useringredientpreferencesservice.Config{
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
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.FilesystemProvider,
						BucketName:        "recipe_media",
						FilesystemConfig: &objectstorage.FilesystemConfig{
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
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.FilesystemProvider,
						BucketName:        "recipe_media",
						FilesystemConfig: &objectstorage.FilesystemConfig{
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
			ValidMeasurementUnitConversions: validmeasurementconversionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettings: servicesettings.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettingConfigurations: servicesettingconfigurations.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeRatings: reciperatingsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			HouseholdInstrumentOwnerships: householdinstrumentownershipsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			OAuth2Clients: oauth2clientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Workers: workersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserNotifications: usernotificationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
		},
	}
}

func buildLocalDevelopmentServiceConfig(local bool) func(context.Context, string) error {
	const localUploadsDir = "artifacts/uploads"
	const localRedisAddr = "localhost:6379"
	return func(ctx context.Context, filePath string) error {
		cfg := buildLocalDevConfig()

		if local {
			cfg.Database.ConnectionDetails = "postgres://dbuser:hunter2@localhost:5432/dinner-done-better?sslmode=disable"
			cfg.Events.Consumers.Redis.QueueAddresses = []string{localRedisAddr}
			cfg.Events.Publishers.Redis.QueueAddresses = []string{localRedisAddr}
			cfg.Services.Users.Uploads.Storage.FilesystemConfig.RootDirectory = localUploadsDir
			cfg.Services.Recipes.Uploads.Storage.FilesystemConfig.RootDirectory = localUploadsDir
			cfg.Services.RecipeSteps.Uploads.Storage.FilesystemConfig.RootDirectory = localUploadsDir
		}

		return saveConfig(ctx, filePath, cfg, true, true)
	}
}

func buildLocaldevKubernetesConfig() *config.APIServiceConfig {
	return &config.APIServiceConfig{
		Routing: routingcfg.Config{
			Provider:               routingcfg.ChiProvider,
			EnableCORSForLocalhost: true,
			SilenceRouteLogging:    false,
		},
		Queues: config.QueuesConfig{
			DataChangesTopicName:              dataChangesTopicName,
			OutboundEmailsTopicName:           outboundEmailsTopicName,
			SearchIndexRequestsTopicName:      searchIndexRequestsTopicName,
			UserDataAggregationTopicName:      userDataAggregationTopicName,
			WebhookExecutionRequestsTopicName: webhookExecutionRequestsTopicName,
		},
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
				Redis: redis.Config{
					QueueAddresses: []string{"redis-master.localdev.svc.cluster.local:6379"},
				},
			},
			Publishers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{"redis-master.localdev.svc.cluster.local:6379"},
				},
			},
		},
		Search: searchcfg.Config{
			Algolia:  &algolia.Config{},
			Provider: searchcfg.AlgoliaProvider,
		},
		Server: http.Config{
			Debug:           true,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: dbconfig.Config{
			OAuth2TokenEncryptionKey: localOAuth2TokenEncryptionKey,
			Debug:                    true,
			RunMigrations:            true,
			LogQueries:               true,
			MaxPingAttempts:          maxAttempts,
			PingWaitPeriod:           time.Second,
			ConnectionDetails:        "postgres://dbuser:hunter2@postgres-postgresql.localdev.svc.cluster.local:5432/dinner-done-better?sslmode=disable",
		},
		Observability: observability.Config{
			Logging: logcfg.Config{
				Level:    logging.DebugLevel,
				Provider: logcfg.ProviderSlog,
			},
			Tracing: tracingcfg.Config{
				Provider: tracingcfg.ProviderOtel,
				Otel: &oteltrace.Config{
					SpanCollectionProbability: 1,
					CollectorEndpoint:         "http://0.0.0.0:4317",
					ServiceName:               "dinner_done_better_service",
				},
			},
			Metrics: metricscfg.Config{
				Provider: tracingcfg.ProviderOtel,
				Otel: &otelgrpc.Config{
					CollectorEndpoint:  "http://0.0.0.0:4317",
					BaseName:           "ddb.api",
					CollectionInterval: 3 * time.Second,
				},
			},
		},
		Services: config.ServicesConfig{
			AuditLogEntries: auditlogentriesservice.Config{},
			Auth: authservice.Config{
				OAuth2: authservice.OAuth2Config{
					Domain:               "http://localhost:9000",
					AccessTokenLifespan:  time.Hour,
					RefreshTokenLifespan: time.Hour,
					Debug:                false,
				},
				SSO: authservice.SSOConfigs{
					Google: authservice.GoogleSSOConfig{
						CallbackURL: "https://app.dinnerdonebetter.dev/auth/google/callback",
					},
				},
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
				DataChangesTopicName:  dataChangesTopicName,
				JWTAudience:           "localhost",
				JWTSigningKey:         base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
				JWTLifetime:           5 * time.Minute,
			},
			DataPrivacy: dataprivacyservice.Config{
				Uploads: uploads.Config{
					Storage: objectstorage.Config{
						FilesystemConfig: &objectstorage.FilesystemConfig{RootDirectory: "/tmp"},
						BucketName:       "userdata",
						Provider:         objectstorage.FilesystemProvider,
					},
					Debug: false,
				},
				DataChangesTopicName:         dataChangesTopicName,
				UserDataAggregationTopicName: userDataAggregationTopicName,
			},
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "avatar",
						Provider:          objectstorage.FilesystemProvider,
						BucketName:        "avatars",
						FilesystemConfig: &objectstorage.FilesystemConfig{
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
			Webhooks: webhooksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstruments: validinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidVessels: validvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparationVessels: validpreparationvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredients: validingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientGroups: validingredientgroupsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparations: validpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserIngredientPreferences: useringredientpreferencesservice.Config{
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
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.FilesystemProvider,
						BucketName:        "recipe_media",
						FilesystemConfig: &objectstorage.FilesystemConfig{
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
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.FilesystemProvider,
						BucketName:        "recipe_media",
						FilesystemConfig: &objectstorage.FilesystemConfig{
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
			ValidMeasurementUnitConversions: validmeasurementconversionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettings: servicesettings.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettingConfigurations: servicesettingconfigurations.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeRatings: reciperatingsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			HouseholdInstrumentOwnerships: householdinstrumentownershipsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			OAuth2Clients: oauth2clientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Workers: workersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserNotifications: usernotificationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
		},
	}
}

func buildIntegrationTestsConfig() *config.APIServiceConfig {
	return &config.APIServiceConfig{
		Routing: routingcfg.Config{
			Provider:               routingcfg.ChiProvider,
			EnableCORSForLocalhost: true,
			SilenceRouteLogging:    false,
		},
		Meta: config.MetaSettings{
			Debug:   false,
			RunMode: testingEnv,
		},
		Queues: config.QueuesConfig{
			DataChangesTopicName:              dataChangesTopicName,
			OutboundEmailsTopicName:           outboundEmailsTopicName,
			SearchIndexRequestsTopicName:      searchIndexRequestsTopicName,
			UserDataAggregationTopicName:      userDataAggregationTopicName,
			WebhookExecutionRequestsTopicName: webhookExecutionRequestsTopicName,
		},
		Events: msgconfig.Config{
			Consumers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
			Publishers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Server: http.Config{
			Debug:           false,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: dbconfig.Config{
			OAuth2TokenEncryptionKey: localOAuth2TokenEncryptionKey,
			Debug:                    true,
			RunMigrations:            true,
			LogQueries:               true,
			MaxPingAttempts:          maxAttempts,
			PingWaitPeriod:           1500 * time.Millisecond,
			ConnectionDetails:        devPostgresDBConnDetails,
		},
		Observability: observability.Config{
			Logging: logcfg.Config{
				Level:    logging.InfoLevel,
				Provider: logcfg.ProviderSlog,
			},
			Tracing: tracingcfg.Config{
				Provider: tracingcfg.ProviderOtel,
				Otel: &oteltrace.Config{
					SpanCollectionProbability: 1,
					CollectorEndpoint:         "http://tracing-server:14268/api/traces",
					ServiceName:               "dinner_done_better_service",
				},
			},
		},
		Services: config.ServicesConfig{
			AuditLogEntries: auditlogentriesservice.Config{},
			Auth: authservice.Config{
				OAuth2: authservice.OAuth2Config{
					Domain:               "http://localhost:9000",
					AccessTokenLifespan:  time.Hour,
					RefreshTokenLifespan: time.Hour,
					Debug:                false,
				},
				Debug:                 false,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
				DataChangesTopicName:  dataChangesTopicName,
				JWTAudience:           "localhost",
				JWTSigningKey:         base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
				JWTLifetime:           5 * time.Minute,
			},
			DataPrivacy: dataprivacyservice.Config{
				Uploads: uploads.Config{
					Storage: objectstorage.Config{
						FilesystemConfig: &objectstorage.FilesystemConfig{RootDirectory: "/tmp"},
						BucketName:       "userdata",
						Provider:         objectstorage.FilesystemProvider,
					},
					Debug: false,
				},
				DataChangesTopicName:         dataChangesTopicName,
				UserDataAggregationTopicName: userDataAggregationTopicName,
			},
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				Uploads: uploads.Config{
					Debug: false,
					Storage: objectstorage.Config{
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
			Webhooks: webhooksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstruments: validinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidVessels: validvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparationVessels: validpreparationvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredients: validingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientGroups: validingredientgroupsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparations: validpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserIngredientPreferences: useringredientpreferencesservice.Config{
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
					Storage: objectstorage.Config{
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
					Storage: objectstorage.Config{
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
			ValidMeasurementUnitConversions: validmeasurementconversionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettings: servicesettings.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettingConfigurations: servicesettingconfigurations.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeRatings: reciperatingsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			HouseholdInstrumentOwnerships: householdinstrumentownershipsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Workers: workersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserNotifications: usernotificationsservice.Config{
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
	envConfigs := map[string]*environmentConfigSet{
		"deploy/kustomize/environments/localdev/configs": {
			renderPretty: true,
			rootConfig:   buildLocaldevKubernetesConfig(),
		},
		"environments/testing/config_files": {
			renderPretty:         true,
			apiServiceConfigPath: "integration-tests-config.json",
			rootConfig:           buildIntegrationTestsConfig(),
		},
		"environments/localdev/config_files": {
			renderPretty: true,
			rootConfig:   buildLocalDevConfig(),
		},
	}

	for p, cfg := range envConfigs {
		if err := cfg.Render(p); err != nil {
			panic(err)
		}
	}
}
