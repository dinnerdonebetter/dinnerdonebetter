package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
)

// BuildFakeRecipeRating builds a faked valid ingredient.
func BuildFakeRecipeRating() *types.RecipeRating {
	return &types.RecipeRating{
		CreatedAt:    BuildFakeTime(),
		Notes:        buildUniqueString(),
		ID:           buildUniqueString(),
		RecipeID:     buildUniqueString(),
		ByUser:       buildUniqueString(),
		Taste:        float32(buildFakeNumber()),
		Instructions: float32(buildFakeNumber()),
		Overall:      float32(buildFakeNumber()),
		Cleanup:      float32(buildFakeNumber()),
		Difficulty:   float32(buildFakeNumber()),
	}
}

// BuildFakeRecipeRatingsList builds a faked RecipeRatingList.
func BuildFakeRecipeRatingsList() *filtering.QueryFilteredResult[types.RecipeRating] {
	var examples []*types.RecipeRating
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeRating())
	}

	return &filtering.QueryFilteredResult[types.RecipeRating]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeRatingUpdateRequestInput builds a faked RecipeRatingUpdateRequestInput from a valid ingredient.
func BuildFakeRecipeRatingUpdateRequestInput() *types.RecipeRatingUpdateRequestInput {
	recipeRating := BuildFakeRecipeRating()
	return converters.ConvertRecipeRatingToRecipeRatingUpdateRequestInput(recipeRating)
}

// BuildFakeRecipeRatingCreationRequestInput builds a faked RecipeRatingCreationRequestInput.
func BuildFakeRecipeRatingCreationRequestInput() *types.RecipeRatingCreationRequestInput {
	recipeRating := BuildFakeRecipeRating()
	return converters.ConvertRecipeRatingToRecipeRatingCreationRequestInput(recipeRating)
}
