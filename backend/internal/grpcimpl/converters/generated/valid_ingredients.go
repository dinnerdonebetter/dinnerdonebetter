package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidIngredientCreationRequestInputToValidIngredient(input *messages.ValidIngredientCreationRequestInput) *messages.ValidIngredient {

output := &messages.ValidIngredient{
    ContainsFish: input.ContainsFish,
    IsStarch: input.IsStarch,
    IsGrain: input.IsGrain,
    IsFruit: input.IsFruit,
    StorageTemperatureInCelsius: input.StorageTemperatureInCelsius,
    ContainsEgg: input.ContainsEgg,
    ContainsSoy: input.ContainsSoy,
    ContainsSesame: input.ContainsSesame,
    ContainsTreeNut: input.ContainsTreeNut,
    PluralName: input.PluralName,
    Name: input.Name,
    Slug: input.Slug,
    ContainsGluten: input.ContainsGluten,
    StorageInstructions: input.StorageInstructions,
    ContainsPeanut: input.ContainsPeanut,
    AnimalFlesh: input.AnimalFlesh,
    ShoppingSuggestions: input.ShoppingSuggestions,
    ContainsShellfish: input.ContainsShellfish,
    AnimalDerived: input.AnimalDerived,
    IsAcid: input.IsAcid,
    Description: input.Description,
    IsSalt: input.IsSalt,
    IsProtein: input.IsProtein,
    IsLiquid: input.IsLiquid,
    Warning: input.Warning,
    IconPath: input.IconPath,
    RestrictToPreparations: input.RestrictToPreparations,
    ContainsDairy: input.ContainsDairy,
    ContainsAlcohol: input.ContainsAlcohol,
    ContainsWheat: input.ContainsWheat,
    IsFat: input.IsFat,
    IsHeat: input.IsHeat,
}

return output
}
func ConvertValidIngredientUpdateRequestInputToValidIngredient(input *messages.ValidIngredientUpdateRequestInput) *messages.ValidIngredient {

output := &messages.ValidIngredient{
    IsProtein: input.IsProtein,
    IsAcid: input.IsAcid,
    ContainsTreeNut: input.ContainsTreeNut,
    Description: input.Description,
    ContainsPeanut: input.ContainsPeanut,
    IsFat: input.IsFat,
    StorageInstructions: input.StorageInstructions,
    AnimalDerived: input.AnimalDerived,
    ContainsFish: input.ContainsFish,
    IsHeat: input.IsHeat,
    IsLiquid: input.IsLiquid,
    IconPath: input.IconPath,
    ShoppingSuggestions: input.ShoppingSuggestions,
    IsStarch: input.IsStarch,
    Slug: input.Slug,
    ContainsWheat: input.ContainsWheat,
    ContainsGluten: input.ContainsGluten,
    IsFruit: input.IsFruit,
    ContainsShellfish: input.ContainsShellfish,
    Warning: input.Warning,
    PluralName: input.PluralName,
    ContainsDairy: input.ContainsDairy,
    IsGrain: input.IsGrain,
    StorageTemperatureInCelsius: input.StorageTemperatureInCelsius,
    ContainsSoy: input.ContainsSoy,
    ContainsAlcohol: input.ContainsAlcohol,
    ContainsEgg: input.ContainsEgg,
    RestrictToPreparations: input.RestrictToPreparations,
    ContainsSesame: input.ContainsSesame,
    AnimalFlesh: input.AnimalFlesh,
    IsSalt: input.IsSalt,
    Name: input.Name,
}

return output
}
