package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	encoding "github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/require"
)

type adminServiceHTTPRoutesTestHelper struct {
	ctx                                      context.Context
	service                                  *service
	exampleUser                              *types.User
	exampleHousehold                         *types.Household
	exampleAccountStatusUpdateInput          *types.UserAccountStatusUpdateInput
	exampleArbitraryQueueMessageRequestInput *types.ArbitraryQueueMessageRequestInput

	req *http.Request
	res *httptest.ResponseRecorder
}

func (helper *adminServiceHTTPRoutesTestHelper) neuterAdminUser() {
	helper.exampleUser.ServiceRole = authorization.ServiceUserRole.String()
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return &types.SessionContextData{
			Requester: types.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
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

	helper.ctx = context.Background()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleUser.ServiceRole = authorization.ServiceAdminRole.String()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleAccountStatusUpdateInput = fakes.BuildFakeUserAccountStatusUpdateInput()
	helper.exampleArbitraryQueueMessageRequestInput = fakes.BuildFakeArbitraryQueueMessageRequestInput()

	var err error
	helper.res = httptest.NewRecorder()
	helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://blah.com", http.NoBody)
	require.NoError(t, err)
	require.NotNil(t, helper.req)

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
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
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	return helper
}
