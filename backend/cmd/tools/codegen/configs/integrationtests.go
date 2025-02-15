package main

import (
	"encoding/base64"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/lib/analytics/config"
	tokenscfg "github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens/config"
	"github.com/dinnerdonebetter/backend/internal/lib/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/lib/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/oteltrace"
	"github.com/dinnerdonebetter/backend/internal/lib/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/lib/routing/config"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/server/http"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/objectstorage"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/dataprivacy"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/users"
	recipemanagement "github.com/dinnerdonebetter/backend/internal/services/eating/handlers/recipe_management"
)

func buildIntegrationTestsConfig() *config.APIServiceConfig {
	return &config.APIServiceConfig{
		Routing: routingcfg.Config{
			Provider: routingcfg.ProviderChi,
			Chi: &chi.Config{
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
					QueueAddresses: []string{dockerComposeWorkerQueueAddress},
				},
			},
			Publisher: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{dockerComposeWorkerQueueAddress},
				},
			},
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		HTTPServer: http.Config{
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
				ServiceName: otelServiceName,
				Level:       logging.InfoLevel,
				Provider:    loggingcfg.ProviderSlog,
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
		TextSearch: textsearchcfg.Config{
			// we're using a noop version of this in dev right now, but it still tries to instantiate a circuit breaker.
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "feature_flagger",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
		},
		FeatureFlags: featureflagscfg.Config{
			// we're using a noop version of this in dev right now, but it still tries to instantiate a circuit breaker.
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "feature_flagger",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
		},
		Analytics: analyticscfg.Config{
			// we're using a noop version of this in dev right now, but it still tries to instantiate a circuit breaker.
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "feature_flagger",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
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
				TokenLifetime:         5 * time.Minute,
				Tokens: tokenscfg.Config{
					Provider:                tokenscfg.ProviderPASETO,
					Audience:                "https://api.dinnerdonebetter.dev",
					Base64EncodedSigningKey: base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
				},
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
			Recipes: recipemanagement.Config{
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
