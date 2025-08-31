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

		createdUser, testClient := createUserAndClientForTest(t)

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

	T.Run("fails for non-admin user", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		createdUser, testClient := createUserAndClientForTest(t)

		status, err := testClient.GetAuthStatus(ctx, &authsvc.GetAuthStatusRequest{})
		require.NoError(t, err)
		require.NotNil(t, status)

		_, err = testClient.AdminUpdateUserStatus(ctx, &identitysvc.AdminUpdateUserStatusRequest{
			TargetUserID: createdUser.ID,
			NewStatus:    identity.BannedUserAccountStatus.String(),
			Reason:       t.Name(),
		})
		require.Error(t, err)
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

		user, testClient := createUserAndClientForTest(t)

		exampleWebhookInput := &webhooks.WebhookCreationRequestInput{
			ContentType: "application/json",
			Method:      http.MethodPost,
			Name:        t.Name(),
			URL:         "https://whatever.gov",
			Events:      []string{webhooks.WebhookCreatedTriggerEvent},
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

		impersonatedCtx := client.ImpersonateUseAndAccountContext(ctx, user.ID, account.Result.ID)

		t.Logf("impersonating user %s and account %s to get webhook %s", user.ID, account.Result.ID, retrievedWebhook.Result.ID)

		webhook, err := adminClient.GetWebhook(impersonatedCtx, &webhookssvc.GetWebhookRequest{WebhookID: retrievedWebhook.Result.ID})
		assert.NoError(t, err)
		assert.NotNil(t, webhook)
	})

	T.Run("standard user should not be able to impersonate others", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(t)
		_, testClient2 := createUserAndClientForTest(t)

		exampleWebhookInput := &webhooks.WebhookCreationRequestInput{
			ContentType: "application/json",
			Method:      http.MethodPost,
			Name:        t.Name(),
			URL:         "https://whatever.gov",
			Events:      []string{webhooks.WebhookCreatedTriggerEvent},
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

		impersonatedCtx := client.ImpersonateUseAndAccountContext(ctx, user.ID, account.Result.ID)

		t.Logf("impersonating user %s and account %s to get webhook %s", user.ID, account.Result.ID, retrievedWebhook.Result.ID)

		webhook, err := testClient2.GetWebhook(impersonatedCtx, &webhookssvc.GetWebhookRequest{WebhookID: retrievedWebhook.Result.ID})
		assert.Error(t, err)
		assert.Nil(t, webhook)
	})
}
