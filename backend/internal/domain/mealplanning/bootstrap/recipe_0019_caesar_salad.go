package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// CaesarSaladRecipe creates the Caesar Salad recipe.
// Source: https://www.seriouseats.com/the-best-caesar-salad-recipe
// This returns three recipes: Garlic Croutons (component), Caesar Dressing (component), and Caesar Salad (main recipe).
func CaesarSaladRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {

	// ==================== GARLIC CROUTONS RECIPE ====================

	// Get preparations
	preheatPrep := enums.Preparations["preheat"]
	combinePrep := enums.Preparations["combine"]
	mixPrep := enums.Preparations["mix"]
	addPrep := enums.Preparations["add"]
	tossPrep := enums.Preparations["toss"]
	seasonPrep := enums.Preparations["season"]
	transferPrep := enums.Preparations["transfer"]
	bakePrep := enums.Preparations["bake"]
	coolPrep := enums.Preparations["cool"]
	blendPrep := enums.Preparations["blend"]
	sprinklePrep := enums.Preparations["sprinkle"]
	inspectPrep := enums.Preparations["inspect"]
	rinsePrep := enums.Preparations["rinse"]
	dryPrep := enums.Preparations["dry"]
	mincePrep := enums.Preparations["mince"]
	pressPrep := enums.Preparations["press"]
	cubePrep := enums.Preparations["cube"]
	gratePrep := enums.Preparations["grate"]

	// Get ingredients
	oliveOil := enums.Ingredients["olive oil"]
	garlic := enums.Ingredients["garlic"]
	heartyBread := enums.Ingredients["hearty bread"]
	parmesanCheese := enums.Ingredients["parmesan cheese"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	eggYolk := enums.Ingredients["egg yolk"]
	lemonJuice := enums.Ingredients["lemon juice"]
	anchovies := enums.Ingredients["anchovies"]
	worcestershire := enums.Ingredients["Worcestershire sauce"]
	canolaOil := enums.Ingredients["canola oil"]
	romaineLettuce := enums.Ingredients["romaine lettuce"]

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
	bareHands := enums.Instruments["bare hands"]
	knife := enums.Instruments["knife"]
	cheeseGrater := enums.Instruments["cheese grater"]

	// Get vessels
	smallBowl := enums.Vessels["small bowl"]
	mediumBowl := enums.Vessels["medium bowl"]
	largeBowl := enums.Vessels["large bowl"]
	fineMeshStrainer := enums.Vessels["fine-mesh strainer"]
	bakingSheet := enums.Vessels["baking sheet"]
	oven := enums.Vessels["oven"]
	immersionBlenderCup := enums.Vessels["immersion blender cup"]
	servingBowl := enums.Vessels["serving bowl"]
	saladSpinner := enums.Vessels["salad spinner"]
	cuttingBoard := enums.Vessels["cutting board"]

	// Get ingredient states for completion conditions
	goldenBrownState := enums.IngredientStates["browned"]
	emulsifiedState := enums.IngredientStates["at desired consistency"]

	// === CROUTONS BRIDGE TABLE ENTRIES ===
	// Preheat preparation bridges
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]

	// Whisk preparation bridges (used in dressing recipe)
	whiskOliveOilVIP := enums.IngredientPreparations[mixPrep.ID][oliveOil.ID]
	whiskWhiskVPI := enums.PreparationInstruments[mixPrep.ID][whisk.ID]
	whiskMediumBowlVPV := enums.PreparationVessels[mixPrep.ID][mediumBowl.ID]

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
	tossRomaineVIP := enums.IngredientPreparations[tossPrep.ID][romaineLettuce.ID]
	tossLargeBowlVPV := enums.PreparationVessels[tossPrep.ID][largeBowl.ID]

	// Season preparation bridges
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	seasonLargeBowlVPV := enums.PreparationVessels[seasonPrep.ID][largeBowl.ID]
	seasonMediumBowlVPV := enums.PreparationVessels[seasonPrep.ID][mediumBowl.ID]

	// Transfer preparation bridges
	transferBreadVIP := enums.IngredientPreparations[transferPrep.ID][heartyBread.ID]
	transferBakingSheetVPV := enums.PreparationVessels[transferPrep.ID][bakingSheet.ID]
	transferMediumBowlVPV := enums.PreparationVessels[transferPrep.ID][mediumBowl.ID]
	transferServingBowlVPV := enums.PreparationVessels[transferPrep.ID][servingBowl.ID]

	// Bake preparation bridges
	bakeBreadVIP := enums.IngredientPreparations[bakePrep.ID][heartyBread.ID]
	bakeBakingSheetVPV := enums.PreparationVessels[bakePrep.ID][bakingSheet.ID]
	bakeOvenVPV := enums.PreparationVessels[bakePrep.ID][oven.ID]

	// Cool preparation bridges
	coolBreadVIP := enums.IngredientPreparations[coolPrep.ID][heartyBread.ID]
	coolBakingSheetVPV := enums.PreparationVessels[coolPrep.ID][bakingSheet.ID]

	// Measurement unit bridges
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	oliveOilCupVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][cupMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
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
				RecipeStepProductRecipeID:        nil, // Will be resolved during recipe creation (references Caesar Dressing recipe)
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

	// ==================== CAESAR DRESSING RECIPE ====================

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
	whiskCanolaOilVIP := enums.IngredientPreparations[mixPrep.ID][canolaOil.ID]

	// Measurement unit bridges for dressing
	eggYolkUnitVIMU := enums.IngredientMeasurementUnits[eggYolk.ID][unitMeasurement.ID]
	lemonJuiceTablespoonVIMU := enums.IngredientMeasurementUnits[lemonJuice.ID][tablespoonMeasurement.ID]
	anchoviesUnitVIMU := enums.IngredientMeasurementUnits[anchovies.ID][unitMeasurement.ID]
	worcestershireTeaspoonVIMU := enums.IngredientMeasurementUnits[worcestershire.ID][teaspoonMeasurement.ID]
	parmesanCupVIMU := enums.IngredientMeasurementUnits[parmesanCheese.ID][cupMeasurement.ID]
	canolaOilCupVIMU := enums.IngredientMeasurementUnits[canolaOil.ID][cupMeasurement.ID]

	// Step -1: Mince garlic
	drStepNeg1 := &mealplanning.RecipeStepCreationRequestInput{
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

	// Step -0.5: Combine olive oil with minced garlic and whisk
	drStepNeg05 := &mealplanning.RecipeStepCreationRequestInput{
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

	// Step -0.25: Press garlic through fine-mesh strainer to get pressed garlic
	drStepNeg025 := &mealplanning.RecipeStepCreationRequestInput{
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

	// Step 0: Combine dressing ingredients in immersion blender cup
	drStep0 := &mealplanning.RecipeStepCreationRequestInput{
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

	// Step 1: Transfer mixture to medium bowl
	drStep1 := &mealplanning.RecipeStepCreationRequestInput{
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

	// Step 2: Whisk in remaining olive oil
	drStep2 := &mealplanning.RecipeStepCreationRequestInput{
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

	// Step 3: Season to taste
	drStep3 := &mealplanning.RecipeStepCreationRequestInput{
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
			drStepNeg1, drStepNeg05, drStepNeg025, drStep0, drStep1, drStep2, drStep3,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
		Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
	}

	// ==================== CAESAR SALAD RECIPE ====================

	// Measurement unit bridges for salad
	romaineUnitVIMU := enums.IngredientMeasurementUnits[romaineLettuce.ID][unitMeasurement.ID]

	// Inspect preparation bridges for salad
	inspectRomaineVIP := enums.IngredientPreparations[inspectPrep.ID][romaineLettuce.ID]
	inspectCuttingBoardSaladVPV := enums.PreparationVessels[inspectPrep.ID][cuttingBoard.ID]
	inspectBareHandsSaladVPI := enums.PreparationInstruments[inspectPrep.ID][bareHands.ID]

	// Rinse preparation bridges for salad
	rinseRomaineVIP := enums.IngredientPreparations[rinsePrep.ID][romaineLettuce.ID]
	rinseLargeBowlSaladVPV := enums.PreparationVessels[rinsePrep.ID][largeBowl.ID]

	// Dry preparation bridges for salad
	dryRomaineVIP := enums.IngredientPreparations[dryPrep.ID][romaineLettuce.ID]
	drySaladSpinnerSaladVPV := enums.PreparationVessels[dryPrep.ID][saladSpinner.ID]

	// Sprinkle preparation bridges
	sprinkleParmesanVIP := enums.IngredientPreparations[sprinklePrep.ID][parmesanCheese.ID]
	sprinkleBreadVIP := enums.IngredientPreparations[sprinklePrep.ID][heartyBread.ID]
	sprinkleServingBowlVPV := enums.PreparationVessels[sprinklePrep.ID][servingBowl.ID]

	// Step -2: Select inner romaine leaves
	slStepNeg2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: inspectPrep.ID,
		Index:         0,
		Notes:         "Select the inner romaine leaves, discarding the outer leaves.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &inspectRomaineVIP.ID,
				ValidIngredientMeasurementUnitID: &romaineUnitVIMU.ID,
				Name:                             "romaine lettuce",
				QuantityNotes:                    "select inner leaves only",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &inspectBareHandsSaladVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &inspectCuttingBoardSaladVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "inner romaine leaves",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step -1: Wash romaine leaves
	slStepNeg1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: rinsePrep.ID,
		Index:         1,
		Notes:         "Wash the inner romaine leaves in several changes of water until no dirt or grit remains.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &rinseRomaineVIP.ID,
				Name:                            "inner romaine leaves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &rinseLargeBowlSaladVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "washed inner romaine leaves",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step -0.5: Dry romaine leaves
	slStepNeg05 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: dryPrep.ID,
		Index:         2,
		Notes:         "Carefully dry the washed romaine leaves using a salad spinner or paper towels.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &dryRomaineVIP.ID,
				Name:                            "washed inner romaine leaves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &drySaladSpinnerSaladVPV.ID,
				Name:                     "salad spinner",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "romaine lettuce, inner leaves only, washed and carefully dried",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 0: Toss lettuce with dressing
	slStep0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: tossPrep.ID,
		Index:         3,
		Notes:         "Toss lettuce with a few tablespoons of dressing, adding more if desired. Large leaves should be torn into smaller pieces, smaller leaves left intact.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossRomaineVIP.ID,
				Name:                            "romaine lettuce, inner leaves only, washed and carefully dried",
				QuantityNotes:                   "large leaves torn into smaller pieces, smaller leaves left intact",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				RecipeStepProductRecipeID:        nil, // Will be resolved during recipe creation (references Caesar Dressing recipe)
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				Name:                             "Caesar dressing",
				QuantityNotes:                    "add more if desired",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
					Max: pointer.To[float32](6),
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
				Name:              "dressed lettuce",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Grate parmesan cheese
	slStep1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: gratePrep.ID,
		Index:         4,
		Notes:         "Finely grate 1 cup parmesan cheese.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &grateParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanCupVIMU.ID,
				Name:                             "parmesan cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
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
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Add cheese and croutons, toss again
	slStep2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: tossPrep.ID,
		Index:         5,
		Notes:         "Once lettuce is coated, add half of remaining cheese and three-quarters of croutons and toss again.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossRomaineVIP.ID,
				Name:                            "dressed lettuce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossParmesanVIP.ID,
				Name:                            "finely grated parmesan cheese",
				QuantityNotes:                   "half of remaining cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				RecipeStepProductRecipeID:        nil, // Will be resolved during recipe creation (references Garlic Parmesan Croutons recipe)
				ValidIngredientPreparationID:     &tossBreadVIP.ID,
				ValidIngredientMeasurementUnitID: &breadCupVIMU.ID,
				Name:                             "garlic parmesan croutons",
				QuantityNotes:                    "three-quarters of croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2.25,
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
				Name:              "tossed Caesar salad",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Transfer to serving bowl
	slStep3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: transferPrep.ID,
		Index:         6,
		Notes:         "Transfer to a salad bowl.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "tossed Caesar salad",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &transferServingBowlVPV.ID,
				Name:                     "salad bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "Caesar salad in serving bowl",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Sprinkle with remaining cheese and croutons
	slStep4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: sprinklePrep.ID,
		Index:         7,
		Notes:         "Sprinkle with remaining cheese and croutons. Serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "Caesar salad in serving bowl",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &sprinkleParmesanVIP.ID,
				Name:                            "finely grated parmesan cheese",
				QuantityNotes:                   "remaining cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				RecipeStepProductRecipeID:        nil, // Will be resolved during recipe creation (references Garlic Parmesan Croutons recipe)
				ValidIngredientPreparationID:     &sprinkleBreadVIP.ID,
				ValidIngredientMeasurementUnitID: &breadCupVIMU.ID,
				Name:                             "garlic parmesan croutons",
				QuantityNotes:                    "remaining croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &sprinkleServingBowlVPV.ID,
				Name:                     "salad bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "Caesar salad",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	saladRecipe := &mealplanning.RecipeCreationRequestInput{
		Name:                "Caesar Salad",
		Slug:                "caesar-salad",
		Source:              "https://www.seriouseats.com/the-best-caesar-salad-recipe",
		Description:         "The crowd-pleasing salad of crisp romaine leaves, crunchy croutons, and a creamy, emulsified dressing with just the right amount of anchovy.",
		YieldsComponentType: mealplanning.MealComponentTypesSalad,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepCreationRequestInput{
			slStepNeg2, slStepNeg1, slStepNeg05, slStep0, slStep1, slStep2, slStep3, slStep4,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
		Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		dressingRecipe,
		croutonsRecipe,
		saladRecipe,
	}
}
