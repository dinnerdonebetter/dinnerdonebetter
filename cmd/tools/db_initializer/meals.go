package main

import (
	"context"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var mealCollection = struct {
	MushroomRisotto,
	GrilledChicken,
	MushroomRisottoWithGrilledChicken,
	Spaghetti,
	Neatballs,
	CapreseSalad,
	GrilledCheeseSandwiches,
	BakedPotato,
	GrilledCheeseSandwichesWithBakedPotato,
	Ramen,
	Lasagna,
	LasagnaAndRamen,
	Tacos,
	TacosAndLasagna,
	TacosAndRamen,
	EggFriedRice,
	TacosAndEggFriedRice,
	MashedPotatoes,
	BakedPotatoAndMashedPotatoes,
	CollardGreens,
	TacosAndCollardGreens,
	SpaghettiWithNeatballs,
	SpaghettiWithNeatballsAndGrilledChicken,
	SpaghettiWithNeatballsAndCapreseSalad *types.MealDatabaseCreationInput
}{
	MushroomRisotto: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "mushroom risotto",
		Description:   "a nice mushroom risotto",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.MushroomRisotto.ID,
		},
	},
	GrilledChicken: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "grilled chicken",
		Description:   "chicken",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.GrilledChicken.ID,
		},
	},
	MushroomRisottoWithGrilledChicken: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "mushroom risotto with grilled chicken",
		Description:   "chicken with a nice mushroom risotto",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.MushroomRisotto.ID,
			recipeCollection.GrilledChicken.ID,
		},
	},
	SpaghettiWithNeatballsAndGrilledChicken: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "spaghetti with grilled chicken",
		Description:   "spaghetti with a nice grilled chicken",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Spaghetti.ID,
			recipeCollection.GrilledChicken.ID,
		},
	},
	Spaghetti: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "spaghetti",
		Description:   "spaghetti",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Spaghetti.ID,
		},
	},
	Neatballs: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "neatballs",
		Description:   "neat",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Neatballs.ID,
		},
	},
	SpaghettiWithNeatballs: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "spaghetti with neatballs",
		Description:   "spaghetti with neatballs",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Spaghetti.ID,
			recipeCollection.Neatballs.ID,
		},
	},
	CapreseSalad: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "a caprese salad",
		Description:   "a nice caprese salad",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.CapreseSalad.ID,
		},
	},
	SpaghettiWithNeatballsAndCapreseSalad: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "spaghetti with neatballs and a caprese salad",
		Description:   "spaghetti with a nice caprese salad",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Spaghetti.ID,
			recipeCollection.Neatballs.ID,
			recipeCollection.CapreseSalad.ID,
		},
	},
	GrilledCheeseSandwiches: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "grilled cheese sandwiches",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.GrilledCheeseSandwiches.ID,
		},
	},
	BakedPotato: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "baked potato",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.BakedPotato.ID,
		},
	},
	GrilledCheeseSandwichesWithBakedPotato: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "grilled cheese sandwiches with baked potato",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.GrilledCheeseSandwiches.ID,
			recipeCollection.BakedPotato.ID,
		},
	},
	Ramen: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "ramen",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Ramen.ID,
		},
	},
	Lasagna: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "lasagna",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Lasagna.ID,
		},
	},
	LasagnaAndRamen: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "lasagna and ramen",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Lasagna.ID,
			recipeCollection.Ramen.ID,
		},
	},
	Tacos: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "tacos",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Tacos.ID,
		},
	},
	TacosAndLasagna: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "tacos and lasagna",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Tacos.ID,
			recipeCollection.Lasagna.ID,
		},
	},
	TacosAndRamen: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "tacos and ramen",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Tacos.ID,
			recipeCollection.Ramen.ID,
		},
	},
	EggFriedRice: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "egg fried rice",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.EggFriedRice.ID,
		},
	},
	TacosAndEggFriedRice: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "tacos and egg fried rice",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Tacos.ID,
			recipeCollection.EggFriedRice.ID,
		},
	},
	MashedPotatoes: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "mashed potatoes",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.MashedPotatoes.ID,
		},
	},
	BakedPotatoAndMashedPotatoes: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "baked potato and mashed potatoes",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.BakedPotato.ID,
			recipeCollection.MashedPotatoes.ID,
		},
	},
	CollardGreens: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "collard greens",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.CollardGreens.ID,
		},
	},
	TacosAndCollardGreens: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "tacos and collard greens",
		Description:   "",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.Tacos.ID,
			recipeCollection.CollardGreens.ID,
		},
	},
}

func scaffoldMeals(ctx context.Context, db database.DataManager) error {
	meals := []*types.MealDatabaseCreationInput{
		mealCollection.MushroomRisotto,
		mealCollection.GrilledChicken,
		mealCollection.MushroomRisottoWithGrilledChicken,
		mealCollection.Spaghetti,
		mealCollection.Neatballs,
		mealCollection.CapreseSalad,
		mealCollection.GrilledCheeseSandwiches,
		mealCollection.BakedPotato,
		mealCollection.GrilledCheeseSandwichesWithBakedPotato,
		mealCollection.Ramen,
		mealCollection.Lasagna,
		mealCollection.LasagnaAndRamen,
		mealCollection.Tacos,
		mealCollection.TacosAndLasagna,
		mealCollection.TacosAndRamen,
		mealCollection.EggFriedRice,
		mealCollection.TacosAndEggFriedRice,
		mealCollection.MashedPotatoes,
		mealCollection.BakedPotatoAndMashedPotatoes,
		mealCollection.CollardGreens,
		mealCollection.TacosAndCollardGreens,
		mealCollection.SpaghettiWithNeatballs,
		mealCollection.SpaghettiWithNeatballsAndGrilledChicken,
		mealCollection.SpaghettiWithNeatballsAndCapreseSalad,
	}

	for _, input := range meals {
		if _, err := db.CreateMeal(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
