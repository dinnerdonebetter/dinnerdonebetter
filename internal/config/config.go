package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"

	analyticsconfig "github.com/prixfixeco/backend/internal/analytics/config"
	dbconfig "github.com/prixfixeco/backend/internal/database/config"
	emailconfig "github.com/prixfixeco/backend/internal/email/config"
	"github.com/prixfixeco/backend/internal/encoding"
	featureflagsconfig "github.com/prixfixeco/backend/internal/featureflags/config"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/routing"
	"github.com/prixfixeco/backend/internal/server"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	householdinvitationsservice "github.com/prixfixeco/backend/internal/services/householdinvitations"
	householdsservice "github.com/prixfixeco/backend/internal/services/households"
	mealplaneventsservice "github.com/prixfixeco/backend/internal/services/mealplanevents"
	"github.com/prixfixeco/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/prixfixeco/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/backend/internal/services/mealplans"
	"github.com/prixfixeco/backend/internal/services/mealplantasks"
	mealsservice "github.com/prixfixeco/backend/internal/services/meals"
	"github.com/prixfixeco/backend/internal/services/recipepreptasks"
	recipesservice "github.com/prixfixeco/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/prixfixeco/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/prixfixeco/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/prixfixeco/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/prixfixeco/backend/internal/services/recipesteps"
	usersservice "github.com/prixfixeco/backend/internal/services/users"
	validingredientmeasurementunitsservice "github.com/prixfixeco/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/prixfixeco/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/prixfixeco/backend/internal/services/validingredients"
	"github.com/prixfixeco/backend/internal/services/validingredientstateingredients"
	"github.com/prixfixeco/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/prixfixeco/backend/internal/services/validinstruments"
	"github.com/prixfixeco/backend/internal/services/validmeasurementconversions"
	validmeaurementunitsservice "github.com/prixfixeco/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/prixfixeco/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/prixfixeco/backend/internal/services/validpreparations"
	webhooksservice "github.com/prixfixeco/backend/internal/services/webhooks"
	websocketsservice "github.com/prixfixeco/backend/internal/services/websockets"
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
		Email         emailconfig.Config        `json:"email" mapstructure:"email" toml:"email,omitempty"`
		Analytics     analyticsconfig.Config    `json:"analytics" mapstructure:"analytics" toml:"analytics,omitempty"`
		Encoding      encoding.Config           `json:"encoding" mapstructure:"encoding" toml:"encoding,omitempty"`
		Routing       routing.Config            `json:"routing" mapstructure:"routing" toml:"routing,omitempty"`
		Database      dbconfig.Config           `json:"database" mapstructure:"database" toml:"database,omitempty"`
		Meta          MetaSettings              `json:"meta" mapstructure:"meta" toml:"meta,omitempty"`
		Events        msgconfig.Config          `json:"events" mapstructure:"events" toml:"events,omitempty"`
		FeatureFlags  featureflagsconfig.Config `json:"featureFlags" mapstructure:"events" toml:"events,omitempty"`
		Server        server.Config             `json:"server" mapstructure:"server" toml:"server,omitempty"`
		Services      ServicesConfigurations    `json:"services" mapstructure:"services" toml:"services,omitempty"`
	}

	// ServicesConfigurations collects the various service configurations.
	ServicesConfigurations struct {
		_                               struct{}
		ValidMeasurementUnits           validmeaurementunitsservice.Config            `json:"validMeasurementUnits" mapstructure:"valid_measurement_units" toml:"valid_measurement_units,omitempty"`
		ValidInstruments                validinstrumentsservice.Config                `json:"validInstruments" mapstructure:"valid_instruments" toml:"valid_instruments,omitempty"`
		ValidIngredients                validingredientsservice.Config                `json:"validIngredients" mapstructure:"valid_ingredients" toml:"valid_ingredients,omitempty"`
		ValidPreparations               validpreparationsservice.Config               `json:"validPreparations" mapstructure:"valid_preparations" toml:"valid_preparations,omitempty"`
		MealPlanEvents                  mealplaneventsservice.Config                  `json:"mealPlanEvents" mapstructure:"meal_plan_events" toml:"meal_plan_events,omitempty"`
		MealPlanOptionVotes             mealplanoptionvotesservice.Config             `json:"mealPlanOptionVotes" mapstructure:"meal_plan_option_votes" toml:"meal_plan_option_votes,omitempty"`
		ValidIngredientPreparations     validingredientpreparationsservice.Config     `json:"validIngredientPreparations" mapstructure:"valid_ingredient_preparations" toml:"valid_ingredient_preparations,omitempty"`
		ValidPreparationInstruments     validpreparationinstrumentsservice.Config     `json:"validPreparationInstruments" mapstructure:"valid_preparation_instruments" toml:"valid_preparation_instruments,omitempty"`
		ValidInstrumentMeasurementUnits validingredientmeasurementunitsservice.Config `json:"validInstrumentMeasurementUnits" mapstructure:"valid_ingredient_measurement_units" toml:"valid_ingredient_measurement_units,omitempty"`
		Meals                           mealsservice.Config                           `json:"meals" mapstructure:"meals" toml:"meals,omitempty"`
		Recipes                         recipesservice.Config                         `json:"recipes" mapstructure:"recipes" toml:"recipes,omitempty"`
		RecipeSteps                     recipestepsservice.Config                     `json:"recipeSteps" mapstructure:"recipe_steps" toml:"recipe_steps,omitempty"`
		RecipeStepProducts              recipestepproductsservice.Config              `json:"recipeStepProducts" mapstructure:"recipe_step_products" toml:"recipe_step_products,omitempty"`
		RecipeStepInstruments           recipestepinstrumentsservice.Config           `json:"recipeStepInstruments" mapstructure:"recipe_step_instruments" toml:"recipe_step_instruments,omitempty"`
		RecipeStepIngredients           recipestepingredientsservice.Config           `json:"recipeStepIngredients" mapstructure:"recipe_step_ingredients" toml:"recipe_step_ingredients,omitempty"`
		RecipeStepCompletionConditions  recipestepcompletionconditionsservice.Config  `json:"recipeStepCompletionConditions" mapstructure:"recipe_step_completion_conditions" toml:"recipe_step_completion_conditions,omitempty"`
		MealPlans                       mealplansservice.Config                       `json:"mealPlans" mapstructure:"meal_plans" toml:"meal_plans,omitempty"`
		MealPlanOptions                 mealplanoptionsservice.Config                 `json:"mealPlanOptions" mapstructure:"meal_plan_options" toml:"meal_plan_options,omitempty"`
		Households                      householdsservice.Config                      `json:"households" mapstructure:"households" toml:"households,omitempty"`
		HouseholdInvitations            householdinvitationsservice.Config            `json:"householdInvitations" mapstructure:"household_invitations" toml:"household_invitations,omitempty"`
		Websockets                      websocketsservice.Config                      `json:"websockets" mapstructure:"websockets" toml:"websockets,omitempty"`
		Webhooks                        webhooksservice.Config                        `json:"webhooks" mapstructure:"webhooks" toml:"webhooks,omitempty"`
		Users                           usersservice.Config                           `json:"users" mapstructure:"users" toml:"users,omitempty"`
		MealPlanTasks                   mealplantasks.Config                          `json:"mealPlanTasks" mapstructure:"meal_plan_tasks" toml:"meal_plan_tasks,omitempty"`
		RecipePrepTasks                 recipepreptasks.Config                        `json:"recipePrepTasks" mapstructure:"recipe_prep_tasks" toml:"recipe_prep_tasks,omitempty"`
		MealPlanGroceryListItems        mealplangrocerylistitems.Config               `json:"mealPlanGroceryListItems" mapstructure:"meal_plan_grocery_list_items" toml:"meal_plan_grocery_list_items,omitempty"`
		ValidMeasurementConversions     validmeasurementconversions.Config            `json:"validMeasurementConversions" mapstructure:"valid_measurement_conversions" toml:"valid_measurement_conversions,omitempty"`
		ValidIngredientStates           validingredientstates.Config                  `json:"validIngredientStates" mapstructure:"valid_ingredient_states" toml:"valid_ingredient_states,omitempty"`
		ValidIngredientStateIngredients validingredientstateingredients.Config        `json:"validIngredientStateIngredients" mapstructure:"valid_ingredient_state_ingredients" toml:"valid_ingredient_state_ingredients,omitempty"`
		Auth                            authservice.Config                            `json:"auth" mapstructure:"auth" toml:"auth,omitempty"`
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

	if err := cfg.Routing.ValidateWithContext(ctx); err != nil {
		result = multierror.Append(fmt.Errorf("error validating Routing portion of config: %w", err), result)
	}

	if err := cfg.Meta.ValidateWithContext(ctx); err != nil {
		result = multierror.Append(fmt.Errorf("error validating Meta portion of config: %w", err), result)
	}

	if err := cfg.Encoding.ValidateWithContext(ctx); err != nil {
		result = multierror.Append(fmt.Errorf("error validating Encoding portion of config: %w", err), result)
	}

	if err := cfg.Analytics.ValidateWithContext(ctx); err != nil {
		result = multierror.Append(fmt.Errorf("error validating Analytics portion of config: %w", err), result)
	}

	if err := cfg.Observability.ValidateWithContext(ctx); err != nil {
		result = multierror.Append(fmt.Errorf("error validating Observability portion of config: %w", err), result)
	}

	if err := cfg.Database.ValidateWithContext(ctx); err != nil {
		result = multierror.Append(fmt.Errorf("error validating Database portion of config: %w", err), result)
	}

	if err := cfg.Server.ValidateWithContext(ctx); err != nil {
		result = multierror.Append(fmt.Errorf("error validating Server portion of config: %w", err), result)
	}

	if err := cfg.Email.ValidateWithContext(ctx); err != nil {
		result = multierror.Append(fmt.Errorf("error validating Email portion of config: %w", err), result)
	}

	if err := cfg.FeatureFlags.ValidateWithContext(ctx); err != nil {
		result = multierror.Append(fmt.Errorf("error validating FeatureFlags portion of config: %w", err), result)
	}

	if validateServices {
		if err := cfg.Services.Auth.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating Auth service portion of config: %w", err), result)
		}

		if err := cfg.Services.Users.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating Auth service portion of config: %w", err), result)
		}

		if err := cfg.Services.Webhooks.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating Webhooks service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidInstruments.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidInstruments service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidIngredients.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidIngredients service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidPreparations.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidPreparations service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidMeasurementUnits.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidMeasurementUnits service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidIngredientPreparations.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidIngredientPreparations service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidIngredientStateIngredients.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidIngredientStateIngredients service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidPreparationInstruments.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidPreparationInstruments service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidInstrumentMeasurementUnits.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidInstrumentMeasurementUnits service portion of config: %w", err), result)
		}

		if err := cfg.Services.Recipes.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating Components service portion of config: %w", err), result)
		}

		if err := cfg.Services.RecipeSteps.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating RecipeSteps service portion of config: %w", err), result)
		}

		if err := cfg.Services.RecipeStepInstruments.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating RecipeStepInstruments service portion of config: %w", err), result)
		}

		if err := cfg.Services.RecipeStepIngredients.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating RecipeStepIngredients service portion of config: %w", err), result)
		}

		if err := cfg.Services.RecipeStepCompletionConditions.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating RecipeStepCompletionConditions service portion of config: %w", err), result)
		}

		if err := cfg.Services.MealPlans.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating MealPlans service portion of config: %w", err), result)
		}

		if err := cfg.Services.MealPlanEvents.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating MealPlanEvents service portion of config: %w", err), result)
		}

		if err := cfg.Services.MealPlanOptions.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating MealPlanOptions service portion of config: %w", err), result)
		}

		if err := cfg.Services.MealPlanOptionVotes.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating MealPlanOptionVotes service portion of config: %w", err), result)
		}

		if err := cfg.Services.RecipePrepTasks.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating RecipePrepTasks service portion of config: %w", err), result)
		}

		if err := cfg.Services.MealPlanGroceryListItems.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating MealPlanGroceryListItems service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidMeasurementConversions.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidMeasurementConversions service portion of config: %w", err), result)
		}

		if err := cfg.Services.ValidIngredientStates.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating ValidIngredientStates service portion of config: %w", err), result)
		}
	}

	if result != nil {
		return result
	}

	return nil
}
