package indexing

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ValidMeasurementUnitSearchSubset represents the subset of values suitable to index for search.
type ValidMeasurementUnitSearchSubset struct {
	_ struct{} `json:"-"`

	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	PluralName  string `json:"pluralName,omitempty"`
}

// ConvertValidMeasurementUnitToValidMeasurementUnitSearchSubset converts a ValidMeasurementUnit to a ValidMeasurementUnitSearchSubset.
func ConvertValidMeasurementUnitToValidMeasurementUnitSearchSubset(x *types.ValidMeasurementUnit) *ValidMeasurementUnitSearchSubset {
	return &ValidMeasurementUnitSearchSubset{
		ID:          x.ID,
		Name:        x.Name,
		PluralName:  x.PluralName,
		Description: x.Description,
	}
}

// MealSearchSubset represents the subset of values suitable to index for search.
type MealSearchSubset struct {
	_ struct{} `json:"-"`

	ID          string          `json:"id,omitempty"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Recipes     []types.NamedID `json:"recipes,omitempty"`
}

func ConvertMealToMealSearchSubset(r *types.Meal) *MealSearchSubset {
	x := &MealSearchSubset{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, component := range r.Components {
		x.Recipes = append(x.Recipes, types.NamedID{ID: component.Recipe.ID, Name: component.Recipe.Name})
	}

	return x
}

// RecipeSearchSubset represents the subset of values suitable to index for search.
type RecipeSearchSubset struct {
	_ struct{} `json:"-"`

	ID          string                    `json:"id,omitempty"`
	Name        string                    `json:"name,omitempty"`
	Description string                    `json:"description,omitempty"`
	Steps       []*RecipeStepSearchSubset `json:"steps,omitempty"`
}

// ConvertRecipeToRecipeSearchSubset converts a Recipe to a RecipeSearchSubset.
func ConvertRecipeToRecipeSearchSubset(r *types.Recipe) *RecipeSearchSubset {
	x := &RecipeSearchSubset{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, step := range r.Steps {
		x.Steps = append(x.Steps, ConvertRecipeStepToRecipeStepSearchSubset(step))
	}

	return x
}

// RecipeStepSearchSubset represents the subset of values suitable to index for search.
type RecipeStepSearchSubset struct {
	_ struct{} `json:"-"`

	Preparation string          `json:"preparation,omitempty"`
	Ingredients []types.NamedID `json:"ingredients,omitempty"`
	Instruments []types.NamedID `json:"instruments,omitempty"`
	Vessels     []types.NamedID `json:"vessels,omitempty"`
}

func ConvertRecipeStepToRecipeStepSearchSubset(x *types.RecipeStep) *RecipeStepSearchSubset {
	stepSubset := &RecipeStepSearchSubset{
		Preparation: x.Preparation.Name,
	}

	for _, ingredient := range x.Ingredients {
		stepSubset.Ingredients = append(stepSubset.Ingredients, types.NamedID{ID: ingredient.ID, Name: ingredient.Name})
	}

	for _, instrument := range x.Instruments {
		stepSubset.Instruments = append(stepSubset.Instruments, types.NamedID{ID: instrument.ID, Name: instrument.Name})
	}

	for _, vessel := range x.Vessels {
		stepSubset.Vessels = append(stepSubset.Vessels, types.NamedID{ID: vessel.ID, Name: vessel.Name})
	}

	return stepSubset
}

// ValidIngredientSearchSubset represents the subset of values suitable to index for search.
type ValidIngredientSearchSubset struct {
	_ struct{} `json:"-"`

	PluralName          string `json:"pluralName,omitempty"`
	Name                string `json:"name,omitempty"`
	ID                  string `json:"id,omitempty"`
	Description         string `json:"description,omitempty"`
	ShoppingSuggestions string `json:"shoppingSuggestions,omitempty"`
}

// ConvertValidIngredientToValidIngredientSearchSubset converts a ValidIngredient to a ValidIngredientSearchSubset.
func ConvertValidIngredientToValidIngredientSearchSubset(x *types.ValidIngredient) *ValidIngredientSearchSubset {
	return &ValidIngredientSearchSubset{
		ID:                  x.ID,
		Name:                x.Name,
		PluralName:          x.PluralName,
		Description:         x.Description,
		ShoppingSuggestions: x.ShoppingSuggestions,
	}
}

// ValidIngredientStateSearchSubset represents the subset of values suitable to index for search.
type ValidIngredientStateSearchSubset struct {
	_ struct{} `json:"-"`

	ID            string `json:"id,omitempty"`
	PastTense     string `json:"pastTense,omitempty"`
	Description   string `json:"description,omitempty"`
	Name          string `json:"name,omitempty"`
	AttributeType string `json:"attributeType,omitempty"`
}

// ConvertValidIngredientStateToValidIngredientStateSearchSubset converts a ValidIngredientState to a ValidIngredientStateSearchSubset.
func ConvertValidIngredientStateToValidIngredientStateSearchSubset(x *types.ValidIngredientState) *ValidIngredientStateSearchSubset {
	return &ValidIngredientStateSearchSubset{
		ID:            x.ID,
		Name:          x.Name,
		PastTense:     x.PastTense,
		Description:   x.Description,
		AttributeType: x.AttributeType,
	}
}

// ValidInstrumentSearchSubset represents the subset of values suitable to index for search.
type ValidInstrumentSearchSubset struct {
	_ struct{} `json:"-"`

	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	PluralName  string `json:"pluralName,omitempty"`
	Description string `json:"description,omitempty"`
}

// ConvertValidInstrumentToValidInstrumentSearchSubset converts a ValidInstrument to a ValidInstrumentSearchSubset.
func ConvertValidInstrumentToValidInstrumentSearchSubset(x *types.ValidInstrument) *ValidInstrumentSearchSubset {
	return &ValidInstrumentSearchSubset{
		ID:          x.ID,
		Name:        x.Name,
		PluralName:  x.PluralName,
		Description: x.Description,
	}
}

// ValidPreparationSearchSubset represents the subset of values suitable to index for search.
type ValidPreparationSearchSubset struct {
	_ struct{} `json:"-"`

	PastTense   string `json:"pastTense,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ConvertValidPreparationToValidPreparationSearchSubset converts a ValidPreparation to a ValidPreparationSearchSubset.
func ConvertValidPreparationToValidPreparationSearchSubset(x *types.ValidPreparation) *ValidPreparationSearchSubset {
	return &ValidPreparationSearchSubset{
		ID:          x.ID,
		Name:        x.Name,
		PastTense:   x.PastTense,
		Description: x.Description,
	}
}

// ValidVesselSearchSubset represents the subset of values suitable to index for search.
type ValidVesselSearchSubset struct {
	_ struct{} `json:"-"`

	ID               string  `json:"id,omitempty"`
	Name             string  `json:"name,omitempty"`
	PluralName       string  `json:"pluralName,omitempty"`
	Description      string  `json:"description,omitempty"`
	CapacityUnitName string  `json:"capacityUnitName"`
	Capacity         float32 `json:"capacity,omitempty"`
}

// ConvertValidVesselToValidVesselSearchSubset converts a ValidVessel to a ValidVesselSearchSubset.
func ConvertValidVesselToValidVesselSearchSubset(x *types.ValidVessel) *ValidVesselSearchSubset {
	return &ValidVesselSearchSubset{
		ID:          x.ID,
		Name:        x.Name,
		PluralName:  x.PluralName,
		Description: x.Description,
	}
}
