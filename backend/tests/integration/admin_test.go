package integration

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/services/webhooks/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdmin_BanningUsers(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		createdUser, testClient := createUserAndClientForTest(t, httpTestServerAddress, grpcTestServerAddress)

		status, err := testClient.GetAuthStatus(ctx, &authsvc.GetAuthStatusRequest{})
		require.NoError(t, err)
		require.NotNil(t, status)

		newStatus := identity.BannedUserAccountStatus.String()

		_, err = adminClient.AdminUpdateUserStatus(ctx, &identitysvc.AdminUpdateUserStatusRequest{
			TargetUserID: createdUser.ID,
			NewStatus:    newStatus,
			Reason:       t.Name(),
		})
		require.NoError(t, err)

		status, err = testClient.GetAuthStatus(ctx, &authsvc.GetAuthStatusRequest{})
		assert.NoError(t, err)
		assert.Equal(t, status.AccountStatus, newStatus)
	})

	T.Run("nonexistent user", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.AdminUpdateUserStatus(ctx, &identitysvc.AdminUpdateUserStatusRequest{
			TargetUserID: nonexistentID,
			NewStatus:    identity.BannedUserAccountStatus.String(),
			Reason:       t.Name(),
		})
		require.Error(t, err)
	})
}

func TestAdmin_UserImpersonation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t, httpTestServerAddress, grpcTestServerAddress)

		// Create webhook.
		exampleWebhookInput := &webhooks.WebhookCreationRequestInput{
			ContentType: "application/json",
			Method:      http.MethodPost,
			Name:        t.Name(),
			URL:         "https://whatever.gov",
			Events:      nil,
		}

		input := converters.ConvertWebhookCreationRequestInputToGRPCWebhookCreationRequestInput(exampleWebhookInput)

		createdWebhook, err := testClient.CreateWebhook(ctx, &webhookssvc.CreateWebhookRequest{Input: input})
		require.NoError(t, err)

		retrievedWebhook, err := testClient.GetWebhook(ctx, &webhookssvc.GetWebhookRequest{WebhookID: createdWebhook.Created.ID})
		require.NoError(t, err)
		require.NotNil(t, retrievedWebhook)

		account, err := testClient.GetActiveAccount(ctx, &authsvc.GetActiveAccountRequest{})
		require.NoError(t, err)
		require.NotNil(t, account)

		user, err := testClient.GetSelf(ctx, &authsvc.GetSelfRequest{})
		require.NoError(t, err)
		require.NotNil(t, user)

		impersonatedCtx := client.ImpersonateUserContext(ctx, user.Result.ID)

		webhook, err := adminClient.GetWebhook(impersonatedCtx, &webhookssvc.GetWebhookRequest{WebhookID: retrievedWebhook.Result.ID})
		assert.NoError(t, err)
		assert.NotNil(t, webhook)
	})

	T.Run("standard user should not be able to impersonate others", func(t *testing.T) {
		t.Parallel()
	})
}

/*

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/platform/fake"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	s.runTest("should be possible to impersonate users without specifying account ID", func(testClients *testClientWrapper) func() {
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

			account, err := testClients.userClient.GetActiveAccount(ctx)
			requireNotNilAndNoProblems(t, account, err)

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

*/
