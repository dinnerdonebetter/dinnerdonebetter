package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ValidPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		c, mockDB := buildTestClient()
		mockDB.ValidPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, nil)

		actual, err := c.ValidPreparationExists(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		c, mockDB := buildTestClient()
		mockDB.ValidPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)

		actual, err := c.GetValidPreparation(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparation, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllValidPreparationsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ValidPreparationDataManager.On("GetAllValidPreparationsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllValidPreparationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := models.DefaultQueryFilter()
		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList()

		c, mockDB := buildTestClient()
		mockDB.ValidPreparationDataManager.On("GetValidPreparations", mock.Anything, filter).Return(exampleValidPreparationList, nil)

		actual, err := c.GetValidPreparations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList()

		c, mockDB := buildTestClient()
		mockDB.ValidPreparationDataManager.On("GetValidPreparations", mock.Anything, filter).Return(exampleValidPreparationList, nil)

		actual, err := c.GetValidPreparations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		c, mockDB := buildTestClient()
		mockDB.ValidPreparationDataManager.On("CreateValidPreparation", mock.Anything, exampleInput).Return(exampleValidPreparation, nil)

		actual, err := c.CreateValidPreparation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparation, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		c, mockDB := buildTestClient()

		mockDB.ValidPreparationDataManager.On("UpdateValidPreparation", mock.Anything, exampleValidPreparation).Return(expected)

		err := c.UpdateValidPreparation(ctx, exampleValidPreparation)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		c, mockDB := buildTestClient()
		mockDB.ValidPreparationDataManager.On("ArchiveValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(expected)

		err := c.ArchiveValidPreparation(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
