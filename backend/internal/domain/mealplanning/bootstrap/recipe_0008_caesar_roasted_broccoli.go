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
	grindPrep := enums.Preparations["grind"]
	linePrep := enums.Preparations["line"]
	adjustPrep := enums.Preparations["adjust"]
	preheatPrep := enums.Preparations["preheat"]
	tossPrep := enums.Preparations["toss"]
	roastPrep := enums.Preparations["roast"]
	topPrep := enums.Preparations["top"]
	transferPrep := enums.Preparations["transfer"]
	zestPrep := enums.Preparations["zest"]

	// Get ingredients
	broccoli := enums.Ingredients["broccoli"]
	oliveOil := enums.Ingredients["olive oil"]
	wholePeppercorns := enums.Ingredients["whole black peppercorns"]
	parmesan := enums.Ingredients["parmesan cheese"]
	lemon := enums.Ingredients["lemon"]
	salt := enums.Ingredients["salt"]
	breadcrumbs := enums.Ingredients["breadcrumbs"]
	aluminumFoilIngredient := enums.Ingredients["aluminum foil"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	poundMeasurement := enums.MeasurementUnits["pound"]
	gramMeasurement := enums.MeasurementUnits["gram"]

	// Get instruments
	mortarAndPestle := enums.Instruments["mortar and pestle"]
	spiceGrinder := enums.Instruments["spice grinder"]
	microplane := enums.Instruments["microplane"]

	// Get vessels
	bakingSheet := enums.Vessels["baking sheet"]
	largeBowl := enums.Vessels["large bowl"]
	servingPlatter := enums.Vessels["serving platter"]
	oven := enums.Vessels["oven"]

	// === BROCCOLI BRIDGE TABLE ENTRIES ===
	transferBroccoliVIP := enums.IngredientPreparations[transferPrep.ID][broccoli.ID]
	transferBakingSheetVPV := enums.PreparationVessels[transferPrep.ID][bakingSheet.ID]
	transferServingPlatterVPV := enums.PreparationVessels[transferPrep.ID][servingPlatter.ID]

	lineAluminumFoilVIP := enums.IngredientPreparations[linePrep.ID][aluminumFoilIngredient.ID]
	lineAluminumFoilVIMU := enums.IngredientMeasurementUnits[aluminumFoilIngredient.ID][unitMeasurement.ID]
	lineBakingSheetVPV := enums.PreparationVessels[linePrep.ID][bakingSheet.ID]

	adjustOvenVPV := enums.PreparationVessels[adjustPrep.ID][oven.ID]
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]
	preheatBakingSheetVPV := enums.PreparationVessels[preheatPrep.ID][bakingSheet.ID]

	tossBroccoliVIP := enums.IngredientPreparations[tossPrep.ID][broccoli.ID]
	tossOliveOilVIP := enums.IngredientPreparations[tossPrep.ID][oliveOil.ID]
	tossSaltVIP := enums.IngredientPreparations[tossPrep.ID][salt.ID]
	tossLargeBowlVPV := enums.PreparationVessels[tossPrep.ID][largeBowl.ID]

	zestLemonVIP := enums.IngredientPreparations[zestPrep.ID][lemon.ID]
	zestMicroplaneVPI := enums.PreparationInstruments[zestPrep.ID][microplane.ID]

	roastBroccoliVIP := enums.IngredientPreparations[roastPrep.ID][broccoli.ID]
	roastBakingSheetVPV := enums.PreparationVessels[roastPrep.ID][bakingSheet.ID]

	topBroccoliVIP := enums.IngredientPreparations[topPrep.ID][broccoli.ID]
	topBreadcrumbsVIP := enums.IngredientPreparations[topPrep.ID][breadcrumbs.ID]
	topParmesanVIP := enums.IngredientPreparations[topPrep.ID][parmesan.ID]
	topServingPlatterVPV := enums.PreparationVessels[topPrep.ID][servingPlatter.ID]

	// Grind preparation bridges
	grindPeppercornsVIP := enums.IngredientPreparations[grindPrep.ID][wholePeppercorns.ID]
	peppercornsGramVIMU := enums.IngredientMeasurementUnits[wholePeppercorns.ID][gramMeasurement.ID]
	grindMortarAndPestleVPI := enums.PreparationInstruments[grindPrep.ID][mortarAndPestle.ID]
	grindSpiceGrinderVPI := enums.PreparationInstruments[grindPrep.ID][spiceGrinder.ID]

	// Measurement unit bridges
	broccoliPoundVIMU := enums.IngredientMeasurementUnits[broccoli.ID][poundMeasurement.ID]
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	parmesanTablespoonVIMU := enums.IngredientMeasurementUnits[parmesan.ID][tablespoonMeasurement.ID]
	lemonUnitVIMU := enums.IngredientMeasurementUnits[lemon.ID][unitMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	breadcrumbsCupVIMU := enums.IngredientMeasurementUnits[breadcrumbs.ID][cupMeasurement.ID]

	// ==================== CAESAR ROASTED BROCCOLI RECIPE STEPS ====================

	// Step 0: Grind whole black peppercorns
	brStep0 := &mealplanning.RecipeStepCreationRequestInput{
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

	// Step 1: Line baking sheet with aluminum foil
	brStep1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        linePrep.ID,
		Index:                1,
		ExplicitInstructions: "Line a rimmed baking sheet with aluminum foil.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &lineAluminumFoilVIP.ID,
				ValidIngredientMeasurementUnitID: &lineAluminumFoilVIMU.ID,
				Name:                             "aluminum foil",
				Quantity: types.Float32RangeWithOptionalMax{
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

	// Broccoli Step 2: Adjust oven rack to upper position
	brStep2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        adjustPrep.ID,
		Index:                2,
		ExplicitInstructions: "Adjust the oven rack to the upper position.",
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &adjustOvenVPV.ID,
				Name:                     "oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "oven with rack adjusted",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 3: Preheat oven to 500°F
	brStep3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        preheatPrep.ID,
		Index:                3,
		ExplicitInstructions: "Preheat the oven to 500°F (260°C).",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260),
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &preheatOvenVPV.ID,
				Name:                            "oven with rack adjusted",
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

	// Broccoli Step 4: Place baking sheet in oven to preheat
	brStep4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        preheatPrep.ID,
		Index:                4,
		ExplicitInstructions: "Place the foil-lined baking sheet on the oven rack to preheat.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260),
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &preheatBakingSheetVPV.ID,
				Name:                            "foil-lined baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "preheated oven",
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

	// Broccoli Step 5: Toss broccoli with olive oil, salt, and pepper
	brStep5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        tossPrep.ID,
		Index:                5,
		ExplicitInstructions: "In a large bowl, toss the broccoli florets with olive oil, salt, and pepper.",
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
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "freshly ground black pepper",
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

	// Broccoli Step 6: Add broccoli to preheated baking sheet
	brStep6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        transferPrep.ID,
		Index:                6,
		ExplicitInstructions: "Carefully add the broccoli to the preheated baking sheet in a single layer.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
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
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
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

	// Broccoli Step 7: Roast broccoli
	brStep7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        roastPrep.ID,
		Index:                7,
		ExplicitInstructions: "Roast until the broccoli is tender and deeply browned in spots, about 20 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](1200), // 20 minutes
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](6),
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

	// Broccoli Step 8: Zest the lemon
	brStep8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        zestPrep.ID,
		Index:                8,
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
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 9: Toss roasted broccoli with lemon zest
	brStep9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        tossPrep.ID,
		Index:                9,
		ExplicitInstructions: "In the now empty bowl, toss the broccoli with 1 teaspoon lemon zest.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "roasted broccoli",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "lemon zest",
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

	// Broccoli Step 10: Transfer broccoli to serving platter
	brStep10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        transferPrep.ID,
		Index:                10,
		ExplicitInstructions: "Transfer the broccoli to a serving platter.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
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

	// Broccoli Step 11: Sprinkle with breadcrumbs and Parmigiano-Reggiano
	// This step references the Caesar Breadcrumbs recipe as a component
	brStep11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        topPrep.ID,
		Index:                11,
		ExplicitInstructions: "Sprinkle with breadcrumbs and Parmigiano-Reggiano.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &topBroccoliVIP.ID,
				ValidIngredientMeasurementUnitID: &broccoliPoundVIMU.ID,
				Name:                             "broccoli on serving platter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:         pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				RecipeStepProductRecipeID:        getRecipeIDBySlug(createdRecipes, "caesar-breadcrumbs"),
				RecipeStepProductRecipeSlug:      pointer.To("caesar-breadcrumbs"),
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
		Steps:             []*mealplanning.RecipeStepCreationRequestInput{brStep0, brStep1, brStep2, brStep3, brStep4, brStep5, brStep6, brStep7, brStep8, brStep9, brStep10, brStep11},
		PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
		Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
		AlsoCreateMeal:    false,
	}

	return []*mealplanning.RecipeCreationRequestInput{
		caesarRoastedBroccoliRecipe,
	}
}
