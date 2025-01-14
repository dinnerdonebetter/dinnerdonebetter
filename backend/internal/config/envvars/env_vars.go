package envvars

/*
This file contains a reference of all valid service environment variables.
*/

const (
	// AnalyticsCircuitBreakerErrorRateEnvVarKey is the environment variable name to set in order to override `config.Analytics.CircuitBreakerConfig.ErrorRate`.
	AnalyticsCircuitBreakerErrorRateEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_CIRCUIT_BREAKER_ERROR_RATE"

	// AnalyticsCircuitBreakerMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.Analytics.CircuitBreakerConfig.MinimumSampleThreshold`.
	AnalyticsCircuitBreakerMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_CIRCUIT_BREAKER_MINIMUM_SAMPLE_THRESHOLD"

	// AnalyticsPosthogCircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.Analytics.Posthog.CircuitBreakerConfig.ErrorRate`.
	AnalyticsPosthogCircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_POSTHOG_CIRCUIT_BREAKING_ERROR_RATE"

	// AnalyticsPosthogCircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.Analytics.Posthog.CircuitBreakerConfig.MinimumSampleThreshold`.
	AnalyticsPosthogCircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_POSTHOG_CIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// AnalyticsPosthogPersonalAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Analytics.Posthog.PersonalAPIKey`.
	AnalyticsPosthogPersonalAPIKeyEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_POSTHOG_PERSONAL_API_KEY"

	// AnalyticsPosthogProjectAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Analytics.Posthog.ProjectAPIKey`.
	AnalyticsPosthogProjectAPIKeyEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_POSTHOG_PROJECT_API_KEY"

	// AnalyticsProviderEnvVarKey is the environment variable name to set in order to override `config.Analytics.Provider`.
	AnalyticsProviderEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_PROVIDER"

	// AnalyticsRudderstackAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Analytics.Rudderstack.APIKey`.
	AnalyticsRudderstackAPIKeyEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_RUDDERSTACK_API_KEY"

	// AnalyticsRudderstackDataPlaneURLEnvVarKey is the environment variable name to set in order to override `config.Analytics.Rudderstack.DataPlaneURL`.
	AnalyticsRudderstackDataPlaneURLEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_RUDDERSTACK_DATA_PLANE_URL"

	// AnalyticsSegmentAPITokenEnvVarKey is the environment variable name to set in order to override `config.Analytics.Segment.APIToken`.
	AnalyticsSegmentAPITokenEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_SEGMENT_API_TOKEN"

	// DatabaseConnectionDetailsDatabaseEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Database`.
	DatabaseConnectionDetailsDatabaseEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_DATABASE"

	// DatabaseConnectionDetailsDisableSslEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.DisableSSL`.
	DatabaseConnectionDetailsDisableSslEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_DISABLE_SSL"

	// DatabaseConnectionDetailsHostEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Host`.
	DatabaseConnectionDetailsHostEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_HOST"

	// DatabaseConnectionDetailsPasswordUnsetEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Password`.
	DatabaseConnectionDetailsPasswordUnsetEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_PASSWORD,unset"

	// DatabaseConnectionDetailsPortEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Port`.
	DatabaseConnectionDetailsPortEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_PORT"

	// DatabaseConnectionDetailsUsernameEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Username`.
	DatabaseConnectionDetailsUsernameEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_USERNAME"

	// DatabaseDebugEnvVarKey is the environment variable name to set in order to override `config.Database.Debug`.
	DatabaseDebugEnvVarKey = "DINNER_DONE_BETTER_DATABASE_DEBUG"

	// DatabaseLogQueriesEnvVarKey is the environment variable name to set in order to override `config.Database.LogQueries`.
	DatabaseLogQueriesEnvVarKey = "DINNER_DONE_BETTER_DATABASE_LOG_QUERIES"

	// DatabaseMaxPingAttemptsEnvVarKey is the environment variable name to set in order to override `config.Database.MaxPingAttempts`.
	DatabaseMaxPingAttemptsEnvVarKey = "DINNER_DONE_BETTER_DATABASE_MAX_PING_ATTEMPTS"

	// DatabaseOauth2TokenEncryptionKeyEnvVarKey is the environment variable name to set in order to override `config.Database.OAuth2TokenEncryptionKey`.
	DatabaseOauth2TokenEncryptionKeyEnvVarKey = "DINNER_DONE_BETTER_DATABASE_OAUTH2_TOKEN_ENCRYPTION_KEY"

	// DatabasePingWaitPeriodEnvVarKey is the environment variable name to set in order to override `config.Database.PingWaitPeriod`.
	DatabasePingWaitPeriodEnvVarKey = "DINNER_DONE_BETTER_DATABASE_PING_WAIT_PERIOD"

	// DatabaseProviderEnvVarKey is the environment variable name to set in order to override `config.Database.Provider`.
	DatabaseProviderEnvVarKey = "DINNER_DONE_BETTER_DATABASE_PROVIDER"

	// DatabaseRunMigrationsEnvVarKey is the environment variable name to set in order to override `config.Database.RunMigrations`.
	DatabaseRunMigrationsEnvVarKey = "DINNER_DONE_BETTER_DATABASE_RUN_MIGRATIONS"

	// EmailCircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.Email.CircuitBreakerConfig.ErrorRate`.
	EmailCircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_EMAIL_CIRCUIT_BREAKING_ERROR_RATE"

	// EmailCircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.Email.CircuitBreakerConfig.MinimumSampleThreshold`.
	EmailCircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_EMAIL_CIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// EmailMailgunDomainEnvVarKey is the environment variable name to set in order to override `config.Email.Mailgun.Domain`.
	EmailMailgunDomainEnvVarKey = "DINNER_DONE_BETTER_EMAIL_MAILGUN_DOMAIN"

	// EmailMailgunPrivateAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Email.Mailgun.PrivateAPIKey`.
	EmailMailgunPrivateAPIKeyEnvVarKey = "DINNER_DONE_BETTER_EMAIL_MAILGUN_PRIVATE_API_KEY"

	// EmailMailjetAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Email.Mailjet.APIKey`.
	EmailMailjetAPIKeyEnvVarKey = "DINNER_DONE_BETTER_EMAIL_MAILJET_API_KEY"

	// EmailMailjetSecretKeyEnvVarKey is the environment variable name to set in order to override `config.Email.Mailjet.SecretKey`.
	EmailMailjetSecretKeyEnvVarKey = "DINNER_DONE_BETTER_EMAIL_MAILJET_SECRET_KEY"

	// EmailProviderEnvVarKey is the environment variable name to set in order to override `config.Email.Provider`.
	EmailProviderEnvVarKey = "DINNER_DONE_BETTER_EMAIL_PROVIDER"

	// EmailSendgridAPITokenEnvVarKey is the environment variable name to set in order to override `config.Email.Sendgrid.APIToken`.
	EmailSendgridAPITokenEnvVarKey = "DINNER_DONE_BETTER_EMAIL_SENDGRID_API_TOKEN"

	// EncodingContentTypeEnvVarKey is the environment variable name to set in order to override `config.Encoding.ContentType`.
	EncodingContentTypeEnvVarKey = "DINNER_DONE_BETTER_ENCODING_CONTENT_TYPE"

	// EventsConsumerProviderEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.Provider`.
	EventsConsumerProviderEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_PROVIDER"

	// EventsConsumerPubsubProjectIDEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.PubSub.ProjectID`.
	EventsConsumerPubsubProjectIDEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_PUBSUB_PROJECT_ID"

	// EventsConsumerRedisPasswordEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.Redis.Password`.
	EventsConsumerRedisPasswordEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_REDIS_PASSWORD"

	// EventsConsumerRedisQueueAddressesEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.Redis.QueueAddresses`.
	EventsConsumerRedisQueueAddressesEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_REDIS_QUEUE_ADDRESSES"

	// EventsConsumerRedisUsernameEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.Redis.Username`.
	EventsConsumerRedisUsernameEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_REDIS_USERNAME"

	// EventsConsumerSqsQueueAddressEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.SQS.QueueAddress`.
	EventsConsumerSqsQueueAddressEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_SQS_QUEUE_ADDRESS"

	// EventsPublisherProviderEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.Provider`.
	EventsPublisherProviderEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_PROVIDER"

	// EventsPublisherPubsubProjectIDEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.PubSub.ProjectID`.
	EventsPublisherPubsubProjectIDEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_PUBSUB_PROJECT_ID"

	// EventsPublisherRedisPasswordEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.Redis.Password`.
	EventsPublisherRedisPasswordEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_REDIS_PASSWORD"

	// EventsPublisherRedisQueueAddressesEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.Redis.QueueAddresses`.
	EventsPublisherRedisQueueAddressesEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_REDIS_QUEUE_ADDRESSES"

	// EventsPublisherRedisUsernameEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.Redis.Username`.
	EventsPublisherRedisUsernameEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_REDIS_USERNAME"

	// EventsPublisherSqsQueueAddressEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.SQS.QueueAddress`.
	EventsPublisherSqsQueueAddressEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_SQS_QUEUE_ADDRESS"

	// FeatureFlagsCircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.CircuitBreakingConfig.ErrorRate`.
	FeatureFlagsCircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_CIRCUIT_BREAKING_ERROR_RATE"

	// FeatureFlagsCircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.CircuitBreakingConfig.MinimumSampleThreshold`.
	FeatureFlagsCircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_CIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// FeatureFlagsLaunchDarklycircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.LaunchDarkly.CircuitBreakerConfig.ErrorRate`.
	FeatureFlagsLaunchDarklycircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_LAUNCH_DARKLYCIRCUIT_BREAKING_ERROR_RATE"

	// FeatureFlagsLaunchDarklycircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.LaunchDarkly.CircuitBreakerConfig.MinimumSampleThreshold`.
	FeatureFlagsLaunchDarklycircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_LAUNCH_DARKLYCIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// FeatureFlagsLaunchDarklyinitTimeoutEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.LaunchDarkly.InitTimeout`.
	FeatureFlagsLaunchDarklyinitTimeoutEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_LAUNCH_DARKLYINIT_TIMEOUT"

	// FeatureFlagsLaunchDarklysdkKeyEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.LaunchDarkly.SDKKey`.
	FeatureFlagsLaunchDarklysdkKeyEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_LAUNCH_DARKLYSDK_KEY"

	// FeatureFlagsPosthogCircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.PostHog.CircuitBreakerConfig.ErrorRate`.
	FeatureFlagsPosthogCircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_POSTHOG_CIRCUIT_BREAKING_ERROR_RATE"

	// FeatureFlagsPosthogCircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.PostHog.CircuitBreakerConfig.MinimumSampleThreshold`.
	FeatureFlagsPosthogCircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_POSTHOG_CIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// FeatureFlagsPosthogPersonalAPIKeyEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.PostHog.PersonalAPIKey`.
	FeatureFlagsPosthogPersonalAPIKeyEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_POSTHOG_PERSONAL_API_KEY"

	// FeatureFlagsPosthogProjectAPIKeyEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.PostHog.ProjectAPIKey`.
	FeatureFlagsPosthogProjectAPIKeyEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_POSTHOG_PROJECT_API_KEY"

	// FeatureFlagsProviderEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.Provider`.
	FeatureFlagsProviderEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_PROVIDER"

	// MetaDebugEnvVarKey is the environment variable name to set in order to override `config.Meta.Debug`.
	MetaDebugEnvVarKey = "DINNER_DONE_BETTER_META_DEBUG"

	// MetaRunModeEnvVarKey is the environment variable name to set in order to override `config.Meta.RunMode`.
	MetaRunModeEnvVarKey = "DINNER_DONE_BETTER_META_RUN_MODE"

	// ObservabilityLoggingLevelEnvVarKey is the environment variable name to set in order to override `config.Observability.Logging.Level`.
	ObservabilityLoggingLevelEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_LOGGING_LEVEL"

	// ObservabilityLoggingOutputFilepathEnvVarKey is the environment variable name to set in order to override `config.Observability.Logging.OutputFilepath`.
	ObservabilityLoggingOutputFilepathEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_LOGGING_OUTPUT_FILEPATH"

	// ObservabilityLoggingProviderEnvVarKey is the environment variable name to set in order to override `config.Observability.Logging.Provider`.
	ObservabilityLoggingProviderEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_LOGGING_PROVIDER"

	// ObservabilityMetricsOtelCollectionIntervalEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.CollectionInterval`.
	ObservabilityMetricsOtelCollectionIntervalEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_COLLECTION_INTERVAL"

	// ObservabilityMetricsOtelCollectionTimeoutEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.CollectionTimeout`.
	ObservabilityMetricsOtelCollectionTimeoutEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_COLLECTION_TIMEOUT"

	// ObservabilityMetricsOtelCollectorEndpointEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.CollectorEndpoint`.
	ObservabilityMetricsOtelCollectorEndpointEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_COLLECTOR_ENDPOINT"

	// ObservabilityMetricsOtelInsecureEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.Insecure`.
	ObservabilityMetricsOtelInsecureEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_INSECURE"

	// ObservabilityMetricsOtelServiceNameEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.ServiceName`.
	ObservabilityMetricsOtelServiceNameEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_SERVICE_NAME"

	// ObservabilityMetricsProviderEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Provider`.
	ObservabilityMetricsProviderEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_PROVIDER"

	// ObservabilityTracingCloudtraceGoogleCloudTraceProjectIDEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.CloudTrace.ProjectID`.
	ObservabilityTracingCloudtraceGoogleCloudTraceProjectIDEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_CLOUDTRACE_GOOGLE_CLOUD_TRACE_PROJECT_ID"

	// ObservabilityTracingOtelgrpcCollectorEndpointEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.Otel.CollectorEndpoint`.
	ObservabilityTracingOtelgrpcCollectorEndpointEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_OTELGRPC_COLLECTOR_ENDPOINT"

	// ObservabilityTracingOtelgrpcInsecureEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.Otel.Insecure`.
	ObservabilityTracingOtelgrpcInsecureEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_OTELGRPC_INSECURE"

	// ObservabilityTracingTracingProviderEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.Provider`.
	ObservabilityTracingTracingProviderEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_TRACING_PROVIDER"

	// ObservabilityTracingTracingServiceNameEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.ServiceName`.
	ObservabilityTracingTracingServiceNameEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_TRACING_SERVICE_NAME"

	// ObservabilityTracingTracingSpanCollectionProbabilityEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.SpanCollectionProbability`.
	ObservabilityTracingTracingSpanCollectionProbabilityEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_TRACING_SPAN_COLLECTION_PROBABILITY"

	// QueuesDataChangesTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.DataChangesTopicName`.
	QueuesDataChangesTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_DATA_CHANGES_TOPIC_NAME"

	// QueuesOutboundEmailsTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.OutboundEmailsTopicName`.
	QueuesOutboundEmailsTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_OUTBOUND_EMAILS_TOPIC_NAME"

	// QueuesSearchIndexRequestsTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.SearchIndexRequestsTopicName`.
	QueuesSearchIndexRequestsTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_SEARCH_INDEX_REQUESTS_TOPIC_NAME"

	// QueuesUserDataAggregationTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.UserDataAggregationTopicName`.
	QueuesUserDataAggregationTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_USER_DATA_AGGREGATION_TOPIC_NAME"

	// QueuesWebhookExecutionRequestsTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.WebhookExecutionRequestsTopicName`.
	QueuesWebhookExecutionRequestsTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_WEBHOOK_EXECUTION_REQUESTS_TOPIC_NAME"

	// RoutingChiEnableCorsForLocalhostEnvVarKey is the environment variable name to set in order to override `config.Routing.ChiConfig.EnableCORSForLocalhost`.
	RoutingChiEnableCorsForLocalhostEnvVarKey = "DINNER_DONE_BETTER_ROUTING_CHI_ENABLE_CORS_FOR_LOCALHOST"

	// RoutingChiServiceNameEnvVarKey is the environment variable name to set in order to override `config.Routing.ChiConfig.ServiceName`.
	RoutingChiServiceNameEnvVarKey = "DINNER_DONE_BETTER_ROUTING_CHI_SERVICE_NAME"

	// RoutingChiSilenceRouteLoggingEnvVarKey is the environment variable name to set in order to override `config.Routing.ChiConfig.SilenceRouteLogging`.
	RoutingChiSilenceRouteLoggingEnvVarKey = "DINNER_DONE_BETTER_ROUTING_CHI_SILENCE_ROUTE_LOGGING"

	// RoutingChiValidDomainsEnvVarKey is the environment variable name to set in order to override `config.Routing.ChiConfig.ValidDomains`.
	RoutingChiValidDomainsEnvVarKey = "DINNER_DONE_BETTER_ROUTING_CHI_VALID_DOMAINS"

	// RoutingProviderEnvVarKey is the environment variable name to set in order to override `config.Routing.Provider`.
	RoutingProviderEnvVarKey = "DINNER_DONE_BETTER_ROUTING_PROVIDER"

	// SearchAlgoliaAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Search.Algolia.APIKey`.
	SearchAlgoliaAPIKeyEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ALGOLIA_API_KEY"

	// SearchAlgoliaAppIDEnvVarKey is the environment variable name to set in order to override `config.Search.Algolia.AppID`.
	SearchAlgoliaAppIDEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ALGOLIA_APP_ID"

	// SearchAlgoliaTimeoutEnvVarKey is the environment variable name to set in order to override `config.Search.Algolia.Timeout`.
	SearchAlgoliaTimeoutEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ALGOLIA_TIMEOUT"

	// SearchCircuitBreakerErrorRateEnvVarKey is the environment variable name to set in order to override `config.Search.CircuitBreakerConfig.ErrorRate`.
	SearchCircuitBreakerErrorRateEnvVarKey = "DINNER_DONE_BETTER_SEARCH_CIRCUIT_BREAKER_ERROR_RATE"

	// SearchCircuitBreakerMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.Search.CircuitBreakerConfig.MinimumSampleThreshold`.
	SearchCircuitBreakerMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_SEARCH_CIRCUIT_BREAKER_MINIMUM_SAMPLE_THRESHOLD"

	// SearchElasticsearchAddressEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.Address`.
	SearchElasticsearchAddressEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_ADDRESS"

	// SearchElasticsearchCaCertEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.CACert`.
	SearchElasticsearchCaCertEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_CA_CERT"

	// SearchElasticsearchIndexOperationTimeoutEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.IndexOperationTimeout`.
	SearchElasticsearchIndexOperationTimeoutEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_INDEX_OPERATION_TIMEOUT"

	// SearchElasticsearchPasswordEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.Password`.
	SearchElasticsearchPasswordEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_PASSWORD"

	// SearchElasticsearchUsernameEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.Username`.
	SearchElasticsearchUsernameEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_USERNAME"

	// SearchProviderEnvVarKey is the environment variable name to set in order to override `config.Search.Provider`.
	SearchProviderEnvVarKey = "DINNER_DONE_BETTER_SEARCH_PROVIDER"

	// ServerAutocertAddressesEnvVarKey is the environment variable name to set in order to override `config.Server.Autocert.Addresses`.
	ServerAutocertAddressesEnvVarKey = "DINNER_DONE_BETTER_SERVER_AUTOCERT_ADDRESSES"

	// ServerDebugEnvVarKey is the environment variable name to set in order to override `config.Server.Debug`.
	ServerDebugEnvVarKey = "DINNER_DONE_BETTER_SERVER_DEBUG"

	// ServerHTTPPortEnvVarKey is the environment variable name to set in order to override `config.Server.HTTPPort`.
	ServerHTTPPortEnvVarKey = "DINNER_DONE_BETTER_SERVER_HTTP_PORT"

	// ServerHTTPSCertificateFilepathEnvVarKey is the environment variable name to set in order to override `config.Server.HTTPSCertificateFile`.
	ServerHTTPSCertificateFilepathEnvVarKey = "DINNER_DONE_BETTER_SERVER_HTTPS_CERTIFICATE_FILEPATH"

	// ServerHTTPSCertificateKeyFilepathEnvVarKey is the environment variable name to set in order to override `config.Server.HTTPSCertificateKeyFile`.
	ServerHTTPSCertificateKeyFilepathEnvVarKey = "DINNER_DONE_BETTER_SERVER_HTTPS_CERTIFICATE_KEY_FILEPATH"

	// ServerStartupDeadlineEnvVarKey is the environment variable name to set in order to override `config.Server.StartupDeadline`.
	ServerStartupDeadlineEnvVarKey = "DINNER_DONE_BETTER_SERVER_STARTUP_DEADLINE"

	// ServiceAuthDebugEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.Debug`.
	ServiceAuthDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_DEBUG"

	// ServiceAuthEnableUserSignupEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.EnableUserSignup`.
	ServiceAuthEnableUserSignupEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_ENABLE_USER_SIGNUP"

	// ServiceAuthJwtLifetimeEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.TokenLifetime`.
	ServiceAuthJwtLifetimeEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_JWT_LIFETIME"

	// ServiceAuthJwtSigningKeyEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.JWTSigningKey`.
	ServiceAuthJwtSigningKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_JWT_SIGNING_KEY"

	// ServiceAuthMinimumPasswordLengthEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.MinimumPasswordLength`.
	ServiceAuthMinimumPasswordLengthEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_MINIMUM_PASSWORD_LENGTH"

	// ServiceAuthMinimumUsernameLengthEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.MinimumUsernameLength`.
	ServiceAuthMinimumUsernameLengthEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_MINIMUM_USERNAME_LENGTH"

	// ServiceAuthOauth2AccessTokenLifespanEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.OAuth2.AccessTokenLifespan`.
	ServiceAuthOauth2AccessTokenLifespanEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_OAUTH2ACCESS_TOKEN_LIFESPAN"

	// ServiceAuthOauth2DebugEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.OAuth2.Debug`.
	ServiceAuthOauth2DebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_OAUTH2DEBUG"

	// ServiceAuthOauth2DomainEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.OAuth2.Domain`.
	ServiceAuthOauth2DomainEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_OAUTH2DOMAIN"

	// ServiceAuthOauth2RefreshTokenLifespanEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.OAuth2.RefreshTokenLifespan`.
	ServiceAuthOauth2RefreshTokenLifespanEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_OAUTH2REFRESH_TOKEN_LIFESPAN"

	// ServiceAuthSsoConfigGoogleCallbackURLEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.SSO.Google.CallbackURL`.
	ServiceAuthSsoConfigGoogleCallbackURLEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_SSO_CONFIG_GOOGLE_CALLBACK_URL"

	// ServiceAuthSsoConfigGoogleClientIDEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.SSO.Google.ClientID`.
	ServiceAuthSsoConfigGoogleClientIDEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_SSO_CONFIG_GOOGLE_CLIENT_ID"

	// ServiceAuthSsoConfigGoogleClientSecretEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.SSO.Google.ClientSecret`.
	ServiceAuthSsoConfigGoogleClientSecretEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_SSO_CONFIG_GOOGLE_CLIENT_SECRET"

	// ServiceAuthTokenAudienceEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.Tokens.Audience`.
	ServiceAuthTokenAudienceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_TOKEN_AUDIENCE"

	// ServiceAuthTokenProviderEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.Tokens.Provider`.
	ServiceAuthTokenProviderEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_TOKEN_PROVIDER"

	// ServiceAuthTokenSigningKeyEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.Tokens.Base64EncodedSigningKey`.
	ServiceAuthTokenSigningKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_TOKEN_SIGNING_KEY"

	// ServiceDataPrivacyUploadsDebugEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Debug`.
	ServiceDataPrivacyUploadsDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_DEBUG"

	// ServiceDataPrivacyUploadsStorageBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.BucketName`.
	ServiceDataPrivacyUploadsStorageBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_BUCKET_NAME"

	// ServiceDataPrivacyUploadsStorageBucketPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.BucketPrefix`.
	ServiceDataPrivacyUploadsStorageBucketPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_BUCKET_PREFIX"

	// ServiceDataPrivacyUploadsStorageFilesystemRootDirectoryEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.FilesystemConfig.RootDirectory`.
	ServiceDataPrivacyUploadsStorageFilesystemRootDirectoryEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_FILESYSTEM_ROOT_DIRECTORY"

	// ServiceDataPrivacyUploadsStorageGcpBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.GCPConfig.BucketName`.
	ServiceDataPrivacyUploadsStorageGcpBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_GCP_BUCKET_NAME"

	// ServiceDataPrivacyUploadsStorageProviderEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.Provider`.
	ServiceDataPrivacyUploadsStorageProviderEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_PROVIDER"

	// ServiceDataPrivacyUploadsStorageS3BucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.S3Config.BucketName`.
	ServiceDataPrivacyUploadsStorageS3BucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_S3_BUCKET_NAME"

	// ServiceDataPrivacyUploadsStorageUploadFilenameKeyEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.UploadFilenameKey`.
	ServiceDataPrivacyUploadsStorageUploadFilenameKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_UPLOAD_FILENAME_KEY"

	// ServiceHouseholdInvitationsDebugEnvVarKey is the environment variable name to set in order to override `config.Services.HouseholdInvitations.Debug`.
	ServiceHouseholdInvitationsDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_HOUSEHOLD_INVITATIONS_DEBUG"

	// ServiceMealPlanningUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.MealPlanning.UseSearchService`.
	ServiceMealPlanningUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_MEAL_PLANNING_USE_SEARCH_SERVICE"

	// ServiceOauth2ClientsCreationDisabledEnvVarKey is the environment variable name to set in order to override `config.Services.OAuth2Clients.OAuth2ClientCreationDisabled`.
	ServiceOauth2ClientsCreationDisabledEnvVarKey = "DINNER_DONE_BETTER_SERVICE_OAUTH2_CLIENTS_CREATION_DISABLED"

	// ServiceRecipesPublicMediaURLPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.PublicMediaURLPrefix`.
	ServiceRecipesPublicMediaURLPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_PUBLIC_MEDIA_URL_PREFIX"

	// ServiceRecipesUploadsDebugEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Debug`.
	ServiceRecipesUploadsDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_DEBUG"

	// ServiceRecipesUploadsStorageBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.BucketName`.
	ServiceRecipesUploadsStorageBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_BUCKET_NAME"

	// ServiceRecipesUploadsStorageBucketPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.BucketPrefix`.
	ServiceRecipesUploadsStorageBucketPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_BUCKET_PREFIX"

	// ServiceRecipesUploadsStorageFilesystemRootDirectoryEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.FilesystemConfig.RootDirectory`.
	ServiceRecipesUploadsStorageFilesystemRootDirectoryEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_FILESYSTEM_ROOT_DIRECTORY"

	// ServiceRecipesUploadsStorageGcpBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.GCPConfig.BucketName`.
	ServiceRecipesUploadsStorageGcpBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_GCP_BUCKET_NAME"

	// ServiceRecipesUploadsStorageProviderEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.Provider`.
	ServiceRecipesUploadsStorageProviderEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_PROVIDER"

	// ServiceRecipesUploadsStorageS3BucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.S3Config.BucketName`.
	ServiceRecipesUploadsStorageS3BucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_S3_BUCKET_NAME"

	// ServiceRecipesUploadsStorageUploadFilenameKeyEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.UploadFilenameKey`.
	ServiceRecipesUploadsStorageUploadFilenameKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_UPLOAD_FILENAME_KEY"

	// ServiceRecipesUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.UseSearchService`.
	ServiceRecipesUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_USE_SEARCH_SERVICE"

	// ServiceUsersPublicMediaURLPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.Users.PublicMediaURLPrefix`.
	ServiceUsersPublicMediaURLPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_PUBLIC_MEDIA_URL_PREFIX"

	// ServiceUsersUploadsDebugEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Debug`.
	ServiceUsersUploadsDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_DEBUG"

	// ServiceUsersUploadsStorageBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.BucketName`.
	ServiceUsersUploadsStorageBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_BUCKET_NAME"

	// ServiceUsersUploadsStorageBucketPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.BucketPrefix`.
	ServiceUsersUploadsStorageBucketPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_BUCKET_PREFIX"

	// ServiceUsersUploadsStorageFilesystemRootDirectoryEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.FilesystemConfig.RootDirectory`.
	ServiceUsersUploadsStorageFilesystemRootDirectoryEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_FILESYSTEM_ROOT_DIRECTORY"

	// ServiceUsersUploadsStorageGcpBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.GCPConfig.BucketName`.
	ServiceUsersUploadsStorageGcpBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_GCP_BUCKET_NAME"

	// ServiceUsersUploadsStorageProviderEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.Provider`.
	ServiceUsersUploadsStorageProviderEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_PROVIDER"

	// ServiceUsersUploadsStorageS3BucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.S3Config.BucketName`.
	ServiceUsersUploadsStorageS3BucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_S3_BUCKET_NAME"

	// ServiceUsersUploadsStorageUploadFilenameKeyEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.UploadFilenameKey`.
	ServiceUsersUploadsStorageUploadFilenameKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_UPLOAD_FILENAME_KEY"

	// ServiceValidEnumerationsUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.ValidEnumerations.UseSearchService`.
	ServiceValidEnumerationsUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_VALID_ENUMERATIONS_USE_SEARCH_SERVICE"

	// ServiceWebhooksDebugEnvVarKey is the environment variable name to set in order to override `config.Services.Webhooks.Debug`.
	ServiceWebhooksDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_WEBHOOKS_DEBUG"

	// ServiceWorkersDataChangesTopicNameEnvVarKey is the environment variable name to set in order to override `config.Services.Workers.DataChangesTopicName`.
	ServiceWorkersDataChangesTopicNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_WORKERS_DATA_CHANGES_TOPIC_NAME"
)
