package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// PeanutButterNoodlesRecipe creates the Peanut Butter Noodles recipe.
// Source: https://cooking.nytimes.com/recipes/1025047-peanut-butter-noodles
func PeanutButterNoodlesRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	submergePrep := enums.Preparations["submerge"]
	boilPrep := enums.Preparations["boil"]
	reservePrep := enums.Preparations["reserve"]
	drainPrep := enums.Preparations["drain"]
	addPrep := enums.Preparations["add"]
	removeFromHeatPrep := enums.Preparations["remove from heat"]
	stirPrep := enums.Preparations["stir"]
	seasonPrep := enums.Preparations["season"]
	topPrep := enums.Preparations["top"]

	// Get ingredients
	spaghetti := enums.Ingredients["spaghetti"]
	instantRamen := enums.Ingredients["instant ramen noodles"]
	salt := enums.Ingredients["salt"]
	water := enums.Ingredients["water"]
	peanutButter := enums.Ingredients["peanut butter"]
	butter := enums.Ingredients["butter"]
	parmesan := enums.Ingredients["parmesan cheese"]
	soySauce := enums.Ingredients["soy sauce"]

	// Get measurement units
	ounceMeasurement := enums.MeasurementUnits["ounce"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	woodenSpoon := enums.Instruments["wooden spoon"]
	spoon := enums.Instruments["spoon"]

	// Get vessels
	pot := enums.Vessels["pot"]
	colander := enums.Vessels["colander"]

	// Get bridge table entries
	submergeSpaghettiVIP := enums.IngredientPreparations[submergePrep.ID][spaghetti.ID]
	submergeInstantRamenVIP := enums.IngredientPreparations[submergePrep.ID][instantRamen.ID]
	submergeWaterVIP := enums.IngredientPreparations[submergePrep.ID][water.ID]
	submergeSaltVIP := enums.IngredientPreparations[submergePrep.ID][salt.ID]
	submergePotVPV := enums.PreparationVessels[submergePrep.ID][pot.ID]

	boilWoodenSpoonVPI := enums.PreparationInstruments[boilPrep.ID][woodenSpoon.ID]

	drainColanderVPV := enums.PreparationVessels[drainPrep.ID][colander.ID]

	addPeanutButterVIP := enums.IngredientPreparations[addPrep.ID][peanutButter.ID]
	addButterVIP := enums.IngredientPreparations[addPrep.ID][butter.ID]
	addParmesanVIP := enums.IngredientPreparations[addPrep.ID][parmesan.ID]
	addSoySauceVIP := enums.IngredientPreparations[addPrep.ID][soySauce.ID]

	stirSpoonVPI := enums.PreparationInstruments[stirPrep.ID][spoon.ID]

	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]

	topParmesanVIP := enums.IngredientPreparations[topPrep.ID][parmesan.ID]

	// Measurement unit bridges
	spaghettiOunceVIMU := enums.IngredientMeasurementUnits[spaghetti.ID][ounceMeasurement.ID]
	instantRamenUnitVIMU := enums.IngredientMeasurementUnits[instantRamen.ID][unitMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	peanutButterTablespoonVIMU := enums.IngredientMeasurementUnits[peanutButter.ID][tablespoonMeasurement.ID]
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]
	parmesanTablespoonVIMU := enums.IngredientMeasurementUnits[parmesan.ID][tablespoonMeasurement.ID]
	soySauceTeaspoonVIMU := enums.IngredientMeasurementUnits[soySauce.ID][teaspoonMeasurement.ID]

	// Ingredient states
	tenderState := enums.IngredientStates["tender"]
	coatedState := enums.IngredientStates["coated"]

	// Step 0: Place pasta in pot with salted water
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        submergePrep.ID,
		Index:                0,
		ExplicitInstructions: "Bring a pot of water to a boil. If using spaghetti, salt the water. Add the spaghetti (or ramen noodles) and cook according to package instructions.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &submergeSpaghettiVIP.ID,
				ValidIngredientMeasurementUnitID: &spaghettiOunceVIMU.ID,
				Name:                             "spaghetti",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Index:       pointer.To[uint16](0),
				OptionIndex: 0,
			},
			{
				ValidIngredientPreparationID:     &submergeInstantRamenVIP.ID,
				ValidIngredientMeasurementUnitID: &instantRamenUnitVIMU.ID,
				Name:                             "package instant ramen (noodles only)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
				Index:       pointer.To[uint16](0),
				OptionIndex: 1,
			},
			{
				ValidIngredientPreparationID:     &submergeWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				QuantityNotes:                    "enough to cover noodles",
				Name:                             "water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{
				ValidIngredientPreparationID: &submergeSaltVIP.ID,
				QuantityNotes:                "if using spaghetti",
				Name:                         "salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 20,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &submergePotVPV.ID,
				Name:                     "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "noodles in water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 1: Boil until al dente
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        boilPrep.ID,
		Index:                1,
		ExplicitInstructions: "Cook the noodles according to package instructions until al dente.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "noodles in water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
				OptionIndex: 0,
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
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: tenderState.ID,
				Notes:             "noodles should be al dente",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "cooked al dente noodles",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 2: Reserve ½ cup cooking water
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        reservePrep.ID,
		Index:                2,
		ExplicitInstructions: "Reserve ½ cup of the cooking water.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "noodle cooking water",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot with noodles",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "reserved cooking water",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
		},
	}

	// Step 3: Drain noodles
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        drainPrep.ID,
		Index:                3,
		ExplicitInstructions: "Drain the noodles and return to the pot. Turn off the heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "cooked al dente noodles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
				OptionIndex: 0,
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "cooked al dente noodles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
				OptionIndex: 1,
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
				Name:              "drained noodles",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Return noodles to pot
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                4,
		ExplicitInstructions: "Return the drained noodles to the pot.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "drained noodles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "noodles in pot",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 5: Remove from heat
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        removeFromHeatPrep.ID,
		Index:                5,
		ExplicitInstructions: "Turn off the heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "noodles in pot",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "noodles in pot off heat",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 6: Add peanut butter, butter, Parmesan, and soy sauce
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                6,
		ExplicitInstructions: "Add the peanut butter, butter, Parmesan, and soy sauce.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "noodles in pot off heat",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addPeanutButterVIP.ID,
				ValidIngredientMeasurementUnitID: &peanutButterTablespoonVIMU.ID,
				Name:                             "creamy peanut butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &addButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				Name:                             "unsalted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanTablespoonVIMU.ID,
				Name:                             "finely grated Parmesan",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &addSoySauceVIP.ID,
				ValidIngredientMeasurementUnitID: &soySauceTeaspoonVIMU.ID,
				Name:                             "soy sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "noodles with sauce ingredients",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 7: Stir vigorously, adding reserved water until sauce is glossy and clings
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                7,
		ExplicitInstructions: "Vigorously stir the noodles for a minute, adding some reserved cooking water, a tablespoon or two at a time, until the sauce is glossy and clings to the noodles.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "noodles with sauce ingredients",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
				OptionIndex: 0,
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "noodles with sauce ingredients",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
				OptionIndex: 1,
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "reserved cooking water",
				QuantityNotes:                   "a tablespoon or two at a time, until sauce is glossy",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &stirSpoonVPI.ID,
				Name:                         "spoon",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: coatedState.ID,
				Notes:             "sauce should be glossy and cling to the noodles",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "peanut butter noodles",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 8: Season to taste with salt
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        seasonPrep.ID,
		Index:                8,
		ExplicitInstructions: "Season to taste with salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "peanut butter noodles",
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
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seasoned peanut butter noodles",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 9: Top with more cheese and serve
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        topPrep.ID,
		Index:                9,
		ExplicitInstructions: "Top with more Parmesan, if you'd like, and serve immediately.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "seasoned peanut butter noodles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID: &topParmesanVIP.ID,
				QuantityNotes:                "for serving, if desired",
				Name:                         "Parmesan",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "peanut butter noodles",
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
			Name:                "Peanut Butter Noodles",
			Slug:                "peanut-butter-noodles",
			Source:              "https://cooking.nytimes.com/recipes/1025047-peanut-butter-noodles",
			Description:         "A quick, creamy pasta dish with peanut butter, butter, Parmesan, and soy sauce. Use spaghetti or instant ramen (without the seasoning packet).",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 1,
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				step0, step1, step2, step3, step4, step5, step6, step7, step8, step9,
			},
			PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
		},
	}
}
