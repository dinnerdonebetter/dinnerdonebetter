package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ChickenVermicelliSoupRecipe creates the Chicken and Vermicelli Soup with Lime recipe.
// Source: https://cooking.nytimes.com/recipes/1026337-chicken-and-vermicelli-soup-with-lime
func ChickenVermicelliSoupRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	meltPrep := enums.Preparations["melt"]
	sautePrep := enums.Preparations["sauté"]
	sprinklePrep := enums.Preparations["sprinkle"]
	addPrep := enums.Preparations["add"]
	cookPrep := enums.Preparations["cook"]
	placePrep := enums.Preparations["place"]
	seasonPrep := enums.Preparations["season"]
	simmerPrep := enums.Preparations["simmer"]
	transferPrep := enums.Preparations["transfer"]
	shredPrep := enums.Preparations["shred"]
	stirPrep := enums.Preparations["stir"]

	// Get ingredients
	butter := enums.Ingredients["butter"]
	onion := enums.Ingredients["onion"]
	turmeric := enums.Ingredients["turmeric"]
	potato := enums.Ingredients["potato"]
	carrot := enums.Ingredients["carrot"]
	tomatoPaste := enums.Ingredients["tomato paste"]
	chickenBreast := enums.Ingredients["chicken breast"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	water := enums.Ingredients["water"]
	vermicelli := enums.Ingredients["vermicelli"]
	parsley := enums.Ingredients["parsley"]
	limeJuice := enums.Ingredients["lime juice"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	woodenSpoon := enums.Instruments["wooden spoon"]
	fork := enums.Instruments["fork"]

	// Get vessels
	largePot := enums.Vessels["pot"]
	mediumBowl := enums.Vessels["medium bowl"]

	// Get bridge table entries
	meltButterVIP := enums.IngredientPreparations[meltPrep.ID][butter.ID]
	meltPotVPV := enums.PreparationVessels[meltPrep.ID][largePot.ID]

	sauteOnionVIP := enums.IngredientPreparations[sautePrep.ID][onion.ID]
	sautePotVPV := enums.PreparationVessels[sautePrep.ID][largePot.ID]

	sprinkleTurmericVIP := enums.IngredientPreparations[sprinklePrep.ID][turmeric.ID]
	sprinklePotVPV := enums.PreparationVessels[sprinklePrep.ID][largePot.ID]

	addPotatoVIP := enums.IngredientPreparations[addPrep.ID][potato.ID]
	addCarrotVIP := enums.IngredientPreparations[addPrep.ID][carrot.ID]
	addPotVPV := enums.PreparationVessels[addPrep.ID][largePot.ID]

	cookTomatoPasteVIP := enums.IngredientPreparations[cookPrep.ID][tomatoPaste.ID]
	cookPotVPV := enums.PreparationVessels[cookPrep.ID][largePot.ID]
	cookWoodenSpoonVPI := enums.PreparationInstruments[cookPrep.ID][woodenSpoon.ID]

	placeChickenVIP := enums.IngredientPreparations[placePrep.ID][chickenBreast.ID]
	placePotVPV := enums.PreparationVessels[placePrep.ID][largePot.ID]

	seasonChickenVIP := enums.IngredientPreparations[seasonPrep.ID][chickenBreast.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonBlackPepperVIP := enums.IngredientPreparations[seasonPrep.ID][blackPepper.ID]
	seasonPotVPV := enums.PreparationVessels[seasonPrep.ID][largePot.ID]

	addWaterVIP := enums.IngredientPreparations[addPrep.ID][water.ID]
	addChickenVIP := enums.IngredientPreparations[addPrep.ID][chickenBreast.ID]
	addVermicelliVIP := enums.IngredientPreparations[addPrep.ID][vermicelli.ID]
	simmerChickenVIP := enums.IngredientPreparations[simmerPrep.ID][chickenBreast.ID]
	simmerPotVPV := enums.PreparationVessels[simmerPrep.ID][largePot.ID]

	transferChickenVIP := enums.IngredientPreparations[transferPrep.ID][chickenBreast.ID]
	transferMediumBowlVPV := enums.PreparationVessels[transferPrep.ID][mediumBowl.ID]

	shredChickenVIP := enums.IngredientPreparations[shredPrep.ID][chickenBreast.ID]
	shredForkVPI := enums.PreparationInstruments[shredPrep.ID][fork.ID]
	shredMediumBowlVPV := enums.PreparationVessels[shredPrep.ID][mediumBowl.ID]

	stirParsleyVIP := enums.IngredientPreparations[stirPrep.ID][parsley.ID]
	stirLimeJuiceVIP := enums.IngredientPreparations[stirPrep.ID][limeJuice.ID]
	stirPotVPV := enums.PreparationVessels[stirPrep.ID][largePot.ID]
	stirWoodenSpoonVPI := enums.PreparationInstruments[stirPrep.ID][woodenSpoon.ID]

	// Measurement unit bridges
	butterTbspVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]
	onionUnitVIMU := enums.IngredientMeasurementUnits[onion.ID][unitMeasurement.ID]
	turmericTspVIMU := enums.IngredientMeasurementUnits[turmeric.ID][teaspoonMeasurement.ID]
	potatoUnitVIMU := enums.IngredientMeasurementUnits[potato.ID][unitMeasurement.ID]
	carrotUnitVIMU := enums.IngredientMeasurementUnits[carrot.ID][unitMeasurement.ID]
	tomatoPasteTbspVIMU := enums.IngredientMeasurementUnits[tomatoPaste.ID][tablespoonMeasurement.ID]
	chickenBreastUnitVIMU := enums.IngredientMeasurementUnits[chickenBreast.ID][unitMeasurement.ID]
	saltTspVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	blackPepperTspVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][teaspoonMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	vermicelliCupVIMU := enums.IngredientMeasurementUnits[vermicelli.ID][cupMeasurement.ID]
	parsleyTbspVIMU := enums.IngredientMeasurementUnits[parsley.ID][tablespoonMeasurement.ID]
	limeJuiceTbspVIMU := enums.IngredientMeasurementUnits[limeJuice.ID][tablespoonMeasurement.ID]

	// Ingredient states
	translucentState := enums.IngredientStates["translucent"]
	combinedState := enums.IngredientStates["combined"]
	tenderState := enums.IngredientStates["tender"]

	// Step 0: Melt butter (or olive oil) in large pot over medium heat
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        meltPrep.ID,
		Index:                0,
		ExplicitInstructions: "In a large pot, melt the butter (or olive oil) over medium heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &meltButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTbspVIMU.ID,
				Name:                             "butter or olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &meltPotVPV.ID,
				Name:                     "large pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with melted fat",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	// Step 1: Add onion and cook until softened and translucent, about 5 minutes
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        sautePrep.ID,
		Index:                1,
		ExplicitInstructions: "Add the onion and cook, stirring frequently, until softened and translucent, about 5 minutes; you don't want the onion to take on any color.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &sauteOnionVIP.ID,
				ValidIngredientMeasurementUnitID: &onionUnitVIMU.ID,
				Name:                             "small or ½ large yellow onion, finely chopped",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &sautePotVPV.ID,
				Name:                            "pot with melted fat",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cooked translucent onion",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot with cooked onion",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: translucentState.ID,
				Notes:             "Onion should be softened and translucent, no color",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
	}

	// Step 2: Sprinkle in turmeric and stir until fragrant, about 30 seconds
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        sprinklePrep.ID,
		Index:                2,
		ExplicitInstructions: "Sprinkle in the turmeric and stir until fragrant, about 30 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &sprinkleTurmericVIP.ID,
				ValidIngredientMeasurementUnitID: &turmericTspVIMU.ID,
				Name:                             "ground turmeric",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &sprinklePotVPV.ID,
				Name:                            "pot with cooked onion",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with onion and turmeric",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: combinedState.ID,
				Notes:             "Turmeric should be fragrant and combined with the onion",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
	}

	// Step 3: Add potato and carrot, stir and cook for 2 minutes
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                3,
		ExplicitInstructions: "Add the potato and carrot, then stir and cook for 2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &addPotatoVIP.ID,
				ValidIngredientMeasurementUnitID: &potatoUnitVIMU.ID,
				Name:                             "medium Yukon Gold potato, diced into small cubes",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addCarrotVIP.ID,
				ValidIngredientMeasurementUnitID: &carrotUnitVIMU.ID,
				Name:                             "large carrot, finely chopped",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &addPotVPV.ID,
				Name:                            "pot with onion and turmeric",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with vegetables",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	// Step 4: Adjust heat to medium-low, add tomato paste, cook until color releases (2-3 min)
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        cookPrep.ID,
		Index:                4,
		ExplicitInstructions: "Adjust heat to medium-low and add the tomato paste. Cook, stirring frequently, until the tomato paste releases its color into the oil, 2 to 3 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120),
			Max: pointer.To[uint32](180),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &cookTomatoPasteVIP.ID,
				ValidIngredientMeasurementUnitID: &tomatoPasteTbspVIMU.ID,
				Name:                             "tomato paste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
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
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &cookPotVPV.ID,
				Name:                            "pot with vegetables",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with cooked tomato paste and vegetables",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	// Step 5: Place chicken on top of vegetables
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        placePrep.ID,
		Index:                5,
		ExplicitInstructions: "Place the chicken breast on top of the vegetables.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &placeChickenVIP.ID,
				ValidIngredientMeasurementUnitID: &chickenBreastUnitVIMU.ID,
				Name:                             "boneless, skinless chicken breast",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &placePotVPV.ID,
				Name:                            "pot with cooked tomato paste and vegetables",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with chicken on vegetables",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	// Step 6: Season with salt and pepper
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        seasonPrep.ID,
		Index:                6,
		ExplicitInstructions: "Season with 2 teaspoons salt and ¼ teaspoon pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonChickenVIP.ID,
				Name:                            "chicken in pot",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTspVIMU.ID,
				Name:                             "Kosher salt (such as Diamond Crystal)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonBlackPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &blackPepperTspVIMU.ID,
				Name:                             "black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &seasonPotVPV.ID,
				Name:                            "pot with chicken on vegetables",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with seasoned chicken and vegetables",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	// Step 7: Add water and stir
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                7,
		ExplicitInstructions: "Add 6 cups of water and stir.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &addWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &addPotVPV.ID,
				Name:                            "pot with seasoned chicken and vegetables",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with chicken and seasoned soup base",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	// Step 8: Partially cover, bring to boil, then cover completely, adjust to low, simmer 40 min
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        simmerPrep.ID,
		Index:                8,
		ExplicitInstructions: "Partially cover, increase heat and bring to a boil, then cover completely, adjust heat to low and simmer gently for 40 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](2400),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &simmerChickenVIP.ID,
				Name:                            "chicken in soup",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &simmerPotVPV.ID,
				Name:                            "pot with chicken and seasoned soup base",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "simmered soup with cooked chicken",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: tenderState.ID,
				Notes:             "Chicken should be cooked through and tender",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
	}

	// Step 9: Transfer chicken to medium bowl and shred with two forks
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        transferPrep.ID,
		Index:                9,
		ExplicitInstructions: "Transfer the chicken to a medium bowl.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferChickenVIP.ID,
				Name:                            "cooked chicken breast",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &transferMediumBowlVPV.ID,
				Name:                     "medium bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "chicken in bowl",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Shred chicken with two forks
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        shredPrep.ID,
		Index:                10,
		ExplicitInstructions: "Shred the chicken with two forks while the soup continues to simmer.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &shredChickenVIP.ID,
				Name:                            "chicken in bowl",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &shredForkVPI.ID,
				Name:                         "fork",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &shredMediumBowlVPV.ID,
				Name:                     "medium bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "shredded chicken",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 11: Add shredded chicken and vermicelli, bring to lively simmer, cover, simmer 10 min
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                11,
		ExplicitInstructions: "Taste the soup and adjust salt as needed. Add the shredded chicken and vermicelli, stir and increase heat to bring the soup to a lively simmer. Cover, adjust heat to low and simmer until the noodles soften and the flavors come together, 10 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addChickenVIP.ID,
				Name:                            "shredded chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addVermicelliVIP.ID,
				ValidIngredientMeasurementUnitID: &vermicelliCupVIMU.ID,
				Name:                             "broken wheat vermicelli noodles, broken angel hair pasta or fideo",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &addPotVPV.ID,
				Name:                            "simmered soup (without chicken)",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with finished soup",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: tenderState.ID,
				Notes:             "Noodles should be softened and flavors combined",
				Ingredients:       []uint64{1},
				Optional:          false,
			},
		},
	}

	// Step 12: Remove from heat, stir in parsley and lime juice, let rest 5-10 min
	step12 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                12,
		ExplicitInstructions: "Remove from the heat, stir in the parsley and lime juice and let sit, covered, for 5 to 10 minutes. Taste and adjust seasoning with more lime juice, salt or pepper as needed.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300),
			Max: pointer.To[uint32](600),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &stirParsleyVIP.ID,
				ValidIngredientMeasurementUnitID: &parsleyTbspVIMU.ID,
				Name:                             "chopped parsley",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &stirLimeJuiceVIP.ID,
				ValidIngredientMeasurementUnitID: &limeJuiceTbspVIMU.ID,
				Name:                             "lime or lemon juice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &stirWoodenSpoonVPI.ID,
				Name:                         "wooden spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &stirPotVPV.ID,
				Name:                            "pot with finished soup",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "chicken and vermicelli soup",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
			Name:                "Chicken and Vermicelli Soup with Lime",
			Slug:                "chicken-vermicelli-soup",
			Source:              "https://cooking.nytimes.com/recipes/1026337-chicken-and-vermicelli-soup-with-lime",
			Description:         "A comforting Middle Eastern–inspired soup with tender chicken, vermicelli noodles, turmeric, and a bright finish of lime. Serve with lime slices and extra parsley.",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 6,
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12,
			},
			PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
		},
	}
}
