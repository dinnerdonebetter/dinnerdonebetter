// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidPreparation struct {
		PastTense                   string                     `json:"pastTense"`
		ArchivedAt                  string                     `json:"archivedAt"`
		CreatedAt                   string                     `json:"createdAt"`
		Description                 string                     `json:"description"`
		IconPath                    string                     `json:"iconPath"`
		ID                          string                     `json:"id"`
		LastUpdatedAt               string                     `json:"lastUpdatedAt"`
		Slug                        string                     `json:"slug"`
		Name                        string                     `json:"name"`
		InstrumentCount             Uint16RangeWithOptionalMax `json:"instrumentCount"`
		IngredientCount             Uint16RangeWithOptionalMax `json:"ingredientCount"`
		VesselCount                 Uint16RangeWithOptionalMax `json:"vesselCount"`
		YieldsNothing               bool                       `json:"yieldsNothing"`
		RestrictToIngredients       bool                       `json:"restrictToIngredients"`
		OnlyForVessels              bool                       `json:"onlyForVessels"`
		TemperatureRequired         bool                       `json:"temperatureRequired"`
		TimeEstimateRequired        bool                       `json:"timeEstimateRequired"`
		ConditionExpressionRequired bool                       `json:"conditionExpressionRequired"`
		ConsumesVessel              bool                       `json:"consumesVessel"`
	}
)
