package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeServiceSettingConfiguration builds a faked valid preparation.
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

// BuildFakeServiceSettingConfigurationList builds a faked ServiceSettingConfigurationList.
func BuildFakeServiceSettingConfigurationList() *types.QueryFilteredResult[types.ServiceSettingConfiguration] {
	var examples []*types.ServiceSettingConfiguration
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeServiceSettingConfiguration())
	}

	return &types.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeServiceSettingConfigurationUpdateRequestInput builds a faked ServiceSettingConfigurationUpdateRequestInput from a valid preparation.
func BuildFakeServiceSettingConfigurationUpdateRequestInput() *types.ServiceSettingConfigurationUpdateRequestInput {
	validPreparation := BuildFakeServiceSettingConfiguration()
	return converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationUpdateRequestInput(validPreparation)
}

// BuildFakeServiceSettingConfigurationCreationRequestInput builds a faked ServiceSettingConfigurationCreationRequestInput.
func BuildFakeServiceSettingConfigurationCreationRequestInput() *types.ServiceSettingConfigurationCreationRequestInput {
	validPreparation := BuildFakeServiceSettingConfiguration()
	return converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationCreationRequestInput(validPreparation)
}
