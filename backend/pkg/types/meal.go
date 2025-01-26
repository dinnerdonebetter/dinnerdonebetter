package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

const (
	// MealComponentTypesUnspecified represents the unspecified meal component type.
	MealComponentTypesUnspecified = "unspecified"
	// MealComponentTypesAmuseBouche represents the amuse-bouche meal component type.
	MealComponentTypesAmuseBouche = "amuse-bouche"
	// MealComponentTypesAppetizer represents the appetizer meal component type.
	MealComponentTypesAppetizer = "appetizer"
	// MealComponentTypesSoup represents the soup meal component type.
	MealComponentTypesSoup = "soup"
	// MealComponentTypesMain represents the main meal component type.
	MealComponentTypesMain = "main"
	// MealComponentTypesSalad represents the salad meal component type.
	MealComponentTypesSalad = "salad"
	// MealComponentTypesBeverage represents the beverage meal component type.
	MealComponentTypesBeverage = "beverage"
	// MealComponentTypesSide represents the side meal component type.
	MealComponentTypesSide = "side"
	// MealComponentTypesDessert represents the dessert meal component type.
	MealComponentTypesDessert = "dessert"

	// MealCreatedServiceEventType indicates a meal was created.
	MealCreatedServiceEventType = "meal_created"
	// MealUpdatedServiceEventType indicates a meal was updated.
	MealUpdatedServiceEventType = "meal_updated"
	// MealArchivedServiceEventType indicates a meal was archived.
	MealArchivedServiceEventType = "meal_archived"
)

func init() {
	gob.Register(new(Meal))
	gob.Register(new(MealCreationRequestInput))
	gob.Register(new(MealUpdateRequestInput))
}

