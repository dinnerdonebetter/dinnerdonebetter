package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// GlazedCarrotsWithBrownButterAndSageRecipe creates the Glazed Carrots with Brown Butter and Sage recipe.
// Source: https://www.seriouseats.com/glazed-carrots-recipe-11856362
func GlazedCarrotsWithBrownButterAndSageRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	meltPrep := enums.Preparations["melt"]
	cookPrep := enums.Preparations["cook"]
	addPrep := enums.Preparations["add"]
	coverPrep := enums.Preparations["cover"]
	boilPrep := enums.Preparations["boil"]
	uncoverPrep := enums.Preparations["uncover"]
	reducePrep := enums.Preparations["reduce"]
	removeFromHeatPrep := enums.Preparations["remove from heat"]
	discardPrep := enums.Preparations["discard"]
	seasonPrep := enums.Preparations["season"]
	sprinklePrep := enums.Preparations["sprinkle"]

	// Get ingredients
	butter := enums.Ingredients["butter"]
	sage := enums.Ingredients["sage"]
	carrot := enums.Ingredients["carrot"]
	appleCider := enums.Ingredients["apple cider"]
	chickenStock := enums.Ingredients["chicken stock"]
	honey := enums.Ingredients["honey"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	appleCiderVinegar := enums.Ingredients["apple cider vinegar"]
	parsley := enums.Ingredients["parsley"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	poundMeasurement := enums.MeasurementUnits["pound"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	spoon := enums.Instruments["spoon"]

	// Get vessels
	pan := enums.Vessels["pan"]
	servingBowl := enums.Vessels["serving bowl"]

	// Get ingredient states for completion conditions
	crispState := enums.IngredientStates["crisp"]
	brownedState := enums.IngredientStates["browned"]
	glazedState := enums.IngredientStates["glazed"]
	tenderState := enums.IngredientStates["tender"]

	// === BRIDGE TABLE ENTRIES ===
	// Melt preparation bridges
	meltButterVIP := enums.IngredientPreparations[meltPrep.ID][butter.ID]
	meltPanVPV := enums.PreparationVessels[meltPrep.ID][pan.ID]
	meltSpoonVPI := enums.PreparationInstruments[meltPrep.ID][spoon.ID]

	// Cook preparation bridges
	cookSageVIP := enums.IngredientPreparations[cookPrep.ID][sage.ID]
	cookPanVPV := enums.PreparationVessels[cookPrep.ID][pan.ID]
	cookSpoonVPI := enums.PreparationInstruments[cookPrep.ID][spoon.ID]

	// Add preparation bridges
	addCarrotVIP := enums.IngredientPreparations[addPrep.ID][carrot.ID]
	addAppleCiderVIP := enums.IngredientPreparations[addPrep.ID][appleCider.ID]
	addChickenStockVIP := enums.IngredientPreparations[addPrep.ID][chickenStock.ID]
	addHoneyVIP := enums.IngredientPreparations[addPrep.ID][honey.ID]
	addSaltVIP := enums.IngredientPreparations[addPrep.ID][salt.ID]
	addBlackPepperVIP := enums.IngredientPreparations[addPrep.ID][blackPepper.ID]
	addPanVPV := enums.PreparationVessels[addPrep.ID][pan.ID]

	// Boil preparation bridges
	boilCarrotVIP := enums.IngredientPreparations[boilPrep.ID][carrot.ID]
	boilPanVPV := enums.PreparationVessels[boilPrep.ID][pan.ID]

	// Cover preparation bridges
	coverPanVPV := enums.PreparationVessels[coverPrep.ID][pan.ID]

	// Uncover preparation bridges
	uncoverPanVPV := enums.PreparationVessels[uncoverPrep.ID][pan.ID]

	// Reduce preparation bridges
	reducePanVPV := enums.PreparationVessels[reducePrep.ID][pan.ID]
	reduceSpoonVPI := enums.PreparationInstruments[reducePrep.ID][spoon.ID]

	// Remove from heat preparation bridges
	removeFromHeatPanVPV := enums.PreparationVessels[removeFromHeatPrep.ID][pan.ID]

	// Discard preparation bridges
	discardSageVIP := enums.IngredientPreparations[discardPrep.ID][sage.ID]
	discardPanVPV := enums.PreparationVessels[discardPrep.ID][pan.ID]

	// Season preparation bridges
	seasonCarrotVIP := enums.IngredientPreparations[seasonPrep.ID][carrot.ID]
	seasonAppleCiderVinegarVIP := enums.IngredientPreparations[seasonPrep.ID][appleCiderVinegar.ID]
	seasonPanVPV := enums.PreparationVessels[seasonPrep.ID][pan.ID]
	seasonSpoonVPI := enums.PreparationInstruments[seasonPrep.ID][spoon.ID]

	// Sprinkle preparation bridges
	sprinkleParsleyVIP := enums.IngredientPreparations[sprinklePrep.ID][parsley.ID]
	sprinkleServingBowlVPV := enums.PreparationVessels[sprinklePrep.ID][servingBowl.ID]

	// Measurement unit bridges
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]
	sageSprigVIMU := enums.IngredientMeasurementUnits[sage.ID][sprigMeasurement.ID]
	carrotPoundVIMU := enums.IngredientMeasurementUnits[carrot.ID][poundMeasurement.ID]
	appleCiderCupVIMU := enums.IngredientMeasurementUnits[appleCider.ID][cupMeasurement.ID]
	chickenStockCupVIMU := enums.IngredientMeasurementUnits[chickenStock.ID][cupMeasurement.ID]
	honeyTablespoonVIMU := enums.IngredientMeasurementUnits[honey.ID][tablespoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	blackPepperTeaspoonVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]
	appleCiderVinegarTeaspoonVIMU := enums.IngredientMeasurementUnits[appleCiderVinegar.ID][teaspoonMeasurement.ID]
	parsleyTablespoonVIMU := enums.IngredientMeasurementUnits[parsley.ID][tablespoonMeasurement.ID]

	// Suppress unused variable warnings
	_ = addAppleCiderVIP
	_ = addChickenStockVIP
	_ = addHoneyVIP
	_ = addSaltVIP
	_ = addBlackPepperVIP

	// ==================== RECIPE STEPS ====================

	// Step 0: Melt butter in skillet over medium heat and cook until browned
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   meltPrep.ID,
		Index:           0,
		Notes:           "In a deep 12-inch stainless-steel skillet, melt butter over medium heat, stirring often, until melted, about 2 minutes. Once melted, continue to cook, stirring constantly, just until milk solids separate and sink to the bottom of the skillet and begin to darken, 2 to 3 minutes. Butter can go from brown to burnt quickly, so keep a close eye on it as you stir it.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](240), // 4 minutes total (2 + 2)
			Max: pointer.To[uint32](300), // 5 minutes total (2 + 3)
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &meltButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				IngredientID:                     &butter.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "unsalted butter, cut into 6 pieces",
				QuantityNotes:                    "about 75g",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
				ValidPreparationInstrumentID: &meltSpoonVPI.ID,
				InstrumentID:                 &spoon.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &meltPanVPV.ID,
				VesselID:                 &pan.ID,
				Name:                     "deep 12-inch stainless-steel skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				IngredientStateID:   brownedState.ID,
				Notes:               "milk solids should separate and sink to the bottom and begin to darken",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "browning butter",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "skillet with browning butter",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 1: Add sage sprigs and cook until crisp
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cookPrep.ID,
		Index:           1,
		Notes:           "Add sage sprigs; cook, stirring constantly, until sage leaves darken and crisp and butter foams and browns, 1 to 2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60),  // 1 minute
			Max: pointer.To[uint32](120), // 2 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &butter.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "browning butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &cookSageVIP.ID,
				ValidIngredientMeasurementUnitID: &sageSprigVIMU.ID,
				IngredientID:                     &sage.ID,
				MeasurementUnitID:                sprigMeasurement.ID,
				Name:                             "5-inch long sage sprigs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidPreparationInstrumentID: &cookSpoonVPI.ID,
				InstrumentID:                 &spoon.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &cookPanVPV.ID,
				VesselID:                        &pan.ID,
				Name:                            "skillet with browning butter",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				IngredientStateID:   crispState.ID,
				Notes:               "sage leaves should darken and crisp",
				Optional:            false,
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				IngredientStateID:   brownedState.ID,
				Notes:               "butter should foam and brown",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "brown butter with sage",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Add carrots and liquids
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           2,
		Notes:           "Quickly add carrots, apple cider, chicken or vegetable stock, honey, salt, and pepper to brown butter in skillet. Carrots should be almost submerged, if not, add a small amount of stock until they are.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &butter.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "brown butter with sage",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addCarrotVIP.ID,
				ValidIngredientMeasurementUnitID: &carrotPoundVIMU.ID,
				IngredientID:                     &carrot.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "medium carrots, peeled and sliced on the bias into 1/2 inch–thick discs",
				QuantityNotes:                    "about 910g",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addAppleCiderVIP.ID,
				ValidIngredientMeasurementUnitID: &appleCiderCupVIMU.ID,
				IngredientID:                     &appleCider.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "apple cider",
				QuantityNotes:                    "about 240ml",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addChickenStockVIP.ID,
				ValidIngredientMeasurementUnitID: &chickenStockCupVIMU.ID,
				IngredientID:                     &chickenStock.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "homemade chicken stock or store-bought low-sodium chicken broth",
				QuantityNotes:                    "about 120ml",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addHoneyVIP.ID,
				ValidIngredientMeasurementUnitID: &honeyTablespoonVIMU.ID,
				IngredientID:                     &honey.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "honey",
				QuantityNotes:                    "about 45ml",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Diamond Crystal kosher salt",
				QuantityNotes:                    "plus more to taste; for table salt use half as much by volume",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addBlackPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &blackPepperTeaspoonVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &addPanVPV.ID,
				VesselID:                 &pan.ID,
				Name:                     "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "carrot mixture in skillet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Bring to a boil over high heat
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   boilPrep.ID,
		Index:           3,
		Notes:           "Bring to a boil over high heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &boilCarrotVIP.ID,
				IngredientID:                    &carrot.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "carrot mixture in skillet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &boilPanVPV.ID,
				VesselID:                 &pan.ID,
				Name:                     "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "boiling carrot mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Cover and continue to boil
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   coverPrep.ID,
		Index:           4,
		Notes:           "Cover, reduce heat to medium-high, and continue to boil, vigorously shaking the skillet occasionally, until carrots are crisp/tender and still firm in the center, about 8 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](480), // 8 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &carrot.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "boiling carrot mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &coverPanVPV.ID,
				VesselID:                 &pan.ID,
				Name:                     "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				IngredientStateID:   tenderState.ID,
				Notes:               "carrots should be crisp/tender and still firm in the center",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "covered boiling carrots",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Uncover and reduce to glaze
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   uncoverPrep.ID,
		Index:           5,
		Notes:           "Reduce heat to medium, uncover (the liquid should look creamy and still almost cover the carrots).",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &carrot.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "covered boiling carrots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step5ID,
				ValidPreparationVesselID: &uncoverPanVPV.ID,
				VesselID:                 &pan.ID,
				Name:                     "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "uncovered carrots in skillet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Continue boiling until reduced to glaze
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   reducePrep.ID,
		Index:           6,
		Notes:           "Continue to boil, vigorously stirring and shaking skillet often, until the mixture is reduced to a glaze that coats and clings to the carrots, 12 to 14 minutes. If the sauce begins to break and you see oily, butter-colored specks, add a splash of water (about 2 tablespoons) and return to a vigorous simmer, stirring constantly, until the mixture looks creamy and homogenous again.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](720), // 12 minutes
			Max: pointer.To[uint32](840), // 14 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &carrot.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "uncovered carrots in skillet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: &reduceSpoonVPI.ID,
				InstrumentID:                 &spoon.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &reducePanVPV.ID,
				VesselID:                 &pan.ID,
				Name:                     "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				IngredientStateID:   glazedState.ID,
				Notes:               "mixture should be reduced to a glaze that coats and clings to the carrots",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "glazed carrots in skillet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Remove from heat
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   removeFromHeatPrep.ID,
		Index:           7,
		Notes:           "Remove from heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &carrot.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "glazed carrots in skillet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step7ID,
				ValidPreparationVesselID: &removeFromHeatPanVPV.ID,
				VesselID:                 &pan.ID,
				Name:                     "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "glazed carrots off heat",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Discard sage sprigs
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   discardPrep.ID,
		Index:           8,
		Notes:           "Discard sage sprigs.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &discardSageVIP.ID,
				IngredientID:                    &sage.ID,
				MeasurementUnitID:               sprigMeasurement.ID,
				Name:                            "sage sprigs from glazed carrots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: &discardPanVPV.ID,
				VesselID:                 &pan.ID,
				Name:                     "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "glazed carrots without sage",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Stir in vinegar and season to taste
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           9,
		Notes:           "Stir in vinegar and season with salt to taste.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonCarrotVIP.ID,
				IngredientID:                    &carrot.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "glazed carrots without sage",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step9ID,
				ValidIngredientPreparationID:     &seasonAppleCiderVinegarVIP.ID,
				ValidIngredientMeasurementUnitID: &appleCiderVinegarTeaspoonVIMU.ID,
				IngredientID:                     &appleCiderVinegar.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "apple cider vinegar, white wine vinegar, or unseasoned rice vinegar",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step9ID,
				ValidPreparationInstrumentID: &seasonSpoonVPI.ID,
				InstrumentID:                 &spoon.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step9ID,
				ValidPreparationVesselID: &seasonPanVPV.ID,
				VesselID:                 &pan.ID,
				Name:                     "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "seasoned glazed carrots",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Sprinkle with herbs and serve
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   sprinklePrep.ID,
		Index:           10,
		Notes:           "Sprinkle with herbs. Serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &carrot.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "seasoned glazed carrots",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &sprinkleParsleyVIP.ID,
				ValidIngredientMeasurementUnitID: &parsleyTablespoonVIMU.ID,
				IngredientID:                     &parsley.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "chopped fresh tender herbs such as flat-leaf parsley, chives, and/or tarragon",
				QuantityNotes:                    "1-2 tablespoons, can use a combination",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
					Max: pointer.To[float32](2),
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: &sprinkleServingBowlVPV.ID,
				VesselID:                 &servingBowl.ID,
				Name:                     "serving bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "Glazed Carrots with Brown Butter and Sage",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](6),
					Max: pointer.To[float32](8),
				},
			},
		},
	}

	recipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Glazed Carrots with Brown Butter and Sage",
		Slug:                "glazed-carrots-with-brown-butter-and-sage",
		Source:              "https://www.seriouseats.com/glazed-carrots-recipe-11856362",
		Description:         "Perfectly tender carrots enhanced with brown butter, sage, and an emulsified buttery gloss. With only one skillet and a handful of ingredients, you can transform humble carrots into a first-class side dish for holidays and weeknights.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 6,
			Max: pointer.To[float32](8),
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{recipe}
}
