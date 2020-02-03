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

func checkRecipeIterationEquality(t *testing.T, expected, actual *models.RecipeIteration) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.RecipeID, actual.RecipeID, "expected RecipeID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeID, actual.RecipeID)
	assert.Equal(t, expected.EndDifficultyRating, actual.EndDifficultyRating, "expected EndDifficultyRating for ID %d to be %v, but it was %v ", expected.ID, expected.EndDifficultyRating, actual.EndDifficultyRating)
	assert.Equal(t, expected.EndComplexityRating, actual.EndComplexityRating, "expected EndComplexityRating for ID %d to be %v, but it was %v ", expected.ID, expected.EndComplexityRating, actual.EndComplexityRating)
	assert.Equal(t, expected.EndTasteRating, actual.EndTasteRating, "expected EndTasteRating for ID %d to be %v, but it was %v ", expected.ID, expected.EndTasteRating, actual.EndTasteRating)
	assert.Equal(t, expected.EndOverallRating, actual.EndOverallRating, "expected EndOverallRating for ID %d to be %v, but it was %v ", expected.ID, expected.EndOverallRating, actual.EndOverallRating)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyRecipeIteration(t *testing.T) *models.RecipeIteration {
	t.Helper()

	x := &models.RecipeIterationCreationInput{
		RecipeID:            uint64(fake.Uint32()),
		EndDifficultyRating: fake.Float32(),
		EndComplexityRating: fake.Float32(),
		EndTasteRating:      fake.Float32(),
		EndOverallRating:    fake.Float32(),
	}
	y, err := todoClient.CreateRecipeIteration(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestRecipeIterations(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe iteration
			expected := &models.RecipeIteration{
				RecipeID:            uint64(fake.Uint32()),
				EndDifficultyRating: fake.Float32(),
				EndComplexityRating: fake.Float32(),
				EndTasteRating:      fake.Float32(),
				EndOverallRating:    fake.Float32(),
			}
			premade, err := todoClient.CreateRecipeIteration(ctx, &models.RecipeIterationCreationInput{
				RecipeID:            expected.RecipeID,
				EndDifficultyRating: expected.EndDifficultyRating,
				EndComplexityRating: expected.EndComplexityRating,
				EndTasteRating:      expected.EndTasteRating,
				EndOverallRating:    expected.EndOverallRating,
			})
			checkValueAndError(t, premade, err)

			// Assert recipe iteration equality
			checkRecipeIterationEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveRecipeIteration(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetRecipeIteration(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkRecipeIterationEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe iterations
			var expected []*models.RecipeIteration
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyRecipeIteration(t))
			}

			// Assert recipe iteration list equality
			actual, err := todoClient.GetRecipeIterations(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeIterations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeIterations),
			)

			// Clean up
			for _, x := range actual.RecipeIterations {
				err = todoClient.ArchiveRecipeIteration(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch recipe iteration
			_, err := todoClient.GetRecipeIteration(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe iteration
			expected := &models.RecipeIteration{
				RecipeID:            uint64(fake.Uint32()),
				EndDifficultyRating: fake.Float32(),
				EndComplexityRating: fake.Float32(),
				EndTasteRating:      fake.Float32(),
				EndOverallRating:    fake.Float32(),
			}
			premade, err := todoClient.CreateRecipeIteration(ctx, &models.RecipeIterationCreationInput{
				RecipeID:            expected.RecipeID,
				EndDifficultyRating: expected.EndDifficultyRating,
				EndComplexityRating: expected.EndComplexityRating,
				EndTasteRating:      expected.EndTasteRating,
				EndOverallRating:    expected.EndOverallRating,
			})
			checkValueAndError(t, premade, err)

			// Fetch recipe iteration
			actual, err := todoClient.GetRecipeIteration(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe iteration equality
			checkRecipeIterationEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveRecipeIteration(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateRecipeIteration(ctx, &models.RecipeIteration{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe iteration
			expected := &models.RecipeIteration{
				RecipeID:            uint64(fake.Uint32()),
				EndDifficultyRating: fake.Float32(),
				EndComplexityRating: fake.Float32(),
				EndTasteRating:      fake.Float32(),
				EndOverallRating:    fake.Float32(),
			}
			premade, err := todoClient.CreateRecipeIteration(tctx, &models.RecipeIterationCreationInput{
				RecipeID:            uint64(fake.Uint32()),
				EndDifficultyRating: fake.Float32(),
				EndComplexityRating: fake.Float32(),
				EndTasteRating:      fake.Float32(),
				EndOverallRating:    fake.Float32(),
			})
			checkValueAndError(t, premade, err)

			// Change recipe iteration
			premade.Update(expected.ToInput())
			err = todoClient.UpdateRecipeIteration(ctx, premade)
			assert.NoError(t, err)

			// Fetch recipe iteration
			actual, err := todoClient.GetRecipeIteration(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe iteration equality
			checkRecipeIterationEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveRecipeIteration(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe iteration
			expected := &models.RecipeIteration{
				RecipeID:            uint64(fake.Uint32()),
				EndDifficultyRating: fake.Float32(),
				EndComplexityRating: fake.Float32(),
				EndTasteRating:      fake.Float32(),
				EndOverallRating:    fake.Float32(),
			}
			premade, err := todoClient.CreateRecipeIteration(ctx, &models.RecipeIterationCreationInput{
				RecipeID:            expected.RecipeID,
				EndDifficultyRating: expected.EndDifficultyRating,
				EndComplexityRating: expected.EndComplexityRating,
				EndTasteRating:      expected.EndTasteRating,
				EndOverallRating:    expected.EndOverallRating,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveRecipeIteration(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
