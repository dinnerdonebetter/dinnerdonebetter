package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// CaesarDressingRecipe creates the Caesar Dressing recipe.
// Source: https://www.seriouseats.com/the-best-caesar-salad-recipe
func CaesarDressingRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	mincePrep := enums.Preparations["mince"]
	pressPrep := enums.Preparations["press"]
	combinePrep := enums.Preparations["combine"]
	blendPrep := enums.Preparations["blend"]
	mixPrep := enums.Preparations["mix"]
	seasonPrep := enums.Preparations["season"]
	transferPrep := enums.Preparations["transfer"]

	// Get ingredients
	oliveOil := enums.Ingredients["olive oil"]
	garlic := enums.Ingredients["garlic"]
	eggYolk := enums.Ingredients["egg yolk"]
	lemonJuice := enums.Ingredients["lemon juice"]
	anchovies := enums.Ingredients["anchovies"]
	worcestershire := enums.Ingredients["Worcestershire sauce"]
	parmesanCheese := enums.Ingredients["parmesan cheese"]
	canolaOil := enums.Ingredients["canola oil"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	cloveMeasurement := enums.MeasurementUnits["clove"]

	// Get instruments
	whisk := enums.Instruments["whisk"]
	stickBlender := enums.Instruments["immersion blender"]
	spoon := enums.Instruments["spoon"]
	knife := enums.Instruments["knife"]

	// Get vessels
	smallBowl := enums.Vessels["small bowl"]
	mediumBowl := enums.Vessels["medium bowl"]
	largeBowl := enums.Vessels["large bowl"]
	fineMeshStrainer := enums.Vessels["fine-mesh strainer"]
	immersionBlenderCup := enums.Vessels["immersion blender cup"]
	cuttingBoard := enums.Vessels["cutting board"]

	// Get ingredient states for completion conditions
	emulsifiedState := enums.IngredientStates["at desired consistency"]

	// Mince preparation bridges for dressing
	minceGarlicDressingVIP := enums.IngredientPreparations[mincePrep.ID][garlic.ID]
	minceKnifeDressingVPI := enums.PreparationInstruments[mincePrep.ID][knife.ID]
	minceCuttingBoardDressingVPV := enums.PreparationVessels[mincePrep.ID][cuttingBoard.ID]

	// Press preparation bridges for dressing (for pressing garlic through strainer)
	pressGarlicDressingVIP := enums.IngredientPreparations[pressPrep.ID][garlic.ID]
	pressSpoonDressingVPI := enums.PreparationInstruments[pressPrep.ID][spoon.ID]
	pressFineMeshStrainerDressingVPV := enums.PreparationVessels[pressPrep.ID][fineMeshStrainer.ID]
	pressLargeBowlDressingVPV := enums.PreparationVessels[pressPrep.ID][largeBowl.ID]

	// Combine preparation bridges for dressing (for garlic oil mixture)
	combineOliveOilDressingVIP := enums.IngredientPreparations[combinePrep.ID][oliveOil.ID]
	combineGarlicDressingVIP := enums.IngredientPreparations[combinePrep.ID][garlic.ID]
	combineSmallBowlDressingVPV := enums.PreparationVessels[combinePrep.ID][smallBowl.ID]
	combineWhiskDressingVPI := enums.PreparationInstruments[combinePrep.ID][whisk.ID]

	// Blend preparation bridges for dressing
	blendEggYolkVIP := enums.IngredientPreparations[blendPrep.ID][eggYolk.ID]
	blendLemonJuiceVIP := enums.IngredientPreparations[blendPrep.ID][lemonJuice.ID]
	blendAnchoviesVIP := enums.IngredientPreparations[blendPrep.ID][anchovies.ID]
	blendWorcestershireVIP := enums.IngredientPreparations[blendPrep.ID][worcestershire.ID]
	blendGarlicVIP := enums.IngredientPreparations[blendPrep.ID][garlic.ID]
	blendParmesanVIP := enums.IngredientPreparations[blendPrep.ID][parmesanCheese.ID]
	blendCanolaOilVIP := enums.IngredientPreparations[blendPrep.ID][canolaOil.ID]
	blendStickBlenderVPI := enums.PreparationInstruments[blendPrep.ID][stickBlender.ID]
	blendImmersionBlenderCupVPV := enums.PreparationVessels[blendPrep.ID][immersionBlenderCup.ID]

	// Whisk bridges for dressing
	whiskOliveOilVIP := enums.IngredientPreparations[mixPrep.ID][oliveOil.ID]
	whiskCanolaOilVIP := enums.IngredientPreparations[mixPrep.ID][canolaOil.ID]
	whiskWhiskVPI := enums.PreparationInstruments[mixPrep.ID][whisk.ID]
	whiskMediumBowlVPV := enums.PreparationVessels[mixPrep.ID][mediumBowl.ID]

	// Season preparation bridges
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	seasonMediumBowlVPV := enums.PreparationVessels[seasonPrep.ID][mediumBowl.ID]

	// Transfer preparation bridges
	transferMediumBowlVPV := enums.PreparationVessels[transferPrep.ID][mediumBowl.ID]

	// Measurement unit bridges for dressing
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	oliveOilCupVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][cupMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
	eggYolkUnitVIMU := enums.IngredientMeasurementUnits[eggYolk.ID][unitMeasurement.ID]
	lemonJuiceTablespoonVIMU := enums.IngredientMeasurementUnits[lemonJuice.ID][tablespoonMeasurement.ID]
	anchoviesUnitVIMU := enums.IngredientMeasurementUnits[anchovies.ID][unitMeasurement.ID]
	worcestershireTeaspoonVIMU := enums.IngredientMeasurementUnits[worcestershire.ID][teaspoonMeasurement.ID]
	parmesanCupVIMU := enums.IngredientMeasurementUnits[parmesanCheese.ID][cupMeasurement.ID]
	canolaOilCupVIMU := enums.IngredientMeasurementUnits[canolaOil.ID][cupMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	pepperTeaspoonVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]

	// Step 0: Mince garlic
	drStep0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: mincePrep.ID,
		Index:         0,
		Notes:         "Mince the garlic cloves.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &minceGarlicDressingVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				Name:                             "garlic cloves",
				QuantityNotes:                    "about 2 teaspoons when minced",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &minceKnifeDressingVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &minceCuttingBoardDressingVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "minced garlic",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &teaspoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 1: Combine olive oil with minced garlic and whisk
	drStep1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: combinePrep.ID,
		Index:         1,
		Notes:         "In a small bowl, combine 3 tablespoons olive oil with minced garlic and whisk for 30 seconds.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &combineOliveOilDressingVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				Name:                             "extra-virgin olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &combineGarlicDressingVIP.ID,
				Name:                            "minced garlic",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &combineWhiskDressingVPI.ID,
				Name:                         "whisk",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &combineSmallBowlDressingVPV.ID,
				Name:                     "small bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "garlic oil mixture",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Press garlic through fine-mesh strainer to get pressed garlic
	drStep2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: pressPrep.ID,
		Index:         2,
		Notes:         "Transfer garlic oil mixture to a fine-mesh strainer set over a large bowl and press with the back of a spoon to extract as much oil as possible, leaving garlic behind. Reserve pressed garlic for dressing.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &pressGarlicDressingVIP.ID,
				Name:                            "garlic oil mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &pressSpoonDressingVPI.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &pressFineMeshStrainerDressingVPV.ID,
				Name:                     "fine-mesh strainer",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationVesselID: &pressLargeBowlDressingVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "garlic-infused olive oil",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &tablespoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				Name:              "pressed garlic",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             1,
				MeasurementUnitID: &teaspoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 3: Combine dressing ingredients in immersion blender cup
	drStep3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: blendPrep.ID,
		Index:         3,
		Notes:         "Combine egg yolk, lemon juice, anchovies, Worcestershire sauce, reserved pressed garlic, and 1/4 cup parmesan cheese in the bottom of a cup that just fits the head of an immersion blender. With blender running, slowly drizzle in canola oil until a smooth emulsion forms.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &blendEggYolkVIP.ID,
				ValidIngredientMeasurementUnitID: &eggYolkUnitVIMU.ID,
				Name:                             "large egg yolk",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &blendLemonJuiceVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonJuiceTablespoonVIMU.ID,
				Name:                             "fresh lemon juice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &blendAnchoviesVIP.ID,
				ValidIngredientMeasurementUnitID: &anchoviesUnitVIMU.ID,
				Name:                             "anchovies",
				QuantityNotes:                    "amount can vary according to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
					Max: pointer.To[float32](6),
				},
			},
			{
				ValidIngredientPreparationID:     &blendWorcestershireVIP.ID,
				ValidIngredientMeasurementUnitID: &worcestershireTeaspoonVIMU.ID,
				Name:                             "Worcestershire sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &blendGarlicVIP.ID,
				Name:                            "pressed garlic",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &blendParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanCupVIMU.ID,
				Name:                             "finely grated parmesan cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ValidIngredientPreparationID:     &blendCanolaOilVIP.ID,
				ValidIngredientMeasurementUnitID: &canolaOilCupVIMU.ID,
				Name:                             "canola oil",
				QuantityNotes:                    "drizzle in slowly while blending",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.33,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &blendStickBlenderVPI.ID,
				Name:                         "immersion blender",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &blendImmersionBlenderCupVPV.ID,
				Name:                     "immersion blender cup",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: emulsifiedState.ID,
				Notes:             "a smooth emulsion should form",
				Ingredients:       []uint64{6}, // canola oil (ingredient index 6)
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "base emulsion",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Transfer mixture to medium bowl
	drStep4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: transferPrep.ID,
		Index:         4,
		Notes:         "Transfer mixture to a medium bowl.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "base emulsion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &transferMediumBowlVPV.ID,
				Name:                     "medium bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "emulsion in bowl",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Whisk in remaining olive oil
	drStep5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: mixPrep.ID,
		Index:         5,
		Notes:         "Whisking constantly, slowly drizzle in remaining 1/4 cup extra-virgin olive oil.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &whiskCanolaOilVIP.ID,
				Name:                            "emulsion in bowl",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &whiskOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilCupVIMU.ID,
				Name:                             "extra-virgin olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &whiskWhiskVPI.ID,
				Name:                         "whisk",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &whiskMediumBowlVPV.ID,
				Name:                     "medium bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "whisked dressing",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Season to taste
	drStep6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         6,
		Notes:         "Season to taste generously with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "whisked dressing",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "kosher salt",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
				Optional: true,
			},
			{
				ValidIngredientPreparationID:     &seasonPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperTeaspoonVIMU.ID,
				Name:                             "freshly ground black pepper",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &seasonMediumBowlVPV.ID,
				Name:                     "medium bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "Caesar dressing",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	dressingRecipe := &mealplanning.RecipeCreationRequestInput{
		Name:                "Caesar Dressing",
		Slug:                "caesar-dressing",
		Source:              "https://www.seriouseats.com/the-best-caesar-salad-recipe",
		Description:         "A modern emulsified Caesar salad dressing made with egg yolk, anchovies, parmesan, and two oils.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
			Max: pointer.To[float32](8),
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  false,
		Steps: []*mealplanning.RecipeStepCreationRequestInput{
			drStep0, drStep1, drStep2, drStep3, drStep4, drStep5, drStep6,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
		Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		dressingRecipe,
	}
}
