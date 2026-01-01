package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// GrilledWholeCauliflowerRecipe creates the Grilled Whole Cauliflower with Teriyaki Sauce recipe.
// Source: https://www.seriouseats.com/grilled-whole-cauliflower-with-teriyaki-sauce-recipe-8678549
func GrilledWholeCauliflowerRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	// ==================== TERIYAKI SAUCE RECIPE ====================
	teriyakiRecipeID := identifiers.New()

	// Get preparations
	addPrep := enums.Preparations["add"]
	whiskPrep := enums.Preparations["whisk"]
	boilPrep := enums.Preparations["boil"]
	reducePrep := enums.Preparations["reduce"]
	stirPrep := enums.Preparations["stir"]

	// Get ingredients
	soySauce := enums.Ingredients["soy sauce"]
	sake := enums.Ingredients["sake"]
	mirin := enums.Ingredients["mirin"]
	sugar := enums.Ingredients["sugar"]
	dashiPowder := enums.Ingredients["dashi powder"]
	chickenFat := enums.Ingredients["rendered chicken fat"]
	sesameOil := enums.Ingredients["toasted sesame oil"]

	// Get measurement units
	cupMeasurement := enums.MeasurementUnits["cup"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	whisk := enums.Instruments["whisk"]
	spoon := enums.Instruments["spoon"]

	// Get vessels
	saucepan := enums.Vessels["saucepan"]

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

	whiskSaucepanVPV := getVPV("whisk", "saucepan")
	whiskWhiskVPI := getVPI("whisk", "whisk")

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
	ts0ID := identifiers.New()
	ts0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              ts0ID,
		BelongsToRecipe: teriyakiRecipeID,
		PreparationID:   addPrep.ID,
		Index:           0,
		Notes:           "In a medium saucepan, add shoyu, sake, mirin, sugar, and powdered dashi.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: ts0ID, ValidIngredientPreparationID: vipID(addSoySauceVIP), ValidIngredientMeasurementUnitID: vimuID(soySauceCupVIMU), IngredientID: &soySauce.ID, MeasurementUnitID: cupMeasurement.ID, Name: "Koikuchi shoyu, tamari, or saishikomi shoyu", QuantityNotes: "about 120ml", Quantity: types.Float32RangeWithOptionalMax{Min: 0.5}},
			{ID: identifiers.New(), BelongsToRecipeStep: ts0ID, ValidIngredientPreparationID: vipID(addSakeVIP), ValidIngredientMeasurementUnitID: vimuID(sakeCupVIMU), IngredientID: &sake.ID, MeasurementUnitID: cupMeasurement.ID, Name: "sake", QuantityNotes: "about 120ml", Quantity: types.Float32RangeWithOptionalMax{Min: 0.5}},
			{ID: identifiers.New(), BelongsToRecipeStep: ts0ID, ValidIngredientPreparationID: vipID(addMirinVIP), ValidIngredientMeasurementUnitID: vimuID(mirinCupVIMU), IngredientID: &mirin.ID, MeasurementUnitID: cupMeasurement.ID, Name: "mirin", QuantityNotes: "about 60ml", Quantity: types.Float32RangeWithOptionalMax{Min: 0.25}},
			{ID: identifiers.New(), BelongsToRecipeStep: ts0ID, ValidIngredientPreparationID: vipID(addSugarVIP), ValidIngredientMeasurementUnitID: vimuID(sugarTablespoonVIMU), IngredientID: &sugar.ID, MeasurementUnitID: tablespoonMeasurement.ID, Name: "granulated sugar", QuantityNotes: "about 60g", Quantity: types.Float32RangeWithOptionalMax{Min: 5}},
			{ID: identifiers.New(), BelongsToRecipeStep: ts0ID, ValidIngredientPreparationID: vipID(addDashiVIP), ValidIngredientMeasurementUnitID: vimuID(dashiTeaspoonVIMU), IngredientID: &dashiPowder.ID, MeasurementUnitID: teaspoonMeasurement.ID, Name: "powdered dashi", QuantityNotes: "about 8g", Quantity: types.Float32RangeWithOptionalMax{Min: 1.5}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts0ID, ValidPreparationVesselID: vpvID(addSaucepanVPV), VesselID: &saucepan.ID, Name: "medium saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts0ID, Name: "sauce ingredients in saucepan", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	ts1ID := identifiers.New()
	ts1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: ts1ID, BelongsToRecipe: teriyakiRecipeID, PreparationID: whiskPrep.ID, Index: 1, Notes: "Whisk until combined.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts1ID, ProductOfRecipeStepIndex: pointer.To[uint64](0), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), IngredientID: &soySauce.ID, MeasurementUnitID: unitMeasurement.ID, Name: "sauce ingredients in saucepan", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts1ID, ValidPreparationInstrumentID: vpiID(whiskWhiskVPI), InstrumentID: &whisk.ID, Name: "whisk", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts1ID, ValidPreparationVesselID: vpvID(whiskSaucepanVPV), VesselID: &saucepan.ID, Name: "saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts1ID, Name: "combined sauce mixture", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	ts2ID := identifiers.New()
	ts2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: ts2ID, BelongsToRecipe: teriyakiRecipeID, PreparationID: boilPrep.ID, Index: 2, Notes: "Bring mixture to a boil.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts2ID, ProductOfRecipeStepIndex: pointer.To[uint64](1), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), IngredientID: &soySauce.ID, MeasurementUnitID: unitMeasurement.ID, Name: "combined sauce mixture", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts2ID, ValidPreparationVesselID: vpvID(boilSaucepanVPV), VesselID: &saucepan.ID, Name: "saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts2ID, Name: "boiling sauce mixture", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	ts3ID := identifiers.New()
	ts3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: ts3ID, BelongsToRecipe: teriyakiRecipeID, PreparationID: reducePrep.ID, Index: 3,
		Notes:                  "Cook over medium heat, swirling pan occasionally, until temperature reaches 225°F (107℃) and sauce thickens and is reduced to a scant 1 cup, 12 to 16 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](720), Max: pointer.To[uint32](960)},
		Ingredients:            []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts3ID, ProductOfRecipeStepIndex: pointer.To[uint64](2), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), IngredientID: &soySauce.ID, MeasurementUnitID: unitMeasurement.ID, Name: "boiling sauce mixture", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments:            []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts3ID, ValidPreparationInstrumentID: vpiID(reduceSpoonVPI), InstrumentID: &spoon.ID, Name: "spoon", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:                []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts3ID, ValidPreparationVesselID: vpvID(reduceSaucepanVPV), VesselID: &saucepan.ID, Name: "saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		CompletionConditions:   []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts3ID, IngredientStateID: thickenedState.ID, Notes: "Sauce should reach 225°F (107℃) and reduce to a scant 1 cup", Optional: false}},
		Products:               []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts3ID, Name: "reduced teriyaki sauce", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &cupMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	ts4ID := identifiers.New()
	ts4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: ts4ID, BelongsToRecipe: teriyakiRecipeID, PreparationID: stirPrep.ID, Index: 4, Notes: "Off heat, stir in rendered chicken fat or butter and sesame oil. Set aside.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: ts4ID, ProductOfRecipeStepIndex: pointer.To[uint64](3), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), IngredientID: &soySauce.ID, MeasurementUnitID: cupMeasurement.ID, Name: "reduced teriyaki sauce", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
			{ID: identifiers.New(), BelongsToRecipeStep: ts4ID, ValidIngredientPreparationID: vipID(stirChickenFatVIP), ValidIngredientMeasurementUnitID: vimuID(chickenFatTablespoonVIMU), IngredientID: &chickenFat.ID, MeasurementUnitID: tablespoonMeasurement.ID, Name: "rendered chicken fat or unsalted butter", Quantity: types.Float32RangeWithOptionalMax{Min: 3}},
			{ID: identifiers.New(), BelongsToRecipeStep: ts4ID, ValidIngredientPreparationID: vipID(stirSesameOilVIP), ValidIngredientMeasurementUnitID: vimuID(sesameOilTablespoonVIMU), IngredientID: &sesameOil.ID, MeasurementUnitID: tablespoonMeasurement.ID, Name: "toasted sesame oil", QuantityNotes: "about 14g", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts4ID, ValidPreparationInstrumentID: vpiID(stirSpoonVPI), InstrumentID: &spoon.ID, Name: "spoon", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts4ID, ValidPreparationVesselID: vpvID(stirSaucepanVPV), VesselID: &saucepan.ID, Name: "saucepan", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: ts4ID, Name: "teriyaki sauce", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &cupMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	teriyakiSauceRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID: teriyakiRecipeID, CreatedByUser: userID, Name: "Teriyaki Sauce", Slug: "teriyaki-sauce",
		Source:              "https://www.seriouseats.com/grilled-whole-cauliflower-with-teriyaki-sauce-recipe-8678549",
		Description:         "A savory, sweet, and umami-rich teriyaki sauce made with shoyu, sake, mirin, sugar, and dashi, finished with chicken fat and sesame oil.",
		YieldsComponentType: mealplanning.MealComponentTypesUnspecified, EstimatedPortions: types.Float32RangeWithOptionalMax{Min: 1},
		PortionName: "cup", PluralPortionName: "cups", EligibleForMeals: false,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{ts0, ts1, ts2, ts3, ts4},
	}

	// ==================== GRILLED CAULIFLOWER RECIPE ====================
	cauliflowerRecipeID := identifiers.New()

	// Additional preparations for cauliflower
	trimPrep := enums.Preparations["trim"]
	slicePrep := enums.Preparations["slice"]
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
	largeBowl := enums.Vessels["large bowl"]
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

	whiskPotVPV := enums.PreparationVessels[whiskPrep.ID][pot.ID]

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
	gc0ID := identifiers.New()
	gc0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc0ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: trimPrep.ID, Index: 0, Notes: "Remove leaves from bottom of each cauliflower head.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc0ID, ValidIngredientPreparationID: vipID(trimCauliflowerVIP), ValidIngredientMeasurementUnitID: vimuID(cauliflowerPoundVIMU), IngredientID: &cauliflower.ID, MeasurementUnitID: poundMeasurement.ID, Name: "large heads cauliflower", QuantityNotes: "about 2 pounds each; 900g each", Quantity: types.Float32RangeWithOptionalMax{Min: 4}}},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc0ID, ValidPreparationInstrumentID: vpiID(trimKnifeVPI), InstrumentID: &knife.ID, Name: "knife", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc0ID, Name: "trimmed cauliflower heads", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &poundMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](4)}}},
	}

	// Step 1: Slice stem off each head
	gc1ID := identifiers.New()
	gc1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc1ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: slicePrep.ID, Index: 1, Notes: "Using a sharp knife, slice stem off of each head so that cauliflower sits evenly on flat surface. Do not cut out the core.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc1ID, ProductOfRecipeStepIndex: pointer.To[uint64](0), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(sliceCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: poundMeasurement.ID, Name: "trimmed cauliflower heads", Quantity: types.Float32RangeWithOptionalMax{Min: 4}}},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc1ID, ValidPreparationInstrumentID: vpiID(sliceKnifeVPI), InstrumentID: &knife.ID, Name: "sharp knife", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc1ID, Name: "prepared cauliflower heads", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &poundMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](4)}}},
	}

	// Step 2: Add water and salt to container
	gc2ID := identifiers.New()
	gc2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc2ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: addPrep.ID, Index: 2, Notes: "In a large container (8-quart) or large stock pot, add water and salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc2ID, ValidIngredientPreparationID: vipID(addWaterVIP), ValidIngredientMeasurementUnitID: vimuID(waterLiterVIMU), IngredientID: &water.ID, MeasurementUnitID: literMeasurement.ID, Name: "water", QuantityNotes: "3 quarts (2.84L)", Quantity: types.Float32RangeWithOptionalMax{Min: 2.84}},
			{ID: identifiers.New(), BelongsToRecipeStep: gc2ID, ValidIngredientPreparationID: vipID(addSaltVIP), ValidIngredientMeasurementUnitID: vimuID(saltCupVIMU), IngredientID: &salt.ID, MeasurementUnitID: cupMeasurement.ID, Name: "kosher salt", QuantityNotes: "3/4 cup plus 2 tablespoons (126g)", Quantity: types.Float32RangeWithOptionalMax{Min: 0.875}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc2ID, ValidPreparationVesselID: vpvID(addPotVPV), VesselID: &pot.ID, Name: "8-quart container or large stock pot", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc2ID, Name: "water and salt in container", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 3: Whisk until salt dissolved
	gc3ID := identifiers.New()
	gc3SaltIngredientID := identifiers.New()
	gc3CompletionConditionID := identifiers.New()
	gc3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc3ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: whiskPrep.ID, Index: 3, Notes: "Whisk water and salt until dissolved.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: gc3SaltIngredientID, BelongsToRecipeStep: gc3ID, ProductOfRecipeStepIndex: pointer.To[uint64](2), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), IngredientID: &water.ID, MeasurementUnitID: unitMeasurement.ID, Name: "water and salt in container", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc3ID, ValidPreparationInstrumentID: vpiID(whiskWhiskVPI), InstrumentID: &whisk.ID, Name: "whisk", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc3ID, ValidPreparationVesselID: vpvID(whiskPotVPV), VesselID: &pot.ID, Name: "container", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc3ID, Name: "saltwater brine", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{{
			ID: gc3CompletionConditionID, BelongsToRecipeStep: gc3ID, IngredientStateID: dissolvedState.ID, Notes: "Salt should be completely dissolved in the water", Optional: false,
			Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStepCompletionCondition: gc3CompletionConditionID, RecipeStepIngredient: gc3SaltIngredientID}},
		}},
	}

	// Step 4: Submerge cauliflower in brine
	gc4ID := identifiers.New()
	gc4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc4ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: submergePrep.ID, Index: 4, Notes: "Place cauliflower in saltwater brine, core side up, making sure that cauliflower is submerged.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc4ID, ProductOfRecipeStepIndex: pointer.To[uint64](1), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(submergeCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: poundMeasurement.ID, Name: "prepared cauliflower heads", Quantity: types.Float32RangeWithOptionalMax{Min: 4}},
			{ID: identifiers.New(), BelongsToRecipeStep: gc4ID, ProductOfRecipeStepIndex: pointer.To[uint64](3), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), IngredientID: &water.ID, MeasurementUnitID: unitMeasurement.ID, Name: "saltwater brine", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc4ID, ValidPreparationVesselID: vpvID(submergePotVPV), VesselID: &pot.ID, Name: "container with brine", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc4ID, Name: "cauliflower submerged in brine", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 5: Brine at room temperature
	gc5ID := identifiers.New()
	gc5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc5ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: brinePrep.ID, Index: 5,
		Notes:                  "Cover and let sit at room temperature for at least 3 hours and up to 6 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](10800), Max: pointer.To[uint32](21600)}, // 3-6 hours
		Ingredients:            []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc5ID, ProductOfRecipeStepIndex: pointer.To[uint64](4), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(brineCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "cauliflower submerged in brine", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:                []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc5ID, ValidPreparationVesselID: vpvID(brinePotVPV), VesselID: &pot.ID, Name: "covered container", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:               []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc5ID, Name: "brined cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 6: Light charcoal in chimney starter
	gc6ID := identifiers.New()
	gc6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc6ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: lightPrep.ID, Index: 6,
		Notes:       "For a charcoal grill: Open bottom vent completely. Light large chimney starter filled with charcoal briquettes (6 quarts). When top coals are partially covered with ash, pour evenly over half of grill.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc6ID, ValidIngredientPreparationID: vipID(lightCharcoalVIP), IngredientID: &charcoal.ID, MeasurementUnitID: unitMeasurement.ID, Name: "charcoal briquettes", QuantityNotes: "6 quarts", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc6ID, ValidPreparationInstrumentID: vpiID(lightChimneyStarterVPI), InstrumentID: &chimneyStarterInstrument.ID, Name: "chimney starter", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc6ID, ValidPreparationVesselID: vpvID(lightChimneyStarterVPV), VesselID: &chimneyStarter.ID, Name: "chimney starter", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc6ID, Name: "lit charcoal", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 7: Preheat grill
	gc7ID := identifiers.New()
	gc7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc7ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: preheatPrep.ID, Index: 7,
		Notes:                  "Set cooking grate in place, cover, and open lid vent completely. Heat grill until hot, about 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](300)},
		Ingredients:            []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc7ID, ProductOfRecipeStepIndex: pointer.To[uint64](6), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), IngredientID: &charcoal.ID, MeasurementUnitID: unitMeasurement.ID, Name: "lit charcoal on grill", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:                []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc7ID, ValidPreparationVesselID: vpvID(preheatGrillVPV), VesselID: &grill.ID, Name: "charcoal grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:               []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc7ID, Name: "preheated grill", Type: mealplanning.RecipeStepProductVesselType, Index: 0, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 8: Remove cauliflower from brine and drain
	gc8ID := identifiers.New()
	gc8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc8ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: drainPrep.ID, Index: 8, Notes: "Remove cauliflower from brine, letting excess liquid drain back into container.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc8ID, ProductOfRecipeStepIndex: pointer.To[uint64](5), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(drainCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "brined cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc8ID, ValidPreparationInstrumentID: vpiID(drainTongsVPI), InstrumentID: &tongs.ID, Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc8ID, Name: "drained brined cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 9: Place cauliflower on cooler side of grill
	gc9ID := identifiers.New()
	gc9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc9ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: placePrep.ID, Index: 9, Notes: "Place both cauliflower heads, stem side down onto cooler side of grill, approximately 2 inches from edge of hot coals or primary burner.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc9ID, ProductOfRecipeStepIndex: pointer.To[uint64](8), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(placeCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "drained brined cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc9ID, ProductOfRecipeStepIndex: pointer.To[uint64](7), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidPreparationVesselID: vpvID(placeGrillingGrateVPV), VesselID: &grillingGrate.ID, Name: "grill grate (cool side)", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc9ID, Name: "cauliflower on grill", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 10: Cover and cook for 20 minutes
	gc10ID := identifiers.New()
	gc10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc10ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: grillPrep.ID, Index: 10,
		Notes:                  "Cover and cook for 20 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](1200)},
		Ingredients:            []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc10ID, ProductOfRecipeStepIndex: pointer.To[uint64](9), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(grillCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "cauliflower on grill", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments:            []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc10ID, ValidPreparationInstrumentID: vpiID(grillTongsVPI), InstrumentID: &tongs.ID, Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:                []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc10ID, ValidPreparationVesselID: vpvID(grillGrillVPV), VesselID: &grill.ID, Name: "covered grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:               []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc10ID, Name: "partially cooked cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 11: Brush first layer of sauce
	gc11ID := identifiers.New()
	gc11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc11ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: brushPrep.ID, Index: 11, Notes: "Uncover grill, and using a heatproof brush, brush one layer of reserved sauce over cauliflower heads.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc11ID, ProductOfRecipeStepIndex: pointer.To[uint64](10), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(brushCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "partially cooked cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
			{ID: identifiers.New(), BelongsToRecipeStep: gc11ID, RecipeStepProductRecipeID: &teriyakiRecipeID, IngredientID: &soySauce.ID, MeasurementUnitID: cupMeasurement.ID, Name: "teriyaki sauce", Quantity: types.Float32RangeWithOptionalMax{Min: 0.33}},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc11ID, ValidPreparationInstrumentID: vpiID(brushBrushVPI), InstrumentID: &brush.ID, Name: "heatproof brush", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc11ID, Name: "sauced cauliflower (first coat)", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 12: Continue cooking until tender
	gc12ID := identifiers.New()
	gc12CauliflowerIngredientID := identifiers.New()
	gc12CompletionConditionID := identifiers.New()
	gc12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc12ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: grillPrep.ID, Index: 12,
		Notes:                  "Cover and continue cooking until thermometer registers 175°F at the thickest part of the core, and cauliflower is tan, but not well browned yet, rotating cauliflower occasionally, 20 to 30 minutes longer.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](1200), Max: pointer.To[uint32](1800)},
		Ingredients:            []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: gc12CauliflowerIngredientID, BelongsToRecipeStep: gc12ID, ProductOfRecipeStepIndex: pointer.To[uint64](11), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(grillCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "sauced cauliflower (first coat)", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc12ID, ValidPreparationInstrumentID: vpiID(grillTongsVPI), InstrumentID: &tongs.ID, Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}},
			{ID: identifiers.New(), BelongsToRecipeStep: gc12ID, ValidPreparationInstrumentID: vpiID(grillThermometerVPI), InstrumentID: &thermometer.ID, Name: "thermometer", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc12ID, ValidPreparationVesselID: vpvID(grillGrillVPV), VesselID: &grill.ID, Name: "covered grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc12ID, Name: "tender cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{{
			ID: gc12CompletionConditionID, BelongsToRecipeStep: gc12ID, IngredientStateID: tenderState.ID, Notes: "Core should register 175°F (79°C), cauliflower should be tan but not well browned", Optional: false,
			Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStepCompletionCondition: gc12CompletionConditionID, RecipeStepIngredient: gc12CauliflowerIngredientID}},
		}},
	}

	// Step 13: Brush second layer of sauce
	gc13ID := identifiers.New()
	gc13 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc13ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: brushPrep.ID, Index: 13, Notes: "Uncover grill and brush cauliflower with a second layer of sauce.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc13ID, ProductOfRecipeStepIndex: pointer.To[uint64](12), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(brushCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "tender cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
			{ID: identifiers.New(), BelongsToRecipeStep: gc13ID, RecipeStepProductRecipeID: &teriyakiRecipeID, IngredientID: &soySauce.ID, MeasurementUnitID: cupMeasurement.ID, Name: "teriyaki sauce", Quantity: types.Float32RangeWithOptionalMax{Min: 0.33}},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc13ID, ValidPreparationInstrumentID: vpiID(brushBrushVPI), InstrumentID: &brush.ID, Name: "heatproof brush", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc13ID, Name: "sauced cauliflower (second coat)", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 14: Flip and place floret-side down over hot side
	gc14ID := identifiers.New()
	gc14 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc14ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: flipPrep.ID, Index: 14, Notes: "Using tongs, flip cauliflower and place floret-side down directly over the hottest part of grill (over the coals).",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc14ID, ProductOfRecipeStepIndex: pointer.To[uint64](13), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(flipCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "sauced cauliflower (second coat)", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc14ID, ValidPreparationInstrumentID: vpiID(flipTongsVPI), InstrumentID: &tongs.ID, Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc14ID, ValidPreparationVesselID: vpvID(flipGrillVPV), VesselID: &grill.ID, Name: "grill (hot side)", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc14ID, Name: "flipped cauliflower on hot side", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 15: Cover and cook until lightly browned
	gc15ID := identifiers.New()
	gc15CauliflowerIngredientID := identifiers.New()
	gc15CompletionConditionID := identifiers.New()
	gc15 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc15ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: grillPrep.ID, Index: 15,
		Notes:                  "Cover and cook until lightly browned, 3 to 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](180), Max: pointer.To[uint32](300)},
		Ingredients:            []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: gc15CauliflowerIngredientID, BelongsToRecipeStep: gc15ID, ProductOfRecipeStepIndex: pointer.To[uint64](14), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(grillCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "flipped cauliflower on hot side", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:                []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc15ID, ValidPreparationVesselID: vpvID(grillGrillVPV), VesselID: &grill.ID, Name: "covered grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:               []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc15ID, Name: "lightly browned cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{{
			ID: gc15CompletionConditionID, BelongsToRecipeStep: gc15ID, IngredientStateID: brownedState.ID, Notes: "Cauliflower should be lightly browned", Optional: false,
			Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStepCompletionCondition: gc15CompletionConditionID, RecipeStepIngredient: gc15CauliflowerIngredientID}},
		}},
	}

	// Step 16: Flip, brush with final layer of sauce
	gc16ID := identifiers.New()
	gc16 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc16ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: brushPrep.ID, Index: 16, Notes: "Uncover grill, flip cauliflower heads stem side down, and brush florets all over with a final layer of sauce.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc16ID, ProductOfRecipeStepIndex: pointer.To[uint64](15), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(brushCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "lightly browned cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
			{ID: identifiers.New(), BelongsToRecipeStep: gc16ID, RecipeStepProductRecipeID: &teriyakiRecipeID, IngredientID: &soySauce.ID, MeasurementUnitID: cupMeasurement.ID, Name: "teriyaki sauce", Quantity: types.Float32RangeWithOptionalMax{Min: 0.33}},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc16ID, ValidPreparationInstrumentID: vpiID(brushBrushVPI), InstrumentID: &brush.ID, Name: "heatproof brush", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc16ID, Name: "sauced cauliflower (final coat)", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 17: Flip and grill until charred
	gc17ID := identifiers.New()
	gc17CauliflowerIngredientID := identifiers.New()
	gc17CompletionConditionID := identifiers.New()
	gc17 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc17ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: grillPrep.ID, Index: 17,
		Notes:                  "Flip cauliflower and place floret side down, cover, and cook until well browned and lightly charred, 3 to 5 minutes longer.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{Min: pointer.To[uint32](180), Max: pointer.To[uint32](300)},
		Ingredients:            []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: gc17CauliflowerIngredientID, BelongsToRecipeStep: gc17ID, ProductOfRecipeStepIndex: pointer.To[uint64](16), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(grillCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "sauced cauliflower (final coat)", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Instruments:            []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc17ID, ValidPreparationInstrumentID: vpiID(grillTongsVPI), InstrumentID: &tongs.ID, Name: "tongs", Quantity: types.Uint32RangeWithOptionalMax{Min: 1}}},
		Vessels:                []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc17ID, ValidPreparationVesselID: vpvID(grillGrillVPV), VesselID: &grill.ID, Name: "covered grill", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:               []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc17ID, Name: "charred cauliflower", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{{
			ID: gc17CompletionConditionID, BelongsToRecipeStep: gc17ID, IngredientStateID: lightlyCharredState.ID, Notes: "Cauliflower should be well browned and lightly charred", Optional: false,
			Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStepCompletionCondition: gc17CompletionConditionID, RecipeStepIngredient: gc17CauliflowerIngredientID}},
		}},
	}

	// Step 18: Transfer to plate
	gc18ID := identifiers.New()
	gc18 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc18ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: transferPrep.ID, Index: 18, Notes: "Transfer cauliflower to a plate.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc18ID, ProductOfRecipeStepIndex: pointer.To[uint64](17), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), ValidIngredientPreparationID: vipID(transferCauliflowerVIP), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "charred cauliflower", Quantity: types.Float32RangeWithOptionalMax{Min: 1}}},
		Vessels:     []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc18ID, ValidPreparationVesselID: vpvID(transferServingPlatterVPV), VesselID: &servingPlatter.ID, Name: "serving plate", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products:    []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc18ID, Name: "cauliflower on plate", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](1)}}},
	}

	// Step 19: Season with togarashi and serve
	gc19ID := identifiers.New()
	gc19 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID: gc19ID, BelongsToRecipe: cauliflowerRecipeID, PreparationID: seasonPrep.ID, Index: 19, Notes: "Season with togarashi. Serve with lemon wedges and remaining sauce.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc19ID, ProductOfRecipeStepIndex: pointer.To[uint64](18), ProductOfRecipeStepProductIndex: pointer.To[uint64](0), IngredientID: &cauliflower.ID, MeasurementUnitID: unitMeasurement.ID, Name: "cauliflower on plate", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
			{ID: identifiers.New(), BelongsToRecipeStep: gc19ID, ValidIngredientPreparationID: vipID(seasonTogarashiVIP), ValidIngredientMeasurementUnitID: vimuID(togarashiTeaspoonVIMU), IngredientID: &togarashi.ID, MeasurementUnitID: teaspoonMeasurement.ID, Name: "shichimi togarashi", Quantity: types.Float32RangeWithOptionalMax{Min: 1}},
		},
		Vessels:  []*mealplanning.RecipeStepVesselDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc19ID, ValidPreparationVesselID: vpvID(seasonServingPlatterVPV), VesselID: &servingPlatter.ID, Name: "serving plate", Quantity: types.Uint16RangeWithOptionalMax{Min: 1}}},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{{ID: identifiers.New(), BelongsToRecipeStep: gc19ID, Name: "Grilled Whole Cauliflower with Teriyaki Sauce", Type: mealplanning.RecipeStepProductIngredientType, Index: 0, MeasurementUnitID: &unitMeasurement.ID, MeasurementQuantity: types.OptionalFloat32Range{Min: pointer.To[float32](4), Max: pointer.To[float32](8)}}},
	}

	// Create prep task for brining ahead of time
	prepTask1ID := identifiers.New()
	prepTask1 := &mealplanning.RecipePrepTaskDatabaseCreationInput{
		ID: prepTask1ID, BelongsToRecipe: cauliflowerRecipeID,
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
		TaskSteps: []*mealplanning.RecipePrepTaskStepDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: gc0ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: gc1ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: gc2ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: gc3ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: gc4ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: gc5ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: true},
		},
	}

	grilledCauliflowerRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID: cauliflowerRecipeID, CreatedByUser: userID, Name: "Grilled Whole Cauliflower with Teriyaki Sauce", Slug: "grilled-whole-cauliflower-with-teriyaki-sauce",
		Source:              "https://www.seriouseats.com/grilled-whole-cauliflower-with-teriyaki-sauce-recipe-8678549",
		Description:         "Burnished, lightly charred domed cauliflower heads slathered in a savory teriyaki sauce. Brining the whole heads ensures deep seasoning, while low-and-slow grilling followed by high-heat charring produces tender, smoky cauliflower with entrée energy.",
		YieldsComponentType: mealplanning.MealComponentTypesMain, EstimatedPortions: types.Float32RangeWithOptionalMax{Min: 4, Max: pointer.To[float32](8)},
		PortionName: "serving", PluralPortionName: "servings", EligibleForMeals: true,
		Steps:     []*mealplanning.RecipeStepDatabaseCreationInput{gc0, gc1, gc2, gc3, gc4, gc5, gc6, gc7, gc8, gc9, gc10, gc11, gc12, gc13, gc14, gc15, gc16, gc17, gc18, gc19},
		PrepTasks: []*mealplanning.RecipePrepTaskDatabaseCreationInput{prepTask1},
	}

	// Suppress unused variable warnings
	_ = largeBowl

	return []*mealplanning.RecipeDatabaseCreationInput{
		teriyakiSauceRecipe,
		grilledCauliflowerRecipe,
	}
}
