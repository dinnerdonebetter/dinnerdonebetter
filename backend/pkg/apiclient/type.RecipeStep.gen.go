// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStep struct {
		ExplicitInstructions    string                          `json:"explicitInstructions"`
		ArchivedAt              string                          `json:"archivedAt"`
		Notes                   string                          `json:"notes"`
		ConditionExpression     string                          `json:"conditionExpression"`
		CreatedAt               string                          `json:"createdAt"`
		LastUpdatedAt           string                          `json:"lastUpdatedAt"`
		BelongsToRecipe         string                          `json:"belongsToRecipe"`
		ID                      string                          `json:"id"`
		Instruments             []RecipeStepInstrument          `json:"instruments"`
		Ingredients             []RecipeStepIngredient          `json:"ingredients"`
		Products                []RecipeStepProduct             `json:"products"`
		Vessels                 []RecipeStepVessel              `json:"vessels"`
		Media                   []RecipeMedia                   `json:"media"`
		CompletionConditions    []RecipeStepCompletionCondition `json:"completionConditions"`
		Preparation             ValidPreparation                `json:"preparation"`
		TemperatureInCelsius    OptionalFloat32Range            `json:"temperatureInCelsius"`
		EstimatedTimeInSeconds  OptionalUint32Range             `json:"estimatedTimeInSeconds"`
		Index                   uint64                          `json:"index"`
		Optional                bool                            `json:"optional"`
		StartTimerAutomatically bool                            `json:"startTimerAutomatically"`
	}
)
