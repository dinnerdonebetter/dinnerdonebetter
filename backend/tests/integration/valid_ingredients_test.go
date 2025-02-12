package integration

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpcimpl/converters"
	"github.com/dinnerdonebetter/backend/internal/lib/fake"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidIngredientForTest(t *testing.T, ctx context.Context, adminClient service.EatingServiceClient) *messages.ValidIngredient {
	t.Helper()

	exampleValidIngredient := fake.BuildFakeForTest[*messages.ValidIngredient](t)
	exampleValidIngredientInput := grpcconverters.ConvertValidIngredientToValidIngredientCreationInput(exampleValidIngredient)

	createdValidIngredient, err := adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
	require.NoError(t, err)
	// assertJSONEquality(t, exampleValidIngredient, createdValidIngredient)

	retrievedValidIngredient, err := adminClient.GetValidIngredient(ctx, &messages.GetValidIngredientRequest{ValidIngredientID: createdValidIngredient.Result.ID})
	requireNotNilAndNoProblems(t, retrievedValidIngredient, err)
	// assertJSONEquality(t, exampleValidIngredient, retrievedValidIngredient)

	return createdValidIngredient.Result
}

func (s *TestSuite) TestValidIngredients_CompleteLifecycle() {
	s.runTest("should CRUD", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredient := createValidIngredientForTest(t, ctx, testClients.adminClient)

			newValidIngredient := fake.BuildFakeForTest[*messages.ValidIngredient](t)
			exampleUpdateInput := grpcconverters.ConvertValidIngredientToValidIngredientUpdateInput(newValidIngredient)

			_, err := testClients.adminClient.UpdateValidIngredient(ctx, exampleUpdateInput)
			assert.NoError(t, err)

			actual, err := testClients.adminClient.GetValidIngredient(ctx, &messages.GetValidIngredientRequest{ValidIngredientID: createdValidIngredient.ID})
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient equality
			// assertJSONEquality(t, newValidIngredient, actual.Result)
			assert.NotNil(t, actual.Result.LastUpdatedAt)

			_, err = testClients.adminClient.ArchiveValidIngredient(ctx, &messages.ArchiveValidIngredientRequest{ValidIngredientID: createdValidIngredient.ID})
			assert.NoError(t, err)
		}
	})
}

/*
func (s *TestSuite) TestValidIngredients_GetRandom() {
	s.runTest("should be able to get a random valid ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidIngredientInput := fake.BuildFakeForTest[*messages.CreateValidIngredientRequest](t)

			createdValidIngredient, err := testClients.adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			assert.NoError(t, err)
			assert.NotNil(t, createdValidIngredient)

			retrievedValidIngredient, err := testClients.userClient.GetRandomValidIngredient(ctx, nil)
			requireNotNilAndNoProblems(t, retrievedValidIngredient, err)
			assert.Equal(t, createdValidIngredient.Result.ID, retrievedValidIngredient.Result.ID)

			deleted, err := testClients.adminClient.ArchiveValidIngredient(ctx, &messages.ArchiveValidIngredientRequest{ValidIngredientID: createdValidIngredient.Result.ID})
			assert.NoError(t, err)
			assert.NotNil(t, deleted)
		}
	})
}

func (s *TestSuite) TestValidIngredients_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
				createdValidIngredient, createdValidIngredientErr := testClients.adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, createdValidIngredientErr)

				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

				expected = append(expected, createdValidIngredient)
			}

			// assert valid ingredient list equality
			actual, err := testClients.adminClient.GetValidIngredients(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidIngredient := range expected {
				assert.NoError(t, testClients.adminClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredients_Searching() {
	s.runTest("should be able to be search for valid ingredients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidIngredient
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredient.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidIngredient.Name
			for i := 0; i < 5; i++ {
				exampleValidIngredient.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
				createdValidIngredient, createdValidIngredientErr := testClients.adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, createdValidIngredientErr)
				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

				expected = append(expected, createdValidIngredient)
			}

			// assert valid ingredient list equality
			actual, err := testClients.adminClient.SearchForValidIngredients(ctx, searchQuery, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Results),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Results),
			)

			for _, createdValidIngredient := range expected {
				assert.NoError(t, testClients.adminClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})
}


*/
