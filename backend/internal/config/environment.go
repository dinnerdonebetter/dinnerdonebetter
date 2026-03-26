package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability"
)

// EnvironmentConfigSet contains a way of rendering a set of every config for a given environment to a given folder.
type EnvironmentConfigSet struct {
	RootConfig                               *APIServiceConfig
	SearchDataIndexSchedulerConfigPath       string
	MealPlanFinalizerConfigPath              string
	MealPlanGroceryListInitializerConfigPath string
	MealPlanTaskCreatorConfigPath            string
	DBCleanerConfigPath                      string
	MobileNotificationSchedulerConfigPath    string
	AsyncMessageHandlerConfigPath            string
	EmailDeliverabilityTestConfigPath        string
	QueueTestJobConfigPath                   string
	APIServiceConfigPath                     string
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

// disableWorkerOtelMetrics turns off runtime and host metrics for worker configs to reduce cardinality.
// It clones the Otel config so the root config (API server, async message handler) is not mutated.
func disableWorkerOtelMetrics(obs *observability.Config) {
	if obs == nil || obs.Metrics.Otel == nil {
		return
	}
	copied := *obs.Metrics.Otel
	copied.EnableRuntimeMetrics = false
	copied.EnableHostMetrics = false
	obs.Metrics.Otel = &copied
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
	edtConfigObservabilityServiceName   = "email_deliverability_test"
	qtConfigObservabilityServiceName    = "queue_test"
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
	disableWorkerOtelMetrics(&dbcConfig.Observability)

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
	disableWorkerOtelMetrics(&mpfConfig.Observability)

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
	disableWorkerOtelMetrics(&mpgliConfig.Observability)

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
	disableWorkerOtelMetrics(&mptcConfig.Observability)

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
	disableWorkerOtelMetrics(&sdisConfig.Observability)

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
	disableWorkerOtelMetrics(&mnsConfig.Observability)

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
		BaseURL:           s.RootConfig.BaseURL,
	}
	amhConfig.Observability.Tracing.ServiceName = amhConfigObservabilityServiceName
	amhConfig.Observability.Metrics.ServiceName = amhConfigObservabilityServiceName
	amhConfig.Observability.Logging.ServiceName = amhConfigObservabilityServiceName
	amhConfig.Observability.Profiling.ServiceName = amhConfigObservabilityServiceName

	edtServiceEnv := "prod"
	if strings.Contains(outputDir, "localdev") {
		edtServiceEnv = "dev"
	} else if strings.Contains(outputDir, "testing") {
		edtServiceEnv = "testing"
	}
	edtConfig := &EmailDeliverabilityTestConfig{
		Observability:         s.RootConfig.Observability,
		Email:                 s.RootConfig.Email,
		RecipientEmailAddress: "verygoodsoftwarenotvirus@protonmail.com",
		ServiceEnvironment:    edtServiceEnv,
	}
	edtConfig.Observability.Tracing.ServiceName = edtConfigObservabilityServiceName
	edtConfig.Observability.Metrics.ServiceName = edtConfigObservabilityServiceName
	edtConfig.Observability.Logging.ServiceName = edtConfigObservabilityServiceName
	edtConfig.Observability.Profiling.ServiceName = edtConfigObservabilityServiceName
	disableWorkerOtelMetrics(&edtConfig.Observability)

	qtConfig := &QueueTestJobConfig{
		Observability: s.RootConfig.Observability,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
		Queues:        s.RootConfig.Queues,
	}
	qtConfig.Observability.Tracing.ServiceName = qtConfigObservabilityServiceName
	qtConfig.Observability.Metrics.ServiceName = qtConfigObservabilityServiceName
	qtConfig.Observability.Logging.ServiceName = qtConfigObservabilityServiceName
	qtConfig.Observability.Profiling.ServiceName = qtConfigObservabilityServiceName
	disableWorkerOtelMetrics(&qtConfig.Observability)

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
			edtConfig,
			qtConfig,
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
		path.Join(outputDir, stringOrDefault(s.EmailDeliverabilityTestConfigPath, "job_email_deliverability_test_config.json")):                 renderJSON(edtConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.QueueTestJobConfigPath, "job_queue_test_config.json")):                                           renderJSON(qtConfig, pretty),
	}

	for p, b := range pathToConfigMap {
		if err := writeFile(p, b); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}
