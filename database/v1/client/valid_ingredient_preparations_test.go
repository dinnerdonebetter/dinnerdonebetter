package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ValidIngredientPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("ValidIngredientPreparationExists", mock.Anything, exampleValidIngredientPreparation.ID).Return(true, nil)

		actual, err := c.ValidIngredientPreparationExists(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(exampleValidIngredientPreparation, nil)

		actual, err := c.GetValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparation, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllValidIngredientPreparationsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("GetAllValidIngredientPreparationsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllValidIngredientPreparationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		results := make(chan []models.ValidIngredientPreparation)

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("GetAllValidIngredientPreparations", mock.Anything, results).Return(nil)

		err := c.GetAllValidIngredientPreparations(ctx, results)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := models.DefaultQueryFilter()
		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("GetValidIngredientPreparations", mock.Anything, filter).Return(exampleValidIngredientPreparationList, nil)

		actual, err := c.GetValidIngredientPreparations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("GetValidIngredientPreparations", mock.Anything, filter).Return(exampleValidIngredientPreparationList, nil)

		actual, err := c.GetValidIngredientPreparations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidIngredientPreparationsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList().ValidIngredientPreparations
		var exampleIDs []uint64
		for _, x := range exampleValidIngredientPreparationList {
			exampleIDs = append(exampleIDs, x.ID)
		}

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("GetValidIngredientPreparationsWithIDs", mock.Anything, defaultLimit, exampleIDs).Return(exampleValidIngredientPreparationList, nil)

		actual, err := c.GetValidIngredientPreparationsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("CreateValidIngredientPreparation", mock.Anything, exampleInput).Return(exampleValidIngredientPreparation, nil)

		actual, err := c.CreateValidIngredientPreparation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparation, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		c, mockDB := buildTestClient()

		mockDB.ValidIngredientPreparationDataManager.On("UpdateValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation).Return(expected)

		err := c.UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("ArchiveValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(expected)

		err := c.ArchiveValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
