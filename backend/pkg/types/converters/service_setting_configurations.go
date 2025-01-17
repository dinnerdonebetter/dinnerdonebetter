package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertServiceSettingConfigurationToServiceSettingConfigurationUpdateRequestInput creates a ServiceSettingConfigurationUpdateRequestInput from a ServiceSettingConfiguration.
func ConvertServiceSettingConfigurationToServiceSettingConfigurationUpdateRequestInput(input *types.ServiceSettingConfiguration) *types.ServiceSettingConfigurationUpdateRequestInput {
	x := &types.ServiceSettingConfigurationUpdateRequestInput{
		Value:              &input.Value,
		Notes:              &input.Notes,
		ServiceSettingID:   &input.ServiceSetting.ID,
		BelongsToUser:      &input.BelongsToUser,
		BelongsToHousehold: &input.BelongsToHousehold,
	}

	return x
}

// ConvertServiceSettingConfigurationCreationRequestInputToServiceSettingConfigurationDatabaseCreationInput creates a ServiceSettingConfigurationDatabaseCreationInput from a ServiceSettingConfigurationCreationRequestInput.
func ConvertServiceSettingConfigurationCreationRequestInputToServiceSettingConfigurationDatabaseCreationInput(input *types.ServiceSettingConfigurationCreationRequestInput) *types.ServiceSettingConfigurationDatabaseCreationInput {
	x := &types.ServiceSettingConfigurationDatabaseCreationInput{
		ID:                 identifiers.New(),
		Value:              input.Value,
		Notes:              input.Notes,
		ServiceSettingID:   input.ServiceSettingID,
		BelongsToUser:      input.BelongsToUser,
		BelongsToHousehold: input.BelongsToHousehold,
	}

	return x
}

// ConvertServiceSettingConfigurationToServiceSettingConfigurationCreationRequestInput builds a ServiceSettingConfigurationCreationRequestInput from a ServiceSettingConfiguration.
func ConvertServiceSettingConfigurationToServiceSettingConfigurationCreationRequestInput(input *types.ServiceSettingConfiguration) *types.ServiceSettingConfigurationCreationRequestInput {
	return &types.ServiceSettingConfigurationCreationRequestInput{
		Value:              input.Value,
		Notes:              input.Notes,
		ServiceSettingID:   input.ServiceSetting.ID,
		BelongsToUser:      input.BelongsToUser,
		BelongsToHousehold: input.BelongsToHousehold,
	}
}

// ConvertServiceSettingConfigurationToServiceSettingConfigurationDatabaseCreationInput builds a ServiceSettingConfigurationDatabaseCreationInput from a ServiceSettingConfiguration.
func ConvertServiceSettingConfigurationToServiceSettingConfigurationDatabaseCreationInput(input *types.ServiceSettingConfiguration) *types.ServiceSettingConfigurationDatabaseCreationInput {
	return &types.ServiceSettingConfigurationDatabaseCreationInput{
		ID:                 input.ID,
		Value:              input.Value,
		Notes:              input.Notes,
		ServiceSettingID:   input.ServiceSetting.ID,
		BelongsToUser:      input.BelongsToUser,
		BelongsToHousehold: input.BelongsToHousehold,
	}
}
