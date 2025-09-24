package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	oauthfakes "github.com/dinnerdonebetter/backend/internal/domain/oauth/fakes"
	managermock "github.com/dinnerdonebetter/backend/internal/domain/oauth/manager/mock"
	grpcfiltering "github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func buildTestService(t *testing.T) (*serviceImpl, *managermock.OAuth2Manager) {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())
	oauthDataManager := &managermock.OAuth2Manager{}

	service := &serviceImpl{
		tracer:           tracer,
		logger:           logger,
		oauthDataManager: oauthDataManager,
	}

	return service, oauthDataManager
}

func TestServiceImpl_ArchiveOAuth2Client(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		oauth2ClientID := "test-oauth2-client-id"

		service, oauthDataManager := buildTestService(t)

		request := &oauthsvc.ArchiveOAuth2ClientRequest{
			OAuth2ClientID: oauth2ClientID,
		}

		oauthDataManager.On("ArchiveOAuth2Client", mock.MatchedBy(func(ctx context.Context) bool { return true }), oauth2ClientID).Return(nil)

		response, err := service.ArchiveOAuth2Client(ctx, request)

		require.NotNil(t, response)
		assert.NoError(t, err)
		assert.NotNil(t, response.ResponseDetails)
		assert.NotEmpty(t, response.ResponseDetails.TraceID)

		mock.AssertExpectationsForObjects(t, oauthDataManager)
	})

	t.Run("with manager error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		oauth2ClientID := "test-oauth2-client-id"
		expectedError := errors.New("test error")

		service, oauthDataManager := buildTestService(t)

		request := &oauthsvc.ArchiveOAuth2ClientRequest{
			OAuth2ClientID: oauth2ClientID,
		}

		oauthDataManager.On("ArchiveOAuth2Client", mock.MatchedBy(func(ctx context.Context) bool { return true }), oauth2ClientID).Return(expectedError)

		response, err := service.ArchiveOAuth2Client(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, grpcStatus.Code())

		mock.AssertExpectationsForObjects(t, oauthDataManager)
	})
}

func TestServiceImpl_CreateOAuth2Client(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		createdOAuth2Client := oauthfakes.BuildFakeOAuth2Client()

		service, oauthDataManager := buildTestService(t)

		request := &oauthsvc.CreateOAuth2ClientRequest{
			Input: &oauthsvc.OAuth2ClientCreationRequestInput{
				Name: "test oauth2 client",
			},
		}

		oauthDataManager.On("CreateOAuth2Client", mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("*oauth.OAuth2ClientCreationRequestInput")).Return(createdOAuth2Client, nil)

		response, err := service.CreateOAuth2Client(ctx, request)

		require.NotNil(t, response)
		assert.NoError(t, err)
		assert.NotNil(t, response.ResponseDetails)
		assert.NotEmpty(t, response.ResponseDetails.TraceID)
		assert.NotNil(t, response.Created)
		assert.Equal(t, createdOAuth2Client.ID, response.Created.ID)
		assert.Equal(t, createdOAuth2Client.Name, response.Created.Name)

		mock.AssertExpectationsForObjects(t, oauthDataManager)
	})

	t.Run("with manager error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedError := errors.New("test error")

		service, oauthDataManager := buildTestService(t)

		request := &oauthsvc.CreateOAuth2ClientRequest{
			Input: &oauthsvc.OAuth2ClientCreationRequestInput{
				Name: "test oauth2 client",
			},
		}

		oauthDataManager.On("CreateOAuth2Client", mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("*oauth.OAuth2ClientCreationRequestInput")).Return((*oauth.OAuth2Client)(nil), expectedError)

		response, err := service.CreateOAuth2Client(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, grpcStatus.Code())

		mock.AssertExpectationsForObjects(t, oauthDataManager)
	})
}

func TestServiceImpl_GetOAuth2Client(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		oauth2Client := oauthfakes.BuildFakeOAuth2Client()

		service, oauthDataManager := buildTestService(t)

		request := &oauthsvc.GetOAuth2ClientRequest{
			OAuth2ClientID: oauth2Client.ID,
		}

		oauthDataManager.On("GetOAuth2Client", mock.MatchedBy(func(ctx context.Context) bool { return true }), oauth2Client.ID).Return(oauth2Client, nil)

		response, err := service.GetOAuth2Client(ctx, request)

		require.NotNil(t, response)
		assert.NoError(t, err)
		assert.NotNil(t, response.ResponseDetails)
		assert.NotEmpty(t, response.ResponseDetails.TraceID)
		assert.NotNil(t, response.Result)
		assert.Equal(t, oauth2Client.ID, response.Result.ID)
		assert.Equal(t, oauth2Client.Name, response.Result.Name)

		mock.AssertExpectationsForObjects(t, oauthDataManager)
	})

	t.Run("with manager error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		oauth2ClientID := "test-oauth2-client-id"
		expectedError := errors.New("test error")

		service, oauthDataManager := buildTestService(t)

		request := &oauthsvc.GetOAuth2ClientRequest{
			OAuth2ClientID: oauth2ClientID,
		}

		oauthDataManager.On("GetOAuth2Client", mock.MatchedBy(func(ctx context.Context) bool { return true }), oauth2ClientID).Return((*oauth.OAuth2Client)(nil), expectedError)

		response, err := service.GetOAuth2Client(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, grpcStatus.Code())

		mock.AssertExpectationsForObjects(t, oauthDataManager)
	})
}

func TestServiceImpl_GetOAuth2Clients(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		oauth2Clients := []*oauth.OAuth2Client{
			oauthfakes.BuildFakeOAuth2Client(),
			oauthfakes.BuildFakeOAuth2Client(),
		}

		service, oauthDataManager := buildTestService(t)

		request := &oauthsvc.GetOAuth2ClientsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: func(x uint32) *uint32 { return &x }(20),
			},
		}

		expectedResult := &filtering.QueryFilteredResult[oauth.OAuth2Client]{
			Data: oauth2Clients,
			Pagination: filtering.Pagination{
				Page:          1,
				Limit:         20,
				FilteredCount: uint64(len(oauth2Clients)),
				TotalCount:    uint64(len(oauth2Clients)),
			},
		}

		oauthDataManager.On("GetOAuth2Clients", mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("*filtering.QueryFilter")).Return(expectedResult, nil)

		response, err := service.GetOAuth2Clients(ctx, request)

		require.NotNil(t, response)
		assert.NoError(t, err)
		assert.NotNil(t, response.ResponseDetails)
		assert.NotEmpty(t, response.ResponseDetails.TraceID)
		assert.NotNil(t, response.Results)
		assert.Len(t, response.Results, len(oauth2Clients))

		for i, result := range response.Results {
			assert.Equal(t, oauth2Clients[i].ID, result.ID)
			assert.Equal(t, oauth2Clients[i].Name, result.Name)
		}

		mock.AssertExpectationsForObjects(t, oauthDataManager)
	})

	t.Run("with manager error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedError := errors.New("test error")

		service, oauthDataManager := buildTestService(t)

		request := &oauthsvc.GetOAuth2ClientsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: func(x uint32) *uint32 { return &x }(20),
			},
		}

		oauthDataManager.On("GetOAuth2Clients", mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("*filtering.QueryFilter")).Return((*filtering.QueryFilteredResult[oauth.OAuth2Client])(nil), expectedError)

		response, err := service.GetOAuth2Clients(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, grpcStatus.Code())

		mock.AssertExpectationsForObjects(t, oauthDataManager)
	})
}
