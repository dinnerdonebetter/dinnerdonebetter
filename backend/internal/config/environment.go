package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

// defaultCookieLifetime is used when no override is provided; must satisfy cookies.Config validation (>= 5 min).
const defaultCookieLifetime = 24 * time.Hour

// EnvironmentConfigSet contains a way of rendering a set of every config for a given environment to a given folder.
type EnvironmentConfigSet struct {
	RootConfig                               *APIServiceConfig
	ConsumerWebappCookiesOverride            *cookies.Config
	AdminWebappCookiesOverride               *cookies.Config
	SearchDataIndexSchedulerConfigPath       string
	MealPlanFinalizerConfigPath              string
	MealPlanGroceryListInitializerConfigPath string
	MealPlanTaskCreatorConfigPath            string
	DBCleanerConfigPath                      string
	MobileNotificationSchedulerConfigPath    string
	AsyncMessageHandlerConfigPath            string
	AdminWebappConfigPath                    string
	ConsumerWebappConfigPath                 string
	APIServiceConfigPath                     string
	ConsumerWebappPortOverride               uint16
	AdminWebappPortOverride                  uint16
}

func stringOrDefault(s, defaultStr string) string {
	if s != "" {
		return s
	}
	return defaultStr
}

func renderJSON(obj any, pretty bool) []byte {
	var (
		b   []byte
		err error
	)
	if pretty {
		b, err = json.MarshalIndent(obj, "", "\t")
	} else {
		b, err = json.Marshal(obj)
	}

	if err != nil {
		panic(err)
	}

	return b
}

func writeFile(p string, content []byte) error {
	//nolint:gosec // I want this to be 644 I think
	return os.WriteFile(p, content, 0o0644)
}

const (
	apiConfigObservabilityServiceName   = "api_server"
	dbcConfigObservabilityServiceName   = "db_cleaner"
	mpfConfigObservabilityServiceName   = "meal_plan_finalizer"
	mpgliConfigObservabilityServiceName = "meal_plan_grocery_list_initializer"
	mptcConfigObservabilityServiceName  = "meal_plan_task_creator"
	sdisConfigObservabilityServiceName  = "search_data_index_scheduler"
	mnsConfigObservabilityServiceName   = "mobile_notification_scheduler"
	amhConfigObservabilityServiceName   = "async_message_handler"
	awaConfigObservabilityServiceName   = "admin_webapp"
	cwaConfigObservabilityServiceName   = "consumer_webapp"
)

