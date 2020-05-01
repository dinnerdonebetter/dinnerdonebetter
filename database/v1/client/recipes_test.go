package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)

		actual, err := c.RecipeExists(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("GetRecipe", mock.Anything, exampleRecipe.ID).Return(exampleRecipe, nil)

		actual, err := c.GetRecipe(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipe, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipesCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("GetAllRecipesCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipesCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipes(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := models.DefaultQueryFilter()
		exampleRecipeList := fakemodels.BuildFakeRecipeList()

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("GetRecipes", mock.Anything, filter).Return(exampleRecipeList, nil)

		actual, err := c.GetRecipes(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		exampleRecipeList := fakemodels.BuildFakeRecipeList()

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("GetRecipes", mock.Anything, filter).Return(exampleRecipeList, nil)

		actual, err := c.GetRecipes(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("CreateRecipe", mock.Anything, exampleInput).Return(exampleRecipe, nil)

		actual, err := c.CreateRecipe(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipe, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()

		mockDB.RecipeDataManager.On("UpdateRecipe", mock.Anything, exampleRecipe).Return(expected)

		err := c.UpdateRecipe(ctx, exampleRecipe)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("ArchiveRecipe", mock.Anything, exampleRecipe.ID, exampleRecipe.BelongsToUser).Return(expected)

		err := c.ArchiveRecipe(ctx, exampleRecipe.ID, exampleRecipe.BelongsToUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
