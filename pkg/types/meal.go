package types

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

const (
	// MealDataType indicates an event is related to a meal.
	MealDataType dataType = "meal"

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

	// MealCreatedCustomerEventType indicates a meal was created.
	MealCreatedCustomerEventType CustomerEventType = "meal_created"
	// MealUpdatedCustomerEventType indicates a meal was updated.
	MealUpdatedCustomerEventType CustomerEventType = "meal_updated"
	// MealArchivedCustomerEventType indicates a meal was archived.
	MealArchivedCustomerEventType CustomerEventType = "meal_archived"
)

func init() {
	gob.Register(new(Meal))
	gob.Register(new(MealList))
	gob.Register(new(MealCreationRequestInput))
	gob.Register(new(MealUpdateRequestInput))
}

type (
	// Meal represents a meal.
	Meal struct {
		_             struct{}
		CreatedAt     time.Time        `json:"createdAt"`
		ArchivedAt    *time.Time       `json:"archivedAt"`
		LastUpdatedAt *time.Time       `json:"lastUpdatedAt"`
		ID            string           `json:"id"`
		Description   string           `json:"description"`
		CreatedByUser string           `json:"createdByUser"`
		Name          string           `json:"name"`
		Components    []*MealComponent `json:"components"`
	}

	// MealComponent is a recipe with some extra data attached to it.
	MealComponent struct {
		ComponentType string `json:"componentType"`
		Recipe        Recipe `json:"recipe"`
	}

	// MealList represents a list of meals.
	MealList struct {
		_ struct{}

		Meals []*Meal `json:"data"`
		Pagination
	}

	// MealCreationRequestInput represents what a user could set as input for creating meals.
	MealCreationRequestInput struct {
		_ struct{}

		ID            string                               `json:"-"`
		Name          string                               `json:"name"`
		Description   string                               `json:"description"`
		CreatedByUser string                               `json:"-"`
		Components    []*MealComponentCreationRequestInput `json:"recipes"`
	}

	// MealComponentCreationRequestInput represents what a user could set as input for creating meal recipes.
	MealComponentCreationRequestInput struct {
		RecipeID      string `json:"recipeID"`
		ComponentType string `json:"mealComponentType"`
	}

	// MealDatabaseCreationInput represents what a user could set as input for creating meals.
	MealDatabaseCreationInput struct {
		_ struct{}

		ID            string                                `json:"id"`
		Name          string                                `json:"name"`
		Description   string                                `json:"description"`
		CreatedByUser string                                `json:"belongsToHousehold"`
		Components    []*MealComponentDatabaseCreationInput `json:"recipes"`
	}

	// MealComponentDatabaseCreationInput represents what a user could set as input for creating meal recipes.
	MealComponentDatabaseCreationInput struct {
		RecipeID      string `json:"recipeID"`
		ComponentType string `json:"mealComponentType"`
	}

	// MealUpdateRequestInput represents what a user could set as input for updating meals.
	MealUpdateRequestInput struct {
		_             struct{}
		Name          *string                            `json:"name"`
		Description   *string                            `json:"description"`
		CreatedByUser *string                            `json:"-"`
		Components    []*MealComponentUpdateRequestInput `json:"recipes"`
	}

	// MealComponentUpdateRequestInput represents what a user could set as input for creating meal recipes.
	MealComponentUpdateRequestInput struct {
		RecipeID      string `json:"recipeID"`
		ComponentType string `json:"mealComponentType"`
	}

	// MealDataManager describes a structure capable of storing meals permanently.
	MealDataManager interface {
		MealExists(ctx context.Context, mealID string) (bool, error)
		GetMeal(ctx context.Context, mealID string) (*Meal, error)
		GetMeals(ctx context.Context, filter *QueryFilter) (*MealList, error)
		SearchForMeals(ctx context.Context, query string, filter *QueryFilter) (*MealList, error)
		CreateMeal(ctx context.Context, input *MealDatabaseCreationInput) (*Meal, error)
		ArchiveMeal(ctx context.Context, mealID, userID string) error
	}

	// MealDataService describes a structure capable of serving traffic related to meals.
	MealDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		SearchHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
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
}

var errOneMainMinimumRequired = errors.New("at least one main required for meal creation")

var _ validation.ValidatableWithContext = (*MealCreationRequestInput)(nil)

// ValidateWithContext validates a MealCreationRequestInput.
func (x *MealCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	var result *multierror.Error

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

	if result != nil {
		return result
	}

	return nil
}

var _ validation.ValidatableWithContext = (*MealCreationRequestInput)(nil)

// ValidateWithContext validates a MealComponentCreationRequestInput.
func (x *MealComponentCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ComponentType, validation.In(
			MealComponentTypesUnspecified,
			MealComponentTypesAmuseBouche,
			MealComponentTypesAppetizer,
			MealComponentTypesSoup,
			MealComponentTypesMain,
			MealComponentTypesSalad,
			MealComponentTypesBeverage,
			MealComponentTypesSide,
			MealComponentTypesDessert,
		)),
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
