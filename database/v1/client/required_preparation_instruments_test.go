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

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("RequiredPreparationInstrumentExists", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(true, nil)

		actual, err := c.RequiredPreparationInstrumentExists(ctx, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(exampleRequiredPreparationInstrument, nil)

		actual, err := c.GetRequiredPreparationInstrument(ctx, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID)
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

func TestClient_GetRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		filter := models.DefaultQueryFilter()
		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, exampleValidPreparation.ID, filter).Return(exampleRequiredPreparationInstrumentList, nil)

		actual, err := c.GetRequiredPreparationInstruments(ctx, exampleValidPreparation.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		filter := (*models.QueryFilter)(nil)
		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()

		c, mockDB := buildTestClient()
		mockDB.RequiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, exampleValidPreparation.ID, filter).Return(exampleRequiredPreparationInstrumentList, nil)

		actual, err := c.GetRequiredPreparationInstruments(ctx, exampleValidPreparation.ID, filter)
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
		mockDB.RequiredPreparationInstrumentDataManager.On("ArchiveRequiredPreparationInstrument", mock.Anything, exampleRequiredPreparationInstrument.BelongsToValidPreparation, exampleRequiredPreparationInstrument.ID).Return(expected)

		err := c.ArchiveRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument.BelongsToValidPreparation, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
