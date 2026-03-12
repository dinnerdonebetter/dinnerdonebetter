package main

import (
	"encoding/base64"
	"time"

	authcfg "github.com/dinnerdonebetter/backend/internal/authentication/config"
	tokenscfg "github.com/dinnerdonebetter/backend/internal/authentication/tokens/config"
	"github.com/dinnerdonebetter/backend/internal/branding"
	"github.com/dinnerdonebetter/backend/internal/config"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/platform/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/platform/analytics/segment"
	"github.com/dinnerdonebetter/backend/internal/platform/circuitbreaking"
	encryptioncfg "github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	emailcfg "github.com/dinnerdonebetter/backend/internal/platform/email/config"
	"github.com/dinnerdonebetter/backend/internal/platform/email/resend"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/platform/featureflags/config"
	"github.com/dinnerdonebetter/backend/internal/platform/featureflags/posthog"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue/pubsub"
	notificationscfg "github.com/dinnerdonebetter/backend/internal/platform/notifications/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	logotelgrpc "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/otelgrpc"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/otelgrpc"
	profilingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/profiling/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/profiling/pyroscope"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/oteltrace"
	"github.com/dinnerdonebetter/backend/internal/platform/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/platform/routing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/search/text/algolia"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/platform/server/grpc"
	"github.com/dinnerdonebetter/backend/internal/platform/server/http"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"
	uploadscfg "github.com/dinnerdonebetter/backend/internal/platform/uploads/config"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads/objectstorage"
	authservice "github.com/dinnerdonebetter/backend/internal/services/auth/handlers/authentication"
	dataprivacycfg "github.com/dinnerdonebetter/backend/internal/services/dataprivacy/config"
	identitycfg "github.com/dinnerdonebetter/backend/internal/services/identity/config"
	mealplanningcfg "github.com/dinnerdonebetter/backend/internal/services/mealplanning/config"
	oauthcfg "github.com/dinnerdonebetter/backend/internal/services/oauth/config"
	uploadedmediacfg "github.com/dinnerdonebetter/backend/internal/services/uploadedmedia/config"
)

const (
	prodGCPProject            = "dinner-done-better-prod"
	prodMediaBucket           = "media.dinnerdonebetter.com"
	prodUserDataBucket        = "userdata.dinnerdonebetter.com"
	prodOtelCollectorEndpoint = "otel-collector-svc.prod.svc.cluster.local:4317"
	prodOAuth2Domain          = "https://dinnerdonebetter.com"
	prodTokensAudience        = "https://http-api.dinnerdonebetter.com" //nolint:gosec // G101: audience URL, not a credential
	iosTeamID                 = "K8R2Q5UWQS"
	iosBundleID               = "com.dinnerdonebetter.ios"
)

