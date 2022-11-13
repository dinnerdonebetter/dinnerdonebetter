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

		DataType                         dataType                        `json:"dataType"`
		EventType                        CustomerEventType               `json:"messageType"`
		ValidInstrument                  *ValidInstrument                `json:"validInstrument,omitempty"`
		ValidIngredient                  *ValidIngredient                `json:"validIngredient,omitempty"`
		ValidPreparation                 *ValidPreparation               `json:"validPreparation,omitempty"`
		ValidIngredientPreparation       *ValidIngredientPreparation     `json:"validIngredientPreparation,omitempty"`
		MealID                           string                          `json:"mealID,omitempty"`
		Meal                             *Meal                           `json:"meal,omitempty"`
		RecipeID                         string                          `json:"recipeID,omitempty"`
		Recipe                           *Recipe                         `json:"recipe,omitempty"`
		RecipeStepID                     string                          `json:"recipeStepID,omitempty"`
		RecipeStep                       *RecipeStep                     `json:"recipeStep,omitempty"`
		RecipeStepProduct                *RecipeStepProduct              `json:"recipeStepProduct,omitempty"`
		RecipeStepInstrument             *RecipeStepInstrument           `json:"recipeStepInstrument,omitempty"`
		RecipeStepIngredient             *RecipeStepIngredient           `json:"recipeStepIngredient,omitempty"`
		MealPlan                         *MealPlan                       `json:"mealPlan,omitempty"`
		MealPlanID                       string                          `json:"mealPlanID,omitempty"`
		MealPlanEvent                    *MealPlanEvent                  `json:"mealPlanEvent,omitempty"`
		MealPlanEventID                  string                          `json:"mealPlanEventID,omitempty"`
		MealPlanOption                   *MealPlanOption                 `json:"mealPlanOption,omitempty"`
		MealPlanOptionID                 string                          `json:"mealPlanOptionID,omitempty"`
		MealPlanOptionVote               *MealPlanOptionVote             `json:"mealPlanOptionVote,omitempty"`
		MealPlanOptionVoteID             string                          `json:"mealPlanOptionVoteID,omitempty"`
		Webhook                          *Webhook                        `json:"webhook,omitempty"`
		Household                        *Household                      `json:"household,omitempty"`
		APIClientID                      string                          `json:"apiClientID,omitempty"`
		HouseholdID                      string                          `json:"householdID,omitempty"`
		HouseholdInvitation              *HouseholdInvitation            `json:"householdInvitation,omitempty"`
		HouseholdInvitationID            string                          `json:"householdInvitationID,omitempty"`
		UserMembership                   *HouseholdUserMembership        `json:"userMembership,omitempty"`
		ValidMeasurementUnitID           string                          `json:"validMeasurementUnitID,omitempty"`
		ValidMeasurementUnit             *ValidMeasurementUnit           `json:"validMeasurementUnit,omitempty"`
		ValidPreparationInstrumentID     string                          `json:"validPreparationInstrumentID,omitempty"`
		ValidPreparationInstrument       *ValidPreparationInstrument     `json:"validPreparationInstrument,omitempty"`
		ValidIngredientMeasurementUnitID string                          `json:"validIngredientMeasurementUnitID,omitempty"`
		ValidIngredientMeasurementUnit   *ValidIngredientMeasurementUnit `json:"validIngredientMeasurementUnit,omitempty"`
		MealPlanTask                     *MealPlanTask                   `json:"mealPlanTask,omitempty"`
		MealPlanTaskID                   string                          `json:"mealPlanTaskID,omitempty"`
		RecipePrepTask                   *RecipePrepTask                 `json:"recipePrepTask,omitempty"`
		RecipePrepTaskID                 string                          `json:"recipePrepTaskID,omitempty"`
		Context                          map[string]string               `json:"context,omitempty"`
		AttributableToUserID             string                          `json:"attributableToUserID,omitempty"`
		AttributableToHouseholdID        string                          `json:"attributableToHouseholdID,omitempty"`
		MealPlanGroceryListItem          *MealPlanGroceryListItem        `json:"mealPlanGroceryListItem,omitempty"`
		MealPlanGroceryListItemID        string                          `json:"mealPlanGroceryListItemID,omitempty"`
		ValidMeasurementConversion       *ValidMeasurementConversion     `json:"validMeasurementConversion,omitempty"`
		ValidMeasurementConversionID     string                          `json:"validMeasurementConversionID,omitempty"`
		ValidIngredientState             *ValidIngredientState           `json:"validIngredientState,omitempty"`
		ValidIngredientStateID           string                          `json:"validIngredientStateID,omitempty"`
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
