package validingredients

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

type validIngredientsServiceHTTPRoutesTestHelper struct {
	ctx                    context.Context
	req                    *http.Request
	res                    *httptest.ResponseRecorder
	service                *service
	exampleUser            *types.User
	exampleAccount         *types.Account
	exampleValidIngredient *types.ValidIngredient
	exampleCreationInput   *types.ValidIngredientCreationRequestInput
	exampleUpdateInput     *types.ValidIngredientUpdateRequestInput
}

func buildTestHelper(t *testing.T) *validIngredientsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &validIngredientsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleAccount = fakes.BuildFakeAccount()
	helper.exampleAccount.BelongsToUser = helper.exampleUser.ID
	helper.exampleValidIngredient = fakes.BuildFakeValidIngredient()
	helper.exampleCreationInput = fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(helper.exampleValidIngredient)
	helper.exampleUpdateInput = fakes.BuildFakeValidIngredientUpdateRequestInputFromValidIngredient(helper.exampleValidIngredient)

	helper.service.validIngredientIDFetcher = func(*http.Request) string {
		return helper.exampleValidIngredient.ID
	}

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
			UserID:                helper.exampleUser.ID,
			Reputation:            helper.exampleUser.ServiceAccountStatus,
			ReputationExplanation: helper.exampleUser.ReputationExplanation,
			ServicePermissions:    authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRoles...),
		},
		ActiveAccountID: helper.exampleAccount.ID,
		AccountPermissions: map[string]authorization.AccountRolePermissionsChecker{
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
