package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func SousVidePorkChopsRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

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

	// Get measurement units
	unitMeasurement := enums.MeasurementUnits["unit"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]

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
	heatWaterBathVPV := enums.PreparationVessels[heatPrep.ID][waterBath.ID]
	heatSousVideCookerVPI := enums.PreparationInstruments[heatPrep.ID][sousVideCooker.ID]

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
	basteSkilletVPV := enums.PreparationVessels[bastePrep.ID][castIronSkillet.ID]
	basteServingPlateVPV := enums.PreparationVessels[bastePrep.ID][servingPlate.ID]

	// Rest preparation bridges
	restPorkChopVIP := enums.IngredientPreparations[restPrep.ID][porkChop.ID]
	restTongsVPI := enums.PreparationInstruments[restPrep.ID][tongs.ID]
	restWireRackVPV := enums.PreparationVessels[restPrep.ID][wireRack.ID]
	restBakingSheetVPV := enums.PreparationVessels[restPrep.ID][bakingSheet.ID]

	// Step 0: Preheat water bath
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           0,
		Notes:           "Place an immersion circulator in a water bath and set the circulator to the desired final temperature. For medium-rare, set to 140°F (60°C). Allow the water bath to come to temperature.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](60), // 140°F = 60°C (medium-rare default)
			Max: pointer.To[float32](60),
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
				ValidPreparationInstrumentID: &heatSousVideCookerVPI.ID,
				InstrumentID:                 &sousVideCooker.ID,
				Name:                         "sous vide cooker",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &heatWaterBathVPV.ID,
				VesselID:                 &waterBath.ID,
				Name:                     "water bath",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "preheated water bath",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
			},
		},
	}

	// Step 1: Season pork chops
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           1,
		Notes:           "Season pork chops generously with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &seasonPorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				IngredientID:                     &porkChop.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "bone-in pork rib chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltGramVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                gramMeasurement.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &seasonPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperGramVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                gramMeasurement.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste: true,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
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
				BelongsToRecipeStep: step1ID,
				Name:                "seasoned pork chops",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 2: Bag pork chops
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   bagPrep.ID,
		Index:           2,
		Notes:           "Place pork chops in vacuum-seal or zipper-lock bags.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &porkChop.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "seasoned pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &bagVacuumBagVPV.ID,
				VesselID:                 &vacuumBag.ID,
				Name:                     "vacuum bag",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "bagged pork chops",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "vacuum bag with pork chops",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 3: Seal bags
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   sealPrep.ID,
		Index:           3,
		Notes:           "Seal bags. If using zipper-lock bags, use the displacement method: seal the bag almost entirely closed, slowly lower into water to press out air, then seal completely just above the waterline.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &porkChop.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "bagged pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &sealVacuumBagVPV.ID,
				VesselID:                        &vacuumBag.ID,
				Name:                            "vacuum bag with pork chops",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "sealed bagged pork chops",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 4: Cook sous vide
	step4ID := identifiers.New()
	step4PorkChopIngredientID := identifiers.New()
	step4CompletionConditionID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   sousVidePrep.ID,
		Index:           4,
		Notes:           "Place sealed bagged pork chops in the preheated water bath and cook for the recommended time. For 1 1/2 inch thick chops at 140°F (60°C), cook for 1 to 4 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](3600),  // 1 hour minimum
			Max: pointer.To[uint32](14400), // 4 hours maximum
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](60), // 140°F = 60°C
			Max: pointer.To[float32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step4PorkChopIngredientID,
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &porkChop.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "sealed bagged pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &sousVideCookerVPI.ID,
				InstrumentID:                 &sousVideCooker.ID,
				Name:                         "sous vide cooker",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &sousVideWaterBathVPV.ID,
				VesselID:                 &waterBath.ID,
				Name:                     "preheated water bath",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "sous vide cooked pork chops",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step4CompletionConditionID,
				BelongsToRecipeStep: step4ID,
				IngredientStateID:   atTemperatureState.ID,
				Notes:               "Pork chops should reach internal temperature and be held for at least 1 hour",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step4CompletionConditionID,
						RecipeStepIngredient:                   step4PorkChopIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 5: Pat dry pork chops
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dryPrep.ID,
		Index:           5,
		Notes:           "Remove pork chops from water bath and bag. Carefully pat dry with paper towels.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &dryPorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				IngredientID:                     &porkChop.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "sous vide cooked pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &dryPaperTowelsVPI.ID,
				InstrumentID:                 &paperTowels.ID,
				Name:                         "paper towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "dried sous vide pork chops",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 6: Heat oil in skillet (Optional - Pan Finish)
	step6ID := identifiers.New()
	step6OilIngredientID := identifiers.New()
	step6CompletionConditionID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           6,
		Optional:        true,
		Notes:           "To finish in a pan: Turn on your vents and open your windows. Add 2 tablespoons canola oil to a heavy cast iron or stainless steel skillet, place it over the hottest burner you have, and preheat skillet until the oil starts to smoke.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](220), // Very high heat
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step6OilIngredientID,
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilTablespoonVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "canola oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
					Max: pointer.To[float32](4),
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &heatSkilletVPV.ID,
				VesselID:                 &castIronSkillet.ID,
				Name:                     "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "heated skillet with smoking oil",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step6CompletionConditionID,
				BelongsToRecipeStep: step6ID,
				IngredientStateID:   smokingState.ID,
				Notes:               "Oil should begin to smoke",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step6CompletionConditionID,
						RecipeStepIngredient:                   step6OilIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 7: Pan-sear first side (Optional - Pan Finish)
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   panSearPrep.ID,
		Index:           7,
		Optional:        true,
		Notes:           "Using your fingers or a set of tongs, gently lay two pork chops in skillet. If desired, add 1 tablespoon butter; for a cleaner-tasting sear, omit butter at this stage. Carefully lift and peek under pork as it cooks to gauge how quickly it is browning. Let it continue to cook until the crust is deep brown and very crisp, about 45 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](45),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &panSearPorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				IngredientID:                     &porkChop.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "dried sous vide pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2, // Two at a time
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: &panSearTongsVPI.ID,
				InstrumentID:                 &tongs.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step7ID,
				ProductOfRecipeStepIndex: pointer.To[uint64](6),
				ValidPreparationVesselID: &panSearSkilletVPV.ID,
				VesselID:                 &castIronSkillet.ID,
				Name:                     "heated skillet with smoking oil",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "partially seared pork chops",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 8: Flip and baste (Optional - Pan Finish)
	step8ID := identifiers.New()
	step8PorkChopIngredientID := identifiers.New()
	step8CompletionConditionID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   bastePrep.ID,
		Index:           8,
		Optional:        true,
		Notes:           "Flip pork chops. If desired, add 1 more tablespoon butter, along with half of the thyme, rosemary, garlic, and/or shallots. Spoon butter over pork chops as they cook, if using. Continue cooking until second side is browned, about 45 seconds longer. When pork is browned, pick it up with a pair of tongs, rotate it sideways, and make sure to brown the edges as well. Transfer cooked pork chops to a wire rack set over a rimmed baking sheet. Discard aromatics. Repeat with remaining pork chops, butter, and aromatics, adding additional oil to skillet if necessary.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](45),
			Max: pointer.To[uint32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step8PorkChopIngredientID,
				BelongsToRecipeStep:              step8ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &bastePorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				IngredientID:                     &porkChop.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "partially seared pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &basteButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				IngredientID:                     &butter.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
					Max: pointer.To[float32](2),
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &basteThymeVIP.ID,
				ValidIngredientMeasurementUnitID: &thymeSprigVIMU.ID,
				IngredientID:                     &thyme.ID,
				MeasurementUnitID:                sprigMeasurement.ID,
				Name:                             "thyme",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &basteRosemaryVIP.ID,
				ValidIngredientMeasurementUnitID: &rosemarySprigVIMU.ID,
				IngredientID:                     &rosemary.ID,
				MeasurementUnitID:                sprigMeasurement.ID,
				Name:                             "rosemary",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &basteGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicUnitVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "garlic cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &basteShallotVIP.ID,
				ValidIngredientMeasurementUnitID: &shallotGramVIMU.ID,
				IngredientID:                     &shallot.ID,
				MeasurementUnitID:                gramMeasurement.ID,
				Name:                             "shallots, thinly sliced",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 30, // About 1 shallot
				},
				Optional: true,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step7ID,
				ValidPreparationInstrumentID: &basteTongsVPI.ID,
				InstrumentID:                 &tongs.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step7ID,
				ValidPreparationInstrumentID: &basteSpoonVPI.ID,
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
				BelongsToRecipeStep:      step7ID,
				ValidPreparationVesselID: &basteSkilletVPV.ID,
				VesselID:                 &castIronSkillet.ID,
				Name:                     "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "seared and basted pork chops",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step8CompletionConditionID,
				BelongsToRecipeStep: step8ID,
				IngredientStateID:   brownedState.ID,
				Notes:               "Both sides and edges should be deep brown and very crisp",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step8CompletionConditionID,
						RecipeStepIngredient:                   step8PorkChopIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 9: Rest on wire rack (Optional - Pan Finish)
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           9,
		Optional:        true,
		Notes:           "Transfer cooked pork chops to a wire rack set over a rimmed baking sheet. Let chops rest for 3 to 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180), // 3 minutes
			Max: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step9ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &restPorkChopVIP.ID,
				ValidIngredientMeasurementUnitID: &porkChopUnitVIMU.ID,
				IngredientID:                     &porkChop.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "seared and basted pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step8ID,
				ValidPreparationInstrumentID: &restTongsVPI.ID,
				InstrumentID:                 &tongs.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: &restWireRackVPV.ID,
				VesselID:                 &wireRack.ID,
				Name:                     "wire rack",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: &restBakingSheetVPV.ID,
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
				BelongsToRecipeStep: step8ID,
				Name:                "rested pork chops",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 10: Reheat drippings and pour over (Optional - Pan Finish)
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   bastePrep.ID,
		Index:           10,
		Optional:        true,
		Notes:           "Just before serving, reheat the drippings in the pan until sizzling-hot, then pour them over pork chops in order to re-crisp their exteriors. Serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &porkChop.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "rested pork chops",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step10ID,
				ValidPreparationInstrumentID: &basteSpoonVPI.ID,
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
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: &basteSkilletVPV.ID,
				VesselID:                 &castIronSkillet.ID,
				Name:                     "cast iron skillet with drippings",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: &basteServingPlateVPV.ID,
				VesselID:                 &servingPlate.ID,
				Name:                     "serving plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "pan-finished sous vide pork chops",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		{
			ID:                  recipeID,
			CreatedByUser:       userID,
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
			Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
				step0, step1, step2, step3, step4,
				step5, step6, step7, step8, step9, step10,
			},
		},
	}
}
