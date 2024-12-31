package envvars

/*
This file contains a reference of all valid service environment variables.
You should be able to suss out what controls what from its name.
*/

const (
	// ServerStartupDeadlineEnvVarKey is the environment variable name to set in order to override `config.Server.StartupDeadline`.
	ServerStartupDeadlineEnvVarKey = "DINNER_DONE_BETTER_SERVER_STARTUP_DEADLINE"

	// ServiceDataPrivacyUploadsStorageGcpBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.GCPConfig.BucketName`.
	ServiceDataPrivacyUploadsStorageGcpBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_GCP_BUCKET_NAME"

	// ServiceValidPreparationsUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.ValidPreparations.UseSearchService`.
	ServiceValidPreparationsUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_VALID_PREPARATIONS_USE_SEARCH_SERVICE"

	// ServiceUsersUploadsStorageGcpBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.GCPConfig.BucketName`.
	ServiceUsersUploadsStorageGcpBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_GCP_BUCKET_NAME"

	// ServiceValidIngredientStatesUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.ValidIngredientStates.UseSearchService`.
	ServiceValidIngredientStatesUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_VALID_INGREDIENT_STATES_USE_SEARCH_SERVICE"

	// EventsConsumerSqsQueueAddressEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.SQS.QueueAddress`.
	EventsConsumerSqsQueueAddressEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_SQS_QUEUE_ADDRESS"

	// RoutingProviderEnvVarKey is the environment variable name to set in order to override `config.Routing.Provider`.
	RoutingProviderEnvVarKey = "DINNER_DONE_BETTER_ROUTING_PROVIDER"

	// ServiceRecipesUploadsStorageS3BucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.S3Config.BucketName`.
	ServiceRecipesUploadsStorageS3BucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_S3_BUCKET_NAME"

	// AnalyticsPosthogProjectAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Analytics.Posthog.ProjectAPIKey`.
	AnalyticsPosthogProjectAPIKeyEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_POSTHOG_PROJECT_API_KEY"

	// AnalyticsRudderstackDataPlaneURLEnvVarKey is the environment variable name to set in order to override `config.Analytics.Rudderstack.DataPlaneURL`.
	AnalyticsRudderstackDataPlaneURLEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_RUDDERSTACK_DATA_PLANE_URL"

	// SearchCircuitBreakerMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.Search.CircuitBreakerConfig.MinimumSampleThreshold`.
	SearchCircuitBreakerMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_SEARCH_CIRCUIT_BREAKER_MINIMUM_SAMPLE_THRESHOLD"

	// ObservabilityTracingTracingSpanCollectionProbabilityEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.SpanCollectionProbability`.
	ObservabilityTracingTracingSpanCollectionProbabilityEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_TRACING_SPAN_COLLECTION_PROBABILITY"

	// RoutingValidDomainsEnvVarKey is the environment variable name to set in order to override `config.Routing.ValidDomains`.
	RoutingValidDomainsEnvVarKey = "DINNER_DONE_BETTER_ROUTING_VALID_DOMAINS"

	// DatabaseConnectionDetailsDisableSslEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.DisableSSL`.
	DatabaseConnectionDetailsDisableSslEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_DISABLE_SSL"

	// DatabaseConnectionDetailsPortEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Port`.
	DatabaseConnectionDetailsPortEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_PORT"

	// QueuesDataChangesTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.DataChangesTopicName`.
	QueuesDataChangesTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_DATA_CHANGES_TOPIC_NAME"

	// QueuesUserDataAggregationTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.UserDataAggregationTopicName`.
	QueuesUserDataAggregationTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_USER_DATA_AGGREGATION_TOPIC_NAME"

	// EmailCircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.Email.CircuitBreakerConfig.MinimumSampleThreshold`.
	EmailCircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_EMAIL_CIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// FeatureFlagsLaunchDarklycircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.LaunchDarkly.CircuitBreakerConfig.MinimumSampleThreshold`.
	FeatureFlagsLaunchDarklycircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_LAUNCH_DARKLYCIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// ObservabilityLoggingProviderEnvVarKey is the environment variable name to set in order to override `config.Observability.Logging.Provider`.
	ObservabilityLoggingProviderEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_LOGGING_PROVIDER"

	// ServiceRecipesUploadsStorageUploadFilenameKeyEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.UploadFilenameKey`.
	ServiceRecipesUploadsStorageUploadFilenameKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_UPLOAD_FILENAME_KEY"

	// FeatureFlagsPosthogPersonalAPIKeyEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.PostHog.PersonalAPIKey`.
	FeatureFlagsPosthogPersonalAPIKeyEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_POSTHOG_PERSONAL_API_KEY"

	// ObservabilityMetricsProviderEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Provider`.
	ObservabilityMetricsProviderEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_PROVIDER"

	// DatabaseConnectionDetailsDatabaseEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Database`.
	DatabaseConnectionDetailsDatabaseEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_DATABASE"

	// ServiceValidMeasurementUnitsUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.ValidMeasurementUnits.UseSearchService`.
	ServiceValidMeasurementUnitsUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_VALID_MEASUREMENT_UNITS_USE_SEARCH_SERVICE"

	// SearchProviderEnvVarKey is the environment variable name to set in order to override `config.Search.Provider`.
	SearchProviderEnvVarKey = "DINNER_DONE_BETTER_SEARCH_PROVIDER"

	// EventsConsumerProviderEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.Provider`.
	EventsConsumerProviderEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_PROVIDER"

	// EncodingContentTypeEnvVarKey is the environment variable name to set in order to override `config.Encoding.ContentType`.
	EncodingContentTypeEnvVarKey = "DINNER_DONE_BETTER_ENCODING_CONTENT_TYPE"

	// ObservabilityLoggingOutputFilepathEnvVarKey is the environment variable name to set in order to override `config.Observability.Logging.OutputFilepath`.
	ObservabilityLoggingOutputFilepathEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_LOGGING_OUTPUT_FILEPATH"

	// ServerHTTPSCertificateFilepathEnvVarKey is the environment variable name to set in order to override `config.Server.HTTPSCertificateFile`.
	ServerHTTPSCertificateFilepathEnvVarKey = "DINNER_DONE_BETTER_SERVER_HTTPS_CERTIFICATE_FILEPATH"

	// ServiceRecipesUploadsDebugEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Debug`.
	ServiceRecipesUploadsDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_DEBUG"

	// ServiceRecipeStepsUploadsStorageBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.RecipeSteps.Uploads.Storage.BucketName`.
	ServiceRecipeStepsUploadsStorageBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPE_STEPS_UPLOADS_STORAGE_BUCKET_NAME"

	// ServiceHouseholdInvitationsDebugEnvVarKey is the environment variable name to set in order to override `config.Services.HouseholdInvitations.Debug`.
	ServiceHouseholdInvitationsDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_HOUSEHOLD_INVITATIONS_DEBUG"

	// QueuesWebhookExecutionRequestsTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.WebhookExecutionRequestsTopicName`.
	QueuesWebhookExecutionRequestsTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_WEBHOOK_EXECUTION_REQUESTS_TOPIC_NAME"

	// AnalyticsCircuitBreakerErrorRateEnvVarKey is the environment variable name to set in order to override `config.Analytics.CircuitBreakerConfig.ErrorRate`.
	AnalyticsCircuitBreakerErrorRateEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_CIRCUIT_BREAKER_ERROR_RATE"

	// DatabaseLogQueriesEnvVarKey is the environment variable name to set in order to override `config.Database.LogQueries`.
	DatabaseLogQueriesEnvVarKey = "DINNER_DONE_BETTER_DATABASE_LOG_QUERIES"

	// ServiceDataPrivacyUploadsStorageFilesystemRootDirectoryEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.FilesystemConfig.RootDirectory`.
	ServiceDataPrivacyUploadsStorageFilesystemRootDirectoryEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_FILESYSTEM_ROOT_DIRECTORY"

	// QueuesOutboundEmailsTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.OutboundEmailsTopicName`.
	QueuesOutboundEmailsTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_OUTBOUND_EMAILS_TOPIC_NAME"

	// DatabaseDebugEnvVarKey is the environment variable name to set in order to override `config.Database.Debug`.
	DatabaseDebugEnvVarKey = "DINNER_DONE_BETTER_DATABASE_DEBUG"

	// ServiceAuthOauth2RefreshTokenLifespanEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.OAuth2.RefreshTokenLifespan`.
	ServiceAuthOauth2RefreshTokenLifespanEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_OAUTH2REFRESH_TOKEN_LIFESPAN"

	// ServiceRecipeStepsUploadsStorageBucketPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.RecipeSteps.Uploads.Storage.BucketPrefix`.
	ServiceRecipeStepsUploadsStorageBucketPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPE_STEPS_UPLOADS_STORAGE_BUCKET_PREFIX"

	// ServiceAuthOauth2DomainEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.OAuth2.Domain`.
	ServiceAuthOauth2DomainEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_OAUTH2DOMAIN"

	// EventsConsumerRedisQueueAddressesEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.Redis.QueueAddresses`.
	EventsConsumerRedisQueueAddressesEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_REDIS_QUEUE_ADDRESSES"

	// DatabaseConnectionDetailsHostEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Host`.
	DatabaseConnectionDetailsHostEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_HOST"

	// ServiceAuthSsoConfigGoogleClientIDEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.SSO.Google.ClientID`.
	ServiceAuthSsoConfigGoogleClientIDEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_SSO_CONFIG_GOOGLE_CLIENT_ID"

	// ServiceAuthJwtSigningKeyEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.JWTSigningKey`.
	ServiceAuthJwtSigningKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_JWT_SIGNING_KEY"

	// FeatureFlagsPosthogCircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.PostHog.CircuitBreakerConfig.MinimumSampleThreshold`.
	FeatureFlagsPosthogCircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_POSTHOG_CIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// SearchElasticsearchIndexOperationTimeoutEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.IndexOperationTimeout`.
	SearchElasticsearchIndexOperationTimeoutEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_INDEX_OPERATION_TIMEOUT"

	// ServiceRecipesUploadsStorageBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.BucketName`.
	ServiceRecipesUploadsStorageBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_BUCKET_NAME"

	// ServiceDataPrivacyUploadsStorageProviderEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.Provider`.
	ServiceDataPrivacyUploadsStorageProviderEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_PROVIDER"

	// ServiceDataPrivacyUploadsStorageUploadFilenameKeyEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.UploadFilenameKey`.
	ServiceDataPrivacyUploadsStorageUploadFilenameKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_UPLOAD_FILENAME_KEY"

	// QueuesSearchIndexRequestsTopicNameEnvVarKey is the environment variable name to set in order to override `config.Queues.SearchIndexRequestsTopicName`.
	QueuesSearchIndexRequestsTopicNameEnvVarKey = "DINNER_DONE_BETTER_QUEUES_SEARCH_INDEX_REQUESTS_TOPIC_NAME"

	// SearchAlgoliaAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Search.Algolia.APIKey`.
	SearchAlgoliaAPIKeyEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ALGOLIA_API_KEY"

	// ServiceUsersUploadsStorageBucketPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.BucketPrefix`.
	ServiceUsersUploadsStorageBucketPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_BUCKET_PREFIX"

	// ServiceWebhooksDebugEnvVarKey is the environment variable name to set in order to override `config.Services.Webhooks.Debug`.
	ServiceWebhooksDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_WEBHOOKS_DEBUG"

	// DatabaseConnectionDetailsUsernameEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Username`.
	DatabaseConnectionDetailsUsernameEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_USERNAME"

	// ServiceRecipesUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.UseSearchService`.
	ServiceRecipesUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_USE_SEARCH_SERVICE"

	// EmailMailjetSecretKeyEnvVarKey is the environment variable name to set in order to override `config.Email.Mailjet.SecretKey`.
	EmailMailjetSecretKeyEnvVarKey = "DINNER_DONE_BETTER_EMAIL_MAILJET_SECRET_KEY"

	// FeatureFlagsPosthogCircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.PostHog.CircuitBreakerConfig.ErrorRate`.
	FeatureFlagsPosthogCircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_POSTHOG_CIRCUIT_BREAKING_ERROR_RATE"

	// EventsPublisherSqsQueueAddressEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.SQS.QueueAddress`.
	EventsPublisherSqsQueueAddressEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_SQS_QUEUE_ADDRESS"

	// ObservabilityTracingTracingServiceNameEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.ServiceName`.
	ObservabilityTracingTracingServiceNameEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_TRACING_SERVICE_NAME"

	// ServiceDataPrivacyUploadsStorageS3BucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.S3Config.BucketName`.
	ServiceDataPrivacyUploadsStorageS3BucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_S3_BUCKET_NAME"

	// FeatureFlagsProviderEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.Provider`.
	FeatureFlagsProviderEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_PROVIDER"

	// ObservabilityTracingOtelgrpcCollectorEndpointEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.Otel.CollectorEndpoint`.
	ObservabilityTracingOtelgrpcCollectorEndpointEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_OTELGRPC_COLLECTOR_ENDPOINT"

	// ServiceRecipeStepsUploadsStorageS3BucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.RecipeSteps.Uploads.Storage.S3Config.BucketName`.
	ServiceRecipeStepsUploadsStorageS3BucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPE_STEPS_UPLOADS_STORAGE_S3_BUCKET_NAME"

	// ServiceUsersUploadsStorageUploadFilenameKeyEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.UploadFilenameKey`.
	ServiceUsersUploadsStorageUploadFilenameKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_UPLOAD_FILENAME_KEY"

	// EventsPublisherRedisPasswordEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.Redis.Password`.
	EventsPublisherRedisPasswordEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_REDIS_PASSWORD"

	// ObservabilityTracingTracingProviderEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.Provider`.
	ObservabilityTracingTracingProviderEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_TRACING_PROVIDER"

	// ServiceRecipesUploadsStorageFilesystemRootDirectoryEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.FilesystemConfig.RootDirectory`.
	ServiceRecipesUploadsStorageFilesystemRootDirectoryEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_FILESYSTEM_ROOT_DIRECTORY"

	// EmailMailgunPrivateAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Email.Mailgun.PrivateAPIKey`.
	EmailMailgunPrivateAPIKeyEnvVarKey = "DINNER_DONE_BETTER_EMAIL_MAILGUN_PRIVATE_API_KEY"

	// DatabaseOauth2TokenEncryptionKeyEnvVarKey is the environment variable name to set in order to override `config.Database.OAuth2TokenEncryptionKey`.
	DatabaseOauth2TokenEncryptionKeyEnvVarKey = "DINNER_DONE_BETTER_DATABASE_OAUTH2_TOKEN_ENCRYPTION_KEY"

	// ServerHTTPPortEnvVarKey is the environment variable name to set in order to override `config.Server.HTTPPort`.
	ServerHTTPPortEnvVarKey = "DINNER_DONE_BETTER_SERVER_HTTP_PORT"

	// ServiceUsersUploadsStorageS3BucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.S3Config.BucketName`.
	ServiceUsersUploadsStorageS3BucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_S3_BUCKET_NAME"

	// AnalyticsSegmentAPITokenEnvVarKey is the environment variable name to set in order to override `config.Analytics.Segment.APIToken`.
	AnalyticsSegmentAPITokenEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_SEGMENT_API_TOKEN"

	// AnalyticsPosthogPersonalAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Analytics.Posthog.PersonalAPIKey`.
	AnalyticsPosthogPersonalAPIKeyEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_POSTHOG_PERSONAL_API_KEY"

	// FeatureFlagsPosthogProjectAPIKeyEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.PostHog.ProjectAPIKey`.
	FeatureFlagsPosthogProjectAPIKeyEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_POSTHOG_PROJECT_API_KEY"

	// ObservabilityMetricsOtelServiceNameEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.ServiceName`.
	ObservabilityMetricsOtelServiceNameEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_SERVICE_NAME"

	// FeatureFlagsLaunchDarklyinitTimeoutEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.LaunchDarkly.InitTimeout`.
	FeatureFlagsLaunchDarklyinitTimeoutEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_LAUNCH_DARKLYINIT_TIMEOUT"

	// ObservabilityTracingCloudtraceGoogleCloudTraceProjectIDEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.CloudTrace.ProjectID`.
	ObservabilityTracingCloudtraceGoogleCloudTraceProjectIDEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_CLOUDTRACE_GOOGLE_CLOUD_TRACE_PROJECT_ID"

	// ServiceRecipeStepsUploadsDebugEnvVarKey is the environment variable name to set in order to override `config.Services.RecipeSteps.Uploads.Debug`.
	ServiceRecipeStepsUploadsDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPE_STEPS_UPLOADS_DEBUG"

	// ServiceAuthSsoConfigGoogleCallbackURLEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.SSO.Google.CallbackURL`.
	ServiceAuthSsoConfigGoogleCallbackURLEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_SSO_CONFIG_GOOGLE_CALLBACK_URL"

	// ServiceDataPrivacyUploadsStorageBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.BucketName`.
	ServiceDataPrivacyUploadsStorageBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_BUCKET_NAME"

	// ServiceValidIngredientsUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.ValidIngredients.UseSearchService`.
	ServiceValidIngredientsUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_VALID_INGREDIENTS_USE_SEARCH_SERVICE"

	// ServiceValidInstrumentsUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.ValidInstruments.UseSearchService`.
	ServiceValidInstrumentsUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_VALID_INSTRUMENTS_USE_SEARCH_SERVICE"

	// ServiceRecipesUploadsStorageBucketPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.BucketPrefix`.
	ServiceRecipesUploadsStorageBucketPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_BUCKET_PREFIX"

	// AnalyticsProviderEnvVarKey is the environment variable name to set in order to override `config.Analytics.Provider`.
	AnalyticsProviderEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_PROVIDER"

	// EmailMailjetAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Email.Mailjet.APIKey`.
	EmailMailjetAPIKeyEnvVarKey = "DINNER_DONE_BETTER_EMAIL_MAILJET_API_KEY"

	// ServerDebugEnvVarKey is the environment variable name to set in order to override `config.Server.Debug`.
	ServerDebugEnvVarKey = "DINNER_DONE_BETTER_SERVER_DEBUG"

	// DatabaseMaxPingAttemptsEnvVarKey is the environment variable name to set in order to override `config.Database.MaxPingAttempts`.
	DatabaseMaxPingAttemptsEnvVarKey = "DINNER_DONE_BETTER_DATABASE_MAX_PING_ATTEMPTS"

	// ObservabilityMetricsOtelInsecureEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.Insecure`.
	ObservabilityMetricsOtelInsecureEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_INSECURE"

	// ObservabilityLoggingLevelEnvVarKey is the environment variable name to set in order to override `config.Observability.Logging.Level`.
	ObservabilityLoggingLevelEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_LOGGING_LEVEL"

	// ServiceAuthEnableUserSignupEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.EnableUserSignup`.
	ServiceAuthEnableUserSignupEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_ENABLE_USER_SIGNUP"

	// AnalyticsPosthogCircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.Analytics.Posthog.CircuitBreakerConfig.ErrorRate`.
	AnalyticsPosthogCircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_POSTHOG_CIRCUIT_BREAKING_ERROR_RATE"

	// EmailCircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.Email.CircuitBreakerConfig.ErrorRate`.
	EmailCircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_EMAIL_CIRCUIT_BREAKING_ERROR_RATE"

	// FeatureFlagsCircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.CircuitBreakingConfig.MinimumSampleThreshold`.
	FeatureFlagsCircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_CIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// FeatureFlagsCircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.CircuitBreakingConfig.ErrorRate`.
	FeatureFlagsCircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_CIRCUIT_BREAKING_ERROR_RATE"

	// ServiceRecipeStepsUploadsStorageUploadFilenameKeyEnvVarKey is the environment variable name to set in order to override `config.Services.RecipeSteps.Uploads.Storage.UploadFilenameKey`.
	ServiceRecipeStepsUploadsStorageUploadFilenameKeyEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPE_STEPS_UPLOADS_STORAGE_UPLOAD_FILENAME_KEY"

	// ServiceAuthDataChangesTopicNameEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.DataChangesTopicName`.
	ServiceAuthDataChangesTopicNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_DATA_CHANGES_TOPIC_NAME"

	// SearchElasticsearchAddressEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.Address`.
	SearchElasticsearchAddressEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_ADDRESS"

	// ServiceAuthOauth2DebugEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.OAuth2.Debug`.
	ServiceAuthOauth2DebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_OAUTH2DEBUG"

	// ServiceRecipeStepsPublicMediaURLPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.RecipeSteps.PublicMediaURLPrefix`.
	ServiceRecipeStepsPublicMediaURLPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPE_STEPS_PUBLIC_MEDIA_URL_PREFIX"

	// ServiceDataPrivacyUploadsDebugEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Debug`.
	ServiceDataPrivacyUploadsDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_DEBUG"

	// ServiceDataPrivacyUploadsStorageBucketPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.DataPrivacy.Uploads.Storage.BucketPrefix`.
	ServiceDataPrivacyUploadsStorageBucketPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_DATA_PRIVACY_UPLOADS_STORAGE_BUCKET_PREFIX"

	// ServiceAuthJwtLifetimeEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.JWTLifetime`.
	ServiceAuthJwtLifetimeEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_JWT_LIFETIME"

	// EventsConsumerRedisPasswordEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.Redis.Password`.
	EventsConsumerRedisPasswordEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_REDIS_PASSWORD"

	// EventsPublisherPubsubProjectIDEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.PubSub.ProjectID`.
	EventsPublisherPubsubProjectIDEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_PUBSUB_PROJECT_ID"

	// ServiceUsersUploadsDebugEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Debug`.
	ServiceUsersUploadsDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_DEBUG"

	// EventsConsumerPubsubProjectIDEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.PubSub.ProjectID`.
	EventsConsumerPubsubProjectIDEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_PUBSUB_PROJECT_ID"

	// ObservabilityMetricsOtelCollectionTimeoutEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.CollectionTimeout`.
	ObservabilityMetricsOtelCollectionTimeoutEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_COLLECTION_TIMEOUT"

	// RoutingServiceNameEnvVarKey is the environment variable name to set in order to override `config.Routing.ServiceName`.
	RoutingServiceNameEnvVarKey = "DINNER_DONE_BETTER_ROUTING_SERVICE_NAME"

	// RoutingSilenceRouteLoggingEnvVarKey is the environment variable name to set in order to override `config.Routing.SilenceRouteLogging`.
	RoutingSilenceRouteLoggingEnvVarKey = "DINNER_DONE_BETTER_ROUTING_SILENCE_ROUTE_LOGGING"

	// ServiceRecipeStepsUploadsStorageGcpBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.RecipeSteps.Uploads.Storage.GCPConfig.BucketName`.
	ServiceRecipeStepsUploadsStorageGcpBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPE_STEPS_UPLOADS_STORAGE_GCP_BUCKET_NAME"

	// EmailSendgridAPITokenEnvVarKey is the environment variable name to set in order to override `config.Email.Sendgrid.APIToken`.
	EmailSendgridAPITokenEnvVarKey = "DINNER_DONE_BETTER_EMAIL_SENDGRID_API_TOKEN"

	// FeatureFlagsLaunchDarklysdkKeyEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.LaunchDarkly.SDKKey`.
	FeatureFlagsLaunchDarklysdkKeyEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_LAUNCH_DARKLYSDK_KEY"

	// SearchAlgoliaAppIDEnvVarKey is the environment variable name to set in order to override `config.Search.Algolia.AppID`.
	SearchAlgoliaAppIDEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ALGOLIA_APP_ID"

	// ObservabilityMetricsOtelCollectorEndpointEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.CollectorEndpoint`.
	ObservabilityMetricsOtelCollectorEndpointEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_COLLECTOR_ENDPOINT"

	// SearchCircuitBreakerErrorRateEnvVarKey is the environment variable name to set in order to override `config.Search.CircuitBreakerConfig.ErrorRate`.
	SearchCircuitBreakerErrorRateEnvVarKey = "DINNER_DONE_BETTER_SEARCH_CIRCUIT_BREAKER_ERROR_RATE"

	// MetaRunModeEnvVarKey is the environment variable name to set in order to override `config.Meta.RunMode`.
	MetaRunModeEnvVarKey = "DINNER_DONE_BETTER_META_RUN_MODE"

	// RoutingEnableCorsForLocalhostEnvVarKey is the environment variable name to set in order to override `config.Routing.EnableCORSForLocalhost`.
	RoutingEnableCorsForLocalhostEnvVarKey = "DINNER_DONE_BETTER_ROUTING_ENABLE_CORS_FOR_LOCALHOST"

	// EventsPublisherRedisUsernameEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.Redis.Username`.
	EventsPublisherRedisUsernameEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_REDIS_USERNAME"

	// ServiceMealsUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.Meals.UseSearchService`.
	ServiceMealsUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_MEALS_USE_SEARCH_SERVICE"

	// ServiceAuthSsoConfigGoogleClientSecretEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.SSO.Google.ClientSecret`.
	ServiceAuthSsoConfigGoogleClientSecretEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_SSO_CONFIG_GOOGLE_CLIENT_SECRET"

	// ServiceRecipesPublicMediaURLPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.PublicMediaURLPrefix`.
	ServiceRecipesPublicMediaURLPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_PUBLIC_MEDIA_URL_PREFIX"

	// EmailMailgunDomainEnvVarKey is the environment variable name to set in order to override `config.Email.Mailgun.Domain`.
	EmailMailgunDomainEnvVarKey = "DINNER_DONE_BETTER_EMAIL_MAILGUN_DOMAIN"

	// EventsConsumerRedisUsernameEnvVarKey is the environment variable name to set in order to override `config.Events.Consumer.Redis.Username`.
	EventsConsumerRedisUsernameEnvVarKey = "DINNER_DONE_BETTER_EVENTS_CONSUMER_REDIS_USERNAME"

	// ServiceValidVesselsUseSearchServiceEnvVarKey is the environment variable name to set in order to override `config.Services.ValidVessels.UseSearchService`.
	ServiceValidVesselsUseSearchServiceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_VALID_VESSELS_USE_SEARCH_SERVICE"

	// ServiceUsersPublicMediaURLPrefixEnvVarKey is the environment variable name to set in order to override `config.Services.Users.PublicMediaURLPrefix`.
	ServiceUsersPublicMediaURLPrefixEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_PUBLIC_MEDIA_URL_PREFIX"

	// ServiceWorkersDataChangesTopicNameEnvVarKey is the environment variable name to set in order to override `config.Services.Workers.DataChangesTopicName`.
	ServiceWorkersDataChangesTopicNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_WORKERS_DATA_CHANGES_TOPIC_NAME"

	// ServiceRecipeStepsUploadsStorageFilesystemRootDirectoryEnvVarKey is the environment variable name to set in order to override `config.Services.RecipeSteps.Uploads.Storage.FilesystemConfig.RootDirectory`.
	ServiceRecipeStepsUploadsStorageFilesystemRootDirectoryEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPE_STEPS_UPLOADS_STORAGE_FILESYSTEM_ROOT_DIRECTORY"

	// EmailProviderEnvVarKey is the environment variable name to set in order to override `config.Email.Provider`.
	EmailProviderEnvVarKey = "DINNER_DONE_BETTER_EMAIL_PROVIDER"

	// ServiceAuthJwtAudienceEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.JWTAudience`.
	ServiceAuthJwtAudienceEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_JWT_AUDIENCE"

	// AnalyticsCircuitBreakerMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.Analytics.CircuitBreakerConfig.MinimumSampleThreshold`.
	AnalyticsCircuitBreakerMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_CIRCUIT_BREAKER_MINIMUM_SAMPLE_THRESHOLD"

	// MetaDebugEnvVarKey is the environment variable name to set in order to override `config.Meta.Debug`.
	MetaDebugEnvVarKey = "DINNER_DONE_BETTER_META_DEBUG"

	// ObservabilityTracingOtelgrpcInsecureEnvVarKey is the environment variable name to set in order to override `config.Observability.Tracing.Otel.Insecure`.
	ObservabilityTracingOtelgrpcInsecureEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_TRACING_OTELGRPC_INSECURE"

	// EventsPublisherProviderEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.Provider`.
	EventsPublisherProviderEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_PROVIDER"

	// DatabaseConnectionDetailsPasswordEnvVarKey is the environment variable name to set in order to override `config.Database.ConnectionDetails.Password`.
	DatabaseConnectionDetailsPasswordEnvVarKey = "DINNER_DONE_BETTER_DATABASE_CONNECTION_DETAILS_PASSWORD"

	// ServiceUsersUploadsStorageProviderEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.Provider`.
	ServiceUsersUploadsStorageProviderEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_PROVIDER"

	// ServiceUsersUploadsStorageBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.BucketName`.
	ServiceUsersUploadsStorageBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_BUCKET_NAME"

	// ServiceAuthOauth2AccessTokenLifespanEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.OAuth2.AccessTokenLifespan`.
	ServiceAuthOauth2AccessTokenLifespanEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_OAUTH2ACCESS_TOKEN_LIFESPAN"

	// ServiceRecipesUploadsStorageProviderEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.Provider`.
	ServiceRecipesUploadsStorageProviderEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_PROVIDER"

	// FeatureFlagsLaunchDarklycircuitBreakingErrorRateEnvVarKey is the environment variable name to set in order to override `config.FeatureFlags.LaunchDarkly.CircuitBreakerConfig.ErrorRate`.
	FeatureFlagsLaunchDarklycircuitBreakingErrorRateEnvVarKey = "DINNER_DONE_BETTER_FEATURE_FLAGS_LAUNCH_DARKLYCIRCUIT_BREAKING_ERROR_RATE"

	// ServiceUsersUploadsStorageFilesystemRootDirectoryEnvVarKey is the environment variable name to set in order to override `config.Services.Users.Uploads.Storage.FilesystemConfig.RootDirectory`.
	ServiceUsersUploadsStorageFilesystemRootDirectoryEnvVarKey = "DINNER_DONE_BETTER_SERVICE_USERS_UPLOADS_STORAGE_FILESYSTEM_ROOT_DIRECTORY"

	// AnalyticsPosthogCircuitBreakingMinimumSampleThresholdEnvVarKey is the environment variable name to set in order to override `config.Analytics.Posthog.CircuitBreakerConfig.MinimumSampleThreshold`.
	AnalyticsPosthogCircuitBreakingMinimumSampleThresholdEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_POSTHOG_CIRCUIT_BREAKING_MINIMUM_SAMPLE_THRESHOLD"

	// DatabasePingWaitPeriodEnvVarKey is the environment variable name to set in order to override `config.Database.PingWaitPeriod`.
	DatabasePingWaitPeriodEnvVarKey = "DINNER_DONE_BETTER_DATABASE_PING_WAIT_PERIOD"

	// ServerHTTPSCertificateKeyFilepathEnvVarKey is the environment variable name to set in order to override `config.Server.HTTPSCertificateKeyFile`.
	ServerHTTPSCertificateKeyFilepathEnvVarKey = "DINNER_DONE_BETTER_SERVER_HTTPS_CERTIFICATE_KEY_FILEPATH"

	// ServiceAuthMinimumUsernameLengthEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.MinimumUsernameLength`.
	ServiceAuthMinimumUsernameLengthEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_MINIMUM_USERNAME_LENGTH"

	// ServiceAuthMinimumPasswordLengthEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.MinimumPasswordLength`.
	ServiceAuthMinimumPasswordLengthEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_MINIMUM_PASSWORD_LENGTH"

	// ServiceRecipeStepsUploadsStorageProviderEnvVarKey is the environment variable name to set in order to override `config.Services.RecipeSteps.Uploads.Storage.Provider`.
	ServiceRecipeStepsUploadsStorageProviderEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPE_STEPS_UPLOADS_STORAGE_PROVIDER"

	// SearchElasticsearchUsernameEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.Username`.
	SearchElasticsearchUsernameEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_USERNAME"

	// ObservabilityMetricsOtelCollectionIntervalEnvVarKey is the environment variable name to set in order to override `config.Observability.Metrics.Otel.CollectionInterval`.
	ObservabilityMetricsOtelCollectionIntervalEnvVarKey = "DINNER_DONE_BETTER_OBSERVABILITY_METRICS_OTEL_COLLECTION_INTERVAL"

	// ServiceAuthDebugEnvVarKey is the environment variable name to set in order to override `config.Services.Auth.Debug`.
	ServiceAuthDebugEnvVarKey = "DINNER_DONE_BETTER_SERVICE_AUTH_DEBUG"

	// AnalyticsRudderstackAPIKeyEnvVarKey is the environment variable name to set in order to override `config.Analytics.Rudderstack.APIKey`.
	AnalyticsRudderstackAPIKeyEnvVarKey = "DINNER_DONE_BETTER_ANALYTICS_RUDDERSTACK_API_KEY"

	// SearchAlgoliaTimeoutEnvVarKey is the environment variable name to set in order to override `config.Search.Algolia.Timeout`.
	SearchAlgoliaTimeoutEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ALGOLIA_TIMEOUT"

	// SearchElasticsearchCaCertEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.CACert`.
	SearchElasticsearchCaCertEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_CA_CERT"

	// EventsPublisherRedisQueueAddressesEnvVarKey is the environment variable name to set in order to override `config.Events.Publisher.Redis.QueueAddresses`.
	EventsPublisherRedisQueueAddressesEnvVarKey = "DINNER_DONE_BETTER_EVENTS_PUBLISHER_REDIS_QUEUE_ADDRESSES"

	// ServiceOauth2ClientsCreationDisabledEnvVarKey is the environment variable name to set in order to override `config.Services.OAuth2Clients.OAuth2ClientCreationDisabled`.
	ServiceOauth2ClientsCreationDisabledEnvVarKey = "DINNER_DONE_BETTER_SERVICE_OAUTH2_CLIENTS_CREATION_DISABLED"

	// ServiceRecipesUploadsStorageGcpBucketNameEnvVarKey is the environment variable name to set in order to override `config.Services.Recipes.Uploads.Storage.GCPConfig.BucketName`.
	ServiceRecipesUploadsStorageGcpBucketNameEnvVarKey = "DINNER_DONE_BETTER_SERVICE_RECIPES_UPLOADS_STORAGE_GCP_BUCKET_NAME"

	// SearchElasticsearchPasswordEnvVarKey is the environment variable name to set in order to override `config.Search.Elasticsearch.Password`.
	SearchElasticsearchPasswordEnvVarKey = "DINNER_DONE_BETTER_SEARCH_ELASTICSEARCH_PASSWORD"

	// DatabaseRunMigrationsEnvVarKey is the environment variable name to set in order to override `config.Database.RunMigrations`.
	DatabaseRunMigrationsEnvVarKey = "DINNER_DONE_BETTER_DATABASE_RUN_MIGRATIONS"
)
