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
	JonesPastHouseholdMealPlanInput,
	JonesCurrentHouseholdMealPlanInput,
	JonesFutureHouseholdMealPlanInput *types.MealPlanDatabaseCreationInput
	JonesPastHouseholdMealPlan,
	JonesCurrentHouseholdMealPlan,
	JonesFutureHouseholdMealPlan *types.MealPlan
}{
	JonesPastHouseholdMealPlanInput: &types.MealPlanDatabaseCreationInput{
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
		VotingDeadline: uint64(time.Now().Add((-24 * time.Hour) * 7).Add(-10 * time.Minute).Unix()),
		StartsAt:       uint64(time.Now().Add((-24 * time.Hour) * 7).Unix()),
		EndsAt:         uint64(time.Now().Add((-24 * time.Hour) * 7).Add((24 * time.Hour) * 7).Unix()),
	},
	JonesCurrentHouseholdMealPlanInput: &types.MealPlanDatabaseCreationInput{
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
	JonesFutureHouseholdMealPlanInput: &types.MealPlanDatabaseCreationInput{
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
		VotingDeadline: uint64(time.Now().Add((24 * time.Hour) * 7).Add(-10 * time.Minute).Unix()),
		StartsAt:       uint64(time.Now().Add((24 * time.Hour) * 7).Unix()),
		EndsAt:         uint64(time.Now().Add((24 * time.Hour) * 7).Add((24 * time.Hour) * 7).Unix()),
	},
}

func scaffoldMealPlans(ctx context.Context, db database.DataManager) error {
	mealPlanCollection.JonesPastHouseholdMealPlanInput.BelongsToHousehold = jonesHouseholdID
	mealPlanCollection.JonesCurrentHouseholdMealPlanInput.BelongsToHousehold = jonesHouseholdID
	mealPlanCollection.JonesFutureHouseholdMealPlanInput.BelongsToHousehold = jonesHouseholdID

	var err error
	mealPlanCollection.JonesPastHouseholdMealPlan, err = db.CreateMealPlan(ctx, mealPlanCollection.JonesPastHouseholdMealPlanInput)
	if err != nil {
		return fmt.Errorf("voting for meal plan: %w", err)
	}

	if _, finalizationErr := db.AttemptToFinalizeCompleteMealPlan(ctx, mealPlanCollection.JonesPastHouseholdMealPlanInput.ID, jonesHouseholdID); finalizationErr != nil {
		return fmt.Errorf("finalizing meal plan: %w", finalizationErr)
	}

	mealPlanCollection.JonesCurrentHouseholdMealPlan, err = db.CreateMealPlan(ctx, mealPlanCollection.JonesCurrentHouseholdMealPlanInput)
	if err != nil {
		return fmt.Errorf("voting for meal plan: %w", err)
	}

	if _, finalizationErr := db.AttemptToFinalizeCompleteMealPlan(ctx, mealPlanCollection.JonesCurrentHouseholdMealPlanInput.ID, jonesHouseholdID); finalizationErr != nil {
		return fmt.Errorf("finalizing meal plan: %w", finalizationErr)
	}

	mealPlanCollection.JonesFutureHouseholdMealPlan, err = db.CreateMealPlan(ctx, mealPlanCollection.JonesFutureHouseholdMealPlanInput)
	if err != nil {
		return fmt.Errorf("voting for meal plan: %w", err)
	}

	if _, finalizationErr := db.AttemptToFinalizeCompleteMealPlan(ctx, mealPlanCollection.JonesFutureHouseholdMealPlanInput.ID, jonesHouseholdID); finalizationErr != nil {
		return fmt.Errorf("finalizing meal plan: %w", finalizationErr)
	}

	x := mealPlanCollection
	_ = x

	return nil
}
