package types

const (
	// FinalizeMealPlansWithExpiredVotingPeriodsChoreType asks the worker to finalize meal plans with expired voting periods.
	FinalizeMealPlansWithExpiredVotingPeriodsChoreType choreType = "finalize_meal_plans_with_expired_voting_periods"
	// CreateMealPlanTasksChoreType asks the worker to finalize meal plans with expired voting periods.
	CreateMealPlanTasksChoreType choreType = "create_meal_plan_tasks"
)

type (
	dataType string

	// CustomerEventType is the type to use/compare against when checking meal plan status.
	CustomerEventType string

	// DataChangeMessage represents an event that asks a worker to write data to the datastore.
	DataChangeMessage struct {
		_ struct{}

		HouseholdInvitation              *HouseholdInvitation            `json:"householdInvitation,omitempty"`
		ValidMeasurementConversion       *ValidMeasurementUnitConversion `json:"validMeasurementConversion,omitempty"`
		ValidInstrument                  *ValidInstrument                `json:"validInstrument,omitempty"`
		ValidIngredient                  *ValidIngredient                `json:"validIngredient,omitempty"`
		ValidPreparation                 *ValidPreparation               `json:"validPreparation,omitempty"`
		ValidIngredientState             *ValidIngredientState           `json:"validIngredientState,omitempty"`
		MealPlanGroceryListItem          *MealPlanGroceryListItem        `json:"mealPlanGroceryListItem,omitempty"`
		Meal                             *Meal                           `json:"meal,omitempty"`
		Context                          map[string]string               `json:"context,omitempty"`
		Recipe                           *Recipe                         `json:"recipe,omitempty"`
		RecipePrepTask                   *RecipePrepTask                 `json:"recipePrepTask,omitempty"`
		RecipeStep                       *RecipeStep                     `json:"recipeStep,omitempty"`
		RecipeStepProduct                *RecipeStepProduct              `json:"recipeStepProduct,omitempty"`
		RecipeStepInstrument             *RecipeStepInstrument           `json:"recipeStepInstrument,omitempty"`
		RecipeStepIngredient             *RecipeStepIngredient           `json:"recipeStepIngredient,omitempty"`
		MealPlan                         *MealPlan                       `json:"mealPlan,omitempty"`
		MealPlanTask                     *MealPlanTask                   `json:"mealPlanTask,omitempty"`
		MealPlanEvent                    *MealPlanEvent                  `json:"mealPlanEvent,omitempty"`
		Household                        *Household                      `json:"household,omitempty"`
		MealPlanOption                   *MealPlanOption                 `json:"mealPlanOption,omitempty"`
		ValidIngredientMeasurementUnit   *ValidIngredientMeasurementUnit `json:"validIngredientMeasurementUnit,omitempty"`
		MealPlanOptionVote               *MealPlanOptionVote             `json:"mealPlanOptionVote,omitempty"`
		ValidPreparationInstrument       *ValidPreparationInstrument     `json:"validPreparationInstrument,omitempty"`
		Webhook                          *Webhook                        `json:"webhook,omitempty"`
		ValidIngredientPreparation       *ValidIngredientPreparation     `json:"validIngredientPreparation,omitempty"`
		ValidMeasurementUnit             *ValidMeasurementUnit           `json:"validMeasurementUnit,omitempty"`
		UserMembership                   *HouseholdUserMembership        `json:"userMembership,omitempty"`
		RecipeStepCompletionCondition    *RecipeStepCompletionCondition  `json:"recipeStepCompletionCondition,omitempty"`
		HouseholdInvitationID            string                          `json:"householdInvitationID,omitempty"`
		HouseholdID                      string                          `json:"householdID,omitempty"`
		ValidMeasurementUnitID           string                          `json:"validMeasurementUnitID,omitempty"`
		APIClientID                      string                          `json:"apiClientID,omitempty"`
		ValidPreparationInstrumentID     string                          `json:"validPreparationInstrumentID,omitempty"`
		MealPlanOptionVoteID             string                          `json:"mealPlanOptionVoteID,omitempty"`
		ValidIngredientMeasurementUnitID string                          `json:"validIngredientMeasurementUnitID,omitempty"`
		MealPlanOptionID                 string                          `json:"mealPlanOptionID,omitempty"`
		MealPlanID                       string                          `json:"mealPlanID,omitempty"`
		MealPlanTaskID                   string                          `json:"mealPlanTaskID,omitempty"`
		RecipeStepID                     string                          `json:"recipeStepID,omitempty"`
		RecipePrepTaskID                 string                          `json:"recipePrepTaskID,omitempty"`
		RecipeID                         string                          `json:"recipeID,omitempty"`
		AttributableToUserID             string                          `json:"attributableToUserID,omitempty"`
		AttributableToHouseholdID        string                          `json:"attributableToHouseholdID,omitempty"`
		MealID                           string                          `json:"mealID,omitempty"`
		MealPlanGroceryListItemID        string                          `json:"mealPlanGroceryListItemID,omitempty"`
		EventType                        CustomerEventType               `json:"messageType"`
		ValidIngredientStateIngredient   *ValidIngredientStateIngredient `json:"validIngredientStateIngredient,omitempty"`
		ValidIngredientStateIngredientID string                          `json:"validIngredientStateIngredientID"`
		ValidMeasurementConversionID     string                          `json:"validMeasurementConversionID,omitempty"`
		DataType                         dataType                        `json:"dataType"`
		ValidIngredientStateID           string                          `json:"validIngredientStateID,omitempty"`
		MealPlanEventID                  string                          `json:"mealPlanEventID,omitempty"`
	}

	choreType string

	// ChoreMessage represents an event that asks a worker to perform a chore.
	ChoreMessage struct {
		_ struct{}

		ChoreType                 choreType `json:"choreType"`
		MealPlanID                string    `json:"mealPlanID,omitempty"`
		AttributableToHouseholdID string    `json:"attributableToHouseholdID,omitempty"`
	}
)
