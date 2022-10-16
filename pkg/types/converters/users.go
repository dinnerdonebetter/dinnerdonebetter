package converters

import (
	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertUserToUserCreationResponse builds a UserCreationResponse from a user.
func ConvertUserToUserCreationResponse(user *types.User) *types.UserCreationResponse {
	return &types.UserCreationResponse{
		CreatedUserID: user.ID,
		Username:      user.Username,
		CreatedAt:     user.CreatedAt,
	}
}

// ConvertUserToUserDatabaseCreationInput builds a UserDatabaseCreationInput from a User.
func ConvertUserToUserDatabaseCreationInput(user *types.User) *types.UserDatabaseCreationInput {
	return &types.UserDatabaseCreationInput{
		ID:              user.ID,
		EmailAddress:    user.EmailAddress,
		Username:        user.Username,
		HashedPassword:  user.HashedPassword,
		TwoFactorSecret: user.TwoFactorSecret,
		BirthDay:        user.BirthDay,
		BirthMonth:      user.BirthMonth,
	}
}
