package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertValidIngredientToValidIngredientUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertValidIngredientToValidIngredientUpdateRequestInput(input *types.ValidIngredient) *types.ValidIngredientUpdateRequestInput {
	x := &types.ValidIngredientUpdateRequestInput{
		Name:                                    &input.Name,
		Description:                             &input.Description,
		Warning:                                 &input.Warning,
		IconPath:                                &input.IconPath,
		ContainsDairy:                           &input.ContainsDairy,
		ContainsPeanut:                          &input.ContainsPeanut,
		ContainsTreeNut:                         &input.ContainsTreeNut,
		ContainsEgg:                             &input.ContainsEgg,
		ContainsWheat:                           &input.ContainsWheat,
		ContainsShellfish:                       &input.ContainsShellfish,
		ContainsSesame:                          &input.ContainsSesame,
		ContainsFish:                            &input.ContainsFish,
		ContainsGluten:                          &input.ContainsGluten,
		AnimalFlesh:                             &input.AnimalFlesh,
		IsMeasuredVolumetrically:                &input.IsMeasuredVolumetrically,
		IsLiquid:                                &input.IsLiquid,
		ContainsSoy:                             &input.ContainsSoy,
		PluralName:                              &input.PluralName,
		AnimalDerived:                           &input.AnimalDerived,
		RestrictToPreparations:                  &input.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: input.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: input.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     &input.StorageInstructions,
		Slug:                                    &input.Slug,
		ContainsAlcohol:                         &input.ContainsAlcohol,
		ShoppingSuggestions:                     &input.ShoppingSuggestions,
	}

	return x
}

// ConvertValidIngredientCreationRequestInputToValidIngredientDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertValidIngredientCreationRequestInputToValidIngredientDatabaseCreationInput(input *types.ValidIngredientCreationRequestInput) *types.ValidIngredientDatabaseCreationInput {
	x := &types.ValidIngredientDatabaseCreationInput{
		ID:                                      identifiers.New(),
		Name:                                    input.Name,
		Description:                             input.Description,
		Warning:                                 input.Warning,
		ContainsEgg:                             input.ContainsEgg,
		ContainsDairy:                           input.ContainsDairy,
		ContainsPeanut:                          input.ContainsPeanut,
		ContainsTreeNut:                         input.ContainsTreeNut,
		ContainsSoy:                             input.ContainsSoy,
		ContainsWheat:                           input.ContainsWheat,
		ContainsShellfish:                       input.ContainsShellfish,
		ContainsSesame:                          input.ContainsSesame,
		ContainsFish:                            input.ContainsFish,
		ContainsGluten:                          input.ContainsGluten,
		AnimalFlesh:                             input.AnimalFlesh,
		IsMeasuredVolumetrically:                input.IsMeasuredVolumetrically,
		IsLiquid:                                input.IsLiquid,
		IconPath:                                input.IconPath,
		PluralName:                              input.PluralName,
		AnimalDerived:                           input.AnimalDerived,
		RestrictToPreparations:                  input.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: input.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: input.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     input.StorageInstructions,
		Slug:                                    input.Slug,
		ContainsAlcohol:                         input.ContainsAlcohol,
		ShoppingSuggestions:                     input.ShoppingSuggestions,
	}

	return x
}

// ConvertValidIngredientToValidIngredientCreationRequestInput builds a ValidIngredientCreationRequestInput from a Ingredient.
func ConvertValidIngredientToValidIngredientCreationRequestInput(validIngredient *types.ValidIngredient) *types.ValidIngredientCreationRequestInput {
	return &types.ValidIngredientCreationRequestInput{
		ID:                                      validIngredient.ID,
		Name:                                    validIngredient.Name,
		Description:                             validIngredient.Description,
		Warning:                                 validIngredient.Warning,
		ContainsEgg:                             validIngredient.ContainsEgg,
		ContainsDairy:                           validIngredient.ContainsDairy,
		ContainsPeanut:                          validIngredient.ContainsPeanut,
		ContainsTreeNut:                         validIngredient.ContainsTreeNut,
		ContainsSoy:                             validIngredient.ContainsSoy,
		ContainsWheat:                           validIngredient.ContainsWheat,
		ContainsShellfish:                       validIngredient.ContainsShellfish,
		ContainsSesame:                          validIngredient.ContainsSesame,
		ContainsFish:                            validIngredient.ContainsFish,
		ContainsGluten:                          validIngredient.ContainsGluten,
		AnimalFlesh:                             validIngredient.AnimalFlesh,
		IsMeasuredVolumetrically:                validIngredient.IsMeasuredVolumetrically,
		IsLiquid:                                validIngredient.IsLiquid,
		IconPath:                                validIngredient.IconPath,
		PluralName:                              validIngredient.PluralName,
		AnimalDerived:                           validIngredient.AnimalDerived,
		RestrictToPreparations:                  validIngredient.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: validIngredient.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: validIngredient.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     validIngredient.StorageInstructions,
		Slug:                                    validIngredient.Slug,
		ContainsAlcohol:                         validIngredient.ContainsAlcohol,
		ShoppingSuggestions:                     validIngredient.ShoppingSuggestions,
	}
}

