package householdinvitations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
)

type householdInvitationsServiceHTTPRoutesTestHelper struct {
	ctx                        context.Context
	req                        *http.Request
	res                        *httptest.ResponseRecorder
	service                    *service
	exampleUser                *types.User
	exampleHousehold           *types.Household
	exampleHouseholdInvitation *types.HouseholdInvitation
	exampleCreationInput       *types.HouseholdInvitationCreationRequestInput
}

func newTestHelper(t *testing.T) *householdInvitationsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &householdInvitationsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleHouseholdInvitation = fakes.BuildFakeHouseholdInvitation()
	helper.exampleCreationInput = converters.ConvertHouseholdInvitationToHouseholdInvitationCreationInput(helper.exampleHouseholdInvitation)

	helper.service.householdIDFetcher = func(*http.Request) string {
		return helper.exampleHousehold.ID
	}
	helper.service.householdInvitationIDFetcher = func(*http.Request) string {
		return helper.exampleHouseholdInvitation.ID
	}

	sessionCtxData := &sessions.ContextData{
		Requester: sessions.RequesterInfo{
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
	helper.service.sessionContextDataFetcher = func(*http.Request) (*sessions.ContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), sessions.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
