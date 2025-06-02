package manager

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

type (
	MealPlanningManager interface {
		ListMeals(ctx context.Context, filter *filtering.QueryFilter) ([]*types.Meal, string, error)
		CreateMeal(ctx context.Context, input *types.MealCreationRequestInput) (*types.Meal, error)
		ReadMeal(ctx context.Context, mealID string) (*types.Meal, error)
		SearchMeals(ctx context.Context, query string, filter *filtering.QueryFilter) ([]*types.Meal, error)
		ArchiveMeal(ctx context.Context, mealID, ownerID string) error

		ListMealPlans(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*types.MealPlan, string, error)
		CreateMealPlan(ctx context.Context, input *types.MealPlanCreationRequestInput) (*types.MealPlan, error)
		ReadMealPlan(ctx context.Context, mealPlanID, ownerID string) (*types.MealPlan, error)
		UpdateMealPlan(ctx context.Context, mealPlanID, ownerID string, input *types.MealPlanUpdateRequestInput) error
		ArchiveMealPlan(ctx context.Context, mealPlanID, ownerID string) error
		FinalizeMealPlan(ctx context.Context, mealPlanID, ownerID string) (bool, error)

		ListMealPlanEvents(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanEvent, string, error)
		CreateMealPlanEvent(ctx context.Context, input *types.MealPlanEventCreationRequestInput) (*types.MealPlanEvent, error)
		ReadMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error)
		UpdateMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string, input *types.MealPlanEventUpdateRequestInput) error
		ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error

		ListMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *filtering.QueryFilter) ([]*types.MealPlanOption, string, error)
		CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionCreationRequestInput) (*types.MealPlanOption, error)
		ReadMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*types.MealPlanOption, error)
		UpdateMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, input *types.MealPlanOptionUpdateRequestInput) error
		ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error

		ListMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *filtering.QueryFilter) ([]*types.MealPlanOptionVote, string, error)
		CreateMealPlanOptionVotes(ctx context.Context, input *types.MealPlanOptionVoteCreationRequestInput) ([]*types.MealPlanOptionVote, error)
		ReadMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error)
		UpdateMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string, input *types.MealPlanOptionVoteUpdateRequestInput) error
		ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error

		ListMealPlanTasksByMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanTask, string, error)
		ReadMealPlanTask(ctx context.Context, mealPlanID, mealPlanTaskID string) (*types.MealPlanTask, error)
		CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskCreationRequestInput) (*types.MealPlanTask, error)
		MealPlanTaskStatusChange(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error

		ListMealPlanGroceryListItemsByMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanGroceryListItem, string, error)
		CreateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemCreationRequestInput) (*types.MealPlanGroceryListItem, error)
		ReadMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error)
		UpdateMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string, input *types.MealPlanGroceryListItemUpdateRequestInput) error
		ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) error

		ListIngredientPreferences(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*types.IngredientPreference, string, error)
		CreateIngredientPreference(ctx context.Context, input *types.IngredientPreferenceCreationRequestInput) ([]*types.IngredientPreference, error)
		UpdateIngredientPreference(ctx context.Context, ingredientPreferenceID, ownerID string, input *types.IngredientPreferenceUpdateRequestInput) error
		ArchiveIngredientPreference(ctx context.Context, ownerID, ingredientPreferenceID string) error

		ListInstrumentOwnerships(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*types.InstrumentOwnership, string, error)
		CreateInstrumentOwnership(ctx context.Context, input *types.InstrumentOwnershipCreationRequestInput) (*types.InstrumentOwnership, error)
		ReadInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) (*types.InstrumentOwnership, error)
		UpdateInstrumentOwnership(ctx context.Context, instrumentOwnershipID, ownerID string, input *types.InstrumentOwnershipUpdateRequestInput) error
		ArchiveInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) error
	}
)
