package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/admin"
	adminsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/admin"
)

func ConvertGRPCAdminUpdateUserStatusRequestToUserAccountStatusUpdateInput(input *adminsvc.AdminUpdateUserStatusRequest) *admin.UserAccountStatusUpdateInput {
	return &admin.UserAccountStatusUpdateInput{
		NewStatus:    input.NewStatus,
		Reason:       input.Reason,
		TargetUserID: input.TargetUserID,
	}
}
