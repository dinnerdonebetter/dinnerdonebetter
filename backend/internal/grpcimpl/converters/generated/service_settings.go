package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertServiceSettingCreationRequestInputToServiceSetting(input *messages.ServiceSettingCreationRequestInput) *messages.ServiceSetting {

output := &messages.ServiceSetting{
    Name: input.Name,
    DefaultValue: input.DefaultValue,
    Type: input.Type,
    Description: input.Description,
    Enumeration: input.Enumeration,
    AdminsOnly: input.AdminsOnly,
}

return output
}
