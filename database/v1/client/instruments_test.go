package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInstrumentID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.Instrument{}

		c, mockDB := buildTestClient()
		mockDB.InstrumentDataManager.On("GetInstrument", mock.Anything, exampleInstrumentID, exampleUserID).Return(expected, nil)

		actual, err := c.GetInstrument(context.Background(), exampleInstrumentID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetInstrumentCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.InstrumentDataManager.On("GetInstrumentCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetInstrumentCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.InstrumentDataManager.On("GetInstrumentCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetInstrumentCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.InstrumentDataManager.On("GetAllInstrumentsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllInstrumentsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetInstruments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.InstrumentList{}

		mockDB.InstrumentDataManager.On("GetInstruments", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetInstruments(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.InstrumentList{}

		mockDB.InstrumentDataManager.On("GetInstruments", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetInstruments(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.InstrumentCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.Instrument{}

		mockDB.InstrumentDataManager.On("CreateInstrument", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateInstrument(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.Instrument{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.InstrumentDataManager.On("UpdateInstrument", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateInstrument(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleInstrumentID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.InstrumentDataManager.On("ArchiveInstrument", mock.Anything, exampleInstrumentID, exampleUserID).Return(expected)

		err := c.ArchiveInstrument(context.Background(), exampleUserID, exampleInstrumentID)
		assert.NoError(t, err)
	})
}
