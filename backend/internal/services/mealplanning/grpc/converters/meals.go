package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	grpctypes "github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
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
		EstimatedPortions: &grpctypes.Float32RangeWithOptionalMax{
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

func ConvertGRPCMealToMeal(input *mealplanningsvc.Meal) *mealplanning.Meal {
	var components []*mealplanning.MealComponent
	for _, component := range input.Components {
		components = append(components, ConvertGRPCMealComponentToMealComponent(component))
	}

	return &mealplanning.Meal{
		CreatedAt:     grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		ArchivedAt:    grpcconverters.ConvertPBTimestampToTimePointer(input.ArchivedAt),
		EstimatedPortions: types.Float32RangeWithOptionalMax{
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

func ConvertGRPCMealComponentToMealComponent(input *mealplanningsvc.MealComponent) *mealplanning.MealComponent {
	return &mealplanning.MealComponent{
		Recipe:        *ConvertGRPCRecipeToRecipe(input.Recipe),
		ComponentType: input.ComponentType,
		RecipeScale:   input.RecipeScale,
	}
}

func ConvertMealPlanGroceryListItemToGRPCMealPlanGroceryListItem(input *mealplanning.MealPlanGroceryListItem) *mealplanningsvc.MealPlanGroceryListItem {
	return &mealplanningsvc.MealPlanGroceryListItem{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		QuantityNeeded: &grpctypes.Float32RangeWithOptionalMax{
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

func ConvertGRPCMealPlanToMealPlan(input *mealplanningsvc.MealPlan) *mealplanning.MealPlan {
	var mealPlanEvents []*mealplanning.MealPlanEvent
	for _, event := range input.Events {
		mealPlanEvents = append(mealPlanEvents, ConvertGRPCMealPlanEventToMealPlanEvent(event))
	}

	return &mealplanning.MealPlan{
		CreatedAt:              grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		LastUpdatedAt:          grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		ArchivedAt:             grpcconverters.ConvertPBTimestampToTimePointer(input.ArchivedAt),
		VotingDeadline:         grpcconverters.ConvertPBTimestampToTime(input.VotingDeadline),
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

func ConvertGRPCMealPlanEventToMealPlanEvent(input *mealplanningsvc.MealPlanEvent) *mealplanning.MealPlanEvent {
	var mealPlanOptions []*mealplanning.MealPlanOption
	for _, option := range input.Options {
		mealPlanOptions = append(mealPlanOptions, ConvertGRPCMealPlanOptionToMealPlanOption(option))
	}

	return &mealplanning.MealPlanEvent{
		CreatedAt:         grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		LastUpdatedAt:     grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		ArchivedAt:        grpcconverters.ConvertPBTimestampToTimePointer(input.ArchivedAt),
		StartsAt:          grpcconverters.ConvertPBTimestampToTime(input.StartsAt),
		EndsAt:            grpcconverters.ConvertPBTimestampToTime(input.EndsAt),
		MealName:          input.MealName,
		Notes:             input.Notes,
		BelongsToMealPlan: input.BelongsToMealPlan,
		ID:                input.ID,
		Options:           mealPlanOptions,
	}
}

func ConvertGRPCMealPlanOptionToMealPlanOption(input *mealplanningsvc.MealPlanOption) *mealplanning.MealPlanOption {
	var votes []*mealplanning.MealPlanOptionVote
	for _, vote := range input.Votes {
		votes = append(votes, ConvertGRPCMealPlanOptionVoteToMealPlanOptionVote(vote))
	}

	return &mealplanning.MealPlanOption{
		CreatedAt:              grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		LastUpdatedAt:          grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		ArchivedAt:             grpcconverters.ConvertPBTimestampToTimePointer(input.ArchivedAt),
		Meal:                   *ConvertGRPCMealToMeal(input.Meal),
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

func ConvertGRPCMealPlanOptionVoteToMealPlanOptionVote(input *mealplanningsvc.MealPlanOptionVote) *mealplanning.MealPlanOptionVote {
	return &mealplanning.MealPlanOptionVote{
		CreatedAt:               grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		LastUpdatedAt:           grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		ArchivedAt:              grpcconverters.ConvertPBTimestampToTimePointer(input.ArchivedAt),
		ID:                      input.ID,
		Notes:                   input.Notes,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		ByUser:                  input.ByUser,
		Rank:                    uint8(input.Rank),
		Abstain:                 input.Abstain,
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

func ConvertGRPCUserIngredientPreferenceToUserIngredientPreference(input *mealplanningsvc.UserIngredientPreference) *mealplanning.UserIngredientPreference {
	return &mealplanning.UserIngredientPreference{
		CreatedAt:     grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		ArchivedAt:    grpcconverters.ConvertPBTimestampToTimePointer(input.ArchivedAt),
		Ingredient:    *ConvertGRPCValidIngredientToValidIngredient(input.Ingredient),
		ID:            input.ID,
		Notes:         input.Notes,
		BelongsToUser: input.BelongsToUser,
		Rating:        int8(input.Rating),
		Allergy:       input.Allergy,
	}
}

func ConvertMealCreationRequestInputToGRPCMealCreationRequestInput(input *mealplanning.MealCreationRequestInput) *mealplanningsvc.MealCreationRequestInput {
	var components []*mealplanningsvc.MealComponentCreationRequestInput
	for _, component := range input.Components {
		components = append(components, ConvertMealComponentCreationRequestInputToGRPCMealComponentCreationRequestInput(component))
	}

	return &mealplanningsvc.MealCreationRequestInput{
		EstimatedPortions: &grpctypes.Float32RangeWithOptionalMax{
			Min: input.EstimatedPortions.Min,
			Max: input.EstimatedPortions.Max,
		},
		Name:                 input.Name,
		Description:          input.Description,
		Components:           components,
		EligibleForMealPlans: input.EligibleForMealPlans,
	}
}

func ConvertMealComponentCreationRequestInputToGRPCMealComponentCreationRequestInput(input *mealplanning.MealComponentCreationRequestInput) *mealplanningsvc.MealComponentCreationRequestInput {
	return &mealplanningsvc.MealComponentCreationRequestInput{
		RecipeID:      input.RecipeID,
		ComponentType: input.ComponentType,
		RecipeScale:   input.RecipeScale,
	}
}

func ConvertGRPCMealCreationRequestInputToMealCreationRequestInput(input *mealplanningsvc.MealCreationRequestInput) *mealplanning.MealCreationRequestInput {
	var components []*mealplanning.MealComponentCreationRequestInput
	for _, component := range input.Components {
		components = append(components, ConvertGRPCMealComponentCreationRequestInputToMealComponentCreationRequestInput(component))
	}

	return &mealplanning.MealCreationRequestInput{
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: input.EstimatedPortions.Min,
			Max: input.EstimatedPortions.Max,
		},
		Name:                 input.Name,
		Description:          input.Description,
		Components:           components,
		EligibleForMealPlans: input.EligibleForMealPlans,
	}
}

func ConvertGRPCMealComponentCreationRequestInputToMealComponentCreationRequestInput(input *mealplanningsvc.MealComponentCreationRequestInput) *mealplanning.MealComponentCreationRequestInput {
	return &mealplanning.MealComponentCreationRequestInput{
		RecipeID:      input.RecipeID,
		ComponentType: input.ComponentType,
		RecipeScale:   input.RecipeScale,
	}
}

func ConvertMealPlanCreationRequestInputToGRPCMealPlanCreationRequestInput(input *mealplanning.MealPlanCreationRequestInput) *mealplanningsvc.MealPlanCreationRequestInput {
	var events []*mealplanningsvc.MealPlanEventCreationRequestInput
	for _, event := range input.Events {
		events = append(events, ConvertMealPlanEventCreationRequestInputToGRPCMealPlanEventCreationRequestInput(event))
	}

	return &mealplanningsvc.MealPlanCreationRequestInput{
		VotingDeadline: grpcconverters.ConvertTimeToPBTimestamp(input.VotingDeadline),
		Notes:          input.Notes,
		ElectionMethod: input.ElectionMethod,
		Events:         events,
	}
}

func ConvertMealPlanEventCreationRequestInputToGRPCMealPlanEventCreationRequestInput(input *mealplanning.MealPlanEventCreationRequestInput) *mealplanningsvc.MealPlanEventCreationRequestInput {
	var options []*mealplanningsvc.MealPlanOptionCreationRequestInput
	for _, option := range input.Options {
		options = append(options, ConvertMealPlanOptionCreationRequestInputToGRPCMealPlanOptionCreationRequestInput(option))
	}

	return &mealplanningsvc.MealPlanEventCreationRequestInput{
		EndsAt:   grpcconverters.ConvertTimeToPBTimestamp(input.EndsAt),
		StartsAt: grpcconverters.ConvertTimeToPBTimestamp(input.StartsAt),
		Notes:    input.Notes,
		MealName: input.MealName,
		Options:  options,
	}
}

func ConvertMealPlanOptionCreationRequestInputToGRPCMealPlanOptionCreationRequestInput(input *mealplanning.MealPlanOptionCreationRequestInput) *mealplanningsvc.MealPlanOptionCreationRequestInput {
	return &mealplanningsvc.MealPlanOptionCreationRequestInput{
		AssignedCook:       input.AssignedCook,
		AssignedDishwasher: input.AssignedDishwasher,
		MealID:             input.MealID,
		Notes:              input.Notes,
		MealScale:          input.MealScale,
	}
}

func ConvertGRPCMealPlanCreationRequestInputToMealPlanCreationRequestInput(input *mealplanningsvc.MealPlanCreationRequestInput) *mealplanning.MealPlanCreationRequestInput {
	var events []*mealplanning.MealPlanEventCreationRequestInput
	for _, event := range input.Events {
		events = append(events, ConvertGRPCMealPlanEventCreationRequestInputToMealPlanEventCreationRequestInput(event))
	}

	return &mealplanning.MealPlanCreationRequestInput{
		VotingDeadline: grpcconverters.ConvertPBTimestampToTime(input.VotingDeadline),
		Notes:          input.Notes,
		ElectionMethod: input.ElectionMethod,
		Events:         events,
	}
}

func ConvertGRPCMealPlanEventCreationRequestInputToMealPlanEventCreationRequestInput(input *mealplanningsvc.MealPlanEventCreationRequestInput) *mealplanning.MealPlanEventCreationRequestInput {
	var options []*mealplanning.MealPlanOptionCreationRequestInput
	for _, option := range input.Options {
		options = append(options, ConvertGRPCMealPlanOptionCreationRequestInputToMealPlanOptionCreationRequestInput(option))
	}

	return &mealplanning.MealPlanEventCreationRequestInput{
		EndsAt:   grpcconverters.ConvertPBTimestampToTime(input.EndsAt),
		StartsAt: grpcconverters.ConvertPBTimestampToTime(input.StartsAt),
		Notes:    input.Notes,
		MealName: input.MealName,
		Options:  options,
	}
}

func ConvertGRPCMealPlanOptionCreationRequestInputToMealPlanOptionCreationRequestInput(input *mealplanningsvc.MealPlanOptionCreationRequestInput) *mealplanning.MealPlanOptionCreationRequestInput {
	return &mealplanning.MealPlanOptionCreationRequestInput{
		AssignedCook:       input.AssignedCook,
		AssignedDishwasher: input.AssignedDishwasher,
		MealID:             input.MealID,
		Notes:              input.Notes,
		MealScale:          input.MealScale,
	}
}

func ConvertGRPCMealPlanOptionVoteCreationRequestInputToMealPlanOptionVoteCreationRequestInput(input *mealplanningsvc.MealPlanOptionVoteCreationRequestInput) *mealplanning.MealPlanOptionVoteCreationRequestInput {
	var votes []*mealplanning.MealPlanOptionVoteCreationInput
	for _, vote := range input.Votes {
		votes = append(votes, ConvertGRPCMealPlanOptionVoteCreationInputToMealPlanOptionVoteCreationInput(vote))
	}

	return &mealplanning.MealPlanOptionVoteCreationRequestInput{
		Votes: votes,
	}
}

func ConvertGRPCMealPlanOptionVoteCreationInputToMealPlanOptionVoteCreationInput(input *mealplanningsvc.MealPlanOptionVoteCreationInput) *mealplanning.MealPlanOptionVoteCreationInput {
	return &mealplanning.MealPlanOptionVoteCreationInput{
		ID:                      input.ID,
		Notes:                   input.Notes,
		ByUser:                  input.ByUser,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		Rank:                    uint8(input.Rank),
		Abstain:                 input.Abstain,
	}
}

func ConvertGRPCMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItemCreationRequestInput(input *mealplanningsvc.MealPlanGroceryListItemCreationRequestInput) *mealplanning.MealPlanGroceryListItemCreationRequestInput {
	return &mealplanning.MealPlanGroceryListItemCreationRequestInput{
		PurchasedMeasurementUnitID: input.PurchasedMeasurementUnitID,
		PurchasedUPC:               input.PurchasedUPC,
		PurchasePrice:              input.PurchasePrice,
		QuantityPurchased:          input.QuantityPurchased,
		Status:                     input.Status,
		BelongsToMealPlan:          input.BelongsToMealPlan,
		ValidIngredientID:          input.ValidIngredientID,
		ValidMeasurementUnitID:     input.ValidMeasurementUnitID,
		StatusExplanation:          input.StatusExplanation,
		QuantityNeeded:             types.Float32RangeWithOptionalMax{},
	}
}

func ConvertGRPCMealPlanTaskCreationRequestInputToMealPlanTaskCreationRequestInput(input *mealplanningsvc.MealPlanTaskCreationRequestInput) *mealplanning.MealPlanTaskCreationRequestInput {
	return &mealplanning.MealPlanTaskCreationRequestInput{
		AssignedToUser:      input.AssignedToUser,
		Status:              input.Status,
		CreationExplanation: input.CreationExplanation,
		StatusExplanation:   input.StatusExplanation,
		MealPlanOptionID:    input.MealPlanOptionID,
		RecipePrepTaskID:    input.RecipePrepTaskID,
	}
}

func ConvertGRPCUserIngredientPreferenceCreationRequestInputToUserIngredientPreferenceCreationRequestInput(input *mealplanningsvc.UserIngredientPreferenceCreationRequestInput) *mealplanning.UserIngredientPreferenceCreationRequestInput {
	return &mealplanning.UserIngredientPreferenceCreationRequestInput{
		ValidIngredientGroupID: input.ValidIngredientGroupID,
		ValidIngredientID:      input.ValidIngredientID,
		Notes:                  input.Notes,
		Rating:                 int8(input.Rating),
		Allergy:                input.Allergy,
	}
}

func ConvertUserIngredientPreferenceCreationRequestInputToGRPCUserIngredientPreferenceCreationRequestInput(input *mealplanning.UserIngredientPreferenceCreationRequestInput) *mealplanningsvc.UserIngredientPreferenceCreationRequestInput {
	return &mealplanningsvc.UserIngredientPreferenceCreationRequestInput{
		ValidIngredientGroupID: input.ValidIngredientGroupID,
		ValidIngredientID:      input.ValidIngredientID,
		Notes:                  input.Notes,
		Rating:                 int32(input.Rating),
		Allergy:                input.Allergy,
	}
}

func ConvertGRPCMealPlanUpdateRequestInputToMealPlanUpdateRequestInput(input *mealplanningsvc.MealPlanUpdateRequestInput) *mealplanning.MealPlanUpdateRequestInput {
	return &mealplanning.MealPlanUpdateRequestInput{
		BelongsToAccount: input.BelongsToAccount,
		Notes:            input.Notes,
		VotingDeadline:   grpcconverters.ConvertPBTimestampToTimePointer(input.VotingDeadline),
	}
}

func ConvertGRPCMealPlanEventUpdateRequestInputToMealPlanEventUpdateRequestInput(input *mealplanningsvc.MealPlanEventUpdateRequestInput) *mealplanning.MealPlanEventUpdateRequestInput {
	return &mealplanning.MealPlanEventUpdateRequestInput{
		Notes:             input.Notes,
		StartsAt:          grpcconverters.ConvertPBTimestampToTimePointer(input.StartsAt),
		MealName:          input.MealName,
		EndsAt:            grpcconverters.ConvertPBTimestampToTimePointer(input.EndsAt),
		BelongsToMealPlan: input.BelongsToMealPlan,
	}
}

func ConvertGRPCMealPlanGroceryListItemUpdateRequestInputToMealPlanGroceryListItemUpdateRequestInput(input *mealplanningsvc.MealPlanGroceryListItemUpdateRequestInput) *mealplanning.MealPlanGroceryListItemUpdateRequestInput {
	return &mealplanning.MealPlanGroceryListItemUpdateRequestInput{
		BelongsToMealPlan:          input.BelongsToMealPlan,
		ValidIngredientID:          input.ValidIngredientID,
		ValidMeasurementUnitID:     input.ValidMeasurementUnitID,
		StatusExplanation:          input.StatusExplanation,
		QuantityPurchased:          input.QuantityPurchased,
		PurchasedMeasurementUnitID: input.PurchasedMeasurementUnitID,
		PurchasedUPC:               input.PurchasedUPC,
		PurchasePrice:              input.PurchasePrice,
		Status:                     input.Status,
		QuantityNeeded: types.Float32RangeWithOptionalMaxUpdateRequestInput{
			Min: input.QuantityNeeded.Min,
			Max: input.QuantityNeeded.Max,
		},
	}
}

func ConvertGRPCMealPlanOptionUpdateRequestInputToMealPlanOptionUpdateRequestInput(input *mealplanningsvc.MealPlanOptionUpdateRequestInput) *mealplanning.MealPlanOptionUpdateRequestInput {
	return &mealplanning.MealPlanOptionUpdateRequestInput{
		MealID:                 input.MealID,
		Notes:                  input.Notes,
		AssignedCook:           input.AssignedCook,
		AssignedDishwasher:     input.AssignedDishwasher,
		MealScale:              input.MealScale,
		BelongsToMealPlanEvent: input.BelongsToMealPlanEvent,
	}
}

func ConvertMealPlanOptionUpdateRequestInputToGRPCMealPlanOptionUpdateRequestInput(input *mealplanning.MealPlanOptionUpdateRequestInput) *mealplanningsvc.MealPlanOptionUpdateRequestInput {
	return &mealplanningsvc.MealPlanOptionUpdateRequestInput{
		MealID:                 input.MealID,
		Notes:                  input.Notes,
		AssignedCook:           input.AssignedCook,
		AssignedDishwasher:     input.AssignedDishwasher,
		MealScale:              input.MealScale,
		BelongsToMealPlanEvent: input.BelongsToMealPlanEvent,
	}
}

func ConvertGRPCMealPlanOptionVoteUpdateRequestInputToMealPlanOptionVoteUpdateRequestInput(input *mealplanningsvc.MealPlanOptionVoteUpdateRequestInput) *mealplanning.MealPlanOptionVoteUpdateRequestInput {
	return &mealplanning.MealPlanOptionVoteUpdateRequestInput{
		Notes:                   input.Notes,
		Rank:                    pointer.To(uint8(pointer.Dereference(input.Rank))),
		Abstain:                 input.Abstain,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
	}
}

func ConvertGRPCMealPlanTaskStatusChangeRequestInputToMealPlanTaskStatusChangeRequestInput(input *mealplanningsvc.MealPlanTaskStatusChangeRequestInput) *mealplanning.MealPlanTaskStatusChangeRequestInput {
	return &mealplanning.MealPlanTaskStatusChangeRequestInput{
		Status:            input.Status,
		StatusExplanation: input.StatusExplanation,
		AssignedToUser:    input.AssignedToUser,
		ID:                input.ID,
	}
}

func ConvertGRPCUserIngredientPreferenceUpdateRequestInputToUserIngredientPreferenceUpdateRequestInput(input *mealplanningsvc.UserIngredientPreferenceUpdateRequestInput) *mealplanning.UserIngredientPreferenceUpdateRequestInput {
	return &mealplanning.UserIngredientPreferenceUpdateRequestInput{
		Notes:        input.Notes,
		IngredientID: input.IngredientID,
		Rating:       pointer.To(int8(pointer.Dereference(input.Rating))),
		Allergy:      input.Allergy,
	}
}

func ConvertAccountInstrumentOwnershipCreationRequestInputToGRPCAccountInstrumentOwnershipCreationRequestInput(input *mealplanning.AccountInstrumentOwnershipCreationRequestInput) *mealplanningsvc.AccountInstrumentOwnershipCreationRequestInput {
	return &mealplanningsvc.AccountInstrumentOwnershipCreationRequestInput{
		Notes:             input.Notes,
		ValidInstrumentID: input.ValidInstrumentID,
		Quantity:          uint32(input.Quantity),
	}
}

func ConvertGRPCAccountInstrumentOwnershipToAccountInstrumentOwnership(input *mealplanningsvc.AccountInstrumentOwnership) *mealplanning.AccountInstrumentOwnership {
	return &mealplanning.AccountInstrumentOwnership{
		CreatedAt:        grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertPBTimestampToTimePointer(input.ArchivedAt),
		LastUpdatedAt:    grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		ID:               input.ID,
		Notes:            input.Notes,
		BelongsToAccount: input.BelongsToAccount,
		Instrument:       *ConvertGRPCValidInstrumentToValidInstrument(input.Instrument),
		Quantity:         uint16(input.Quantity),
	}
}

func ConvertGRPCAccountInstrumentOwnershipCreationRequestInputToAccountInstrumentOwnershipCreationRequestInput(input *mealplanningsvc.AccountInstrumentOwnershipCreationRequestInput) *mealplanning.AccountInstrumentOwnershipCreationRequestInput {
	return &mealplanning.AccountInstrumentOwnershipCreationRequestInput{
		Notes:             input.Notes,
		ValidInstrumentID: input.ValidInstrumentID,
		Quantity:          uint16(input.Quantity),
	}
}

func ConvertAccountInstrumentOwnershipToGRPCAccountInstrumentOwnership(input *mealplanning.AccountInstrumentOwnership) *mealplanningsvc.AccountInstrumentOwnership {
	return &mealplanningsvc.AccountInstrumentOwnership{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		Instrument:       ConvertValidInstrumentToGRPCValidInstrument(&input.Instrument),
		ID:               input.ID,
		Notes:            input.Notes,
		BelongsToAccount: input.BelongsToAccount,
		Quantity:         uint32(input.Quantity),
	}
}

func ConvertGRPCAccountInstrumentOwnershipUpdateRequestInputToAccountInstrumentOwnershipUpdateRequestInput(input *mealplanningsvc.AccountInstrumentOwnershipUpdateRequestInput) *mealplanning.AccountInstrumentOwnershipUpdateRequestInput {
	var quantity *uint16
	if input.Quantity != nil {
		quantity = pointer.To(uint16(*input.Quantity))
	}

	return &mealplanning.AccountInstrumentOwnershipUpdateRequestInput{
		Notes:             input.Notes,
		Quantity:          quantity,
		ValidInstrumentID: input.ValidInstrumentID,
	}
}

func ConvertMealPlanOptionVoteCreationRequestInputToGRPCMealPlanOptionVoteCreationRequestInput(input *mealplanning.MealPlanOptionVoteCreationRequestInput) *mealplanningsvc.MealPlanOptionVoteCreationRequestInput {
	var votes []*mealplanningsvc.MealPlanOptionVoteCreationInput
	for _, vote := range input.Votes {
		votes = append(votes, ConvertMealPlanOptionVoteCreationInputToGRPCMealPlanOptionVoteCreationInput(vote))
	}

	return &mealplanningsvc.MealPlanOptionVoteCreationRequestInput{Votes: votes}
}

func ConvertMealPlanOptionVoteCreationInputToGRPCMealPlanOptionVoteCreationInput(input *mealplanning.MealPlanOptionVoteCreationInput) *mealplanningsvc.MealPlanOptionVoteCreationInput {
	return &mealplanningsvc.MealPlanOptionVoteCreationInput{
		ID:                      input.ID,
		Notes:                   input.Notes,
		ByUser:                  input.ByUser,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		Rank:                    uint32(input.Rank),
		Abstain:                 input.Abstain,
	}
}
