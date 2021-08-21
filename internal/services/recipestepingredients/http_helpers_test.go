package recipestepingredients

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

type recipeStepIngredientsServiceHTTPRoutesTestHelper struct {
	ctx                         context.Context
	req                         *http.Request
	res                         *httptest.ResponseRecorder
	service                     *service
	exampleUser                 *types.User
	exampleHousehold            *types.Household
	exampleRecipe               *types.Recipe
	exampleRecipeStep           *types.RecipeStep
	exampleRecipeStepIngredient *types.RecipeStepIngredient
	exampleCreationInput        *types.RecipeStepIngredientCreationInput
	exampleUpdateInput          *types.RecipeStepIngredientUpdateInput
}

func buildTestHelper(t *testing.T) *recipeStepIngredientsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &recipeStepIngredientsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleRecipe = fakes.BuildFakeRecipe()
	helper.exampleRecipe.BelongsToHousehold = helper.exampleHousehold.ID
	helper.exampleRecipeStep = fakes.BuildFakeRecipeStep()
	helper.exampleRecipeStep.BelongsToRecipe = helper.exampleRecipe.ID
	helper.exampleRecipeStepIngredient = fakes.BuildFakeRecipeStepIngredient()
	helper.exampleRecipeStepIngredient.BelongsToRecipeStep = helper.exampleRecipeStep.ID
	helper.exampleCreationInput = fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(helper.exampleRecipeStepIngredient)
	helper.exampleUpdateInput = fakes.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(helper.exampleRecipeStepIngredient)

	helper.service.recipeIDFetcher = func(*http.Request) uint64 {
		return helper.exampleRecipe.ID
	}

	helper.service.recipeStepIDFetcher = func(*http.Request) uint64 {
		return helper.exampleRecipeStep.ID
	}

	helper.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
		return helper.exampleRecipeStepIngredient.ID
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
