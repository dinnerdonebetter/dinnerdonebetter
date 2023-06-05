package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	analyticsconfig "github.com/dinnerdonebetter/backend/internal/analytics/config"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	emailconfig "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	featureflagsconfig "github.com/dinnerdonebetter/backend/internal/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	"github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	"github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	"github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientgroupsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/validmeasurementconversions"
	validmeaurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	vendorproxyservice "github.com/dinnerdonebetter/backend/internal/services/vendorproxy"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
	websocketsservice "github.com/dinnerdonebetter/backend/internal/services/websockets"

	"github.com/hashicorp/go-multierror"
)

const (
	// DevelopmentRunMode is the run mode for a development environment.
	DevelopmentRunMode runMode = "development"
	// TestingRunMode is the run mode for a testing environment.
	TestingRunMode runMode = "testing"
	// ProductionRunMode is the run mode for a production environment.
	ProductionRunMode runMode = "production"
)

type (
	// runMode describes what method of operation the server is under.
	runMode string

	// CloserFunc calls all io.Closers in the service.
	CloserFunc func()

	// InstanceConfig configures an instance of the service. It is composed of all the other setting structs.
	InstanceConfig struct {
		_             struct{}
		Observability observability.Config      `json:"observability" mapstructure:"observability" toml:"observability,omitempty"`
		Email         emailconfig.Config        `json:"email"         mapstructure:"email"         toml:"email,omitempty"`
		Analytics     analyticsconfig.Config    `json:"analytics"     mapstructure:"analytics"     toml:"analytics,omitempty"`
		FeatureFlags  featureflagsconfig.Config `json:"featureFlags"  mapstructure:"events"        toml:"events,omitempty"`
		Encoding      encoding.Config           `json:"encoding"      mapstructure:"encoding"      toml:"encoding,omitempty"`
		Routing       routing.Config            `json:"routing"       mapstructure:"routing"       toml:"routing,omitempty"`
		Meta          MetaSettings              `json:"meta"          mapstructure:"meta"          toml:"meta,omitempty"`
		Events        msgconfig.Config          `json:"events"        mapstructure:"events"        toml:"events,omitempty"`
		Server        http.Config               `json:"server"        mapstructure:"server"        toml:"server,omitempty"`
		Database      dbconfig.Config           `json:"database"      mapstructure:"database"      toml:"database,omitempty"`
		Services      ServicesConfig            `json:"services"      mapstructure:"services"      toml:"services,omitempty"`
	}

	// ServicesConfig collects the various service configurations.
	ServicesConfig struct {
		_                               struct{}
		MealPlanEvents                  mealplaneventsservice.Config                  `json:"mealPlanEvents"                  mapstructure:"meal_plan_events"                   toml:"meal_plan_events,omitempty"`
		ValidIngredientGroups           validingredientgroupsservice.Config           `json:"validIngredientGroups"           mapstructure:"valid_ingredient_groups"            toml:"valid_ingredient_groups,omitempty"`
		ValidInstruments                validinstrumentsservice.Config                `json:"validInstruments"                mapstructure:"valid_instruments"                  toml:"valid_instruments,omitempty"`
		RecipeStepCompletionConditions  recipestepcompletionconditionsservice.Config  `json:"recipeStepCompletionConditions"  mapstructure:"recipe_step_completion_conditions"  toml:"recipe_step_completion_conditions,omitempty"`
		ValidIngredientPreparations     validingredientpreparationsservice.Config     `json:"validIngredientPreparations"     mapstructure:"valid_ingredient_preparations"      toml:"valid_ingredient_preparations,omitempty"`
		ValidPreparations               validpreparationsservice.Config               `json:"validPreparations"               mapstructure:"valid_preparations"                 toml:"valid_preparations,omitempty"`
		MealPlanOptionVotes             mealplanoptionvotesservice.Config             `json:"mealPlanOptionVotes"             mapstructure:"meal_plan_option_votes"             toml:"meal_plan_option_votes,omitempty"`
		MealPlans                       mealplansservice.Config                       `json:"mealPlans"                       mapstructure:"meal_plans"                         toml:"meal_plans,omitempty"`
		RecipeStepVessels               recipestepvesselsservice.Config               `json:"recipeStepVessels"               mapstructure:"recipe_step_vessels"                toml:"recipe_step_vessels,omitempty"`
		ValidInstrumentMeasurementUnits validingredientmeasurementunitsservice.Config `json:"validInstrumentMeasurementUnits" mapstructure:"valid_ingredient_measurement_units" toml:"valid_ingredient_measurement_units,omitempty"`
		RecipeStepInstruments           recipestepinstrumentsservice.Config           `json:"recipeStepInstruments"           mapstructure:"recipe_step_instruments"            toml:"recipe_step_instruments,omitempty"`
		RecipeStepIngredients           recipestepingredientsservice.Config           `json:"recipeStepIngredients"           mapstructure:"recipe_step_ingredients"            toml:"recipe_step_ingredients,omitempty"`
		ValidMeasurementUnits           validmeaurementunitsservice.Config            `json:"validMeasurementUnits"           mapstructure:"valid_measurement_units"            toml:"valid_measurement_units,omitempty"`
		RecipeStepProducts              recipestepproductsservice.Config              `json:"recipeStepProducts"              mapstructure:"recipe_step_products"               toml:"recipe_step_products,omitempty"`
		Households                      householdsservice.Config                      `json:"households"                      mapstructure:"households"                         toml:"households,omitempty"`
		VendorProxy                     vendorproxyservice.Config                     `json:"vendorProxy"                     mapstructure:"vendor_proxy"                       toml:"vendor_proxy,omitempty"`
		ValidIngredients                validingredientsservice.Config                `json:"validIngredients"                mapstructure:"valid_ingredients"                  toml:"valid_ingredients,omitempty"`
		ValidPreparationInstruments     validpreparationinstrumentsservice.Config     `json:"validPreparationInstruments"     mapstructure:"valid_preparation_instruments"      toml:"valid_preparation_instruments,omitempty"`
		MealPlanOptions                 mealplanoptionsservice.Config                 `json:"mealPlanOptions"                 mapstructure:"meal_plan_options"                  toml:"meal_plan_options,omitempty"`
		Meals                           mealsservice.Config                           `json:"meals"                           mapstructure:"meals"                              toml:"meals,omitempty"`
		ValidIngredientStateIngredients validingredientstateingredients.Config        `json:"validIngredientStateIngredients" mapstructure:"valid_ingredient_state_ingredients" toml:"valid_ingredient_state_ingredients,omitempty"`
		Websockets                      websocketsservice.Config                      `json:"websockets"                      mapstructure:"websockets"                         toml:"websockets,omitempty"`
		ValidIngredientStates           validingredientstates.Config                  `json:"validIngredientStates"           mapstructure:"valid_ingredient_states"            toml:"valid_ingredient_states,omitempty"`
		ValidMeasurementConversions     validmeasurementconversions.Config            `json:"validMeasurementConversions"     mapstructure:"valid_measurement_conversions"      toml:"valid_measurement_conversions,omitempty"`
		MealPlanTasks                   mealplantasks.Config                          `json:"mealPlanTasks"                   mapstructure:"meal_plan_tasks"                    toml:"meal_plan_tasks,omitempty"`
		RecipePrepTasks                 recipepreptasks.Config                        `json:"recipePrepTasks"                 mapstructure:"recipe_prep_tasks"                  toml:"recipe_prep_tasks,omitempty"`
		MealPlanGroceryListItems        mealplangrocerylistitems.Config               `json:"mealPlanGroceryListItems"        mapstructure:"meal_plan_grocery_list_items"       toml:"meal_plan_grocery_list_items,omitempty"`
		ServiceSettings                 servicesettingsservice.Config                 `json:"serviceSettings"                 mapstructure:"service_settings"                   toml:"service_settings,omitempty"`
		ServiceSettingConfigurations    servicesettingconfigurationsservice.Config    `json:"serviceSettingConfigurations"    mapstructure:"service_setting_configurations"     toml:"service_setting_configurations,omitempty"`
		UserIngredientPreferences       useringredientpreferencesservice.Config       `json:"userIngredientPreferences"       mapstructure:"user_ingredient_preferences"        toml:"user_ingredient_preferences,omitempty"`
		Users                           usersservice.Config                           `json:"users"                           mapstructure:"users"                              toml:"users,omitempty"`
		RecipeSteps                     recipestepsservice.Config                     `json:"recipeSteps"                     mapstructure:"recipe_steps"                       toml:"recipe_steps,omitempty"`
		Recipes                         recipesservice.Config                         `json:"recipes"                         mapstructure:"recipes"                            toml:"recipes,omitempty"`
		Webhooks                        webhooksservice.Config                        `json:"webhooks"                        mapstructure:"webhooks"                           toml:"webhooks,omitempty"`
		HouseholdInvitations            householdinvitationsservice.Config            `json:"householdInvitations"            mapstructure:"household_invitations"              toml:"household_invitations,omitempty"`
		Auth                            authservice.Config                            `json:"auth"                            mapstructure:"auth"                               toml:"auth,omitempty"`
	}
)

