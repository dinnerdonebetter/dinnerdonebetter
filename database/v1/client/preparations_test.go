package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		examplePreparationID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.Preparation{}

		c, mockDB := buildTestClient()
		mockDB.PreparationDataManager.On("GetPreparation", mock.Anything, examplePreparationID, exampleUserID).Return(expected, nil)

		actual, err := c.GetPreparation(context.Background(), examplePreparationID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetPreparationCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.PreparationDataManager.On("GetPreparationCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetPreparationCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.PreparationDataManager.On("GetPreparationCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetPreparationCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllPreparationsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.PreparationDataManager.On("GetAllPreparationsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllPreparationsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetPreparations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.PreparationList{}

		mockDB.PreparationDataManager.On("GetPreparations", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetPreparations(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.PreparationList{}

		mockDB.PreparationDataManager.On("GetPreparations", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetPreparations(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreatePreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.PreparationCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.Preparation{}

		mockDB.PreparationDataManager.On("CreatePreparation", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreatePreparation(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdatePreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.Preparation{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.PreparationDataManager.On("UpdatePreparation", mock.Anything, exampleInput).Return(expected)

		err := c.UpdatePreparation(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchivePreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		examplePreparationID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.PreparationDataManager.On("ArchivePreparation", mock.Anything, examplePreparationID, exampleUserID).Return(expected)

		err := c.ArchivePreparation(context.Background(), exampleUserID, examplePreparationID)
		assert.NoError(t, err)
	})
}
