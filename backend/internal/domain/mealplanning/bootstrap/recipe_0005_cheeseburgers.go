package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ClassicSmashBurgersRecipe creates the Classic Smashed Burgers recipe from Serious Eats.
func ClassicSmashBurgersRecipe(userID string, enums *Enumerations) *mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	heatPrep := enums.Preparations["heat"]
	dividePrep := enums.Preparations["divide"]
	formPrep := enums.Preparations["form"]
	seasonPrep := enums.Preparations["season"]
	toastPrep := enums.Preparations["toast"]
	smashPrep := enums.Preparations["smash"]
	panSearPrep := enums.Preparations["pan-sear"]
	flipPrep := enums.Preparations["flip"]
	topPrep := enums.Preparations["top"]
	assemblePrep := enums.Preparations["assemble"]

	// Get ingredients
	groundBeef := enums.Ingredients["ground beef"]
	vegetableOil := enums.Ingredients["vegetable oil"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	americanCheese := enums.Ingredients["American cheese"]
	burgerBun := enums.Ingredients["burger bun"]

	// Get measurement units
	ounceMeasurement := enums.MeasurementUnits["ounce"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	sliceMeasurement := enums.MeasurementUnits["slice"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	wideSpatula := enums.Instruments["wide spatula"]
	bareHands := enums.Instruments["bare hands"]

	// Get vessels
	castIronSkillet := enums.Vessels["cast iron skillet"]
	servingPlate := enums.Vessels["serving plate"]

	// Get ingredient states for completion conditions
	smokingState := enums.IngredientStates["smoking"]
	brownedState := enums.IngredientStates["browned"]

	// Get bridge table entries
	// Heat preparation bridges
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][vegetableOil.ID]
	oilTeaspoonVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][teaspoonMeasurement.ID]
	heatSkilletVPV := enums.PreparationVessels[heatPrep.ID][castIronSkillet.ID]

	// Divide preparation bridges
	divideBeefVIP := enums.IngredientPreparations[dividePrep.ID][groundBeef.ID]
	beefOunceVIMU := enums.IngredientMeasurementUnits[groundBeef.ID][ounceMeasurement.ID]
	divideBareHandsVPI := enums.PreparationInstruments[dividePrep.ID][bareHands.ID]

	// Form preparation bridges
	formBeefVIP := enums.IngredientPreparations[formPrep.ID][groundBeef.ID]
	formBareHandsVPI := enums.PreparationInstruments[formPrep.ID][bareHands.ID]

	// Season preparation bridges
	seasonBeefVIP := enums.IngredientPreparations[seasonPrep.ID][groundBeef.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]

	// Toast preparation bridges
	toastBunVIP := enums.IngredientPreparations[toastPrep.ID][burgerBun.ID]
	bunUnitVIMU := enums.IngredientMeasurementUnits[burgerBun.ID][unitMeasurement.ID]

	// Smash preparation bridges
	smashBeefVIP := enums.IngredientPreparations[smashPrep.ID][groundBeef.ID]
	smashSpatulaVPI := enums.PreparationInstruments[smashPrep.ID][wideSpatula.ID]
	smashSkilletVPV := enums.PreparationVessels[smashPrep.ID][castIronSkillet.ID]

	// Pan-sear preparation bridges
	panSearSpatulaVPI := enums.PreparationInstruments[panSearPrep.ID][wideSpatula.ID]
	panSearSkilletVPV := enums.PreparationVessels[panSearPrep.ID][castIronSkillet.ID]

	// Flip preparation bridges
	flipSpatulaVPI := enums.PreparationInstruments[flipPrep.ID][wideSpatula.ID]
	flipSkilletVPV := enums.PreparationVessels[flipPrep.ID][castIronSkillet.ID]

	// Top preparation bridges
	topCheeseVIP := enums.IngredientPreparations[topPrep.ID][americanCheese.ID]
	cheeseSliceVIMU := enums.IngredientMeasurementUnits[americanCheese.ID][sliceMeasurement.ID]
	topSkilletVPV := enums.PreparationVessels[topPrep.ID][castIronSkillet.ID]

	// Assemble preparation bridges
	assembleBunVIP := enums.IngredientPreparations[assemblePrep.ID][burgerBun.ID]
	assembleServingPlateVPV := enums.PreparationVessels[assemblePrep.ID][servingPlate.ID]

	// Step 0: Add oil to skillet and preheat
	step0ID := identifiers.New()
	step0OilIngredientID := identifiers.New()
	step0CompletionConditionID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           0,
		Notes:           "Add oil to a 12-inch cast iron skillet and wipe around with a paper towel. Set skillet over medium heat and allow to preheat for about 5 minutes, then increase heat to high until smoking.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step0OilIngredientID,
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilTeaspoonVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "vegetable oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &heatSkilletVPV.ID,
				VesselID:                 &castIronSkillet.ID,
				Name:                     "12-inch cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "hot skillet with smoking oil",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step0CompletionConditionID,
				BelongsToRecipeStep: step0ID,
				IngredientStateID:   smokingState.ID,
				Notes:               "Oil should be smoking",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step0CompletionConditionID,
						RecipeStepIngredient:                   step0OilIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 1: Divide ground beef into portions
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dividePrep.ID,
		Index:           1,
		Notes:           "Divide the ground beef into four 4-ounce portions.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &divideBeefVIP.ID,
				ValidIngredientMeasurementUnitID: &beefOunceVIMU.ID,
				IngredientID:                     &groundBeef.ID,
				MeasurementUnitID:                ounceMeasurement.ID,
				Name:                             "ground beef",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 16,
					Max: pointer.To[float32](20),
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidPreparationInstrumentID: &divideBareHandsVPI.ID,
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
				Name:                "portioned ground beef",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &ounceMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](16),
				},
			},
		},
	}

	// Step 2: Form beef into pucks
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   formPrep.ID,
		Index:           2,
		Notes:           "Gently form each portion of ground beef into a cylindrical puck about 2 inches tall, pressing together just until meat holds its shape without falling apart.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &formBeefVIP.ID,
				ValidIngredientMeasurementUnitID: &beefOunceVIMU.ID,
				IngredientID:                     &groundBeef.ID,
				MeasurementUnitID:                ounceMeasurement.ID,
				Name:                             "portioned ground beef",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 16,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step2ID,
				ValidPreparationInstrumentID: &formBareHandsVPI.ID,
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
				BelongsToRecipeStep: step2ID,
				Name:                "beef pucks",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 3: Season pucks
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           3,
		Notes:           "Season the beef pucks generously on all sides with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &seasonBeefVIP.ID,
				ValidIngredientMeasurementUnitID: &beefOunceVIMU.ID,
				IngredientID:                     &groundBeef.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "beef pucks",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
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
				BelongsToRecipeStep:              step3ID,
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
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "seasoned beef pucks",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 4: Toast buns
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   toastPrep.ID,
		Index:           4,
		Notes:           "Lightly toast the burger buns.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ValidIngredientPreparationID:     &toastBunVIP.ID,
				ValidIngredientMeasurementUnitID: &bunUnitVIMU.ID,
				IngredientID:                     &burgerBun.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "burger buns",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "toasted burger buns",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 5: Add pucks to skillet and smash
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   smashPrep.ID,
		Index:           5,
		Notes:           "Add 2 beef pucks to the hot skillet and, using a firm, stiff metal spatula, press down on each one until they're roughly 4 to 4 1/2 inches in diameter and 1/2-inch thick. It helps to use a second spatula to apply downward pressure to the first if you are having trouble smashing them hard enough.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &smashBeefVIP.ID,
				ValidIngredientMeasurementUnitID: &beefOunceVIMU.ID,
				IngredientID:                     &groundBeef.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "seasoned beef pucks",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &smashSpatulaVPI.ID,
				InstrumentID:                 &wideSpatula.ID,
				Name:                         "firm, stiff metal spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step5ID,
				ProductOfRecipeStepIndex: pointer.To[uint64](0),
				ValidPreparationVesselID: &smashSkilletVPV.ID,
				VesselID:                 &castIronSkillet.ID,
				Name:                     "hot skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "smashed burger patties",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 6: Sear first side until golden brown
	step6ID := identifiers.New()
	step6PattyIngredientID := identifiers.New()
	step6CompletionConditionID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   panSearPrep.ID,
		Index:           6,
		Notes:           "Cook without moving until a golden brown crust develops, about 1 1/2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](90), // 1.5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step6PattyIngredientID,
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &groundBeef.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "smashed burger patties",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: &panSearSpatulaVPI.ID,
				InstrumentID:                 &wideSpatula.ID,
				Name:                         "metal spatula",
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
				BelongsToRecipeStep: step6ID,
				Name:                "seared burger patties (first side)",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step6CompletionConditionID,
				BelongsToRecipeStep: step6ID,
				IngredientStateID:   brownedState.ID,
				Notes:               "A golden brown crust should develop on the bottom",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step6CompletionConditionID,
						RecipeStepIngredient:                   step6PattyIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 7: Flip patties
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   flipPrep.ID,
		Index:           7,
		Notes:           "Use the edge of the spatula to carefully scrape up and flip the patties one at a time, making sure to get all browned bits removed from the skillet.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &groundBeef.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "seared burger patties (first side)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step7ID,
				ValidPreparationInstrumentID: &flipSpatulaVPI.ID,
				InstrumentID:                 &wideSpatula.ID,
				Name:                         "metal spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step7ID,
				ValidPreparationVesselID: &flipSkilletVPV.ID,
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
				BelongsToRecipeStep: step7ID,
				Name:                "flipped burger patties",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 8: Top with cheese (optional)
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   topPrep.ID,
		Index:           8,
		Optional:        true,
		Notes:           "If using cheese, add a slice to each patty now.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &groundBeef.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "flipped burger patties",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &topCheeseVIP.ID,
				ValidIngredientMeasurementUnitID: &cheeseSliceVIMU.ID,
				IngredientID:                     &americanCheese.ID,
				MeasurementUnitID:                sliceMeasurement.ID,
				Name:                             "cheese slices",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: &topSkilletVPV.ID,
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
				Name:                "burger patties with cheese",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 9: Finish cooking second side
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   panSearPrep.ID,
		Index:           9,
		Notes:           "Continue to cook until patties are cooked to desired doneness—about 30 seconds longer for medium-rare.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &groundBeef.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "burger patties with cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step9ID,
				ValidPreparationVesselID: &panSearSkilletVPV.ID,
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
				BelongsToRecipeStep: step9ID,
				Name:                "cooked smash burger patties",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 10: Assemble burgers
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   assemblePrep.ID,
		Index:           10,
		Notes:           "Transfer patties to toasted buns, topping buns and/or patties as desired, close burgers, and serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &assembleBunVIP.ID,
				ValidIngredientMeasurementUnitID: &bunUnitVIMU.ID,
				IngredientID:                     &burgerBun.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "toasted burger buns",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &groundBeef.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "cooked smash burger patties",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: &assembleServingPlateVPV.ID,
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
				Name:                "classic smash burger",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	return &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Classic Smashed Burgers",
		Slug:                "classic-smashed-burgers",
		Source:              "https://www.seriouseats.com/classic-smashed-burgers-recipe",
		Description:         "Classic smashed cheeseburgers with maximum juiciness and a deep-brown, beefy crust. Smashing down on the burger patties within the first 30 seconds of hitting a hot skillet ensures maximum juiciness and a flavorful, well-browned crust.",
		YieldsComponentType: mealplanning.MealComponentTypesMain,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "burger",
		PluralPortionName: "burgers",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10,
		},
	}
}
