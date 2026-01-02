package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// MixedGreenSaladRecipe creates the Mixed Green Salad (Misticanza alla Romana) recipe.
// Source: https://www.seriouseats.com/roman-mixed-green-salad-misticanza-recipe
func MixedGreenSaladRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	trimPrep := enums.Preparations["trim"]
	slicePrep := enums.Preparations["slice"]
	rinsePrep := enums.Preparations["rinse"]
	dryPrep := enums.Preparations["dry"]
	mixPrep := enums.Preparations["mix"]
	tossPrep := enums.Preparations["toss"]
	seasonPrep := enums.Preparations["season"]

	// Get ingredients
	lettuce := enums.Ingredients["lettuce"]
	radicchio := enums.Ingredients["radicchio"]
	endive := enums.Ingredients["endive"]
	frisee := enums.Ingredients["frisée"]
	kale := enums.Ingredients["kale"]
	dandelionGreens := enums.Ingredients["dandelion greens"]
	purslane := enums.Ingredients["purslane"]
	fennelFronds := enums.Ingredients["fennel fronds"]
	parsley := enums.Ingredients["parsley"]
	tarragon := enums.Ingredients["tarragon"]
	chervil := enums.Ingredients["chervil"]
	basil := enums.Ingredients["basil"]
	mint := enums.Ingredients["mint"]
	oliveOil := enums.Ingredients["olive oil"]
	lemon := enums.Ingredients["lemon"]
	salt := enums.Ingredients["salt"]

	// Get measurement units
	cupMeasurement := enums.MeasurementUnits["cup"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	knife := enums.Instruments["knife"]
	bareHands := enums.Instruments["bare hands"]

	// Get vessels
	cuttingBoard := enums.Vessels["cutting board"]
	largeBowl := enums.Vessels["large bowl"]
	saladSpinner := enums.Vessels["salad spinner"]
	servingBowl := enums.Vessels["serving bowl"]

	// Get bridge table entries
	// Trim
	trimLettuceVIP := enums.IngredientPreparations[trimPrep.ID][lettuce.ID]
	trimRadicchioVIP := enums.IngredientPreparations[trimPrep.ID][radicchio.ID]
	trimEndiveVIP := enums.IngredientPreparations[trimPrep.ID][endive.ID]
	trimFriseeVIP := enums.IngredientPreparations[trimPrep.ID][frisee.ID]
	trimKaleVIP := enums.IngredientPreparations[trimPrep.ID][kale.ID]
	trimDandelionGreensVIP := enums.IngredientPreparations[trimPrep.ID][dandelionGreens.ID]
	trimPurslaneVIP := enums.IngredientPreparations[trimPrep.ID][purslane.ID]
	trimFennelFrondsVIP := enums.IngredientPreparations[trimPrep.ID][fennelFronds.ID]
	trimParsleyVIP := enums.IngredientPreparations[trimPrep.ID][parsley.ID]
	trimTarragonVIP := enums.IngredientPreparations[trimPrep.ID][tarragon.ID]
	trimChervilVIP := enums.IngredientPreparations[trimPrep.ID][chervil.ID]
	trimBasilVIP := enums.IngredientPreparations[trimPrep.ID][basil.ID]
	trimMintVIP := enums.IngredientPreparations[trimPrep.ID][mint.ID]
	trimCuttingBoardVPV := enums.PreparationVessels[trimPrep.ID][cuttingBoard.ID]
	trimBareHandsVPI := enums.PreparationInstruments[trimPrep.ID][bareHands.ID]

	// Slice
	sliceLettuceVIP := enums.IngredientPreparations[slicePrep.ID][lettuce.ID]
	sliceRadicchioVIP := enums.IngredientPreparations[slicePrep.ID][radicchio.ID]
	sliceEndiveVIP := enums.IngredientPreparations[slicePrep.ID][endive.ID]
	sliceCuttingBoardVPV := enums.PreparationVessels[slicePrep.ID][cuttingBoard.ID]
	sliceKnifeVPI := enums.PreparationInstruments[slicePrep.ID][knife.ID]

	// Rinse
	rinseLettuceVIP := enums.IngredientPreparations[rinsePrep.ID][lettuce.ID]
	rinseRadicchioVIP := enums.IngredientPreparations[rinsePrep.ID][radicchio.ID]
	rinseEndiveVIP := enums.IngredientPreparations[rinsePrep.ID][endive.ID]
	rinseFriseeVIP := enums.IngredientPreparations[rinsePrep.ID][frisee.ID]
	rinseKaleVIP := enums.IngredientPreparations[rinsePrep.ID][kale.ID]
	rinseDandelionGreensVIP := enums.IngredientPreparations[rinsePrep.ID][dandelionGreens.ID]
	rinsePurslaneVIP := enums.IngredientPreparations[rinsePrep.ID][purslane.ID]
	rinseFennelFrondsVIP := enums.IngredientPreparations[rinsePrep.ID][fennelFronds.ID]
	rinseParsleyVIP := enums.IngredientPreparations[rinsePrep.ID][parsley.ID]
	rinseTarragonVIP := enums.IngredientPreparations[rinsePrep.ID][tarragon.ID]
	rinseChervilVIP := enums.IngredientPreparations[rinsePrep.ID][chervil.ID]
	rinseBasilVIP := enums.IngredientPreparations[rinsePrep.ID][basil.ID]
	rinseMintVIP := enums.IngredientPreparations[rinsePrep.ID][mint.ID]
	rinseLargeBowlVPV := enums.PreparationVessels[rinsePrep.ID][largeBowl.ID]

	// Dry
	dryLettuceVIP := enums.IngredientPreparations[dryPrep.ID][lettuce.ID]
	dryRadicchioVIP := enums.IngredientPreparations[dryPrep.ID][radicchio.ID]
	dryEndiveVIP := enums.IngredientPreparations[dryPrep.ID][endive.ID]
	dryFriseeVIP := enums.IngredientPreparations[dryPrep.ID][frisee.ID]
	dryKaleVIP := enums.IngredientPreparations[dryPrep.ID][kale.ID]
	dryDandelionGreensVIP := enums.IngredientPreparations[dryPrep.ID][dandelionGreens.ID]
	dryPurslaneVIP := enums.IngredientPreparations[dryPrep.ID][purslane.ID]
	dryFennelFrondsVIP := enums.IngredientPreparations[dryPrep.ID][fennelFronds.ID]
	dryParsleyVIP := enums.IngredientPreparations[dryPrep.ID][parsley.ID]
	dryTarragonVIP := enums.IngredientPreparations[dryPrep.ID][tarragon.ID]
	dryChervilVIP := enums.IngredientPreparations[dryPrep.ID][chervil.ID]
	dryBasilVIP := enums.IngredientPreparations[dryPrep.ID][basil.ID]
	dryMintVIP := enums.IngredientPreparations[dryPrep.ID][mint.ID]
	drySaladSpinnerVPV := enums.PreparationVessels[dryPrep.ID][saladSpinner.ID]

	// Mix
	mixLettuceVIP := enums.IngredientPreparations[mixPrep.ID][lettuce.ID]
	mixRadicchioVIP := enums.IngredientPreparations[mixPrep.ID][radicchio.ID]
	mixEndiveVIP := enums.IngredientPreparations[mixPrep.ID][endive.ID]
	mixFriseeVIP := enums.IngredientPreparations[mixPrep.ID][frisee.ID]
	mixKaleVIP := enums.IngredientPreparations[mixPrep.ID][kale.ID]
	mixDandelionGreensVIP := enums.IngredientPreparations[mixPrep.ID][dandelionGreens.ID]
	mixPurslaneVIP := enums.IngredientPreparations[mixPrep.ID][purslane.ID]
	mixFennelFrondsVIP := enums.IngredientPreparations[mixPrep.ID][fennelFronds.ID]
	mixParsleyVIP := enums.IngredientPreparations[mixPrep.ID][parsley.ID]
	mixTarragonVIP := enums.IngredientPreparations[mixPrep.ID][tarragon.ID]
	mixChervilVIP := enums.IngredientPreparations[mixPrep.ID][chervil.ID]
	mixBasilVIP := enums.IngredientPreparations[mixPrep.ID][basil.ID]
	mixMintVIP := enums.IngredientPreparations[mixPrep.ID][mint.ID]
	mixLargeBowlVPV := enums.PreparationVessels[mixPrep.ID][largeBowl.ID]
	mixBareHandsVPI := enums.PreparationInstruments[mixPrep.ID][bareHands.ID]

	// Toss
	tossOliveOilVIP := enums.IngredientPreparations[tossPrep.ID][oliveOil.ID]
	tossServingBowlVPV := enums.PreparationVessels[tossPrep.ID][servingBowl.ID]
	tossBareHandsVPI := enums.PreparationInstruments[tossPrep.ID][bareHands.ID]

	// Season
	seasonLemonVIP := enums.IngredientPreparations[seasonPrep.ID][lemon.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonServingBowlVPV := enums.PreparationVessels[seasonPrep.ID][servingBowl.ID]

	// Measurement unit bridges
	lettuceCupVIMU := enums.IngredientMeasurementUnits[lettuce.ID][cupMeasurement.ID]
	radicchioCupVIMU := enums.IngredientMeasurementUnits[radicchio.ID][cupMeasurement.ID]
	endiveCupVIMU := enums.IngredientMeasurementUnits[endive.ID][cupMeasurement.ID]
	friseeCupVIMU := enums.IngredientMeasurementUnits[frisee.ID][cupMeasurement.ID]
	kaleCupVIMU := enums.IngredientMeasurementUnits[kale.ID][cupMeasurement.ID]
	dandelionGreensCupVIMU := enums.IngredientMeasurementUnits[dandelionGreens.ID][cupMeasurement.ID]
	purslaneCupVIMU := enums.IngredientMeasurementUnits[purslane.ID][cupMeasurement.ID]
	fennelFrondsCupVIMU := enums.IngredientMeasurementUnits[fennelFronds.ID][cupMeasurement.ID]
	parsleyCupVIMU := enums.IngredientMeasurementUnits[parsley.ID][cupMeasurement.ID]
	tarragonCupVIMU := enums.IngredientMeasurementUnits[tarragon.ID][cupMeasurement.ID]
	chervilCupVIMU := enums.IngredientMeasurementUnits[chervil.ID][cupMeasurement.ID]
	basilCupVIMU := enums.IngredientMeasurementUnits[basil.ID][cupMeasurement.ID]
	mintCupVIMU := enums.IngredientMeasurementUnits[mint.ID][cupMeasurement.ID]
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	lemonTablespoonVIMU := enums.IngredientMeasurementUnits[lemon.ID][tablespoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]

	// Step 0: Pick over the leafy vegetables, discarding any wilted or damaged leaves
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   trimPrep.ID,
		Index:           0,
		Notes:           "Pick over the leafy vegetables, discarding any wilted or damaged leaves.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimLettuceVIP.ID,
				ValidIngredientMeasurementUnitID: &lettuceCupVIMU.ID,
				IngredientID:                     &lettuce.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "lettuce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimRadicchioVIP.ID,
				ValidIngredientMeasurementUnitID: &radicchioCupVIMU.ID,
				IngredientID:                     &radicchio.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "radicchio",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimEndiveVIP.ID,
				ValidIngredientMeasurementUnitID: &endiveCupVIMU.ID,
				IngredientID:                     &endive.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "endive",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimFriseeVIP.ID,
				ValidIngredientMeasurementUnitID: &friseeCupVIMU.ID,
				IngredientID:                     &frisee.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "frisée",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimKaleVIP.ID,
				ValidIngredientMeasurementUnitID: &kaleCupVIMU.ID,
				IngredientID:                     &kale.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "kale",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimDandelionGreensVIP.ID,
				ValidIngredientMeasurementUnitID: &dandelionGreensCupVIMU.ID,
				IngredientID:                     &dandelionGreens.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "dandelion greens",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimPurslaneVIP.ID,
				ValidIngredientMeasurementUnitID: &purslaneCupVIMU.ID,
				IngredientID:                     &purslane.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "purslane",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimFennelFrondsVIP.ID,
				ValidIngredientMeasurementUnitID: &fennelFrondsCupVIMU.ID,
				IngredientID:                     &fennelFronds.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "fennel fronds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimParsleyVIP.ID,
				ValidIngredientMeasurementUnitID: &parsleyCupVIMU.ID,
				IngredientID:                     &parsley.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "parsley",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimTarragonVIP.ID,
				ValidIngredientMeasurementUnitID: &tarragonCupVIMU.ID,
				IngredientID:                     &tarragon.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "tarragon",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimChervilVIP.ID,
				ValidIngredientMeasurementUnitID: &chervilCupVIMU.ID,
				IngredientID:                     &chervil.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "chervil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimBasilVIP.ID,
				ValidIngredientMeasurementUnitID: &basilCupVIMU.ID,
				IngredientID:                     &basil.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "basil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &trimMintVIP.ID,
				ValidIngredientMeasurementUnitID: &mintCupVIMU.ID,
				IngredientID:                     &mint.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "mint",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
				ValidPreparationInstrumentID: &trimBareHandsVPI.ID,
				InstrumentID:                 &bareHands.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &trimCuttingBoardVPV.ID,
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
				Name:                "inspected lettuce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected radicchio",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               1,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected endive",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               2,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected frisée",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               3,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected kale",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               4,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected dandelion greens",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               5,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected purslane",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               6,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected fennel fronds",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               7,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected parsley",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               8,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected tarragon",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               9,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected chervil",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               10,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected basil",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               11,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "inspected mint",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               12,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	// Step 1: Cut lettuce leaves free of their cores, pick tender herbs from stems, and quarter, core, and slice tight leafy heads like radicchio and endive
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   slicePrep.ID,
		Index:           1,
		Notes:           "Cut lettuce leaves free of their cores, pick tender herbs from stems, and quarter, core, and slice tight leafy heads like radicchio and endive.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &sliceLettuceVIP.ID,
				IngredientID:                    &lettuce.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected lettuce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &sliceRadicchioVIP.ID,
				IngredientID:                    &radicchio.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected radicchio",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](2),
				ValidIngredientPreparationID:    &sliceEndiveVIP.ID,
				IngredientID:                    &endive.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected endive",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidPreparationInstrumentID: &sliceKnifeVPI.ID,
				InstrumentID:                 &knife.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step1ID,
				ValidPreparationVesselID: &sliceCuttingBoardVPV.ID,
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
				Name:                "prepared lettuce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "prepared radicchio",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               1,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "prepared endive",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               2,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Wash everything in several changes of water until no dirt or grit remains
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   rinsePrep.ID,
		Index:           2,
		Notes:           "Wash everything in several changes of water until no dirt or grit remains.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &rinseLettuceVIP.ID,
				IngredientID:                    &lettuce.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "prepared lettuce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &rinseRadicchioVIP.ID,
				IngredientID:                    &radicchio.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "prepared radicchio",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](2),
				ValidIngredientPreparationID:    &rinseEndiveVIP.ID,
				IngredientID:                    &endive.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "prepared endive",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](3),
				ValidIngredientPreparationID:    &rinseFriseeVIP.ID,
				IngredientID:                    &frisee.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected frisée",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](4),
				ValidIngredientPreparationID:    &rinseKaleVIP.ID,
				IngredientID:                    &kale.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected kale",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](5),
				ValidIngredientPreparationID:    &rinseDandelionGreensVIP.ID,
				IngredientID:                    &dandelionGreens.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected dandelion greens",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](6),
				ValidIngredientPreparationID:    &rinsePurslaneVIP.ID,
				IngredientID:                    &purslane.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected purslane",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](7),
				ValidIngredientPreparationID:    &rinseFennelFrondsVIP.ID,
				IngredientID:                    &fennelFronds.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected fennel fronds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](8),
				ValidIngredientPreparationID:    &rinseParsleyVIP.ID,
				IngredientID:                    &parsley.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected parsley",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](9),
				ValidIngredientPreparationID:    &rinseTarragonVIP.ID,
				IngredientID:                    &tarragon.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected tarragon",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](10),
				ValidIngredientPreparationID:    &rinseChervilVIP.ID,
				IngredientID:                    &chervil.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected chervil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](11),
				ValidIngredientPreparationID:    &rinseBasilVIP.ID,
				IngredientID:                    &basil.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected basil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](12),
				ValidIngredientPreparationID:    &rinseMintVIP.ID,
				IngredientID:                    &mint.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "inspected mint",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &rinseLargeBowlVPV.ID,
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
				BelongsToRecipeStep: step2ID,
				Name:                "washed lettuce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed radicchio",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               1,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed endive",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               2,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed frisée",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               3,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed kale",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               4,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed dandelion greens",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               5,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed purslane",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               6,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed fennel fronds",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               7,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed parsley",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               8,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed tarragon",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               9,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed chervil",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               10,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed basil",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               11,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "washed mint",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               12,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	// Step 3: Dry well in a salad spinner
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dryPrep.ID,
		Index:           3,
		Notes:           "Dry well in a salad spinner.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &dryLettuceVIP.ID,
				IngredientID:                    &lettuce.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed lettuce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &dryRadicchioVIP.ID,
				IngredientID:                    &radicchio.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed radicchio",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](2),
				ValidIngredientPreparationID:    &dryEndiveVIP.ID,
				IngredientID:                    &endive.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed endive",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](3),
				ValidIngredientPreparationID:    &dryFriseeVIP.ID,
				IngredientID:                    &frisee.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed frisée",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](4),
				ValidIngredientPreparationID:    &dryKaleVIP.ID,
				IngredientID:                    &kale.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed kale",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](5),
				ValidIngredientPreparationID:    &dryDandelionGreensVIP.ID,
				IngredientID:                    &dandelionGreens.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed dandelion greens",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](6),
				ValidIngredientPreparationID:    &dryPurslaneVIP.ID,
				IngredientID:                    &purslane.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed purslane",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](7),
				ValidIngredientPreparationID:    &dryFennelFrondsVIP.ID,
				IngredientID:                    &fennelFronds.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed fennel fronds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](8),
				ValidIngredientPreparationID:    &dryParsleyVIP.ID,
				IngredientID:                    &parsley.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed parsley",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](9),
				ValidIngredientPreparationID:    &dryTarragonVIP.ID,
				IngredientID:                    &tarragon.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed tarragon",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](10),
				ValidIngredientPreparationID:    &dryChervilVIP.ID,
				IngredientID:                    &chervil.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed chervil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](11),
				ValidIngredientPreparationID:    &dryBasilVIP.ID,
				IngredientID:                    &basil.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed basil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](12),
				ValidIngredientPreparationID:    &dryMintVIP.ID,
				IngredientID:                    &mint.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "washed mint",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &drySaladSpinnerVPV.ID,
				VesselID:                 &saladSpinner.ID,
				Name:                     "salad spinner",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried lettuce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried radicchio",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               1,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried endive",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               2,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried frisée",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               3,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried kale",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               4,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried dandelion greens",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               5,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried purslane",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               6,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried fennel fronds",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               7,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried parsley",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               8,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried tarragon",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               9,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried chervil",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               10,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried basil",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               11,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "dried mint",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               12,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	// Step 4: Combine all the dried greens into a mixed greens mixture
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   mixPrep.ID,
		Index:           4,
		Notes:           "Combine all the dried greens and herbs together in a large bowl.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &mixLettuceVIP.ID,
				IngredientID:                    &lettuce.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried lettuce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &mixRadicchioVIP.ID,
				IngredientID:                    &radicchio.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried radicchio",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](2),
				ValidIngredientPreparationID:    &mixEndiveVIP.ID,
				IngredientID:                    &endive.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried endive",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](3),
				ValidIngredientPreparationID:    &mixFriseeVIP.ID,
				IngredientID:                    &frisee.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried frisée",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](4),
				ValidIngredientPreparationID:    &mixKaleVIP.ID,
				IngredientID:                    &kale.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried kale",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](5),
				ValidIngredientPreparationID:    &mixDandelionGreensVIP.ID,
				IngredientID:                    &dandelionGreens.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried dandelion greens",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](6),
				ValidIngredientPreparationID:    &mixPurslaneVIP.ID,
				IngredientID:                    &purslane.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried purslane",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](7),
				ValidIngredientPreparationID:    &mixFennelFrondsVIP.ID,
				IngredientID:                    &fennelFronds.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried fennel fronds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](8),
				ValidIngredientPreparationID:    &mixParsleyVIP.ID,
				IngredientID:                    &parsley.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried parsley",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](9),
				ValidIngredientPreparationID:    &mixTarragonVIP.ID,
				IngredientID:                    &tarragon.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried tarragon",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](10),
				ValidIngredientPreparationID:    &mixChervilVIP.ID,
				IngredientID:                    &chervil.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried chervil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](11),
				ValidIngredientPreparationID:    &mixBasilVIP.ID,
				IngredientID:                    &basil.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried basil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](12),
				ValidIngredientPreparationID:    &mixMintVIP.ID,
				IngredientID:                    &mint.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "dried mint",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &mixBareHandsVPI.ID,
				InstrumentID:                 &bareHands.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &mixLargeBowlVPV.ID,
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
				BelongsToRecipeStep: step4ID,
				Name:                "mixed greens",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](8),
				},
			},
		},
	}

	// Step 5: In a large serving bowl, gently toss salad with just enough olive oil to gently coat leaves
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   tossPrep.ID,
		Index:           5,
		Notes:           "In a large serving bowl, gently toss salad with just enough olive oil to gently coat leaves.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &lettuce.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "mixed greens",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &tossOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				IngredientID:                     &oliveOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "extra-virgin olive oil",
				QuantityNotes:                    "just enough to gently coat leaves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
					Max: pointer.To[float32](3),
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &tossBareHandsVPI.ID,
				InstrumentID:                 &bareHands.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step5ID,
				ValidPreparationVesselID: &tossServingBowlVPV.ID,
				VesselID:                 &servingBowl.ID,
				Name:                     "large serving bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "dressed greens",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "serving bowl with salad",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Add splash of lemon juice and salt to taste, tossing to combine. Serve.
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           6,
		Notes:           "Add a splash of lemon juice and salt to taste, tossing to combine. Serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &lettuce.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "dressed greens",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &seasonLemonVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonTablespoonVIMU.ID,
				IngredientID:                     &lemon.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "fresh lemon juice",
				QuantityNotes:                    "a splash, to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
					Max: pointer.To[float32](1),
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "kosher or sea salt",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
					Max: pointer.To[float32](0.5),
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &seasonServingBowlVPV.ID,
				VesselID:                        &servingBowl.ID,
				Name:                            "serving bowl with salad",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "mixed green salad",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	mixedGreenSaladRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Mixed Green Salad",
		Slug:                "mixed-green-salad",
		Source:              "https://www.seriouseats.com/roman-mixed-green-salad-misticanza-recipe",
		Description:         "A Roman-inspired mixed green salad (Misticanza alla Romana) made with a variety of fresh leafy greens and tender herbs, dressed simply with extra-virgin olive oil, fresh lemon juice, and salt.",
		YieldsComponentType: mealplanning.MealComponentTypesSalad,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		mixedGreenSaladRecipe,
	}
}
