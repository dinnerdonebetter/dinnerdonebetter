// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidIngredient struct {
		LastUpdatedAt               string               `json:"lastUpdatedAt"`
		Warning                     string               `json:"warning"`
		ArchivedAt                  string               `json:"archivedAt"`
		StorageInstructions         string               `json:"storageInstructions"`
		Slug                        string               `json:"slug"`
		ShoppingSuggestions         string               `json:"shoppingSuggestions"`
		PluralName                  string               `json:"pluralName"`
		Name                        string               `json:"name"`
		CreatedAt                   string               `json:"createdAt"`
		ID                          string               `json:"id"`
		IconPath                    string               `json:"iconPath"`
		Description                 string               `json:"description"`
		StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
		ContainsSesame              bool                 `json:"containsSesame"`
		IsProtein                   bool                 `json:"isProtein"`
		ContainsTreeNut             bool                 `json:"containsTreeNut"`
		ContainsSoy                 bool                 `json:"containsSoy"`
		ContainsShellfish           bool                 `json:"containsShellfish"`
		IsAcid                      bool                 `json:"isAcid"`
		IsFat                       bool                 `json:"isFat"`
		IsFruit                     bool                 `json:"isFruit"`
		IsGrain                     bool                 `json:"isGrain"`
		IsHeat                      bool                 `json:"isHeat"`
		IsLiquid                    bool                 `json:"isLiquid"`
		ContainsWheat               bool                 `json:"containsWheat"`
		IsSalt                      bool                 `json:"isSalt"`
		IsStarch                    bool                 `json:"isStarch"`
		AnimalDerived               bool                 `json:"animalDerived"`
		ContainsPeanut              bool                 `json:"containsPeanut"`
		ContainsGluten              bool                 `json:"containsGluten"`
		RestrictToPreparations      bool                 `json:"restrictToPreparations"`
		ContainsFish                bool                 `json:"containsFish"`
		ContainsEgg                 bool                 `json:"containsEgg"`
		ContainsDairy               bool                 `json:"containsDairy"`
		ContainsAlcohol             bool                 `json:"containsAlcohol"`
		AnimalFlesh                 bool                 `json:"animalFlesh"`
	}
)
