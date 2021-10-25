package integration

import (
	"context"
	"testing"
	"time"

	"gitlab.com/prixfixe/prixfixe/pkg/client/httpclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

func checkMealPlanEquality(t *testing.T, expected, actual *types.MealPlan) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.State, actual.State, "expected State for meal plan %s to be %v, but it was %v", expected.ID, expected.State, actual.State)
	assert.Equal(t, expected.StartsAt, actual.StartsAt, "expected StartsAt for meal plan %s to be %v, but it was %v", expected.ID, expected.StartsAt, actual.StartsAt)
	assert.Equal(t, expected.EndsAt, actual.EndsAt, "expected EndsAt for meal plan %s to be %v, but it was %v", expected.ID, expected.EndsAt, actual.EndsAt)
	assert.NotZero(t, actual.CreatedOn)
}

// convertMealPlanToMealPlanUpdateInput creates an MealPlanUpdateRequestInput struct from a meal plan.
func convertMealPlanToMealPlanUpdateInput(x *types.MealPlan) *types.MealPlanUpdateRequestInput {
	return &types.MealPlanUpdateRequestInput{
		Notes:    x.Notes,
		State:    x.State,
		StartsAt: x.StartsAt,
		EndsAt:   x.EndsAt,
	}
}

func createMealPlanWithNotificationChannel(ctx context.Context, t *testing.T, notificationsChan chan *types.DataChangeMessage, client *httpclient.Client) *types.MealPlan {
	t.Helper()

	var n *types.DataChangeMessage

	t.Log("creating meal plan")
	exampleMealPlan := fakes.BuildFakeMealPlan()
	for i := range exampleMealPlan.Options {
		_, _, createdRecipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, client)
		exampleMealPlan.Options[i].RecipeID = createdRecipe.ID
	}

	exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
	createdMealPlanID, err := client.CreateMealPlan(ctx, exampleMealPlanInput)
	require.NoError(t, err)
	t.Logf("meal plan %q created", createdMealPlanID)

	n = <-notificationsChan
	assert.Equal(t, n.DataType, types.MealPlanDataType)
	require.NotNil(t, n.MealPlan)
	checkMealPlanEquality(t, exampleMealPlan, n.MealPlan)

	createdMealPlan, err := client.GetMealPlan(ctx, createdMealPlanID)
	requireNotNilAndNoProblems(t, createdMealPlan, err)
	checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

	return createdMealPlan
}

func createMealPlanWhilePolling(ctx context.Context, t *testing.T, client *httpclient.Client) *types.MealPlan {
	t.Helper()

	var checkFunc func() bool

	t.Log("creating meal plan")
	exampleMealPlan := fakes.BuildFakeMealPlan()
	for i := range exampleMealPlan.Options {
		_, _, createdRecipe := createRecipeWithPolling(ctx, t, client)
		exampleMealPlan.Options[i].RecipeID = createdRecipe.ID
	}

	exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
	createdMealPlanID, err := client.CreateMealPlan(ctx, exampleMealPlanInput)
	require.NoError(t, err)
	t.Logf("meal plan %q created", createdMealPlanID)

	var createdMealPlan *types.MealPlan
	checkFunc = func() bool {
		createdMealPlan, err = client.GetMealPlan(ctx, createdMealPlanID)
		return assert.NotNil(t, createdMealPlan) && assert.NoError(t, err)
	}
	assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
	checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

	return createdMealPlan
}

func (s *TestSuite) TestMealPlans_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			createdMealPlan := createMealPlanWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			t.Log("changing meal plan")
			newMealPlan := fakes.BuildFakeMealPlan()
			newMealPlan.Options = createdMealPlan.Options
			createdMealPlan.Update(convertMealPlanToMealPlanUpdateInput(newMealPlan))
			assert.NoError(t, testClients.main.UpdateMealPlan(ctx, createdMealPlan))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.MealPlanDataType)

			t.Log("fetching changed meal plan")
			actual, err := testClients.main.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan equality
			checkMealPlanEquality(t, newMealPlan, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanWhilePolling(ctx, t, testClients.main)

			// change meal plan
			newMealPlan := fakes.BuildFakeMealPlan()
			newMealPlan.Options = createdMealPlan.Options
			createdMealPlan.Update(convertMealPlanToMealPlanUpdateInput(newMealPlan))
			assert.NoError(t, testClients.main.UpdateMealPlan(ctx, createdMealPlan))

			time.Sleep(time.Second)

			// retrieve changed meal plan
			var (
				actual *types.MealPlan
				err    error
			)
			checkFunc = func() bool {
				actual, err = testClients.main.GetMealPlan(ctx, createdMealPlan.ID)
				return assert.NotNil(t, createdMealPlan) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan equality
			checkMealPlanEquality(t, newMealPlan, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}

func (s *TestSuite) TestMealPlans_Listing() {
	s.runForCookieClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			t.Log("creating meal plans")
			var expected []*types.MealPlan
			for i := 0; i < 5; i++ {
				createdMealPlan := createMealPlanWithNotificationChannel(ctx, t, notificationsChan, testClients.main)
				expected = append(expected, createdMealPlan)
			}

			// assert meal plan list equality
			actual, err := testClients.main.GetMealPlans(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.MealPlans),
				"expected %d to be <= %d",
				len(expected),
				len(actual.MealPlans),
			)

			t.Log("cleaning up")
			for _, createdMealPlan := range expected {
				assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlan.ID))
			}
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating meal plans")
			var expected []*types.MealPlan
			for i := 0; i < 5; i++ {
				createdMealPlan := createMealPlanWhilePolling(ctx, t, testClients.main)
				expected = append(expected, createdMealPlan)
			}

			// assert meal plan list equality
			actual, err := testClients.main.GetMealPlans(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.MealPlans),
				"expected %d to be <= %d",
				len(expected),
				len(actual.MealPlans),
			)

			t.Log("cleaning up")
			for _, createdMealPlan := range expected {
				assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlan.ID))
			}
		}
	})
}
