package integration

import (
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (s *TestSuite) TestAdmin_Returns404WhenModifyingUserAccountStatus() {
	s.runTest("should not be possible to ban a user that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			input := fake.BuildFakeForTest[*apiclient.UserAccountStatusUpdateInput](t)
			input.TargetUserID = nonexistentID

			// Ban userClient.
			_, err := testClients.adminClient.AdminUpdateUserStatus(ctx, input)
			assert.Error(t, err)
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

			input := &apiclient.UserAccountStatusUpdateInput{
				TargetUserID: user.ID,
				NewStatus:    string(types.BannedUserAccountStatus),
				Reason:       "testing",
			}

			_, err = testClients.adminClient.AdminUpdateUserStatus(ctx, input)
			assert.Error(t, err)

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
			exampleWebhookInput := &apiclient.WebhookCreationRequestInput{
				ContentType: "application/json",
				Method:      http.MethodPost,
				Name:        t.Name(),
				URL:         "https://whatever.gov",
				Events:      nil,
			}

			createdWebhook, err := testClients.userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			createdWebhook, err = testClients.userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			household, err := testClients.userClient.GetActiveHousehold(ctx)
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

	s.runTest("plain user should not be able to impersonate users", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			// t.SkipNow() // DELETEME: address with new client modifications

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, userClient := createUserAndClientForTest(ctx, t, nil)

			// Create webhook.
			exampleWebhookInput := &apiclient.WebhookCreationRequestInput{
				ContentType: "application/json",
				Method:      http.MethodPost,
				Name:        t.Name(),
				URL:         "https://whatever.gov",
				Events: []string{
					"webhook_created",
				},
			}

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
