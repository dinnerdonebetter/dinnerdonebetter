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

func checkInstrumentEquality(t *testing.T, expected, actual *models.Instrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Variant, actual.Variant, "expected Variant for ID %d to be %v, but it was %v ", expected.ID, expected.Variant, actual.Variant)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for ID %d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Icon, actual.Icon, "expected Icon for ID %d to be %v, but it was %v ", expected.ID, expected.Icon, actual.Icon)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyInstrument(t *testing.T) *models.Instrument {
	t.Helper()

	x := &models.InstrumentCreationInput{
		Name:        fake.Word(),
		Variant:     fake.Word(),
		Description: fake.Word(),
		Icon:        fake.Word(),
	}
	y, err := todoClient.CreateInstrument(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestInstruments(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create instrument
			expected := &models.Instrument{
				Name:        fake.Word(),
				Variant:     fake.Word(),
				Description: fake.Word(),
				Icon:        fake.Word(),
			}
			premade, err := todoClient.CreateInstrument(ctx, &models.InstrumentCreationInput{
				Name:        expected.Name,
				Variant:     expected.Variant,
				Description: expected.Description,
				Icon:        expected.Icon,
			})
			checkValueAndError(t, premade, err)

			// Assert instrument equality
			checkInstrumentEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveInstrument(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetInstrument(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkInstrumentEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create instruments
			var expected []*models.Instrument
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyInstrument(t))
			}

			// Assert instrument list equality
			actual, err := todoClient.GetInstruments(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Instruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Instruments),
			)

			// Clean up
			for _, x := range actual.Instruments {
				err = todoClient.ArchiveInstrument(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch instrument
			_, err := todoClient.GetInstrument(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create instrument
			expected := &models.Instrument{
				Name:        fake.Word(),
				Variant:     fake.Word(),
				Description: fake.Word(),
				Icon:        fake.Word(),
			}
			premade, err := todoClient.CreateInstrument(ctx, &models.InstrumentCreationInput{
				Name:        expected.Name,
				Variant:     expected.Variant,
				Description: expected.Description,
				Icon:        expected.Icon,
			})
			checkValueAndError(t, premade, err)

			// Fetch instrument
			actual, err := todoClient.GetInstrument(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert instrument equality
			checkInstrumentEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveInstrument(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateInstrument(ctx, &models.Instrument{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create instrument
			expected := &models.Instrument{
				Name:        fake.Word(),
				Variant:     fake.Word(),
				Description: fake.Word(),
				Icon:        fake.Word(),
			}
			premade, err := todoClient.CreateInstrument(tctx, &models.InstrumentCreationInput{
				Name:        fake.Word(),
				Variant:     fake.Word(),
				Description: fake.Word(),
				Icon:        fake.Word(),
			})
			checkValueAndError(t, premade, err)

			// Change instrument
			premade.Update(expected.ToInput())
			err = todoClient.UpdateInstrument(ctx, premade)
			assert.NoError(t, err)

			// Fetch instrument
			actual, err := todoClient.GetInstrument(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert instrument equality
			checkInstrumentEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveInstrument(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create instrument
			expected := &models.Instrument{
				Name:        fake.Word(),
				Variant:     fake.Word(),
				Description: fake.Word(),
				Icon:        fake.Word(),
			}
			premade, err := todoClient.CreateInstrument(ctx, &models.InstrumentCreationInput{
				Name:        expected.Name,
				Variant:     expected.Variant,
				Description: expected.Description,
				Icon:        expected.Icon,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveInstrument(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
