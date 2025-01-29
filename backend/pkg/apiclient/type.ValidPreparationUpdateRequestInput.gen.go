// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
ValidPreparationUpdateRequestInput struct {
   ConditionExpressionRequired bool `json:"conditionExpressionRequired"`
 ConsumesVessel bool `json:"consumesVessel"`
 Description string `json:"description"`
 IconPath string `json:"iconPath"`
 IngredientCount Uint16RangeWithOptionalMaxUpdateRequestInput `json:"ingredientCount"`
 InstrumentCount Uint16RangeWithOptionalMaxUpdateRequestInput `json:"instrumentCount"`
 Name string `json:"name"`
 OnlyForVessels bool `json:"onlyForVessels"`
 PastTense string `json:"pastTense"`
 RestrictToIngredients bool `json:"restrictToIngredients"`
 Slug string `json:"slug"`
 TemperatureRequired bool `json:"temperatureRequired"`
 TimeEstimateRequired bool `json:"timeEstimateRequired"`
 VesselCount Uint16RangeWithOptionalMaxUpdateRequestInput `json:"vesselCount"`
 YieldsNothing bool `json:"yieldsNothing"`

}
)
