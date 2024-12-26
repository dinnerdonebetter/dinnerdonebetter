package main

import (
	"encoding/base64"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics/otelgrpc"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing/oteltrace"
	"github.com/dinnerdonebetter/backend/internal/pkg/testutils"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/search/text/algolia"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
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

func buildLocalDevConfig() *config.APIServiceConfig {
	return &config.APIServiceConfig{
		Routing: routingcfg.Config{
			ServiceName:            otelServiceName,
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
			Consumer: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
			Publisher: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
		},
		Search: textsearchcfg.Config{
			Algolia:  &algolia.Config{},
			Provider: textsearchcfg.AlgoliaProvider,
		},
		Server: http.Config{
			Debug:           true,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: databasecfg.Config{
			OAuth2TokenEncryptionKey: localOAuth2TokenEncryptionKey,
			Debug:                    true,
			RunMigrations:            true,
			LogQueries:               true,
			MaxPingAttempts:          maxAttempts,
			PingWaitPeriod:           time.Second,
			ConnectionDetails:        devPostgresDBConnectionDetails,
		},
		Observability: observability.Config{
			Logging: loggingcfg.Config{
				Level:          logging.DebugLevel,
				Provider:       loggingcfg.ProviderSlog,
				OutputFilepath: "/var/log/dinnerdonebetter/api-service.log",
			},
			Metrics: metricscfg.Config{
				Otel: &otelgrpc.Config{
					ServiceName:        otelServiceName,
					Insecure:           true,
					CollectorEndpoint:  "otel_collector:4317",
					CollectionInterval: time.Second,
				},
				Provider: metricscfg.ProviderOtel,
			},
			Tracing: tracingcfg.Config{
				Provider:                  tracingcfg.ProviderOtel,
				ServiceName:               otelServiceName,
				SpanCollectionProbability: 1,
				Otel: &oteltrace.Config{
					Insecure:          true,
					CollectorEndpoint: "otel_collector:4317",
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
