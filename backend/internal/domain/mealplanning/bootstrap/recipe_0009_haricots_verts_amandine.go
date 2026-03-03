package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// HaricotsVertsAmandineRecipe creates the Haricots Verts Amandine (French-Style Green Beans With Almonds) recipe.
// Source: https://www.seriouseats.com/green-beans-amandine-french-almondine-recipe
func HaricotsVertsAmandineRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	grindPrep := enums.Preparations["grind"]
	boilPrep := enums.Preparations["boil"]
	blanchPrep := enums.Preparations["blanch"]
	shockPrep := enums.Preparations["shock"]
	drainPrep := enums.Preparations["drain"]
	dryPrep := enums.Preparations["dry"]
	heatPrep := enums.Preparations["heat"]
	toastPrep := enums.Preparations["toast"]
	cookPrep := enums.Preparations["cook"]
	stirPrep := enums.Preparations["stir"]
	adjustPrep := enums.Preparations["adjust"]
	removeFromHeatPrep := enums.Preparations["remove from heat"]
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
	wholePeppercorns := enums.Ingredients["whole black peppercorns"]
	water := enums.Ingredients["water"]
	iceCubes := enums.Ingredients["ice cubes"]

	// Get measurement units
	gramMeasurement := enums.MeasurementUnits["gram"]
	poundMeasurement := enums.MeasurementUnits["pound"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	ounceMeasurement := enums.MeasurementUnits["ounce"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	trayMeasurement := enums.MeasurementUnits["tray"]
	cupMeasurement := enums.MeasurementUnits["cup"]

	// Get instruments
	mortarAndPestle := enums.Instruments["mortar and pestle"]
	spiceGrinder := enums.Instruments["spice grinder"]
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

	// Toast (for almonds)
	toastAlmondsVIP := enums.IngredientPreparations[toastPrep.ID][sliveredAlmonds.ID]
	toastMediumSkilletVPV := enums.PreparationVessels[toastPrep.ID][mediumSkillet.ID]
	toastSpatulaVPI := enums.PreparationInstruments[toastPrep.ID][rubberSpatula.ID]

	// Adjust (for heat)
	adjustMediumSkilletVPV := enums.PreparationVessels[adjustPrep.ID][mediumSkillet.ID]

	// Remove from heat
	removeFromHeatMediumSkilletVPV := enums.PreparationVessels[removeFromHeatPrep.ID][mediumSkillet.ID]

	// Cook
	cookGarlicVIP := enums.IngredientPreparations[cookPrep.ID][garlic.ID]
	cookShallotVIP := enums.IngredientPreparations[cookPrep.ID][shallot.ID]
	cookSkilletVPV := enums.PreparationVessels[cookPrep.ID][mediumSkillet.ID]
	cookSpatulaVPI := enums.PreparationInstruments[cookPrep.ID][rubberSpatula.ID]

	// Stir
	stirLemonVIP := enums.IngredientPreparations[stirPrep.ID][lemon.ID]
	stirWaterVIP := enums.IngredientPreparations[stirPrep.ID][water.ID]
	stirIceCubesVIP := enums.IngredientPreparations[stirPrep.ID][iceCubes.ID]
	stirSkilletVPV := enums.PreparationVessels[stirPrep.ID][mediumSkillet.ID]
	stirLargeBowlVPV := enums.PreparationVessels[stirPrep.ID][largeBowl.ID]

	// Emulsify
	emulsifySkilletVPV := enums.PreparationVessels[emulsifyPrep.ID][mediumSkillet.ID]

	// Season
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
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
	waterTablespoonVIMU := enums.IngredientMeasurementUnits[water.ID][tablespoonMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	iceCubesTrayVIMU := enums.IngredientMeasurementUnits[iceCubes.ID][trayMeasurement.ID]

	// Grind preparation bridges
	grindPeppercornsVIP := enums.IngredientPreparations[grindPrep.ID][wholePeppercorns.ID]
	peppercornsGramVIMU := enums.IngredientMeasurementUnits[wholePeppercorns.ID][gramMeasurement.ID]
	grindMortarAndPestleVPI := enums.PreparationInstruments[grindPrep.ID][mortarAndPestle.ID]
	grindSpiceGrinderVPI := enums.PreparationInstruments[grindPrep.ID][spiceGrinder.ID]

	// Step 0: Grind whole black peppercorns
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        grindPrep.ID,
		Index:                0,
		ExplicitInstructions: "Using a mortar and pestle or spice grinder, coarsely grind the whole black peppercorns.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &grindPeppercornsVIP.ID,
				ValidIngredientMeasurementUnitID: &peppercornsGramVIMU.ID,
				Name:                             "whole black peppercorns",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &grindMortarAndPestleVPI.ID,
				Name:                         "mortar and pestle",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
				Index:       pointer.To[uint16](0),
				OptionIndex: 0,
			},
			{
				ValidPreparationInstrumentID: &grindSpiceGrinderVPI.ID,
				Name:                         "spice grinder",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
				Index:       pointer.To[uint16](0),
				OptionIndex: 1,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "freshly ground black pepper",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Bring a large pot of salted water to a boil
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        boilPrep.ID,
		Index:                1,
		ExplicitInstructions: "Bring a large pot of salted water to a boil.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &boilWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterTablespoonVIMU.ID,
				Name:                             "water",
				QuantityNotes:                    "enough to cover green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &boilSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "kosher salt",
				QuantityNotes:                    "generously salted",
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
				Name:              "boiling salted water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot with boiling salted water",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Prepare an ice bath
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                2,
		ExplicitInstructions: "Prepare an ice bath by filling a large bowl with ice and cold water.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &stirWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				Name:                             "cold water",
				QuantityNotes:                    "enough to submerge green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ValidIngredientPreparationID:     &stirIceCubesVIP.ID,
				ValidIngredientMeasurementUnitID: &iceCubesTrayVIMU.ID,
				Name:                             "ice cubes",
				QuantityNotes:                    "about 1 tray",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &stirLargeBowlVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "ice bath",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: (OPTIONAL) Trim the green beans
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        trimPrep.ID,
		Index:                3,
		Optional:             true,
		ExplicitInstructions: "Trim the ends off the green beans. This step is optional if using pre-trimmed green beans.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &trimGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				Name:                             "green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &trimKnifeVPI.ID,
				Name:                         "knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &trimCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "trimmed green beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Blanch the green beans
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        blanchPrep.ID,
		Index:                4,
		ExplicitInstructions: "Add the green beans to the boiling water and cook until tender-crisp, about 3 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180), // 3 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &blanchGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				Name:                             "trimmed green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "boiling salted water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &blanchSpiderVPI.ID,
				Name:                         "wire mesh spider",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &blanchPotVPV.ID,
				Name:                            "pot with boiling salted water",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "blanched green beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Transfer to ice bath (shock)
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        shockPrep.ID,
		Index:                5,
		ExplicitInstructions: "Transfer to the ice bath using a wire mesh spider or tongs. Allow to chill completely.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &shockGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				Name:                             "blanched green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &shockSpiderVPI.ID,
				Name:                         "wire mesh spider",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &shockLargeBowlVPV.ID,
				Name:                            "ice bath",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "chilled green beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Drain green beans
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        drainPrep.ID,
		Index:                6,
		ExplicitInstructions: "Drain the green beans thoroughly.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &drainGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				Name:                             "chilled green beans",
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
				Name:              "drained green beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Dry green beans
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        dryPrep.ID,
		Index:                7,
		ExplicitInstructions: "Dry thoroughly with kitchen towels or paper towels.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &dryGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				Name:                             "drained green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &dryKitchenTowelsVPI.ID,
				Name:                         "kitchen towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationInstrumentID: &dryPaperTowelsVPI.ID,
				Name:                         "paper towels",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "dried green beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Heat butter and almonds in skillet
	brownedState := enums.IngredientStates["browned"]
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        heatPrep.ID,
		Index:                8,
		ExplicitInstructions: "In a medium skillet, heat the butter and almonds over medium-low heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &heatButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				Name:                             "unsalted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ValidIngredientPreparationID:     &heatAlmondsVIP.ID,
				ValidIngredientMeasurementUnitID: &almondsOunceVIMU.ID,
				Name:                             "slivered almonds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &heatSpatulaVPI.ID,
				Name:                         "rubber spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &heatMediumSkilletVPV.ID,
				Name:                     "medium skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "butter and almonds in skillet",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "heated skillet with butter and almonds",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Toast almonds until deeply browned
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        toastPrep.ID,
		Index:                9,
		ExplicitInstructions: "Cook, stirring frequently, until the almonds are deeply browned and nutty, about 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &toastAlmondsVIP.ID,
				ValidIngredientMeasurementUnitID: &almondsOunceVIMU.ID,
				Name:                             "butter and almonds in skillet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &toastSpatulaVPI.ID,
				Name:                         "rubber spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &toastMediumSkilletVPV.ID,
				Name:                            "heated skillet with butter and almonds",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "toasted almonds in brown butter",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "heated skillet with brown butter",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: brownedState.ID,
				Notes:             "Butter should be browned and nutty",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
			{
				IngredientStateID: toastedState.ID,
				Notes:             "Almonds should be deeply browned and nutty",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
	}

	// Step 10: Add garlic and shallot and cook
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        cookPrep.ID,
		Index:                10,
		ExplicitInstructions: "Add the garlic and shallot and cook, stirring, until lightly browned, about 2 minutes longer.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120), // 2 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "toasted almonds in brown butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &cookGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicUnitVIMU.ID,
				Name:                             "garlic, thinly sliced",
				QuantityNotes:                    "2 medium cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &cookShallotVIP.ID,
				ValidIngredientMeasurementUnitID: &shallotUnitVIMU.ID,
				Name:                             "shallot, thinly sliced",
				QuantityNotes:                    "1 medium",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &cookSpatulaVPI.ID,
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
				ValidPreparationVesselID:        &cookSkilletVPV.ID,
				Name:                            "heated skillet with brown butter",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "almond mixture with garlic and shallot",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet with aromatics",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 11: Add lemon juice and water
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                11,
		ExplicitInstructions: "Add lemon juice, along with a tablespoon or two of water.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "almond mixture with garlic and shallot",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &stirLemonVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonTablespoonVIMU.ID,
				Name:                             "lemon juice",
				QuantityNotes:                    "juice from 1 lemon",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1.5,
				},
			},
			{
				ValidIngredientPreparationID:     &stirWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterTablespoonVIMU.ID,
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
					Max: pointer.To[float32](2),
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &stirSkilletVPV.ID,
				Name:                            "skillet with aromatics",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "sauce mixture",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet with sauce",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 12: Increase heat to high
	step12 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        adjustPrep.ID,
		Index:                12,
		ExplicitInstructions: "Increase the heat to high.",
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &adjustMediumSkilletVPV.ID,
				Name:                            "skillet with sauce",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "skillet with sauce over high heat",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 13: Emulsify the sauce
	desiredConsistencyState := enums.IngredientStates["at desired consistency"]
	step13 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        emulsifyPrep.ID,
		Index:                13,
		ExplicitInstructions: "Stir and shake the pan rapidly to emulsify, about 30 seconds. The sauce should have a glossy sheen and not appear watery or greasy. If it's still watery, continue to simmer and shake. If it looks greasy, add another tablespoon of water to re-emulsify.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30), // 30 seconds
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "sauce mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &emulsifySkilletVPV.ID,
				Name:                            "skillet with sauce over high heat",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "emulsified brown butter sauce",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet with emulsified sauce",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: desiredConsistencyState.ID,
				Notes:             "Sauce should have a glossy sheen and be emulsified, not watery or greasy",
				Ingredients:       []uint64{0}, // Index of sauce ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 14: Remove from heat
	step14 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        removeFromHeatPrep.ID,
		Index:                14,
		ExplicitInstructions: "Remove from heat.",
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &removeFromHeatMediumSkilletVPV.ID,
				Name:                            "skillet with emulsified sauce",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "skillet with emulsified sauce, off heat",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 15: Season the sauce
	step15 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        seasonPrep.ID,
		Index:                15,
		ExplicitInstructions: "Season to taste with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "emulsified brown butter sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "kosher salt",
				QuantityNotes:                    "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "freshly ground black pepper",
				QuantityNotes:                   "to taste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &seasonSkilletVPV.ID,
				Name:                            "skillet with emulsified sauce, off heat",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seasoned brown butter sauce with almonds",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet with seasoned sauce",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 16: Add beans and toss
	step16 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        tossPrep.ID,
		Index:                16,
		ExplicitInstructions: "Add the beans to the pan with the sauce and toss to coat and combine.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](15),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "seasoned brown butter sauce with almonds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &tossGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				Name:                             "dried green beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](15),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &tossSkilletVPV.ID,
				Name:                            "skillet with seasoned sauce",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "green beans tossed in sauce",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "skillet with green beans",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 17: Cook until heated through
	step17 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        cookPrep.ID,
		Index:                17,
		ExplicitInstructions: "Cook over medium heat, tossing, until heated through, about 1 minute.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60), // 1 minute
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](16),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "green beans tossed in sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](16),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &cookSkilletVPV.ID,
				Name:                            "skillet with green beans",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "heated green beans amandine",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 18: Transfer to serving platter
	step18 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        transferPrep.ID,
		Index:                18,
		ExplicitInstructions: "Transfer to serving platter.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](17),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &transferGreenBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &greenBeansPoundVIMU.ID,
				Name:                             "heated green beans amandine",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &transferServingPlatterVPV.ID,
				Name:                     "serving platter",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "haricots verts amandine",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Create prep task for blanching beans ahead of time
	prepTask1 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
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
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 1, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 2, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 3, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 4, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 5, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 6, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 7, SatisfiesRecipeStep: true},
		},
	}

	// Create prep task for making sauce ahead of time
	prepTask2 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
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
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 8, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 9, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 10, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 11, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 12, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 13, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 14, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 15, SatisfiesRecipeStep: true},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
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
			Steps:             []*mealplanning.RecipeStepCreationRequestInput{step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13, step14, step15, step16, step17, step18},
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{prepTask1, prepTask2},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
