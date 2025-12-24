package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// SoySauceBraisedChickenThighsRecipe creates the Soy Sauce–Braised Chicken Thighs recipe.
// Source: https://www.seriouseats.com/soy-sauce-braised-chicken-thighs-recipe-8737800
func SoySauceBraisedChickenThighsRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	combinePrep := enums.Preparations["combine"]
	dryPrep := enums.Preparations["dry"]
	seasonPrep := enums.Preparations["season"]
	transferPrep := enums.Preparations["transfer"]
	preheatPrep := enums.Preparations["preheat"]
	heatPrep := enums.Preparations["heat"]
	panSearPrep := enums.Preparations["pan-sear"]
	flipPrep := enums.Preparations["flip"]
	sautPrep := enums.Preparations["sauté"]
	simmerPrep := enums.Preparations["simmer"]
	braisePrep := enums.Preparations["braise"]

	// Get ingredients
	salt := enums.Ingredients["salt"]
	msg := enums.Ingredients["MSG"]
	fiveSpice := enums.Ingredients["five spice powder"]
	darkBrownSugar := enums.Ingredients["dark brown sugar"]
	whitePepper := enums.Ingredients["white pepper"]
	chickenThighs := enums.Ingredients["chicken thigh"]
	vegetableOil := enums.Ingredients["vegetable oil"]
	scallions := enums.Ingredients["scallions"]
	ginger := enums.Ingredients["ginger"]
	garlic := enums.Ingredients["garlic"]
	starAnise := enums.Ingredients["star anise"]
	cassiaBark := enums.Ingredients["cassia bark"]
	lightSoySauce := enums.Ingredients["light soy sauce"]
	shaoxingWine := enums.Ingredients["Shaoxing wine"]
	water := enums.Ingredients["water"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	poundMeasurement := enums.MeasurementUnits["pound"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	cloveMeasurement := enums.MeasurementUnits["clove"]

	// Get instruments
	paperTowels := enums.Instruments["paper towels"]
	whisk := enums.Instruments["whisk"]
	bareHands := enums.Instruments["bare hands"]
	tongs := enums.Instruments["tongs"]
	woodenSpoon := enums.Instruments["wooden spoon"]
	thermometer := enums.Instruments["instant-read thermometer"]

	// Get vessels
	smallBowl := enums.Vessels["small bowl"]
	wireRack := enums.Vessels["wire rack"]
	bakingSheet := enums.Vessels["baking sheet"]
	castIronSkillet := enums.Vessels["cast iron skillet"]
	largePlate := enums.Vessels["large plate"]
	oven := enums.Vessels["oven"]

	// Get bridge table entries
	// Combine
	combineSaltVIP := enums.IngredientPreparations[combinePrep.ID][salt.ID]
	combineMSGVIP := enums.IngredientPreparations[combinePrep.ID][msg.ID]
	combineFiveSpiceVIP := enums.IngredientPreparations[combinePrep.ID][fiveSpice.ID]
	combineDarkBrownSugarVIP := enums.IngredientPreparations[combinePrep.ID][darkBrownSugar.ID]
	combineWhitePepperVIP := enums.IngredientPreparations[combinePrep.ID][whitePepper.ID]
	combineSmallBowlVPV := enums.PreparationVessels[combinePrep.ID][smallBowl.ID]
	combineWhiskVPI := enums.PreparationInstruments[combinePrep.ID][whisk.ID]

	// Dry
	dryChickenVIP := enums.IngredientPreparations[dryPrep.ID][chickenThighs.ID]
	dryPaperTowelsVPI := enums.PreparationInstruments[dryPrep.ID][paperTowels.ID]

	// Season
	seasonChickenVIP := enums.IngredientPreparations[seasonPrep.ID][chickenThighs.ID]
	seasonBareHandsVPI := enums.PreparationInstruments[seasonPrep.ID][bareHands.ID]

	// Transfer
	transferChickenVIP := enums.IngredientPreparations[transferPrep.ID][chickenThighs.ID]
	transferWireRackVPV := enums.PreparationVessels[transferPrep.ID][wireRack.ID]
	transferBakingSheetVPV := enums.PreparationVessels[transferPrep.ID][bakingSheet.ID]
	transferLargePlateVPV := enums.PreparationVessels[transferPrep.ID][largePlate.ID]
	transferSkilletVPV := enums.PreparationVessels[transferPrep.ID][castIronSkillet.ID]

	// Preheat
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]

	// Heat
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][vegetableOil.ID]
	heatSkilletVPV := enums.PreparationVessels[heatPrep.ID][castIronSkillet.ID]

	// Pan-sear
	panSearChickenVIP := enums.IngredientPreparations[panSearPrep.ID][chickenThighs.ID]
	panSearSkilletVPV := enums.PreparationVessels[panSearPrep.ID][castIronSkillet.ID]
	panSearTongsVPI := enums.PreparationInstruments[panSearPrep.ID][tongs.ID]

	// Flip
	flipChickenVIP := enums.IngredientPreparations[flipPrep.ID][chickenThighs.ID]
	flipSkilletVPV := enums.PreparationVessels[flipPrep.ID][castIronSkillet.ID]
	flipTongsVPI := enums.PreparationInstruments[flipPrep.ID][tongs.ID]

	// Sauté
	sautScallionsVIP := enums.IngredientPreparations[sautPrep.ID][scallions.ID]
	sautGingerVIP := enums.IngredientPreparations[sautPrep.ID][ginger.ID]
	sautGarlicVIP := enums.IngredientPreparations[sautPrep.ID][garlic.ID]
	sautFiveSpiceVIP := enums.IngredientPreparations[sautPrep.ID][fiveSpice.ID]
	sautDarkBrownSugarVIP := enums.IngredientPreparations[sautPrep.ID][darkBrownSugar.ID]
	sautSkilletVPV := enums.PreparationVessels[sautPrep.ID][castIronSkillet.ID]
	sautWoodenSpoonVPI := enums.PreparationInstruments[sautPrep.ID][woodenSpoon.ID]

	// Simmer
	simmerStarAniseVIP := enums.IngredientPreparations[simmerPrep.ID][starAnise.ID]
	simmerCassiaBarkVIP := enums.IngredientPreparations[simmerPrep.ID][cassiaBark.ID]
	simmerLightSoySauceVIP := enums.IngredientPreparations[simmerPrep.ID][lightSoySauce.ID]
	simmerShaoxingWineVIP := enums.IngredientPreparations[simmerPrep.ID][shaoxingWine.ID]
	simmerWaterVIP := enums.IngredientPreparations[simmerPrep.ID][water.ID]
	simmerSkilletVPV := enums.PreparationVessels[simmerPrep.ID][castIronSkillet.ID]

	// Braise
	braiseChickenVIP := enums.IngredientPreparations[braisePrep.ID][chickenThighs.ID]
	braiseSkilletVPV := enums.PreparationVessels[braisePrep.ID][castIronSkillet.ID]
	braiseOvenVPV := enums.PreparationVessels[braisePrep.ID][oven.ID]
	braiseThermometerVPI := enums.PreparationInstruments[braisePrep.ID][thermometer.ID]

	// Measurement unit bridges
	saltTablespoonVIMU := enums.IngredientMeasurementUnits[salt.ID][tablespoonMeasurement.ID]
	msgTeaspoonVIMU := enums.IngredientMeasurementUnits[msg.ID][teaspoonMeasurement.ID]
	fiveSpiceTeaspoonVIMU := enums.IngredientMeasurementUnits[fiveSpice.ID][teaspoonMeasurement.ID]
	darkBrownSugarTablespoonVIMU := enums.IngredientMeasurementUnits[darkBrownSugar.ID][tablespoonMeasurement.ID]
	whitePepperTeaspoonVIMU := enums.IngredientMeasurementUnits[whitePepper.ID][teaspoonMeasurement.ID]
	chickenThighsPoundVIMU := enums.IngredientMeasurementUnits[chickenThighs.ID][poundMeasurement.ID]
	vegetableOilTablespoonVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][tablespoonMeasurement.ID]
	scallionsUnitVIMU := enums.IngredientMeasurementUnits[scallions.ID][unitMeasurement.ID]
	gingerUnitVIMU := enums.IngredientMeasurementUnits[ginger.ID][unitMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
	starAniseUnitVIMU := enums.IngredientMeasurementUnits[starAnise.ID][unitMeasurement.ID]
	cassiaBarkUnitVIMU := enums.IngredientMeasurementUnits[cassiaBark.ID][unitMeasurement.ID]
	lightSoySauceCupVIMU := enums.IngredientMeasurementUnits[lightSoySauce.ID][cupMeasurement.ID]
	shaoxingWineCupVIMU := enums.IngredientMeasurementUnits[shaoxingWine.ID][cupMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]

	// Step 0: In a small bowl, whisk together salt, MSG, 1/2 teaspoon five spice powder, 3 tablespoons dark brown sugar, and ground white pepper to combine. Set aside.
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   combinePrep.ID,
		Index:           0,
		Notes:           "In a small bowl, whisk together salt, MSG, 1/2 teaspoon five spice powder, 3 tablespoons dark brown sugar, and ground white pepper to combine. Set aside.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &combineSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTablespoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "Diamond Crystal kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &combineMSGVIP.ID,
				ValidIngredientMeasurementUnitID: &msgTeaspoonVIMU.ID,
				IngredientID:                     &msg.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "MSG",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &combineFiveSpiceVIP.ID,
				ValidIngredientMeasurementUnitID: &fiveSpiceTeaspoonVIMU.ID,
				IngredientID:                     &fiveSpice.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "five spice powder",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &combineDarkBrownSugarVIP.ID,
				ValidIngredientMeasurementUnitID: &darkBrownSugarTablespoonVIMU.ID,
				IngredientID:                     &darkBrownSugar.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "dark brown sugar",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &combineWhitePepperVIP.ID,
				ValidIngredientMeasurementUnitID: &whitePepperTeaspoonVIMU.ID,
				IngredientID:                     &whitePepper.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "ground white pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
				ValidPreparationInstrumentID: &combineWhiskVPI.ID,
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
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &combineSmallBowlVPV.ID,
				VesselID:                 &smallBowl.ID,
				Name:                     "small bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "dry brine mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Using a paper towel, pat chicken thighs dry
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dryPrep.ID,
		Index:           1,
		Notes:           "Using a paper towel, pat chicken thighs dry.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &dryChickenVIP.ID,
				ValidIngredientMeasurementUnitID: &chickenThighsPoundVIMU.ID,
				IngredientID:                     &chickenThighs.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "bone-in, skin-on chicken thighs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidPreparationInstrumentID: &dryPaperTowelsVPI.ID,
				InstrumentID:                 &paperTowels.ID,
				Name:                         "paper towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "dried chicken thighs",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 2: Season chicken generously on all sides with salt mixture
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           2,
		Notes:           "Season chicken generously on all sides with salt mixture.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "dried chicken thighs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &salt.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "dry brine mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step2ID,
				ValidPreparationInstrumentID: &seasonBareHandsVPI.ID,
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
				Name:                "seasoned chicken thighs",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 3: Transfer to a wire rack set in a rimmed 13- by 18-inch baking sheet
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           3,
		Notes:           "Transfer to a wire rack set in a rimmed 13- by 18-inch baking sheet.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "seasoned chicken thighs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &transferWireRackVPV.ID,
				VesselID:                 &wireRack.ID,
				Name:                     "wire rack",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &transferBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "rimmed 13- by 18-inch baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "chicken on wire rack",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: When ready to cook, adjust oven rack to middle position and preheat to 300°F (150°C)
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   preheatPrep.ID,
		Index:           4,
		Notes:           "When ready to cook, adjust oven rack to middle position and preheat to 300°F (150°C).",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &preheatOvenVPV.ID,
				VesselID:                 &oven.ID,
				Name:                     "oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "preheated oven at 300°F",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Using paper towels, pat chicken dry
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dryPrep.ID,
		Index:           5,
		Notes:           "Using paper towels, pat chicken dry.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &dryChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "chicken on wire rack",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step5ID,
				ValidPreparationInstrumentID: &dryPaperTowelsVPI.ID,
				InstrumentID:                 &paperTowels.ID,
				Name:                         "paper towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "dried chicken thighs",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 6: In a large cast iron or carbon steel skillet set over medium heat, heat vegetable oil until shimmering
	step6ID := identifiers.New()
	step6OilIngredientID := identifiers.New()
	step6CompletionConditionID := identifiers.New()
	shimmeringState := enums.IngredientStates["shimmering"]
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           6,
		Notes:           "In a large cast iron or carbon steel skillet set over medium heat, heat vegetable oil until shimmering.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step6OilIngredientID,
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &vegetableOilTablespoonVIMU.ID,
				IngredientID:                     &vegetableOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "neutral oil, such as vegetable",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &heatSkilletVPV.ID,
				VesselID:                 &castIronSkillet.ID,
				Name:                     "large cast iron or carbon steel skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "heated skillet with oil",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step6CompletionConditionID,
				BelongsToRecipeStep: step6ID,
				IngredientStateID:   shimmeringState.ID,
				Notes:               "Oil should shimmer when viewed",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step6CompletionConditionID,
						RecipeStepIngredient:                   step6OilIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 7: Working in batches if necessary, add chicken, skin-side-down, and cook without moving until well-browned and crispy, 4 to 6 minutes
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   panSearPrep.ID,
		Index:           7,
		Notes:           "Working in batches if necessary, add chicken, skin-side-down, and cook without moving until well-browned and crispy, 4 to 6 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](240), // 4 minutes
			Max: pointer.To[uint32](360), // 6 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &panSearChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "dried chicken thighs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step7ID,
				ValidPreparationInstrumentID: &panSearTongsVPI.ID,
				InstrumentID:                 &tongs.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &panSearSkilletVPV.ID,
				VesselID:                        &castIronSkillet.ID,
				Name:                            "heated skillet with oil",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "seared chicken thighs (skin-side)",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "skillet with seared chicken",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Flip chicken and cook lightly on second side, about 2 minutes
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   flipPrep.ID,
		Index:           8,
		Notes:           "Flip chicken and cook lightly on second side, about 2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120), // 2 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &flipChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "seared chicken thighs (skin-side)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step8ID,
				ValidPreparationInstrumentID: &flipTongsVPI.ID,
				InstrumentID:                 &tongs.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &flipSkilletVPV.ID,
				VesselID:                        &castIronSkillet.ID,
				Name:                            "skillet with seared chicken",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "seared chicken thighs (both sides)",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "skillet with seared chicken",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Transfer chicken to a large plate and set aside
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           9,
		Notes:           "Transfer chicken to a large plate and set aside.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "seared chicken thighs (both sides)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step9ID,
				ValidPreparationVesselID: &transferLargePlateVPV.ID,
				VesselID:                 &largePlate.ID,
				Name:                     "large plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "seared chicken on plate",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Reduce heat to medium-low. Add scallions, ginger, garlic, 1 teaspoon five spice powder, and 3 tablespoons dark brown sugar and cook, stirring, until vegetables are softened and starting to brown, 3 to 5 minutes
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   sautPrep.ID,
		Index:           10,
		Notes:           "Reduce heat to medium-low. Add scallions, ginger, garlic, 1 teaspoon five spice powder, and 3 tablespoons dark brown sugar and cook, stirring, until vegetables are softened and starting to brown, 3 to 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180), // 3 minutes
			Max: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &sautScallionsVIP.ID,
				ValidIngredientMeasurementUnitID: &scallionsUnitVIMU.ID,
				IngredientID:                     &scallions.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "scallions, green and white parts cut into 2-inch segments",
				QuantityNotes:                    "8 scallions",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 8,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &sautGingerVIP.ID,
				ValidIngredientMeasurementUnitID: &gingerUnitVIMU.ID,
				IngredientID:                     &ginger.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "ginger, peeled and thinly sliced",
				QuantityNotes:                    "One 2-inch piece",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &sautGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                cloveMeasurement.ID,
				Name:                             "garlic",
				QuantityNotes:                    "5 medium cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &sautFiveSpiceVIP.ID,
				ValidIngredientMeasurementUnitID: &fiveSpiceTeaspoonVIMU.ID,
				IngredientID:                     &fiveSpice.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "five spice powder",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &sautDarkBrownSugarVIP.ID,
				ValidIngredientMeasurementUnitID: &darkBrownSugarTablespoonVIMU.ID,
				IngredientID:                     &darkBrownSugar.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "dark brown sugar",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step10ID,
				ValidPreparationInstrumentID: &sautWoodenSpoonVPI.ID,
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &sautSkilletVPV.ID,
				VesselID:                        &castIronSkillet.ID,
				Name:                            "skillet with seared chicken",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "cooked aromatics",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "skillet with aromatics",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 11: Add star anise, cassia bark or cinnamon stick, soy sauce, Shaoxing wine, and water, and bring to a simmer over medium heat
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   simmerPrep.ID,
		Index:           11,
		Notes:           "Add star anise, cassia bark or cinnamon stick, soy sauce, Shaoxing wine, and water, and bring to a simmer over medium heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &scallions.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "cooked aromatics",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &simmerStarAniseVIP.ID,
				ValidIngredientMeasurementUnitID: &starAniseUnitVIMU.ID,
				IngredientID:                     &starAnise.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "star anise",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &simmerCassiaBarkVIP.ID,
				ValidIngredientMeasurementUnitID: &cassiaBarkUnitVIMU.ID,
				IngredientID:                     &cassiaBark.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "cassia bark or cinnamon stick",
				QuantityNotes:                    "One 2-inch piece",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &simmerLightSoySauceVIP.ID,
				ValidIngredientMeasurementUnitID: &lightSoySauceCupVIMU.ID,
				IngredientID:                     &lightSoySauce.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "light soy sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &simmerShaoxingWineVIP.ID,
				ValidIngredientMeasurementUnitID: &shaoxingWineCupVIMU.ID,
				IngredientID:                     &shaoxingWine.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "Shaoxing wine",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &simmerWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				IngredientID:                     &water.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &simmerSkilletVPV.ID,
				VesselID:                        &castIronSkillet.ID,
				Name:                            "skillet with aromatics",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "simmering braising liquid",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "skillet with braising liquid",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 12: Return chicken to pan skin-side-up, leaving the skin above the liquid but submerging most of the meat
	step12ID := identifiers.New()
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           12,
		Notes:           "Return chicken to pan skin-side-up, leaving the skin above the liquid but submerging most of the meat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "seared chicken on plate",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &transferSkilletVPV.ID,
				VesselID:                        &castIronSkillet.ID,
				Name:                            "skillet with braising liquid",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "chicken in braising liquid",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "skillet with chicken in braising liquid",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 13: Transfer skillet to oven and cook uncovered until chicken is cooked through and tender and registers at least 175°F (79°C), about 30 minutes
	step13ID := identifiers.New()
	step13 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step13ID,
		BelongsToRecipe: recipeID,
		PreparationID:   braisePrep.ID,
		Index:           13,
		Notes:           "Transfer skillet to oven and cook uncovered until chicken is cooked through and tender and registers at least 175°F (79°C), about 30 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](1800), // 30 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &braiseChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "chicken in braising liquid",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step13ID,
				ValidPreparationInstrumentID: &braiseThermometerVPI.ID,
				InstrumentID:                 &thermometer.ID,
				Name:                         "instant-read thermometer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &braiseSkilletVPV.ID,
				VesselID:                        &castIronSkillet.ID,
				Name:                            "skillet with chicken in braising liquid",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &braiseOvenVPV.ID,
				VesselID:                        &oven.ID,
				Name:                            "preheated oven at 300°F",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "braised chicken thighs",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "skillet with braised chicken",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 14: Remove pan from oven
	step14ID := identifiers.New()
	step14 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step14ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           14,
		Notes:           "Remove pan from oven.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "braised chicken thighs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &transferSkilletVPV.ID,
				VesselID:                        &castIronSkillet.ID,
				Name:                            "skillet with braised chicken",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step14ID,
				Name:                "soy sauce braised chicken thighs",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Create prep task for dry-brining chicken ahead of time
	prepTask1ID := identifiers.New()
	prepTask1 := &mealplanning.RecipePrepTaskDatabaseCreationInput{
		ID:                          prepTask1ID,
		BelongsToRecipe:             recipeID,
		Name:                        "Dry-brine chicken thighs",
		Description:                 "The chicken can be dry-brined at least 8 hours and up to 72 hours in advance. Store in the refrigerator uncovered on a wire rack set in a rimmed baking sheet.",
		Notes:                       "Dry-brining yields tender, juicy meat with crackly skin.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the seasoned chicken on a wire rack set in a rimmed baking sheet in the refrigerator, uncovered, for at least 8 hours and up to 72 hours.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeWireRack,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 28800,                      // 8 hours
			Max: pointer.To[uint32](259200), // 72 hours
		},
		TaskSteps: []*mealplanning.RecipePrepTaskStepDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: step0ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step1ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step2ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step3ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: true},
		},
	}

	soySauceBraisedChickenThighsRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Soy Sauce–Braised Chicken Thighs",
		Slug:                "soy-sauce-braised-chicken-thighs",
		Source:              "https://www.seriouseats.com/soy-sauce-braised-chicken-thighs-recipe-8737800",
		Description:         "Bathed in a fragrant blend of soy sauce, brown sugar, and warm spices, these tender chicken thighs evoke the flavors of classic Cantonese soy sauce chicken. The chicken is seared until crispy, then gently braised in the oven.",
		YieldsComponentType: mealplanning.MealComponentTypesMain,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13, step14,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskDatabaseCreationInput{
			prepTask1,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		soySauceBraisedChickenThighsRecipe,
	}
}
