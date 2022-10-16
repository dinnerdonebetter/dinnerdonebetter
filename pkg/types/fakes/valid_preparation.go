package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
)

// BuildFakeValidPreparation builds a faked valid preparation.
func BuildFakeValidPreparation() *types.ValidPreparation {
	return &types.ValidPreparation{
		ID:                       BuildFakeID(),
		Name:                     buildUniqueString(),
		Description:              buildUniqueString(),
		IconPath:                 buildUniqueString(),
		YieldsNothing:            fake.Bool(),
		RestrictToIngredients:    fake.Bool(),
		ZeroIngredientsAllowable: fake.Bool(),
		PastTense:                buildUniqueString(),
		CreatedAt:                fake.Date(),
	}
}

// BuildFakeValidPreparationList builds a faked ValidPreparationList.
func BuildFakeValidPreparationList() *types.ValidPreparationList {
	var examples []*types.ValidPreparation
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidPreparation())
	}

	return &types.ValidPreparationList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidPreparations: examples,
	}
}

// BuildFakeValidPreparationUpdateRequestInput builds a faked ValidPreparationUpdateRequestInput from a valid preparation.
func BuildFakeValidPreparationUpdateRequestInput() *types.ValidPreparationUpdateRequestInput {
	validPreparation := BuildFakeValidPreparation()
	return &types.ValidPreparationUpdateRequestInput{
		Name:                     &validPreparation.Name,
		Description:              &validPreparation.Description,
		IconPath:                 &validPreparation.IconPath,
		YieldsNothing:            &validPreparation.YieldsNothing,
		RestrictToIngredients:    &validPreparation.RestrictToIngredients,
		ZeroIngredientsAllowable: &validPreparation.ZeroIngredientsAllowable,
		PastTense:                &validPreparation.PastTense,
	}
}

// BuildFakeValidPreparationCreationRequestInput builds a faked ValidPreparationCreationRequestInput.
func BuildFakeValidPreparationCreationRequestInput() *types.ValidPreparationCreationRequestInput {
	validPreparation := BuildFakeValidPreparation()
	return converters.ConvertValidPreparationToValidPreparationCreationRequestInput(validPreparation)
}
