package config

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/analytics/segment"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	emailcfg "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/email/sendgrid"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/routing/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/authentication"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipes"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipesteps"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAPIServerConfigFromGoogleCloudRunEnvironment(t *testing.T) {
	baseConfig := &APIServiceConfig{
		Observability: observability.Config{},
		Email: emailcfg.Config{
			Sendgrid: &sendgrid.Config{},
			Provider: emailcfg.ProviderSendgrid,
		},
		Analytics: analyticscfg.Config{
			Provider: analyticscfg.ProviderSegment,
			Segment:  &segment.Config{},
		},
		Encoding: encoding.Config{ContentType: "application/json"},
		Routing:  routingcfg.Config{},
		Database: databasecfg.Config{},
		Meta:     MetaSettings{},
		Events:   msgconfig.Config{},
		Server: http.Config{
			StartupDeadline: time.Second,
		},
		Services: ServicesConfig{
			Recipes: recipesservice.Config{
				PublicMediaURLPrefix: t.Name(),
			},
			RecipeSteps: recipestepsservice.Config{
				PublicMediaURLPrefix: t.Name(),
			},
			Auth: authservice.Config{
				MinimumPasswordLength: 8,
				MinimumUsernameLength: 8,
			},
		},
	}

	f, err := os.CreateTemp("", "testing.json")
	require.NoError(t, err)

	require.NoError(t, json.NewEncoder(f).Encode(baseConfig))

	t.Setenv(ConfigurationFilePathEnvVarKey, f.Name())
	t.Setenv(gcpPortEnvVarKey, "1234")
	t.Setenv(gcpDatabaseSocketDirEnvVarKey, "/example/blah")
	t.Setenv(gcpDatabaseUserEnvVarKey, "user")
	t.Setenv(gcpDatabaseUserPasswordEnvVarKey, "hunter2")
	t.Setenv(gcpDatabaseNameEnvVarKey, "fake_db_name")
	t.Setenv(gcpDatabaseInstanceConnNameEnvVarKey, "fake_conn_name")
	t.Setenv(gcpSendgridTokenEnvVarKey, "fake_sendgrid_token")
	t.Setenv(gcpPostHogKeyEnvVarKey, "fake_posthog_api_key")
	t.Setenv(gcpSegmentTokenEnvVarKey, "fake_segment_token")
	t.Setenv(gcpAlgoliaAPIKeyEnvVarKey, "fake_algolia_api_key")
	t.Setenv(gcpAlgoliaAppIDEnvVarKey, "fake_algolia_app_id")
	t.Setenv(gcpOauth2TokenEncryptionKeyEnvVarKey, "fake_oauth2_token_encryption_key")
	t.Setenv(gcpGoogleSSOClientIDEnvVarKey, "fake_google_sso_client_id")
	t.Setenv(gcpGoogleSSOClientSecretEnvVarKey, "fake_google_sso_client_secret")
	t.Setenv(gcpDataChangesTopicNameEnvVarKey, "data_changes")
	t.Setenv(gcpOutboundEmailsTopicNameEnvVarKey, "outbound_emails")
	t.Setenv(gcpSearchIndexingTopicNameEnvVarKey, "search_indexing")
	t.Setenv(gcpWebhookExecutionTopicNameEnvVarKey, "webhook_execution_requests")
	t.Setenv(gcpUserAggregatorTopicName, "data_aggregator")

	ctx := context.Background()

	cfg, err := GetAPIServiceConfigFromGoogleCloudRunEnvironment(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}
