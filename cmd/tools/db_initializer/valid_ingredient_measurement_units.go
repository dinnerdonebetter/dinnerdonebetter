package main

import (
	"context"
	"fmt"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var validIngredientMeasurementUnitCollection = struct {
	ChickenBreastGrams,
	WaterGrams,
	WaterMilliliter,
	OnionUnits,
	OnionGrams,
	GarlicGrams,
	GarlicCloves,
	BlackPepperGrams,
	OliveOilGrams,
	OliveOilMilliliter,
	CoffeeGrams,
	PastaGrams,
	TomatoGrams,
	AllPurposeFlourGrams,
	SaltGrams,
	SaltTeaspoon,
	SaltTablespoon,
	BakingPowderGrams,
	BakingPowderTeaspoon,
	BakingPowderTablespoon,
	VegetableShorteningGrams,
	VegetableShorteningTeaspoon,
	VegetableShorteningTablespoon,
	HotWaterGrams,
	HotWaterTeaspoon,
	HotWaterTablespoon,
	MozzarellaGrams,
	_ *types.ValidIngredientMeasurementUnitDatabaseCreationInput
}{
	ChickenBreastGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_cbg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_chicken_breast",
	},
	WaterGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_wtrg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_water",
	},
	WaterMilliliter: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_wtrml",
		ValidMeasurementUnitID: "vmu_milliliter",
		ValidIngredientID:      "vi_water",
	},
	OnionUnits: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_onionu",
		ValidMeasurementUnitID: "vmu_unit",
		ValidIngredientID:      "vi_onion",
	},
	OnionGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_oniong",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_onion",
	},
	GarlicGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_garlicg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_garlic",
	},
	GarlicCloves: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_garlicc",
		ValidMeasurementUnitID: "vmu_clove",
		ValidIngredientID:      "vi_garlic",
	},
	BlackPepperGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_bpeppg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_black_pepper",
	},
	OliveOilGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_olioilg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_olive_oil",
	},
	OliveOilMilliliter: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_olioilml",
		ValidMeasurementUnitID: "vmu_milliliter",
		ValidIngredientID:      "vi_olive_oil",
	},
	CoffeeGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_coffeeg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_coffee",
	},
	PastaGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_pastag",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_pasta",
	},
	TomatoGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_tomatog",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_tomato",
	},
	AllPurposeFlourGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_apfg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_ap_flour",
	},
	SaltGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_saltg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_salt",
	},
	SaltTeaspoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_salttsp",
		ValidMeasurementUnitID: "vmu_teaspoon",
		ValidIngredientID:      "vi_salt",
	},
	SaltTablespoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_salttbsp",
		ValidMeasurementUnitID: "vmu_tablespoon",
		ValidIngredientID:      "vi_salt",
	},
	BakingPowderGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_bakpwdg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_baking_powder",
	},
	BakingPowderTeaspoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_bakpwdtsp",
		ValidMeasurementUnitID: "vmu_teaspoon",
		ValidIngredientID:      "vi_baking_powder",
	},
	BakingPowderTablespoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_bakpwdtbsp",
		ValidMeasurementUnitID: "vmu_tablespoon",
		ValidIngredientID:      "vi_baking_powder",
	},
	VegetableShorteningGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_veg_oil_g",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_vegetable_shortening",
	},
	VegetableShorteningTeaspoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_veg_oil_tsp",
		ValidMeasurementUnitID: "vmu_teaspoon",
		ValidIngredientID:      "vi_vegetable_shortening",
	},
	VegetableShorteningTablespoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_veg_oil_tbsp",
		ValidMeasurementUnitID: "vmu_tablespoon",
		ValidIngredientID:      "vi_vegetable_shortening",
	},
	HotWaterGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_hotwtr_g",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_hot_water",
	},
	HotWaterTeaspoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_hotwtr_tsp",
		ValidMeasurementUnitID: "vmu_teaspoon",
		ValidIngredientID:      "vi_hot_water",
	},
	HotWaterTablespoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_hotwtr_tbsp",
		ValidMeasurementUnitID: "vmu_tablespoon",
		ValidIngredientID:      "vi_hot_water",
	},
	MozzarellaGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     "vimu_mozzg",
		ValidMeasurementUnitID: "vmu_gram",
		ValidIngredientID:      "vi_mozzarella",
	},
}

func scaffoldValidIngredientMeasurementUnits(ctx context.Context, db database.DataManager) error {
	validIngredientMeasurementUnits := []*types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		validIngredientMeasurementUnitCollection.ChickenBreastGrams,
		validIngredientMeasurementUnitCollection.WaterGrams,
		validIngredientMeasurementUnitCollection.WaterMilliliter,
		validIngredientMeasurementUnitCollection.OnionUnits,
		validIngredientMeasurementUnitCollection.OnionGrams,
		validIngredientMeasurementUnitCollection.GarlicGrams,
		validIngredientMeasurementUnitCollection.GarlicCloves,
		validIngredientMeasurementUnitCollection.BlackPepperGrams,
		validIngredientMeasurementUnitCollection.OliveOilGrams,
		validIngredientMeasurementUnitCollection.OliveOilMilliliter,
		validIngredientMeasurementUnitCollection.CoffeeGrams,
		validIngredientMeasurementUnitCollection.PastaGrams,
		validIngredientMeasurementUnitCollection.TomatoGrams,
		validIngredientMeasurementUnitCollection.AllPurposeFlourGrams,
		validIngredientMeasurementUnitCollection.SaltGrams,
		validIngredientMeasurementUnitCollection.SaltTeaspoon,
		validIngredientMeasurementUnitCollection.SaltTablespoon,
		validIngredientMeasurementUnitCollection.BakingPowderGrams,
		validIngredientMeasurementUnitCollection.BakingPowderTeaspoon,
		validIngredientMeasurementUnitCollection.BakingPowderTablespoon,
		validIngredientMeasurementUnitCollection.VegetableShorteningGrams,
		validIngredientMeasurementUnitCollection.VegetableShorteningTeaspoon,
		validIngredientMeasurementUnitCollection.VegetableShorteningTablespoon,
		validIngredientMeasurementUnitCollection.HotWaterGrams,
		validIngredientMeasurementUnitCollection.HotWaterTeaspoon,
		validIngredientMeasurementUnitCollection.HotWaterTablespoon,
		validIngredientMeasurementUnitCollection.MozzarellaGrams,
	}

	for i, input := range validIngredientMeasurementUnits {
		if _, err := db.CreateValidIngredientMeasurementUnit(ctx, input); err != nil {
			return fmt.Errorf("creating preparation instrument #%d: %w", i, err)
		}
	}

	return nil
}
