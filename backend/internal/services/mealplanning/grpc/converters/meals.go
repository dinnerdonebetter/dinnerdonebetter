package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
)

func ConvertMealPlanTaskToGRPCMealPlanTask(input *mealplanning.MealPlanTaskDatabaseCreationEstimate) *mealplanningsvc.MealPlanTask {
	return &mealplanningsvc.MealPlanTask{
		CreationExplanation: input.CreationExplanation,
	}
}

func ConvertMealPlanOptionToGRPCMealPlanOption(input *mealplanning.MealPlanOption) *mealplanningsvc.MealPlanOption {
	var votes []*mealplanningsvc.MealPlanOptionVote
	for _, vote := range input.Votes {
		votes = append(votes, ConvertMealPlanOptionVoteToGRPCMealPlanOptionVote(vote))
	}

	return &mealplanningsvc.MealPlanOption{
		CreatedAt:              grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:          grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:             grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		Meal:                   ConvertMealToGRPCMeal(&input.Meal),
		ID:                     input.ID,
		Notes:                  input.Notes,
		BelongsToMealPlanEvent: input.BelongsToMealPlanEvent,
		AssignedDishwasher:     input.AssignedDishwasher,
		AssignedCook:           input.AssignedCook,
		Votes:                  votes,
		MealScale:              input.MealScale,
		Chosen:                 input.Chosen,
		TieBroken:              input.TieBroken,
	}
}

func ConvertMealPlanOptionVoteToGRPCMealPlanOptionVote(input *mealplanning.MealPlanOptionVote) *mealplanningsvc.MealPlanOptionVote {
	return &mealplanningsvc.MealPlanOptionVote{
		CreatedAt:               grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:           grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:              grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		ID:                      input.ID,
		Notes:                   input.Notes,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		ByUser:                  input.ByUser,
		Rank:                    uint32(input.Rank),
		Abstain:                 input.Abstain,
	}
}

func ConvertMealToGRPCMeal(input *mealplanning.Meal) *mealplanningsvc.Meal {
	var components []*mealplanningsvc.MealComponent
	for _, component := range input.Components {
		components = append(components, ConvertMealComponentToGRPCMealComponent(component))
	}

	return &mealplanningsvc.Meal{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		EstimatedPortions: &types.Float32RangeWithOptionalMax{
			Max: input.EstimatedPortions.Max,
			Min: input.EstimatedPortions.Min,
		},
		ID:                   input.ID,
		Description:          input.Description,
		CreatedByUser:        input.CreatedByUser,
		Name:                 input.Name,
		Components:           components,
		EligibleForMealPlans: input.EligibleForMealPlans,
	}
}

func ConvertMealComponentToGRPCMealComponent(input *mealplanning.MealComponent) *mealplanningsvc.MealComponent {
	return &mealplanningsvc.MealComponent{
		Recipe:        ConvertRecipeToGRPCRecipe(&input.Recipe),
		ComponentType: input.ComponentType,
		RecipeScale:   input.RecipeScale,
	}
}
