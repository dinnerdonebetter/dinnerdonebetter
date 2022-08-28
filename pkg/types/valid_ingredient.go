package types

import (
	"context"
	"encoding/gob"
	"net/http"

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
	gob.Register(new(ValidIngredientList))
	gob.Register(new(ValidIngredientCreationRequestInput))
	gob.Register(new(ValidIngredientUpdateRequestInput))
}

type (
	// ValidIngredient represents a valid ingredient.
	ValidIngredient struct {
		_                                       struct{}
		LastUpdatedAt                           *uint64 `json:"lastUpdatedAt"`
		ArchivedAt                              *uint64 `json:"archivedAt"`
		Name                                    string  `json:"name"`
		Description                             string  `json:"description"`
		Warning                                 string  `json:"warning"`
		ID                                      string  `json:"id"`
		IconPath                                string  `json:"iconPath"`
		PluralName                              string  `json:"pluralName"`
		CreatedAt                               uint64  `json:"createdAt"`
		MinimumIdealStorageTemperatureInCelsius float32 `json:"minimumIdealStorageTemperatureInCelsius"`
		MaximumIdealStorageTemperatureInCelsius float32 `json:"maximumIdealStorageTemperatureInCelsius"`
		ContainsTreeNut                         bool    `json:"containsTreeNut"`
		ContainsShellfish                       bool    `json:"containsShellfish"`
		ContainsGluten                          bool    `json:"containsGluten"`
		AnimalFlesh                             bool    `json:"animalFlesh"`
		IsMeasuredVolumetrically                bool    `json:"isMeasuredVolumetrically"`
		IsLiquid                                bool    `json:"isLiquid"`
		ContainsPeanut                          bool    `json:"containsPeanut"`
		ContainsDairy                           bool    `json:"containsDairy"`
		ContainsEgg                             bool    `json:"containsEgg"`
		ContainsWheat                           bool    `json:"containsWheat"`
		ContainsSoy                             bool    `json:"containsSoy"`
		AnimalDerived                           bool    `json:"animalDerived"`
		RestrictToPreparations                  bool    `json:"restrictToPreparations"`
		ContainsSesame                          bool    `json:"containsSesame"`
		ContainsFish                            bool    `json:"containsFish"`
	}

	// ValidIngredientList represents a list of valid ingredients.
	ValidIngredientList struct {
		_                struct{}
		ValidIngredients []*ValidIngredient `json:"data"`
		Pagination
	}

	// ValidIngredientCreationRequestInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientCreationRequestInput struct {
		_                                       struct{}
		ID                                      string  `json:"-"`
		Name                                    string  `json:"name"`
		Description                             string  `json:"description"`
		Warning                                 string  `json:"warning"`
		IconPath                                string  `json:"iconPath"`
		PluralName                              string  `json:"pluralName"`
		MinimumIdealStorageTemperatureInCelsius float32 `json:"minimumIdealStorageTemperatureInCelsius"`
		MaximumIdealStorageTemperatureInCelsius float32 `json:"maximumIdealStorageTemperatureInCelsius"`
		ContainsEgg                             bool    `json:"containsEgg"`
		ContainsWheat                           bool    `json:"containsWheat"`
		ContainsShellfish                       bool    `json:"containsShellfish"`
		ContainsTreeNut                         bool    `json:"containsTreeNut"`
		ContainsFish                            bool    `json:"containsFish"`
		ContainsGluten                          bool    `json:"containsGluten"`
		AnimalFlesh                             bool    `json:"animalFlesh"`
		IsMeasuredVolumetrically                bool    `json:"isMeasuredVolumetrically"`
		IsLiquid                                bool    `json:"isLiquid"`
		ContainsSoy                             bool    `json:"containsSoy"`
		ContainsPeanut                          bool    `json:"containsPeanut"`
		AnimalDerived                           bool    `json:"animalDerived"`
		RestrictToPreparations                  bool    `json:"restrictToPreparations"`
		ContainsDairy                           bool    `json:"containsDairy"`
		ContainsSesame                          bool    `json:"containsSesame"`
	}

	// ValidIngredientDatabaseCreationInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientDatabaseCreationInput struct {
		_                                       struct{}
		ID                                      string  `json:"id"`
		Name                                    string  `json:"name"`
		Description                             string  `json:"description"`
		Warning                                 string  `json:"warning"`
		IconPath                                string  `json:"iconPath"`
		PluralName                              string  `json:"pluralName"`
		MinimumIdealStorageTemperatureInCelsius float32 `json:"minimumIdealStorageTemperatureInCelsius"`
		MaximumIdealStorageTemperatureInCelsius float32 `json:"maximumIdealStorageTemperatureInCelsius"`
		ContainsEgg                             bool    `json:"containsEgg"`
		ContainsWheat                           bool    `json:"containsWheat"`
		ContainsShellfish                       bool    `json:"containsShellfish"`
		ContainsTreeNut                         bool    `json:"containsTreeNut"`
		ContainsFish                            bool    `json:"containsFish"`
		ContainsGluten                          bool    `json:"containsGluten"`
		AnimalFlesh                             bool    `json:"animalFlesh"`
		IsMeasuredVolumetrically                bool    `json:"isMeasuredVolumetrically"`
		IsLiquid                                bool    `json:"isLiquid"`
		ContainsSoy                             bool    `json:"containsSoy"`
		ContainsPeanut                          bool    `json:"containsPeanut"`
		AnimalDerived                           bool    `json:"animalDerived"`
		RestrictToPreparations                  bool    `json:"restrictToPreparations"`
		ContainsDairy                           bool    `json:"containsDairy"`
		ContainsSesame                          bool    `json:"containsSesame"`
	}

	// ValidIngredientUpdateRequestInput represents what a user could set as input for updating valid ingredients.
	ValidIngredientUpdateRequestInput struct {
		_                                       struct{}
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
	}

	// ValidIngredientDataManager describes a structure capable of storing valid ingredients permanently.
	ValidIngredientDataManager interface {
		ValidIngredientExists(ctx context.Context, validIngredientID string) (bool, error)
		GetValidIngredient(ctx context.Context, validIngredientID string) (*ValidIngredient, error)
		GetRandomValidIngredient(ctx context.Context) (*ValidIngredient, error)
		GetValidIngredients(ctx context.Context, filter *QueryFilter) (*ValidIngredientList, error)
		SearchForValidIngredients(ctx context.Context, query string) ([]*ValidIngredient, error)
		SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string) ([]*ValidIngredient, error)
		CreateValidIngredient(ctx context.Context, input *ValidIngredientDatabaseCreationInput) (*ValidIngredient, error)
		UpdateValidIngredient(ctx context.Context, updated *ValidIngredient) error
		ArchiveValidIngredient(ctx context.Context, validIngredientID string) error
	}

	// ValidIngredientDataService describes a structure capable of serving traffic related to valid ingredients.
	ValidIngredientDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
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

	if input.MinimumIdealStorageTemperatureInCelsius != nil && *input.MinimumIdealStorageTemperatureInCelsius != x.MinimumIdealStorageTemperatureInCelsius {
		x.MinimumIdealStorageTemperatureInCelsius = *input.MinimumIdealStorageTemperatureInCelsius
	}

	if input.MaximumIdealStorageTemperatureInCelsius != nil && *input.MaximumIdealStorageTemperatureInCelsius != x.MaximumIdealStorageTemperatureInCelsius {
		x.MaximumIdealStorageTemperatureInCelsius = *input.MaximumIdealStorageTemperatureInCelsius
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

// ValidIngredientUpdateRequestInputFromValidIngredient creates a DatabaseCreationInput from a CreationInput.
func ValidIngredientUpdateRequestInputFromValidIngredient(input *ValidIngredient) *ValidIngredientUpdateRequestInput {
	x := &ValidIngredientUpdateRequestInput{
		Name:                                    &input.Name,
		Description:                             &input.Description,
		Warning:                                 &input.Warning,
		IconPath:                                &input.IconPath,
		ContainsDairy:                           &input.ContainsDairy,
		ContainsPeanut:                          &input.ContainsPeanut,
		ContainsTreeNut:                         &input.ContainsTreeNut,
		ContainsEgg:                             &input.ContainsEgg,
		ContainsWheat:                           &input.ContainsWheat,
		ContainsShellfish:                       &input.ContainsShellfish,
		ContainsSesame:                          &input.ContainsSesame,
		ContainsFish:                            &input.ContainsFish,
		ContainsGluten:                          &input.ContainsGluten,
		AnimalFlesh:                             &input.AnimalFlesh,
		IsMeasuredVolumetrically:                &input.IsMeasuredVolumetrically,
		IsLiquid:                                &input.IsLiquid,
		ContainsSoy:                             &input.ContainsSoy,
		PluralName:                              &input.PluralName,
		AnimalDerived:                           &input.AnimalDerived,
		RestrictToPreparations:                  &input.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: &input.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: &input.MaximumIdealStorageTemperatureInCelsius,
	}

	return x
}

// ValidIngredientDatabaseCreationInputFromValidIngredientCreationInput creates a DatabaseCreationInput from a CreationInput.
func ValidIngredientDatabaseCreationInputFromValidIngredientCreationInput(input *ValidIngredientCreationRequestInput) *ValidIngredientDatabaseCreationInput {
	x := &ValidIngredientDatabaseCreationInput{
		Name:                                    input.Name,
		Description:                             input.Description,
		Warning:                                 input.Warning,
		ContainsEgg:                             input.ContainsEgg,
		ContainsDairy:                           input.ContainsDairy,
		ContainsPeanut:                          input.ContainsPeanut,
		ContainsTreeNut:                         input.ContainsTreeNut,
		ContainsSoy:                             input.ContainsSoy,
		ContainsWheat:                           input.ContainsWheat,
		ContainsShellfish:                       input.ContainsShellfish,
		ContainsSesame:                          input.ContainsSesame,
		ContainsFish:                            input.ContainsFish,
		ContainsGluten:                          input.ContainsGluten,
		AnimalFlesh:                             input.AnimalFlesh,
		IsMeasuredVolumetrically:                input.IsMeasuredVolumetrically,
		IsLiquid:                                input.IsLiquid,
		IconPath:                                input.IconPath,
		PluralName:                              input.PluralName,
		AnimalDerived:                           input.AnimalDerived,
		RestrictToPreparations:                  input.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: input.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: input.MaximumIdealStorageTemperatureInCelsius,
	}

	return x
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
