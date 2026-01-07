package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// GrilledPorkTenderloinRecipe creates the Grilled Pork Tenderloin recipe.
// Source: https://www.seriouseats.com/grilled-pork-tenderloin-recipe-7505776
func GrilledPorkTenderloinRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
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
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: trimPrep.ID,
		Index:         0,
		Notes:         "Trim silverskin from pork tenderloins. This step is optional and can be done ahead of time.",
		Optional:      true,
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &trimPorkVIP.ID,
				ValidIngredientMeasurementUnitID: &porkTenderloinPoundVIMU.ID,
				Name:                             "pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
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
				Name:              "trimmed pork tenderloins",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 1: Sprinkle pork tenderloins all over with salt
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         1,
		Notes:         "Sprinkle pork tenderloins all over with salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonPorkVIP.ID,
				Name:                            "trimmed pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
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
				Name:              "salted pork tenderloins",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 2: Place on a wire rack set over a rimmed baking sheet
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: transferPrep.ID,
		Index:         2,
		Notes:         "Place on a wire rack set over a rimmed baking sheet.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferPorkVIP.ID,
				Name:                            "salted pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &transferWireRackVPV.ID,
				Name:                     "wire rack",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationVesselID: &transferBakingSheetVPV.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pork tenderloins on wire rack",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:              "pork tenderloins",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             1,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 3: When ready to cook, light a chimney full of charcoal. When all charcoal is lit and covered with gray ash, pour out and arrange coals on one side of coal grate and set grilling grate in place. Alternatively, set half the burners of a gas grill to high heat. Cover grill and allow to preheat for 5 minutes.
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: preheatPrep.ID,
		Index:         3,
		Notes:         "When ready to cook, light a chimney full of charcoal. When all charcoal is lit and covered with gray ash, pour out and arrange coals on one side of coal grate and set grilling grate in place. Alternatively, set half the burners of a gas grill to high heat. Cover grill and allow to preheat for 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &preheatGrillVPV.ID,
				Name:                     "grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "preheated grill",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Clean grilling grate
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: cleanPrep.ID,
		Index:         4,
		Notes:         "Clean grilling grate.",
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &cleanGrillingGrateVPV.ID,
				Name:                     "grilling grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &cleanGrillBrushVPI.ID,
				Name:                         "grill brush",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "cleaned grilling grate",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Oil grilling grate
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: oilPrep.ID,
		Index:         5,
		Notes:         "Oil grilling grate.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &oilVegetableOilVIP.ID,
				ValidIngredientMeasurementUnitID: &vegetableOilTablespoonVIMU.ID,
				Name:                             "vegetable oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &oilGrillingGrateVPV.ID,
				Name:                            "cleaned grilling grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &oilBrushVPI.ID,
				Name:                         "brush",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "oiled grilling grate",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Season tenderloins all over with pepper
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         6,
		Notes:         "Season tenderloins all over with pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &seasonPorkVIP.ID,
				Name:                            "pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonBlackPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &blackPepperTeaspoonVIMU.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
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
				Name:              "seasoned pork tenderloins",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 7: Set over hot side of grill and cook, turning often, until well browned on all sides, about 15 minutes
	brownedState := enums.IngredientStates["browned"]
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: grillPrep.ID,
		Index:         7,
		Notes:         "Set over hot side of grill and cook, turning often, until well browned on all sides, about 15 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &grillPorkVIP.ID,
				Name:                            "seasoned pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &grillTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &grillGrillVPV.ID,
				Name:                            "preheated grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &grillGrillingGrateVPV.ID,
				Name:                            "oiled grilling grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "browned pork tenderloins",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: brownedState.ID,
				Notes:             "Pork should be well browned on all sides",
				Ingredients:       []uint64{0}, // Index of pork ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 8: Move tenderloins to cooler side of grill and continue to cook, turning often, until an instant-read thermometer inserted in the center registers 120 to 130°F (49 to 54°C) for medium-rare or 130 to 140°F (54 to 60°C) for medium
	atTemperatureState := enums.IngredientStates["at temperature"]
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: grillPrep.ID,
		Index:         8,
		Notes:         "Move tenderloins to cooler side of grill and continue to cook, turning often, until an instant-read thermometer inserted in the center registers 120 to 130°F (49 to 54°C) for medium-rare or 130 to 140°F (54 to 60°C) for medium.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &grillPorkVIP.ID,
				Name:                            "browned pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &grillTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationInstrumentID: &grillThermometerVPI.ID,
				Name:                         "instant-read thermometer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &grillGrillVPV.ID,
				Name:                            "preheated grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cooked pork tenderloins",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: atTemperatureState.ID,
				Notes:             "Internal temperature should register 120 to 130°F (49 to 54°C) for medium-rare or 130 to 140°F (54 to 60°C) for medium",
				Ingredients:       []uint64{0}, // Index of pork ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 9: Transfer pork to a carving board
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: transferPrep.ID,
		Index:         9,
		Notes:         "Transfer pork to a carving board.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferPorkVIP.ID,
				Name:                            "cooked pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &transferCarvingBoardVPV.ID,
				Name:                     "carving board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pork tenderloins on carving board",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:              "cooked pork tenderloins",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             1,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 10: Let rest for 10 minutes
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: restPrep.ID,
		Index:         10,
		Notes:         "Let rest for 10 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &restPorkVIP.ID,
				Name:                            "cooked pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &restCarvingBoardVPV.ID,
				Name:                            "pork tenderloins on carving board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "rested pork tenderloins",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 11: Carve pork tenderloins and serve as desired
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: carvePrep.ID,
		Index:         11,
		Notes:         "Carve pork tenderloins and serve as desired.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &carvePorkVIP.ID,
				Name:                            "rested pork tenderloins",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &carveCarvingKnifeVPI.ID,
				Name:                         "carving knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &carveCarvingBoardVPV.ID,
				Name:                            "pork tenderloins on carving board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "grilled pork tenderloin",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Create prep task for dry-brining pork ahead of time
	prepTask1 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
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
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 1, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 2, SatisfiesRecipeStep: true},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
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
			Steps:             []*mealplanning.RecipeStepCreationRequestInput{step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11},
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{prepTask1},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
