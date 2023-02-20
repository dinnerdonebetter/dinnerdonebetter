package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidIngredientDataType indicates an event is related to a valid ingredient.
	ValidIngredientDataType dataType = "valid_ingredient"

	// ValidIngredientCreatedCustomerEventType indicates a valid ingredient was created.
	ValidIngredientCreatedCustomerEventType CustomerEventType = "valid_ingredient_created"
	// ValidIngredientUpdatedCustomerEventType indicates a valid ingredient was updated.
	ValidIngredientUpdatedCustomerEventType CustomerEventType = "valid_ingredient_updated"
	// ValidIngredientArchivedCustomerEventType indicates a valid ingredient was archived.
	ValidIngredientArchivedCustomerEventType CustomerEventType = "valid_ingredient_archived"
)

func init() {
	gob.Register(new(ValidIngredient))
	gob.Register(new(ValidIngredientCreationRequestInput))
	gob.Register(new(ValidIngredientUpdateRequestInput))
}

type (
	// ValidIngredient represents a valid ingredient.
	ValidIngredient struct {
		_ struct{}

		CreatedAt                               time.Time  `json:"createdAt"`
		LastUpdatedAt                           *time.Time `json:"lastUpdatedAt"`
		ArchivedAt                              *time.Time `json:"archivedAt"`
		MaximumIdealStorageTemperatureInCelsius *float32   `json:"maximumIdealStorageTemperatureInCelsius"`
		MinimumIdealStorageTemperatureInCelsius *float32   `json:"minimumIdealStorageTemperatureInCelsius"`
		IconPath                                string     `json:"iconPath"`
		Warning                                 string     `json:"warning"`
		PluralName                              string     `json:"pluralName"`
		StorageInstructions                     string     `json:"storageInstructions"`
		Name                                    string     `json:"name"`
		ID                                      string     `json:"id"`
		Description                             string     `json:"description"`
		Slug                                    string     `json:"slug"`
		ShoppingSuggestions                     string     `json:"shoppingSuggestions"`
		ContainsShellfish                       bool       `json:"containsShellfish"`
		IsMeasuredVolumetrically                bool       `json:"isMeasuredVolumetrically"`
		IsLiquid                                bool       `json:"isLiquid"`
		ContainsPeanut                          bool       `json:"containsPeanut"`
		ContainsTreeNut                         bool       `json:"containsTreeNut"`
		ContainsEgg                             bool       `json:"containsEgg"`
		ContainsWheat                           bool       `json:"containsWheat"`
		ContainsSoy                             bool       `json:"containsSoy"`
		AnimalDerived                           bool       `json:"animalDerived"`
		RestrictToPreparations                  bool       `json:"restrictToPreparations"`
		ContainsSesame                          bool       `json:"containsSesame"`
		ContainsFish                            bool       `json:"containsFish"`
		ContainsGluten                          bool       `json:"containsGluten"`
		ContainsDairy                           bool       `json:"containsDairy"`
		ContainsAlcohol                         bool       `json:"containsAlcohol"`
		AnimalFlesh                             bool       `json:"animalFlesh"`
	}

	// NullableValidIngredient represents a nullable valid ingredient.
	NullableValidIngredient struct {
		_ struct{}

		CreatedAt                               *time.Time
		LastUpdatedAt                           *time.Time
		ArchivedAt                              *time.Time
		ID                                      *string
		Warning                                 *string
		Description                             *string
		IconPath                                *string
		PluralName                              *string
		StorageInstructions                     *string
		Name                                    *string
		MaximumIdealStorageTemperatureInCelsius *float32
		MinimumIdealStorageTemperatureInCelsius *float32
		ContainsShellfish                       *bool
		ContainsDairy                           *bool
		AnimalFlesh                             *bool
		IsMeasuredVolumetrically                *bool
		IsLiquid                                *bool
		ContainsPeanut                          *bool
		ContainsTreeNut                         *bool
		ContainsEgg                             *bool
		ContainsWheat                           *bool
		ContainsSoy                             *bool
		AnimalDerived                           *bool
		RestrictToPreparations                  *bool
		ContainsSesame                          *bool
		ContainsFish                            *bool
		ContainsGluten                          *bool
		Slug                                    *string
		ContainsAlcohol                         *bool
		ShoppingSuggestions                     *string
	}

	// ValidIngredientCreationRequestInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientCreationRequestInput struct {
		_ struct{}

		MinimumIdealStorageTemperatureInCelsius *float32 `json:"minimumIdealStorageTemperatureInCelsius"`
		MaximumIdealStorageTemperatureInCelsius *float32 `json:"maximumIdealStorageTemperatureInCelsius"`
		Warning                                 string   `json:"warning"`
		IconPath                                string   `json:"iconPath"`
		PluralName                              string   `json:"pluralName"`
		StorageInstructions                     string   `json:"storageInstructions"`
		Name                                    string   `json:"name"`
		Description                             string   `json:"description"`
		Slug                                    string   `json:"slug"`
		ShoppingSuggestions                     string   `json:"shoppingSuggestions"`
		IsMeasuredVolumetrically                bool     `json:"isMeasuredVolumetrically"`
		ContainsFish                            bool     `json:"containsFish"`
		ContainsShellfish                       bool     `json:"containsShellfish"`
		AnimalFlesh                             bool     `json:"animalFlesh"`
		ContainsEgg                             bool     `json:"containsEgg"`
		IsLiquid                                bool     `json:"isLiquid"`
		ContainsSoy                             bool     `json:"containsSoy"`
		ContainsPeanut                          bool     `json:"containsPeanut"`
		AnimalDerived                           bool     `json:"animalDerived"`
		RestrictToPreparations                  bool     `json:"restrictToPreparations"`
		ContainsDairy                           bool     `json:"containsDairy"`
		ContainsSesame                          bool     `json:"containsSesame"`
		ContainsTreeNut                         bool     `json:"containsTreeNut"`
		ContainsWheat                           bool     `json:"containsWheat"`
		ContainsAlcohol                         bool     `json:"containsAlcohol"`
		ContainsGluten                          bool     `json:"containsGluten"`
	}

	// ValidIngredientDatabaseCreationInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientDatabaseCreationInput struct {
		_ struct{}

		MinimumIdealStorageTemperatureInCelsius *float32
		MaximumIdealStorageTemperatureInCelsius *float32
		ID                                      string
		Warning                                 string
		IconPath                                string
		PluralName                              string
		StorageInstructions                     string
		Name                                    string
		Description                             string
		Slug                                    string
		ShoppingSuggestions                     string
		IsMeasuredVolumetrically                bool
		ContainsFish                            bool
		ContainsShellfish                       bool
		AnimalFlesh                             bool
		ContainsEgg                             bool
		IsLiquid                                bool
		ContainsSoy                             bool
		ContainsPeanut                          bool
		AnimalDerived                           bool
		RestrictToPreparations                  bool
		ContainsDairy                           bool
		ContainsSesame                          bool
		ContainsTreeNut                         bool
		ContainsWheat                           bool
		ContainsAlcohol                         bool
		ContainsGluten                          bool
	}

	// ValidIngredientUpdateRequestInput represents what a user could set as input for updating valid ingredients.
	ValidIngredientUpdateRequestInput struct {
		_ struct{}

		Name                                    *string  `json:"name"`
		Description                             *string  `json:"description"`
		Warning                                 *string  `json:"warning"`
		IconPath                                *string  `json:"iconPath"`
		ContainsDairy                           *bool    `json:"containsDairy"`
		ContainsPeanut                          *bool    `json:"containsPeanut"`
		ContainsTreeNut                         *bool    `json:"containsTreeNut"`
		ContainsEgg                             *bool    `json:"containsEgg"`
		ContainsWheat                           *bool    `json:"containsWheat"`
		ContainsShellfish                       *bool    `json:"containsShellfish"`
		ContainsSesame                          *bool    `json:"containsSesame"`
		ContainsFish                            *bool    `json:"containsFish"`
		ContainsGluten                          *bool    `json:"containsGluten"`
		AnimalFlesh                             *bool    `json:"animalFlesh"`
		IsMeasuredVolumetrically                *bool    `json:"isMeasuredVolumetrically"`
		IsLiquid                                *bool    `json:"isLiquid"`
		ContainsSoy                             *bool    `json:"containsSoy"`
		PluralName                              *string  `json:"pluralName"`
		AnimalDerived                           *bool    `json:"animalDerived"`
		RestrictToPreparations                  *bool    `json:"restrictToPreparations"`
		MinimumIdealStorageTemperatureInCelsius *float32 `json:"minimumIdealStorageTemperatureInCelsius"`
		MaximumIdealStorageTemperatureInCelsius *float32 `json:"maximumIdealStorageTemperatureInCelsius"`
		StorageInstructions                     *string  `json:"storageInstructions"`
		Slug                                    *string  `json:"slug"`
		ContainsAlcohol                         *bool    `json:"containsAlcohol"`
		ShoppingSuggestions                     *string  `json:"shoppingSuggestions"`
	}

	// ValidIngredientDataManager describes a structure capable of storing valid ingredients permanently.
	ValidIngredientDataManager interface {
		ValidIngredientExists(ctx context.Context, validIngredientID string) (bool, error)
		GetValidIngredient(ctx context.Context, validIngredientID string) (*ValidIngredient, error)
		GetRandomValidIngredient(ctx context.Context) (*ValidIngredient, error)
		GetValidIngredients(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidIngredient], error)
		SearchForValidIngredients(ctx context.Context, query string, filter *QueryFilter) ([]*ValidIngredient, error)
		SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string, filter *QueryFilter) ([]*ValidIngredient, error)
		SearchForValidIngredientsForIngredientState(ctx context.Context, ingredientStateID, query string, filter *QueryFilter) ([]*ValidIngredient, error)
		CreateValidIngredient(ctx context.Context, input *ValidIngredientDatabaseCreationInput) (*ValidIngredient, error)
		UpdateValidIngredient(ctx context.Context, updated *ValidIngredient) error
		ArchiveValidIngredient(ctx context.Context, validIngredientID string) error
	}

	// ValidIngredientDataService describes a structure capable of serving traffic related to valid ingredients.
	ValidIngredientDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		ForValidIngredientStateHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		RandomHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidIngredientUpdateRequestInput with a valid ingredient.
