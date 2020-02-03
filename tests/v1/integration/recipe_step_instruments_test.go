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

func checkRecipeStepInstrumentEquality(t *testing.T, expected, actual *models.RecipeStepInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, *expected.InstrumentID, *actual.InstrumentID, "expected InstrumentID to be %v, but it was %v ", expected.InstrumentID, actual.InstrumentID)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for ID %d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyRecipeStepInstrument(t *testing.T) *models.RecipeStepInstrument {
	t.Helper()

	x := &models.RecipeStepInstrumentCreationInput{
		InstrumentID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		RecipeStepID: uint64(fake.Uint32()),
		Notes:        fake.Word(),
	}
	y, err := todoClient.CreateRecipeStepInstrument(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestRecipeStepInstruments(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step instrument
			expected := &models.RecipeStepInstrument{
				InstrumentID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				RecipeStepID: uint64(fake.Uint32()),
				Notes:        fake.Word(),
			}
			premade, err := todoClient.CreateRecipeStepInstrument(ctx, &models.RecipeStepInstrumentCreationInput{
				InstrumentID: expected.InstrumentID,
				RecipeStepID: expected.RecipeStepID,
				Notes:        expected.Notes,
			})
			checkValueAndError(t, premade, err)

			// Assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveRecipeStepInstrument(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetRecipeStepInstrument(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepInstrumentEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step instruments
			var expected []*models.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyRecipeStepInstrument(t))
			}

			// Assert recipe step instrument list equality
			actual, err := todoClient.GetRecipeStepInstruments(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepInstruments),
			)

			// Clean up
			for _, x := range actual.RecipeStepInstruments {
				err = todoClient.ArchiveRecipeStepInstrument(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch recipe step instrument
			_, err := todoClient.GetRecipeStepInstrument(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step instrument
			expected := &models.RecipeStepInstrument{
				InstrumentID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				RecipeStepID: uint64(fake.Uint32()),
				Notes:        fake.Word(),
			}
			premade, err := todoClient.CreateRecipeStepInstrument(ctx, &models.RecipeStepInstrumentCreationInput{
				InstrumentID: expected.InstrumentID,
				RecipeStepID: expected.RecipeStepID,
				Notes:        expected.Notes,
			})
			checkValueAndError(t, premade, err)

			// Fetch recipe step instrument
			actual, err := todoClient.GetRecipeStepInstrument(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveRecipeStepInstrument(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateRecipeStepInstrument(ctx, &models.RecipeStepInstrument{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step instrument
			expected := &models.RecipeStepInstrument{
				InstrumentID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				RecipeStepID: uint64(fake.Uint32()),
				Notes:        fake.Word(),
			}
			premade, err := todoClient.CreateRecipeStepInstrument(tctx, &models.RecipeStepInstrumentCreationInput{
				InstrumentID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				RecipeStepID: uint64(fake.Uint32()),
				Notes:        fake.Word(),
			})
			checkValueAndError(t, premade, err)

			// Change recipe step instrument
			premade.Update(expected.ToInput())
			err = todoClient.UpdateRecipeStepInstrument(ctx, premade)
			assert.NoError(t, err)

			// Fetch recipe step instrument
			actual, err := todoClient.GetRecipeStepInstrument(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveRecipeStepInstrument(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step instrument
			expected := &models.RecipeStepInstrument{
				InstrumentID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
				RecipeStepID: uint64(fake.Uint32()),
				Notes:        fake.Word(),
			}
			premade, err := todoClient.CreateRecipeStepInstrument(ctx, &models.RecipeStepInstrumentCreationInput{
				InstrumentID: expected.InstrumentID,
				RecipeStepID: expected.RecipeStepID,
				Notes:        expected.Notes,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveRecipeStepInstrument(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
