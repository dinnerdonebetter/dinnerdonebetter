// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidPreparationUpdateRequestInput struct {
		Name                        string                                       `json:"name"`
		Slug                        string                                       `json:"slug"`
		Description                 string                                       `json:"description"`
		IconPath                    string                                       `json:"iconPath"`
		PastTense                   string                                       `json:"pastTense"`
		InstrumentCount             Uint16RangeWithOptionalMaxUpdateRequestInput `json:"instrumentCount"`
		IngredientCount             Uint16RangeWithOptionalMaxUpdateRequestInput `json:"ingredientCount"`
		VesselCount                 Uint16RangeWithOptionalMaxUpdateRequestInput `json:"vesselCount"`
		ConsumesVessel              bool                                         `json:"consumesVessel"`
		ConditionExpressionRequired bool                                         `json:"conditionExpressionRequired"`
		RestrictToIngredients       bool                                         `json:"restrictToIngredients"`
		OnlyForVessels              bool                                         `json:"onlyForVessels"`
		TemperatureRequired         bool                                         `json:"temperatureRequired"`
		TimeEstimateRequired        bool                                         `json:"timeEstimateRequired"`
		YieldsNothing               bool                                         `json:"yieldsNothing"`
	}
)
