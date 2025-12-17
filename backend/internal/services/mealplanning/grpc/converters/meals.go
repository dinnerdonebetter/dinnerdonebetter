package grpcconverters

import (
	"log"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	grpctypes "github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertStringToMealPlanTaskStatus(s string) mealplanningsvc.MealPlanTaskStatus {
	switch s {
	case mealplanning.MealPlanTaskStatusPostponed:
		return mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_POSTPONED
	case mealplanning.MealPlanTaskStatusIgnored:
		return mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_IGNORED
	case mealplanning.MealPlanTaskStatusCanceled:
		return mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_CANCELED
	case mealplanning.MealPlanTaskStatusFinished:
		return mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_FINISHED
	case mealplanning.MealPlanTaskStatusUnfinished:
		return mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_UNFINISHED
	default:
		log.Printf("UNKNOWN MEALPLAN TASK_STATUS: %s", s)
		return mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_UNFINISHED
	}
}

func ConvertMealPlanTaskStatusToString(s mealplanningsvc.MealPlanTaskStatus) string {
	switch s {
	case mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_POSTPONED:
		return mealplanning.MealPlanTaskStatusPostponed
	case mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_IGNORED:
		return mealplanning.MealPlanTaskStatusIgnored
	case mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_CANCELED:
		return mealplanning.MealPlanTaskStatusCanceled
	case mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_FINISHED:
		return mealplanning.MealPlanTaskStatusFinished
	case mealplanningsvc.MealPlanTaskStatus_MEAL_PLAN_TASK_STATUS_UNFINISHED:
		return mealplanning.MealPlanTaskStatusUnfinished
	default:
		log.Printf("UNKNOWN MEALPLAN TASK_STATUS: %s", s)
		return mealplanning.MealPlanTaskStatusUnfinished
	}
}

func ConvertMealPlanTaskToGRPCMealPlanTask(input *mealplanning.MealPlanTask) *mealplanningsvc.MealPlanTask {
	return &mealplanningsvc.MealPlanTask{
		RecipePrepTask:      ConvertRecipePrepTaskToGRPCRecipePrepTask(&input.RecipePrepTask),
		CreatedAt:           grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		CompletedAt:         grpcconverters.ConvertTimePointerToPBTimestamp(input.CompletedAt),
		AssignedToUser:      input.AssignedToUser,
		Id:                  input.ID,
		Status:              ConvertStringToMealPlanTaskStatus(input.Status),
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
		Id:                     input.ID,
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
		Id:                      input.ID,
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
		Id:                   input.ID,
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
		ComponentType: ConvertStringToMealComponentType(input.ComponentType),
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
		ID:                   input.Id,
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
		ComponentType: ConvertMealComponentTypeToString(input.ComponentType),
		RecipeScale:   input.RecipeScale,
	}
}

func ConvertStringToMealPlanGroceryListItemStatus(s string) mealplanningsvc.MealPlanGroceryListItemStatus {
	switch s {
	case mealplanning.MealPlanGroceryListItemStatusAlreadyOwned:
		return mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_ALREADY_OWNED
	case mealplanning.MealPlanGroceryListItemStatusNeeds:
		return mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_NEEDS
	case mealplanning.MealPlanGroceryListItemStatusUnavailable:
		return mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_UNAVAILABLE
	case mealplanning.MealPlanGroceryListItemStatusAcquired:
		return mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_ACQUIRED
	case mealplanning.MealPlanGroceryListItemStatusUnknown:
		return mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_UNKNOWN
	default:
		log.Printf("UNKNOWN MealPlanGroceryListItemStatus: %q", s)
		return mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_UNKNOWN
	}
}

