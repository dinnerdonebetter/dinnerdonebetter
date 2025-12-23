package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func PanSearedButterBastedSteakRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	dryPrep := enums.Preparations["dry"]
	seasonPrep := enums.Preparations["season"]
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

	// Get measurement units
	unitMeasurement := enums.MeasurementUnits["unit"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	milliliterMeasurement := enums.MeasurementUnits["milliliter"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]

	// Get instruments
	paperTowels := enums.Instruments["paper towels"]
	tongs := enums.Instruments["tongs"]
	spoon := enums.Instruments["spoon"]
	thermometer := enums.Instruments["instant-read thermometer"]

	// Get vessels
	sheetPan := enums.Vessels["sheet pan"]
	castIronSkillet := enums.Vessels["cast iron skillet"]
	servingPlate := enums.Vessels["serving plate"]

	// Get ingredient states for completion conditions
	_ = enums.IngredientStates["dry"] // dry state validated but not used
	smokingState := enums.IngredientStates["smoking"]
	atTemperatureState := enums.IngredientStates["at temperature"]

	// Get bridge table entries
	// Dry preparation bridges
	dryRibeyeVIP := enums.IngredientPreparations[dryPrep.ID][ribeye.ID]
	ribeyeUnitVIMU := enums.IngredientMeasurementUnits[ribeye.ID][unitMeasurement.ID]
	dryPaperTowelsVPI := enums.PreparationInstruments[dryPrep.ID][paperTowels.ID]

	// Season preparation bridges
	_ = enums.IngredientPreparations[seasonPrep.ID][ribeye.ID] // validated but not used
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]
	seasonSheetPanVPV := enums.PreparationVessels[seasonPrep.ID][sheetPan.ID]

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
	basteShallotVIP := enums.IngredientPreparations[bastePrep.ID][shallot.ID]
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
	step0ID := identifiers.New()
	step0IngredientID := identifiers.New()
	_ = identifiers.New() // completion condition ID placeholder
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:            step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dryPrep.ID,
		Index:           0,
		Notes:           "Carefully pat steak dry with paper towels.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step0IngredientID,
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &dryRibeyeVIP.ID,
				ValidIngredientMeasurementUnitID: &ribeyeUnitVIMU.ID,
				IngredientID:                     &ribeye.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                          identifiers.New(),
				BelongsToRecipeStep:         step0ID,
				ValidPreparationInstrumentID: &dryPaperTowelsVPI.ID,
				InstrumentID:                &paperTowels.ID,
				Name:                        "paper towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "dried ribeye steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Season the steak
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:            step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           1,
		Notes:           "Season liberally on all sides, including edges, with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex:   pointer.To[uint64](0),
				IngredientID:                     &ribeye.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "dried ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
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
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                        identifiers.New(),
				BelongsToRecipeStep:       step1ID,
				ValidPreparationVesselID:  &seasonSheetPanVPV.ID,
				VesselID:                  &sheetPan.ID,
				Name:                      "sheet pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "seasoned ribeye steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Rest the steak (optional - at room temperature or refrigerated)
	step2ID := identifiers.New()
	step2IngredientID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:            step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           2,
		Optional:        true,
		Notes:           "If desired, let steak rest at room temperature for 45 minutes, or refrigerated, loosely covered, up to 3 days.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](2700), // 45 minutes minimum
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step2IngredientID,
				BelongsToRecipeStep:              step2ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex:   pointer.To[uint64](0),
				IngredientID:                     &ribeye.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "seasoned ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                        identifiers.New(),
				BelongsToRecipeStep:       step2ID,
				ValidPreparationVesselID:  &restSheetPanVPV.ID,
				VesselID:                  &sheetPan.ID,
				Name:                      "sheet pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "rested seasoned ribeye steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Heat oil until smoking
	step3ID := identifiers.New()
	step3OilIngredientID := identifiers.New()
	step3CompletionConditionID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:            step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           3,
		Notes:           "In a 12-inch heavy-bottomed cast iron skillet, heat oil over high heat until just beginning to smoke.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](200), // High heat, approximately 200°C
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step3OilIngredientID,
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilMilliliterVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                milliliterMeasurement.ID,
				Name:                             "vegetable or canola oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 60, // 1/4 cup = 60 ml
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                        identifiers.New(),
				BelongsToRecipeStep:       step3ID,
				ValidPreparationVesselID:  &heatSkilletVPV.ID,
				VesselID:                  &castIronSkillet.ID,
				Name:                      "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "heated smoking oil",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &milliliterMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](60),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step3CompletionConditionID,
				BelongsToRecipeStep: step3ID,
				IngredientStateID:   smokingState.ID,
				Notes:               "Oil should be just beginning to smoke",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step3CompletionConditionID,
						RecipeStepIngredient:                   step3OilIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 4: Pan-sear the steak
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:            step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   panSearPrep.ID,
		Index:           4,
		Notes:           "Carefully add steak to the hot skillet and cook, flipping frequently, until a pale golden-brown crust starts to develop, about 4 minutes total.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](240), // 4 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex:   pointer.To[uint64](0),
				IngredientID:                     &ribeye.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "rested seasoned ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:   pointer.To[uint64](0),
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                milliliterMeasurement.ID,
				Name:                             "heated smoking oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 60,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                          identifiers.New(),
				BelongsToRecipeStep:         step4ID,
				ValidPreparationInstrumentID: &panSearTongsVPI.ID,
				InstrumentID:                &tongs.ID,
				Name:                        "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                        identifiers.New(),
				BelongsToRecipeStep:       step4ID,
				ValidPreparationVesselID:  &panSearSkilletVPV.ID,
				VesselID:                  &castIronSkillet.ID,
				Name:                      "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "pan-seared ribeye steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Baste the steak
	step5ID := identifiers.New()
	step5SteakIngredientID := identifiers.New()
	step5CompletionConditionID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:            step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   bastePrep.ID,
		Index:           5,
		Notes:           "Add butter, herbs (if using), and shallot (if using) to skillet and continue to cook, flipping steak occasionally and basting any light spots with foaming butter. If butter begins to smoke excessively or steak begins to burn, reduce heat to medium. To baste, tilt pan slightly so that butter collects by handle. Use a spoon to pick up butter and pour it over steak, aiming at light spots. Continue flipping and basting until an instant-read thermometer inserted into thickest part of tenderloin side registers 120 to 125°F (49 to 52°C) for medium-rare or 130°F (54°C) for medium, 8 to 10 minutes total.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](480), // 8 minutes
			Max: pointer.To[uint32](600), // 10 minutes
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](49), // 120°F = 49°C
			Max: pointer.To[float32](54), // 130°F = 54°C
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step5SteakIngredientID,
				BelongsToRecipeStep:              step5ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:   pointer.To[uint64](0),
				IngredientID:                     &ribeye.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "pan-seared ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &basteButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterGramVIMU.ID,
				IngredientID:                     &butter.ID,
				MeasurementUnitID:                gramMeasurement.ID,
				Name:                             "unsalted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 45, // 3 tablespoons = 45g
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &basteThymeVIP.ID,
				ValidIngredientMeasurementUnitID: &thymeSprigVIMU.ID,
				IngredientID:                     &thyme.ID,
				MeasurementUnitID:                sprigMeasurement.ID,
				Name:                             "thyme",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6,
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &basteRosemaryVIP.ID,
				ValidIngredientMeasurementUnitID: &rosemarySprigVIMU.ID,
				IngredientID:                     &rosemary.ID,
				MeasurementUnitID:                sprigMeasurement.ID,
				Name:                             "rosemary",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6,
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &basteShallotVIP.ID,
				ValidIngredientMeasurementUnitID: &shallotGramVIMU.ID,
				IngredientID:                     &shallot.ID,
				MeasurementUnitID:                gramMeasurement.ID,
				Name:                             "finely sliced shallots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0, // About 1 large shallot, optional
				},
				Optional: true,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                          identifiers.New(),
				BelongsToRecipeStep:         step5ID,
				ValidPreparationInstrumentID: &basteSpoonVPI.ID,
				InstrumentID:                &spoon.ID,
				Name:                        "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                          identifiers.New(),
				BelongsToRecipeStep:         step5ID,
				ValidPreparationInstrumentID: &basteThermometerVPI.ID,
				InstrumentID:                &thermometer.ID,
				Name:                        "instant-read thermometer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                          identifiers.New(),
				BelongsToRecipeStep:         step5ID,
				ValidPreparationInstrumentID: &basteTongsVPI.ID,
				InstrumentID:                &tongs.ID,
				Name:                        "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                        identifiers.New(),
				BelongsToRecipeStep:       step5ID,
				ValidPreparationVesselID:  &basteSkilletVPV.ID,
				VesselID:                  &castIronSkillet.ID,
				Name:                      "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "butter-basted ribeye steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step5CompletionConditionID,
				BelongsToRecipeStep: step5ID,
				IngredientStateID:   atTemperatureState.ID,
				Notes:               "Steak internal temperature should reach 120-125°F (49-52°C) for medium-rare or 130°F (54°C) for medium",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step5CompletionConditionID,
						RecipeStepIngredient:                   step5SteakIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 6: Rest the steak
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:            step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           6,
		Notes:           "Immediately transfer steak to a large heatproof plate and pour pan juices on top. Let rest 5 to 10 minutes. Carve and serve.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
			Max: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:   pointer.To[uint64](0),
				IngredientID:                     &ribeye.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "butter-basted ribeye steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                          identifiers.New(),
				BelongsToRecipeStep:         step6ID,
				ValidPreparationInstrumentID: &restTongsVPI.ID,
				InstrumentID:                &tongs.ID,
				Name:                        "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                        identifiers.New(),
				BelongsToRecipeStep:       step6ID,
				ValidPreparationVesselID:  &restPlateVPV.ID,
				VesselID:                  &servingPlate.ID,
				Name:                      "serving plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "rested pan-seared butter-basted ribeye steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		{
			ID:                  recipeID,
			CreatedByUser:       userID,
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
			Steps:             []*mealplanning.RecipeStepDatabaseCreationInput{step0, step1, step2, step3, step4, step5, step6},
		},
	}
}