// EncodeToFile renders your config to a file given your favorite encoder.
func (cfg *InstanceConfig) EncodeToFile(path string, marshaller func(v any) ([]byte, error)) error {
	if cfg == nil {
		return errors.New("nil config")
	}

	byteSlice, err := marshaller(*cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, byteSlice, 0o600)
}

// ValidateWithContext validates a InstanceConfig struct.
func (cfg *InstanceConfig) ValidateWithContext(ctx context.Context, validateServices bool) error {
	var result *multierror.Error

	validatorsToRun := map[string]func(context.Context) error{
		"Routing":                              cfg.Routing.ValidateWithContext,
		"Meta":                                 cfg.Meta.ValidateWithContext,
		"Encoding":                             cfg.Encoding.ValidateWithContext,
		"Analytics":                            cfg.Analytics.ValidateWithContext,
		"Observability":                        cfg.Observability.ValidateWithContext,
		"Database":                             cfg.Database.ValidateWithContext,
		"Server":                               cfg.Server.ValidateWithContext,
		"Email":                                cfg.Email.ValidateWithContext,
		"FeatureFlags":                         cfg.FeatureFlags.ValidateWithContext,
		"Services.Auth":                        cfg.Services.Auth.ValidateWithContext,
		"Services.Users":                       cfg.Services.Users.ValidateWithContext,
		"Services.Webhooks":                    cfg.Services.Webhooks.ValidateWithContext,
		"Services.ValidInstruments":            cfg.Services.ValidInstruments.ValidateWithContext,
		"Services.ValidIngredients":            cfg.Services.ValidIngredients.ValidateWithContext,
		"Services.ValidIngredientGroups":       cfg.Services.ValidIngredientGroups.ValidateWithContext,
		"Services.ValidPreparations":           cfg.Services.ValidPreparations.ValidateWithContext,
		"Services.ValidMeasurementUnits":       cfg.Services.ValidMeasurementUnits.ValidateWithContext,
		"Services.ValidIngredientPreparations": cfg.Services.ValidIngredientPreparations.ValidateWithContext,
		"Services.ValidIngredientStateIngredients": cfg.Services.ValidIngredientStateIngredients.ValidateWithContext,
		"Services.ValidPreparationInstruments":     cfg.Services.ValidPreparationInstruments.ValidateWithContext,
		"Services.ValidInstrumentMeasurementUnits": cfg.Services.ValidInstrumentMeasurementUnits.ValidateWithContext,
		"Services.Recipes":                         cfg.Services.Recipes.ValidateWithContext,
		"Services.RecipeSteps":                     cfg.Services.RecipeSteps.ValidateWithContext,
		"Services.RecipeStepInstruments":           cfg.Services.RecipeStepInstruments.ValidateWithContext,
		"Services.RecipeStepVessels":               cfg.Services.RecipeStepVessels.ValidateWithContext,
		"Services.RecipeStepIngredients":           cfg.Services.RecipeStepIngredients.ValidateWithContext,
		"Services.RecipeStepCompletionConditions":  cfg.Services.RecipeStepCompletionConditions.ValidateWithContext,
		"Services.MealPlans":                       cfg.Services.MealPlans.ValidateWithContext,
		"Services.MealPlanEvents":                  cfg.Services.MealPlanEvents.ValidateWithContext,
		"Services.MealPlanOptions":                 cfg.Services.MealPlanOptions.ValidateWithContext,
		"Services.MealPlanOptionVotes":             cfg.Services.MealPlanOptionVotes.ValidateWithContext,
		"Services.RecipePrepTasks":                 cfg.Services.RecipePrepTasks.ValidateWithContext,
		"Services.MealPlanGroceryListItems":        cfg.Services.MealPlanGroceryListItems.ValidateWithContext,
		"Services.ValidMeasurementConversions":     cfg.Services.ValidMeasurementConversions.ValidateWithContext,
		"Services.ValidIngredientStates":           cfg.Services.ValidIngredientStates.ValidateWithContext,
		"Services.VendorProxy":                     cfg.Services.VendorProxy.ValidateWithContext,
		"Services.ServiceSettings":                 cfg.Services.ServiceSettings.ValidateWithContext,
		"Services.ServiceSettingConfigurations":    cfg.Services.ServiceSettingConfigurations.ValidateWithContext,
		"Services.UserIngredientPreferences":       cfg.Services.UserIngredientPreferences.ValidateWithContext,
	}

	for name, validator := range validatorsToRun {
		if strings.HasPrefix(name, "Services") && !validateServices {
			continue
		}
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}
