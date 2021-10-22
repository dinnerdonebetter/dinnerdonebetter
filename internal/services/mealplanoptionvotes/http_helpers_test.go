package mealplanoptionvotes

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

type mealPlanOptionVotesServiceHTTPRoutesTestHelper struct {
	ctx                       context.Context
	req                       *http.Request
	res                       *httptest.ResponseRecorder
	service                   *service
	exampleUser               *types.User
	exampleHousehold          *types.Household
	exampleMealPlanOptionVote *types.MealPlanOptionVote
	exampleCreationInput      *types.MealPlanOptionVoteCreationRequestInput
	exampleUpdateInput        *types.MealPlanOptionVoteUpdateRequestInput
}

func buildTestHelper(t *testing.T) *mealPlanOptionVotesServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &mealPlanOptionVotesServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleMealPlanOptionVote = fakes.BuildFakeMealPlanOptionVote()
	helper.exampleMealPlanOptionVote.BelongsToHousehold = helper.exampleHousehold.ID
	helper.exampleCreationInput = fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(helper.exampleMealPlanOptionVote)
	helper.exampleUpdateInput = fakes.BuildFakeMealPlanOptionVoteUpdateRequestInputFromMealPlanOptionVote(helper.exampleMealPlanOptionVote)

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

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), types.SessionContextDataKey, sessionCtxData))

	helper.res = httptest.NewRecorder()

	return helper
}
