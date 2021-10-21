package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// BuildFakeValidIngredient builds a faked valid ingredient.
func BuildFakeValidIngredient() *types.ValidIngredient {
	return &types.ValidIngredient{
		ID:                 ksuid.New().String(),
		Name:               fake.Word(),
		Variant:            fake.Word(),
		Description:        fake.Word(),
		Warning:            fake.Word(),
		ContainsEgg:        fake.Bool(),
		ContainsDairy:      fake.Bool(),
		ContainsPeanut:     fake.Bool(),
		ContainsTreeNut:    fake.Bool(),
		ContainsSoy:        fake.Bool(),
		ContainsWheat:      fake.Bool(),
		ContainsShellfish:  fake.Bool(),
		ContainsSesame:     fake.Bool(),
		ContainsFish:       fake.Bool(),
		ContainsGluten:     fake.Bool(),
		AnimalFlesh:        fake.Bool(),
		AnimalDerived:      fake.Bool(),
		MeasurableByVolume: fake.Bool(),
		Icon:               fake.Word(),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
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
		Name:               validIngredient.Name,
		Variant:            validIngredient.Variant,
		Description:        validIngredient.Description,
		Warning:            validIngredient.Warning,
		ContainsEgg:        validIngredient.ContainsEgg,
		ContainsDairy:      validIngredient.ContainsDairy,
		ContainsPeanut:     validIngredient.ContainsPeanut,
		ContainsTreeNut:    validIngredient.ContainsTreeNut,
		ContainsSoy:        validIngredient.ContainsSoy,
		ContainsWheat:      validIngredient.ContainsWheat,
		ContainsShellfish:  validIngredient.ContainsShellfish,
		ContainsSesame:     validIngredient.ContainsSesame,
		ContainsFish:       validIngredient.ContainsFish,
		ContainsGluten:     validIngredient.ContainsGluten,
		AnimalFlesh:        validIngredient.AnimalFlesh,
		AnimalDerived:      validIngredient.AnimalDerived,
		MeasurableByVolume: validIngredient.MeasurableByVolume,
		Icon:               validIngredient.Icon,
	}
}

// BuildFakeValidIngredientUpdateRequestInputFromValidIngredient builds a faked ValidIngredientUpdateRequestInput from a valid ingredient.
func BuildFakeValidIngredientUpdateRequestInputFromValidIngredient(validIngredient *types.ValidIngredient) *types.ValidIngredientUpdateRequestInput {
	return &types.ValidIngredientUpdateRequestInput{
		Name:               validIngredient.Name,
		Variant:            validIngredient.Variant,
		Description:        validIngredient.Description,
		Warning:            validIngredient.Warning,
		ContainsEgg:        validIngredient.ContainsEgg,
		ContainsDairy:      validIngredient.ContainsDairy,
		ContainsPeanut:     validIngredient.ContainsPeanut,
		ContainsTreeNut:    validIngredient.ContainsTreeNut,
		ContainsSoy:        validIngredient.ContainsSoy,
		ContainsWheat:      validIngredient.ContainsWheat,
		ContainsShellfish:  validIngredient.ContainsShellfish,
		ContainsSesame:     validIngredient.ContainsSesame,
		ContainsFish:       validIngredient.ContainsFish,
		ContainsGluten:     validIngredient.ContainsGluten,
		AnimalFlesh:        validIngredient.AnimalFlesh,
		AnimalDerived:      validIngredient.AnimalDerived,
		MeasurableByVolume: validIngredient.MeasurableByVolume,
		Icon:               validIngredient.Icon,
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
		ID:                 validIngredient.ID,
		Name:               validIngredient.Name,
		Variant:            validIngredient.Variant,
		Description:        validIngredient.Description,
		Warning:            validIngredient.Warning,
		ContainsEgg:        validIngredient.ContainsEgg,
		ContainsDairy:      validIngredient.ContainsDairy,
		ContainsPeanut:     validIngredient.ContainsPeanut,
		ContainsTreeNut:    validIngredient.ContainsTreeNut,
		ContainsSoy:        validIngredient.ContainsSoy,
		ContainsWheat:      validIngredient.ContainsWheat,
		ContainsShellfish:  validIngredient.ContainsShellfish,
		ContainsSesame:     validIngredient.ContainsSesame,
		ContainsFish:       validIngredient.ContainsFish,
		ContainsGluten:     validIngredient.ContainsGluten,
		AnimalFlesh:        validIngredient.AnimalFlesh,
		AnimalDerived:      validIngredient.AnimalDerived,
		MeasurableByVolume: validIngredient.MeasurableByVolume,
		Icon:               validIngredient.Icon,
	}
}

// BuildFakeValidIngredientDatabaseCreationInput builds a faked ValidIngredientDatabaseCreationInput.
func BuildFakeValidIngredientDatabaseCreationInput() *types.ValidIngredientDatabaseCreationInput {
	validIngredient := BuildFakeValidIngredient()
	return BuildFakeValidIngredientDatabaseCreationInputFromValidIngredient(validIngredient)
}

// BuildFakeValidIngredientDatabaseCreationInputFromValidIngredient builds a faked ValidIngredientDatabaseCreationInput from a valid ingredient.
func BuildFakeValidIngredientDatabaseCreationInputFromValidIngredient(validIngredient *types.ValidIngredient) *types.ValidIngredientDatabaseCreationInput {
	return &types.ValidIngredientDatabaseCreationInput{
		ID:                 validIngredient.ID,
		Name:               validIngredient.Name,
		Variant:            validIngredient.Variant,
		Description:        validIngredient.Description,
		Warning:            validIngredient.Warning,
		ContainsEgg:        validIngredient.ContainsEgg,
		ContainsDairy:      validIngredient.ContainsDairy,
		ContainsPeanut:     validIngredient.ContainsPeanut,
		ContainsTreeNut:    validIngredient.ContainsTreeNut,
		ContainsSoy:        validIngredient.ContainsSoy,
		ContainsWheat:      validIngredient.ContainsWheat,
		ContainsShellfish:  validIngredient.ContainsShellfish,
		ContainsSesame:     validIngredient.ContainsSesame,
		ContainsFish:       validIngredient.ContainsFish,
		ContainsGluten:     validIngredient.ContainsGluten,
		AnimalFlesh:        validIngredient.AnimalFlesh,
		AnimalDerived:      validIngredient.AnimalDerived,
		MeasurableByVolume: validIngredient.MeasurableByVolume,
		Icon:               validIngredient.Icon,
	}
}
