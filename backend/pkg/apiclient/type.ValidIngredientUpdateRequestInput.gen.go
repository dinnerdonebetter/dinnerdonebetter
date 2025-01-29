// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
ValidIngredientUpdateRequestInput struct {
   AnimalDerived bool `json:"animalDerived"`
 AnimalFlesh bool `json:"animalFlesh"`
 ContainsAlcohol bool `json:"containsAlcohol"`
 ContainsDairy bool `json:"containsDairy"`
 ContainsEgg bool `json:"containsEgg"`
 ContainsFish bool `json:"containsFish"`
 ContainsGluten bool `json:"containsGluten"`
 ContainsPeanut bool `json:"containsPeanut"`
 ContainsSesame bool `json:"containsSesame"`
 ContainsShellfish bool `json:"containsShellfish"`
 ContainsSoy bool `json:"containsSoy"`
 ContainsTreeNut bool `json:"containsTreeNut"`
 ContainsWheat bool `json:"containsWheat"`
 Description string `json:"description"`
 IconPath string `json:"iconPath"`
 IsAcid bool `json:"isAcid"`
 IsFat bool `json:"isFat"`
 IsFruit bool `json:"isFruit"`
 IsGrain bool `json:"isGrain"`
 IsHeat bool `json:"isHeat"`
 IsLiquid bool `json:"isLiquid"`
 IsProtein bool `json:"isProtein"`
 IsSalt bool `json:"isSalt"`
 IsStarch bool `json:"isStarch"`
 Name string `json:"name"`
 PluralName string `json:"pluralName"`
 RestrictToPreparations bool `json:"restrictToPreparations"`
 ShoppingSuggestions string `json:"shoppingSuggestions"`
 Slug string `json:"slug"`
 StorageInstructions string `json:"storageInstructions"`
 StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
 Warning string `json:"warning"`

}
)
