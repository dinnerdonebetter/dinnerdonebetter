package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, expected.IngredientCount, actual.IngredientCount, "expected IngredientCount for valid preparation %s to be %v, but it was %v", expected.ID, expected.IngredientCount, actual.IngredientCount)
	assert.Equal(t, expected.InstrumentCount, actual.InstrumentCount, "expected InstrumentCount for valid preparation %s to be %v, but it was %v", expected.ID, expected.InstrumentCount, actual.InstrumentCount)
	assert.Equal(t, expected.VesselCount, actual.VesselCount, "expected VesselCount for valid preparation %s to be %v, but it was %v", expected.ID, expected.VesselCount, actual.VesselCount)
	assert.Equal(t, expected.TemperatureRequired, actual.TemperatureRequired, "expected TemperatureRequired for valid preparation %s to be %v, but it was %v", expected.ID, expected.TemperatureRequired, actual.TemperatureRequired)
	assert.Equal(t, expected.TimeEstimateRequired, actual.TimeEstimateRequired, "expected TimeEstimateRequired for valid preparation %s to be %v, but it was %v", expected.ID, expected.TimeEstimateRequired, actual.TimeEstimateRequired)
	assert.Equal(t, expected.ConditionExpressionRequired, actual.ConditionExpressionRequired, "expected ConditionExpressionRequired for valid preparation %s to be %v, but it was %v", expected.ID, expected.ConditionExpressionRequired, actual.ConditionExpressionRequired)
	assert.Equal(t, expected.ConsumesVessel, actual.ConsumesVessel, "expected ConsumesVessel for valid preparation %s to be %v, but it was %v", expected.ID, expected.ConsumesVessel, actual.ConsumesVessel)
	assert.Equal(t, expected.OnlyForVessels, actual.OnlyForVessels, "expected OnlyForVessels for valid preparation %s to be %v, but it was %v", expected.ID, expected.OnlyForVessels, actual.OnlyForVessels)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for valid preparation %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.NotZero(t, actual.CreatedAt)
}

func createValidPreparationForTest(t *testing.T, ctx context.Context, vessel *types.ValidPreparation, adminClient *apiclient.Client) *types.ValidPreparation {
	t.Helper()

	exampleValidPreparation := vessel
	if exampleValidPreparation == nil {
		exampleValidPreparation = fakes.BuildFakeValidPreparation()
	}

	exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
	createdValidPreparation, err := adminClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
	require.NoError(t, err)
	checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

	createdValidPreparation, err = adminClient.GetValidPreparation(ctx, createdValidPreparation.ID)
	requireNotNilAndNoProblems(t, createdValidPreparation, err)
	checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

	return createdValidPreparation
}

func (s *TestSuite) TestValidPreparations_CompleteLifecycle() {
	s.runTest("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			createdValidPreparation := createValidPreparationForTest(t, ctx, exampleValidPreparation, testClients.adminClient)

			newValidPreparation := fakes.BuildFakeValidPreparation()
			updateInput := converters.ConvertValidPreparationToValidPreparationUpdateRequestInput(newValidPreparation)
			createdValidPreparation.Update(updateInput)
			assert.NoError(t, testClients.adminClient.UpdateValidPreparation(ctx, createdValidPreparation.ID, updateInput))

			actual, err := testClients.adminClient.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation equality
			checkValidPreparationEquality(t, newValidPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.adminClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_GetRandom() {
	s.runTest("should be able to get a random valid preparation", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.adminClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.adminClient.GetRandomValidPreparation(ctx)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			assert.NoError(t, testClients.adminClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparations_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidPreparation
			for i := 0; i < 5; i++ {
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
				createdValidPreparation, createdValidPreparationErr := testClients.adminClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, createdValidPreparationErr)

				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				expected = append(expected, createdValidPreparation)
			}

			// assert valid preparation list equality
			actual, err := testClients.adminClient.GetValidPreparations(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidPreparation := range expected {
				assert.NoError(t, testClients.adminClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidPreparations_Searching() {
	s.runTest("should be able to be search for valid preparations", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidPreparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparation.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidPreparation.Name
			for i := 0; i < 5; i++ {
				exampleValidPreparation.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
				createdValidPreparation, createdValidPreparationErr := testClients.adminClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, createdValidPreparationErr)
				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				expected = append(expected, createdValidPreparation)
			}

			filter := types.DefaultQueryFilter()
			filter.Limit = pointer.To(uint8(20))

			// assert valid preparation list equality
			actual, err := testClients.adminClient.SearchForValidPreparations(ctx, searchQuery, filter)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidPreparation := range expected {
				assert.NoError(t, testClients.adminClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
			}
		}
	})
}
