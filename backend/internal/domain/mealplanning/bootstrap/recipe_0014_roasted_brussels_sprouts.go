package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// RoastedBrusselsSproutsRecipe creates the Roasted Brussels Sprouts recipe.
// Source: https://www.seriouseats.com/roasted-brussels-sprouts-and-shallots-with-balsamic-vinegar-thanksgiving-recipe
func RoastedBrusselsSproutsRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
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
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: trimPrep.ID,
		Index:         0,
		ExplicitInstructions: "Trim the bottoms, remove the outer leaves, and cut the Brussels sprouts in half. This step is optional and can be done up to 2 days ahead.",
		Optional:      true,
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &trimBrusselsSproutsVIP.ID,
				ValidIngredientMeasurementUnitID: &brusselsSproutsPoundVIMU.ID,
				Name:                             "Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &trimChefsKnifeVPI.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "trimmed Brussels sprouts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 1: Cut trimmed Brussels sprouts in half
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: halvePrep.ID,
		Index:         1,
		ExplicitInstructions: "Cut the trimmed Brussels sprouts in half.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &halveBrusselsSproutsVIP.ID,
				Name:                            "trimmed Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &halveChefsKnifeVPI.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "halved Brussels sprouts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 2: Preheat baking sheets in oven to 500°F (260°C)
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: preheatPrep.ID,
		Index:         2,
		ExplicitInstructions: "Preheat the baking sheets in the oven to 500°F (260°C).",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260), // 500°F = 260°C
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &preheatOvenVPV.ID,
				Name:                     "oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationVesselID: &preheatBakingSheetVPV.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "preheated oven with baking sheets at 500°F",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: In a large bowl, add sprouts, 3 tablespoons of the olive oil, and salt and pepper to taste and toss to combine
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: tossPrep.ID,
		Index:         3,
		ExplicitInstructions: "In a large bowl, add the sprouts, 3 tablespoons of the olive oil, and salt and pepper to taste and toss to combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossBrusselsSproutsVIP.ID,
				Name:                            "halved Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ValidIngredientPreparationID:     &tossOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				Name:                             "extra-virgin olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ValidIngredientPreparationID:     &tossSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &tossBlackPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &blackPepperTeaspoonVIMU.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &tossBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
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
				Name:              "tossed Brussels sprouts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 4: Remove the preheated baking sheets from the oven. Place Brussels sprouts on the sheets in a single even layer. Return the sheets to the oven.
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       placePrep.ID,
		Index:                4,
		ExplicitInstructions: "Remove the preheated baking sheets from the oven. Place the Brussels sprouts on the sheets in a single even layer, shaking the sheets to distribute evenly. Return the sheets to the oven.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &placeBrusselsSproutsVIP.ID,
				Name:                            "tossed Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &placeOvenMittVPI.ID,
				Name:                         "oven mitt",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationInstrumentID: &placeTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &placeBakingSheetVPV.ID,
				Name:                            "preheated baking sheets at 500°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &returnOvenVPV.ID,
				Name:                            "preheated oven at 500°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "Brussels sprouts on baking sheets",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				Name:  "baking sheets with Brussels sprouts in oven",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Roast for 10 minutes
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: roastPrep.ID,
		Index:         5,
		ExplicitInstructions: "Roast for 10 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &roastBrusselsSproutsVIP.ID,
				Name:                            "Brussels sprouts on baking sheets",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &roastBakingSheetVPV.ID,
				Name:                            "baking sheet with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &roastOvenVPV.ID,
				Name:                            "preheated oven with baking sheets at 500°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "partially roasted Brussels sprouts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 6: Slice shallots thinly
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       slicePrep.ID,
		Index:                6,
		ExplicitInstructions: "Slice 8 medium shallots thinly.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &sliceShallotsVIP.ID,
				ValidIngredientMeasurementUnitID: &shallotsUnitVIMU.ID,
				Name:                             "medium shallots",
				QuantityNotes:                    "8 medium shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &sliceKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &sliceCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "sliced shallots",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](8),
				},
			},
		},
	}

	// Step 7: While Brussels sprouts roast, in the now-empty bowl, toss shallots, remaining 1 tablespoon olive oil, and salt and pepper to taste to combine
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: tossPrep.ID,
		Index:         7,
		ExplicitInstructions: "While the Brussels sprouts roast, in the now-empty bowl, toss the shallots, remaining 1 tablespoon olive oil, and salt and pepper to taste to combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossShallotsVIP.ID,
				Name:                            "sliced shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ValidIngredientPreparationID:     &tossOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				Name:                             "extra-virgin olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &tossSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ValidIngredientPreparationID:     &tossBlackPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &blackPepperTeaspoonVIMU.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &tossBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "now-empty large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "tossed shallots",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](8),
				},
			},
		},
	}

	// Step 8: Working quickly and carefully, add the shallot mixture to the sheets and stir with the Brussels sprouts to combine.
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       stirPrep.ID,
		Index:                8,
		ExplicitInstructions: "Working quickly and carefully, add the shallot mixture to the sheets and stir with the Brussels sprouts to combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirShallotsVIP.ID,
				Name:                            "tossed shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirBrusselsSproutsVIP.ID,
				Name:                            "partially roasted Brussels sprouts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &stirBakingSheetVPV.ID,
				Name:                            "baking sheet with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "Brussels sprouts and shallots combined",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 9: Rotate and swap pans top to bottom in oven. Continue to bake until Brussels sprouts are deeply charred and fully tender and shallots begin to brown, 10 to 15 minutes more.
	brownedState := enums.IngredientStates["browned"]
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: roastPrep.ID,
		Index:         9,
		ExplicitInstructions: "Rotate and swap the pans top to bottom in the oven. Continue to bake until the Brussels sprouts are deeply charred and fully tender and the shallots begin to brown, 10 to 15 minutes more.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
			Max: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &roastBrusselsSproutsVIP.ID,
				Name:                            "Brussels sprouts and shallots combined",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &rotateBakingSheetVPV.ID,
				Name:                            "baking sheet with Brussels sprouts in oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &roastOvenVPV.ID,
				Name:                            "preheated oven with baking sheets at 500°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "roasted Brussels sprouts and shallots",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: brownedState.ID,
				Notes:             "Brussels sprouts should be deeply charred and fully tender, and shallots should begin to brown",
				Ingredients:       []uint64{0}, // Index of combined ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 10: Immediately after removing sheets from oven, drizzle sprouts with balsamic vinegar and shake to coat
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       drizzlePrep.ID,
		Index:                10,
		ExplicitInstructions: "Immediately after removing the sheets from the oven, drizzle the sprouts with balsamic vinegar and shake to coat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &drizzleBrusselsSproutsVIP.ID,
				Name:                            "roasted Brussels sprouts and shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ValidIngredientPreparationID:     &drizzleBalsamicVinegarVIP.ID,
				ValidIngredientMeasurementUnitID: &balsamicVinegarTablespoonVIMU.ID,
				Name:                             "balsamic vinegar",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &drizzleBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &drizzleBakingSheetVPV.ID,
				Name:                            "baking sheet with roasted Brussels sprouts and shallots",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "balsamic-glazed Brussels sprouts and shallots",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 11: Season to taste with more salt and pepper if desired and serve
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         11,
		ExplicitInstructions: "Season to taste with more salt and pepper if desired and serve.",
		Optional:      true,
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonBrusselsSproutsVIP.ID,
				Name:                            "balsamic-glazed Brussels sprouts and shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonBlackPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &blackPepperTeaspoonVIMU.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &seasonBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "roasted Brussels sprouts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Create prep task for trimming and halving Brussels sprouts ahead of time
	prepTask1 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
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
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 0, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 1, SatisfiesRecipeStep: true},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
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
			Steps:             []*mealplanning.RecipeStepCreationRequestInput{step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11},
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{prepTask1},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
