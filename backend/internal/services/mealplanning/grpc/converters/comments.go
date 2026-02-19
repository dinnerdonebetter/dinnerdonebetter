package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/comments"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
)

// ConvertCommentToGRPCComment converts a domain Comment to proto Comment.
func ConvertCommentToGRPCComment(input *comments.Comment) *mealplanningsvc.Comment {
	if input == nil {
		return nil
	}
	return &mealplanningsvc.Comment{
		Id:              input.ID,
		Content:         input.Content,
		TargetType:      input.TargetType,
		ReferencedId:    input.ReferencedID,
		ParentCommentId: input.ParentCommentID,
		BelongsToUser:   input.BelongsToUser,
		CreatedAt:       grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:   grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
	}
}
