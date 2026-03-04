package main

import (
	"encoding/base64"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
)

const (
	defaultHTTPPort = 8000
	defaultGRPCPort = 8001
	maxAttempts     = 50
	otelServiceName = "api_server"

	/* #nosec G101 */
	debugCookieHashKey = "HEREISA32CHARSECRETWHICHISMADEUP"

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	// message provider topics.
	dataChangesTopicName              = "data_changes"
	outboundEmailsTopicName           = "outbound_emails"
	searchIndexRequestsTopicName      = "search_index_requests"
	mobileNotificationsTopicName      = "mobile_notifications"
	userDataAggregationTopicName      = "user_data_aggregation_requests"
	webhookExecutionRequestsTopicName = "webhook_execution_requests"
)

var (
	contentTypeJSON         = encoding.ContentTypeToString(encoding.ContentTypeJSON)
	localdevConsumerCookies = buildLocaldevConsumerCookies()
)

func buildLocaldevConsumerCookies() *cookies.Config {
	key := base64.StdEncoding.EncodeToString([]byte(debugCookieHashKey))
	return &cookies.Config{
		CookieName:            "consumer_session",
		Base64EncodedHashKey:  key,
		Base64EncodedBlockKey: key,
	}
}

func main() {
	// localdev config is generated to two locations:
	// - config_files/ for docker-compose usage
	// - kustomize/configs/ for Kubernetes usage (hostnames overridden via env vars)
	localdevConfig := buildLocalDevConfig()

	envConfigs := map[string]*config.EnvironmentConfigSet{
		"deploy/environments/localdev/config_files": {
			RootConfig:                    localdevConfig,
			ConsumerWebappCookiesOverride: localdevConsumerCookies,
			ConsumerWebappPortOverride:    8889, // matches consumer.sh proxy.app_port
		},
		"deploy/environments/localdev/kustomize/configs": {
			RootConfig:                    localdevConfig,
			ConsumerWebappCookiesOverride: localdevConsumerCookies,
			ConsumerWebappPortOverride:    8889, // matches consumer.sh proxy.app_port
		},
		"deploy/environments/testing/config_files": {
			APIServiceConfigPath: "integration-tests-config.json",
			RootConfig:           buildIntegrationTestsConfig(),
		},
		"deploy/environments/prod/kustomize/configs": {
			RootConfig: buildProdConfig(),
		},
	}

	for p, cfg := range envConfigs {
		if err := cfg.Render(p, true, true); err != nil {
			panic(fmt.Errorf("validating config %s: %w", p, err))
		}
	}
}
