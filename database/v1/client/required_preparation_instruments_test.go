package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleRequiredPreparationInstrumentID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.RequiredPreparationInstrument{}

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleRequiredPreparationInstrumentID, exampleUserID).Return(expected, nil)

		actual, err := c.GetRequiredPreparationInstrument(context.Background(), exampleRequiredPreparationInstrumentID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRequiredPreparationInstrumentCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrumentCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRequiredPreparationInstrumentCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrumentCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRequiredPreparationInstrumentCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllRequiredPreparationInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetAllRequiredPreparationInstrumentsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllRequiredPreparationInstrumentsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RequiredPreparationInstrumentList{}

		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetRequiredPreparationInstruments(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.RequiredPreparationInstrumentList{}

		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetRequiredPreparationInstruments(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RequiredPreparationInstrumentCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.RequiredPreparationInstrument{}

		mockDB.RequiredPreparationInstrumentDataManager.On("CreateRequiredPreparationInstrument", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateRequiredPreparationInstrument(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.RequiredPreparationInstrument{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.RequiredPreparationInstrumentDataManager.On("UpdateRequiredPreparationInstrument", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateRequiredPreparationInstrument(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleRequiredPreparationInstrumentID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("ArchiveRequiredPreparationInstrument", mock.Anything, exampleRequiredPreparationInstrumentID, exampleUserID).Return(expected)

		err := c.ArchiveRequiredPreparationInstrument(context.Background(), exampleUserID, exampleRequiredPreparationInstrumentID)
		assert.NoError(t, err)
	})
}
