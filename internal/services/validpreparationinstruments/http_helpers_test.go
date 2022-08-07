package validpreparationinstruments

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

type validPreparationInstrumentsServiceHTTPRoutesTestHelper struct {
	ctx                               context.Context
	req                               *http.Request
	res                               *httptest.ResponseRecorder
	service                           *service
	exampleUser                       *types.User
	exampleHousehold                  *types.Household
	exampleValidPreparation           *types.ValidPreparation
	exampleValidInstrument            *types.ValidInstrument
	exampleValidPreparationInstrument *types.ValidPreparationInstrument
	exampleCreationInput              *types.ValidPreparationInstrumentCreationRequestInput
	exampleUpdateInput                *types.ValidPreparationInstrumentUpdateRequestInput
}

func buildTestHelper(t *testing.T) *validPreparationInstrumentsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &validPreparationInstrumentsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleValidPreparation = fakes.BuildFakeValidPreparation()
	helper.exampleValidInstrument = fakes.BuildFakeValidInstrument()
	helper.exampleValidPreparationInstrument = fakes.BuildFakeValidPreparationInstrument()
	helper.exampleCreationInput = fakes.BuildFakeValidPreparationInstrumentCreationRequestInputFromValidPreparationInstrument(helper.exampleValidPreparationInstrument)
	helper.exampleUpdateInput = fakes.BuildFakeValidPreparationInstrumentUpdateRequestInputFromValidPreparationInstrument(helper.exampleValidPreparationInstrument)

	helper.service.validPreparationInstrumentIDFetcher = func(*http.Request) string {
		return helper.exampleValidPreparationInstrument.ID
	}

	helper.service.validPreparationIDFetcher = func(*http.Request) string {
		return helper.exampleValidPreparation.ID
	}

	helper.service.validInstrumentIDFetcher = func(*http.Request) string {
		return helper.exampleValidInstrument.ID
	}

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
			UserID:                   helper.exampleUser.ID,
			AccountStatus:            helper.exampleUser.AccountStatus,
			AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
			ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRoles...),
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
