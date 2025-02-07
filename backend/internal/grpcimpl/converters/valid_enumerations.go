package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func ConvertUpdateValidIngredientRequestToValidIngredient(req *messages.UpdateValidIngredientRequest) *types.ValidIngredient {
	return &types.ValidIngredient{
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Max: req.Input.StorageTemperatureInCelsius.Max,
			Min: req.Input.StorageTemperatureInCelsius.Min,
		},
		IconPath:               req.Input.IconPath,
		Warning:                req.Input.Warning,
		PluralName:             req.Input.PluralName,
		StorageInstructions:    req.Input.StorageInstructions,
		Name:                   req.Input.Name,
		ID:                     req.ValidIngredientID,
		Description:            req.Input.Description,
		Slug:                   req.Input.Slug,
		ShoppingSuggestions:    req.Input.ShoppingSuggestions,
		ContainsShellfish:      req.Input.ContainsShellfish,
		IsLiquid:               req.Input.IsLiquid,
		ContainsPeanut:         req.Input.ContainsPeanut,
		ContainsTreeNut:        req.Input.ContainsTreeNut,
		ContainsEgg:            req.Input.ContainsEgg,
		ContainsWheat:          req.Input.ContainsWheat,
		ContainsSoy:            req.Input.ContainsSoy,
		AnimalDerived:          req.Input.AnimalDerived,
		RestrictToPreparations: req.Input.RestrictToPreparations,
		ContainsSesame:         req.Input.ContainsSesame,
		ContainsFish:           req.Input.ContainsFish,
		ContainsGluten:         req.Input.ContainsGluten,
		ContainsDairy:          req.Input.ContainsDairy,
		ContainsAlcohol:        req.Input.ContainsAlcohol,
		AnimalFlesh:            req.Input.AnimalFlesh,
		IsStarch:               req.Input.IsStarch,
		IsProtein:              req.Input.IsProtein,
		IsGrain:                req.Input.IsGrain,
		IsFruit:                req.Input.IsFruit,
		IsSalt:                 req.Input.IsSalt,
		IsFat:                  req.Input.IsFat,
		IsAcid:                 req.Input.IsAcid,
		IsHeat:                 req.Input.IsHeat,
	}
}

func ConvertValidIngredientToProtobuf(validIngredient *types.ValidIngredient) *messages.ValidIngredient {
	return &messages.ValidIngredient{
		CreatedAt:     ConvertTimeToPBTimestamp(validIngredient.CreatedAt),
		LastUpdatedAt: ConvertTimePointerToPBTimestamp(validIngredient.LastUpdatedAt),
		ArchivedAt:    ConvertTimePointerToPBTimestamp(validIngredient.ArchivedAt),
		StorageTemperatureInCelsius: &messages.OptionalFloat32Range{
			Min: validIngredient.StorageTemperatureInCelsius.Min,
			Max: validIngredient.StorageTemperatureInCelsius.Max,
		},
		StorageInstructions:    validIngredient.StorageInstructions,
		Warning:                validIngredient.Warning,
		PluralName:             validIngredient.PluralName,
		IconPath:               validIngredient.IconPath,
		Name:                   validIngredient.Name,
		ID:                     validIngredient.ID,
		Description:            validIngredient.Description,
		Slug:                   validIngredient.Slug,
		ShoppingSuggestions:    validIngredient.ShoppingSuggestions,
		ContainsEgg:            validIngredient.ContainsEgg,
		ContainsAlcohol:        validIngredient.ContainsAlcohol,
		ContainsPeanut:         validIngredient.ContainsPeanut,
		ContainsWheat:          validIngredient.ContainsWheat,
		ContainsSoy:            validIngredient.ContainsSoy,
		AnimalDerived:          validIngredient.AnimalDerived,
		RestrictToPreparations: validIngredient.RestrictToPreparations,
		ContainsSesame:         validIngredient.ContainsSesame,
		ContainsFish:           validIngredient.ContainsFish,
		ContainsGluten:         validIngredient.ContainsGluten,
		ContainsDairy:          validIngredient.ContainsDairy,
		ContainsTreeNut:        validIngredient.ContainsTreeNut,
		AnimalFlesh:            validIngredient.AnimalFlesh,
		IsStarch:               validIngredient.IsStarch,
		IsProtein:              validIngredient.IsProtein,
		IsGrain:                validIngredient.IsGrain,
		IsFruit:                validIngredient.IsFruit,
		IsSalt:                 validIngredient.IsSalt,
		IsFat:                  validIngredient.IsFat,
		IsAcid:                 validIngredient.IsAcid,
		IsHeat:                 validIngredient.IsHeat,
		IsLiquid:               validIngredient.IsLiquid,
		ContainsShellfish:      validIngredient.ContainsShellfish,
	}
}
