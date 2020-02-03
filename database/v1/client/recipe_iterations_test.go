package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleRecipeIterationID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.RecipeIteration{}

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipeIterationID, exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeIteration(context.Background(), exampleRecipeIterationID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeIterationCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetRecipeIterationCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeIterationCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetRecipeIterationCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeIterationCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllRecipeIterationsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetAllRecipeIterationsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllRecipeIterationsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeIterations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeIterationList{}

		mockDB.RecipeIterationDataManager.On("GetRecipeIterations", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeIterationList{}

		mockDB.RecipeIterationDataManager.On("GetRecipeIterations", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeIterations(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeIterationCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.RecipeIteration{}

		mockDB.RecipeIterationDataManager.On("CreateRecipeIteration", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateRecipeIteration(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeIteration{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.RecipeIterationDataManager.On("UpdateRecipeIteration", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateRecipeIteration(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleRecipeIterationID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("ArchiveRecipeIteration", mock.Anything, exampleRecipeIterationID, exampleUserID).Return(expected)

		err := c.ArchiveRecipeIteration(context.Background(), exampleUserID, exampleRecipeIterationID)
		assert.NoError(t, err)
	})
}
