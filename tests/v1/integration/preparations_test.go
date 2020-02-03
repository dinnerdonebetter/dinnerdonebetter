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

func checkPreparationEquality(t *testing.T, expected, actual *models.Preparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Variant, actual.Variant, "expected Variant for ID %d to be %v, but it was %v ", expected.ID, expected.Variant, actual.Variant)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for ID %d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.AllergyWarning, actual.AllergyWarning, "expected AllergyWarning for ID %d to be %v, but it was %v ", expected.ID, expected.AllergyWarning, actual.AllergyWarning)
	assert.Equal(t, expected.Icon, actual.Icon, "expected Icon for ID %d to be %v, but it was %v ", expected.ID, expected.Icon, actual.Icon)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyPreparation(t *testing.T) *models.Preparation {
	t.Helper()

	x := &models.PreparationCreationInput{
		Name:           fake.Word(),
		Variant:        fake.Word(),
		Description:    fake.Word(),
		AllergyWarning: fake.Word(),
		Icon:           fake.Word(),
	}
	y, err := todoClient.CreatePreparation(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestPreparations(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create preparation
			expected := &models.Preparation{
				Name:           fake.Word(),
				Variant:        fake.Word(),
				Description:    fake.Word(),
				AllergyWarning: fake.Word(),
				Icon:           fake.Word(),
			}
			premade, err := todoClient.CreatePreparation(ctx, &models.PreparationCreationInput{
				Name:           expected.Name,
				Variant:        expected.Variant,
				Description:    expected.Description,
				AllergyWarning: expected.AllergyWarning,
				Icon:           expected.Icon,
			})
			checkValueAndError(t, premade, err)

			// Assert preparation equality
			checkPreparationEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchivePreparation(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetPreparation(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkPreparationEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create preparations
			var expected []*models.Preparation
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyPreparation(t))
			}

			// Assert preparation list equality
			actual, err := todoClient.GetPreparations(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Preparations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Preparations),
			)

			// Clean up
			for _, x := range actual.Preparations {
				err = todoClient.ArchivePreparation(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch preparation
			_, err := todoClient.GetPreparation(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create preparation
			expected := &models.Preparation{
				Name:           fake.Word(),
				Variant:        fake.Word(),
				Description:    fake.Word(),
				AllergyWarning: fake.Word(),
				Icon:           fake.Word(),
			}
			premade, err := todoClient.CreatePreparation(ctx, &models.PreparationCreationInput{
				Name:           expected.Name,
				Variant:        expected.Variant,
				Description:    expected.Description,
				AllergyWarning: expected.AllergyWarning,
				Icon:           expected.Icon,
			})
			checkValueAndError(t, premade, err)

			// Fetch preparation
			actual, err := todoClient.GetPreparation(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert preparation equality
			checkPreparationEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchivePreparation(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdatePreparation(ctx, &models.Preparation{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create preparation
			expected := &models.Preparation{
				Name:           fake.Word(),
				Variant:        fake.Word(),
				Description:    fake.Word(),
				AllergyWarning: fake.Word(),
				Icon:           fake.Word(),
			}
			premade, err := todoClient.CreatePreparation(tctx, &models.PreparationCreationInput{
				Name:           fake.Word(),
				Variant:        fake.Word(),
				Description:    fake.Word(),
				AllergyWarning: fake.Word(),
				Icon:           fake.Word(),
			})
			checkValueAndError(t, premade, err)

			// Change preparation
			premade.Update(expected.ToInput())
			err = todoClient.UpdatePreparation(ctx, premade)
			assert.NoError(t, err)

			// Fetch preparation
			actual, err := todoClient.GetPreparation(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert preparation equality
			checkPreparationEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchivePreparation(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create preparation
			expected := &models.Preparation{
				Name:           fake.Word(),
				Variant:        fake.Word(),
				Description:    fake.Word(),
				AllergyWarning: fake.Word(),
				Icon:           fake.Word(),
			}
			premade, err := todoClient.CreatePreparation(ctx, &models.PreparationCreationInput{
				Name:           expected.Name,
				Variant:        expected.Variant,
				Description:    expected.Description,
				AllergyWarning: expected.AllergyWarning,
				Icon:           expected.Icon,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchivePreparation(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
