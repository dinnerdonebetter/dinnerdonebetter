package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeValidIngredient builds a faked valid ingredient.
func BuildFakeValidIngredient() *types.ValidIngredient {
	return &types.ValidIngredient{
		ID:                                      ksuid.New().String(),
		Name:                                    buildUniqueString(),
		Description:                             buildUniqueString(),
		Warning:                                 buildUniqueString(),
		ContainsEgg:                             fake.Bool(),
		ContainsDairy:                           fake.Bool(),
		ContainsPeanut:                          fake.Bool(),
		ContainsTreeNut:                         fake.Bool(),
		ContainsSoy:                             fake.Bool(),
		ContainsWheat:                           fake.Bool(),
		ContainsShellfish:                       fake.Bool(),
		ContainsSesame:                          fake.Bool(),
		ContainsFish:                            fake.Bool(),
		ContainsGluten:                          fake.Bool(),
		AnimalFlesh:                             fake.Bool(),
		IsMeasuredVolumetrically:                fake.Bool(),
		IsLiquid:                                fake.Bool(),
		IconPath:                                buildUniqueString(),
		CreatedAt:                               fake.Date(),
		PluralName:                              buildUniqueString(),
		AnimalDerived:                           fake.Bool(),
		RestrictToPreparations:                  fake.Bool(),
		MinimumIdealStorageTemperatureInCelsius: func(f float32) *float32 { return &f }(fake.Float32()),
		MaximumIdealStorageTemperatureInCelsius: func(f float32) *float32 { return &f }(fake.Float32()),
		StorageInstructions:                     buildUniqueString(),
	}
}

// BuildFakeValidIngredientList builds a faked ValidIngredientList.
func BuildFakeValidIngredientList() *types.ValidIngredientList {
	var examples []*types.ValidIngredient
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredient())
	}

	return &types.ValidIngredientList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidIngredients: examples,
	}
}

// BuildFakeValidIngredientUpdateRequestInput builds a faked ValidIngredientUpdateRequestInput from a valid ingredient.
func BuildFakeValidIngredientUpdateRequestInput() *types.ValidIngredientUpdateRequestInput {
	validIngredient := BuildFakeValidIngredient()
	return &types.ValidIngredientUpdateRequestInput{
		Name:                                    &validIngredient.Name,
		Description:                             &validIngredient.Description,
		Warning:                                 &validIngredient.Warning,
		ContainsEgg:                             &validIngredient.ContainsEgg,
		ContainsDairy:                           &validIngredient.ContainsDairy,
		ContainsPeanut:                          &validIngredient.ContainsPeanut,
		ContainsTreeNut:                         &validIngredient.ContainsTreeNut,
		ContainsSoy:                             &validIngredient.ContainsSoy,
		ContainsWheat:                           &validIngredient.ContainsWheat,
		ContainsShellfish:                       &validIngredient.ContainsShellfish,
		ContainsSesame:                          &validIngredient.ContainsSesame,
		ContainsFish:                            &validIngredient.ContainsFish,
		ContainsGluten:                          &validIngredient.ContainsGluten,
		AnimalFlesh:                             &validIngredient.AnimalFlesh,
		IsMeasuredVolumetrically:                &validIngredient.IsMeasuredVolumetrically,
		IsLiquid:                                &validIngredient.IsLiquid,
		IconPath:                                &validIngredient.IconPath,
		PluralName:                              &validIngredient.PluralName,
		AnimalDerived:                           &validIngredient.AnimalDerived,
		RestrictToPreparations:                  &validIngredient.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: validIngredient.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: validIngredient.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     &validIngredient.StorageInstructions,
	}
}

// BuildFakeValidIngredientUpdateRequestInputFromValidIngredient builds a faked ValidIngredientUpdateRequestInput from a valid ingredient.
func BuildFakeValidIngredientUpdateRequestInputFromValidIngredient(validIngredient *types.ValidIngredient) *types.ValidIngredientUpdateRequestInput {
	return &types.ValidIngredientUpdateRequestInput{
		Name:                                    &validIngredient.Name,
		Description:                             &validIngredient.Description,
		Warning:                                 &validIngredient.Warning,
		ContainsEgg:                             &validIngredient.ContainsEgg,
		ContainsDairy:                           &validIngredient.ContainsDairy,
		ContainsPeanut:                          &validIngredient.ContainsPeanut,
		ContainsTreeNut:                         &validIngredient.ContainsTreeNut,
		ContainsSoy:                             &validIngredient.ContainsSoy,
		ContainsWheat:                           &validIngredient.ContainsWheat,
		ContainsShellfish:                       &validIngredient.ContainsShellfish,
		ContainsSesame:                          &validIngredient.ContainsSesame,
		ContainsFish:                            &validIngredient.ContainsFish,
		ContainsGluten:                          &validIngredient.ContainsGluten,
		AnimalFlesh:                             &validIngredient.AnimalFlesh,
		IsMeasuredVolumetrically:                &validIngredient.IsMeasuredVolumetrically,
		IsLiquid:                                &validIngredient.IsLiquid,
		IconPath:                                &validIngredient.IconPath,
		PluralName:                              &validIngredient.PluralName,
		AnimalDerived:                           &validIngredient.AnimalDerived,
		RestrictToPreparations:                  &validIngredient.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: validIngredient.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: validIngredient.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     &validIngredient.StorageInstructions,
	}
}

// BuildFakeValidIngredientCreationRequestInput builds a faked ValidIngredientCreationRequestInput.
func BuildFakeValidIngredientCreationRequestInput() *types.ValidIngredientCreationRequestInput {
	validIngredient := BuildFakeValidIngredient()
	return BuildFakeValidIngredientCreationRequestInputFromValidIngredient(validIngredient)
}

// BuildFakeValidIngredientCreationRequestInputFromValidIngredient builds a faked ValidIngredientCreationRequestInput from a valid ingredient.
func BuildFakeValidIngredientCreationRequestInputFromValidIngredient(validIngredient *types.ValidIngredient) *types.ValidIngredientCreationRequestInput {
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
		MinimumIdealStorageTemperatureInCelsius: *validIngredient.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: *validIngredient.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     validIngredient.StorageInstructions,
	}
}

// BuildFakeValidIngredientDatabaseCreationInputFromValidIngredient builds a faked ValidIngredientDatabaseCreationInput from a valid ingredient.
func BuildFakeValidIngredientDatabaseCreationInputFromValidIngredient(validIngredient *types.ValidIngredient) *types.ValidIngredientDatabaseCreationInput {
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
	}
}
