package integration

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated/v2"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkUserNotificationEquality(t *testing.T, expected, actual *types.UserNotification) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Content, actual.Content, "expected Content for userClient notification %s to be %v, but it was %v", expected.ID, expected.Content, actual.Content)
	assert.Equal(t, expected.Status, actual.Status, "expected Status for userClient notification %s to be %v, but it was %v", expected.ID, expected.Status, actual.Status)
	assert.NotZero(t, actual.CreatedAt)
}

func createUserNotificationForTest(t *testing.T, ctx context.Context, forUser *types.User, userClient *apiclient.Client, adminClient *apiclient.Client) *types.UserNotification {
	t.Helper()

	exampleUserNotification := fakes.BuildFakeUserNotification()
	exampleUserNotification.BelongsToUser = forUser.ID
	exampleUserNotificationInput := converters.ConvertUserNotificationToUserNotificationCreationRequestInput(exampleUserNotification)
	createdUserNotification, err := adminClient.CreateUserNotification(ctx, exampleUserNotificationInput)
	require.NoError(t, err)
	checkUserNotificationEquality(t, exampleUserNotification, createdUserNotification)

	createdUserNotification, err = userClient.GetUserNotification(ctx, createdUserNotification.ID)
	requireNotNilAndNoProblems(t, createdUserNotification, err)
	checkUserNotificationEquality(t, exampleUserNotification, createdUserNotification)

	return createdUserNotification
}

func (s *TestSuite) TestUserNotifications_CompleteLifecycle() {
	s.runTest("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, userClient := createUserAndClientForTest(ctx, t, nil)

			exampleUserNotification := fakes.BuildFakeUserNotification()
			exampleUserNotification.BelongsToUser = user.ID
			exampleUserNotificationInput := converters.ConvertUserNotificationToUserNotificationCreationRequestInput(exampleUserNotification)
			_, err := userClient.CreateUserNotification(ctx, exampleUserNotificationInput)
			require.Error(t, err)

			createdUserNotification := createUserNotificationForTest(t, ctx, user, userClient, testClients.adminClient)

			createdUserNotification.Status = types.UserNotificationStatusTypeRead
			updateInput := converters.ConvertUserNotificationToUserNotificationUpdateRequestInput(createdUserNotification)
			createdUserNotification.Update(updateInput)
			require.NoError(t, userClient.UpdateUserNotification(ctx, createdUserNotification.ID, updateInput))

			actual, err := userClient.GetUserNotification(ctx, createdUserNotification.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert userClient notification equality
			checkUserNotificationEquality(t, createdUserNotification, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			updateInput.Status = pointer.To(types.UserNotificationStatusTypeRead)
			assert.NoError(t, userClient.UpdateUserNotification(ctx, createdUserNotification.ID, updateInput))
		}
	})
}

func (s *TestSuite) TestUserNotifications_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, userClient := createUserAndClientForTest(ctx, t, nil)

			var expected []*types.UserNotification
			for i := 0; i < 5; i++ {
				exampleUserNotification := fakes.BuildFakeUserNotification()
				exampleUserNotification.BelongsToUser = user.ID
				exampleUserNotificationInput := converters.ConvertUserNotificationToUserNotificationCreationRequestInput(exampleUserNotification)
				createdUserNotification, createdUserNotificationErr := testClients.adminClient.CreateUserNotification(ctx, exampleUserNotificationInput)
				require.NoError(t, createdUserNotificationErr)

				checkUserNotificationEquality(t, exampleUserNotification, createdUserNotification)

				expected = append(expected, createdUserNotification)
			}

			// assert userClient notification list equality
			actual, err := userClient.GetUserNotifications(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdUserNotification := range expected {
				createdUserNotification.Status = types.UserNotificationStatusTypeRead
				updateInput := converters.ConvertUserNotificationToUserNotificationUpdateRequestInput(createdUserNotification)
				assert.NoError(t, userClient.UpdateUserNotification(ctx, createdUserNotification.ID, updateInput))
			}
		}
	})
}
