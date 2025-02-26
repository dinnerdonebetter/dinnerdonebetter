package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeServiceSettingConfiguration builds a faked service setting.
func BuildFakeServiceSettingConfiguration() *types.ServiceSettingConfiguration {
	return &types.ServiceSettingConfiguration{
		ID:                 BuildFakeID(),
		Value:              buildUniqueString(),
		Notes:              buildUniqueString(),
		ServiceSetting:     *BuildFakeServiceSetting(),
		BelongsToUser:      buildUniqueString(),
		BelongsToHousehold: buildUniqueString(),
		CreatedAt:          BuildFakeTime(),
	}
}

// BuildFakeServiceSettingConfigurationsList builds a faked ServiceSettingConfigurationList.
func BuildFakeServiceSettingConfigurationsList() *filtering.QueryFilteredResult[types.ServiceSettingConfiguration] {
	var examples []*types.ServiceSettingConfiguration
	for range exampleQuantity {
		examples = append(examples, BuildFakeServiceSettingConfiguration())
	}

	return &filtering.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeServiceSettingConfigurationUpdateRequestInput builds a faked ServiceSettingConfigurationUpdateRequestInput from a service setting.
func BuildFakeServiceSettingConfigurationUpdateRequestInput() *types.ServiceSettingConfigurationUpdateRequestInput {
	serviceSetting := BuildFakeServiceSettingConfiguration()
	return converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationUpdateRequestInput(serviceSetting)
}

// BuildFakeServiceSettingConfigurationCreationRequestInput builds a faked ServiceSettingConfigurationCreationRequestInput.
func BuildFakeServiceSettingConfigurationCreationRequestInput() *types.ServiceSettingConfigurationCreationRequestInput {
	serviceSetting := BuildFakeServiceSettingConfiguration()
	return converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationCreationRequestInput(serviceSetting)
}
