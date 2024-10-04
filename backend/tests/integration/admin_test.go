package integration

import (
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestAdmin_Returns404WhenModifyingUserAccountStatus() {
	s.runTest("should not be possible to ban a userClient that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			input := fakes.BuildFakeUserAccountStatusUpdateInput()
			input.TargetUserID = nonexistentID

			// Ban userClient.
			assert.Error(t, testClients.adminClient.UpdateUserAccountStatus(ctx, input))
		}
	})
}

func (s *TestSuite) TestAdmin_BanningUsers() {
	s.runTest("should be possible to ban users", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, userClient := createUserAndClientForTest(ctx, t, nil)

			// Assert that userClient can access service
			_, err := userClient.GetWebhooks(ctx, nil)
			require.NoError(t, err)

			input := &types.UserAccountStatusUpdateInput{
				TargetUserID: user.ID,
				NewStatus:    string(types.BannedUserAccountStatus),
				Reason:       "testing",
			}

			assert.NoError(t, testClients.adminClient.UpdateUserAccountStatus(ctx, input))

			// Assert userClient can no longer access service
			_, err = userClient.GetWebhooks(ctx, nil)
			assert.Error(t, err)

			// Clean up.
			assert.NoError(t, testClients.adminClient.ArchiveUser(ctx, user.ID))
		}
	})
}

func (s *TestSuite) TestAdmin_ImpersonatingUsers() {
	s.runTest("should be possible to impersonate users without specifying household ID", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			t.SkipNow() // DELETEME: fix

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)

			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			household, err := testClients.userClient.GetCurrentHousehold(ctx)
			requireNotNilAndNoProblems(t, household, err)

			user, err := testClients.userClient.GetSelf(ctx)
			requireNotNilAndNoProblems(t, user, err)

			// impersonate userClient)
			require.NoError(t, testClients.adminClient.SetOptions(apiclient.ImpersonatingUser(user.ID)))

			webhook, err := testClients.adminClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)

			// Clean up.
			assert.NoError(t, testClients.userClient.ArchiveWebhook(ctx, createdWebhook.ID))
		}
	})

	s.runTest("plain userClient should not be able to impersonate users", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			// t.SkipNow() // DELETEME: address with new client modifications

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, userClient := createUserAndClientForTest(ctx, t, nil)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)

			createdWebhook, err := userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			createdWebhook, err = userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			// impersonate userClient)
			require.NoError(t, testClients.userClient.SetOptions(apiclient.ImpersonatingUser(user.ID)))

			webhook, err := testClients.adminClient.GetWebhook(ctx, createdWebhook.ID)
			assert.Nil(t, webhook)
			assert.Error(t, err)

			// Clean up.
			assert.NoError(t, userClient.ArchiveWebhook(ctx, createdWebhook.ID))
		}
	})
}
