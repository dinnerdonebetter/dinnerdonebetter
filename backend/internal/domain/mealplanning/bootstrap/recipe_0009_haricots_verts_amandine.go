package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// HaricotsVertsAmandineRecipe creates the Haricots Verts Amandine (French-Style Green Beans With Almonds) recipe.
// Source: https://www.seriouseats.com/green-beans-amandine-french-almondine-recipe
func HaricotsVertsAmandineRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	boilPrep := enums.Preparations["boil"]
	blanchPrep := enums.Preparations["blanch"]
	shockPrep := enums.Preparations["shock"]
	drainPrep := enums.Preparations["drain"]
	dryPrep := enums.Preparations["dry"]
	heatPrep := enums.Preparations["heat"]
	cookPrep := enums.Preparations["cook"]
	stirPrep := enums.Preparations["stir"]
	emulsifyPrep := enums.Preparations["emulsify"]
	seasonPrep := enums.Preparations["season"]
	tossPrep := enums.Preparations["toss"]
	transferPrep := enums.Preparations["transfer"]
	trimPrep := enums.Preparations["trim"]

	// Get ingredients
	greenBeans := enums.Ingredients["green beans"]
	butter := enums.Ingredients["butter"]
	sliveredAlmonds := enums.Ingredients["slivered almonds"]
	garlic := enums.Ingredients["garlic"]
	shallot := enums.Ingredients["shallot"]
	lemon := enums.Ingredients["lemon"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	water := enums.Ingredients["water"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	ounceMeasurement := enums.MeasurementUnits["ounce"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]

	// Get instruments
	wireMeshSpider := enums.Instruments["wire mesh spider"]
	paperTowels := enums.Instruments["paper towels"]
	kitchenTowels := enums.Instruments["kitchen towels"]
	rubberSpatula := enums.Instruments["rubber spatula"]
	knife := enums.Instruments["knife"]

	// Get vessels
	pot := enums.Vessels["pot"]
	largeBowl := enums.Vessels["large bowl"]
	mediumSkillet := enums.Vessels["medium skillet"]
	servingPlatter := enums.Vessels["serving platter"]
	colander := enums.Vessels["colander"]
	cuttingBoard := enums.Vessels["cutting board"]

	// Get ingredient states for completion conditions
	toastedState := enums.IngredientStates["toasted"]

	// Get bridge table entries
	// Boil
	boilWaterVIP := enums.IngredientPreparations[boilPrep.ID][water.ID]
	boilSaltVIP := enums.IngredientPreparations[boilPrep.ID][salt.ID]
	boilPotVPV := enums.PreparationVessels[boilPrep.ID][pot.ID]

	// Blanch
	blanchGreenBeansVIP := enums.IngredientPreparations[blanchPrep.ID][greenBeans.ID]
	blanchPotVPV := enums.PreparationVessels[blanchPrep.ID][pot.ID]
	blanchSpiderVPI := enums.PreparationInstruments[blanchPrep.ID][wireMeshSpider.ID]

	// Shock
	shockGreenBeansVIP := enums.IngredientPreparations[shockPrep.ID][greenBeans.ID]
	shockLargeBowlVPV := enums.PreparationVessels[shockPrep.ID][largeBowl.ID]
	shockSpiderVPI := enums.PreparationInstruments[shockPrep.ID][wireMeshSpider.ID]

	// Drain
	drainGreenBeansVIP := enums.IngredientPreparations[drainPrep.ID][greenBeans.ID]
	drainColanderVPV := enums.PreparationVessels[drainPrep.ID][colander.ID]

	// Dry
	dryGreenBeansVIP := enums.IngredientPreparations[dryPrep.ID][greenBeans.ID]
	dryPaperTowelsVPI := enums.PreparationInstruments[dryPrep.ID][paperTowels.ID]
	dryKitchenTowelsVPI := enums.PreparationInstruments[dryPrep.ID][kitchenTowels.ID]

	// Trim
	trimGreenBeansVIP := enums.IngredientPreparations[trimPrep.ID][greenBeans.ID]
	trimKnifeVPI := enums.PreparationInstruments[trimPrep.ID][knife.ID]
	trimCuttingBoardVPV := enums.PreparationVessels[trimPrep.ID][cuttingBoard.ID]

	// Heat
	heatButterVIP := enums.IngredientPreparations[heatPrep.ID][butter.ID]
	heatAlmondsVIP := enums.IngredientPreparations[heatPrep.ID][sliveredAlmonds.ID]
	heatMediumSkilletVPV := enums.PreparationVessels[heatPrep.ID][mediumSkillet.ID]
	heatSpatulaVPI := enums.PreparationInstruments[heatPrep.ID][rubberSpatula.ID]

	// Cook
	cookGarlicVIP := enums.IngredientPreparations[cookPrep.ID][garlic.ID]
	cookShallotVIP := enums.IngredientPreparations[cookPrep.ID][shallot.ID]
	cookSkilletVPV := enums.PreparationVessels[cookPrep.ID][mediumSkillet.ID]
	cookSpatulaVPI := enums.PreparationInstruments[cookPrep.ID][rubberSpatula.ID]

	// Stir
	stirLemonVIP := enums.IngredientPreparations[stirPrep.ID][lemon.ID]
	stirWaterVIP := enums.IngredientPreparations[stirPrep.ID][water.ID]
	stirSkilletVPV := enums.PreparationVessels[stirPrep.ID][mediumSkillet.ID]
	stirLargeBowlVPV := enums.PreparationVessels[stirPrep.ID][largeBowl.ID]

	// Emulsify
	emulsifySkilletVPV := enums.PreparationVessels[emulsifyPrep.ID][mediumSkillet.ID]

	// Season
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	seasonSkilletVPV := enums.PreparationVessels[seasonPrep.ID][mediumSkillet.ID]

	// Toss
	tossGreenBeansVIP := enums.IngredientPreparations[tossPrep.ID][greenBeans.ID]
	tossSkilletVPV := enums.PreparationVessels[tossPrep.ID][mediumSkillet.ID]

	// Transfer
	transferGreenBeansVIP := enums.IngredientPreparations[transferPrep.ID][greenBeans.ID]
	transferServingPlatterVPV := enums.PreparationVessels[transferPrep.ID][servingPlatter.ID]

	// Measurement unit bridges
	greenBeansPoundVIMU := enums.IngredientMeasurementUnits[greenBeans.ID][poundMeasurement.ID]
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]
	almondsOunceVIMU := enums.IngredientMeasurementUnits[sliveredAlmonds.ID][ounceMeasurement.ID]
	garlicUnitVIMU := enums.IngredientMeasurementUnits[garlic.ID][unitMeasurement.ID]
	shallotUnitVIMU := enums.IngredientMeasurementUnits[shallot.ID][unitMeasurement.ID]
	lemonTablespoonVIMU := enums.IngredientMeasurementUnits[lemon.ID][tablespoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	pepperTeaspoonVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]
	waterTablespoonVIMU := enums.IngredientMeasurementUnits[water.ID][tablespoonMeasurement.ID]

	// Step 0: Bring a large pot of salted water to a boil
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   boilPrep.ID,
		Index:           0,
		Notes:           "Bring a large pot of salted water to a boil.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &boilWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterTablespoonVIMU.ID,
				IngredientID:                     &water.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "water",
				QuantityNotes:                    "enough to cover green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &boilSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "kosher salt",
				QuantityNotes:                    "generously salted",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &boilPotVPV.ID,
				VesselID:                 &pot.ID,
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
				Name:                "boiling salted water",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Prepare an ice bath
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           1,
		Notes:           "Prepare an ice bath by filling a large bowl with ice and cold water.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &stirWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterTablespoonVIMU.ID,
				IngredientID:                     &water.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "ice water",
				QuantityNotes:                    "enough to submerge green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step1ID,
				ValidPreparationVesselID: &stirLargeBowlVPV.ID,
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
				BelongsToRecipeStep: step1ID,
				Name:                "ice bath",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: (OPTIONAL) Trim the green beans
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   trimPrep.ID,
		Index:           2,
		Optional:        true,
		Notes:           "Trim the ends off the green beans. This step is optional if using pre-trimmed green beans.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &trimGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				IngredientID:                     &greenBeans.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step2ID,
				ValidPreparationInstrumentID: &trimKnifeVPI.ID,
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
				BelongsToRecipeStep:      step2ID,
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
				BelongsToRecipeStep: step2ID,
				Name:                "trimmed green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Blanch the green beans
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   blanchPrep.ID,
		Index:           3,
		Notes:           "Add green beans to boiling water and cook until tender-crisp, about 3 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180), // 3 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &blanchGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				IngredientID:                     &greenBeans.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "trimmed green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &water.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "boiling salted water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step3ID,
				ValidPreparationInstrumentID: &blanchSpiderVPI.ID,
				InstrumentID:                 &wireMeshSpider.ID,
				Name:                         "wire mesh spider",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &blanchPotVPV.ID,
				VesselID:                 &pot.ID,
				Name:                     "large pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "blanched green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Transfer to ice bath (shock)
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   shockPrep.ID,
		Index:           4,
		Notes:           "Transfer to ice bath using a wire mesh spider or tongs. Allow to chill completely.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &shockGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				IngredientID:                     &greenBeans.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "blanched green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &shockSpiderVPI.ID,
				InstrumentID:                 &wireMeshSpider.ID,
				Name:                         "wire mesh spider",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &shockLargeBowlVPV.ID,
				VesselID:                        &largeBowl.ID,
				Name:                            "ice bath",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "chilled green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Drain green beans
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   drainPrep.ID,
		Index:           5,
		Notes:           "Drain the green beans thoroughly.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step5ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &drainGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				IngredientID:                     &greenBeans.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "chilled green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step5ID,
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
				BelongsToRecipeStep: step5ID,
				Name:                "drained green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Dry green beans
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dryPrep.ID,
		Index:           6,
		Notes:           "Dry thoroughly with kitchen towels or paper towels.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &dryGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				IngredientID:                     &greenBeans.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "drained green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: &dryKitchenTowelsVPI.ID,
				InstrumentID:                 &kitchenTowels.ID,
				Name:                         "kitchen towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
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
				BelongsToRecipeStep: step6ID,
				Name:                "dried green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Heat butter and almonds in skillet, toast until deeply browned (combined step)
	step7ID := identifiers.New()
	step7AlmondsIngredientID := identifiers.New()
	step7CompletionConditionID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           7,
		Notes:           "In a medium skillet, heat butter and almonds over medium-low heat and cook, stirring frequently, until almonds are deeply browned and nutty, about 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ValidIngredientPreparationID:     &heatButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				IngredientID:                     &butter.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "unsalted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               step7AlmondsIngredientID,
				BelongsToRecipeStep:              step7ID,
				ValidIngredientPreparationID:     &heatAlmondsVIP.ID,
				ValidIngredientMeasurementUnitID: &almondsOunceVIMU.ID,
				IngredientID:                     &sliveredAlmonds.ID,
				MeasurementUnitID:                ounceMeasurement.ID,
				Name:                             "slivered almonds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step7ID,
				ValidPreparationInstrumentID: &heatSpatulaVPI.ID,
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
				BelongsToRecipeStep:      step7ID,
				ValidPreparationVesselID: &heatMediumSkilletVPV.ID,
				VesselID:                 &mediumSkillet.ID,
				Name:                     "medium skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "toasted almonds in brown butter",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "heated skillet with brown butter",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step7CompletionConditionID,
				BelongsToRecipeStep: step7ID,
				IngredientStateID:   toastedState.ID,
				Notes:               "Almonds should be deeply browned and nutty",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step7CompletionConditionID,
						RecipeStepIngredient:                   step7AlmondsIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 8: Add garlic and shallot and cook
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cookPrep.ID,
		Index:           8,
		Notes:           "Add garlic and shallot and cook, stirring, until lightly browned, about 2 minutes longer.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120), // 2 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &sliveredAlmonds.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "toasted almonds in brown butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &cookGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicUnitVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "garlic, thinly sliced",
				QuantityNotes:                    "2 medium cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &cookShallotVIP.ID,
				ValidIngredientMeasurementUnitID: &shallotUnitVIMU.ID,
				IngredientID:                     &shallot.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "shallot, thinly sliced",
				QuantityNotes:                    "1 medium",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step8ID,
				ValidPreparationInstrumentID: &cookSpatulaVPI.ID,
				InstrumentID:                 &rubberSpatula.ID,
				Name:                         "rubber spatula",
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
				ValidPreparationVesselID:        &cookSkilletVPV.ID,
				VesselID:                        &mediumSkillet.ID,
				Name:                            "heated skillet with brown butter",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "almond mixture with garlic and shallot",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "skillet with aromatics",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Add lemon juice and water
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           9,
		Notes:           "Add lemon juice, along with a tablespoon or two of water.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &sliveredAlmonds.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "almond mixture with garlic and shallot",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step9ID,
				ValidIngredientPreparationID:     &stirLemonVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonTablespoonVIMU.ID,
				IngredientID:                     &lemon.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "lemon juice",
				QuantityNotes:                    "juice from 1 lemon",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1.5,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step9ID,
				ValidIngredientPreparationID:     &stirWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterTablespoonVIMU.ID,
				IngredientID:                     &water.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
					Max: pointer.To[float32](2),
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &stirSkilletVPV.ID,
				VesselID:                        &mediumSkillet.ID,
				Name:                            "skillet with aromatics",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "sauce mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "skillet with sauce",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Emulsify the sauce
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   emulsifyPrep.ID,
		Index:           10,
		Notes:           "Increase heat to high and stir and shake pan rapidly to emulsify, about 30 seconds. The sauce should have a glossy sheen and not appear watery or greasy. If it's still watery, continue to simmer and shake. If it looks greasy, add another tablespoon of water to re-emulsify.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30), // 30 seconds
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &butter.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "sauce mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &emulsifySkilletVPV.ID,
				VesselID:                        &mediumSkillet.ID,
				Name:                            "skillet with sauce",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "emulsified brown butter sauce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "skillet with emulsified sauce",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 11: Season the sauce
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           11,
		Notes:           "When sauce is ready, remove from heat and season to taste with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &butter.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "emulsified brown butter sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "kosher salt",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &seasonPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperTeaspoonVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "freshly ground black pepper",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &seasonSkilletVPV.ID,
				VesselID:                        &mediumSkillet.ID,
				Name:                            "skillet with emulsified sauce",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "seasoned brown butter sauce with almonds",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "skillet with seasoned sauce",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 12: Add beans and toss
	step12ID := identifiers.New()
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   tossPrep.ID,
		Index:           12,
		Notes:           "Add beans to pan with sauce and toss to coat and combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &butter.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "seasoned brown butter sauce with almonds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step12ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &tossGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				IngredientID:                     &greenBeans.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "dried green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &tossSkilletVPV.ID,
				VesselID:                        &mediumSkillet.ID,
				Name:                            "skillet with seasoned sauce",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "green beans tossed in sauce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "skillet with green beans",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 13: Cook until heated through
	step13ID := identifiers.New()
	step13 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step13ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cookPrep.ID,
		Index:           13,
		Notes:           "Return to medium heat and cook, tossing, until heated through, about 1 minute.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60), // 1 minute
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &greenBeans.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "green beans tossed in sauce",
				Quantity: types.Float32RangeWithOptionalMax{
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
				ValidPreparationVesselID:        &cookSkilletVPV.ID,
				VesselID:                        &mediumSkillet.ID,
				Name:                            "skillet with green beans",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "heated green beans amandine",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 14: Transfer to serving platter
	step14ID := identifiers.New()
	step14 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step14ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           14,
		Notes:           "Serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step14ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &transferGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				IngredientID:                     &greenBeans.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "heated green beans amandine",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step14ID,
				ValidPreparationVesselID: &transferServingPlatterVPV.ID,
				VesselID:                 &servingPlatter.ID,
				Name:                     "serving platter",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step14ID,
				Name:                "haricots verts amandine",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Create prep task for blanching beans ahead of time
	prepTask1ID := identifiers.New()
	prepTask1 := &mealplanning.RecipePrepTaskDatabaseCreationInput{
		ID:                          prepTask1ID,
		BelongsToRecipe:             recipeID,
		Name:                        "Blanch green beans",
		Description:                 "The green beans can be blanched and dried several days in advance. Store in an airtight container in the refrigerator.",
		Notes:                       "Blanching ahead of time makes day-of preparation much faster.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the blanched and dried green beans in an airtight container in the refrigerator for up to 3 days.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeAirtightContainer,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 0,
		},
		TaskSteps: []*mealplanning.RecipePrepTaskStepDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: step0ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step1ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step2ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step3ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step4ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step5ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step6ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: true},
		},
	}

	// Create prep task for making sauce ahead of time
	prepTask2ID := identifiers.New()
	prepTask2 := &mealplanning.RecipePrepTaskDatabaseCreationInput{
		ID:                          prepTask2ID,
		BelongsToRecipe:             recipeID,
		Name:                        "Prepare brown butter almond sauce",
		Description:                 "The sauce can be prepared several days in advance. To finish, reheat the sauce in a skillet over high heat with a tablespoon of water until it melts back into a liquid.",
		Notes:                       "Store the sauce separately from the blanched beans.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the sauce in an airtight container in the refrigerator for up to 3 days.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeAirtightContainer,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 0,
		},
		TaskSteps: []*mealplanning.RecipePrepTaskStepDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: step7ID, BelongsToRecipePrepTask: prepTask2ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step8ID, BelongsToRecipePrepTask: prepTask2ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step9ID, BelongsToRecipePrepTask: prepTask2ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step10ID, BelongsToRecipePrepTask: prepTask2ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step11ID, BelongsToRecipePrepTask: prepTask2ID, SatisfiesRecipeStep: true},
		},
	}

	haricotsVertsAmandineRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Haricots Verts Amandine",
		Slug:                "haricots-verts-amandine",
		Source:              "https://www.seriouseats.com/green-beans-amandine-french-almondine-recipe",
		Description:         "The classic French side dish of green beans with almonds, featuring blanched tender-crisp green beans tossed in a brown butter sauce with deeply toasted almonds, garlic, shallots, and a bright lemon finish.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
			Max: pointer.To[float32](6),
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13, step14,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskDatabaseCreationInput{
			prepTask1, prepTask2,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		haricotsVertsAmandineRecipe,
	}
}
