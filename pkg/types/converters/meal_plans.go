package converters

import (
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertMealPlanToMealPlanUpdateRequestInput creates a MealPlanUpdateRequestInput from a MealPlan.
func ConvertMealPlanToMealPlanUpdateRequestInput(input *types.MealPlan) *types.MealPlanUpdateRequestInput {
	x := &types.MealPlanUpdateRequestInput{
		BelongsToHousehold: &input.BelongsToHousehold,
		Notes:              &input.Notes,
		VotingDeadline:     &input.VotingDeadline,
	}

	return x
}

// ConvertMealPlanCreationRequestInputToMealPlanDatabaseCreationInput creates a MealPlanDatabaseCreationInput from a MealPlanCreationRequestInput.
func ConvertMealPlanCreationRequestInputToMealPlanDatabaseCreationInput(input *types.MealPlanCreationRequestInput) *types.MealPlanDatabaseCreationInput {
	events := []*types.MealPlanEventDatabaseCreationInput{}
	for _, e := range input.Events {
		events = append(events, ConvertMealPlanEventCreationRequestInputToMealPlanEventDatabaseCreationInput(e))
	}

	x := &types.MealPlanDatabaseCreationInput{
		Notes:          input.Notes,
		VotingDeadline: input.VotingDeadline,
		Events:         events,
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
		ID:                 mealPlan.ID,
		Notes:              mealPlan.Notes,
		VotingDeadline:     mealPlan.VotingDeadline,
		Events:             events,
		BelongsToHousehold: mealPlan.BelongsToHousehold,
	}
}

// ConvertMealPlanToMealPlanDatabaseCreationInput builds a MealPlanDatabaseCreationInput from a MealPlan.
func ConvertMealPlanToMealPlanDatabaseCreationInput(mealPlan *types.MealPlan) *types.MealPlanDatabaseCreationInput {
	events := []*types.MealPlanEventDatabaseCreationInput{}
	for _, event := range mealPlan.Events {
		events = append(events, ConvertMealPlanEventToMealPlanEventDatabaseCreationInput(event))
	}

	return &types.MealPlanDatabaseCreationInput{
		ID:                 mealPlan.ID,
		Notes:              mealPlan.Notes,
		VotingDeadline:     mealPlan.VotingDeadline,
		Events:             events,
		BelongsToHousehold: mealPlan.BelongsToHousehold,
	}
}
