package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// CaesarSaladRecipe creates the Caesar Salad recipe.
// Source: https://www.seriouseats.com/the-best-caesar-salad-recipe
// Note: This recipe references the Caesar Dressing and Garlic Parmesan Croutons recipes, which must be created first.
// The createdRecipes map should contain the "caesar-dressing" and "garlic-parmesan-croutons" recipes keyed by their slugs.
func CaesarSaladRecipe(enums *Enumerations, createdRecipes map[string]*mealplanning.Recipe) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	tossPrep := enums.Preparations["toss"]
	sprinklePrep := enums.Preparations["sprinkle"]
	inspectPrep := enums.Preparations["inspect"]
	rinsePrep := enums.Preparations["rinse"]
	dryPrep := enums.Preparations["dry"]
	gratePrep := enums.Preparations["grate"]
	transferPrep := enums.Preparations["transfer"]

	// Get ingredients
	oliveOil := enums.Ingredients["olive oil"]
	heartyBread := enums.Ingredients["hearty bread"]
	parmesanCheese := enums.Ingredients["parmesan cheese"]
	romaineLettuce := enums.Ingredients["romaine lettuce"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	bareHands := enums.Instruments["bare hands"]
	cheeseGrater := enums.Instruments["cheese grater"]

	// Get vessels
	largeBowl := enums.Vessels["large bowl"]
	servingBowl := enums.Vessels["serving bowl"]
	saladSpinner := enums.Vessels["salad spinner"]
	cuttingBoard := enums.Vessels["cutting board"]

	// === SALAD BRIDGE TABLE ENTRIES ===
	// Toss preparation bridges
	tossRomaineVIP := enums.IngredientPreparations[tossPrep.ID][romaineLettuce.ID]
	tossParmesanVIP := enums.IngredientPreparations[tossPrep.ID][parmesanCheese.ID]
	tossBreadVIP := enums.IngredientPreparations[tossPrep.ID][heartyBread.ID]

	// Grate preparation bridges
	grateParmesanVIP := enums.IngredientPreparations[gratePrep.ID][parmesanCheese.ID]
	grateCheeseGraterVPI := enums.PreparationInstruments[gratePrep.ID][cheeseGrater.ID]
	grateCuttingBoardVPV := enums.PreparationVessels[gratePrep.ID][cuttingBoard.ID]

	// Inspect preparation bridges for salad
	inspectRomaineVIP := enums.IngredientPreparations[inspectPrep.ID][romaineLettuce.ID]
	inspectCuttingBoardSaladVPV := enums.PreparationVessels[inspectPrep.ID][cuttingBoard.ID]
	inspectBareHandsSaladVPI := enums.PreparationInstruments[inspectPrep.ID][bareHands.ID]

	// Rinse preparation bridges for salad
	rinseRomaineVIP := enums.IngredientPreparations[rinsePrep.ID][romaineLettuce.ID]
	rinseLargeBowlSaladVPV := enums.PreparationVessels[rinsePrep.ID][largeBowl.ID]

	// Dry preparation bridges for salad
	dryRomaineVIP := enums.IngredientPreparations[dryPrep.ID][romaineLettuce.ID]
	drySaladSpinnerSaladVPV := enums.PreparationVessels[dryPrep.ID][saladSpinner.ID]

	// Sprinkle preparation bridges
	sprinkleParmesanVIP := enums.IngredientPreparations[sprinklePrep.ID][parmesanCheese.ID]
	sprinkleBreadVIP := enums.IngredientPreparations[sprinklePrep.ID][heartyBread.ID]
	sprinkleServingBowlVPV := enums.PreparationVessels[sprinklePrep.ID][servingBowl.ID]

	// Transfer preparation bridges
	transferServingBowlVPV := enums.PreparationVessels[transferPrep.ID][servingBowl.ID]

	// Measurement unit bridges for salad
	romaineUnitVIMU := enums.IngredientMeasurementUnits[romaineLettuce.ID][unitMeasurement.ID]
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	breadCupVIMU := enums.IngredientMeasurementUnits[heartyBread.ID][cupMeasurement.ID]
	parmesanCupVIMU := enums.IngredientMeasurementUnits[parmesanCheese.ID][cupMeasurement.ID]

	// ==================== CAESAR SALAD RECIPE STEPS ====================

	// Step 0: Select inner romaine leaves
	slStep0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        inspectPrep.ID,
		Index:                0,
		ExplicitInstructions: "Select the inner romaine leaves, discarding the outer leaves.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &inspectRomaineVIP.ID,
				ValidIngredientMeasurementUnitID: &romaineUnitVIMU.ID,
				Name:                             "romaine lettuce",
				QuantityNotes:                    "select inner leaves only",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &inspectBareHandsSaladVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &inspectCuttingBoardSaladVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "inner romaine leaves",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 1: Wash romaine leaves
	slStep1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        rinsePrep.ID,
		Index:                1,
		ExplicitInstructions: "Wash the inner romaine leaves in several changes of water until no dirt or grit remains.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &rinseRomaineVIP.ID,
				Name:                            "inner romaine leaves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &rinseLargeBowlSaladVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "washed inner romaine leaves",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				Name:  "large bowl",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 2: Dry romaine leaves
	slStep2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        dryPrep.ID,
		Index:                2,
		ExplicitInstructions: "Carefully dry the washed romaine leaves using a salad spinner or paper towels.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &dryRomaineVIP.ID,
				Name:                            "washed inner romaine leaves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &drySaladSpinnerSaladVPV.ID,
				Name:                     "salad spinner",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "romaine lettuce, inner leaves only, washed and carefully dried",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 3: Toss lettuce with dressing
	slStep3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        tossPrep.ID,
		Index:                3,
		ExplicitInstructions: "Toss the lettuce with a few tablespoons of dressing, adding more if desired. Large leaves should be torn into smaller pieces, smaller leaves left intact.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossRomaineVIP.ID,
				Name:                            "romaine lettuce, inner leaves only, washed and carefully dried",
				QuantityNotes:                   "large leaves torn into smaller pieces, smaller leaves left intact",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				// RecipeStepProductRecipeID references the "Caesar Dressing" recipe (slug: "caesar-dressing")
				// The product "Caesar dressing" is from step 6 (index 6), product index 0
				// Note: ProductOfRecipeStepIndex refers to the step index in the OTHER recipe, not this one
				ProductOfRecipeStepIndex:         pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				RecipeStepProductRecipeID:        getRecipeIDBySlug(createdRecipes, "caesar-dressing"),
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				Name:                             "Caesar dressing",
				QuantityNotes:                    "add more if desired",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
					Max: pointer.To[float32](6),
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "dressed lettuce",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large bowl",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 4: Grate parmesan cheese
	slStep4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        gratePrep.ID,
		Index:                4,
		ExplicitInstructions: "Finely grate 1 cup parmesan cheese.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &grateParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanCupVIMU.ID,
				Name:                             "parmesan cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &grateCheeseGraterVPI.ID,
				Name:                         "cheese grater",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &grateCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "finely grated parmesan cheese",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Add cheese and croutons, toss again
	slStep5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        tossPrep.ID,
		Index:                5,
		ExplicitInstructions: "Once the lettuce is coated, add half of the remaining cheese and three-quarters of the croutons and toss again.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossRomaineVIP.ID,
				Name:                            "dressed lettuce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossParmesanVIP.ID,
				Name:                            "finely grated parmesan cheese",
				QuantityNotes:                   "half of remaining cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				// RecipeStepProductRecipeID references the "Garlic Parmesan Croutons" recipe (slug: "garlic-parmesan-croutons")
				// The product "garlic parmesan croutons" is from step 9 (index 9), product index 0
				// Note: ProductOfRecipeStepIndex refers to the step index in the OTHER recipe, not this one
				ProductOfRecipeStepIndex:         pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				RecipeStepProductRecipeID:        getRecipeIDBySlug(createdRecipes, "garlic-parmesan-croutons"),
				ValidIngredientPreparationID:     &tossBreadVIP.ID,
				ValidIngredientMeasurementUnitID: &breadCupVIMU.ID,
				Name:                             "garlic parmesan croutons",
				QuantityNotes:                    "three-quarters of croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "tossed Caesar salad",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Transfer to serving bowl
	slStep6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        transferPrep.ID,
		Index:                6,
		ExplicitInstructions: "Transfer to a salad bowl.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "tossed Caesar salad",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &transferServingBowlVPV.ID,
				Name:                     "salad bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "Caesar salad in serving bowl",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Sprinkle with remaining cheese and croutons
	slStep7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        sprinklePrep.ID,
		Index:                7,
		ExplicitInstructions: "Sprinkle with the remaining cheese and croutons. Serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "Caesar salad in serving bowl",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &sprinkleParmesanVIP.ID,
				Name:                            "finely grated parmesan cheese",
				QuantityNotes:                   "remaining cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				// RecipeStepProductRecipeID references the "Garlic Parmesan Croutons" recipe (slug: "garlic-parmesan-croutons")
				// The product "garlic parmesan croutons" is from step 9 (index 9), product index 0
				// Note: ProductOfRecipeStepIndex refers to the step index in the OTHER recipe, not this one
				ProductOfRecipeStepIndex:         pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				RecipeStepProductRecipeID:        getRecipeIDBySlug(createdRecipes, "garlic-parmesan-croutons"),
				ValidIngredientPreparationID:     &sprinkleBreadVIP.ID,
				ValidIngredientMeasurementUnitID: &breadCupVIMU.ID,
				Name:                             "garlic parmesan croutons",
				QuantityNotes:                    "remaining croutons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &sprinkleServingBowlVPV.ID,
				Name:                     "salad bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "Caesar salad",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	saladRecipe := &mealplanning.RecipeCreationRequestInput{
		Name:                "Caesar Salad",
		Slug:                "caesar-salad",
		Source:              "https://www.seriouseats.com/the-best-caesar-salad-recipe",
		Description:         "The crowd-pleasing salad of crisp romaine leaves, crunchy croutons, and a creamy, emulsified dressing with just the right amount of anchovy.",
		YieldsComponentType: mealplanning.MealComponentTypesSalad,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepCreationRequestInput{
			slStep0, slStep1, slStep2, slStep3, slStep4, slStep5, slStep6, slStep7,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
		Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		saladRecipe,
	}
}
