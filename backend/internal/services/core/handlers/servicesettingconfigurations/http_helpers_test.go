package servicesettingconfigurations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
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
	exampleHousehold                       *types.Household
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
	helper.exampleHousehold = fakes.BuildFakeHousehold()
	helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID
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
		ActiveHouseholdID: helper.exampleHousehold.ID,
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
			helper.exampleHousehold.ID: authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdMemberRole.String()),
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
