// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidIngredientCreationRequestInput struct {
		Description                 string               `json:"description"`
		Warning                     string               `json:"warning"`
		StorageInstructions         string               `json:"storageInstructions"`
		Slug                        string               `json:"slug"`
		ShoppingSuggestions         string               `json:"shoppingSuggestions"`
		PluralName                  string               `json:"pluralName"`
		Name                        string               `json:"name"`
		IconPath                    string               `json:"iconPath"`
		StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
		ContainsSesame              bool                 `json:"containsSesame"`
		IsLiquid                    bool                 `json:"isLiquid"`
		ContainsTreeNut             bool                 `json:"containsTreeNut"`
		ContainsWheat               bool                 `json:"containsWheat"`
		ContainsShellfish           bool                 `json:"containsShellfish"`
		AnimalDerived               bool                 `json:"animalDerived"`
		IsAcid                      bool                 `json:"isAcid"`
		IsFat                       bool                 `json:"isFat"`
		IsFruit                     bool                 `json:"isFruit"`
		IsGrain                     bool                 `json:"isGrain"`
		IsHeat                      bool                 `json:"isHeat"`
		ContainsSoy                 bool                 `json:"containsSoy"`
		IsProtein                   bool                 `json:"isProtein"`
		IsSalt                      bool                 `json:"isSalt"`
		IsStarch                    bool                 `json:"isStarch"`
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
