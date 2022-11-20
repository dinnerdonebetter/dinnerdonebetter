package integration

import (
	"fmt"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func checkValidPreparationEquality(t *testing.T, expected, actual *types.ValidPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid preparation %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid preparation %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid preparation %s to be %v, but it was %v", expected.ID, expected.IconPath, actual.IconPath)
	assert.Equal(t, expected.PastTense, actual.PastTense, "expected PastTense for valid preparation %s to be %v, but it was %v", expected.ID, expected.PastTense, actual.PastTense)
	assert.Equal(t, expected.YieldsNothing, actual.YieldsNothing, "expected YieldsNothing for valid preparation %s to be %v, but it was %v", expected.ID, expected.YieldsNothing, actual.YieldsNothing)
	assert.Equal(t, expected.RestrictToIngredients, actual.RestrictToIngredients, "expected RestrictToIngredients for valid preparation %s to be %v, but it was %v", expected.ID, expected.RestrictToIngredients, actual.RestrictToIngredients)
	assert.Equal(t, expected.ZeroIngredientsAllowable, actual.ZeroIngredientsAllowable, "expected ZeroIngredientsAllowable for valid preparation %s to be %v, but it was %v", expected.ID, expected.ZeroIngredientsAllowable, actual.ZeroIngredientsAllowable)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for valid preparation %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestValidPreparations_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.admin.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("changing valid preparation")
			newValidPreparation := fakes.BuildFakeValidPreparation()
			createdValidPreparation.Update(converters.ConvertValidPreparationToValidPreparationUpdateRequestInput(newValidPreparation))
			assert.NoError(t, testClients.admin.UpdateValidPreparation(ctx, createdValidPreparation))

			t.Log("fetching changed valid preparation")
			actual, err := testClients.admin.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation equality
			checkValidPreparationEquality(t, newValidPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up valid preparation")
			assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_GetRandom() {
	s.runForEachClient("should be able to get a random valid preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.admin.GetRandomValidPreparation(ctx)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			t.Log("cleaning up valid preparation")
			assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparations")
			var expected []*types.ValidPreparation
			for i := 0; i < 5; i++ {
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
				createdValidPreparation, createdValidPreparationErr := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, createdValidPreparationErr)
				t.Logf("valid preparation %q created", createdValidPreparation.ID)

				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				expected = append(expected, createdValidPreparation)
			}

			// assert valid preparation list equality
			actual, err := testClients.admin.GetValidPreparations(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			t.Log("cleaning up")
			for _, createdValidPreparation := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidPreparations_Searching() {
	s.runForEachClient("should be able to be search for valid preparations", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparations")
			var expected []*types.ValidPreparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparation.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidPreparation.Name
			for i := 0; i < 5; i++ {
				exampleValidPreparation.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
				createdValidPreparation, createdValidPreparationErr := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, createdValidPreparationErr)
				t.Logf("valid preparation %q created", createdValidPreparation.ID)
				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				expected = append(expected, createdValidPreparation)
			}

			exampleLimit := uint8(20)

			// assert valid preparation list equality
			actual, err := testClients.admin.SearchValidPreparations(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdValidPreparation := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})
}
