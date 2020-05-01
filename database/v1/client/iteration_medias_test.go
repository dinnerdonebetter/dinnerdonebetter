package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_IterationMediaExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("IterationMediaExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(true, nil)

		actual, err := c.IterationMediaExists(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(exampleIterationMedia, nil)

		actual, err := c.GetIterationMedia(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleIterationMedia, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllIterationMediasCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("GetAllIterationMediasCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllIterationMediasCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetIterationMedias(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		filter := models.DefaultQueryFilter()
		exampleIterationMediaList := fakemodels.BuildFakeIterationMediaList()

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("GetIterationMedias", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, filter).Return(exampleIterationMediaList, nil)

		actual, err := c.GetIterationMedias(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleIterationMediaList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		filter := (*models.QueryFilter)(nil)
		exampleIterationMediaList := fakemodels.BuildFakeIterationMediaList()

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("GetIterationMedias", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, filter).Return(exampleIterationMediaList, nil)

		actual, err := c.GetIterationMedias(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleIterationMediaList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("CreateIterationMedia", mock.Anything, exampleInput).Return(exampleIterationMedia, nil)

		actual, err := c.CreateIterationMedia(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleIterationMedia, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()

		c, mockDB := buildTestClient()

		mockDB.IterationMediaDataManager.On("UpdateIterationMedia", mock.Anything, exampleIterationMedia).Return(expected)

		err := c.UpdateIterationMedia(ctx, exampleIterationMedia)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("ArchiveIterationMedia", mock.Anything, exampleIterationMedia.BelongsToRecipeIteration, exampleIterationMedia.ID).Return(expected)

		err := c.ArchiveIterationMedia(ctx, exampleIterationMedia.BelongsToRecipeIteration, exampleIterationMedia.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
