package households

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
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/require"
)

type householdsServiceHTTPRoutesTestHelper struct {
	ctx              context.Context
	req              *http.Request
	res              *httptest.ResponseRecorder
	service          *service
	exampleUser      *types.User
	exampleHousehold *types.Household
}

func buildTestHelper(t *testing.T) *householdsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &householdsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID

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
	helper.service.householdIDFetcher = func(req *http.Request) string {
		return helper.exampleHousehold.ID
	}
	helper.service.userIDFetcher = func(req *http.Request) string {
		return helper.exampleUser.ID
	}

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
