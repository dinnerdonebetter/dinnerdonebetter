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
	coverPrep := enums.Preparations["cover"]
	addPrep := enums.Preparations["add"]
	boilPrep := enums.Preparations["boil"]
	reducePrep := enums.Preparations["reduce"]
	simmerPrep := enums.Preparations["simmer"]
	seasonPrep := enums.Preparations["season"]
	drainPrep := enums.Preparations["drain"]
	reservePrep := enums.Preparations["reserve"]
	measurePrep := enums.Preparations["measure"]
	discardPrep := enums.Preparations["discard"]
	heatPrep := enums.Preparations["heat"]
	sautPrep := enums.Preparations["sauté"]
	stirPrep := enums.Preparations["stir"]
	smashPrep := enums.Preparations["smash"]
	thinPrep := enums.Preparations["thin"]

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
	potatoMasher := enums.Instruments["potato masher"]

	// Get vessels
	largePot := enums.Vessels["pot"]
	largeSkillet := enums.Vessels["cast iron skillet"]
	largeBowl := enums.Vessels["large bowl"]

	// Get bridge table entries
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

	// Reserve
	reserveLargeBowlVPV := enums.PreparationVessels[reservePrep.ID][largeBowl.ID]

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

	// Thin
	thinPintoBeansVIP := enums.IngredientPreparations[thinPrep.ID][pintoBeans.ID]
	thinWaterVIP := enums.IngredientPreparations[thinPrep.ID][water.ID]
	thinLargeSkilletVPV := enums.PreparationVessels[thinPrep.ID][largeSkillet.ID]

	// Measurement unit bridges
	pintoBeansPoundVIMU := enums.IngredientMeasurementUnits[pintoBeans.ID][poundMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	epazoteSprigVIMU := enums.IngredientMeasurementUnits[epazote.ID][sprigMeasurement.ID]
	whiteOnionUnitVIMU := enums.IngredientMeasurementUnits[whiteOnion.ID][unitMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	lardTablespoonVIMU := enums.IngredientMeasurementUnits[lard.ID][tablespoonMeasurement.ID]

	// Step 0: In a large pot, cover the beans with cold water by at least 2 inches
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   coverPrep.ID,
		Index:           0,
		Notes:           "In a large pot, cover the beans with cold water by at least 2 inches.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
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
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step0ID,
				ValidIngredientPreparationID:     &coverWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				IngredientID:                    &water.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "cold water",
				QuantityNotes:                   "Enough to cover beans by at least 2 inches",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4, // Approximate cups
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
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
				BelongsToRecipeStep: step0ID,
				Name:                "beans covered with water",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Add herb sprigs, the whole onion half, and garlic cloves
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           1,
		Notes:           "Add herb sprigs, the whole onion half, and garlic cloves.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &addEpazoteVIP.ID,
				ValidIngredientMeasurementUnitID: &epazoteSprigVIMU.ID,
				IngredientID:                     &epazote.ID,
				MeasurementUnitID:                sprigMeasurement.ID,
				Name:                             "fresh epazote or oregano",
				QuantityNotes:                    "2 sprigs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &addWhiteOnionVIP.ID,
				ValidIngredientMeasurementUnitID: &whiteOnionUnitVIMU.ID,
				IngredientID:                     &whiteOnion.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "medium white onion, 1/2 left whole",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &addGarlicVIP.ID,
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
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
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
				BelongsToRecipeStep: step1ID,
				Name:                "beans with aromatics in pot",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Bring to a boil over high heat
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   boilPrep.ID,
		Index:           2,
		Notes:           "Bring to a boil over high heat.",
		Ingredients:     []*mealplanning.RecipeStepIngredientDatabaseCreationInput{},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
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
				BelongsToRecipeStep: step2ID,
				Name:                "boiling beans with aromatics",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Reduce heat to simmer
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   reducePrep.ID,
		Index:           3,
		Notes:           "Reduce heat to simmer.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
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
				BelongsToRecipeStep: step3ID,
				Name:                "beans ready to simmer",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Simmer until beans are very tender, about 1 to 2 hours
	step4ID := identifiers.New()
	step4BeansIngredientID := identifiers.New()
	step4CompletionConditionID := identifiers.New()
	tenderState := enums.IngredientStates["tender"]
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   simmerPrep.ID,
		Index:           4,
		Notes:           "Simmer until beans are very tender, about 1 to 2 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](3600), // 1 hour
			Max: pointer.To[uint32](7200), // 2 hours
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step4BeansIngredientID,
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
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
				BelongsToRecipeStep:             step4ID,
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
				BelongsToRecipeStep: step4ID,
				Name:                "very tender cooked beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "pot with cooked beans and cooking liquid",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step4CompletionConditionID,
				BelongsToRecipeStep: step4ID,
				IngredientStateID:   tenderState.ID,
				Notes:               "Beans should be very tender",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step4CompletionConditionID,
						RecipeStepIngredient:                   step4BeansIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 5: Season with salt
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           5,
		Notes:           "Season with salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
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
				BelongsToRecipeStep:              step5ID,
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
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
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
				BelongsToRecipeStep: step5ID,
				Name:                "seasoned cooked beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 6: Drain beans, reserving bean-cooking liquid
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   drainPrep.ID,
		Index:           6,
		Notes:           "Drain beans, reserving bean-cooking liquid.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
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
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
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
				ValidPreparationVesselID: &reserveLargeBowlVPV.ID,
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
				BelongsToRecipeStep: step6ID,
				Name:                "drained cooked beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "reserved bean-cooking liquid",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               1,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Measure out 3 cups of beans (if you have more, reserve the rest for another use)
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   measurePrep.ID,
		Index:           7,
		Notes:           "You should have about 3 cups of cooked beans; if you have more, measure out 3 cups of beans and reserve the rest for another use.",
		Optional:        true,
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
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
				BelongsToRecipeStep: step7ID,
				Name:                "3 cups of cooked beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 8: Discard herb sprigs, onion, and garlic
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   discardPrep.ID,
		Index:           8,
		Notes:           "Discard herb sprigs, onion, and garlic.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
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
				BelongsToRecipeStep:              step8ID,
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
				BelongsToRecipeStep:              step8ID,
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
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
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

	// Step 9: In a large skillet, heat lard, bacon drippings, or oil until shimmering, or butter until foaming, over medium-high heat
	step9ID := identifiers.New()
	step9FatIngredientID := identifiers.New()
	step9CompletionConditionID := identifiers.New()
	shimmeringState := enums.IngredientStates["shimmering"]
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           9,
		Notes:           "In a large skillet, heat lard until shimmering over medium-high heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step9FatIngredientID,
				BelongsToRecipeStep:              step9ID,
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
				BelongsToRecipeStep:      step9ID,
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
				BelongsToRecipeStep: step9ID,
				Name:                "heated fat in skillet",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step9CompletionConditionID,
				BelongsToRecipeStep: step9ID,
				IngredientStateID:   shimmeringState.ID,
				Notes:               "Lard, bacon drippings, or oil should shimmer; butter should foam",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step9CompletionConditionID,
						RecipeStepIngredient:                   step9FatIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 10: Add minced onion and cook, stirring occasionally, until translucent and lightly golden, about 7 minutes
	step10ID := identifiers.New()
	step10OnionIngredientID := identifiers.New()
	step10CompletionConditionID := identifiers.New()
	translucentState := enums.IngredientStates["translucent"]
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   sautPrep.ID,
		Index:           10,
		Notes:           "Add minced onion and cook, stirring occasionally, until translucent and lightly golden, about 7 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](420), // 7 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step10OnionIngredientID,
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &sautWhiteOnionVIP.ID,
				ValidIngredientMeasurementUnitID: &whiteOnionUnitVIMU.ID,
				IngredientID:                     &whiteOnion.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "minced white onion",
				QuantityNotes:                    "1/2 minced (about 1/2 cup; 26 g)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
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
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step10CompletionConditionID,
				BelongsToRecipeStep: step10ID,
				IngredientStateID:   translucentState.ID,
				Notes:               "Onion should be translucent and lightly golden",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step10CompletionConditionID,
						RecipeStepIngredient:                   step10OnionIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 11: Stir in beans and cook for 2 minutes
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           11,
		Notes:           "Stir in beans and cook for 2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120), // 2 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
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
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
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
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
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
				BelongsToRecipeStep: step11ID,
				Name:                "beans and onion in skillet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 12: Add 1/4 cup of reserved bean-cooking liquid
	step12ID := identifiers.New()
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           12,
		Notes:           "Add 1/4 cup of reserved bean-cooking liquid.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
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
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
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
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 13: Using bean masher, potato masher, or back of a wooden spoon, smash the beans to form a chunky purée; alternatively, use a stick blender to make a smoother purée
	step13ID := identifiers.New()
	step13 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step13ID,
		BelongsToRecipe: recipeID,
		PreparationID:   smashPrep.ID,
		Index:           13,
		Notes:           "Using bean masher, potato masher, or back of a wooden spoon, smash the beans to form a chunky purée; alternatively, use a stick blender to make a smoother purée.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
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
				BelongsToRecipeStep:          step13ID,
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
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
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
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 14: Thin with more bean cooking water until desired consistency is reached. If refried beans become too wet, simmer, stirring, until thickened; if they become too dry, add more bean-cooking liquid, 1 tablespoon at a time, as needed.
	step14ID := identifiers.New()
	step14BeansIngredientID := identifiers.New()
	step14CompletionConditionID := identifiers.New()
	desiredConsistencyState := enums.IngredientStates["at desired consistency"]
	step14 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step14ID,
		BelongsToRecipe: recipeID,
		PreparationID:   thinPrep.ID,
		Index:           14,
		Notes:           "Thin with more bean cooking water until desired consistency is reached. If refried beans become too wet, simmer, stirring, until thickened; if they become too dry, add more bean-cooking liquid, 1 tablespoon at a time, as needed.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step14BeansIngredientID,
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &thinPintoBeansVIP.ID,
				IngredientID:                    &pintoBeans.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "mashed beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &thinWaterVIP.ID,
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
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &thinLargeSkilletVPV.ID,
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
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step14CompletionConditionID,
				BelongsToRecipeStep: step14ID,
				IngredientStateID:   desiredConsistencyState.ID,
				Notes:               "Refried beans should reach desired consistency - not too wet or too dry",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step14CompletionConditionID,
						RecipeStepIngredient:                   step14BeansIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 15: Season with salt and serve
	step15ID := identifiers.New()
	step15 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step15ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           15,
		Notes:           "Season with salt and serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step15ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
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
				BelongsToRecipeStep:              step15ID,
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
				BelongsToRecipeStep:             step15ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
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
				Quantity: types.OptionalFloat32Range{
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
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13, step14, step15,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		refriedBeansRecipe,
	}
}
