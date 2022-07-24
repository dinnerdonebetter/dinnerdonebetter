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
	Mozzarella *types.ValidIngredientDatabaseCreationInput
}{
	ChickenBreast: &types.ValidIngredientDatabaseCreationInput{
		ID:          "vi_chicken_breast",
		Name:        "chicken breast",
		AnimalFlesh: true,
	},
	Water: &types.ValidIngredientDatabaseCreationInput{
		ID:       "vi_water",
		Name:     "water",
		IsLiquid: true,
	},
	Onion: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_onion",
		Name: "onion",
	},
	Garlic: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_garlic",
		Name: "garlic",
	},
	BlackPepper: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_black_pepper",
		Name: "black pepper",
	},
	OliveOil: &types.ValidIngredientDatabaseCreationInput{
		ID:                       "vi_olive_oil",
		Name:                     "olive oil",
		IsMeasuredVolumetrically: true,
		IsLiquid:                 true,
	},
	Coffee: &types.ValidIngredientDatabaseCreationInput{
		ID:                       "vi_coffee",
		Name:                     "brewed coffee",
		IsMeasuredVolumetrically: true,
		IsLiquid:                 true,
	},
	Pasta: &types.ValidIngredientDatabaseCreationInput{
		ID:             "vi_pasta",
		Name:           "pasta",
		ContainsGluten: true,
	},
	Tomato: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_tomato",
		Name: "tomato",
	},
	AllPurposeFlour: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_ap_flour",
		Name: "AllPurposeFLour",
	},
	Salt: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_salt",
		Name: "Salt",
	},
	BakingPowder: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_baking_powder",
		Name: "BakingPowder",
	},
	VegetableShortening: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_vegetable_shortening",
		Name: "VegetableShortening",
	},
	HotWater: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_hot_water",
		Name: "HotWater",
	},
	Mozzarella: &types.ValidIngredientDatabaseCreationInput{
		ID:   "vi_mozzarella",
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
