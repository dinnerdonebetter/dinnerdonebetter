package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gitlab.com/prixfixe/prixfixe/internal/v1/config"

	"github.com/spf13/viper"
)

const (
	defaultPort                      = 8888
	oneDay                           = 24 * time.Hour
	debugCookieSecret                = "HEREISA32CHARSECRETWHICHISMADEUP"
	defaultFrontendFilepath          = "/frontend"
	postgresDBConnDetails            = "postgres://dbuser:hunter2@database:5432/prixfixe?sslmode=disable"
	metaDebug                        = "meta.debug"
	metaRunMode                      = "meta.run_mode"
	metaStartupDeadline              = "meta.startup_deadline"
	serverHTTPPort                   = "server.http_port"
	serverDebug                      = "server.debug"
	frontendDebug                    = "frontend.debug"
	frontendStaticFilesDir           = "frontend.static_files_directory"
	frontendCacheStatics             = "frontend.cache_static_files"
	authDebug                        = "auth.debug"
	authCookieDomain                 = "auth.cookie_domain"
	authCookieSecret                 = "auth.cookie_secret"
	authCookieLifetime               = "auth.cookie_lifetime"
	authSecureCookiesOnly            = "auth.secure_cookies_only"
	authEnableUserSignup             = "auth.enable_user_signup"
	metricsProvider                  = "metrics.metrics_provider"
	metricsTracer                    = "metrics.tracing_provider"
	metricsDBCollectionInterval      = "metrics.database_metrics_collection_interval"
	metricsRuntimeCollectionInterval = "metrics.runtime_metrics_collection_interval"
	dbDebug                          = "database.debug"
	dbProvider                       = "database.provider"
	dbDeets                          = "database.connection_details"
	dbCreateDummyUser                = "database.create_dummy_user"
	validInstrumentSearchIndexPath   = "search.valid_instruments_index_path"
	validIngredientSearchIndexPath   = "search.valid_ingredients_index_path"
	validPreparationSearchIndexPath  = "search.valid_preparations_index_path"

	// run modes
	developmentEnv = "development"
	testingEnv     = "testing"

	// database providers
	postgres = "postgres"

	// search index paths
	defaultValidInstrumentsSearchIndexPath  = "valid_instruments.bleve"
	defaultValidIngredientsSearchIndexPath  = "valid_ingredients.bleve"
	defaultValidPreparationsSearchIndexPath = "valid_preparations.bleve"
)

type configFunc func(filepath string) error

var (
	files = map[string]configFunc{
		"environments/dev/config.toml":                                      developmentConfig,
		"environments/local/config.toml":                                    localConfig,
		"environments/testing/config_files/coverage.toml":                   coverageConfig,
		"environments/testing/config_files/frontend-tests.toml":             frontendTestsConfig,
		"environments/testing/config_files/integration-tests-postgres.toml": buildIntegrationTestForDBImplementation(postgres, postgresDBConnDetails),
	}
)

func exampleMetricsConfiguration(cfg *viper.Viper) {
	cfg.Set(metricsProvider, "prometheus")
	cfg.Set(metricsTracer, "jaeger")
	cfg.Set(metricsDBCollectionInterval, time.Second)
	cfg.Set(metricsRuntimeCollectionInterval, time.Second)
}

func removeEmptyCookieSecretSettingFromFile(path string) error {
	configAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	configToWrite := strings.Replace(string(configAsBytes), "  cookie_secret = \"\"\n", "", 1)

	return ioutil.WriteFile(path, []byte(configToWrite), os.FileMode(0644))
}

func developmentConfig(path string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaDebug, true)
	cfg.Set(metaRunMode, developmentEnv)
	cfg.Set(metaStartupDeadline, time.Minute)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authCookieDomain, "prixfixe.dev")
	cfg.Set(authCookieLifetime, oneDay)
	cfg.Set(authCookieSecret, "")
	cfg.Set(authSecureCookiesOnly, true)
	cfg.Set(authEnableUserSignup, true)

	// exampleMetricsConfiguration

	cfg.Set(dbProvider, postgres)
	cfg.Set(dbDeets, "postgresql://prixfixe_dev:vfhfFBwoCoDWTY86bVYa9znk1xcp19IO@database.prixfixe.dev:25060/dev_prixfixe?sslmode=require")

	cfg.Set(validInstrumentSearchIndexPath, "/etc/prixfixe_search_indices/valid_instruments.bleve")
	cfg.Set(validIngredientSearchIndexPath, "/etc/prixfixe_search_indices/valid_ingredients.bleve")
	cfg.Set(validPreparationSearchIndexPath, "/etc/prixfixe_search_indices/valid_preparations.bleve")

	if err := cfg.WriteConfigAs(path); err != nil {
		return err
	}

	return removeEmptyCookieSecretSettingFromFile(path)
}

