package oauth2clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityfakes "github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth/fakes"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"
)

type oauth2ClientsServiceHTTPRoutesTestHelper struct {
	ctx                 context.Context
	req                 *http.Request
	res                 *httptest.ResponseRecorder
	service             *service
	exampleUser         *identity.User
	exampleAccount      *identity.Account
	exampleOAuth2Client *oauth.OAuth2Client
	exampleInput        *oauth.OAuth2ClientCreationRequestInput
}

func buildTestHelper(t *testing.T) *oauth2ClientsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &oauth2ClientsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService(t)
	helper.exampleUser = identityfakes.BuildFakeUser()
	helper.exampleAccount = identityfakes.BuildFakeAccount()
	helper.exampleAccount.BelongsToUser = helper.exampleUser.ID
	helper.exampleOAuth2Client = fakes.BuildFakeOAuth2Client()
	helper.exampleInput = converters.ConvertOAuth2ClientToOAuth2ClientCreationInput(helper.exampleOAuth2Client)

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
	helper.service.urlClientIDExtractor = func(*http.Request) string {
		return helper.exampleOAuth2Client.ID
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), sessions.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
