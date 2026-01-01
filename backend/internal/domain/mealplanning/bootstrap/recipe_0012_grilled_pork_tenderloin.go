package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// GrilledPorkTenderloinRecipe creates the Grilled Pork Tenderloin recipe.
// Source: https://www.seriouseats.com/grilled-pork-tenderloin-recipe-7505776
func GrilledPorkTenderloinRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	trimPrep := enums.Preparations["trim"]
	seasonPrep := enums.Preparations["season"]
	transferPrep := enums.Preparations["transfer"]
	preheatPrep := enums.Preparations["preheat"]
	cleanPrep := enums.Preparations["clean"]
	oilPrep := enums.Preparations["oil"]
	grillPrep := enums.Preparations["grill"]
	restPrep := enums.Preparations["rest"]
	carvePrep := enums.Preparations["carve"]

	// Get ingredients
	porkTenderloin := enums.Ingredients["pork tenderloin"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	vegetableOil := enums.Ingredients["vegetable oil"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]

	// Get instruments
	chefsKnife := enums.Instruments["knife"]
	bareHands := enums.Instruments["bare hands"]
	grillBrush := enums.Instruments["grill brush"]
	brush := enums.Instruments["brush"]
	tongs := enums.Instruments["tongs"]
	thermometer := enums.Instruments["instant-read thermometer"]
	carvingKnife := enums.Instruments["carving knife"]

	// Get vessels
	wireRack := enums.Vessels["wire rack"]
	bakingSheet := enums.Vessels["baking sheet"]
	grill := enums.Vessels["grill"]
	grillingGrate := enums.Vessels["grilling grate"]
	carvingBoard := enums.Vessels["carving board"]

	// Get bridge table entries
	// Trim
	trimPorkVIP := enums.IngredientPreparations[trimPrep.ID][porkTenderloin.ID]
	trimChefsKnifeVPI := enums.PreparationInstruments[trimPrep.ID][chefsKnife.ID]

	// Season
	seasonPorkVIP := enums.IngredientPreparations[seasonPrep.ID][porkTenderloin.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonBlackPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	seasonBareHandsVPI := enums.PreparationInstruments[seasonPrep.ID][bareHands.ID]

	// Transfer
	transferPorkVIP := enums.IngredientPreparations[transferPrep.ID][porkTenderloin.ID]
	transferWireRackVPV := enums.PreparationVessels[transferPrep.ID][wireRack.ID]
	transferBakingSheetVPV := enums.PreparationVessels[transferPrep.ID][bakingSheet.ID]
	transferCarvingBoardVPV := enums.PreparationVessels[transferPrep.ID][carvingBoard.ID]

	// Preheat
	preheatGrillVPV := enums.PreparationVessels[preheatPrep.ID][grill.ID]

	// Clean
	cleanGrillingGrateVPV := enums.PreparationVessels[cleanPrep.ID][grillingGrate.ID]
	cleanGrillBrushVPI := enums.PreparationInstruments[cleanPrep.ID][grillBrush.ID]

	// Oil
	oilVegetableOilVIP := enums.IngredientPreparations[oilPrep.ID][vegetableOil.ID]
	oilGrillingGrateVPV := enums.PreparationVessels[oilPrep.ID][grillingGrate.ID]
	oilBrushVPI := enums.PreparationInstruments[oilPrep.ID][brush.ID]

	// Grill
	grillPorkVIP := enums.IngredientPreparations[grillPrep.ID][porkTenderloin.ID]
	grillGrillVPV := enums.PreparationVessels[grillPrep.ID][grill.ID]
	grillGrillingGrateVPV := enums.PreparationVessels[grillPrep.ID][grillingGrate.ID]
	grillTongsVPI := enums.PreparationInstruments[grillPrep.ID][tongs.ID]
	grillThermometerVPI := enums.PreparationInstruments[grillPrep.ID][thermometer.ID]

	// Rest
	restPorkVIP := enums.IngredientPreparations[restPrep.ID][porkTenderloin.ID]
	restCarvingBoardVPV := enums.PreparationVessels[restPrep.ID][carvingBoard.ID]

	// Carve
	carvePorkVIP := enums.IngredientPreparations[carvePrep.ID][porkTenderloin.ID]
	carveCarvingBoardVPV := enums.PreparationVessels[carvePrep.ID][carvingBoard.ID]
	carveCarvingKnifeVPI := enums.PreparationInstruments[carvePrep.ID][carvingKnife.ID]

	// Measurement unit bridges
	porkTenderloinPoundVIMU := enums.IngredientMeasurementUnits[porkTenderloin.ID][poundMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	blackPepperTeaspoonVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]
	vegetableOilTablespoonVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][tablespoonMeasurement.ID]

	// Step 0: Trim silverskin from pork tenderloins (optional prep step)
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   trimPrep.ID,
		Index:           0,
		Notes:           "Trim silverskin from pork tenderloins. This step is optional and can be done ahead of time.",
		Optional:        true,
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimPorkVIP.ID,
				ValidIngredientMeasurementUnitID: &porkTenderloinPoundVIMU.ID,
				IngredientID:                     &porkTenderloin.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
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
				Name:                "trimmed pork tenderloins",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 1: Sprinkle pork tenderloins all over with salt
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           1,
		Notes:           "Sprinkle pork tenderloins all over with salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonPorkVIP.ID,
				IngredientID:                    &porkTenderloin.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "trimmed pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
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
				Name:                "salted pork tenderloins",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 2: Place on a wire rack set over a rimmed baking sheet
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           2,
		Notes:           "Place on a wire rack set over a rimmed baking sheet.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferPorkVIP.ID,
				IngredientID:                    &porkTenderloin.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "salted pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &transferWireRackVPV.ID,
				VesselID:                 &wireRack.ID,
				Name:                     "wire rack",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
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
				BelongsToRecipeStep: step2ID,
				Name:                "pork tenderloins on wire rack",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: When ready to cook, light a chimney full of charcoal. When all charcoal is lit and covered with gray ash, pour out and arrange coals on one side of coal grate and set grilling grate in place. Alternatively, set half the burners of a gas grill to high heat. Cover grill and allow to preheat for 5 minutes.
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   preheatPrep.ID,
		Index:           3,
		Notes:           "When ready to cook, light a chimney full of charcoal. When all charcoal is lit and covered with gray ash, pour out and arrange coals on one side of coal grate and set grilling grate in place. Alternatively, set half the burners of a gas grill to high heat. Cover grill and allow to preheat for 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &preheatGrillVPV.ID,
				VesselID:                 &grill.ID,
				Name:                     "grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "preheated grill",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Clean grilling grate
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cleanPrep.ID,
		Index:           4,
		Notes:           "Clean grilling grate.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &cleanGrillingGrateVPV.ID,
				VesselID:                 &grillingGrate.ID,
				Name:                     "grilling grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &cleanGrillBrushVPI.ID,
				InstrumentID:                 &grillBrush.ID,
				Name:                         "grill brush",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "cleaned grilling grate",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Oil grilling grate
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   oilPrep.ID,
		Index:           5,
		Notes:           "Oil grilling grate.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &oilVegetableOilVIP.ID,
				ValidIngredientMeasurementUnitID: &vegetableOilTablespoonVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "vegetable oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &oilGrillingGrateVPV.ID,
				VesselID:                        &grillingGrate.ID,
				Name:                            "cleaned grilling grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &oilBrushVPI.ID,
				InstrumentID:                 &brush.ID,
				Name:                         "brush",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "oiled grilling grate",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Season tenderloins all over with pepper
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           6,
		Notes:           "Season tenderloins all over with pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonPorkVIP.ID,
				IngredientID:                    &porkTenderloin.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "pork tenderloins on wire rack",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &seasonBlackPepperVIP.ID,
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
				BelongsToRecipeStep:          step6ID,
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
				BelongsToRecipeStep: step6ID,
				Name:                "seasoned pork tenderloins",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 7: Set over hot side of grill and cook, turning often, until well browned on all sides, about 15 minutes
	step7ID := identifiers.New()
	step7PorkIngredientID := identifiers.New()
	step7CompletionConditionID := identifiers.New()
	brownedState := enums.IngredientStates["browned"]
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   grillPrep.ID,
		Index:           7,
		Notes:           "Set over hot side of grill and cook, turning often, until well browned on all sides, about 15 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step7PorkIngredientID,
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &grillPorkVIP.ID,
				IngredientID:                    &porkTenderloin.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "seasoned pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step7ID,
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
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &grillGrillVPV.ID,
				VesselID:                        &grill.ID,
				Name:                            "preheated grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &grillGrillingGrateVPV.ID,
				VesselID:                        &grillingGrate.ID,
				Name:                            "oiled grilling grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "browned pork tenderloins",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step7CompletionConditionID,
				BelongsToRecipeStep: step7ID,
				IngredientStateID:   brownedState.ID,
				Notes:               "Pork should be well browned on all sides",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step7CompletionConditionID,
						RecipeStepIngredient:                   step7PorkIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 8: Move tenderloins to cooler side of grill and continue to cook, turning often, until an instant-read thermometer inserted in the center registers 120 to 130°F (49 to 54°C) for medium-rare or 130 to 140°F (54 to 60°C) for medium
	step8ID := identifiers.New()
	step8PorkIngredientID := identifiers.New()
	step8CompletionConditionID := identifiers.New()
	atTemperatureState := enums.IngredientStates["at temperature"]
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   grillPrep.ID,
		Index:           8,
		Notes:           "Move tenderloins to cooler side of grill and continue to cook, turning often, until an instant-read thermometer inserted in the center registers 120 to 130°F (49 to 54°C) for medium-rare or 130 to 140°F (54 to 60°C) for medium.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step8PorkIngredientID,
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &grillPorkVIP.ID,
				IngredientID:                    &porkTenderloin.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "browned pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step8ID,
				ValidPreparationInstrumentID: &grillTongsVPI.ID,
				InstrumentID:                 &tongs.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step8ID,
				ValidPreparationInstrumentID: &grillThermometerVPI.ID,
				InstrumentID:                 &thermometer.ID,
				Name:                         "instant-read thermometer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &grillGrillVPV.ID,
				VesselID:                        &grill.ID,
				Name:                            "preheated grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "cooked pork tenderloins",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step8CompletionConditionID,
				BelongsToRecipeStep: step8ID,
				IngredientStateID:   atTemperatureState.ID,
				Notes:               "Internal temperature should register 120 to 130°F (49 to 54°C) for medium-rare or 130 to 140°F (54 to 60°C) for medium",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step8CompletionConditionID,
						RecipeStepIngredient:                   step8PorkIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 9: Transfer pork to a carving board
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           9,
		Notes:           "Transfer pork to a carving board.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferPorkVIP.ID,
				IngredientID:                    &porkTenderloin.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "cooked pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step9ID,
				ValidPreparationVesselID: &transferCarvingBoardVPV.ID,
				VesselID:                 &carvingBoard.ID,
				Name:                     "carving board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "pork tenderloins on carving board",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Let rest for 10 minutes
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           10,
		Notes:           "Let rest for 10 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &restPorkVIP.ID,
				IngredientID:                    &porkTenderloin.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "pork tenderloins on carving board",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &restCarvingBoardVPV.ID,
				VesselID:                        &carvingBoard.ID,
				Name:                            "pork tenderloins on carving board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "rested pork tenderloins",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 11: Carve pork tenderloins and serve as desired
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   carvePrep.ID,
		Index:           11,
		Notes:           "Carve pork tenderloins and serve as desired.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &carvePorkVIP.ID,
				IngredientID:                    &porkTenderloin.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "rested pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step11ID,
				ValidPreparationInstrumentID: &carveCarvingKnifeVPI.ID,
				InstrumentID:                 &carvingKnife.ID,
				Name:                         "carving knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &carveCarvingBoardVPV.ID,
				VesselID:                        &carvingBoard.ID,
				Name:                            "pork tenderloins on carving board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "grilled pork tenderloin",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Create prep task for dry-brining pork ahead of time
	prepTask1ID := identifiers.New()
	prepTask1 := &mealplanning.RecipePrepTaskDatabaseCreationInput{
		ID:                          prepTask1ID,
		BelongsToRecipe:             recipeID,
		Name:                        "Dry-brine pork tenderloins",
		Description:                 "The pork tenderloins can be dry-brined at least 45 minutes and up to 24 hours in advance. Store in the refrigerator uncovered on a wire rack set in a rimmed baking sheet.",
		Notes:                       "Dry-brining improves flavor, moisture retention, and subsequent juiciness, and ensures deeper browning on the exterior.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the salted pork tenderloins on a wire rack set in a rimmed baking sheet in the refrigerator, uncovered, for at least 45 minutes and up to 24 hours.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeWireRack,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 2700,                      // 45 minutes
			Max: pointer.To[uint32](86400), // 24 hours
		},
		TaskSteps: []*mealplanning.RecipePrepTaskStepDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: step1ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step2ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: true},
		},
	}

	grilledPorkTenderloinRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Grilled Pork Tenderloin",
		Slug:                "grilled-pork-tenderloin",
		Source:              "https://www.seriouseats.com/grilled-pork-tenderloin-recipe-7505776",
		Description:         "Salting in advance and grilling over high-heat are the key steps to great grilled pork tenderloin. The meat is grilled on the hot side until well-browned, then moved to the cooler side to finish cooking to the desired temperature.",
		YieldsComponentType: mealplanning.MealComponentTypesMain,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
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
		grilledPorkTenderloinRecipe,
	}
}
