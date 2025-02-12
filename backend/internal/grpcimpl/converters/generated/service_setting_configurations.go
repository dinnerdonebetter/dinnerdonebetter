package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertServiceSettingConfigurationCreationRequestInputToServiceSettingConfiguration(input *messages.ServiceSettingConfigurationCreationRequestInput) *messages.ServiceSettingConfiguration {

output := &messages.ServiceSettingConfiguration{
    Value: input.Value,
    Notes: input.Notes,
    BelongsToUser: input.BelongsToUser,
    BelongsToHousehold: input.BelongsToHousehold,
}

return output
}
func ConvertServiceSettingConfigurationUpdateRequestInputToServiceSettingConfiguration(input *messages.ServiceSettingConfigurationUpdateRequestInput) *messages.ServiceSettingConfiguration {

output := &messages.ServiceSettingConfiguration{
    Notes: input.Notes,
    BelongsToUser: input.BelongsToUser,
    BelongsToHousehold: input.BelongsToHousehold,
    Value: input.Value,
}

return output
}
