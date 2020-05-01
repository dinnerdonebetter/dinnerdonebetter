package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeStepExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)

		actual, err := c.RecipeStepExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("GetRecipeStep", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(exampleRecipeStep, nil)

		actual, err := c.GetRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStep, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeStepsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("GetAllRecipeStepsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipeStepsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := models.DefaultQueryFilter()
		exampleRecipeStepList := fakemodels.BuildFakeRecipeStepList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("GetRecipeSteps", mock.Anything, exampleRecipe.ID, filter).Return(exampleRecipeStepList, nil)

		actual, err := c.GetRecipeSteps(ctx, exampleRecipe.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := (*models.QueryFilter)(nil)
		exampleRecipeStepList := fakemodels.BuildFakeRecipeStepList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("GetRecipeSteps", mock.Anything, exampleRecipe.ID, filter).Return(exampleRecipeStepList, nil)

		actual, err := c.GetRecipeSteps(ctx, exampleRecipe.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("CreateRecipeStep", mock.Anything, exampleInput).Return(exampleRecipeStep, nil)

		actual, err := c.CreateRecipeStep(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStep, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()

		c, mockDB := buildTestClient()

		mockDB.RecipeStepDataManager.On("UpdateRecipeStep", mock.Anything, exampleRecipeStep).Return(expected)

		err := c.UpdateRecipeStep(ctx, exampleRecipeStep)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepDataManager.On("ArchiveRecipeStep", mock.Anything, exampleRecipeStep.BelongsToRecipe, exampleRecipeStep.ID).Return(expected)

		err := c.ArchiveRecipeStep(ctx, exampleRecipeStep.BelongsToRecipe, exampleRecipeStep.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
