package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// PanSearedSalmonFilletsRecipe creates the Crispy Pan-Seared Salmon Fillets recipe.
// Source: https://www.seriouseats.com/crispy-pan-seared-salmon-fillets-recipe
func PanSearedSalmonFilletsRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	dryPrep := enums.Preparations["dry"]
	seasonPrep := enums.Preparations["season"]
	heatPrep := enums.Preparations["heat"]
	panSearPrep := enums.Preparations["pan-sear"]
	pressPrep := enums.Preparations["press"]
	flipPrep := enums.Preparations["flip"]
	transferPrep := enums.Preparations["transfer"]
	drainPrep := enums.Preparations["drain"]

	// Get ingredients
	salmonFillet := enums.Ingredients["salmon fillet"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	vegetableOil := enums.Ingredients["vegetable oil"]

	// Get measurement units
	ounceMeasurement := enums.MeasurementUnits["ounce"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]

	// Get instruments
	paperTowels := enums.Instruments["paper towels"]
	bareHands := enums.Instruments["bare hands"]
	fishSpatula := enums.Instruments["fish spatula"]
	thermometer := enums.Instruments["instant-read thermometer"]
	fork := enums.Instruments["fork"]

	// Get vessels
	skillet := enums.Vessels["cast iron skillet"]
	plate := enums.Vessels["large plate"]

	// Get bridge table entries
	// Dry
	drySalmonVIP := enums.IngredientPreparations[dryPrep.ID][salmonFillet.ID]
	dryPaperTowelsVPI := enums.PreparationInstruments[dryPrep.ID][paperTowels.ID]

	// Season
	seasonSalmonVIP := enums.IngredientPreparations[seasonPrep.ID][salmonFillet.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonBlackPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	seasonBareHandsVPI := enums.PreparationInstruments[seasonPrep.ID][bareHands.ID]

	// Heat
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][vegetableOil.ID]
	heatSkilletVPV := enums.PreparationVessels[heatPrep.ID][skillet.ID]

	// Pan-sear
	panSearSalmonVIP := enums.IngredientPreparations[panSearPrep.ID][salmonFillet.ID]
	panSearSkilletVPV := enums.PreparationVessels[panSearPrep.ID][skillet.ID]
	panSearFishSpatulaVPI := enums.PreparationInstruments[panSearPrep.ID][fishSpatula.ID]
	panSearThermometerVPI := enums.PreparationInstruments[panSearPrep.ID][thermometer.ID]

	// Press
	pressSalmonVIP := enums.IngredientPreparations[pressPrep.ID][salmonFillet.ID]
	pressSkilletVPV := enums.PreparationVessels[pressPrep.ID][skillet.ID]
	pressFishSpatulaVPI := enums.PreparationInstruments[pressPrep.ID][fishSpatula.ID]

	// Flip
	flipSalmonVIP := enums.IngredientPreparations[flipPrep.ID][salmonFillet.ID]
	flipSkilletVPV := enums.PreparationVessels[flipPrep.ID][skillet.ID]
	flipFishSpatulaVPI := enums.PreparationInstruments[flipPrep.ID][fishSpatula.ID]
	flipForkVPI := enums.PreparationInstruments[flipPrep.ID][fork.ID]

	// Transfer
	transferSalmonVIP := enums.IngredientPreparations[transferPrep.ID][salmonFillet.ID]
	transferPlateVPV := enums.PreparationVessels[transferPrep.ID][plate.ID]

	// Drain
	drainSalmonVIP := enums.IngredientPreparations[drainPrep.ID][salmonFillet.ID]
	drainPlateVPV := enums.PreparationVessels[drainPrep.ID][plate.ID]

	// Measurement unit bridges
	salmonFilletOunceVIMU := enums.IngredientMeasurementUnits[salmonFillet.ID][ounceMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	blackPepperTeaspoonVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]
	vegetableOilTablespoonVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][tablespoonMeasurement.ID]

	// Step 0: Press salmon fillets between paper towels to dry surfaces thoroughly
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dryPrep.ID,
		Index:           0,
		Notes:           "Press salmon fillets between paper towels to dry surfaces thoroughly.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &drySalmonVIP.ID,
				ValidIngredientMeasurementUnitID: &salmonFilletOunceVIMU.ID,
				IngredientID:                     &salmonFillet.ID,
				MeasurementUnitID:                ounceMeasurement.ID,
				Name:                             "skin-on salmon fillets",
				QuantityNotes:                    "4 fillets, about 6 ounces (170 g) each",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 24, // 4 fillets × 6 ounces
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
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
				BelongsToRecipeStep: step0ID,
				Name:                "dried salmon fillets",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &ounceMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](24),
				},
			},
		},
	}

	// Step 1: Season on all sides with salt and pepper
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           1,
		Notes:           "Season on all sides with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonSalmonVIP.ID,
				IngredientID:                    &salmonFillet.ID,
				MeasurementUnitID:               ounceMeasurement.ID,
				Name:                            "dried salmon fillets",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 24,
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
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
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
				Name:                "seasoned salmon fillets",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &ounceMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](24),
				},
			},
		},
	}

	// Step 2: In a large stainless, cast iron, or carbon steel skillet, heat oil over medium-high heat until shimmering
	step2ID := identifiers.New()
	step2OilIngredientID := identifiers.New()
	step2CompletionConditionID := identifiers.New()
	shimmeringState := enums.IngredientStates["shimmering"]
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           2,
		Notes:           "In a large stainless, cast iron, or carbon steel skillet, heat oil over medium-high heat until shimmering.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step2OilIngredientID,
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &vegetableOilTablespoonVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "vegetable, canola, or light olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &heatSkilletVPV.ID,
				VesselID:                 &skillet.ID,
				Name:                     "large stainless, cast iron, or carbon steel skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "heated skillet with shimmering oil",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step2CompletionConditionID,
				BelongsToRecipeStep: step2ID,
				IngredientStateID:   shimmeringState.ID,
				Notes:               "Oil should shimmer when viewed",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step2CompletionConditionID,
						RecipeStepIngredient:                   step2OilIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 3: Reduce heat to medium-low, then add a salmon fillet, skin side down. Press firmly in place for 10 seconds, using the back of a flexible fish spatula, to prevent the skin from buckling. Add remaining fillets one at a time, pressing each with spatula for 10 seconds, until all fillets are in the pan.
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   pressPrep.ID,
		Index:           3,
		Notes:           "Reduce heat to medium-low, then add a salmon fillet, skin side down. Press firmly in place for 10 seconds, using the back of a flexible fish spatula, to prevent the skin from buckling. Add remaining fillets one at a time, pressing each with spatula for 10 seconds, until all fillets are in the pan.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](40), // 4 fillets × 10 seconds
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &pressSalmonVIP.ID,
				IngredientID:                    &salmonFillet.ID,
				MeasurementUnitID:               ounceMeasurement.ID,
				Name:                            "seasoned salmon fillets",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 24,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step3ID,
				ValidPreparationInstrumentID: &pressFishSpatulaVPI.ID,
				InstrumentID:                 &fishSpatula.ID,
				Name:                         "flexible fish spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &pressSkilletVPV.ID,
				VesselID:                        &skillet.ID,
				Name:                            "heated skillet with shimmering oil",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "salmon fillets pressed in skillet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &ounceMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](24),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "skillet with pressed salmon",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Cook, pressing gently on back of fillets occasionally to ensure good contact with skin, until skin releases easily from pan, about 4 minutes. Continue to cook until salmon registers 110°F (43°C) in the very center for rare, 120°F (49°C) for medium-rare, or 130°F (54°C) for medium, 5 to 7 minutes total.
	step4ID := identifiers.New()
	step4SalmonIngredientID := identifiers.New()
	step4CompletionConditionID := identifiers.New()
	atTemperatureState := enums.IngredientStates["at temperature"]
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   panSearPrep.ID,
		Index:           4,
		Notes:           "Cook, pressing gently on back of fillets occasionally to ensure good contact with skin, until skin releases easily from pan, about 4 minutes. If skin shows resistance when you attempt to lift a corner with spatula, allow it to continue to cook until it lifts easily. Continue to cook until salmon registers 110°F (43°C) in the very center for rare, 120°F (49°C) for medium-rare, or 130°F (54°C) for medium, 5 to 7 minutes total.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
			Max: pointer.To[uint32](420), // 7 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step4SalmonIngredientID,
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &panSearSalmonVIP.ID,
				IngredientID:                    &salmonFillet.ID,
				MeasurementUnitID:               ounceMeasurement.ID,
				Name:                            "salmon fillets pressed in skillet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 24,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &panSearFishSpatulaVPI.ID,
				InstrumentID:                 &fishSpatula.ID,
				Name:                         "flexible fish spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &panSearThermometerVPI.ID,
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
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &panSearSkilletVPV.ID,
				VesselID:                        &skillet.ID,
				Name:                            "skillet with pressed salmon",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "cooked salmon fillets (skin side)",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &ounceMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](24),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "skillet with cooked salmon",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step4CompletionConditionID,
				BelongsToRecipeStep: step4ID,
				IngredientStateID:   atTemperatureState.ID,
				Notes:               "Internal temperature should register 110°F (43°C) for rare, 120°F (49°C) for medium-rare, or 130°F (54°C) for medium",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step4CompletionConditionID,
						RecipeStepIngredient:                   step4SalmonIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 5: Using spatula and a fork, flip salmon fillets
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   flipPrep.ID,
		Index:           5,
		Notes:           "Using spatula and a fork, flip salmon fillets.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &flipSalmonVIP.ID,
				IngredientID:                    &salmonFillet.ID,
				MeasurementUnitID:               ounceMeasurement.ID,
				Name:                            "cooked salmon fillets (skin side)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 24,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &flipFishSpatulaVPI.ID,
				InstrumentID:                 &fishSpatula.ID,
				Name:                         "flexible fish spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &flipForkVPI.ID,
				InstrumentID:                 &fork.ID,
				Name:                         "fork",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &flipSkilletVPV.ID,
				VesselID:                        &skillet.ID,
				Name:                            "skillet with cooked salmon",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "flipped salmon fillets",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &ounceMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](24),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "skillet with flipped salmon",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Cook on second side for 15 seconds
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   panSearPrep.ID,
		Index:           6,
		Notes:           "Cook on second side for 15 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](15),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &panSearSalmonVIP.ID,
				IngredientID:                    &salmonFillet.ID,
				MeasurementUnitID:               ounceMeasurement.ID,
				Name:                            "flipped salmon fillets",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 24,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &panSearSkilletVPV.ID,
				VesselID:                        &skillet.ID,
				Name:                            "skillet with flipped salmon",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "fully cooked salmon fillets",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &ounceMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](24),
				},
			},
		},
	}

	// Step 7: Transfer to a paper towel–lined plate
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           7,
		Notes:           "Transfer to a paper towel–lined plate to drain.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferSalmonVIP.ID,
				IngredientID:                    &salmonFillet.ID,
				MeasurementUnitID:               ounceMeasurement.ID,
				Name:                            "fully cooked salmon fillets",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 24,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step7ID,
				ValidPreparationVesselID: &transferPlateVPV.ID,
				VesselID:                 &plate.ID,
				Name:                     "paper towel–lined plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "salmon fillets on plate",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Drain excess oil
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   drainPrep.ID,
		Index:           8,
		Notes:           "Drain excess oil.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &drainSalmonVIP.ID,
				IngredientID:                    &salmonFillet.ID,
				MeasurementUnitID:               ounceMeasurement.ID,
				Name:                            "salmon fillets on plate",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 24,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &drainPlateVPV.ID,
				VesselID:                        &plate.ID,
				Name:                            "salmon fillets on plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "crispy pan-seared salmon fillets",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &ounceMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](24),
				},
			},
		},
	}

	panSearedSalmonFilletsRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Crispy Pan-Seared Salmon Fillets",
		Slug:                "crispy-pan-seared-salmon-fillets",
		Source:              "https://www.seriouseats.com/crispy-pan-seared-salmon-fillets-recipe",
		Description:         "How to simultaneously achieve extra-crunchy skin and perfectly tender fish. The key is to cook the salmon most of the way through with the skin side down in order to insulate the delicate flesh from the direct heat of the pan.",
		YieldsComponentType: mealplanning.MealComponentTypesMain,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		panSearedSalmonFilletsRecipe,
	}
}
