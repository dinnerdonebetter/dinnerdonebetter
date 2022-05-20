package main

import (
	"context"
	"fmt"

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
				Day:               3,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.SpaghettiWithNeatballsAndCapreseSalad.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               3,
			},
			{
				ID:                ksuid.New().String(),
				MealID:            mealCollection.SpaghettiWithNeatballsAndGrilledChicken.ID,
				Notes:             "",
				MealName:          "dinner",
				BelongsToMealPlan: "",
				Day:               3,
			},
		},
		VotingDeadline: 0,
		StartsAt:       0,
		EndsAt:         0,
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
