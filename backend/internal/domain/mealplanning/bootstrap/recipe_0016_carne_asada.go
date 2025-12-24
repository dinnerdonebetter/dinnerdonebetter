package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// CarneAsadaRecipe creates the Best Carne Asada recipe.
// Source: https://www.seriouseats.com/carne-asada-food-lab-recipe-kenji
func CarneAsadaRecipe(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	recipeID := identifiers.New()

	// Get preparations
	microwavePrep := enums.Preparations["microwave"]
	transferPrep := enums.Preparations["transfer"]
	addPrep := enums.Preparations["add"]
	blendPrep := enums.Preparations["blend"]
	seasonPrep := enums.Preparations["season"]
	refrigeratePrep := enums.Preparations["refrigerate"]
	removePrep := enums.Preparations["remove"]
	lightPrep := enums.Preparations["light"]
	pourPrep := enums.Preparations["pour"]
	setPrep := enums.Preparations["set"]
	coverPrep := enums.Preparations["cover"]
	preheatPrep := enums.Preparations["preheat"]
	cleanPrep := enums.Preparations["clean"]
	oilPrep := enums.Preparations["oil"]
	grillPrep := enums.Preparations["grill"]
	restPrep := enums.Preparations["rest"]
	slicePrep := enums.Preparations["slice"]
	toastPrep := enums.Preparations["toast"]
	grindPrep := enums.Preparations["grind"]
	dividePrep := enums.Preparations["divide"]
	marinatePrep := enums.Preparations["marinate"]
	removeAirPrep := enums.Preparations["remove air"]
	unrefrigeratePrep := enums.Preparations["unrefrigerate"]

	// Get ingredients
	anchoChiles := enums.Ingredients["dried ancho chile"]
	guajilloChiles := enums.Ingredients["dried guajillo chile"]
	chipotlePeppers := enums.Ingredients["chipotle peppers in adobo"]
	orangeJuice := enums.Ingredients["orange juice"]
	limeJuice := enums.Ingredients["lime juice"]
	oliveOil := enums.Ingredients["olive oil"]
	soySauce := enums.Ingredients["soy sauce"]
	fishSauce := enums.Ingredients["Asian fish sauce"]
	darkBrownSugar := enums.Ingredients["dark brown sugar"]
	cilantro := enums.Ingredients["cilantro"]
	garlic := enums.Ingredients["garlic"]
	cuminSeeds := enums.Ingredients["cumin seeds"]
	corianderSeeds := enums.Ingredients["coriander seeds"]
	salt := enums.Ingredients["salt"]
	skirtSteak := enums.Ingredients["skirt steak"]

	// Get measurement units
	unitMeasurement := enums.MeasurementUnits["unit"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	poundMeasurement := enums.MeasurementUnits["pound"]
	cloveMeasurement := enums.MeasurementUnits["clove"]

	// Get instruments
	blender := enums.Instruments["blender"]
	chimneyStarterInstrument := enums.Instruments["chimney starter"]
	grillBrush := enums.Instruments["grill brush"]
	tongs := enums.Instruments["tongs"]
	thermometer := enums.Instruments["instant-read thermometer"]
	carvingKnife := enums.Instruments["carving knife"]

	// Get vessels
	microwaveSafePlate := enums.Vessels["microwave-safe plate"]
	blenderJar := enums.Vessels["blender jar"]
	largeBowl := enums.Vessels["large bowl"]
	sealedContainer := enums.Vessels["sealed container"]
	zipperLockBag := enums.Vessels["zipper-lock bag"]
	refrigerator := enums.Vessels["refrigerator"]
	chimneyStarter := enums.Vessels["chimney starter"]
	grill := enums.Vessels["grill"]
	charcoalGrate := enums.Vessels["charcoal grate"]
	cookingGrate := enums.Vessels["cooking grate"]
	grillingGrate := enums.Vessels["grilling grate"]
	cuttingBoard := enums.Vessels["cutting board"]
	smallSkillet := enums.Vessels["small skillet"]

	// Get ingredient states for completion conditions
	atTemperatureState := enums.IngredientStates["at temperature"]
	brownedState := enums.IngredientStates["browned"]
	pliableState := enums.IngredientStates["pliable"]

	// Get bridge table entries
	// Microwave
	microwaveAnchoChilesVIP := enums.IngredientPreparations[microwavePrep.ID][anchoChiles.ID]
	microwaveGuajilloChilesVIP := enums.IngredientPreparations[microwavePrep.ID][guajilloChiles.ID]
	microwaveMicrowaveSafePlateVPV := enums.PreparationVessels[microwavePrep.ID][microwaveSafePlate.ID]

	// Transfer
	transferSkirtSteakVIP := enums.IngredientPreparations[transferPrep.ID][skirtSteak.ID]
	transferCuttingBoardVPV := enums.PreparationVessels[transferPrep.ID][cuttingBoard.ID]

	// Add
	addChipotlePeppersVIP := enums.IngredientPreparations[addPrep.ID][chipotlePeppers.ID]
	addOrangeJuiceVIP := enums.IngredientPreparations[addPrep.ID][orangeJuice.ID]
	addLimeJuiceVIP := enums.IngredientPreparations[addPrep.ID][limeJuice.ID]
	addOliveOilVIP := enums.IngredientPreparations[addPrep.ID][oliveOil.ID]
	addSoySauceVIP := enums.IngredientPreparations[addPrep.ID][soySauce.ID]
	addFishSauceVIP := enums.IngredientPreparations[addPrep.ID][fishSauce.ID]
	addDarkBrownSugarVIP := enums.IngredientPreparations[addPrep.ID][darkBrownSugar.ID]
	addCilantroVIP := enums.IngredientPreparations[addPrep.ID][cilantro.ID]
	addGarlicVIP := enums.IngredientPreparations[addPrep.ID][garlic.ID]
	addCuminSeedsVIP := enums.IngredientPreparations[addPrep.ID][cuminSeeds.ID]
	addCorianderSeedsVIP := enums.IngredientPreparations[addPrep.ID][corianderSeeds.ID]
	addSaltVIP := enums.IngredientPreparations[addPrep.ID][salt.ID]

	// Blend
	blendAnchoChilesVIP := enums.IngredientPreparations[blendPrep.ID][anchoChiles.ID]
	blendBlenderJarVPV := enums.PreparationVessels[blendPrep.ID][blenderJar.ID]
	blendBlenderVPI := enums.PreparationInstruments[blendPrep.ID][blender.ID]

	// Season
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonBlenderJarVPV := enums.PreparationVessels[seasonPrep.ID][blenderJar.ID]

	// Refrigerate
	refrigerateSealedContainerVPV := enums.PreparationVessels[refrigeratePrep.ID][sealedContainer.ID]
	refrigerateRefrigeratorVPV := enums.PreparationVessels[refrigeratePrep.ID][refrigerator.ID]
	refrigerateSkirtSteakVIP := enums.IngredientPreparations[refrigeratePrep.ID][skirtSteak.ID]
	refrigerateZipperLockBagVPV := enums.PreparationVessels[refrigeratePrep.ID][zipperLockBag.ID]

	// Remove
	removeSkirtSteakVIP := enums.IngredientPreparations[removePrep.ID][skirtSteak.ID]
	removeZipperLockBagVPV := enums.PreparationVessels[removePrep.ID][zipperLockBag.ID]

	// Light
	lightChimneyStarterVPV := enums.PreparationVessels[lightPrep.ID][chimneyStarter.ID]
	lightChimneyStarterVPI := enums.PreparationInstruments[lightPrep.ID][chimneyStarterInstrument.ID]

	// Pour
	pourChimneyStarterVPV := enums.PreparationVessels[pourPrep.ID][chimneyStarter.ID]
	pourCharcoalGrateVPV := enums.PreparationVessels[pourPrep.ID][charcoalGrate.ID]

	// Set
	setCookingGrateVPV := enums.PreparationVessels[setPrep.ID][cookingGrate.ID]
	setGrillVPV := enums.PreparationVessels[setPrep.ID][grill.ID]

	// Cover
	coverGrillVPV := enums.PreparationVessels[coverPrep.ID][grill.ID]

	// Preheat
	preheatGrillVPV := enums.PreparationVessels[preheatPrep.ID][grill.ID]

	// Clean
	cleanGrillingGrateVPV := enums.PreparationVessels[cleanPrep.ID][grillingGrate.ID]
	cleanGrillBrushVPI := enums.PreparationInstruments[cleanPrep.ID][grillBrush.ID]

	// Oil
	oilOliveOilVIP := enums.IngredientPreparations[oilPrep.ID][oliveOil.ID]
	oilGrillingGrateVPV := enums.PreparationVessels[oilPrep.ID][grillingGrate.ID]

	// Grill
	grillSkirtSteakVIP := enums.IngredientPreparations[grillPrep.ID][skirtSteak.ID]
	grillGrillVPV := enums.PreparationVessels[grillPrep.ID][grill.ID]
	grillGrillingGrateVPV := enums.PreparationVessels[grillPrep.ID][grillingGrate.ID]
	grillTongsVPI := enums.PreparationInstruments[grillPrep.ID][tongs.ID]
	grillThermometerVPI := enums.PreparationInstruments[grillPrep.ID][thermometer.ID]

	// Rest
	restSkirtSteakVIP := enums.IngredientPreparations[restPrep.ID][skirtSteak.ID]
	restCuttingBoardVPV := enums.PreparationVessels[restPrep.ID][cuttingBoard.ID]

	// Slice
	sliceSkirtSteakVIP := enums.IngredientPreparations[slicePrep.ID][skirtSteak.ID]
	sliceCarvingKnifeVPI := enums.PreparationInstruments[slicePrep.ID][carvingKnife.ID]

	// Toast
	toastCuminSeedsVIP := enums.IngredientPreparations[toastPrep.ID][cuminSeeds.ID]
	toastCorianderSeedsVIP := enums.IngredientPreparations[toastPrep.ID][corianderSeeds.ID]
	toastSmallSkilletVPV := enums.PreparationVessels[toastPrep.ID][smallSkillet.ID]

	// Grind
	grindCuminSeedsVIP := enums.IngredientPreparations[grindPrep.ID][cuminSeeds.ID]
	grindBlenderJarVPV := enums.PreparationVessels[grindPrep.ID][blenderJar.ID]
	grindBlenderVPI := enums.PreparationInstruments[grindPrep.ID][blender.ID]

	// Divide
	divideMarinadeVIP := enums.IngredientPreparations[dividePrep.ID][anchoChiles.ID] // Using ancho chiles as proxy for marinade
	divideLargeBowlVPV := enums.PreparationVessels[dividePrep.ID][largeBowl.ID]
	divideSealedContainerVPV := enums.PreparationVessels[dividePrep.ID][sealedContainer.ID]

	// Marinate
	marinateSkirtSteakVIP := enums.IngredientPreparations[marinatePrep.ID][skirtSteak.ID]
	marinateLargeBowlVPV := enums.PreparationVessels[marinatePrep.ID][largeBowl.ID]
	marinateZipperLockBagVPV := enums.PreparationVessels[marinatePrep.ID][zipperLockBag.ID]

	// Remove air
	removeAirZipperLockBagVPV := enums.PreparationVessels[removeAirPrep.ID][zipperLockBag.ID]

	// Unrefrigerate
	unrefrigerateSealedContainerVPV := enums.PreparationVessels[unrefrigeratePrep.ID][sealedContainer.ID]
	unrefrigerateRefrigeratorVPV := enums.PreparationVessels[unrefrigeratePrep.ID][refrigerator.ID]

	// Measurement unit bridges
	anchoChilesUnitVIMU := enums.IngredientMeasurementUnits[anchoChiles.ID][unitMeasurement.ID]
	guajilloChilesUnitVIMU := enums.IngredientMeasurementUnits[guajilloChiles.ID][unitMeasurement.ID]
	chipotlePeppersUnitVIMU := enums.IngredientMeasurementUnits[chipotlePeppers.ID][unitMeasurement.ID]
	orangeJuiceCupVIMU := enums.IngredientMeasurementUnits[orangeJuice.ID][cupMeasurement.ID]
	limeJuiceTablespoonVIMU := enums.IngredientMeasurementUnits[limeJuice.ID][tablespoonMeasurement.ID]
	oliveOilTablespoonVIMU := enums.IngredientMeasurementUnits[oliveOil.ID][tablespoonMeasurement.ID]
	soySauceTablespoonVIMU := enums.IngredientMeasurementUnits[soySauce.ID][tablespoonMeasurement.ID]
	fishSauceTablespoonVIMU := enums.IngredientMeasurementUnits[fishSauce.ID][tablespoonMeasurement.ID]
	darkBrownSugarTablespoonVIMU := enums.IngredientMeasurementUnits[darkBrownSugar.ID][tablespoonMeasurement.ID]
	cilantroUnitVIMU := enums.IngredientMeasurementUnits[cilantro.ID][unitMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
	cuminSeedsTablespoonVIMU := enums.IngredientMeasurementUnits[cuminSeeds.ID][tablespoonMeasurement.ID]
	corianderSeedsTeaspoonVIMU := enums.IngredientMeasurementUnits[corianderSeeds.ID][teaspoonMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	skirtSteakPoundVIMU := enums.IngredientMeasurementUnits[skirtSteak.ID][poundMeasurement.ID]

	// Step 0: Toast cumin and coriander seeds until fragrant (optional prep step - can be done ahead)
	step0ID := identifiers.New()
	step0CuminSeedsIngredientID := identifiers.New()
	step0CorianderSeedsIngredientID := identifiers.New()
	step0 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0ID,
		BelongsToRecipe: recipeID,
		PreparationID:   toastPrep.ID,
		Index:           0,
		Notes:           "Toast cumin and coriander seeds until fragrant. This step is optional and can be done ahead of time.",
		Optional:        true,
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step0CuminSeedsIngredientID,
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &toastCuminSeedsVIP.ID,
				ValidIngredientMeasurementUnitID: &cuminSeedsTablespoonVIMU.ID,
				IngredientID:                     &cuminSeeds.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "cumin seeds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               step0CorianderSeedsIngredientID,
				BelongsToRecipeStep:              step0ID,
				ValidIngredientPreparationID:     &toastCorianderSeedsVIP.ID,
				ValidIngredientMeasurementUnitID: &corianderSeedsTeaspoonVIMU.ID,
				IngredientID:                     &corianderSeeds.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "coriander seeds",
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
				Name:                "toasted cumin and coriander seeds",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 0b: Grind toasted cumin and coriander seeds (optional prep step - can be done ahead)
	step0bID := identifiers.New()
	step0b := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step0bID,
		BelongsToRecipe: recipeID,
		PreparationID:   grindPrep.ID,
		Index:           1,
		Notes:           "Grind toasted cumin and coriander seeds. This step is optional and can be done ahead of time.",
		Optional:        true,
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step0bID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &grindCuminSeedsVIP.ID,
				IngredientID:                    &cuminSeeds.ID,
				MeasurementUnitID:               tablespoonMeasurement.ID,
				Name:                            "toasted cumin and coriander seeds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step0bID,
				ValidPreparationVesselID: &grindBlenderJarVPV.ID,
				VesselID:                 &blenderJar.ID,
				Name:                     "blender jar",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step0bID,
				ValidPreparationInstrumentID: &grindBlenderVPI.ID,
				InstrumentID:                 &blender.ID,
				Name:                         "blender",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step0bID,
				Name:                "toasted and ground cumin and coriander seeds",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 1: Microwave dried ancho and guajillo chiles on a microwave-safe plate until pliable and toasty-smelling, 10 to 20 seconds
	step1ID := identifiers.New()
	step1ChilesIngredientID := identifiers.New()
	step1CompletionConditionID := identifiers.New()
	step1 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step1ID,
		BelongsToRecipe: recipeID,
		PreparationID:   microwavePrep.ID,
		Index:           2,
		Notes:           "Place dried ancho and guajillo chiles on a microwave-safe plate and microwave until pliable and toasty-smelling, 10 to 20 seconds.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](10),
			Max: pointer.To[uint32](20),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               step1ChilesIngredientID,
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &microwaveAnchoChilesVIP.ID,
				ValidIngredientMeasurementUnitID: &anchoChilesUnitVIMU.ID,
				IngredientID:                     &anchoChiles.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "dried ancho chiles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step1ID,
				ValidIngredientPreparationID:     &microwaveGuajilloChilesVIP.ID,
				ValidIngredientMeasurementUnitID: &guajilloChilesUnitVIMU.ID,
				IngredientID:                     &guajilloChiles.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "dried guajillo chiles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step1ID,
				ValidPreparationVesselID: &microwaveMicrowaveSafePlateVPV.ID,
				VesselID:                 &microwaveSafePlate.ID,
				Name:                     "microwave-safe plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step1ID,
				Name:                "toasted and pliable chiles",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](6),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step1CompletionConditionID,
				BelongsToRecipeStep: step1ID,
				IngredientStateID:   pliableState.ID,
				Notes:               "Chiles should be pliable and toasty-smelling",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step1CompletionConditionID,
						RecipeStepIngredient:                   step1ChilesIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 2: Transfer chiles to blender jar, add all marinade ingredients, and blend until a smooth sauce has formed, about 1 minute
	step2ID := identifiers.New()
	step2 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step2ID,
		BelongsToRecipe: recipeID,
		PreparationID:   blendPrep.ID,
		Index:           3,
		Notes:           "Transfer chiles to blender jar, add chipotle peppers, orange juice, lime juice, olive oil, soy sauce, fish sauce, brown sugar, cilantro, garlic, cumin seed, and coriander seed, then blend until a smooth sauce has formed, about 1 minute.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](60),
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &blendAnchoChilesVIP.ID,
				IngredientID:                    &anchoChiles.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "toasted and pliable ancho chiles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &blendAnchoChilesVIP.ID,
				IngredientID:                    &guajilloChiles.ID,
				MeasurementUnitID:               unitMeasurement.ID,
				Name:                            "toasted and pliable guajillo chiles",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addChipotlePeppersVIP.ID,
				ValidIngredientMeasurementUnitID: &chipotlePeppersUnitVIMU.ID,
				IngredientID:                     &chipotlePeppers.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "chipotle peppers in adobo",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addOrangeJuiceVIP.ID,
				ValidIngredientMeasurementUnitID: &orangeJuiceCupVIMU.ID,
				IngredientID:                     &orangeJuice.ID,
				MeasurementUnitID:                cupMeasurement.ID,
				Name:                             "fresh orange juice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addLimeJuiceVIP.ID,
				ValidIngredientMeasurementUnitID: &limeJuiceTablespoonVIMU.ID,
				IngredientID:                     &limeJuice.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "fresh lime juice",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addOliveOilVIP.ID,
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
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addSoySauceVIP.ID,
				ValidIngredientMeasurementUnitID: &soySauceTablespoonVIMU.ID,
				IngredientID:                     &soySauce.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "soy sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addFishSauceVIP.ID,
				ValidIngredientMeasurementUnitID: &fishSauceTablespoonVIMU.ID,
				IngredientID:                     &fishSauce.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "Asian fish sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addDarkBrownSugarVIP.ID,
				ValidIngredientMeasurementUnitID: &darkBrownSugarTablespoonVIMU.ID,
				IngredientID:                     &darkBrownSugar.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "dark brown sugar",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addCilantroVIP.ID,
				ValidIngredientMeasurementUnitID: &cilantroUnitVIMU.ID,
				IngredientID:                     &cilantro.ID,
				MeasurementUnitID:                unitMeasurement.ID,
				Name:                             "cilantro, leaves and tender stems only",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step2ID,
				ValidIngredientPreparationID:     &addGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				IngredientID:                     &garlic.ID,
				MeasurementUnitID:                cloveMeasurement.ID,
				Name:                             "garlic cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addCuminSeedsVIP.ID,
				IngredientID:                    &cuminSeeds.ID,
				MeasurementUnitID:               tablespoonMeasurement.ID,
				Name:                            "toasted and ground cumin seeds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step2ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addCorianderSeedsVIP.ID,
				IngredientID:                    &corianderSeeds.ID,
				MeasurementUnitID:               teaspoonMeasurement.ID,
				Name:                            "toasted and ground coriander seeds",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step2ID,
				ValidPreparationVesselID: &blendBlenderJarVPV.ID,
				VesselID:                 &blenderJar.ID,
				Name:                     "blender jar",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step2ID,
				ValidPreparationInstrumentID: &blendBlenderVPI.ID,
				InstrumentID:                 &blender.ID,
				Name:                         "blender",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step2ID,
				Name:                "smooth marinade sauce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1.5),
				},
			},
		},
	}

	// Step 3: Season to taste with salt
	step3ID := identifiers.New()
	step3 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step3ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           4,
		Notes:           "Season to taste with salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step3ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonSaltVIP.ID,
				IngredientID:                    &salt.ID,
				MeasurementUnitID:               teaspoonMeasurement.ID,
				Name:                            "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step3ID,
				ValidPreparationVesselID: &seasonBlenderJarVPV.ID,
				VesselID:                 &blenderJar.ID,
				Name:                     "blender jar",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step3ID,
				Name:                "seasoned marinade sauce",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &cupMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1.5),
				},
			},
		},
	}

	// Step 4: Divide the marinade - transfer half to a large bowl and the other half to a sealed container
	step4ID := identifiers.New()
	step4 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step4ID,
		BelongsToRecipe: recipeID,
		PreparationID:   dividePrep.ID,
		Index:           5,
		Notes:           "Transfer half of the marinade to a large bowl and the other half to a sealed container.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step4ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &divideMarinadeVIP.ID,
				IngredientID:                    &anchoChiles.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "seasoned marinade sauce",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.75,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &divideLargeBowlVPV.ID,
				VesselID:                 &largeBowl.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step4ID,
				ValidPreparationVesselID: &divideSealedContainerVPV.ID,
				VesselID:                 &sealedContainer.ID,
				Name:                     "sealed container",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "marinade in large bowl",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step4ID,
				Name:                "marinade in sealed container",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               1,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Set aside the sealed container in the refrigerator
	step5ID := identifiers.New()
	step5 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step5ID,
		BelongsToRecipe: recipeID,
		PreparationID:   refrigeratePrep.ID,
		Index:           6,
		Notes:           "Set aside the sealed container in the refrigerator.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step5ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &refrigerateSealedContainerVPV.ID,
				VesselID:                        &sealedContainer.ID,
				Name:                            "sealed container with marinade",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step5ID,
				ValidPreparationVesselID: &refrigerateRefrigeratorVPV.ID,
				VesselID:                 &refrigerator.ID,
				Name:                     "refrigerator",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step5ID,
				Name:                "refrigerated marinade in sealed container",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 6: Season the marinade in the bowl with an extra 2 teaspoons of salt
	step6ID := identifiers.New()
	step6 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step6ID,
		BelongsToRecipe: recipeID,
		PreparationID:   seasonPrep.ID,
		Index:           7,
		Notes:           "Add an extra 2 teaspoons of salt to the marinade in the bowl. It should taste slightly saltier than is comfortable to taste.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step6ID,
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				IngredientID:                     &salt.ID,
				MeasurementUnitID:                teaspoonMeasurement.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step6ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &divideLargeBowlVPV.ID,
				VesselID:                        &largeBowl.ID,
				Name:                            "large bowl with marinade",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step6ID,
				Name:                "salted marinade in large bowl",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Marinate steak pieces one at a time in the bowl, turning to coat, then transfer to a gallon-sized zipper-lock bag
	step7ID := identifiers.New()
	step7 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step7ID,
		BelongsToRecipe: recipeID,
		PreparationID:   marinatePrep.ID,
		Index:           8,
		Notes:           "Add 1 piece of steak to bowl and turn to coat. Repeat with remaining steaks, adding them all to the same zipper-lock bag with the top folded over to prevent excess sauce and meat juices from contaminating the seal.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step7ID,
				ValidIngredientPreparationID:     &marinateSkirtSteakVIP.ID,
				ValidIngredientMeasurementUnitID: &skirtSteakPoundVIMU.ID,
				IngredientID:                     &skirtSteak.ID,
				MeasurementUnitID:                poundMeasurement.ID,
				Name:                             "skirt steak, trimmed and cut with the grain into 5- to 6-inch lengths",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step7ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &marinateLargeBowlVPV.ID,
				VesselID:                        &largeBowl.ID,
				Name:                            "large bowl with salted marinade",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step7ID,
				ValidPreparationVesselID: &marinateZipperLockBagVPV.ID,
				VesselID:                 &zipperLockBag.ID,
				Name:                     "gallon-sized zipper-lock bag",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step7ID,
				Name:                "steak in zipper-lock bag",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Pour any excess marinade over the steaks
	step8ID := identifiers.New()
	step8 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step8ID,
		BelongsToRecipe: recipeID,
		PreparationID:   pourPrep.ID,
		Index:           9,
		Notes:           "Pour any excess marinade over the steaks.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addSaltVIP.ID,
				IngredientID:                    &salt.ID,
				MeasurementUnitID:               cupMeasurement.ID,
				Name:                            "excess marinade",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step8ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &pourChimneyStarterVPV.ID,
				VesselID:                        &zipperLockBag.ID,
				Name:                            "zipper-lock bag with steak",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step8ID,
				Name:                "steak in bag with excess marinade",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Squeeze all air out of the bag and seal
	step9ID := identifiers.New()
	step9 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step9ID,
		BelongsToRecipe: recipeID,
		PreparationID:   removeAirPrep.ID,
		Index:           10,
		Notes:           "Squeeze all air out of the bag and seal.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step9ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &removeAirZipperLockBagVPV.ID,
				VesselID:                        &zipperLockBag.ID,
				Name:                            "zipper-lock bag with steak and marinade",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step9ID,
				Name:                "sealed zipper-lock bag with steak",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 10: Refrigerate for at least 3 hours or up to 12 hours
	step10ID := identifiers.New()
	step10 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step10ID,
		BelongsToRecipe: recipeID,
		PreparationID:   refrigeratePrep.ID,
		Index:           11,
		Notes:           "Refrigerate for at least 3 hours or up to 12 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](10800), // 3 hours
			Max: pointer.To[uint32](43200), // 12 hours
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &refrigerateSkirtSteakVIP.ID,
				IngredientID:                    &skirtSteak.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "steak in sealed bag",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step10ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](10),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &refrigerateZipperLockBagVPV.ID,
				VesselID:                        &zipperLockBag.ID,
				Name:                            "sealed zipper-lock bag with steak",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step10ID,
				ValidPreparationVesselID: &refrigerateRefrigeratorVPV.ID,
				VesselID:                 &refrigerator.ID,
				Name:                     "refrigerator",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step10ID,
				Name:                "refrigerated marinated steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 11: When ready to cook, remove the extra marinade from the fridge to allow it to warm up a little
	step11ID := identifiers.New()
	step11 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step11ID,
		BelongsToRecipe: recipeID,
		PreparationID:   unrefrigeratePrep.ID,
		Index:           12,
		Notes:           "When ready to cook, remove the extra marinade from the fridge to allow it to warm up a little.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &unrefrigerateSealedContainerVPV.ID,
				VesselID:                        &sealedContainer.ID,
				Name:                            "sealed container with marinade from refrigerator",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step11ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &unrefrigerateRefrigeratorVPV.ID,
				VesselID:                        &refrigerator.ID,
				Name:                            "refrigerator",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step11ID,
				Name:                "removed marinade container",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 12: Light one chimney full of charcoal
	step12ID := identifiers.New()
	step12 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step12ID,
		BelongsToRecipe: recipeID,
		PreparationID:   lightPrep.ID,
		Index:           13,
		Notes:           "Light one chimney full of charcoal. When all the charcoal is lit and covered with gray ash, pour out and arrange the coals on one side of the charcoal grate.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step12ID,
				ValidPreparationVesselID: &lightChimneyStarterVPV.ID,
				VesselID:                 &chimneyStarter.ID,
				Name:                     "chimney starter",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step12ID,
				ValidPreparationInstrumentID: &lightChimneyStarterVPI.ID,
				InstrumentID:                 &chimneyStarterInstrument.ID,
				Name:                         "chimney starter",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step12ID,
				Name:                "lit chimney starter with charcoal",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 13: Pour out and arrange coals on one side of the charcoal grate
	step13ID := identifiers.New()
	step13 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step13ID,
		BelongsToRecipe: recipeID,
		PreparationID:   pourPrep.ID,
		Index:           14,
		Notes:           "Pour out and arrange coals on one side of the charcoal grate.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step13ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &pourChimneyStarterVPV.ID,
				VesselID:                        &chimneyStarter.ID,
				Name:                            "lit chimney starter",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step13ID,
				ValidPreparationVesselID: &pourCharcoalGrateVPV.ID,
				VesselID:                 &charcoalGrate.ID,
				Name:                     "charcoal grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step13ID,
				Name:                "charcoal arranged on one side of grate",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 14: Set cooking grate in place
	step14ID := identifiers.New()
	step14 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step14ID,
		BelongsToRecipe: recipeID,
		PreparationID:   setPrep.ID,
		Index:           15,
		Notes:           "Set cooking grate in place. Alternatively, set half the burners on a gas grill to the highest heat setting.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step14ID,
				ValidPreparationVesselID: &setCookingGrateVPV.ID,
				VesselID:                 &cookingGrate.ID,
				Name:                     "cooking grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step14ID,
				ValidPreparationVesselID: &setGrillVPV.ID,
				VesselID:                 &grill.ID,
				Name:                     "grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step14ID,
				Name:                "cooking grate set in place on grill",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 15: Cover grill and preheat for 5 minutes. Alternatively, for gas grill, preheat for 10 minutes.
	step15ID := identifiers.New()
	step15 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step15ID,
		BelongsToRecipe: recipeID,
		PreparationID:   preheatPrep.ID,
		Index:           16,
		Notes:           "Cover grill and preheat for 5 minutes. Alternatively, for gas grill, preheat for 10 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
			Max: pointer.To[uint32](600), // 10 minutes
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step15ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](15),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &coverGrillVPV.ID,
				VesselID:                        &grill.ID,
				Name:                            "grill with cooking grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step15ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](15),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &preheatGrillVPV.ID,
				VesselID:                        &grill.ID,
				Name:                            "covered grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step15ID,
				Name:                "preheated grill",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 16: Clean and oil the grilling grate
	step16ID := identifiers.New()
	step16 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step16ID,
		BelongsToRecipe: recipeID,
		PreparationID:   cleanPrep.ID,
		Index:           17,
		Notes:           "Clean and oil the grilling grate.",
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step16ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](16),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &cleanGrillingGrateVPV.ID,
				VesselID:                        &grillingGrate.ID,
				Name:                            "grilling grate on preheated grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step16ID,
				ValidPreparationInstrumentID: &cleanGrillBrushVPI.ID,
				InstrumentID:                 &grillBrush.ID,
				Name:                         "grill brush",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step16ID,
				Name:                "cleaned grilling grate",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 17: Oil the grilling grate
	step17ID := identifiers.New()
	step17 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step17ID,
		BelongsToRecipe: recipeID,
		PreparationID:   oilPrep.ID,
		Index:           18,
		Notes:           "Oil the grilling grate.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                               identifiers.New(),
				BelongsToRecipeStep:              step17ID,
				ValidIngredientPreparationID:     &oilOliveOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oliveOilTablespoonVIMU.ID,
				IngredientID:                     &oliveOil.ID,
				MeasurementUnitID:                tablespoonMeasurement.ID,
				Name:                             "olive oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step17ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](17),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &oilGrillingGrateVPV.ID,
				VesselID:                        &grillingGrate.ID,
				Name:                            "cleaned grilling grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step17ID,
				Name:                "oiled grilling grate",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 18: Remove steaks from marinade and wipe off excess
	step18ID := identifiers.New()
	step18 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step18ID,
		BelongsToRecipe: recipeID,
		PreparationID:   removePrep.ID,
		Index:           19,
		Notes:           "Remove steaks from marinade and wipe off excess.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step18ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &removeSkirtSteakVIP.ID,
				IngredientID:                    &skirtSteak.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "refrigerated marinated steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step18ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &removeZipperLockBagVPV.ID,
				VesselID:                        &zipperLockBag.ID,
				Name:                            "zipper-lock bag with marinated steak",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step18ID,
				Name:                "wiped steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 19: Place steaks directly over the hot side of the grill and cook, turning occasionally, until well charred and center registers 110°F (43°C), 5 to 10 minutes total
	step19ID := identifiers.New()
	step19SteakIngredientID := identifiers.New()
	step19CompletionConditionID := identifiers.New()
	step19 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step19ID,
		BelongsToRecipe: recipeID,
		PreparationID:   grillPrep.ID,
		Index:           20,
		Notes:           "Place steaks directly over the hot side of the grill. If using a gas grill, cover; if using a charcoal grill, leave exposed. Cook, turning occasionally, until steaks are well charred on outside and center registers 110°F (43°C) on an instant-read thermometer, 5 to 10 minutes total.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
			Max: pointer.To[uint32](600), // 10 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              step19SteakIngredientID,
				BelongsToRecipeStep:             step19ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](19),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &grillSkirtSteakVIP.ID,
				IngredientID:                    &skirtSteak.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "wiped steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step19ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](18),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &grillGrillingGrateVPV.ID,
				VesselID:                        &grillingGrate.ID,
				Name:                            "oiled grilling grate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step19ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](16),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &grillGrillVPV.ID,
				VesselID:                        &grill.ID,
				Name:                            "preheated grill",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step19ID,
				ValidPreparationInstrumentID: &grillTongsVPI.ID,
				InstrumentID:                 &tongs.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step19ID,
				ValidPreparationInstrumentID: &grillThermometerVPI.ID,
				InstrumentID:                 &thermometer.ID,
				Name:                         "instant-read thermometer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step19ID,
				Name:                "grilled steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{
			{
				ID:                  step19CompletionConditionID,
				BelongsToRecipeStep: step19ID,
				IngredientStateID:   brownedState.ID,
				Notes:               "Steaks should be well charred on the outside",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step19CompletionConditionID,
						RecipeStepIngredient:                   step19SteakIngredientID,
					},
				},
				Optional: false,
			},
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step19ID,
				IngredientStateID:   atTemperatureState.ID,
				Notes:               "Center should register 110°F (43°C) on an instant-read thermometer",
				Ingredients: []*mealplanning.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
					{
						ID:                                     identifiers.New(),
						BelongsToRecipeStepCompletionCondition: step19CompletionConditionID,
						RecipeStepIngredient:                   step19SteakIngredientID,
					},
				},
				Optional: false,
			},
		},
	}

	// Step 20: Transfer to a cutting board and allow to rest for 5 minutes
	step20ID := identifiers.New()
	step20 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step20ID,
		BelongsToRecipe: recipeID,
		PreparationID:   transferPrep.ID,
		Index:           21,
		Notes:           "Transfer to a cutting board.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step20ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](20),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &transferSkirtSteakVIP.ID,
				IngredientID:                    &skirtSteak.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "grilled steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                       identifiers.New(),
				BelongsToRecipeStep:      step20ID,
				ValidPreparationVesselID: &transferCuttingBoardVPV.ID,
				VesselID:                 &cuttingBoard.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step20ID,
				Name:                "steak on cutting board",
				Type:                mealplanning.RecipeStepProductVesselType,
				Index:               0,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 21: Allow to rest for 5 minutes
	step21ID := identifiers.New()
	step21 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step21ID,
		BelongsToRecipe: recipeID,
		PreparationID:   restPrep.ID,
		Index:           22,
		Notes:           "Allow to rest for 5 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](300), // 5 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step21ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](21),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &restSkirtSteakVIP.ID,
				IngredientID:                    &skirtSteak.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "steak on cutting board",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step21ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](21),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &restCuttingBoardVPV.ID,
				VesselID:                        &cuttingBoard.ID,
				Name:                            "cutting board with steak",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step21ID,
				Name:                "rested steak",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 22: Slice thinly against the grain and serve immediately
	step22ID := identifiers.New()
	step22 := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:              step22ID,
		BelongsToRecipe: recipeID,
		PreparationID:   slicePrep.ID,
		Index:           23,
		Notes:           "Slice thinly against the grain and serve immediately, passing extra salsa, lime wedges, avocado, onions, cilantro, and tortillas on the side.",
		Ingredients: []*mealplanning.RecipeStepIngredientDatabaseCreationInput{
			{
				ID:                              identifiers.New(),
				BelongsToRecipeStep:             step22ID,
				ProductOfRecipeStepIndex:        pointer.To[uint64](22),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &sliceSkirtSteakVIP.ID,
				IngredientID:                    &skirtSteak.ID,
				MeasurementUnitID:               poundMeasurement.ID,
				Name:                            "rested steak",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{
			{
				ID:                           identifiers.New(),
				BelongsToRecipeStep:          step22ID,
				ValidPreparationInstrumentID: &sliceCarvingKnifeVPI.ID,
				InstrumentID:                 &carvingKnife.ID,
				Name:                         "carving knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductDatabaseCreationInput{
			{
				ID:                  identifiers.New(),
				BelongsToRecipeStep: step22ID,
				Name:                "sliced carne asada",
				Type:                mealplanning.RecipeStepProductIngredientType,
				Index:               0,
				MeasurementUnitID:   &poundMeasurement.ID,
				Quantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Create prep task for toasting and grinding spices ahead of time
	prepTask0ID := identifiers.New()
	prepTask0 := &mealplanning.RecipePrepTaskDatabaseCreationInput{
		ID:                          prepTask0ID,
		BelongsToRecipe:             recipeID,
		Name:                        "Toast and grind cumin and coriander seeds",
		Description:                 "The cumin and coriander seeds can be toasted and ground ahead of time. Store in an airtight container at room temperature.",
		Notes:                       "Toasting and grinding spices enhances their flavor. This step is optional but recommended.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the toasted and ground spices in an airtight container at room temperature.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeAirtightContainer,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](18), // Room temperature
			Max: pointer.To[float32](25),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 0,
			Max: pointer.To[uint32](604800), // Up to 7 days
		},
		TaskSteps: []*mealplanning.RecipePrepTaskStepDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: step0ID, BelongsToRecipePrepTask: prepTask0ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step0bID, BelongsToRecipePrepTask: prepTask0ID, SatisfiesRecipeStep: true},
		},
	}

	// Create prep task for marinating steak ahead of time
	prepTask1ID := identifiers.New()
	prepTask1 := &mealplanning.RecipePrepTaskDatabaseCreationInput{
		ID:                          prepTask1ID,
		BelongsToRecipe:             recipeID,
		Name:                        "Marinate skirt steak",
		Description:                 "The skirt steak can be marinated for at least 3 hours and up to 12 hours in advance. Store in the refrigerator in a zipper-lock bag.",
		Notes:                       "Marinating improves flavor and tenderness.",
		Optional:                    false,
		ExplicitStorageInstructions: "Store the marinated steak in a sealed zipper-lock bag in the refrigerator for at least 3 hours and up to 12 hours.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeAirtightContainer,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 10800,                     // 3 hours
			Max: pointer.To[uint32](43200), // 12 hours
		},
		TaskSteps: []*mealplanning.RecipePrepTaskStepDatabaseCreationInput{
			{ID: identifiers.New(), BelongsToRecipeStep: step6ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step7ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step8ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step9ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: false},
			{ID: identifiers.New(), BelongsToRecipeStep: step10ID, BelongsToRecipePrepTask: prepTask1ID, SatisfiesRecipeStep: true},
		},
	}

	carneAsadaRecipe := &mealplanning.RecipeDatabaseCreationInput{
		ID:                  recipeID,
		CreatedByUser:       userID,
		Name:                "The Best Carne Asada",
		Slug:                "the-best-carne-asada",
		Source:              "https://www.seriouseats.com/carne-asada-food-lab-recipe-kenji",
		Description:         "The best carne asada combines dried chiles, citrus, spices, and skirt steak, which is then grilled over ripping-hot heat. The marinade should have a good balance of flavors, with no single ingredient overwhelming any other.",
		YieldsComponentType: mealplanning.MealComponentTypesMain,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: 6,
		},
		PortionName:       "serving",
		PluralPortionName: "servings",
		EligibleForMeals:  true,
		Steps: []*mealplanning.RecipeStepDatabaseCreationInput{
			step0, step0b, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13, step14, step15, step16, step17, step18, step19, step20, step21, step22,
		},
		PrepTasks: []*mealplanning.RecipePrepTaskDatabaseCreationInput{
			prepTask0,
			prepTask1,
		},
	}

	return []*mealplanning.RecipeDatabaseCreationInput{
		carneAsadaRecipe,
	}
}
