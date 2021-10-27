package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

type adminServiceHTTPRoutesTestHelper struct {
	ctx              context.Context
	service          *service
	exampleUser      *types.User
	exampleHousehold *types.Household
	exampleInput     *types.UserReputationUpdateInput

	req *http.Request
	res *httptest.ResponseRecorder
}

func (helper *adminServiceHTTPRoutesTestHelper) neuterAdminUser() {
	helper.exampleUser.ServiceRoles = []string{authorization.ServiceUserRole.String()}
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return &types.SessionContextData{
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
		}, nil
	}
}

func buildTestHelper(t *testing.T) *adminServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &adminServiceHTTPRoutesTestHelper{}

	helper.service = buildTestService(t)

	var err error
	helper.ctx, err = helper.service.sessionManager.Load(context.Background(), "")
	require.NoError(t, err)

	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleUser.ServiceRoles = []string{authorization.ServiceAdminRole.String()}
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleInput = fakes.BuildFakeUserReputationUpdateInput()

	helper.res = httptest.NewRecorder()
	helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://blah.com", nil)
	require.NoError(t, err)
	require.NotNil(t, helper.req)

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

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}
	helper.service.userIDFetcher = func(req *http.Request) string {
		return helper.exampleUser.ID
	}

	return helper
}
