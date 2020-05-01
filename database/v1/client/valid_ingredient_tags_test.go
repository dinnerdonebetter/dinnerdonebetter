package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ValidIngredientTagExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientTagDataManager.On("ValidIngredientTagExists", mock.Anything, exampleValidIngredientTag.ID).Return(true, nil)

		actual, err := c.ValidIngredientTagExists(ctx, exampleValidIngredientTag.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidIngredientTag(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(exampleValidIngredientTag, nil)

		actual, err := c.GetValidIngredientTag(ctx, exampleValidIngredientTag.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientTag, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllValidIngredientTagsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientTagDataManager.On("GetAllValidIngredientTagsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllValidIngredientTagsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetValidIngredientTags(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := models.DefaultQueryFilter()
		exampleValidIngredientTagList := fakemodels.BuildFakeValidIngredientTagList()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientTagDataManager.On("GetValidIngredientTags", mock.Anything, filter).Return(exampleValidIngredientTagList, nil)

		actual, err := c.GetValidIngredientTags(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientTagList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		exampleValidIngredientTagList := fakemodels.BuildFakeValidIngredientTagList()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientTagDataManager.On("GetValidIngredientTags", mock.Anything, filter).Return(exampleValidIngredientTagList, nil)

		actual, err := c.GetValidIngredientTags(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientTagList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateValidIngredientTag(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientTagDataManager.On("CreateValidIngredientTag", mock.Anything, exampleInput).Return(exampleValidIngredientTag, nil)

		actual, err := c.CreateValidIngredientTag(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientTag, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateValidIngredientTag(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		c, mockDB := buildTestClient()

		mockDB.ValidIngredientTagDataManager.On("UpdateValidIngredientTag", mock.Anything, exampleValidIngredientTag).Return(expected)

		err := c.UpdateValidIngredientTag(ctx, exampleValidIngredientTag)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveValidIngredientTag(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		c, mockDB := buildTestClient()
		mockDB.ValidIngredientTagDataManager.On("ArchiveValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(expected)

		err := c.ArchiveValidIngredientTag(ctx, exampleValidIngredientTag.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
