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

func checkRecipeEquality(t *testing.T, expected, actual *models.Recipe) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Source, actual.Source, "expected Source for ID %d to be %v, but it was %v ", expected.ID, expected.Source, actual.Source)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for ID %d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, *expected.InspiredByRecipeID, *actual.InspiredByRecipeID, "expected InspiredByRecipeID to be %v, but it was %v ", expected.InspiredByRecipeID, actual.InspiredByRecipeID)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyRecipe(t *testing.T) *models.Recipe {
	t.Helper()

	x := &models.RecipeCreationInput{
		Name:               fake.Word(),
		Source:             fake.Word(),
		Description:        fake.Word(),
		InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
	}
	y, err := todoClient.CreateRecipe(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestRecipes(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe
			expected := &models.Recipe{
				Name:               fake.Word(),
				Source:             fake.Word(),
				Description:        fake.Word(),
				InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			}
			premade, err := todoClient.CreateRecipe(ctx, &models.RecipeCreationInput{
				Name:               expected.Name,
				Source:             expected.Source,
				Description:        expected.Description,
				InspiredByRecipeID: expected.InspiredByRecipeID,
			})
			checkValueAndError(t, premade, err)

			// Assert recipe equality
			checkRecipeEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveRecipe(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetRecipe(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkRecipeEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipes
			var expected []*models.Recipe
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyRecipe(t))
			}

			// Assert recipe list equality
			actual, err := todoClient.GetRecipes(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Recipes),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Recipes),
			)

			// Clean up
			for _, x := range actual.Recipes {
				err = todoClient.ArchiveRecipe(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch recipe
			_, err := todoClient.GetRecipe(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe
			expected := &models.Recipe{
				Name:               fake.Word(),
				Source:             fake.Word(),
				Description:        fake.Word(),
				InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			}
			premade, err := todoClient.CreateRecipe(ctx, &models.RecipeCreationInput{
				Name:               expected.Name,
				Source:             expected.Source,
				Description:        expected.Description,
				InspiredByRecipeID: expected.InspiredByRecipeID,
			})
			checkValueAndError(t, premade, err)

			// Fetch recipe
			actual, err := todoClient.GetRecipe(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe equality
			checkRecipeEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveRecipe(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateRecipe(ctx, &models.Recipe{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe
			expected := &models.Recipe{
				Name:               fake.Word(),
				Source:             fake.Word(),
				Description:        fake.Word(),
				InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			}
			premade, err := todoClient.CreateRecipe(tctx, &models.RecipeCreationInput{
				Name:               fake.Word(),
				Source:             fake.Word(),
				Description:        fake.Word(),
				InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			})
			checkValueAndError(t, premade, err)

			// Change recipe
			premade.Update(expected.ToInput())
			err = todoClient.UpdateRecipe(ctx, premade)
			assert.NoError(t, err)

			// Fetch recipe
			actual, err := todoClient.GetRecipe(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe equality
			checkRecipeEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveRecipe(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create recipe
			expected := &models.Recipe{
				Name:               fake.Word(),
				Source:             fake.Word(),
				Description:        fake.Word(),
				InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			}
			premade, err := todoClient.CreateRecipe(ctx, &models.RecipeCreationInput{
				Name:               expected.Name,
				Source:             expected.Source,
				Description:        expected.Description,
				InspiredByRecipeID: expected.InspiredByRecipeID,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveRecipe(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
