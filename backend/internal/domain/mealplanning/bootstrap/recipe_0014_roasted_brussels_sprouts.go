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
	tossPrep := enums.Preparations["toss"]
	adjustPrep := enums.Preparations["adjust"]
	placePrep := enums.Preparations["place"]
	preheatPrep := enums.Preparations["preheat"]
	removePrep := enums.Preparations["remove"]
	dividePrep := enums.Preparations["divide"]
	shakePrep := enums.Preparations["shake"]
	returnPrep := enums.Preparations["return"]
	roastPrep := enums.Preparations["roast"]
	stirPrep := enums.Preparations["stir"]
	rotatePrep := enums.Preparations["rotate"]
	swapPrep := enums.Preparations["swap"]
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
	dishTowel := enums.Instruments["dish towel"]

	// Get vessels
	bakingSheet := enums.Vessels["baking sheet"]
	oven := enums.Vessels["oven"]
	largeBowl := enums.Vessels["large bowl"]

	// Get bridge table entries
	// Trim
	trimBrusselsSproutsVIP := enums.IngredientPreparations[trimPrep.ID][brusselsSprouts.ID]
	trimChefsKnifeVPI := enums.PreparationInstruments[trimPrep.ID][chefsKnife.ID]

	// Halve
	halveBrusselsSproutsVIP := enums.IngredientPreparations[halvePrep.ID][brusselsSprouts.ID]
	halveChefsKnifeVPI := enums.PreparationInstruments[halvePrep.ID][chefsKnife.ID]

	// Toss
	tossBrusselsSproutsVIP := enums.IngredientPreparations[tossPrep.ID][brusselsSprouts.ID]
	tossOliveOilVIP := enums.IngredientPreparations[tossPrep.ID][oliveOil.ID]
	tossSaltVIP := enums.IngredientPreparations[tossPrep.ID][salt.ID]
	tossBlackPepperVIP := enums.IngredientPreparations[tossPrep.ID][blackPepper.ID]
	tossShallotsVIP := enums.IngredientPreparations[tossPrep.ID][shallots.ID]
	tossLargeBowlVPV := enums.PreparationVessels[tossPrep.ID][largeBowl.ID]
	tossBareHandsVPI := enums.PreparationInstruments[tossPrep.ID][bareHands.ID]

	// Adjust
	adjustOvenVPV := enums.PreparationVessels[adjustPrep.ID][oven.ID]

	// Place
	placeBakingSheetVPV := enums.PreparationVessels[placePrep.ID][bakingSheet.ID]
	placeOvenVPV := enums.PreparationVessels[placePrep.ID][oven.ID]

	// Preheat
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]
	preheatBakingSheetVPV := enums.PreparationVessels[preheatPrep.ID][bakingSheet.ID]

	// Remove
	removeBakingSheetVPV := enums.PreparationVessels[removePrep.ID][bakingSheet.ID]
	removeOvenMittVPI := enums.PreparationInstruments[removePrep.ID][ovenMitt.ID]
	removeDishTowelVPI := enums.PreparationInstruments[removePrep.ID][dishTowel.ID]

	// Divide
	divideBrusselsSproutsVIP := enums.IngredientPreparations[dividePrep.ID][brusselsSprouts.ID]
	divideShallotsVIP := enums.IngredientPreparations[dividePrep.ID][shallots.ID]
	divideBakingSheetVPV := enums.PreparationVessels[dividePrep.ID][bakingSheet.ID]

	// Shake
	shakeBareHandsVPI := enums.PreparationInstruments[shakePrep.ID][bareHands.ID]

	// Return
	returnOvenVPV := enums.PreparationVessels[returnPrep.ID][oven.ID]

	// Roast
	roastBrusselsSproutsVIP := enums.IngredientPreparations[roastPrep.ID][brusselsSprouts.ID]
	roastShallotsVIP := enums.IngredientPreparations[roastPrep.ID][shallots.ID]
	roastBakingSheetVPV := enums.PreparationVessels[roastPrep.ID][bakingSheet.ID]
	roastOvenVPV := enums.PreparationVessels[roastPrep.ID][oven.ID]

	// Stir
	stirBrusselsSproutsVIP := enums.IngredientPreparations[stirPrep.ID][brusselsSprouts.ID]
	stirBakingSheetVPV := enums.PreparationVessels[stirPrep.ID][bakingSheet.ID]

	// Rotate
	rotateBakingSheetVPV := enums.PreparationVessels[rotatePrep.ID][bakingSheet.ID]

	// Swap
	swapBakingSheetVPV := enums.PreparationVessels[swapPrep.ID][bakingSheet.ID]

	// Drizzle
	drizzleBalsamicVinegarVIP := enums.IngredientPreparations[drizzlePrep.ID][balsamicVinegar.ID]
	drizzleBrusselsSproutsVIP := enums.IngredientPreparations[drizzlePrep.ID][brusselsSprouts.ID]
	drizzleBakingSheetVPV := enums.PreparationVessels[drizzlePrep.ID][bakingSheet.ID]

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
				Quantity: types.OptionalFloat32Range{
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
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 2: Adjust oven racks to upper and lower middle positions
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   adjustPrep.ID,
		Index:           2,
		Notes:           "Adjust oven racks to upper and lower middle positions.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &adjustOvenVPV.ID,
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
				BelongsToRecipeStep: step2ID,
				Name:                "oven with adjusted racks",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Place a rimmed baking sheet on each rack
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   placePrep.ID,
		Index:           3,
		Notes:           "Place a rimmed baking sheet on each rack.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &placeBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &placeOvenVPV.ID,
				VesselID:                        &oven.ID,
				Name:                            "oven with adjusted racks",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "baking sheets in oven",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Preheat oven with sheets to 500°F (260°C)
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   preheatPrep.ID,
		Index:           4,
		Notes:           "Preheat oven with sheets to 500°F (260°C).",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260), // 500°F = 260°C
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &preheatOvenVPV.ID,
				VesselID:                        &oven.ID,
				Name:                            "baking sheets in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &preheatBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheets in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "preheated oven with baking sheets at 500°F",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: In a large bowl, add sprouts, 3 tablespoons of the olive oil, and salt and pepper to taste and toss to combine
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   tossPrep.ID,
		Index:           5,
		Notes:           "In a large bowl, add sprouts, 3 tablespoons of the olive oil, and salt and pepper to taste and toss to combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
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
				BelongsToRecipeStep:              step5ID,
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
				BelongsToRecipeStep:              step5ID,
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
				BelongsToRecipeStep:              step5ID,
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
				BelongsToRecipeStep:          step5ID,
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
				BelongsToRecipeStep:      step5ID,
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
				BelongsToRecipeStep: step5ID,
				Name:                "tossed Brussels sprouts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 6: Once the oven has reached temperature, working quickly, remove the baking sheets with a dish towel or oven mitt and divide Brussels sprouts mixture evenly between both trays, shaking sheets to distribute into a single even layer. Return pans to the oven.
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dividePrep.ID,
		Index:           6,
		Notes:           "Once the oven has reached temperature, working quickly, remove the baking sheets with a dish towel or oven mitt and divide Brussels sprouts mixture evenly between both trays, shaking sheets to distribute into a single even layer. Return pans to the oven.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &divideBrusselsSproutsVIP.ID,
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
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: &removeOvenMittVPI.ID,
				InstrumentID:                 &ovenMitt.ID,
				Name:                         "oven mitt",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: &removeDishTowelVPI.ID,
				InstrumentID:                 &dishTowel.ID,
				Name:                         "dish towel",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: &shakeBareHandsVPI.ID,
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
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &removeBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "preheated oven with baking sheets at 500°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &divideBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &returnOvenVPV.ID,
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
				BelongsToRecipeStep: step6ID,
				Name:                "Brussels sprouts on baking sheets",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "baking sheets with Brussels sprouts in oven",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Roast for 10 minutes
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   roastPrep.ID,
		Index:           7,
		Notes:           "Roast for 10 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
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
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &roastBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheets with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
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
				BelongsToRecipeStep: step7ID,
				Name:                "partially roasted Brussels sprouts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 8: While Brussels sprouts roast, in the now-empty bowl, toss shallots, remaining 1 tablespoon olive oil, and salt and pepper to taste to combine
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   tossPrep.ID,
		Index:           8,
		Notes:           "While Brussels sprouts roast, in the now-empty bowl, toss shallots, remaining 1 tablespoon olive oil, and salt and pepper to taste to combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &tossShallotsVIP.ID,
				ValidIngredientMeasurementUnitID: &shallotsUnitVIMU.ID,
				IngredientID:                     &shallots.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "medium shallots, sliced thinly",
				QuantityNotes:                    "8 medium shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
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
				BelongsToRecipeStep:              step8ID,
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
				BelongsToRecipeStep:              step8ID,
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
				BelongsToRecipeStep:          step8ID,
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
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
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
				BelongsToRecipeStep: step8ID,
				Name:                "tossed shallots",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](8),
				},
			},
		},
	}

	// Step 9: Working quickly and carefully, divide the shallot mixture evenly between the two sheets and stir with the Brussels sprouts to combine
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           9,
		Notes:           "Working quickly and carefully, divide the shallot mixture evenly between the two sheets and stir with the Brussels sprouts to combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &divideShallotsVIP.ID,
				IngredientID:                    &shallots.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "tossed shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
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
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &stirBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheets with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "Brussels sprouts and shallots combined",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 10: Rotate and swap pans top to bottom in oven. Continue to bake until Brussels sprouts are deeply charred and fully tender and shallots begin to brown, 10 to 15 minutes more.
	step10ID := identifiers.New()
	step10BrusselsSproutsIngredientID := identifiers.New()
	step10ShallotsIngredientID := identifiers.New()
	step10CompletionConditionID := identifiers.New()
	brownedState := enums.IngredientStates["browned"]
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   roastPrep.ID,
		Index:           10,
		Notes:           "Rotate and swap pans top to bottom in oven. Continue to bake until Brussels sprouts are deeply charred and fully tender and shallots begin to brown, 10 to 15 minutes more.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
			Max: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step10BrusselsSproutsIngredientID,
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &roastBrusselsSproutsVIP.ID,
				IngredientID:                    &brusselsSprouts.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "Brussels sprouts and shallots combined",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                              step10ShallotsIngredientID,
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &roastShallotsVIP.ID,
				IngredientID:                    &shallots.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "Brussels sprouts and shallots combined",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &rotateBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheets with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &swapBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheets with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
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
				BelongsToRecipeStep: step10ID,
				Name:                "roasted Brussels sprouts and shallots",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step10CompletionConditionID,
				BelongsToRecipeStep: step10ID,
				IngredientStateID:   brownedState.ID,
				Notes:               "Brussels sprouts should be deeply charred and fully tender, and shallots should begin to brown",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step10CompletionConditionID,
						RecipeStepIngredient:                   step10BrusselsSproutsIngredientID,
					},
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step10CompletionConditionID,
						RecipeStepIngredient:                   step10ShallotsIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 11: Immediately after removing sheets from oven, drizzle sprouts with balsamic vinegar and shake to coat
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   drizzlePrep.ID,
		Index:           11,
		Notes:           "Immediately after removing sheets from oven, drizzle sprouts with balsamic vinegar and shake to coat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
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
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &drizzleBalsamicVinegarVIP.ID,
				ValidIngredientMeasurementUnitID: &balsamicVinegarTablespoonVIMU.ID,
				IngredientID:                     &balsamicVinegar.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "balsamic vinegar or aged sherry vinegar",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step11ID,
				ValidPreparationInstrumentID: &shakeBareHandsVPI.ID,
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
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &drizzleBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "baking sheets with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "balsamic-glazed Brussels sprouts and shallots",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 12: Season to taste with more salt and pepper if desired and serve
	step12ID := identifiers.New()
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           12,
		Notes:           "Season to taste with more salt and pepper if desired and serve.",
		Optional:        true,
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
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
				BelongsToRecipeStep:              step12ID,
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
				BelongsToRecipeStep:              step12ID,
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
				BelongsToRecipeStep:          step12ID,
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
				BelongsToRecipeStep: step12ID,
				Name:                "roasted Brussels sprouts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
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
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskDatabaseCreationInput{
			prepTask1,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		roastedBrusselsSproutsRecipe,
	}
}
