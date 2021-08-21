package types

import (
	"context"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/search"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidIngredientsSearchIndexName is the name of the index used to search through valid ingredients.
	ValidIngredientsSearchIndexName search.IndexName = "valid_ingredients"
)

type (
	// ValidIngredient represents a valid ingredient.
	ValidIngredient struct {
		ArchivedOn        *uint64 `json:"archivedOn"`
		LastUpdatedOn     *uint64 `json:"lastUpdatedOn"`
		ExternalID        string  `json:"externalID"`
		Name              string  `json:"name"`
		Variant           string  `json:"variant"`
		Description       string  `json:"description"`
		Warning           string  `json:"warning"`
		IconPath          string  `json:"iconPath"`
		CreatedOn         uint64  `json:"createdOn"`
		ID                uint64  `json:"id"`
		ContainsSoy       bool    `json:"containsSoy"`
		ContainsTreeNut   bool    `json:"containsTreeNut"`
		ContainsShellfish bool    `json:"containsShellfish"`
		ContainsSesame    bool    `json:"containsSesame"`
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
		ValidIngredients []*ValidIngredient `json:"validIngredients"`
		Pagination
	}

	// ValidIngredientCreationInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientCreationInput struct {
		Name              string `json:"name"`
		Variant           string `json:"variant"`
		Description       string `json:"description"`
		Warning           string `json:"warning"`
		IconPath          string `json:"iconPath"`
		ContainsDairy     bool   `json:"containsDairy"`
		ContainsPeanut    bool   `json:"containsPeanut"`
		ContainsTreeNut   bool   `json:"containsTreeNut"`
		ContainsSoy       bool   `json:"containsSoy"`
		ContainsEgg       bool   `json:"containsEgg"`
		ContainsShellfish bool   `json:"containsShellfish"`
		ContainsSesame    bool   `json:"containsSesame"`
		ContainsFish      bool   `json:"containsFish"`
		ContainsGluten    bool   `json:"containsGluten"`
		AnimalFlesh       bool   `json:"animalFlesh"`
		AnimalDerived     bool   `json:"animalDerived"`
		Volumetric        bool   `json:"volumetric"`
		ContainsWheat     bool   `json:"containsWheat"`
	}

	// ValidIngredientUpdateInput represents what a user could set as input for updating valid ingredients.
	ValidIngredientUpdateInput struct {
		Name              string `json:"name"`
		Variant           string `json:"variant"`
		Description       string `json:"description"`
		Warning           string `json:"warning"`
		IconPath          string `json:"iconPath"`
		ContainsDairy     bool   `json:"containsDairy"`
		ContainsPeanut    bool   `json:"containsPeanut"`
		ContainsTreeNut   bool   `json:"containsTreeNut"`
		ContainsSoy       bool   `json:"containsSoy"`
		ContainsEgg       bool   `json:"containsEgg"`
		ContainsShellfish bool   `json:"containsShellfish"`
		ContainsSesame    bool   `json:"containsSesame"`
		ContainsFish      bool   `json:"containsFish"`
		ContainsGluten    bool   `json:"containsGluten"`
		AnimalFlesh       bool   `json:"animalFlesh"`
		AnimalDerived     bool   `json:"animalDerived"`
		Volumetric        bool   `json:"volumetric"`
		ContainsWheat     bool   `json:"containsWheat"`
	}

	// ValidIngredientDataManager describes a structure capable of storing valid ingredients permanently.
	ValidIngredientDataManager interface {
		ValidIngredientExists(ctx context.Context, validIngredientID uint64) (bool, error)
		GetValidIngredient(ctx context.Context, validIngredientID uint64) (*ValidIngredient, error)
		GetAllValidIngredientsCount(ctx context.Context) (uint64, error)
		GetAllValidIngredients(ctx context.Context, resultChannel chan []*ValidIngredient, bucketSize uint16) error
		GetValidIngredients(ctx context.Context, filter *QueryFilter) (*ValidIngredientList, error)
		GetValidIngredientsWithIDs(ctx context.Context, limit uint8, validPreparationID uint64, ids []uint64) ([]*ValidIngredient, error)
		CreateValidIngredient(ctx context.Context, input *ValidIngredientCreationInput, createdByUser uint64) (*ValidIngredient, error)
		UpdateValidIngredient(ctx context.Context, updated *ValidIngredient, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveValidIngredient(ctx context.Context, validIngredientID, archivedBy uint64) error
		GetAuditLogEntriesForValidIngredient(ctx context.Context, validIngredientID uint64) ([]*AuditLogEntry, error)
	}

	// ValidIngredientDataService describes a structure capable of serving traffic related to valid ingredients.
	ValidIngredientDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		SearchForValidIngredients(ctx context.Context, sessionCtxData *SessionContextData, validPreparation uint64, query string, filter *QueryFilter) ([]*ValidIngredient, error)
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidIngredientUpdateInput with a valid ingredient.
func (x *ValidIngredient) Update(input *ValidIngredientUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.Name != x.Name {
		out = append(out, &FieldChangeSummary{
			FieldName: "Name",
			OldValue:  x.Name,
			NewValue:  input.Name,
		})

		x.Name = input.Name
	}

	if input.Variant != x.Variant {
		out = append(out, &FieldChangeSummary{
			FieldName: "Variant",
			OldValue:  x.Variant,
			NewValue:  input.Variant,
		})

		x.Variant = input.Variant
	}

	if input.Description != x.Description {
		out = append(out, &FieldChangeSummary{
			FieldName: "Description",
			OldValue:  x.Description,
			NewValue:  input.Description,
		})

		x.Description = input.Description
	}

	if input.Warning != x.Warning {
		out = append(out, &FieldChangeSummary{
			FieldName: "Warning",
			OldValue:  x.Warning,
			NewValue:  input.Warning,
		})

		x.Warning = input.Warning
	}

	if input.ContainsEgg != x.ContainsEgg {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsEgg",
			OldValue:  x.ContainsEgg,
			NewValue:  input.ContainsEgg,
		})

		x.ContainsEgg = input.ContainsEgg
	}

	if input.ContainsDairy != x.ContainsDairy {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsDairy",
			OldValue:  x.ContainsDairy,
			NewValue:  input.ContainsDairy,
		})

		x.ContainsDairy = input.ContainsDairy
	}

	if input.ContainsPeanut != x.ContainsPeanut {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsPeanut",
			OldValue:  x.ContainsPeanut,
			NewValue:  input.ContainsPeanut,
		})

		x.ContainsPeanut = input.ContainsPeanut
	}

	if input.ContainsTreeNut != x.ContainsTreeNut {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsTreeNut",
			OldValue:  x.ContainsTreeNut,
			NewValue:  input.ContainsTreeNut,
		})

		x.ContainsTreeNut = input.ContainsTreeNut
	}

	if input.ContainsSoy != x.ContainsSoy {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsSoy",
			OldValue:  x.ContainsSoy,
			NewValue:  input.ContainsSoy,
		})

		x.ContainsSoy = input.ContainsSoy
	}

	if input.ContainsWheat != x.ContainsWheat {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsWheat",
			OldValue:  x.ContainsWheat,
			NewValue:  input.ContainsWheat,
		})

		x.ContainsWheat = input.ContainsWheat
	}

	if input.ContainsShellfish != x.ContainsShellfish {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsShellfish",
			OldValue:  x.ContainsShellfish,
			NewValue:  input.ContainsShellfish,
		})

		x.ContainsShellfish = input.ContainsShellfish
	}

	if input.ContainsSesame != x.ContainsSesame {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsSesame",
			OldValue:  x.ContainsSesame,
			NewValue:  input.ContainsSesame,
		})

		x.ContainsSesame = input.ContainsSesame
	}

	if input.ContainsFish != x.ContainsFish {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsFish",
			OldValue:  x.ContainsFish,
			NewValue:  input.ContainsFish,
		})

		x.ContainsFish = input.ContainsFish
	}

	if input.ContainsGluten != x.ContainsGluten {
		out = append(out, &FieldChangeSummary{
			FieldName: "ContainsGluten",
			OldValue:  x.ContainsGluten,
			NewValue:  input.ContainsGluten,
		})

		x.ContainsGluten = input.ContainsGluten
	}

	if input.AnimalFlesh != x.AnimalFlesh {
		out = append(out, &FieldChangeSummary{
			FieldName: "AnimalFlesh",
			OldValue:  x.AnimalFlesh,
			NewValue:  input.AnimalFlesh,
		})

		x.AnimalFlesh = input.AnimalFlesh
	}

	if input.AnimalDerived != x.AnimalDerived {
		out = append(out, &FieldChangeSummary{
			FieldName: "AnimalDerived",
			OldValue:  x.AnimalDerived,
			NewValue:  input.AnimalDerived,
		})

		x.AnimalDerived = input.AnimalDerived
	}

	if input.Volumetric != x.Volumetric {
		out = append(out, &FieldChangeSummary{
			FieldName: "Volumetric",
			OldValue:  x.Volumetric,
			NewValue:  input.Volumetric,
		})

		x.Volumetric = input.Volumetric
	}

	if input.IconPath != x.IconPath {
		out = append(out, &FieldChangeSummary{
			FieldName: "IconPath",
			OldValue:  x.IconPath,
			NewValue:  input.IconPath,
		})

		x.IconPath = input.IconPath
	}

	return out
}

var _ validation.ValidatableWithContext = (*ValidIngredientCreationInput)(nil)

// ValidateWithContext validates a ValidIngredientCreationInput.
func (x *ValidIngredientCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientUpdateInput)(nil)

// ValidateWithContext validates a ValidIngredientUpdateInput.
func (x *ValidIngredientUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}
