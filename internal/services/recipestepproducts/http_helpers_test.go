package recipestepproducts

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

type recipeStepProductsServiceHTTPRoutesTestHelper struct {
	ctx                      context.Context
	req                      *http.Request
	res                      *httptest.ResponseRecorder
	service                  *service
	exampleUser              *types.User
	exampleHousehold         *types.Household
	exampleRecipe            *types.Recipe
	exampleRecipeStep        *types.RecipeStep
	exampleRecipeStepProduct *types.RecipeStepProduct
	exampleCreationInput     *types.RecipeStepProductCreationInput
	exampleUpdateInput       *types.RecipeStepProductUpdateInput
}

func buildTestHelper(t *testing.T) *recipeStepProductsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &recipeStepProductsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleRecipe = fakes.BuildFakeRecipe()
	helper.exampleRecipe.BelongsToHousehold = helper.exampleHousehold.ID
	helper.exampleRecipeStep = fakes.BuildFakeRecipeStep()
	helper.exampleRecipeStep.BelongsToRecipe = helper.exampleRecipe.ID
	helper.exampleRecipeStepProduct = fakes.BuildFakeRecipeStepProduct()
	helper.exampleRecipeStepProduct.BelongsToRecipeStep = helper.exampleRecipeStep.ID
	helper.exampleCreationInput = fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(helper.exampleRecipeStepProduct)
	helper.exampleUpdateInput = fakes.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(helper.exampleRecipeStepProduct)

	helper.service.recipeIDFetcher = func(*http.Request) uint64 {
		return helper.exampleRecipe.ID
	}

	helper.service.recipeStepIDFetcher = func(*http.Request) uint64 {
		return helper.exampleRecipeStep.ID
	}

	helper.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
		return helper.exampleRecipeStepProduct.ID
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
