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

func checkRecipeStepEventEquality(t *testing.T, expected, actual *models.RecipeStepEvent) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.EventType, actual.EventType, "expected EventType for ID %d to be %v, but it was %v ", expected.ID, expected.EventType, actual.EventType)
	assert.Equal(t, expected.Done, actual.Done, "expected Done for ID %d to be %v, but it was %v ", expected.ID, expected.Done, actual.Done)
	assert.Equal(t, expected.RecipeIterationID, actual.RecipeIterationID, "expected RecipeIterationID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeIterationID, actual.RecipeIterationID)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyRecipeStepEvent(t *testing.T) *models.RecipeStepEvent {
	t.Helper()

	x := &models.RecipeStepEventCreationInput{
		EventType:         fake.Word(),
		Done:              fake.Bool(),
		RecipeIterationID: uint64(fake.Uint32()),
		RecipeStepID:      uint64(fake.Uint32()),
	}
	y, err := todoClient.CreateRecipeStepEvent(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestRecipeStepEvents(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step event
			expected := &models.RecipeStepEvent{
				EventType:         fake.Word(),
				Done:              fake.Bool(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepEvent(ctx, &models.RecipeStepEventCreationInput{
				EventType:         expected.EventType,
				Done:              expected.Done,
				RecipeIterationID: expected.RecipeIterationID,
				RecipeStepID:      expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Assert recipe step event equality
			checkRecipeStepEventEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveRecipeStepEvent(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetRecipeStepEvent(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepEventEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step events
			var expected []*models.RecipeStepEvent
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyRecipeStepEvent(t))
			}

			// Assert recipe step event list equality
			actual, err := todoClient.GetRecipeStepEvents(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepEvents),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepEvents),
			)

			// Clean up
			for _, x := range actual.RecipeStepEvents {
				err = todoClient.ArchiveRecipeStepEvent(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch recipe step event
			_, err := todoClient.GetRecipeStepEvent(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step event
			expected := &models.RecipeStepEvent{
				EventType:         fake.Word(),
				Done:              fake.Bool(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepEvent(ctx, &models.RecipeStepEventCreationInput{
				EventType:         expected.EventType,
				Done:              expected.Done,
				RecipeIterationID: expected.RecipeIterationID,
				RecipeStepID:      expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Fetch recipe step event
			actual, err := todoClient.GetRecipeStepEvent(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step event equality
			checkRecipeStepEventEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveRecipeStepEvent(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateRecipeStepEvent(ctx, &models.RecipeStepEvent{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step event
			expected := &models.RecipeStepEvent{
				EventType:         fake.Word(),
				Done:              fake.Bool(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepEvent(tctx, &models.RecipeStepEventCreationInput{
				EventType:         fake.Word(),
				Done:              fake.Bool(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      uint64(fake.Uint32()),
			})
			checkValueAndError(t, premade, err)

			// Change recipe step event
			premade.Update(expected.ToInput())
			err = todoClient.UpdateRecipeStepEvent(ctx, premade)
			assert.NoError(t, err)

			// Fetch recipe step event
			actual, err := todoClient.GetRecipeStepEvent(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step event equality
			checkRecipeStepEventEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveRecipeStepEvent(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step event
			expected := &models.RecipeStepEvent{
				EventType:         fake.Word(),
				Done:              fake.Bool(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStepEvent(ctx, &models.RecipeStepEventCreationInput{
				EventType:         expected.EventType,
				Done:              expected.Done,
				RecipeIterationID: expected.RecipeIterationID,
				RecipeStepID:      expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveRecipeStepEvent(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
