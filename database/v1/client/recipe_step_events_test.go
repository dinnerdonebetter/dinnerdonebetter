package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleRecipeStepEventID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.RecipeStepEvent{}

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipeStepEventID, exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepEvent(context.Background(), exampleRecipeStepEventID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeStepEventCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetRecipeStepEventCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepEventCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetRecipeStepEventCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepEventCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllRecipeStepEventsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetAllRecipeStepEventsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllRecipeStepEventsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeStepEvents(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepEventList{}

		mockDB.RecipeStepEventDataManager.On("GetRecipeStepEvents", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepEvents(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepEventList{}

		mockDB.RecipeStepEventDataManager.On("GetRecipeStepEvents", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepEvents(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStepEventCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepEvent{}

		mockDB.RecipeStepEventDataManager.On("CreateRecipeStepEvent", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateRecipeStepEvent(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStepEvent{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.RecipeStepEventDataManager.On("UpdateRecipeStepEvent", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateRecipeStepEvent(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleRecipeStepEventID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("ArchiveRecipeStepEvent", mock.Anything, exampleRecipeStepEventID, exampleUserID).Return(expected)

		err := c.ArchiveRecipeStepEvent(context.Background(), exampleUserID, exampleRecipeStepEventID)
		assert.NoError(t, err)
	})
}
