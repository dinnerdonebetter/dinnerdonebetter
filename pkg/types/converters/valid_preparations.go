package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertValidPreparationToValidPreparationUpdateRequestInput creates a ValidPreparationUpdateRequestInput from a ValidPreparation.
func ConvertValidPreparationToValidPreparationUpdateRequestInput(input *types.ValidPreparation) *types.ValidPreparationUpdateRequestInput {
	x := &types.ValidPreparationUpdateRequestInput{
		Name:                     &input.Name,
		Description:              &input.Description,
		IconPath:                 &input.IconPath,
		YieldsNothing:            &input.YieldsNothing,
		RestrictToIngredients:    &input.RestrictToIngredients,
		ZeroIngredientsAllowable: &input.ZeroIngredientsAllowable,
		Slug:                     &input.Slug,
		PastTense:                &input.PastTense,
	}

	return x
}

// ConvertValidPreparationCreationRequestInputToValidPreparationDatabaseCreationInput creates a ValidPreparationDatabaseCreationInput from a ValidPreparationCreationRequestInput.
func ConvertValidPreparationCreationRequestInputToValidPreparationDatabaseCreationInput(input *types.ValidPreparationCreationRequestInput) *types.ValidPreparationDatabaseCreationInput {
	x := &types.ValidPreparationDatabaseCreationInput{
		ID:                       identifiers.New(),
		Name:                     input.Name,
		Description:              input.Description,
		IconPath:                 input.IconPath,
		YieldsNothing:            input.YieldsNothing,
		RestrictToIngredients:    input.RestrictToIngredients,
		ZeroIngredientsAllowable: input.ZeroIngredientsAllowable,
		Slug:                     input.Slug,
		PastTense:                input.PastTense,
	}

	return x
}

// ConvertValidPreparationToValidPreparationCreationRequestInput builds a ValidPreparationCreationRequestInput from a ValidPreparation.
func ConvertValidPreparationToValidPreparationCreationRequestInput(validPreparation *types.ValidPreparation) *types.ValidPreparationCreationRequestInput {
	return &types.ValidPreparationCreationRequestInput{
		Name:                     validPreparation.Name,
		Description:              validPreparation.Description,
		IconPath:                 validPreparation.IconPath,
		YieldsNothing:            validPreparation.YieldsNothing,
		RestrictToIngredients:    validPreparation.RestrictToIngredients,
		ZeroIngredientsAllowable: validPreparation.ZeroIngredientsAllowable,
		Slug:                     validPreparation.Slug,
		PastTense:                validPreparation.PastTense,
	}
}

// ConvertValidPreparationToValidPreparationDatabaseCreationInput builds a ValidPreparationDatabaseCreationInput from a ValidPreparation.
func ConvertValidPreparationToValidPreparationDatabaseCreationInput(validPreparation *types.ValidPreparation) *types.ValidPreparationDatabaseCreationInput {
	return &types.ValidPreparationDatabaseCreationInput{
		ID:                       validPreparation.ID,
		Name:                     validPreparation.Name,
		Description:              validPreparation.Description,
		IconPath:                 validPreparation.IconPath,
		YieldsNothing:            validPreparation.YieldsNothing,
		RestrictToIngredients:    validPreparation.RestrictToIngredients,
		ZeroIngredientsAllowable: validPreparation.ZeroIngredientsAllowable,
		Slug:                     validPreparation.Slug,
		PastTense:                validPreparation.PastTense,
	}
}
