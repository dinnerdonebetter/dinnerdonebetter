package types

const (
	// FinalizeMealPlansWithExpiredVotingPeriodsChoreType asks the worker to finalize meal plans with expired voting periods.
	FinalizeMealPlansWithExpiredVotingPeriodsChoreType choreType = "finalize_meal_plans_with_expired_voting_periods"
)

type (
	dataType string

	// PreWriteMessage represents an event that asks a worker to write data to the datastore.
	PreWriteMessage struct {
		_ struct{}

		DataType                   dataType                                         `json:"dataType"`
		ValidInstrument            *ValidInstrumentDatabaseCreationInput            `json:"validInstrument,omitempty"`
		ValidIngredient            *ValidIngredientDatabaseCreationInput            `json:"validIngredient,omitempty"`
		ValidPreparation           *ValidPreparationDatabaseCreationInput           `json:"validPreparation,omitempty"`
		ValidIngredientPreparation *ValidIngredientPreparationDatabaseCreationInput `json:"validIngredientPreparation,omitempty"`
		MealID                     string                                           `json:"mealID,omitempty"`
		Meal                       *MealDatabaseCreationInput                       `json:"meal,omitempty"`
		RecipeID                   string                                           `json:"recipeID,omitempty"`
		Recipe                     *RecipeDatabaseCreationInput                     `json:"recipe,omitempty"`
		RecipeStepID               string                                           `json:"recipeStepID,omitempty"`
		RecipeStep                 *RecipeStepDatabaseCreationInput                 `json:"recipeStep,omitempty"`
		RecipeStepInstrument       *RecipeStepInstrumentDatabaseCreationInput       `json:"recipeStepInstrument,omitempty"`
		RecipeStepIngredient       *RecipeStepIngredientDatabaseCreationInput       `json:"recipeStepIngredient,omitempty"`
		RecipeStepProduct          *RecipeStepProductDatabaseCreationInput          `json:"recipeStepProduct,omitempty"`
		MealPlan                   *MealPlanDatabaseCreationInput                   `json:"mealPlan,omitempty"`
		MealPlanID                 string                                           `json:"mealPlanID,omitempty"`
		MealPlanOption             *MealPlanOptionDatabaseCreationInput             `json:"mealPlanOption,omitempty"`
		MealPlanOptionID           string                                           `json:"mealPlanOptionID,omitempty"`
		MealPlanOptionVote         *MealPlanOptionVoteDatabaseCreationInput         `json:"mealPlanOptionVote,omitempty"`
		Webhook                    *WebhookDatabaseCreationInput                    `json:"webhook,omitempty"`
		UserMembership             *HouseholdUserMembershipDatabaseCreationInput    `json:"userMembership,omitempty"`
		HouseholdInvitation        *HouseholdInvitationDatabaseCreationInput        `json:"householdInvitation,omitempty"`
		AttributableToUserID       string                                           `json:"attributableToUserID,omitempty"`
		AttributableToHouseholdID  string                                           `json:"attributableToHouseholdID,omitempty"`
	}

	// PreUpdateMessage represents an event that asks a worker to update data in the datastore.
	PreUpdateMessage struct {
		_ struct{}

		DataType                   dataType                    `json:"dataType"`
		ValidInstrument            *ValidInstrument            `json:"validInstrument,omitempty"`
		ValidIngredient            *ValidIngredient            `json:"validIngredient,omitempty"`
		ValidPreparation           *ValidPreparation           `json:"validPreparation,omitempty"`
		ValidIngredientPreparation *ValidIngredientPreparation `json:"validIngredientPreparation,omitempty"`
		MealID                     string                      `json:"mealID,omitempty"`
		Meal                       *Meal                       `json:"meal,omitempty"`
		RecipeID                   string                      `json:"recipeID,omitempty"`
		Recipe                     *Recipe                     `json:"recipe,omitempty"`
		RecipeStepID               string                      `json:"recipeStepID,omitempty"`
		RecipeStep                 *RecipeStep                 `json:"recipeStep,omitempty"`
		RecipeStepInstrument       *RecipeStepInstrument       `json:"recipeStepInstrument,omitempty"`
		RecipeStepIngredient       *RecipeStepIngredient       `json:"recipeStepIngredient,omitempty"`
		RecipeStepProduct          *RecipeStepProduct          `json:"recipeStepProduct,omitempty"`
		MealPlan                   *MealPlan                   `json:"mealPlan,omitempty"`
		MealPlanID                 string                      `json:"mealPlanID,omitempty"`
		MealPlanOption             *MealPlanOption             `json:"mealPlanOption,omitempty"`
		MealPlanOptionID           string                      `json:"mealPlanOptionID,omitempty"`
		MealPlanOptionVote         *MealPlanOptionVote         `json:"mealPlanOptionVote,omitempty"`
		HouseholdInvitation        *HouseholdInvitation        `json:"householdInvitation,omitempty"`
		AttributableToUserID       string                      `json:"attributableToUserID,omitempty"`
		AttributableToHouseholdID  string                      `json:"attributableToHouseholdID,omitempty"`
	}

	// PreArchiveMessage represents an event that asks a worker to archive data in the datastore.
	PreArchiveMessage struct {
		_ struct{}

		DataType                     dataType `json:"dataType"`
		ValidInstrumentID            string   `json:"validInstrumentID,omitempty"`
		ValidIngredientID            string   `json:"validIngredientID,omitempty"`
		ValidPreparationID           string   `json:"validPreparationID,omitempty"`
		ValidIngredientPreparationID string   `json:"validIngredientPreparationID,omitempty"`
		MealID                       string   `json:"mealID,omitempty"`
		RecipeID                     string   `json:"recipeID,omitempty"`
		RecipeStepID                 string   `json:"recipeStepID,omitempty"`
		RecipeStepInstrumentID       string   `json:"recipeStepInstrumentID,omitempty"`
		RecipeStepIngredientID       string   `json:"recipeStepIngredientID,omitempty"`
		RecipeStepProductID          string   `json:"recipeStepProductID,omitempty"`
		MealPlanID                   string   `json:"mealPlanID,omitempty"`
		MealPlanOptionID             string   `json:"mealPlanOptionID,omitempty"`
		MealPlanOptionVoteID         string   `json:"mealPlanOptionVoteID,omitempty"`
		WebhookID                    string   `json:"webhookID,omitempty"`
		HouseholdInvitationID        string   `json:"HouseholdInvitationID,omitempty"`
		AttributableToUserID         string   `json:"attributableToUserID,omitempty"`
		AttributableToHouseholdID    string   `json:"attributableToHouseholdID,omitempty"`
	}

	// DataChangeMessage represents an event that asks a worker to write data to the datastore.
	DataChangeMessage struct {
		_ struct{}

		DataType                   dataType                    `json:"dataType"`
		MessageType                string                      `json:"messageType"`
		ValidInstrument            *ValidInstrument            `json:"validInstrument,omitempty"`
		ValidIngredient            *ValidIngredient            `json:"validIngredient,omitempty"`
		ValidPreparation           *ValidPreparation           `json:"validPreparation,omitempty"`
		ValidIngredientPreparation *ValidIngredientPreparation `json:"validIngredientPreparation,omitempty"`
		MealID                     string                      `json:"mealID,omitempty"`
		Meal                       *Meal                       `json:"meal,omitempty"`
		RecipeID                   string                      `json:"recipeID,omitempty"`
		Recipe                     *Recipe                     `json:"recipe,omitempty"`
		RecipeStepID               string                      `json:"recipeStepID,omitempty"`
		RecipeStep                 *RecipeStep                 `json:"recipeStep,omitempty"`
		RecipeStepInstrument       *RecipeStepInstrument       `json:"recipeStepInstrument,omitempty"`
		RecipeStepIngredient       *RecipeStepIngredient       `json:"recipeStepIngredient,omitempty"`
		RecipeStepProduct          *RecipeStepProduct          `json:"recipeStepProduct,omitempty"`
		MealPlan                   *MealPlan                   `json:"mealPlan,omitempty"`
		MealPlanID                 string                      `json:"mealPlanID,omitempty"`
		MealPlanOption             *MealPlanOption             `json:"mealPlanOption,omitempty"`
		MealPlanOptionID           string                      `json:"mealPlanOptionID,omitempty"`
		MealPlanOptionVote         *MealPlanOptionVote         `json:"mealPlanOptionVote,omitempty"`
		MealPlanOptionVoteID       string                      `json:"mealPlanOptionVoteID,omitempty"`
		Webhook                    *Webhook                    `json:"webhook,omitempty"`
		Household                  *Household                  `json:"household,omitempty"`
		HouseholdID                string                      `json:"hosueholdID,omitempty"`
		HouseholdInvitation        *HouseholdInvitation        `json:"householdInvitation,omitempty"`
		UserMembership             *HouseholdUserMembership    `json:"userMembership,omitempty"`
		Context                    map[string]string           `json:"context,omitempty"`
		AttributableToUserID       string                      `json:"attributableToUserID,omitempty"`
		AttributableToHouseholdID  string                      `json:"attributableToHouseholdID,omitempty"`
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
