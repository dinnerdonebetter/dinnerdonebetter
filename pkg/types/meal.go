package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealDataType indicates an event is related to a meal.
	MealDataType dataType = "meal"

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
		ArchivedOn    *uint64   `json:"archivedOn"`
		LastUpdatedOn *uint64   `json:"lastUpdatedOn"`
		ID            string    `json:"id"`
		Description   string    `json:"description"`
		CreatedByUser string    `json:"createdByUser"`
		Name          string    `json:"name"`
		Recipes       []*Recipe `json:"recipes"`
		CreatedOn     uint64    `json:"createdOn"`
	}

	// MealRecipe is a recipe with some extra data attached to it.
	MealRecipe struct {
		Recipe        *Recipe `json:"recipe"`
		ComponentType string  `json:"componentType"`
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

		ID            string   `json:"-"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		CreatedByUser string   `json:"-"`
		Recipes       []string `json:"recipes"`
	}

	// MealDatabaseCreationInput represents what a user could set as input for creating meals.
	MealDatabaseCreationInput struct {
		_ struct{}

		ID            string   `json:"id"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		CreatedByUser string   `json:"belongsToHousehold"`
		Recipes       []string `json:"recipes"`
	}

	// MealUpdateRequestInput represents what a user could set as input for updating meals.
	MealUpdateRequestInput struct {
		_             struct{}
		Name          *string  `json:"name"`
		Description   *string  `json:"description"`
		CreatedByUser *string  `json:"-"`
		Recipes       []string `json:"recipes"`
	}

	// MealDataManager describes a structure capable of storing meals permanently.
	MealDataManager interface {
		MealExists(ctx context.Context, mealID string) (bool, error)
		GetMeal(ctx context.Context, mealID string) (*Meal, error)
		GetTotalMealCount(ctx context.Context) (uint64, error)
		GetMeals(ctx context.Context, filter *QueryFilter) (*MealList, error)
		GetMealsWithIDs(ctx context.Context, userID string, limit uint8, ids []string) ([]*Meal, error)
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

var _ validation.ValidatableWithContext = (*MealCreationRequestInput)(nil)

// ValidateWithContext validates a MealCreationRequestInput.
func (x *MealCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Recipes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealDatabaseCreationInput.
func (x *MealDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Recipes, validation.Required),
		validation.Field(&x.CreatedByUser, validation.Required),
	)
}

// MealDatabaseCreationInputFromMealCreationInput creates a DatabaseCreationInput from a CreationInput.
func MealDatabaseCreationInputFromMealCreationInput(input *MealCreationRequestInput) *MealDatabaseCreationInput {
	x := &MealDatabaseCreationInput{
		ID:            input.ID,
		Name:          input.Name,
		Description:   input.Description,
		CreatedByUser: input.CreatedByUser,
		Recipes:       input.Recipes,
	}

	return x
}

var _ validation.ValidatableWithContext = (*MealUpdateRequestInput)(nil)

// ValidateWithContext validates a MealUpdateRequestInput.
func (x *MealUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.Recipes, validation.Required),
		validation.Field(&x.CreatedByUser, validation.Required),
	)
}
