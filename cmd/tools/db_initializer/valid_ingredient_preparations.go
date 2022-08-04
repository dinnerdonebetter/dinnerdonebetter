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
		ID:                 "vip_cb_slice",
		ValidPreparationID: "vprep_slice",
		ValidIngredientID:  "vi_chicken_breast",
	},
	WaterBoil: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_water_boil",
		ValidPreparationID: "vprep_boil",
		ValidIngredientID:  "vi_water",
	},
	OnionSlice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_onion_slice",
		ValidPreparationID: "vprep_slice",
		ValidIngredientID:  "vi_onion",
	},
	OnionDice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_onion_dice",
		ValidPreparationID: "vprep_dice",
		ValidIngredientID:  "vi_onion",
	},
	GarlicSlice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_garlic_mince",
		ValidPreparationID: "vprep_slice",
		ValidIngredientID:  "vi_garlic",
	},
	GarlicMince: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_garlic_mince",
		ValidPreparationID: "vprep_mince",
		ValidIngredientID:  "vi_garlic",
	},
	BlackPepperGrind: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_black_pepper_grind",
		ValidPreparationID: "vprep_grind",
		ValidIngredientID:  "vi_black_pepper",
	},
	PastaBoil: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_pasta_boil",
		ValidPreparationID: "vprep_boil",
		ValidIngredientID:  "vi_pasta",
	},
	TomatoSlice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_tomato_slice",
		ValidPreparationID: "vprep_slice",
		ValidIngredientID:  "vi_tomato",
	},
	TomatoDice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_tomato_dice",
		ValidPreparationID: "vprep_dice",
		ValidIngredientID:  "vi_tomato",
	},
	MozzarellaSlice: &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 "vip_mozzarella_slice",
		ValidPreparationID: "vprep_slice",
		ValidIngredientID:  "vi_mozzarella",
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
