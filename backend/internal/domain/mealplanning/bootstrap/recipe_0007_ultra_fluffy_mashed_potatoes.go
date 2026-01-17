package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// UltraFluffyMashedPotatoesRecipe creates the Ultra-Fluffy Mashed Potatoes recipe.
// Source: https://www.seriouseats.com/ultra-fluffy-mashed-potatoes-recipe
func UltraFluffyMashedPotatoesRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
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
	slicePrep := enums.Preparations["slice"]
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

	// Season preparation bridges for pot (seasoning water)
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]

	// Boil preparation bridges
	boilPotatoVIP := enums.IngredientPreparations[boilPrep.ID][potato.ID]

	// Drain preparation bridges
	drainPotatoVIP := enums.IngredientPreparations[drainPrep.ID][potato.ID]
	drainColanderVPV := enums.PreparationVessels[drainPrep.ID][colander.ID]

	// Rest preparation bridges
	restPotatoVIP := enums.IngredientPreparations[restPrep.ID][potato.ID]
	restColanderVPV := enums.PreparationVessels[restPrep.ID][colander.ID]

	// Rice preparation bridges
	ricePotatoVIP := enums.IngredientPreparations[ricePrep.ID][potato.ID]
	ricePotatoRicerVPI := enums.PreparationInstruments[ricePrep.ID][potatoRicer.ID]

	// Slice preparation bridges
	sliceButterVIP := enums.IngredientPreparations[slicePrep.ID][butter.ID]
	sliceKnifeVPI := enums.PreparationInstruments[slicePrep.ID][knife.ID]
	sliceCuttingBoardVPV := enums.PreparationVessels[slicePrep.ID][cuttingBoard.ID]

	// Fold preparation bridges
	foldPotatoVIP := enums.IngredientPreparations[foldPrep.ID][potato.ID]
	foldRubberSpatulaVPI := enums.PreparationInstruments[foldPrep.ID][rubberSpatula.ID]

	// Season preparation bridges (for final seasoning step)
	seasonPotatoVIP := enums.IngredientPreparations[seasonPrep.ID][potato.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]

	// Simmer preparation bridges
	simmerMilkVIP := enums.IngredientPreparations[simmerPrep.ID][milk.ID]

	// Measurement unit bridges
	potatoPoundVIMU := enums.IngredientMeasurementUnits[potato.ID][poundMeasurement.ID]
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	milkCupVIMU := enums.IngredientMeasurementUnits[milk.ID][cupMeasurement.ID]
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]

	// Step 0: Peel potatoes
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: peelPrep.ID,
		Index:         0,
		Notes:         "Peel the russet potatoes.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &peelPotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				Name:                             "russet potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &peelVegetablePeelerVPI.ID,
				Name:                         "vegetable peeler",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &peelCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "peeled potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 1: Cut potatoes into 1-2 inch cubes
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: cubePrep.ID,
		Index:         1,
		Notes:         "Cut the peeled potatoes into 1- or 2-inch cubes.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &cubePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				Name:                             "peeled potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &cubeKnifeVPI.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &cubeCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cubed potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 2: Rinse potatoes in pot of cold water
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: rinsePrep.ID,
		Index:         2,
		Notes:         "Transfer potatoes to a pot of cold water and rinse, changing the water 2 or 3 times until it runs clear.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &rinsePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				Name:                             "cubed potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &rinsePotVPV.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "rinsed potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 3: Cover potatoes with fresh cold water
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: submergePrep.ID,
		Index:         3,
		Notes:         "Cover potatoes with fresh cold water.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &submergePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				Name:                             "rinsed potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &submergeWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				QuantityNotes:                    "enough to cover potatoes",
				Name:                             "cold water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8, // enough to cover potatoes
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "potatoes in water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 4: Season the water generously with salt
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         4,
		Notes:         "Season the water generously with salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "potatoes in water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltGramVIMU.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 15, // generous amount
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "potatoes in salted water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 5: Boil and simmer until tender
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: boilPrep.ID,
		Index:         5,
		Notes:         "Set over medium-high heat and bring to a boil, then reduce heat to maintain a gentle simmer. Cook until potatoes are completely tender, about 15 minutes after reaching a simmer.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](900), // 15 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &boilPotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				Name:                             "potatoes in salted water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "boiled potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 6: Drain potatoes in a colander
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: drainPrep.ID,
		Index:         6,
		Notes:         "Drain potatoes in a colander.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &drainPotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				Name:                             "boiled potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &drainColanderVPV.ID,
				Name:                     "colander",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "drained potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 7: Rinse under hot running water
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: rinsePrep.ID,
		Index:         7,
		Notes:         "Rinse potatoes under hot running water for 30 seconds to wash away excess starch.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &rinsePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				Name:                             "drained potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &rinseColanderVPV.ID,
				Name:                     "colander",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "rinsed hot potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 8: Allow potatoes to steam/rest
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: restPrep.ID,
		Index:         8,
		Notes:         "Allow potatoes to steam for 1 minute to remove excess moisture.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60), // 1 minute
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &restPotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				Name:                             "rinsed hot potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &restColanderVPV.ID,
				Name:                     "colander",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "steamed potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 9: Pass potatoes through ricer
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: ricePrep.ID,
		Index:         9,
		Notes:         "Set a ricer or food mill over the now-empty pot and pass potatoes through.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &ricePotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoPoundVIMU.ID,
				Name:                             "steamed potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &ricePotatoRicerVPI.ID,
				Name:                         "potato ricer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "riced potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 10: Cut butter into 1/2-inch pats
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: slicePrep.ID,
		Index:         10,
		Notes:         "Cut room temperature unsalted butter into 1/2-inch pats.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &sliceButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				Name:                             "unsalted butter, room temperature",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6, // 6 tablespoons
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &sliceKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &sliceCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "butter pats (1/2-inch)",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &tablespoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](6),
				},
			},
		},
	}

	// Step 11: Fold in butter
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: foldPrep.ID,
		Index:         11,
		Notes:         "Add butter pats and gently fold into potatoes.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "riced potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "butter pats (1/2-inch)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6, // 6 tablespoons
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &foldRubberSpatulaVPI.ID,
				Name:                         "rubber spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "buttered potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 12: Simmer milk
	step12 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: simmerPrep.ID,
		Index:         12,
		Notes:         "Mound potatoes into the center of the pot and pour milk all around. Set over medium heat and bring milk to a simmer.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "buttered potatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &simmerMilkVIP.ID,
				ValidIngredientMeasurementUnitID: &milkCupVIMU.ID,
				Name:                             "whole milk",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5, // 1/2 cup, plus more as needed
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "potatoes with simmering milk",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 13: Fold simmered milk into potatoes
	step13 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: foldPrep.ID,
		Index:         13,
		Notes:         "Gently fold the simmered milk into the potatoes. If looser potatoes are desired, add additional milk in a similar fashion around the mashed potato mass and bring it to a simmer before folding into potatoes.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &foldPotatoVIP.ID,
				Name:                            "potatoes with simmering milk",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &foldRubberSpatulaVPI.ID,
				Name:                         "rubber spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "mashed potatoes with milk",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 14: Season with salt and pepper
	step14 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         14,
		Notes:         "Season with salt and freshly ground black pepper, then serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonPotatoVIP.ID,
				Name:                            "mashed potatoes with milk",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperGramVIMU.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "ultra-fluffy mashed potatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Prep task: Prepare potatoes up to 24 hours ahead
	prepTask1 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
		Name:                        "Prepare and soak potatoes",
		Description:                 "Steps 1-4 (cube, rinse, submerge, and season potatoes) can be done up to 24 hours ahead of time if left submerged in water in the refrigerator.",
		Notes:                       "Preparing the potatoes ahead saves time on the day of cooking.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the cubed and seasoned potatoes in the pot, covered with water, in the refrigerator for up to 24 hours.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeCovered,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4), // Refrigerator temperature
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 0,
			Max: pointer.To[uint32](86400), // 24 hours
		},
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 1, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 2, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 3, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 4, SatisfiesRecipeStep: true},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
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
			Steps:             []*mealplanning.RecipeStepCreationRequestInput{step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13, step14},
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{prepTask1},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