func (s *EnvironmentConfigSet) Render(outputDir string, pretty, validate bool) error {
	if err := os.MkdirAll(outputDir, 0o0750); err != nil {
		return err
	}
	errs := &multierror.Error{}

	// Ensure API server config has the correct observability name before writing.
	s.RootConfig.Observability.Tracing.ServiceName = apiConfigObservabilityServiceName
	s.RootConfig.Observability.Metrics.ServiceName = apiConfigObservabilityServiceName
	s.RootConfig.Observability.Logging.ServiceName = apiConfigObservabilityServiceName
	s.RootConfig.Observability.Profiling.ServiceName = apiConfigObservabilityServiceName
	if s.RootConfig.Routing.Chi != nil {
		s.RootConfig.Routing.Chi.ServiceName = apiConfigObservabilityServiceName
	}

	// write files
	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.APIServiceConfigPath, "api_service_config.json")),
		renderJSON(s.RootConfig, pretty),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	dbcConfig := &DBCleanerConfig{
		Observability: s.RootConfig.Observability,
		Database:      s.RootConfig.Database,
	}
	dbcConfig.Observability.Tracing.ServiceName = dbcConfigObservabilityServiceName
	dbcConfig.Observability.Metrics.ServiceName = dbcConfigObservabilityServiceName
	dbcConfig.Observability.Logging.ServiceName = dbcConfigObservabilityServiceName
	dbcConfig.Observability.Profiling.ServiceName = dbcConfigObservabilityServiceName

	mpfConfig := &MealPlanFinalizerConfig{
		Observability: s.RootConfig.Observability,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
		Queues:        s.RootConfig.Queues,
	}
	mpfConfig.Observability.Tracing.ServiceName = mpfConfigObservabilityServiceName
	mpfConfig.Observability.Metrics.ServiceName = mpfConfigObservabilityServiceName
	mpfConfig.Observability.Logging.ServiceName = mpfConfigObservabilityServiceName
	mpfConfig.Observability.Profiling.ServiceName = mpfConfigObservabilityServiceName

	mpgliConfig := &MealPlanGroceryListInitializerConfig{
		Observability: s.RootConfig.Observability,
		Analytics:     s.RootConfig.Analytics,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
		Queues:        s.RootConfig.Queues,
	}
	mpgliConfig.Observability.Tracing.ServiceName = mpgliConfigObservabilityServiceName
	mpgliConfig.Observability.Metrics.ServiceName = mpgliConfigObservabilityServiceName
	mpgliConfig.Observability.Logging.ServiceName = mpgliConfigObservabilityServiceName
	mpgliConfig.Observability.Profiling.ServiceName = mpgliConfigObservabilityServiceName

	mptcConfig := &MealPlanTaskCreatorConfig{
		Observability: s.RootConfig.Observability,
		Analytics:     s.RootConfig.Analytics,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
		Queues:        s.RootConfig.Queues,
	}
	mptcConfig.Observability.Tracing.ServiceName = mptcConfigObservabilityServiceName
	mptcConfig.Observability.Metrics.ServiceName = mptcConfigObservabilityServiceName
	mptcConfig.Observability.Logging.ServiceName = mptcConfigObservabilityServiceName
	mptcConfig.Observability.Profiling.ServiceName = mptcConfigObservabilityServiceName

	sdisConfig := &SearchDataIndexSchedulerConfig{
		Observability: s.RootConfig.Observability,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
		Queues:        s.RootConfig.Queues,
	}
	sdisConfig.Observability.Tracing.ServiceName = sdisConfigObservabilityServiceName
	sdisConfig.Observability.Metrics.ServiceName = sdisConfigObservabilityServiceName
	sdisConfig.Observability.Logging.ServiceName = sdisConfigObservabilityServiceName
	sdisConfig.Observability.Profiling.ServiceName = sdisConfigObservabilityServiceName

	mnsConfig := &MobileNotificationSchedulerConfig{
		Observability: s.RootConfig.Observability,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
		Queues:        s.RootConfig.Queues,
	}
	mnsConfig.Observability.Tracing.ServiceName = mnsConfigObservabilityServiceName
	mnsConfig.Observability.Metrics.ServiceName = mnsConfigObservabilityServiceName
	mnsConfig.Observability.Logging.ServiceName = mnsConfigObservabilityServiceName
	mnsConfig.Observability.Profiling.ServiceName = mnsConfigObservabilityServiceName

	amhConfig := &AsyncMessageHandlerConfig{
		Storage:           s.RootConfig.Services.DataPrivacy.Uploads.Storage,
		Queues:            s.RootConfig.Queues,
		Email:             s.RootConfig.Email,
		Analytics:         s.RootConfig.Analytics,
		Search:            s.RootConfig.TextSearch,
		Events:            s.RootConfig.Events,
		Observability:     s.RootConfig.Observability,
		Database:          s.RootConfig.Database,
		PushNotifications: s.RootConfig.PushNotifications,
	}
	amhConfig.Observability.Tracing.ServiceName = amhConfigObservabilityServiceName
	amhConfig.Observability.Metrics.ServiceName = amhConfigObservabilityServiceName
	amhConfig.Observability.Logging.ServiceName = amhConfigObservabilityServiceName
	amhConfig.Observability.Profiling.ServiceName = amhConfigObservabilityServiceName

	awaHTTPServer := s.RootConfig.HTTPServer
	if s.AdminWebappPortOverride != 0 {
		awaHTTPServer.Port = s.AdminWebappPortOverride
	}
	awaCookies := cookies.Config{
		CookieName:            "admin_webapp",
		Base64EncodedHashKey:  " ",
		Base64EncodedBlockKey: " ",
		Lifetime:              defaultCookieLifetime,
	}
	if s.AdminWebappCookiesOverride != nil {
		awaCookies = *s.AdminWebappCookiesOverride
	}
	awaConfig := &AdminWebappConfig{
		Cookies:       awaCookies,
		Encoding:      s.RootConfig.Encoding,
		Observability: s.RootConfig.Observability,
		Meta:          s.RootConfig.Meta,
		Routing:       s.RootConfig.Routing,
		HTTPServer:    awaHTTPServer,
	}
	awaConfig.Observability.Tracing.ServiceName = awaConfigObservabilityServiceName
	awaConfig.Observability.Metrics.ServiceName = awaConfigObservabilityServiceName
	awaConfig.Observability.Logging.ServiceName = awaConfigObservabilityServiceName
	awaConfig.Observability.Profiling.ServiceName = awaConfigObservabilityServiceName
	if awaConfig.Routing.Chi != nil {
		chiCopy := *awaConfig.Routing.Chi
		chiCopy.ServiceName = awaConfigObservabilityServiceName
		awaConfig.Routing.Chi = &chiCopy
	}

	cwaHTTPServer := s.RootConfig.HTTPServer
	if s.ConsumerWebappPortOverride != 0 {
		cwaHTTPServer.Port = s.ConsumerWebappPortOverride
	}
	cwaCookies := cookies.Config{
		CookieName:            "consumer_webapp",
		Base64EncodedHashKey:  " ",
		Base64EncodedBlockKey: " ",
		Lifetime:              defaultCookieLifetime,
	}
	if s.ConsumerWebappCookiesOverride != nil {
		cwaCookies = *s.ConsumerWebappCookiesOverride
	}
	cwaConfig := &ConsumerWebappConfig{
		Cookies:       cwaCookies,
		Encoding:      s.RootConfig.Encoding,
		Observability: s.RootConfig.Observability,
		Meta:          s.RootConfig.Meta,
		Routing:       s.RootConfig.Routing,
		HTTPServer:    cwaHTTPServer,
	}
	cwaConfig.Observability.Tracing.ServiceName = cwaConfigObservabilityServiceName
	cwaConfig.Observability.Metrics.ServiceName = cwaConfigObservabilityServiceName
	cwaConfig.Observability.Logging.ServiceName = cwaConfigObservabilityServiceName
	cwaConfig.Observability.Profiling.ServiceName = cwaConfigObservabilityServiceName
	if cwaConfig.Routing.Chi != nil {
		chiCopy := *cwaConfig.Routing.Chi
		chiCopy.ServiceName = cwaConfigObservabilityServiceName
		cwaConfig.Routing.Chi = &chiCopy
	}

	if validate {
		allConfigs := []validation.ValidatableWithContext{
			s.RootConfig,
			dbcConfig,
			mpfConfig,
			mpgliConfig,
			mptcConfig,
			sdisConfig,
			mnsConfig,
			amhConfig,
			awaConfig,
			cwaConfig,
		}
		for i, cfg := range allConfigs {
			if err := cfg.ValidateWithContext(context.Background()); err != nil {
				errs = multierror.Append(errs, fmt.Errorf("validating config %d: %w", i, err))
				continue
			}
		}
	}

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	pathToConfigMap := map[string][]byte{
		path.Join(outputDir, stringOrDefault(s.DBCleanerConfigPath, "job_db_cleaner_config.json")):                                              renderJSON(dbcConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanFinalizerConfigPath, "job_meal_plan_finalizer_config.json")):                             renderJSON(mpfConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanGroceryListInitializerConfigPath, "job_meal_plan_grocery_list_initializer_config.json")): renderJSON(mpgliConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanTaskCreatorConfigPath, "job_meal_plan_task_creator_config.json")):                        renderJSON(mptcConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.SearchDataIndexSchedulerConfigPath, "job_search_data_index_scheduler_config.json")):              renderJSON(sdisConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MobileNotificationSchedulerConfigPath, "job_mobile_notification_scheduler_config.json")):         renderJSON(mnsConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.AsyncMessageHandlerConfigPath, "async_message_handler_config.json")):                             renderJSON(amhConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.AdminWebappConfigPath, "admin_webapp_config.json")):                                              renderJSON(awaConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.ConsumerWebappConfigPath, "consumer_webapp_config.json")):                                        renderJSON(cwaConfig, pretty),
	}

	for p, b := range pathToConfigMap {
		if err := writeFile(p, b); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}
