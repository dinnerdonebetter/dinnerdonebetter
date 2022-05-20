package main

import (
	"context"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var mealCollection = struct {
	MushroomRisottoWithGrilledChicken,
	SpaghettiWithNeatballsAndGrilledChicken,
	SpaghettiWithNeatballsAndCapreseSalad *types.MealDatabaseCreationInput
}{
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
			recipeCollection.SpaghettiWithNeatballs.ID,
			recipeCollection.GrilledChicken.ID,
		},
	},
	SpaghettiWithNeatballsAndCapreseSalad: &types.MealDatabaseCreationInput{
		ID:            ksuid.New().String(),
		Name:          "spaghetti with neatballs and a caprese salad",
		Description:   "spaghetti with a nice caprese salad",
		CreatedByUser: userCollection.MomJones.ID,
		Recipes: []string{
			recipeCollection.SpaghettiWithNeatballs.ID,
			recipeCollection.CapreseSalad.ID,
		},
	},
}

func scaffoldMeals(ctx context.Context, db database.DataManager) error {
	meals := []*types.MealDatabaseCreationInput{
		mealCollection.MushroomRisottoWithGrilledChicken,
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
