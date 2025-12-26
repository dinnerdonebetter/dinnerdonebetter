package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// StovetopMacAndCheeseRecipe creates the Ultra-Gooey Stovetop Mac and Cheese recipe.
// Source: https://www.seriouseats.com/the-food-labs-ultra-gooey-stovetop-mac-cheese
func StovetopMacAndCheeseRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

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

	// Get vessels
	saucepan := enums.Vessels["saucepan"]
	largeBowl := enums.Vessels["large bowl"]
	mediumBowl := enums.Vessels["medium bowl"]
	colander := enums.Vessels["colander"]

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
	removeFromHeatMacaroniVIP := enums.IngredientPreparations[removeFromHeatPrep.ID][elbowMacaroni.ID]
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
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   submergePrep.ID,
		Index:           0,
		Notes:           "Place the macaroni in a large saucepan and cover it with salted water by 2 inches.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &submergeMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				IngredientID:                     &elbowMacaroni.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "elbow macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &submergeWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID].ID,
				IngredientID:                     &water.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				QuantityNotes:                    "enough to cover macaroni by 2 inches",
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &submergeSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID].ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				QuantityNotes:                    "to taste",
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &submergeSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "large saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "macaroni in salted water",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Bring to a boil, stirring occasionally
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   boilPrep.ID,
		Index:           1,
		Notes:           "Bring to a boil over high heat, stirring occasionally to keep the pasta from sticking.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &boilMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				IngredientID:                     &elbowMacaroni.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "macaroni in salted water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidPreparationInstrumentID: &boilWoodenSpoonVPI.ID,
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step1ID,
				ValidPreparationVesselID: &boilSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "boiling macaroni",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Cover the pan
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   coverPrep.ID,
		Index:           2,
		Notes:           "Cover the pan.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &elbowMacaroni.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "boiling macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &coverSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "covered macaroni",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Remove the pan from the heat
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   removeFromHeatPrep.ID,
		Index:           3,
		Notes:           "Remove the pan from the heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &removeFromHeatMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				IngredientID:                     &elbowMacaroni.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "covered macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &removeFromHeatSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "covered macaroni off heat",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Let stand until barely al dente
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           4,
		Notes:           "Let stand until the pasta is barely al dente, about 8 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](480), // 8 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &restMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				IngredientID:                     &elbowMacaroni.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "covered macaroni off heat",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &restSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
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
				Notes:               "pasta should be barely al dente",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "al dente macaroni",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Mix together evaporated milk, eggs, hot sauce, and mustard
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   mixPrep.ID,
		Index:           5,
		Notes:           "Meanwhile, mix together the evaporated milk, eggs, hot sauce, and mustard in a bowl until homogeneous.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &mixEvaporatedMilkVIP.ID,
				ValidIngredientMeasurementUnitID: &evaporatedMilkFluidOzVIMU.ID,
				IngredientID:                     &evaporatedMilk.ID,
				MeasurementUnitID:                fluidOunceMeasurement.ID,
				Name:                             "evaporated milk",
				QuantityNotes:                    "one 12-ounce can",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 12,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &mixEggsVIP.ID,
				ValidIngredientMeasurementUnitID: &eggsUnitVIMU.ID,
				IngredientID:                     &eggs.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "large eggs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &mixHotSauceVIP.ID,
				ValidIngredientMeasurementUnitID: &hotSauceTeaspoonVIMU.ID,
				IngredientID:                     &hotSauce.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "hot sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ValidIngredientPreparationID:     &mixGroundMustardVIP.ID,
				ValidIngredientMeasurementUnitID: &groundMustardTeaspoonVIMU.ID,
				IngredientID:                     &groundMustard.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "ground mustard",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &mixWhiskVPI.ID,
				InstrumentID:                 &whisk.ID,
				Name:                         "whisk",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step5ID,
				ValidPreparationVesselID: &mixMediumBowlVPV.ID,
				VesselID:                 &mediumBowl.ID,
				Name:                     "bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "milk mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Toss cheeses with cornstarch
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   tossPrep.ID,
		Index:           6,
		Notes:           "Toss the cheeses with the cornstarch in a large bowl until thoroughly combined.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &tossCheddarVIP.ID,
				ValidIngredientMeasurementUnitID: &cheddarPoundVIMU.ID,
				IngredientID:                     &cheddarCheese.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "extra-sharp cheddar, grated",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &tossAmericanVIP.ID,
				ValidIngredientMeasurementUnitID: &americanOunceVIMU.ID,
				IngredientID:                     &americanCheese.ID,
				MeasurementUnitID:                ounceMeasurement.ID,
				Name:                             "American cheese, cut into 1/2-inch cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &tossCornstarchVIP.ID,
				ValidIngredientMeasurementUnitID: &cornstarchTablespoonVIMU.ID,
				IngredientID:                     &cornstarch.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "cornstarch",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &tossLargeBowlVPV.ID,
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
				Name:                "cheese mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Drain pasta
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   drainPrep.ID,
		Index:           7,
		Notes:           "When the pasta is cooked, drain it.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &drainMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				IngredientID:                     &elbowMacaroni.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "al dente macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step7ID,
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
				BelongsToRecipeStep: step7ID,
				Name:                "drained macaroni",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Return pasta to saucepan
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           8,
		Notes:           "Return the drained pasta to the saucepan and place over low heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &elbowMacaroni.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "drained macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: &addSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "macaroni in saucepan",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Add butter and stir until melted
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           9,
		Notes:           "Add the butter and stir until melted.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step9ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &stirMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				IngredientID:                     &elbowMacaroni.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "macaroni in saucepan",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step9ID,
				ValidIngredientPreparationID:     &stirButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				IngredientID:                     &butter.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "unsalted butter, cut into 4 chunks",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8, // 8 tablespoons = 1 stick
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step9ID,
				ValidPreparationVesselID: &stirSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				IngredientStateID:   meltedState.ID,
				Notes:               "butter should be completely melted",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "buttered macaroni",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Add milk and cheese mixtures, cook until creamy
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cookPrep.ID,
		Index:           10,
		Notes:           "Add the milk mixture and cheese mixture and cook, stirring constantly, until the cheese is completely melted and the mixture is hot and creamy.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &cookMacaroniVIP.ID,
				ValidIngredientMeasurementUnitID: &macaroniPoundVIMU.ID,
				IngredientID:                     &elbowMacaroni.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "buttered macaroni",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &cookEvaporatedMilkVIP.ID,
				IngredientID:                    &evaporatedMilk.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "milk mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &cookCheddarVIP.ID,
				IngredientID:                    &cheddarCheese.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "cheese mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step10ID,
				ValidPreparationInstrumentID: &cookWoodenSpoonVPI.ID,
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: &cookSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				IngredientStateID:   meltedState.ID,
				Notes:               "cheese should be completely melted and mixture hot and creamy",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "mac and cheese",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 11: Season to taste
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           11,
		Notes:           "Season to taste with salt and more hot sauce. Serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &elbowMacaroni.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "mac and cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID].ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "salt",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				Optional: true,
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &seasonHotSauceVIP.ID,
				ValidIngredientMeasurementUnitID: &hotSauceTeaspoonVIMU.ID,
				IngredientID:                     &hotSauce.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "hot sauce",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step11ID,
				ValidPreparationVesselID: &seasonSaucepanVPV.ID,
				VesselID:                 &saucepan.ID,
				Name:                     "saucepan",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "stovetop mac and cheese",
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
			Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
				step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11,
			},
		},
	}
}
