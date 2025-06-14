package users

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
)

type (
	// AdminUserRepository contains administrative User functions.
	AdminUserRepository interface {
		UpdateUserAccountStatus(ctx context.Context, userID string, input *UserAccountStatusUpdateInput) error
	}

	UserRepository interface {
		GetUser(ctx context.Context, userID string) (*User, error)
		GetUserByUsername(ctx context.Context, username string) (*User, error)
		GetAdminUserByUsername(ctx context.Context, username string) (*User, error)
		GetUsers(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[User], error)
		GetUserByEmail(ctx context.Context, email string) (*User, error)
		SearchForUsersByUsername(ctx context.Context, usernameQuery string) ([]*User, error)
		CreateUser(ctx context.Context, input *UserDatabaseCreationInput) (*User, error)
		UpdateUserAvatar(ctx context.Context, userID, newAvatarContent string) error
		UpdateUserUsername(ctx context.Context, userID, newUsername string) error
		UpdateUserEmailAddress(ctx context.Context, userID, newEmailAddress string) error
		UpdateUserDetails(ctx context.Context, userID string, input *UserDetailsDatabaseUpdateInput) error
		UpdateUserPassword(ctx context.Context, userID, newHash string) error
		ArchiveUser(ctx context.Context, userID string) error
		GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID string) (*User, error)
		MarkUserTwoFactorSecretAsVerified(ctx context.Context, userID string) error
		MarkUserTwoFactorSecretAsUnverified(ctx context.Context, userID, newSecret string) error
		GetEmailAddressVerificationTokenForUser(ctx context.Context, userID string) (string, error)
		GetUserByEmailAddressVerificationToken(ctx context.Context, token string) (*User, error)
		MarkUserEmailAddressAsVerified(ctx context.Context, userID, token string) error
		MarkUserEmailAddressAsUnverified(ctx context.Context, userID string) error
		GetUserIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error)
		MarkUserAsIndexed(ctx context.Context, userID string) error
	}
)
