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
		_                 struct{}
		LastUpdatedOn     *uint64 `json:"lastUpdatedOn"`
		ArchivedOn        *uint64 `json:"archivedOn"`
		Name              string  `json:"name"`
		Description       string  `json:"description"`
		Warning           string  `json:"warning"`
		ID                string  `json:"id"`
		IconPath          string  `json:"iconPath"`
		Variant           string  `json:"variant"`
		CreatedOn         uint64  `json:"createdOn"`
		ContainsSesame    bool    `json:"containsSesame"`
		ContainsSoy       bool    `json:"containsSoy"`
		ContainsShellfish bool    `json:"containsShellfish"`
		ContainsTreeNut   bool    `json:"containsTreeNut"`
		ContainsFish      bool    `json:"containsFish"`
		ContainsGluten    bool    `json:"containsGluten"`
		AnimalFlesh       bool    `json:"animalFlesh"`
		AnimalDerived     bool    `json:"animalDerived"`
		Volumetric        bool    `json:"volumetric"`
		ContainsPeanut    bool    `json:"containsPeanut"`
		ContainsDairy     bool    `json:"containsDairy"`
		ContainsEgg       bool    `json:"containsEgg"`
		ContainsWheat     bool    `json:"containsWheat"`
	}

	// ValidIngredientList represents a list of valid ingredients.
	ValidIngredientList struct {
		_                struct{}
		ValidIngredients []*ValidIngredient `json:"validIngredients"`
		Pagination
	}

	// ValidIngredientCreationRequestInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientCreationRequestInput struct {
		_                 struct{}
		ID                string `json:"-"`
		Name              string `json:"name"`
		Variant           string `json:"variant"`
		Description       string `json:"description"`
		Warning           string `json:"warning"`
		IconPath          string `json:"iconPath"`
		ContainsDairy     bool   `json:"containsDairy"`
		ContainsPeanut    bool   `json:"containsPeanut"`
		ContainsTreeNut   bool   `json:"containsTreeNut"`
		ContainsEgg       bool   `json:"containsEgg"`
		ContainsWheat     bool   `json:"containsWheat"`
		ContainsShellfish bool   `json:"containsShellfish"`
		ContainsSesame    bool   `json:"containsSesame"`
		ContainsFish      bool   `json:"containsFish"`
		ContainsGluten    bool   `json:"containsGluten"`
		AnimalFlesh       bool   `json:"animalFlesh"`
		AnimalDerived     bool   `json:"animalDerived"`
		Volumetric        bool   `json:"volumetric"`
		ContainsSoy       bool   `json:"containsSoy"`
	}

	// ValidIngredientDatabaseCreationInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientDatabaseCreationInput struct {
		_                 struct{}
		ID                string `json:"id"`
		Name              string `json:"name"`
		Variant           string `json:"variant"`
		Description       string `json:"description"`
		Warning           string `json:"warning"`
		IconPath          string `json:"iconPath"`
		ContainsDairy     bool   `json:"containsDairy"`
		ContainsPeanut    bool   `json:"containsPeanut"`
		ContainsTreeNut   bool   `json:"containsTreeNut"`
		ContainsEgg       bool   `json:"containsEgg"`
		ContainsWheat     bool   `json:"containsWheat"`
		ContainsShellfish bool   `json:"containsShellfish"`
		ContainsSesame    bool   `json:"containsSesame"`
		ContainsFish      bool   `json:"containsFish"`
		ContainsGluten    bool   `json:"containsGluten"`
		AnimalFlesh       bool   `json:"animalFlesh"`
		AnimalDerived     bool   `json:"animalDerived"`
		Volumetric        bool   `json:"volumetric"`
		ContainsSoy       bool   `json:"containsSoy"`
	}

	// ValidIngredientUpdateRequestInput represents what a user could set as input for updating valid ingredients.
	ValidIngredientUpdateRequestInput struct {
		_                 struct{}
		Name              string `json:"name"`
		Variant           string `json:"variant"`
		Description       string `json:"description"`
		Warning           string `json:"warning"`
		IconPath          string `json:"iconPath"`
		ContainsDairy     bool   `json:"containsDairy"`
		ContainsPeanut    bool   `json:"containsPeanut"`
		ContainsTreeNut   bool   `json:"containsTreeNut"`
		ContainsEgg       bool   `json:"containsEgg"`
		ContainsWheat     bool   `json:"containsWheat"`
		ContainsShellfish bool   `json:"containsShellfish"`
		ContainsSesame    bool   `json:"containsSesame"`
		ContainsFish      bool   `json:"containsFish"`
		ContainsGluten    bool   `json:"containsGluten"`
		AnimalFlesh       bool   `json:"animalFlesh"`
		AnimalDerived     bool   `json:"animalDerived"`
		Volumetric        bool   `json:"volumetric"`
		ContainsSoy       bool   `json:"containsSoy"`
	}

	// ValidIngredientDataManager describes a structure capable of storing valid ingredients permanently.
	ValidIngredientDataManager interface {
		ValidIngredientExists(ctx context.Context, validIngredientID string) (bool, error)
		GetValidIngredient(ctx context.Context, validIngredientID string) (*ValidIngredient, error)
		GetTotalValidIngredientCount(ctx context.Context) (uint64, error)
		GetValidIngredients(ctx context.Context, filter *QueryFilter) (*ValidIngredientList, error)
		GetValidIngredientsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*ValidIngredient, error)
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
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidIngredientUpdateRequestInput with a valid ingredient.
func (x *ValidIngredient) Update(input *ValidIngredientUpdateRequestInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Variant != "" && input.Variant != x.Variant {
		x.Variant = input.Variant
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.Warning != "" && input.Warning != x.Warning {
		x.Warning = input.Warning
	}

	if input.ContainsEgg != x.ContainsEgg {
		x.ContainsEgg = input.ContainsEgg
	}

	if input.ContainsDairy != x.ContainsDairy {
		x.ContainsDairy = input.ContainsDairy
	}

	if input.ContainsPeanut != x.ContainsPeanut {
		x.ContainsPeanut = input.ContainsPeanut
	}

	if input.ContainsTreeNut != x.ContainsTreeNut {
		x.ContainsTreeNut = input.ContainsTreeNut
	}

	if input.ContainsSoy != x.ContainsSoy {
		x.ContainsSoy = input.ContainsSoy
	}

	if input.ContainsWheat != x.ContainsWheat {
		x.ContainsWheat = input.ContainsWheat
	}

	if input.ContainsShellfish != x.ContainsShellfish {
		x.ContainsShellfish = input.ContainsShellfish
	}

	if input.ContainsSesame != x.ContainsSesame {
		x.ContainsSesame = input.ContainsSesame
	}

	if input.ContainsFish != x.ContainsFish {
		x.ContainsFish = input.ContainsFish
	}

	if input.ContainsGluten != x.ContainsGluten {
		x.ContainsGluten = input.ContainsGluten
	}

	if input.AnimalFlesh != x.AnimalFlesh {
		x.AnimalFlesh = input.AnimalFlesh
	}

	if input.AnimalDerived != x.AnimalDerived {
		x.AnimalDerived = input.AnimalDerived
	}

	if input.Volumetric != x.Volumetric {
		x.Volumetric = input.Volumetric
	}

	if input.IconPath != "" && input.IconPath != x.IconPath {
		x.IconPath = input.IconPath
	}
}

var _ validation.ValidatableWithContext = (*ValidIngredientCreationRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientCreationRequestInput.
func (x *ValidIngredientCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Variant, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.Warning, validation.Required),
		validation.Field(&x.IconPath, validation.Required),
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
		validation.Field(&x.Variant, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.Warning, validation.Required),
		validation.Field(&x.IconPath, validation.Required),
	)
}

