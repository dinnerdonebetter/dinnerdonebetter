// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidPreparation struct {
		PastTense                   string                     `json:"pastTense"`
		Name                        string                     `json:"name"`
		Slug                        string                     `json:"slug"`
		LastUpdatedAt               string                     `json:"lastUpdatedAt"`
		Description                 string                     `json:"description"`
		IconPath                    string                     `json:"iconPath"`
		ID                          string                     `json:"id"`
		ArchivedAt                  string                     `json:"archivedAt"`
		CreatedAt                   string                     `json:"createdAt"`
		InstrumentCount             Uint16RangeWithOptionalMax `json:"instrumentCount"`
		IngredientCount             Uint16RangeWithOptionalMax `json:"ingredientCount"`
		VesselCount                 Uint16RangeWithOptionalMax `json:"vesselCount"`
		ConditionExpressionRequired bool                       `json:"conditionExpressionRequired"`
		OnlyForVessels              bool                       `json:"onlyForVessels"`
		RestrictToIngredients       bool                       `json:"restrictToIngredients"`
		ConsumesVessel              bool                       `json:"consumesVessel"`
		TemperatureRequired         bool                       `json:"temperatureRequired"`
		TimeEstimateRequired        bool                       `json:"timeEstimateRequired"`
		YieldsNothing               bool                       `json:"yieldsNothing"`
	}
)