func buildProdConfig() *config.APIServiceConfig {
	gcpMediaStorage := objectstorage.Config{
		Provider:          objectstorage.GCPCloudStorageProvider,
		BucketName:        prodMediaBucket,
		BucketPrefix:      "avatars/",
		UploadFilenameKey: "avatar",
		GCP: &objectstorage.GCPConfig{
			BucketName: prodMediaBucket,
		},
	}

	gcpUserDataStorage := objectstorage.Config{
		Provider:   objectstorage.GCPCloudStorageProvider,
		BucketName: prodUserDataBucket,
		GCP: &objectstorage.GCPConfig{
			BucketName: prodUserDataBucket,
		},
	}

	pubsubConfig := msgconfig.MessageQueueConfig{
		Provider: msgconfig.ProviderPubSub,
		PubSub: pubsub.Config{
			ProjectID: prodGCPProject,
		},
	}

	prodObservabilityConfig := observability.Config{
		Logging: loggingcfg.Config{
			ServiceName: otelServiceName,
			Level:       logging.InfoLevel,
			Provider:    loggingcfg.ProviderOtelSlog,
			OtelSlog: &logotelgrpc.Config{
				CollectorEndpoint: prodOtelCollectorEndpoint,
				Insecure:          true,
				Timeout:           2 * time.Second,
			},
		},
		Metrics: metricscfg.Config{
			ServiceName: otelServiceName,
			Otel: &otelgrpc.Config{
				Insecure:             true,
				CollectorEndpoint:    prodOtelCollectorEndpoint,
				CollectionInterval:   30 * time.Second,
				EnableRuntimeMetrics: true,
				EnableHostMetrics:    true,
			},
			Provider: metricscfg.ProviderOtel,
		},
		Tracing: tracingcfg.Config{
			Provider:                  tracingcfg.ProviderOtel,
			ServiceName:               otelServiceName,
			SpanCollectionProbability: 1.0,
			Otel: &oteltrace.Config{
				Insecure:          true,
				CollectorEndpoint: prodOtelCollectorEndpoint,
			},
		},
		Profiling: profilingcfg.Config{
			ServiceName: otelServiceName,
			Provider:    profilingcfg.ProviderPyroscope,
			Pyroscope: &pyroscope.Config{
				ServerAddress: "https://profiles-prod-001.grafana.net",
				UploadRate:    15 * time.Second,
			},
		},
	}

	return &config.APIServiceConfig{
		Routing: routingcfg.Config{
			Provider: routingcfg.ProviderChi,
			Chi: &chi.Config{
				ServiceName:            otelServiceName,
				EnableCORSForLocalhost: false,
				SilenceRouteLogging:    false,
			},
		},
		Queues: msgconfig.QueuesConfig{
			DataChangesTopicName:              dataChangesTopicName,
			OutboundEmailsTopicName:           outboundEmailsTopicName,
			SearchIndexRequestsTopicName:      searchIndexRequestsTopicName,
			MobileNotificationsTopicName:      mobileNotificationsTopicName,
			UserDataAggregationTopicName:      userDataAggregationTopicName,
			WebhookExecutionRequestsTopicName: webhookExecutionRequestsTopicName,
		},
		Meta: config.MetaSettings{
			Debug:   false,
			RunMode: "production",
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Events: msgconfig.Config{
			Consumer:  pubsubConfig,
			Publisher: pubsubConfig,
		},
		GRPCServer: grpc.Config{
			Port: defaultGRPCPort,
		},
		HTTPServer: http.Config{
			Debug:           false,
			Port:            defaultHTTPPort,
			StartupDeadline: 60 * time.Second,
		},
		Database: databasecfg.Config{
			Provider:                     databasecfg.ProviderPostgres,
			Encryption:                   encryptioncfg.Config{Provider: encryptioncfg.ProviderSalsa20},
			OAuth2TokenEncryptionKey:     "",
			UserDeviceTokenEncryptionKey: "",
			Debug:                        false,
			RunMigrations:                true,
			LogQueries:                   false,
			MaxPingAttempts:              maxAttempts,
			PingWaitPeriod:               time.Second,
			ReadConnection: databasecfg.ConnectionDetails{
				Username:   "api_db_user",
				Database:   "dinner-done-better",
				Port:       5432,
				DisableSSL: false,
			},
			WriteConnection: databasecfg.ConnectionDetails{
				Username:   "api_db_user",
				Database:   "dinner-done-better",
				Port:       5432,
				DisableSSL: false,
			},
		},
		Observability: prodObservabilityConfig,
		Email: emailcfg.Config{
			Provider: emailcfg.ProviderResend,
			Resend: &resend.Config{
				APIToken: "placeholder", // overridden by env from api-service-config secret
			},
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "prod_emailer",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
		},
		Analytics: analyticscfg.Config{
			Provider: analyticscfg.ProviderSegment,
			Segment:  &segment.Config{APIToken: "placeholder"}, // overridden by env from CSI secret
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "prod_analytics",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
		},
		TextSearch: textsearchcfg.Config{
			Provider: textsearchcfg.AlgoliaProvider,
			Algolia:  &algolia.Config{},
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "prod_text_searcher",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
		},
		FeatureFlags: featureflagscfg.Config{
			Provider: featureflagscfg.ProviderPostHog,
			PostHog:  &posthog.Config{ProjectAPIKey: "placeholder"}, // overridden by env from CSI secret
			CircuitBreaker: circuitbreaking.Config{
				Name:                   "feature_flagger",
				ErrorRate:              .5,
				MinimumSampleThreshold: 100,
			},
		},
		Auth: authcfg.Config{
			SSO: authcfg.SSOConfigs{Google: authcfg.GoogleSSOConfig{}},
			Passkey: authcfg.PasskeyConfig{
				RPID:          "dinnerdonebetter.com",
				RPDisplayName: branding.CompanyName,
				RPOrigins:     []string{"https://dinnerdonebetter.com", "https://www.dinnerdonebetter.com", "https://admin.dinnerdonebetter.com"},
			},
			Tokens: tokenscfg.Config{
				Provider:                tokenscfg.ProviderPASETO,
				Audience:                prodTokensAudience,
				Base64EncodedSigningKey: base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
			},
			Debug:                 false,
			EnableUserSignup:      true,
			MinimumUsernameLength: 3,
			MinimumPasswordLength: 8,
		},
		Services: config.ServicesConfig{
			Auth: authservice.Config{
				OAuth2: authservice.OAuth2Config{
					Domain:               prodOAuth2Domain,
					AccessTokenLifespan:  time.Hour,
					RefreshTokenLifespan: time.Hour,
					Debug:                false,
				},
				SSO: authservice.SSOConfigs{
					Google: authservice.GoogleSSOConfig{},
				},
				Debug:                 false,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
				TokenLifetime:         5 * time.Minute,
				Tokens: tokenscfg.Config{
					Provider:                tokenscfg.ProviderPASETO,
					Audience:                prodTokensAudience,
					Base64EncodedSigningKey: base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
				},
			},
			DataPrivacy: dataprivacycfg.Config{
				Uploads: uploadscfg.Config{
					Storage: gcpUserDataStorage,
					Debug:   false,
				},
			},
			Users: identitycfg.Config{
				PublicMediaURLPrefix: "https://" + prodMediaBucket + "/avatars",
				Uploads: uploadscfg.Config{
					Storage: gcpMediaStorage,
					Debug:   false,
				},
			},
			UploadedMedia: uploadedmediacfg.Config{
				Uploads: uploadscfg.Config{
					Storage: gcpMediaStorage,
					Debug:   false,
				},
			},
			MealPlanning: mealplanningcfg.Config{
				UseSearchService: true,
			},
			OAuth2Clients: oauthcfg.Config{
				OAuth2ClientCreationDisabled: true,
			},
		},
		PushNotifications: notificationscfg.Config{
			Provider: notificationscfg.ProviderAPNsFCM,
			APNs: &notificationscfg.APNsConfig{
				AuthKeyPath: "/mnt/apns/apns-auth-key.p8", // mounted from K8s secret apns-credentials
				TeamID:      iosTeamID,
				BundleID:    iosBundleID,
				Production:  true,
			},
			FCM: &notificationscfg.FCMConfig{
				// CredentialsPath empty: uses Application Default Credentials (GCP workload identity)
			},
		},
	}
}
