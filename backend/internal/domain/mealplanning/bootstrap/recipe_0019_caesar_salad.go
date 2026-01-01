package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// CaesarSaladRecipe creates the Caesar Salad recipe.
// Source: https://www.seriouseats.com/the-best-caesar-salad-recipe
// This returns three recipes: Garlic Croutons (component), Caesar Dressing (component), and Caesar Salad (main recipe).
func CaesarSaladRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	// ==================== GARLIC CROUTONS RECIPE ====================
	croutonsRecipeID := identifiers.New()

	// Get preparations
	preheatPrep := enums.Preparations["preheat"]
	combinePrep := enums.Preparations["combine"]
	mixPrep := enums.Preparations["mix"]
	strainPrep := enums.Preparations["strain"]
	addPrep := enums.Preparations["add"]
	tossPrep := enums.Preparations["toss"]
	seasonPrep := enums.Preparations["season"]
	transferPrep := enums.Preparations["transfer"]
	bakePrep := enums.Preparations["bake"]
	coolPrep := enums.Preparations["cool"]
	blendPrep := enums.Preparations["blend"]
	sprinklePrep := enums.Preparations["sprinkle"]

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
	stickBlender := enums.Instruments["stick blender"]
	spoon := enums.Instruments["spoon"]

	// Get vessels
	smallBowl := enums.Vessels["small bowl"]
	mediumBowl := enums.Vessels["medium bowl"]
	largeBowl := enums.Vessels["large bowl"]
	fineMeshStrainer := enums.Vessels["fine-mesh strainer"]
	bakingSheet := enums.Vessels["baking sheet"]
	oven := enums.Vessels["oven"]
	immersionBlenderCup := enums.Vessels["immersion blender cup"]
	servingBowl := enums.Vessels["serving bowl"]

	// Get ingredient states for completion conditions
	goldenBrownState := enums.IngredientStates["browned"]
	emulsifiedState := enums.IngredientStates["at desired consistency"]

	// === CROUTONS BRIDGE TABLE ENTRIES ===
	// Preheat preparation bridges
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]

	// Combine preparation bridges
	combineOliveOilVIP := enums.IngredientPreparations[combinePrep.ID][oliveOil.ID]
	combineGarlicVIP := enums.IngredientPreparations[combinePrep.ID][garlic.ID]
	combineSmallBowlVPV := enums.PreparationVessels[combinePrep.ID][smallBowl.ID]
	combineWhiskVPI := enums.PreparationInstruments[combinePrep.ID][whisk.ID]

	// Whisk preparation bridges
	whiskOliveOilVIP := enums.IngredientPreparations[mixPrep.ID][oliveOil.ID]
	whiskWhiskVPI := enums.PreparationInstruments[mixPrep.ID][whisk.ID]
	whiskMediumBowlVPV := enums.PreparationVessels[mixPrep.ID][mediumBowl.ID]

	// Strain preparation bridges
	strainOliveOilVIP := enums.IngredientPreparations[strainPrep.ID][oliveOil.ID]
	strainSpoonVPI := enums.PreparationInstruments[strainPrep.ID][spoon.ID]
	strainFineMeshStrainerVPV := enums.PreparationVessels[strainPrep.ID][fineMeshStrainer.ID]
	strainLargeBowlVPV := enums.PreparationVessels[strainPrep.ID][largeBowl.ID]

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
	garlicTeaspoonVIMU := enums.IngredientMeasurementUnits[garlic.ID][teaspoonMeasurement.ID]
	breadCupVIMU := enums.IngredientMeasurementUnits[heartyBread.ID][cupMeasurement.ID]
	parmesanTablespoonVIMU := enums.IngredientMeasurementUnits[parmesanCheese.ID][tablespoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	pepperTeaspoonVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]

	// ==================== CROUTONS RECIPE STEPS ====================

	// Step 0: Preheat oven to 375°F
	crStep0ID := identifiers.New()
	crStep0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep0ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   preheatPrep.ID,
		Index:           0,
		Notes:           "Adjust oven rack to middle position and preheat oven to 375°F (190°C).",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](190),
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep0ID,
				ValidPreparationVesselID: &preheatOvenVPV.ID,
				VesselID:                 &oven.ID,
				Name:                     "oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep0ID,
				Name:                "preheated oven",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
			},
		},
	}

	// Step 1: Combine olive oil with minced garlic and whisk
	crStep1ID := identifiers.New()
	crStep1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep1ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   combinePrep.ID,
		Index:           1,
		Notes:           "In a small bowl, combine 3 tablespoons olive oil with minced garlic and whisk for 30 seconds.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              crStep1ID,
				ValidIngredientPreparationID:     &combineOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				IngredientID:                     &oliveOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "extra-virgin olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              crStep1ID,
				ValidIngredientPreparationID:     &combineGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                cloveMeasurement.ID,
				Name:                             "garlic, minced",
				QuantityNotes:                    "about 2 teaspoons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          crStep1ID,
				ValidPreparationInstrumentID: &combineWhiskVPI.ID,
				InstrumentID:                 &whisk.ID,
				Name:                         "whisk",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep1ID,
				ValidPreparationVesselID: &combineSmallBowlVPV.ID,
				VesselID:                 &smallBowl.ID,
				Name:                     "small bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep1ID,
				Name:                "garlic oil mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Strain garlic oil through fine-mesh strainer
	crStep2ID := identifiers.New()
	crStep2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep2ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   strainPrep.ID,
		Index:           2,
		Notes:           "Transfer garlic oil mixture to a fine-mesh strainer set over a large bowl and press with the back of a spoon to extract as much oil as possible, leaving garlic behind. Reserve pressed garlic separately.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &strainOliveOilVIP.ID,
				IngredientID:                    &oliveOil.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "garlic oil mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          crStep2ID,
				ValidPreparationInstrumentID: &strainSpoonVPI.ID,
				InstrumentID:                 &spoon.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep2ID,
				ValidPreparationVesselID: &strainFineMeshStrainerVPV.ID,
				VesselID:                 &fineMeshStrainer.ID,
				Name:                     "fine-mesh strainer",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep2ID,
				ValidPreparationVesselID: &strainLargeBowlVPV.ID,
				VesselID:                 &largeBowl.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep2ID,
				Name:                "garlic-infused oil",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep2ID,
				Name:                "pressed garlic",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               1,
				MeasurementUnitID:   &teaspoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 3: Add bread cubes to garlic oil and toss to coat
	crStep3ID := identifiers.New()
	crStep3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep3ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   addPrep.ID,
		Index:           3,
		Notes:           "Add bread cubes to garlic oil and toss to coat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &oliveOil.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "garlic-infused oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              crStep3ID,
				ValidIngredientPreparationID:     &addBreadVIP.ID,
				ValidIngredientMeasurementUnitID: &breadCupVIMU.ID,
				IngredientID:                     &heartyBread.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "hearty bread, cut into 3/4-inch cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &addLargeBowlVPV.ID,
				VesselID:                        &largeBowl.ID,
				Name:                            "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep3ID,
				Name:                "oiled bread cubes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Add parmesan, toss, and season
	crStep4ID := identifiers.New()
	crStep4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep4ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   tossPrep.ID,
		Index:           4,
		Notes:           "Add 2 tablespoons parmesan cheese, toss again, and season to taste with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossBreadVIP.ID,
				IngredientID:                    &heartyBread.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "oiled bread cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              crStep4ID,
				ValidIngredientPreparationID:     &tossParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanTablespoonVIMU.ID,
				IngredientID:                     &parmesanCheese.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "finely grated parmesan cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep4ID,
				ValidPreparationVesselID: &tossLargeBowlVPV.ID,
				VesselID:                 &largeBowl.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep4ID,
				Name:                "cheesy bread cubes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Season bread cubes with salt and pepper
	crStep5ID := identifiers.New()
	crStep5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep5ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   seasonPrep.ID,
		Index:           5,
		Notes:           "Season to taste with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &heartyBread.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "cheesy bread cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              crStep5ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "kosher salt",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              crStep5ID,
				ValidIngredientPreparationID:     &seasonPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperTeaspoonVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "freshly ground black pepper",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep5ID,
				ValidPreparationVesselID: &seasonLargeBowlVPV.ID,
				VesselID:                 &largeBowl.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep5ID,
				Name:                "seasoned bread cubes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Transfer to rimmed baking sheet
	crStep6ID := identifiers.New()
	crStep6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep6ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   transferPrep.ID,
		Index:           6,
		Notes:           "Transfer to a rimmed baking sheet.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferBreadVIP.ID,
				IngredientID:                    &heartyBread.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "seasoned bread cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep6ID,
				ValidPreparationVesselID: &transferBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep6ID,
				Name:                "bread cubes on baking sheet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Bake until pale golden brown and crisp
	crStep7ID := identifiers.New()
	crStep7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep7ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   bakePrep.ID,
		Index:           7,
		Notes:           "Bake until croutons are pale golden brown and crisp, about 15 minutes.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](190),
		},
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &bakeBreadVIP.ID,
				IngredientID:                    &heartyBread.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "bread cubes on baking sheet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep7ID,
				ValidPreparationVesselID: &bakeBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &bakeOvenVPV.ID,
				VesselID:                        &oven.ID,
				Name:                            "preheated oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep7ID,
				IngredientStateID:   goldenBrownState.ID,
				Notes:               "croutons should be pale golden brown and crisp",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep7ID,
				Name:                "baked croutons",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Toss with more parmesan
	crStep8ID := identifiers.New()
	crStep8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep8ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   tossPrep.ID,
		Index:           8,
		Notes:           "Remove from oven and toss with 2 more tablespoons parmesan.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossBreadVIP.ID,
				IngredientID:                    &heartyBread.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "baked croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              crStep8ID,
				ValidIngredientPreparationID:     &tossParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanTablespoonVIMU.ID,
				IngredientID:                     &parmesanCheese.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "finely grated parmesan cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep8ID,
				ValidPreparationVesselID: &tossLargeBowlVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep8ID,
				Name:                "parmesan croutons",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Allow to cool
	crStep9ID := identifiers.New()
	crStep9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              crStep9ID,
		BelongsToRecipe: croutonsRecipeID,
		PreparationID:   coolPrep.ID,
		Index:           9,
		Notes:           "Allow croutons to cool.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             crStep9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &coolBreadVIP.ID,
				IngredientID:                    &heartyBread.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "parmesan croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      crStep9ID,
				ValidPreparationVesselID: &coolBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: crStep9ID,
				Name:                "garlic parmesan croutons",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	croutonsRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  croutonsRecipeID,
		CreatedByUser:       userID,
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
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			crStep0, crStep1, crStep2, crStep3, crStep4, crStep5, crStep6, crStep7, crStep8, crStep9,
		},
	}

	// ==================== CAESAR DRESSING RECIPE ====================
	dressingRecipeID := identifiers.New()

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

	// Step 0: Combine dressing ingredients in immersion blender cup
	drStep0ID := identifiers.New()
	drStep0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              drStep0ID,
		BelongsToRecipe: dressingRecipeID,
		PreparationID:   blendPrep.ID,
		Index:           0,
		Notes:           "Combine egg yolk, lemon juice, anchovies, Worcestershire sauce, pressed garlic, and 1/4 cup parmesan cheese in the bottom of a cup that just fits the head of an immersion blender. With blender running, slowly drizzle in canola oil until a smooth emulsion forms.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep0ID,
				ValidIngredientPreparationID:     &blendEggYolkVIP.ID,
				ValidIngredientMeasurementUnitID: &eggYolkUnitVIMU.ID,
				IngredientID:                     &eggYolk.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "large egg yolk",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep0ID,
				ValidIngredientPreparationID:     &blendLemonJuiceVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonJuiceTablespoonVIMU.ID,
				IngredientID:                     &lemonJuice.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "fresh lemon juice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep0ID,
				ValidIngredientPreparationID:     &blendAnchoviesVIP.ID,
				ValidIngredientMeasurementUnitID: &anchoviesUnitVIMU.ID,
				IngredientID:                     &anchovies.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "anchovies",
				QuantityNotes:                    "amount can vary according to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
					Max: pointer.To[float32](6),
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep0ID,
				ValidIngredientPreparationID:     &blendWorcestershireVIP.ID,
				ValidIngredientMeasurementUnitID: &worcestershireTeaspoonVIMU.ID,
				IngredientID:                     &worcestershire.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Worcestershire sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep0ID,
				ValidIngredientPreparationID:     &blendGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicTeaspoonVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "pressed garlic (reserved from croutons)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep0ID,
				ValidIngredientPreparationID:     &blendParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanCupVIMU.ID,
				IngredientID:                     &parmesanCheese.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "finely grated parmesan cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep0ID,
				ValidIngredientPreparationID:     &blendCanolaOilVIP.ID,
				ValidIngredientMeasurementUnitID: &canolaOilCupVIMU.ID,
				IngredientID:                     &canolaOil.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "canola oil",
				QuantityNotes:                    "drizzle in slowly while blending",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.33,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          drStep0ID,
				ValidPreparationInstrumentID: &blendStickBlenderVPI.ID,
				InstrumentID:                 &stickBlender.ID,
				Name:                         "immersion blender",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      drStep0ID,
				ValidPreparationVesselID: &blendImmersionBlenderCupVPV.ID,
				VesselID:                 &immersionBlenderCup.ID,
				Name:                     "immersion blender cup",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: drStep0ID,
				IngredientStateID:   emulsifiedState.ID,
				Notes:               "a smooth emulsion should form",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: drStep0ID,
				Name:                "base emulsion",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Transfer mixture to medium bowl
	drStep1ID := identifiers.New()
	drStep1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              drStep1ID,
		BelongsToRecipe: dressingRecipeID,
		PreparationID:   transferPrep.ID,
		Index:           1,
		Notes:           "Transfer mixture to a medium bowl.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             drStep1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &canolaOil.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "base emulsion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      drStep1ID,
				ValidPreparationVesselID: &transferMediumBowlVPV.ID,
				VesselID:                 &mediumBowl.ID,
				Name:                     "medium bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: drStep1ID,
				Name:                "emulsion in bowl",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Whisk in remaining olive oil
	drStep2ID := identifiers.New()
	drStep2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              drStep2ID,
		BelongsToRecipe: dressingRecipeID,
		PreparationID:   mixPrep.ID,
		Index:           2,
		Notes:           "Whisking constantly, slowly drizzle in remaining 1/4 cup extra-virgin olive oil.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             drStep2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &whiskCanolaOilVIP.ID,
				IngredientID:                    &canolaOil.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "emulsion in bowl",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep2ID,
				ValidIngredientPreparationID:     &whiskOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilCupVIMU.ID,
				IngredientID:                     &oliveOil.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "extra-virgin olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          drStep2ID,
				ValidPreparationInstrumentID: &whiskWhiskVPI.ID,
				InstrumentID:                 &whisk.ID,
				Name:                         "whisk",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      drStep2ID,
				ValidPreparationVesselID: &whiskMediumBowlVPV.ID,
				VesselID:                 &mediumBowl.ID,
				Name:                     "medium bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: drStep2ID,
				Name:                "whisked dressing",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Season to taste
	drStep3ID := identifiers.New()
	drStep3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              drStep3ID,
		BelongsToRecipe: dressingRecipeID,
		PreparationID:   seasonPrep.ID,
		Index:           3,
		Notes:           "Season to taste generously with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             drStep3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &canolaOil.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "whisked dressing",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep3ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "kosher salt",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              drStep3ID,
				ValidIngredientPreparationID:     &seasonPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperTeaspoonVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "freshly ground black pepper",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      drStep3ID,
				ValidPreparationVesselID: &seasonMediumBowlVPV.ID,
				VesselID:                 &mediumBowl.ID,
				Name:                     "medium bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: drStep3ID,
				Name:                "Caesar dressing",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	dressingRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  dressingRecipeID,
		CreatedByUser:       userID,
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
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			drStep0, drStep1, drStep2, drStep3,
		},
	}

	// ==================== CAESAR SALAD RECIPE ====================
	saladRecipeID := identifiers.New()

	// Measurement unit bridges for salad
	romaineUnitVIMU := enums.IngredientMeasurementUnits[romaineLettuce.ID][unitMeasurement.ID]

	// Sprinkle preparation bridges
	sprinkleParmesanVIP := enums.IngredientPreparations[sprinklePrep.ID][parmesanCheese.ID]
	sprinkleBreadVIP := enums.IngredientPreparations[sprinklePrep.ID][heartyBread.ID]
	sprinkleServingBowlVPV := enums.PreparationVessels[sprinklePrep.ID][servingBowl.ID]

	// Step 0: Toss lettuce with dressing
	slStep0ID := identifiers.New()
	slStep0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              slStep0ID,
		BelongsToRecipe: saladRecipeID,
		PreparationID:   tossPrep.ID,
		Index:           0,
		Notes:           "Toss lettuce with a few tablespoons of dressing, adding more if desired.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              slStep0ID,
				ValidIngredientPreparationID:     &tossRomaineVIP.ID,
				ValidIngredientMeasurementUnitID: &romaineUnitVIMU.ID,
				IngredientID:                     &romaineLettuce.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "romaine lettuce, inner leaves only, washed and carefully dried",
				QuantityNotes:                    "large leaves torn into smaller pieces, smaller leaves left intact",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              slStep0ID,
				RecipeStepProductRecipeID:        &dressingRecipeID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				IngredientID:                     &oliveOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "Caesar dressing",
				QuantityNotes:                    "add more if desired",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
					Max: pointer.To[float32](6),
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      slStep0ID,
				ValidPreparationVesselID: &tossLargeBowlVPV.ID,
				VesselID:                 &largeBowl.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: slStep0ID,
				Name:                "dressed lettuce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Add cheese and croutons, toss again
	slStep1ID := identifiers.New()
	slStep1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              slStep1ID,
		BelongsToRecipe: saladRecipeID,
		PreparationID:   tossPrep.ID,
		Index:           1,
		Notes:           "Once lettuce is coated, add half of remaining cheese and three-quarters of croutons and toss again.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             slStep1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossRomaineVIP.ID,
				IngredientID:                    &romaineLettuce.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "dressed lettuce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              slStep1ID,
				ValidIngredientPreparationID:     &tossParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanCupVIMU.ID,
				IngredientID:                     &parmesanCheese.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "finely grated parmesan cheese",
				QuantityNotes:                    "half of remaining cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              slStep1ID,
				RecipeStepProductRecipeID:        &croutonsRecipeID,
				ValidIngredientPreparationID:     &tossBreadVIP.ID,
				ValidIngredientMeasurementUnitID: &breadCupVIMU.ID,
				IngredientID:                     &heartyBread.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "garlic parmesan croutons",
				QuantityNotes:                    "three-quarters of croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      slStep1ID,
				ValidPreparationVesselID: &tossLargeBowlVPV.ID,
				VesselID:                 &largeBowl.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: slStep1ID,
				Name:                "tossed Caesar salad",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Transfer to serving bowl and sprinkle with remaining cheese and croutons
	slStep2ID := identifiers.New()
	slStep2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              slStep2ID,
		BelongsToRecipe: saladRecipeID,
		PreparationID:   transferPrep.ID,
		Index:           2,
		Notes:           "Transfer to a salad bowl and sprinkle with remaining cheese and croutons. Serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             slStep2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &romaineLettuce.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "tossed Caesar salad",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      slStep2ID,
				ValidPreparationVesselID: &transferServingBowlVPV.ID,
				VesselID:                 &servingBowl.ID,
				Name:                     "salad bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: slStep2ID,
				Name:                "Caesar salad in serving bowl",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Sprinkle with remaining cheese and croutons
	slStep3ID := identifiers.New()
	slStep3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              slStep3ID,
		BelongsToRecipe: saladRecipeID,
		PreparationID:   sprinklePrep.ID,
		Index:           3,
		Notes:           "Sprinkle with remaining cheese and croutons. Serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             slStep3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &romaineLettuce.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "Caesar salad in serving bowl",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              slStep3ID,
				ValidIngredientPreparationID:     &sprinkleParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanCupVIMU.ID,
				IngredientID:                     &parmesanCheese.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "finely grated parmesan cheese",
				QuantityNotes:                    "remaining cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              slStep3ID,
				RecipeStepProductRecipeID:        &croutonsRecipeID,
				ValidIngredientPreparationID:     &sprinkleBreadVIP.ID,
				ValidIngredientMeasurementUnitID: &breadCupVIMU.ID,
				IngredientID:                     &heartyBread.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "garlic parmesan croutons",
				QuantityNotes:                    "remaining croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      slStep3ID,
				ValidPreparationVesselID: &sprinkleServingBowlVPV.ID,
				VesselID:                 &servingBowl.ID,
				Name:                     "salad bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: slStep3ID,
				Name:                "Caesar salad",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	saladRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  saladRecipeID,
		CreatedByUser:       userID,
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
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			slStep0, slStep1, slStep2, slStep3,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		croutonsRecipe,
		dressingRecipe,
		saladRecipe,
	}
}
