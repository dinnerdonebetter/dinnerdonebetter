package grpc

import (
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityfakes "github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	grpcfiltering "github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServiceImpl_ArchiveAccount(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccountID := identityfakes.BuildFakeID()

		identityDataManager.On("ArchiveAccount", testutils.ContextMatcher, exampleAccountID, mock.AnythingOfType("string")).Return(nil)

		request := &identitysvc.ArchiveAccountRequest{
			AccountID: exampleAccountID,
		}

		result, err := service.ArchiveAccount(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.ArchiveAccountRequest{
			AccountID: identityfakes.BuildFakeID(),
		}

		result, err := service.ArchiveAccount(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccountID := identityfakes.BuildFakeID()

		identityDataManager.On("ArchiveAccount", testutils.ContextMatcher, exampleAccountID, mock.AnythingOfType("string")).Return(errors.New("database error"))

		request := &identitysvc.ArchiveAccountRequest{
			AccountID: exampleAccountID,
		}

		result, err := service.ArchiveAccount(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_CreateAccount(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccount := identityfakes.BuildFakeAccount()

		identityDataManager.On("CreateAccount", testutils.ContextMatcher, mock.MatchedBy(func(input *identity.AccountCreationRequestInput) bool {
			return input.Name == exampleAccount.Name &&
				input.BelongsToUser != ""
		})).Return(exampleAccount, nil)

		request := &identitysvc.CreateAccountRequest{
			Input: &identitysvc.AccountCreationRequestInput{
				Name:         exampleAccount.Name,
				ContactPhone: exampleAccount.ContactPhone,
				AddressLine1: exampleAccount.AddressLine1,
				City:         exampleAccount.City,
				State:        exampleAccount.State,
				ZipCode:      exampleAccount.ZipCode,
				Country:      exampleAccount.Country,
			},
		}

		result, err := service.CreateAccount(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
		assert.NotNil(t, result.Created)
		assert.Equal(t, exampleAccount.ID, result.Created.ID)
		assert.Equal(t, exampleAccount.Name, result.Created.Name)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.CreateAccountRequest{
			Input: &identitysvc.AccountCreationRequestInput{
				Name: "Test Account",
			},
		}

		result, err := service.CreateAccount(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		identityDataManager.On("CreateAccount", testutils.ContextMatcher, mock.AnythingOfType("*identity.AccountCreationRequestInput")).Return((*identity.Account)(nil), errors.New("creation error"))

		request := &identitysvc.CreateAccountRequest{
			Input: &identitysvc.AccountCreationRequestInput{
				Name: "Test Account",
			},
		}

		result, err := service.CreateAccount(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_CreateAccountInvitation(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitation := identityfakes.BuildFakeAccountInvitation()

		identityDataManager.On("CreateAccountInvitation", testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.MatchedBy(func(input *identity.AccountInvitationCreationRequestInput) bool {
			return input.ToEmail == exampleInvitation.ToEmail &&
				input.ToName == exampleInvitation.ToName
		})).Return(exampleInvitation, nil)

		request := &identitysvc.CreateAccountInvitationRequest{
			Input: &identitysvc.AccountInvitationCreationRequestInput{
				ToEmail: exampleInvitation.ToEmail,
				ToName:  exampleInvitation.ToName,
				Note:    exampleInvitation.Note,
			},
		}

		result, err := service.CreateAccountInvitation(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
		assert.NotNil(t, result.Created)
		assert.Equal(t, exampleInvitation.ID, result.Created.ID)
		assert.Equal(t, exampleInvitation.ToEmail, result.Created.ToEmail)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.CreateAccountInvitationRequest{
			Input: &identitysvc.AccountInvitationCreationRequestInput{
				ToEmail: "test@example.com",
			},
		}

		result, err := service.CreateAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		identityDataManager.On("CreateAccountInvitation", testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*identity.AccountInvitationCreationRequestInput")).Return((*identity.AccountInvitation)(nil), errors.New("creation error"))

		request := &identitysvc.CreateAccountInvitationRequest{
			Input: &identitysvc.AccountInvitationCreationRequestInput{
				ToEmail: "test@example.com",
			},
		}

		result, err := service.CreateAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_GetAccount(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccount := identityfakes.BuildFakeAccount()

		identityDataManager.On("GetAccount", testutils.ContextMatcher, exampleAccount.ID).Return(exampleAccount, nil)

		request := &identitysvc.GetAccountRequest{
			AccountID: exampleAccount.ID,
		}

		result, err := service.GetAccount(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
		assert.NotNil(t, result.Result)
		assert.Equal(t, exampleAccount.ID, result.Result.ID)
		assert.Equal(t, exampleAccount.Name, result.Result.Name)
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccountID := identityfakes.BuildFakeID()

		identityDataManager.On("GetAccount", testutils.ContextMatcher, exampleAccountID).Return((*identity.Account)(nil), errors.New("database error"))

		request := &identitysvc.GetAccountRequest{
			AccountID: exampleAccountID,
		}

		result, err := service.GetAccount(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_GetAccounts(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccounts := []*identity.Account{
			identityfakes.BuildFakeAccount(),
			identityfakes.BuildFakeAccount(),
		}

		// TODO: wtf
		identityDataManager.On("GetAccounts", testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*filtering.QueryFilter")).Return(exampleAccounts, "", nil)

		pageSize := uint32(25)
		request := &identitysvc.GetAccountsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &pageSize,
			},
		}

		result, err := service.GetAccounts(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
		assert.Len(t, result.Result, len(exampleAccounts))
		assert.Equal(t, exampleAccounts[0].ID, result.Result[0].ID)
		assert.Equal(t, exampleAccounts[1].ID, result.Result[1].ID)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		pageSize := uint32(25)
		request := &identitysvc.GetAccountsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &pageSize,
			},
		}

		result, err := service.GetAccounts(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		identityDataManager.On("GetAccounts", testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*filtering.QueryFilter")).Return(([]*identity.Account)(nil), "", errors.New("database error"))

		pageSize := uint32(25)
		request := &identitysvc.GetAccountsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &pageSize,
			},
		}

		result, err := service.GetAccounts(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_SetDefaultAccount(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccountID := identityfakes.BuildFakeID()

		identityDataManager.On("SetDefaultAccount", testutils.ContextMatcher, mock.AnythingOfType("string"), exampleAccountID).Return(nil)

		request := &identitysvc.SetDefaultAccountRequest{
			AccountID: exampleAccountID,
		}

		result, err := service.SetDefaultAccount(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
		assert.True(t, result.Success)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.SetDefaultAccountRequest{
			AccountID: identityfakes.BuildFakeID(),
		}

		result, err := service.SetDefaultAccount(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccountID := identityfakes.BuildFakeID()

		identityDataManager.On("SetDefaultAccount", testutils.ContextMatcher, mock.AnythingOfType("string"), exampleAccountID).Return(errors.New("update error"))

		request := &identitysvc.SetDefaultAccountRequest{
			AccountID: exampleAccountID,
		}

		result, err := service.SetDefaultAccount(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_TransferAccountOwnership(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccountID := identityfakes.BuildFakeID()

		identityDataManager.On("TransferAccountOwnership", testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*identity.AccountOwnershipTransferInput")).Return(nil)

		request := &identitysvc.TransferAccountOwnershipRequest{
			AccountID: exampleAccountID,
			Input: &identitysvc.AccountOwnershipTransferInput{
				CurrentOwner: identityfakes.BuildFakeID(),
				NewOwner:     identityfakes.BuildFakeID(),
				Reason:       "Transfer for testing",
			},
		}

		result, err := service.TransferAccountOwnership(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
		assert.True(t, result.Success)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.TransferAccountOwnershipRequest{
			AccountID: identityfakes.BuildFakeID(),
			Input: &identitysvc.AccountOwnershipTransferInput{
				CurrentOwner: identityfakes.BuildFakeID(),
				NewOwner:     identityfakes.BuildFakeID(),
			},
		}

		result, err := service.TransferAccountOwnership(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		identityDataManager.On("TransferAccountOwnership", testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*identity.AccountOwnershipTransferInput")).Return(errors.New("transfer error"))

		request := &identitysvc.TransferAccountOwnershipRequest{
			AccountID: identityfakes.BuildFakeID(),
			Input: &identitysvc.AccountOwnershipTransferInput{
				CurrentOwner: identityfakes.BuildFakeID(),
				NewOwner:     identityfakes.BuildFakeID(),
			},
		}

		result, err := service.TransferAccountOwnership(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_UpdateAccount(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleAccountID := identityfakes.BuildFakeID()

		identityDataManager.On("UpdateAccount", testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*identity.AccountUpdateRequestInput")).Return(nil)

		request := &identitysvc.UpdateAccountRequest{
			AccountID: exampleAccountID,
			Input: &identitysvc.AccountUpdateRequestInput{
				Name:         pointer.To("Updated Account Name"),
				ContactPhone: pointer.To("555-0123"),
			},
		}

		result, err := service.UpdateAccount(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.UpdateAccountRequest{
			AccountID: identityfakes.BuildFakeID(),
			Input: &identitysvc.AccountUpdateRequestInput{
				Name: pointer.To("Updated Account Name"),
			},
		}

		result, err := service.UpdateAccount(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		identityDataManager.On("UpdateAccount", testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*identity.AccountUpdateRequestInput")).Return(errors.New("update error"))

		request := &identitysvc.UpdateAccountRequest{
			AccountID: identityfakes.BuildFakeID(),
			Input: &identitysvc.AccountUpdateRequestInput{
				Name: pointer.To("Updated Account Name"),
			},
		}

		result, err := service.UpdateAccount(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_UpdateAccountMemberPermissions(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleUserID := identityfakes.BuildFakeID()

		identityDataManager.On("UpdateAccountMemberPermissions", testutils.ContextMatcher, exampleUserID, mock.AnythingOfType("string"), mock.AnythingOfType("*identity.ModifyUserPermissionsInput")).Return(nil)

		request := &identitysvc.UpdateAccountMemberPermissionsRequest{
			UserID: exampleUserID,
			Input: &identitysvc.ModifyUserPermissionsInput{
				NewRole: "account_admin",
				Reason:  "Promotion for good work",
			},
		}

		result, err := service.UpdateAccountMemberPermissions(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.UpdateAccountMemberPermissionsRequest{
			UserID: identityfakes.BuildFakeID(),
			Input: &identitysvc.ModifyUserPermissionsInput{
				NewRole: "account_admin",
			},
		}

		result, err := service.UpdateAccountMemberPermissions(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		identityDataManager.On("UpdateAccountMemberPermissions", testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*identity.ModifyUserPermissionsInput")).Return(errors.New("update error"))

		request := &identitysvc.UpdateAccountMemberPermissionsRequest{
			UserID: identityfakes.BuildFakeID(),
			Input: &identitysvc.ModifyUserPermissionsInput{
				NewRole: "account_admin",
			},
		}

		result, err := service.UpdateAccountMemberPermissions(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_ArchiveUserMembership(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleUserID := identityfakes.BuildFakeID()
		exampleAccountID := identityfakes.BuildFakeID()

		identityDataManager.On("ArchiveUserMembership", testutils.ContextMatcher, exampleUserID, exampleAccountID).Return(nil)

		request := &identitysvc.ArchiveUserMembershipRequest{
			UserID:    exampleUserID,
			AccountID: exampleAccountID,
		}

		result, err := service.ArchiveUserMembership(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleUserID := identityfakes.BuildFakeID()
		exampleAccountID := identityfakes.BuildFakeID()

		identityDataManager.On("ArchiveUserMembership", testutils.ContextMatcher, exampleUserID, exampleAccountID).Return(errors.New("archive error"))

		request := &identitysvc.ArchiveUserMembershipRequest{
			UserID:    exampleUserID,
			AccountID: exampleAccountID,
		}

		result, err := service.ArchiveUserMembership(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}