func localConfig(path string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaDebug, true)
	cfg.Set(metaRunMode, developmentEnv)
	cfg.Set(metaStartupDeadline, time.Minute)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authCookieDomain, "localhost")
	cfg.Set(authCookieSecret, debugCookieSecret)
	cfg.Set(authCookieLifetime, oneDay)
	cfg.Set(authSecureCookiesOnly, false)
	cfg.Set(authEnableUserSignup, true)

	exampleMetricsConfiguration(cfg)

	cfg.Set(dbProvider, postgres)
	cfg.Set(dbCreateDummyUser, true)
	cfg.Set(dbDeets, postgresDBConnDetails)

	cfg.Set(validInstrumentSearchIndexPath, defaultValidInstrumentsSearchIndexPath)
	cfg.Set(validIngredientSearchIndexPath, defaultValidIngredientsSearchIndexPath)
	cfg.Set(validPreparationSearchIndexPath, defaultValidPreparationsSearchIndexPath)

	if err := cfg.WriteConfigAs(path); err != nil {
		return err
	}

	return removeEmptyCookieSecretSettingFromFile(path)
}

func coverageConfig(filepath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaRunMode, testingEnv)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendDebug, true)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, false)
	cfg.Set(authCookieSecret, debugCookieSecret)

	cfg.Set(dbDebug, false)
	cfg.Set(dbProvider, postgres)
	cfg.Set(dbDeets, postgresDBConnDetails)

	cfg.Set(validInstrumentSearchIndexPath, defaultValidInstrumentsSearchIndexPath)
	cfg.Set(validIngredientSearchIndexPath, defaultValidIngredientsSearchIndexPath)
	cfg.Set(validPreparationSearchIndexPath, defaultValidPreparationsSearchIndexPath)

	return cfg.WriteConfigAs(filepath)
}

func buildIntegrationTestForDBImplementation(dbprov, dbDetails string) configFunc {
	return func(filepath string) error {
		cfg := config.BuildConfig()

		cfg.Set(metaRunMode, testingEnv)
		cfg.Set(metaDebug, false)

		sd := time.Minute
		cfg.Set(metaStartupDeadline, sd)

		cfg.Set(serverHTTPPort, defaultPort)
		cfg.Set(serverDebug, true)

		cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
		cfg.Set(authCookieSecret, debugCookieSecret)

		cfg.Set(metricsProvider, "prometheus")
		cfg.Set(metricsTracer, "jaeger")

		cfg.Set(dbDebug, false)
		cfg.Set(dbProvider, dbprov)
		cfg.Set(dbDeets, dbDetails)

		cfg.Set(validInstrumentSearchIndexPath, defaultValidInstrumentsSearchIndexPath)
		cfg.Set(validIngredientSearchIndexPath, defaultValidIngredientsSearchIndexPath)
		cfg.Set(validPreparationSearchIndexPath, defaultValidPreparationsSearchIndexPath)

		return cfg.WriteConfigAs(filepath)
	}
}

func frontendTestsConfig(filePath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaRunMode, testingEnv)
	cfg.Set(metaStartupDeadline, time.Minute)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendDebug, true)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, true)
	cfg.Set(authCookieDomain, "localhost")
	cfg.Set(authCookieSecret, debugCookieSecret)
	cfg.Set(authCookieLifetime, oneDay)
	cfg.Set(authSecureCookiesOnly, false)
	cfg.Set(authEnableUserSignup, true)

	cfg.Set(metricsProvider, "prometheus")
	cfg.Set(metricsTracer, "jaeger")
	cfg.Set(metricsDBCollectionInterval, time.Second)
	cfg.Set(metricsRuntimeCollectionInterval, time.Second)

	cfg.Set(dbDebug, true)
	cfg.Set(dbProvider, postgres)
	cfg.Set(dbDeets, postgresDBConnDetails)

	cfg.Set(validInstrumentSearchIndexPath, defaultValidInstrumentsSearchIndexPath)
	cfg.Set(validIngredientSearchIndexPath, defaultValidIngredientsSearchIndexPath)
	cfg.Set(validPreparationSearchIndexPath, defaultValidPreparationsSearchIndexPath)

	if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
		return fmt.Errorf("error writing developmentEnv config: %w", writeErr)
	}

	return nil
}

func main() {
	for filepath, fun := range files {
		if err := fun(filepath); err != nil {
			log.Fatal(err)
		}
	}
}
