// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepCreationRequestInput struct {
		PreparationID           string                                              `json:"preparationID"`
		ExplicitInstructions    string                                              `json:"explicitInstructions"`
		Notes                   string                                              `json:"notes"`
		ConditionExpression     string                                              `json:"conditionExpression"`
		CompletionConditions    []RecipeStepCompletionConditionCreationRequestInput `json:"completionConditions"`
		Vessels                 []RecipeStepVesselCreationRequestInput              `json:"vessels"`
		Ingredients             []RecipeStepIngredientCreationRequestInput          `json:"ingredients"`
		Products                []RecipeStepProductCreationRequestInput             `json:"products"`
		Instruments             []RecipeStepInstrumentCreationRequestInput          `json:"instruments"`
		TemperatureInCelsius    OptionalFloat32Range                                `json:"temperatureInCelsius"`
		EstimatedTimeInSeconds  OptionalUint32Range                                 `json:"estimatedTimeInSeconds"`
		Index                   uint64                                              `json:"index"`
		Optional                bool                                                `json:"optional"`
		StartTimerAutomatically bool                                                `json:"startTimerAutomatically"`
	}
)
