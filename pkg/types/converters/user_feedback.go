package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertUserFeedbackCreationRequestInputToUserFeedbackDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertUserFeedbackCreationRequestInputToUserFeedbackDatabaseCreationInput(x *types.UserFeedbackCreationRequestInput) *types.UserFeedbackDatabaseCreationInput {
	out := &types.UserFeedbackDatabaseCreationInput{
		ID:       identifiers.New(),
		Context:  x.Context,
		Prompt:   x.Prompt,
		Feedback: x.Feedback,
		Rating:   x.Rating,
	}

	return out
}

// ConvertUserFeedbackToUserFeedbackCreationRequestInput builds a UserFeedbackCreationRequestInput from a Ingredient.
func ConvertUserFeedbackToUserFeedbackCreationRequestInput(x *types.UserFeedback) *types.UserFeedbackCreationRequestInput {
	return &types.UserFeedbackCreationRequestInput{
		Context:  x.Context,
		Prompt:   x.Prompt,
		Feedback: x.Feedback,
		Rating:   x.Rating,
	}
}

// ConvertUserFeedbackToUserFeedbackDatabaseCreationInput builds a UserFeedbackDatabaseCreationInput from a UserFeedback.
func ConvertUserFeedbackToUserFeedbackDatabaseCreationInput(x *types.UserFeedback) *types.UserFeedbackDatabaseCreationInput {
	return &types.UserFeedbackDatabaseCreationInput{
		ID:       x.ID,
		Context:  x.Context,
		Prompt:   x.Prompt,
		Feedback: x.Feedback,
		Rating:   x.Rating,
		ByUser:   x.ByUser,
	}
}