func (x *ValidIngredient) Update(input *ValidIngredientUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.Warning != nil && *input.Warning != x.Warning {
		x.Warning = *input.Warning
	}

	if input.ContainsEgg != nil && *input.ContainsEgg != x.ContainsEgg {
		x.ContainsEgg = *input.ContainsEgg
	}

	if input.ContainsDairy != nil && *input.ContainsDairy != x.ContainsDairy {
		x.ContainsDairy = *input.ContainsDairy
	}

	if input.ContainsPeanut != nil && *input.ContainsPeanut != x.ContainsPeanut {
		x.ContainsPeanut = *input.ContainsPeanut
	}

	if input.ContainsTreeNut != nil && *input.ContainsTreeNut != x.ContainsTreeNut {
		x.ContainsTreeNut = *input.ContainsTreeNut
	}

	if input.ContainsSoy != nil && *input.ContainsSoy != x.ContainsSoy {
		x.ContainsSoy = *input.ContainsSoy
	}

	if input.ContainsWheat != nil && *input.ContainsWheat != x.ContainsWheat {
		x.ContainsWheat = *input.ContainsWheat
	}

	if input.ContainsShellfish != nil && *input.ContainsShellfish != x.ContainsShellfish {
		x.ContainsShellfish = *input.ContainsShellfish
	}

	if input.ContainsSesame != nil && *input.ContainsSesame != x.ContainsSesame {
		x.ContainsSesame = *input.ContainsSesame
	}

	if input.ContainsFish != nil && *input.ContainsFish != x.ContainsFish {
		x.ContainsFish = *input.ContainsFish
	}

	if input.ContainsGluten != nil && *input.ContainsGluten != x.ContainsGluten {
		x.ContainsGluten = *input.ContainsGluten
	}

	if input.AnimalFlesh != nil && *input.AnimalFlesh != x.AnimalFlesh {
		x.AnimalFlesh = *input.AnimalFlesh
	}

	if input.IsMeasuredVolumetrically != nil && *input.IsMeasuredVolumetrically != x.IsMeasuredVolumetrically {
		x.IsMeasuredVolumetrically = *input.IsMeasuredVolumetrically
	}

	if input.IsLiquid != nil && *input.IsLiquid != x.IsLiquid {
		x.IsLiquid = *input.IsLiquid
	}

	if input.IconPath != nil && *input.IconPath != x.IconPath {
		x.IconPath = *input.IconPath
	}

	if input.PluralName != nil && *input.PluralName != x.PluralName {
		x.PluralName = *input.PluralName
	}

	if input.AnimalDerived != nil && *input.AnimalDerived != x.AnimalDerived {
		x.AnimalDerived = *input.AnimalDerived
	}

	if input.RestrictToPreparations != nil && *input.RestrictToPreparations != x.RestrictToPreparations {
		x.RestrictToPreparations = *input.RestrictToPreparations
	}

	if input.MinimumIdealStorageTemperatureInCelsius != nil && input.MinimumIdealStorageTemperatureInCelsius != x.MinimumIdealStorageTemperatureInCelsius {
		x.MinimumIdealStorageTemperatureInCelsius = input.MinimumIdealStorageTemperatureInCelsius
	}

	if input.MaximumIdealStorageTemperatureInCelsius != nil && input.MaximumIdealStorageTemperatureInCelsius != x.MaximumIdealStorageTemperatureInCelsius {
		x.MaximumIdealStorageTemperatureInCelsius = input.MaximumIdealStorageTemperatureInCelsius
	}

	if input.StorageInstructions != nil && *input.StorageInstructions != x.StorageInstructions {
		x.StorageInstructions = *input.StorageInstructions
	}

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
	}

	if input.ContainsAlcohol != nil && *input.ContainsAlcohol != x.ContainsAlcohol {
		x.ContainsAlcohol = *input.ContainsAlcohol
	}

	if input.ShoppingSuggestions != nil && *input.ShoppingSuggestions != x.ShoppingSuggestions {
		x.ShoppingSuggestions = *input.ShoppingSuggestions
	}
}

var _ validation.ValidatableWithContext = (*ValidIngredientCreationRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientCreationRequestInput.
func (x *ValidIngredientCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidIngredientDatabaseCreationInput.
func (x *ValidIngredientDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientUpdateRequestInput.
func (x *ValidIngredientUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
