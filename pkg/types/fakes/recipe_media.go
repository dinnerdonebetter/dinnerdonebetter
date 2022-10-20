package fakes

import (
	"fmt"

	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
)

// BuildFakeRecipeMedia builds a faked valid preparation.
func BuildFakeRecipeMedia() *types.RecipeMedia {
	return &types.RecipeMedia{
		ID:                  BuildFakeID(),
		BelongsToRecipe:     nil,
		BelongsToRecipeStep: nil,
		MimeType:            fake.FileMimeType(),
		InternalPath:        fmt.Sprintf("%s.%s", BuildFakePassword(), fake.FileExtension()),
		ExternalPath:        "",
		CreatedAt:           fake.Date(),
	}
}

// BuildFakeRecipeMediaList builds a faked RecipeMediaList.
func BuildFakeRecipeMediaList() *types.RecipeMediaList {
	var examples []*types.RecipeMedia
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeMedia())
	}

	return &types.RecipeMediaList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		RecipeMedia: examples,
	}
}

// BuildFakeRecipeMediaUpdateRequestInput builds a faked RecipeMediaUpdateRequestInput from a valid preparation.
func BuildFakeRecipeMediaUpdateRequestInput() *types.RecipeMediaUpdateRequestInput {
	validPreparation := BuildFakeRecipeMedia()
	return &types.RecipeMediaUpdateRequestInput{
		BelongsToRecipe:     validPreparation.BelongsToRecipe,
		BelongsToRecipeStep: validPreparation.BelongsToRecipeStep,
		MimeType:            &validPreparation.MimeType,
		InternalPath:        &validPreparation.InternalPath,
		ExternalPath:        &validPreparation.ExternalPath,
	}
}

// BuildFakeRecipeMediaCreationRequestInput builds a faked RecipeMediaCreationRequestInput.
func BuildFakeRecipeMediaCreationRequestInput() *types.RecipeMediaCreationRequestInput {
	validPreparation := BuildFakeRecipeMedia()
	return converters.ConvertRecipeMediaToRecipeMediaCreationRequestInput(validPreparation)
}
