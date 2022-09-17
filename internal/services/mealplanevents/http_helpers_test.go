package mealplanevents

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

type mealPlanEventsServiceHTTPRoutesTestHelper struct {
	ctx                  context.Context
	req                  *http.Request
	res                  *httptest.ResponseRecorder
	service              *service
	exampleUser          *types.User
	exampleHousehold     *types.Household
	exampleMealPlan      *types.MealPlan
	exampleMealPlanEvent *types.MealPlanEvent
	exampleCreationInput *types.MealPlanEventCreationRequestInput
	exampleUpdateInput   *types.MealPlanEventUpdateRequestInput
}

func buildTestHelper(t *testing.T) *mealPlanEventsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &mealPlanEventsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleMealPlan = fakes.BuildFakeMealPlan()
	helper.exampleMealPlanEvent = fakes.BuildFakeMealPlanEvent()
	helper.exampleMealPlanEvent.BelongsToMealPlan = helper.exampleMealPlan.ID
	helper.exampleCreationInput = fakes.BuildFakeMealPlanEventCreationRequestInputFromMealPlanEvent(helper.exampleMealPlanEvent)
	helper.exampleUpdateInput = fakes.BuildFakeMealPlanEventUpdateRequestInputFromMealPlanEvent(helper.exampleMealPlanEvent)

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanEventIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanEvent.ID
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
