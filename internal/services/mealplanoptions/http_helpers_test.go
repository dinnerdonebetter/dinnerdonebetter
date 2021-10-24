package mealplanoptions

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

type mealPlanOptionsServiceHTTPRoutesTestHelper struct {
	ctx                   context.Context
	req                   *http.Request
	res                   *httptest.ResponseRecorder
	service               *service
	exampleUser           *types.User
	exampleHousehold      *types.Household
	exampleMealPlan       *types.MealPlan
	exampleMealPlanOption *types.MealPlanOption
	exampleCreationInput  *types.MealPlanOptionCreationRequestInput
	exampleUpdateInput    *types.MealPlanOptionUpdateRequestInput
}

func buildTestHelper(t *testing.T) *mealPlanOptionsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &mealPlanOptionsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleMealPlan = fakes.BuildFakeMealPlan()
	helper.exampleMealPlan.BelongsToHousehold = helper.exampleHousehold.ID
	helper.exampleMealPlanOption = fakes.BuildFakeMealPlanOption()
	helper.exampleMealPlanOption.BelongsToMealPlan = helper.exampleMealPlan.ID
	helper.exampleCreationInput = fakes.BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(helper.exampleMealPlanOption)
	helper.exampleUpdateInput = fakes.BuildFakeMealPlanOptionUpdateRequestInputFromMealPlanOption(helper.exampleMealPlanOption)

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanOptionIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanOption.ID
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

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), types.SessionContextDataKey, sessionCtxData))

	helper.res = httptest.NewRecorder()

	return helper
}
