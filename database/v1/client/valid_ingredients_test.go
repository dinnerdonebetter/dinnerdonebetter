package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ValidIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, nil)

		actual, err := c.ValidIngredientExists(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)

		actual, err := c.GetValidIngredient(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredient, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllValidIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientDataManager.On("GetAllValidIngredientsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllValidIngredientsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := models.DefaultQueryFilter()
		exampleValidIngredientList := fakemodels.BuildFakeValidIngredientList()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientDataManager.On("GetValidIngredients", mock.Anything, filter).Return(exampleValidIngredientList, nil)

		actual, err := c.GetValidIngredients(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		exampleValidIngredientList := fakemodels.BuildFakeValidIngredientList()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientDataManager.On("GetValidIngredients", mock.Anything, filter).Return(exampleValidIngredientList, nil)

		actual, err := c.GetValidIngredients(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientDataManager.On("CreateValidIngredient", mock.Anything, exampleInput).Return(exampleValidIngredient, nil)

		actual, err := c.CreateValidIngredient(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredient, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		c, mockDB := buildTestClient()

		mockDB.ValidIngredientDataManager.On("UpdateValidIngredient", mock.Anything, exampleValidIngredient).Return(expected)

		err := c.UpdateValidIngredient(ctx, exampleValidIngredient)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientDataManager.On("ArchiveValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(expected)

		err := c.ArchiveValidIngredient(ctx, exampleValidIngredient.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
