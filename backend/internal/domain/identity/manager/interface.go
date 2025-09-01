package manager

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
)

type (
	IdentityDataManager interface {
		AcceptAccountInvitation(ctx context.Context, accountID, accountInvitationID string, input *identity.AccountInvitationUpdateRequestInput) error
		RejectAccountInvitation(ctx context.Context, accountID, accountInvitationID string, input *identity.AccountInvitationUpdateRequestInput) error
		CancelAccountInvitation(ctx context.Context, accountID, accountInvitationID, note string) error
		ArchiveAccount(ctx context.Context, accountID, ownerID string) error
		ArchiveUserMembership(ctx context.Context, userID, accountID string) error
		ArchiveUser(ctx context.Context, userID string) error
		CreateAccount(ctx context.Context, input *identity.AccountCreationRequestInput) (*identity.Account, error)
		CreateAccountInvitation(ctx context.Context, userID, accountID string, input *identity.AccountInvitationCreationRequestInput) (*identity.AccountInvitation, error)
		CreateUser(ctx context.Context, registrationInput *identity.UserRegistrationInput) (*identity.UserCreationResponse, error)
		GetAccount(ctx context.Context, accountID string) (*identity.Account, error)
		GetAccountInvitation(ctx context.Context, accountID, accountInvitationID string) (*identity.AccountInvitation, error)
		GetAccounts(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.Account, string, error)
		GetReceivedAccountInvitations(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.AccountInvitation, string, error)
		GetSentAccountInvitations(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.AccountInvitation, string, error)
		GetUser(ctx context.Context, userID string) (*identity.User, error)
		GetUsers(ctx context.Context, filter *filtering.QueryFilter) ([]*identity.User, string, error)
		SearchForUsers(ctx context.Context, query string, useSearchService bool, filter *filtering.QueryFilter) ([]*identity.User, string, error)
		SetDefaultAccount(ctx context.Context, userID, accountID string) error
		TransferAccountOwnership(ctx context.Context, accountID string, input *identity.AccountOwnershipTransferInput) error
		UpdateAccount(ctx context.Context, accountID string, input *identity.AccountUpdateRequestInput) error
		UpdateAccountMemberPermissions(ctx context.Context, userID, accountID string, input *identity.ModifyUserPermissionsInput) error
		UpdateUserDetails(ctx context.Context, userID string, input *identity.UserDetailsUpdateRequestInput) error
		UpdateUserEmailAddress(ctx context.Context, userID, newEmail string) error
		UpdateUserUsername(ctx context.Context, userID, newUsername string) error
		UploadUserAvatar(ctx context.Context, userID, base64EncodedImageData string) error
		AdminUpdateUserStatus(ctx context.Context, input *identity.UserAccountStatusUpdateInput) error
	}
)