// ConvertValidIngredientToValidIngredientDatabaseCreationInput builds a ValidIngredientDatabaseCreationInput from a Ingredient.
func ConvertValidIngredientToValidIngredientDatabaseCreationInput(validIngredient *types.ValidIngredient) *types.ValidIngredientDatabaseCreationInput {
	return &types.ValidIngredientDatabaseCreationInput{
		ID:                                      validIngredient.ID,
		Name:                                    validIngredient.Name,
		Description:                             validIngredient.Description,
		Warning:                                 validIngredient.Warning,
		ContainsEgg:                             validIngredient.ContainsEgg,
		ContainsDairy:                           validIngredient.ContainsDairy,
		ContainsPeanut:                          validIngredient.ContainsPeanut,
		ContainsTreeNut:                         validIngredient.ContainsTreeNut,
		ContainsSoy:                             validIngredient.ContainsSoy,
		ContainsWheat:                           validIngredient.ContainsWheat,
		ContainsShellfish:                       validIngredient.ContainsShellfish,
		ContainsSesame:                          validIngredient.ContainsSesame,
		ContainsFish:                            validIngredient.ContainsFish,
		ContainsGluten:                          validIngredient.ContainsGluten,
		AnimalFlesh:                             validIngredient.AnimalFlesh,
		IsMeasuredVolumetrically:                validIngredient.IsMeasuredVolumetrically,
		IsLiquid:                                validIngredient.IsLiquid,
		IconPath:                                validIngredient.IconPath,
		PluralName:                              validIngredient.PluralName,
		AnimalDerived:                           validIngredient.AnimalDerived,
		RestrictToPreparations:                  validIngredient.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: validIngredient.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: validIngredient.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     validIngredient.StorageInstructions,
		Slug:                                    validIngredient.Slug,
		ContainsAlcohol:                         validIngredient.ContainsAlcohol,
		ShoppingSuggestions:                     validIngredient.ShoppingSuggestions,
	}
}

// ConvertNullableValidIngredientToValidIngredient converts a NullableValidIngredient to a ValidIngredient.
func ConvertNullableValidIngredientToValidIngredient(input *types.NullableValidIngredient) *types.ValidIngredient {
	return &types.ValidIngredient{
		CreatedAt:                               *input.CreatedAt,
		LastUpdatedAt:                           input.LastUpdatedAt,
		ArchivedAt:                              input.ArchivedAt,
		ID:                                      *input.ID,
		Warning:                                 *input.Warning,
		Description:                             *input.Description,
		IconPath:                                *input.IconPath,
		PluralName:                              *input.PluralName,
		StorageInstructions:                     *input.StorageInstructions,
		Name:                                    *input.Name,
		MaximumIdealStorageTemperatureInCelsius: input.MaximumIdealStorageTemperatureInCelsius,
		MinimumIdealStorageTemperatureInCelsius: input.MinimumIdealStorageTemperatureInCelsius,
		ContainsShellfish:                       *input.ContainsShellfish,
		ContainsDairy:                           *input.ContainsDairy,
		AnimalFlesh:                             *input.AnimalFlesh,
		IsMeasuredVolumetrically:                *input.IsMeasuredVolumetrically,
		IsLiquid:                                *input.IsLiquid,
		ContainsPeanut:                          *input.ContainsPeanut,
		ContainsTreeNut:                         *input.ContainsTreeNut,
		ContainsEgg:                             *input.ContainsEgg,
		ContainsWheat:                           *input.ContainsWheat,
		ContainsSoy:                             *input.ContainsSoy,
		AnimalDerived:                           *input.AnimalDerived,
		RestrictToPreparations:                  *input.RestrictToPreparations,
		ContainsSesame:                          *input.ContainsSesame,
		ContainsFish:                            *input.ContainsFish,
		ContainsGluten:                          *input.ContainsGluten,
		Slug:                                    *input.Slug,
		ContainsAlcohol:                         *input.ContainsAlcohol,
		ShoppingSuggestions:                     *input.ShoppingSuggestions,
	}
}

// ConvertValidIngredientToNullableValidIngredient converts a NullableValidIngredient to a ValidIngredient.
func ConvertValidIngredientToNullableValidIngredient(input *types.ValidIngredient) *types.NullableValidIngredient {
	return &types.NullableValidIngredient{
		CreatedAt:                               &input.CreatedAt,
		LastUpdatedAt:                           input.LastUpdatedAt,
		ArchivedAt:                              input.ArchivedAt,
		ID:                                      &input.ID,
		Warning:                                 &input.Warning,
		Description:                             &input.Description,
		IconPath:                                &input.IconPath,
		PluralName:                              &input.PluralName,
		StorageInstructions:                     &input.StorageInstructions,
		Name:                                    &input.Name,
		MaximumIdealStorageTemperatureInCelsius: input.MaximumIdealStorageTemperatureInCelsius,
		MinimumIdealStorageTemperatureInCelsius: input.MinimumIdealStorageTemperatureInCelsius,
		ContainsShellfish:                       &input.ContainsShellfish,
		ContainsDairy:                           &input.ContainsDairy,
		AnimalFlesh:                             &input.AnimalFlesh,
		IsMeasuredVolumetrically:                &input.IsMeasuredVolumetrically,
		IsLiquid:                                &input.IsLiquid,
		ContainsPeanut:                          &input.ContainsPeanut,
		ContainsTreeNut:                         &input.ContainsTreeNut,
		ContainsEgg:                             &input.ContainsEgg,
		ContainsWheat:                           &input.ContainsWheat,
		ContainsSoy:                             &input.ContainsSoy,
		AnimalDerived:                           &input.AnimalDerived,
		RestrictToPreparations:                  &input.RestrictToPreparations,
		ContainsSesame:                          &input.ContainsSesame,
		ContainsFish:                            &input.ContainsFish,
		ContainsGluten:                          &input.ContainsGluten,
		Slug:                                    &input.Slug,
		ContainsAlcohol:                         &input.ContainsAlcohol,
		ShoppingSuggestions:                     &input.ShoppingSuggestions,
	}
}
