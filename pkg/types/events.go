package types

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
		RecipeID                   string                                           `json:"recipeID"`
		Recipe                     *RecipeDatabaseCreationInput                     `json:"recipe,omitempty"`
		RecipeStepID               string                                           `json:"recipeStepID"`
		RecipeStep                 *RecipeStepDatabaseCreationInput                 `json:"recipeStep,omitempty"`
		RecipeStepInstrument       *RecipeStepInstrumentDatabaseCreationInput       `json:"recipeStepInstrument,omitempty"`
		RecipeStepIngredient       *RecipeStepIngredientDatabaseCreationInput       `json:"recipeStepIngredient,omitempty"`
		RecipeStepProduct          *RecipeStepProductDatabaseCreationInput          `json:"recipeStepProduct,omitempty"`
		MealPlan                   *MealPlanDatabaseCreationInput                   `json:"mealPlan,omitempty"`
		MealPlanID                 string                                           `json:"mealPlanID"`
		MealPlanOption             *MealPlanOptionDatabaseCreationInput             `json:"mealPlanOption,omitempty"`
		MealPlanOptionID           string                                           `json:"mealPlanOptionID"`
		MealPlanOptionVote         *MealPlanOptionVoteDatabaseCreationInput         `json:"mealPlanOptionVote,omitempty"`
		Webhook                    *WebhookDatabaseCreationInput                    `json:"webhook,omitempty"`
		UserMembership             *HouseholdUserMembershipDatabaseCreationInput    `json:"userMembership,omitempty"`
		HouseholdInvitation        *HouseholdInvitationDatabaseCreationInput        `json:"householdInvitation,omitempty"`
		AttributableToUserID       string                                           `json:"attributableToUserID"`
		AttributableToHouseholdID  string                                           `json:"attributableToHouseholdID"`
	}

	// PreUpdateMessage represents an event that asks a worker to update data in the datastore.
	PreUpdateMessage struct {
		_ struct{}

		DataType                   dataType                    `json:"dataType"`
		ValidInstrument            *ValidInstrument            `json:"validInstrument,omitempty"`
		ValidIngredient            *ValidIngredient            `json:"validIngredient,omitempty"`
		ValidPreparation           *ValidPreparation           `json:"validPreparation,omitempty"`
		ValidIngredientPreparation *ValidIngredientPreparation `json:"validIngredientPreparation,omitempty"`
		RecipeID                   string                      `json:"recipeID"`
		Recipe                     *Recipe                     `json:"recipe,omitempty"`
		RecipeStepID               string                      `json:"recipeStepID"`
		RecipeStep                 *RecipeStep                 `json:"recipeStep,omitempty"`
		RecipeStepInstrument       *RecipeStepInstrument       `json:"recipeStepInstrument,omitempty"`
		RecipeStepIngredient       *RecipeStepIngredient       `json:"recipeStepIngredient,omitempty"`
		RecipeStepProduct          *RecipeStepProduct          `json:"recipeStepProduct,omitempty"`
		MealPlan                   *MealPlan                   `json:"mealPlan,omitempty"`
		MealPlanID                 string                      `json:"mealPlanID"`
		MealPlanOption             *MealPlanOption             `json:"mealPlanOption,omitempty"`
		MealPlanOptionID           string                      `json:"mealPlanOptionID"`
		MealPlanOptionVote         *MealPlanOptionVote         `json:"mealPlanOptionVote,omitempty"`
		HouseholdInvitation        *HouseholdInvitation        `json:"householdInvitation,omitempty"`
		AttributableToUserID       string                      `json:"attributableToUserID"`
		AttributableToHouseholdID  string                      `json:"attributableToHouseholdID"`
	}

	// PreArchiveMessage represents an event that asks a worker to archive data in the datastore.
	PreArchiveMessage struct {
		_ struct{}

		DataType                     dataType `json:"dataType"`
		ValidInstrumentID            string   `json:"validInstrumentID"`
		ValidIngredientID            string   `json:"validIngredientID"`
		ValidPreparationID           string   `json:"validPreparationID"`
		ValidIngredientPreparationID string   `json:"validIngredientPreparationID"`
		RecipeID                     string   `json:"recipeID"`
		RecipeStepID                 string   `json:"recipeStepID"`
		RecipeStepInstrumentID       string   `json:"recipeStepInstrumentID"`
		RecipeStepIngredientID       string   `json:"recipeStepIngredientID"`
		RecipeStepProductID          string   `json:"recipeStepProductID"`
		MealPlanID                   string   `json:"mealPlanID"`
		MealPlanOptionID             string   `json:"mealPlanOptionID"`
		MealPlanOptionVoteID         string   `json:"mealPlanOptionVoteID"`
		WebhookID                    string   `json:"webhookID"`
		HouseholdInvitationID        string   `json:"HouseholdInvitationID"`
		AttributableToUserID         string   `json:"attributableToUserID"`
		AttributableToHouseholdID    string   `json:"attributableToHouseholdID"`
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
		RecipeID                   string                      `json:"recipeID"`
		Recipe                     *Recipe                     `json:"recipe,omitempty"`
		RecipeStepID               string                      `json:"recipeStepID"`
		RecipeStep                 *RecipeStep                 `json:"recipeStep,omitempty"`
		RecipeStepInstrument       *RecipeStepInstrument       `json:"recipeStepInstrument,omitempty"`
		RecipeStepIngredient       *RecipeStepIngredient       `json:"recipeStepIngredient,omitempty"`
		RecipeStepProduct          *RecipeStepProduct          `json:"recipeStepProduct,omitempty"`
		MealPlan                   *MealPlan                   `json:"mealPlan,omitempty"`
		MealPlanID                 string                      `json:"mealPlanID"`
		MealPlanOption             *MealPlanOption             `json:"mealPlanOption,omitempty"`
		MealPlanOptionID           string                      `json:"mealPlanOptionID"`
		MealPlanOptionVote         *MealPlanOptionVote         `json:"mealPlanOptionVote,omitempty"`
		MealPlanOptionVoteID       string                      `json:"mealPlanOptionVoteID"`
		Webhook                    *Webhook                    `json:"webhook,omitempty"`
		HouseholdInvitation        *HouseholdInvitation        `json:"householdInvitation,omitempty"`
		UserMembership             *HouseholdUserMembership    `json:"userMembership,omitempty"`
		Context                    map[string]string           `json:"context,omitempty"`
		AttributableToUserID       string                      `json:"attributableToUserID"`
		AttributableToHouseholdID  string                      `json:"attributableToHouseholdID"`
	}
)
