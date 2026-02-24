package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func SousVideChickenBreastRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	grindPrep := enums.Preparations["grind"]
	heatPrep := enums.Preparations["heat"]
	seasonPrep := enums.Preparations["season"]
	bagPrep := enums.Preparations["bag"]
	sousVidePrep := enums.Preparations["sous-vide"]
	removePrep := enums.Preparations["remove"]
	dryPrep := enums.Preparations["dry"]
	panSearPrep := enums.Preparations["pan-sear"]
	restPrep := enums.Preparations["rest"]

	// Get ingredients
	chickenBreast := enums.Ingredients["chicken breast"]
	salt := enums.Ingredients["salt"]
	wholePeppercorns := enums.Ingredients["whole black peppercorns"]
	thyme := enums.Ingredients["thyme"]
	rosemary := enums.Ingredients["rosemary"]
	vegetableOil := enums.Ingredients["vegetable oil"]
	paperTowelsIngredient := enums.Ingredients["paper towels"]

	// Get measurement units
	gramMeasurement := enums.MeasurementUnits["gram"]
	milliliterMeasurement := enums.MeasurementUnits["milliliter"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]
	unitMeasurement := enums.MeasurementUnits["unit"]

	// Get instruments
	mortarAndPestle := enums.Instruments["mortar and pestle"]
	spiceGrinder := enums.Instruments["spice grinder"]
	sousVideCooker := enums.Instruments["sous vide cooker"]
	spatula := enums.Instruments["spatula"]
	tongs := enums.Instruments["tongs"]
	bareHands := enums.Instruments["bare hands"]
	stovetop := enums.Instruments["stovetop"]

	// Get vessels
	waterBath := enums.Vessels["water bath"]
	plasticBag := enums.Vessels["plastic bag"]
	vacuumBag := enums.Vessels["vacuum bag"]
	castIronSkillet := enums.Vessels["cast iron skillet"]
	servingPlate := enums.Vessels["serving plate"]

	// Get ingredient states for completion conditions
	atTemperatureState := enums.IngredientStates["at temperature"]
	smokingState := enums.IngredientStates["smoking"]

	// Get bridge table entries
	// Grind preparation bridges
	grindPeppercornsVIP := enums.IngredientPreparations[grindPrep.ID][wholePeppercorns.ID]
	peppercornsGramVIMU := enums.IngredientMeasurementUnits[wholePeppercorns.ID][gramMeasurement.ID]
	grindMortarAndPestleVPI := enums.PreparationInstruments[grindPrep.ID][mortarAndPestle.ID]
	grindSpiceGrinderVPI := enums.PreparationInstruments[grindPrep.ID][spiceGrinder.ID]

	// Heat preparation bridges (for preheating water bath)
	heatSousVideCookerVPI := enums.PreparationInstruments[heatPrep.ID][sousVideCooker.ID]
	heatWaterBathVPV := enums.PreparationVessels[heatPrep.ID][waterBath.ID]

	// Season preparation bridges
	seasonChickenVIP := enums.IngredientPreparations[seasonPrep.ID][chickenBreast.ID]
	seasonSaltVIP := enums.IngredientPreparations[seasonPrep.ID][salt.ID]
	chickenGramVIMU := enums.IngredientMeasurementUnits[chickenBreast.ID][gramMeasurement.ID]
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	seasonBareHandsVPI := enums.PreparationInstruments[seasonPrep.ID][bareHands.ID]

	// Bag preparation bridges
	bagThymeVIP := enums.IngredientPreparations[bagPrep.ID][thyme.ID]
	bagRosemaryVIP := enums.IngredientPreparations[bagPrep.ID][rosemary.ID]
	thymeSprigVIMU := enums.IngredientMeasurementUnits[thyme.ID][sprigMeasurement.ID]
	rosemarySprigVIMU := enums.IngredientMeasurementUnits[rosemary.ID][sprigMeasurement.ID]
	bagPlasticBagVPV := enums.PreparationVessels[bagPrep.ID][plasticBag.ID]
	bagVacuumBagVPV := enums.PreparationVessels[bagPrep.ID][vacuumBag.ID]

	// Remove preparation bridges (for removing chicken from bag)
	removeChickenVIP := enums.IngredientPreparations[removePrep.ID][chickenBreast.ID]
	removeTongsVPI := enums.PreparationInstruments[removePrep.ID][tongs.ID]

	// Dry preparation bridges (for drying chicken)
	dryChickenVIP := enums.IngredientPreparations[dryPrep.ID][chickenBreast.ID]
	dryPaperTowelsVIP := enums.IngredientPreparations[dryPrep.ID][paperTowelsIngredient.ID]
	paperTowelsUnitVIMU := enums.IngredientMeasurementUnits[paperTowelsIngredient.ID][unitMeasurement.ID]
	dryBareHandsVPI := enums.PreparationInstruments[dryPrep.ID][bareHands.ID]

	// Heat preparation bridges (for heating oil in skillet)
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][vegetableOil.ID]
	oilMilliliterVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][milliliterMeasurement.ID]
	heatStovetopVPI := enums.PreparationInstruments[heatPrep.ID][stovetop.ID]
	heatSkilletVPV := enums.PreparationVessels[heatPrep.ID][castIronSkillet.ID]

	// Pan-sear preparation bridges (for finishing)
	panSearSpatulaVPI := enums.PreparationInstruments[panSearPrep.ID][spatula.ID]
	panSearTongsVPI := enums.PreparationInstruments[panSearPrep.ID][tongs.ID]
	panSearSkilletVPV := enums.PreparationVessels[panSearPrep.ID][castIronSkillet.ID]

	// Rest preparation bridges
	restTongsVPI := enums.PreparationInstruments[restPrep.ID][tongs.ID]
	restPlateVPV := enums.PreparationVessels[restPrep.ID][servingPlate.ID]

	// Step 0: Grind whole black peppercorns
	step0 := &mealplanning.RecipeStepCreationRequestInput{
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

	// Step 1: Preheat water bath
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        heatPrep.ID,
		Index:                1,
		ExplicitInstructions: "Preheat a water bath to 150°F (66°C) using a sous vide cooker. This temperature produces tender and juicy chicken, ideal for chicken salad when served cold.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](66), // 150°F = 66°C
			Max: pointer.To[float32](66), // 150°F = 66°C
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &heatSousVideCookerVPI.ID,
				Name:                         "sous vide cooker",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &heatWaterBathVPV.ID,
				Name:                     "water bath",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:  "preheated water bath",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 0,
			},
		},
	}

	// Step 2: Season chicken
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        seasonPrep.ID,
		Index:                2,
		ExplicitInstructions: "Season the chicken generously with salt and pepper.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &seasonChickenVIP.ID,
				ValidIngredientMeasurementUnitID: &chickenGramVIMU.ID,
				Name:                             "boneless chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 900,
				},
			},
			{
				ValidIngredientPreparationID:     &seasonSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltGramVIMU.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste: true,
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 0,
				},
				ToTaste: true,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &seasonBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seasoned boneless chicken breasts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](900),
				},
			},
		},
	}

	// Step 3: Bag chicken
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        bagPrep.ID,
		Index:                3,
		ExplicitInstructions: "Place the chicken in zipper-lock bags or vacuum bags and add thyme or rosemary sprigs, if using. If using zipper-lock bags: Remove air by closing the bags, leaving the last inch of the top unsealed. Slowly lower into the preheated water bath, sealing the bag completely just before it fully submerges. If using vacuum bags: Seal according to the manufacturer's instructions.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "seasoned boneless chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 900,
				},
			},
			{
				ValidIngredientPreparationID:     &bagThymeVIP.ID,
				ValidIngredientMeasurementUnitID: &thymeSprigVIMU.ID,
				Name:                             "thyme",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Optional: true,
			},
			{
				ValidIngredientPreparationID:     &bagRosemaryVIP.ID,
				ValidIngredientMeasurementUnitID: &rosemarySprigVIMU.ID,
				Name:                             "rosemary",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4,
				},
				Optional: true,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &bagPlasticBagVPV.ID,
				Name:                     "zipper-lock bag",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationVesselID: &bagVacuumBagVPV.ID,
				Name:                     "vacuum bag",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "bagged seasoned boneless chicken breasts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](900),
				},
			},
		},
	}

	// Step 4: Cook sous vide
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        sousVidePrep.ID,
		Index:                4,
		ExplicitInstructions: "Add the bagged chicken to the preheated water bath and cook at 150°F (66°C) for 1 to 4 hours.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](3600),  // 1 hour minimum
			Max: pointer.To[uint32](14400), // 4 hours maximum
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](66), // 150°F = 66°C
			Max: pointer.To[float32](66), // 150°F = 66°C
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "bagged seasoned boneless chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 900,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "preheated sous vide cooker",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex: pointer.To[uint64](1),
				Name:                     "preheated water bath",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "sous vide cooked boneless chicken breasts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](900),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: atTemperatureState.ID,
				Notes:             "Chicken should reach 150°F (66°C) and be held at that temperature for at least 1 hour",
				Ingredients:       []uint64{0}, // Index of the chicken ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 5: Remove chicken from water bath and bag
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        removePrep.ID,
		Index:                5,
		Optional:             true,
		ExplicitInstructions: "Turn on your vents and open your windows. Remove the chicken from the water bath and bag. Discard herbs, if using.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &removeChickenVIP.ID,
				Name:                            "sous vide cooked boneless chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 900,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &removeTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "unbagged sous vide boneless chicken breasts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](900),
				},
			},
		},
	}

	// Step 6: Dry chicken with paper towels
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        dryPrep.ID,
		Index:                6,
		Optional:             true,
		ExplicitInstructions: "Carefully pat the chicken dry with paper towels.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidIngredientPreparationID:    &dryChickenVIP.ID,
				Name:                            "unbagged sous vide boneless chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 900,
				},
			},
			{
				ValidIngredientPreparationID:     &dryPaperTowelsVIP.ID,
				ValidIngredientMeasurementUnitID: &paperTowelsUnitVIMU.ID,
				Name:                             "paper towels",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &dryBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "dried sous vide boneless chicken breasts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](900),
				},
			},
		},
	}

	// Step 7: Heat oil in cast iron skillet
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        heatPrep.ID,
		Index:                7,
		Optional:             true,
		ExplicitInstructions: "Heat the oil in a heavy cast iron or stainless steel skillet over medium-high heat until shimmering.",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](180), // Medium-high heat
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilMilliliterVIMU.ID,
				Name:                             "vegetable oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 15, // Enough to coat the pan
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &heatStovetopVPI.ID,
				Name:                         "stovetop",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &heatSkilletVPV.ID,
				Name:                     "cast iron skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "heated shimmering oil",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &milliliterMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](15),
				},
			},
			{
				Name:  "cast iron skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				ItemQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: smokingState.ID,
				Notes:             "Oil should be shimmering",
				Ingredients:       []uint64{0}, // Index of the oil ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 8: Pan-sear the chicken
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        panSearPrep.ID,
		Index:                8,
		Optional:             true,
		ExplicitInstructions: "Gently lay the chicken in the skillet using your fingers or a set of tongs. Hold the chicken down flat in the pan with a flexible metal spatula or your fingers (be careful of splattering oil). Cook until golden brown and crisp, about 2 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120), // 2 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "dried sous vide boneless chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 900,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "heated shimmering oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 15,
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &panSearSpatulaVPI.ID,
				Name:                         "flexible metal spatula",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationInstrumentID: &panSearTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
				Index:       pointer.To[uint16](1),
				OptionIndex: 0,
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID:        &panSearSkilletVPV.ID,
				Name:                            "cast iron skillet",
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "pan-seared sous vide boneless chicken breasts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](900),
				},
			},
			{
				Name:  "cast iron skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				ItemQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: atTemperatureState.ID,
				Notes:             "Chicken should be golden brown and crisp",
				Ingredients:       []uint64{0}, // Index of the chicken ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 9: Rest and serve
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID:        restPrep.ID,
		Index:                9,
		ExplicitInstructions: "Remove from the pan and let rest until cool enough to handle, about 2 minutes. Slice the chicken and serve immediately.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](120), // 2 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "pan-seared sous vide boneless chicken breasts",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 900,
				},
				Optional: true, // From step 7 if pan finishing chosen
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &restTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &restPlateVPV.ID,
				Name:                     "serving plate",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "sliced sous vide boneless chicken breasts",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](900),
				},
			},
		},
	}

	// Create prep task for seasoning and bagging chicken ahead of time
	prepTask1 := &mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{
		Name:                        "Season and bag chicken breasts",
		Description:                 "The chicken breasts can be seasoned and sealed in bags up to 24 hours ahead of time. Store in the refrigerator in sealed bags.",
		Notes:                       "Preparing the chicken ahead saves time on the day of cooking.",
		Optional:                    true,
		ExplicitStorageInstructions: "Store the sealed chicken breasts in the refrigerator for up to 24 hours.",
		StorageType:                 mealplanning.RecipePrepTaskStorageTypeAirtightContainer,
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: pointer.To[float32](4), // Refrigerator temperature
		},
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMax{
			Min: 0,
			Max: pointer.To[uint32](86400), // 24 hours
		},
		RecipeSteps: []*mealplanning.RecipePrepTaskStepWithinRecipeCreationRequestInput{
			{BelongsToRecipeStepIndex: 2, SatisfiesRecipeStep: true},
			{BelongsToRecipeStepIndex: 3, SatisfiesRecipeStep: true},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
			Name:                "Sous Vide Chicken Breast",
			Slug:                "sous-vide-chicken-breast",
			Source:              "https://www.seriouseats.com/sous-vide-chicken-breast-recipe",
			Description:         "",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 2,
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps:             []*mealplanning.RecipeStepCreationRequestInput{step0, step1, step2, step3, step4, step5, step6, step7, step8, step9},
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{prepTask1},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
