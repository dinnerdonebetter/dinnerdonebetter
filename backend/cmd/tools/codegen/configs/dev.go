package main

import (
	"time"

	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/analytics/segment"
	"github.com/dinnerdonebetter/backend/internal/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	emailcfg "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/email/sendgrid"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics/otelgrpc"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing/cloudtrace"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/search/text/algolia"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/dataprivacy"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	validingredientstatesservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	validmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	validvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validvessels"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
)

func buildDevEnvironmentServerConfig() *config.APIServiceConfig {
	emailConfig := emailcfg.Config{
		Provider: emailcfg.ProviderSendgrid,
		Sendgrid: &sendgrid.Config{},
	}

	analyticsConfig := analyticscfg.Config{
		Provider: analyticscfg.ProviderSegment,
		Segment:  &segment.Config{APIToken: ""},
	}

	cfg := &config.APIServiceConfig{
		Routing: routingcfg.Config{
			ServiceName:            otelServiceName,
			Provider:               routingcfg.ChiProvider,
			EnableCORSForLocalhost: true,
			SilenceRouteLogging:    false,
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
				Provider: msgconfig.ProviderPubSub,
			},
			Publisher: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderPubSub,
			},
		},
		Email:     emailConfig,
		Analytics: analyticsConfig,
		Server: http.Config{
			Debug:           true,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Search: textsearchcfg.Config{
			Algolia:  &algolia.Config{},
			Provider: textsearchcfg.AlgoliaProvider,
		},
		Database: databasecfg.Config{
			Debug:           true,
			LogQueries:      true,
			RunMigrations:   true,
			MaxPingAttempts: maxAttempts,
			PingWaitPeriod:  time.Second,
			ConnectionDetails: databasecfg.ConnectionDetails{
				Username:   "api_db_user",
				Password:   "",
				Database:   "dinner-done-better",
				Host:       "",
				Port:       5432,
				DisableSSL: false,
			},
		},
		Observability: observability.Config{
			Logging: loggingcfg.Config{
				Level:    logging.DebugLevel,
				Provider: loggingcfg.ProviderSlog,
			},
			Metrics: metricscfg.Config{
				Provider: tracingcfg.ProviderCloudTrace,
				Otel: &otelgrpc.Config{
					ServiceName:        otelServiceName,
					CollectorEndpoint:  "localhost:4317",
					CollectionInterval: 3 * time.Second,
				},
			},
			Tracing: tracingcfg.Config{
				Provider:                  tracingcfg.ProviderCloudTrace,
				ServiceName:               otelServiceName,
				SpanCollectionProbability: 1,
				CloudTrace: &cloudtrace.Config{
					ProjectID: "dinner-done-better-dev",
				},
			},
		},
		Services: config.ServicesConfig{
			Auth: &authservice.Config{
				OAuth2: authservice.OAuth2Config{
					Domain:               "https://dinnerdonebetter.dev",
					AccessTokenLifespan:  time.Hour,
					RefreshTokenLifespan: time.Hour,
					Debug:                false,
				},
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
				JWTAudience:           "https://api.dinnerdonebetter.dev",
				JWTLifetime:           5 * time.Minute,
			},
			DataPrivacy: &dataprivacyservice.Config{
				Uploads: uploads.Config{
					Storage: objectstorage.Config{
						GCPConfig:  &objectstorage.GCPConfig{BucketName: "userdata.dinnerdonebetter.dev"},
						BucketName: "userdata.dinnerdonebetter.dev",
						Provider:   objectstorage.GCPCloudStorageProvider,
					},
					Debug: false,
				},
			},
			Users: &usersservice.Config{
				PublicMediaURLPrefix: "https://media.dinnerdonebetter.dev/avatars",
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "avatar",
						Provider:          objectstorage.GCPCloudStorageProvider,
						BucketName:        "media.dinnerdonebetter.dev",
						BucketPrefix:      "avatars/",
						GCPConfig: &objectstorage.GCPConfig{
							BucketName: "media.dinnerdonebetter.dev",
						},
					},
				},
			},
			Recipes: &recipesservice.Config{
				// note, this should effectively be "https://media.dinnerdonebetter.dev" + bucket prefix
				UseSearchService:     true,
				PublicMediaURLPrefix: "https://media.dinnerdonebetter.dev/recipe_media",
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.GCPCloudStorageProvider,
						BucketName:        "media.dinnerdonebetter.dev",
						BucketPrefix:      "recipe_media/",
						GCPConfig: &objectstorage.GCPConfig{
							BucketName: "media.dinnerdonebetter.dev",
						},
					},
				},
			},
			RecipeSteps: &recipestepsservice.Config{
				// note, this should effectively be "https://media.dinnerdonebetter.dev" + bucket prefix
				PublicMediaURLPrefix: "https://media.dinnerdonebetter.dev/recipe_media",
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.GCPCloudStorageProvider,
						BucketName:        "media.dinnerdonebetter.dev",
						BucketPrefix:      "recipe_media/",
						GCPConfig: &objectstorage.GCPConfig{
							BucketName: "media.dinnerdonebetter.dev",
						},
					},
				},
			},
			ValidIngredients: &validingredientsservice.Config{
				UseSearchService: true,
			},
			ValidIngredientStates: &validingredientstatesservice.Config{
				UseSearchService: true,
			},
			ValidInstruments: &validinstrumentsservice.Config{
				UseSearchService: true,
			},
			ValidVessels: &validvesselsservice.Config{
				UseSearchService: true,
			},
			ValidMeasurementUnits: &validmeasurementunitsservice.Config{
				UseSearchService: true,
			},
			ValidPreparations: &validpreparationsservice.Config{
				UseSearchService: true,
			},
			Meals: &mealsservice.Config{
				UseSearchService: true,
			},
		},
	}

	return cfg
}
