package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeTagExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeTagDataManager.On("RecipeTagExists", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(true, nil)

		actual, err := c.RecipeTagExists(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeTag(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(exampleRecipeTag, nil)

		actual, err := c.GetRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeTag, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeTagsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeTagDataManager.On("GetAllRecipeTagsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipeTagsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeTags(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := models.DefaultQueryFilter()
		exampleRecipeTagList := fakemodels.BuildFakeRecipeTagList()

		c, mockDB := buildTestClient()
		mockDB.RecipeTagDataManager.On("GetRecipeTags", mock.Anything, exampleRecipe.ID, filter).Return(exampleRecipeTagList, nil)

		actual, err := c.GetRecipeTags(ctx, exampleRecipe.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeTagList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)
		exampleRecipeTagList := fakemodels.BuildFakeRecipeTagList()

		c, mockDB := buildTestClient()
		mockDB.RecipeTagDataManager.On("GetRecipeTags", mock.Anything, exampleRecipe.ID, filter).Return(exampleRecipeTagList, nil)

		actual, err := c.GetRecipeTags(ctx, exampleRecipe.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeTagList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipeTag(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

		c, mockDB := buildTestClient()
		mockDB.RecipeTagDataManager.On("CreateRecipeTag", mock.Anything, exampleInput).Return(exampleRecipeTag, nil)

		actual, err := c.CreateRecipeTag(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeTag, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipeTag(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()

		c, mockDB := buildTestClient()

		mockDB.RecipeTagDataManager.On("UpdateRecipeTag", mock.Anything, exampleRecipeTag).Return(expected)

		err := c.UpdateRecipeTag(ctx, exampleRecipeTag)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipeTag(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()

		c, mockDB := buildTestClient()
		mockDB.RecipeTagDataManager.On("ArchiveRecipeTag", mock.Anything, exampleRecipeTag.BelongsToRecipe, exampleRecipeTag.ID).Return(expected)

		err := c.ArchiveRecipeTag(ctx, exampleRecipeTag.BelongsToRecipe, exampleRecipeTag.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
