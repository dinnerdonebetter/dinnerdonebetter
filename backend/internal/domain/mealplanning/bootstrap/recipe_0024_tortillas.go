package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// TortillasRecipe creates the Simple Tortillas recipe.
// Source: https://www.kingarthurbaking.com/recipes/simple-tortillas-recipe
func TortillasRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	mixPrep := enums.Preparations["mix"]
	addPrep := enums.Preparations["add"]
	stirPrep := enums.Preparations["stir"]
	kneadPrep := enums.Preparations["knead"]
	dividePrep := enums.Preparations["divide"]
	formPrep := enums.Preparations["form"]
	restPrep := enums.Preparations["rest"]
	coverPrep := enums.Preparations["cover"]
	preheatPrep := enums.Preparations["preheat"]
	rollPrep := enums.Preparations["roll"]
	cookPrep := enums.Preparations["cook"]
	transferPrep := enums.Preparations["transfer"]

	// Get ingredients
	flour := enums.Ingredients["flour"]
	bakingPowder := enums.Ingredients["baking powder"]
	salt := enums.Ingredients["salt"]
	lard := enums.Ingredients["lard"]
	water := enums.Ingredients["water"]

	// Get measurement units
	cupMeasurement := enums.MeasurementUnits["cup"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	whisk := enums.Instruments["whisk"]
	pastryBlender := enums.Instruments["pastry blender"]
	bareHands := enums.Instruments["bare hands"]
	fork := enums.Instruments["fork"]
	rollingPin := enums.Instruments["rolling pin"]

	// Get vessels
	mediumBowl := enums.Vessels["medium bowl"]
	countertop := enums.Vessels["countertop"]
	castIronSkillet := enums.Vessels["cast iron skillet"]
	kitchenTowel := enums.Vessels["kitchen towel"]

	// Get ingredient states
	pliableState := enums.IngredientStates["pliable"]

	// Helper to safely get VIP bridge
	getVIP := func(prep, ingName string) *mealplanning.ValidIngredientPreparation {
		p := enums.Preparations[prep]
		i := enums.Ingredients[ingName]
		if p == nil || i == nil {
			return nil
		}
		prepMap := enums.IngredientPreparations[p.ID]
		if prepMap == nil {
			return nil
		}
		return prepMap[i.ID]
	}

	// Helper to safely get VPV bridge
	getVPV := func(prep, vesselName string) *mealplanning.ValidPreparationVessel {
		p := enums.Preparations[prep]
		v := enums.Vessels[vesselName]
		if p == nil || v == nil {
			return nil
		}
		prepMap := enums.PreparationVessels[p.ID]
		if prepMap == nil {
			return nil
		}
		return prepMap[v.ID]
	}

	// Helper to safely get VPI bridge
	getVPI := func(prep, instName string) *mealplanning.ValidPreparationInstrument {
		p := enums.Preparations[prep]
		i := enums.Instruments[instName]
		if p == nil || i == nil {
			return nil
		}
		prepMap := enums.PreparationInstruments[p.ID]
		if prepMap == nil {
			return nil
		}
		return prepMap[i.ID]
	}

	// Helper to safely get VIMU bridge
	getVIMU := func(ingName, unitName string) *mealplanning.ValidIngredientMeasurementUnit {
		i := enums.Ingredients[ingName]
		u := enums.MeasurementUnits[unitName]
		if i == nil || u == nil {
			return nil
		}
		ingMap := enums.IngredientMeasurementUnits[i.ID]
		if ingMap == nil {
			return nil
		}
		return ingMap[u.ID]
	}

	// Helper to safely get MealPlanTaskID pointer from VIP
	vipID := func(v *mealplanning.ValidIngredientPreparation) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Helper to safely get MealPlanTaskID pointer from VPV
	vpvID := func(v *mealplanning.ValidPreparationVessel) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Helper to safely get MealPlanTaskID pointer from VPI
	vpiID := func(v *mealplanning.ValidPreparationInstrument) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Helper to safely get MealPlanTaskID pointer from VIMU
	vimuID := func(v *mealplanning.ValidIngredientMeasurementUnit) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Bridge entries
	mixFlourVIP := getVIP("mix", "flour")
	mixBakingPowderVIP := getVIP("mix", "baking powder")
	mixSaltVIP := getVIP("mix", "salt")
	mixLardVIP := getVIP("mix", "lard")
	mixMediumBowlVPV := getVPV("mix", "medium bowl")
	mixWhiskVPI := getVPI("mix", "whisk")
	mixPastryBlenderVPI := getVPI("mix", "pastry blender")

	addLardVIP := getVIP("add", "lard")
	addWaterVIP := getVIP("add", "water")
	addMediumBowlVPV := getVPV("add", "medium bowl")

	stirMediumBowlVPV := getVPV("stir", "medium bowl")
	stirForkVPI := getVPI("stir", "fork")

	kneadFlourVIP := getVIP("knead", "flour")
	kneadCountertopVPV := getVPV("knead", "countertop")
	kneadBareHandsVPI := getVPI("knead", "bare hands")

	divideFlourVIP := getVIP("divide", "flour")
	divideCountertopVPV := getVPV("divide", "countertop")
	divideBareHandsVPI := getVPI("divide", "bare hands")

	formFlourVIP := getVIP("form", "flour")
	formCountertopVPV := getVPV("form", "countertop")
	formBareHandsVPI := getVPI("form", "bare hands")

	restFlourVIP := getVIP("rest", "flour")
	restCountertopVPV := getVPV("rest", "countertop")

	coverKitchenTowelVPV := getVPV("cover", "kitchen towel")

	preheatCastIronSkilletVPV := getVPV("preheat", "cast iron skillet")

	rollFlourVIP := getVIP("roll", "flour")
	rollCountertopVPV := getVPV("roll", "countertop")
	rollRollingPinVPI := getVPI("roll", "rolling pin")

	cookFlourVIP := getVIP("cook", "flour")
	cookCastIronSkilletVPV := getVPV("cook", "cast iron skillet")

	transferFlourVIP := getVIP("transfer", "flour")
	transferKitchenTowelVPV := getVPV("transfer", "kitchen towel")

	// Measurement unit bridges
	flourCupVIMU := getVIMU("flour", "cup")
	bakingPowderTeaspoonVIMU := getVIMU("baking powder", "teaspoon")
	saltTeaspoonVIMU := getVIMU("salt", "teaspoon")
	lardTablespoonVIMU := getVIMU("lard", "tablespoon")
	waterCupVIMU := getVIMU("water", "cup")

	// ==================== RECIPE STEPS ====================

	// Step 0: Whisk together the flour, baking powder, and salt
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   mixPrep.ID,
		Index:           0,
		Notes:           "In a medium-sized bowl, whisk together the flour, baking powder, and salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     vipID(mixFlourVIP),
				ValidIngredientMeasurementUnitID: vimuID(flourCupVIMU),
				IngredientID:                     &flour.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "all-purpose flour",
				QuantityNotes:                    "300g, plus additional as needed",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 2.5},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     vipID(mixBakingPowderVIP),
				ValidIngredientMeasurementUnitID: vimuID(bakingPowderTeaspoonVIMU),
				IngredientID:                     &bakingPowder.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "baking powder",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     vipID(mixSaltVIP),
				ValidIngredientMeasurementUnitID: vimuID(saltTeaspoonVIMU),
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "table salt",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 0.5},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
				ValidPreparationInstrumentID: vpiID(mixWhiskVPI),
				InstrumentID:                 &whisk.ID,
				Name:                         "whisk",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: vpvID(mixMediumBowlVPV),
				VesselID:                 &mediumBowl.ID,
				Name:                     "medium-sized bowl",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "dry ingredient mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "bowl with dry ingredients",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 1: Add the lard (fat) to the dry ingredients
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           1,
		Notes:           "Add the lard (or butter, shortening, or vegetable oil).",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "dry ingredient mixture",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     vipID(addLardVIP),
				ValidIngredientMeasurementUnitID: vimuID(lardTablespoonVIMU),
				IngredientID:                     &lard.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "lard",
				QuantityNotes:                    "57g, room temperature (or butter, shortening, or vegetable oil)",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 4},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(addMediumBowlVPV),
				VesselID:                        &mediumBowl.ID,
				Name:                            "bowl with dry ingredients",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "dry ingredients with fat",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "bowl with dry ingredients and fat",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 2: Work the fat into the flour
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   mixPrep.ID,
		Index:           2,
		Notes:           "Use your fingers or a pastry blender to work the fat into the flour until it disappears. Coating most of the flour with fat inhibits gluten formation, making the tortillas easier to roll out.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(mixLardVIP),
				IngredientID:                    &lard.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "dry ingredients with fat",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step2ID,
				ValidPreparationInstrumentID: vpiID(mixPastryBlenderVPI),
				InstrumentID:                 &pastryBlender.ID,
				Name:                         "pastry blender or fingers",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(mixMediumBowlVPV),
				VesselID:                        &mediumBowl.ID,
				Name:                            "bowl with dry ingredients and fat",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "flour-fat mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "bowl with flour-fat mixture",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 3: Add hot water and stir to bring the dough together
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           3,
		Notes:           "Pour in the lesser amount of hot water (110°F to 120°F).",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "flour-fat mixture",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     vipID(addWaterVIP),
				ValidIngredientMeasurementUnitID: vimuID(waterCupVIMU),
				IngredientID:                     &water.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "hot water",
				QuantityNotes:                    "200g to 227g, about 110°F to 120°F",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 0.875, Max: pointer.To[float32](1)},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(addMediumBowlVPV),
				VesselID:                        &mediumBowl.ID,
				Name:                            "bowl with flour-fat mixture",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "flour mixture with water",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "bowl with flour mixture and water",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 4: Stir briskly to bring the dough together
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           4,
		Notes:           "Stir briskly with a fork or whisk to bring the dough together into a shaggy mass. Stir in additional water as needed to bring the dough together.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "flour mixture with water",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: vpiID(stirForkVPI),
				InstrumentID:                 &fork.ID,
				Name:                         "fork or whisk",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(stirMediumBowlVPV),
				VesselID:                        &mediumBowl.ID,
				Name:                            "bowl with flour mixture and water",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "shaggy dough",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 5: Turn the dough out onto a lightly floured counter and knead briefly
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   kneadPrep.ID,
		Index:           5,
		Notes:           "Turn the dough out onto a lightly floured counter and knead briefly, just until the dough forms a ball. If the dough is very sticky, gradually add a bit more flour.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(kneadFlourVIP),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "shaggy dough",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: vpiID(kneadBareHandsVPI),
				InstrumentID:                 &bareHands.ID,
				Name:                         "hands",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step5ID,
				ValidPreparationVesselID: vpvID(kneadCountertopVPV),
				VesselID:                 &countertop.ID,
				Name:                     "lightly floured counter",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "kneaded dough ball",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "countertop with dough",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 6: Divide the dough into 8 pieces
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dividePrep.ID,
		Index:           6,
		Notes:           "Divide the dough into 8 pieces.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(divideFlourVIP),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "kneaded dough ball",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: vpiID(divideBareHandsVPI),
				InstrumentID:                 &bareHands.ID,
				Name:                         "hands or knife",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(divideCountertopVPV),
				VesselID:                        &countertop.ID,
				Name:                            "countertop with dough",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "dough pieces",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](8)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "countertop with dough pieces",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 7: Round the pieces into balls and flatten slightly
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   formPrep.ID,
		Index:           7,
		Notes:           "Round the pieces into balls and flatten slightly. If you wish, coat each ball lightly in oil before covering; this ensures the dough doesn't dry out.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(formFlourVIP),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "dough pieces",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 8},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step7ID,
				ValidPreparationInstrumentID: vpiID(formBareHandsVPI),
				InstrumentID:                 &bareHands.ID,
				Name:                         "hands",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(formCountertopVPV),
				VesselID:                        &countertop.ID,
				Name:                            "countertop with dough pieces",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "flattened dough balls",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](8)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "countertop with dough balls",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 8: Cover the dough balls
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   coverPrep.ID,
		Index:           8,
		Notes:           "Cover the dough balls with a clean kitchen towel.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "flattened dough balls",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 8},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				VesselID:                        &countertop.ID,
				Name:                            "countertop with dough balls",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: vpvID(coverKitchenTowelVPV),
				VesselID:                 &kitchenTowel.ID,
				Name:                     "clean kitchen towel",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "covered dough balls",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](8)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "countertop with covered dough",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 9: Allow the dough balls to rest for about 30 minutes
	step9ID := identifiers.New()
	step9DoughIngredientID := identifiers.New()
	step9CompletionConditionID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           9,
		Notes:           "Allow them to rest, covered, for about 30 minutes. The resting period improves the texture of the dough by giving the flour time to absorb the water. The tortillas will roll out more easily if you include the rest.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](1800), // 30 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step9DoughIngredientID,
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(restFlourVIP),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "covered dough balls",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 8},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(restCountertopVPV),
				VesselID:                        &countertop.ID,
				Name:                            "countertop with covered dough",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step9CompletionConditionID,
				BelongsToRecipeStep: step9ID,
				IngredientStateID:   pliableState.ID,
				Notes:               "dough should be pliable and easier to roll out",
				Optional:            false,
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step9CompletionConditionID,
						RecipeStepIngredient:                   step9DoughIngredientID,
					},
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "rested dough balls",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](8)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "countertop with rested dough",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 10: Preheat an ungreased cast iron griddle or skillet
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   preheatPrep.ID,
		Index:           10,
		Notes:           "While the dough rests, preheat an ungreased cast iron griddle or skillet over medium high heat, about 400°F.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](204), // 400°F = ~204°C
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: vpvID(preheatCastIronSkilletVPV),
				VesselID:                 &castIronSkillet.ID,
				Name:                     "ungreased cast iron griddle or skillet",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "preheated skillet",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
			},
		},
	}

	// Step 11: Roll one piece of dough into a round about 8" in diameter
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   rollPrep.ID,
		Index:           11,
		Notes:           "Working with one piece of dough at a time, roll into a round about 8\" in diameter. Keep the remaining dough covered while you work.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(rollFlourVIP),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "rested dough balls",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 8},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step11ID,
				ValidPreparationInstrumentID: vpiID(rollRollingPinVPI),
				InstrumentID:                 &rollingPin.ID,
				Name:                         "rolling pin",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(rollCountertopVPV),
				VesselID:                        &countertop.ID,
				Name:                            "countertop with rested dough",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "raw tortilla rounds",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](8)},
			},
		},
	}

	// Step 12: Cook the tortilla in the ungreased pan
	step12ID := identifiers.New()
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cookPrep.ID,
		Index:           12,
		Notes:           "Cook the tortilla in the ungreased pan for about 30 seconds on each side. Repeat with the remaining dough balls.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30),
			Max: pointer.To[uint32](60),
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](204), // 400°F
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(cookFlourVIP),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "raw tortilla rounds",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 8},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        vpvID(cookCastIronSkilletVPV),
				VesselID:                        &castIronSkillet.ID,
				Name:                            "preheated skillet",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "cooked tortillas",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](8)},
			},
		},
	}

	// Step 13: Wrap the tortilla in a clean cloth
	step13ID := identifiers.New()
	step13 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step13ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           13,
		Notes:           "Wrap the tortilla in a clean cloth when it comes off the griddle, to keep it pliable.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(transferFlourVIP),
				IngredientID:                    &flour.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "cooked tortillas",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 8},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step13ID,
				ValidPreparationVesselID: vpvID(transferKitchenTowelVPV),
				VesselID:                 &kitchenTowel.ID,
				Name:                     "clean kitchen towel",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "Tortillas",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](8)},
			},
		},
	}

	recipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Tortillas",
		Slug:                "tortillas",
		Source:              "https://www.kingarthurbaking.com/recipes/simple-tortillas-recipe",
		Description:         "This recipe for soft flour tortillas is quick and easy. Soft and tender, with just a little bit of \"chew,\" you can have these on the table in under an hour. The dough can also be made the day before and allowed to rest in the fridge overnight.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 8,
		},
		PortionName:       "tortilla",
		PluralPortionName: "tortillas",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{recipe}
}
