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
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing/config"
	"github.com/dinnerdonebetter/backend/internal/search/text/algolia"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/dataprivacy"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
)

func buildLocalDevConfig() *config.APIServiceConfig {
	return &config.APIServiceConfig{
		Routing: routingcfg.Config{
			Provider: routingcfg.ProviderChi,
			ChiConfig: &chi.Config{
				ServiceName:            otelServiceName,
				EnableCORSForLocalhost: true,
				SilenceRouteLogging:    false,
			},
		},
		Queues: msgconfig.QueuesConfig{
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
			ConnectionDetails:        localdevPostgresDBConnectionDetails,
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
			},
			Users: usersservice.Config{
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
			Recipes: recipesservice.Config{
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
		},
	}
}
