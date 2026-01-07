package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// GarlicParmesanCroutonsRecipe creates the Garlic Parmesan Croutons recipe.
// Source: https://www.seriouseats.com/the-best-caesar-salad-recipe
// Note: This recipe references the Caesar Dressing recipe for "garlic-infused olive oil", which must be created first.
func GarlicParmesanCroutonsRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	preheatPrep := enums.Preparations["preheat"]
	addPrep := enums.Preparations["add"]
	tossPrep := enums.Preparations["toss"]
	seasonPrep := enums.Preparations["season"]
	transferPrep := enums.Preparations["transfer"]
	bakePrep := enums.Preparations["bake"]
	coolPrep := enums.Preparations["cool"]
	cubePrep := enums.Preparations["cube"]
	gratePrep := enums.Preparations["grate"]

	// Get ingredients
	oliveOil := enums.Ingredients["olive oil"]
	heartyBread := enums.Ingredients["hearty bread"]
	parmesanCheese := enums.Ingredients["parmesan cheese"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	knife := enums.Instruments["knife"]
	cheeseGrater := enums.Instruments["cheese grater"]

	// Get vessels
	largeBowl := enums.Vessels["large bowl"]
	bakingSheet := enums.Vessels["baking sheet"]
	oven := enums.Vessels["oven"]
	cuttingBoard := enums.Vessels["cutting board"]

	// Get ingredient states for completion conditions
	goldenBrownState := enums.IngredientStates["browned"]

	// === CROUTONS BRIDGE TABLE ENTRIES ===
	// Preheat preparation bridges
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]

	// Cube preparation bridges
	cubeBreadVIP := enums.IngredientPreparations[cubePrep.ID][heartyBread.ID]
	cubeKnifeVPI := enums.PreparationInstruments[cubePrep.ID][knife.ID]
	cubeCuttingBoardVPV := enums.PreparationVessels[cubePrep.ID][cuttingBoard.ID]

	// Grate preparation bridges
	grateParmesanVIP := enums.IngredientPreparations[gratePrep.ID][parmesanCheese.ID]
	grateCheeseGraterVPI := enums.PreparationInstruments[gratePrep.ID][cheeseGrater.ID]
	grateCuttingBoardVPV := enums.PreparationVessels[gratePrep.ID][cuttingBoard.ID]

	// Add preparation bridges
	addBreadVIP := enums.IngredientPreparations[addPrep.ID][heartyBread.ID]
	addLargeBowlVPV := enums.PreparationVessels[addPrep.ID][largeBowl.ID]

	// Toss preparation bridges
	tossBreadVIP := enums.IngredientPreparations[tossPrep.ID][heartyBread.ID]
	tossParmesanVIP := enums.IngredientPreparations[tossPrep.ID][parmesanCheese.ID]
	tossLargeBowlVPV := enums.PreparationVessels[tossPrep.ID][largeBowl.ID]

	// Season preparation bridges
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	seasonLargeBowlVPV := enums.PreparationVessels[seasonPrep.ID][largeBowl.ID]

	// Transfer preparation bridges
	transferBreadVIP := enums.IngredientPreparations[transferPrep.ID][heartyBread.ID]
	transferBakingSheetVPV := enums.PreparationVessels[transferPrep.ID][bakingSheet.ID]

	// Bake preparation bridges
	bakeBreadVIP := enums.IngredientPreparations[bakePrep.ID][heartyBread.ID]
	bakeBakingSheetVPV := enums.PreparationVessels[bakePrep.ID][bakingSheet.ID]
	bakeOvenVPV := enums.PreparationVessels[bakePrep.ID][oven.ID]

	// Cool preparation bridges
	coolBreadVIP := enums.IngredientPreparations[coolPrep.ID][heartyBread.ID]
	coolBakingSheetVPV := enums.PreparationVessels[coolPrep.ID][bakingSheet.ID]

	// Measurement unit bridges
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	breadCupVIMU := enums.IngredientMeasurementUnits[heartyBread.ID][cupMeasurement.ID]
	parmesanTablespoonVIMU := enums.IngredientMeasurementUnits[parmesanCheese.ID][tablespoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	pepperTeaspoonVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]

	// ==================== CROUTONS RECIPE STEPS ====================

	// Step 0: Preheat oven to 375°F
	crStep0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: preheatPrep.ID,
		Index:         0,
		Notes:         "Adjust oven rack to middle position and preheat oven to 375°F (190°C).",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](190),
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &preheatOvenVPV.ID,
				Name:                     "oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "preheated oven",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	// Step 1: Cut bread into 3/4-inch cubes
	crStep1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: cubePrep.ID,
		Index:         1,
		Notes:         "Cut the hearty bread into 3/4-inch cubes.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &cubeBreadVIP.ID,
				ValidIngredientMeasurementUnitID: &breadCupVIMU.ID,
				Name:                             "hearty bread",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &cubeKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &cubeCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "hearty bread, cut into 3/4-inch cubes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 2: Add bread cubes to garlic-infused olive oil and toss to coat
	crStep2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: addPrep.ID,
		Index:         2,
		Notes:         "Add bread cubes to garlic-infused olive oil from the dressing recipe and toss to coat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				// RecipeStepProductRecipeID should reference the "Caesar Dressing" recipe (slug: "caesar-dressing")
				// This needs to be resolved by looking up the recipe by name or slug during recipe creation
				// The product should be "garlic-infused olive oil" from step 2 (index 2), product index 0
				RecipeStepProductRecipeID:        nil,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				Name:                             "garlic-infused olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addBreadVIP.ID,
				Name:                            "hearty bread, cut into 3/4-inch cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &addLargeBowlVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "oiled bread cubes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Grate parmesan cheese
	crStep3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: gratePrep.ID,
		Index:         3,
		Notes:         "Finely grate 4 tablespoons parmesan cheese.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &grateParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanTablespoonVIMU.ID,
				Name:                             "parmesan cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &grateCheeseGraterVPI.ID,
				Name:                         "cheese grater",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &grateCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "finely grated parmesan cheese",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &tablespoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 4: Add parmesan, toss, and season
	crStep4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: tossPrep.ID,
		Index:         4,
		Notes:         "Add 2 tablespoons grated parmesan cheese, toss again, and season to taste with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossBreadVIP.ID,
				Name:                            "oiled bread cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossParmesanVIP.ID,
				Name:                            "finely grated parmesan cheese",
				QuantityNotes:                   "2 tablespoons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &tossLargeBowlVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cheesy bread cubes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Season bread cubes with salt and pepper
	crStep5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         5,
		Notes:         "Season to taste with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "cheesy bread cubes",
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
				ValidPreparationVesselID: &seasonLargeBowlVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seasoned bread cubes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Transfer to rimmed baking sheet
	crStep6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: transferPrep.ID,
		Index:         6,
		Notes:         "Transfer to a rimmed baking sheet.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferBreadVIP.ID,
				Name:                            "seasoned bread cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &transferBakingSheetVPV.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "bread cubes on baking sheet",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Bake until pale golden brown and crisp
	crStep7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: bakePrep.ID,
		Index:         7,
		Notes:         "Bake until croutons are pale golden brown and crisp, about 15 minutes.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](190),
		},
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &bakeBreadVIP.ID,
				Name:                            "bread cubes on baking sheet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &bakeBakingSheetVPV.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &bakeOvenVPV.ID,
				Name:                            "preheated oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: goldenBrownState.ID,
				Notes:             "croutons should be pale golden brown and crisp",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "baked croutons",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Toss with more parmesan
	crStep8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: tossPrep.ID,
		Index:         8,
		Notes:         "Remove from oven and toss with 2 more tablespoons grated parmesan.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossBreadVIP.ID,
				Name:                            "baked croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossParmesanVIP.ID,
				Name:                            "finely grated parmesan cheese",
				QuantityNotes:                   "2 tablespoons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &tossLargeBowlVPV.ID,
				Name:                     "baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "parmesan croutons",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Allow to cool
	crStep9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: coolPrep.ID,
		Index:         9,
		Notes:         "Allow croutons to cool.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &coolBreadVIP.ID,
				Name:                            "parmesan croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &coolBakingSheetVPV.ID,
				Name:                     "baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "garlic parmesan croutons",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	croutonsRecipe := &mealplanning.RecipeCreationRequestInput{
		Name:                "Garlic Parmesan Croutons",
		Slug:                "garlic-parmesan-croutons",
		Source:              "https://www.seriouseats.com/the-best-caesar-salad-recipe",
		Description:         "Homemade croutons with garlic-infused olive oil and parmesan cheese, perfect for Caesar salad.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "cup",
		PluralPortionName: "cups",
		EligibleForMeals:  false,
		Steps: []*mealplanning.RecipeStepCreationRequestInput{
			crStep0, crStep1, crStep2, crStep3, crStep4, crStep5, crStep6, crStep7, crStep8, crStep9,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
		Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		croutonsRecipe,
	}
}
