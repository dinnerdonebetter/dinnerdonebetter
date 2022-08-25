package main

import (
	"context"
	"fmt"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var validIngredientPreparationCollection = struct {
	ChickenBreastSlice,
	WaterBoil,
	OnionSlice,
	OnionDice,
	GarlicSlice,
	GarlicMince,
	BlackPepperGrind,
	PastaBoil,
	TomatoSlice,
	TomatoDice,
	MozzarellaSlice,
	_ *types.ValidIngredientPreparationDatabaseCreationInput
}{
	ChickenBreastSlice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_cb_slice"),
		ValidPreparationID: padID("vprep_slice"),
		ValidIngredientID:  padID("vi_chicken_breast"),
	},
	WaterBoil: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_water_boil"),
		ValidPreparationID: padID("vprep_boil"),
		ValidIngredientID:  padID("vi_water"),
	},
	OnionSlice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_onion_slice"),
		ValidPreparationID: padID("vprep_slice"),
		ValidIngredientID:  padID("vi_onion"),
	},
	OnionDice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_onion_dice"),
		ValidPreparationID: padID("vprep_dice"),
		ValidIngredientID:  padID("vi_onion"),
	},
	GarlicSlice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_garlic_mince"),
		ValidPreparationID: padID("vprep_slice"),
		ValidIngredientID:  padID("vi_garlic"),
	},
	GarlicMince: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_garlic_mince"),
		ValidPreparationID: padID("vprep_mince"),
		ValidIngredientID:  padID("vi_garlic"),
	},
	BlackPepperGrind: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_bp_grind"),
		ValidPreparationID: padID("vprep_grind"),
		ValidIngredientID:  padID("vi_black_pepper"),
	},
	PastaBoil: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_pasta_boil"),
		ValidPreparationID: padID("vprep_boil"),
		ValidIngredientID:  padID("vi_pasta"),
	},
	TomatoSlice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_tomato_slice"),
		ValidPreparationID: padID("vprep_slice"),
		ValidIngredientID:  padID("vi_tomato"),
	},
	TomatoDice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_tomato_dice"),
		ValidPreparationID: padID("vprep_dice"),
		ValidIngredientID:  padID("vi_tomato"),
	},
	MozzarellaSlice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 padID("vip_mozzarella_slice"),
		ValidPreparationID: padID("vprep_slice"),
		ValidIngredientID:  padID("vi_mozzarella"),
	},
}

func scaffoldValidIngredientPreparations(ctx context.Context, db database.DataManager) error {
	validIngredientPreparations := []*types.ValidIngredientPreparationDatabaseCreationInput{
		validIngredientPreparationCollection.ChickenBreastSlice,
		validIngredientPreparationCollection.WaterBoil,
		validIngredientPreparationCollection.OnionSlice,
		validIngredientPreparationCollection.OnionDice,
		validIngredientPreparationCollection.GarlicMince,
		validIngredientPreparationCollection.BlackPepperGrind,
		validIngredientPreparationCollection.PastaBoil,
		validIngredientPreparationCollection.TomatoSlice,
		validIngredientPreparationCollection.TomatoDice,
		validIngredientPreparationCollection.MozzarellaSlice,
	}

	for i, input := range validIngredientPreparations {
		if _, err := db.CreateValidIngredientPreparation(ctx, input); err != nil {
			return fmt.Errorf("creating preparation instrument #%d: %w", i, err)
		}
	}

	return nil
}
