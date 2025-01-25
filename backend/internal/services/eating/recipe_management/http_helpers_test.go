package recipemanagement

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	encoding "github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
)

type recipesServiceHTTPRoutesTestHelper struct {
	ctx                        context.Context
	req                        *http.Request
	res                        *httptest.ResponseRecorder
	service                    *service
	exampleUser                *types.User
	exampleHousehold           *types.Household
	exampleRecipe              *types.Recipe
	exampleRecipeCreationInput *types.RecipeCreationRequestInput
	exampleRecipeUpdateInput   *types.RecipeUpdateRequestInput

	exampleRecipeStep              *types.RecipeStep
	exampleRecipeStepCreationInput *types.RecipeStepCreationRequestInput
	exampleRecipeStepUpdateInput   *types.RecipeStepUpdateRequestInput

	exampleRecipeStepProduct              *types.RecipeStepProduct
	exampleRecipeStepProductCreationInput *types.RecipeStepProductCreationRequestInput
	exampleRecipeStepProductUpdateInput   *types.RecipeStepProductUpdateRequestInput

	exampleRecipeStepCompletionCondition                               *types.RecipeStepCompletionCondition
	exampleRecipeStepCompletionConditionForExistingRecipeCreationInput *types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput
	exampleRecipeStepCompletionConditionUpdateInput                    *types.RecipeStepCompletionConditionUpdateRequestInput

	exampleRecipePrepTask              *types.RecipePrepTask
	exampleRecipePrepTaskCreationInput *types.RecipePrepTaskCreationRequestInput
	exampleRecipePrepTaskUpdateInput   *types.RecipePrepTaskUpdateRequestInput

	exampleRecipeRating              *types.RecipeRating
	exampleRecipeRatingCreationInput *types.RecipeRatingCreationRequestInput
	exampleRecipeRatingUpdateInput   *types.RecipeRatingUpdateRequestInput

	exampleRecipeStepIngredient              *types.RecipeStepIngredient
	exampleRecipeStepIngredientCreationInput *types.RecipeStepIngredientCreationRequestInput
	exampleRecipeStepIngredientUpdateInput   *types.RecipeStepIngredientUpdateRequestInput

	exampleRecipeStepInstrument              *types.RecipeStepInstrument
	exampleRecipeStepInstrumentCreationInput *types.RecipeStepInstrumentCreationRequestInput
	exampleRecipeStepInstrumentUpdateInput   *types.RecipeStepInstrumentUpdateRequestInput

	exampleRecipeStepVessel              *types.RecipeStepVessel
	exampleRecipeStepVesselCreationInput *types.RecipeStepVesselCreationRequestInput
	exampleRecipeStepVesselUpdateInput   *types.RecipeStepVesselUpdateRequestInput
}

