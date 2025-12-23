package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func SousVideChickenBreastRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	heatPrep := enums.Preparations["heat"]
	seasonPrep := enums.Preparations["season"]
	bagPrep := enums.Preparations["bag"]
	sousVidePrep := enums.Preparations["sous-vide"]
	dryPrep := enums.Preparations["dry"]
	panSearPrep := enums.Preparations["pan-sear"]
	grillPrep := enums.Preparations["grill"]
	restPrep := enums.Preparations["rest"]

	// Get ingredients
	chickenBreast := enums.Ingredients["chicken breast"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	thyme := enums.Ingredients["thyme"]
	rosemary := enums.Ingredients["rosemary"]
	vegetableOil := enums.Ingredients["vegetable oil"]

	// Get measurement units
	unitMeasurement := enums.MeasurementUnits["unit"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	milliliterMeasurement := enums.MeasurementUnits["milliliter"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]

	// Get instruments
	sousVideCooker := enums.Instruments["sous vide cooker"]
	paperTowels := enums.Instruments["paper towels"]
	spatula := enums.Instruments["spatula"]
	tongs := enums.Instruments["tongs"]
	bareHands := enums.Instruments["bare hands"]

	// Get vessels
	waterBath := enums.Vessels["water bath"]
	plasticBag := enums.Vessels["plastic bag"]
	vacuumBag := enums.Vessels["vacuum bag"]
	castIronSkillet := enums.Vessels["cast iron skillet"]
	grillVessel := enums.Vessels["grill"]
	servingPlate := enums.Vessels["serving plate"]

	// Get ingredient states for completion conditions
	atTemperatureState := enums.IngredientStates["at temperature"]

	// Get bridge table entries
	// Heat preparation bridges (for preheating water bath)
	heatWaterBathVPV := enums.PreparationVessels[heatPrep.ID][waterBath.ID]

	// Season preparation bridges
	seasonChickenVIP := enums.IngredientPreparations[seasonPrep.ID][chickenBreast.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	chickenUnitVIMU := enums.IngredientMeasurementUnits[chickenBreast.ID][unitMeasurement.ID]
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]
	seasonBareHandsVPI := enums.PreparationInstruments[seasonPrep.ID][bareHands.ID]

	// Bag preparation bridges
	_ = enums.IngredientPreparations[bagPrep.ID][chickenBreast.ID] // validated but not used
	bagThymeVIP := enums.IngredientPreparations[bagPrep.ID][thyme.ID]
	bagRosemaryVIP := enums.IngredientPreparations[bagPrep.ID][rosemary.ID]
	thymeSprigVIMU := enums.IngredientMeasurementUnits[thyme.ID][sprigMeasurement.ID]
	rosemarySprigVIMU := enums.IngredientMeasurementUnits[rosemary.ID][sprigMeasurement.ID]
	bagPlasticBagVPV := enums.PreparationVessels[bagPrep.ID][plasticBag.ID]
	bagVacuumBagVPV := enums.PreparationVessels[bagPrep.ID][vacuumBag.ID]

	// Sous vide preparation bridges
	_ = enums.IngredientPreparations[sousVidePrep.ID][chickenBreast.ID] // validated but not used
	sousVideCookerVPI := enums.PreparationInstruments[sousVidePrep.ID][sousVideCooker.ID]
	_ = enums.PreparationVessels[sousVidePrep.ID][waterBath.ID] // validated but not used

	// Dry preparation bridges (for patting dry after sous vide)
	_ = enums.IngredientPreparations[dryPrep.ID][chickenBreast.ID] // validated but not used
	_ = enums.PreparationInstruments[dryPrep.ID][paperTowels.ID]   // validated but not used

	// Pan-sear preparation bridges (for finishing)
	_ = enums.IngredientPreparations[panSearPrep.ID][chickenBreast.ID] // validated but not used
	panSearOilVIP := enums.IngredientPreparations[panSearPrep.ID][vegetableOil.ID]
	oilMilliliterVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][milliliterMeasurement.ID]
	panSearPaperTowelsVPI := enums.PreparationInstruments[panSearPrep.ID][paperTowels.ID]
	panSearSpatulaVPI := enums.PreparationInstruments[panSearPrep.ID][spatula.ID]
	panSearTongsVPI := enums.PreparationInstruments[panSearPrep.ID][tongs.ID]
	panSearSkilletVPV := enums.PreparationVessels[panSearPrep.ID][castIronSkillet.ID]

	// Grill preparation bridges (for finishing)
	_ = enums.IngredientPreparations[grillPrep.ID][chickenBreast.ID] // validated but not used
	grillPaperTowelsVPI := enums.PreparationInstruments[grillPrep.ID][paperTowels.ID]
	grillTongsVPI := enums.PreparationInstruments[grillPrep.ID][tongs.ID]
	grillGrillVPV := enums.PreparationVessels[grillPrep.ID][grillVessel.ID]

	// Rest preparation bridges
	restTongsVPI := enums.PreparationInstruments[restPrep.ID][tongs.ID]
	restPlateVPV := enums.PreparationVessels[restPrep.ID][servingPlate.ID]

	// Step 0: Preheat water bath
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           0,
		Notes:           "Preheat a water bath to 150°F (66°C) using a sous vide cooker. This temperature produces tender and juicy chicken, ideal for chicken salad when served cold.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](66), // 150°F = 66°C
			Max: pointer.To[float32](66), // 150°F = 66°C
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
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

	// Step 1: Season chicken
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           1,
		Notes:           "Season chicken generously with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &seasonChickenVIP.ID,
				ValidIngredientMeasurementUnitID: &chickenUnitVIMU.ID,
				IngredientID:                     &chickenBreast.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
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
				Name:                "seasoned chicken breasts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 2: Bag chicken
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   bagPrep.ID,
		Index:           2,
		Notes:           "Place chicken in zipper-lock bags or vacuum bags and add thyme or rosemary sprigs, if using. If using zipper-lock bags: Remove air by closing bags, leaving the last inch of the top unsealed. Slowly lower into preheated water bath, sealing bag completely just before it fully submerges. If using vacuum bags: Seal according to manufacturer's instructions.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &chickenBreast.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "seasoned chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &bagThymeVIP.ID,
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
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &bagRosemaryVIP.ID,
				ValidIngredientMeasurementUnitID: &rosemarySprigVIMU.ID,
				IngredientID:                     &rosemary.ID,
				MeasurementUnitID:                sprigMeasurement.ID,
				Name:                             "rosemary",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &bagPlasticBagVPV.ID,
				VesselID:                 &plasticBag.ID,
				Name:                     "zipper-lock bag",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
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
				Name:                "bagged seasoned chicken breasts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 3: Cook sous vide
	step3ID := identifiers.New()
	step3ChickenIngredientID := identifiers.New()
	step3CompletionConditionID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   sousVidePrep.ID,
		Index:           3,
		Notes:           "Add bagged chicken to preheated water bath and cook at 150°F (66°C) for 1 to 4 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](3600),  // 1 hour minimum
			Max: pointer.To[uint32](14400), // 4 hours maximum
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](66), // 150°F = 66°C
			Max: pointer.To[float32](66), // 150°F = 66°C
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step3ChickenIngredientID,
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &chickenBreast.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "bagged seasoned chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step3ID,
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
				BelongsToRecipeStep:      step3ID,
				ProductOfRecipeStepIndex: pointer.To[uint64](0),
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
				BelongsToRecipeStep: step3ID,
				Name:                "sous vide cooked chicken breasts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step3CompletionConditionID,
				BelongsToRecipeStep: step3ID,
				IngredientStateID:   atTemperatureState.ID,
				Notes:               "Chicken should reach 150°F (66°C) and be held at that temperature for at least 1 hour",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step3CompletionConditionID,
						RecipeStepIngredient:                   step3ChickenIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 4a: Finish in pan (optional)
	step4aID := identifiers.New()
	step4aCompletionConditionID := identifiers.New()
	step4a := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4aID,
		BelongsToRecipe: recipeID,
		PreparationID:   panSearPrep.ID,
		Index:           4,
		Optional:        true,
		Notes:           "Turn on your vents and open your windows. Remove chicken from water bath and bag. Discard herbs, if using. Carefully pat chicken dry with paper towels. Heat the oil in a heavy cast iron or stainless steel skillet over medium-high heat until shimmering. Gently lay chicken in skillet, skin side down, using your fingers or a set of tongs. Hold chicken down flat in pan with a flexible metal spatula or your fingers (be careful of splattering oil). Cook until golden brown and crisp, about 2 minutes. Remove from pan and let rest until cool enough to handle, about 2 minutes. Remove bones with your fingertips.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120), // 2 minutes
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](180), // Medium-high heat
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4aID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &chickenBreast.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "sous vide cooked chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4aID,
				ValidIngredientPreparationID:     &panSearOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilMilliliterVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                milliliterMeasurement.ID,
				Name:                             "vegetable oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 15, // Enough to coat the pan
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4aID,
				ValidPreparationInstrumentID: &panSearPaperTowelsVPI.ID,
				InstrumentID:                 &paperTowels.ID,
				Name:                         "paper towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4aID,
				ValidPreparationInstrumentID: &panSearSpatulaVPI.ID,
				InstrumentID:                 &spatula.ID,
				Name:                         "flexible metal spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4aID,
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
				BelongsToRecipeStep:      step4aID,
				ValidPreparationVesselID: &panSearSkilletVPV.ID,
				VesselID:                 &castIronSkillet.ID,
				Name:                     "cast iron or stainless steel skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4aID,
				Name:                "pan-seared sous vide chicken breasts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step4aCompletionConditionID,
				BelongsToRecipeStep: step4aID,
				IngredientStateID:   atTemperatureState.ID,
				Notes:               "Skin should be golden brown and crisp",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step4aCompletionConditionID,
						RecipeStepIngredient:                   "", // Will be set below
					},
				},
				Optional: false,
			},
		},
	}

	// Fix step 4a completion condition - need to reference the actual ingredient ID
	step4aChickenIngredientID := step4a.Ingredients[0].ID
	step4a.CompletionConditions[0].Ingredients[0].RecipeStepIngredient = step4aChickenIngredientID

	// Step 4b: Finish on grill (optional)
	step4bID := identifiers.New()
	step4b := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4bID,
		BelongsToRecipe: recipeID,
		PreparationID:   grillPrep.ID,
		Index:           5,
		Optional:        true,
		Notes:           "Light one half chimney full of charcoal. When all the charcoal is lit and covered with gray ash, pour out and arrange coals on one side of charcoal grate. Set cooking grate in place, cover grill, and allow to preheat for 5 minutes. Alternatively, set half the burners on a gas grill to medium-high heat setting, cover, and preheat for 10 minutes. Clean and oil the grilling grate. Remove chicken from water bath and bag. Discard herbs, if using. Carefully pat chicken dry with paper towels. Gently lay chicken on hot side of grill, skin side down. Cook until nicely marked, about 2 minutes. Rotate 90 degrees and cook until chicken is marked again and skin is crisp, about 1 1/2 minutes longer. Remove from grill and let rest until cool enough to handle, about 2 minutes. Remove bones with your fingertips.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](210), // 3.5 minutes total
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](190), // Medium-high heat
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4bID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &chickenBreast.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "sous vide cooked chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4bID,
				ValidPreparationInstrumentID: &grillPaperTowelsVPI.ID,
				InstrumentID:                 &paperTowels.ID,
				Name:                         "paper towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4bID,
				ValidPreparationInstrumentID: &grillTongsVPI.ID,
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
				BelongsToRecipeStep:      step4bID,
				ValidPreparationVesselID: &grillGrillVPV.ID,
				VesselID:                 &grillVessel.ID,
				Name:                     "grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4bID,
				Name:                "grilled sous vide chicken breasts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 6: Rest and serve
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           6,
		Notes:           "Slice chicken and serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &chickenBreast.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "pan-seared sous vide chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
				Optional: true, // From step 4a if pan finishing chosen
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &chickenBreast.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "grilled sous vide chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
				Optional: true, // From step 4b if grill finishing chosen
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
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
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &restPlateVPV.ID,
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
				BelongsToRecipeStep: step6ID,
				Name:                "sliced sous vide chicken breasts",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		{
			ID:                  recipeID,
			CreatedByUser:       userID,
			Name:                "Sous Vide Chicken Breast",
			Slug:                "sous-vide-chicken-breast",
			Source:              "https://www.seriouseats.com/the-food-lab-complete-guide-to-sous-vide-chicken-breast",
			Description:         "",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 2,
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps:             []*mealplanning.RecipeStepDatabaseCreationInput{step0, step1, step2, step3, step4a, step4b, step6},
		},
	}
}
