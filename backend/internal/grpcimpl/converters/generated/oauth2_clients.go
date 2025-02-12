package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertOAuth2ClientCreationRequestInputToOAuth2Client(input *messages.OAuth2ClientCreationRequestInput) *messages.OAuth2Client {

output := &messages.OAuth2Client{
    Name: input.Name,
    Description: input.Description,
}

return output
}
