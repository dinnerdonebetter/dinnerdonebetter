package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func PanSearedButterBastedSteakRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {

	// Get preparations
	dryPrep := enums.Preparations["dry"]
	seasonPrep := enums.Preparations["season"]
	slicePrep := enums.Preparations["slice"]
	restPrep := enums.Preparations["rest"]
	heatPrep := enums.Preparations["heat"]
	panSearPrep := enums.Preparations["pan-sear"]
	bastePrep := enums.Preparations["baste"]

	// Get ingredients
	ribeye := enums.Ingredients["ribeye steak"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	vegetableOil := enums.Ingredients["vegetable oil"]
	butter := enums.Ingredients["butter"]
	thyme := enums.Ingredients["thyme"]
	rosemary := enums.Ingredients["rosemary"]
	shallot := enums.Ingredients["shallot"]
	paperTowels := enums.Ingredients["paper towels"]

	// Get measurement units
	unitMeasurement := enums.MeasurementUnits["unit"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	milliliterMeasurement := enums.MeasurementUnits["milliliter"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]

	// Get instruments
	bareHands := enums.Instruments["bare hands"]
	knife := enums.Instruments["knife"]
	tongs := enums.Instruments["tongs"]
	spoon := enums.Instruments["spoon"]
	thermometer := enums.Instruments["instant-read thermometer"]

	// Get vessels
	sheetPan := enums.Vessels["sheet pan"]
	cuttingBoard := enums.Vessels["cutting board"]
	castIronSkillet := enums.Vessels["cast iron skillet"]
	servingPlate := enums.Vessels["serving plate"]

	// Get ingredient states for completion conditions
	smokingState := enums.IngredientStates["smoking"]
	atTemperatureState := enums.IngredientStates["at temperature"]

	// Get bridge table entries
	// Dry preparation bridges
	dryRibeyeVIP := enums.IngredientPreparations[dryPrep.ID][ribeye.ID]
	ribeyeGramVIMU := enums.IngredientMeasurementUnits[ribeye.ID][gramMeasurement.ID]
	dryPaperTowelsVIP := enums.IngredientPreparations[dryPrep.ID][paperTowels.ID]
	paperTowelsUnitVIMU := enums.IngredientMeasurementUnits[paperTowels.ID][unitMeasurement.ID]
	dryBareHandsVPI := enums.PreparationInstruments[dryPrep.ID][bareHands.ID]

	// Season preparation bridges
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]
	seasonSheetPanVPV := enums.PreparationVessels[seasonPrep.ID][sheetPan.ID]

	// Slice preparation bridges
	sliceShallotVIP := enums.IngredientPreparations[slicePrep.ID][shallot.ID]
	sliceKnifeVPI := enums.PreparationInstruments[slicePrep.ID][knife.ID]
	sliceBareHandsVPI := enums.PreparationInstruments[slicePrep.ID][bareHands.ID]
	sliceCuttingBoardVPV := enums.PreparationVessels[slicePrep.ID][cuttingBoard.ID]

	// Rest preparation bridges (for optional rest step)
	restSheetPanVPV := enums.PreparationVessels[restPrep.ID][sheetPan.ID]

	// Heat preparation bridges
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][vegetableOil.ID]
	oilMilliliterVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][milliliterMeasurement.ID]
	heatSkilletVPV := enums.PreparationVessels[heatPrep.ID][castIronSkillet.ID]

	// Pan-sear preparation bridges
	panSearTongsVPI := enums.PreparationInstruments[panSearPrep.ID][tongs.ID]
	panSearSkilletVPV := enums.PreparationVessels[panSearPrep.ID][castIronSkillet.ID]

	// Baste preparation bridges
	basteButterVIP := enums.IngredientPreparations[bastePrep.ID][butter.ID]
	basteThymeVIP := enums.IngredientPreparations[bastePrep.ID][thyme.ID]
	basteRosemaryVIP := enums.IngredientPreparations[bastePrep.ID][rosemary.ID]
	butterGramVIMU := enums.IngredientMeasurementUnits[butter.ID][gramMeasurement.ID]
	thymeSprigVIMU := enums.IngredientMeasurementUnits[thyme.ID][sprigMeasurement.ID]
	rosemarySprigVIMU := enums.IngredientMeasurementUnits[rosemary.ID][sprigMeasurement.ID]
	shallotGramVIMU := enums.IngredientMeasurementUnits[shallot.ID][gramMeasurement.ID]
	basteSpoonVPI := enums.PreparationInstruments[bastePrep.ID][spoon.ID]
	basteThermometerVPI := enums.PreparationInstruments[bastePrep.ID][thermometer.ID]
	basteTongsVPI := enums.PreparationInstruments[bastePrep.ID][tongs.ID]
	basteSkilletVPV := enums.PreparationVessels[bastePrep.ID][castIronSkillet.ID]

	// Final rest preparation bridges
	restTongsVPI := enums.PreparationInstruments[restPrep.ID][tongs.ID]
	restPlateVPV := enums.PreparationVessels[restPrep.ID][servingPlate.ID]

	// Step 0: Pat dry the steak
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: dryPrep.ID,
		Index:         0,
		Notes:         "Carefully pat steak dry with paper towels using your bare hands.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &dryRibeyeVIP.ID,
				ValidIngredientMeasurementUnitID: &ribeyeGramVIMU.ID,
				Name:                             "bone-in ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 700,                      // 24 ounces = ~680g, rounded to 700g
					Max: pointer.To[float32](900), // 32 ounces = ~907g, rounded to 900g
				},
			},
			{
				ValidIngredientPreparationID:     &dryPaperTowelsVIP.ID,
				ValidIngredientMeasurementUnitID: &paperTowelsUnitVIMU.ID,
				Name:                             "paper towels",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &dryBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "dried bone-in ribeye steak",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](700),
					Max: pointer.To[float32](900),
				},
			},
		},
	}

	// Step 1: Season the steak
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         1,
		Notes:         "Season liberally on all sides, including edges, with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "dried bone-in ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 700,
					Max: pointer.To[float32](900),
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltGramVIMU.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste:     true,
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
			},
			{
				ValidIngredientPreparationID:     &seasonPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperGramVIMU.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste:     true,
				Index:       pointer.To[uint16](2),
				OptionIndex: 0,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &seasonSheetPanVPV.ID,
				Name:                     "sheet pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seasoned bone-in ribeye steak",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](700),
					Max: pointer.To[float32](900),
				},
			},
		},
	}

	// Step 2: Rest the steak (optional - at room temperature or refrigerated)
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: restPrep.ID,
		Index:         2,
		Optional:      true,
		Notes:         "If desired, let steak rest at room temperature for 45 minutes, or refrigerated, loosely covered, up to 3 days.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](2700),   // 45 minutes minimum
			Max: pointer.To[uint32](259200), // 3 days maximum
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "seasoned bone-in ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 700,
					Max: pointer.To[float32](900),
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &restSheetPanVPV.ID,
				Name:                     "sheet pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "rested seasoned bone-in ribeye steak",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](700),
					Max: pointer.To[float32](900),
				},
			},
		},
	}

	// Step 3: Slice shallots (optional)
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: slicePrep.ID,
		Index:         3,
		Optional:      true,
		Notes:         "Finely slice shallot into thin slices (about 28g, or 1 large shallot).",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &sliceShallotVIP.ID,
				ValidIngredientMeasurementUnitID: &shallotGramVIMU.ID,
				Name:                             "shallot",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 28, // About 1 large shallot
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
			{
				ValidPreparationInstrumentID: &sliceBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
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
				Name:              "finely sliced shallots",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](28),
				},
			},
		},
	}

	// Step 4: Heat oil until smoking
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: heatPrep.ID,
		Index:         4,
		Notes:         "In a 12-inch heavy-bottomed cast iron skillet, heat oil over high heat until just beginning to smoke.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](200), // High heat, approximately 200°C
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilMilliliterVIMU.ID,
				Name:                             "vegetable or canola oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 60, // 1/4 cup = 60 ml
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &heatSkilletVPV.ID,
				Name:                     "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "heated smoking oil",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &milliliterMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](60),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: smokingState.ID,
				Notes:             "Oil should be just beginning to smoke",
				Ingredients:       []uint64{0}, // Index of the oil ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 5: Pan-sear the steak
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: panSearPrep.ID,
		Index:         5,
		Notes:         "Carefully add steak to the hot skillet and cook, flipping frequently, until a pale golden-brown crust starts to develop, about 4 minutes total.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](240), // 4 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "rested seasoned bone-in ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 700,
					Max: pointer.To[float32](900),
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "heated smoking oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 60,
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &panSearTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &panSearSkilletVPV.ID,
				Name:                     "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "pan-seared bone-in ribeye steak",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](700),
					Max: pointer.To[float32](900),
				},
			},
		},
	}

	// Step 6: Baste the steak
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: bastePrep.ID,
		Index:         6,
		Notes:         "Add butter, herbs (if using), and shallot (if using) to skillet and continue to cook, flipping steak occasionally and basting any light spots with foaming butter. If butter begins to smoke excessively or steak begins to burn, reduce heat to medium. To baste, tilt pan slightly so that butter collects by handle. Use a spoon to pick up butter and pour it over steak, aiming at light spots. Continue flipping and basting until an instant-read thermometer inserted into thickest part of tenderloin side registers 120 to 125°F (49 to 52°C) for medium-rare or 130°F (54°C) for medium, 8 to 10 minutes total.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](480), // 8 minutes
			Max: pointer.To[uint32](600), // 10 minutes
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](49), // 120°F = 49°C
			Max: pointer.To[float32](54), // 130°F = 54°C
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "pan-seared bone-in ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 700,
					Max: pointer.To[float32](900),
				},
			},
			{
				ValidIngredientPreparationID:     &basteButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterGramVIMU.ID,
				Name:                             "unsalted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 45, // 3 tablespoons = 45g
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
			},
			{
				ValidIngredientPreparationID:     &basteThymeVIP.ID,
				ValidIngredientMeasurementUnitID: &thymeSprigVIMU.ID,
				Name:                             "thyme",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6,
				},
				Optional:    true,
				Index:       pointer.To[uint16](2),
				OptionIndex: 0,
			},
			{
				ValidIngredientPreparationID:     &basteRosemaryVIP.ID,
				ValidIngredientMeasurementUnitID: &rosemarySprigVIMU.ID,
				Name:                             "rosemary",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6,
				},
				Optional:    true,
				Index:       pointer.To[uint16](3),
				OptionIndex: 0,
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "finely sliced shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 28,
				},
				Optional:    true,
				Index:       pointer.To[uint16](4),
				OptionIndex: 0,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &basteSpoonVPI.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationInstrumentID: &basteThermometerVPI.ID,
				Name:                         "instant-read thermometer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
			},
			{
				ValidPreparationInstrumentID: &basteTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
				Index:       pointer.To[uint16](2),
				OptionIndex: 0,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &basteSkilletVPV.ID,
				Name:                     "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "butter-basted bone-in ribeye steak",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](700),
					Max: pointer.To[float32](900),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: atTemperatureState.ID,
				Notes:             "Steak internal temperature should reach 120-125°F (49-52°C) for medium-rare or 130°F (54°C) for medium",
				Ingredients:       []uint64{0}, // Index of the steak ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 7: Rest the steak
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: restPrep.ID,
		Index:         7,
		Notes:         "Immediately transfer steak to a large heatproof plate and pour pan juices on top. Let rest 5 to 10 minutes. Carve and serve.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
			Max: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "butter-basted bone-in ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 700,
					Max: pointer.To[float32](900),
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &restTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &restPlateVPV.ID,
				Name:                     "serving plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "rested pan-seared butter-basted bone-in ribeye steak",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](700),
					Max: pointer.To[float32](900),
				},
			},
		},
	}

	steps := []*mealplanning.RecipeStepCreationRequestInput{
		step0,
		step1,
		step2,
		step3,
		step4,
		step5,
		step6,
		step7,
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
			Name:                "Pan-Seared Butter-Basted Steak",
			Slug:                "pan-seared-butter-basted-steak",
			Source:              "https://www.seriouseats.com/butter-basted-pan-seared-steaks-recipe",
			Description:         "Thick and meaty pan-seared steak, slicked with butter and infused with flavor from aromatics. This recipe is designed for very large steaks, at least one and a half inches thick and weighing 24 to 32 ounces (700 to 900g) with the bone in.",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 2,
				Max: pointer.To[float32](3),
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps:             steps,
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
