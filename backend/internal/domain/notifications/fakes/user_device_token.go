package fakes

import (
	types "github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
)

// BuildFakeUserDeviceToken builds a faked user device token.
func BuildFakeUserDeviceToken() *types.UserDeviceToken {
	return &types.UserDeviceToken{
		CreatedAt:     BuildFakeTime(),
		ID:            BuildFakeID(),
		DeviceToken:   buildUniqueString() + "_token",
		Platform:      types.UserDeviceTokenPlatformIOS,
		BelongsToUser: BuildFakeID(),
	}
}

// BuildFakeUserDeviceTokenDatabaseCreationInput builds a faked UserDeviceTokenDatabaseCreationInput.
func BuildFakeUserDeviceTokenDatabaseCreationInput() *types.UserDeviceTokenDatabaseCreationInput {
	token := BuildFakeUserDeviceToken()
	return &types.UserDeviceTokenDatabaseCreationInput{
		ID:            token.ID,
		DeviceToken:   token.DeviceToken,
		Platform:      token.Platform,
		BelongsToUser: token.BelongsToUser,
	}
}

// BuildFakeUserDeviceTokensList builds a faked list of user device tokens.
func BuildFakeUserDeviceTokensList() *filtering.QueryFilteredResult[types.UserDeviceToken] {
	var tokens []*types.UserDeviceToken
	for range exampleQuantity {
		tokens = append(tokens, BuildFakeUserDeviceToken())
	}

	return &filtering.QueryFilteredResult[types.UserDeviceToken]{
		Pagination: filtering.Pagination{
			Cursor:          BuildFakeID(),
			MaxResponseSize: 50,
			FilteredCount:   exampleQuantity / 2,
			TotalCount:      exampleQuantity,
		},
		Data: tokens,
	}
}
