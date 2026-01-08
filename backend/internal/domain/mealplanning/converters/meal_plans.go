package converters

import (
	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

// ConvertMealPlanToMealPlanUpdateRequestInput creates a MealPlanUpdateRequestInput from a MealPlan.
func ConvertMealPlanToMealPlanUpdateRequestInput(input *types.MealPlan) *types.MealPlanUpdateRequestInput {
	x := &types.MealPlanUpdateRequestInput{
		BelongsToAccount: &input.BelongsToAccount,
		Notes:            &input.Notes,
		VotingDeadline:   &input.VotingDeadline,
	}

	return x
}

// ConvertMealPlanCreationRequestInputToMealPlanDatabaseCreationInput creates a MealPlanDatabaseCreationInput from a MealPlanCreationRequestInput.
func ConvertMealPlanCreationRequestInputToMealPlanDatabaseCreationInput(input *types.MealPlanCreationRequestInput) *types.MealPlanDatabaseCreationInput {
	mealPlanID := identifiers.New()
	events := []*types.MealPlanEventDatabaseCreationInput{}
	for _, e := range input.Events {
		eventInput := ConvertMealPlanEventCreationRequestInputToMealPlanEventDatabaseCreationInput(e)
		eventInput.BelongsToMealPlan = mealPlanID
		events = append(events, eventInput)
	}

	x := &types.MealPlanDatabaseCreationInput{
		ID:             mealPlanID,
		Notes:          input.Notes,
		CreatedByUser:  identifiers.New(),
		VotingDeadline: input.VotingDeadline,
		Events:         events,
		ElectionMethod: input.ElectionMethod,
	}

	return x
}

// ConvertMealPlanToMealPlanCreationRequestInput builds a MealPlanCreationRequestInput from a MealPlan.
func ConvertMealPlanToMealPlanCreationRequestInput(mealPlan *types.MealPlan) *types.MealPlanCreationRequestInput {
	events := []*types.MealPlanEventCreationRequestInput{}
	for _, evt := range mealPlan.Events {
		events = append(events, ConvertMealPlanEventToMealPlanEventCreationRequestInput(evt))
	}

	return &types.MealPlanCreationRequestInput{
		Notes:          mealPlan.Notes,
		VotingDeadline: mealPlan.VotingDeadline,
		Events:         events,
		ElectionMethod: mealPlan.ElectionMethod,
	}
}

// ConvertMealPlanToMealPlanDatabaseCreationInput builds a MealPlanDatabaseCreationInput from a MealPlan.
func ConvertMealPlanToMealPlanDatabaseCreationInput(mealPlan *types.MealPlan) *types.MealPlanDatabaseCreationInput {
	events := []*types.MealPlanEventDatabaseCreationInput{}
	for _, event := range mealPlan.Events {
		events = append(events, ConvertMealPlanEventToMealPlanEventDatabaseCreationInput(event))
	}

	selections := []*types.MealPlanRecipeOptionSelectionDatabaseCreationInput{}
	for _, selection := range mealPlan.Selections {
		selections = append(selections, ConvertMealPlanRecipeOptionSelectionToMealPlanRecipeOptionSelectionDatabaseCreationInput(selection))
	}

	return &types.MealPlanDatabaseCreationInput{
		ID:               mealPlan.ID,
		Notes:            mealPlan.Notes,
		VotingDeadline:   mealPlan.VotingDeadline,
		Events:           events,
		ElectionMethod:   mealPlan.ElectionMethod,
		BelongsToAccount: mealPlan.BelongsToAccount,
		CreatedByUser:    mealPlan.CreatedByUser,
		Selections:       selections,
	}
}
