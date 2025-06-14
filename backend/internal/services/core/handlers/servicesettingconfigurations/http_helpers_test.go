package servicesettingconfigurations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
)

type serviceSettingConfigurationsServiceHTTPRoutesTestHelper struct {
	ctx                                    context.Context
	req                                    *http.Request
	res                                    *httptest.ResponseRecorder
	service                                *service
	exampleUser                            *types.User
	exampleAccount                         *types.Account
	exampleServiceSettingConfiguration     *types.ServiceSettingConfiguration
	exampleServiceSettingConfigurationList *filtering.QueryFilteredResult[types.ServiceSettingConfiguration]
	exampleCreationInput                   *types.ServiceSettingConfigurationCreationRequestInput
	exampleUpdateInput                     *types.ServiceSettingConfigurationUpdateRequestInput
}

func buildTestHelper(t *testing.T) *serviceSettingConfigurationsServiceHTTPRoutesTestHelper {
	t.Helper()

	helper := &serviceSettingConfigurationsServiceHTTPRoutesTestHelper{}

	helper.ctx = context.Background()
	helper.service = buildTestService()
	helper.exampleUser = fakes.BuildFakeUser()
	helper.exampleAccount = fakes.BuildFakeAccount()
	helper.exampleAccount.BelongsToUser = helper.exampleUser.ID
	helper.exampleServiceSettingConfiguration = fakes.BuildFakeServiceSettingConfiguration()
	helper.exampleServiceSettingConfigurationList = fakes.BuildFakeServiceSettingConfigurationsList()
	helper.exampleCreationInput = converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationCreationRequestInput(helper.exampleServiceSettingConfiguration)
	helper.exampleUpdateInput = converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationUpdateRequestInput(helper.exampleServiceSettingConfiguration)

	helper.service.serviceSettingConfigurationIDFetcher = func(*http.Request) string {
		return helper.exampleServiceSettingConfiguration.ID
	}

	helper.service.serviceSettingNameFetcher = func(*http.Request) string {
		return helper.exampleServiceSettingConfiguration.ServiceSetting.Name
	}

	sessionCtxData := &sessions.ContextData{
		Requester: sessions.RequesterInfo{
			UserID:                   helper.exampleUser.ID,
			AccountStatus:            helper.exampleUser.AccountStatus,
			AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
			ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
		},
		ActiveAccountID: helper.exampleAccount.ID,
		AccountPermissions: map[string]authorization.AccountRolePermissionsChecker{
			helper.exampleAccount.ID: authorization.NewAccountRolePermissionChecker(authorization.AccountMemberRole.String()),
		},
	}

	helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)
	helper.service.sessionContextDataFetcher = func(*http.Request) (*sessions.ContextData, error) {
		return sessionCtxData, nil
	}

	req := testutils.BuildTestRequest(t)

	helper.req = req.WithContext(context.WithValue(req.Context(), sessions.SessionContextDataKey, sessionCtxData))
	helper.res = httptest.NewRecorder()

	return helper
}
