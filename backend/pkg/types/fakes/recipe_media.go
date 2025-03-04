package fakes

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeRecipeMedia builds a faked valid preparation.
func BuildFakeRecipeMedia() *types.RecipeMedia {
	return &types.RecipeMedia{
		ID:                  BuildFakeID(),
		BelongsToRecipe:     nil,
		BelongsToRecipeStep: nil,
		MimeType:            fake.FileMimeType(),
		InternalPath:        fmt.Sprintf("%s.%s", buildFakePassword(), fake.FileExtension()),
		ExternalPath:        "",
		CreatedAt:           BuildFakeTime(),
	}
}

// BuildFakeRecipeMediaList builds a faked RecipeMediaList.
func BuildFakeRecipeMediaList() *filtering.QueryFilteredResult[types.RecipeMedia] {
	var examples []*types.RecipeMedia
	for range exampleQuantity {
		examples = append(examples, BuildFakeRecipeMedia())
	}

	return &filtering.QueryFilteredResult[types.RecipeMedia]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
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
