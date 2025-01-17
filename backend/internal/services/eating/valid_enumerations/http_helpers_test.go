package validenumerations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
)

type validEnumerationsServiceHTTPRoutesTestHelper struct {
	_ struct{}

	ctx              context.Context
	req              *http.Request
	res              *httptest.ResponseRecorder
	service          *service
	exampleUser      *types.User
	exampleHousehold *types.Household

	exampleValidIngredient              *types.ValidIngredient
	exampleValidIngredientCreationInput *types.ValidIngredientCreationRequestInput
	exampleValidIngredientUpdateInput   *types.ValidIngredientUpdateRequestInput

	exampleValidIngredientState              *types.ValidIngredientState
	exampleValidIngredientStateCreationInput *types.ValidIngredientStateCreationRequestInput
	exampleValidIngredientStateUpdateInput   *types.ValidIngredientStateUpdateRequestInput

	exampleValidIngredientStateIngredient              *types.ValidIngredientStateIngredient
	exampleValidIngredientStateIngredientCreationInput *types.ValidIngredientStateIngredientCreationRequestInput
	exampleValidIngredientStateIngredientUpdateInput   *types.ValidIngredientStateIngredientUpdateRequestInput

	exampleValidIngredientGroup              *types.ValidIngredientGroup
	exampleValidIngredientGroupCreationInput *types.ValidIngredientGroupCreationRequestInput
	exampleValidIngredientGroupUpdateInput   *types.ValidIngredientGroupUpdateRequestInput

	exampleValidInstrument              *types.ValidInstrument
	exampleValidInstrumentCreationInput *types.ValidInstrumentCreationRequestInput
	exampleValidInstrumentUpdateInput   *types.ValidInstrumentUpdateRequestInput

	exampleValidMeasurementUnit              *types.ValidMeasurementUnit
	exampleValidMeasurementUnitCreationInput *types.ValidMeasurementUnitCreationRequestInput
	exampleValidMeasurementUnitUpdateInput   *types.ValidMeasurementUnitUpdateRequestInput

	exampleValidMeasurementUnitConversion              *types.ValidMeasurementUnitConversion
	exampleValidMeasurementUnitConversionCreationInput *types.ValidMeasurementUnitConversionCreationRequestInput
	exampleValidMeasurementUnitConversionUpdateInput   *types.ValidMeasurementUnitConversionUpdateRequestInput

	exampleValidPreparation              *types.ValidPreparation
	exampleValidPreparationCreationInput *types.ValidPreparationCreationRequestInput
	exampleValidPreparationUpdateInput   *types.ValidPreparationUpdateRequestInput

	exampleValidPreparationInstrument              *types.ValidPreparationInstrument
	exampleValidPreparationInstrumentCreationInput *types.ValidPreparationInstrumentCreationRequestInput
	exampleValidPreparationInstrumentUpdateInput   *types.ValidPreparationInstrumentUpdateRequestInput

	exampleValidVessel              *types.ValidVessel
	exampleValidVesselCreationInput *types.ValidVesselCreationRequestInput
	exampleValidVesselUpdateInput   *types.ValidVesselUpdateRequestInput

	exampleValidPreparationVessel              *types.ValidPreparationVessel
	exampleValidPreparationVesselCreationInput *types.ValidPreparationVesselCreationRequestInput
	exampleValidPreparationVesselUpdateInput   *types.ValidPreparationVesselUpdateRequestInput

	exampleValidIngredientPreparation              *types.ValidIngredientPreparation
	exampleValidIngredientPreparationCreationInput *types.ValidIngredientPreparationCreationRequestInput
	exampleValidIngredientPreparationUpdateInput   *types.ValidIngredientPreparationUpdateRequestInput

	exampleValidIngredientMeasurementUnit              *types.ValidIngredientMeasurementUnit
	exampleValidIngredientMeasurementUnitCreationInput *types.ValidIngredientMeasurementUnitCreationRequestInput
	exampleValidIngredientMeasurementUnitUpdateInput   *types.ValidIngredientMeasurementUnitUpdateRequestInput
}

