package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var mealPlanCollection = struct {
	JonesHouseholdMealPlan *types.MealPlanDatabaseCreationInput
}{
	JonesHouseholdMealPlan: &types.MealPlanDatabaseCreationInput{
		ID:                 ksuid.New().String(),
		Status:             "finalized",
		BelongsToHousehold: jonesHouseholdID,
		Notes:              "",
		Options: []*types.MealPlanOptionDatabaseCreationInput{
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.MushroomRisottoWithGrilledChicken.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Tuesday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.SpaghettiWithNeatballsAndCapreseSalad.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Tuesday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.EggFriedRice.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Tuesday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.TacosAndCollardGreens.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Wednesday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.LasagnaAndRamen.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Wednesday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.GrilledCheeseSandwiches.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Wednesday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.GrilledCheeseSandwichesWithBakedPotato.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Thursday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.EggFriedRice.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Thursday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.Lasagna.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Thursday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.TacosAndEggFriedRice.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Friday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.BakedPotatoAndMashedPotatoes.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Friday,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.Ramen.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               time.Friday,
			},
		},
		VotingDeadline: uint64(time.Now().Add(-10 * time.Minute).Unix()),
		StartsAt:       uint64(time.Now().Unix()),
		EndsAt:         uint64(time.Now().Add((24 * time.Hour) * 7).Unix()),
	},
}

func scaffoldMealPlans(ctx context.Context, db database.DataManager) error {
	mealPlanCollection.JonesHouseholdMealPlan.BelongsToHousehold = jonesHouseholdID

	mealPlans := []*types.MealPlanDatabaseCreationInput{
		mealPlanCollection.JonesHouseholdMealPlan,
	}

	for _, input := range mealPlans {
		if _, err := db.CreateMealPlan(ctx, input); err != nil {
			return fmt.Errorf("creating meal plan: %w", err)
		}

		if _, err := db.AttemptToFinalizeCompleteMealPlan(ctx, input.ID, jonesHouseholdID); err != nil {
			return fmt.Errorf("finalizing meal plan: %w", err)
		}
	}

	return nil
}
