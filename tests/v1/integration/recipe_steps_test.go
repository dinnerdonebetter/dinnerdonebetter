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

func checkRecipeStepEquality(t *testing.T, expected, actual *models.RecipeStep) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Index, actual.Index, "expected Index for ID %d to be %v, but it was %v ", expected.ID, expected.Index, actual.Index)
	assert.Equal(t, expected.PreparationID, actual.PreparationID, "expected PreparationID for ID %d to be %v, but it was %v ", expected.ID, expected.PreparationID, actual.PreparationID)
	assert.Equal(t, expected.PrerequisiteStep, actual.PrerequisiteStep, "expected PrerequisiteStep for ID %d to be %v, but it was %v ", expected.ID, expected.PrerequisiteStep, actual.PrerequisiteStep)
	assert.Equal(t, expected.MinEstimatedTimeInSeconds, actual.MinEstimatedTimeInSeconds, "expected MinEstimatedTimeInSeconds for ID %d to be %v, but it was %v ", expected.ID, expected.MinEstimatedTimeInSeconds, actual.MinEstimatedTimeInSeconds)
	assert.Equal(t, expected.MaxEstimatedTimeInSeconds, actual.MaxEstimatedTimeInSeconds, "expected MaxEstimatedTimeInSeconds for ID %d to be %v, but it was %v ", expected.ID, expected.MaxEstimatedTimeInSeconds, actual.MaxEstimatedTimeInSeconds)
	assert.Equal(t, *expected.TemperatureInCelsius, *actual.TemperatureInCelsius, "expected TemperatureInCelsius to be %v, but it was %v ", expected.TemperatureInCelsius, actual.TemperatureInCelsius)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for ID %d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.RecipeID, actual.RecipeID, "expected RecipeID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeID, actual.RecipeID)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyRecipeStep(t *testing.T) *models.RecipeStep {
	t.Helper()

	x := &models.RecipeStepCreationInput{
		Index:                     uint(fake.Uint32()),
		PreparationID:             uint64(fake.Uint32()),
		PrerequisiteStep:          uint64(fake.Uint32()),
		MinEstimatedTimeInSeconds: fake.Uint32(),
		MaxEstimatedTimeInSeconds: fake.Uint32(),
		TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
		Notes:                     fake.Word(),
		RecipeID:                  uint64(fake.Uint32()),
	}
	y, err := todoClient.CreateRecipeStep(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestRecipeSteps(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step
			expected := &models.RecipeStep{
				Index:                     uint(fake.Uint32()),
				PreparationID:             uint64(fake.Uint32()),
				PrerequisiteStep:          uint64(fake.Uint32()),
				MinEstimatedTimeInSeconds: fake.Uint32(),
				MaxEstimatedTimeInSeconds: fake.Uint32(),
				TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
				Notes:                     fake.Word(),
				RecipeID:                  uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStep(ctx, &models.RecipeStepCreationInput{
				Index:                     expected.Index,
				PreparationID:             expected.PreparationID,
				PrerequisiteStep:          expected.PrerequisiteStep,
				MinEstimatedTimeInSeconds: expected.MinEstimatedTimeInSeconds,
				MaxEstimatedTimeInSeconds: expected.MaxEstimatedTimeInSeconds,
				TemperatureInCelsius:      expected.TemperatureInCelsius,
				Notes:                     expected.Notes,
				RecipeID:                  expected.RecipeID,
			})
			checkValueAndError(t, premade, err)

			// Assert recipe step equality
			checkRecipeStepEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveRecipeStep(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetRecipeStep(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe steps
			var expected []*models.RecipeStep
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyRecipeStep(t))
			}

			// Assert recipe step list equality
			actual, err := todoClient.GetRecipeSteps(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeSteps),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeSteps),
			)

			// Clean up
			for _, x := range actual.RecipeSteps {
				err = todoClient.ArchiveRecipeStep(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch recipe step
			_, err := todoClient.GetRecipeStep(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step
			expected := &models.RecipeStep{
				Index:                     uint(fake.Uint32()),
				PreparationID:             uint64(fake.Uint32()),
				PrerequisiteStep:          uint64(fake.Uint32()),
				MinEstimatedTimeInSeconds: fake.Uint32(),
				MaxEstimatedTimeInSeconds: fake.Uint32(),
				TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
				Notes:                     fake.Word(),
				RecipeID:                  uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStep(ctx, &models.RecipeStepCreationInput{
				Index:                     expected.Index,
				PreparationID:             expected.PreparationID,
				PrerequisiteStep:          expected.PrerequisiteStep,
				MinEstimatedTimeInSeconds: expected.MinEstimatedTimeInSeconds,
				MaxEstimatedTimeInSeconds: expected.MaxEstimatedTimeInSeconds,
				TemperatureInCelsius:      expected.TemperatureInCelsius,
				Notes:                     expected.Notes,
				RecipeID:                  expected.RecipeID,
			})
			checkValueAndError(t, premade, err)

			// Fetch recipe step
			actual, err := todoClient.GetRecipeStep(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step equality
			checkRecipeStepEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveRecipeStep(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateRecipeStep(ctx, &models.RecipeStep{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step
			expected := &models.RecipeStep{
				Index:                     uint(fake.Uint32()),
				PreparationID:             uint64(fake.Uint32()),
				PrerequisiteStep:          uint64(fake.Uint32()),
				MinEstimatedTimeInSeconds: fake.Uint32(),
				MaxEstimatedTimeInSeconds: fake.Uint32(),
				TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
				Notes:                     fake.Word(),
				RecipeID:                  uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStep(tctx, &models.RecipeStepCreationInput{
				Index:                     uint(fake.Uint32()),
				PreparationID:             uint64(fake.Uint32()),
				PrerequisiteStep:          uint64(fake.Uint32()),
				MinEstimatedTimeInSeconds: fake.Uint32(),
				MaxEstimatedTimeInSeconds: fake.Uint32(),
				TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
				Notes:                     fake.Word(),
				RecipeID:                  uint64(fake.Uint32()),
			})
			checkValueAndError(t, premade, err)

			// Change recipe step
			premade.Update(expected.ToInput())
			err = todoClient.UpdateRecipeStep(ctx, premade)
			assert.NoError(t, err)

			// Fetch recipe step
			actual, err := todoClient.GetRecipeStep(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step equality
			checkRecipeStepEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveRecipeStep(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe step
			expected := &models.RecipeStep{
				Index:                     uint(fake.Uint32()),
				PreparationID:             uint64(fake.Uint32()),
				PrerequisiteStep:          uint64(fake.Uint32()),
				MinEstimatedTimeInSeconds: fake.Uint32(),
				MaxEstimatedTimeInSeconds: fake.Uint32(),
				TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
				Notes:                     fake.Word(),
				RecipeID:                  uint64(fake.Uint32()),
			}
			premade, err := todoClient.CreateRecipeStep(ctx, &models.RecipeStepCreationInput{
				Index:                     expected.Index,
				PreparationID:             expected.PreparationID,
				PrerequisiteStep:          expected.PrerequisiteStep,
				MinEstimatedTimeInSeconds: expected.MinEstimatedTimeInSeconds,
				MaxEstimatedTimeInSeconds: expected.MaxEstimatedTimeInSeconds,
				TemperatureInCelsius:      expected.TemperatureInCelsius,
				Notes:                     expected.Notes,
				RecipeID:                  expected.RecipeID,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveRecipeStep(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
