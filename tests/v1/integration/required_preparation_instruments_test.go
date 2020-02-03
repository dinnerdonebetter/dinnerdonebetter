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

func checkRequiredPreparationInstrumentEquality(t *testing.T, expected, actual *models.RequiredPreparationInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.InstrumentID, actual.InstrumentID, "expected InstrumentID for ID %d to be %v, but it was %v ", expected.ID, expected.InstrumentID, actual.InstrumentID)
	assert.Equal(t, expected.PreparationID, actual.PreparationID, "expected PreparationID for ID %d to be %v, but it was %v ", expected.ID, expected.PreparationID, actual.PreparationID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for ID %d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyRequiredPreparationInstrument(t *testing.T) *models.RequiredPreparationInstrument {
	t.Helper()

	x := &models.RequiredPreparationInstrumentCreationInput{
		InstrumentID:  uint64(fake.Uint32()),
		PreparationID: uint64(fake.Uint32()),
		Notes:         fake.Word(),
	}
	y, err := todoClient.CreateRequiredPreparationInstrument(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestRequiredPreparationInstruments(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create required preparation instrument
			expected := &models.RequiredPreparationInstrument{
				InstrumentID:  uint64(fake.Uint32()),
				PreparationID: uint64(fake.Uint32()),
				Notes:         fake.Word(),
			}
			premade, err := todoClient.CreateRequiredPreparationInstrument(ctx, &models.RequiredPreparationInstrumentCreationInput{
				InstrumentID:  expected.InstrumentID,
				PreparationID: expected.PreparationID,
				Notes:         expected.Notes,
			})
			checkValueAndError(t, premade, err)

			// Assert required preparation instrument equality
			checkRequiredPreparationInstrumentEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveRequiredPreparationInstrument(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetRequiredPreparationInstrument(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkRequiredPreparationInstrumentEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create required preparation instruments
			var expected []*models.RequiredPreparationInstrument
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyRequiredPreparationInstrument(t))
			}

			// Assert required preparation instrument list equality
			actual, err := todoClient.GetRequiredPreparationInstruments(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RequiredPreparationInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RequiredPreparationInstruments),
			)

			// Clean up
			for _, x := range actual.RequiredPreparationInstruments {
				err = todoClient.ArchiveRequiredPreparationInstrument(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch required preparation instrument
			_, err := todoClient.GetRequiredPreparationInstrument(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create required preparation instrument
			expected := &models.RequiredPreparationInstrument{
				InstrumentID:  uint64(fake.Uint32()),
				PreparationID: uint64(fake.Uint32()),
				Notes:         fake.Word(),
			}
			premade, err := todoClient.CreateRequiredPreparationInstrument(ctx, &models.RequiredPreparationInstrumentCreationInput{
				InstrumentID:  expected.InstrumentID,
				PreparationID: expected.PreparationID,
				Notes:         expected.Notes,
			})
			checkValueAndError(t, premade, err)

			// Fetch required preparation instrument
			actual, err := todoClient.GetRequiredPreparationInstrument(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert required preparation instrument equality
			checkRequiredPreparationInstrumentEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveRequiredPreparationInstrument(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateRequiredPreparationInstrument(ctx, &models.RequiredPreparationInstrument{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create required preparation instrument
			expected := &models.RequiredPreparationInstrument{
				InstrumentID:  uint64(fake.Uint32()),
				PreparationID: uint64(fake.Uint32()),
				Notes:         fake.Word(),
			}
			premade, err := todoClient.CreateRequiredPreparationInstrument(tctx, &models.RequiredPreparationInstrumentCreationInput{
				InstrumentID:  uint64(fake.Uint32()),
				PreparationID: uint64(fake.Uint32()),
				Notes:         fake.Word(),
			})
			checkValueAndError(t, premade, err)

			// Change required preparation instrument
			premade.Update(expected.ToInput())
			err = todoClient.UpdateRequiredPreparationInstrument(ctx, premade)
			assert.NoError(t, err)

			// Fetch required preparation instrument
			actual, err := todoClient.GetRequiredPreparationInstrument(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert required preparation instrument equality
			checkRequiredPreparationInstrumentEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveRequiredPreparationInstrument(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create required preparation instrument
			expected := &models.RequiredPreparationInstrument{
				InstrumentID:  uint64(fake.Uint32()),
				PreparationID: uint64(fake.Uint32()),
				Notes:         fake.Word(),
			}
			premade, err := todoClient.CreateRequiredPreparationInstrument(ctx, &models.RequiredPreparationInstrumentCreationInput{
				InstrumentID:  expected.InstrumentID,
				PreparationID: expected.PreparationID,
				Notes:         expected.Notes,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveRequiredPreparationInstrument(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
