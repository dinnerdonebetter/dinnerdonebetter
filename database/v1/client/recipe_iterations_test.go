package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeIterationExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, nil)

		actual, err := c.RecipeIterationExists(ctx, exampleRecipe.ID, exampleRecipeIteration.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(exampleRecipeIteration, nil)

		actual, err := c.GetRecipeIteration(ctx, exampleRecipe.ID, exampleRecipeIteration.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIteration, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeIterationsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetAllRecipeIterationsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipeIterationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeIterations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		results := make(chan []models.RecipeIteration)

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetAllRecipeIterations", mock.Anything, results).Return(nil)

		err := c.GetAllRecipeIterations(ctx, results)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeIterations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := models.DefaultQueryFilter()
		exampleRecipeIterationList := fakemodels.BuildFakeRecipeIterationList()

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetRecipeIterations", mock.Anything, exampleRecipe.ID, filter).Return(exampleRecipeIterationList, nil)

		actual, err := c.GetRecipeIterations(ctx, exampleRecipe.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)
		exampleRecipeIterationList := fakemodels.BuildFakeRecipeIterationList()

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetRecipeIterations", mock.Anything, exampleRecipe.ID, filter).Return(exampleRecipeIterationList, nil)

		actual, err := c.GetRecipeIterations(ctx, exampleRecipe.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeIterationsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIterationList := fakemodels.BuildFakeRecipeIterationList().RecipeIterations
		var exampleIDs []uint64
		for _, x := range exampleRecipeIterationList {
			exampleIDs = append(exampleIDs, x.ID)
		}

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("GetRecipeIterationsWithIDs", mock.Anything, exampleRecipe.ID, defaultLimit, exampleIDs).Return(exampleRecipeIterationList, nil)

		actual, err := c.GetRecipeIterationsWithIDs(ctx, exampleRecipe.ID, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("CreateRecipeIteration", mock.Anything, exampleInput).Return(exampleRecipeIteration, nil)

		actual, err := c.CreateRecipeIteration(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIteration, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()

		c, mockDB := buildTestClient()

		mockDB.RecipeIterationDataManager.On("UpdateRecipeIteration", mock.Anything, exampleRecipeIteration).Return(expected)

		err := c.UpdateRecipeIteration(ctx, exampleRecipeIteration)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()

		c, mockDB := buildTestClient()
		mockDB.RecipeIterationDataManager.On("ArchiveRecipeIteration", mock.Anything, exampleRecipeIteration.BelongsToRecipe, exampleRecipeIteration.ID).Return(expected)

		err := c.ArchiveRecipeIteration(ctx, exampleRecipeIteration.BelongsToRecipe, exampleRecipeIteration.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
