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

func checkIngredientEquality(t *testing.T, expected, actual *models.Ingredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Variant, actual.Variant, "expected Variant for ID %d to be %v, but it was %v ", expected.ID, expected.Variant, actual.Variant)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for ID %d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Warning, actual.Warning, "expected Warning for ID %d to be %v, but it was %v ", expected.ID, expected.Warning, actual.Warning)
	assert.Equal(t, expected.ContainsEgg, actual.ContainsEgg, "expected ContainsEgg for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsEgg, actual.ContainsEgg)
	assert.Equal(t, expected.ContainsDairy, actual.ContainsDairy, "expected ContainsDairy for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsDairy, actual.ContainsDairy)
	assert.Equal(t, expected.ContainsPeanut, actual.ContainsPeanut, "expected ContainsPeanut for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsPeanut, actual.ContainsPeanut)
	assert.Equal(t, expected.ContainsTreeNut, actual.ContainsTreeNut, "expected ContainsTreeNut for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsTreeNut, actual.ContainsTreeNut)
	assert.Equal(t, expected.ContainsSoy, actual.ContainsSoy, "expected ContainsSoy for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsSoy, actual.ContainsSoy)
	assert.Equal(t, expected.ContainsWheat, actual.ContainsWheat, "expected ContainsWheat for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsWheat, actual.ContainsWheat)
	assert.Equal(t, expected.ContainsShellfish, actual.ContainsShellfish, "expected ContainsShellfish for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsShellfish, actual.ContainsShellfish)
	assert.Equal(t, expected.ContainsSesame, actual.ContainsSesame, "expected ContainsSesame for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsSesame, actual.ContainsSesame)
	assert.Equal(t, expected.ContainsFish, actual.ContainsFish, "expected ContainsFish for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsFish, actual.ContainsFish)
	assert.Equal(t, expected.ContainsGluten, actual.ContainsGluten, "expected ContainsGluten for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsGluten, actual.ContainsGluten)
	assert.Equal(t, expected.AnimalFlesh, actual.AnimalFlesh, "expected AnimalFlesh for ID %d to be %v, but it was %v ", expected.ID, expected.AnimalFlesh, actual.AnimalFlesh)
	assert.Equal(t, expected.AnimalDerived, actual.AnimalDerived, "expected AnimalDerived for ID %d to be %v, but it was %v ", expected.ID, expected.AnimalDerived, actual.AnimalDerived)
	assert.Equal(t, expected.ConsideredStaple, actual.ConsideredStaple, "expected ConsideredStaple for ID %d to be %v, but it was %v ", expected.ID, expected.ConsideredStaple, actual.ConsideredStaple)
	assert.Equal(t, expected.Icon, actual.Icon, "expected Icon for ID %d to be %v, but it was %v ", expected.ID, expected.Icon, actual.Icon)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyIngredient(t *testing.T) *models.Ingredient {
	t.Helper()

	x := &models.IngredientCreationInput{
		Name:              fake.Word(),
		Variant:           fake.Word(),
		Description:       fake.Word(),
		Warning:           fake.Word(),
		ContainsEgg:       fake.Bool(),
		ContainsDairy:     fake.Bool(),
		ContainsPeanut:    fake.Bool(),
		ContainsTreeNut:   fake.Bool(),
		ContainsSoy:       fake.Bool(),
		ContainsWheat:     fake.Bool(),
		ContainsShellfish: fake.Bool(),
		ContainsSesame:    fake.Bool(),
		ContainsFish:      fake.Bool(),
		ContainsGluten:    fake.Bool(),
		AnimalFlesh:       fake.Bool(),
		AnimalDerived:     fake.Bool(),
		ConsideredStaple:  fake.Bool(),
		Icon:              fake.Word(),
	}
	y, err := todoClient.CreateIngredient(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestIngredients(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create ingredient
			expected := &models.Ingredient{
				Name:              fake.Word(),
				Variant:           fake.Word(),
				Description:       fake.Word(),
				Warning:           fake.Word(),
				ContainsEgg:       fake.Bool(),
				ContainsDairy:     fake.Bool(),
				ContainsPeanut:    fake.Bool(),
				ContainsTreeNut:   fake.Bool(),
				ContainsSoy:       fake.Bool(),
				ContainsWheat:     fake.Bool(),
				ContainsShellfish: fake.Bool(),
				ContainsSesame:    fake.Bool(),
				ContainsFish:      fake.Bool(),
				ContainsGluten:    fake.Bool(),
				AnimalFlesh:       fake.Bool(),
				AnimalDerived:     fake.Bool(),
				ConsideredStaple:  fake.Bool(),
				Icon:              fake.Word(),
			}
			premade, err := todoClient.CreateIngredient(ctx, &models.IngredientCreationInput{
				Name:              expected.Name,
				Variant:           expected.Variant,
				Description:       expected.Description,
				Warning:           expected.Warning,
				ContainsEgg:       expected.ContainsEgg,
				ContainsDairy:     expected.ContainsDairy,
				ContainsPeanut:    expected.ContainsPeanut,
				ContainsTreeNut:   expected.ContainsTreeNut,
				ContainsSoy:       expected.ContainsSoy,
				ContainsWheat:     expected.ContainsWheat,
				ContainsShellfish: expected.ContainsShellfish,
				ContainsSesame:    expected.ContainsSesame,
				ContainsFish:      expected.ContainsFish,
				ContainsGluten:    expected.ContainsGluten,
				AnimalFlesh:       expected.AnimalFlesh,
				AnimalDerived:     expected.AnimalDerived,
				ConsideredStaple:  expected.ConsideredStaple,
				Icon:              expected.Icon,
			})
			checkValueAndError(t, premade, err)

			// Assert ingredient equality
			checkIngredientEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveIngredient(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetIngredient(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkIngredientEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create ingredients
			var expected []*models.Ingredient
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyIngredient(t))
			}

			// Assert ingredient list equality
			actual, err := todoClient.GetIngredients(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Ingredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Ingredients),
			)

			// Clean up
			for _, x := range actual.Ingredients {
				err = todoClient.ArchiveIngredient(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch ingredient
			_, err := todoClient.GetIngredient(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create ingredient
			expected := &models.Ingredient{
				Name:              fake.Word(),
				Variant:           fake.Word(),
				Description:       fake.Word(),
				Warning:           fake.Word(),
				ContainsEgg:       fake.Bool(),
				ContainsDairy:     fake.Bool(),
				ContainsPeanut:    fake.Bool(),
				ContainsTreeNut:   fake.Bool(),
				ContainsSoy:       fake.Bool(),
				ContainsWheat:     fake.Bool(),
				ContainsShellfish: fake.Bool(),
				ContainsSesame:    fake.Bool(),
				ContainsFish:      fake.Bool(),
				ContainsGluten:    fake.Bool(),
				AnimalFlesh:       fake.Bool(),
				AnimalDerived:     fake.Bool(),
				ConsideredStaple:  fake.Bool(),
				Icon:              fake.Word(),
			}
			premade, err := todoClient.CreateIngredient(ctx, &models.IngredientCreationInput{
				Name:              expected.Name,
				Variant:           expected.Variant,
				Description:       expected.Description,
				Warning:           expected.Warning,
				ContainsEgg:       expected.ContainsEgg,
				ContainsDairy:     expected.ContainsDairy,
				ContainsPeanut:    expected.ContainsPeanut,
				ContainsTreeNut:   expected.ContainsTreeNut,
				ContainsSoy:       expected.ContainsSoy,
				ContainsWheat:     expected.ContainsWheat,
				ContainsShellfish: expected.ContainsShellfish,
				ContainsSesame:    expected.ContainsSesame,
				ContainsFish:      expected.ContainsFish,
				ContainsGluten:    expected.ContainsGluten,
				AnimalFlesh:       expected.AnimalFlesh,
				AnimalDerived:     expected.AnimalDerived,
				ConsideredStaple:  expected.ConsideredStaple,
				Icon:              expected.Icon,
			})
			checkValueAndError(t, premade, err)

			// Fetch ingredient
			actual, err := todoClient.GetIngredient(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert ingredient equality
			checkIngredientEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveIngredient(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateIngredient(ctx, &models.Ingredient{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create ingredient
			expected := &models.Ingredient{
				Name:              fake.Word(),
				Variant:           fake.Word(),
				Description:       fake.Word(),
				Warning:           fake.Word(),
				ContainsEgg:       fake.Bool(),
				ContainsDairy:     fake.Bool(),
				ContainsPeanut:    fake.Bool(),
				ContainsTreeNut:   fake.Bool(),
				ContainsSoy:       fake.Bool(),
				ContainsWheat:     fake.Bool(),
				ContainsShellfish: fake.Bool(),
				ContainsSesame:    fake.Bool(),
				ContainsFish:      fake.Bool(),
				ContainsGluten:    fake.Bool(),
				AnimalFlesh:       fake.Bool(),
				AnimalDerived:     fake.Bool(),
				ConsideredStaple:  fake.Bool(),
				Icon:              fake.Word(),
			}
			premade, err := todoClient.CreateIngredient(tctx, &models.IngredientCreationInput{
				Name:              fake.Word(),
				Variant:           fake.Word(),
				Description:       fake.Word(),
				Warning:           fake.Word(),
				ContainsEgg:       fake.Bool(),
				ContainsDairy:     fake.Bool(),
				ContainsPeanut:    fake.Bool(),
				ContainsTreeNut:   fake.Bool(),
				ContainsSoy:       fake.Bool(),
				ContainsWheat:     fake.Bool(),
				ContainsShellfish: fake.Bool(),
				ContainsSesame:    fake.Bool(),
				ContainsFish:      fake.Bool(),
				ContainsGluten:    fake.Bool(),
				AnimalFlesh:       fake.Bool(),
				AnimalDerived:     fake.Bool(),
				ConsideredStaple:  fake.Bool(),
				Icon:              fake.Word(),
			})
			checkValueAndError(t, premade, err)

			// Change ingredient
			premade.Update(expected.ToInput())
			err = todoClient.UpdateIngredient(ctx, premade)
			assert.NoError(t, err)

			// Fetch ingredient
			actual, err := todoClient.GetIngredient(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert ingredient equality
			checkIngredientEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveIngredient(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create ingredient
			expected := &models.Ingredient{
				Name:              fake.Word(),
				Variant:           fake.Word(),
				Description:       fake.Word(),
				Warning:           fake.Word(),
				ContainsEgg:       fake.Bool(),
				ContainsDairy:     fake.Bool(),
				ContainsPeanut:    fake.Bool(),
				ContainsTreeNut:   fake.Bool(),
				ContainsSoy:       fake.Bool(),
				ContainsWheat:     fake.Bool(),
				ContainsShellfish: fake.Bool(),
				ContainsSesame:    fake.Bool(),
				ContainsFish:      fake.Bool(),
				ContainsGluten:    fake.Bool(),
				AnimalFlesh:       fake.Bool(),
				AnimalDerived:     fake.Bool(),
				ConsideredStaple:  fake.Bool(),
				Icon:              fake.Word(),
			}
			premade, err := todoClient.CreateIngredient(ctx, &models.IngredientCreationInput{
				Name:              expected.Name,
				Variant:           expected.Variant,
				Description:       expected.Description,
				Warning:           expected.Warning,
				ContainsEgg:       expected.ContainsEgg,
				ContainsDairy:     expected.ContainsDairy,
				ContainsPeanut:    expected.ContainsPeanut,
				ContainsTreeNut:   expected.ContainsTreeNut,
				ContainsSoy:       expected.ContainsSoy,
				ContainsWheat:     expected.ContainsWheat,
				ContainsShellfish: expected.ContainsShellfish,
				ContainsSesame:    expected.ContainsSesame,
				ContainsFish:      expected.ContainsFish,
				ContainsGluten:    expected.ContainsGluten,
				AnimalFlesh:       expected.AnimalFlesh,
				AnimalDerived:     expected.AnimalDerived,
				ConsideredStaple:  expected.ConsideredStaple,
				Icon:              expected.Icon,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveIngredient(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
