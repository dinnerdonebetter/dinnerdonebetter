// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepUpdateRequestInput struct {
		BelongsToRecipe         string               `json:"belongsToRecipe"`
		ConditionExpression     string               `json:"conditionExpression"`
		ExplicitInstructions    string               `json:"explicitInstructions"`
		Notes                   string               `json:"notes"`
		Preparation             ValidPreparation     `json:"preparation"`
		EstimatedTimeInSeconds  OptionalUint32Range  `json:"estimatedTimeInSeconds"`
		TemperatureInCelsius    OptionalFloat32Range `json:"temperatureInCelsius"`
		Index                   uint64               `json:"index"`
		Optional                bool                 `json:"optional"`
		StartTimerAutomatically bool                 `json:"startTimerAutomatically"`
	}
)
