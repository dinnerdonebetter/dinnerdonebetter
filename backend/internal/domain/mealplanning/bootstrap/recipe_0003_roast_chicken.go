package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func PerfectRoastChickenRecipe(userID string, enums *Enumerations) *mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	mixPrep := enums.Preparations["mix"]
	seasonPrep := enums.Preparations["season"]
	trussPrep := enums.Preparations["truss"]
	dryBrinePrep := enums.Preparations["dry-brine"]
	heatPrep := enums.Preparations["heat"]
	rubPrep := enums.Preparations["rub"]
	panSearPrep := enums.Preparations["pan-sear"]
	roastPrep := enums.Preparations["roast"]
	restPrep := enums.Preparations["rest"]

	// Get ingredients
	wholeChicken := enums.Ingredients["whole chicken"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	bakingPowder := enums.Ingredients["baking powder"]
	vegetableOil := enums.Ingredients["vegetable oil"]

	// Get measurement units
	unitMeasurement := enums.MeasurementUnits["unit"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	milliliterMeasurement := enums.MeasurementUnits["milliliter"]

	// Get instruments
	butchersTwine := enums.Instruments["butcher's twine"]
	tongs := enums.Instruments["tongs"]
	thermometer := enums.Instruments["instant-read thermometer"]
	bareHands := enums.Instruments["bare hands"]

	// Get vessels
	smallBowl := enums.Vessels["small bowl"]
	wireRack := enums.Vessels["wire rack"]
	bakingSheet := enums.Vessels["baking sheet"]
	stainlessSteelSkillet := enums.Vessels["stainless steel skillet"]
	carvingBoard := enums.Vessels["carving board"]

	// Get ingredient states for completion conditions
	brownedState := enums.IngredientStates["browned"]
	atTemperatureState := enums.IngredientStates["at temperature"]

	// Get bridge table entries
	// Mix preparation bridges
	mixSaltVIP := enums.IngredientPreparations[mixPrep.ID][salt.ID]
	mixPepperVIP := enums.IngredientPreparations[mixPrep.ID][blackPepper.ID]
	mixBakingPowderVIP := enums.IngredientPreparations[mixPrep.ID][bakingPowder.ID]
	saltTablespoonVIMU := enums.IngredientMeasurementUnits[salt.ID][tablespoonMeasurement.ID]
	pepperTeaspoonVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]
	bakingPowderTeaspoonVIMU := enums.IngredientMeasurementUnits[bakingPowder.ID][teaspoonMeasurement.ID]
	mixSmallBowlVPV := enums.PreparationVessels[mixPrep.ID][smallBowl.ID]

	// Season preparation bridges
	seasonChickenVIP := enums.IngredientPreparations[seasonPrep.ID][wholeChicken.ID]
	chickenUnitVIMU := enums.IngredientMeasurementUnits[wholeChicken.ID][unitMeasurement.ID]
	seasonBareHandsVPI := enums.PreparationInstruments[seasonPrep.ID][bareHands.ID]

	// Truss preparation bridges
	_ = enums.IngredientPreparations[trussPrep.ID][wholeChicken.ID] // validated but not used
	trussTwineVPI := enums.PreparationInstruments[trussPrep.ID][butchersTwine.ID]

	// Dry-brine preparation bridges
	_ = enums.IngredientPreparations[dryBrinePrep.ID][wholeChicken.ID] // validated but not used
	dryBrineWireRackVPV := enums.PreparationVessels[dryBrinePrep.ID][wireRack.ID]
	dryBrineBakingSheetVPV := enums.PreparationVessels[dryBrinePrep.ID][bakingSheet.ID]

	// Heat preparation bridges
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][vegetableOil.ID]
	oilTablespoonVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][tablespoonMeasurement.ID]
	heatSkilletVPV := enums.PreparationVessels[heatPrep.ID][stainlessSteelSkillet.ID]

	// Rub preparation bridges
	_ = enums.IngredientPreparations[rubPrep.ID][wholeChicken.ID] // validated but not used
	rubOilVIP := enums.IngredientPreparations[rubPrep.ID][vegetableOil.ID]
	rubBareHandsVPI := enums.PreparationInstruments[rubPrep.ID][bareHands.ID]

	// Pan-sear preparation bridges
	_ = enums.IngredientPreparations[panSearPrep.ID][wholeChicken.ID] // validated but not used
	panSearTongsVPI := enums.PreparationInstruments[panSearPrep.ID][tongs.ID]
	panSearSkilletVPV := enums.PreparationVessels[panSearPrep.ID][stainlessSteelSkillet.ID]

	// Roast preparation bridges
	_ = enums.IngredientPreparations[roastPrep.ID][wholeChicken.ID] // validated but not used
	roastThermometerVPI := enums.PreparationInstruments[roastPrep.ID][thermometer.ID]
	roastSkilletVPV := enums.PreparationVessels[roastPrep.ID][stainlessSteelSkillet.ID]

	// Rest preparation bridges
	_ = enums.IngredientPreparations[restPrep.ID][wholeChicken.ID] // validated but not used
	restCarvingBoardVPV := enums.PreparationVessels[restPrep.ID][carvingBoard.ID]

	// Step 0: Mix the seasoning
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   mixPrep.ID,
		Index:           0,
		Notes:           "In a small bowl, thoroughly mix the salt with black pepper and baking powder (if using).",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &mixSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTablespoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "Diamond Crystal kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1, // 1 tablespoon = 9g
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &mixPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperTeaspoonVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &mixBakingPowderVIP.ID,
				ValidIngredientMeasurementUnitID: &bakingPowderTeaspoonVIMU.ID,
				IngredientID:                     &bakingPowder.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "baking powder",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &mixSmallBowlVPV.ID,
				VesselID:                 &smallBowl.ID,
				Name:                     "small bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "seasoning mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &gramMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](12), // approximately 12g total
				},
			},
		},
	}

	// Step 1: Season the chicken
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           1,
		Notes:           "Season chicken all over, inside and out, with salt mixture (or just plain salt if not using pepper and baking powder).",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &seasonChickenVIP.ID,
				ValidIngredientMeasurementUnitID: &chickenUnitVIMU.ID,
				IngredientID:                     &wholeChicken.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "large whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                             identifiers.New(),
				BelongsToRecipeStep:            step1ID,
				ProductOfRecipeStepIndex:       pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				MeasurementUnitID:              gramMeasurement.ID,
				Name:                           "seasoning mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 12,
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
				ID:                identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "seasoned whole chicken",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Truss the chicken
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   trussPrep.ID,
		Index:           2,
		Notes:           "Set chicken, breast side up, on work surface and tuck wings behind back. Using butcher's twine, run the center of the twine under the tip of the tail end and truss chicken by tying drumsticks together at their bony ends, securing the legs and the tip of the tail together in a bundle. Criss-cross the twine and pass along the crevasse where the legs meet the breast; pass twine over wings to hold them into place, then tie securely around the stump of the neck.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                             identifiers.New(),
				BelongsToRecipeStep:            step2ID,
				ProductOfRecipeStepIndex:       pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                   &wholeChicken.ID,
				MeasurementUnitID:              unitMeasurement.ID,
				Name:                           "seasoned whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                            identifiers.New(),
				BelongsToRecipeStep:           step2ID,
				ValidPreparationInstrumentID: &trussTwineVPI.ID,
				InstrumentID:                 &butchersTwine.ID,
				Name:                         "butcher's twine",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "trussed seasoned whole chicken",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Dry-brine (refrigerate)
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dryBrinePrep.ID,
		Index:           3,
		Notes:           "Place chicken, back side down, on a wire rack set in a rimmed baking sheet and refrigerate, uncovered, at least 1 hour and up to 2 days.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](3600),   // 1 hour minimum
			Max: pointer.To[uint32](172800), // 2 days maximum
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                             identifiers.New(),
				BelongsToRecipeStep:            step3ID,
				ProductOfRecipeStepIndex:       pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                   &wholeChicken.ID,
				MeasurementUnitID:              unitMeasurement.ID,
				Name:                           "trussed seasoned whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &dryBrineWireRackVPV.ID,
				VesselID:                 &wireRack.ID,
				Name:                     "wire rack",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &dryBrineBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dry-brined whole chicken",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Heat oil in skillet (and preheat oven)
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           4,
		Notes:           "Adjust oven rack to middle position and preheat oven to 425°F (220°C). In a 10- or 12-inch stainless steel skillet, heat oil over medium-high heat until shimmering.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](220), // 425°F
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilTablespoonVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "neutral-flavored oil (vegetable or canola)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &heatSkilletVPV.ID,
				VesselID:                 &stainlessSteelSkillet.ID,
				Name:                     "stainless steel skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "heated oil in skillet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &milliliterMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](15),
				},
			},
		},
	}

	// Step 5: Rub chicken with oil
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   rubPrep.ID,
		Index:           5,
		Notes:           "Rub chicken lightly with oil.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                             identifiers.New(),
				BelongsToRecipeStep:            step5ID,
				ProductOfRecipeStepIndex:       pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                   &wholeChicken.ID,
				MeasurementUnitID:              unitMeasurement.ID,
				Name:                           "dry-brined whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &rubOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilTablespoonVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "neutral-flavored oil for rubbing",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &rubBareHandsVPI.ID,
				InstrumentID:                 &bareHands.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "oiled whole chicken",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Brown chicken legs
	step6ID := identifiers.New()
	step6ChickenIngredientID := identifiers.New()
	step6CompletionConditionID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   panSearPrep.ID,
		Index:           6,
		Notes:           "Set chicken on its side in the skillet so that the full thigh and drumstick are in contact with the pan; the wing will also be touching, but the breast should have little to no contact with the skillet. Cook until leg is well browned, 8 to 10 minutes, then flip bird so other leg is touching pan and repeat; lower heat at any point if chicken skin begins to burn.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](960),  // 16 minutes minimum (8 min per side)
			Max: pointer.To[uint32](1200), // 20 minutes maximum (10 min per side)
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                             step6ChickenIngredientID,
				BelongsToRecipeStep:            step6ID,
				ProductOfRecipeStepIndex:       pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                   &wholeChicken.ID,
				MeasurementUnitID:              unitMeasurement.ID,
				Name:                           "oiled whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                             identifiers.New(),
				BelongsToRecipeStep:            step6ID,
				ProductOfRecipeStepIndex:       pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                   &vegetableOil.ID,
				MeasurementUnitID:              milliliterMeasurement.ID,
				Name:                           "heated oil in skillet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 15,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                            identifiers.New(),
				BelongsToRecipeStep:           step6ID,
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
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &panSearSkilletVPV.ID,
				VesselID:                 &stainlessSteelSkillet.ID,
				Name:                     "stainless steel skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "leg-browned whole chicken",
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
				ID:                  step6CompletionConditionID,
				BelongsToRecipeStep: step6ID,
				IngredientStateID:   brownedState.ID,
				Notes:               "Both chicken legs should be well browned",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step6CompletionConditionID,
						RecipeStepIngredient:                   step6ChickenIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 7: Roast the chicken
	step7ID := identifiers.New()
	step7ChickenIngredientID := identifiers.New()
	step7CompletionConditionID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   roastPrep.ID,
		Index:           7,
		Notes:           "Using hands and spatula if needed, rotate chicken so it is breast side up in the skillet and transfer to oven. Roast until breast registers 150°F (65°C) in the center of its thickest part and thighs register 165°F (75°C) near (but not touching) the bone, about 40 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](2400), // 40 minutes
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](65), // 150°F for breast
			Max: pointer.To[float32](75), // 165°F for thighs
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                             step7ChickenIngredientID,
				BelongsToRecipeStep:            step7ID,
				ProductOfRecipeStepIndex:       pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                   &wholeChicken.ID,
				MeasurementUnitID:              unitMeasurement.ID,
				Name:                           "leg-browned whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                            identifiers.New(),
				BelongsToRecipeStep:           step7ID,
				ValidPreparationInstrumentID: &roastThermometerVPI.ID,
				InstrumentID:                 &thermometer.ID,
				Name:                         "instant-read thermometer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step7ID,
				ValidPreparationVesselID: &roastSkilletVPV.ID,
				VesselID:                 &stainlessSteelSkillet.ID,
				Name:                     "stainless steel skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "roasted whole chicken",
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
				ID:                  step7CompletionConditionID,
				BelongsToRecipeStep: step7ID,
				IngredientStateID:   atTemperatureState.ID,
				Notes:               "Breast should register 150°F (65°C) and thighs 165°F (75°C)",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step7CompletionConditionID,
						RecipeStepIngredient:                   step7ChickenIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 8: Rest the chicken
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           8,
		Notes:           "Remove from oven and transfer chicken to a carving board. Let rest 10 to 20 minutes, then carve and serve.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600),  // 10 minutes
			Max: pointer.To[uint32](1200), // 20 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                             identifiers.New(),
				BelongsToRecipeStep:            step8ID,
				ProductOfRecipeStepIndex:       pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                   &wholeChicken.ID,
				MeasurementUnitID:              unitMeasurement.ID,
				Name:                           "roasted whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: &restCarvingBoardVPV.ID,
				VesselID:                 &carvingBoard.ID,
				Name:                     "carving board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "rested roast chicken",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	return &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Perfect Roast Chicken",
		Slug:                "perfect-roast-chicken",
		Source:              "https://www.seriouseats.com/perfect-roast-chicken-recipe-8384377",
		Description:         "A dry-brine fully seasons the chicken and allows the skin to dehydrate, improving browning and crisping during cooking. Measuring the internal temperature of the chicken to determine doneness leads to more reliable and superior results than going by time. Rubbing the chicken skin with oil before roasting, instead of basting with watery drippings, ensures even browning and a crisp skin.",
		YieldsComponentType: mealplanning.MealComponentTypesMain,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps:             []*mealplanning.RecipeStepDatabaseCreationInput{step0, step1, step2, step3, step4, step5, step6, step7, step8},
	}
}

