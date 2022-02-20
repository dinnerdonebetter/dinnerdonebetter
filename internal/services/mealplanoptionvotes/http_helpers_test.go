package mealplanoptionvotes

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

type mealPlanOptionVotesServiceHTTPRoutesTestHelper struct {
	ctx                        context.Context
	req                        *http.Request
	res                        *httptest.ResponseRecorder
	service                    *service
	exampleUser                *types.User
	exampleHousehold           *types.Household
	exampleMealPlan            *types.MealPlan
	exampleMealPlanOption      *types.MealPlanOption
	exampleMealPlanOptionVote  *types.MealPlanOptionVote
	exampleCreationInput       *types.MealPlanOptionVoteCreationRequestInput
	exampleUpdateInput         *types.MealPlanOptionVoteUpdateRequestInput
	exampleMealPlanOptionVotes []*types.MealPlanOptionVote
}

func buildTestHelper(t *testing.T) *mealPlanOptionVotesServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &mealPlanOptionVotesServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleMealPlan = fakes.BuildFakeMealPlan()
	helper.exampleMealPlan.BelongsToHousehold = helper.exampleHousehold.ID
	helper.exampleMealPlanOption = fakes.BuildFakeMealPlanOption()
	helper.exampleMealPlanOption.BelongsToMealPlan = helper.exampleMealPlan.ID
	helper.exampleMealPlanOptionVote = fakes.BuildFakeMealPlanOptionVote()
	helper.exampleMealPlanOptionVote.BelongsToMealPlanOption = helper.exampleMealPlanOption.ID
	helper.exampleMealPlanOptionVotes = []*types.MealPlanOptionVote{helper.exampleMealPlanOptionVote}
	helper.exampleCreationInput = fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(helper.exampleMealPlanOptionVote)
	helper.exampleUpdateInput = fakes.BuildFakeMealPlanOptionVoteUpdateRequestInputFromMealPlanOptionVote(helper.exampleMealPlanOptionVote)

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanOptionIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanOption.ID
	}

	helper.service.mealPlanOptionVoteIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanOptionVote.ID
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
