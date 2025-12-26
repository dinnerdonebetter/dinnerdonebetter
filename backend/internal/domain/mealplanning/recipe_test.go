package mealplanning

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipe_FindStepForIndex(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: t.Name(),
				},
			},
		}

		assert.NotNil(t, x.FindStepForIndex(0))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		assert.Nil(t, x.FindStepForIndex(0))
	})
}

func TestRecipe_FindStepByID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: t.Name(),
				},
			},
		}

		assert.NotNil(t, x.FindStepByID(t.Name()))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		assert.Nil(t, x.FindStepByID("whatever"))
	})
}

func TestRecipe_FindStepForRecipeStepProductID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: "whatever",
					Products: []*RecipeStepProduct{
						{
							ID: t.Name(),
						},
					},
				},
			},
		}

		assert.NotNil(t, x.FindStepForRecipeStepProductID(t.Name()))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		assert.Nil(t, x.FindStepForRecipeStepProductID("whatever"))
	})
}

func TestRecipe_GetRelatedRecipeIDs(T *testing.T) {
	T.Parallel()

	T.Run("returns related recipe IDs from ingredients", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: "step-1",
					Ingredients: []*RecipeStepIngredient{
						{RecipeStepProductRecipeID: pointer.To("recipe-1")},
						{RecipeStepProductRecipeID: pointer.To("recipe-2")},
					},
				},
				{
					ID: "step-2",
					Ingredients: []*RecipeStepIngredient{
						{RecipeStepProductRecipeID: pointer.To("recipe-3")},
					},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Equal(t, []string{"recipe-1", "recipe-2", "recipe-3"}, ids)
	})

	T.Run("deduplicates recipe IDs", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: "step-1",
					Ingredients: []*RecipeStepIngredient{
						{RecipeStepProductRecipeID: pointer.To("recipe-1")},
						{RecipeStepProductRecipeID: pointer.To("recipe-2")},
					},
				},
				{
					ID: "step-2",
					Ingredients: []*RecipeStepIngredient{
						{RecipeStepProductRecipeID: pointer.To("recipe-1")}, // duplicate
						{RecipeStepProductRecipeID: pointer.To("recipe-3")},
					},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Equal(t, []string{"recipe-1", "recipe-2", "recipe-3"}, ids)
	})

	T.Run("skips nil and empty values", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: "step-1",
					Ingredients: []*RecipeStepIngredient{
						{RecipeStepProductRecipeID: pointer.To("recipe-1")},
						{RecipeStepProductRecipeID: nil},
						{RecipeStepProductRecipeID: pointer.To("")},
						{RecipeStepProductRecipeID: pointer.To("recipe-2")},
					},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Equal(t, []string{"recipe-1", "recipe-2"}, ids)
	})

	T.Run("returns empty slice when no related recipes", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: "step-1",
					Ingredients: []*RecipeStepIngredient{
						{RecipeStepProductRecipeID: nil},
						{Name: "regular ingredient"},
					},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Empty(t, ids)
	})

	T.Run("returns empty slice for recipe with no steps", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		ids := x.GetRelatedRecipeIDs()
		assert.Empty(t, ids)
	})

	T.Run("returns empty slice for steps with no ingredients", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID:          "step-1",
					Ingredients: nil,
				},
				{
					ID:          "step-2",
					Ingredients: []*RecipeStepIngredient{},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Empty(t, ids)
	})
}

func TestRecipe_FindStepIndexByID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: t.Name(),
				},
			},
		}

		x.FindStepIndexByID(t.Name())
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		assert.Equal(t, -1, x.FindStepIndexByID("whatever"))
	})
}

func TestRecipe_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}
		input := &RecipeUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.EligibleForMeals = pointer.To(true)

		x.Update(input)
	})
}

func TestRecipeCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationRequestInput{
			Name:                t.Name(),
			Source:              t.Name(),
			Slug:                t.Name(),
			PortionName:         t.Name(),
			PluralPortionName:   t.Name(),
			Description:         t.Name(),
			YieldsComponentType: MealComponentTypesMain,
			InspiredByRecipeID:  pointer.To(t.Name()),
			Steps: []*RecipeStepCreationRequestInput{
				buildValidRecipeStepCreationRequestInput(),
				buildValidRecipeStepCreationRequestInput(),
			},
			PrepTasks: []*RecipePrepTaskWithinRecipeCreationRequestInput{
				{
					RecipeSteps: []*RecipePrepTaskStepWithinRecipeCreationRequestInput{
						{
							BelongsToRecipeStepIndex: 0,
						},
					},
				},
			},
			EstimatedPortions: types.Float32RangeWithOptionalMax{Min: fake.Float32()},
		}

		assert.NoError(t, x.ValidateWithContext(t.Context()))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})

	T.Run("with overreferenced task steps", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationRequestInput{
			Name:                t.Name(),
			Source:              t.Name(),
			Slug:                t.Name(),
			PortionName:         t.Name(),
			PluralPortionName:   t.Name(),
			Description:         t.Name(),
			YieldsComponentType: MealComponentTypesMain,
			InspiredByRecipeID:  pointer.To(t.Name()),
			Steps: []*RecipeStepCreationRequestInput{
				buildValidRecipeStepCreationRequestInput(),
				buildValidRecipeStepCreationRequestInput(),
			},
			PrepTasks: []*RecipePrepTaskWithinRecipeCreationRequestInput{
				{
					RecipeSteps: []*RecipePrepTaskStepWithinRecipeCreationRequestInput{
						{
							BelongsToRecipeStepIndex: 0,
						},
						{
							BelongsToRecipeStepIndex: 0,
						},
					},
				},
			},
			EstimatedPortions: types.Float32RangeWithOptionalMax{Min: fake.Float32()},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipeDatabaseCreationInput_FindStepByIndex(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Index: 0,
				},
			},
		}

		assert.NotNil(t, x.FindStepByIndex(0))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Index: 0,
				},
			},
		}

		assert.Nil(t, x.FindStepByIndex(1))
	})
}

func TestRecipeDatabaseCreationInput_FindStepByID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID: t.Name(),
				},
			},
		}

		assert.NotNil(t, x.FindStepByID(t.Name()))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID: t.Name(),
				},
			},
		}

		assert.Nil(t, x.FindStepByID("whatever"))
	})
}

func TestRecipeDatabaseCreationInput_GetAllValidIngredientPreparationIDs(T *testing.T) {
	T.Parallel()

	T.Run("returns IDs from multiple steps", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{ValidIngredientPreparationID: pointer.To("vip-1")},
						{ValidIngredientPreparationID: pointer.To("vip-2")},
					},
				},
				{
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{ValidIngredientPreparationID: pointer.To("vip-3")},
					},
				},
			},
		}

		ids := x.GetAllValidIngredientPreparationIDs()
		assert.Equal(t, []string{"vip-1", "vip-2", "vip-3"}, ids)
	})

	T.Run("skips nil and empty values", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{ValidIngredientPreparationID: pointer.To("vip-1")},
						{ValidIngredientPreparationID: nil},
						{ValidIngredientPreparationID: pointer.To("")},
						{ValidIngredientPreparationID: pointer.To("vip-2")},
					},
				},
			},
		}

		ids := x.GetAllValidIngredientPreparationIDs()
		assert.Equal(t, []string{"vip-1", "vip-2"}, ids)
	})

	T.Run("returns empty slice when no IDs present", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{ValidIngredientPreparationID: nil},
					},
				},
			},
		}

		ids := x.GetAllValidIngredientPreparationIDs()
		assert.Empty(t, ids)
	})
}

