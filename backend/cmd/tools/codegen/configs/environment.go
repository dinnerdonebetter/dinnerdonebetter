package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path"

	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/hashicorp/go-multierror"
)

type environmentConfigSet struct {
	rootConfig *config.APIServiceConfig
	apiServiceConfigPath,
	dbCleanerConfigPath,
	emailProberConfigPath,
	mealPlanFinalizerConfigPath,
	mealPlanGroceryListInitializerConfigPath,
	mealPlanTaskCreatorConfigPath,
	searchDataIndexSchedulerConfigPath string
}

func stringOrDefault(s1, s2 string) string {
	if s1 != "" {
		return s1
	}
	return s2
}

func renderJSON(obj any) []byte {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(obj); err != nil {
		panic(err)
	}

	return b.Bytes()
}

func writeFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0o0644)
}

func (s *environmentConfigSet) Render(outputDir string) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}
	errs := &multierror.Error{}

	// write files
	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.apiServiceConfigPath, "api_service_config.json")),
		renderJSON(s.rootConfig),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.dbCleanerConfigPath, "job_db_cleaner_config.json")),
		renderJSON(&config.DBCleanerConfig{
			Observability: s.rootConfig.Observability,
			Database:      s.rootConfig.Database,
		}),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.emailProberConfigPath, "job_email_prober_config.json")),
		renderJSON(&config.EmailProberConfig{
			Observability: s.rootConfig.Observability,
			Email:         s.rootConfig.Email,
			Database:      s.rootConfig.Database,
		}),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.mealPlanFinalizerConfigPath, "job_meal_plan_grocery_list_initializer_config.json")),
		renderJSON(&config.MealPlanFinalizerConfig{
			Observability: s.rootConfig.Observability,
			Events:        s.rootConfig.Events,
			Database:      s.rootConfig.Database,
		}),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.mealPlanGroceryListInitializerConfigPath, "job_meal_plan_grocery_list_initializer_config.json")),
		renderJSON(&config.MealPlanGroceryListInitializerConfig{
			Observability: s.rootConfig.Observability,
			Analytics:     s.rootConfig.Analytics,
			Events:        s.rootConfig.Events,
			Database:      s.rootConfig.Database,
		}),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.mealPlanTaskCreatorConfigPath, "job_meal_plan_task_creator_config.json")),
		renderJSON(&config.MealPlanTaskCreatorConfig{
			Observability: s.rootConfig.Observability,
			Analytics:     s.rootConfig.Analytics,
			Events:        s.rootConfig.Events,
			Database:      s.rootConfig.Database,
		}),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := writeFile(
		path.Join(outputDir, stringOrDefault(s.searchDataIndexSchedulerConfigPath, "job_search_data_index_scheduler_config.json")),
		renderJSON(&config.SearchDataIndexSchedulerConfig{
			Observability: s.rootConfig.Observability,
			Events:        s.rootConfig.Events,
			Database:      s.rootConfig.Database,
		}),
	); err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}
