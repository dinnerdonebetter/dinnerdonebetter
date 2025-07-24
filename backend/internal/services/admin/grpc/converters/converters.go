package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/services/admin"
)

func ConvertGRPCAdminUpdateUserStatusRequestToUserAccountStatusUpdateInput(input *admin.AdminUpdateUserStatusRequest) *identity.UserAccountStatusUpdateInput {
	return &identity.UserAccountStatusUpdateInput{
		NewStatus:    input.NewStatus,
		Reason:       input.Reason,
		TargetUserID: input.TargetUserID,
	}
}
