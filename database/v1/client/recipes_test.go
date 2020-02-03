package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetRecipe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleRecipeID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.Recipe{}

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("GetRecipe", mock.Anything, exampleRecipeID, exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipe(context.Background(), exampleRecipeID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("GetRecipeCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("GetRecipeCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllRecipesCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("GetAllRecipesCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllRecipesCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipes(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeList{}

		mockDB.RecipeDataManager.On("GetRecipes", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipes(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeList{}

		mockDB.RecipeDataManager.On("GetRecipes", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipes(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.Recipe{}

		mockDB.RecipeDataManager.On("CreateRecipe", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateRecipe(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.Recipe{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.RecipeDataManager.On("UpdateRecipe", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateRecipe(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveRecipe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleRecipeID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.RecipeDataManager.On("ArchiveRecipe", mock.Anything, exampleRecipeID, exampleUserID).Return(expected)

		err := c.ArchiveRecipe(context.Background(), exampleUserID, exampleRecipeID)
		assert.NoError(t, err)
	})
}
