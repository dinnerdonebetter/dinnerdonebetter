package oauth2clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessioncontext"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
)

type oauth2ClientsServiceHTTPRoutesTestHelper struct {
	ctx                 context.Context
	req                 *http.Request
	res                 *httptest.ResponseRecorder
	service             *service
	exampleUser         *types.User
	exampleHousehold    *types.Household
	exampleOAuth2Client *types.OAuth2Client
	exampleInput        *types.OAuth2ClientCreationRequestInput
}

func buildTestHelper(t *testing.T) *oauth2ClientsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &oauth2ClientsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService(t)
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleOAuth2Client = fakes.BuildFakeOAuth2Client()
	helper.exampleInput = converters.ConvertOAuth2ClientToOAuth2ClientCreationInput(helper.exampleOAuth2Client)

	sessionCtxData := &sessioncontext.SessionContextData{
		Requester: sessioncontext.RequesterInfo{
			UserID:                   helper.exampleUser.ID,
			AccountStatus:            helper.exampleUser.AccountStatus,
			AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
			ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
		},
		ActiveHouseholdID: helper.exampleHousehold.ID,
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
			helper.exampleHousehold.ID: authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdMemberRole.String()),
		},
	}

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*sessioncontext.SessionContextData, error) {
		return sessionCtxData, nil
	}
	helper.service.urlClientIDExtractor = func(*http.Request) string {
		return helper.exampleOAuth2Client.ID
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), sessioncontext.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
