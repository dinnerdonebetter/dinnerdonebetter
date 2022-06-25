package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeValidIngredient builds a faked valid ingredient.
func BuildFakeValidIngredient() *types.ValidIngredient {
	return &types.ValidIngredient{
		ID:                ksuid.New().String(),
		Name:              fake.LoremIpsumSentence(exampleQuantity),
		Description:       fake.LoremIpsumSentence(exampleQuantity),
		Warning:           fake.LoremIpsumSentence(exampleQuantity),
		ContainsEgg:       fake.Bool(),
		ContainsDairy:     fake.Bool(),
		ContainsPeanut:    fake.Bool(),
		ContainsTreeNut:   fake.Bool(),
		ContainsSoy:       fake.Bool(),
		ContainsWheat:     fake.Bool(),
		ContainsShellfish: fake.Bool(),
		ContainsSesame:    fake.Bool(),
		ContainsFish:      fake.Bool(),
		ContainsGluten:    fake.Bool(),
		AnimalFlesh:       fake.Bool(),
		Volumetric:        fake.Bool(),
		IconPath:          fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:         uint64(uint32(fake.Date().Unix())),
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
		Name:              validIngredient.Name,
		Description:       validIngredient.Description,
		Warning:           validIngredient.Warning,
		ContainsEgg:       validIngredient.ContainsEgg,
		ContainsDairy:     validIngredient.ContainsDairy,
		ContainsPeanut:    validIngredient.ContainsPeanut,
		ContainsTreeNut:   validIngredient.ContainsTreeNut,
		ContainsSoy:       validIngredient.ContainsSoy,
		ContainsWheat:     validIngredient.ContainsWheat,
		ContainsShellfish: validIngredient.ContainsShellfish,
		ContainsSesame:    validIngredient.ContainsSesame,
		ContainsFish:      validIngredient.ContainsFish,
		ContainsGluten:    validIngredient.ContainsGluten,
		AnimalFlesh:       validIngredient.AnimalFlesh,
		Volumetric:        validIngredient.Volumetric,
		IconPath:          validIngredient.IconPath,
	}
}

// BuildFakeValidIngredientUpdateRequestInputFromValidIngredient builds a faked ValidIngredientUpdateRequestInput from a valid ingredient.
func BuildFakeValidIngredientUpdateRequestInputFromValidIngredient(validIngredient *types.ValidIngredient) *types.ValidIngredientUpdateRequestInput {
	return &types.ValidIngredientUpdateRequestInput{
		Name:              validIngredient.Name,
		Description:       validIngredient.Description,
		Warning:           validIngredient.Warning,
		ContainsEgg:       validIngredient.ContainsEgg,
		ContainsDairy:     validIngredient.ContainsDairy,
		ContainsPeanut:    validIngredient.ContainsPeanut,
		ContainsTreeNut:   validIngredient.ContainsTreeNut,
		ContainsSoy:       validIngredient.ContainsSoy,
		ContainsWheat:     validIngredient.ContainsWheat,
		ContainsShellfish: validIngredient.ContainsShellfish,
		ContainsSesame:    validIngredient.ContainsSesame,
		ContainsFish:      validIngredient.ContainsFish,
		ContainsGluten:    validIngredient.ContainsGluten,
		AnimalFlesh:       validIngredient.AnimalFlesh,
		Volumetric:        validIngredient.Volumetric,
		IconPath:          validIngredient.IconPath,
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
		ID:                validIngredient.ID,
		Name:              validIngredient.Name,
		Description:       validIngredient.Description,
		Warning:           validIngredient.Warning,
		ContainsEgg:       validIngredient.ContainsEgg,
		ContainsDairy:     validIngredient.ContainsDairy,
		ContainsPeanut:    validIngredient.ContainsPeanut,
		ContainsTreeNut:   validIngredient.ContainsTreeNut,
		ContainsSoy:       validIngredient.ContainsSoy,
		ContainsWheat:     validIngredient.ContainsWheat,
		ContainsShellfish: validIngredient.ContainsShellfish,
		ContainsSesame:    validIngredient.ContainsSesame,
		ContainsFish:      validIngredient.ContainsFish,
		ContainsGluten:    validIngredient.ContainsGluten,
		AnimalFlesh:       validIngredient.AnimalFlesh,
		Volumetric:        validIngredient.Volumetric,
		IconPath:          validIngredient.IconPath,
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
		ID:                validIngredient.ID,
		Name:              validIngredient.Name,
		Description:       validIngredient.Description,
		Warning:           validIngredient.Warning,
		ContainsEgg:       validIngredient.ContainsEgg,
		ContainsDairy:     validIngredient.ContainsDairy,
		ContainsPeanut:    validIngredient.ContainsPeanut,
		ContainsTreeNut:   validIngredient.ContainsTreeNut,
		ContainsSoy:       validIngredient.ContainsSoy,
		ContainsWheat:     validIngredient.ContainsWheat,
		ContainsShellfish: validIngredient.ContainsShellfish,
		ContainsSesame:    validIngredient.ContainsSesame,
		ContainsFish:      validIngredient.ContainsFish,
		ContainsGluten:    validIngredient.ContainsGluten,
		AnimalFlesh:       validIngredient.AnimalFlesh,
		Volumetric:        validIngredient.Volumetric,
		IconPath:          validIngredient.IconPath,
	}
}
