package main

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var validIngredientCollection = struct {
	ChickenBreast,
	Water,
	Onion,
	Garlic,
	BlackPepper,
	OliveOil,
	Coffee,
	Pasta,
	Tomato,
	AllPurposeFlour,
	Salt,
	BakingPowder,
	VegetableShortening,
	HotWater,
	Mozzarella,
	_ *types.ValidIngredientDatabaseCreationInput
}{
	ChickenBreast: &types.ValidIngredientDatabaseCreationInput{
		ID:          padID("vi_chicken_breast"),
		Name:        "chicken breast",
		AnimalFlesh: true,
	},
	Water: &types.ValidIngredientDatabaseCreationInput{
		ID:       padID("vi_water"),
		Name:     "water",
		IsLiquid: true,
	},
	Onion: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_onion"),
		Name: "onion",
	},
	Garlic: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_garlic"),
		Name: "garlic",
	},
	BlackPepper: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_black_pepper"),
		Name: "black pepper",
	},
	OliveOil: &types.ValidIngredientDatabaseCreationInput{
		ID:                       padID("vi_olive_oil"),
		Name:                     "olive oil",
		IsMeasuredVolumetrically: true,
		IsLiquid:                 true,
	},
	Coffee: &types.ValidIngredientDatabaseCreationInput{
		ID:                       padID("vi_coffee"),
		Name:                     "brewed coffee",
		IsMeasuredVolumetrically: true,
		IsLiquid:                 true,
	},
	Pasta: &types.ValidIngredientDatabaseCreationInput{
		ID:             padID("vi_pasta"),
		Name:           "pasta",
		ContainsGluten: true,
	},
	Tomato: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_tomato"),
		Name: "tomato",
	},
	AllPurposeFlour: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_ap_flour"),
		Name: "all purpose flour",
	},
	Salt: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_salt"),
		Name: "salt",
	},
	BakingPowder: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_baking_powder"),
		Name: "baking powder",
	},
	VegetableShortening: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_vegetable_shortening"),
		Name: "vegetable shortening",
	},
	HotWater: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_hot_water"),
		Name: "hot water",
	},
	Mozzarella: &types.ValidIngredientDatabaseCreationInput{
		ID:   padID("vi_mozzarella"),
		Name: "mozzarella",
	},
}

func scaffoldValidIngredients(ctx context.Context, db database.DataManager) error {
	validIngredients := []*types.ValidIngredientDatabaseCreationInput{
		validIngredientCollection.ChickenBreast,
		validIngredientCollection.Water,
		validIngredientCollection.Onion,
		validIngredientCollection.Garlic,
		validIngredientCollection.BlackPepper,
		validIngredientCollection.OliveOil,
		validIngredientCollection.Coffee,
		validIngredientCollection.Pasta,
		validIngredientCollection.Tomato,
		validIngredientCollection.AllPurposeFlour,
		validIngredientCollection.Salt,
		validIngredientCollection.BakingPowder,
		validIngredientCollection.VegetableShortening,
		validIngredientCollection.HotWater,
		validIngredientCollection.Mozzarella,
	}

	for _, input := range validIngredients {
		if _, err := db.CreateValidIngredient(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
