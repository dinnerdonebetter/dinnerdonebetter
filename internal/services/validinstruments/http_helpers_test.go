package validinstruments

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"
)

type validInstrumentsServiceHTTPRoutesTestHelper struct {
	ctx                    context.Context
	req                    *http.Request
	res                    *httptest.ResponseRecorder
	service                *service
	exampleUser            *types.User
	exampleHousehold       *types.Household
	exampleValidInstrument *types.ValidInstrument
	exampleCreationInput   *types.ValidInstrumentCreationInput
	exampleUpdateInput     *types.ValidInstrumentUpdateInput
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
	helper.exampleCreationInput = fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(helper.exampleValidInstrument)
	helper.exampleUpdateInput = fakes.BuildFakeValidInstrumentUpdateInputFromValidInstrument(helper.exampleValidInstrument)

	helper.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
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
		HouseholdPermissions: map[uint64]authorization.HouseholdRolePermissionsChecker{
			helper.exampleHousehold.ID: authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdMemberRole.String()),
		},
	}

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), types.SessionContextDataKey, sessionCtxData))

	helper.res = httptest.NewRecorder()

	return helper
}
