package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleRecipeStepProductID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.RecipeStepProduct{}

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipeStepProductID, exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepProduct(context.Background(), exampleRecipeStepProductID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeStepProductCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetRecipeStepProductCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepProductCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetRecipeStepProductCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepProductCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllRecipeStepProductsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetAllRecipeStepProductsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllRecipeStepProductsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepProductList{}

		mockDB.RecipeStepProductDataManager.On("GetRecipeStepProducts", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepProducts(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepProductList{}

		mockDB.RecipeStepProductDataManager.On("GetRecipeStepProducts", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepProducts(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStepProductCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepProduct{}

		mockDB.RecipeStepProductDataManager.On("CreateRecipeStepProduct", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateRecipeStepProduct(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStepProduct{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.RecipeStepProductDataManager.On("UpdateRecipeStepProduct", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateRecipeStepProduct(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleRecipeStepProductID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("ArchiveRecipeStepProduct", mock.Anything, exampleRecipeStepProductID, exampleUserID).Return(expected)

		err := c.ArchiveRecipeStepProduct(context.Background(), exampleUserID, exampleRecipeStepProductID)
		assert.NoError(t, err)
	})
}
