// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepCreationRequestInput struct {
		PreparationID           string                                              `json:"preparationID"`
		ExplicitInstructions    string                                              `json:"explicitInstructions"`
		Notes                   string                                              `json:"notes"`
		ConditionExpression     string                                              `json:"conditionExpression"`
		Ingredients             []RecipeStepIngredientCreationRequestInput          `json:"ingredients"`
		Instruments             []RecipeStepInstrumentCreationRequestInput          `json:"instruments"`
		Products                []RecipeStepProductCreationRequestInput             `json:"products"`
		CompletionConditions    []RecipeStepCompletionConditionCreationRequestInput `json:"completionConditions"`
		Vessels                 []RecipeStepVesselCreationRequestInput              `json:"vessels"`
		TemperatureInCelsius    OptionalFloat32Range                                `json:"temperatureInCelsius"`
		EstimatedTimeInSeconds  OptionalUint32Range                                 `json:"estimatedTimeInSeconds"`
		Index                   uint64                                              `json:"index"`
		StartTimerAutomatically bool                                                `json:"startTimerAutomatically"`
		Optional                bool                                                `json:"optional"`
	}
)
