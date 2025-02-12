package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertPasswordResetTokenCreationRequestInputToPasswordResetToken(input *messages.PasswordResetTokenCreationRequestInput) *messages.PasswordResetToken {

output := &messages.PasswordResetToken{
}

return output
}
