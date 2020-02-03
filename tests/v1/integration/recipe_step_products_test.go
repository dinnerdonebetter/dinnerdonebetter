package integration

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
)

func checkRecipeStepProductEquality(t *testing.T, expected, actual *models.RecipeStepProduct) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyRecipeStepProduct(t *testing.T) *models.RecipeStepProduct {
	t.Helper()

	x := &models.RecipeStepProductCreationInput{
		Name:         fake.Word(),
		RecipeStepID: uint64(fake.Uint32()),
	}
	y, err := todoClient.CreateRecipeStepProduct(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestRecipeStepProducts(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step product
			expected := &models.RecipeStepProduct{
				Name:         fake.Word(),
				RecipeStepID: uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepProduct(ctx, &models.RecipeStepProductCreationInput{
				Name:         expected.Name,
				RecipeStepID: expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Assert recipe step product equality
			checkRecipeStepProductEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveRecipeStepProduct(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetRecipeStepProduct(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepProductEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step products
			var expected []*models.RecipeStepProduct
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyRecipeStepProduct(t))
			}

			// Assert recipe step product list equality
			actual, err := todoClient.GetRecipeStepProducts(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepProducts),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepProducts),
			)

			// Clean up
			for _, x := range actual.RecipeStepProducts {
				err = todoClient.ArchiveRecipeStepProduct(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch recipe step product
			_, err := todoClient.GetRecipeStepProduct(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step product
			expected := &models.RecipeStepProduct{
				Name:         fake.Word(),
				RecipeStepID: uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepProduct(ctx, &models.RecipeStepProductCreationInput{
				Name:         expected.Name,
				RecipeStepID: expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Fetch recipe step product
			actual, err := todoClient.GetRecipeStepProduct(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step product equality
			checkRecipeStepProductEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveRecipeStepProduct(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateRecipeStepProduct(ctx, &models.RecipeStepProduct{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step product
			expected := &models.RecipeStepProduct{
				Name:         fake.Word(),
				RecipeStepID: uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepProduct(tctx, &models.RecipeStepProductCreationInput{
				Name:         fake.Word(),
				RecipeStepID: uint64(fake.Uint32()),
			})
			checkValueAndError(t, premade, err)

			// Change recipe step product
			premade.Update(expected.ToInput())
			err = todoClient.UpdateRecipeStepProduct(ctx, premade)
			assert.NoError(t, err)

			// Fetch recipe step product
			actual, err := todoClient.GetRecipeStepProduct(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step product equality
			checkRecipeStepProductEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveRecipeStepProduct(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step product
			expected := &models.RecipeStepProduct{
				Name:         fake.Word(),
				RecipeStepID: uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepProduct(ctx, &models.RecipeStepProductCreationInput{
				Name:         expected.Name,
				RecipeStepID: expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveRecipeStepProduct(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
