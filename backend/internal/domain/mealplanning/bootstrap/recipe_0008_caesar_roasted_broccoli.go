package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// CaesarRoastedBroccoliRecipe creates the Caesar Roasted Broccoli recipe.
// Source: https://www.seriouseats.com/caesar-roasted-broccoli-recipe-8672043
// This returns two recipes: Caesar Breadcrumbs (component) and Caesar Roasted Broccoli (main recipe).
func CaesarRoastedBroccoliRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	// ==================== CAESAR BREADCRUMBS RECIPE ====================
	breadcrumbsRecipeID := identifiers.New()

	// Get preparations
	meltPrep := enums.Preparations["melt"]
	stirPrep := enums.Preparations["stir"]
	cookPrep := enums.Preparations["cook"]
	toastPrep := enums.Preparations["toast"]
	seasonPrep := enums.Preparations["season"]
	transferPrep := enums.Preparations["transfer"]

	// Get ingredients for breadcrumbs
	saltedButter := enums.Ingredients["salted butter"]
	breadcrumbs := enums.Ingredients["breadcrumbs"]
	anchovyPaste := enums.Ingredients["anchovy paste"]
	garlic := enums.Ingredients["garlic"]
	lemon := enums.Ingredients["lemon"]
	salt := enums.Ingredients["salt"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	rubberSpatula := enums.Instruments["rubber spatula"]

	// Get vessels
	smallNonstickSkillet := enums.Vessels["small nonstick skillet"]
	smallBowl := enums.Vessels["small bowl"]

	// Get bridge table entries for breadcrumbs
	meltButterVIP := enums.IngredientPreparations[meltPrep.ID][saltedButter.ID]
	meltSkilletVPV := enums.PreparationVessels[meltPrep.ID][smallNonstickSkillet.ID]

	stirAnchovyVIP := enums.IngredientPreparations[stirPrep.ID][anchovyPaste.ID]
	stirGarlicVIP := enums.IngredientPreparations[stirPrep.ID][garlic.ID]
	stirBreadcrumbsVIP := enums.IngredientPreparations[stirPrep.ID][breadcrumbs.ID]
	stirLemonVIP := enums.IngredientPreparations[stirPrep.ID][lemon.ID]
	stirSkilletVPV := enums.PreparationVessels[stirPrep.ID][smallNonstickSkillet.ID]
	stirSpatulaVPI := enums.PreparationInstruments[stirPrep.ID][rubberSpatula.ID]

	cookSkilletVPV := enums.PreparationVessels[cookPrep.ID][smallNonstickSkillet.ID]

	toastBreadcrumbsVIP := enums.IngredientPreparations[toastPrep.ID][breadcrumbs.ID]
	toastSkilletVPV := enums.PreparationVessels[toastPrep.ID][smallNonstickSkillet.ID]
	toastSpatulaVPI := enums.PreparationInstruments[toastPrep.ID][rubberSpatula.ID]

	seasonBreadcrumbsVIP := enums.IngredientPreparations[seasonPrep.ID][breadcrumbs.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]

	transferBreadcrumbsVIP := enums.IngredientPreparations[transferPrep.ID][breadcrumbs.ID]
	transferSmallBowlVPV := enums.PreparationVessels[transferPrep.ID][smallBowl.ID]

	// Measurement unit bridges for breadcrumbs
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[saltedButter.ID][tablespoonMeasurement.ID]
	breadcrumbsCupVIMU := enums.IngredientMeasurementUnits[breadcrumbs.ID][cupMeasurement.ID]
	anchovyTeaspoonVIMU := enums.IngredientMeasurementUnits[anchovyPaste.ID][teaspoonMeasurement.ID]
	garlicUnitVIMU := enums.IngredientMeasurementUnits[garlic.ID][unitMeasurement.ID]
	lemonTeaspoonVIMU := enums.IngredientMeasurementUnits[lemon.ID][teaspoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]

	// Breadcrumbs Step 0: Melt butter in a small nonstick skillet
	bcStep0ID := identifiers.New()
	bcStep0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              bcStep0ID,
		BelongsToRecipe: breadcrumbsRecipeID,
		PreparationID:   meltPrep.ID,
		Index:           0,
		Notes:           "In a small nonstick skillet, melt butter over medium-low heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              bcStep0ID,
				ValidIngredientPreparationID:     &meltButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				IngredientID:                     &saltedButter.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "salted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      bcStep0ID,
				ValidPreparationVesselID: &meltSkilletVPV.ID,
				VesselID:                 &smallNonstickSkillet.ID,
				Name:                     "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: bcStep0ID,
				Name:                "melted butter",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &tablespoonMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Breadcrumbs Step 1: Stir in anchovy paste and garlic
	bcStep1ID := identifiers.New()
	bcStep1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              bcStep1ID,
		BelongsToRecipe: breadcrumbsRecipeID,
		PreparationID:   stirPrep.ID,
		Index:           1,
		Notes:           "Stir in anchovy paste and garlic.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             bcStep1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &saltedButter.ID,
				MeasurementUnitID:               tablespoonMeasurement.ID,
				Name:                            "melted butter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              bcStep1ID,
				ValidIngredientPreparationID:     &stirAnchovyVIP.ID,
				ValidIngredientMeasurementUnitID: &anchovyTeaspoonVIMU.ID,
				IngredientID:                     &anchovyPaste.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "anchovy paste",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              bcStep1ID,
				ValidIngredientPreparationID:     &stirGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicUnitVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "garlic, minced",
				QuantityNotes:                    "1 small clove",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          bcStep1ID,
				ValidPreparationInstrumentID: &stirSpatulaVPI.ID,
				InstrumentID:                 &rubberSpatula.ID,
				Name:                         "flexible spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      bcStep1ID,
				ValidPreparationVesselID: &stirSkilletVPV.ID,
				VesselID:                 &smallNonstickSkillet.ID,
				Name:                     "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: bcStep1ID,
				Name:                "butter with anchovy and garlic",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Breadcrumbs Step 2: Cook until fragrant
	bcStep2ID := identifiers.New()
	bcStep2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              bcStep2ID,
		BelongsToRecipe: breadcrumbsRecipeID,
		PreparationID:   cookPrep.ID,
		Index:           2,
		Notes:           "Cook until fragrant, about 1 minute.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             bcStep2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &saltedButter.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "butter with anchovy and garlic",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      bcStep2ID,
				ValidPreparationVesselID: &cookSkilletVPV.ID,
				VesselID:                 &smallNonstickSkillet.ID,
				Name:                     "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: bcStep2ID,
				Name:                "fragrant butter mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Breadcrumbs Step 3: Add breadcrumbs and toss to coat
	bcStep3ID := identifiers.New()
	bcStep3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              bcStep3ID,
		BelongsToRecipe: breadcrumbsRecipeID,
		PreparationID:   stirPrep.ID,
		Index:           3,
		Notes:           "Add breadcrumbs and, using a flexible spatula, toss to coat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             bcStep3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &saltedButter.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "fragrant butter mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              bcStep3ID,
				ValidIngredientPreparationID:     &stirBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				IngredientID:                     &breadcrumbs.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "plain breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          bcStep3ID,
				ValidPreparationInstrumentID: &stirSpatulaVPI.ID,
				InstrumentID:                 &rubberSpatula.ID,
				Name:                         "flexible spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      bcStep3ID,
				ValidPreparationVesselID: &stirSkilletVPV.ID,
				VesselID:                 &smallNonstickSkillet.ID,
				Name:                     "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: bcStep3ID,
				Name:                "coated breadcrumbs",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	// Breadcrumbs Step 4: Toast breadcrumbs until golden brown
	bcStep4ID := identifiers.New()
	bcStep4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              bcStep4ID,
		BelongsToRecipe: breadcrumbsRecipeID,
		PreparationID:   toastPrep.ID,
		Index:           4,
		Notes:           "Cook, stirring constantly until breadcrumbs are golden brown, about 3 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              bcStep4ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &toastBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				IngredientID:                     &breadcrumbs.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "coated breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          bcStep4ID,
				ValidPreparationInstrumentID: &toastSpatulaVPI.ID,
				InstrumentID:                 &rubberSpatula.ID,
				Name:                         "flexible spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      bcStep4ID,
				ValidPreparationVesselID: &toastSkilletVPV.ID,
				VesselID:                 &smallNonstickSkillet.ID,
				Name:                     "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: bcStep4ID,
				Name:                "toasted breadcrumbs",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	// Breadcrumbs Step 5: Stir in lemon zest (off heat)
	bcStep5ID := identifiers.New()
	bcStep5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              bcStep5ID,
		BelongsToRecipe: breadcrumbsRecipeID,
		PreparationID:   stirPrep.ID,
		Index:           5,
		Notes:           "Off heat, stir in 1/2 teaspoon lemon zest.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             bcStep5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &breadcrumbs.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "toasted breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              bcStep5ID,
				ValidIngredientPreparationID:     &stirLemonVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonTeaspoonVIMU.ID,
				IngredientID:                     &lemon.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "lemon zest",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          bcStep5ID,
				ValidPreparationInstrumentID: &stirSpatulaVPI.ID,
				InstrumentID:                 &rubberSpatula.ID,
				Name:                         "flexible spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      bcStep5ID,
				ValidPreparationVesselID: &stirSkilletVPV.ID,
				VesselID:                 &smallNonstickSkillet.ID,
				Name:                     "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: bcStep5ID,
				Name:                "breadcrumbs with lemon zest",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	// Breadcrumbs Step 6: Season breadcrumbs with salt
	bcStep6ID := identifiers.New()
	bcStep6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              bcStep6ID,
		BelongsToRecipe: breadcrumbsRecipeID,
		PreparationID:   seasonPrep.ID,
		Index:           6,
		Notes:           "Season with salt to taste.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              bcStep6ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &seasonBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				IngredientID:                     &breadcrumbs.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "breadcrumbs with lemon zest",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              bcStep6ID,
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
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      bcStep6ID,
				ValidPreparationVesselID: &stirSkilletVPV.ID,
				VesselID:                 &smallNonstickSkillet.ID,
				Name:                     "small nonstick skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: bcStep6ID,
				Name:                "seasoned caesar breadcrumbs",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	// Breadcrumbs Step 7: Transfer breadcrumbs to bowl and let cool
	bcStep7ID := identifiers.New()
	bcStep7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              bcStep7ID,
		BelongsToRecipe: breadcrumbsRecipeID,
		PreparationID:   transferPrep.ID,
		Index:           7,
		Notes:           "Transfer to a bowl and let cool completely.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              bcStep7ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &transferBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				IngredientID:                     &breadcrumbs.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "seasoned caesar breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      bcStep7ID,
				ValidPreparationVesselID: &transferSmallBowlVPV.ID,
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
				BelongsToRecipeStep: bcStep7ID,
				Name:                "caesar breadcrumbs",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.25),
				},
			},
		},
	}

	caesarBreadcrumbsRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  breadcrumbsRecipeID,
		CreatedByUser:       userID,
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
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			bcStep0, bcStep1, bcStep2, bcStep3, bcStep4, bcStep5, bcStep6, bcStep7,
		},
	}

	// ==================== CAESAR ROASTED BROCCOLI RECIPE ====================
	broccoliRecipeID := identifiers.New()

	// Additional preparations for broccoli
	linePrep := enums.Preparations["line"]
	preheatPrep := enums.Preparations["preheat"]
	tossPrep := enums.Preparations["toss"]
	roastPrep := enums.Preparations["roast"]
	topPrep := enums.Preparations["top"]

	// Additional ingredients for broccoli
	broccoli := enums.Ingredients["broccoli"]
	oliveOil := enums.Ingredients["olive oil"]
	blackPepper := enums.Ingredients["black pepper"]
	parmesan := enums.Ingredients["parmesan cheese"]

	// Additional measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	gramMeasurement := enums.MeasurementUnits["gram"]

	// Additional instruments
	aluminumFoil := enums.Instruments["aluminum foil"]

	// Additional vessels
	bakingSheet := enums.Vessels["baking sheet"]
	largeBowl := enums.Vessels["large bowl"]
	servingPlatter := enums.Vessels["serving platter"]
	oven := enums.Vessels["oven"]

	// Bridge table entries for broccoli
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

	// Additional measurement unit bridges
	broccoliPoundVIMU := enums.IngredientMeasurementUnits[broccoli.ID][poundMeasurement.ID]
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]
	parmesanTablespoonVIMU := enums.IngredientMeasurementUnits[parmesan.ID][tablespoonMeasurement.ID]

	// Broccoli Step 0: Line baking sheet with aluminum foil
	brStep0ID := identifiers.New()
	brStep0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              brStep0ID,
		BelongsToRecipe: broccoliRecipeID,
		PreparationID:   linePrep.ID,
		Index:           0,
		Notes:           "Line a rimmed baking sheet with aluminum foil.",
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          brStep0ID,
				ValidPreparationInstrumentID: &lineFoilVPI.ID,
				InstrumentID:                 &aluminumFoil.ID,
				Name:                         "aluminum foil",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      brStep0ID,
				ValidPreparationVesselID: &lineBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: brStep0ID,
				Name:                "foil-lined baking sheet",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 1: Preheat oven to 500°F
	brStep1ID := identifiers.New()
	brStep1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              brStep1ID,
		BelongsToRecipe: broccoliRecipeID,
		PreparationID:   preheatPrep.ID,
		Index:           1,
		Notes:           "Adjust oven rack to upper position and preheat oven to 500°F (260°C).",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260),
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      brStep1ID,
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
				BelongsToRecipeStep: brStep1ID,
				Name:                "preheated oven",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 2: Place baking sheet in oven to preheat
	brStep2ID := identifiers.New()
	brStep2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              brStep2ID,
		BelongsToRecipe: broccoliRecipeID,
		PreparationID:   preheatPrep.ID,
		Index:           2,
		Notes:           "Place the foil-lined baking sheet on oven rack to preheat.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260),
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      brStep2ID,
				ValidPreparationVesselID: &preheatBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "foil-lined baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: brStep2ID,
				Name:                "preheated baking sheet",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 3: Toss broccoli with olive oil, salt, and pepper
	brStep3ID := identifiers.New()
	brStep3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              brStep3ID,
		BelongsToRecipe: broccoliRecipeID,
		PreparationID:   tossPrep.ID,
		Index:           3,
		Notes:           "In a large bowl, toss broccoli florets with olive oil, salt, and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep3ID,
				ValidIngredientPreparationID:     &tossBroccoliVIP.ID,
				ValidIngredientMeasurementUnitID: &broccoliPoundVIMU.ID,
				IngredientID:                     &broccoli.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "broccoli florets",
				QuantityNotes:                    "cut into 1 1/2 to 2-inch pieces",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep3ID,
				ValidIngredientPreparationID:     &tossOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				IngredientID:                     &oliveOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "extra-virgin olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep3ID,
				ValidIngredientPreparationID:     &tossSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Diamond Crystal kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep3ID,
				ValidIngredientPreparationID:     &tossPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperGramVIMU.ID,
				IngredientID:                     &blackPepper.ID,
				MeasurementUnitID:                gramMeasurement.ID,
				Name:                             "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      brStep3ID,
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
				BelongsToRecipeStep: brStep3ID,
				Name:                "seasoned broccoli",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 4: Add broccoli to preheated baking sheet
	brStep4ID := identifiers.New()
	brStep4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              brStep4ID,
		BelongsToRecipe: broccoliRecipeID,
		PreparationID:   transferPrep.ID,
		Index:           4,
		Notes:           "Carefully add broccoli to preheated baking sheet in a single layer.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep4ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &transferBroccoliVIP.ID,
				ValidIngredientMeasurementUnitID: &broccoliPoundVIMU.ID,
				IngredientID:                     &broccoli.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "seasoned broccoli",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      brStep4ID,
				ValidPreparationVesselID: &transferBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "preheated baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: brStep4ID,
				Name:                "broccoli on baking sheet",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 5: Roast broccoli
	brStep5ID := identifiers.New()
	brStep5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              brStep5ID,
		BelongsToRecipe: broccoliRecipeID,
		PreparationID:   roastPrep.ID,
		Index:           5,
		Notes:           "Roast until broccoli is tender and deeply browned in spots, about 20 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](1200), // 20 minutes
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](260),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep5ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &roastBroccoliVIP.ID,
				ValidIngredientMeasurementUnitID: &broccoliPoundVIMU.ID,
				IngredientID:                     &broccoli.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "broccoli on baking sheet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      brStep5ID,
				ValidPreparationVesselID: &roastBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: brStep5ID,
				Name:                "roasted broccoli",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 6: Toss roasted broccoli with lemon zest
	brStep6ID := identifiers.New()
	brStep6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              brStep6ID,
		BelongsToRecipe: broccoliRecipeID,
		PreparationID:   tossPrep.ID,
		Index:           6,
		Notes:           "In the now empty bowl, toss broccoli with 1 teaspoon lemon zest.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             brStep6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &broccoli.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "roasted broccoli",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep6ID,
				ValidIngredientPreparationID:     &tossLemonVIP.ID,
				ValidIngredientMeasurementUnitID: &lemonTeaspoonVIMU.ID,
				IngredientID:                     &lemon.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "lemon zest",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      brStep6ID,
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
				BelongsToRecipeStep: brStep6ID,
				Name:                "broccoli with lemon zest",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 7: Transfer broccoli to serving platter
	brStep7ID := identifiers.New()
	brStep7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              brStep7ID,
		BelongsToRecipe: broccoliRecipeID,
		PreparationID:   transferPrep.ID,
		Index:           7,
		Notes:           "Transfer broccoli to a serving platter.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             brStep7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &broccoli.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "broccoli with lemon zest",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      brStep7ID,
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
				BelongsToRecipeStep: brStep7ID,
				Name:                "broccoli on serving platter",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Broccoli Step 8: Sprinkle with breadcrumbs and Parmigiano-Reggiano
	// This step references the Caesar Breadcrumbs recipe as a component
	brStep8ID := identifiers.New()
	brStep8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              brStep8ID,
		BelongsToRecipe: broccoliRecipeID,
		PreparationID:   topPrep.ID,
		Index:           8,
		Notes:           "Sprinkle with breadcrumbs and Parmigiano-Reggiano and serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep8ID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0),
				ValidIngredientPreparationID:     &topBroccoliVIP.ID,
				ValidIngredientMeasurementUnitID: &broccoliPoundVIMU.ID,
				IngredientID:                     &broccoli.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "broccoli on serving platter",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				// This ingredient references the Caesar Breadcrumbs recipe
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep8ID,
				RecipeStepProductRecipeID:        &breadcrumbsRecipeID,
				ProductOfRecipeStepIndex:         pointer.To[uint64](7), // Final step of breadcrumbs recipe
				ProductOfRecipeStepProductIndex:  pointer.To[uint64](0), // First product of that step
				ValidIngredientPreparationID:     &topBreadcrumbsVIP.ID,
				ValidIngredientMeasurementUnitID: &breadcrumbsCupVIMU.ID,
				IngredientID:                     &breadcrumbs.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "caesar breadcrumbs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              brStep8ID,
				ValidIngredientPreparationID:     &topParmesanVIP.ID,
				ValidIngredientMeasurementUnitID: &parmesanTablespoonVIMU.ID,
				IngredientID:                     &parmesan.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "grated Parmigiano-Reggiano cheese",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      brStep8ID,
				ValidPreparationVesselID: &topServingPlatterVPV.ID,
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
				BelongsToRecipeStep: brStep8ID,
				Name:                "caesar roasted broccoli",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	caesarRoastedBroccoliRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  broccoliRecipeID,
		CreatedByUser:       userID,
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
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			brStep0, brStep1, brStep2, brStep3, brStep4, brStep5, brStep6, brStep7, brStep8,
		},
	}

	// Return both recipes - breadcrumbs first since broccoli depends on it
	return []*mealplanning.RecipeDatabaseCreationInput{
		caesarBreadcrumbsRecipe,
		caesarRoastedBroccoliRecipe,
	}
}
