package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ValidInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		c, mockDB := buildTestClient()
		mockDB.ValidInstrumentDataManager.On("ValidInstrumentExists", mock.Anything, exampleValidInstrument.ID).Return(true, nil)

		actual, err := c.ValidInstrumentExists(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		c, mockDB := buildTestClient()
		mockDB.ValidInstrumentDataManager.On("GetValidInstrument", mock.Anything, exampleValidInstrument.ID).Return(exampleValidInstrument, nil)

		actual, err := c.GetValidInstrument(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllValidInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ValidInstrumentDataManager.On("GetAllValidInstrumentsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllValidInstrumentsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := models.DefaultQueryFilter()
		exampleValidInstrumentList := fakemodels.BuildFakeValidInstrumentList()

		c, mockDB := buildTestClient()
		mockDB.ValidInstrumentDataManager.On("GetValidInstruments", mock.Anything, filter).Return(exampleValidInstrumentList, nil)

		actual, err := c.GetValidInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		exampleValidInstrumentList := fakemodels.BuildFakeValidInstrumentList()

		c, mockDB := buildTestClient()
		mockDB.ValidInstrumentDataManager.On("GetValidInstruments", mock.Anything, filter).Return(exampleValidInstrumentList, nil)

		actual, err := c.GetValidInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
		exampleInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		c, mockDB := buildTestClient()
		mockDB.ValidInstrumentDataManager.On("CreateValidInstrument", mock.Anything, exampleInput).Return(exampleValidInstrument, nil)

		actual, err := c.CreateValidInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		c, mockDB := buildTestClient()

		mockDB.ValidInstrumentDataManager.On("UpdateValidInstrument", mock.Anything, exampleValidInstrument).Return(expected)

		err := c.UpdateValidInstrument(ctx, exampleValidInstrument)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		c, mockDB := buildTestClient()
		mockDB.ValidInstrumentDataManager.On("ArchiveValidInstrument", mock.Anything, exampleValidInstrument.ID).Return(expected)

		err := c.ArchiveValidInstrument(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
