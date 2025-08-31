package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidIngredientForTest(t *testing.T) *mealplanning.ValidIngredient {
	t.Helper()

	ctx := t.Context()

	creationRequestInput := fakes.BuildFakeValidIngredientCreationRequestInput()
	convertedInput := grpcconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(creationRequestInput)

	created, err := adminClient.CreateValidIngredient(ctx, &mealplanningsvc.CreateValidIngredientRequest{
		Input: convertedInput,
	})
	require.NoError(t, err)
	assert.NotNil(t, created)

	return grpcconverters.ConvertGRPCValidIngredientToValidIngredient(created.Result)
}

func TestValidIngredients_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidIngredientForTest(t)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidIngredientCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(creationRequestInput)

		c := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)
		created, err := c.CreateValidIngredient(ctx, &mealplanningsvc.CreateValidIngredientRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidIngredientCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(creationRequestInput)
		// this is not allowed
		convertedInput.Name = ""

		created, err := adminClient.CreateValidIngredient(ctx, &mealplanningsvc.CreateValidIngredientRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		creationRequestInput := fakes.BuildFakeValidIngredientCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(creationRequestInput)

		created, err := testClient.CreateValidIngredient(ctx, &mealplanningsvc.CreateValidIngredientRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestValidIngredients_Reading(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		retrieved, err := testClient.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientID: created.ID})
		assert.NoError(t, err)

		converted := grpcconverters.ConvertGRPCValidIngredientToValidIngredient(retrieved.Result)

		assertRoughEquality(t, created, converted, "CreatedAt", "LastUpdatedAt", "ArchivedAt")
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

		_, err := c.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestValidIngredients_Updating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		updateInput := fakes.BuildFakeValidIngredientUpdateRequestInput()
		created.Update(updateInput)

		response, err := adminClient.UpdateValidIngredient(ctx, &mealplanningsvc.UpdateValidIngredientRequest{
			ValidIngredientID: created.ID,
			Input:             grpcconverters.ConvertValidIngredientUpdateRequestInputToGRPCValidIngredientUpdateRequestInput(updateInput),
		})
		assert.NoError(t, err)

		updated := response.Result
		// Ensure UpdatedAt was set
		require.NotNil(t, updated.LastUpdatedAt)

		assertRoughEquality(t, created, updated, "CreatedAt", "LastUpdatedAt", "ArchivedAt")
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		updateInput := fakes.BuildFakeValidIngredientUpdateRequestInput()
		created.Update(updateInput)

		c := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

		_, err := c.UpdateValidIngredient(ctx, &mealplanningsvc.UpdateValidIngredientRequest{
			ValidIngredientID: created.ID,
			Input:             grpcconverters.ConvertValidIngredientUpdateRequestInputToGRPCValidIngredientUpdateRequestInput(updateInput),
		})
		assert.Error(t, err)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()

		/*
			there's no way to provide invalid input to this method, but
			I want to make it explicit that tests should be written the moment that changes
		*/
	})

	T.Run("non-admin users are forbidden from updating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		created := createValidIngredientForTest(t)

		response, err := testClient.UpdateValidIngredient(ctx, &mealplanningsvc.UpdateValidIngredientRequest{
			ValidIngredientID: created.ID,
			Input: &mealplanningsvc.ValidIngredientUpdateRequestInput{
				Name: pointer.To("doesn't matter"),
			},
		})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidIngredients_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		_, err := adminClient.ArchiveValidIngredient(ctx, &mealplanningsvc.ArchiveValidIngredientRequest{ValidIngredientID: created.ID})
		assert.NoError(t, err)

		x, err := adminClient.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientID: created.ID})
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

		_, err := c.ArchiveValidIngredient(ctx, &mealplanningsvc.ArchiveValidIngredientRequest{ValidIngredientID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.ArchiveValidIngredient(ctx, &mealplanningsvc.ArchiveValidIngredientRequest{ValidIngredientID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("non-admin users are forbidden from archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)
		_, testClient := createUserAndClientForTest(T)

		_, err := testClient.ArchiveValidIngredient(ctx, &mealplanningsvc.ArchiveValidIngredientRequest{ValidIngredientID: created.ID})
		assert.Error(t, err)
	})
}

func TestValidIngredients_GetRandom(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// in case we haven't already
		createValidIngredientForTest(t)

		response, err := testClient.GetRandomValidIngredient(ctx, &mealplanningsvc.GetRandomValidIngredientRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t, grpcTestServerAddress)

		response, err := c.GetRandomValidIngredient(ctx, &mealplanningsvc.GetRandomValidIngredientRequest{})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

/*

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

*/