func TestRecipeDatabaseCreationInput_GetAllValidIngredientMeasurementUnitIDs(T *testing.T) {
	T.Parallel()

	T.Run("returns IDs from multiple steps", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{ValidIngredientMeasurementUnitID: pointer.To("vimu-1")},
						{ValidIngredientMeasurementUnitID: pointer.To("vimu-2")},
					},
				},
				{
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{ValidIngredientMeasurementUnitID: pointer.To("vimu-3")},
					},
				},
			},
		}

		ids := x.GetAllValidIngredientMeasurementUnitIDs()
		assert.Equal(t, []string{"vimu-1", "vimu-2", "vimu-3"}, ids)
	})

	T.Run("skips nil and empty values", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{ValidIngredientMeasurementUnitID: pointer.To("vimu-1")},
						{ValidIngredientMeasurementUnitID: nil},
						{ValidIngredientMeasurementUnitID: pointer.To("")},
						{ValidIngredientMeasurementUnitID: pointer.To("vimu-2")},
					},
				},
			},
		}

		ids := x.GetAllValidIngredientMeasurementUnitIDs()
		assert.Equal(t, []string{"vimu-1", "vimu-2"}, ids)
	})

	T.Run("returns empty slice when no IDs present", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{ValidIngredientMeasurementUnitID: nil},
					},
				},
			},
		}

		ids := x.GetAllValidIngredientMeasurementUnitIDs()
		assert.Empty(t, ids)
	})
}

func TestRecipeDatabaseCreationInput_GetAllValidPreparationInstrumentIDs(T *testing.T) {
	T.Parallel()

	T.Run("returns IDs from multiple steps", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Instruments: []*RecipeStepInstrumentDatabaseCreationInput{
						{ValidPreparationInstrumentID: pointer.To("vpi-1")},
						{ValidPreparationInstrumentID: pointer.To("vpi-2")},
					},
				},
				{
					Instruments: []*RecipeStepInstrumentDatabaseCreationInput{
						{ValidPreparationInstrumentID: pointer.To("vpi-3")},
					},
				},
			},
		}

		ids := x.GetAllValidPreparationInstrumentIDs()
		assert.Equal(t, []string{"vpi-1", "vpi-2", "vpi-3"}, ids)
	})

	T.Run("skips nil and empty values", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Instruments: []*RecipeStepInstrumentDatabaseCreationInput{
						{ValidPreparationInstrumentID: pointer.To("vpi-1")},
						{ValidPreparationInstrumentID: nil},
						{ValidPreparationInstrumentID: pointer.To("")},
						{ValidPreparationInstrumentID: pointer.To("vpi-2")},
					},
				},
			},
		}

		ids := x.GetAllValidPreparationInstrumentIDs()
		assert.Equal(t, []string{"vpi-1", "vpi-2"}, ids)
	})

	T.Run("returns empty slice when no IDs present", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Instruments: []*RecipeStepInstrumentDatabaseCreationInput{
						{ValidPreparationInstrumentID: nil},
					},
				},
			},
		}

		ids := x.GetAllValidPreparationInstrumentIDs()
		assert.Empty(t, ids)
	})
}

func TestRecipeDatabaseCreationInput_GetAllValidPreparationVesselIDs(T *testing.T) {
	T.Parallel()

	T.Run("returns IDs from multiple steps", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Vessels: []*RecipeStepVesselDatabaseCreationInput{
						{ValidPreparationVesselID: pointer.To("vpv-1")},
						{ValidPreparationVesselID: pointer.To("vpv-2")},
					},
				},
				{
					Vessels: []*RecipeStepVesselDatabaseCreationInput{
						{ValidPreparationVesselID: pointer.To("vpv-3")},
					},
				},
			},
		}

		ids := x.GetAllValidPreparationVesselIDs()
		assert.Equal(t, []string{"vpv-1", "vpv-2", "vpv-3"}, ids)
	})

	T.Run("skips nil and empty values", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Vessels: []*RecipeStepVesselDatabaseCreationInput{
						{ValidPreparationVesselID: pointer.To("vpv-1")},
						{ValidPreparationVesselID: nil},
						{ValidPreparationVesselID: pointer.To("")},
						{ValidPreparationVesselID: pointer.To("vpv-2")},
					},
				},
			},
		}

		ids := x.GetAllValidPreparationVesselIDs()
		assert.Equal(t, []string{"vpv-1", "vpv-2"}, ids)
	})

	T.Run("returns empty slice when no IDs present", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Vessels: []*RecipeStepVesselDatabaseCreationInput{
						{ValidPreparationVesselID: nil},
					},
				},
			},
		}

		ids := x.GetAllValidPreparationVesselIDs()
		assert.Empty(t, ids)
	})
}

