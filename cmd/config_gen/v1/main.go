package main

import (
	"log"
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
	dbCreateDummyUser                = "database.create_dummy_user"
	dbProvider                       = "database.provider"
	dbDeets                          = "database.connection_details"
	postgres                         = "postgres"
)

type configFunc func(filepath string) error

var (
	files = map[string]configFunc{
		"config_files/coverage.toml":                   coverageConfig,
		"config_files/local.toml":                      localConfig,
		"config_files/development.toml":                developmentConfig,
		"config_files/integration-tests-postgres.toml": buildIntegrationTestForDBImplementation(postgres, postgresDBConnDetails),
		"config_files/production.toml":                 productionConfig,
	}
)

func exampleMetricsConfiguration(cfg *viper.Viper) {
	cfg.Set(metricsProvider, "prometheus")
	cfg.Set(metricsTracer, "jaeger")
	cfg.Set(metricsDBCollectionInterval, time.Second)
	cfg.Set(metricsRuntimeCollectionInterval, time.Second)
}

func developmentConfig(filepath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaDebug, true)
	cfg.Set(metaRunMode, "development")
	cfg.Set(metaStartupDeadline, time.Minute)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authCookieDomain, "prixfixe.dev")
	cfg.Set(authCookieLifetime, oneDay)
	cfg.Set(authSecureCookiesOnly, true)
	cfg.Set(authEnableUserSignup, true)

	// exampleMetricsConfiguration

	cfg.Set(dbProvider, "postgres")
	cfg.Set(dbDeets, "postgresql://prixfixe_dev:vfhfFBwoCoDWTY86bVYa9znk1xcp19IO@database.prixfixe.dev:25060/dev_prixfixe?sslmode=require")

	return cfg.WriteConfigAs(filepath)
}

func localConfig(filepath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaDebug, true)
	cfg.Set(metaRunMode, "development")
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

	// exampleMetricsConfiguration

	cfg.Set(dbProvider, "postgres")
	cfg.Set(dbCreateDummyUser, true)
	cfg.Set(dbDeets, postgresDBConnDetails)

	return cfg.WriteConfigAs(filepath)
}

func coverageConfig(filepath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaRunMode, "development")

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendDebug, true)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, false)
	cfg.Set(authCookieSecret, debugCookieSecret)

	cfg.Set(dbDebug, false)
	cfg.Set(dbProvider, "postgres")
	cfg.Set(dbDeets, postgresDBConnDetails)

	return cfg.WriteConfigAs(filepath)
}

func productionConfig(filepath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaDebug, false)
	cfg.Set(metaRunMode, "production")
	cfg.Set(metaStartupDeadline, time.Minute)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, false)

	cfg.Set(frontendDebug, false)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, false)
	cfg.Set(authCookieDomain, "")
	cfg.Set(authCookieSecret, debugCookieSecret)
	cfg.Set(authCookieLifetime, oneDay)
	cfg.Set(authSecureCookiesOnly, false)
	cfg.Set(authEnableUserSignup, true)

	exampleMetricsConfiguration(cfg)

	cfg.Set(dbDebug, false)
	cfg.Set(dbProvider, "postgres")
	cfg.Set(dbDeets, postgresDBConnDetails)

	return cfg.WriteConfigAs(filepath)
}

func buildIntegrationTestForDBImplementation(dbprov, dbDetails string) configFunc {
	return func(filepath string) error {
		cfg := config.BuildConfig()

		cfg.Set(metaDebug, false)
		cfg.Set(metaRunMode, "testing")

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

		return cfg.WriteConfigAs(filepath)
	}
}

func main() {
	for filepath, fun := range files {
		if err := fun(filepath); err != nil {
			log.Fatal(err)
		}
	}
}
