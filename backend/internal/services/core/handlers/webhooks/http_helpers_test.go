package webhooks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
)

type webhooksServiceHTTPRoutesTestHelper struct {
	ctx                              context.Context
	req                              *http.Request
	res                              *httptest.ResponseRecorder
	service                          *service
	exampleUser                      *types.User
	exampleAccount                   *types.Account
	exampleWebhook                   *types.Webhook
	exampleWebhookTriggerEvent       *types.WebhookTriggerEvent
	exampleCreationInput             *types.WebhookCreationRequestInput
	exampleTriggerEventCreationInput *types.WebhookTriggerEventCreationRequestInput
}

func newTestHelper(t *testing.T) *webhooksServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &webhooksServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleAccount = fakes.BuildFakeAccount()
	helper.exampleAccount.BelongsToUser = helper.exampleUser.ID
	helper.exampleWebhook = fakes.BuildFakeWebhook()
	helper.exampleWebhook.BelongsToAccount = helper.exampleAccount.ID
	helper.exampleWebhookTriggerEvent = fakes.BuildFakeWebhookTriggerEvent()
	helper.exampleWebhookTriggerEvent.BelongsToWebhook = helper.exampleWebhook.ID
	helper.exampleCreationInput = converters.ConvertWebhookToWebhookCreationRequestInput(helper.exampleWebhook)
	helper.exampleTriggerEventCreationInput = converters.ConvertWebhookTriggerEventToWebhookTriggerEventCreationRequestInput(fakes.BuildFakeWebhookTriggerEvent())

	helper.service.webhookIDFetcher = func(*http.Request) string {
		return helper.exampleWebhook.ID
	}

	helper.service.webhookTriggerEventIDFetcher = func(*http.Request) string {
		return helper.exampleWebhookTriggerEvent.ID
	}

	sessionCtxData := &sessions.ContextData{
		Requester: sessions.RequesterInfo{
			UserID:                   helper.exampleUser.ID,
			AccountStatus:            helper.exampleUser.AccountStatus,
			AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
			ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
		},
		ActiveAccountID: helper.exampleAccount.ID,
		AccountPermissions: map[string]authorization.AccountRolePermissionsChecker{
			helper.exampleAccount.ID: authorization.NewAccountRolePermissionChecker(authorization.AccountMemberRole.String()),
		},
	}

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*sessions.ContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), sessions.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
