package main

import (
	"encoding/base64"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/lib/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens/config"
	"github.com/dinnerdonebetter/backend/internal/lib/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	logotelgrpc "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/otelgrpc"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/otelgrpc"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/oteltrace"
	"github.com/dinnerdonebetter/backend/internal/lib/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/lib/routing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/algolia"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/server/http"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/objectstorage"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/core/dataprivacy"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/core/users"
	recipemanagement "github.com/dinnerdonebetter/backend/internal/services/eating/recipe_management"
	"github.com/dinnerdonebetter/backend/internal/testutils"
)

const (
	dockerComposeWorkerQueueAddress = "worker_queue:6379"
	localOAuth2TokenEncryptionKey   = debugCookieHashKey
)

var (
	localdevPostgresDBConnectionDetails = databasecfg.ConnectionDetails{
		Username:   "dbuser",
		Password:   "hunter2",
		Database:   "dinner-done-better",
		Host:       "pgdatabase",
		Port:       5432,
		DisableSSL: true,
	}
)

func buildLocalDevConfig() *config.APIServiceConfig {
	return &config.APIServiceConfig{
		Routing: routingcfg.Config{
			Provider: routingcfg.ProviderChi,
			Chi: &chi.Config{
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
		FeatureFlags: featureflagscfg.Config{
			// we're using a noop version of this in localdev right now, but it still tries to instantiate a circuit breaker.
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "feature_flagger",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
		},
		Analytics: analyticscfg.Config{
			// we're using a noop version of this in localdev right now, but it still tries to instantiate a circuit breaker.
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "feature_flagger",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
		},
		TextSearch: textsearchcfg.Config{
			Algolia:  &algolia.Config{},
			Provider: textsearchcfg.AlgoliaProvider,
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "dev_text_searcher",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
		},
		HTTPServer: http.Config{
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
				ServiceName: otelServiceName,
				Level:       logging.DebugLevel,
				Provider:    loggingcfg.ProviderOtelSlog,
				OtelSlog: &logotelgrpc.Config{
					CollectorEndpoint: "otel_collector:4317",
					Insecure:          true,
					Timeout:           time.Second * 3,
				},
			},
			Metrics: metricscfg.Config{
				ServiceName: otelServiceName,
				Otel: &otelgrpc.Config{
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
			Recipes: recipemanagement.Config{
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
