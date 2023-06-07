package indexing

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
)

type NamedID struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// RecipeSearchSubset is a subset of Recipe fields for search indexing.
type RecipeSearchSubset struct {
	ID          string                    `json:"id,omitempty"`
	Name        string                    `json:"name,omitempty"`
	Description string                    `json:"description,omitempty"`
	Steps       []*RecipeStepSearchSubset `json:"steps,omitempty"`
}

func RecipeSearchSubsetFromRecipe(r *types.Recipe) *RecipeSearchSubset {
	x := &RecipeSearchSubset{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, step := range r.Steps {
		x.Steps = append(x.Steps, RecipeStepSearchSubsetFromRecipeStep(step))
	}

	return x
}

// RecipeStepSearchSubset is a subset of RecipeStep fields for search indexing.
type RecipeStepSearchSubset struct {
	Preparation string    `json:"preparation,omitempty"`
	Ingredients []NamedID `json:"ingredients,omitempty"`
	Instruments []NamedID `json:"instruments,omitempty"`
	Vessels     []NamedID `json:"vessels,omitempty"`
}

func RecipeStepSearchSubsetFromRecipeStep(x *types.RecipeStep) *RecipeStepSearchSubset {
	stepSubset := &RecipeStepSearchSubset{
		Preparation: x.Preparation.Name,
	}

	for _, ingredient := range x.Ingredients {
		stepSubset.Ingredients = append(stepSubset.Ingredients, NamedID{ID: ingredient.ID, Name: ingredient.Name})
	}

	for _, instrument := range x.Instruments {
		stepSubset.Instruments = append(stepSubset.Instruments, NamedID{ID: instrument.ID, Name: instrument.Name})
	}

	for _, vessel := range x.Vessels {
		stepSubset.Vessels = append(stepSubset.Vessels, NamedID{ID: vessel.ID, Name: vessel.Name})
	}

	return stepSubset
}

// MealSearchSubset is a subset of Meal fields for search indexing.
type MealSearchSubset struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Recipes     []NamedID `json:"recipes,omitempty"`
}

func MealSearchSubsetFromMeal(r *types.Meal) *MealSearchSubset {
	x := &MealSearchSubset{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, component := range r.Components {
		x.Recipes = append(x.Recipes, NamedID{ID: component.Recipe.ID, Name: component.Recipe.Name})
	}

	return x
}

// ValidIngredientMeasurementUnitSearchSubset is a subset of ValidIngredientMeasurementUnit fields for search indexing.
type ValidIngredientMeasurementUnitSearchSubset struct {
	ID              string  `json:"id,omitempty"`
	Notes           string  `json:"notes,omitempty"`
	MeasurementUnit NamedID `json:"measurementUnit,omitempty"`
	Ingredient      NamedID `json:"ingredient,omitempty"`
}

func ValidIngredientMeasurementUnitSearchSubsetFromValidIngredientMeasurementUnit(x *types.ValidIngredientMeasurementUnit) *ValidIngredientMeasurementUnitSearchSubset {
	y := &ValidIngredientMeasurementUnitSearchSubset{
		ID:              x.ID,
		Notes:           x.Notes,
		MeasurementUnit: NamedID{ID: x.MeasurementUnit.ID, Name: x.MeasurementUnit.Name},
		Ingredient:      NamedID{ID: x.Ingredient.ID, Name: x.Ingredient.Name},
	}

	return y
}

// ValidMeasurementUnitConversionSearchSubset is a subset of ValidMeasurementUnitConversion fields for search indexing.
type ValidMeasurementUnitConversionSearchSubset struct {
	ID                  string  `json:"id,omitempty"`
	Notes               string  `json:"notes,omitempty"`
	FromMeasurementUnit NamedID `json:"fromMeasurementUnit,omitempty"`
	ToMeasurementUnit   NamedID `json:"toMeasurementUnit,omitempty"`
}

func ValidMeasurementUnitConversionSearchSubsetFromValidMeasurementUnitConversion(x *types.ValidMeasurementUnitConversion) *ValidMeasurementUnitConversionSearchSubset {
	y := &ValidMeasurementUnitConversionSearchSubset{
		ID:                  x.ID,
		Notes:               x.Notes,
		FromMeasurementUnit: NamedID{ID: x.From.ID, Name: x.From.Name},
		ToMeasurementUnit:   NamedID{ID: x.To.ID, Name: x.To.Name},
	}

	return y
}

// ValidPreparationInstrumentSearchSubset is a subset of ValidPreparationInstrument fields for search indexing.
type ValidPreparationInstrumentSearchSubset struct {
	ID          string  `json:"id,omitempty"`
	Notes       string  `json:"notes,omitempty"`
	Instrument  NamedID `json:"instrument,omitempty"`
	Preparation NamedID `json:"preparation,omitempty"`
}

func ValidPreparationInstrumentSearchSubsetFromValidPreparationInstrument(x *types.ValidPreparationInstrument) *ValidPreparationInstrumentSearchSubset {
	y := &ValidPreparationInstrumentSearchSubset{
		ID:          x.ID,
		Notes:       x.Notes,
		Instrument:  NamedID{ID: x.Instrument.ID, Name: x.Instrument.Name},
		Preparation: NamedID{ID: x.Preparation.ID, Name: x.Preparation.Name},
	}

	return y
}
