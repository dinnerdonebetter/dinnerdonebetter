package validingredientpreparations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
)

type validIngredientPreparationsServiceHTTPRoutesTestHelper struct {
	ctx                               context.Context
	req                               *http.Request
	res                               *httptest.ResponseRecorder
	service                           *service
	exampleUser                       *types.User
	exampleHousehold                  *types.Household
	exampleValidIngredient            *types.ValidIngredient
	exampleValidPreparation           *types.ValidPreparation
	exampleValidIngredientPreparation *types.ValidIngredientPreparation
	exampleCreationInput              *types.ValidIngredientPreparationCreationRequestInput
	exampleUpdateInput                *types.ValidIngredientPreparationUpdateRequestInput
}

func buildTestHelper(t *testing.T) *validIngredientPreparationsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &validIngredientPreparationsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleValidIngredient = fakes.BuildFakeValidIngredient()
	helper.exampleValidPreparation = fakes.BuildFakeValidPreparation()
	helper.exampleValidIngredientPreparation = fakes.BuildFakeValidIngredientPreparation()
	helper.exampleCreationInput = converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(helper.exampleValidIngredientPreparation)
	helper.exampleUpdateInput = converters.ConvertValidIngredientPreparationToValidIngredientPreparationUpdateRequestInput(helper.exampleValidIngredientPreparation)

	helper.service.validIngredientPreparationIDFetcher = func(*http.Request) string {
		return helper.exampleValidIngredientPreparation.ID
	}

	helper.service.validIngredientIDFetcher = func(*http.Request) string {
		return helper.exampleValidIngredient.ID
	}

	helper.service.validPreparationIDFetcher = func(*http.Request) string {
		return helper.exampleValidPreparation.ID
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
