package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// GochujangButterPastaRecipe creates the Gochujang Butter Pasta recipe.
func GochujangButterPastaRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	boilPrep := enums.Preparations["boil"]
	addPrep := enums.Preparations["add"]
	drainPrep := enums.Preparations["drain"]
	meltPrep := enums.Preparations["melt"]
	stirPrep := enums.Preparations["stir"]
	reducePrep := enums.Preparations["reduce"]
	removeFromHeatPrep := enums.Preparations["remove from heat"]
	mincePrep := enums.Preparations["mince"]

	// Get ingredients
	pasta := enums.Ingredients["pasta"]
	butter := enums.Ingredients["butter"]
	garlic := enums.Ingredients["garlic"]
	salt := enums.Ingredients["salt"]
	gochujangPaste := enums.Ingredients["gochujang paste"]
	honey := enums.Ingredients["honey"]
	vinegar := enums.Ingredients["vinegar"]
	water := enums.Ingredients["water"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	cloveMeasurement := enums.MeasurementUnits["clove"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	woodenSpoon := enums.Instruments["wooden spoon"]
	spoon := enums.Instruments["spoon"]
	knife := enums.Instruments["knife"]

	// Get vessels
	pot := enums.Vessels["pot"]
	skillet := enums.Vessels["pan"]
	colander := enums.Vessels["colander"]
	smallBowl := enums.Vessels["small bowl"]
	cuttingBoard := enums.Vessels["cutting board"]

	// Get bridge table entries
	boilPastaVIP := enums.IngredientPreparations[boilPrep.ID][pasta.ID]
	boilWaterVIP := enums.IngredientPreparations[boilPrep.ID][water.ID]
	boilSaltVIP := enums.IngredientPreparations[boilPrep.ID][salt.ID]
	boilPotVPV := enums.PreparationVessels[boilPrep.ID][pot.ID]
	boilWoodenSpoonVPI := enums.PreparationInstruments[boilPrep.ID][woodenSpoon.ID]

	addPastaVIP := enums.IngredientPreparations[addPrep.ID][pasta.ID]
	addPotVPV := enums.PreparationVessels[addPrep.ID][pot.ID]

	drainPastaVIP := enums.IngredientPreparations[drainPrep.ID][pasta.ID]
	drainColanderVPV := enums.PreparationVessels[drainPrep.ID][colander.ID]
	drainPotVPV := enums.PreparationVessels[drainPrep.ID][pot.ID]
	drainSmallBowlVPV := enums.PreparationVessels[drainPrep.ID][smallBowl.ID]

	meltButterVIP := enums.IngredientPreparations[meltPrep.ID][butter.ID]
	meltSkilletVPV := enums.PreparationVessels[meltPrep.ID][skillet.ID]
	meltSpoonVPI := enums.PreparationInstruments[meltPrep.ID][spoon.ID]

	addGarlicVIP := enums.IngredientPreparations[addPrep.ID][garlic.ID]
	addSaltVIP := enums.IngredientPreparations[addPrep.ID][salt.ID]
	addSpoonVPI := enums.PreparationInstruments[addPrep.ID][spoon.ID]
	addSkilletVPV := enums.PreparationVessels[addPrep.ID][skillet.ID]
	addGochujangVIP := enums.IngredientPreparations[addPrep.ID][gochujangPaste.ID]
	addHoneyVIP := enums.IngredientPreparations[addPrep.ID][honey.ID]
	addVinegarVIP := enums.IngredientPreparations[addPrep.ID][vinegar.ID]

	reduceSpoonVPI := enums.PreparationInstruments[reducePrep.ID][spoon.ID]

	stirSauceVIP := enums.IngredientPreparations[stirPrep.ID][gochujangPaste.ID]
	stirPastaVIP := enums.IngredientPreparations[stirPrep.ID][pasta.ID]
	stirButterVIP := enums.IngredientPreparations[stirPrep.ID][butter.ID]
	stirSaltVIP := enums.IngredientPreparations[stirPrep.ID][salt.ID]
	stirPotVPV := enums.PreparationVessels[stirPrep.ID][pot.ID]
	stirSpoonVPI := enums.PreparationInstruments[stirPrep.ID][spoon.ID]

	minceGarlicVIP := enums.IngredientPreparations[mincePrep.ID][garlic.ID]
	minceKnifeVPI := enums.PreparationInstruments[mincePrep.ID][knife.ID]
	minceCuttingBoardVPV := enums.PreparationVessels[mincePrep.ID][cuttingBoard.ID]

	// Measurement unit bridges
	pastaPoundVIMU := enums.IngredientMeasurementUnits[pasta.ID][poundMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
	gochujangCupVIMU := enums.IngredientMeasurementUnits[gochujangPaste.ID][cupMeasurement.ID]
	honeyCupVIMU := enums.IngredientMeasurementUnits[honey.ID][cupMeasurement.ID]
	vinegarCupVIMU := enums.IngredientMeasurementUnits[vinegar.ID][cupMeasurement.ID]
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]

	// Ingredient states
	tenderState := enums.IngredientStates["tender"]
	desiredConsistencyState := enums.IngredientStates["at desired consistency"]

	// Step 0: Mince garlic (can be done while pasta water heats)
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        mincePrep.ID,
		Index:                0,
		ExplicitInstructions: "Finely chop the garlic cloves.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &minceGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				Name:                             "garlic cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 12,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &minceKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &minceCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "finely chopped garlic",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.33),
				},
			},
		},
	}

	// Step 1: Boil water
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        boilPrep.ID,
		Index:                1,
		ExplicitInstructions: "Bring a large pot of water to a boil. Salt generously.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &boilWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				QuantityNotes:                    "enough to cover pasta",
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ValidIngredientPreparationID: &boilSaltVIP.ID,
				QuantityNotes:                "generously",
				Name:                         "salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &boilPotVPV.ID,
				Name:                     "large pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "salted boiling water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
			{
				Name:  "large pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 2: Add spaghetti
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                2,
		ExplicitInstructions: "Add the spaghetti to the boiling water.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &addPastaVIP.ID,
				ValidIngredientMeasurementUnitID: &pastaPoundVIMU.ID,
				Name:                             "spaghetti or other long pasta",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &addPotVPV.ID,
				Name:                            "large pot with salted boiling water",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "pasta in salted water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 3: Cook until al dente
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        boilPrep.ID,
		Index:                3,
		ExplicitInstructions: "Cook according to package instructions until al dente.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &boilPastaVIP.ID,
				ValidIngredientMeasurementUnitID: &pastaPoundVIMU.ID,
				Name:                             "pasta in salted water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &boilWoodenSpoonVPI.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "large pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: tenderState.ID,
				Notes:             "pasta should be al dente",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cooked al dente pasta",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 4: Drain pasta, reserving 1 cup cooking water
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        drainPrep.ID,
		Index:                4,
		ExplicitInstructions: "Before draining, scoop out 1 cup of the cooking water and set aside. Drain the spaghetti and return to its pot.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &drainPastaVIP.ID,
				ValidIngredientMeasurementUnitID: &pastaPoundVIMU.ID,
				Name:                             "cooked al dente pasta",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &drainPotVPV.ID,
				Name:                            "large pot with pasta and cooking water",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationVesselID: &drainColanderVPV.ID,
				Name:                     "colander",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationVesselID: &drainSmallBowlVPV.ID,
				Name:                     "small bowl or measuring cup",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "drained spaghetti",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:              "reserved pasta water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             1,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Return pasta to pot
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                5,
		ExplicitInstructions: "Return the drained spaghetti to its pot.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "drained spaghetti",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "large pot (empty after draining)",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "spaghetti in pot",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 6: Melt butter in skillet (while pasta cooks)
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        meltPrep.ID,
		Index:                6,
		ExplicitInstructions: "melt 4 tablespoons of the butter in a skillet over medium-low.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &meltButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				Name:                             "unsalted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &meltSpoonVPI.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &meltSkilletVPV.ID,
				Name:                     "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "melted butter",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 7: Add garlic and salt, cook until garlic softens
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                7,
		ExplicitInstructions: "Add the garlic and season generously with salt. Cook, stirring occasionally, until the garlic starts to soften but not brown, 1 to 3 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60),
			Max: pointer.To[uint32](180),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "melted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &addGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				Name:                             "finely chopped garlic",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 12,
				},
			},
			{
				ValidIngredientPreparationID: &addSaltVIP.ID,
				QuantityNotes:                "generously",
				Name:                         "salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &addSpoonVPI.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: tenderState.ID,
				Notes:             "garlic should soften but not brown",
				Ingredients:       []uint64{1},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "garlic in melted butter",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 8: Add gochujang, honey, vinegar
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                8,
		ExplicitInstructions: "Stir in the gochujang, honey and vinegar.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "garlic in melted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addGochujangVIP.ID,
				ValidIngredientMeasurementUnitID: &gochujangCupVIMU.ID,
				Name:                             "gochujang paste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ValidIngredientPreparationID:     &addHoneyVIP.ID,
				ValidIngredientMeasurementUnitID: &honeyCupVIMU.ID,
				Name:                             "honey",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ValidIngredientPreparationID:     &addVinegarVIP.ID,
				ValidIngredientMeasurementUnitID: &vinegarCupVIMU.ID,
				Name:                             "sherry vinegar or rice vinegar",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &addSpoonVPI.ID,
				Name:                         "spoon or spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &addSkilletVPV.ID,
				Name:                            "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "gochujang butter mixture",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 9: Bring to simmer and reduce
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        reducePrep.ID,
		Index:                9,
		ExplicitInstructions: "Bring to a simmer over medium-high. Cook, stirring constantly, until the mixture reduces significantly, 3 to 4 minutes; when you drag a spatula across the bottom of the pan, it should leave behind a trail that stays put for about 3 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180),
			Max: pointer.To[uint32](240),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "gochujang butter mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &reduceSpoonVPI.ID,
				Name:                         "spoon or spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: desiredConsistencyState.ID,
				Notes:             "sauce should leave a trail that stays put for about 3 seconds when dragged with spatula",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "reduced gochujang butter sauce",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 10: Remove sauce from heat
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        removeFromHeatPrep.ID,
		Index:                10,
		ExplicitInstructions: "Remove the sauce from the heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "reduced gochujang butter sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "gochujang butter sauce off heat",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 11: Transfer sauce to pasta, add remaining butter, stir and season
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                11,
		ExplicitInstructions: "Transfer the sauce to the pot with the spaghetti and add the remaining 2 tablespoons butter. Vigorously stir until the butter melts. Add splashes of the pasta cooking water, as needed, to thin out the sauce. Taste and season with salt and pepper. Top with the cilantro or scallions (if using) and serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &stirPastaVIP.ID,
				ValidIngredientMeasurementUnitID: &pastaPoundVIMU.ID,
				Name:                             "spaghetti in pot",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirSauceVIP.ID,
				Name:                            "gochujang butter sauce off heat",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &stirButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				Name:                             "remaining unsalted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "reserved pasta water",
				QuantityNotes:                   "splashes as needed to thin sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
			},
			{
				ValidIngredientPreparationID: &stirSaltVIP.ID,
				QuantityNotes:                "to taste",
				Name:                         "salt and pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				Optional: true,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &stirSpoonVPI.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &stirPotVPV.ID,
				Name:                            "large pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: desiredConsistencyState.ID,
				Notes:             "butter should be melted and sauce should coat the pasta",
				Ingredients:       []uint64{2},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "gochujang butter pasta",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	prepTask1 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
		Name:                        "Mince garlic",
		Description:                 "Finely chop the garlic cloves. Minced garlic keeps 3 to 4 days in an airtight container in the refrigerator.",
		Notes:                       "Having the garlic ready ahead of time speeds up the sauce-making step.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the minced garlic in an airtight container in the refrigerator for up to 3 days.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeAirtightContainer,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 0,
			Max: pointer.To[uint32](259200), // 3 days
		},
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 0, SatisfiesRecipeStep: true},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
			Name:                "Gochujang Butter Pasta",
			Slug:                "gochujang-butter-pasta",
			Description:         "Spaghetti with a spicy-sweet gochujang butter sauce, garlic, honey, and vinegar. Top with cilantro or scallions.",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 4,
			},
			Source:            "https://cooking.nytimes.com/recipes/1024066-gochujang-buttered-noodles",
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11,
			},
			PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{prepTask1},
			Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
		},
	}
}