func TestRecipeDatabaseCreationInput_GetRelatedRecipeIDs(T *testing.T) {
	T.Parallel()

	T.Run("returns related recipe IDs from ingredients", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID: "step-1",
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{RecipeStepProductRecipeID: pointer.To("recipe-1")},
						{RecipeStepProductRecipeID: pointer.To("recipe-2")},
					},
				},
				{
					ID: "step-2",
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{RecipeStepProductRecipeID: pointer.To("recipe-3")},
					},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Equal(t, []string{"recipe-1", "recipe-2", "recipe-3"}, ids)
	})

	T.Run("deduplicates recipe IDs", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID: "step-1",
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{RecipeStepProductRecipeID: pointer.To("recipe-1")},
						{RecipeStepProductRecipeID: pointer.To("recipe-2")},
					},
				},
				{
					ID: "step-2",
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{RecipeStepProductRecipeID: pointer.To("recipe-1")}, // duplicate
						{RecipeStepProductRecipeID: pointer.To("recipe-3")},
					},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Equal(t, []string{"recipe-1", "recipe-2", "recipe-3"}, ids)
	})

	T.Run("skips nil and empty values", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID: "step-1",
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{RecipeStepProductRecipeID: pointer.To("recipe-1")},
						{RecipeStepProductRecipeID: nil},
						{RecipeStepProductRecipeID: pointer.To("")},
						{RecipeStepProductRecipeID: pointer.To("recipe-2")},
					},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Equal(t, []string{"recipe-1", "recipe-2"}, ids)
	})

	T.Run("returns empty slice when no related recipes", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID: "step-1",
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{RecipeStepProductRecipeID: nil},
						{Name: "regular ingredient"},
					},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Empty(t, ids)
	})

	T.Run("returns empty slice for recipe with no steps", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{}

		ids := x.GetRelatedRecipeIDs()
		assert.Empty(t, ids)
	})

	T.Run("returns empty slice for steps with no ingredients", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID:          "step-1",
					Ingredients: nil,
				},
				{
					ID:          "step-2",
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Empty(t, ids)
	})

	T.Run("handles mixed ingredients with and without recipe references", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID: "step-1",
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{Name: "salt", RecipeStepProductRecipeID: nil},
						{Name: "caesar breadcrumbs", RecipeStepProductRecipeID: pointer.To("breadcrumbs-recipe-id")},
						{Name: "pepper", RecipeStepProductRecipeID: nil},
					},
				},
				{
					ID: "step-2",
					Ingredients: []*RecipeStepIngredientDatabaseCreationInput{
						{Name: "garlic butter", RecipeStepProductRecipeID: pointer.To("garlic-butter-recipe-id")},
						{Name: "olive oil", RecipeStepProductRecipeID: nil},
					},
				},
			},
		}

		ids := x.GetRelatedRecipeIDs()
		assert.Equal(t, []string{"breadcrumbs-recipe-id", "garlic-butter-recipe-id"}, ids)
	})
}

func TestRecipeDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			ID:            t.Name(),
			Name:          t.Name(),
			CreatedByUser: t.Name(),
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipeUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeUpdateRequestInput{
			Name:               pointer.To(t.Name()),
			Source:             pointer.To(t.Name()),
			Description:        pointer.To(t.Name()),
			InspiredByRecipeID: pointer.To(t.Name()),
			EstimatedPortions:  types.Float32RangeWithOptionalMaxUpdateRequestInput{Min: pointer.To(fake.Float32())},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeUpdateRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}
