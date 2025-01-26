package authentication

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/require"
)

type authServiceHTTPRoutesTestHelper struct {
	ctx                 context.Context
	req                 *http.Request
	res                 *httptest.ResponseRecorder
	sessionCtxData      *types.SessionContextData
	service             *service
	exampleUser         *types.User
	exampleHousehold    *types.Household
	examplePermCheckers map[string]authorization.HouseholdRolePermissionsChecker
	exampleLoginInput   *types.UserLoginInput
}

func (helper *authServiceHTTPRoutesTestHelper) setContextFetcher(t *testing.T) {
	t.Helper()

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
			UserID:                   helper.exampleUser.ID,
			AccountStatus:            helper.exampleUser.AccountStatus,
			AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
			ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
		},
		ActiveHouseholdID:    helper.exampleHousehold.ID,
		HouseholdPermissions: helper.examplePermCheckers,
	}

	helper.sessionCtxData = sessionCtxData
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}
}

func buildTestHelper(t *testing.T) *authServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &authServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService(t)
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleLoginInput = fakes.BuildFakeUserLoginInputFromUser(helper.exampleUser)

	helper.examplePermCheckers = map[string]authorization.HouseholdRolePermissionsChecker{
		helper.exampleHousehold.ID: authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdMemberRole.String()),
	}

	helper.setContextFetcher(t)

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	var err error

	helper.res = httptest.NewRecorder()
	helper.req, err = http.NewRequestWithContext(
		helper.ctx,
		http.MethodGet,
		"https://whatever.whocares.gov",
		http.NoBody,
	)
	require.NotNil(t, helper.req)
	require.NoError(t, err)

	return helper
}
