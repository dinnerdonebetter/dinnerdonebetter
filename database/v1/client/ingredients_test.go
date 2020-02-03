package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleIngredientID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.Ingredient{}

		c, mockDB := buildTestClient()
		mockDB.IngredientDataManager.On("GetIngredient", mock.Anything, exampleIngredientID, exampleUserID).Return(expected, nil)

		actual, err := c.GetIngredient(context.Background(), exampleIngredientID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetIngredientCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.IngredientDataManager.On("GetIngredientCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetIngredientCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.IngredientDataManager.On("GetIngredientCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetIngredientCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.IngredientDataManager.On("GetAllIngredientsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllIngredientsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetIngredients(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.IngredientList{}

		mockDB.IngredientDataManager.On("GetIngredients", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetIngredients(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.IngredientList{}

		mockDB.IngredientDataManager.On("GetIngredients", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetIngredients(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.IngredientCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.Ingredient{}

		mockDB.IngredientDataManager.On("CreateIngredient", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateIngredient(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.Ingredient{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.IngredientDataManager.On("UpdateIngredient", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateIngredient(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleIngredientID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.IngredientDataManager.On("ArchiveIngredient", mock.Anything, exampleIngredientID, exampleUserID).Return(expected)

		err := c.ArchiveIngredient(context.Background(), exampleUserID, exampleIngredientID)
		assert.NoError(t, err)
	})
}