type (
	// Meal represents a meal.
	Meal struct {
		_ struct{} `json:"-"`

		CreatedAt            time.Time                   `json:"createdAt"`
		ArchivedAt           *time.Time                  `json:"archivedAt"`
		LastUpdatedAt        *time.Time                  `json:"lastUpdatedAt"`
		EstimatedPortions    Float32RangeWithOptionalMax `json:"estimatedPortions"`
		ID                   string                      `json:"id"`
		Description          string                      `json:"description"`
		CreatedByUser        string                      `json:"createdByUser"`
		Name                 string                      `json:"name"`
		Components           []*MealComponent            `json:"components"`
		EligibleForMealPlans bool                        `json:"eligibleForMealPlans"`
	}

	// MealComponent is a recipe with some extra data attached to it.
	MealComponent struct {
		_ struct{} `json:"-"`

		ComponentType string  `json:"componentType"`
		Recipe        Recipe  `json:"recipe"`
		RecipeScale   float32 `json:"recipeScale"`
	}

	// MealCreationRequestInput represents what a user could set as input for creating meals.
	MealCreationRequestInput struct {
		_ struct{} `json:"-"`

		EstimatedPortions    Float32RangeWithOptionalMax          `json:"estimatedPortions"`
		Name                 string                               `json:"name"`
		Description          string                               `json:"description"`
		Components           []*MealComponentCreationRequestInput `json:"components"`
		EligibleForMealPlans bool                                 `json:"eligibleForMealPlans"`
	}

	// MealComponentCreationRequestInput represents what a user could set as input for creating meal recipes.
	MealComponentCreationRequestInput struct {
		_ struct{} `json:"-"`

		RecipeID      string  `json:"recipeID"`
		ComponentType string  `json:"componentType"`
		RecipeScale   float32 `json:"recipeScale"`
	}

	// MealDatabaseCreationInput represents what a user could set as input for creating meals.
	MealDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		EstimatedPortions    Float32RangeWithOptionalMax           `json:"-"`
		ID                   string                                `json:"-"`
		Name                 string                                `json:"-"`
		Description          string                                `json:"-"`
		CreatedByUser        string                                `json:"-"`
		Components           []*MealComponentDatabaseCreationInput `json:"-"`
		EligibleForMealPlans bool                                  `json:"-"`
	}

	// MealComponentDatabaseCreationInput represents what a user could set as input for creating meal recipes.
	MealComponentDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		RecipeID      string  `json:"-"`
		ComponentType string  `json:"-"`
		RecipeScale   float32 `json:"-"`
	}

	// MealUpdateRequestInput represents what a user could set as input for updating meals.
	MealUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name                 *string                                       `json:"name,omitempty"`
		Description          *string                                       `json:"description,omitempty"`
		CreatedByUser        *string                                       `json:"-"`
		EstimatedPortions    Float32RangeWithOptionalMaxUpdateRequestInput `json:"estimatedPortions"`
		EligibleForMealPlans *bool                                         `json:"eligibleForMealPlans"`
		Components           []*MealComponentUpdateRequestInput            `json:"recipes,omitempty"`
	}

	// MealComponentUpdateRequestInput represents what a user could set as input for creating meal recipes.
	MealComponentUpdateRequestInput struct {
		_ struct{} `json:"-"`

		RecipeID      *string  `json:"recipeID"`
		ComponentType *string  `json:"componentType"`
		RecipeScale   *float32 `json:"recipeScale"`
	}

	// MealDataManager describes a structure capable of storing meals permanently.
	MealDataManager interface {
		MealExists(ctx context.Context, mealID string) (bool, error)
		GetMeal(ctx context.Context, mealID string) (*Meal, error)
		GetMeals(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[Meal], error)
		GetMealsCreatedByUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[Meal], error)
		SearchForMeals(ctx context.Context, query string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[Meal], error)
		CreateMeal(ctx context.Context, input *MealDatabaseCreationInput) (*Meal, error)
		MarkMealAsIndexed(ctx context.Context, mealID string) error
		ArchiveMeal(ctx context.Context, mealID, userID string) error
		GetMealIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error)
		GetMealsWithIDs(ctx context.Context, ids []string) ([]*Meal, error)
	}

	// MealDataService describes a structure capable of serving traffic related to meals.
	MealDataService interface {
		ListMealsHandler(http.ResponseWriter, *http.Request)
		CreateMealHandler(http.ResponseWriter, *http.Request)
		ReadMealHandler(http.ResponseWriter, *http.Request)
		SearchMealsHandler(http.ResponseWriter, *http.Request)
		ArchiveMealHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an MealUpdateRequestInput with a meal.
func (x *Meal) Update(input *MealUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.EstimatedPortions.Min != nil && *input.EstimatedPortions.Min != x.EstimatedPortions.Min {
		x.EstimatedPortions.Min = *input.EstimatedPortions.Min
	}

	if input.EstimatedPortions.Max != nil && input.EstimatedPortions.Max != x.EstimatedPortions.Max {
		x.EstimatedPortions.Max = input.EstimatedPortions.Max
	}

	if input.EligibleForMealPlans != nil && *input.EligibleForMealPlans != x.EligibleForMealPlans {
		x.EligibleForMealPlans = *input.EligibleForMealPlans
	}
}

// Update merges an MealComponentUpdateRequestInput with a meal.
func (x *MealComponent) Update(input *MealComponentUpdateRequestInput) {
	if input.RecipeID != nil && *input.RecipeID != x.Recipe.ID {
		x.Recipe = Recipe{ID: *input.RecipeID}
	}

	if input.ComponentType != nil && *input.ComponentType != x.ComponentType {
		x.ComponentType = *input.ComponentType
	}

	if input.RecipeScale != nil && *input.RecipeScale != x.RecipeScale {
		x.RecipeScale = *input.RecipeScale
	}
}

var _ validation.ValidatableWithContext = (*MealCreationRequestInput)(nil)

// ValidateWithContext validates a MealCreationRequestInput.
func (x *MealCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	atLeastOneMain := false
	for _, component := range x.Components {
		if component.ComponentType == MealComponentTypesMain {
			atLeastOneMain = true
		}

		if componentValidationErr := component.ValidateWithContext(ctx); componentValidationErr != nil {
			result = multierror.Append(result, componentValidationErr)
		}
	}

	if !atLeastOneMain {
		result = multierror.Append(result, errOneMainMinimumRequired)
	}

	if validationErr := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Components, validation.Required),
	); validationErr != nil {
		result = multierror.Append(result, validationErr)
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*MealCreationRequestInput)(nil)

// ValidateWithContext validates a MealComponentCreationRequestInput.
func (x *MealComponentCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ComponentType,
			validation.Required,
			validation.In(
				MealComponentTypesUnspecified,
				MealComponentTypesAmuseBouche,
				MealComponentTypesAppetizer,
				MealComponentTypesSoup,
				MealComponentTypesMain,
				MealComponentTypesSalad,
				MealComponentTypesBeverage,
				MealComponentTypesSide,
				MealComponentTypesDessert,
			),
		),
	)
}

var _ validation.ValidatableWithContext = (*MealDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealDatabaseCreationInput.
func (x *MealDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Components, validation.Required),
		validation.Field(&x.CreatedByUser, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealUpdateRequestInput)(nil)

// ValidateWithContext validates a MealUpdateRequestInput.
func (x *MealUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.Components, validation.Required),
		validation.Field(&x.CreatedByUser, validation.Required),
	)
}
