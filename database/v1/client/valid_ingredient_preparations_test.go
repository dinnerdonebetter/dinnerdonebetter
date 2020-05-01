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

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("ValidIngredientPreparationExists", mock.Anything, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID).Return(true, nil)

		actual, err := c.ValidIngredientPreparationExists(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID).Return(exampleValidIngredientPreparation, nil)

		actual, err := c.GetValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
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

func TestClient_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := models.DefaultQueryFilter()
		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("GetValidIngredientPreparations", mock.Anything, exampleValidIngredient.ID, filter).Return(exampleValidIngredientPreparationList, nil)

		actual, err := c.GetValidIngredientPreparations(ctx, exampleValidIngredient.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := (*models.QueryFilter)(nil)
		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientPreparationDataManager.On("GetValidIngredientPreparations", mock.Anything, exampleValidIngredient.ID, filter).Return(exampleValidIngredientPreparationList, nil)

		actual, err := c.GetValidIngredientPreparations(ctx, exampleValidIngredient.ID, filter)
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
		mockDB.ValidIngredientPreparationDataManager.On("ArchiveValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.BelongsToValidIngredient, exampleValidIngredientPreparation.ID).Return(expected)

		err := c.ArchiveValidIngredientPreparation(ctx, exampleValidIngredientPreparation.BelongsToValidIngredient, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
