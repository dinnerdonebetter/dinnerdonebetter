package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// RefriedBeansRecipe creates the Perfect Frijoles Refritos (Mexican Refried Beans) recipe.
// Source: https://www.seriouseats.com/perfect-refried-beans
func RefriedBeansRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	halvePrep := enums.Preparations["halve"]
	mincePrep := enums.Preparations["mince"]
	peelPrep := enums.Preparations["peel"]
	coverPrep := enums.Preparations["cover"]
	addPrep := enums.Preparations["add"]
	boilPrep := enums.Preparations["boil"]
	reducePrep := enums.Preparations["reduce"]
	simmerPrep := enums.Preparations["simmer"]
	seasonPrep := enums.Preparations["season"]
	drainPrep := enums.Preparations["drain"]
	measurePrep := enums.Preparations["measure"]
	discardPrep := enums.Preparations["discard"]
	heatPrep := enums.Preparations["heat"]
	sautPrep := enums.Preparations["sauté"]
	stirPrep := enums.Preparations["stir"]
	smashPrep := enums.Preparations["smash"]
	dilutePrep := enums.Preparations["dilute"]

	// Get ingredients
	pintoBeans := enums.Ingredients["pinto beans"]
	water := enums.Ingredients["water"]
	epazote := enums.Ingredients["epazote"]
	whiteOnion := enums.Ingredients["onion"]
	garlic := enums.Ingredients["garlic"]
	salt := enums.Ingredients["salt"]
	lard := enums.Ingredients["lard"]

	// Get measurement units
	poundMeasurement := enums.MeasurementUnits["pound"]
	cupMeasurement := enums.MeasurementUnits["cup"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]
	cloveMeasurement := enums.MeasurementUnits["clove"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	chefsKnife := enums.Instruments["knife"]
	bareHands := enums.Instruments["bare hands"]
	potatoMasher := enums.Instruments["potato masher"]

	// Get vessels
	cuttingBoard := enums.Vessels["cutting board"]
	largePot := enums.Vessels["pot"]
	largeSkillet := enums.Vessels["cast iron skillet"]
	largeBowl := enums.Vessels["large bowl"]

	// Get bridge table entries
	// Halve
	halveWhiteOnionVIP := enums.IngredientPreparations[halvePrep.ID][whiteOnion.ID]
	halveChefsKnifeVPI := enums.PreparationInstruments[halvePrep.ID][chefsKnife.ID]
	halveCuttingBoardVPV := enums.PreparationVessels[halvePrep.ID][cuttingBoard.ID]

	// Mince
	minceWhiteOnionVIP := enums.IngredientPreparations[mincePrep.ID][whiteOnion.ID]
	minceChefsKnifeVPI := enums.PreparationInstruments[mincePrep.ID][chefsKnife.ID]
	minceCuttingBoardVPV := enums.PreparationVessels[mincePrep.ID][cuttingBoard.ID]

	// Peel
	peelGarlicVIP := enums.IngredientPreparations[peelPrep.ID][garlic.ID]
	peelBareHandsVPI := enums.PreparationInstruments[peelPrep.ID][bareHands.ID]

	// Cover
	coverPintoBeansVIP := enums.IngredientPreparations[coverPrep.ID][pintoBeans.ID]
	coverWaterVIP := enums.IngredientPreparations[coverPrep.ID][water.ID]
	coverLargePotVPV := enums.PreparationVessels[coverPrep.ID][largePot.ID]

	// Add
	addEpazoteVIP := enums.IngredientPreparations[addPrep.ID][epazote.ID]
	addWhiteOnionVIP := enums.IngredientPreparations[addPrep.ID][whiteOnion.ID]
	addGarlicVIP := enums.IngredientPreparations[addPrep.ID][garlic.ID]
	addWaterVIP := enums.IngredientPreparations[addPrep.ID][water.ID]
	addLargePotVPV := enums.PreparationVessels[addPrep.ID][largePot.ID]

	// Boil
	boilLargePotVPV := enums.PreparationVessels[boilPrep.ID][largePot.ID]

	// Reduce
	reduceLargePotVPV := enums.PreparationVessels[reducePrep.ID][largePot.ID]

	// Simmer
	simmerPintoBeansVIP := enums.IngredientPreparations[simmerPrep.ID][pintoBeans.ID]
	simmerLargePotVPV := enums.PreparationVessels[simmerPrep.ID][largePot.ID]

	// Season
	seasonPintoBeansVIP := enums.IngredientPreparations[seasonPrep.ID][pintoBeans.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	seasonLargePotVPV := enums.PreparationVessels[seasonPrep.ID][largePot.ID]
	seasonLargeSkilletVPV := enums.PreparationVessels[seasonPrep.ID][largeSkillet.ID]

	// Drain
	drainPintoBeansVIP := enums.IngredientPreparations[drainPrep.ID][pintoBeans.ID]
	drainLargePotVPV := enums.PreparationVessels[drainPrep.ID][largePot.ID]
	drainLargeBowlVPV := enums.PreparationVessels[drainPrep.ID][largeBowl.ID]

	// Reserve (not used in this recipe, but bridge table entry exists)

	// Measure
	measurePintoBeansVIP := enums.IngredientPreparations[measurePrep.ID][pintoBeans.ID]
	measureLargeBowlVPV := enums.PreparationVessels[measurePrep.ID][largeBowl.ID]

	// Discard
	discardEpazoteVIP := enums.IngredientPreparations[discardPrep.ID][epazote.ID]
	discardWhiteOnionVIP := enums.IngredientPreparations[discardPrep.ID][whiteOnion.ID]
	discardGarlicVIP := enums.IngredientPreparations[discardPrep.ID][garlic.ID]
	discardLargePotVPV := enums.PreparationVessels[discardPrep.ID][largePot.ID]

	// Heat
	heatLardVIP := enums.IngredientPreparations[heatPrep.ID][lard.ID]
	heatLargeSkilletVPV := enums.PreparationVessels[heatPrep.ID][largeSkillet.ID]

	// Sauté
	sautWhiteOnionVIP := enums.IngredientPreparations[sautPrep.ID][whiteOnion.ID]
	sautLargeSkilletVPV := enums.PreparationVessels[sautPrep.ID][largeSkillet.ID]

	// Stir
	stirPintoBeansVIP := enums.IngredientPreparations[stirPrep.ID][pintoBeans.ID]
	stirWhiteOnionVIP := enums.IngredientPreparations[stirPrep.ID][whiteOnion.ID]
	stirLargeSkilletVPV := enums.PreparationVessels[stirPrep.ID][largeSkillet.ID]

	// Smash
	smashPintoBeansVIP := enums.IngredientPreparations[smashPrep.ID][pintoBeans.ID]
	smashLargeSkilletVPV := enums.PreparationVessels[smashPrep.ID][largeSkillet.ID]
	smashPotatoMasherVPI := enums.PreparationInstruments[smashPrep.ID][potatoMasher.ID]

	// Dilute
	dilutePintoBeansVIP := enums.IngredientPreparations[dilutePrep.ID][pintoBeans.ID]
	diluteWaterVIP := enums.IngredientPreparations[dilutePrep.ID][water.ID]
	diluteLargeSkilletVPV := enums.PreparationVessels[dilutePrep.ID][largeSkillet.ID]

	// Measurement unit bridges
	pintoBeansPoundVIMU := enums.IngredientMeasurementUnits[pintoBeans.ID][poundMeasurement.ID]
	waterCupVIMU := enums.IngredientMeasurementUnits[water.ID][cupMeasurement.ID]
	epazoteSprigVIMU := enums.IngredientMeasurementUnits[epazote.ID][sprigMeasurement.ID]
	whiteOnionUnitVIMU := enums.IngredientMeasurementUnits[whiteOnion.ID][unitMeasurement.ID]
	garlicCloveVIMU := enums.IngredientMeasurementUnits[garlic.ID][cloveMeasurement.ID]
	saltTeaspoonVIMU := enums.IngredientMeasurementUnits[salt.ID][teaspoonMeasurement.ID]
	lardTablespoonVIMU := enums.IngredientMeasurementUnits[lard.ID][tablespoonMeasurement.ID]

	// Step 0: Halve the onion
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        halvePrep.ID,
		Index:                0,
		ExplicitInstructions: "Halve the medium white onion.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &halveWhiteOnionVIP.ID,
				ValidIngredientMeasurementUnitID: &whiteOnionUnitVIMU.ID,
				Name:                             "medium white onion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &halveChefsKnifeVPI.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &halveCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "halved onion (2 halves)",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 1: In a large pot, cover the beans with cold water by at least 2 inches
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        coverPrep.ID,
		Index:                1,
		ExplicitInstructions: "In a large pot, cover the beans with cold water by at least 2 inches.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ValidIngredientPreparationID:     &coverPintoBeansVIP.ID,
				ValidIngredientMeasurementUnitID: &pintoBeansPoundVIMU.ID,
				Name:                             "dried pinto beans",
				QuantityNotes:                    "1/2 pound (227 g)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{

				ValidIngredientPreparationID:     &coverWaterVIP.ID,
				ValidIngredientMeasurementUnitID: &waterCupVIMU.ID,
				Name:                             "cold water",
				QuantityNotes:                    "Enough to cover beans by at least 2 inches",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4, // Approximate cups
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ValidPreparationVesselID: &coverLargePotVPV.ID,
				Name:                     "large pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:  "beans covered with water",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 2: Peel garlic cloves
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        peelPrep.ID,
		Index:                2,
		ExplicitInstructions: "Peel 2 medium cloves of garlic.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ValidIngredientPreparationID:     &peelGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				Name:                             "medium cloves garlic",
				QuantityNotes:                    "2 medium cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{

				ValidPreparationInstrumentID: &peelBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "peeled garlic cloves",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cloveMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](2),
				},
			},
		},
	}

	// Step 3: Add herb sprigs, the whole onion half, and peeled garlic cloves
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                3,
		ExplicitInstructions: "Add the herb sprigs, the whole onion half, and the peeled garlic cloves.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ValidIngredientPreparationID:     &addEpazoteVIP.ID,
				ValidIngredientMeasurementUnitID: &epazoteSprigVIMU.ID,
				Name:                             "fresh epazote",
				QuantityNotes:                    "2 sprigs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addWhiteOnionVIP.ID,
				Name:                            "onion half (left whole)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &addGarlicVIP.ID,
				Name:                            "peeled garlic cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &addLargePotVPV.ID,
				Name:                            "large pot with beans covered with water",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:  "beans with aromatics in pot",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Bring to a boil over high heat
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        boilPrep.ID,
		Index:                4,
		ExplicitInstructions: "Bring to a boil over high heat.",
		Ingredients:          []*mealplanning.RecipeStepIngredientCreationRequestInput{},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &boilLargePotVPV.ID,
				Name:                            "beans with aromatics in pot",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:  "boiling beans with aromatics",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Reduce heat to simmer
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        reducePrep.ID,
		Index:                5,
		ExplicitInstructions: "Reduce the heat to a simmer.",
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &reduceLargePotVPV.ID,
				Name:                            "boiling beans with aromatics",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:  "beans ready to simmer",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:              "beans ready to simmer",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             1,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 6: Simmer until beans are very tender, about 1 to 2 hours
	tenderState := enums.IngredientStates["tender"]
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        simmerPrep.ID,
		Index:                6,
		ExplicitInstructions: "Simmer until the beans are very tender, about 1 to 2 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](3600), // 1 hour
			Max: pointer.To[uint32](7200), // 2 hours
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &simmerPintoBeansVIP.ID,
				Name:                            "beans ready to simmer",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &simmerLargePotVPV.ID,
				Name:                            "beans ready to simmer",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "very tender cooked beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{

				Name:  "pot with cooked beans and cooking liquid",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: tenderState.ID,
				Notes:             "Beans should be very tender",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
	}

	// Step 7: Season with salt
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        seasonPrep.ID,
		Index:                7,
		ExplicitInstructions: "Season with salt.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonPintoBeansVIP.ID,
				Name:                            "very tender cooked beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{

				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &seasonLargePotVPV.ID,
				Name:                            "pot with cooked beans and cooking liquid",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "seasoned cooked beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 8: Drain beans, reserving bean-cooking liquid
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        drainPrep.ID,
		Index:                8,
		ExplicitInstructions: "Drain the beans, reserving the bean-cooking liquid.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &drainPintoBeansVIP.ID,
				Name:                            "seasoned cooked beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &drainLargePotVPV.ID,
				Name:                            "pot with cooked beans and cooking liquid",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{

				ValidPreparationVesselID: &drainLargeBowlVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "drained cooked beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
			{

				Name:              "reserved bean-cooking liquid",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             1,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 9: Measure out 3 cups of beans (if you have more, reserve the rest for another use)
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        measurePrep.ID,
		Index:                9,
		ExplicitInstructions: "You should have about 3 cups of cooked beans; if you have more, measure out 3 cups of beans and reserve the rest for another use.",
		Optional:             true,
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &measurePintoBeansVIP.ID,
				Name:                            "drained cooked beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ValidPreparationVesselID: &measureLargeBowlVPV.ID,
				Name:                     "large bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "3 cups of cooked beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 10: Discard herb sprigs, onion, and garlic
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        discardPrep.ID,
		Index:                10,
		ExplicitInstructions: "Discard the herb sprigs, onion, and garlic.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ValidIngredientPreparationID:     &discardEpazoteVIP.ID,
				ValidIngredientMeasurementUnitID: &epazoteSprigVIMU.ID,
				Name:                             "herb sprigs",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
			{

				ValidIngredientPreparationID:     &discardWhiteOnionVIP.ID,
				ValidIngredientMeasurementUnitID: &whiteOnionUnitVIMU.ID,
				Name:                             "whole onion half",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
			{

				ValidIngredientPreparationID:     &discardGarlicVIP.ID,
				ValidIngredientMeasurementUnitID: &garlicCloveVIMU.ID,
				Name:                             "garlic cloves",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 2,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidPreparationVesselID:        &discardLargePotVPV.ID,
				Name:                            "pot with cooked beans and cooking liquid",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "pot with cooked beans and herbs removed",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 11: Mince one half of the onion
	step11 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        mincePrep.ID,
		Index:                11,
		ExplicitInstructions: "Mince one half of the halved onion.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &minceWhiteOnionVIP.ID,
				Name:                            "onion half",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{

				ValidPreparationInstrumentID: &minceChefsKnifeVPI.ID,
				Name:                         "chef's knife",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ValidPreparationVesselID: &minceCuttingBoardVPV.ID,
				Name:                     "cutting board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "minced white onion",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
		},
	}

	// Step 12: In a large skillet, heat lard until shimmering over medium-high heat
	shimmeringState := enums.IngredientStates["shimmering"]
	step12 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        heatPrep.ID,
		Index:                12,
		ExplicitInstructions: "In a large skillet, heat the lard until shimmering over medium-high heat.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &heatLardVIP.ID,
				ValidIngredientMeasurementUnitID: &lardTablespoonVIMU.ID,
				Name:                             "lard",
				QuantityNotes:                    "6 tablespoons (77 g)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 6,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ValidPreparationVesselID: &heatLargeSkilletVPV.ID,
				Name:                     "large skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:  "heated fat in skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: shimmeringState.ID,
				Notes:             "Lard, bacon drippings, or oil should shimmer; butter should foam",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
	}

	// Step 13: Add minced onion and cook, stirring occasionally, until translucent and lightly golden, about 7 minutes
	translucentState := enums.IngredientStates["translucent"]
	step13 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        sautPrep.ID,
		Index:                13,
		ExplicitInstructions: "Add the minced onion and cook, stirring occasionally, until translucent and lightly golden, about 7 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](420), // 7 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](11),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &sautWhiteOnionVIP.ID,
				Name:                            "minced white onion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &sautLargeSkilletVPV.ID,
				Name:                            "heated fat in skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "cooked minced onion",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](0.5),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: translucentState.ID,
				Notes:             "Onion should be translucent and lightly golden",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
	}

	// Step 14: Stir in beans and cook for 2 minutes
	step14 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        stirPrep.ID,
		Index:                14,
		ExplicitInstructions: "Stir in the beans and cook for 2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120), // 2 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirPintoBeansVIP.ID,
				Name:                            "3 cups of cooked beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](13),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &stirWhiteOnionVIP.ID,
				Name:                            "cooked minced onion",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](12),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &stirLargeSkilletVPV.ID,
				Name:                            "heated fat in skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "beans and onion in skillet",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 15: Add 1/4 cup of reserved bean-cooking liquid
	step15 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        addPrep.ID,
		Index:                15,
		ExplicitInstructions: "Add 1/4 cup of the reserved bean-cooking liquid.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &addWaterVIP.ID,
				Name:                            "reserved bean-cooking liquid",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "skillet with beans and onion",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "beans with liquid in skillet",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 16: Using bean masher, potato masher, or back of a wooden spoon, smash the beans to form a chunky purée; alternatively, use a stick blender to make a smoother purée
	step16 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        smashPrep.ID,
		Index:                16,
		ExplicitInstructions: "Using a bean masher, potato masher, or the back of a wooden spoon, smash the beans to form a chunky purée; alternatively, use a stick blender to make a smoother purée.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &smashPintoBeansVIP.ID,
				Name:                            "beans with liquid in skillet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{

				ValidPreparationInstrumentID: &smashPotatoMasherVPI.ID,
				Name:                         "potato masher",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](14),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &smashLargeSkilletVPV.ID,
				Name:                            "skillet with beans and onion",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "mashed beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](3),
				},
			},
		},
	}

	// Step 17: Thin with more bean cooking water until desired consistency is reached. If refried beans become too wet, simmer, stirring, until thickened; if they become too dry, add more bean-cooking liquid, 1 tablespoon at a time, as needed.
	desiredConsistencyState := enums.IngredientStates["at desired consistency"]
	step17 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        dilutePrep.ID,
		Index:                17,
		ExplicitInstructions: "Thin with more bean cooking water until the desired consistency is reached. If the refried beans become too wet, simmer, stirring, until thickened; if they become too dry, add more bean-cooking liquid, 1 tablespoon at a time, as needed.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](16),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &dilutePintoBeansVIP.ID,
				Name:                            "mashed beans",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 3,
				},
			},
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				ValidIngredientPreparationID:    &diluteWaterVIP.ID,
				Name:                            "reserved bean-cooking liquid",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.25,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](15),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &diluteLargeSkilletVPV.ID,
				Name:                            "skillet with beans and onion",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "refried beans at desired consistency",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: desiredConsistencyState.ID,
				Notes:             "Refried beans should reach desired consistency - not too wet or too dry",
				Ingredients:       []uint64{0},
				Optional:          false,
			},
		},
	}

	// Step 18: Season with salt and serve
	step18 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        seasonPrep.ID,
		Index:                18,
		ExplicitInstructions: "Season with salt and serve.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](17),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &seasonPintoBeansVIP.ID,
				Name:                            "refried beans at desired consistency",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
			},
			{

				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltTeaspoonVIMU.ID,
				Name:                             "Kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0.5,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{

				ProductOfRecipeStepIndex:        pointer.To[uint64](17),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &seasonLargeSkilletVPV.ID,
				Name:                            "skillet with beans and onion",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{

				Name:              "refried beans",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &cupMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](4),
				},
			},
		},
	}

	prepTask1 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
		Name:                        "Cook dried pinto beans",
		Description:                 "Halve onion, cover beans with water, peel garlic, add aromatics, bring to a boil, simmer until very tender (1-2 hours), season with salt, drain reserving liquid, and discard aromatics. Cooked beans keep up to 5 days refrigerated.",
		Notes:                       "Reserve the bean-cooking liquid separately; you will need it for the refrying stage.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the drained cooked beans and reserved bean-cooking liquid in separate airtight containers in the refrigerator for up to 5 days.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeAirtightContainer,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 0,
			Max: pointer.To[uint32](432000), // 5 days
		},
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 0, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 1, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 2, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 3, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 4, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 5, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 6, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 7, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 8, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 9, SatisfiesRecipeStep: false},
			{BelongsToRecipeStepIndex: 10, SatisfiesRecipeStep: true},
		},
	}

	prepTask2 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
		Name:                        "Mince onion",
		Description:                 "Mince the remaining onion half for the refrying stage. Minced onion keeps 3-5 days in the fridge.",
		Notes:                       "This is the second half of the onion used in the bean cooking stage.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the minced onion in an airtight container in the refrigerator for up to 3 days.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeAirtightContainer,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4),
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 0,
			Max: pointer.To[uint32](259200), // 3 days
		},
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 11, SatisfiesRecipeStep: true},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
			Name:                "Perfect Frijoles Refritos (Mexican Refried Beans)",
			Slug:                "refried-beans",
			Source:              "https://www.seriouseats.com/perfect-refried-beans",
			Description:         "Use this master recipe to make perfect refried beans in any style: chunky or smooth; with black beans or pintos; and using your choice of cooking fat. By offering choices, including bean type, fat type, and mashing technique, this recipe makes it possible to get exactly the style of refried beans you want.",
			YieldsComponentType: mealplanning.MealComponentTypesSide,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 4,
			},
			PortionName:       "cup",
			PluralPortionName: "cups",
			EligibleForMeals:  true,
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10, step11, step12, step13, step14, step15, step16, step17, step18,
			},
			PrepTasks: []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{prepTask1, prepTask2},
			Media:     []*mealplanning.RecipeMediaCreationRequestInput{},
		},
	}
}
