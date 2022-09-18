package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkMealPlanEventEquality(t *testing.T, expected, actual *types.MealPlanEvent) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for meal plan event %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.StartsAt, actual.StartsAt, "expected StartsAt for meal plan event %s to be %v, but it was %v", expected.ID, expected.StartsAt, actual.StartsAt)
	assert.Equal(t, expected.EndsAt, actual.EndsAt, "expected EndsAt for meal plan event %s to be %v, but it was %v", expected.ID, expected.EndsAt, actual.EndsAt)
	assert.Equal(t, expected.MealName, actual.MealName, "expected MealName for meal plan event %s to be %v, but it was %v", expected.ID, expected.MealName, actual.MealName)
	assert.Equal(t, expected.BelongsToMealPlan, actual.BelongsToMealPlan, "expected BelongsToMealPlan for meal plan event %s to be %v, but it was %v", expected.ID, expected.BelongsToMealPlan, actual.BelongsToMealPlan)
	assert.NotZero(t, actual.CreatedAt)
}

// convertMealPlanEventToMealPlanEventUpdateInput creates an MealPlanEventUpdateRequestInput struct from a meal plan event.
func convertMealPlanEventToMealPlanEventUpdateInput(x *types.MealPlanEvent) *types.MealPlanEventUpdateRequestInput {
	return &types.MealPlanEventUpdateRequestInput{
		Notes:             &x.Notes,
		StartsAt:          &x.StartsAt,
		EndsAt:            &x.EndsAt,
		MealName:          &x.MealName,
		BelongsToMealPlan: x.BelongsToMealPlan,
	}
}

func (s *TestSuite) TestMealPlanEvents_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanForTest(ctx, t, nil, testClients.admin, testClients.user)

			require.NotNil(t, createdMealPlan)
			require.NotEmpty(t, createdMealPlan.Events)
			require.NotNil(t, createdMealPlan.Events[0])
			createdMealPlanEvent := createdMealPlan.Events[0]

			t.Log("changing meal plan event")
			newMealPlanEvent := fakes.BuildFakeMealPlanEvent()
			newMealPlanEvent.BelongsToMealPlan = createdMealPlan.ID

			createdMealPlanEvent.Update(convertMealPlanEventToMealPlanEventUpdateInput(newMealPlanEvent))
			assert.NoError(t, testClients.user.UpdateMealPlanEvent(ctx, createdMealPlanEvent))

			t.Log("fetching changed meal plan event")
			actual, err := testClients.user.GetMealPlanEvent(ctx, createdMealPlan.ID, createdMealPlanEvent.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan event equality
			checkMealPlanEventEquality(t, newMealPlanEvent, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up meal plan event")
			assert.NoError(t, testClients.user.ArchiveMealPlanEvent(ctx, createdMealPlan.ID, createdMealPlanEvent.ID))

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.user.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}

func (s *TestSuite) TestMealPlanEvents_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleMealPlan := fakes.BuildFakeMealPlan()
			exampleMealPlan.Events = nil
			createdMealPlan := createMealPlanForTest(ctx, t, exampleMealPlan, testClients.admin, testClients.user)

			t.Log("creating meal plan events")
			var expected []*types.MealPlanEvent
			for i := 0; i < 5; i++ {
				exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
				exampleMealPlanEvent.Options = nil
				exampleMealPlanEvent.BelongsToMealPlan = createdMealPlan.ID

				exampleMealPlanEventInput := fakes.BuildFakeMealPlanEventCreationRequestInputFromMealPlanEvent(exampleMealPlanEvent)
				createdMealPlanEvent, err := testClients.user.CreateMealPlanEvent(ctx, createdMealPlan.ID, exampleMealPlanEventInput)
				require.NoError(t, err)
				t.Logf("meal plan event %q created", createdMealPlanEvent.ID)

				rawBytes, err := json.MarshalIndent(createdMealPlanEvent, "", "\t")
				require.NoError(t, err)
				t.Log(string(rawBytes))

				checkMealPlanEventEquality(t, exampleMealPlanEvent, createdMealPlanEvent)

				createdMealPlanEvent, err = testClients.user.GetMealPlanEvent(ctx, createdMealPlan.ID, createdMealPlanEvent.ID)
				requireNotNilAndNoProblems(t, createdMealPlanEvent, err)
				require.Equal(t, createdMealPlan.ID, createdMealPlanEvent.BelongsToMealPlan)

				expected = append(expected, createdMealPlanEvent)
			}

			// assert meal plan event list equality
			actual, err := testClients.user.GetMealPlanEvents(ctx, createdMealPlan.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.MealPlanEvents),
				"expected %d to be <= %d",
				len(expected),
				len(actual.MealPlanEvents),
			)

			t.Log("cleaning up")
			for _, createdMealPlanEvent := range expected {
				assert.NoError(t, testClients.user.ArchiveMealPlanEvent(ctx, createdMealPlan.ID, createdMealPlanEvent.ID))
			}

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.user.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}
