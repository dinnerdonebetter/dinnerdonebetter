package validinstruments

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

type validInstrumentsServiceHTTPRoutesTestHelper struct {
	ctx                    context.Context
	req                    *http.Request
	res                    *httptest.ResponseRecorder
	service                *service
	exampleUser            *types.User
	exampleHousehold       *types.Household
	exampleValidInstrument *types.ValidInstrument
	exampleCreationInput   *types.ValidInstrumentCreationRequestInput
	exampleUpdateInput     *types.ValidInstrumentUpdateRequestInput
}

func buildTestHelper(t *testing.T) *validInstrumentsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &validInstrumentsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleValidInstrument = fakes.BuildFakeValidInstrument()
	helper.exampleCreationInput = fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(helper.exampleValidInstrument)
	helper.exampleUpdateInput = fakes.BuildFakeValidInstrumentUpdateRequestInputFromValidInstrument(helper.exampleValidInstrument)

	helper.service.validInstrumentIDFetcher = func(*http.Request) string {
		return helper.exampleValidInstrument.ID
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

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), types.SessionContextDataKey, sessionCtxData))

	helper.res = httptest.NewRecorder()

	return helper
}
