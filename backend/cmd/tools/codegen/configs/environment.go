package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/hashicorp/go-multierror"
)

type environmentConfigSet struct {
	rootConfig                               *config.APIServiceConfig
	apiServiceConfigPath                     string
	dbCleanerConfigPath                      string
	emailProberConfigPath                    string
	mealPlanFinalizerConfigPath              string
	mealPlanGroceryListInitializerConfigPath string
	mealPlanTaskCreatorConfigPath            string
	searchDataIndexSchedulerConfigPath       string
	renderPretty                             bool
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

func (s *environmentConfigSet) Render(outputDir string, validate bool) error {
	if err := os.MkdirAll(outputDir, 0o0750); err != nil {
		return err
	}
	errs := &multierror.Error{}

	// write files
	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.apiServiceConfigPath, "api_service_config.json")),
		renderJSON(s.rootConfig, s.renderPretty),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	dbcConfig := &config.DBCleanerConfig{
		Observability: s.rootConfig.Observability,
		Database:      s.rootConfig.Database,
	}
	empConfig := &config.EmailProberConfig{
		Observability: s.rootConfig.Observability,
		Email:         s.rootConfig.Email,
		Database:      s.rootConfig.Database,
	}
	mpfConfig := &config.MealPlanFinalizerConfig{
		Observability: s.rootConfig.Observability,
		Events:        s.rootConfig.Events,
		Database:      s.rootConfig.Database,
	}
	mpgliConfig := &config.MealPlanGroceryListInitializerConfig{
		Observability: s.rootConfig.Observability,
		Analytics:     s.rootConfig.Analytics,
		Events:        s.rootConfig.Events,
		Database:      s.rootConfig.Database,
	}
	mptcConfig := &config.MealPlanTaskCreatorConfig{
		Observability: s.rootConfig.Observability,
		Analytics:     s.rootConfig.Analytics,
		Events:        s.rootConfig.Events,
		Database:      s.rootConfig.Database,
	}
	sdisConfig := &config.SearchDataIndexSchedulerConfig{
		Observability: s.rootConfig.Observability,
		Events:        s.rootConfig.Events,
		Database:      s.rootConfig.Database,
	}

	if validate {
		allConfigs := []validation.ValidatableWithContext{
			s.rootConfig,
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

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.dbCleanerConfigPath, "job_db_cleaner_config.json")),
		renderJSON(dbcConfig, s.renderPretty),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.emailProberConfigPath, "job_email_prober_config.json")),
		renderJSON(empConfig, s.renderPretty),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.mealPlanFinalizerConfigPath, "job_meal_plan_finalizer_config.json")),
		renderJSON(mpfConfig, s.renderPretty),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.mealPlanGroceryListInitializerConfigPath, "job_meal_plan_grocery_list_initializer_config.json")),
		renderJSON(mpgliConfig, s.renderPretty),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.mealPlanTaskCreatorConfigPath, "job_meal_plan_task_creator_config.json")),
		renderJSON(mptcConfig, s.renderPretty),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.searchDataIndexSchedulerConfigPath, "job_search_data_index_scheduler_config.json")),
		renderJSON(sdisConfig, s.renderPretty),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}