// ValidIngredientDatabaseCreationInputFromValidIngredientCreationInput creates a DatabaseCreationInput from a CreationInput.
func ValidIngredientDatabaseCreationInputFromValidIngredientCreationInput(input *ValidIngredientCreationRequestInput) *ValidIngredientDatabaseCreationInput {
	x := &ValidIngredientDatabaseCreationInput{
		Name:              input.Name,
		Variant:           input.Variant,
		Description:       input.Description,
		Warning:           input.Warning,
		ContainsEgg:       input.ContainsEgg,
		ContainsDairy:     input.ContainsDairy,
		ContainsPeanut:    input.ContainsPeanut,
		ContainsTreeNut:   input.ContainsTreeNut,
		ContainsSoy:       input.ContainsSoy,
		ContainsWheat:     input.ContainsWheat,
		ContainsShellfish: input.ContainsShellfish,
		ContainsSesame:    input.ContainsSesame,
		ContainsFish:      input.ContainsFish,
		ContainsGluten:    input.ContainsGluten,
		AnimalFlesh:       input.AnimalFlesh,
		AnimalDerived:     input.AnimalDerived,
		Volumetric:        input.Volumetric,
		IconPath:          input.IconPath,
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
		validation.Field(&x.Variant, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.Warning, validation.Required),
		validation.Field(&x.IconPath, validation.Required),
	)
}
