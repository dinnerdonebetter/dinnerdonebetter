package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// RefriedBeansRecipe creates the Perfect Frijoles Refritos (Mexican Refried Beans) recipe.
// Source: https://www.seriouseats.com/perfect-refried-beans
func RefriedBeansRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	halvePrep := enums.Preparations["halve"]
	mincePrep := enums.Preparations["mince"]
	peelPrep := enums.Preparations["peel"]
	coverPrep := enums.Preparations["cover"]
	addPrep := enums.Preparations["add"]
	boilPrep := enums.Preparations["boil"]
	reducePrep := enums.Preparations["reduce"]
	simmerPrep := enums.Preparations["simmer"]
	seasonPrep := enums.Preparations["season"]
	drainPrep := enums.Preparations["drain"]
	measurePrep := enums.Preparations["measure"]
	discardPrep := enums.Preparations["discard"]
	heatPrep := enums.Preparations["heat"]
	sautPrep := enums.Preparations["sauté"]
	stirPrep := enums.Preparations["stir"]
	smashPrep := enums.Preparations["smash"]
	dilutePrep := enums.Preparations["dilute"]

	// Get ingredients
	pintoBeans := enums.Ingredients["pinto beans"]
	water := enums.Ingredients["water"]
	epazote := enums.Ingredients["epazote"]
	whiteOnion := enums.Ingredients["onion"]
	garlic := enums.Ingredients["garlic"]
	salt := enums.Ingredients["salt"]
	lard := enums.Ingredients["lard"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]
	cloveMeasurement := enums.MeasurementUnits["clove"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	chefsKnife := enums.Instruments["knife"]
	bareHands := enums.Instruments["bare hands"]
	potatoMasher := enums.Instruments["potato masher"]

	// Get vessels
	cuttingBoard := enums.Vessels["cutting board"]
	largePot := enums.Vessels["pot"]
	largeSkillet := enums.Vessels["cast iron skillet"]
	largeBowl := enums.Vessels["large bowl"]

	// Get bridge table entries
	// Halve
	halveWhiteOnionVIP := enums.IngredientPreparations[halvePrep.ID][whiteOnion.ID]
	halveChefsKnifeVPI := enums.PreparationInstruments[halvePrep.ID][chefsKnife.ID]
	halveCuttingBoardVPV := enums.PreparationVessels[halvePrep.ID][cuttingBoard.ID]

	// Mince
	minceWhiteOnionVIP := enums.IngredientPreparations[mincePrep.ID][whiteOnion.ID]
	minceChefsKnifeVPI := enums.PreparationInstruments[mincePrep.ID][chefsKnife.ID]
	minceCuttingBoardVPV := enums.PreparationVessels[mincePrep.ID][cuttingBoard.ID]

	// Peel
	peelGarlicVIP := enums.IngredientPreparations[peelPrep.ID][garlic.ID]
	peelBareHandsVPI := enums.PreparationInstruments[peelPrep.ID][bareHands.ID]

	// Cover
	coverPintoBeansVIP := enums.IngredientPreparations[coverPrep.ID][pintoBeans.ID]
	coverWaterVIP := enums.IngredientPreparations[coverPrep.ID][water.ID]
	coverLargePotVPV := enums.PreparationVessels[coverPrep.ID][largePot.ID]

	// Add
	addEpazoteVIP := enums.IngredientPreparations[addPrep.ID][epazote.ID]
	addWhiteOnionVIP := enums.IngredientPreparations[addPrep.ID][whiteOnion.ID]
	addGarlicVIP := enums.IngredientPreparations[addPrep.ID][garlic.ID]
	addWaterVIP := enums.IngredientPreparations[addPrep.ID][water.ID]
	addLargePotVPV := enums.PreparationVessels[addPrep.ID][largePot.ID]

	// Boil
	boilLargePotVPV := enums.PreparationVessels[boilPrep.ID][largePot.ID]

	// Reduce
	reduceLargePotVPV := enums.PreparationVessels[reducePrep.ID][largePot.ID]

	// Simmer
	simmerPintoBeansVIP := enums.IngredientPreparations[simmerPrep.ID][pintoBeans.ID]
	simmerLargePotVPV := enums.PreparationVessels[simmerPrep.ID][largePot.ID]

	// Season
	seasonPintoBeansVIP := enums.IngredientPreparations[seasonPrep.ID][pintoBeans.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonLargePotVPV := enums.PreparationVessels[seasonPrep.ID][largePot.ID]
	seasonLargeSkilletVPV := enums.PreparationVessels[seasonPrep.ID][largeSkillet.ID]

	// Drain
	drainPintoBeansVIP := enums.IngredientPreparations[drainPrep.ID][pintoBeans.ID]
	drainLargePotVPV := enums.PreparationVessels[drainPrep.ID][largePot.ID]
	drainLargeBowlVPV := enums.PreparationVessels[drainPrep.ID][largeBowl.ID]

	// Reserve (not used in this recipe, but bridge table entry exists)

	// Measure
	measurePintoBeansVIP := enums.IngredientPreparations[measurePrep.ID][pintoBeans.ID]
	measureLargeBowlVPV := enums.PreparationVessels[measurePrep.ID][largeBowl.ID]

	// Discard
	discardEpazoteVIP := enums.IngredientPreparations[discardPrep.ID][epazote.ID]
	discardWhiteOnionVIP := enums.IngredientPreparations[discardPrep.ID][whiteOnion.ID]
	discardGarlicVIP := enums.IngredientPreparations[discardPrep.ID][garlic.ID]
	discardLargePotVPV := enums.PreparationVessels[discardPrep.ID][largePot.ID]

	// Heat
	heatLardVIP := enums.IngredientPreparations[heatPrep.ID][lard.ID]
	heatLargeSkilletVPV := enums.PreparationVessels[heatPrep.ID][largeSkillet.ID]

	// Sauté
	sautWhiteOnionVIP := enums.IngredientPreparations[sautPrep.ID][whiteOnion.ID]
	sautLargeSkilletVPV := enums.PreparationVessels[sautPrep.ID][largeSkillet.ID]

	// Stir
	stirPintoBeansVIP := enums.IngredientPreparations[stirPrep.ID][pintoBeans.ID]
	stirWhiteOnionVIP := enums.IngredientPreparations[stirPrep.ID][whiteOnion.ID]
	stirLargeSkilletVPV := enums.PreparationVessels[stirPrep.ID][largeSkillet.ID]

	// Smash
	smashPintoBeansVIP := enums.IngredientPreparations[smashPrep.ID][pintoBeans.ID]
	smashLargeSkilletVPV := enums.PreparationVessels[smashPrep.ID][largeSkillet.ID]
	smashPotatoMasherVPI := enums.PreparationInstruments[smashPrep.ID][potatoMasher.ID]

	// Dilute
	dilutePintoBeansVIP := enums.IngredientPreparations[dilutePrep.ID][pintoBeans.ID]
	diluteWaterVIP := enums.IngredientPreparations[dilutePrep.ID][water.ID]
	diluteLargeSkilletVPV := enums.PreparationVessels[dilutePrep.ID][largeSkillet.ID]

	// Measurement unit bridges
	pintoBeansPoundVIMU := enums.IngredientMeasurementUnits[pintoBeans.ID][poundMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	epazoteSprigVIMU := enums.IngredientMeasurementUnits[epazote.ID][sprigMeasurement.ID]
	whiteOnionUnitVIMU := enums.IngredientMeasurementUnits[whiteOnion.ID][unitMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	lardTablespoonVIMU := enums.IngredientMeasurementUnits[lard.ID][tablespoonMeasurement.ID]

	// Step 0: Halve the onion
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   halvePrep.ID,
		Index:           0,
		Notes:           "Halve the medium white onion.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &halveWhiteOnionVIP.ID,
				ValidIngredientMeasurementUnitID: &whiteOnionUnitVIMU.ID,
				IngredientID:                     &whiteOnion.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "medium white onion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
				ValidPreparationInstrumentID: &halveChefsKnifeVPI.ID,
				InstrumentID:                 &chefsKnife.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &halveCuttingBoardVPV.ID,
				VesselID:                 &cuttingBoard.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "halved onion (2 halves)",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 1: In a large pot, cover the beans with cold water by at least 2 inches
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   coverPrep.ID,
		Index:           1,
		Notes:           "In a large pot, cover the beans with cold water by at least 2 inches.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &coverPintoBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &pintoBeansPoundVIMU.ID,
				IngredientID:                     &pintoBeans.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "dried pinto or black beans",
				QuantityNotes:                    "1/2 pound (227 g)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &coverWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				IngredientID:                     &water.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "cold water",
				QuantityNotes:                    "Enough to cover beans by at least 2 inches",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4, // Approximate cups
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step1ID,
				ValidPreparationVesselID: &coverLargePotVPV.ID,
				VesselID:                 &largePot.ID,
				Name:                     "large pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "beans covered with water",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Peel garlic cloves
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   peelPrep.ID,
		Index:           2,
		Notes:           "Peel 2 medium cloves garlic.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &peelGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                cloveMeasurement.ID,
				Name:                             "medium cloves garlic",
				QuantityNotes:                    "2 medium cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step2ID,
				ValidPreparationInstrumentID: &peelBareHandsVPI.ID,
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
				Name:                "peeled garlic cloves",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cloveMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 3: Add herb sprigs, the whole onion half, and peeled garlic cloves
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           3,
		Notes:           "Add herb sprigs, the whole onion half, and peeled garlic cloves.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &addEpazoteVIP.ID,
				ValidIngredientMeasurementUnitID: &epazoteSprigVIMU.ID,
				IngredientID:                     &epazote.ID,
				MeasurementUnitID:                sprigMeasurement.ID,
				Name:                             "fresh epazote",
				QuantityNotes:                    "2 sprigs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addWhiteOnionVIP.ID,
				IngredientID:                    &whiteOnion.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "onion half (left whole)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addGarlicVIP.ID,
				IngredientID:                    &garlic.ID,
				MeasurementUnitID:               cloveMeasurement.ID,
				Name:                            "peeled garlic cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &addLargePotVPV.ID,
				VesselID:                        &largePot.ID,
				Name:                            "large pot with beans covered with water",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "beans with aromatics in pot",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Bring to a boil over high heat
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   boilPrep.ID,
		Index:           4,
		Notes:           "Bring to a boil over high heat.",
		Ingredients:     []*mealplanning.RecipeStepIngredientDatabaseCreationInput{},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &boilLargePotVPV.ID,
				VesselID:                        &largePot.ID,
				Name:                            "beans with aromatics in pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "boiling beans with aromatics",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Reduce heat to simmer
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   reducePrep.ID,
		Index:           5,
		Notes:           "Reduce heat to simmer.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &reduceLargePotVPV.ID,
				VesselID:                        &largePot.ID,
				Name:                            "boiling beans with aromatics",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "beans ready to simmer",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Simmer until beans are very tender, about 1 to 2 hours
	step6ID := identifiers.New()
	step6BeansIngredientID := identifiers.New()
	step6CompletionConditionID := identifiers.New()
	tenderState := enums.IngredientStates["tender"]
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   simmerPrep.ID,
		Index:           6,
		Notes:           "Simmer until beans are very tender, about 1 to 2 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](3600), // 1 hour
			Max: pointer.To[uint32](7200), // 2 hours
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step6BeansIngredientID,
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &simmerPintoBeansVIP.ID,
				IngredientID:                    &pintoBeans.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "beans ready to simmer",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &simmerLargePotVPV.ID,
				VesselID:                        &largePot.ID,
				Name:                            "beans ready to simmer",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "very tender cooked beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "pot with cooked beans and cooking liquid",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step6CompletionConditionID,
				BelongsToRecipeStep: step6ID,
				IngredientStateID:   tenderState.ID,
				Notes:               "Beans should be very tender",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step6CompletionConditionID,
						RecipeStepIngredient:                   step6BeansIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 7: Season with salt
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           7,
		Notes:           "Season with salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonPintoBeansVIP.ID,
				IngredientID:                    &pintoBeans.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "very tender cooked beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
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
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &seasonLargePotVPV.ID,
				VesselID:                        &largePot.ID,
				Name:                            "pot with cooked beans and cooking liquid",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "seasoned cooked beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 8: Drain beans, reserving bean-cooking liquid
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   drainPrep.ID,
		Index:           8,
		Notes:           "Drain beans, reserving bean-cooking liquid.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &drainPintoBeansVIP.ID,
				IngredientID:                    &pintoBeans.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "seasoned cooked beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &drainLargePotVPV.ID,
				VesselID:                        &largePot.ID,
				Name:                            "pot with cooked beans and cooking liquid",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &drainLargeBowlVPV.ID,
				VesselID:                 &largeBowl.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "drained cooked beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "reserved bean-cooking liquid",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               1,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Measure out 3 cups of beans (if you have more, reserve the rest for another use)
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   measurePrep.ID,
		Index:           9,
		Notes:           "You should have about 3 cups of cooked beans; if you have more, measure out 3 cups of beans and reserve the rest for another use.",
		Optional:        true,
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &measurePintoBeansVIP.ID,
				IngredientID:                    &pintoBeans.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "drained cooked beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step7ID,
				ValidPreparationVesselID: &measureLargeBowlVPV.ID,
				VesselID:                 &largeBowl.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "3 cups of cooked beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 10: Discard herb sprigs, onion, and garlic
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   discardPrep.ID,
		Index:           10,
		Notes:           "Discard herb sprigs, onion, and garlic.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &discardEpazoteVIP.ID,
				ValidIngredientMeasurementUnitID: &epazoteSprigVIMU.ID,
				IngredientID:                     &epazote.ID,
				MeasurementUnitID:                sprigMeasurement.ID,
				Name:                             "herb sprigs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &discardWhiteOnionVIP.ID,
				ValidIngredientMeasurementUnitID: &whiteOnionUnitVIMU.ID,
				IngredientID:                     &whiteOnion.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "whole onion half",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &discardGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                cloveMeasurement.ID,
				Name:                             "garlic cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &discardLargePotVPV.ID,
				VesselID:                        &largePot.ID,
				Name:                            "pot with cooked beans and cooking liquid",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
	}

	// Step 11: Mince one half of the onion
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   mincePrep.ID,
		Index:           11,
		Notes:           "Mince one half of the halved onion.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &minceWhiteOnionVIP.ID,
				IngredientID:                    &whiteOnion.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                           "onion half",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step10ID,
				ValidPreparationInstrumentID: &minceChefsKnifeVPI.ID,
				InstrumentID:                 &chefsKnife.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: &minceCuttingBoardVPV.ID,
				VesselID:                 &cuttingBoard.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "minced white onion",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
		},
	}

	// Step 12: In a large skillet, heat lard until shimmering over medium-high heat
	step12ID := identifiers.New()
	step12FatIngredientID := identifiers.New()
	step12CompletionConditionID := identifiers.New()
	shimmeringState := enums.IngredientStates["shimmering"]
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           12,
		Notes:           "In a large skillet, heat lard until shimmering over medium-high heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step12FatIngredientID,
				BelongsToRecipeStep:              step12ID,
				ValidIngredientPreparationID:     &heatLardVIP.ID,
				ValidIngredientMeasurementUnitID: &lardTablespoonVIMU.ID,
				IngredientID:                     &lard.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "lard",
				QuantityNotes:                    "6 tablespoons (77 g)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step11ID,
				ValidPreparationVesselID: &heatLargeSkilletVPV.ID,
				VesselID:                 &largeSkillet.ID,
				Name:                     "large skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "heated fat in skillet",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step12CompletionConditionID,
				BelongsToRecipeStep: step12ID,
				IngredientStateID:   shimmeringState.ID,
				Notes:               "Lard, bacon drippings, or oil should shimmer; butter should foam",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step12CompletionConditionID,
						RecipeStepIngredient:                   step12FatIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 13: Add minced onion and cook, stirring occasionally, until translucent and lightly golden, about 7 minutes
	step13ID := identifiers.New()
	step13OnionIngredientID := identifiers.New()
	step13CompletionConditionID := identifiers.New()
	translucentState := enums.IngredientStates["translucent"]
	step13 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step13ID,
		BelongsToRecipe: recipeID,
		PreparationID:   sautPrep.ID,
		Index:           13,
		Notes:           "Add minced onion and cook, stirring occasionally, until translucent and lightly golden, about 7 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](420), // 7 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step13OnionIngredientID,
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &sautWhiteOnionVIP.ID,
				IngredientID:                    &whiteOnion.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "minced white onion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &sautLargeSkilletVPV.ID,
				VesselID:                        &largeSkillet.ID,
				Name:                            "heated fat in skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "cooked minced onion",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step13CompletionConditionID,
				BelongsToRecipeStep: step13ID,
				IngredientStateID:   translucentState.ID,
				Notes:               "Onion should be translucent and lightly golden",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step13CompletionConditionID,
						RecipeStepIngredient:                   step13OnionIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 14: Stir in beans and cook for 2 minutes
	step14ID := identifiers.New()
	step14 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step14ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           14,
		Notes:           "Stir in beans and cook for 2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120), // 2 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirPintoBeansVIP.ID,
				IngredientID:                    &pintoBeans.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "3 cups of cooked beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirWhiteOnionVIP.ID,
				IngredientID:                    &whiteOnion.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "cooked minced onion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &stirLargeSkilletVPV.ID,
				VesselID:                        &largeSkillet.ID,
				Name:                            "heated fat in skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "beans and onion in skillet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 15: Add 1/4 cup of reserved bean-cooking liquid
	step15ID := identifiers.New()
	step15 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step15ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           15,
		Notes:           "Add 1/4 cup of reserved bean-cooking liquid.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step15ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &addWaterVIP.ID,
				IngredientID:                    &water.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "reserved bean-cooking liquid",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step15ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				VesselID:                        &largeSkillet.ID,
				Name:                            "skillet with beans and onion",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "beans with liquid in skillet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 16: Using bean masher, potato masher, or back of a wooden spoon, smash the beans to form a chunky purée; alternatively, use a stick blender to make a smoother purée
	step16ID := identifiers.New()
	step16 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step16ID,
		BelongsToRecipe: recipeID,
		PreparationID:   smashPrep.ID,
		Index:           16,
		Notes:           "Using bean masher, potato masher, or back of a wooden spoon, smash the beans to form a chunky purée; alternatively, use a stick blender to make a smoother purée.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step16ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &smashPintoBeansVIP.ID,
				IngredientID:                    &pintoBeans.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "beans with liquid in skillet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step16ID,
				ValidPreparationInstrumentID: &smashPotatoMasherVPI.ID,
				InstrumentID:                 &potatoMasher.ID,
				Name:                         "potato masher",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step16ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &smashLargeSkilletVPV.ID,
				VesselID:                        &largeSkillet.ID,
				Name:                            "skillet with beans and onion",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "mashed beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 17: Thin with more bean cooking water until desired consistency is reached. If refried beans become too wet, simmer, stirring, until thickened; if they become too dry, add more bean-cooking liquid, 1 tablespoon at a time, as needed.
	step17ID := identifiers.New()
	step17BeansIngredientID := identifiers.New()
	step17CompletionConditionID := identifiers.New()
	desiredConsistencyState := enums.IngredientStates["at desired consistency"]
	step17 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step17ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dilutePrep.ID,
		Index:           17,
		Notes:           "Thin with more bean cooking water until desired consistency is reached. If refried beans become too wet, simmer, stirring, until thickened; if they become too dry, add more bean-cooking liquid, 1 tablespoon at a time, as needed.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step17BeansIngredientID,
				BelongsToRecipeStep:             step17ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](16),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &dilutePintoBeansVIP.ID,
				IngredientID:                    &pintoBeans.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "mashed beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step17ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &diluteWaterVIP.ID,
				IngredientID:                    &water.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "reserved bean-cooking liquid",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step17ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](15),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &diluteLargeSkilletVPV.ID,
				VesselID:                        &largeSkillet.ID,
				Name:                            "skillet with beans and onion",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step14ID,
				Name:                "refried beans at desired consistency",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step17CompletionConditionID,
				BelongsToRecipeStep: step17ID,
				IngredientStateID:   desiredConsistencyState.ID,
				Notes:               "Refried beans should reach desired consistency - not too wet or too dry",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step17CompletionConditionID,
						RecipeStepIngredient:                   step17BeansIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 18: Season with salt and serve
	step18ID := identifiers.New()
	step18 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step18ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           18,
		Notes:           "Season with salt and serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step18ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](17),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonPintoBeansVIP.ID,
				IngredientID:                    &pintoBeans.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "refried beans at desired consistency",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step18ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step18ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](17),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &seasonLargeSkilletVPV.ID,
				VesselID:                        &largeSkillet.ID,
				Name:                            "skillet with beans and onion",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step15ID,
				Name:                "refried beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	refriedBeansRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Perfect Frijoles Refritos (Mexican Refried Beans)",
		Slug:                "refried-beans",
		Source:              "https://www.seriouseats.com/perfect-refried-beans",
		Description:         "Use this master recipe to make perfect refried beans in any style: chunky or smooth; with black beans or pintos; and using your choice of cooking fat. By offering choices, including bean type, fat type, and mashing technique, this recipe makes it possible to get exactly the style of refried beans you want.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "cup",
		PluralPortionName: "cups",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13, step14, step15, step16, step17, step18,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		refriedBeansRecipe,
	}
}
