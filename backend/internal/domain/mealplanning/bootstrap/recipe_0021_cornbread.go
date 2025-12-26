package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// CornbreadRecipe creates the Cornbread recipe.
// Source: https://www.kingarthurbaking.com/recipes/cornbread-recipe
func CornbreadRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	preheatPrep := enums.Preparations["preheat"]
	greasePrep := enums.Preparations["grease"]
	mixPrep := enums.Preparations["mix"]
	pourPrep := enums.Preparations["pour"]
	combinePrep := enums.Preparations["combine"]
	bakePrep := enums.Preparations["bake"]
	coolPrep := enums.Preparations["cool"]

	// Get ingredients
	flour := enums.Ingredients["flour"]
	cornmeal := enums.Ingredients["cornmeal"]
	sugar := enums.Ingredients["sugar"]
	bakingPowder := enums.Ingredients["baking powder"]
	bakingSoda := enums.Ingredients["baking soda"]
	salt := enums.Ingredients["salt"]
	milk := enums.Ingredients["milk"]
	butter := enums.Ingredients["butter"]
	vegetableOil := enums.Ingredients["vegetable oil"]
	eggs := enums.Ingredients["eggs"]

	// Get measurement units
	cupMeasurement := enums.MeasurementUnits["cup"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	gramMeasurement := enums.MeasurementUnits["gram"]

	// Get instruments
	whisk := enums.Instruments["whisk"]
	spoon := enums.Instruments["spoon"]

	// Get vessels
	oven := enums.Vessels["oven"]
	bakingPan := enums.Vessels["baking pan"]
	mediumBowl := enums.Vessels["medium bowl"]
	largeBowl := enums.Vessels["large bowl"]
	wireRack := enums.Vessels["wire rack"]

	// Get ingredient states for completion conditions
	combinedState := enums.IngredientStates["combined"]
	bakedState := enums.IngredientStates["baked"]

	// === BRIDGE TABLE ENTRIES ===
	// Preheat preparation bridges
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]

	// Grease preparation bridges
	greaseButterVIP := enums.IngredientPreparations[greasePrep.ID][butter.ID]
	greaseBakingPanVPV := enums.PreparationVessels[greasePrep.ID][bakingPan.ID]

	// Mix preparation bridges
	mixFlourVIP := enums.IngredientPreparations[mixPrep.ID][flour.ID]
	mixCornmealVIP := enums.IngredientPreparations[mixPrep.ID][cornmeal.ID]
	mixSugarVIP := enums.IngredientPreparations[mixPrep.ID][sugar.ID]
	mixBakingPowderVIP := enums.IngredientPreparations[mixPrep.ID][bakingPowder.ID]
	mixBakingSodaVIP := enums.IngredientPreparations[mixPrep.ID][bakingSoda.ID]
	mixSaltVIP := enums.IngredientPreparations[mixPrep.ID][salt.ID]
	mixMilkVIP := enums.IngredientPreparations[mixPrep.ID][milk.ID]
	mixButterVIP := enums.IngredientPreparations[mixPrep.ID][butter.ID]
	mixVegetableOilVIP := enums.IngredientPreparations[mixPrep.ID][vegetableOil.ID]
	mixEggsVIP := enums.IngredientPreparations[mixPrep.ID][eggs.ID]
	mixMediumBowlVPV := enums.PreparationVessels[mixPrep.ID][mediumBowl.ID]
	mixLargeBowlVPV := enums.PreparationVessels[mixPrep.ID][largeBowl.ID]
	mixWhiskVPI := enums.PreparationInstruments[mixPrep.ID][whisk.ID]

	// Pour preparation bridges
	pourBakingPanVPV := enums.PreparationVessels[pourPrep.ID][bakingPan.ID]

	// Combine preparation bridges
	combineMediumBowlVPV := enums.PreparationVessels[combinePrep.ID][mediumBowl.ID]

	// Mix preparation bridges for step 5 (continued)
	mixSpoonVPI := enums.PreparationInstruments[mixPrep.ID][spoon.ID]

	// Bake preparation bridges
	bakeOvenVPV := enums.PreparationVessels[bakePrep.ID][oven.ID]
	bakeBakingPanVPV := enums.PreparationVessels[bakePrep.ID][bakingPan.ID]

	// Cool preparation bridges
	coolWireRackVPV := enums.PreparationVessels[coolPrep.ID][wireRack.ID]
	coolBakingPanVPV := enums.PreparationVessels[coolPrep.ID][bakingPan.ID]

	// Measurement unit bridges
	flourCupVIMU := enums.IngredientMeasurementUnits[flour.ID][cupMeasurement.ID]
	flourGramVIMU := enums.IngredientMeasurementUnits[flour.ID][gramMeasurement.ID]
	cornmealCupVIMU := enums.IngredientMeasurementUnits[cornmeal.ID][cupMeasurement.ID]
	cornmealGramVIMU := enums.IngredientMeasurementUnits[cornmeal.ID][gramMeasurement.ID]
	sugarCupVIMU := enums.IngredientMeasurementUnits[sugar.ID][cupMeasurement.ID]
	sugarGramVIMU := enums.IngredientMeasurementUnits[sugar.ID][gramMeasurement.ID]
	bakingPowderTeaspoonVIMU := enums.IngredientMeasurementUnits[bakingPowder.ID][teaspoonMeasurement.ID]
	bakingSodaTeaspoonVIMU := enums.IngredientMeasurementUnits[bakingSoda.ID][teaspoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	milkCupVIMU := enums.IngredientMeasurementUnits[milk.ID][cupMeasurement.ID]
	milkGramVIMU := enums.IngredientMeasurementUnits[milk.ID][gramMeasurement.ID]
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]
	butterGramVIMU := enums.IngredientMeasurementUnits[butter.ID][gramMeasurement.ID]
	vegetableOilCupVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][cupMeasurement.ID]
	vegetableOilGramVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][gramMeasurement.ID]
	eggsUnitVIMU := enums.IngredientMeasurementUnits[eggs.ID][unitMeasurement.ID]

	// Suppress unused variable warnings
	_ = flourGramVIMU
	_ = cornmealGramVIMU
	_ = sugarGramVIMU
	_ = milkGramVIMU
	_ = butterGramVIMU
	_ = vegetableOilGramVIMU
	_ = coolBakingPanVPV
	_ = mixLargeBowlVPV

	// ==================== RECIPE STEPS ====================

	// Step 0: Preheat the oven to 375°F
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   preheatPrep.ID,
		Index:           0,
		Notes:           "Preheat the oven to 375°F.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](190), // 375°F = ~190°C
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &preheatOvenVPV.ID,
				VesselID:                 &oven.ID,
				Name:                     "oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "preheated oven",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
			},
		},
	}

	// Step 1: Lightly grease a 9" square or round pan
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   greasePrep.ID,
		Index:           1,
		Notes:           "Lightly grease a 9\" square or round pan.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidIngredientPreparationID: &greaseButterVIP.ID,
				IngredientID:                 &butter.ID,
				MeasurementUnitID:            unitMeasurement.ID,
				Name:                         "butter for greasing",
				QuantityNotes:                "as needed",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step1ID,
				ValidPreparationVesselID: &greaseBakingPanVPV.ID,
				VesselID:                 &bakingPan.ID,
				Name:                     `9" square or round pan`,
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "greased baking pan",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
			},
		},
	}

	// Step 2: Whisk together the dry ingredients in a medium bowl
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   mixPrep.ID,
		Index:           2,
		Notes:           "In a medium bowl, mix together the flour, cornmeal, sugar, baking powder, baking soda, and salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &mixFlourVIP.ID,
				ValidIngredientMeasurementUnitID: &flourCupVIMU.ID,
				IngredientID:                     &flour.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "all-purpose flour",
				QuantityNotes:                    "210g",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1.75,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &mixCornmealVIP.ID,
				ValidIngredientMeasurementUnitID: &cornmealCupVIMU.ID,
				IngredientID:                     &cornmeal.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "yellow cornmeal",
				QuantityNotes:                    "156g",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &mixSugarVIP.ID,
				ValidIngredientMeasurementUnitID: &sugarCupVIMU.ID,
				IngredientID:                     &sugar.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "granulated sugar",
				QuantityNotes:                    "50g",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &mixBakingPowderVIP.ID,
				ValidIngredientMeasurementUnitID: &bakingPowderTeaspoonVIMU.ID,
				IngredientID:                     &bakingPowder.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "baking powder",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &mixBakingSodaVIP.ID,
				ValidIngredientMeasurementUnitID: &bakingSodaTeaspoonVIMU.ID,
				IngredientID:                     &bakingSoda.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "baking soda",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &mixSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "table salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step2ID,
				ValidPreparationInstrumentID: &mixWhiskVPI.ID,
				InstrumentID:                 &whisk.ID,
				Name:                         "whisk",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &mixMediumBowlVPV.ID,
				VesselID:                 &mediumBowl.ID,
				Name:                     "medium bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "dry ingredient mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "bowl with dry ingredients",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 3: Whisk together the wet ingredients in another bowl
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   mixPrep.ID,
		Index:           3,
		Notes:           "In another bowl or large measuring cup, mix together the milk, melted butter, vegetable oil, and egg.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &mixMilkVIP.ID,
				ValidIngredientMeasurementUnitID: &milkCupVIMU.ID,
				IngredientID:                     &milk.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "milk, lukewarm",
				QuantityNotes:                    "283g",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &mixButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				IngredientID:                     &butter.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "unsalted butter, melted, cooled",
				QuantityNotes:                    "57g",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &mixVegetableOilVIP.ID,
				ValidIngredientMeasurementUnitID: &vegetableOilCupVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "vegetable oil",
				QuantityNotes:                    "50g",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &mixEggsVIP.ID,
				ValidIngredientMeasurementUnitID: &eggsUnitVIMU.ID,
				IngredientID:                     &eggs.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "large egg",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step3ID,
				ValidPreparationInstrumentID: &mixWhiskVPI.ID,
				InstrumentID:                 &whisk.ID,
				Name:                         "whisk",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &mixMediumBowlVPV.ID,
				VesselID:                 &largeBowl.ID,
				Name:                     "bowl or large measuring cup",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "wet ingredient mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "bowl with wet ingredients",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 4: Combine the liquid mixture with the flour mixture
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   combinePrep.ID,
		Index:           4,
		Notes:           "Pour the liquid all at once into the flour mixture.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &milk.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "wet ingredient mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "dry ingredient mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &combineMediumBowlVPV.ID,
				VesselID:                        &mediumBowl.ID,
				Name:                            "bowl with dry ingredients",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "combined wet and dry ingredients",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "bowl with combined ingredients",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 5: Mix until just combined
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   mixPrep.ID,
		Index:           5,
		Notes:           "Mix quickly and gently until just combined. Don't over mix: stir the batter just enough to bring it together and evenly moisten the ingredients.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "combined wet and dry ingredients",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &mixSpoonVPI.ID,
				InstrumentID:                 &spoon.ID,
				Name:                         "spoon or spatula",
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
				ValidPreparationVesselID:        &mixMediumBowlVPV.ID,
				VesselID:                        &mediumBowl.ID,
				Name:                            "bowl with combined ingredients",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				IngredientStateID:   combinedState.ID,
				Notes:               "batter should be just combined with no dry flour visible",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "cornbread batter",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Spread the batter into the prepared pan
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   pourPrep.ID,
		Index:           6,
		Notes:           "Spread the batter into the prepared pan.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "cornbread batter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &pourBakingPanVPV.ID,
				VesselID:                        &bakingPan.ID,
				Name:                            "greased baking pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "unbaked cornbread in pan",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "pan with cornbread batter",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 7: Bake for 20 to 25 minutes
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   bakePrep.ID,
		Index:           7,
		Notes:           "Bake the bread for 20 to 25 minutes, until the edges just begin to pull away from the pan and a cake tester or paring knife inserted in the center comes out clean.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](1200), // 20 minutes
			Max: pointer.To[uint32](1500), // 25 minutes
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](190), // 375°F
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "unbaked cornbread in pan",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &bakeOvenVPV.ID,
				VesselID:                        &oven.ID,
				Name:                            "preheated oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &bakeBakingPanVPV.ID,
				VesselID:                        &bakingPan.ID,
				Name:                            "pan with cornbread batter",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				IngredientStateID:   bakedState.ID,
				Notes:               "edges should pull away from the pan and a cake tester inserted in the center should come out clean",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "baked cornbread",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "pan with baked cornbread",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 8: Cool on a rack for 5 minutes
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   coolPrep.ID,
		Index:           8,
		Notes:           "Remove the bread from the oven and cool it on a rack for 5 minutes before cutting; serve warm.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "baked cornbread",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: &coolWireRackVPV.ID,
				VesselID:                 &wireRack.ID,
				Name:                     "wire rack",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "Cornbread",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](12),
				},
			},
		},
	}

	recipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Cornbread",
		Slug:                "cornbread",
		Source:              "https://www.kingarthurbaking.com/recipes/cornbread-recipe",
		Description:         "This cornbread is a rare compromise between Southern and Northern cornbreads: it's tender and moist, with pleasing corn flavor and just the right amount of sweetness.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 12,
		},
		PortionName:       "piece",
		PluralPortionName: "pieces",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{recipe}
}
