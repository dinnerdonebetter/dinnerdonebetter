package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func SousVidePorkChopsRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	heatPrep := enums.Preparations["heat"]
	seasonPrep := enums.Preparations["season"]
	bagPrep := enums.Preparations["bag"]
	sealPrep := enums.Preparations["seal"]
	sousVidePrep := enums.Preparations["sous-vide"]
	dryPrep := enums.Preparations["dry"]
	panSearPrep := enums.Preparations["pan-sear"]
	bastePrep := enums.Preparations["baste"]
	restPrep := enums.Preparations["rest"]

	// Get ingredients
	porkChop := enums.Ingredients["pork chop"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	vegetableOil := enums.Ingredients["vegetable oil"]
	butter := enums.Ingredients["butter"]
	thyme := enums.Ingredients["thyme"]
	rosemary := enums.Ingredients["rosemary"]
	garlic := enums.Ingredients["garlic"]
	shallot := enums.Ingredients["shallot"]
	water := enums.Ingredients["water"]

	// Get measurement units
	unitMeasurement := enums.MeasurementUnits["unit"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]
	quartMeasurement := enums.MeasurementUnits["quart"]

	// Get instruments
	sousVideCooker := enums.Instruments["sous vide cooker"]
	paperTowels := enums.Instruments["paper towels"]
	tongs := enums.Instruments["tongs"]
	spoon := enums.Instruments["spoon"]
	bareHands := enums.Instruments["bare hands"]

	// Get vessels
	waterBath := enums.Vessels["water bath"]
	vacuumBag := enums.Vessels["vacuum bag"]
	castIronSkillet := enums.Vessels["cast iron skillet"]
	wireRack := enums.Vessels["wire rack"]
	bakingSheet := enums.Vessels["baking sheet"]
	servingPlate := enums.Vessels["serving plate"]

	// Get ingredient states for completion conditions
	atTemperatureState := enums.IngredientStates["at temperature"]
	smokingState := enums.IngredientStates["smoking"]
	brownedState := enums.IngredientStates["browned"]

	// Get bridge table entries
	// Heat preparation bridges (for preheating water bath)
	heatSousVideCookerVPI := enums.PreparationInstruments[heatPrep.ID][sousVideCooker.ID]
	heatWaterVIP := enums.IngredientPreparations[heatPrep.ID][water.ID]
	waterQuartVIMU := enums.IngredientMeasurementUnits[water.ID][quartMeasurement.ID]

	// Season preparation bridges
	seasonPorkChopVIP := enums.IngredientPreparations[seasonPrep.ID][porkChop.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	porkChopUnitVIMU := enums.IngredientMeasurementUnits[porkChop.ID][unitMeasurement.ID]
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]
	seasonBareHandsVPI := enums.PreparationInstruments[seasonPrep.ID][bareHands.ID]

	// Bag preparation bridges
	bagVacuumBagVPV := enums.PreparationVessels[bagPrep.ID][vacuumBag.ID]

	// Seal preparation bridges
	sealVacuumBagVPV := enums.PreparationVessels[sealPrep.ID][vacuumBag.ID]

	// Sous vide preparation bridges
	sousVideCookerVPI := enums.PreparationInstruments[sousVidePrep.ID][sousVideCooker.ID]
	sousVideWaterBathVPV := enums.PreparationVessels[sousVidePrep.ID][waterBath.ID]

	// Dry preparation bridges
	dryPorkChopVIP := enums.IngredientPreparations[dryPrep.ID][porkChop.ID]
	dryPaperTowelsVPI := enums.PreparationInstruments[dryPrep.ID][paperTowels.ID]

	// Heat preparation bridges (for heating skillet)
	heatSkilletVPV := enums.PreparationVessels[heatPrep.ID][castIronSkillet.ID]
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][vegetableOil.ID]
	oilTablespoonVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][tablespoonMeasurement.ID]

	// Pan-sear preparation bridges
	panSearPorkChopVIP := enums.IngredientPreparations[panSearPrep.ID][porkChop.ID]
	panSearTongsVPI := enums.PreparationInstruments[panSearPrep.ID][tongs.ID]
	panSearSkilletVPV := enums.PreparationVessels[panSearPrep.ID][castIronSkillet.ID]

	// Baste preparation bridges
	bastePorkChopVIP := enums.IngredientPreparations[bastePrep.ID][porkChop.ID]
	basteButterVIP := enums.IngredientPreparations[bastePrep.ID][butter.ID]
	basteThymeVIP := enums.IngredientPreparations[bastePrep.ID][thyme.ID]
	basteRosemaryVIP := enums.IngredientPreparations[bastePrep.ID][rosemary.ID]
	basteGarlicVIP := enums.IngredientPreparations[bastePrep.ID][garlic.ID]
	basteShallotVIP := enums.IngredientPreparations[bastePrep.ID][shallot.ID]
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]
	thymeSprigVIMU := enums.IngredientMeasurementUnits[thyme.ID][sprigMeasurement.ID]
	rosemarySprigVIMU := enums.IngredientMeasurementUnits[rosemary.ID][sprigMeasurement.ID]
	garlicUnitVIMU := enums.IngredientMeasurementUnits[garlic.ID][unitMeasurement.ID]
	shallotGramVIMU := enums.IngredientMeasurementUnits[shallot.ID][gramMeasurement.ID]
	basteTongsVPI := enums.PreparationInstruments[bastePrep.ID][tongs.ID]
	basteSpoonVPI := enums.PreparationInstruments[bastePrep.ID][spoon.ID]
	basteServingPlateVPV := enums.PreparationVessels[bastePrep.ID][servingPlate.ID]

	// Rest preparation bridges
	restPorkChopVIP := enums.IngredientPreparations[restPrep.ID][porkChop.ID]
	restTongsVPI := enums.PreparationInstruments[restPrep.ID][tongs.ID]
	restWireRackVPV := enums.PreparationVessels[restPrep.ID][wireRack.ID]
	restBakingSheetVPV := enums.PreparationVessels[restPrep.ID][bakingSheet.ID]

	// Step 0: Preheat water bath
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       heatPrep.ID,
		Index:                0,
		ExplicitInstructions: "Place an immersion circulator in a water bath and set the circulator to the desired final temperature. For medium-rare, set to 140°F (60°C). Allow the water bath to come to temperature.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](60), // 140°F = 60°C (medium-rare default)
			Max: pointer.To[float32](60),
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &heatSousVideCookerVPI.ID,
				Name:                         "sous vide cooker",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &heatWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterQuartVIMU.ID,
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 12,
					Max: pointer.To[float32](12),
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "sous vide cooker",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
			{
				Name:              "heated water bath",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             1,
				MeasurementUnitID: &quartMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](12),
					Max: pointer.To[float32](12),
				},
			},
		},
	}

	// Step 1: Season pork chops
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       seasonPrep.ID,
		Index:                1,
		ExplicitInstructions: "Season the pork chops generously with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &seasonPorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				Name:                             "bone-in pork rib chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltGramVIMU.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste: true,
			},
			{
				ValidIngredientPreparationID:     &seasonPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperGramVIMU.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste: true,
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
				Name:              "seasoned pork chops",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 2: Bag pork chops
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       bagPrep.ID,
		Index:                2,
		ExplicitInstructions: "Place the pork chops in vacuum-seal or zipper-lock bags.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "seasoned pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &bagVacuumBagVPV.ID,
				Name:                     "vacuum bag",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "bagged pork chops",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
			{
				Name:  "vacuum bag with pork chops",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 3: Seal bags
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       sealPrep.ID,
		Index:                3,
		ExplicitInstructions: "Seal the bags. If using zipper-lock bags, use the displacement method: seal the bag almost entirely closed, slowly lower into the water to press out air, then seal completely just above the waterline.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "bagged pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &sealVacuumBagVPV.ID,
				Name:                            "vacuum bag with pork chops",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "sealed bagged pork chops",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 4: Cook sous vide
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       sousVidePrep.ID,
		Index:                4,
		ExplicitInstructions: "Place the sealed bagged pork chops in the preheated water bath and cook for the recommended time. For 1 1/2 inch thick chops at 140°F (60°C), cook for 1 to 4 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](3600),  // 1 hour minimum
			Max: pointer.To[uint32](14400), // 4 hours maximum
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](60), // 140°F = 60°C
			Max: pointer.To[float32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "sealed bagged pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "water bath",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationInstrumentID:    &sousVideCookerVPI.ID,
				Name:                            "sous vide cooker",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex: pointer.To[uint64](0),
				ValidPreparationVesselID: &sousVideWaterBathVPV.ID,
				Name:                     "preheated water bath",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "sous vide cooked pork chops",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: atTemperatureState.ID,
				Notes:             "Pork chops should reach internal temperature and be held for at least 1 hour",
				Ingredients:       []uint64{0}, // Index of the pork chop ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 5: Pat dry pork chops
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       dryPrep.ID,
		Index:                5,
		ExplicitInstructions: "Remove the pork chops from the water bath and bag. Carefully pat dry with paper towels.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &dryPorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				Name:                             "sous vide cooked pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &dryPaperTowelsVPI.ID,
				Name:                         "paper towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "dried sous vide pork chops",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 6: Heat oil in skillet (Optional - Pan Finish)
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       heatPrep.ID,
		Index:                6,
		Optional:             true,
		ExplicitInstructions: "To finish in a pan: Turn on your vents and open your windows. Add 2 tablespoons canola oil to a heavy cast iron or stainless steel skillet, place it over the hottest burner you have, and preheat the skillet until the oil starts to smoke.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](220), // Very high heat
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilTablespoonVIMU.ID,
				Name:                             "canola oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
					Max: pointer.To[float32](4),
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
				Name:  "smoking oil",
				Type:  mealplanning.RecipeStepProductIngredientType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
					Max: pointer.To[float32](4),
				},
				MeasurementUnitID: &tablespoonMeasurement.ID,
			},
			{
				Name:  "heated skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
					Max: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: smokingState.ID,
				Notes:             "Oil should begin to smoke",
				Ingredients:       []uint64{0}, // Index of the oil ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 7: Pan-sear first side (Optional - Pan Finish)
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       panSearPrep.ID,
		Index:                7,
		Optional:             true,
		ExplicitInstructions: "Using your fingers or a set of tongs, gently lay two pork chops in the skillet. If desired, add 1 tablespoon butter; for a cleaner-tasting sear, omit butter at this stage. Carefully lift and peek under the pork as it cooks to gauge how quickly it is browning. Let it continue to cook until the crust is deep brown and very crisp, about 45 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](45),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &panSearPorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				Name:                             "dried sous vide pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2, // Two at a time
				},
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
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &panSearSkilletVPV.ID,
				Name:                            "heated skillet with smoking oil",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "partially seared pork chops",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				Name:  "heated cast iron skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
					Max: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Flip and baste (Optional - Pan Finish)
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       bastePrep.ID,
		Index:                8,
		Optional:             true,
		ExplicitInstructions: "Flip the pork chops. If desired, add 1 more tablespoon butter, along with half of the thyme, rosemary, garlic, and/or shallots. Spoon the butter over the pork chops as they cook, if using. Continue cooking until the second side is browned, about 45 seconds longer. When the pork is browned, pick it up with a pair of tongs, rotate it sideways, and make sure to brown the edges as well. Transfer the cooked pork chops to a wire rack set over a rimmed baking sheet. Discard aromatics. Repeat with the remaining pork chops, butter, and aromatics, adding additional oil to the skillet if necessary.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](45),
			Max: pointer.To[uint32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &bastePorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				Name:                             "partially seared pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &basteButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				Name:                             "butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
					Max: pointer.To[float32](2),
				},
				Optional: true,
			},
			{
				ValidIngredientPreparationID:     &basteThymeVIP.ID,
				ValidIngredientMeasurementUnitID: &thymeSprigVIMU.ID,
				Name:                             "thyme",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Optional: true,
			},
			{
				ValidIngredientPreparationID:     &basteRosemaryVIP.ID,
				ValidIngredientMeasurementUnitID: &rosemarySprigVIMU.ID,
				Name:                             "rosemary",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Optional: true,
			},
			{
				ValidIngredientPreparationID:     &basteGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicUnitVIMU.ID,
				Name:                             "garlic cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
				Optional: true,
			},
			{
				ValidIngredientPreparationID:     &basteShallotVIP.ID,
				ValidIngredientMeasurementUnitID: &shallotGramVIMU.ID,
				Name:                             "shallots, thinly sliced",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 30, // About 1 shallot
				},
				Optional: true,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &basteTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationInstrumentID: &basteSpoonVPI.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seared and basted pork chops",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
			{
				Name:  "heated cast iron skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
					Max: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: brownedState.ID,
				Notes:             "Both sides and edges should be deep brown and very crisp",
				Ingredients:       []uint64{0}, // Index of the pork chop ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 9: Rest on wire rack (Optional - Pan Finish)
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       restPrep.ID,
		Index:                9,
		Optional:             true,
		ExplicitInstructions: "Transfer the cooked pork chops to a wire rack set over a rimmed baking sheet. Let the chops rest for 3 to 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180), // 3 minutes
			Max: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &restPorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				Name:                             "seared and basted pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
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
				ValidPreparationVesselID: &restWireRackVPV.ID,
				Name:                     "wire rack",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationVesselID: &restBakingSheetVPV.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "rested pork chops",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 10: Reheat drippings and pour over (Optional - Pan Finish)
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:       bastePrep.ID,
		Index:                10,
		Optional:             true,
		ExplicitInstructions: "Just before serving, reheat the drippings in the pan until sizzling-hot, then pour them over the pork chops in order to re-crisp their exteriors. Serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "rested pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
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
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "cast iron skillet with drippings",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationVesselID: &basteServingPlateVPV.ID,
				Name:                     "serving plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "pan-finished sous vide pork chops",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
			Name:                "Sous Vide Pork Chops",
			Slug:                "sous-vide-pork-chops",
			Source:              "https://www.seriouseats.com/sous-vide-pork-chops-recipe",
			Description:         "Using an immersion sous vide cooker is the easy, foolproof way to guarantee extra-juicy pork chops. Cooking sous vide ensures pork chops are perfectly cooked from edge to edge by maintaining a precise water temperature that precludes overcooking and preserves moisture. The method allows for greater control over texture by adjusting the cooking temperature, offering options from a pink and tender medium-rare to a traditional well-done chop. A high-heat finish, in a skillet or on the grill, gives the chops a crisp, browned crust and keeps the interior juicy.",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 4,
			},
			PortionName:       "pork chop",
			PluralPortionName: "pork chops",
			EligibleForMeals:  true,
			Steps:             []*mealplanning.RecipeStepCreationRequestInput{step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10},
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
