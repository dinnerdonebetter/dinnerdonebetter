package recipesteps

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

type recipeStepsServiceHTTPRoutesTestHelper struct {
	ctx                  context.Context
	req                  *http.Request
	res                  *httptest.ResponseRecorder
	service              *service
	exampleUser          *types.User
	exampleAccount       *types.Account
	exampleRecipe        *types.Recipe
	exampleRecipeStep    *types.RecipeStep
	exampleCreationInput *types.RecipeStepCreationInput
	exampleUpdateInput   *types.RecipeStepUpdateInput
}

func buildTestHelper(t *testing.T) *recipeStepsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &recipeStepsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleAccount = fakes.BuildFakeAccount()
	helper.exampleAccount.BelongsToUser = helper.exampleUser.ID
	helper.exampleRecipe = fakes.BuildFakeRecipe()
	helper.exampleRecipe.BelongsToAccount = helper.exampleAccount.ID
	helper.exampleRecipeStep = fakes.BuildFakeRecipeStep()
	helper.exampleRecipeStep.BelongsToRecipe = helper.exampleRecipe.ID
	helper.exampleCreationInput = fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(helper.exampleRecipeStep)
	helper.exampleUpdateInput = fakes.BuildFakeRecipeStepUpdateInputFromRecipeStep(helper.exampleRecipeStep)

	helper.service.recipeIDFetcher = func(*http.Request) uint64 {
		return helper.exampleRecipe.ID
	}

	helper.service.recipeStepIDFetcher = func(*http.Request) uint64 {
		return helper.exampleRecipeStep.ID
	}

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
			UserID:                helper.exampleUser.ID,
			Reputation:            helper.exampleUser.ServiceAccountStatus,
			ReputationExplanation: helper.exampleUser.ReputationExplanation,
			ServicePermissions:    authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRoles...),
		},
		ActiveAccountID: helper.exampleAccount.ID,
		AccountPermissions: map[uint64]authorization.AccountRolePermissionsChecker{
			helper.exampleAccount.ID: authorization.NewAccountRolePermissionChecker(authorization.AccountMemberRole.String()),
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
