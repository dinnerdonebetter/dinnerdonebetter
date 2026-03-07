package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// OnePanPastaRecipe creates the One-Pan Pasta with Tomatoes and Greens recipe.
func OnePanPastaRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	halvePrep := enums.Preparations["halve"]
	zestPrep := enums.Preparations["zest"]
	boilPrep := enums.Preparations["boil"]
	addPrep := enums.Preparations["add"]
	coverPrep := enums.Preparations["cover"]
	simmerPrep := enums.Preparations["simmer"]
	cookPrep := enums.Preparations["cook"]
	seasonPrep := enums.Preparations["season"]
	topPrep := enums.Preparations["top"]

	// Get ingredients
	tomato := enums.Ingredients["tomato"]
	lemon := enums.Ingredients["lemon"]
	water := enums.Ingredients["water"]
	pasta := enums.Ingredients["pasta"]
	oliveOil := enums.Ingredients["olive oil"]
	salt := enums.Ingredients["salt"]
	kale := enums.Ingredients["kale"]
	spinach := enums.Ingredients["spinach"]
	blackPepper := enums.Ingredients["black pepper"]
	parmesan := enums.Ingredients["parmesan cheese"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	quartMeasurement := enums.MeasurementUnits["quart"]
	cupMeasurement := enums.MeasurementUnits["cup"]

	// Get instruments
	microplane := enums.Instruments["microplane"]
	tongs := enums.Instruments["tongs"]
	knife := enums.Instruments["knife"]

	// Get vessels
	pot := enums.Vessels["pot"]
	sautePan := enums.Vessels["sauté pan"]
	cuttingBoard := enums.Vessels["cutting board"]

	// Get bridge table entries
	halveTomatoVIP := enums.IngredientPreparations[halvePrep.ID][tomato.ID]
	halveKnifeVPI := enums.PreparationInstruments[halvePrep.ID][knife.ID]
	halveCuttingBoardVPV := enums.PreparationVessels[halvePrep.ID][cuttingBoard.ID]

	zestLemonVIP := enums.IngredientPreparations[zestPrep.ID][lemon.ID]
	zestMicroplaneVPI := enums.PreparationInstruments[zestPrep.ID][microplane.ID]

	boilWaterVIP := enums.IngredientPreparations[boilPrep.ID][water.ID]
	boilPotVPV := enums.PreparationVessels[boilPrep.ID][pot.ID]

	addPastaVIP := enums.IngredientPreparations[addPrep.ID][pasta.ID]
	addTomatoVIP := enums.IngredientPreparations[addPrep.ID][tomato.ID]
	addOliveOilVIP := enums.IngredientPreparations[addPrep.ID][oliveOil.ID]
	addSaltVIP := enums.IngredientPreparations[addPrep.ID][salt.ID]
	addWaterVIP := enums.IngredientPreparations[addPrep.ID][water.ID]
	addKaleVIP := enums.IngredientPreparations[addPrep.ID][kale.ID]
	addSpinachVIP := enums.IngredientPreparations[addPrep.ID][spinach.ID]
	addSautePanVPV := enums.PreparationVessels[addPrep.ID][sautePan.ID]

	coverSautePanVPV := enums.PreparationVessels[coverPrep.ID][sautePan.ID]

	simmerPastaVIP := enums.IngredientPreparations[simmerPrep.ID][pasta.ID]
	simmerTongsVPI := enums.PreparationInstruments[simmerPrep.ID][tongs.ID]
	simmerSautePanVPV := enums.PreparationVessels[simmerPrep.ID][sautePan.ID]

	cookPastaVIP := enums.IngredientPreparations[cookPrep.ID][pasta.ID]
	cookTongsVPI := enums.PreparationInstruments[cookPrep.ID][tongs.ID]
	cookSautePanVPV := enums.PreparationVessels[cookPrep.ID][sautePan.ID]

	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	seasonSautePanVPV := enums.PreparationVessels[seasonPrep.ID][sautePan.ID]

	topParmesanVIP := enums.IngredientPreparations[topPrep.ID][parmesan.ID]
	topSautePanVPV := enums.PreparationVessels[topPrep.ID][sautePan.ID]

	// Measurement unit bridges (2 pints cherry tomatoes ≈ 4 cups)
	tomatoCupVIMU := enums.IngredientMeasurementUnits[tomato.ID][cupMeasurement.ID]
	lemonUnitVIMU := enums.IngredientMeasurementUnits[lemon.ID][unitMeasurement.ID]
	waterQuartVIMU := enums.IngredientMeasurementUnits[water.ID][quartMeasurement.ID]
	pastaPoundVIMU := enums.IngredientMeasurementUnits[pasta.ID][poundMeasurement.ID]
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	kaleCupVIMU := enums.IngredientMeasurementUnits[kale.ID][cupMeasurement.ID]
	spinachCupVIMU := enums.IngredientMeasurementUnits[spinach.ID][cupMeasurement.ID]

	// Ingredient states
	tenderState := enums.IngredientStates["tender"]
	desiredConsistencyState := enums.IngredientStates["at desired consistency"]

	// Step 0: Halve cherry tomatoes
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        halvePrep.ID,
		Index:                0,
		ExplicitInstructions: "Halve the cherry tomatoes.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &halveTomatoVIP.ID,
				ValidIngredientMeasurementUnitID: &tomatoCupVIMU.ID,
				Name:                             "cherry tomatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &halveKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &halveCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "halved cherry tomatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 1: Zest the lemons
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        zestPrep.ID,
		Index:                1,
		ExplicitInstructions: "Zest the lemons.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &zestLemonVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonUnitVIMU.ID,
				Name:                             "lemons",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &zestMicroplaneVPI.ID,
				Name:                         "microplane",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "lemon zest",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Bring water to a boil
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        boilPrep.ID,
		Index:                2,
		ExplicitInstructions: "Bring just over a quart of water to a boil in a pot.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &boilWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterQuartVIMU.ID,
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &boilPotVPV.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "boiling water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &quartMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1.25),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 3: Place spaghetti, tomatoes, lemon zest, oil, and salt in the pan
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                3,
		ExplicitInstructions: "Place spaghetti, halved tomatoes, lemon zest, olive oil, and 2 teaspoons kosher salt in a large, dry, shallow pan. The pan should be large enough that the dry spaghetti can lie flat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &addPastaVIP.ID,
				ValidIngredientMeasurementUnitID: &pastaPoundVIMU.ID,
				Name:                             "spaghetti",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addTomatoVIP.ID,
				Name:                            "halved cherry tomatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "lemon zest",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				Name:                             "olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 7,
				},
			},
			{
				ValidIngredientPreparationID:     &addSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &addSautePanVPV.ID,
				Name:                     "large shallow pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "dry spaghetti with tomatoes, zest, oil, and salt",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large shallow pan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 4: Add boiling water to the pan
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                4,
		ExplicitInstructions: "Carefully add the boiling water to the pan with the spaghetti.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addWaterVIP.ID,
				Name:                            "boiling water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "large shallow pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "spaghetti in water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large shallow pan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 5: Cover and bring to a boil
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        coverPrep.ID,
		Index:                5,
		ExplicitInstructions: "Cover the pan and bring to a boil.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "spaghetti in water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &coverSautePanVPV.ID,
				Name:                            "large shallow pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "covered spaghetti at boil",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large shallow pan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 6: Uncover and simmer, stirring with tongs
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        simmerPrep.ID,
		Index:                6,
		ExplicitInstructions: "Remove the lid and simmer for about 6 minutes, using tongs to move the spaghetti around now and then so it doesn't stick.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](360),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &simmerPastaVIP.ID,
				Name:                            "covered spaghetti at boil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &simmerTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &simmerSautePanVPV.ID,
				Name:                            "large shallow pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "simmered spaghetti",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large shallow pan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 7: Add kale or spinach
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                7,
		ExplicitInstructions: "Add kale or spinach (leaves only, washed and chopped) and continue cooking.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "simmered spaghetti",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addKaleVIP.ID,
				ValidIngredientMeasurementUnitID: &kaleCupVIMU.ID,
				Name:                             "kale (leaves only, washed and chopped)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
			},
			{
				ValidIngredientPreparationID:     &addSpinachVIP.ID,
				ValidIngredientMeasurementUnitID: &spinachCupVIMU.ID,
				Name:                             "spinach (leaves only, washed and chopped)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 1,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "large shallow pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "spaghetti with greens",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large shallow pan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 8: Cook until liquid reduces to sauce and pasta is tender
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        cookPrep.ID,
		Index:                8,
		ExplicitInstructions: "Continue cooking until the remaining liquid has reduced to a sauce and the pasta is cooked through.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &cookPastaVIP.ID,
				Name:                            "spaghetti with greens",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &cookTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &cookSautePanVPV.ID,
				Name:                            "large shallow pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: tenderState.ID,
				Notes:             "pasta should be cooked through",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
			{
				IngredientStateID: desiredConsistencyState.ID,
				Notes:             "liquid should have reduced to a sauce",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "one-pan pasta",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large shallow pan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 9: Season and top with Parmesan
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        seasonPrep.ID,
		Index:                9,
		ExplicitInstructions: "Taste, season with salt and pepper to taste, and top with Parmesan.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "one-pan pasta",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID: &seasonSaltVIP.ID,
				QuantityNotes:                "to taste",
				Name:                         "salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
			},
			{
				ValidIngredientPreparationID: &seasonPepperVIP.ID,
				QuantityNotes:                "to taste",
				Name:                         "black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &seasonSautePanVPV.ID,
				Name:                            "large shallow pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seasoned one-pan pasta",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "large shallow pan",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 10: Top with Parmesan
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        topPrep.ID,
		Index:                10,
		ExplicitInstructions: "Top with Parmesan for serving.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "seasoned one-pan pasta",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID: &topParmesanVIP.ID,
				QuantityNotes:                "for serving",
				Name:                         "Parmesan",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &topSautePanVPV.ID,
				Name:                            "large shallow pan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "one-pan pasta with Parmesan",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
			Name:                "One-Pan Pasta with Tomatoes and Greens",
			Slug:                "one-pan-pasta-tomatoes-greens",
			Source:              "https://cooking.nytimes.com/recipes/1018322-one-pot-spaghetti-with-cherry-tomatoes-and-kale",
			Description:         "A simple one-pan pasta where spaghetti cooks directly with cherry tomatoes, lemon zest, and olive oil. Add kale or spinach near the end and top with Parmesan.",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 4,
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10,
			},
			PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
		},
	}
}
