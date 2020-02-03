package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleRecipeStepInstrumentID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.RecipeStepInstrument{}

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipeStepInstrumentID, exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepInstrument(context.Background(), exampleRecipeStepInstrumentID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeStepInstrumentCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetRecipeStepInstrumentCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepInstrumentCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetRecipeStepInstrumentCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepInstrumentCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllRecipeStepInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetAllRecipeStepInstrumentsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllRecipeStepInstrumentsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepInstrumentList{}

		mockDB.RecipeStepInstrumentDataManager.On("GetRecipeStepInstruments", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepInstruments(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepInstrumentList{}

		mockDB.RecipeStepInstrumentDataManager.On("GetRecipeStepInstruments", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRecipeStepInstruments(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStepInstrumentCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.RecipeStepInstrument{}

		mockDB.RecipeStepInstrumentDataManager.On("CreateRecipeStepInstrument", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateRecipeStepInstrument(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RecipeStepInstrument{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.RecipeStepInstrumentDataManager.On("UpdateRecipeStepInstrument", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateRecipeStepInstrument(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleRecipeStepInstrumentID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("ArchiveRecipeStepInstrument", mock.Anything, exampleRecipeStepInstrumentID, exampleUserID).Return(expected)

		err := c.ArchiveRecipeStepInstrument(context.Background(), exampleUserID, exampleRecipeStepInstrumentID)
		assert.NoError(t, err)
	})
}
