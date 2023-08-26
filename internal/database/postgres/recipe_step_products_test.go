package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuerier_RecipeStepProductExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepProductExists(ctx, "", exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepProductExists(ctx, exampleRecipeID, "", exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepProductExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProduct(ctx, "", exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProduct(ctx, exampleRecipeID, "", exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProduct(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_getRecipeStepProductsForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with missing recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.getRecipeStepProductsForRecipe(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProducts(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProducts(ctx, exampleRecipeID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepProduct(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepProduct(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepProduct(ctx, "", exampleRecipeStepProduct.ID))
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepProduct(ctx, exampleRecipeStepID, ""))
	})
}
