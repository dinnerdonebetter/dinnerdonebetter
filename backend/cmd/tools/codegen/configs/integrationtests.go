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
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing/oteltrace"
	"github.com/dinnerdonebetter/backend/internal/pkg/testutils"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/core/dataprivacy"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/core/users"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipes"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipesteps"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
)

func buildIntegrationTestsConfig() *config.APIServiceConfig {
	return &config.APIServiceConfig{
		Routing: routingcfg.Config{
			Provider: routingcfg.ProviderChi,
			ChiConfig: &chi.Config{
				ServiceName:            otelServiceName,
				EnableCORSForLocalhost: true,
				SilenceRouteLogging:    false,
			},
		},
		Meta: config.MetaSettings{
			Debug:   false,
			RunMode: testingEnv,
		},
		Queues: msgconfig.QueuesConfig{
			DataChangesTopicName:              dataChangesTopicName,
			OutboundEmailsTopicName:           outboundEmailsTopicName,
			SearchIndexRequestsTopicName:      searchIndexRequestsTopicName,
			UserDataAggregationTopicName:      userDataAggregationTopicName,
			WebhookExecutionRequestsTopicName: webhookExecutionRequestsTopicName,
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
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Server: http.Config{
			Debug:           false,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: databasecfg.Config{
			OAuth2TokenEncryptionKey: localOAuth2TokenEncryptionKey,
			Debug:                    true,
			RunMigrations:            true,
			LogQueries:               true,
			MaxPingAttempts:          maxAttempts,
			PingWaitPeriod:           1500 * time.Millisecond,
			ConnectionDetails:        localdevPostgresDBConnectionDetails,
		},
		Observability: observability.Config{
			Logging: loggingcfg.Config{
				Level:    logging.InfoLevel,
				Provider: loggingcfg.ProviderSlog,
			},
			Tracing: tracingcfg.Config{
				Provider:                  tracingcfg.ProviderOtel,
				SpanCollectionProbability: 1,
				ServiceName:               otelServiceName,
				Otel: &oteltrace.Config{
					CollectorEndpoint: "http://tracing-server:14268/api/traces",
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
				Debug:                 false,
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
					Debug: false,
					Storage: objectstorage.Config{
						Provider:   "memory",
						BucketName: "avatars",
						S3Config:   nil,
					},
				},
			},
			Recipes: recipesservice.Config{
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
		},
	}
}
