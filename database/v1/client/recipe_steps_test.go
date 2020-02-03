package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleRecipeStepID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.RecipeStep{}

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("GetRecipeStep", mock.Anything, exampleRecipeStepID, exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStep(context.Background(), exampleRecipeStepID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeStepCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("GetRecipeStepCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("GetRecipeStepCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllRecipeStepsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("GetAllRecipeStepsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllRecipeStepsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepList{}

		mockDB.RecipeStepDataManager.On("GetRecipeSteps", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepList{}

		mockDB.RecipeStepDataManager.On("GetRecipeSteps", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeSteps(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStepCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.RecipeStep{}

		mockDB.RecipeStepDataManager.On("CreateRecipeStep", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateRecipeStep(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStep{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.RecipeStepDataManager.On("UpdateRecipeStep", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateRecipeStep(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleRecipeStepID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("ArchiveRecipeStep", mock.Anything, exampleRecipeStepID, exampleUserID).Return(expected)

		err := c.ArchiveRecipeStep(context.Background(), exampleUserID, exampleRecipeStepID)
		assert.NoError(t, err)
	})
}
