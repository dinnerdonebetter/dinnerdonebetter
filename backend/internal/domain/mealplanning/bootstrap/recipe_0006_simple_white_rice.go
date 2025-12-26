package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// SimpleWhiteRiceRecipe creates the Simple White Rice recipe.
// Source: https://www.seriouseats.com/essentials-how-to-cook-rice
func SimpleWhiteRiceRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	simmerPrep := enums.Preparations["simmer"]
	stirPrep := enums.Preparations["stir"]
	coverPrep := enums.Preparations["cover"]
	restPrep := enums.Preparations["rest"]
	fluffPrep := enums.Preparations["fluff"]

	// Get ingredients
	rice := enums.Ingredients["rice"]
	water := enums.Ingredients["water"]
	salt := enums.Ingredients["salt"]
	oliveOil := enums.Ingredients["olive oil"]

	// Get measurement units
	cupMeasurement := enums.MeasurementUnits["cup"]
	pinchMeasurement := enums.MeasurementUnits["pinch"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]

	// Get instruments
	fork := enums.Instruments["fork"]
	woodenSpoon := enums.Instruments["wooden spoon"]

	// Get vessels
	saucepan := enums.Vessels["saucepan"]

	// Get bridge table entries
	// Simmer preparation bridges
	simmerRiceVIP := enums.IngredientPreparations[simmerPrep.ID][rice.ID]
	simmerWaterVIP := enums.IngredientPreparations[simmerPrep.ID][water.ID]
	simmerSaltVIP := enums.IngredientPreparations[simmerPrep.ID][salt.ID]
	simmerOliveOilVIP := enums.IngredientPreparations[simmerPrep.ID][oliveOil.ID]
	simmerSaucepanVPV := enums.PreparationVessels[simmerPrep.ID][saucepan.ID]

	// Stir preparation bridges
	stirRiceVIP := enums.IngredientPreparations[stirPrep.ID][rice.ID]
	stirWoodenSpoonVPI := enums.PreparationInstruments[stirPrep.ID][woodenSpoon.ID]
	stirSaucepanVPV := enums.PreparationVessels[stirPrep.ID][saucepan.ID]

	// Cover preparation bridges
	coverSaucepanVPV := enums.PreparationVessels[coverPrep.ID][saucepan.ID]

	// Rest preparation bridges
	restRiceVIP := enums.IngredientPreparations[restPrep.ID][rice.ID]
	restSaucepanVPV := enums.PreparationVessels[restPrep.ID][saucepan.ID]

	// Fluff preparation bridges
	fluffRiceVIP := enums.IngredientPreparations[fluffPrep.ID][rice.ID]
	fluffForkVPI := enums.PreparationInstruments[fluffPrep.ID][fork.ID]
	fluffSaucepanVPV := enums.PreparationVessels[fluffPrep.ID][saucepan.ID]

	// Measurement unit bridges
	riceCupVIMU := enums.IngredientMeasurementUnits[rice.ID][cupMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	saltPinchVIMU := enums.IngredientMeasurementUnits[salt.ID][pinchMeasurement.ID]
	oliveOilTbspVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]

	// Step 0: Combine all ingredients in saucepan and bring to a simmer
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   simmerPrep.ID,
		Index:           0,
		Notes:           "Combine all ingredients in a 2-quart saucepan and bring to a simmer.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &simmerRiceVIP.ID,
				ValidIngredientMeasurementUnitID: &riceCupVIMU.ID,
				IngredientID:                     &rice.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "long-grain white rice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &simmerWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				IngredientID:                     &water.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1.75, // 1 3/4 cups
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &simmerSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltPinchVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                pinchMeasurement.ID,
				Name:                             "salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &simmerOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTbspVIMU.ID,
				IngredientID:                     &oliveOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &simmerSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "2-quart saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "simmering rice mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2.75),
				},
			},
		},
	}

	// Step 1: Stir everything when it reaches a lively simmer
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           1,
		Notes:           "As soon as the water reaches a lively simmer, give everything a good stir.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &stirRiceVIP.ID,
				ValidIngredientMeasurementUnitID: &riceCupVIMU.ID,
				IngredientID:                     &rice.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "simmering rice mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2.75,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidPreparationInstrumentID: &stirWoodenSpoonVPI.ID,
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step1ID,
				ValidPreparationVesselID: &stirSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "stirred rice mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2.75),
				},
			},
		},
	}

	// Step 2: Cover the pot and lower heat
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   coverPrep.ID,
		Index:           2,
		Notes:           "Cover the pot and lower the heat as much as possible. Cook for 15 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &rice.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "stirred rice mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2.75,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &coverSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "cooked rice in covered pot",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 3: Remove from heat and rest
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           3,
		Notes:           "Turn off the burner, remove the pot from the heat, and let it sit for at least 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &restRiceVIP.ID,
				ValidIngredientMeasurementUnitID: &riceCupVIMU.ID,
				IngredientID:                     &rice.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "cooked rice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &restSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "rested rice",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 4: Fluff with fork and serve
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   fluffPrep.ID,
		Index:           4,
		Notes:           "Fluff with a fork and serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &fluffRiceVIP.ID,
				ValidIngredientMeasurementUnitID: &riceCupVIMU.ID,
				IngredientID:                     &rice.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "rested rice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &fluffForkVPI.ID,
				InstrumentID:                 &fork.ID,
				Name:                         "fork",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &fluffSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "simple white rice",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		{
			ID:                  recipeID,
			CreatedByUser:       userID,
			Name:                "Simple White Rice",
			Slug:                "simple-white-rice",
			Source:              "https://www.seriouseats.com/essentials-how-to-cook-rice",
			Description:         "",
			YieldsComponentType: mealplanning.MealComponentTypesSide,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 3,
				Max: pointer.To[float32](4),
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
				step0, step1, step2, step3, step4,
			},
		},
	}
}
