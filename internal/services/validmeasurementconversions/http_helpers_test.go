package validmeasurementconversions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prixfixeco/backend/internal/authorization"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	testutils "github.com/prixfixeco/backend/tests/utils"
)

type validMeasurementConversionsServiceHTTPRoutesTestHelper struct {
	ctx                               context.Context
	req                               *http.Request
	res                               *httptest.ResponseRecorder
	service                           *service
	exampleUser                       *types.User
	exampleHousehold                  *types.Household
	exampleValidMeasurementUnit       *types.ValidMeasurementUnit
	exampleValidMeasurementConversion *types.ValidMeasurementUnitConversion
	exampleCreationInput              *types.ValidMeasurementUnitConversionCreationRequestInput
	exampleUpdateInput                *types.ValidMeasurementUnitConversionUpdateRequestInput
}

func buildTestHelper(t *testing.T) *validMeasurementConversionsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &validMeasurementConversionsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleValidMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
	helper.exampleValidMeasurementConversion = fakes.BuildFakeValidMeasurementConversion()
	helper.exampleCreationInput = converters.ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(helper.exampleValidMeasurementConversion)
	helper.exampleUpdateInput = converters.ConvertValidMeasurementConversionToValidMeasurementConversionUpdateRequestInput(helper.exampleValidMeasurementConversion)

	helper.service.validMeasurementConversionIDFetcher = func(*http.Request) string {
		return helper.exampleValidMeasurementConversion.ID
	}

	helper.service.validMeasurementUnitIDFetcher = func(*http.Request) string {
		return helper.exampleValidMeasurementUnit.ID
	}

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

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), types.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
