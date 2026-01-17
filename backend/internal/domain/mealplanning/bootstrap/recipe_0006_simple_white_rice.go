package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// SimpleWhiteRiceRecipe creates the Simple White Rice recipe.
// Source: https://www.seriouseats.com/essentials-how-to-cook-rice
func SimpleWhiteRiceRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	rinsePrep := enums.Preparations["rinse"]
	simmerPrep := enums.Preparations["simmer"]
	stirPrep := enums.Preparations["stir"]
	coverPrep := enums.Preparations["cover"]
	removeFromHeatPrep := enums.Preparations["remove from heat"]
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
	smallBowl := enums.Vessels["small bowl"]

	// Get bridge table entries
	// Rinse preparation bridges
	rinseRiceVIP := enums.IngredientPreparations[rinsePrep.ID][rice.ID]
	rinseWaterVIP := enums.IngredientPreparations[rinsePrep.ID][water.ID]
	rinseBowlVPV := enums.PreparationVessels[rinsePrep.ID][smallBowl.ID]

	// Simmer preparation bridges
	simmerRiceVIP := enums.IngredientPreparations[simmerPrep.ID][rice.ID]
	simmerWaterVIP := enums.IngredientPreparations[simmerPrep.ID][water.ID]
	simmerSaltVIP := enums.IngredientPreparations[simmerPrep.ID][salt.ID]
	simmerOliveOilVIP := enums.IngredientPreparations[simmerPrep.ID][oliveOil.ID]
	simmerSaucepanVPV := enums.PreparationVessels[simmerPrep.ID][saucepan.ID]

	// Stir preparation bridges
	stirRiceVIP := enums.IngredientPreparations[stirPrep.ID][rice.ID]
	stirWoodenSpoonVPI := enums.PreparationInstruments[stirPrep.ID][woodenSpoon.ID]

	// Rest preparation bridges
	restRiceVIP := enums.IngredientPreparations[restPrep.ID][rice.ID]

	// Fluff preparation bridges
	fluffRiceVIP := enums.IngredientPreparations[fluffPrep.ID][rice.ID]
	fluffForkVPI := enums.PreparationInstruments[fluffPrep.ID][fork.ID]

	// Measurement unit bridges
	riceCupVIMU := enums.IngredientMeasurementUnits[rice.ID][cupMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	saltPinchVIMU := enums.IngredientMeasurementUnits[salt.ID][pinchMeasurement.ID]
	oliveOilTbspVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]

	// Step 0: Rinse rice until water runs clear
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: rinsePrep.ID,
		Index:         0,
		Notes:         "Rinse the rice in a bowl with cold water, swishing it around with your hand. Drain and repeat until the water runs clear.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &rinseRiceVIP.ID,
				ValidIngredientMeasurementUnitID: &riceCupVIMU.ID,
				Name:                             "long-grain white rice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &rinseWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				Name:                             "cold water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &rinseBowlVPV.ID,
				Name:                     "bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "rinsed rice",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Combine all ingredients in saucepan and bring to a simmer
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: simmerPrep.ID,
		Index:         1,
		Notes:         "Combine all ingredients in a 2-quart saucepan and bring to a simmer.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &simmerRiceVIP.ID,
				ValidIngredientMeasurementUnitID: &riceCupVIMU.ID,
				Name:                             "rinsed rice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &simmerWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1.75, // 1 3/4 cups
				},
			},
			{
				ValidIngredientPreparationID:     &simmerSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltPinchVIMU.ID,
				Name:                             "salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &simmerOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTbspVIMU.ID,
				Name:                             "olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &simmerSaucepanVPV.ID,
				Name:                     "2-quart saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "simmering rice mixture",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2.75),
				},
			},
			{
				Name:  "2-quart saucepan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 2: Stir everything when it reaches a lively simmer
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: stirPrep.ID,
		Index:         2,
		Notes:         "As soon as the water reaches a lively simmer, give everything a good stir.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &stirRiceVIP.ID,
				ValidIngredientMeasurementUnitID: &riceCupVIMU.ID,
				Name:                             "simmering rice mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2.75,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &stirWoodenSpoonVPI.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "2-quart saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "stirred rice mixture",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2.75),
				},
			},
			{
				Name:  "2-quart saucepan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 3: Cover the pot and lower heat
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: coverPrep.ID,
		Index:         3,
		Notes:         "Cover the pot and lower the heat as much as possible. Cook for 15 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "stirred rice mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2.75,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "2-quart saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cooked rice in covered pot",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				Name:  "2-quart saucepan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 4: Remove from heat
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: removeFromHeatPrep.ID,
		Index:         4,
		Notes:         "Turn off the burner and remove the pot from the heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "cooked rice in covered pot",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "2-quart saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cooked rice removed from heat",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				Name:  "2-quart saucepan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 5: Rest
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: restPrep.ID,
		Index:         5,
		Notes:         "Let the pot sit for at least 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &restRiceVIP.ID,
				ValidIngredientMeasurementUnitID: &riceCupVIMU.ID,
				Name:                             "cooked rice removed from heat",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "2-quart saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "rested rice",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				Name:  "2-quart saucepan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 6: Fluff with fork and serve
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: fluffPrep.ID,
		Index:         6,
		Notes:         "Fluff with a fork and serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &fluffRiceVIP.ID,
				ValidIngredientMeasurementUnitID: &riceCupVIMU.ID,
				Name:                             "rested rice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &fluffForkVPI.ID,
				Name:                         "fork",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "2-quart saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "simple white rice",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
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
			Steps:             []*mealplanning.RecipeStepCreationRequestInput{step0, step1, step2, step3, step4, step5, step6},
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
