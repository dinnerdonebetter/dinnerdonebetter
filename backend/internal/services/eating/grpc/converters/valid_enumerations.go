package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

func ConvertGRPCValidIngredientRequestCreationInputToValidIngredientRequestCreationInput(request *messages.CreateValidIngredientRequest) *types.ValidIngredientCreationRequestInput {
	return &types.ValidIngredientCreationRequestInput{
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: request.StorageTemperatureInCelsius.Max,
			Min: request.StorageTemperatureInCelsius.Min,
		},
		Warning:                request.Warning,
		IconPath:               request.IconPath,
		PluralName:             request.PluralName,
		StorageInstructions:    request.StorageInstructions,
		Name:                   request.Name,
		Description:            request.Description,
		Slug:                   request.Slug,
		ShoppingSuggestions:    request.ShoppingSuggestions,
		ContainsFish:           request.ContainsFish,
		ContainsShellfish:      request.ContainsShellfish,
		AnimalFlesh:            request.AnimalFlesh,
		ContainsEgg:            request.ContainsEgg,
		IsLiquid:               request.IsLiquid,
		ContainsSoy:            request.ContainsSoy,
		ContainsPeanut:         request.ContainsPeanut,
		AnimalDerived:          request.AnimalDerived,
		RestrictToPreparations: request.RestrictToPreparations,
		ContainsDairy:          request.ContainsDairy,
		ContainsSesame:         request.ContainsSesame,
		ContainsTreeNut:        request.ContainsTreeNut,
		ContainsWheat:          request.ContainsWheat,
		ContainsAlcohol:        request.ContainsAlcohol,
		ContainsGluten:         request.ContainsGluten,
		IsStarch:               request.IsStarch,
		IsProtein:              request.IsProtein,
		IsGrain:                request.IsGrain,
		IsFruit:                request.IsFruit,
		IsSalt:                 request.IsSalt,
		IsFat:                  request.IsFat,
		IsAcid:                 request.IsAcid,
		IsHeat:                 request.IsHeat,
	}
}

func ConvertValidIngredientToGRPCValidIngredient(ingredient *types.ValidIngredient) *messages.ValidIngredient {
	return &messages.ValidIngredient{
		CreatedAt:     ConvertTimeToPBTimestamp(ingredient.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(ingredient.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(ingredient.ArchivedAt),
		StorageTemperatureInCelsius: &messages.OptionalFloat32Range{
			Max: ingredient.StorageTemperatureInCelsius.Max,
			Min: ingredient.StorageTemperatureInCelsius.Min,
		},
		StorageInstructions:    ingredient.StorageInstructions,
		Warning:                ingredient.Warning,
		PluralName:             ingredient.PluralName,
		IconPath:               ingredient.IconPath,
		Name:                   ingredient.Name,
		ID:                     ingredient.ID,
		Description:            ingredient.Description,
		Slug:                   ingredient.Slug,
		ShoppingSuggestions:    ingredient.ShoppingSuggestions,
		ContainsEgg:            ingredient.ContainsEgg,
		ContainsAlcohol:        ingredient.ContainsAlcohol,
		ContainsPeanut:         ingredient.ContainsPeanut,
		ContainsWheat:          ingredient.ContainsWheat,
		ContainsSoy:            ingredient.ContainsSoy,
		AnimalDerived:          ingredient.AnimalDerived,
		RestrictToPreparations: ingredient.RestrictToPreparations,
		ContainsSesame:         ingredient.ContainsSesame,
		ContainsFish:           ingredient.ContainsFish,
		ContainsGluten:         ingredient.ContainsGluten,
		ContainsDairy:          ingredient.ContainsDairy,
		ContainsTreeNut:        ingredient.ContainsTreeNut,
		AnimalFlesh:            ingredient.AnimalFlesh,
		IsStarch:               ingredient.IsStarch,
		IsProtein:              ingredient.IsProtein,
		IsGrain:                ingredient.IsGrain,
		IsFruit:                ingredient.IsFruit,
		IsSalt:                 ingredient.IsSalt,
		IsFat:                  ingredient.IsFat,
		IsAcid:                 ingredient.IsAcid,
		IsHeat:                 ingredient.IsHeat,
		IsLiquid:               ingredient.IsLiquid,
		ContainsShellfish:      ingredient.ContainsShellfish,
	}
}
