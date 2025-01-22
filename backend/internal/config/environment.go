package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

// EnvironmentConfigSet contains a way of rendering a set of every config for a given environment to a given folder.
type EnvironmentConfigSet struct {
	RootConfig                               *APIServiceConfig
	APIServiceConfigPath                     string
	DBCleanerConfigPath                      string
	EmailProberConfigPath                    string
	MealPlanFinalizerConfigPath              string
	MealPlanGroceryListInitializerConfigPath string
	MealPlanTaskCreatorConfigPath            string
	SearchDataIndexSchedulerConfigPath       string
	AsyncMessageHandlerConfigPath            string
	AdminWebappConfigPath                    string
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
	dbcConfigObservabilityServiceName   = "db_cleaner"
	empConfigObservabilityServiceName   = "email_prober"
	mpfConfigObservabilityServiceName   = "meal_plan_finalizer"
	mpgliConfigObservabilityServiceName = "meal_plan_grocery_list_initializer"
	mptcConfigObservabilityServiceName  = "meal_plan_task_creator"
	sdisConfigObservabilityServiceName  = "search_data_index_scheduler"
	amhConfigObservabilityServiceName   = "async_message_handler"
	awaConfigObservabilityServiceName   = "admin_webapp"
)

func (s *EnvironmentConfigSet) Render(outputDir string, pretty, validate bool) error {
	if err := os.MkdirAll(outputDir, 0o0750); err != nil {
		return err
	}
	errs := &multierror.Error{}

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

	empConfig := &EmailProberConfig{
		Observability: s.RootConfig.Observability,
		Email:         s.RootConfig.Email,
		Database:      s.RootConfig.Database,
	}
	empConfig.Observability.Tracing.ServiceName = empConfigObservabilityServiceName
	empConfig.Observability.Metrics.ServiceName = empConfigObservabilityServiceName
	empConfig.Observability.Logging.ServiceName = empConfigObservabilityServiceName

	mpfConfig := &MealPlanFinalizerConfig{
		Observability: s.RootConfig.Observability,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
		Queues:        s.RootConfig.Queues,
	}
	mpfConfig.Observability.Tracing.ServiceName = mpfConfigObservabilityServiceName
	mpfConfig.Observability.Metrics.ServiceName = mpfConfigObservabilityServiceName
	mpfConfig.Observability.Logging.ServiceName = mpfConfigObservabilityServiceName

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

	sdisConfig := &SearchDataIndexSchedulerConfig{
		Observability: s.RootConfig.Observability,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
		Queues:        s.RootConfig.Queues,
	}
	sdisConfig.Observability.Tracing.ServiceName = sdisConfigObservabilityServiceName
	sdisConfig.Observability.Metrics.ServiceName = sdisConfigObservabilityServiceName
	sdisConfig.Observability.Logging.ServiceName = sdisConfigObservabilityServiceName

	amhConfig := &AsyncMessageHandlerConfig{
		Storage:       s.RootConfig.Services.DataPrivacy.Uploads.Storage,
		Queues:        s.RootConfig.Queues,
		Email:         s.RootConfig.Email,
		Analytics:     s.RootConfig.Analytics,
		Search:        s.RootConfig.Search,
		Events:        s.RootConfig.Events,
		Observability: s.RootConfig.Observability,
		Database:      s.RootConfig.Database,
	}
	amhConfig.Observability.Tracing.ServiceName = amhConfigObservabilityServiceName
	amhConfig.Observability.Metrics.ServiceName = amhConfigObservabilityServiceName
	amhConfig.Observability.Logging.ServiceName = amhConfigObservabilityServiceName

	awaConfig := &AdminWebappConfig{
		// No equivalent for the following configs, I'll just have to rely on env vars to set them for now:
		//	- Cookies
		//	- APIServiceConnection
		//	- APIClientCache
		Cookies: cookies.Config{
			CookieName:            " ",
			Base64EncodedHashKey:  " ",
			Base64EncodedBlockKey: " ",
		},
		Encoding:      s.RootConfig.Encoding,
		Observability: s.RootConfig.Observability,
		Meta:          s.RootConfig.Meta,
		Routing:       s.RootConfig.Routing,
		HTTPServer:    s.RootConfig.HTTPServer,
	}
	awaConfig.Observability.Tracing.ServiceName = awaConfigObservabilityServiceName
	awaConfig.Observability.Metrics.ServiceName = awaConfigObservabilityServiceName
	awaConfig.Observability.Logging.ServiceName = awaConfigObservabilityServiceName

	if validate {
		allConfigs := []validation.ValidatableWithContext{
			s.RootConfig,
			dbcConfig,
			empConfig,
			mpfConfig,
			mpgliConfig,
			mptcConfig,
			sdisConfig,
			amhConfig,
			awaConfig,
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
		path.Join(outputDir, stringOrDefault(s.EmailProberConfigPath, "job_email_prober_config.json")):                                          renderJSON(empConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanFinalizerConfigPath, "job_meal_plan_finalizer_config.json")):                             renderJSON(mpfConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanGroceryListInitializerConfigPath, "job_meal_plan_grocery_list_initializer_config.json")): renderJSON(mpgliConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanTaskCreatorConfigPath, "job_meal_plan_task_creator_config.json")):                        renderJSON(mptcConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.SearchDataIndexSchedulerConfigPath, "job_search_data_index_scheduler_config.json")):              renderJSON(sdisConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.AsyncMessageHandlerConfigPath, "async_message_handler_config.json")):                             renderJSON(amhConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.AdminWebappConfigPath, "admin_webapp_config.json")):                                              renderJSON(awaConfig, pretty),
	}

	for p, b := range pathToConfigMap {
		if err := writeFile(p, b); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}
