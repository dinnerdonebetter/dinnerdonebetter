package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// StirFriedGreenBeansRecipe creates the Stir-Fried Green Beans recipe.
// Source: https://www.seriouseats.com/stir-fried-green-beans-recipe
func StirFriedGreenBeansRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
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

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	cloveMeasurement := enums.MeasurementUnits["clove"]
	unitMeasurement := enums.MeasurementUnits["unit"]

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

	// Helper to safely get ID pointer from VIP
	vipID := func(v *mealplanning.ValidIngredientPreparation) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Helper to safely get ID pointer from VPV
	vpvID := func(v *mealplanning.ValidPreparationVessel) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Helper to safely get ID pointer from VPI
	vpiID := func(v *mealplanning.ValidPreparationInstrument) *string {
		if v == nil {
			return nil
		}
		return &v.ID
	}

	// Helper to safely get ID pointer from VIMU
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
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        trimPrep.ID,
		Index:                0,
		ExplicitInstructions: "Stem the green beans, removing the tough stem end from each bean.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     vipID(trimGreenBeansVIP),
				ValidIngredientMeasurementUnitID: vimuID(greenBeansPoundVIMU),
				Name:                             "green beans, washed and dried",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 0.5},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: vpiID(trimKnifeVPI),
				Name:                         "knife",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: vpvID(trimCuttingBoardVPV),
				Name:                     "cutting board",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "stemmed green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](0.5)},
			},
		},
	}

	// Step 1: Snap the green beans in half
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        snapPrep.ID,
		Index:                1,
		ExplicitInstructions: "Snap each green bean in half; set aside.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(snapGreenBeansVIP),
				Name:                            "stemmed green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 0.5},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: vpvID(snapCuttingBoardVPV),
				Name:                     "cutting board",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "snapped green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](0.5)},
			},
		},
	}

	// Step 2: Smash the garlic cloves
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        smashPrep.ID,
		Index:                2,
		ExplicitInstructions: "Smash the garlic cloves with the blunt side of the cleaver until the cloves are flattened.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     vipID(smashGarlicVIP),
				ValidIngredientMeasurementUnitID: vimuID(garlicCloveVIMU),
				Name:                             "garlic cloves, peeled",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: vpiID(smashCleaverVPI),
				Name:                         "cleaver",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: vpvID(smashCuttingBoardVPV),
				Name:                     "cutting board",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "smashed garlic",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cloveMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](2)},
			},
		},
	}

	// Step 3: Preheat the wok until smoking
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        preheatPrep.ID,
		Index:                3,
		ExplicitInstructions: "Place your wok over the highest flame and let heat for 3 to 4 minutes, until the wok is smoking.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](180),
			Max: pointer.To[uint32](240),
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: vpvID(preheatWokVPV),
				Name:                     "wok",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: smokingState.ID,
				Notes:             "wok should be smoking hot",
				Ingredients:       []uint64{},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "preheated wok",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	// Step 4: Swirl in the oil
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        swirPrep.ID,
		Index:                4,
		ExplicitInstructions: "Quickly swirl in the oil or cooking fat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     vipID(swirlLardVIP),
				ValidIngredientMeasurementUnitID: vimuID(lardTablespoonVIMU),
				Name:                             "lard",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        vpvID(swirlWokVPV),
				Name:                            "preheated wok",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "oil in wok",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				Name:  "wok with oil",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 5: Add the smashed garlic
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                5,
		ExplicitInstructions: "Toss in the smashed garlic.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "oil in wok",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(addGarlicVIP),
				Name:                            "smashed garlic",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(addWokVPV),
				Name:                            "wok with oil",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "garlic in wok",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				Name:  "wok with garlic",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 6: Stir the garlic in the oil
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                6,
		ExplicitInstructions: "Move the garlic around in the oil for about 5 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](5),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(stirGarlicVIP),
				Name:                            "garlic in wok",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: vpiID(stirSpatulaVPI),
				Name:                         "spatula or wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(stirWokVPV),
				Name:                            "wok with garlic",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "stirred garlic in wok",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				Name:  "wok with stirred garlic",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 7: Toss all the green beans into the wok
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        tossPrep.ID,
		Index:                7,
		ExplicitInstructions: "Toss all the green beans into the wok.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "stirred garlic in wok",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(tossGreenBeansVIP),
				Name:                            "snapped green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 0.5},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(tossWokVPV),
				Name:                            "wok with stirred garlic",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "green beans in wok",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				Name:  "wok with green beans",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 8: Stir-fry the vegetables
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                8,
		ExplicitInstructions: "Rapidly stir the vegetables with a spatula or wooden spoon. Stir-fry for 2 to 4 minutes, until the green beans are still crisp and barely raw in the center. To check the doneness, taste a green bean after 2 minutes into the cooking.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120),
			Max: pointer.To[uint32](240),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(stirGreenBeansVIP),
				Name:                            "green beans in wok",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: vpiID(stirSpatulaVPI),
				Name:                         "spatula",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
			{
				ValidPreparationInstrumentID: vpiID(stirWoodenSpoonVPI),
				Name:                         "wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(stirWokVPV),
				Name:                            "wok with green beans",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: crispState.ID,
				Notes:             "green beans should still be crisp and barely raw in the center",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "stir-fried green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				Name:  "wok with stir-fried green beans",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 9: Add salt and stock
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                9,
		ExplicitInstructions: "Toss in the salt and the stock.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "stir-fried green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ValidIngredientPreparationID:     vipID(addSaltVIP),
				ValidIngredientMeasurementUnitID: vimuID(saltTeaspoonVIMU),
				Name:                             "salt, or to taste",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 0.5},
			},
			{
				ValidIngredientPreparationID:     vipID(addStockVIP),
				ValidIngredientMeasurementUnitID: vimuID(stockTablespoonVIMU),
				Name:                             "vegetable stock",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 5},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(addWokVPV),
				Name:                            "wok with stir-fried green beans",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "green beans with salt and stock",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				Name:  "wok with seasoned green beans",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 10: Stir to coat the green beans
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                10,
		ExplicitInstructions: "Stir to coat the green beans with the liquid.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(stirGreenBeansVIP),
				Name:                            "green beans with salt and stock",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: vpiID(stirSpatulaVPI),
				Name:                         "spatula or wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(stirWokVPV),
				Name:                            "wok with seasoned green beans",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "coated green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				Name:  "wok with coated green beans",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 11: Cover the wok and turn off the heat
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        coverPrep.ID,
		Index:                11,
		ExplicitInstructions: "Place a lid on the wok and turn off the heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "coated green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(coverWokVPV),
				Name:                            "wok with coated green beans",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "covered green beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				Name:  "covered wok",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
	}

	// Step 12: Let the stock finish cooking the green beans
	step12 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        restPrep.ID,
		Index:                12,
		ExplicitInstructions: "Let the stock gently finish cooking the green beans. The green beans will be done when they are still crisp and just cooked through in the center. Serve piping hot.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    vipID(restGreenBeansVIP),
				Name:                            "covered green beans",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        vpvID(restWokVPV),
				Name:                            "covered wok",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: crispState.ID,
				Notes:             "green beans should still be crisp and just cooked through in the center",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:                "Stir-Fried Green Beans",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](2), Max: pointer.To[float32](3)},
			},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
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
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12,
			},
			PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
		},
	}
}
