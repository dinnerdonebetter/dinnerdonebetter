package recipestepcompletionconditions

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

type recipeStepCompletionConditionsServiceHTTPRoutesTestHelper struct {
	ctx                                  context.Context
	req                                  *http.Request
	res                                  *httptest.ResponseRecorder
	service                              *service
	exampleUser                          *types.User
	exampleHousehold                     *types.Household
	exampleRecipe                        *types.Recipe
	exampleRecipeStep                    *types.RecipeStep
	exampleRecipeStepCompletionCondition *types.RecipeStepCompletionCondition
	exampleCreationInput                 *types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput
	exampleUpdateInput                   *types.RecipeStepCompletionConditionUpdateRequestInput
}

func buildTestHelper(t *testing.T) *recipeStepCompletionConditionsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &recipeStepCompletionConditionsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleRecipe = fakes.BuildFakeRecipe()
	helper.exampleRecipe.CreatedByUser = helper.exampleHousehold.ID
	helper.exampleRecipeStep = fakes.BuildFakeRecipeStep()
	helper.exampleRecipeStep.BelongsToRecipe = helper.exampleRecipe.ID
	helper.exampleRecipeStepCompletionCondition = fakes.BuildFakeRecipeStepCompletionCondition()
	helper.exampleRecipeStepCompletionCondition.BelongsToRecipeStep = helper.exampleRecipeStep.ID
	helper.exampleCreationInput = converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionForExistingRecipeCreationRequestInput(helper.exampleRecipeStepCompletionCondition)
	helper.exampleUpdateInput = converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionUpdateRequestInput(helper.exampleRecipeStepCompletionCondition)

	helper.service.recipeIDFetcher = func(*http.Request) string {
		return helper.exampleRecipe.ID
	}

	helper.service.recipeStepIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStep.ID
	}

	helper.service.recipeStepCompletionConditionIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStepCompletionCondition.ID
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
