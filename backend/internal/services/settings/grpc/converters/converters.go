package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

func ConvertGPRCServiceSettingCreationRequestInputToServiceSettingDatabaseCreationInput(input *settingssvc.ServiceSettingCreationRequestInput) *settings.ServiceSettingDatabaseCreationInput {
	return &settings.ServiceSettingDatabaseCreationInput{
		DefaultValue: &input.DefaultValue,
		Name:         input.Name,
		Type:         input.Type,
		Description:  input.Description,
		Enumeration:  input.Enumeration,
		AdminsOnly:   input.AdminsOnly,
	}
}

func ConvertServiceSettingToGRPCServiceSetting(input *settings.ServiceSetting) *settingssvc.ServiceSetting {
	x := &settingssvc.ServiceSetting{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		ArchivedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		Name:          input.Name,
		ID:            input.ID,
		Type:          input.Type,
		Description:   input.Description,
		Enumeration:   input.Enumeration,
		AdminsOnly:    input.AdminsOnly,
	}

	if input.DefaultValue != nil {
		x.DefaultValue = *input.DefaultValue
	}

	return x
}

func ConvertGRPCServiceSettingConfigurationCreationRequestInputToServiceSettingConfigurationDatabaseCreationInput(input *settingssvc.ServiceSettingConfigurationCreationRequestInput) *settings.ServiceSettingConfigurationDatabaseCreationInput {
	return &settings.ServiceSettingConfigurationDatabaseCreationInput{
		ID:               identifiers.New(),
		Value:            input.Value,
		Notes:            input.Notes,
		ServiceSettingID: input.ServiceSettingID,
		BelongsToUser:    input.BelongsToUser,
		BelongsToAccount: input.BelongsToAccount,
	}
}

func ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(input *settings.ServiceSettingConfiguration) *settingssvc.ServiceSettingConfiguration {
	return &settingssvc.ServiceSettingConfiguration{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		ID:               input.ID,
		Value:            input.Value,
		Notes:            input.Notes,
		BelongsToUser:    input.BelongsToUser,
		BelongsToAccount: input.BelongsToAccount,
		ServiceSetting:   ConvertServiceSettingToGRPCServiceSetting(&input.ServiceSetting),
	}
}

func ConvertGRPCServiceSettingConfigurationUpdateRequestInputToServiceSettingConfigurationUpdateRequestInputTo(input *settingssvc.ServiceSettingConfigurationUpdateRequestInput) *settings.ServiceSettingConfigurationUpdateRequestInput {
	return &settings.ServiceSettingConfigurationUpdateRequestInput{
		Value:            input.Value,
		Notes:            input.Notes,
		ServiceSettingID: input.ServiceSettingID,
		BelongsToUser:    input.BelongsToUser,
		BelongsToAccount: input.BelongsToAccount,
	}
}
