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
	s.runForEachClient("should not be possible to ban a user that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			input := fakes.BuildFakeUserAccountStatusUpdateInput()
			input.TargetUserID = nonexistentID

			// Ban user.
			assert.Error(t, testClients.admin.UpdateUserAccountStatus(ctx, input))
		}
	})
}

func (s *TestSuite) TestAdmin_BanningUsers() {
	s.runForEachClient("should be possible to ban users", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, _, userClient, _ := createUserAndClientForTest(ctx, t, nil)

			// Assert that user can access service
			_, err := userClient.GetWebhooks(ctx, nil)
			require.NoError(t, err)

			input := &types.UserAccountStatusUpdateInput{
				TargetUserID: user.ID,
				NewStatus:    string(types.BannedUserAccountStatus),
				Reason:       "testing",
			}

			assert.NoError(t, testClients.admin.UpdateUserAccountStatus(ctx, input))

			// Assert user can no longer access service
			_, err = userClient.GetWebhooks(ctx, nil)
			assert.Error(t, err)

			// Clean up.
			assert.NoError(t, testClients.admin.ArchiveUser(ctx, user.ID))
		}
	})
}

func (s *TestSuite) TestAdmin_ImpersonatingUsers() {
	s.runForCookieClient("should be possible to impersonate users without specifying household ID", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			t.SkipNow() // DELETEME: fix

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)

			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			household, err := testClients.user.GetCurrentHousehold(ctx)
			requireNotNilAndNoProblems(t, household, err)

			user, err := testClients.user.GetSelf(ctx)
			requireNotNilAndNoProblems(t, user, err)

			// impersonate user)
			require.NoError(t, testClients.admin.SetOptions(apiclient.ImpersonatingUser(user.ID)))

			webhook, err := testClients.admin.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)

			// Clean up.
			assert.NoError(t, testClients.user.ArchiveWebhook(ctx, createdWebhook.ID))
		}
	})

	// The OAuth2 client always uses the default household
	s.runForCookieClient("should be possible to impersonate users while specifying household ID", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			t.SkipNow() // TODO: figure out what is going wrong here

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdHousehold, err := testClients.user.CreateHousehold(ctx, fakes.BuildFakeHouseholdCreationInput())
			requireNotNilAndNoProblems(t, createdHousehold, err)

			require.NoError(t, testClients.user.SwitchActiveHousehold(ctx, createdHousehold.ID))

			user, err := testClients.user.GetSelf(ctx)
			requireNotNilAndNoProblems(t, user, err)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)

			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			// impersonate user)
			require.NoError(t, testClients.admin.SetOptions(apiclient.ImpersonatingUser(user.ID), apiclient.ImpersonatingHousehold(createdHousehold.ID)))

			webhook, err := testClients.admin.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)

			// Clean up.
			assert.NoError(t, testClients.user.ArchiveWebhook(ctx, createdWebhook.ID))
		}
	})

	s.runForCookieClient("plain user should not be able to impersonate users", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			// t.SkipNow() // DELETEME: address with new client modifications

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			user, _, userClient, _ := createUserAndClientForTest(ctx, t, nil)

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)

			createdWebhook, err := userClient.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			createdWebhook, err = userClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			// impersonate user)
			require.NoError(t, testClients.user.SetOptions(apiclient.ImpersonatingUser(user.ID)))

			webhook, err := testClients.admin.GetWebhook(ctx, createdWebhook.ID)
			assert.Nil(t, webhook)
			assert.Error(t, err)

			// Clean up.
			assert.NoError(t, userClient.ArchiveWebhook(ctx, createdWebhook.ID))
		}
	})
}
