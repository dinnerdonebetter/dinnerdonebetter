package integration

import (
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkValidIngredientGroupEquality(t *testing.T, expected, actual *types.ValidIngredientGroup) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid ingredient group %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid ingredient group %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for valid ingredient group %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestValidIngredientGroups_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredient group")
			exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
			exampleValidIngredientGroupInput := converters.ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(exampleValidIngredientGroup)
			createdValidIngredientGroup, err := testClients.admin.CreateValidIngredientGroup(ctx, exampleValidIngredientGroupInput)
			require.NoError(t, err)
			t.Logf("valid ingredient group %q created", createdValidIngredientGroup.ID)
			checkValidIngredientGroupEquality(t, exampleValidIngredientGroup, createdValidIngredientGroup)

			createdValidIngredientGroup, err = testClients.admin.GetValidIngredientGroup(ctx, createdValidIngredientGroup.ID)
			requireNotNilAndNoProblems(t, createdValidIngredientGroup, err)
			checkValidIngredientGroupEquality(t, exampleValidIngredientGroup, createdValidIngredientGroup)

			t.Log("changing valid ingredient group")
			newValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
			createdValidIngredientGroup.Update(converters.ConvertValidIngredientGroupToValidIngredientGroupUpdateRequestInput(newValidIngredientGroup))
			assert.NoError(t, testClients.admin.UpdateValidIngredientGroup(ctx, createdValidIngredientGroup))

			t.Log("fetching changed valid ingredient group")
			actual, err := testClients.admin.GetValidIngredientGroup(ctx, createdValidIngredientGroup.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient group equality
			checkValidIngredientGroupEquality(t, newValidIngredientGroup, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up valid ingredient group")
			assert.NoError(t, testClients.admin.ArchiveValidIngredientGroup(ctx, createdValidIngredientGroup.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientGroups_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredient groups")
			var expected []*types.ValidIngredientGroup
			for i := 0; i < 5; i++ {
				exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
				exampleValidIngredientGroupInput := converters.ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(exampleValidIngredientGroup)
				createdValidIngredientGroup, createdValidIngredientGroupErr := testClients.admin.CreateValidIngredientGroup(ctx, exampleValidIngredientGroupInput)
				require.NoError(t, createdValidIngredientGroupErr)
				t.Logf("valid ingredient group %q created", createdValidIngredientGroup.ID)

				checkValidIngredientGroupEquality(t, exampleValidIngredientGroup, createdValidIngredientGroup)

				expected = append(expected, createdValidIngredientGroup)
			}

			// assert valid ingredient group list equality
			actual, err := testClients.admin.GetValidIngredientGroups(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			t.Log("cleaning up")
			for _, createdValidIngredientGroup := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredientGroup(ctx, createdValidIngredientGroup.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredientGroups_Searching() {
	s.runForEachClient("should be able to be search for valid ingredient groups", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredient groups")
			var expected []*types.ValidIngredientGroup
			exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
			exampleValidIngredientGroup.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidIngredientGroup.Name
			for i := 0; i < 5; i++ {
				exampleValidIngredientGroup.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidIngredientGroupInput := converters.ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(exampleValidIngredientGroup)
				createdValidIngredientGroup, createdValidIngredientGroupErr := testClients.admin.CreateValidIngredientGroup(ctx, exampleValidIngredientGroupInput)
				require.NoError(t, createdValidIngredientGroupErr)
				checkValidIngredientGroupEquality(t, exampleValidIngredientGroup, createdValidIngredientGroup)
				t.Logf("valid ingredient group %q created", createdValidIngredientGroup.ID)

				expected = append(expected, createdValidIngredientGroup)
			}

			exampleLimit := uint8(20)

			// assert valid ingredient group list equality
			actual, err := testClients.admin.SearchValidIngredientGroups(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdValidIngredientGroup := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredientGroup(ctx, createdValidIngredientGroup.ID))
			}
		}
	})
}
