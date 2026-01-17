package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// CaesarRoastedBroccoliRecipe creates the Caesar Roasted Broccoli recipe.
// Source: https://www.seriouseats.com/caesar-roasted-broccoli-recipe-8672043
// Note: This recipe references the Caesar Breadcrumbs recipe, which must be created first.
// The createdRecipes map should contain the "caesar-breadcrumbs" recipe keyed by its slug.
func CaesarRoastedBroccoliRecipe(enums *Enumerations, createdRecipes map[string]*mealplanning.Recipe) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	linePrep := enums.Preparations["line"]
	preheatPrep := enums.Preparations["preheat"]
	tossPrep := enums.Preparations["toss"]
	roastPrep := enums.Preparations["roast"]
	topPrep := enums.Preparations["top"]
	transferPrep := enums.Preparations["transfer"]

	// Get ingredients
	broccoli := enums.Ingredients["broccoli"]
	oliveOil := enums.Ingredients["olive oil"]
	blackPepper := enums.Ingredients["black pepper"]
	parmesan := enums.Ingredients["parmesan cheese"]
	lemon := enums.Ingredients["lemon"]
	salt := enums.Ingredients["salt"]
	breadcrumbs := enums.Ingredients["breadcrumbs"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	poundMeasurement := enums.MeasurementUnits["pound"]
	gramMeasurement := enums.MeasurementUnits["gram"]

	// Get instruments
	aluminumFoil := enums.Instruments["aluminum foil"]

	// Get vessels
	bakingSheet := enums.Vessels["baking sheet"]
	largeBowl := enums.Vessels["large bowl"]
	servingPlatter := enums.Vessels["serving platter"]
	oven := enums.Vessels["oven"]

	// === BROCCOLI BRIDGE TABLE ENTRIES ===
	transferBroccoliVIP := enums.IngredientPreparations[transferPrep.ID][broccoli.ID]
	transferBakingSheetVPV := enums.PreparationVessels[transferPrep.ID][bakingSheet.ID]
	transferServingPlatterVPV := enums.PreparationVessels[transferPrep.ID][servingPlatter.ID]

	lineFoilVPI := enums.PreparationInstruments[linePrep.ID][aluminumFoil.ID]
	lineBakingSheetVPV := enums.PreparationVessels[linePrep.ID][bakingSheet.ID]

	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]
	preheatBakingSheetVPV := enums.PreparationVessels[preheatPrep.ID][bakingSheet.ID]

	tossBroccoliVIP := enums.IngredientPreparations[tossPrep.ID][broccoli.ID]
	tossOliveOilVIP := enums.IngredientPreparations[tossPrep.ID][oliveOil.ID]
	tossSaltVIP := enums.IngredientPreparations[tossPrep.ID][salt.ID]
	tossPepperVIP := enums.IngredientPreparations[tossPrep.ID][blackPepper.ID]
	tossLemonVIP := enums.IngredientPreparations[tossPrep.ID][lemon.ID]
	tossLargeBowlVPV := enums.PreparationVessels[tossPrep.ID][largeBowl.ID]

	roastBroccoliVIP := enums.IngredientPreparations[roastPrep.ID][broccoli.ID]
	roastBakingSheetVPV := enums.PreparationVessels[roastPrep.ID][bakingSheet.ID]

	topBroccoliVIP := enums.IngredientPreparations[topPrep.ID][broccoli.ID]
	topBreadcrumbsVIP := enums.IngredientPreparations[topPrep.ID][breadcrumbs.ID]
	topParmesanVIP := enums.IngredientPreparations[topPrep.ID][parmesan.ID]
	topServingPlatterVPV := enums.PreparationVessels[topPrep.ID][servingPlatter.ID]

	// Measurement unit bridges
	broccoliPoundVIMU := enums.IngredientMeasurementUnits[broccoli.ID][poundMeasurement.ID]
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]
	parmesanTablespoonVIMU := enums.IngredientMeasurementUnits[parmesan.ID][tablespoonMeasurement.ID]
	lemonTeaspoonVIMU := enums.IngredientMeasurementUnits[lemon.ID][teaspoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	breadcrumbsCupVIMU := enums.IngredientMeasurementUnits[breadcrumbs.ID][cupMeasurement.ID]

	// ==================== CAESAR ROASTED BROCCOLI RECIPE STEPS ====================

	// Step 0: Line baking sheet with aluminum foil
	brStep0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: linePrep.ID,
		Index:         0,
		Notes:         "Line a rimmed baking sheet with aluminum foil.",
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &lineFoilVPI.ID,
				Name:                         "aluminum foil",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &lineBakingSheetVPV.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "foil-lined baking sheet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 1: Preheat oven to 500°F
	brStep1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: preheatPrep.ID,
		Index:         1,
		Notes:         "Adjust oven rack to upper position and preheat oven to 500°F (260°C).",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260),
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &preheatOvenVPV.ID,
				Name:                     "oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "preheated oven",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 2: Place baking sheet in oven to preheat
	brStep2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: preheatPrep.ID,
		Index:         2,
		Notes:         "Place the foil-lined baking sheet on oven rack to preheat.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260),
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &preheatBakingSheetVPV.ID,
				Name:                            "foil-lined baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "preheated baking sheet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 3: Toss broccoli with olive oil, salt, and pepper
	brStep3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: tossPrep.ID,
		Index:         3,
		Notes:         "In a large bowl, toss broccoli florets with olive oil, salt, and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &tossBroccoliVIP.ID,
				ValidIngredientMeasurementUnitID: &broccoliPoundVIMU.ID,
				Name:                             "broccoli florets",
				QuantityNotes:                    "cut into 1 1/2 to 2-inch pieces",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &tossOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				Name:                             "extra-virgin olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ValidIngredientPreparationID:     &tossSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
			{
				ValidIngredientPreparationID:     &tossPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperGramVIMU.ID,
				Name:                             "freshly ground black pepper",
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
				Name:              "seasoned broccoli",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 4: Add broccoli to preheated baking sheet
	brStep4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: transferPrep.ID,
		Index:         4,
		Notes:         "Carefully add broccoli to preheated baking sheet in a single layer.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &transferBroccoliVIP.ID,
				ValidIngredientMeasurementUnitID: &broccoliPoundVIMU.ID,
				Name:                             "seasoned broccoli",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &transferBakingSheetVPV.ID,
				Name:                            "preheated baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "broccoli on baking sheet",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 5: Roast broccoli
	brStep5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: roastPrep.ID,
		Index:         5,
		Notes:         "Roast until broccoli is tender and deeply browned in spots, about 20 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](1200), // 20 minutes
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &roastBroccoliVIP.ID,
				ValidIngredientMeasurementUnitID: &broccoliPoundVIMU.ID,
				Name:                             "broccoli on baking sheet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &roastBakingSheetVPV.ID,
				Name:                     "baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "roasted broccoli",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 6: Toss roasted broccoli with lemon zest
	brStep6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: tossPrep.ID,
		Index:         6,
		Notes:         "In the now empty bowl, toss broccoli with 1 teaspoon lemon zest.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "roasted broccoli",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &tossLemonVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonTeaspoonVIMU.ID,
				Name:                             "lemon zest",
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
				Name:              "broccoli with lemon zest",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 7: Transfer broccoli to serving platter
	brStep7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: transferPrep.ID,
		Index:         7,
		Notes:         "Transfer broccoli to a serving platter.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "broccoli with lemon zest",
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
				Name:              "broccoli on serving platter",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 8: Sprinkle with breadcrumbs and Parmigiano-Reggiano
	// This step references the Caesar Breadcrumbs recipe as a component
	// Note: RecipeStepProductRecipeID will need to be set when creating the recipe,
	// as it references the breadcrumbs recipe that will be created first
	brStep8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: topPrep.ID,
		Index:         8,
		Notes:         "Sprinkle with breadcrumbs and Parmigiano-Reggiano and serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &topBroccoliVIP.ID,
				ValidIngredientMeasurementUnitID: &broccoliPoundVIMU.ID,
				Name:                             "broccoli on serving platter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				// RecipeStepProductRecipeID references the "Caesar Breadcrumbs" recipe (slug: "caesar-breadcrumbs")
				// The product "caesar breadcrumbs" is from step 7 (index 7), product index 0
				ProductOfRecipeStepIndex:         pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				RecipeStepProductRecipeID:        getRecipeIDBySlug(createdRecipes, "caesar-breadcrumbs"),
				ValidIngredientPreparationID:     &topBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				Name:                             "caesar breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ValidIngredientPreparationID:     &topParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanTablespoonVIMU.ID,
				Name:                             "grated Parmigiano-Reggiano cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &topServingPlatterVPV.ID,
				Name:                     "serving platter",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "caesar roasted broccoli",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	caesarRoastedBroccoliRecipe := &mealplanning.RecipeCreationRequestInput{
		Name:                "Caesar Roasted Broccoli",
		Slug:                "caesar-roasted-broccoli",
		Source:              "https://www.seriouseats.com/caesar-roasted-broccoli-recipe-8672043",
		Description:         "Dress up sweet and nutty roasted broccoli with savory, crisp Caesar-flavored breadcrumbs.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps:             []*mealplanning.RecipeStepCreationRequestInput{brStep0, brStep1, brStep2, brStep3, brStep4, brStep5, brStep6, brStep7, brStep8},
		PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
		Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
		AlsoCreateMeal:    false,
	}

	return []*mealplanning.RecipeCreationRequestInput{
		caesarRoastedBroccoliRecipe,
	}
}
