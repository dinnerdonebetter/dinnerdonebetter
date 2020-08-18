package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeStepEventExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("RecipeStepEventExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(true, nil)

		actual, err := c.RecipeStepEventExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(exampleRecipeStepEvent, nil)

		actual, err := c.GetRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepEvent, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeStepEventsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetAllRecipeStepEventsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipeStepEventsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeStepEvents(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		results := make(chan []models.RecipeStepEvent)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetAllRecipeStepEvents", mock.Anything, results).Return(nil)

		err := c.GetAllRecipeStepEvents(ctx, results)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepEvents(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := models.DefaultQueryFilter()
		exampleRecipeStepEventList := fakemodels.BuildFakeRecipeStepEventList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetRecipeStepEvents", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepEventList, nil)

		actual, err := c.GetRecipeStepEvents(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepEventList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)
		exampleRecipeStepEventList := fakemodels.BuildFakeRecipeStepEventList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetRecipeStepEvents", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepEventList, nil)

		actual, err := c.GetRecipeStepEvents(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepEventList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepEventsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEventList := fakemodels.BuildFakeRecipeStepEventList().RecipeStepEvents
		var exampleIDs []uint64
		for _, x := range exampleRecipeStepEventList {
			exampleIDs = append(exampleIDs, x.ID)
		}

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("GetRecipeStepEventsWithIDs", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs).Return(exampleRecipeStepEventList, nil)

		actual, err := c.GetRecipeStepEventsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepEventList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("CreateRecipeStepEvent", mock.Anything, exampleInput).Return(exampleRecipeStepEvent, nil)

		actual, err := c.CreateRecipeStepEvent(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepEvent, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()

		c, mockDB := buildTestClient()

		mockDB.RecipeStepEventDataManager.On("UpdateRecipeStepEvent", mock.Anything, exampleRecipeStepEvent).Return(expected)

		err := c.UpdateRecipeStepEvent(ctx, exampleRecipeStepEvent)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepEventDataManager.On("ArchiveRecipeStepEvent", mock.Anything, exampleRecipeStepEvent.BelongsToRecipeStep, exampleRecipeStepEvent.ID).Return(expected)

		err := c.ArchiveRecipeStepEvent(ctx, exampleRecipeStepEvent.BelongsToRecipeStep, exampleRecipeStepEvent.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
