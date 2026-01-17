package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

func PerfectRoastChickenRecipe(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	// Get preparations
	grindPrep := enums.Preparations["grind"]
	mixPrep := enums.Preparations["mix"]
	seasonPrep := enums.Preparations["season"]
	trussPrep := enums.Preparations["truss"]
	dryBrinePrep := enums.Preparations["dry brine"]
	preheatPrep := enums.Preparations["preheat"]
	heatPrep := enums.Preparations["heat"]
	rubPrep := enums.Preparations["rub"]
	panSearPrep := enums.Preparations["pan-sear"]
	roastPrep := enums.Preparations["roast"]
	restPrep := enums.Preparations["rest"]

	// Get ingredients
	wholeChicken := enums.Ingredients["whole chicken"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	wholePeppercorns := enums.Ingredients["whole black peppercorns"]
	bakingPowder := enums.Ingredients["baking powder"]
	vegetableOil := enums.Ingredients["vegetable oil"]

	// Get measurement units
	unitMeasurement := enums.MeasurementUnits["unit"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	milliliterMeasurement := enums.MeasurementUnits["milliliter"]

	// Get instruments
	butchersTwine := enums.Instruments["butcher's twine"]
	tongs := enums.Instruments["tongs"]
	thermometer := enums.Instruments["instant-read thermometer"]
	bareHands := enums.Instruments["bare hands"]
	mortarAndPestle := enums.Instruments["mortar and pestle"]
	spiceGrinder := enums.Instruments["spice grinder"]

	// Get vessels
	smallBowl := enums.Vessels["small bowl"]
	wireRack := enums.Vessels["wire rack"]
	bakingSheet := enums.Vessels["baking sheet"]
	oven := enums.Vessels["oven"]
	stainlessSteelSkillet := enums.Vessels["stainless steel skillet"]
	carvingBoard := enums.Vessels["carving board"]

	// Get ingredient states for completion conditions
	brownedState := enums.IngredientStates["browned"]
	atTemperatureState := enums.IngredientStates["at temperature"]

	// Get bridge table entries
	// Grind preparation bridges
	grindPeppercornsVIP := enums.IngredientPreparations[grindPrep.ID][wholePeppercorns.ID]
	peppercornsGramVIMU := enums.IngredientMeasurementUnits[wholePeppercorns.ID][gramMeasurement.ID]
	grindMortarAndPestleVPI := enums.PreparationInstruments[grindPrep.ID][mortarAndPestle.ID]
	grindSpiceGrinderVPI := enums.PreparationInstruments[grindPrep.ID][spiceGrinder.ID]

	// Mix preparation bridges
	mixSaltVIP := enums.IngredientPreparations[mixPrep.ID][salt.ID]
	mixPepperVIP := enums.IngredientPreparations[mixPrep.ID][blackPepper.ID]
	mixBakingPowderVIP := enums.IngredientPreparations[mixPrep.ID][bakingPowder.ID]
	saltGramVIMU := enums.IngredientMeasurementUnits[salt.ID][gramMeasurement.ID]
	pepperGramVIMU := enums.IngredientMeasurementUnits[blackPepper.ID][gramMeasurement.ID]
	bakingPowderGramVIMU := enums.IngredientMeasurementUnits[bakingPowder.ID][gramMeasurement.ID]
	mixSmallBowlVPV := enums.PreparationVessels[mixPrep.ID][smallBowl.ID]

	// Season preparation bridges
	seasonChickenVIP := enums.IngredientPreparations[seasonPrep.ID][wholeChicken.ID]
	chickenUnitVIMU := enums.IngredientMeasurementUnits[wholeChicken.ID][unitMeasurement.ID]
	seasonBareHandsVPI := enums.PreparationInstruments[seasonPrep.ID][bareHands.ID]

	// Truss preparation bridges
	trussTwineVPI := enums.PreparationInstruments[trussPrep.ID][butchersTwine.ID]

	// Dry-brine preparation bridges
	dryBrineWireRackVPV := enums.PreparationVessels[dryBrinePrep.ID][wireRack.ID]
	dryBrineBakingSheetVPV := enums.PreparationVessels[dryBrinePrep.ID][bakingSheet.ID]

	// Preheat preparation bridges
	preheatOvenVPV := enums.PreparationVessels[preheatPrep.ID][oven.ID]

	// Heat preparation bridges
	heatOilVIP := enums.IngredientPreparations[heatPrep.ID][vegetableOil.ID]
	oilTablespoonVIMU := enums.IngredientMeasurementUnits[vegetableOil.ID][tablespoonMeasurement.ID]
	heatSkilletVPV := enums.PreparationVessels[heatPrep.ID][stainlessSteelSkillet.ID]

	// Rub preparation bridges
	rubOilVIP := enums.IngredientPreparations[rubPrep.ID][vegetableOil.ID]
	rubBareHandsVPI := enums.PreparationInstruments[rubPrep.ID][bareHands.ID]

	// Pan-sear preparation bridges
	panSearTongsVPI := enums.PreparationInstruments[panSearPrep.ID][tongs.ID]

	// Roast preparation bridges
	roastThermometerVPI := enums.PreparationInstruments[roastPrep.ID][thermometer.ID]
	roastOvenVPV := enums.PreparationVessels[roastPrep.ID][oven.ID]

	// Rest preparation bridges
	restCarvingBoardVPV := enums.PreparationVessels[restPrep.ID][carvingBoard.ID]

	// Step 0: Grind whole black peppercorns (optional)
	step0 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: grindPrep.ID,
		Index:         0,
		Optional:      true,
		Notes:         "Using a mortar and pestle or spice grinder, coarsely grind the whole black peppercorns.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &grindPeppercornsVIP.ID,
				ValidIngredientMeasurementUnitID: &peppercornsGramVIMU.ID,
				Name:                             "whole black peppercorns",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1, // approximately 0.5 teaspoon
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

	// Step 1: Mix the seasoning
	step1 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: mixPrep.ID,
		Index:         1,
		Notes:         "In a small bowl, thoroughly mix the salt with black pepper and baking powder (if using).",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &mixSaltVIP.ID,
				ValidIngredientMeasurementUnitID: &saltGramVIMU.ID,
				Name:                             "kosher salt",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 9, // 1 tablespoon = 9g
				},
				Index: pointer.To[uint16](0),
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](0),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "freshly ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1, // 0.5 teaspoon ≈ 1g
				},
				Optional: true,
				Index:    pointer.To[uint16](1),
			},
			{
				ValidIngredientPreparationID:     &mixPepperVIP.ID,
				ValidIngredientMeasurementUnitID: &pepperGramVIMU.ID,
				Name:                             "pre-ground black pepper",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1, // 0.5 teaspoon ≈ 1g
				},
				Optional:    true,
				Index:       pointer.To[uint16](1),
				OptionIndex: 1,
			},
			{
				ValidIngredientPreparationID:     &mixBakingPowderVIP.ID,
				ValidIngredientMeasurementUnitID: &bakingPowderGramVIMU.ID,
				Name:                             "baking powder",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 4, // 1 teaspoon ≈ 4g
				},
				Optional: true,
				Index:    pointer.To[uint16](2),
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &mixSmallBowlVPV.ID,
				Name:                     "small bowl",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "seasoning mixture",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &gramMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](12), // approximately 12g total
				},
			},
		},
	}

	// Step 2: Season the chicken
	step2 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: seasonPrep.ID,
		Index:         2,
		Notes:         "Season chicken all over, inside and out, with salt mixture (or just plain salt if not using pepper and baking powder).",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &seasonChickenVIP.ID,
				ValidIngredientMeasurementUnitID: &chickenUnitVIMU.ID,
				Name:                             "large whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](1),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "seasoning mixture",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 12,
				},
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
				Name:              "seasoned whole chicken",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 3: Truss the chicken
	step3 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: trussPrep.ID,
		Index:         3,
		Notes:         "Set chicken, breast side up, on work surface and tuck wings behind back. Using butcher's twine, run the center of the twine under the tip of the tail end and truss chicken by tying drumsticks together at their bony ends, securing the legs and the tip of the tail together in a bundle. Criss-cross the twine and pass along the crevasse where the legs meet the breast; pass twine over wings to hold them into place, then tie securely around the stump of the neck.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](2),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "seasoned whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &trussTwineVPI.ID,
				Name:                         "butcher's twine",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "trussed seasoned whole chicken",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 4: Dry-brine (refrigerate)
	step4 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: dryBrinePrep.ID,
		Index:         4,
		Notes:         "Place chicken, back side down, on a wire rack set in a rimmed baking sheet and refrigerate, uncovered, at least 1 hour and up to 2 days.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](3600),   // 1 hour minimum
			Max: pointer.To[uint32](172800), // 2 days maximum
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](3),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "trussed seasoned whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &dryBrineWireRackVPV.ID,
				Name:                     "wire rack",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidPreparationVesselID: &dryBrineBakingSheetVPV.ID,
				Name:                     "rimmed baking sheet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "dry-brined whole chicken",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 5: Preheat the oven
	step5 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: preheatPrep.ID,
		Index:         5,
		Notes:         "Adjust oven rack to middle position and preheat oven to 425°F (220°C).",
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](220), // 425°F
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &preheatOvenVPV.ID,
				Name:                     "oven",
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
			},
		},
	}

	// Step 6: Heat oil in skillet
	step6 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: heatPrep.ID,
		Index:         6,
		Notes:         "In a 10- or 12-inch stainless steel skillet, heat oil over medium-high heat until shimmering.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     &heatOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilTablespoonVIMU.ID,
				Name:                             "neutral-flavored oil (vegetable or canola)",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &heatSkilletVPV.ID,
				Name:                     "stainless steel skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "pre-heated oil",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &milliliterMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](15),
				},
			},
			{
				Name:  "pre-heated skillet",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
					Max: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 7: Rub chicken with oil
	step7 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: rubPrep.ID,
		Index:         7,
		Notes:         "Rub chicken lightly with oil.",
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](4),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "dry-brined whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ValidIngredientPreparationID:     &rubOilVIP.ID,
				ValidIngredientMeasurementUnitID: &oilTablespoonVIMU.ID,
				Name:                             "canola oil",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &rubBareHandsVPI.ID,
				Name:                         "bare hands",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "oiled whole chicken",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	// Step 8: Brown chicken legs
	step8 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: panSearPrep.ID,
		Index:         8,
		Notes:         "Set chicken on its side in the skillet so that the full thigh and drumstick are in contact with the pan; the wing will also be touching, but the breast should have little to no contact with the skillet. Cook until leg is well browned, 8 to 10 minutes, then flip bird so other leg is touching pan and repeat; lower heat at any point if chicken skin begins to burn.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](960),  // 16 minutes minimum (8 min per side)
			Max: pointer.To[uint32](1200), // 20 minutes maximum (10 min per side)
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](7),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "oiled whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "heated oil in skillet",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 15,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &panSearTongsVPI.ID,
				Name:                         "tongs",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](6),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "stainless steel skillet",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
					Max: pointer.To[uint16](1),
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "leg-browned whole chicken",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
			{
				Name:  "stainless steel skillet with browned chicken",
				Type:  mealplanning.RecipeStepProductVesselType,
				Index: 1,
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: brownedState.ID,
				Notes:             "Both chicken legs should be well browned",
				Ingredients:       []uint64{0}, // Index of the chicken ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 9: Roast the chicken
	step9 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: roastPrep.ID,
		Index:         9,
		Notes:         "Using hands and spatula if needed, rotate chicken so it is breast side up in the skillet and transfer to oven. Roast until breast registers 150°F (65°C) in the center of its thickest part and thighs register 165°F (75°C) near (but not touching) the bone, about 40 minutes.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](2400), // 40 minutes
		},
		TemperatureInCelsius: types.OptionalFloat32Range{
			Min: pointer.To[float32](65), // 150°F for breast
			Max: pointer.To[float32](75), // 165°F for thighs
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "leg-browned whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: &roastThermometerVPI.ID,
				Name:                         "instant-read thermometer",
				Quantity: types.Uint32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](5),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				ValidPreparationVesselID:        &roastOvenVPV.ID,
				Name:                            "preheated oven",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](8),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](1),
				Name:                            "stainless steel skillet with browned chicken",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "roasted whole chicken",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
		CompletionConditions: []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{
			{
				IngredientStateID: atTemperatureState.ID,
				Notes:             "Breast should register 150°F (65°C) and thighs 165°F (75°C)",
				Ingredients:       []uint64{0}, // Index of the chicken ingredient in the step
				Optional:          false,
			},
		},
	}

	// Step 10: Rest the chicken
	step10 := &mealplanning.RecipeStepCreationRequestInput{
		PreparationID: restPrep.ID,
		Index:         10,
		Notes:         "Remove from oven and transfer chicken to a carving board. Let rest 10 to 20 minutes, then carve and serve.",
		EstimatedTimeInSeconds: types.OptionalUint32Range{
			Min: pointer.To[uint32](600),  // 10 minutes
			Max: pointer.To[uint32](1200), // 20 minutes
		},
		Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
			{
				ProductOfRecipeStepIndex:        pointer.To[uint64](9),
				ProductOfRecipeStepProductIndex: pointer.To[uint64](0),
				Name:                            "roasted whole chicken",
				Quantity: types.Float32RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Vessels: []*mealplanning.RecipeStepVesselCreationRequestInput{
			{
				ValidPreparationVesselID: &restCarvingBoardVPV.ID,
				Name:                     "carving board",
				Quantity: types.Uint16RangeWithOptionalMax{
					Min: 1,
				},
			},
		},
		Products: []*mealplanning.RecipeStepProductCreationRequestInput{
			{
				Name:              "rested roast chicken",
				Type:              mealplanning.RecipeStepProductIngredientType,
				Index:             0,
				MeasurementUnitID: &unitMeasurement.ID,
				MeasurementQuantity: types.OptionalFloat32Range{
					Min: pointer.To[float32](1),
				},
			},
		},
	}

	return []*mealplanning.RecipeCreationRequestInput{
		{
			Name:                "Perfect Roast Chicken",
			Slug:                "perfect-roast-chicken",
			Source:              "https://www.seriouseats.com/perfect-roast-chicken-recipe-8384377",
			Description:         "A dry-brine fully seasons the chicken and allows the skin to dehydrate, improving browning and crisping during cooking. Measuring the internal temperature of the chicken to determine doneness leads to more reliable and superior results than going by time. Rubbing the chicken skin with oil before roasting, instead of basting with watery drippings, ensures even browning and a crisp skin.",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Min: 4,
			},
			PortionName:       "serving",
			PluralPortionName: "servings",
			EligibleForMeals:  true,
			Steps:             []*mealplanning.RecipeStepCreationRequestInput{step0, step1, step2, step3, step4, step5, step6, step7, step8, step9, step10},
			PrepTasks:         []*mealplanning.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Media:             []*mealplanning.RecipeMediaCreationRequestInput{},
			AlsoCreateMeal:    false,
		},
	}
}
