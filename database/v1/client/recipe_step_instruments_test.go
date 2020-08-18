package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeStepInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("RecipeStepInstrumentExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(true, nil)

		actual, err := c.RecipeStepInstrumentExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(exampleRecipeStepInstrument, nil)

		actual, err := c.GetRecipeStepInstrument(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrument, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeStepInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetAllRecipeStepInstrumentsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipeStepInstrumentsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		results := make(chan []models.RecipeStepInstrument)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetAllRecipeStepInstruments", mock.Anything, results).Return(nil)

		err := c.GetAllRecipeStepInstruments(ctx, results)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := models.DefaultQueryFilter()
		exampleRecipeStepInstrumentList := fakemodels.BuildFakeRecipeStepInstrumentList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetRecipeStepInstruments", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepInstrumentList, nil)

		actual, err := c.GetRecipeStepInstruments(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)
		exampleRecipeStepInstrumentList := fakemodels.BuildFakeRecipeStepInstrumentList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetRecipeStepInstruments", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepInstrumentList, nil)

		actual, err := c.GetRecipeStepInstruments(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepInstrumentsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrumentList := fakemodels.BuildFakeRecipeStepInstrumentList().RecipeStepInstruments
		var exampleIDs []uint64
		for _, x := range exampleRecipeStepInstrumentList {
			exampleIDs = append(exampleIDs, x.ID)
		}

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("GetRecipeStepInstrumentsWithIDs", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs).Return(exampleRecipeStepInstrumentList, nil)

		actual, err := c.GetRecipeStepInstrumentsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("CreateRecipeStepInstrument", mock.Anything, exampleInput).Return(exampleRecipeStepInstrument, nil)

		actual, err := c.CreateRecipeStepInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrument, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()

		c, mockDB := buildTestClient()

		mockDB.RecipeStepInstrumentDataManager.On("UpdateRecipeStepInstrument", mock.Anything, exampleRecipeStepInstrument).Return(expected)

		err := c.UpdateRecipeStepInstrument(ctx, exampleRecipeStepInstrument)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepInstrumentDataManager.On("ArchiveRecipeStepInstrument", mock.Anything, exampleRecipeStepInstrument.BelongsToRecipeStep, exampleRecipeStepInstrument.ID).Return(expected)

		err := c.ArchiveRecipeStepInstrument(ctx, exampleRecipeStepInstrument.BelongsToRecipeStep, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
