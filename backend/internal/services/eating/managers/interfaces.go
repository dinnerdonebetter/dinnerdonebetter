package managers

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
		SearchMeals(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.Meal, error)
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

	ValidEnumerationsManager interface {
		SearchValidIngredientGroups(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, error)
		ListValidIngredientGroups(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, error)
		CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupCreationRequestInput) (*types.ValidIngredientGroup, error)
		ReadValidIngredientGroup(ctx context.Context, validIngredientGroupID string) (*types.ValidIngredientGroup, error)
		UpdateValidIngredientGroup(ctx context.Context, validIngredientGroupID string, input *types.ValidIngredientGroupUpdateRequestInput) error
		ArchiveValidIngredientGroup(ctx context.Context, validIngredientGroupID string) error

		ListValidIngredientMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, error)
		CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitCreationRequestInput) (*types.ValidIngredientMeasurementUnit, error)
		ReadValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error)
		UpdateValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string, input *types.ValidIngredientMeasurementUnitUpdateRequestInput) error
		ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error
		SearchValidIngredientMeasurementUnitsByIngredient(ctx context.Context, query string) ([]*types.ValidIngredientMeasurementUnit, error)
		SearchValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, query string) ([]*types.ValidIngredientMeasurementUnit, error)

		ListValidIngredientPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, error)
		CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationRequestInput) (*types.ValidIngredientPreparation, error)
		ReadValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error)
		UpdateValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string, input *types.ValidIngredientPreparationUpdateRequestInput) error
		ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error
		SearchValidIngredientPreparationsByIngredient(ctx context.Context, query string) ([]*types.ValidIngredientPreparation, error)
		SearchValidIngredientPreparationsByPreparation(ctx context.Context, query string) ([]*types.ValidIngredientPreparation, error)

		SearchValidIngredients(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error)
		ListValidIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error)
		CreateValidIngredient(ctx context.Context, input *types.ValidIngredientCreationRequestInput) (*types.ValidIngredient, error)
		ReadValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error)
		RandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error)
		UpdateValidIngredient(ctx context.Context, validIngredientID string, input *types.ValidIngredientUpdateRequestInput) error
		ArchiveValidIngredient(ctx context.Context, validIngredientID string) error
		SearchValidIngredientsByPreparationAndIngredientName(ctx context.Context, query string) ([]*types.ValidIngredient, error)

		ListValidIngredientStateIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, error)
		CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientCreationRequestInput) (*types.ValidIngredientStateIngredient, error)
		ReadValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error)
		UpdateValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string, input *types.ValidIngredientStateIngredientUpdateRequestInput) error
		ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error
		SearchValidIngredientStateIngredientsByIngredient(ctx context.Context, query string) ([]*types.ValidIngredientStateIngredient, error)
		SearchValidIngredientStateIngredientsByIngredientState(ctx context.Context, query string) ([]*types.ValidIngredientStateIngredient, error)

		SearchValidIngredientStates(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, error)
		ListValidIngredientStates(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, error)
		CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateCreationRequestInput) (*types.ValidIngredientState, error)
		ReadValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error)
		UpdateValidIngredientState(ctx context.Context, validIngredientStateID string, input *types.ValidIngredientStateUpdateRequestInput) error
		ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error

		SearchValidMeasurementUnits(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error)
		SearchValidMeasurementUnitsByIngredientID(ctx context.Context, query string) ([]*types.ValidMeasurementUnit, error)
		ListValidMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error)
		CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitCreationRequestInput) (*types.ValidMeasurementUnit, error)
		ReadValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error)
		UpdateValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string, input *types.ValidMeasurementUnitUpdateRequestInput) error
		ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error

		SearchValidInstruments(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidInstrument, error)
		ListValidInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidInstrument, error)
		CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentCreationRequestInput) (*types.ValidInstrument, error)
		ReadValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error)
		RandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error)
		UpdateValidInstrument(ctx context.Context, validInstrumentID string, input *types.ValidInstrumentUpdateRequestInput) error
		ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error

		ValidMeasurementUnitConversionsFromMeasurementUnit(ctx context.Context) ([]*types.ValidMeasurementUnitConversion, error)
		ValidMeasurementUnitConversionsToMeasurementUnit(ctx context.Context) ([]*types.ValidMeasurementUnitConversion, error)
		CreateValidMeasurementUnitConversion(ctx context.Context, input *types.ValidMeasurementUnitConversionCreationRequestInput) (*types.ValidMeasurementUnitConversion, error)
		ReadValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) (*types.ValidMeasurementUnitConversion, error)
		UpdateValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string, input *types.ValidMeasurementUnitConversionUpdateRequestInput) error
		ArchiveValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) error

		ListValidPreparationInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, error)
		CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationRequestInput) (*types.ValidPreparationInstrument, error)
		ReadValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error)
		UpdateValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string, input *types.ValidPreparationInstrumentUpdateRequestInput) error
		ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error
		SearchValidPreparationInstrumentsByPreparation(ctx context.Context, query string) ([]*types.ValidPreparationInstrument, error)
		SearchValidPreparationInstrumentsByInstrument(ctx context.Context, query string) ([]*types.ValidPreparationInstrument, error)

		SearchValidPreparations(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidPreparation, error)
		ListValidPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparation, error)
		CreateValidPreparation(ctx context.Context, input *types.ValidPreparationCreationRequestInput) (*types.ValidPreparation, error)
		ReadValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error)
		RandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error)
		UpdateValidPreparation(ctx context.Context, validPreparationID string, input *types.ValidPreparationUpdateRequestInput) error
		ArchiveValidPreparation(ctx context.Context, validPreparationID string) error

		ListValidPreparationVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, error)
		CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselCreationRequestInput) (*types.ValidPreparationVessel, error)
		ReadValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error)
		UpdateValidPreparationVessel(ctx context.Context, validPreparationVesselID string, input *types.ValidPreparationVesselUpdateRequestInput) error
		ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error
		SearchValidPreparationVesselsByPreparation(ctx context.Context, query string) ([]*types.ValidPreparationVessel, error)
		SearchValidPreparationVesselsByVessel(ctx context.Context, query string) ([]*types.ValidPreparationVessel, error)

		SearchValidVessels(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidVessel, error)
		ListValidVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidVessel, error)
		CreateValidVessel(ctx context.Context, input *types.ValidVesselCreationRequestInput) (*types.ValidVessel, error)
		ReadValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error)
		RandomValidVessel(ctx context.Context) (*types.ValidVessel, error)
		UpdateValidVessel(ctx context.Context, validVesselID string, input *types.ValidVesselUpdateRequestInput) error
		ArchiveValidVessel(ctx context.Context, validVesselID string) error
	}
)
