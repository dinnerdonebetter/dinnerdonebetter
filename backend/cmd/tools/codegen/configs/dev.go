package main

import (
	"time"

	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/analytics/segment"
	tokenscfg "github.com/dinnerdonebetter/backend/internal/authentication/tokens/config"
	"github.com/dinnerdonebetter/backend/internal/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	emailcfg "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/email/sendgrid"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/messagequeue/pubsub"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics/otelgrpc"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing/oteltrace"
	"github.com/dinnerdonebetter/backend/internal/pkg/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing/config"
	"github.com/dinnerdonebetter/backend/internal/search/text/algolia"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/core/dataprivacy"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/core/users"
	mealplanningservice "github.com/dinnerdonebetter/backend/internal/services/eating/meal_planning"
	recipemanagement "github.com/dinnerdonebetter/backend/internal/services/eating/recipe_management"
	validenumerations "github.com/dinnerdonebetter/backend/internal/services/eating/valid_enumerations"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
)

const (
	gcpProjectID = "dinner-done-better-dev"
)

func buildDevEnvironmentServerConfig() *config.APIServiceConfig {
	cfg := &config.APIServiceConfig{
		Routing: routingcfg.Config{
			Provider: routingcfg.ProviderChi,
			ChiConfig: &chi.Config{
				ServiceName:            otelServiceName,
				EnableCORSForLocalhost: true,
				SilenceRouteLogging:    false,
			},
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
				PubSub:   pubsub.Config{ProjectID: gcpProjectID},
			},
			Publisher: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderPubSub,
				PubSub:   pubsub.Config{ProjectID: gcpProjectID},
			},
		},
		Email: emailcfg.Config{
			Provider: emailcfg.ProviderSendgrid,
			CircuitBreakerConfig: &circuitbreaking.Config{
				Name:                   "dev_emailer",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
			Sendgrid: &sendgrid.Config{},
		},
		Analytics: analyticscfg.Config{
			Provider: analyticscfg.ProviderSegment,
			CircuitBreakerConfig: &circuitbreaking.Config{
				Name:                   "dev_analytics",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
			Segment: &segment.Config{APIToken: ""},
		},
		Server: http.Config{
			Debug:           true,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Search: textsearchcfg.Config{
			Algolia: &algolia.Config{},
			CircuitBreakerConfig: &circuitbreaking.Config{
				Name:                   "dev_text_searcher",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
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
				Level:          logging.DebugLevel,
				Provider:       loggingcfg.ProviderSlog,
				OutputFilepath: "/var/log/application/service.log",
			},
			Metrics: metricscfg.Config{
				Provider: tracingcfg.ProviderOtel,
				Otel: &otelgrpc.Config{
					ServiceName:        otelServiceName,
					CollectorEndpoint:  "localhost:4317",
					CollectionInterval: 1 * time.Second,
				},
			},
			Tracing: tracingcfg.Config{
				Provider:                  tracingcfg.ProviderOtel,
				ServiceName:               otelServiceName,
				SpanCollectionProbability: 1,
				Otel: &oteltrace.Config{
					CollectorEndpoint: "localhost:4317",
					Insecure:          true,
				},
			},
		},
		Services: config.ServicesConfig{
			Auth: authservice.Config{
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
				Tokens: tokenscfg.Config{
					Provider:                tokenscfg.ProviderPASETO,
					Audience:                "https://api.dinnerdonebetter.dev",
					Base64EncodedSigningKey: "",
				},
				TokenLifetime: 5 * time.Minute,
			},
			DataPrivacy: dataprivacyservice.Config{
				Uploads: uploads.Config{
					Storage: objectstorage.Config{
						GCPConfig:  &objectstorage.GCPConfig{BucketName: "userdata.dinnerdonebetter.dev"},
						BucketName: "userdata.dinnerdonebetter.dev",
						Provider:   objectstorage.GCPCloudStorageProvider,
					},
					Debug: false,
				},
			},
			Users: usersservice.Config{
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
			Recipes: recipemanagement.Config{
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
			ValidEnumerations: validenumerations.Config{},
			MealPlanning: mealplanningservice.Config{
				UseSearchService: true,
			},
		},
	}

	return cfg
}