func ConvertMealPlanGroceryListItemStatusToString(s mealplanningsvc.MealPlanGroceryListItemStatus) string {
	switch s {
	case mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_ALREADY_OWNED:
		return mealplanning.MealPlanGroceryListItemStatusAlreadyOwned
	case mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_NEEDS:
		return mealplanning.MealPlanGroceryListItemStatusNeeds
	case mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_UNAVAILABLE:
		return mealplanning.MealPlanGroceryListItemStatusUnavailable
	case mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_ACQUIRED:
		return mealplanning.MealPlanGroceryListItemStatusAcquired
	case mealplanningsvc.MealPlanGroceryListItemStatus_MEAL_PLAN_GROCERY_LIST_ITEM_STATUS_UNKNOWN:
		return mealplanning.MealPlanGroceryListItemStatusUnknown
	default:
		log.Printf("UNKNOWN MealPlanGroceryListItemStatus: %q", s)
		return mealplanning.MealPlanGroceryListItemStatusUnknown
	}
}

func ConvertMealPlanGroceryListItemToGRPCMealPlanGroceryListItem(input *mealplanning.MealPlanGroceryListItem) *mealplanningsvc.MealPlanGroceryListItem {
	var purchasedMeasurementUnit *mealplanningsvc.ValidMeasurementUnit
	if input.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnit = ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(input.PurchasedMeasurementUnit)
	}

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
		PurchasedMeasurementUnit: purchasedMeasurementUnit,
		PurchasedUpc:             input.PurchasedUPC,
		Status:                   ConvertStringToMealPlanGroceryListItemStatus(input.Status),
		StatusExplanation:        input.StatusExplanation,
		Id:                       input.ID,
		BelongsToMealPlan:        input.BelongsToMealPlan,
		PurchasePrice:            input.PurchasePrice,
		QuantityPurchased:        input.QuantityPurchased,
	}
}

