package users

import (
	"context"
	"database/sql"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	"github.com/dinnerdonebetter/backend/internal/domain/users"
)

type UserRepository struct {
	db *sql.DB
}

func (u *UserRepository) GetUser(ctx context.Context, userID string) (*users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetUserByUsername(ctx context.Context, username string) (*users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetAdminUserByUsername(ctx context.Context, username string) (*users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetUsers(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[users.User], error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) SearchForUsersByUsername(ctx context.Context, usernameQuery string) ([]*users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) CreateUser(ctx context.Context, input *users.UserDatabaseCreationInput) (*users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) UpdateUserAvatar(ctx context.Context, userID, newAvatarContent string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) UpdateUserUsername(ctx context.Context, userID, newUsername string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) UpdateUserEmailAddress(ctx context.Context, userID, newEmailAddress string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) UpdateUserDetails(ctx context.Context, userID string, input *users.UserDetailsDatabaseUpdateInput) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) UpdateUserPassword(ctx context.Context, userID, newHash string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) ArchiveUser(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID string) (*users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) MarkUserTwoFactorSecretAsVerified(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) MarkUserTwoFactorSecretAsUnverified(ctx context.Context, userID, newSecret string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetEmailAddressVerificationTokenForUser(ctx context.Context, userID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetUserByEmailAddressVerificationToken(ctx context.Context, token string) (*users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) MarkUserEmailAddressAsVerified(ctx context.Context, userID, token string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) MarkUserEmailAddressAsUnverified(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetUserIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) MarkUserAsIndexed(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(db *sql.DB) users.UserRepository {
	return UserRepository{db: db}
}
