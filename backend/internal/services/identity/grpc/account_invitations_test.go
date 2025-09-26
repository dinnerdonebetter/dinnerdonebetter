package grpc

import (
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityfakes "github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	grpcfiltering "github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServiceImpl_AcceptAccountInvitation(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitationID := identityfakes.BuildFakeID()

		identityDataManager.EXPECT().AcceptAccountInvitation(testutils.ContextMatcher, mock.AnythingOfType("string"), exampleInvitationID, mock.AnythingOfType("*identity.AccountInvitationUpdateRequestInput")).Return(nil)

		request := &identitysvc.AcceptAccountInvitationRequest{
			AccountInvitationID: exampleInvitationID,
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Token: "invitation-token",
				Note:  "Accepting invitation",
			},
		}

		result, err := service.AcceptAccountInvitation(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.AcceptAccountInvitationRequest{
			AccountInvitationID: identityfakes.BuildFakeID(),
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Token: "invitation-token",
			},
		}

		result, err := service.AcceptAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitationID := identityfakes.BuildFakeID()

		identityDataManager.EXPECT().AcceptAccountInvitation(testutils.ContextMatcher, mock.AnythingOfType("string"), exampleInvitationID, mock.AnythingOfType("*identity.AccountInvitationUpdateRequestInput")).Return(errors.New("accept error"))

		request := &identitysvc.AcceptAccountInvitationRequest{
			AccountInvitationID: exampleInvitationID,
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Token: "invitation-token",
			},
		}

		result, err := service.AcceptAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_RejectAccountInvitation(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitationID := identityfakes.BuildFakeID()

		identityDataManager.EXPECT().RejectAccountInvitation(testutils.ContextMatcher, mock.AnythingOfType("string"), exampleInvitationID, mock.AnythingOfType("*identity.AccountInvitationUpdateRequestInput")).Return(nil)

		request := &identitysvc.RejectAccountInvitationRequest{
			AccountInvitationID: exampleInvitationID,
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Token: "invitation-token",
				Note:  "Rejecting invitation",
			},
		}

		result, err := service.RejectAccountInvitation(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.RejectAccountInvitationRequest{
			AccountInvitationID: identityfakes.BuildFakeID(),
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Token: "invitation-token",
			},
		}

		result, err := service.RejectAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitationID := identityfakes.BuildFakeID()

		identityDataManager.EXPECT().RejectAccountInvitation(testutils.ContextMatcher, mock.AnythingOfType("string"), exampleInvitationID, mock.AnythingOfType("*identity.AccountInvitationUpdateRequestInput")).Return(errors.New("reject error"))

		request := &identitysvc.RejectAccountInvitationRequest{
			AccountInvitationID: exampleInvitationID,
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Token: "invitation-token",
			},
		}

		result, err := service.RejectAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_CancelAccountInvitation(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitationID := identityfakes.BuildFakeID()

		identityDataManager.EXPECT().CancelAccountInvitation(testutils.ContextMatcher, mock.AnythingOfType("string"), exampleInvitationID, "Cancelling invitation").Return(nil)

		request := &identitysvc.CancelAccountInvitationRequest{
			AccountInvitationID: exampleInvitationID,
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Note: "Cancelling invitation",
			},
		}

		result, err := service.CancelAccountInvitation(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.CancelAccountInvitationRequest{
			AccountInvitationID: identityfakes.BuildFakeID(),
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Note: "Cancelling invitation",
			},
		}

		result, err := service.CancelAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitationID := identityfakes.BuildFakeID()

		identityDataManager.EXPECT().CancelAccountInvitation(testutils.ContextMatcher, mock.AnythingOfType("string"), exampleInvitationID, "Cancelling invitation").Return(errors.New("cancel error"))

		request := &identitysvc.CancelAccountInvitationRequest{
			AccountInvitationID: exampleInvitationID,
			Input: &identitysvc.AccountInvitationUpdateRequestInput{
				Note: "Cancelling invitation",
			},
		}

		result, err := service.CancelAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_GetAccountInvitation(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitation := identityfakes.BuildFakeAccountInvitation()

		identityDataManager.EXPECT().GetAccountInvitation(testutils.ContextMatcher, mock.AnythingOfType("string"), exampleInvitation.ID).Return(exampleInvitation, nil)

		request := &identitysvc.GetAccountInvitationRequest{
			AccountInvitationID: exampleInvitation.ID,
		}

		result, err := service.GetAccountInvitation(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
		assert.NotNil(t, result.Result)
		assert.Equal(t, exampleInvitation.ID, result.Result.ID)
		assert.Equal(t, exampleInvitation.ToEmail, result.Result.ToEmail)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		request := &identitysvc.GetAccountInvitationRequest{
			AccountInvitationID: identityfakes.BuildFakeID(),
		}

		result, err := service.GetAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitationID := identityfakes.BuildFakeID()

		identityDataManager.EXPECT().GetAccountInvitation(testutils.ContextMatcher, mock.AnythingOfType("string"), exampleInvitationID).Return(nil, errors.New("get error"))

		request := &identitysvc.GetAccountInvitationRequest{
			AccountInvitationID: exampleInvitationID,
		}

		result, err := service.GetAccountInvitation(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_GetReceivedAccountInvitations(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitations := []*identity.AccountInvitation{
			identityfakes.BuildFakeAccountInvitation(),
			identityfakes.BuildFakeAccountInvitation(),
		}

		identityDataManager.EXPECT().GetReceivedAccountInvitations(testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*filtering.QueryFilter")).Return(exampleInvitations, "", nil)

		pageSize := uint32(25)
		request := &identitysvc.GetReceivedAccountInvitationsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &pageSize,
			},
		}

		result, err := service.GetReceivedAccountInvitations(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
		assert.Len(t, result.Result, len(exampleInvitations))
		assert.Equal(t, exampleInvitations[0].ID, result.Result[0].ID)
		assert.Equal(t, exampleInvitations[1].ID, result.Result[1].ID)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		pageSize := uint32(25)
		request := &identitysvc.GetReceivedAccountInvitationsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &pageSize,
			},
		}

		result, err := service.GetReceivedAccountInvitations(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		identityDataManager.EXPECT().GetReceivedAccountInvitations(testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*filtering.QueryFilter")).Return(nil, "", errors.New("get error"))

		pageSize := uint32(25)
		request := &identitysvc.GetReceivedAccountInvitationsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &pageSize,
			},
		}

		result, err := service.GetReceivedAccountInvitations(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}

func TestServiceImpl_GetSentAccountInvitations(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		exampleInvitations := []*identity.AccountInvitation{
			identityfakes.BuildFakeAccountInvitation(),
			identityfakes.BuildFakeAccountInvitation(),
		}

		identityDataManager.EXPECT().GetSentAccountInvitations(testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*filtering.QueryFilter")).Return(exampleInvitations, "", nil)

		pageSize := uint32(25)
		request := &identitysvc.GetSentAccountInvitationsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &pageSize,
			},
		}

		result, err := service.GetSentAccountInvitations(t.Context(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.ResponseDetails)
		assert.Len(t, result.Result, len(exampleInvitations))
		assert.Equal(t, exampleInvitations[0].ID, result.Result[0].ID)
		assert.Equal(t, exampleInvitations[1].ID, result.Result[1].ID)
	})

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		service := buildTestServiceWithSessionError(t)

		pageSize := uint32(25)
		request := &identitysvc.GetSentAccountInvitationsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &pageSize,
			},
		}

		result, err := service.GetSentAccountInvitations(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcErr.Code())
	})

	t.Run("with error from data manager", func(t *testing.T) {
		t.Parallel()

		service, identityDataManager := buildTestService(t)

		identityDataManager.EXPECT().GetSentAccountInvitations(testutils.ContextMatcher, mock.AnythingOfType("string"), mock.AnythingOfType("*filtering.QueryFilter")).Return(nil, "", errors.New("get error"))

		pageSize := uint32(25)
		request := &identitysvc.GetSentAccountInvitationsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &pageSize,
			},
		}

		result, err := service.GetSentAccountInvitations(t.Context(), request)

		assert.Error(t, err)
		assert.Nil(t, result)

		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, grpcErr.Code())
	})
}
