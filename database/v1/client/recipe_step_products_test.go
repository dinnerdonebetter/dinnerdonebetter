package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_RecipeStepProductExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("RecipeStepProductExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(true, nil)

		actual, err := c.RecipeStepProductExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(exampleRecipeStepProduct, nil)

		actual, err := c.GetRecipeStepProduct(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProduct, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeStepProductsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetAllRecipeStepProductsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllRecipeStepProductsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		results := make(chan []models.RecipeStepProduct)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetAllRecipeStepProducts", mock.Anything, results).Return(nil)

		err := c.GetAllRecipeStepProducts(ctx, results)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := models.DefaultQueryFilter()
		exampleRecipeStepProductList := fakemodels.BuildFakeRecipeStepProductList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetRecipeStepProducts", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepProductList, nil)

		actual, err := c.GetRecipeStepProducts(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProductList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		filter := (*models.QueryFilter)(nil)
		exampleRecipeStepProductList := fakemodels.BuildFakeRecipeStepProductList()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetRecipeStepProducts", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, filter).Return(exampleRecipeStepProductList, nil)

		actual, err := c.GetRecipeStepProducts(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProductList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetRecipeStepProductsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepProductList := fakemodels.BuildFakeRecipeStepProductList().RecipeStepProducts
		var exampleIDs []uint64
		for _, x := range exampleRecipeStepProductList {
			exampleIDs = append(exampleIDs, x.ID)
		}

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("GetRecipeStepProductsWithIDs", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs).Return(exampleRecipeStepProductList, nil)

		actual, err := c.GetRecipeStepProductsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProductList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("CreateRecipeStepProduct", mock.Anything, exampleInput).Return(exampleRecipeStepProduct, nil)

		actual, err := c.CreateRecipeStepProduct(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProduct, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()

		c, mockDB := buildTestClient()

		mockDB.RecipeStepProductDataManager.On("UpdateRecipeStepProduct", mock.Anything, exampleRecipeStepProduct).Return(expected)

		err := c.UpdateRecipeStepProduct(ctx, exampleRecipeStepProduct)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()

		c, mockDB := buildTestClient()
		mockDB.RecipeStepProductDataManager.On("ArchiveRecipeStepProduct", mock.Anything, exampleRecipeStepProduct.BelongsToRecipeStep, exampleRecipeStepProduct.ID).Return(expected)

		err := c.ArchiveRecipeStepProduct(ctx, exampleRecipeStepProduct.BelongsToRecipeStep, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
