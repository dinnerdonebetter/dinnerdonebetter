package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ChanaMasalaRecipe creates the Chana Masala recipe.
// Source: https://cooking.nytimes.com/recipes/1024429-chana-masala
func ChanaMasalaRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	meltPrep := enums.Preparations["melt"]
	addPrep := enums.Preparations["add"]
	simmerPrep := enums.Preparations["simmer"]
	smashPrep := enums.Preparations["smash"]
	topPrep := enums.Preparations["top"]
	dicePrep := enums.Preparations["dice"]
	mincePrep := enums.Preparations["mince"]
	chopPrep := enums.Preparations["chop"]

	// Get ingredients
	ghee := enums.Ingredients["ghee"]
	garlic := enums.Ingredients["garlic"]
	ginger := enums.Ingredients["ginger"]
	onion := enums.Ingredients["onion"]
	freshThaiChile := enums.Ingredients["fresh Thai chile"]
	cuminSeeds := enums.Ingredients["cumin seeds"]
	turmeric := enums.Ingredients["turmeric"]
	groundCoriander := enums.Ingredients["ground coriander"]
	chiliPowder := enums.Ingredients["chili powder"]
	tomato := enums.Ingredients["tomato"]
	salt := enums.Ingredients["salt"]
	chickpeas := enums.Ingredients["chickpeas"]
	chickenStock := enums.Ingredients["chicken stock"]
	garamMasala := enums.Ingredients["garam masala"]
	cilantro := enums.Ingredients["cilantro"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	cloveMeasurement := enums.MeasurementUnits["clove"]

	// Get instruments
	woodenSpoon := enums.Instruments["wooden spoon"]
	spoon := enums.Instruments["spoon"]
	knife := enums.Instruments["knife"]

	// Get vessels
	pot := enums.Vessels["pot"]
	cuttingBoard := enums.Vessels["cutting board"]

	// Get bridge table entries
	meltGheeVIP := enums.IngredientPreparations[meltPrep.ID][ghee.ID]
	meltPotVPV := enums.PreparationVessels[meltPrep.ID][pot.ID]
	meltSpoonVPI := enums.PreparationInstruments[meltPrep.ID][spoon.ID]

	addGarlicVIP := enums.IngredientPreparations[addPrep.ID][garlic.ID]
	addOnionVIP := enums.IngredientPreparations[addPrep.ID][onion.ID]
	addChilesVIP := enums.IngredientPreparations[addPrep.ID][freshThaiChile.ID]
	addCuminVIP := enums.IngredientPreparations[addPrep.ID][cuminSeeds.ID]
	addTurmericVIP := enums.IngredientPreparations[addPrep.ID][turmeric.ID]
	addCorianderVIP := enums.IngredientPreparations[addPrep.ID][groundCoriander.ID]
	addChiliPowderVIP := enums.IngredientPreparations[addPrep.ID][chiliPowder.ID]
	addTomatoVIP := enums.IngredientPreparations[addPrep.ID][tomato.ID]
	addSaltVIP := enums.IngredientPreparations[addPrep.ID][salt.ID]
	addChickpeasVIP := enums.IngredientPreparations[addPrep.ID][chickpeas.ID]
	addStockVIP := enums.IngredientPreparations[addPrep.ID][chickenStock.ID]
	addPotVPV := enums.PreparationVessels[addPrep.ID][pot.ID]
	addWoodenSpoonVPI := enums.PreparationInstruments[addPrep.ID][woodenSpoon.ID]

	simmerChickpeasVIP := enums.IngredientPreparations[simmerPrep.ID][chickpeas.ID]
	simmerPotVPV := enums.PreparationVessels[simmerPrep.ID][pot.ID]

	smashChickpeasVIP := enums.IngredientPreparations[smashPrep.ID][chickpeas.ID]
	smashPotVPV := enums.PreparationVessels[smashPrep.ID][pot.ID]
	smashSpoonVPI := enums.PreparationInstruments[smashPrep.ID][spoon.ID]

	topGaramMasalaVIP := enums.IngredientPreparations[topPrep.ID][garamMasala.ID]
	topCilantroVIP := enums.IngredientPreparations[topPrep.ID][cilantro.ID]
	topGingerVIP := enums.IngredientPreparations[topPrep.ID][ginger.ID]
	topPotVPV := enums.PreparationVessels[topPrep.ID][pot.ID]

	diceOnionVIP := enums.IngredientPreparations[dicePrep.ID][onion.ID]
	chopTomatoVIP := enums.IngredientPreparations[chopPrep.ID][tomato.ID]
	diceKnifeVPI := enums.PreparationInstruments[dicePrep.ID][knife.ID]
	diceCuttingBoardVPV := enums.PreparationVessels[dicePrep.ID][cuttingBoard.ID]

	minceGarlicVIP := enums.IngredientPreparations[mincePrep.ID][garlic.ID]
	minceGingerVIP := enums.IngredientPreparations[mincePrep.ID][ginger.ID]
	minceKnifeVPI := enums.PreparationInstruments[mincePrep.ID][knife.ID]
	minceCuttingBoardVPV := enums.PreparationVessels[mincePrep.ID][cuttingBoard.ID]

	chopChilesVIP := enums.IngredientPreparations[chopPrep.ID][freshThaiChile.ID]
	chopCilantroVIP := enums.IngredientPreparations[chopPrep.ID][cilantro.ID]
	chopKnifeVPI := enums.PreparationInstruments[chopPrep.ID][knife.ID]
	chopCuttingBoardVPV := enums.PreparationVessels[chopPrep.ID][cuttingBoard.ID]

	// Measurement unit bridges
	gheeTbspVIMU := enums.IngredientMeasurementUnits[ghee.ID][tablespoonMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
	gingerUnitVIMU := enums.IngredientMeasurementUnits[ginger.ID][unitMeasurement.ID]
	onionUnitVIMU := enums.IngredientMeasurementUnits[onion.ID][unitMeasurement.ID]
	chileUnitVIMU := enums.IngredientMeasurementUnits[freshThaiChile.ID][unitMeasurement.ID]
	cuminTspVIMU := enums.IngredientMeasurementUnits[cuminSeeds.ID][teaspoonMeasurement.ID]
	turmericTspVIMU := enums.IngredientMeasurementUnits[turmeric.ID][teaspoonMeasurement.ID]
	corianderTspVIMU := enums.IngredientMeasurementUnits[groundCoriander.ID][teaspoonMeasurement.ID]
	chiliPowderTspVIMU := enums.IngredientMeasurementUnits[chiliPowder.ID][teaspoonMeasurement.ID]
	tomatoUnitVIMU := enums.IngredientMeasurementUnits[tomato.ID][unitMeasurement.ID]
	saltTspVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	chickpeasCupVIMU := enums.IngredientMeasurementUnits[chickpeas.ID][cupMeasurement.ID]
	stockCupVIMU := enums.IngredientMeasurementUnits[chickenStock.ID][cupMeasurement.ID]
	garamMasalaTspVIMU := enums.IngredientMeasurementUnits[garamMasala.ID][teaspoonMeasurement.ID]
	cilantroTbspVIMU := enums.IngredientMeasurementUnits[cilantro.ID][tablespoonMeasurement.ID]

	// Ingredient states
	translucentState := enums.IngredientStates["translucent"]
	desiredConsistencyState := enums.IngredientStates["at desired consistency"]

	// Step 0: Dice the onion
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        dicePrep.ID,
		Index:                0,
		ExplicitInstructions: "Peel and finely chop the medium red onion.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &diceOnionVIP.ID,
				ValidIngredientMeasurementUnitID: &onionUnitVIMU.ID,
				Name:                             "medium red onion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &diceKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &diceCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "finely chopped onion",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Mince garlic and ginger
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        mincePrep.ID,
		Index:                1,
		ExplicitInstructions: "Mince the garlic and grate or mince the ginger (from a peeled 2-inch piece) to yield about 1 tablespoon each.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &minceGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				Name:                             "garlic cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ValidIngredientPreparationID:     &minceGingerVIP.ID,
				ValidIngredientMeasurementUnitID: &gingerUnitVIMU.ID,
				Name:                             "fresh ginger (2-inch piece)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
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
				Name:              "minced garlic and ginger",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &tablespoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 2: Chop tomatoes
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        chopPrep.ID,
		Index:                2,
		ExplicitInstructions: "Finely chop the Roma tomatoes.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &chopTomatoVIP.ID,
				ValidIngredientMeasurementUnitID: &tomatoUnitVIMU.ID,
				Name:                             "Roma tomatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &chopKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &chopCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "finely chopped tomatoes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 2b: Chop chiles
	step2b := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        chopPrep.ID,
		Index:                3,
		ExplicitInstructions: "Chop the Thai green or bird's eye chiles.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &chopChilesVIP.ID,
				ValidIngredientMeasurementUnitID: &chileUnitVIMU.ID,
				Name:                             "Thai green or bird's eye chiles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &chopKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &chopCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "chopped green chiles",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 4: Chop cilantro for garnish
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        chopPrep.ID,
		Index:                4,
		ExplicitInstructions: "Chop the cilantro leaves and tender stems.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &chopCilantroVIP.ID,
				ValidIngredientMeasurementUnitID: &cilantroTbspVIMU.ID,
				Name:                             "cilantro leaves and tender stems",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &chopKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &chopCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "chopped cilantro",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &tablespoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 5: Melt ghee in pot over medium heat
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        meltPrep.ID,
		Index:                5,
		ExplicitInstructions: "In a medium pot, melt the ghee over medium heat. (Alternatively, use 2 tablespoons neutral oil such as vegetable oil.)",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &meltGheeVIP.ID,
				ValidIngredientMeasurementUnitID: &gheeTbspVIMU.ID,
				Name:                             "ghee or neutral oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
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
				ValidPreparationVesselID: &meltPotVPV.ID,
				Name:                     "medium pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "melted ghee in pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Add garlic, ginger, and onion; cook until onion softens
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                6,
		ExplicitInstructions: "Stir in the garlic, ginger, and onion. Continue cooking, stirring occasionally, until the onion softens, 5 to 7 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300),
			Max: pointer.To[uint32](420),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addGarlicVIP.ID,
				Name:                            "minced garlic and ginger",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addOnionVIP.ID,
				Name:                            "finely chopped onion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &addWoodenSpoonVPI.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &addPotVPV.ID,
				Name:                            "pot with melted ghee",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: translucentState.ID,
				Notes:             "onion should be softened and translucent",
				Ingredients:       []uint64{1},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cooked aromatics",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot with cooked aromatics",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Add chiles and spices; stir for 30 seconds
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                7,
		ExplicitInstructions: "Stir in the green chiles, cumin seeds, turmeric, ground coriander, and chile powder. Continue stirring for 30 seconds so the spices don't burn.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30),
			Max: pointer.To[uint32](30),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addChilesVIP.ID,
				Name:                            "chopped green chiles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addCuminVIP.ID,
				ValidIngredientMeasurementUnitID: &cuminTspVIMU.ID,
				Name:                             "cumin seeds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addTurmericVIP.ID,
				ValidIngredientMeasurementUnitID: &turmericTspVIMU.ID,
				Name:                             "ground turmeric",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ValidIngredientPreparationID:     &addCorianderVIP.ID,
				ValidIngredientMeasurementUnitID: &corianderTspVIMU.ID,
				Name:                             "ground coriander",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ValidIngredientPreparationID:     &addChiliPowderVIP.ID,
				ValidIngredientMeasurementUnitID: &chiliPowderTspVIMU.ID,
				Name:                             "Kashmiri or other hot red chile powder",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &addWoodenSpoonVPI.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &addPotVPV.ID,
				Name:                            "pot with cooked aromatics",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "spiced aromatics",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot with spiced aromatics",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Add tomatoes and salt; cook until jammy
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                8,
		ExplicitInstructions: "Add the tomatoes and their juices and salt. Increase the heat to high and cook, stirring often, until the mixture is jammy, 5 to 7 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300),
			Max: pointer.To[uint32](420),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addTomatoVIP.ID,
				Name:                            "finely chopped tomatoes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTspVIMU.ID,
				Name:                             "fine sea salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &addWoodenSpoonVPI.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &addPotVPV.ID,
				Name:                            "pot with spiced aromatics",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: desiredConsistencyState.ID,
				Notes:             "tomato mixture should be thick and jammy",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "jammy tomato masala",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot with jammy tomato masala",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Add chickpeas and stock; bring to boil
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                9,
		ExplicitInstructions: "Stir in the chickpeas and stock (or water). Bring to a boil.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &addChickpeasVIP.ID,
				ValidIngredientMeasurementUnitID: &chickpeasCupVIMU.ID,
				Name:                             "drained chickpeas (2 (15-ounce) cans or 3 cups cooked)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ValidIngredientPreparationID:     &addStockVIP.ID,
				ValidIngredientMeasurementUnitID: &stockCupVIMU.ID,
				Name:                             "unsalted chicken or vegetable stock, or water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &addWoodenSpoonVPI.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &addPotVPV.ID,
				Name:                            "pot with jammy tomato masala",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "chickpeas in masala",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](5),
				},
			},
			{
				Name:  "pot with chickpeas and stock",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Bring to boil, then reduce heat and simmer until thickened
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        simmerPrep.ID,
		Index:                10,
		ExplicitInstructions: "Bring to a boil, then reduce heat and simmer until the mixture has thickened slightly, 5 to 7 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300),
			Max: pointer.To[uint32](420),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &simmerChickpeasVIP.ID,
				Name:                            "chickpeas in masala",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &simmerPotVPV.ID,
				Name:                            "pot with chickpeas and stock",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with simmering chickpeas",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 11: Smash chickpeas to thicken
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        smashPrep.ID,
		Index:                11,
		ExplicitInstructions: "With the back of a spoon, smash some of the chickpeas against the inside of the pot to thicken the mixture; continue smashing until it reaches the desired thickness.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &smashChickpeasVIP.ID,
				Name:                            "chickpeas in masala",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &smashSpoonVPI.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &smashPotVPV.ID,
				Name:                            "pot with simmering chickpeas",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: desiredConsistencyState.ID,
				Notes:             "mixture should reach desired thickness",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "thickened chana masala",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot with chana masala",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 12: Sprinkle with garam masala and top with cilantro and ginger
	step12 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        topPrep.ID,
		Index:                12,
		ExplicitInstructions: "Sprinkle with garam masala and top with cilantro and ginger. If desired, serve with rice or roti and lemon wedges alongside.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "thickened chana masala",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &topGaramMasalaVIP.ID,
				ValidIngredientMeasurementUnitID: &garamMasalaTspVIMU.ID,
				Name:                             "garam masala",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &topCilantroVIP.ID,
				Name:                            "chopped cilantro",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &topGingerVIP.ID,
				ValidIngredientMeasurementUnitID: &gingerUnitVIMU.ID,
				Name:                             "fresh ginger, peeled and sliced into matchsticks (for serving)",
				QuantityNotes:                    "optional",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &topPotVPV.ID,
				Name:                            "pot with chana masala",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "chana masala",
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
			Name:                "Chana Masala",
			Slug:                "chana-masala",
			Description:         "A classic Indian chickpea curry with a fragrant tomato-based masala, warm spices, and fresh cilantro. Serve with rice or roti.",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 4,
			},
			Source:            "https://cooking.nytimes.com/recipes/1024429-chana-masala",
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				step0, step1, step2, step2b, step4, step5, step6, step7, step8, step9, step10, step11, step12,
			},
			PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
		},
	}
}
