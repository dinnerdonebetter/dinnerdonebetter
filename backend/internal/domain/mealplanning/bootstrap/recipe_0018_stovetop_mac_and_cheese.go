package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// StovetopMacAndCheeseRecipe creates the Ultra-Gooey Stovetop Mac and Cheese recipe.
// Source: https://www.seriouseats.com/the-food-labs-ultra-gooey-stovetop-mac-cheese
func StovetopMacAndCheeseRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	submergePrep := enums.Preparations["submerge"]
	boilPrep := enums.Preparations["boil"]
	coverPrep := enums.Preparations["cover"]
	restPrep := enums.Preparations["rest"]
	mixPrep := enums.Preparations["mix"]
	removeFromHeatPrep := enums.Preparations["remove from heat"]
	tossPrep := enums.Preparations["toss"]
	drainPrep := enums.Preparations["drain"]
	addPrep := enums.Preparations["add"]
	stirPrep := enums.Preparations["stir"]
	cookPrep := enums.Preparations["cook"]
	seasonPrep := enums.Preparations["season"]
	cutPrep := enums.Preparations["cut"]
	gratePrep := enums.Preparations["grate"]

	// Get ingredients
	elbowMacaroni := enums.Ingredients["elbow macaroni"]
	salt := enums.Ingredients["salt"]
	water := enums.Ingredients["water"]
	evaporatedMilk := enums.Ingredients["evaporated milk"]
	eggs := enums.Ingredients["eggs"]
	hotSauce := enums.Ingredients["hot sauce"]
	groundMustard := enums.Ingredients["ground mustard"]
	cheddarCheese := enums.Ingredients["cheddar cheese"]
	americanCheese := enums.Ingredients["American cheese"]
	cornstarch := enums.Ingredients["cornstarch"]
	butter := enums.Ingredients["butter"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	fluidOunceMeasurement := enums.MeasurementUnits["fluid ounce"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	ounceMeasurement := enums.MeasurementUnits["ounce"]
	cupMeasurement := enums.MeasurementUnits["cup"]

	// Get instruments
	whisk := enums.Instruments["whisk"]
	woodenSpoon := enums.Instruments["wooden spoon"]
	knife := enums.Instruments["knife"]
	microplane := enums.Instruments["microplane"]

	// Get vessels
	saucepan := enums.Vessels["saucepan"]
	largeBowl := enums.Vessels["large bowl"]
	mediumBowl := enums.Vessels["medium bowl"]
	colander := enums.Vessels["colander"]
	cuttingBoard := enums.Vessels["cutting board"]

	// Get bridge table entries
	// Submerge preparation bridges
	submergeMacaroniVIP := enums.IngredientPreparations[submergePrep.ID][elbowMacaroni.ID]
	submergeWaterVIP := enums.IngredientPreparations[submergePrep.ID][water.ID]
	submergeSaltVIP := enums.IngredientPreparations[submergePrep.ID][salt.ID]
	submergeSaucepanVPV := enums.PreparationVessels[submergePrep.ID][saucepan.ID]

	// Boil preparation bridges
	boilMacaroniVIP := enums.IngredientPreparations[boilPrep.ID][elbowMacaroni.ID]
	boilWoodenSpoonVPI := enums.PreparationInstruments[boilPrep.ID][woodenSpoon.ID]
	boilSaucepanVPV := enums.PreparationVessels[boilPrep.ID][saucepan.ID]

	// Cover preparation bridges
	coverSaucepanVPV := enums.PreparationVessels[coverPrep.ID][saucepan.ID]

	// Rest preparation bridges
	restMacaroniVIP := enums.IngredientPreparations[restPrep.ID][elbowMacaroni.ID]
	restSaucepanVPV := enums.PreparationVessels[restPrep.ID][saucepan.ID]

	// Mix preparation bridges
	mixEvaporatedMilkVIP := enums.IngredientPreparations[mixPrep.ID][evaporatedMilk.ID]
	mixEggsVIP := enums.IngredientPreparations[mixPrep.ID][eggs.ID]
	mixHotSauceVIP := enums.IngredientPreparations[mixPrep.ID][hotSauce.ID]
	mixGroundMustardVIP := enums.IngredientPreparations[mixPrep.ID][groundMustard.ID]
	mixWhiskVPI := enums.PreparationInstruments[mixPrep.ID][whisk.ID]
	mixMediumBowlVPV := enums.PreparationVessels[mixPrep.ID][mediumBowl.ID]

	// Remove from heat preparation bridges
	removeFromHeatSaucepanVPV := enums.PreparationVessels[removeFromHeatPrep.ID][saucepan.ID]

	// Toss preparation bridges
	tossCheddarVIP := enums.IngredientPreparations[tossPrep.ID][cheddarCheese.ID]
	tossAmericanVIP := enums.IngredientPreparations[tossPrep.ID][americanCheese.ID]
	tossCornstarchVIP := enums.IngredientPreparations[tossPrep.ID][cornstarch.ID]
	tossLargeBowlVPV := enums.PreparationVessels[tossPrep.ID][largeBowl.ID]

	// Drain preparation bridges
	drainMacaroniVIP := enums.IngredientPreparations[drainPrep.ID][elbowMacaroni.ID]
	drainColanderVPV := enums.PreparationVessels[drainPrep.ID][colander.ID]

	// Add preparation bridges
	addSaucepanVPV := enums.PreparationVessels[addPrep.ID][saucepan.ID]

	// Stir preparation bridges
	stirMacaroniVIP := enums.IngredientPreparations[stirPrep.ID][elbowMacaroni.ID]
	stirButterVIP := enums.IngredientPreparations[stirPrep.ID][butter.ID]
	stirSaucepanVPV := enums.PreparationVessels[stirPrep.ID][saucepan.ID]

	// Cook preparation bridges
	cookMacaroniVIP := enums.IngredientPreparations[cookPrep.ID][elbowMacaroni.ID]
	cookCheddarVIP := enums.IngredientPreparations[cookPrep.ID][cheddarCheese.ID]
	cookEvaporatedMilkVIP := enums.IngredientPreparations[cookPrep.ID][evaporatedMilk.ID]
	cookWoodenSpoonVPI := enums.PreparationInstruments[cookPrep.ID][woodenSpoon.ID]
	cookSaucepanVPV := enums.PreparationVessels[cookPrep.ID][saucepan.ID]

	// Season preparation bridges
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonHotSauceVIP := enums.IngredientPreparations[seasonPrep.ID][hotSauce.ID]
	seasonSaucepanVPV := enums.PreparationVessels[seasonPrep.ID][saucepan.ID]

	// Cut preparation bridges
	cutButterVIP := enums.IngredientPreparations[cutPrep.ID][butter.ID]
	cutAmericanCheeseVIP := enums.IngredientPreparations[cutPrep.ID][americanCheese.ID]
	cutKnifeVPI := enums.PreparationInstruments[cutPrep.ID][knife.ID]
	cutCuttingBoardVPV := enums.PreparationVessels[cutPrep.ID][cuttingBoard.ID]

	// Grate preparation bridges
	grateCheddarVIP := enums.IngredientPreparations[gratePrep.ID][cheddarCheese.ID]
	grateMicroplaneVPI := enums.PreparationInstruments[gratePrep.ID][microplane.ID]
	grateCuttingBoardVPV := enums.PreparationVessels[gratePrep.ID][cuttingBoard.ID]

	// Measurement unit bridges
	macaroniPoundVIMU := enums.IngredientMeasurementUnits[elbowMacaroni.ID][poundMeasurement.ID]
	evaporatedMilkFluidOzVIMU := enums.IngredientMeasurementUnits[evaporatedMilk.ID][fluidOunceMeasurement.ID]
	eggsUnitVIMU := enums.IngredientMeasurementUnits[eggs.ID][unitMeasurement.ID]
	hotSauceTeaspoonVIMU := enums.IngredientMeasurementUnits[hotSauce.ID][teaspoonMeasurement.ID]
	groundMustardTeaspoonVIMU := enums.IngredientMeasurementUnits[groundMustard.ID][teaspoonMeasurement.ID]
	cheddarPoundVIMU := enums.IngredientMeasurementUnits[cheddarCheese.ID][poundMeasurement.ID]
	americanOunceVIMU := enums.IngredientMeasurementUnits[americanCheese.ID][ounceMeasurement.ID]
	cornstarchTablespoonVIMU := enums.IngredientMeasurementUnits[cornstarch.ID][tablespoonMeasurement.ID]
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]

	// Get ingredient states for completion conditions
	tenderState := enums.IngredientStates["tender"]
	meltedState := enums.IngredientStates["at desired consistency"]

	// Step 0: Place macaroni in saucepan and cover with salted water
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: submergePrep.ID,
		Index:         0,
		Notes:         "Place the macaroni in a large saucepan and cover it with salted water by 2 inches.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &submergeMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				Name:                             "elbow macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &submergeWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID].ID,
				QuantityNotes:                    "enough to cover macaroni by 2 inches",
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ValidIngredientPreparationID:     &submergeSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID].ID,
				QuantityNotes:                    "to taste",
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &submergeSaucepanVPV.ID,
				Name:                     "large saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "macaroni in salted water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Bring to a boil, stirring occasionally
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: boilPrep.ID,
		Index:         1,
		Notes:         "Bring to a boil over high heat, stirring occasionally to keep the pasta from sticking.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &boilMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				Name:                             "macaroni in salted water",
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
				ValidPreparationVesselID: &boilSaucepanVPV.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "boiling macaroni",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Cover the pan
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: coverPrep.ID,
		Index:         2,
		Notes:         "Cover the pan.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "boiling macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &coverSaucepanVPV.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "covered macaroni",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Remove the pan from the heat
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: removeFromHeatPrep.ID,
		Index:         3,
		Notes:         "Remove the pan from the heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "covered macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &removeFromHeatSaucepanVPV.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "covered macaroni off heat",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Let stand until barely al dente
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: restPrep.ID,
		Index:         4,
		Notes:         "Let stand until the pasta is barely al dente, about 8 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](480), // 8 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &restMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				Name:                             "covered macaroni off heat",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &restSaucepanVPV.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: tenderState.ID,
				Notes:             "pasta should be barely al dente",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "al dente macaroni",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Mix together evaporated milk, eggs, hot sauce, and mustard
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: mixPrep.ID,
		Index:         5,
		Notes:         "Meanwhile, mix together the evaporated milk, eggs, hot sauce, and mustard in a bowl until homogeneous.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &mixEvaporatedMilkVIP.ID,
				ValidIngredientMeasurementUnitID: &evaporatedMilkFluidOzVIMU.ID,
				Name:                             "evaporated milk",
				QuantityNotes:                    "one 12-ounce can",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 12,
				},
			},
			{
				ValidIngredientPreparationID:     &mixEggsVIP.ID,
				ValidIngredientMeasurementUnitID: &eggsUnitVIMU.ID,
				Name:                             "large eggs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &mixHotSauceVIP.ID,
				ValidIngredientMeasurementUnitID: &hotSauceTeaspoonVIMU.ID,
				Name:                             "hot sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &mixGroundMustardVIP.ID,
				ValidIngredientMeasurementUnitID: &groundMustardTeaspoonVIMU.ID,
				Name:                             "ground mustard",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &mixWhiskVPI.ID,
				Name:                         "whisk",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &mixMediumBowlVPV.ID,
				Name:                     "bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "milk mixture",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5a: Grate cheddar cheese
	step5a := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: gratePrep.ID,
		Index:         6,
		Notes:         "Grate the cheddar cheese.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &grateCheddarVIP.ID,
				ValidIngredientMeasurementUnitID: &cheddarPoundVIMU.ID,
				Name:                             "extra-sharp cheddar cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &grateMicroplaneVPI.ID,
				Name:                         "microplane",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &grateCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "grated cheddar cheese",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5b: Cut American cheese into cubes
	step5b := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: cutPrep.ID,
		Index:         7,
		Notes:         "Cut the American cheese into 1/2-inch cubes.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &cutAmericanCheeseVIP.ID,
				ValidIngredientMeasurementUnitID: &americanOunceVIMU.ID,
				Name:                             "American cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &cutKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &cutCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "American cheese, cut into 1/2-inch cubes",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &ounceMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](8),
				},
			},
		},
	}

	// Step 6: Toss cheeses with cornstarch
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: tossPrep.ID,
		Index:         8,
		Notes:         "Toss the cheeses with the cornstarch in a large bowl until thoroughly combined.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossCheddarVIP.ID,
				Name:                            "grated cheddar cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &tossAmericanVIP.ID,
				Name:                            "American cheese, cut into 1/2-inch cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ValidIngredientPreparationID:     &tossCornstarchVIP.ID,
				ValidIngredientMeasurementUnitID: &cornstarchTablespoonVIMU.ID,
				Name:                             "cornstarch",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &tossLargeBowlVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cheese mixture",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Drain pasta
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: drainPrep.ID,
		Index:         9,
		Notes:         "When the pasta is cooked, drain it.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &drainMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				Name:                             "al dente macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
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
				Name:              "drained macaroni",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Return pasta to saucepan
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: addPrep.ID,
		Index:         10,
		Notes:         "Return the drained pasta to the saucepan and place over low heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "drained macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &addSaucepanVPV.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "macaroni in saucepan",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8a: Cut butter into chunks
	step8a := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: cutPrep.ID,
		Index:         11,
		Notes:         "Cut the butter into 4 chunks.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &cutButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				Name:                             "unsalted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8, // 8 tablespoons = 1 stick
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &cutKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &cutCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "butter, cut into 4 chunks",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &tablespoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](8)},
			},
		},
	}

	// Step 9: Add butter and stir until melted
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: stirPrep.ID,
		Index:         12,
		Notes:         "Add the butter and stir until melted.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &stirMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				Name:                             "macaroni in saucepan",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirButterVIP.ID,
				Name:                            "butter, cut into 4 chunks",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8, // 8 tablespoons = 1 stick
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &stirSaucepanVPV.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: meltedState.ID,
				Notes:             "butter should be completely melted",
				Ingredients:       []uint64{1},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "buttered macaroni",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Add milk and cheese mixtures, cook until creamy
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: cookPrep.ID,
		Index:         13,
		Notes:         "Add the milk mixture and cheese mixture and cook, stirring constantly, until the cheese is completely melted and the mixture is hot and creamy.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &cookMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				Name:                             "buttered macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &cookEvaporatedMilkVIP.ID,
				Name:                            "milk mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &cookCheddarVIP.ID,
				Name:                            "cheese mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &cookWoodenSpoonVPI.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &cookSaucepanVPV.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: meltedState.ID,
				Notes:             "cheese should be completely melted and mixture hot and creamy",
				Ingredients:       []uint64{2},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "mac and cheese",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 11: Season to taste
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         14,
		Notes:         "Season to taste with salt and more hot sauce. Serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "mac and cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID].ID,
				Name:                             "salt",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				Optional: true,
			},
			{
				ValidIngredientPreparationID:     &seasonHotSauceVIP.ID,
				ValidIngredientMeasurementUnitID: &hotSauceTeaspoonVIMU.ID,
				Name:                             "hot sauce",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &seasonSaucepanVPV.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "stovetop mac and cheese",
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
			Name:                "Stovetop Mac and Cheese",
			Slug:                "stovetop-mac-and-cheese",
			Source:              "https://www.seriouseats.com/the-food-labs-ultra-gooey-stovetop-mac-cheese",
			Description:         "A homemade stovetop mac and cheese recipe that's ultra-gooey with the consistency of the blue box but with the complex flavor of real cheese.",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 4,
				Max: pointer.To[float32](6),
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				step0, step1, step2, step3, step4, step5, step5a, step5b, step6, step7, step8, step8a, step9, step10, step11,
			},
			PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
		},
	}
}
