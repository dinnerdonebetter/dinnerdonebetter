package fakes

import (
	"time"

	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeAdvancedPrepStep builds a faked meal plan event.
func BuildFakeAdvancedPrepStep() *types.AdvancedPrepStep {
	now := time.Now().Add(0).Truncate(time.Second).UTC()
	inOneWeek := now.Add((time.Hour * 24) * 7).Add(0).Truncate(time.Second).UTC()

	return &types.AdvancedPrepStep{
		ID:                   BuildFakeID(),
		CannotCompleteBefore: now,
		CannotCompleteAfter:  inOneWeek,
		MealPlanOption:       *BuildFakeMealPlanOption(),
		RecipeStep:           *BuildFakeRecipeStep(),
		CreatedAt:            fake.Date(),
		CompletedAt:          nil,
	}
}

// BuildFakeAdvancedPrepStepList builds a faked AdvancedPrepStepList.
func BuildFakeAdvancedPrepStepList() *types.AdvancedPrepStepList {
	var examples []*types.AdvancedPrepStep
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeAdvancedPrepStep())
	}

	return &types.AdvancedPrepStepList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		AdvancedPrepSteps: examples,
	}
}
