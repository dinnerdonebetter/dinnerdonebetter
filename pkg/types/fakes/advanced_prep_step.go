package fakes

import (
	"time"

	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeAdvancedPrepStep builds a faked advanced prep step.
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
		Status:               "unfinished",
		StatusExplanation:    buildUniqueString(),
		CreationExplanation:  buildUniqueString(),
		SettledAt:            nil,
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

// BuildFakeAdvancedPrepStepDatabaseCreationInputs builds a faked AdvancedPrepStepList.
func BuildFakeAdvancedPrepStepDatabaseCreationInputs() []*types.AdvancedPrepStepDatabaseCreationInput {
	var examples []*types.AdvancedPrepStepDatabaseCreationInput
	for i := 0; i < exampleQuantity; i++ {
		now := time.Now().Add(0).Truncate(time.Second).UTC()
		inOneWeek := now.Add((time.Hour * 24) * 7).Add(0).Truncate(time.Second).UTC()

		examples = append(examples, &types.AdvancedPrepStepDatabaseCreationInput{
			CompletedAt:          nil,
			MealPlanOptionID:     "",
			RecipeStepID:         "",
			ID:                   BuildFakeID(),
			CannotCompleteBefore: now,
			CannotCompleteAfter:  inOneWeek,
			Status:               "unfinished",
			StatusExplanation:    buildUniqueString(),
			CreationExplanation:  buildUniqueString(),
		})
	}

	return examples
}

// BuildFakeAdvancedPrepStepStatusChangeRequestInput builds a faked advanced prep step.
func BuildFakeAdvancedPrepStepStatusChangeRequestInput() *types.AdvancedPrepStepStatusChangeRequestInput {
	return &types.AdvancedPrepStepStatusChangeRequestInput{
		ID:                BuildFakeID(),
		Status:            "unfinished",
		StatusExplanation: buildUniqueString(),
	}
}
