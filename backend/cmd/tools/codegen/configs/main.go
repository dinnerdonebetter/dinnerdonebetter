package main

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
)

const (
	defaultPort = 8000
	/* #nosec G101 */
	debugCookieHashKey = "HEREISA32CHARSECRETWHICHISMADEUP"

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	otelServiceName = "dinner_done_better_api"

	// message provider topics.
	dataChangesTopicName              = "data_changes"
	outboundEmailsTopicName           = "outbound_emails"
	searchIndexRequestsTopicName      = "search_index_requests"
	userDataAggregationTopicName      = "user_data_aggregation_requests"
	webhookExecutionRequestsTopicName = "webhook_execution_requests"

	maxAttempts = 50

	contentTypeJSON               = "application/json"
	workerQueueAddress            = "worker_queue:6379"
	localOAuth2TokenEncryptionKey = debugCookieHashKey
)

var (
	devPostgresDBConnectionDetails = databasecfg.ConnectionDetails{
		Username:   "dbuser",
		Password:   "hunter2",
		Database:   "dinner-done-better",
		Host:       "pgdatabase",
		Port:       5432,
		DisableSSL: true,
	}
)

func main() {
	envConfigs := map[string]*config.EnvironmentConfigSet{
		"deploy/kustomize/environments/localdev/configs": {
			RootConfig: buildLocaldevKubernetesConfig(),
		},
		"environments/testing/config_files": {
			APIServiceConfigPath: "integration-tests-config.json",
			RootConfig:           buildIntegrationTestsConfig(),
		},
		"environments/localdev/config_files": {
			RootConfig: buildLocalDevConfig(),
		},
		"environments/dev/config_files": {
			APIServiceConfigPath: "service-config.json",
			RootConfig:           buildDevEnvironmentServerConfig(),
		},
	}

	for p, cfg := range envConfigs {
		// we don't want to validate the cloud env configs because they
		// need secrets they extract from cloud providers.
		shouldRenderPrettyAndValidate := p != "environments/dev/config_files"

		if err := cfg.Render(p, shouldRenderPrettyAndValidate, shouldRenderPrettyAndValidate); err != nil {
			panic(fmt.Errorf("validating config %s: %w", p, err))
		}
	}
}
