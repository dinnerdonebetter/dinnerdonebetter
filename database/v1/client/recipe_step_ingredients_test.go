package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleRecipeStepIngredientID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.RecipeStepIngredient{}

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipeStepIngredientID, exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepIngredient(context.Background(), exampleRecipeStepIngredientID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeStepIngredientCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("GetRecipeStepIngredientCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepIngredientCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("GetRecipeStepIngredientCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepIngredientCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllRecipeStepIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("GetAllRecipeStepIngredientsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllRecipeStepIngredientsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepIngredientList{}

		mockDB.RecipeStepIngredientDataManager.On("GetRecipeStepIngredients", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepIngredients(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepIngredientList{}

		mockDB.RecipeStepIngredientDataManager.On("GetRecipeStepIngredients", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepIngredients(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStepIngredientCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepIngredient{}

		mockDB.RecipeStepIngredientDataManager.On("CreateRecipeStepIngredient", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateRecipeStepIngredient(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStepIngredient{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.RecipeStepIngredientDataManager.On("UpdateRecipeStepIngredient", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateRecipeStepIngredient(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleRecipeStepIngredientID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.RecipeStepIngredientDataManager.On("ArchiveRecipeStepIngredient", mock.Anything, exampleRecipeStepIngredientID, exampleUserID).Return(expected)

		err := c.ArchiveRecipeStepIngredient(context.Background(), exampleUserID, exampleRecipeStepIngredientID)
		assert.NoError(t, err)
	})
}
