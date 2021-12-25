package webhooks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

type webhooksServiceHTTPRoutesTestHelper struct {
	ctx                  context.Context
	req                  *http.Request
	res                  *httptest.ResponseRecorder
	service              *service
	exampleUser          *types.User
	exampleHousehold     *types.Household
	exampleWebhook       *types.Webhook
	exampleCreationInput *types.WebhookCreationRequestInput
}

func newTestHelper(t *testing.T) *webhooksServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &webhooksServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleWebhook = fakes.BuildFakeWebhook()
	helper.exampleWebhook.BelongsToHousehold = helper.exampleHousehold.ID
	helper.exampleCreationInput = fakes.BuildFakeWebhookCreationInputFromWebhook(helper.exampleWebhook)

	helper.service.webhookIDFetcher = func(*http.Request) string {
		return helper.exampleWebhook.ID
	}

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
			UserID:                helper.exampleUser.ID,
			Reputation:            helper.exampleUser.ServiceHouseholdStatus,
			ReputationExplanation: helper.exampleUser.ReputationExplanation,
			ServicePermissions:    authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRoles...),
		},
		ActiveHouseholdID: helper.exampleHousehold.ID,
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
			helper.exampleHousehold.ID: authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdMemberRole.String()),
		},
	}

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), types.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
