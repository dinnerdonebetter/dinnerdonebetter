package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// GrilledWholeCauliflowerRecipe creates the Grilled Whole Cauliflower with Teriyaki Sauce recipe.
// Source: https://www.seriouseats.com/grilled-whole-cauliflower-with-teriyaki-sauce-recipe-8678549
func GrilledWholeCauliflowerRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// ==================== TERIYAKI SAUCE RECIPE ====================

	// Get preparations
	addPrep := enums.Preparations["add"]
	combinePrep := enums.Preparations["combine"]
	boilPrep := enums.Preparations["boil"]
	reducePrep := enums.Preparations["reduce"]
	stirPrep := enums.Preparations["stir"]

	// Get measurement units
	cupMeasurement := enums.MeasurementUnits["cup"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get ingredient states
	thickenedState := enums.IngredientStates["at desired consistency"]

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

	// Bridge entries for teriyaki sauce
	addSoySauceVIP := getVIP("add", "soy sauce")
	addSakeVIP := getVIP("add", "sake")
	addMirinVIP := getVIP("add", "mirin")
	addSugarVIP := getVIP("add", "sugar")
	addDashiVIP := getVIP("add", "dashi powder")
	addSaucepanVPV := getVPV("add", "saucepan")

	combineSaucepanVPV := getVPV("combine", "saucepan")

	boilSaucepanVPV := getVPV("boil", "saucepan")

	reduceSaucepanVPV := getVPV("reduce", "saucepan")
	reduceSpoonVPI := getVPI("reduce", "spoon")

	stirChickenFatVIP := getVIP("stir", "rendered chicken fat")
	stirSesameOilVIP := getVIP("stir", "toasted sesame oil")
	stirSaucepanVPV := getVPV("stir", "saucepan")
	stirSpoonVPI := getVPI("stir", "spoon")

	// Measurement unit bridges
	soySauceCupVIMU := getVIMU("soy sauce", "cup")
	sakeCupVIMU := getVIMU("sake", "cup")
	mirinCupVIMU := getVIMU("mirin", "cup")
	sugarTablespoonVIMU := getVIMU("sugar", "tablespoon")
	dashiTeaspoonVIMU := getVIMU("dashi powder", "teaspoon")
	chickenFatTablespoonVIMU := getVIMU("rendered chicken fat", "tablespoon")
	sesameOilTablespoonVIMU := getVIMU("toasted sesame oil", "tablespoon")

	// === TERIYAKI SAUCE STEPS ===
	ts0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: addPrep.ID,
		Index:         0,
		Notes:         "In a medium saucepan, add shoyu, sake, mirin, sugar, and powdered dashi.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ValidIngredientPreparationID: vipID(addSoySauceVIP), ValidIngredientMeasurementUnitID: vimuID(soySauceCupVIMU), Name: "soy sauce", QuantityNotes: "about 120ml", Quantity: types.Float32RangeWithOptionalMax{Min: 0.5}},
			{ValidIngredientPreparationID: vipID(addSakeVIP), ValidIngredientMeasurementUnitID: vimuID(sakeCupVIMU), Name: "sake", QuantityNotes: "about 120ml", Quantity: types.Float32RangeWithOptionalMax{Min: 0.5}},
			{ValidIngredientPreparationID: vipID(addMirinVIP), ValidIngredientMeasurementUnitID: vimuID(mirinCupVIMU), Name: "mirin", QuantityNotes: "about 60ml", Quantity: types.Float32RangeWithOptionalMax{Min: 0.25}},
			{ValidIngredientPreparationID: vipID(addSugarVIP), ValidIngredientMeasurementUnitID: vimuID(sugarTablespoonVIMU), Name: "granulated sugar", QuantityNotes: "about 60g", Quantity: types.Float32RangeWithOptionalMax{Min: 5}},
			{ValidIngredientPreparationID: vipID(addDashiVIP), ValidIngredientMeasurementUnitID: vimuID(dashiTeaspoonVIMU), Name: "powdered dashi", QuantityNotes: "about 8g", Quantity: types.Float32RangeWithOptionalMax{Min: 1.5}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(addSaucepanVPV), Name: "medium saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "sauce ingredients in saucepan", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	ts1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: combinePrep.ID,
		Index:         1,
		Notes:         "Combine until mixed.",
		Ingredients:   []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](0), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), Name: "sauce ingredients in saucepan", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:       []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(combineSaucepanVPV), Name: "saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:      []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "combined sauce mixture", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	ts2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: boilPrep.ID, Index: 2, Notes: "Bring mixture to a boil.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](1), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), Name: "combined sauce mixture", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(boilSaucepanVPV), Name: "saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "boiling sauce mixture", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	ts3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:          reducePrep.ID,
		Index:                  3,
		Notes:                  "Cook over medium heat, swirling pan occasionally, until temperature reaches 225°F (107℃) and sauce thickens and is reduced to a scant 1 cup, 12 to 16 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](720), Max: pointer.To[uint32](960)},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "boiling sauce mixture",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{ValidPreparationInstrumentID: vpiID(reduceSpoonVPI), Name: "spoon", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:              []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(reduceSaucepanVPV), Name: "saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{{IngredientStateID: thickenedState.ID, Notes: "Sauce should reach 225°F (107℃) and reduce to a scant 1 cup", Ingredients: []uint64{0}, Optional: false}},
		Products:             []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "reduced teriyaki sauce", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &cupMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	ts4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: stirPrep.ID, Index: 4, Notes: "Off heat, stir in butter and sesame oil. Set aside.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](3), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), Name: "reduced teriyaki sauce", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
			{ValidIngredientPreparationID: vipID(stirChickenFatVIP), ValidIngredientMeasurementUnitID: vimuID(chickenFatTablespoonVIMU), Name: "unsalted butter", Quantity: types.Float32RangeWithOptionalMax{Min: 3}},
			{ValidIngredientPreparationID: vipID(stirSesameOilVIP), ValidIngredientMeasurementUnitID: vimuID(sesameOilTablespoonVIMU), Name: "toasted sesame oil", QuantityNotes: "about 14g", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{{ValidPreparationInstrumentID: vpiID(stirSpoonVPI), Name: "spoon", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(stirSaucepanVPV), Name: "saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "teriyaki sauce", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &cupMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	teriyakiSauceRecipe := &mealplanning.RecipeCreationRequestInput{
		Name:                "Teriyaki Sauce",
		Slug:                "teriyaki-sauce",
		Source:              "https://www.seriouseats.com/grilled-whole-cauliflower-with-teriyaki-sauce-recipe-8678549",
		Description:         "A savory, sweet, and umami-rich teriyaki sauce made with shoyu, sake, mirin, sugar, and dashi, finished with chicken fat and sesame oil.",
		YieldsComponentType: mealplanning.MealComponentTypesUnspecified,
		EstimatedPortions:   types.Float32RangeWithOptionalMax{Min: 1},
		PortionName:         "cup",
		PluralPortionName:   "cups",
		EligibleForMeals:    false,
		Steps:               []*mealplanning.RecipeStepCreationRequestInput{ts0, ts1, ts2, ts3, ts4},
		PrepTasks:           []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
		Media:               []*mealplanning.RecipeMediaCreationRequestInput{},
	}

	// ==================== GRILLED CAULIFLOWER RECIPE ====================

	// Additional preparations for cauliflower
	trimPrep := enums.Preparations["trim"]
	slicePrep := enums.Preparations["slice"]
	dissolvePrep := enums.Preparations["dissolve"]
	submergePrep := enums.Preparations["submerge"]
	brinePrep := enums.Preparations["brine"]
	lightPrep := enums.Preparations["light"]
	preheatPrep := enums.Preparations["preheat"]
	drainPrep := enums.Preparations["drain"]
	placePrep := enums.Preparations["place"]
	grillPrep := enums.Preparations["grill"]
	brushPrep := enums.Preparations["brush"]
	flipPrep := enums.Preparations["flip"]
	transferPrep := enums.Preparations["transfer"]
	seasonPrep := enums.Preparations["season"]

	// Cauliflower ingredients
	cauliflower := enums.Ingredients["cauliflower"]
	water := enums.Ingredients["water"]
	salt := enums.Ingredients["salt"]
	// lemon wedges are served separately, not as a recipe step ingredient
	togarashi := enums.Ingredients["shichimi togarashi"]
	charcoal := enums.Ingredients["charcoal briquettes"]

	// Measurement units
	literMeasurement := enums.MeasurementUnits["liter"]
	poundMeasurement := enums.MeasurementUnits["pound"]

	// Instruments
	knife := enums.Instruments["knife"]
	tongs := enums.Instruments["tongs"]
	brush := enums.Instruments["brush"]
	thermometer := enums.Instruments["instant-read thermometer"]
	chimneyStarterInstrument := enums.Instruments["chimney starter"]

	// Vessels
	pot := enums.Vessels["pot"]
	grill := enums.Vessels["grill"]
	grillingGrate := enums.Vessels["grilling grate"]
	chimneyStarter := enums.Vessels["chimney starter"]
	servingPlatter := enums.Vessels["serving platter"]

	// Ingredient states
	tenderState := enums.IngredientStates["tender"]
	brownedState := enums.IngredientStates["browned"]
	dissolvedState := enums.IngredientStates["dissolved"]
	lightlyCharredState := enums.IngredientStates["lightly charred"]

	// Bridge entries for cauliflower steps
	trimCauliflowerVIP := enums.IngredientPreparations[trimPrep.ID][cauliflower.ID]
	trimKnifeVPI := enums.PreparationInstruments[trimPrep.ID][knife.ID]

	sliceCauliflowerVIP := enums.IngredientPreparations[slicePrep.ID][cauliflower.ID]
	sliceKnifeVPI := enums.PreparationInstruments[slicePrep.ID][knife.ID]

	addWaterVIP := enums.IngredientPreparations[addPrep.ID][water.ID]
	addSaltVIP := enums.IngredientPreparations[addPrep.ID][salt.ID]
	addPotVPV := enums.PreparationVessels[addPrep.ID][pot.ID]

	dissolvePotVPV := enums.PreparationVessels[dissolvePrep.ID][pot.ID]

	submergeCauliflowerVIP := enums.IngredientPreparations[submergePrep.ID][cauliflower.ID]
	submergePotVPV := enums.PreparationVessels[submergePrep.ID][pot.ID]

	brineCauliflowerVIP := enums.IngredientPreparations[brinePrep.ID][cauliflower.ID]
	brinePotVPV := enums.PreparationVessels[brinePrep.ID][pot.ID]

	lightCharcoalVIP := enums.IngredientPreparations[lightPrep.ID][charcoal.ID]
	lightChimneyStarterVPV := enums.PreparationVessels[lightPrep.ID][chimneyStarter.ID]
	lightChimneyStarterVPI := enums.PreparationInstruments[lightPrep.ID][chimneyStarterInstrument.ID]

	preheatGrillVPV := enums.PreparationVessels[preheatPrep.ID][grill.ID]

	drainCauliflowerVIP := enums.IngredientPreparations[drainPrep.ID][cauliflower.ID]
	drainTongsVPI := enums.PreparationInstruments[drainPrep.ID][tongs.ID]

	placeCauliflowerVIP := enums.IngredientPreparations[placePrep.ID][cauliflower.ID]
	placeGrillingGrateVPV := enums.PreparationVessels[placePrep.ID][grillingGrate.ID]

	grillCauliflowerVIP := enums.IngredientPreparations[grillPrep.ID][cauliflower.ID]
	grillGrillVPV := enums.PreparationVessels[grillPrep.ID][grill.ID]
	grillTongsVPI := enums.PreparationInstruments[grillPrep.ID][tongs.ID]
	grillThermometerVPI := enums.PreparationInstruments[grillPrep.ID][thermometer.ID]

	brushCauliflowerVIP := enums.IngredientPreparations[brushPrep.ID][cauliflower.ID]
	brushBrushVPI := enums.PreparationInstruments[brushPrep.ID][brush.ID]

	flipCauliflowerVIP := enums.IngredientPreparations[flipPrep.ID][cauliflower.ID]
	flipGrillVPV := enums.PreparationVessels[flipPrep.ID][grill.ID]
	flipTongsVPI := enums.PreparationInstruments[flipPrep.ID][tongs.ID]

	transferCauliflowerVIP := enums.IngredientPreparations[transferPrep.ID][cauliflower.ID]
	transferServingPlatterVPV := enums.PreparationVessels[transferPrep.ID][servingPlatter.ID]

	seasonTogarashiVIP := enums.IngredientPreparations[seasonPrep.ID][togarashi.ID]
	seasonServingPlatterVPV := enums.PreparationVessels[seasonPrep.ID][servingPlatter.ID]

	// Measurement unit bridges
	cauliflowerPoundVIMU := enums.IngredientMeasurementUnits[cauliflower.ID][poundMeasurement.ID]
	waterLiterVIMU := enums.IngredientMeasurementUnits[water.ID][literMeasurement.ID]
	saltCupVIMU := enums.IngredientMeasurementUnits[salt.ID][cupMeasurement.ID]
	togarashiTeaspoonVIMU := enums.IngredientMeasurementUnits[togarashi.ID][teaspoonMeasurement.ID]

	// === GRILLED CAULIFLOWER STEPS ===
	// Step 0: Remove leaves from cauliflower
	gc0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: trimPrep.ID,
		Index:         0,
		Notes:         "Remove leaves from bottom of each cauliflower head.",
		Ingredients:   []*mealplanning.RecipeStepIngredientCreationRequestInput{{ValidIngredientPreparationID: vipID(trimCauliflowerVIP), ValidIngredientMeasurementUnitID: vimuID(cauliflowerPoundVIMU), Name: "large heads cauliflower", QuantityNotes: "about 2 pounds each; 900g each", Quantity: types.Float32RangeWithOptionalMax{Min: 4}}},
		Instruments:   []*mealplanning.RecipeStepInstrumentCreationRequestInput{{ValidPreparationInstrumentID: vpiID(trimKnifeVPI), Name: "knife", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:      []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "trimmed cauliflower heads", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &poundMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](4)}}},
	}

	// Step 1: Slice stem off each head
	gc1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: slicePrep.ID,
		Index:         1,
		Notes:         "Using a sharp knife, slice stem off of each head so that cauliflower sits evenly on flat surface. Do not cut out the core.",
		Ingredients:   []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](0), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(sliceCauliflowerVIP), Name: "trimmed cauliflower heads", Quantity: types.Float32RangeWithOptionalMax{Min: 4}}},
		Instruments:   []*mealplanning.RecipeStepInstrumentCreationRequestInput{{ValidPreparationInstrumentID: vpiID(sliceKnifeVPI), Name: "sharp knife", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:      []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "prepared cauliflower heads", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &poundMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](4)}}},
	}

	// Step 2: Add water and salt to container
	gc2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: addPrep.ID,
		Index:         2,
		Notes:         "In a large container (8-quart) or large stock pot, add water and salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ValidIngredientPreparationID: vipID(addWaterVIP), ValidIngredientMeasurementUnitID: vimuID(waterLiterVIMU), Name: "water", QuantityNotes: "3 quarts (2.84L)", Quantity: types.Float32RangeWithOptionalMax{Min: 2.84}},
			{ValidIngredientPreparationID: vipID(addSaltVIP), ValidIngredientMeasurementUnitID: vimuID(saltCupVIMU), Name: "kosher salt", QuantityNotes: "3/4 cup plus 2 tablespoons (126g)", Quantity: types.Float32RangeWithOptionalMax{Min: 0.875}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(addPotVPV), Name: "8-quart container or large stock pot", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "water and salt in container", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 3: Dissolve salt in water
	gc3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: dissolvePrep.ID, Index: 3, Notes: "Dissolve salt in water.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](2), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), Name: "water and salt in container", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(dissolvePotVPV), Name: "container", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "saltwater brine", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{{
			IngredientStateID: dissolvedState.ID, Notes: "Salt should be completely dissolved in the water", Ingredients: []uint64{0}, Optional: false,
		}},
	}

	// Step 4: Submerge cauliflower in brine
	gc4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: submergePrep.ID, Index: 4, Notes: "Place cauliflower in saltwater brine, core side up, making sure that cauliflower is submerged.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](1), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(submergeCauliflowerVIP), Name: "prepared cauliflower heads", Quantity: types.Float32RangeWithOptionalMax{Min: 4}},
			{ProductOfRecipeStepIndex: pointer.To[uint64](3), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), Name: "saltwater brine", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(submergePotVPV), Name: "container with brine", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "cauliflower submerged in brine", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 5: Brine at room temperature
	gc5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:          brinePrep.ID,
		Index:                  5,
		Notes:                  "Cover and let sit at room temperature for at least 3 hours and up to 6 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](10800), Max: pointer.To[uint32](21600)}, // 3-6 hours
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](4), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(brineCauliflowerVIP), Name: "cauliflower submerged in brine", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(brinePotVPV), Name: "covered container", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "brined cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 6: Light charcoal in chimney starter
	gc6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: lightPrep.ID, Index: 6,
		Notes:       "For a charcoal grill: Open bottom vent completely. Light large chimney starter filled with charcoal briquettes (6 quarts). When top coals are partially covered with ash, pour evenly over half of grill.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{{ValidIngredientPreparationID: vipID(lightCharcoalVIP), Name: "charcoal briquettes", QuantityNotes: "6 quarts", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{{ValidPreparationInstrumentID: vpiID(lightChimneyStarterVPI), Name: "chimney starter", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(lightChimneyStarterVPV), Name: "chimney starter", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "lit charcoal", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 7: Preheat grill
	gc7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: preheatPrep.ID, Index: 7,
		Notes:                  "Set cooking grate in place, cover, and open lid vent completely. Heat grill until hot, about 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](300)},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](6), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), Name: "lit charcoal on grill", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(preheatGrillVPV), Name: "charcoal grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "preheated grill", Type: mealplanning.RecipeStepProductVesselType, Index: 0, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 8: Remove cauliflower from brine and drain
	gc8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: drainPrep.ID, Index: 8, Notes: "Remove cauliflower from brine, letting excess liquid drain back into container.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](5), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(drainCauliflowerVIP), Name: "brined cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{{ValidPreparationInstrumentID: vpiID(drainTongsVPI), Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "drained brined cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 9: Place cauliflower on cooler side of grill
	gc9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: placePrep.ID, Index: 9, Notes: "Place both cauliflower heads, stem side down onto cooler side of grill, approximately 2 inches from edge of hot coals or primary burner.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](8), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(placeCauliflowerVIP), Name: "drained brined cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](7), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidPreparationVesselID: vpvID(placeGrillingGrateVPV), Name: "grill grate (cool side)", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "cauliflower on grill", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 10: Cover and cook for 20 minutes
	gc10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:          grillPrep.ID,
		Index:                  10,
		Notes:                  "Cover and cook for 20 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](1200)},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](9), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(grillCauliflowerVIP), Name: "cauliflower on grill", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{ValidPreparationInstrumentID: vpiID(grillTongsVPI), Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(grillGrillVPV), Name: "covered grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "partially cooked cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 11: Brush first layer of sauce
	gc11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: brushPrep.ID, Index: 11, Notes: "Uncover grill, and using a heatproof brush, brush one layer of reserved sauce over cauliflower heads.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex: pointer.To[uint64](10), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(brushCauliflowerVIP), Name: "partially cooked cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				RecipeStepProductRecipeID: nil, // Will be resolved during recipe creation - references Teriyaki Sauce recipe
				Name:                      "teriyaki sauce", Quantity: types.Float32RangeWithOptionalMax{Min: 0.33},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{{ValidPreparationInstrumentID: vpiID(brushBrushVPI), Name: "heatproof brush", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "sauced cauliflower (first coat)", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 12: Continue cooking until tender
	gc12 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: grillPrep.ID, Index: 12,
		Notes:                  "Cover and continue cooking until thermometer registers 175°F at the thickest part of the core, and cauliflower is tan, but not well browned yet, rotating cauliflower occasionally, 20 to 30 minutes longer.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](1200), Max: pointer.To[uint32](1800)},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](11), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(grillCauliflowerVIP), Name: "sauced cauliflower (first coat)", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{ValidPreparationInstrumentID: vpiID(grillTongsVPI), Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}},
			{ValidPreparationInstrumentID: vpiID(grillThermometerVPI), Name: "thermometer", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(grillGrillVPV), Name: "covered grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "tender cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{{
			IngredientStateID: tenderState.ID, Notes: "Core should register 175°F (79°C), cauliflower should be tan but not well browned", Ingredients: []uint64{0}, Optional: false,
		}},
	}

	// Step 13: Brush second layer of sauce
	gc13 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: brushPrep.ID, Index: 13, Notes: "Uncover grill and brush cauliflower with a second layer of sauce.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](12), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(brushCauliflowerVIP), Name: "tender cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
			{RecipeStepProductRecipeID: nil, // Will be resolved during recipe creation - references Teriyaki Sauce recipe
				Name: "teriyaki sauce", Quantity: types.Float32RangeWithOptionalMax{Min: 0.33}},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{{ValidPreparationInstrumentID: vpiID(brushBrushVPI), Name: "heatproof brush", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "sauced cauliflower (second coat)", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 14: Flip and place floret-side down over hot side
	gc14 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: flipPrep.ID, Index: 14, Notes: "Using tongs, flip cauliflower and place floret-side down directly over the hottest part of grill (over the coals).",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](13), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(flipCauliflowerVIP), Name: "sauced cauliflower (second coat)", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{{ValidPreparationInstrumentID: vpiID(flipTongsVPI), Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(flipGrillVPV), Name: "grill (hot side)", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "flipped cauliflower on hot side", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 15: Cover and cook until lightly browned
	gc15 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: grillPrep.ID, Index: 15,
		Notes:                  "Cover and cook until lightly browned, 3 to 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](180), Max: pointer.To[uint32](300)},
		Ingredients:            []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](14), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(grillCauliflowerVIP), Name: "flipped cauliflower on hot side", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:                []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(grillGrillVPV), Name: "covered grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:               []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "lightly browned cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{{
			IngredientStateID: brownedState.ID, Notes: "Cauliflower should be lightly browned", Ingredients: []uint64{0}, Optional: false,
		}},
	}

	// Step 16: Flip, brush with final layer of sauce
	gc16 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: brushPrep.ID, Index: 16, Notes: "Uncover grill, flip cauliflower heads stem side down, and brush florets all over with a final layer of sauce.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](15), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(brushCauliflowerVIP), Name: "lightly browned cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
			{RecipeStepProductRecipeID: nil, // Will be resolved during recipe creation - references Teriyaki Sauce recipe
				Name: "teriyaki sauce", Quantity: types.Float32RangeWithOptionalMax{Min: 0.33}},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{ValidPreparationInstrumentID: vpiID(brushBrushVPI), Name: "heatproof brush", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "sauced cauliflower (final coat)", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 17: Flip and grill until charred
	gc17 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: grillPrep.ID, Index: 17,
		Notes:                  "Flip cauliflower and place floret side down, cover, and cook until well browned and lightly charred, 3 to 5 minutes longer.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](180), Max: pointer.To[uint32](300)},
		Ingredients:            []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](16), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(grillCauliflowerVIP), Name: "sauced cauliflower (final coat)", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments:            []*mealplanning.RecipeStepInstrumentCreationRequestInput{{ValidPreparationInstrumentID: vpiID(grillTongsVPI), Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:                []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(grillGrillVPV), Name: "covered grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:               []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "charred cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{{
			IngredientStateID: lightlyCharredState.ID, Notes: "Cauliflower should be well browned and lightly charred", Ingredients: []uint64{0}, Optional: false,
		}},
	}

	// Step 18: Transfer to plate
	gc18 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: transferPrep.ID, Index: 18, Notes: "Transfer cauliflower to a plate.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{{ProductOfRecipeStepIndex: pointer.To[uint64](17), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(transferCauliflowerVIP), Name: "charred cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(transferServingPlatterVPV), Name: "serving plate", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "cauliflower on plate", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 19: Season with togarashi and serve
	gc19 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID, Index: 19, Notes: "Season with togarashi. Serve with lemon wedges and remaining sauce.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{ProductOfRecipeStepIndex: pointer.To[uint64](18), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), Name: "cauliflower on plate", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
			{ValidIngredientPreparationID: vipID(seasonTogarashiVIP), ValidIngredientMeasurementUnitID: vimuID(togarashiTeaspoonVIMU), Name: "shichimi togarashi", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselCreationRequestInput{{ValidPreparationVesselID: vpvID(seasonServingPlatterVPV), Name: "serving plate", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{{Name: "Grilled Whole Cauliflower with Teriyaki Sauce", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](4), Max: pointer.To[float32](8)}}},
	}

	// Create prep task for brining ahead of time
	prepTask1 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
		Name:                        "Brine cauliflower",
		Description:                 "The cauliflower can be brined ahead of time to ensure deep seasoning throughout.",
		Notes:                       "Brining in a 4-5% salt solution for about 3 hours produces optimal results, but you can let the cauliflower sit for up to 6 hours.",
		Optional:                    true,
		ExplicitStorageInstructions: "Cover and store the brining cauliflower at room temperature for 3-6 hours.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeCovered,
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 10800,                     // 3 hours
			Max: pointer.To[uint32](21600), // 6 hours
		},
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 0, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 1, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 2, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 3, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 4, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 5, SatisfiesRecipeStep: true},
		},
	}

	grilledCauliflowerRecipe := &mealplanning.RecipeCreationRequestInput{
		Name:                "Grilled Whole Cauliflower with Teriyaki Sauce",
		Slug:                "grilled-whole-cauliflower-with-teriyaki-sauce",
		Source:              "https://www.seriouseats.com/grilled-whole-cauliflower-with-teriyaki-sauce-recipe-8678549",
		Description:         "Burnished, lightly charred domed cauliflower heads slathered in a savory teriyaki sauce. Brining the whole heads ensures deep seasoning, while low-and-slow grilling followed by high-heat charring produces tender, smoky cauliflower with entrée energy.",
		YieldsComponentType: mealplanning.MealComponentTypesMain,
		EstimatedPortions:   types.Float32RangeWithOptionalMax{Min: 4, Max: pointer.To[float32](8)},
		PortionName:         "serving",
		PluralPortionName:   "servings",
		EligibleForMeals:    true,
		Steps:               []*mealplanning.RecipeStepCreationRequestInput{gc0, gc1, gc2, gc3, gc4, gc5, gc6, gc7, gc8, gc9, gc10, gc11, gc12, gc13, gc14, gc15, gc16, gc17, gc18, gc19},
		PrepTasks:           []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{prepTask1},
		Media:               []*mealplanning.RecipeMediaCreationRequestInput{},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		teriyakiSauceRecipe,
		grilledCauliflowerRecipe,
	}
}
