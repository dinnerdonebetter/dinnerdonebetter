package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ClassicSmashBurgersRecipe creates the Classic Smashed Burgers recipe from Serious Eats.
func ClassicSmashBurgersRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	heatPrep := enums.Preparations["heat"]
	formPrep := enums.Preparations["form"]
	toastPrep := enums.Preparations["toast"]
	smashPrep := enums.Preparations["smash"]
	panSearPrep := enums.Preparations["pan-sear"]
	flipPrep := enums.Preparations["flip"]
	topPrep := enums.Preparations["top"]
	assemblePrep := enums.Preparations["assemble"]

	// Get ingredients
	groundBeef := enums.Ingredients["ground beef"]
	vegetableOil := enums.Ingredients["vegetable oil"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	americanCheese := enums.Ingredients["American cheese"]
	burgerBun := enums.Ingredients["burger bun"]

	// Get measurement units
	ounceMeasurement := enums.MeasurementUnits["ounce"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	sliceMeasurement := enums.MeasurementUnits["slice"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	wideSpatula := enums.Instruments["wide spatula"]
	bareHands := enums.Instruments["bare hands"]

	// Get vessels
	castIronSkillet := enums.Vessels["cast iron skillet"]
	servingPlate := enums.Vessels["serving plate"]

	// Get ingredient states for completion conditions
	smokingState := enums.IngredientStates["smoking"]
	brownedState := enums.IngredientStates["browned"]

	// Get bridge table entries
	// Heat preparation bridges
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][vegetableOil.ID]
	oilTeaspoonVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][teaspoonMeasurement.ID]
	heatSkilletVPV := enums.PreparationVessels[heatPrep.ID][castIronSkillet.ID]

	// Form preparation bridges
	beefOunceVIMU := enums.IngredientMeasurementUnits[groundBeef.ID][ounceMeasurement.ID]
	formBeefVIP := enums.IngredientPreparations[formPrep.ID][groundBeef.ID]
	formBareHandsVPI := enums.PreparationInstruments[formPrep.ID][bareHands.ID]

	// Season preparation bridges (for measurement units)
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]

	// Toast preparation bridges
	toastBunVIP := enums.IngredientPreparations[toastPrep.ID][burgerBun.ID]
	bunUnitVIMU := enums.IngredientMeasurementUnits[burgerBun.ID][unitMeasurement.ID]
	toastSkilletVPV := enums.PreparationVessels[toastPrep.ID][castIronSkillet.ID]

	// Smash preparation bridges
	smashBeefVIP := enums.IngredientPreparations[smashPrep.ID][groundBeef.ID]
	smashSpatulaVPI := enums.PreparationInstruments[smashPrep.ID][wideSpatula.ID]
	smashSkilletVPV := enums.PreparationVessels[smashPrep.ID][castIronSkillet.ID]

	// Pan-sear preparation bridges
	panSearSpatulaVPI := enums.PreparationInstruments[panSearPrep.ID][wideSpatula.ID]
	panSearSkilletVPV := enums.PreparationVessels[panSearPrep.ID][castIronSkillet.ID]

	// Flip preparation bridges
	flipSpatulaVPI := enums.PreparationInstruments[flipPrep.ID][wideSpatula.ID]
	flipSkilletVPV := enums.PreparationVessels[flipPrep.ID][castIronSkillet.ID]

	// Top preparation bridges
	topCheeseVIP := enums.IngredientPreparations[topPrep.ID][americanCheese.ID]
	cheeseSliceVIMU := enums.IngredientMeasurementUnits[americanCheese.ID][sliceMeasurement.ID]
	topSkilletVPV := enums.PreparationVessels[topPrep.ID][castIronSkillet.ID]

	// Assemble preparation bridges
	assembleBunVIP := enums.IngredientPreparations[assemblePrep.ID][burgerBun.ID]
	assembleServingPlateVPV := enums.PreparationVessels[assemblePrep.ID][servingPlate.ID]

	// Step 0: Add oil to skillet and preheat
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        heatPrep.ID,
		Index:                0,
		ExplicitInstructions: "Add oil to a 12-inch cast iron skillet and wipe around with a paper towel. Set the skillet over medium heat and allow to preheat for about 5 minutes, then increase the heat to high until smoking.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilTeaspoonVIMU.ID,
				Name:                             "vegetable oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &heatSkilletVPV.ID,
				Name:                     "12-inch cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "hot skillet with smoking oil",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: smokingState.ID,
				Notes:             "Oil should be smoking",
				Ingredients:       []uint64{0}, // Index of the oil ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 1: Divide, form, and season beef into patties
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        formPrep.ID,
		Index:                1,
		ExplicitInstructions: "Divide the ground beef into four 4-ounce portions. Gently form each portion into a cylindrical puck about 2 inches tall, pressing together just until the meat holds its shape without falling apart. Season generously on all sides with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &formBeefVIP.ID,
				ValidIngredientMeasurementUnitID: &beefOunceVIMU.ID,
				Name:                             "ground beef",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 16,
					Max: pointer.To[float32](20),
				},
			},
			{
				ValidIngredientMeasurementUnitID: &saltGramVIMU.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste: true,
			},
			{
				ValidIngredientMeasurementUnitID: &pepperGramVIMU.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste: true,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &formBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "formed 4-ounce beef patties",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 2: Toast buns
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        toastPrep.ID,
		Index:                2,
		ExplicitInstructions: "Lightly toast the burger buns.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &toastBunVIP.ID,
				ValidIngredientMeasurementUnitID: &bunUnitVIMU.ID,
				Name:                             "burger buns",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &toastSkilletVPV.ID,
				Name:                            "hot skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "toasted burger buns",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 3: Add pucks to skillet and smash
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        smashPrep.ID,
		Index:                3,
		ExplicitInstructions: "Add 2 beef pucks to the hot skillet and, using a firm, stiff metal spatula, press down on each one until they're roughly 4 to 4 1/2 inches in diameter and 1/2-inch thick. It helps to use a second spatula to apply downward pressure to the first if you are having trouble smashing them hard enough.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &smashBeefVIP.ID,
				ValidIngredientMeasurementUnitID: &beefOunceVIMU.ID,
				Name:                             "formed 4-ounce beef patties",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &smashSpatulaVPI.ID,
				Name:                         "metal spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &smashSkilletVPV.ID,
				Name:                            "hot skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "smashed burger patties",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 4: Sear first side until golden brown
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        panSearPrep.ID,
		Index:                4,
		ExplicitInstructions: "Cook without moving until a golden brown crust develops, about 1 1/2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](90), // 1.5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "smashed burger patties",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &panSearSpatulaVPI.ID,
				Name:                         "metal spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &panSearSkilletVPV.ID,
				Name:                            "hot skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seared burger patties (first side)",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: brownedState.ID,
				Notes:             "A golden brown crust should develop on the bottom",
				Ingredients:       []uint64{0}, // Index of the patty ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 5: Flip patties
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        flipPrep.ID,
		Index:                5,
		ExplicitInstructions: "Use the edge of the spatula to carefully scrape up and flip the patties one at a time, making sure to get all browned bits removed from the skillet.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "seared burger patties (first side)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &flipSpatulaVPI.ID,
				Name:                         "metal spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &flipSkilletVPV.ID,
				Name:                            "hot skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "flipped burger patties",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 6: Top with cheese (optional)
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        topPrep.ID,
		Index:                6,
		Optional:             true,
		ExplicitInstructions: "If using cheese, add a slice to each patty now.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "flipped burger patties",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ValidIngredientPreparationID:     &topCheeseVIP.ID,
				ValidIngredientMeasurementUnitID: &cheeseSliceVIMU.ID,
				Name:                             "cheese slices",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &topSkilletVPV.ID,
				Name:                            "hot skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "burger patties with cheese",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 7: Finish cooking second side
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        panSearPrep.ID,
		Index:                7,
		ExplicitInstructions: "Continue to cook until the patties are cooked to desired doneness—about 30 seconds longer for medium-rare.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "burger patties with cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &panSearSkilletVPV.ID,
				Name:                            "hot skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cooked smash burger patties",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	// Step 8: Assemble burgers
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        assemblePrep.ID,
		Index:                8,
		ExplicitInstructions: "Transfer the patties to the toasted buns, topping buns and/or patties as desired, close the burgers, and serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &assembleBunVIP.ID,
				ValidIngredientMeasurementUnitID: &bunUnitVIMU.ID,
				Name:                             "toasted burger buns",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "cooked smash burger patties",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &assembleServingPlateVPV.ID,
				Name:                     "serving plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "classic smash burger",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
			Name:                "Classic Smashed Burgers",
			Slug:                "classic-smashed-burgers",
			Source:              "https://www.seriouseats.com/classic-smashed-burgers-recipe",
			Description:         "Classic smashed cheeseburgers with maximum juiciness and a deep-brown, beefy crust. Smashing down on the burger patties within the first 30 seconds of hitting a hot skillet ensures maximum juiciness and a flavorful, well-browned crust.",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 4,
			},
			PortionName:       "burger",
			PluralPortionName: "burgers",
			EligibleForMeals:  true,
			Steps:             []*mealplanning.RecipeStepCreationRequestInput{step0, step1, step2, step3, step4, step5, step6, step7, step8},
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
