package main

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"

	"github.com/segmentio/ksuid"
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
	Mozzarella *types.ValidIngredientDatabaseCreationInput
}{
	ChickenBreast: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "chicken breast",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    false,
		AnimalFlesh:       true,
		AnimalDerived:     true,
		Volumetric:        false,
		ContainsSoy:       false,
	},
	Water: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "water",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    false,
		AnimalFlesh:       false,
		AnimalDerived:     false,
		Volumetric:        false,
		ContainsSoy:       false,
	},
	Onion: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "onion",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    false,
		AnimalFlesh:       false,
		AnimalDerived:     false,
		Volumetric:        false,
		ContainsSoy:       false,
	},
	Garlic: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "garlic",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    false,
		AnimalFlesh:       false,
		AnimalDerived:     false,
		Volumetric:        false,
		ContainsSoy:       false,
	},
	BlackPepper: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "black pepper",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    false,
		AnimalFlesh:       false,
		AnimalDerived:     false,
		Volumetric:        false,
		ContainsSoy:       false,
	},
	OliveOil: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "olive oil",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    false,
		AnimalFlesh:       false,
		AnimalDerived:     false,
		Volumetric:        false,
		ContainsSoy:       false,
	},
	Coffee: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "coffee",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    false,
		AnimalFlesh:       false,
		AnimalDerived:     false,
		Volumetric:        false,
		ContainsSoy:       false,
	},
	Pasta: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "pasta",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    true,
		AnimalFlesh:       false,
		AnimalDerived:     false,
		Volumetric:        false,
		ContainsSoy:       false,
	},
	Tomato: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "tomato",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    false,
		AnimalFlesh:       false,
		AnimalDerived:     false,
		Volumetric:        false,
		ContainsSoy:       false,
	},
	Mozzarella: &types.ValidIngredientDatabaseCreationInput{
		ID:                ksuid.New().String(),
		Name:              "mozzarella",
		Variant:           "",
		Description:       "",
		Warning:           "",
		IconPath:          "",
		ContainsDairy:     false,
		ContainsPeanut:    false,
		ContainsTreeNut:   false,
		ContainsEgg:       false,
		ContainsWheat:     false,
		ContainsShellfish: false,
		ContainsSesame:    false,
		ContainsFish:      false,
		ContainsGluten:    false,
		AnimalFlesh:       false,
		AnimalDerived:     true,
		Volumetric:        false,
		ContainsSoy:       false,
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
		validIngredientCollection.Mozzarella,
	}

	for _, input := range validIngredients {
		if _, err := db.CreateValidIngredient(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
