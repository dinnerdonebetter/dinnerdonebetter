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

func checkRecipeStepIngredientEquality(t *testing.T, expected, actual *models.RecipeStepIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, *expected.IngredientID, *actual.IngredientID, "expected IngredientID to be %v, but it was %v ", expected.IngredientID, actual.IngredientID)
	assert.Equal(t, expected.QuantityType, actual.QuantityType, "expected QuantityType for ID %d to be %v, but it was %v ", expected.ID, expected.QuantityType, actual.QuantityType)
	assert.Equal(t, expected.QuantityValue, actual.QuantityValue, "expected QuantityValue for ID %d to be %v, but it was %v ", expected.ID, expected.QuantityValue, actual.QuantityValue)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for ID %d to be %v, but it was %v ", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.Equal(t, expected.ProductOfRecipe, actual.ProductOfRecipe, "expected ProductOfRecipe for ID %d to be %v, but it was %v ", expected.ID, expected.ProductOfRecipe, actual.ProductOfRecipe)
	assert.Equal(t, expected.IngredientNotes, actual.IngredientNotes, "expected IngredientNotes for ID %d to be %v, but it was %v ", expected.ID, expected.IngredientNotes, actual.IngredientNotes)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyRecipeStepIngredient(t *testing.T) *models.RecipeStepIngredient {
	t.Helper()

	x := &models.RecipeStepIngredientCreationInput{
		IngredientID:    func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		QuantityType:    fake.Word(),
		QuantityValue:   fake.Float32(),
		QuantityNotes:   fake.Word(),
		ProductOfRecipe: fake.Bool(),
		IngredientNotes: fake.Word(),
		RecipeStepID:    uint64(fake.Uint32()),
	}
	y, err := todoClient.CreateRecipeStepIngredient(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestRecipeStepIngredients(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step ingredient
			expected := &models.RecipeStepIngredient{
				IngredientID:    func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				QuantityType:    fake.Word(),
				QuantityValue:   fake.Float32(),
				QuantityNotes:   fake.Word(),
				ProductOfRecipe: fake.Bool(),
				IngredientNotes: fake.Word(),
				RecipeStepID:    uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepIngredient(ctx, &models.RecipeStepIngredientCreationInput{
				IngredientID:    expected.IngredientID,
				QuantityType:    expected.QuantityType,
				QuantityValue:   expected.QuantityValue,
				QuantityNotes:   expected.QuantityNotes,
				ProductOfRecipe: expected.ProductOfRecipe,
				IngredientNotes: expected.IngredientNotes,
				RecipeStepID:    expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveRecipeStepIngredient(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetRecipeStepIngredient(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepIngredientEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step ingredients
			var expected []*models.RecipeStepIngredient
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyRecipeStepIngredient(t))
			}

			// Assert recipe step ingredient list equality
			actual, err := todoClient.GetRecipeStepIngredients(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepIngredients),
			)

			// Clean up
			for _, x := range actual.RecipeStepIngredients {
				err = todoClient.ArchiveRecipeStepIngredient(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch recipe step ingredient
			_, err := todoClient.GetRecipeStepIngredient(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step ingredient
			expected := &models.RecipeStepIngredient{
				IngredientID:    func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				QuantityType:    fake.Word(),
				QuantityValue:   fake.Float32(),
				QuantityNotes:   fake.Word(),
				ProductOfRecipe: fake.Bool(),
				IngredientNotes: fake.Word(),
				RecipeStepID:    uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepIngredient(ctx, &models.RecipeStepIngredientCreationInput{
				IngredientID:    expected.IngredientID,
				QuantityType:    expected.QuantityType,
				QuantityValue:   expected.QuantityValue,
				QuantityNotes:   expected.QuantityNotes,
				ProductOfRecipe: expected.ProductOfRecipe,
				IngredientNotes: expected.IngredientNotes,
				RecipeStepID:    expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Fetch recipe step ingredient
			actual, err := todoClient.GetRecipeStepIngredient(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveRecipeStepIngredient(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateRecipeStepIngredient(ctx, &models.RecipeStepIngredient{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step ingredient
			expected := &models.RecipeStepIngredient{
				IngredientID:    func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				QuantityType:    fake.Word(),
				QuantityValue:   fake.Float32(),
				QuantityNotes:   fake.Word(),
				ProductOfRecipe: fake.Bool(),
				IngredientNotes: fake.Word(),
				RecipeStepID:    uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepIngredient(tctx, &models.RecipeStepIngredientCreationInput{
				IngredientID:    func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				QuantityType:    fake.Word(),
				QuantityValue:   fake.Float32(),
				QuantityNotes:   fake.Word(),
				ProductOfRecipe: fake.Bool(),
				IngredientNotes: fake.Word(),
				RecipeStepID:    uint64(fake.Uint32()),
			})
			checkValueAndError(t, premade, err)

			// Change recipe step ingredient
			premade.Update(expected.ToInput())
			err = todoClient.UpdateRecipeStepIngredient(ctx, premade)
			assert.NoError(t, err)

			// Fetch recipe step ingredient
			actual, err := todoClient.GetRecipeStepIngredient(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveRecipeStepIngredient(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step ingredient
			expected := &models.RecipeStepIngredient{
				IngredientID:    func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				QuantityType:    fake.Word(),
				QuantityValue:   fake.Float32(),
				QuantityNotes:   fake.Word(),
				ProductOfRecipe: fake.Bool(),
				IngredientNotes: fake.Word(),
				RecipeStepID:    uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepIngredient(ctx, &models.RecipeStepIngredientCreationInput{
				IngredientID:    expected.IngredientID,
				QuantityType:    expected.QuantityType,
				QuantityValue:   expected.QuantityValue,
				QuantityNotes:   expected.QuantityNotes,
				ProductOfRecipe: expected.ProductOfRecipe,
				IngredientNotes: expected.IngredientNotes,
				RecipeStepID:    expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveRecipeStepIngredient(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
