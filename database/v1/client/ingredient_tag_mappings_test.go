package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_IngredientTagMappingExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		c, mockDB := buildTestClient()
		mockDB.IngredientTagMappingDataManager.On("IngredientTagMappingExists", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(true, nil)

		actual, err := c.IngredientTagMappingExists(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetIngredientTagMapping(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		c, mockDB := buildTestClient()
		mockDB.IngredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(exampleIngredientTagMapping, nil)

		actual, err := c.GetIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleIngredientTagMapping, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllIngredientTagMappingsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.IngredientTagMappingDataManager.On("GetAllIngredientTagMappingsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllIngredientTagMappingsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetIngredientTagMappings(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := models.DefaultQueryFilter()
		exampleIngredientTagMappingList := fakemodels.BuildFakeIngredientTagMappingList()

		c, mockDB := buildTestClient()
		mockDB.IngredientTagMappingDataManager.On("GetIngredientTagMappings", mock.Anything, exampleValidIngredient.ID, filter).Return(exampleIngredientTagMappingList, nil)

		actual, err := c.GetIngredientTagMappings(ctx, exampleValidIngredient.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleIngredientTagMappingList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := (*models.QueryFilter)(nil)
		exampleIngredientTagMappingList := fakemodels.BuildFakeIngredientTagMappingList()

		c, mockDB := buildTestClient()
		mockDB.IngredientTagMappingDataManager.On("GetIngredientTagMappings", mock.Anything, exampleValidIngredient.ID, filter).Return(exampleIngredientTagMappingList, nil)

		actual, err := c.GetIngredientTagMappings(ctx, exampleValidIngredient.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleIngredientTagMappingList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateIngredientTagMapping(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		c, mockDB := buildTestClient()
		mockDB.IngredientTagMappingDataManager.On("CreateIngredientTagMapping", mock.Anything, exampleInput).Return(exampleIngredientTagMapping, nil)

		actual, err := c.CreateIngredientTagMapping(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleIngredientTagMapping, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateIngredientTagMapping(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()

		c, mockDB := buildTestClient()

		mockDB.IngredientTagMappingDataManager.On("UpdateIngredientTagMapping", mock.Anything, exampleIngredientTagMapping).Return(expected)

		err := c.UpdateIngredientTagMapping(ctx, exampleIngredientTagMapping)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveIngredientTagMapping(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()

		c, mockDB := buildTestClient()
		mockDB.IngredientTagMappingDataManager.On("ArchiveIngredientTagMapping", mock.Anything, exampleIngredientTagMapping.BelongsToValidIngredient, exampleIngredientTagMapping.ID).Return(expected)

		err := c.ArchiveIngredientTagMapping(ctx, exampleIngredientTagMapping.BelongsToValidIngredient, exampleIngredientTagMapping.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