func buildTestHelper(t *testing.T) *recipesServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &recipesServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleRecipe = fakes.BuildFakeRecipe()
	helper.exampleRecipe.CreatedByUser = helper.exampleUser.ID
	helper.exampleRecipeCreationInput = converters.ConvertRecipeToRecipeCreationRequestInput(helper.exampleRecipe)
	helper.exampleRecipeUpdateInput = converters.ConvertRecipeToRecipeUpdateRequestInput(helper.exampleRecipe)

	helper.exampleRecipeStep = fakes.BuildFakeRecipeStep()
	helper.exampleRecipeStep.BelongsToRecipe = helper.exampleRecipe.ID
	helper.exampleRecipeStepCompletionCondition = fakes.BuildFakeRecipeStepCompletionCondition()
	helper.exampleRecipeStepCompletionCondition.BelongsToRecipeStep = helper.exampleRecipeStep.ID
	helper.exampleRecipeStepCompletionConditionForExistingRecipeCreationInput = converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionForExistingRecipeCreationRequestInput(helper.exampleRecipeStepCompletionCondition)
	helper.exampleRecipeStepCompletionConditionUpdateInput = converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionUpdateRequestInput(helper.exampleRecipeStepCompletionCondition)

	helper.exampleRecipePrepTask = fakes.BuildFakeRecipePrepTask()
	helper.exampleRecipePrepTask.BelongsToRecipe = helper.exampleRecipe.ID
	helper.exampleRecipePrepTaskCreationInput = converters.ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(helper.exampleRecipePrepTask)
	helper.exampleRecipePrepTaskUpdateInput = fakes.BuildFakeRecipePrepTaskUpdateRequestInputFromRecipePrepTask(helper.exampleRecipePrepTask)

	helper.exampleRecipeRating = fakes.BuildFakeRecipeRating()
	helper.exampleRecipeRating.RecipeID = helper.exampleRecipe.ID
	helper.exampleRecipeRatingCreationInput = converters.ConvertRecipeRatingToRecipeRatingCreationRequestInput(helper.exampleRecipeRating)
	helper.exampleRecipeRatingUpdateInput = converters.ConvertRecipeRatingToRecipeRatingUpdateRequestInput(helper.exampleRecipeRating)

	helper.exampleRecipeStepIngredient = fakes.BuildFakeRecipeStepIngredient()
	helper.exampleRecipeStepIngredient.BelongsToRecipeStep = helper.exampleRecipeStep.ID
	helper.exampleRecipeStepIngredientCreationInput = converters.ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(helper.exampleRecipeStepIngredient)
	helper.exampleRecipeStepIngredientUpdateInput = converters.ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput(helper.exampleRecipeStepIngredient)

	helper.exampleRecipeStepInstrument = fakes.BuildFakeRecipeStepInstrument()
	helper.exampleRecipeStepInstrument.BelongsToRecipeStep = helper.exampleRecipeStep.ID
	helper.exampleRecipeStepInstrumentCreationInput = converters.ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(helper.exampleRecipeStepInstrument)
	helper.exampleRecipeStepInstrumentUpdateInput = converters.ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput(helper.exampleRecipeStepInstrument)

	helper.exampleRecipeStepProduct = fakes.BuildFakeRecipeStepProduct()
	helper.exampleRecipeStepProduct.BelongsToRecipeStep = helper.exampleRecipeStep.ID
	helper.exampleRecipeStepProductCreationInput = converters.ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(helper.exampleRecipeStepProduct)
	helper.exampleRecipeStepProductUpdateInput = converters.ConvertRecipeStepProductToRecipeStepProductUpdateRequestInput(helper.exampleRecipeStepProduct)

	helper.exampleRecipeStep = fakes.BuildFakeRecipeStep()
	helper.exampleRecipeStep.BelongsToRecipe = helper.exampleRecipe.ID
	helper.exampleRecipeStepCreationInput = converters.ConvertRecipeStepToRecipeStepCreationRequestInput(helper.exampleRecipeStep)
	helper.exampleRecipeStepUpdateInput = converters.ConvertRecipeStepToRecipeStepUpdateRequestInput(helper.exampleRecipeStep)

	helper.exampleRecipeStepVessel = fakes.BuildFakeRecipeStepVessel()
	helper.exampleRecipeStepVessel.BelongsToRecipeStep = helper.exampleRecipeStep.ID
	helper.exampleRecipeStepVesselCreationInput = converters.ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(helper.exampleRecipeStepVessel)
	helper.exampleRecipeStepVesselUpdateInput = converters.ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput(helper.exampleRecipeStepVessel)

	// ID fetchers

	helper.service.recipeIDFetcher = func(*http.Request) string {
		return helper.exampleRecipe.ID
	}

	helper.service.recipeStepIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStep.ID
	}

	helper.service.recipeStepCompletionConditionIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStepCompletionCondition.ID
	}

	helper.service.recipePrepTaskIDFetcher = func(*http.Request) string {
		return helper.exampleRecipePrepTask.ID
	}

	helper.service.recipeRatingIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeRating.ID
	}

	helper.service.recipeStepIngredientIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStepIngredient.ID
	}

	helper.service.recipeStepInstrumentIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStepInstrument.ID
	}

	helper.service.recipeStepProductIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStepProduct.ID
	}

	helper.service.recipeStepVesselIDFetcher = func(*http.Request) string {
		return helper.exampleRecipeStepVessel.ID
	}

	// auth stuff

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
