package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// RoastedBrusselsSproutsRecipe creates the Roasted Brussels Sprouts recipe.
// Source: https://www.seriouseats.com/roasted-brussels-sprouts-and-shallots-with-balsamic-vinegar-thanksgiving-recipe
func RoastedBrusselsSproutsRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	trimPrep := enums.Preparations["trim"]
	halvePrep := enums.Preparations["halve"]
	slicePrep := enums.Preparations["slice"]
	tossPrep := enums.Preparations["toss"]
	preheatPrep := enums.Preparations["preheat"]
	placePrep := enums.Preparations["place"]
	returnPrep := enums.Preparations["return"]
	roastPrep := enums.Preparations["roast"]
	stirPrep := enums.Preparations["stir"]
	rotatePrep := enums.Preparations["rotate"]
	drizzlePrep := enums.Preparations["drizzle"]
	seasonPrep := enums.Preparations["season"]

	// Get ingredients
	brusselsSprouts := enums.Ingredients["Brussels sprouts"]
	oliveOil := enums.Ingredients["olive oil"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	shallots := enums.Ingredients["shallot"]
	balsamicVinegar := enums.Ingredients["balsamic vinegar"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	chefsKnife := enums.Instruments["knife"]
	bareHands := enums.Instruments["bare hands"]
	ovenMitt := enums.Instruments["oven mitt"]
	tongs := enums.Instruments["tongs"]

	// Get vessels
	bakingSheet := enums.Vessels["baking sheet"]
	oven := enums.Vessels["oven"]
	largeBowl := enums.Vessels["large bowl"]
	cuttingBoard := enums.Vessels["cutting board"]

	// Get bridge table entries
	// Trim
	trimBrusselsSproutsVIP := enums.IngredientPreparations[trimPrep.ID][brusselsSprouts.ID]
	trimChefsKnifeVPI := enums.PreparationInstruments[trimPrep.ID][chefsKnife.ID]

	// Halve
	halveBrusselsSproutsVIP := enums.IngredientPreparations[halvePrep.ID][brusselsSprouts.ID]
	halveChefsKnifeVPI := enums.PreparationInstruments[halvePrep.ID][chefsKnife.ID]

	// Slice
	sliceShallotsVIP := enums.IngredientPreparations[slicePrep.ID][shallots.ID]
	sliceKnifeVPI := enums.PreparationInstruments[slicePrep.ID][chefsKnife.ID]
	sliceCuttingBoardVPV := enums.PreparationVessels[slicePrep.ID][cuttingBoard.ID]

	// Toss
	tossBrusselsSproutsVIP := enums.IngredientPreparations[tossPrep.ID][brusselsSprouts.ID]
	tossOliveOilVIP := enums.IngredientPreparations[tossPrep.ID][oliveOil.ID]
	tossSaltVIP := enums.IngredientPreparations[tossPrep.ID][salt.ID]
	tossBlackPepperVIP := enums.IngredientPreparations[tossPrep.ID][blackPepper.ID]
	tossShallotsVIP := enums.IngredientPreparations[tossPrep.ID][shallots.ID]
	tossLargeBowlVPV := enums.PreparationVessels[tossPrep.ID][largeBowl.ID]
	tossBareHandsVPI := enums.PreparationInstruments[tossPrep.ID][bareHands.ID]

	// Preheat
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]
	preheatBakingSheetVPV := enums.PreparationVessels[preheatPrep.ID][bakingSheet.ID]

	// Place
	placeBrusselsSproutsVIP := enums.IngredientPreparations[placePrep.ID][brusselsSprouts.ID]
	placeBakingSheetVPV := enums.PreparationVessels[placePrep.ID][bakingSheet.ID]
	placeTongsVPI := enums.PreparationInstruments[placePrep.ID][tongs.ID]
	placeOvenMittVPI := enums.PreparationInstruments[placePrep.ID][ovenMitt.ID]

	// Return
	returnOvenVPV := enums.PreparationVessels[returnPrep.ID][oven.ID]

	// Roast
	roastBrusselsSproutsVIP := enums.IngredientPreparations[roastPrep.ID][brusselsSprouts.ID]
	roastBakingSheetVPV := enums.PreparationVessels[roastPrep.ID][bakingSheet.ID]
	roastOvenVPV := enums.PreparationVessels[roastPrep.ID][oven.ID]

	// Stir
	stirBrusselsSproutsVIP := enums.IngredientPreparations[stirPrep.ID][brusselsSprouts.ID]
	stirShallotsVIP := enums.IngredientPreparations[stirPrep.ID][shallots.ID]
	stirBakingSheetVPV := enums.PreparationVessels[stirPrep.ID][bakingSheet.ID]

	// Rotate
	rotateBakingSheetVPV := enums.PreparationVessels[rotatePrep.ID][bakingSheet.ID]

	// Drizzle
	drizzleBalsamicVinegarVIP := enums.IngredientPreparations[drizzlePrep.ID][balsamicVinegar.ID]
	drizzleBrusselsSproutsVIP := enums.IngredientPreparations[drizzlePrep.ID][brusselsSprouts.ID]
	drizzleBakingSheetVPV := enums.PreparationVessels[drizzlePrep.ID][bakingSheet.ID]
	drizzleBareHandsVPI := enums.PreparationInstruments[drizzlePrep.ID][bareHands.ID]

	// Season
	seasonBrusselsSproutsVIP := enums.IngredientPreparations[seasonPrep.ID][brusselsSprouts.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonBlackPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	seasonBareHandsVPI := enums.PreparationInstruments[seasonPrep.ID][bareHands.ID]

	// Measurement unit bridges
	brusselsSproutsPoundVIMU := enums.IngredientMeasurementUnits[brusselsSprouts.ID][poundMeasurement.ID]
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	blackPepperTeaspoonVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]
	shallotsUnitVIMU := enums.IngredientMeasurementUnits[shallots.ID][unitMeasurement.ID]
	balsamicVinegarTablespoonVIMU := enums.IngredientMeasurementUnits[balsamicVinegar.ID][tablespoonMeasurement.ID]

	// Step 0: Trim bottoms, remove outer leaves, and cut Brussels sprouts in half (optional prep step)
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   trimPrep.ID,
		Index:           0,
		Notes:           "Trim bottoms, remove outer leaves, and cut Brussels sprouts in half. This step is optional and can be done up to 2 days ahead.",
		Optional:        true,
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimBrusselsSproutsVIP.ID,
				ValidIngredientMeasurementUnitID: &brusselsSproutsPoundVIMU.ID,
				IngredientID:                     &brusselsSprouts.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
				ValidPreparationInstrumentID: &trimChefsKnifeVPI.ID,
				InstrumentID:                 &chefsKnife.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "trimmed Brussels sprouts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 1: Cut trimmed Brussels sprouts in half
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   halvePrep.ID,
		Index:           1,
		Notes:           "Cut trimmed Brussels sprouts in half.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &halveBrusselsSproutsVIP.ID,
				IngredientID:                    &brusselsSprouts.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "trimmed Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidPreparationInstrumentID: &halveChefsKnifeVPI.ID,
				InstrumentID:                 &chefsKnife.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "halved Brussels sprouts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 2: Preheat baking sheets in oven to 500°F (260°C)
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   preheatPrep.ID,
		Index:           2,
		Notes:           "Preheat baking sheets in oven to 500°F (260°C).",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260), // 500°F = 260°C
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &preheatOvenVPV.ID,
				VesselID:                 &oven.ID,
				Name:                     "oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &preheatBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "preheated oven with baking sheets at 500°F",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: In a large bowl, add sprouts, 3 tablespoons of the olive oil, and salt and pepper to taste and toss to combine
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   tossPrep.ID,
		Index:           3,
		Notes:           "In a large bowl, add sprouts, 3 tablespoons of the olive oil, and salt and pepper to taste and toss to combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossBrusselsSproutsVIP.ID,
				IngredientID:                    &brusselsSprouts.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "halved Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &tossOliveOilVIP.ID,
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
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &tossSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &tossBlackPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &blackPepperTeaspoonVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step3ID,
				ValidPreparationInstrumentID: &tossBareHandsVPI.ID,
				InstrumentID:                 &bareHands.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
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
				BelongsToRecipeStep: step3ID,
				Name:                "tossed Brussels sprouts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 4: Remove the preheated baking sheets from the oven. Place Brussels sprouts on the sheets in a single even layer. Return the sheets to the oven.
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   placePrep.ID,
		Index:           4,
		Notes:           "Remove the preheated baking sheets from the oven. Place Brussels sprouts on the sheets in a single even layer, shaking sheets to distribute evenly. Return the sheets to the oven.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &placeBrusselsSproutsVIP.ID,
				IngredientID:                    &brusselsSprouts.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "tossed Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &placeOvenMittVPI.ID,
				InstrumentID:                 &ovenMitt.ID,
				Name:                         "oven mitt",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &placeTongsVPI.ID,
				InstrumentID:                 &tongs.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &placeBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "preheated baking sheets at 500°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &returnOvenVPV.ID,
				VesselID:                        &oven.ID,
				Name:                            "preheated oven at 500°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "Brussels sprouts on baking sheets",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "baking sheets with Brussels sprouts in oven",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Roast for 10 minutes
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   roastPrep.ID,
		Index:           5,
		Notes:           "Roast for 10 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &roastBrusselsSproutsVIP.ID,
				IngredientID:                    &brusselsSprouts.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "Brussels sprouts on baking sheets",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &roastBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheet with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &roastOvenVPV.ID,
				VesselID:                        &oven.ID,
				Name:                            "preheated oven with baking sheets at 500°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "partially roasted Brussels sprouts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 6: Slice shallots thinly
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   slicePrep.ID,
		Index:           6,
		Notes:           "Slice 8 medium shallots thinly.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &sliceShallotsVIP.ID,
				ValidIngredientMeasurementUnitID: &shallotsUnitVIMU.ID,
				IngredientID:                     &shallots.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "medium shallots",
				QuantityNotes:                    "8 medium shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: &sliceKnifeVPI.ID,
				InstrumentID:                 &chefsKnife.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &sliceCuttingBoardVPV.ID,
				VesselID:                 &cuttingBoard.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "sliced shallots",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](8),
				},
			},
		},
	}

	// Step 7: While Brussels sprouts roast, in the now-empty bowl, toss shallots, remaining 1 tablespoon olive oil, and salt and pepper to taste to combine
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   tossPrep.ID,
		Index:           7,
		Notes:           "While Brussels sprouts roast, in the now-empty bowl, toss shallots, remaining 1 tablespoon olive oil, and salt and pepper to taste to combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossShallotsVIP.ID,
				IngredientID:                    &shallots.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "sliced shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ValidIngredientPreparationID:     &tossOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				IngredientID:                     &oliveOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "extra-virgin olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ValidIngredientPreparationID:     &tossSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ValidIngredientPreparationID:     &tossBlackPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &blackPepperTeaspoonVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step7ID,
				ValidPreparationInstrumentID: &tossBareHandsVPI.ID,
				InstrumentID:                 &bareHands.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				VesselID:                        &largeBowl.ID,
				Name:                            "now-empty large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "tossed shallots",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](8),
				},
			},
		},
	}

	// Step 8: Working quickly and carefully, add the shallot mixture to the sheets and stir with the Brussels sprouts to combine.
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           8,
		Notes:           "Working quickly and carefully, add the shallot mixture to the sheets and stir with the Brussels sprouts to combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirShallotsVIP.ID,
				IngredientID:                    &shallots.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "tossed shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirBrusselsSproutsVIP.ID,
				IngredientID:                    &brusselsSprouts.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "partially roasted Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &stirBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheet with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "Brussels sprouts and shallots combined",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 9: Rotate and swap pans top to bottom in oven. Continue to bake until Brussels sprouts are deeply charred and fully tender and shallots begin to brown, 10 to 15 minutes more.
	step9ID := identifiers.New()
	step9CombinedIngredientID := identifiers.New()
	step9CompletionConditionID := identifiers.New()
	brownedState := enums.IngredientStates["browned"]
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   roastPrep.ID,
		Index:           9,
		Notes:           "Rotate and swap pans top to bottom in oven. Continue to bake until Brussels sprouts are deeply charred and fully tender and shallots begin to brown, 10 to 15 minutes more.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
			Max: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step9CombinedIngredientID,
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &roastBrusselsSproutsVIP.ID,
				IngredientID:                    &brusselsSprouts.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "Brussels sprouts and shallots combined",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &rotateBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheet with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &roastOvenVPV.ID,
				VesselID:                        &oven.ID,
				Name:                            "preheated oven with baking sheets at 500°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "roasted Brussels sprouts and shallots",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step9CompletionConditionID,
				BelongsToRecipeStep: step9ID,
				IngredientStateID:   brownedState.ID,
				Notes:               "Brussels sprouts should be deeply charred and fully tender, and shallots should begin to brown",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step9CompletionConditionID,
						RecipeStepIngredient:                   step9CombinedIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 10: Immediately after removing sheets from oven, drizzle sprouts with balsamic vinegar and shake to coat
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   drizzlePrep.ID,
		Index:           10,
		Notes:           "Immediately after removing sheets from oven, drizzle sprouts with balsamic vinegar and shake to coat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &drizzleBrusselsSproutsVIP.ID,
				IngredientID:                    &brusselsSprouts.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "roasted Brussels sprouts and shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &drizzleBalsamicVinegarVIP.ID,
				ValidIngredientMeasurementUnitID: &balsamicVinegarTablespoonVIMU.ID,
				IngredientID:                     &balsamicVinegar.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "balsamic vinegar",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step10ID,
				ValidPreparationInstrumentID: &drizzleBareHandsVPI.ID,
				InstrumentID:                 &bareHands.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &drizzleBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheet with roasted Brussels sprouts and shallots",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "balsamic-glazed Brussels sprouts and shallots",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 11: Season to taste with more salt and pepper if desired and serve
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           11,
		Notes:           "Season to taste with more salt and pepper if desired and serve.",
		Optional:        true,
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonBrusselsSproutsVIP.ID,
				IngredientID:                    &brusselsSprouts.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "balsamic-glazed Brussels sprouts and shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &seasonBlackPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &blackPepperTeaspoonVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step11ID,
				ValidPreparationInstrumentID: &seasonBareHandsVPI.ID,
				InstrumentID:                 &bareHands.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "roasted Brussels sprouts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Create prep task for trimming and halving Brussels sprouts ahead of time
	prepTask1ID := identifiers.New()
	prepTask1 := &mealplanning.RecipePrepTaskDatabaseCreationInput{
		ID:                          prepTask1ID,
		BelongsToRecipe:             recipeID,
		Name:                        "Trim and halve Brussels sprouts",
		Description:                 "The Brussels sprouts and shallots can be cut and refrigerated, before tossing with oil, for up to 2 days.",
		Notes:                       "Preparing the vegetables ahead saves time on the day of cooking.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the trimmed and halved Brussels sprouts in the refrigerator for up to 2 days.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeAirtightContainer,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 0,
			Max: pointer.To[uint32](172800), // 2 days
		},
		TaskSteps: []*mealplanning.RecipePrepTaskStepDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: step0ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step1ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: true},
		},
	}

	roastedBrusselsSproutsRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Roasted Brussels Sprouts",
		Slug:                "roasted-brussels-sprouts",
		Source:              "https://www.seriouseats.com/roasted-brussels-sprouts-and-shallots-with-balsamic-vinegar-thanksgiving-recipe",
		Description:         "A last-minute drizzle of balsamic vinegar adds a tart glaze to these crispy sprouts. Extremely high heat, plus a preheated roasting pan, gives the Brussels sprouts sweet flavor and a nutty char.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 8,
			Max: pointer.To[float32](12),
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskDatabaseCreationInput{
			prepTask1,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		roastedBrusselsSproutsRecipe,
	}
}