func ConvertGRPCMealPlanGroceryListItemToMealPlanGroceryListItem(input *mealplanningsvc.MealPlanGroceryListItem) *mealplanning.MealPlanGroceryListItem {
	var purchasedMeasurementUnit *mealplanning.ValidMeasurementUnit
	if input.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnit = ConvertGRPCValidMeasurementUnitToValidMeasurementUnit(input.PurchasedMeasurementUnit)
	}

	return &mealplanning.MealPlanGroceryListItem{
		CreatedAt:     grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		ArchivedAt:    grpcconverters.ConvertPBTimestampToTimePointer(input.ArchivedAt),
		QuantityNeeded: types.Float32RangeWithOptionalMax{
			Max: input.QuantityNeeded.Max,
			Min: input.QuantityNeeded.Min,
		},
		Ingredient:               *ConvertGRPCValidIngredientToValidIngredient(input.Ingredient),
		MeasurementUnit:          *ConvertGRPCValidMeasurementUnitToValidMeasurementUnit(input.MeasurementUnit),
		PurchasedMeasurementUnit: purchasedMeasurementUnit,
		PurchasedUPC:             input.PurchasedUpc,
		Status:                   ConvertMealPlanGroceryListItemStatusToString(input.Status),
		StatusExplanation:        input.StatusExplanation,
		ID:                       input.Id,
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
		Id:                input.ID,
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
		ElectionMethod:         ConvertStringToMealPlanElectionMethod(input.ElectionMethod),
		Status:                 input.Status,
		Notes:                  input.Notes,
		Id:                     input.ID,
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
		ElectionMethod:         ConvertMealPlanElectionMethodToString(input.ElectionMethod),
		Status:                 input.Status,
		Notes:                  input.Notes,
		ID:                     input.Id,
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
		ID:                input.Id,
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
		ID:                     input.Id,
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
		ID:                      input.Id,
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
		Id:            input.ID,
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
		ID:            input.Id,
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

func ConvertStringToMealComponentType(s string) mealplanningsvc.MealComponentType {
	switch s {
	case mealplanning.MealComponentTypesAmuseBouche:
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_AMUSE_BOUCHE
	case mealplanning.MealComponentTypesAppetizer:
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_APPETIZER
	case mealplanning.MealComponentTypesSoup:
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_SOUP
	case mealplanning.MealComponentTypesMain:
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_MAIN
	case mealplanning.MealComponentTypesSalad:
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_SALAD
	case mealplanning.MealComponentTypesBeverage:
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_BEVERAGE
	case mealplanning.MealComponentTypesSide:
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_SIDE
	case mealplanning.MealComponentTypesDessert:
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_DESSERT
	case mealplanning.MealComponentTypesUnspecified:
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_UNSPECIFIED
	default:
		log.Printf("UNKNOWN MEAL COMPONENT TYPE: %q", s)
		return mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_UNSPECIFIED
	}
}

func ConvertMealComponentTypeToString(s mealplanningsvc.MealComponentType) string {
	switch s {
	case mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_AMUSE_BOUCHE:
		return mealplanning.MealComponentTypesAmuseBouche
	case mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_APPETIZER:
		return mealplanning.MealComponentTypesAppetizer
	case mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_SOUP:
		return mealplanning.MealComponentTypesSoup
	case mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_MAIN:
		return mealplanning.MealComponentTypesMain
	case mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_SALAD:
		return mealplanning.MealComponentTypesSalad
	case mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_BEVERAGE:
		return mealplanning.MealComponentTypesBeverage
	case mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_SIDE:
		return mealplanning.MealComponentTypesSide
	case mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_DESSERT:
		return mealplanning.MealComponentTypesDessert
	case mealplanningsvc.MealComponentType_MEAL_COMPONENT_TYPE_UNSPECIFIED:
		return mealplanning.MealComponentTypesUnspecified
	default:
		log.Printf("UNKNOWN MEAL COMPONENT TYPE: %q", s)
		return mealplanning.MealComponentTypesUnspecified
	}
}

func ConvertMealComponentCreationRequestInputToGRPCMealComponentCreationRequestInput(input *mealplanning.MealComponentCreationRequestInput) *mealplanningsvc.MealComponentCreationRequestInput {
	return &mealplanningsvc.MealComponentCreationRequestInput{
		RecipeId:      input.RecipeID,
		ComponentType: ConvertStringToMealComponentType(input.ComponentType),
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
		RecipeID:      input.RecipeId,
		ComponentType: ConvertMealComponentTypeToString(input.ComponentType),
		RecipeScale:   input.RecipeScale,
	}
}

func ConvertStringToMealPlanElectionMethod(s string) mealplanningsvc.MealPlanElectionMethod {
	switch s {
	case mealplanning.MealPlanElectionMethodSchulze:
		return mealplanningsvc.MealPlanElectionMethod_MEAL_PLAN_ELECTION_METHOD_SCHULZE
	case mealplanning.MealPlanElectionMethodInstantRunoff:
		return mealplanningsvc.MealPlanElectionMethod_MEAL_PLAN_ELECTION_METHOD_INSTANT_RUNOFF
	default:
		log.Printf("UNKNOWN MEAL COMPONENT TYPE: %q", s)
		return mealplanningsvc.MealPlanElectionMethod_MEAL_PLAN_ELECTION_METHOD_SCHULZE
	}
}

func ConvertMealPlanElectionMethodToString(s mealplanningsvc.MealPlanElectionMethod) string {
	switch s {
	case mealplanningsvc.MealPlanElectionMethod_MEAL_PLAN_ELECTION_METHOD_SCHULZE:
		return mealplanning.MealPlanElectionMethodSchulze
	case mealplanningsvc.MealPlanElectionMethod_MEAL_PLAN_ELECTION_METHOD_INSTANT_RUNOFF:
		return mealplanning.MealPlanElectionMethodInstantRunoff
	default:
		log.Printf("UNKNOWN MEAL COMPONENT TYPE: %q", s)
		return mealplanning.MealPlanElectionMethodSchulze
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
		ElectionMethod: ConvertStringToMealPlanElectionMethod(input.ElectionMethod),
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
		MealId:             input.MealID,
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
		ElectionMethod: ConvertMealPlanElectionMethodToString(input.ElectionMethod),
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
		MealID:             input.MealId,
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
		ID:                      input.Id,
		Notes:                   input.Notes,
		ByUser:                  input.ByUser,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		Rank:                    uint8(input.Rank),
		Abstain:                 input.Abstain,
	}
}

func ConvertGRPCMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItemCreationRequestInput(input *mealplanningsvc.MealPlanGroceryListItemCreationRequestInput) *mealplanning.MealPlanGroceryListItemCreationRequestInput {
	return &mealplanning.MealPlanGroceryListItemCreationRequestInput{
		PurchasedMeasurementUnitID: input.PurchasedMeasurementUnitId,
		PurchasedUPC:               input.PurchasedUpc,
		PurchasePrice:              input.PurchasePrice,
		QuantityPurchased:          input.QuantityPurchased,
		Status:                     ConvertMealPlanGroceryListItemStatusToString(input.Status),
		BelongsToMealPlan:          input.BelongsToMealPlan,
		ValidIngredientID:          input.ValidIngredientId,
		ValidMeasurementUnitID:     input.ValidMeasurementUnitId,
		StatusExplanation:          input.StatusExplanation,
		QuantityNeeded: types.Float32RangeWithOptionalMax{
			Max: input.QuantityNeeded.Max,
			Min: input.QuantityNeeded.Min,
		},
	}
}

func ConvertMealPlanGroceryListItemCreationRequestInputToGRPCMealPlanGroceryListItemCreationRequestInput(input *mealplanning.MealPlanGroceryListItemCreationRequestInput) *mealplanningsvc.MealPlanGroceryListItemCreationRequestInput {
	return &mealplanningsvc.MealPlanGroceryListItemCreationRequestInput{
		PurchasedMeasurementUnitId: input.PurchasedMeasurementUnitID,
		PurchasedUpc:               input.PurchasedUPC,
		PurchasePrice:              input.PurchasePrice,
		QuantityPurchased:          input.QuantityPurchased,
		Status:                     ConvertStringToMealPlanGroceryListItemStatus(input.Status),
		BelongsToMealPlan:          input.BelongsToMealPlan,
		ValidIngredientId:          input.ValidIngredientID,
		ValidMeasurementUnitId:     input.ValidMeasurementUnitID,
		StatusExplanation:          input.StatusExplanation,
		QuantityNeeded: &grpctypes.Float32RangeWithOptionalMax{
			Max: input.QuantityNeeded.Max,
			Min: input.QuantityNeeded.Min,
		},
	}
}

func ConvertGRPCMealPlanTaskCreationRequestInputToMealPlanTaskCreationRequestInput(input *mealplanningsvc.MealPlanTaskCreationRequestInput) *mealplanning.MealPlanTaskCreationRequestInput {
	return &mealplanning.MealPlanTaskCreationRequestInput{
		AssignedToUser:      input.AssignedToUser,
		Status:              ConvertMealPlanTaskStatusToString(input.Status),
		CreationExplanation: input.CreationExplanation,
		StatusExplanation:   input.StatusExplanation,
		MealPlanOptionID:    input.MealPlanOptionId,
		RecipePrepTaskID:    input.RecipePrepTaskId,
	}
}

func ConvertGRPCUserIngredientPreferenceCreationRequestInputToUserIngredientPreferenceCreationRequestInput(input *mealplanningsvc.UserIngredientPreferenceCreationRequestInput) *mealplanning.UserIngredientPreferenceCreationRequestInput {
	return &mealplanning.UserIngredientPreferenceCreationRequestInput{
		ValidIngredientGroupID: input.ValidIngredientGroupId,
		ValidIngredientID:      input.ValidIngredientId,
		Notes:                  input.Notes,
		Rating:                 int8(input.Rating),
		Allergy:                input.Allergy,
	}
}

func ConvertUserIngredientPreferenceCreationRequestInputToGRPCUserIngredientPreferenceCreationRequestInput(input *mealplanning.UserIngredientPreferenceCreationRequestInput) *mealplanningsvc.UserIngredientPreferenceCreationRequestInput {
	return &mealplanningsvc.UserIngredientPreferenceCreationRequestInput{
		ValidIngredientGroupId: input.ValidIngredientGroupID,
		ValidIngredientId:      input.ValidIngredientID,
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

func ConvertMealPlanEventUpdateRequestInputToGRPCMealPlanEventUpdateRequestInput(input *mealplanning.MealPlanEventUpdateRequestInput) *mealplanningsvc.MealPlanEventUpdateRequestInput {
	var startsAt, endsAt *timestamppb.Timestamp
	if input.StartsAt != nil {
		startsAt = grpcconverters.ConvertTimeToPBTimestamp(*input.StartsAt)
	}

	if input.EndsAt != nil {
		endsAt = grpcconverters.ConvertTimeToPBTimestamp(*input.EndsAt)
	}

	return &mealplanningsvc.MealPlanEventUpdateRequestInput{
		Notes:             input.Notes,
		StartsAt:          startsAt,
		MealName:          input.MealName,
		EndsAt:            endsAt,
		BelongsToMealPlan: input.BelongsToMealPlan,
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
	var status *string
	if input.Status != nil {
		status = pointer.To(ConvertMealPlanGroceryListItemStatusToString(*input.Status))
	}

	return &mealplanning.MealPlanGroceryListItemUpdateRequestInput{
		BelongsToMealPlan:          input.BelongsToMealPlan,
		ValidIngredientID:          input.ValidIngredientId,
		ValidMeasurementUnitID:     input.ValidMeasurementUnitId,
		StatusExplanation:          input.StatusExplanation,
		QuantityPurchased:          input.QuantityPurchased,
		PurchasedMeasurementUnitID: input.PurchasedMeasurementUnitId,
		PurchasedUPC:               input.PurchasedUpc,
		PurchasePrice:              input.PurchasePrice,
		Status:                     status,
		QuantityNeeded: types.Float32RangeWithOptionalMaxUpdateRequestInput{
			Min: input.QuantityNeeded.Min,
			Max: input.QuantityNeeded.Max,
		},
	}
}

func ConvertGRPCMealPlanOptionUpdateRequestInputToMealPlanOptionUpdateRequestInput(input *mealplanningsvc.MealPlanOptionUpdateRequestInput) *mealplanning.MealPlanOptionUpdateRequestInput {
	return &mealplanning.MealPlanOptionUpdateRequestInput{
		MealID:                 input.MealId,
		Notes:                  input.Notes,
		AssignedCook:           input.AssignedCook,
		AssignedDishwasher:     input.AssignedDishwasher,
		MealScale:              input.MealScale,
		BelongsToMealPlanEvent: input.BelongsToMealPlanEvent,
	}
}

func ConvertMealPlanOptionUpdateRequestInputToGRPCMealPlanOptionUpdateRequestInput(input *mealplanning.MealPlanOptionUpdateRequestInput) *mealplanningsvc.MealPlanOptionUpdateRequestInput {
	return &mealplanningsvc.MealPlanOptionUpdateRequestInput{
		MealId:                 input.MealID,
		Notes:                  input.Notes,
		AssignedCook:           input.AssignedCook,
		AssignedDishwasher:     input.AssignedDishwasher,
		MealScale:              input.MealScale,
		BelongsToMealPlanEvent: input.BelongsToMealPlanEvent,
	}
}

func ConvertMealPlanOptionVoteUpdateRequestInputToGRPCMealPlanOptionVoteUpdateRequestInput(input *mealplanning.MealPlanOptionVoteUpdateRequestInput) *mealplanningsvc.MealPlanOptionVoteUpdateRequestInput {
	return &mealplanningsvc.MealPlanOptionVoteUpdateRequestInput{
		Notes:                   input.Notes,
		Rank:                    grpcconverters.ConvertUint8PointerToUint32Pointer(input.Rank),
		Abstain:                 input.Abstain,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
	}
}

func ConvertGRPCMealPlanOptionVoteUpdateRequestInputToMealPlanOptionVoteUpdateRequestInput(input *mealplanningsvc.MealPlanOptionVoteUpdateRequestInput) *mealplanning.MealPlanOptionVoteUpdateRequestInput {
	return &mealplanning.MealPlanOptionVoteUpdateRequestInput{
		Notes:                   input.Notes,
		Rank:                    grpcconverters.ConvertUint32PointerToUint8Pointer(input.Rank),
		Abstain:                 input.Abstain,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
	}
}

func ConvertGRPCMealPlanTaskStatusChangeRequestInputToMealPlanTaskStatusChangeRequestInput(input *mealplanningsvc.MealPlanTaskStatusChangeRequestInput) *mealplanning.MealPlanTaskStatusChangeRequestInput {
	var status *string
	if input.Status != nil {
		status = pointer.To(ConvertMealPlanTaskStatusToString(*input.Status))
	}

	return &mealplanning.MealPlanTaskStatusChangeRequestInput{
		Status:            status,
		StatusExplanation: input.StatusExplanation,
		AssignedToUser:    input.AssignedToUser,
		ID:                input.Id,
	}
}

func ConvertGRPCUserIngredientPreferenceUpdateRequestInputToUserIngredientPreferenceUpdateRequestInput(input *mealplanningsvc.UserIngredientPreferenceUpdateRequestInput) *mealplanning.UserIngredientPreferenceUpdateRequestInput {
	return &mealplanning.UserIngredientPreferenceUpdateRequestInput{
		Notes:        input.Notes,
		IngredientID: input.IngredientId,
		Rating:       pointer.To(int8(pointer.Dereference(input.Rating))),
		Allergy:      input.Allergy,
	}
}

func ConvertAccountInstrumentOwnershipCreationRequestInputToGRPCAccountInstrumentOwnershipCreationRequestInput(input *mealplanning.AccountInstrumentOwnershipCreationRequestInput) *mealplanningsvc.AccountInstrumentOwnershipCreationRequestInput {
	return &mealplanningsvc.AccountInstrumentOwnershipCreationRequestInput{
		Notes:             input.Notes,
		ValidInstrumentId: input.ValidInstrumentID,
		Quantity:          uint32(input.Quantity),
	}
}

func ConvertGRPCAccountInstrumentOwnershipToAccountInstrumentOwnership(input *mealplanningsvc.AccountInstrumentOwnership) *mealplanning.AccountInstrumentOwnership {
	return &mealplanning.AccountInstrumentOwnership{
		CreatedAt:        grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertPBTimestampToTimePointer(input.ArchivedAt),
		LastUpdatedAt:    grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		ID:               input.Id,
		Notes:            input.Notes,
		BelongsToAccount: input.BelongsToAccount,
		Instrument:       *ConvertGRPCValidInstrumentToValidInstrument(input.Instrument),
		Quantity:         uint16(input.Quantity),
	}
}

func ConvertGRPCAccountInstrumentOwnershipCreationRequestInputToAccountInstrumentOwnershipCreationRequestInput(input *mealplanningsvc.AccountInstrumentOwnershipCreationRequestInput) *mealplanning.AccountInstrumentOwnershipCreationRequestInput {
	return &mealplanning.AccountInstrumentOwnershipCreationRequestInput{
		Notes:             input.Notes,
		ValidInstrumentID: input.ValidInstrumentId,
		Quantity:          uint16(input.Quantity),
	}
}

func ConvertAccountInstrumentOwnershipToGRPCAccountInstrumentOwnership(input *mealplanning.AccountInstrumentOwnership) *mealplanningsvc.AccountInstrumentOwnership {
	return &mealplanningsvc.AccountInstrumentOwnership{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		Instrument:       ConvertValidInstrumentToGRPCValidInstrument(&input.Instrument),
		Id:               input.ID,
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
		ValidInstrumentID: input.ValidInstrumentId,
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
		Id:                      input.ID,
		Notes:                   input.Notes,
		ByUser:                  input.ByUser,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		Rank:                    uint32(input.Rank),
		Abstain:                 input.Abstain,
	}
}

func ConvertMealPlanTaskCreationRequestInputToGRPCMealPlanTaskCreationRequestInput(input *mealplanning.MealPlanTaskCreationRequestInput) *mealplanningsvc.MealPlanTaskCreationRequestInput {
	return &mealplanningsvc.MealPlanTaskCreationRequestInput{
		AssignedToUser:      input.AssignedToUser,
		Status:              ConvertStringToMealPlanTaskStatus(input.Status),
		CreationExplanation: input.CreationExplanation,
		StatusExplanation:   input.StatusExplanation,
		MealPlanOptionId:    input.MealPlanOptionID,
		RecipePrepTaskId:    input.RecipePrepTaskID,
	}
}

func ConvertGRPCMealPlanTaskToMealPlanTask(input *mealplanningsvc.MealPlanTask) *mealplanning.MealPlanTask {
	return &mealplanning.MealPlanTask{
		RecipePrepTask:      *ConvertGRPCRecipePrepTaskToRecipePrepTask(input.RecipePrepTask),
		CreatedAt:           grpcconverters.ConvertPBTimestampToTime(input.CreatedAt),
		CompletedAt:         grpcconverters.ConvertPBTimestampToTimePointer(input.CompletedAt),
		LastUpdatedAt:       grpcconverters.ConvertPBTimestampToTimePointer(input.LastUpdatedAt),
		AssignedToUser:      input.AssignedToUser,
		ID:                  input.Id,
		Status:              ConvertMealPlanTaskStatusToString(input.Status),
		CreationExplanation: input.CreationExplanation,
		StatusExplanation:   input.StatusExplanation,
		MealPlanOption:      *ConvertGRPCMealPlanOptionToMealPlanOption(input.MealPlanOption),
	}
}

func ConvertMealListItemToGRPCMealListItem(input *mealplanning.MealListItem) *mealplanningsvc.MealListItem {
	return &mealplanningsvc.MealListItem{
		CreatedAt:         timestamppb.New(input.CreatedAt),
		LastUpdatedAt:     grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:        grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		Id:                input.ID,
		MealId:            input.Meal.ID,
		Notes:             input.Notes,
		BelongsToMealList: input.BelongsToMealList,
		Meal:              ConvertMealToGRPCMeal(&input.Meal),
	}
}

func ConvertMealListToGRPCMealList(input *mealplanning.MealList) *mealplanningsvc.MealList {
	var items []*mealplanningsvc.MealListItem
	for _, item := range input.Items {
		items = append(items, ConvertMealListItemToGRPCMealListItem(item))
	}

	return &mealplanningsvc.MealList{
		CreatedAt:     timestamppb.New(input.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		Id:            input.ID,
		Name:          input.Name,
		Description:   input.Description,
		BelongsToUser: input.BelongsToUser,
		Items:         items,
	}
}

func ConvertGRPCMealListCreationRequestInputToMealListCreationRequestInput(input *mealplanningsvc.MealListCreationRequestInput) *mealplanning.MealListCreationRequestInput {
	return &mealplanning.MealListCreationRequestInput{
		Name:        input.Name,
		Description: input.Description,
	}
}

func ConvertGRPCMealListItemCreationRequestInputToMealListItemCreationRequestInput(input *mealplanningsvc.MealListItemCreationRequestInput) *mealplanning.MealListItemCreationRequestInput {
	return &mealplanning.MealListItemCreationRequestInput{
		MealID: input.MealId,
		Notes:  input.Notes,
	}
}

func ConvertGRPCMealListUpdateRequestInputToMealListUpdateRequestInput(input *mealplanningsvc.MealListUpdateRequestInput) *mealplanning.MealListUpdateRequestInput {
	if input == nil {
		return nil
	}

	var name *string
	if input.Name != nil {
		name = pointer.To(input.GetName())
	}

	var desc *string
	if input.Description != nil {
		desc = pointer.To(input.GetDescription())
	}

	return &mealplanning.MealListUpdateRequestInput{
		Name:        name,
		Description: desc,
	}
}

func ConvertGRPCMealListItemUpdateRequestInputToMealListItemUpdateRequestInput(input *mealplanningsvc.MealListItemUpdateRequestInput) *mealplanning.MealListItemUpdateRequestInput {
	if input == nil {
		return nil
	}

	var notes *string
	if input.Notes != nil {
		notes = pointer.To(input.GetNotes())
	}

	return &mealplanning.MealListItemUpdateRequestInput{
		Notes: notes,
	}
}
