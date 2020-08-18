package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RequiredPreparationInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("RequiredPreparationInstrumentExists", mock.Anything, exampleRequiredPreparationInstrument.ID).Return(true, nil)

		actual, err := c.RequiredPreparationInstrumentExists(ctx, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleRequiredPreparationInstrument.ID).Return(exampleRequiredPreparationInstrument, nil)

		actual, err := c.GetRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrument, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRequiredPreparationInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetAllRequiredPreparationInstrumentsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRequiredPreparationInstrumentsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		results := make(chan []models.RequiredPreparationInstrument)

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetAllRequiredPreparationInstruments", mock.Anything, results).Return(nil)

		err := c.GetAllRequiredPreparationInstruments(ctx, results)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := models.DefaultQueryFilter()
		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, filter).Return(exampleRequiredPreparationInstrumentList, nil)

		actual, err := c.GetRequiredPreparationInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, filter).Return(exampleRequiredPreparationInstrumentList, nil)

		actual, err := c.GetRequiredPreparationInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRequiredPreparationInstrumentsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList().RequiredPreparationInstruments
		var exampleIDs []uint64
		for _, x := range exampleRequiredPreparationInstrumentList {
			exampleIDs = append(exampleIDs, x.ID)
		}

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrumentsWithIDs", mock.Anything, defaultLimit, exampleIDs).Return(exampleRequiredPreparationInstrumentList, nil)

		actual, err := c.GetRequiredPreparationInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("CreateRequiredPreparationInstrument", mock.Anything, exampleInput).Return(exampleRequiredPreparationInstrument, nil)

		actual, err := c.CreateRequiredPreparationInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrument, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		c, mockDB := buildTestClient()

		mockDB.RequiredPreparationInstrumentDataManager.On("UpdateRequiredPreparationInstrument", mock.Anything, exampleRequiredPreparationInstrument).Return(expected)

		err := c.UpdateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("ArchiveRequiredPreparationInstrument", mock.Anything, exampleRequiredPreparationInstrument.ID).Return(expected)

		err := c.ArchiveRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
