package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// CaesarBreadcrumbsRecipe creates the Caesar Breadcrumbs recipe.
// Source: https://www.seriouseats.com/caesar-roasted-broccoli-recipe-8672043
func CaesarBreadcrumbsRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	meltPrep := enums.Preparations["melt"]
	stirPrep := enums.Preparations["stir"]
	mixPrep := enums.Preparations["mix"]
	coatPrep := enums.Preparations["coat"]
	seasonPrep := enums.Preparations["season"]
	transferPrep := enums.Preparations["transfer"]
	zestPrep := enums.Preparations["zest"]

	// Get ingredients for breadcrumbs
	saltedButter := enums.Ingredients["salted butter"]
	breadcrumbs := enums.Ingredients["breadcrumbs"]
	anchovyPaste := enums.Ingredients["anchovy paste"]
	garlic := enums.Ingredients["garlic"]
	salt := enums.Ingredients["salt"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	rubberSpatula := enums.Instruments["rubber spatula"]
	microplane := enums.Instruments["microplane"]

	// Get vessels
	smallNonstickSkillet := enums.Vessels["small nonstick skillet"]
	smallBowl := enums.Vessels["small bowl"]

	// Get bridge table entries for breadcrumbs
	meltButterVIP := enums.IngredientPreparations[meltPrep.ID][saltedButter.ID]
	meltSkilletVPV := enums.PreparationVessels[meltPrep.ID][smallNonstickSkillet.ID]

	stirAnchovyVIP := enums.IngredientPreparations[stirPrep.ID][anchovyPaste.ID]
	stirGarlicVIP := enums.IngredientPreparations[stirPrep.ID][garlic.ID]
	stirBreadcrumbsVIP := enums.IngredientPreparations[stirPrep.ID][breadcrumbs.ID]
	stirSpatulaVPI := enums.PreparationInstruments[stirPrep.ID][rubberSpatula.ID]

	zestLemonVIP := enums.IngredientPreparations[zestPrep.ID][lemon.ID]
	zestMicroplaneVPI := enums.PreparationInstruments[zestPrep.ID][microplane.ID]

	mixButterVIP := enums.IngredientPreparations[mixPrep.ID][saltedButter.ID]

	coatBreadcrumbsVIP := enums.IngredientPreparations[coatPrep.ID][breadcrumbs.ID]
	coatSpatulaVPI := enums.PreparationInstruments[coatPrep.ID][rubberSpatula.ID]

	seasonBreadcrumbsVIP := enums.IngredientPreparations[seasonPrep.ID][breadcrumbs.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]

	transferBreadcrumbsVIP := enums.IngredientPreparations[transferPrep.ID][breadcrumbs.ID]
	transferSmallBowlVPV := enums.PreparationVessels[transferPrep.ID][smallBowl.ID]

	// Measurement unit bridges for breadcrumbs
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[saltedButter.ID][tablespoonMeasurement.ID]
	breadcrumbsCupVIMU := enums.IngredientMeasurementUnits[breadcrumbs.ID][cupMeasurement.ID]
	anchovyTeaspoonVIMU := enums.IngredientMeasurementUnits[anchovyPaste.ID][teaspoonMeasurement.ID]
	garlicUnitVIMU := enums.IngredientMeasurementUnits[garlic.ID][unitMeasurement.ID]
	lemonUnitVIMU := enums.IngredientMeasurementUnits[lemon.ID][unitMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]

	// Breadcrumbs Step 0: Melt butter in a small nonstick skillet
	bcStep0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        meltPrep.ID,
		Index:                0,
		ExplicitInstructions: "In a small nonstick skillet, melt the butter over medium-low heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &meltButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				Name:                             "salted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &meltSkilletVPV.ID,
				Name:                     "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "melted butter",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &tablespoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "small nonstick skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Breadcrumbs Step 1: Stir in anchovy paste and garlic
	bcStep1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                1,
		ExplicitInstructions: "Stir in the anchovy paste and garlic.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "melted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &stirAnchovyVIP.ID,
				ValidIngredientMeasurementUnitID: &anchovyTeaspoonVIMU.ID,
				Name:                             "anchovy paste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &stirGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicUnitVIMU.ID,
				Name:                             "garlic, minced",
				QuantityNotes:                    "1 small clove",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &stirSpatulaVPI.ID,
				Name:                         "flexible spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "butter with anchovy and garlic",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "small nonstick skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Breadcrumbs Step 2: Mix until fragrant
	bcStep2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        mixPrep.ID,
		Index:                2,
		ExplicitInstructions: "Cook until fragrant, about 1 minute.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &mixButterVIP.ID,
				Name:                            "butter with anchovy and garlic",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "fragrant butter mixture",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "small nonstick skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Breadcrumbs Step 3: Add breadcrumbs and toss to coat
	bcStep3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                3,
		ExplicitInstructions: "Add the breadcrumbs and, using a flexible spatula, toss to coat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "fragrant butter mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &stirBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				Name:                             "plain breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &stirSpatulaVPI.ID,
				Name:                         "flexible spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "coated breadcrumbs",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				Name:  "small nonstick skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Breadcrumbs Step 4: Coat breadcrumbs until golden brown
	bcStep4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        coatPrep.ID,
		Index:                4,
		ExplicitInstructions: "Cook, stirring constantly until the breadcrumbs are golden brown, about 3 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &coatBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				Name:                             "coated breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &coatSpatulaVPI.ID,
				Name:                         "flexible spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "toasted breadcrumbs",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				Name:  "small nonstick skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Breadcrumbs Step 5: Zest the lemon
	bcStep5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        zestPrep.ID,
		Index:                5,
		ExplicitInstructions: "Using a microplane, zest the lemon.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &zestLemonVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonUnitVIMU.ID,
				Name:                             "lemon",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &zestMicroplaneVPI.ID,
				Name:                         "microplane",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "lemon zest",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &teaspoonMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
		},
	}

	// Breadcrumbs Step 6: Stir in lemon zest (off heat)
	bcStep6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                6,
		ExplicitInstructions: "Off heat, stir in 1/2 teaspoon lemon zest.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "toasted breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "lemon zest",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &stirSpatulaVPI.ID,
				Name:                         "flexible spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "breadcrumbs with lemon zest",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
			{
				Name:  "small nonstick skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Breadcrumbs Step 7: Season breadcrumbs with salt
	bcStep7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        seasonPrep.ID,
		Index:                7,
		ExplicitInstructions: "Season with salt to taste.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &seasonBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				Name:                             "breadcrumbs with lemon zest",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
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
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seasoned caesar breadcrumbs",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	// Breadcrumbs Step 8: Transfer breadcrumbs to bowl and let cool
	bcStep8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        transferPrep.ID,
		Index:                8,
		ExplicitInstructions: "Transfer to a bowl and let cool completely.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &transferBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				Name:                             "seasoned caesar breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &transferSmallBowlVPV.ID,
				Name:                     "small bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "caesar breadcrumbs",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	caesarBreadcrumbsRecipe := &mealplanning.RecipeCreationRequestInput{
		Name:                "Caesar Breadcrumbs",
		Slug:                "caesar-breadcrumbs",
		Source:              "https://www.seriouseats.com/caesar-roasted-broccoli-recipe-8672043",
		Description:         "Savory, crisp Caesar-flavored breadcrumbs with anchovy, garlic, and lemon zest.",
		YieldsComponentType: mealplanning.MealComponentTypesAmuseBouche, // Component type for recipe components
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 0.25,
		},
		PortionName:       "cup",
		PluralPortionName: "cups",
		EligibleForMeals:  false, // This is a component, not a standalone meal
		Steps:             []*mealplanning.RecipeStepCreationRequestInput{bcStep0, bcStep1, bcStep2, bcStep3, bcStep4, bcStep5, bcStep6, bcStep7, bcStep8},
		PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
		Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
		AlsoCreateMeal:    false,
	}

	return []*mealplanning.RecipeCreationRequestInput{
		caesarBreadcrumbsRecipe,
	}
}
