package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

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
	empConfig := &EmailProberConfig{
		Observability: s.RootConfig.Observability,
		Email:         s.RootConfig.Email,
		Database:      s.RootConfig.Database,
	}
	mpfConfig := &MealPlanFinalizerConfig{
		Observability: s.RootConfig.Observability,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
	}
	mpgliConfig := &MealPlanGroceryListInitializerConfig{
		Observability: s.RootConfig.Observability,
		Analytics:     s.RootConfig.Analytics,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
	}
	mptcConfig := &MealPlanTaskCreatorConfig{
		Observability: s.RootConfig.Observability,
		Analytics:     s.RootConfig.Analytics,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
	}
	sdisConfig := &SearchDataIndexSchedulerConfig{
		Observability: s.RootConfig.Observability,
		Events:        s.RootConfig.Events,
		Database:      s.RootConfig.Database,
	}

	if validate {
		allConfigs := []validation.ValidatableWithContext{
			s.RootConfig,
			dbcConfig,
			empConfig,
			mpfConfig,
			mpgliConfig,
			mptcConfig,
			sdisConfig,
		}
		for i, cfg := range allConfigs {
			if err := cfg.ValidateWithContext(context.Background()); err != nil {
				errs = multierror.Append(errs, fmt.Errorf("validating config %d: %v", i, err))
				continue
			}
		}
	}

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	pathToConfigMap := map[string][]byte{
		path.Join(outputDir, stringOrDefault(s.DBCleanerConfigPath, "job_db_cleaner_config.json")):                                              renderJSON(dbcConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.DBCleanerConfigPath, "job_db_cleaner_config.json")):                                              renderJSON(dbcConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.EmailProberConfigPath, "job_email_prober_config.json")):                                          renderJSON(empConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanFinalizerConfigPath, "job_meal_plan_finalizer_config.json")):                             renderJSON(mpfConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanGroceryListInitializerConfigPath, "job_meal_plan_grocery_list_initializer_config.json")): renderJSON(mpgliConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanTaskCreatorConfigPath, "job_meal_plan_task_creator_config.json")):                        renderJSON(mptcConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.SearchDataIndexSchedulerConfigPath, "job_search_data_index_scheduler_config.json")):              renderJSON(sdisConfig, pretty),
	}

	for p, b := range pathToConfigMap {
		if err := writeFile(p, b); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}
