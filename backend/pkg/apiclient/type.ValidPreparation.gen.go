// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
ValidPreparation struct {
   ArchivedAt string `json:"archivedAt"`
 ConditionExpressionRequired bool `json:"conditionExpressionRequired"`
 ConsumesVessel bool `json:"consumesVessel"`
 CreatedAt string `json:"createdAt"`
 Description string `json:"description"`
 IconPath string `json:"iconPath"`
 ID string `json:"id"`
 IngredientCount Uint16RangeWithOptionalMax `json:"ingredientCount"`
 InstrumentCount Uint16RangeWithOptionalMax `json:"instrumentCount"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Name string `json:"name"`
 OnlyForVessels bool `json:"onlyForVessels"`
 PastTense string `json:"pastTense"`
 RestrictToIngredients bool `json:"restrictToIngredients"`
 Slug string `json:"slug"`
 TemperatureRequired bool `json:"temperatureRequired"`
 TimeEstimateRequired bool `json:"timeEstimateRequired"`
 VesselCount Uint16RangeWithOptionalMax `json:"vesselCount"`
 YieldsNothing bool `json:"yieldsNothing"`

}
)
