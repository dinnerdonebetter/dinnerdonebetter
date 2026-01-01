package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// StirFriedGreenBeansRecipe creates the Stir-Fried Green Beans recipe.
// Source: https://www.seriouseats.com/stir-fried-green-beans-recipe
func StirFriedGreenBeansRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	trimPrep := enums.Preparations["trim"]
	snapPrep := enums.Preparations["snap"]
	smashPrep := enums.Preparations["smash"]
	preheatPrep := enums.Preparations["preheat"]
	swirPrep := enums.Preparations["swirl"]
	addPrep := enums.Preparations["add"]
	stirPrep := enums.Preparations["stir"]
	tossPrep := enums.Preparations["toss"]
	coverPrep := enums.Preparations["cover"]
	restPrep := enums.Preparations["rest"]

	// Get ingredients
	greenBeans := enums.Ingredients["green beans"]
	garlic := enums.Ingredients["garlic"]
	vegetableStock := enums.Ingredients["vegetable stock"]
	salt := enums.Ingredients["salt"]
	lard := enums.Ingredients["lard"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	cloveMeasurement := enums.MeasurementUnits["clove"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	knife := enums.Instruments["knife"]
	cleaver := enums.Instruments["cleaver"]
	spatula := enums.Instruments["spatula"]
	woodenSpoon := enums.Instruments["wooden spoon"]

	// Get vessels
	cuttingBoard := enums.Vessels["cutting board"]
	wok := enums.Vessels["wok"]

	// Get ingredient states
	smokingState := enums.IngredientStates["smoking"]
	crispState := enums.IngredientStates["crisp"]

	// Helper to safely get VIP bridge
	getVIP := func(prep, ingName string) *mealplanning.ValidIngredientPreparation {
		p := enums.Preparations[prep]
		i := enums.Ingredients[ingName]
		if p == nil || i == nil {
			return nil
		}
		prepMap := enums.IngredientPreparations[p.ID]
		if prepMap == nil {
			return nil
		}
		return prepMap[i.ID]
	}

	// Helper to safely get VPV bridge
	getVPV := func(prep, vesselName string) *mealplanning.ValidPreparationVessel {
		p := enums.Preparations[prep]
		v := enums.Vessels[vesselName]
		if p == nil || v == nil {
			return nil
		}
		prepMap := enums.PreparationVessels[p.ID]
		if prepMap == nil {
			return nil
		}
		return prepMap[v.ID]
	}

	// Helper to safely get VPI bridge
	getVPI := func(prep, instName string) *mealplanning.ValidPreparationInstrument {
		p := enums.Preparations[prep]
		i := enums.Instruments[instName]
		if p == nil || i == nil {
			return nil
		}
		prepMap := enums.PreparationInstruments[p.ID]
		if prepMap == nil {
			return nil
		}
		return prepMap[i.ID]
	}

	// Helper to safely get VIMU bridge
	getVIMU := func(ingName, unitName string) *mealplanning.ValidIngredientMeasurementUnit {
		i := enums.Ingredients[ingName]
		u := enums.MeasurementUnits[unitName]
		if i == nil || u == nil {
			return nil
		}
		ingMap := enums.IngredientMeasurementUnits[i.ID]
		if ingMap == nil {
			return nil
		}
		return ingMap[u.ID]
	}

	// Helper to safely get MealPlanTaskID pointer from VIP
	vipID := func(v *mealplanning.ValidIngredientPreparation) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Helper to safely get MealPlanTaskID pointer from VPV
	vpvID := func(v *mealplanning.ValidPreparationVessel) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Helper to safely get MealPlanTaskID pointer from VPI
	vpiID := func(v *mealplanning.ValidPreparationInstrument) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Helper to safely get MealPlanTaskID pointer from VIMU
	vimuID := func(v *mealplanning.ValidIngredientMeasurementUnit) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Bridge entries
	trimGreenBeansVIP := getVIP("trim", "green beans")
	trimKnifeVPI := getVPI("trim", "knife")
	trimCuttingBoardVPV := getVPV("trim", "cutting board")

	snapGreenBeansVIP := getVIP("snap", "green beans")
	snapCuttingBoardVPV := getVPV("snap", "cutting board")

	smashGarlicVIP := getVIP("smash", "garlic")
	smashCleaverVPI := getVPI("smash", "cleaver")
	smashCuttingBoardVPV := getVPV("smash", "cutting board")

	preheatWokVPV := getVPV("preheat", "wok")

	swirlLardVIP := getVIP("swirl", "lard")
	swirlWokVPV := getVPV("swirl", "wok")

	addGarlicVIP := getVIP("add", "garlic")
	addSaltVIP := getVIP("add", "salt")
	addStockVIP := getVIP("add", "vegetable stock")
	addWokVPV := getVPV("add", "wok")

	stirGarlicVIP := getVIP("stir", "garlic")
	stirGreenBeansVIP := getVIP("stir", "green beans")
	stirWokVPV := getVPV("stir", "wok")
	stirSpatulaVPI := getVPI("stir", "spatula")
	stirWoodenSpoonVPI := getVPI("stir", "wooden spoon")

	tossGreenBeansVIP := getVIP("toss", "green beans")
	tossWokVPV := getVPV("toss", "wok")

	coverWokVPV := getVPV("cover", "wok")

	restGreenBeansVIP := getVIP("rest", "green beans")
	restWokVPV := getVPV("rest", "wok")

	// Measurement unit bridges
	greenBeansPoundVIMU := getVIMU("green beans", "pound")
	garlicCloveVIMU := getVIMU("garlic", "clove")
	stockTablespoonVIMU := getVIMU("vegetable stock", "tablespoon")
	saltTeaspoonVIMU := getVIMU("salt", "teaspoon")
	lardTablespoonVIMU := getVIMU("lard", "tablespoon")

	// ==================== RECIPE STEPS ====================

	// Step 0: Stem the green beans
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   trimPrep.ID,
		Index:           0,
		Notes:           "Stem the green beans, removing the tough stem end from each bean.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     vipID(trimGreenBeansVIP),
				ValidIngredientMeasurementUnitID: vimuID(greenBeansPoundVIMU),
				IngredientID:                     &greenBeans.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "green beans, washed and dried",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 0.5},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0ID,
				ValidPreparationInstrumentID: vpiID(trimKnifeVPI),
				InstrumentID:                 &knife.ID,
				Name:                         "knife",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: vpvID(trimCuttingBoardVPV),
				VesselID:                 &cuttingBoard.ID,
				Name:                     "cutting board",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "stemmed green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](0.5)},
			},
		},
	}

	// Step 1: Snap the green beans in half
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   snapPrep.ID,
		Index:           1,
		Notes:           "Snap each green bean in half; set aside.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(snapGreenBeansVIP),
				IngredientID:                    &greenBeans.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "stemmed green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 0.5},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step1ID,
				ValidPreparationVesselID: vpvID(snapCuttingBoardVPV),
				VesselID:                 &cuttingBoard.ID,
				Name:                     "cutting board",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "snapped green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](0.5)},
			},
		},
	}

	// Step 2: Smash the garlic cloves
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   smashPrep.ID,
		Index:           2,
		Notes:           "Smash the garlic cloves with the blunt side of the cleaver until the cloves are flattened.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     vipID(smashGarlicVIP),
				ValidIngredientMeasurementUnitID: vimuID(garlicCloveVIMU),
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                cloveMeasurement.ID,
				Name:                             "garlic cloves, peeled",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step2ID,
				ValidPreparationInstrumentID: vpiID(smashCleaverVPI),
				InstrumentID:                 &cleaver.ID,
				Name:                         "cleaver",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: vpvID(smashCuttingBoardVPV),
				VesselID:                 &cuttingBoard.ID,
				Name:                     "cutting board",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "smashed garlic",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cloveMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](2)},
			},
		},
	}

	// Step 3: Preheat the wok until smoking
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   preheatPrep.ID,
		Index:           3,
		Notes:           "Place your wok over the highest flame and let heat for 3 to 4 minutes, until the wok is smoking.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180),
			Max: pointer.To[uint32](240),
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: vpvID(preheatWokVPV),
				VesselID:                 &wok.ID,
				Name:                     "wok",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				IngredientStateID:   smokingState.ID,
				Notes:               "wok should be smoking hot",
				Optional:            false,
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "preheated wok",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
			},
		},
	}

	// Step 4: Swirl in the oil
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   swirPrep.ID,
		Index:           4,
		Notes:           "Quickly swirl in the oil or cooking fat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ValidIngredientPreparationID:     vipID(swirlLardVIP),
				ValidIngredientMeasurementUnitID: vimuID(lardTablespoonVIMU),
				IngredientID:                     &lard.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "lard",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        vpvID(swirlWokVPV),
				VesselID:                        &wok.ID,
				Name:                            "preheated wok",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "oil in wok",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "wok with oil",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 5: Add the smashed garlic
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           5,
		Notes:           "Toss in the smashed garlic.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &lard.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "oil in wok",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(addGarlicVIP),
				IngredientID:                    &garlic.ID,
				MeasurementUnitID:               cloveMeasurement.ID,
				Name:                            "smashed garlic",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(addWokVPV),
				VesselID:                        &wok.ID,
				Name:                            "wok with oil",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "garlic in wok",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "wok with garlic",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 6: Stir the garlic in the oil
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           6,
		Notes:           "Move the garlic around in the oil for about 5 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](5),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(stirGarlicVIP),
				IngredientID:                    &garlic.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "garlic in wok",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step6ID,
				ValidPreparationInstrumentID: vpiID(stirSpatulaVPI),
				InstrumentID:                 &spatula.ID,
				Name:                         "spatula or wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(stirWokVPV),
				VesselID:                        &wok.ID,
				Name:                            "wok with garlic",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "stirred garlic in wok",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "wok with stirred garlic",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 7: Toss all the green beans into the wok
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   tossPrep.ID,
		Index:           7,
		Notes:           "Toss all the green beans into the wok.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &garlic.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "stirred garlic in wok",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(tossGreenBeansVIP),
				IngredientID:                    &greenBeans.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "snapped green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 0.5},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(tossWokVPV),
				VesselID:                        &wok.ID,
				Name:                            "wok with stirred garlic",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "green beans in wok",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "wok with green beans",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 8: Stir-fry the vegetables
	step8ID := identifiers.New()
	step8GreenBeansIngredientID := identifiers.New()
	step8CompletionConditionID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           8,
		Notes:           "Rapidly stir the vegetables with a spatula or wooden spoon. Stir-fry for 2 to 4 minutes, until the green beans are still crisp and barely raw in the center. To check the doneness, taste a green bean after 2 minutes into the cooking.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120),
			Max: pointer.To[uint32](240),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step8GreenBeansIngredientID,
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(stirGreenBeansVIP),
				IngredientID:                    &greenBeans.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "green beans in wok",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step8ID,
				ValidPreparationInstrumentID: vpiID(stirSpatulaVPI),
				InstrumentID:                 &spatula.ID,
				Name:                         "spatula",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step8ID,
				ValidPreparationInstrumentID: vpiID(stirWoodenSpoonVPI),
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(stirWokVPV),
				VesselID:                        &wok.ID,
				Name:                            "wok with green beans",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step8CompletionConditionID,
				BelongsToRecipeStep: step8ID,
				IngredientStateID:   crispState.ID,
				Notes:               "green beans should still be crisp and barely raw in the center",
				Optional:            false,
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step8CompletionConditionID,
						RecipeStepIngredient:                   step8GreenBeansIngredientID,
					},
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "stir-fried green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "wok with stir-fried green beans",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 9: Add salt and stock
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           9,
		Notes:           "Toss in the salt and the stock.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &greenBeans.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "stir-fried green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step9ID,
				ValidIngredientPreparationID:     vipID(addSaltVIP),
				ValidIngredientMeasurementUnitID: vimuID(saltTeaspoonVIMU),
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "salt, or to taste",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 0.5},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step9ID,
				ValidIngredientPreparationID:     vipID(addStockVIP),
				ValidIngredientMeasurementUnitID: vimuID(stockTablespoonVIMU),
				IngredientID:                     &vegetableStock.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "meat or vegetable stock",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 5},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(addWokVPV),
				VesselID:                        &wok.ID,
				Name:                            "wok with stir-fried green beans",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "green beans with salt and stock",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "wok with seasoned green beans",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 10: Stir to coat the green beans
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   stirPrep.ID,
		Index:           10,
		Notes:           "Stir to coat the green beans with the liquid.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(stirGreenBeansVIP),
				IngredientID:                    &greenBeans.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "green beans with salt and stock",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step10ID,
				ValidPreparationInstrumentID: vpiID(stirSpatulaVPI),
				InstrumentID:                 &spatula.ID,
				Name:                         "spatula or wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(stirWokVPV),
				VesselID:                        &wok.ID,
				Name:                            "wok with seasoned green beans",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "coated green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "wok with coated green beans",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 11: Cover the wok and turn off the heat
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   coverPrep.ID,
		Index:           11,
		Notes:           "Place a lid on the wok and turn off the heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &greenBeans.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "coated green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(coverWokVPV),
				VesselID:                        &wok.ID,
				Name:                            "wok with coated green beans",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "covered green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "covered wok",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
			},
		},
	}

	// Step 12: Let the stock finish cooking the green beans
	step12ID := identifiers.New()
	step12GreenBeansIngredientID := identifiers.New()
	step12CompletionConditionID := identifiers.New()
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           12,
		Notes:           "Let the stock gently finish cooking the green beans. The green beans will be done when they are still crisp and just cooked through in the center. Serve piping hot.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step12GreenBeansIngredientID,
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(restGreenBeansVIP),
				IngredientID:                    &greenBeans.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "covered green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(restWokVPV),
				VesselID:                        &wok.ID,
				Name:                            "covered wok",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step12CompletionConditionID,
				BelongsToRecipeStep: step12ID,
				IngredientStateID:   crispState.ID,
				Notes:               "green beans should still be crisp and just cooked through in the center",
				Optional:            false,
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step12CompletionConditionID,
						RecipeStepIngredient:                   step12GreenBeansIngredientID,
					},
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "Stir-Fried Green Beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](2), Max: pointer.To[float32](3)},
			},
		},
	}

	recipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Stir-Fried Green Beans",
		Slug:                "stir-fried-green-beans",
		Source:              "https://www.seriouseats.com/stir-fried-green-beans-recipe",
		Description:         "Green beans, when prepared correctly, are outrageously good. With their crispness still intact after being quickly stir-fried in a hot wok, they are a perfect vessel for your favorite sauces and dips.",
		YieldsComponentType: mealplanning.MealComponentTypesSide,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 2,
			Max: pointer.To[float32](3),
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{recipe}
}
