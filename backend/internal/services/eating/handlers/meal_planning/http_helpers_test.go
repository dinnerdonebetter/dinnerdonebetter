package mealplanning

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessioncontext"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
)

type mealsServiceHTTPRoutesTestHelper struct {
	ctx                                              context.Context
	exampleMealPlanEventUpdateInput                  *types.MealPlanEventUpdateRequestInput
	exampleHousehold                                 *types.Household
	exampleMealPlanTask                              *types.MealPlanTask
	exampleUser                                      *types.User
	exampleMealPlanGroceryListItem                   *types.MealPlanGroceryListItem
	exampleMeal                                      *types.Meal
	exampleCreationInput                             *types.MealCreationRequestInput
	exampleUpdateInput                               *types.MealUpdateRequestInput
	exampleMealPlan                                  *types.MealPlan
	exampleMealPlanCreationInput                     *types.MealPlanCreationRequestInput
	exampleMealPlanUpdateInput                       *types.MealPlanUpdateRequestInput
	exampleMealPlanOption                            *types.MealPlanOption
	exampleMealPlanEventCreationInput                *types.MealPlanEventCreationRequestInput
	req                                              *http.Request
	service                                          *service
	res                                              *httptest.ResponseRecorder
	exampleMealPlanEvent                             *types.MealPlanEvent
	exampleMealPlanOptionCreationInput               *types.MealPlanOptionCreationRequestInput
	exampleMealPlanOptionUpdateInput                 *types.MealPlanOptionUpdateRequestInput
	exampleMealPlanOptionVote                        *types.MealPlanOptionVote
	exampleMealPlanOptionVoteCreationInput           *types.MealPlanOptionVoteCreationRequestInput
	exampleMealPlanOptionVoteUpdateInput             *types.MealPlanOptionVoteUpdateRequestInput
	exampleHouseholdInstrumentOwnershipUpdateInput   *types.HouseholdInstrumentOwnershipUpdateRequestInput
	exampleUserIngredientPreference                  *types.UserIngredientPreference
	exampleUserIngredientPreferenceCreationInput     *types.UserIngredientPreferenceCreationRequestInput
	exampleUserIngredientPreferenceUpdateInput       *types.UserIngredientPreferenceUpdateRequestInput
	exampleHouseholdInstrumentOwnership              *types.HouseholdInstrumentOwnership
	exampleHouseholdInstrumentOwnershipCreationInput *types.HouseholdInstrumentOwnershipCreationRequestInput
	exampleMealPlanOptionVotes                       []*types.MealPlanOptionVote
}

func buildTestHelper(t *testing.T) *mealsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &mealsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
	helper.exampleMeal = fakes.BuildFakeMeal()
	helper.exampleMeal.CreatedByUser = helper.exampleHousehold.ID
	helper.exampleCreationInput = converters.ConvertMealToMealCreationRequestInput(helper.exampleMeal)
	helper.exampleUpdateInput = converters.ConvertMealToMealUpdateRequestInput(helper.exampleMeal)

	helper.exampleMealPlan = fakes.BuildFakeMealPlan()
	helper.exampleMealPlan.BelongsToHousehold = helper.exampleHousehold.ID
	helper.exampleMealPlanCreationInput = converters.ConvertMealPlanToMealPlanCreationRequestInput(helper.exampleMealPlan)
	helper.exampleMealPlanUpdateInput = converters.ConvertMealPlanToMealPlanUpdateRequestInput(helper.exampleMealPlan)

	helper.exampleMealPlanEvent = fakes.BuildFakeMealPlanEvent()
	helper.exampleMealPlanEvent.BelongsToMealPlan = helper.exampleMealPlan.ID
	helper.exampleMealPlanEventCreationInput = converters.ConvertMealPlanEventToMealPlanEventCreationRequestInput(helper.exampleMealPlanEvent)
	helper.exampleMealPlanEventUpdateInput = converters.ConvertMealPlanEventToMealPlanEventUpdateRequestInput(helper.exampleMealPlanEvent)

	helper.exampleMealPlanGroceryListItem = fakes.BuildFakeMealPlanGroceryListItem()

	helper.exampleMealPlanOption = fakes.BuildFakeMealPlanOption()
	helper.exampleMealPlanOption.BelongsToMealPlanEvent = helper.exampleMealPlanEvent.ID
	helper.exampleMealPlanOptionCreationInput = converters.ConvertMealPlanOptionToMealPlanOptionCreationRequestInput(helper.exampleMealPlanOption)
	helper.exampleMealPlanOptionUpdateInput = converters.ConvertMealPlanOptionToMealPlanOptionUpdateRequestInput(helper.exampleMealPlanOption)

	helper.exampleMealPlanOptionVote = fakes.BuildFakeMealPlanOptionVote()
	helper.exampleMealPlanOptionVote.BelongsToMealPlanOption = helper.exampleMealPlanOption.ID
	helper.exampleMealPlanOptionVotes = []*types.MealPlanOptionVote{helper.exampleMealPlanOptionVote}
	helper.exampleMealPlanOptionVoteCreationInput = converters.ConvertMealPlanOptionVoteToMealPlanOptionVoteCreationRequestInput(helper.exampleMealPlanOptionVote)
	helper.exampleMealPlanOptionVoteUpdateInput = converters.ConvertMealPlanOptionVoteToMealPlanOptionVoteUpdateRequestInput(helper.exampleMealPlanOptionVote)

	helper.exampleMealPlanTask = fakes.BuildFakeMealPlanTask()

	helper.exampleUserIngredientPreference = fakes.BuildFakeUserIngredientPreference()
	helper.exampleUserIngredientPreferenceCreationInput = converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(helper.exampleUserIngredientPreference)
	helper.exampleUserIngredientPreferenceUpdateInput = converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceUpdateRequestInput(helper.exampleUserIngredientPreference)

	helper.exampleHouseholdInstrumentOwnership = fakes.BuildFakeHouseholdInstrumentOwnership()
	helper.exampleHouseholdInstrumentOwnershipCreationInput = converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipCreationRequestInput(helper.exampleHouseholdInstrumentOwnership)
	helper.exampleHouseholdInstrumentOwnershipUpdateInput = converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipUpdateRequestInput(helper.exampleHouseholdInstrumentOwnership)

	// ID fetchers

	helper.service.mealIDFetcher = func(*http.Request) string {
		return helper.exampleMeal.ID
	}

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanEventIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanEvent.ID
	}

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanGroceryListItemIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanGroceryListItem.ID
	}

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanEventIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanEvent.ID
	}

	helper.service.mealPlanOptionIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanOption.ID
	}

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanEventIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanEvent.ID
	}

	helper.service.mealPlanOptionIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanOption.ID
	}

	helper.service.mealPlanOptionVoteIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanOptionVote.ID
	}

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlan.ID
	}

	helper.service.mealPlanTaskIDFetcher = func(*http.Request) string {
		return helper.exampleMealPlanTask.ID
	}

	helper.service.userIngredientPreferenceIDFetcher = func(*http.Request) string {
		return helper.exampleUserIngredientPreference.ID
	}

	helper.service.householdInstrumentOwnershipIDFetcher = func(*http.Request) string {
		return helper.exampleHouseholdInstrumentOwnership.ID
	}

	// auth stuff

	sessionCtxData := &sessioncontext.SessionContextData{
		Requester: sessioncontext.RequesterInfo{
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
	helper.service.sessionContextDataFetcher = func(*http.Request) (*sessioncontext.SessionContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), sessioncontext.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
