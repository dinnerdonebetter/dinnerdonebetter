package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// UltraFluffyMashedPotatoesRecipe creates the Ultra-Fluffy Mashed Potatoes recipe.
// Source: https://www.seriouseats.com/ultra-fluffy-mashed-potatoes-recipe
func UltraFluffyMashedPotatoesRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	peelPrep := enums.Preparations["peel"]
	cubePrep := enums.Preparations["cube"]
	rinsePrep := enums.Preparations["rinse"]
	submergePrep := enums.Preparations["submerge"]
	seasonPrep := enums.Preparations["season"]
	boilPrep := enums.Preparations["boil"]
	drainPrep := enums.Preparations["drain"]
	restPrep := enums.Preparations["rest"]
	ricePrep := enums.Preparations["rice"]
	foldPrep := enums.Preparations["fold"]
	simmerPrep := enums.Preparations["simmer"]

	// Get ingredients
	potato := enums.Ingredients["potato"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	butter := enums.Ingredients["butter"]
	milk := enums.Ingredients["milk"]
	water := enums.Ingredients["water"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	vegetablePeeler := enums.Instruments["vegetable peeler"]
	knife := enums.Instruments["knife"]
	potatoRicer := enums.Instruments["potato ricer"]
	rubberSpatula := enums.Instruments["rubber spatula"]

	// Get vessels
	cuttingBoard := enums.Vessels["cutting board"]
	pot := enums.Vessels["pot"]
	colander := enums.Vessels["colander"]

	// Get bridge table entries
	// Peel preparation bridges
	peelPotatoVIP := enums.IngredientPreparations[peelPrep.ID][potato.ID]
	peelVegetablePeelerVPI := enums.PreparationInstruments[peelPrep.ID][vegetablePeeler.ID]
	peelCuttingBoardVPV := enums.PreparationVessels[peelPrep.ID][cuttingBoard.ID]

	// Cube preparation bridges
	cubePotatoVIP := enums.IngredientPreparations[cubePrep.ID][potato.ID]
	cubeKnifeVPI := enums.PreparationInstruments[cubePrep.ID][knife.ID]
	cubeCuttingBoardVPV := enums.PreparationVessels[cubePrep.ID][cuttingBoard.ID]

	// Rinse preparation bridges
	rinsePotatoVIP := enums.IngredientPreparations[rinsePrep.ID][potato.ID]
	rinsePotVPV := enums.PreparationVessels[rinsePrep.ID][pot.ID]
	rinseColanderVPV := enums.PreparationVessels[rinsePrep.ID][colander.ID]

	// Submerge preparation bridges
	submergePotatoVIP := enums.IngredientPreparations[submergePrep.ID][potato.ID]
	submergeWaterVIP := enums.IngredientPreparations[submergePrep.ID][water.ID]
	submergePotVPV := enums.PreparationVessels[submergePrep.ID][pot.ID]

	// Season preparation bridges for pot (seasoning water)
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPotVPV := enums.PreparationVessels[seasonPrep.ID][pot.ID]

	// Boil preparation bridges
	boilPotatoVIP := enums.IngredientPreparations[boilPrep.ID][potato.ID]
	boilPotVPV := enums.PreparationVessels[boilPrep.ID][pot.ID]

	// Drain preparation bridges
	drainPotatoVIP := enums.IngredientPreparations[drainPrep.ID][potato.ID]
	drainColanderVPV := enums.PreparationVessels[drainPrep.ID][colander.ID]

	// Rest preparation bridges
	restPotatoVIP := enums.IngredientPreparations[restPrep.ID][potato.ID]
	restColanderVPV := enums.PreparationVessels[restPrep.ID][colander.ID]

	// Rice preparation bridges
	ricePotatoVIP := enums.IngredientPreparations[ricePrep.ID][potato.ID]
	ricePotatoRicerVPI := enums.PreparationInstruments[ricePrep.ID][potatoRicer.ID]
	ricePotVPV := enums.PreparationVessels[ricePrep.ID][pot.ID]

	// Fold preparation bridges
	foldPotatoVIP := enums.IngredientPreparations[foldPrep.ID][potato.ID]
	foldButterVIP := enums.IngredientPreparations[foldPrep.ID][butter.ID]
	foldRubberSpatulaVPI := enums.PreparationInstruments[foldPrep.ID][rubberSpatula.ID]
	foldPotVPV := enums.PreparationVessels[foldPrep.ID][pot.ID]

	// Season preparation bridges (for final seasoning step)
	seasonPotatoVIP := enums.IngredientPreparations[seasonPrep.ID][potato.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	// seasonPotVPV already defined above for seasoning water

	// Simmer preparation bridges
	simmerMilkVIP := enums.IngredientPreparations[simmerPrep.ID][milk.ID]
	simmerPotVPV := enums.PreparationVessels[simmerPrep.ID][pot.ID]

	// Measurement unit bridges
	potatoPoundVIMU := enums.IngredientMeasurementUnits[potato.ID][poundMeasurement.ID]
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	milkCupVIMU := enums.IngredientMeasurementUnits[milk.ID][cupMeasurement.ID]
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]

	// Step 0: Peel potatoes
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   peelPrep.ID,
		Index:           0,
		Notes:           "Peel the russet potatoes.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &peelPotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				IngredientID:                     &potato.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "russet potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
				ValidPreparationInstrumentID: &peelVegetablePeelerVPI.ID,
				InstrumentID:                 &vegetablePeeler.ID,
				Name:                         "vegetable peeler",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &peelCuttingBoardVPV.ID,
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
				Name:                "peeled potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 1: Cut potatoes into 1-2 inch cubes
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cubePrep.ID,
		Index:           1,
		Notes:           "Cut the peeled potatoes into 1- or 2-inch cubes.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &cubePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				IngredientID:                     &potato.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "peeled potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidPreparationInstrumentID: &cubeKnifeVPI.ID,
				InstrumentID:                 &knife.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step1ID,
				ValidPreparationVesselID: &cubeCuttingBoardVPV.ID,
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
				BelongsToRecipeStep: step1ID,
				Name:                "cubed potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 2: Rinse potatoes in pot of cold water
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   rinsePrep.ID,
		Index:           2,
		Notes:           "Transfer potatoes to a pot of cold water and rinse, changing the water 2 or 3 times until it runs clear.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &rinsePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				IngredientID:                     &potato.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "cubed potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &rinsePotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "rinsed potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 3: Cover potatoes with fresh cold water
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   submergePrep.ID,
		Index:           3,
		Notes:           "Cover potatoes with fresh cold water.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &submergePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				IngredientID:                     &potato.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "rinsed potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &submergeWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				IngredientID:                     &water.ID,
				QuantityNotes:                    "enough to cover potatoes",
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "cold water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8, // enough to cover potatoes
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &submergePotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "potatoes in water",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 4: Season the water generously with salt
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           4,
		Notes:           "Season the water generously with salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &potato.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "potatoes in water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltGramVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                gramMeasurement.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 15, // generous amount
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &seasonPotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "potatoes in salted water",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 5: Boil and simmer until tender
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   boilPrep.ID,
		Index:           5,
		Notes:           "Set over medium-high heat and bring to a boil, then reduce heat to maintain a gentle simmer. Cook until potatoes are completely tender, about 15 minutes after reaching a simmer.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &boilPotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				IngredientID:                     &potato.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "potatoes in salted water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step5ID,
				ValidPreparationVesselID: &boilPotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "boiled potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 6: Drain potatoes in a colander
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   drainPrep.ID,
		Index:           6,
		Notes:           "Drain potatoes in a colander.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &drainPotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				IngredientID:                     &potato.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "boiled potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &drainColanderVPV.ID,
				VesselID:                 &colander.ID,
				Name:                     "colander",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "drained potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 7: Rinse under hot running water
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   rinsePrep.ID,
		Index:           7,
		Notes:           "Rinse potatoes under hot running water for 30 seconds to wash away excess starch.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &rinsePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				IngredientID:                     &potato.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "drained potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step7ID,
				ValidPreparationVesselID: &rinseColanderVPV.ID,
				VesselID:                 &colander.ID,
				Name:                     "colander",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "rinsed hot potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 8: Allow potatoes to steam/rest
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           8,
		Notes:           "Allow potatoes to steam for 1 minute to remove excess moisture.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60), // 1 minute
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &restPotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				IngredientID:                     &potato.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "rinsed hot potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: &restColanderVPV.ID,
				VesselID:                 &colander.ID,
				Name:                     "colander",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "steamed potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 9: Pass potatoes through ricer
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   ricePrep.ID,
		Index:           9,
		Notes:           "Set a ricer or food mill over the now-empty pot and pass potatoes through.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step9ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &ricePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				IngredientID:                     &potato.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "steamed potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step9ID,
				ValidPreparationInstrumentID: &ricePotatoRicerVPI.ID,
				InstrumentID:                 &potatoRicer.ID,
				Name:                         "potato ricer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step9ID,
				ValidPreparationVesselID: &ricePotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "riced potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Fold in butter
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   foldPrep.ID,
		Index:           10,
		Notes:           "Add butter and gently fold into potatoes.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &foldPotatoVIP.ID,
				IngredientID:                    &potato.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "riced potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &foldButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				IngredientID:                     &butter.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "unsalted butter, room temperature, cut into 1/2-inch pats",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6, // 6 tablespoons
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step10ID,
				ValidPreparationInstrumentID: &foldRubberSpatulaVPI.ID,
				InstrumentID:                 &rubberSpatula.ID,
				Name:                         "rubber spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: &foldPotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "buttered potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 11: Simmer milk
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   simmerPrep.ID,
		Index:           11,
		Notes:           "Mound potatoes into the center of the pot and pour milk all around. Set over medium heat and bring milk to a simmer.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &foldPotatoVIP.ID,
				IngredientID:                    &potato.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "buttered potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &simmerMilkVIP.ID,
				ValidIngredientMeasurementUnitID: &milkCupVIMU.ID,
				IngredientID:                     &milk.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "whole milk",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5, // 1/2 cup, plus more as needed
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step11ID,
				ValidPreparationVesselID: &simmerPotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "potatoes with simmering milk",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 12: Fold simmered milk into potatoes
	step12ID := identifiers.New()
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   foldPrep.ID,
		Index:           12,
		Notes:           "Gently fold the simmered milk into the potatoes. If looser potatoes are desired, add additional milk in a similar fashion around the mashed potato mass and bring it to a simmer before folding into potatoes.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &foldPotatoVIP.ID,
				IngredientID:                    &potato.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "potatoes with simmering milk",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step12ID,
				ValidPreparationInstrumentID: &foldRubberSpatulaVPI.ID,
				InstrumentID:                 &rubberSpatula.ID,
				Name:                         "rubber spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step12ID,
				ValidPreparationVesselID: &foldPotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "mashed potatoes with milk",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 13: Season with salt and pepper
	step13ID := identifiers.New()
	step13 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step13ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           13,
		Notes:           "Season with salt and freshly ground black pepper, then serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonPotatoVIP.ID,
				IngredientID:                    &potato.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "mashed potatoes with milk",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step13ID,
				ValidIngredientPreparationID:     &seasonPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperGramVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                gramMeasurement.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step13ID,
				ValidPreparationVesselID: &seasonPotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "ultra-fluffy mashed potatoes",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		{
			ID:                  recipeID,
			CreatedByUser:       userID,
			Name:                "Ultra-Fluffy Mashed Potatoes",
			Slug:                "ultra-fluffy-mashed-potatoes",
			Source:              "https://www.seriouseats.com/ultra-fluffy-mashed-potatoes-recipe",
			Description:         "For the fluffiest and lightest mashed potatoes, use Russets and rinse off excess potato starch before and after cooking.",
			YieldsComponentType: mealplanning.MealComponentTypesSide,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 4,
				Max: pointer.To[float32](6),
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
				step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13,
			},
		},
	}
}
