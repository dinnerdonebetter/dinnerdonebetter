package managers

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
)

type (
	DataManager interface {
		AcceptAccountInvitation(ctx context.Context, accountInvitationID string, input *identity.AccountInvitationUpdateRequestInput) error
		ArchiveAccount(ctx context.Context, accountID, ownerID string) error
		ArchiveUserMembership(ctx context.Context, userID, accountID string) error
		ArchiveUser(ctx context.Context, userID string) error
		CancelAccountInvitation(ctx context.Context, accountInvitationID, note string) error
		CreateAccount(ctx context.Context, input *identity.AccountCreationRequestInput) error
		CreateAccountInvitation(ctx context.Context, input *identity.AccountInvitationCreationRequestInput) error
		CreateUser(ctx context.Context, input *identity.UserRegistrationInput) (*identity.User, error)
		GetAccount(ctx context.Context) (*identity.Account, error)
		GetAccountInvitation(ctx context.Context, accountInvitationID string) (*identity.AccountInvitation, error)
		GetAccounts(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.Account, string, error)
		GetReceivedAccountInvitations(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*identity.Account, error)
		GetSentAccountInvitations(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*identity.Account, error)
		GetUser(ctx context.Context, userID string) (*identity.User, error)
		RejectAccountInvitation(ctx context.Context, accountInvitationID, note string) error
		GetUsers(ctx context.Context, filter *filtering.QueryFilter) ([]*identity.User, string, error)
		SearchForUsers(ctx context.Context, filter *filtering.QueryFilter) (*identity.Account, error)
		SetDefaultAccount(ctx context.Context, accountID string) error
		TransferAccountOwnership(ctx context.Context) error
		UpdateAccount(ctx context.Context, input *identity.AccountUpdateRequestInput) error
		UpdateAccountMemberPermissions(ctx context.Context, input *identity.ModifyUserPermissionsInput) error
		UpdateUserDetails(ctx context.Context, input *identity.UserDetailsUpdateRequestInput) error
		UpdateUserEmailAddress(ctx context.Context, newEmail string) error
		UpdateUserUsername(ctx context.Context, newUsername string) error
		UploadUserAvatar(ctx context.Context, newAvatar []byte) error
		AdminUpdateUserStatus(ctx context.Context, input identity.UserAccountStatusUpdateInput) error
	}
)
