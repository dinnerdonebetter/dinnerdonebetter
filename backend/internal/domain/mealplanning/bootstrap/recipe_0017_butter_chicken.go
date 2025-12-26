package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ButterChickenRecipe creates the Butter Chicken recipe.
// Source: https://www.seriouseats.com/stovetop-butter-chicken
func ButterChickenRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	toastPrep := enums.Preparations["toast"]
	grindPrep := enums.Preparations["grind"]
	combinePrep := enums.Preparations["combine"]
	coatPrep := enums.Preparations["coat"]
	transferPrep := enums.Preparations["transfer"]
	linePrep := enums.Preparations["line"]
	soakPrep := enums.Preparations["soak"]
	microwavePrep := enums.Preparations["microwave"]
	heatPrep := enums.Preparations["heat"]
	cookPrep := enums.Preparations["cook"]
	stirPrep := enums.Preparations["stir"]
	addPrep := enums.Preparations["add"]
	crushPrep := enums.Preparations["crush"]
	simmerPrep := enums.Preparations["simmer"]
	preheatPrep := enums.Preparations["preheat"]
	broilPrep := enums.Preparations["broil"]
	blendPrep := enums.Preparations["blend"]

	// Get ingredients
	kasuriMethi := enums.Ingredients["kasuri methi"]
	fenugreekSeeds := enums.Ingredients["fenugreek seeds"]
	yogurt := enums.Ingredients["yogurt"]
	garamMasala := enums.Ingredients["garam masala"]
	salt := enums.Ingredients["salt"]
	kalaNamak := enums.Ingredients["kala namak"]
	ginger := enums.Ingredients["ginger"]
	chickenThighs := enums.Ingredients["boneless skinless chicken thighs"]
	chilesDeArbol := enums.Ingredients["dried chile de arbol"]
	brownCardamom := enums.Ingredients["brown cardamom"]
	greenCardamom := enums.Ingredients["green cardamom"]
	wholeCloves := enums.Ingredients["whole cloves"]
	cannedTomatoes := enums.Ingredients["fire-roasted canned tomatoes"]
	cashews := enums.Ingredients["raw cashews"]
	water := enums.Ingredients["water"]
	canolaOil := enums.Ingredients["canola oil"]
	whiteOnion := enums.Ingredients["white onion"]
	bakingSoda := enums.Ingredients["baking soda"]
	garlic := enums.Ingredients["garlic"]
	heavyCream := enums.Ingredients["heavy cream"]
	butter := enums.Ingredients["butter"]

	// Get measurement units
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	poundMeasurement := enums.MeasurementUnits["pound"]
	ounceMeasurement := enums.MeasurementUnits["ounce"]
	unitMeasurement := enums.MeasurementUnits["unit"]
	cloveMeasurement := enums.MeasurementUnits["clove"]

	// Get instruments
	spiceGrinder := enums.Instruments["spice grinder"]
	mortarAndPestle := enums.Instruments["mortar and pestle"]
	woodenSpoon := enums.Instruments["wooden spoon"]
	stickBlender := enums.Instruments["stick blender"]
	bareHands := enums.Instruments["bare hands"]
	aluminumFoil := enums.Instruments["aluminum foil"]
	blenderInst := enums.Instruments["blender"]
	tongs := enums.Instruments["tongs"]

	// Get vessels
	smallSkillet := enums.Vessels["small skillet"]
	mediumBowl := enums.Vessels["medium bowl"]
	bakingSheet := enums.Vessels["baking sheet"]
	microwaveSafeBowl := enums.Vessels["microwave-safe bowl"]
	dutchOven := enums.Vessels["dutch oven"]
	oven := enums.Vessels["oven"]
	servingBowl := enums.Vessels["serving bowl"]

	// Get bridge table entries for marinade ingredients
	toastKasuriMethiVIP := enums.IngredientPreparations[toastPrep.ID][kasuriMethi.ID]
	grindKasuriMethiVIP := enums.IngredientPreparations[grindPrep.ID][kasuriMethi.ID]
	combineYogurtVIP := enums.IngredientPreparations[combinePrep.ID][yogurt.ID]
	combineGaramMasalaVIP := enums.IngredientPreparations[combinePrep.ID][garamMasala.ID]
	combineKalaNamakVIP := enums.IngredientPreparations[combinePrep.ID][kalaNamak.ID]
	combineGingerVIP := enums.IngredientPreparations[combinePrep.ID][ginger.ID]
	combineKasuriMethiVIP := enums.IngredientPreparations[combinePrep.ID][kasuriMethi.ID]
	coatChickenVIP := enums.IngredientPreparations[coatPrep.ID][chickenThighs.ID]
	transferChickenVIP := enums.IngredientPreparations[transferPrep.ID][chickenThighs.ID]

	// Bridge table entries for sauce ingredients
	toastChilesVIP := enums.IngredientPreparations[toastPrep.ID][chilesDeArbol.ID]
	toastBrownCardamomVIP := enums.IngredientPreparations[toastPrep.ID][brownCardamom.ID]
	toastWholeCloveVIP := enums.IngredientPreparations[toastPrep.ID][wholeCloves.ID]
	grindChilesVIP := enums.IngredientPreparations[grindPrep.ID][chilesDeArbol.ID]
	grindBrownCardamomVIP := enums.IngredientPreparations[grindPrep.ID][brownCardamom.ID]
	grindWholeCloveVIP := enums.IngredientPreparations[grindPrep.ID][wholeCloves.ID]
	grindGaramMasalaVIP := enums.IngredientPreparations[grindPrep.ID][garamMasala.ID]
	grindSaltVIP := enums.IngredientPreparations[grindPrep.ID][salt.ID]
	soakCashewsVIP := enums.IngredientPreparations[soakPrep.ID][cashews.ID]
	soakWaterVIP := enums.IngredientPreparations[soakPrep.ID][water.ID]
	microwaveCashewsVIP := enums.IngredientPreparations[microwavePrep.ID][cashews.ID]
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][canolaOil.ID]
	cookOnionVIP := enums.IngredientPreparations[cookPrep.ID][whiteOnion.ID]
	cookBakingSodaVIP := enums.IngredientPreparations[cookPrep.ID][bakingSoda.ID]
	cookGingerVIP := enums.IngredientPreparations[cookPrep.ID][ginger.ID]
	cookGarlicVIP := enums.IngredientPreparations[cookPrep.ID][garlic.ID]
	addCashewsVIP := enums.IngredientPreparations[addPrep.ID][cashews.ID]
	addTomatoesVIP := enums.IngredientPreparations[addPrep.ID][cannedTomatoes.ID]
	addWaterVIP := enums.IngredientPreparations[addPrep.ID][water.ID]
	crushTomatoesVIP := enums.IngredientPreparations[crushPrep.ID][cannedTomatoes.ID]
	simmerTomatoesVIP := enums.IngredientPreparations[simmerPrep.ID][cannedTomatoes.ID]
	broilChickenVIP := enums.IngredientPreparations[broilPrep.ID][chickenThighs.ID]
	blendTomatoesVIP := enums.IngredientPreparations[blendPrep.ID][cannedTomatoes.ID]
	blendButterVIP := enums.IngredientPreparations[blendPrep.ID][butter.ID]
	blendCreamVIP := enums.IngredientPreparations[blendPrep.ID][heavyCream.ID]
	addButterVIP := enums.IngredientPreparations[addPrep.ID][butter.ID]
	addCreamVIP := enums.IngredientPreparations[addPrep.ID][heavyCream.ID]
	addChickenVIP := enums.IngredientPreparations[addPrep.ID][chickenThighs.ID]
	stirOnionVIP := enums.IngredientPreparations[stirPrep.ID][whiteOnion.ID]

	// Preparation-Instrument bridges
	toastSmallSkilletVPV := enums.PreparationVessels[toastPrep.ID][smallSkillet.ID]
	grindSpiceGrinderVPI := enums.PreparationInstruments[grindPrep.ID][spiceGrinder.ID]
	grindMortarVPI := enums.PreparationInstruments[grindPrep.ID][mortarAndPestle.ID]
	combineMediumBowlVPV := enums.PreparationVessels[combinePrep.ID][mediumBowl.ID]
	coatBareHandsVPI := enums.PreparationInstruments[coatPrep.ID][bareHands.ID]
	coatMediumBowlVPV := enums.PreparationVessels[coatPrep.ID][mediumBowl.ID]
	transferBakingSheetVPV := enums.PreparationVessels[transferPrep.ID][bakingSheet.ID]
	lineAluminumFoilVPI := enums.PreparationInstruments[linePrep.ID][aluminumFoil.ID]
	lineBakingSheetVPV := enums.PreparationVessels[linePrep.ID][bakingSheet.ID]
	soakMicrowaveBowlVPV := enums.PreparationVessels[soakPrep.ID][microwaveSafeBowl.ID]
	microwaveMicrowaveBowlVPV := enums.PreparationVessels[microwavePrep.ID][microwaveSafeBowl.ID]
	heatDutchOvenVPV := enums.PreparationVessels[heatPrep.ID][dutchOven.ID]
	cookWoodenSpoonVPI := enums.PreparationInstruments[cookPrep.ID][woodenSpoon.ID]
	cookDutchOvenVPV := enums.PreparationVessels[cookPrep.ID][dutchOven.ID]
	stirWoodenSpoonVPI := enums.PreparationInstruments[stirPrep.ID][woodenSpoon.ID]
	stirDutchOvenVPV := enums.PreparationVessels[stirPrep.ID][dutchOven.ID]
	addDutchOvenVPV := enums.PreparationVessels[addPrep.ID][dutchOven.ID]
	addWoodenSpoonVPI := enums.PreparationInstruments[addPrep.ID][woodenSpoon.ID]
	crushWoodenSpoonVPI := enums.PreparationInstruments[crushPrep.ID][woodenSpoon.ID]
	crushDutchOvenVPV := enums.PreparationVessels[crushPrep.ID][dutchOven.ID]
	simmerDutchOvenVPV := enums.PreparationVessels[simmerPrep.ID][dutchOven.ID]
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]
	broilBakingSheetVPV := enums.PreparationVessels[broilPrep.ID][bakingSheet.ID]
	broilOvenVPV := enums.PreparationVessels[broilPrep.ID][oven.ID]
	blendStickBlenderVPI := enums.PreparationInstruments[blendPrep.ID][stickBlender.ID]
	blendBlenderVPI := enums.PreparationInstruments[blendPrep.ID][blenderInst.ID]
	blendDutchOvenVPV := enums.PreparationVessels[blendPrep.ID][dutchOven.ID]
	transferDutchOvenVPV := enums.PreparationVessels[transferPrep.ID][dutchOven.ID]
	transferServingBowlVPV := enums.PreparationVessels[transferPrep.ID][servingBowl.ID]

	// Measurement unit bridges
	kasuriMethiTablespoonVIMU := enums.IngredientMeasurementUnits[kasuriMethi.ID][tablespoonMeasurement.ID]
	fenugreekSeedsTeaspoonVIMU := enums.IngredientMeasurementUnits[fenugreekSeeds.ID][teaspoonMeasurement.ID]
	yogurtCupVIMU := enums.IngredientMeasurementUnits[yogurt.ID][cupMeasurement.ID]
	garamMasalaTablespoonVIMU := enums.IngredientMeasurementUnits[garamMasala.ID][tablespoonMeasurement.ID]
	kalaNamakTeaspoonVIMU := enums.IngredientMeasurementUnits[kalaNamak.ID][teaspoonMeasurement.ID]
	gingerUnitVIMU := enums.IngredientMeasurementUnits[ginger.ID][unitMeasurement.ID]
	chickenThighsPoundVIMU := enums.IngredientMeasurementUnits[chickenThighs.ID][poundMeasurement.ID]
	chilesDeArbolUnitVIMU := enums.IngredientMeasurementUnits[chilesDeArbol.ID][unitMeasurement.ID]
	brownCardamomUnitVIMU := enums.IngredientMeasurementUnits[brownCardamom.ID][unitMeasurement.ID]
	greenCardamomUnitVIMU := enums.IngredientMeasurementUnits[greenCardamom.ID][unitMeasurement.ID]
	wholeClovesUnitVIMU := enums.IngredientMeasurementUnits[wholeCloves.ID][unitMeasurement.ID]
	cannedTomatoesOunceVIMU := enums.IngredientMeasurementUnits[cannedTomatoes.ID][ounceMeasurement.ID]
	cashewsOunceVIMU := enums.IngredientMeasurementUnits[cashews.ID][ounceMeasurement.ID]
	canolaOilTablespoonVIMU := enums.IngredientMeasurementUnits[canolaOil.ID][tablespoonMeasurement.ID]
	whiteOnionUnitVIMU := enums.IngredientMeasurementUnits[whiteOnion.ID][unitMeasurement.ID]
	bakingSodaTeaspoonVIMU := enums.IngredientMeasurementUnits[bakingSoda.ID][teaspoonMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
	heavyCreamCupVIMU := enums.IngredientMeasurementUnits[heavyCream.ID][cupMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	butterTablespoonVIMU := enums.IngredientMeasurementUnits[butter.ID][tablespoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]

	// Ignore unused variables for now
	_ = crushTomatoesVIP
	_ = addButterVIP
	_ = addCreamVIP
	_ = stirOnionVIP
	_ = grindMortarVPI
	_ = stirDutchOvenVPV
	_ = crushDutchOvenVPV
	_ = blendBlenderVPI
	_ = transferDutchOvenVPV
	_ = fenugreekSeedsTeaspoonVIMU
	_ = greenCardamomUnitVIMU
	_ = grindBrownCardamomVIP
	_ = grindWholeCloveVIP
	_ = tongs
	_ = stirWoodenSpoonVPI
	_ = crushWoodenSpoonVPI

	// Step 0: Toast fenugreek leaves in a small skillet over medium heat until fragrant
	step0ID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   toastPrep.ID,
		Index:           0,
		Notes:           "In a small skillet, toast fenugreek leaves (or fenugreek seeds, if using) over medium heat, tossing them constantly, until quite fragrant, about 30 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &toastKasuriMethiVIP.ID,
				ValidIngredientMeasurementUnitID: &kasuriMethiTablespoonVIMU.ID,
				IngredientID:                     &kasuriMethi.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "kasuri methi (fenugreek leaves)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0ID,
				ValidPreparationVesselID: &toastSmallSkilletVPV.ID,
				VesselID:                 &smallSkillet.ID,
				Name:                     "small skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0ID,
				Name:                "toasted fenugreek",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &tablespoonMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Grind toasted fenugreek to fine powder
	step1ID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   grindPrep.ID,
		Index:           1,
		Notes:           "Transfer toasted leaves to spice grinder or mortar and pestle and grind to fine powder. Set aside.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step1ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &grindKasuriMethiVIP.ID,
				IngredientID:                    &kasuriMethi.ID,
				MeasurementUnitID:               tablespoonMeasurement.ID,
				Name:                            "toasted fenugreek",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step1ID,
				ValidPreparationInstrumentID: &grindSpiceGrinderVPI.ID,
				InstrumentID:                 &spiceGrinder.ID,
				Name:                         "spice grinder",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "ground fenugreek for marinade",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &tablespoonMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Create aluminum foil boat on baking sheet
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   linePrep.ID,
		Index:           2,
		Notes:           "In the center of a rimmed baking sheet, create a roughly 9- by 13-inch aluminum-foil boat with 1-inch sides, and set aside.",
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step2ID,
				ValidPreparationInstrumentID: &lineAluminumFoilVPI.ID,
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
				BelongsToRecipeStep:      step2ID,
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
				BelongsToRecipeStep: step2ID,
				Name:                "prepared baking sheet with foil boat",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Combine marinade ingredients in medium bowl
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   combinePrep.ID,
		Index:           3,
		Notes:           "In a medium mixing bowl, stir together yogurt, garam masala, salt, black salt, grated ginger, and ground fenugreek leaves.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &combineYogurtVIP.ID,
				ValidIngredientMeasurementUnitID: &yogurtCupVIMU.ID,
				IngredientID:                     &yogurt.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "plain Greek yogurt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &combineGaramMasalaVIP.ID,
				ValidIngredientMeasurementUnitID: &garamMasalaTablespoonVIMU.ID,
				IngredientID:                     &garamMasala.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "garam masala",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &combineKalaNamakVIP.ID,
				ValidIngredientMeasurementUnitID: &kalaNamakTeaspoonVIMU.ID,
				IngredientID:                     &kalaNamak.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "kala namak (black salt)",
				QuantityNotes:                    "optional",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientPreparationID:     &combineGingerVIP.ID,
				ValidIngredientMeasurementUnitID: &gingerUnitVIMU.ID,
				IngredientID:                     &ginger.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "fresh ginger, peeled and finely grated",
				QuantityNotes:                    "1-inch piece",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &combineKasuriMethiVIP.ID,
				IngredientID:                    &kasuriMethi.ID,
				MeasurementUnitID:               tablespoonMeasurement.ID,
				Name:                            "ground fenugreek for marinade",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step3ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Diamond Crystal kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &combineMediumBowlVPV.ID,
				VesselID:                 &mediumBowl.ID,
				Name:                     "medium mixing bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "yogurt marinade",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "bowl with marinade",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Coat chicken with marinade
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   coatPrep.ID,
		Index:           4,
		Notes:           "Add chicken thigh pieces to bowl and, using clean hands, toss with marinade until evenly coated.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step4ID,
				ValidIngredientPreparationID:     &coatChickenVIP.ID,
				ValidIngredientMeasurementUnitID: &chickenThighsPoundVIMU.ID,
				IngredientID:                     &chickenThighs.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "boneless, skinless chicken thighs, cut into bite-size 1-inch pieces",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 2},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &yogurt.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "yogurt marinade",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 0.5},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step4ID,
				ValidPreparationInstrumentID: &coatBareHandsVPI.ID,
				InstrumentID:                 &bareHands.ID,
				Name:                         "bare hands",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &coatMediumBowlVPV.ID,
				VesselID:                        &mediumBowl.ID,
				Name:                            "bowl with marinade",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "marinated chicken",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](2)},
			},
		},
	}

	// Step 5: Transfer chicken to prepared baking sheet
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           5,
		Notes:           "Transfer chicken to prepared baking sheet, arranging pieces in a single, even layer in the aluminum-foil boat.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "marinated chicken",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &transferBakingSheetVPV.ID,
				VesselID:                        &bakingSheet.ID,
				Name:                            "prepared baking sheet with foil boat",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "chicken on baking sheet",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 6: Toast spices for sauce
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   toastPrep.ID,
		Index:           6,
		Notes:           "Add fenugreek leaves (or seeds, if using), chiles de arbol, brown cardamom (or green cardamom, if using), and clove to small skillet and place it over medium heat. Toast, tossing frequently, until spices are quite fragrant, about 1 to 2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60),
			Max: pointer.To[uint32](120),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &toastKasuriMethiVIP.ID,
				ValidIngredientMeasurementUnitID: &kasuriMethiTablespoonVIMU.ID,
				IngredientID:                     &kasuriMethi.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "kasuri methi (fenugreek leaves)",
				QuantityNotes:                    "1 tablespoon plus 2 teaspoons",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1.67},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &toastChilesVIP.ID,
				ValidIngredientMeasurementUnitID: &chilesDeArbolUnitVIMU.ID,
				IngredientID:                     &chilesDeArbol.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "whole dried chiles de arbol",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 2},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &toastBrownCardamomVIP.ID,
				ValidIngredientMeasurementUnitID: &brownCardamomUnitVIMU.ID,
				IngredientID:                     &brownCardamom.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "brown cardamom pod",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &toastWholeCloveVIP.ID,
				ValidIngredientMeasurementUnitID: &wholeClovesUnitVIMU.ID,
				IngredientID:                     &wholeCloves.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "whole clove",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step6ID,
				ValidPreparationVesselID: &toastSmallSkilletVPV.ID,
				VesselID:                 &smallSkillet.ID,
				Name:                     "small skillet",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "toasted spices for sauce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 7: Grind spices for sauce
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   grindPrep.ID,
		Index:           7,
		Notes:           "Transfer spices to spice grinder or mortar and pestle along with garam masala and salt and grind to a fine powder. Set aside.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &grindChilesVIP.ID,
				IngredientID:                    &chilesDeArbol.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "toasted spices for sauce",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ValidIngredientPreparationID:     &grindGaramMasalaVIP.ID,
				ValidIngredientMeasurementUnitID: &garamMasalaTablespoonVIMU.ID,
				IngredientID:                     &garamMasala.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "garam masala",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ValidIngredientPreparationID:     &grindSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "Diamond Crystal kosher salt",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step7ID,
				ValidPreparationInstrumentID: &grindSpiceGrinderVPI.ID,
				InstrumentID:                 &spiceGrinder.ID,
				Name:                         "spice grinder",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "ground spice mixture for sauce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 8: Soak cashews with water
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   soakPrep.ID,
		Index:           8,
		Notes:           "In a small, microwave-safe bowl, combine cashew nuts and 1/4 cup water.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &soakCashewsVIP.ID,
				ValidIngredientMeasurementUnitID: &cashewsOunceVIMU.ID,
				IngredientID:                     &cashews.ID,
				MeasurementUnitID:                ounceMeasurement.ID,
				Name:                             "raw cashews",
				QuantityNotes:                    "about 12 to 15",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step8ID,
				ValidIngredientPreparationID:     &soakWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				IngredientID:                     &water.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "water",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 0.25},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step8ID,
				ValidPreparationVesselID: &soakMicrowaveBowlVPV.ID,
				VesselID:                 &microwaveSafeBowl.ID,
				Name:                     "microwave-safe bowl",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "cashews in water",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "bowl with cashews",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 9: Microwave cashews
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   microwavePrep.ID,
		Index:           9,
		Notes:           "Microwave on high until cashews look plump and have softened slightly, about 1 minute. Set aside.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &microwaveCashewsVIP.ID,
				IngredientID:                    &cashews.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "cashews in water",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &microwaveMicrowaveBowlVPV.ID,
				VesselID:                        &microwaveSafeBowl.ID,
				Name:                            "bowl with cashews",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "softened cashews with soaking liquid",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 10: Heat oil in Dutch oven
	step10ID := identifiers.New()
	shimmeringState := enums.IngredientStates["shimmering"]
	step10OilIngredientID := identifiers.New()
	step10CompletionConditionID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   heatPrep.ID,
		Index:           10,
		Notes:           "In a Dutch oven, heat canola oil over medium-high heat until shimmering.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step10OilIngredientID,
				BelongsToRecipeStep:              step10ID,
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &canolaOilTablespoonVIMU.ID,
				IngredientID:                     &canolaOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "canola oil or other neutral-flavored oil",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: &heatDutchOvenVPV.ID,
				VesselID:                 &dutchOven.ID,
				Name:                     "Dutch oven",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "Dutch oven with hot oil",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step10CompletionConditionID,
				BelongsToRecipeStep: step10ID,
				IngredientStateID:   shimmeringState.ID,
				Notes:               "Oil should shimmer when viewed",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step10CompletionConditionID,
						RecipeStepIngredient:                   step10OilIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 11: Cook onions with baking soda
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cookPrep.ID,
		Index:           11,
		Notes:           "Add onions and baking soda and, using a wooden spoon, stir to coat onions in oil and distribute baking soda. Cook, stirring occasionally, until onions have completely broken down, most of their moisture has cooked off, and they begin to brown, 14 to 17 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](840),
			Max: pointer.To[uint32](1020),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &cookOnionVIP.ID,
				ValidIngredientMeasurementUnitID: &whiteOnionUnitVIMU.ID,
				IngredientID:                     &whiteOnion.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "medium white onion, peeled and cut into 1/2-inch dice",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step11ID,
				ValidIngredientPreparationID:     &cookBakingSodaVIP.ID,
				ValidIngredientMeasurementUnitID: &bakingSodaTeaspoonVIMU.ID,
				IngredientID:                     &bakingSoda.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "baking soda",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 0.25},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step11ID,
				ValidPreparationInstrumentID: &cookWoodenSpoonVPI.ID,
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &cookDutchOvenVPV.ID,
				VesselID:                        &dutchOven.ID,
				Name:                            "Dutch oven with hot oil",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "browned onions",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "Dutch oven with browned onions",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 12: Cook ginger and garlic
	step12ID := identifiers.New()
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cookPrep.ID,
		Index:           12,
		Notes:           "Reduce heat to medium low. Add ginger and garlic to pot and cook, stirring constantly, until quite fragrant, about 1 minute.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step12ID,
				ValidIngredientPreparationID:     &cookGingerVIP.ID,
				ValidIngredientMeasurementUnitID: &gingerUnitVIMU.ID,
				IngredientID:                     &ginger.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "fresh ginger, peeled and thinly sliced",
				QuantityNotes:                    "1-inch piece",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step12ID,
				ValidIngredientPreparationID:     &cookGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                cloveMeasurement.ID,
				Name:                             "garlic, smashed and roughly chopped",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 4},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step12ID,
				ValidPreparationInstrumentID: &cookWoodenSpoonVPI.ID,
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step12ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &cookDutchOvenVPV.ID,
				VesselID:                        &dutchOven.ID,
				Name:                            "Dutch oven with browned onions",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "cooked aromatics",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "Dutch oven with aromatics",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 13: Add spice mixture to onions
	step13ID := identifiers.New()
	step13 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step13ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           13,
		Notes:           "Using a wooden spoon, push onions into center of pot to form a mound. Add ground spice mixture to the mounded onions to prevent spices from scorching. Cook, stirring constantly, until onions are coated in spices and mixture is very fragrant, about 30 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](30),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				IngredientID:                    &garamMasala.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "ground spice mixture for sauce",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step13ID,
				ValidPreparationInstrumentID: &addWoodenSpoonVPI.ID,
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &addDutchOvenVPV.ID,
				VesselID:                        &dutchOven.ID,
				Name:                            "Dutch oven with aromatics",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "spiced onion mixture",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "Dutch oven with spiced onions",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 14: Add cashews, tomatoes, and water
	step14ID := identifiers.New()
	step14 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step14ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           14,
		Notes:           "Add cashews and their soaking liquid, scraping up any bits stuck to the bottom of the pot. Add tomatoes and their juices plus 1 cup water and, using the back of wooden spoon, crush tomatoes.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addCashewsVIP.ID,
				IngredientID:                    &cashews.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "softened cashews with soaking liquid",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step14ID,
				ValidIngredientPreparationID:     &addTomatoesVIP.ID,
				ValidIngredientMeasurementUnitID: &cannedTomatoesOunceVIMU.ID,
				IngredientID:                     &cannedTomatoes.ID,
				MeasurementUnitID:                ounceMeasurement.ID,
				Name:                             "whole fire-roasted canned tomatoes with their juices",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 28},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step14ID,
				ValidIngredientPreparationID:     &addWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				IngredientID:                     &water.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "water",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step14ID,
				ValidPreparationInstrumentID: &addWoodenSpoonVPI.ID,
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step14ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &addDutchOvenVPV.ID,
				VesselID:                        &dutchOven.ID,
				Name:                            "Dutch oven with spiced onions",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step14ID,
				Name:                "sauce base",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step14ID,
				Name:                "Dutch oven with sauce base",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 15: Simmer sauce
	step15ID := identifiers.New()
	step15 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step15ID,
		BelongsToRecipe: recipeID,
		PreparationID:   simmerPrep.ID,
		Index:           15,
		Notes:           "Bring to a boil, then reduce heat to maintain gentle simmer. Cook, stirring occasionally, until tomatoes are completely broken down and liquid has reduced, about 40 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](2400),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step15ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &simmerTomatoesVIP.ID,
				IngredientID:                    &cannedTomatoes.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "sauce base",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step15ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &simmerDutchOvenVPV.ID,
				VesselID:                        &dutchOven.ID,
				Name:                            "Dutch oven with sauce base",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step15ID,
				Name:                "simmered sauce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step15ID,
				Name:                "Dutch oven with simmered sauce",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 16: Preheat broiler
	step16ID := identifiers.New()
	step16 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step16ID,
		BelongsToRecipe: recipeID,
		PreparationID:   preheatPrep.ID,
		Index:           16,
		Notes:           "Meanwhile, adjust oven rack to about 3 inches below broiler element and preheat broiler on high.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step16ID,
				ValidPreparationVesselID: &preheatOvenVPV.ID,
				VesselID:                 &oven.ID,
				Name:                     "oven",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step16ID,
				Name:                "preheated broiler",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 17: Broil chicken
	step17ID := identifiers.New()
	step17 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step17ID,
		BelongsToRecipe: recipeID,
		PreparationID:   broilPrep.ID,
		Index:           17,
		Notes:           "Transfer chicken to broiler. Cook, checking the chicken frequently to ensure it's not burning, until chicken is charred in spots and is fully cooked through, about 14 minutes. Remove chicken from broiler and set aside.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](840),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step17ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &broilChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "chicken on baking sheet",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step17ID,
				ValidPreparationVesselID: &broilBakingSheetVPV.ID,
				VesselID:                 &bakingSheet.ID,
				Name:                     "baking sheet with chicken",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step17ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](16),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &broilOvenVPV.ID,
				VesselID:                        &oven.ID,
				Name:                            "preheated broiler",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step17ID,
				Name:                "broiled chicken",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](2)},
			},
		},
	}

	// Step 18: Blend sauce
	step18ID := identifiers.New()
	step18 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step18ID,
		BelongsToRecipe: recipeID,
		PreparationID:   blendPrep.ID,
		Index:           18,
		Notes:           "Using an immersion blender and off the heat, blend contents of Dutch oven until completely smooth, about 2 minutes. Alternatively, transfer contents of pot to blender and blend until completely smooth.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step18ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](15),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &blendTomatoesVIP.ID,
				IngredientID:                    &cannedTomatoes.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "simmered sauce",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step18ID,
				ValidPreparationInstrumentID: &blendStickBlenderVPI.ID,
				InstrumentID:                 &stickBlender.ID,
				Name:                         "immersion blender",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step18ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](15),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &blendDutchOvenVPV.ID,
				VesselID:                        &dutchOven.ID,
				Name:                            "Dutch oven with simmered sauce",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step18ID,
				Name:                "blended sauce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step18ID,
				Name:                "Dutch oven with blended sauce",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 19: Add butter and cream
	step19ID := identifiers.New()
	step19 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step19ID,
		BelongsToRecipe: recipeID,
		PreparationID:   blendPrep.ID,
		Index:           19,
		Notes:           "Add butter and cream, and blend until completely smooth and emulsified, about 2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step19ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](18),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &blendTomatoesVIP.ID,
				IngredientID:                    &cannedTomatoes.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "blended sauce",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step19ID,
				ValidIngredientPreparationID:     &blendButterVIP.ID,
				ValidIngredientMeasurementUnitID: &butterTablespoonVIMU.ID,
				IngredientID:                     &butter.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "unsalted butter, cut into 4 pieces",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 4},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step19ID,
				ValidIngredientPreparationID:     &blendCreamVIP.ID,
				ValidIngredientMeasurementUnitID: &heavyCreamCupVIMU.ID,
				IngredientID:                     &heavyCream.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "heavy cream",
				Quantity:                         types.Float32RangeWithOptionalMax{Min: 0.5},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step19ID,
				ValidPreparationInstrumentID: &blendStickBlenderVPI.ID,
				InstrumentID:                 &stickBlender.ID,
				Name:                         "immersion blender",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step19ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](18),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &blendDutchOvenVPV.ID,
				VesselID:                        &dutchOven.ID,
				Name:                            "Dutch oven with blended sauce",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step19ID,
				Name:                "makhani sauce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step19ID,
				Name:                "Dutch oven with makhani sauce",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 20: Add chicken to sauce
	step20ID := identifiers.New()
	step20 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step20ID,
		BelongsToRecipe: recipeID,
		PreparationID:   addPrep.ID,
		Index:           20,
		Notes:           "Add reserved broiled chicken along with any juices in the sheet pan to sauce and stir until chicken is well incorporated and warmed through.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step20ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](17),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "broiled chicken",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 2},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step20ID,
				ValidPreparationInstrumentID: &addWoodenSpoonVPI.ID,
				InstrumentID:                 &woodenSpoon.ID,
				Name:                         "wooden spoon",
				Quantity:                     types.Uint32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step20ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](19),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &addDutchOvenVPV.ID,
				VesselID:                        &dutchOven.ID,
				Name:                            "Dutch oven with makhani sauce",
				Quantity:                        types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step20ID,
				Name:                "butter chicken",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step20ID,
				Name:                "Dutch oven with butter chicken",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Step 21: Transfer to serving bowl and serve
	step21ID := identifiers.New()
	step21 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step21ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           21,
		Notes:           "Ladle chicken and sauce into serving bowl and drizzle with additional heavy cream. Serve immediately with rice alongside.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step21ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](20),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferChickenVIP.ID,
				IngredientID:                    &chickenThighs.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "butter chicken",
				Quantity:                        types.Float32RangeWithOptionalMax{Min: 1},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step21ID,
				ValidPreparationVesselID: &transferServingBowlVPV.ID,
				VesselID:                 &servingBowl.ID,
				Name:                     "serving bowl",
				Quantity:                 types.Uint16RangeWithOptionalMax{Min: 1},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step21ID,
				Name:                "finished butter chicken",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &unitMeasurement.ID,
				Quantity:            types.OptionalFloat32Range{Min: pointer.To[float32](1)},
			},
		},
	}

	// Create the main recipe
	butterChickenRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "Butter Chicken",
		Slug:                "butter-chicken",
		Source:              "https://www.seriouseats.com/stovetop-butter-chicken",
		Description:         "The sauce for this butter chicken is simmered on the stovetop and takes on a rich tomato flavor, while broiled marinated chicken adds char to the final dish.",
		YieldsComponentType: mealplanning.MealComponentTypesMain,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 4,
			Max: pointer.To[float32](6),
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step1, step2, step3, step4, step5, step6, step7, step8, step9,
			step10, step11, step12, step13, step14, step15, step16, step17, step18, step19,
			step20, step21,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		butterChickenRecipe,
	}
}
