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

func checkIterationMediaEquality(t *testing.T, expected, actual *models.IterationMedia) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Path, actual.Path, "expected Path for ID %d to be %v, but it was %v ", expected.ID, expected.Path, actual.Path)
	assert.Equal(t, expected.Mimetype, actual.Mimetype, "expected Mimetype for ID %d to be %v, but it was %v ", expected.ID, expected.Mimetype, actual.Mimetype)
	assert.Equal(t, expected.RecipeIterationID, actual.RecipeIterationID, "expected RecipeIterationID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeIterationID, actual.RecipeIterationID)
	assert.Equal(t, *expected.RecipeStepID, *actual.RecipeStepID, "expected RecipeStepID to be %v, but it was %v ", expected.RecipeStepID, actual.RecipeStepID)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyIterationMedia(t *testing.T) *models.IterationMedia {
	t.Helper()

	x := &models.IterationMediaCreationInput{
		Path:              fake.Word(),
		Mimetype:          fake.Word(),
		RecipeIterationID: uint64(fake.Uint32()),
		RecipeStepID:      func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
	}
	y, err := todoClient.CreateIterationMedia(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestIterationMedias(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create iteration media
			expected := &models.IterationMedia{
				Path:              fake.Word(),
				Mimetype:          fake.Word(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			}
			premade, err := todoClient.CreateIterationMedia(ctx, &models.IterationMediaCreationInput{
				Path:              expected.Path,
				Mimetype:          expected.Mimetype,
				RecipeIterationID: expected.RecipeIterationID,
				RecipeStepID:      expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Assert iteration media equality
			checkIterationMediaEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveIterationMedia(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetIterationMedia(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkIterationMediaEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create iteration medias
			var expected []*models.IterationMedia
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyIterationMedia(t))
			}

			// Assert iteration media list equality
			actual, err := todoClient.GetIterationMedias(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.IterationMedias),
				"expected %d to be <= %d",
				len(expected),
				len(actual.IterationMedias),
			)

			// Clean up
			for _, x := range actual.IterationMedias {
				err = todoClient.ArchiveIterationMedia(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch iteration media
			_, err := todoClient.GetIterationMedia(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create iteration media
			expected := &models.IterationMedia{
				Path:              fake.Word(),
				Mimetype:          fake.Word(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			}
			premade, err := todoClient.CreateIterationMedia(ctx, &models.IterationMediaCreationInput{
				Path:              expected.Path,
				Mimetype:          expected.Mimetype,
				RecipeIterationID: expected.RecipeIterationID,
				RecipeStepID:      expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Fetch iteration media
			actual, err := todoClient.GetIterationMedia(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert iteration media equality
			checkIterationMediaEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveIterationMedia(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateIterationMedia(ctx, &models.IterationMedia{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create iteration media
			expected := &models.IterationMedia{
				Path:              fake.Word(),
				Mimetype:          fake.Word(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			}
			premade, err := todoClient.CreateIterationMedia(tctx, &models.IterationMediaCreationInput{
				Path:              fake.Word(),
				Mimetype:          fake.Word(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			})
			checkValueAndError(t, premade, err)

			// Change iteration media
			premade.Update(expected.ToInput())
			err = todoClient.UpdateIterationMedia(ctx, premade)
			assert.NoError(t, err)

			// Fetch iteration media
			actual, err := todoClient.GetIterationMedia(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert iteration media equality
			checkIterationMediaEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveIterationMedia(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create iteration media
			expected := &models.IterationMedia{
				Path:              fake.Word(),
				Mimetype:          fake.Word(),
				RecipeIterationID: uint64(fake.Uint32()),
				RecipeStepID:      func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			}
			premade, err := todoClient.CreateIterationMedia(ctx, &models.IterationMediaCreationInput{
				Path:              expected.Path,
				Mimetype:          expected.Mimetype,
				RecipeIterationID: expected.RecipeIterationID,
				RecipeStepID:      expected.RecipeStepID,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveIterationMedia(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
