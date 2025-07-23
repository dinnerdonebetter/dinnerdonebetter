package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
)

func ConvertMealPlanTaskToGRPCMealPlanTask(input *mealplanning.MealPlanTask) *mealplanningsvc.MealPlanTask {
	return &mealplanningsvc.MealPlanTask{
		RecipePrepTask:      ConvertRecipePrepTaskToGRPCRecipePrepTask(&input.RecipePrepTask),
		CreatedAt:           grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		CompletedAt:         grpcconverters.ConvertTimePointerToPBTimestamp(input.CompletedAt),
		AssignedToUser:      input.AssignedToUser,
		ID:                  input.ID,
		Status:              input.Status,
		CreationExplanation: input.CreationExplanation,
		StatusExplanation:   input.StatusExplanation,
		MealPlanOption:      ConvertMealPlanOptionToGRPCMealPlanOption(&input.MealPlanOption),
	}
}

func ConvertMealPlanTaskDatabaseCreationEstimateToGRPCMealPlanTask(input *mealplanning.MealPlanTaskDatabaseCreationEstimate) *mealplanningsvc.MealPlanTask {
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

func ConvertMealPlanGroceryListItemToGRPCMealPlanGroceryListItem(input *mealplanning.MealPlanGroceryListItem) *mealplanningsvc.MealPlanGroceryListItem {
	return &mealplanningsvc.MealPlanGroceryListItem{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		QuantityNeeded: &types.Float32RangeWithOptionalMax{
			Max: input.QuantityNeeded.Max,
			Min: input.QuantityNeeded.Min,
		},
		Ingredient:               ConvertValidIngredientToGRPCValidIngredient(&input.Ingredient),
		MeasurementUnit:          ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(&input.MeasurementUnit),
		PurchasedMeasurementUnit: ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(input.PurchasedMeasurementUnit),
		PurchasedUPC:             input.PurchasedUPC,
		Status:                   input.Status,
		StatusExplanation:        input.StatusExplanation,
		ID:                       input.ID,
		BelongsToMealPlan:        input.BelongsToMealPlan,
		PurchasePrice:            input.PurchasePrice,
		QuantityPurchased:        input.QuantityPurchased,
	}
}

func ConvertMealPlanEventToGRPCMealPlanEvent(input *mealplanning.MealPlanEvent) *mealplanningsvc.MealPlanEvent {
	var mealPlanOptions []*mealplanningsvc.MealPlanOption
	for _, option := range input.Options {
		mealPlanOptions = append(mealPlanOptions, ConvertMealPlanOptionToGRPCMealPlanOption(option))
	}

	return &mealplanningsvc.MealPlanEvent{
		CreatedAt:         grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:     grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:        grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		StartsAt:          grpcconverters.ConvertTimeToPBTimestamp(input.StartsAt),
		EndsAt:            grpcconverters.ConvertTimeToPBTimestamp(input.EndsAt),
		MealName:          input.MealName,
		Notes:             input.Notes,
		BelongsToMealPlan: input.BelongsToMealPlan,
		ID:                input.ID,
		Options:           mealPlanOptions,
	}
}

func ConvertMealPlanToGRPCMealPlan(input *mealplanning.MealPlan) *mealplanningsvc.MealPlan {
	var mealPlanEvents []*mealplanningsvc.MealPlanEvent
	for _, event := range input.Events {
		mealPlanEvents = append(mealPlanEvents, ConvertMealPlanEventToGRPCMealPlanEvent(event))
	}

	return &mealplanningsvc.MealPlan{
		CreatedAt:              grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:          grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:             grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		VotingDeadline:         grpcconverters.ConvertTimeToPBTimestamp(input.VotingDeadline),
		ElectionMethod:         input.ElectionMethod,
		Status:                 input.Status,
		Notes:                  input.Notes,
		ID:                     input.ID,
		BelongsToAccount:       input.BelongsToAccount,
		CreatedByUser:          input.CreatedByUser,
		Events:                 mealPlanEvents,
		GroceryListInitialized: input.GroceryListInitialized,
		TasksCreated:           input.TasksCreated,
	}
}

func ConvertUserIngredientPreferenceToGRPCUserIngredientPreference(input *mealplanning.UserIngredientPreference) *mealplanningsvc.UserIngredientPreference {
	return &mealplanningsvc.UserIngredientPreference{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		Ingredient:    ConvertValidIngredientToGRPCValidIngredient(&input.Ingredient),
		ID:            input.ID,
		Notes:         input.Notes,
		BelongsToUser: input.BelongsToUser,
		Rating:        int32(input.Rating),
		Allergy:       input.Allergy,
	}
}
