package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertServiceSettingToServiceSettingUpdateRequestInput creates a ServiceSettingUpdateRequestInput from a ServiceSetting.
func ConvertServiceSettingToServiceSettingUpdateRequestInput(input *types.ServiceSetting) *types.ServiceSettingUpdateRequestInput {
	x := &types.ServiceSettingUpdateRequestInput{
		Name:         &input.Name,
		Type:         &input.Type,
		Description:  &input.Description,
		DefaultValue: input.DefaultValue,
		AdminsOnly:   &input.AdminsOnly,
	}

	return x
}

// ConvertServiceSettingCreationRequestInputToServiceSettingDatabaseCreationInput creates a ServiceSettingDatabaseCreationInput from a ServiceSettingCreationRequestInput.
func ConvertServiceSettingCreationRequestInputToServiceSettingDatabaseCreationInput(input *types.ServiceSettingCreationRequestInput) *types.ServiceSettingDatabaseCreationInput {
	x := &types.ServiceSettingDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         input.Name,
		Type:         input.Type,
		Description:  input.Description,
		DefaultValue: input.DefaultValue,
		AdminsOnly:   input.AdminsOnly,
	}

	return x
}

// ConvertServiceSettingToServiceSettingCreationRequestInput builds a ServiceSettingCreationRequestInput from a ServiceSetting.
func ConvertServiceSettingToServiceSettingCreationRequestInput(input *types.ServiceSetting) *types.ServiceSettingCreationRequestInput {
	return &types.ServiceSettingCreationRequestInput{
		Name:         input.Name,
		Type:         input.Type,
		Description:  input.Description,
		DefaultValue: input.DefaultValue,
		AdminsOnly:   input.AdminsOnly,
	}
}

// ConvertServiceSettingToServiceSettingDatabaseCreationInput builds a ServiceSettingDatabaseCreationInput from a ServiceSetting.
func ConvertServiceSettingToServiceSettingDatabaseCreationInput(input *types.ServiceSetting) *types.ServiceSettingDatabaseCreationInput {
	return &types.ServiceSettingDatabaseCreationInput{
		Name:         input.Name,
		Type:         input.Type,
		Description:  input.Description,
		DefaultValue: input.DefaultValue,
		AdminsOnly:   input.AdminsOnly,
	}
}