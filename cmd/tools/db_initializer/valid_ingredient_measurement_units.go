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
		ID:                     padID("vimu_cbg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_chicken_breast"),
	},
	WaterGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_wtrg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_water"),
	},
	WaterMilliliter: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_wtrml"),
		ValidMeasurementUnitID: padID("vmu_milliliter"),
		ValidIngredientID:      padID("vi_water"),
	},
	OnionUnits: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_onionu"),
		ValidMeasurementUnitID: padID("vmu_unit"),
		ValidIngredientID:      padID("vi_onion"),
	},
	OnionGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_oniong"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_onion"),
	},
	GarlicGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_garlicg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_garlic"),
	},
	GarlicCloves: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_garlicc"),
		ValidMeasurementUnitID: padID("vmu_clove"),
		ValidIngredientID:      padID("vi_garlic"),
	},
	BlackPepperGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_bpeppg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_black_pepper"),
	},
	OliveOilGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_olioilg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_olive_oil"),
	},
	OliveOilMilliliter: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_olioilml"),
		ValidMeasurementUnitID: padID("vmu_milliliter"),
		ValidIngredientID:      padID("vi_olive_oil"),
	},
	CoffeeGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_coffeeg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_coffee"),
	},
	PastaGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_pastag"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_pasta"),
	},
	TomatoGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_tomatog"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_tomato"),
	},
	AllPurposeFlourGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_apfg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_ap_flour"),
	},
	SaltGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_saltg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_salt"),
	},
	SaltTeaspoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_salttsp"),
		ValidMeasurementUnitID: padID("vmu_teaspoon"),
		ValidIngredientID:      padID("vi_salt"),
	},
	SaltTablespoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_salttbsp"),
		ValidMeasurementUnitID: padID("vmu_tablespoon"),
		ValidIngredientID:      padID("vi_salt"),
	},
	BakingPowderGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_bakpwdg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_baking_powder"),
	},
	BakingPowderTeaspoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_bakpwdtsp"),
		ValidMeasurementUnitID: padID("vmu_teaspoon"),
		ValidIngredientID:      padID("vi_baking_powder"),
	},
	BakingPowderTablespoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_bakpwdtbsp"),
		ValidMeasurementUnitID: padID("vmu_tablespoon"),
		ValidIngredientID:      padID("vi_baking_powder"),
	},
	VegetableShorteningGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_veg_oil_g"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_vegetable_shortening"),
	},
	VegetableShorteningTeaspoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_veg_oil_tsp"),
		ValidMeasurementUnitID: padID("vmu_teaspoon"),
		ValidIngredientID:      padID("vi_vegetable_shortening"),
	},
	VegetableShorteningTablespoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_veg_oil_tbsp"),
		ValidMeasurementUnitID: padID("vmu_tablespoon"),
		ValidIngredientID:      padID("vi_vegetable_shortening"),
	},
	HotWaterGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_hotwtr_g"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_hot_water"),
	},
	HotWaterTeaspoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_hotwtr_tsp"),
		ValidMeasurementUnitID: padID("vmu_teaspoon"),
		ValidIngredientID:      padID("vi_hot_water"),
	},
	HotWaterTablespoon: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_hotwtr_tbsp"),
		ValidMeasurementUnitID: padID("vmu_tablespoon"),
		ValidIngredientID:      padID("vi_hot_water"),
	},
	MozzarellaGrams: &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     padID("vimu_mozzg"),
		ValidMeasurementUnitID: padID("vmu_gram"),
		ValidIngredientID:      padID("vi_mozzarella"),
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