func buildTestHelper(t *testing.T) *validEnumerationsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &validEnumerationsServiceHTTPRoutesTestHelper{}

	// basleine data

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID

	helper.exampleValidIngredientGroup = fakes.BuildFakeValidIngredientGroup()
	helper.exampleValidIngredientGroupCreationInput = converters.ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(helper.exampleValidIngredientGroup)
	helper.exampleValidIngredientGroupUpdateInput = converters.ConvertValidIngredientGroupToValidIngredientGroupUpdateRequestInput(helper.exampleValidIngredientGroup)

	helper.exampleValidIngredient = fakes.BuildFakeValidIngredient()
	helper.exampleValidIngredientState = fakes.BuildFakeValidIngredientState()
	helper.exampleValidIngredientStateIngredient = fakes.BuildFakeValidIngredientStateIngredient()
	helper.exampleValidIngredientStateIngredientCreationInput = converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(helper.exampleValidIngredientStateIngredient)
	helper.exampleValidIngredientStateIngredientUpdateInput = converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientUpdateRequestInput(helper.exampleValidIngredientStateIngredient)

	helper.exampleValidInstrument = fakes.BuildFakeValidInstrument()
	helper.exampleValidInstrumentCreationInput = converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(helper.exampleValidInstrument)
	helper.exampleValidInstrumentUpdateInput = converters.ConvertValidInstrumentToValidInstrumentUpdateRequestInput(helper.exampleValidInstrument)

	helper.exampleValidMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
	helper.exampleValidMeasurementUnitConversion = fakes.BuildFakeValidMeasurementUnitConversion()
	helper.exampleValidMeasurementUnitConversionCreationInput = converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionCreationRequestInput(helper.exampleValidMeasurementUnitConversion)
	helper.exampleValidMeasurementUnitConversionUpdateInput = converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionUpdateRequestInput(helper.exampleValidMeasurementUnitConversion)

	helper.exampleValidMeasurementUnitCreationInput = converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(helper.exampleValidMeasurementUnit)
	helper.exampleValidMeasurementUnitUpdateInput = converters.ConvertValidMeasurementUnitToValidMeasurementUnitUpdateRequestInput(helper.exampleValidMeasurementUnit)

	helper.exampleValidPreparation = fakes.BuildFakeValidPreparation()
	helper.exampleValidPreparationInstrument = fakes.BuildFakeValidPreparationInstrument()
	helper.exampleValidPreparationInstrumentCreationInput = converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(helper.exampleValidPreparationInstrument)
	helper.exampleValidPreparationInstrumentUpdateInput = converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentUpdateRequestInput(helper.exampleValidPreparationInstrument)

	helper.exampleValidVessel = fakes.BuildFakeValidVessel()
	helper.exampleValidPreparationVessel = fakes.BuildFakeValidPreparationVessel()
	helper.exampleValidPreparationVesselCreationInput = converters.ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput(helper.exampleValidPreparationVessel)
	helper.exampleValidPreparationVesselUpdateInput = converters.ConvertValidPreparationVesselToValidPreparationVesselUpdateRequestInput(helper.exampleValidPreparationVessel)

	helper.exampleValidVesselCreationInput = converters.ConvertValidVesselToValidVesselCreationRequestInput(helper.exampleValidVessel)
	helper.exampleValidVesselUpdateInput = converters.ConvertValidVesselToValidVesselUpdateRequestInput(helper.exampleValidVessel)

	helper.exampleValidIngredientStateCreationInput = converters.ConvertValidIngredientStateToValidIngredientStateCreationRequestInput(helper.exampleValidIngredientState)
	helper.exampleValidIngredientStateUpdateInput = converters.ConvertValidIngredientStateToValidIngredientStateUpdateRequestInput(helper.exampleValidIngredientState)

	helper.exampleValidIngredientCreationInput = converters.ConvertValidIngredientToValidIngredientCreationRequestInput(helper.exampleValidIngredient)
	helper.exampleValidIngredientUpdateInput = converters.ConvertValidIngredientToValidIngredientUpdateRequestInput(helper.exampleValidIngredient)

	helper.exampleValidIngredientPreparation = fakes.BuildFakeValidIngredientPreparation()
	helper.exampleValidIngredientPreparationCreationInput = converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(helper.exampleValidIngredientPreparation)
	helper.exampleValidIngredientPreparationUpdateInput = converters.ConvertValidIngredientPreparationToValidIngredientPreparationUpdateRequestInput(helper.exampleValidIngredientPreparation)

	helper.exampleValidIngredientMeasurementUnit = fakes.BuildFakeValidIngredientMeasurementUnit()
	helper.exampleValidIngredientMeasurementUnitCreationInput = converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(helper.exampleValidIngredientMeasurementUnit)
	helper.exampleValidIngredientMeasurementUnitUpdateInput = converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput(helper.exampleValidIngredientMeasurementUnit)

	// ID fetchers

	helper.service.validIngredientGroupIDFetcher = func(*http.Request) string { return helper.exampleValidIngredientGroup.ID }
	helper.service.validIngredientStateIngredientIDFetcher = func(*http.Request) string { return helper.exampleValidIngredientStateIngredient.ID }
	helper.service.validIngredientIDFetcher = func(*http.Request) string { return helper.exampleValidIngredient.ID }
	helper.service.validIngredientStateIDFetcher = func(*http.Request) string { return helper.exampleValidIngredientState.ID }
	helper.service.validInstrumentIDFetcher = func(*http.Request) string { return helper.exampleValidInstrument.ID }
	helper.service.validMeasurementUnitConversionIDFetcher = func(*http.Request) string { return helper.exampleValidMeasurementUnitConversion.ID }
	helper.service.validMeasurementUnitIDFetcher = func(*http.Request) string { return helper.exampleValidMeasurementUnit.ID }
	helper.service.validPreparationInstrumentIDFetcher = func(*http.Request) string { return helper.exampleValidPreparationInstrument.ID }
	helper.service.validPreparationIDFetcher = func(*http.Request) string { return helper.exampleValidPreparation.ID }
	helper.service.validPreparationVesselIDFetcher = func(*http.Request) string { return helper.exampleValidPreparationVessel.ID }
	helper.service.validVesselIDFetcher = func(*http.Request) string { return helper.exampleValidVessel.ID }
	helper.service.validIngredientPreparationIDFetcher = func(*http.Request) string { return helper.exampleValidIngredientPreparation.ID }
	helper.service.validIngredientMeasurementUnitIDFetcher = func(*http.Request) string { return helper.exampleValidIngredientMeasurementUnit.ID }

	// auth shit

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

	// finishing touches

	req := testutils.BuildTestRequest(t)
	helper.req = req.WithContext(context.WithValue(req.Context(), types.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
