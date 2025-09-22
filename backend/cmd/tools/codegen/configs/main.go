package main

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/config"
)

const (
	defaultHTTPPort = 8000
	defaultGRPCPort = 8001
	maxAttempts     = 50
	otelServiceName = "dinner_done_better"

	/* #nosec G101 */
	debugCookieHashKey = "HEREISA32CHARSECRETWHICHISMADEUP"

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	// message provider topics.
	dataChangesTopicName              = "data_changes"
	outboundEmailsTopicName           = "outbound_emails"
	searchIndexRequestsTopicName      = "search_index_requests"
	userDataAggregationTopicName      = "user_data_aggregation_requests"
	webhookExecutionRequestsTopicName = "webhook_execution_requests"

	contentTypeJSON = "application/json"
)

func main() {
	devOutputPath := "deploy/environments/dev/kustomize/configs"

	envConfigs := map[string]*config.EnvironmentConfigSet{
		devOutputPath: {
			RootConfig: buildDevEnvironmentServerConfig(),
		},
		"deploy/environments/localdev/config_files": {
			RootConfig: buildLocalDevConfig(),
		},
		"deploy/environments/testing/config_files": {
			APIServiceConfigPath: "integration-tests-config.json",
			RootConfig:           buildIntegrationTestsConfig(),
		},
	}

	for p, cfg := range envConfigs {
		// we don't want to validate the cloud env configs because they use env vars and cluster secrets to load values
		shouldRenderPrettyAndValidate := p != devOutputPath

		if err := cfg.Render(p, true, shouldRenderPrettyAndValidate); err != nil {
			panic(fmt.Errorf("validating config %s: %w", p, err))
		}
	}
}
