package grpc

import (
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
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func buildTestService(t *testing.T) (*serviceImpl, *managermock.OAuth2Manager) {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())
	oauthManager := &managermock.OAuth2Manager{}

	service := &serviceImpl{
		tracer:           tracer,
		logger:           logger,
		oauthDataManager: oauthManager,
	}

	return service, oauthManager
}

func TestServiceImpl_CreateOAuth2Client(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service, mockManager := buildTestService(t)

		fakeClient := oauthfakes.BuildFakeOAuth2Client()
		fakeInput := oauthfakes.BuildFakeOAuth2ClientCreationRequestInput()

		mockManager.On("CreateOAuth2Client", testutils.ContextMatcher, mock.AnythingOfType("*oauth.OAuth2ClientCreationRequestInput")).Return(fakeClient, nil)

		request := &oauthsvc.CreateOAuth2ClientRequest{
			Input: &oauthsvc.OAuth2ClientCreationRequestInput{
				Name: fakeInput.Name,
			},
		}

		response, err := service.CreateOAuth2Client(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.Created)
		assert.NotNil(t, response.ResponseDetails)
		assert.Equal(t, fakeClient.ID, response.Created.ID)
		assert.Equal(t, fakeClient.Name, response.Created.Name)

		mock.AssertExpectationsForObjects(t, mockManager)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service, mockManager := buildTestService(t)

		fakeInput := oauthfakes.BuildFakeOAuth2ClientCreationRequestInput()

		mockManager.On("CreateOAuth2Client", testutils.ContextMatcher, mock.AnythingOfType("*oauth.OAuth2ClientCreationRequestInput")).Return((*oauth.OAuth2Client)(nil), errors.New("manager error"))

		request := &oauthsvc.CreateOAuth2ClientRequest{
			Input: &oauthsvc.OAuth2ClientCreationRequestInput{
				Name: fakeInput.Name,
			},
		}

		response, err := service.CreateOAuth2Client(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockManager)
	})
}

func TestServiceImpl_GetOAuth2Client(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service, mockManager := buildTestService(t)

		fakeClient := oauthfakes.BuildFakeOAuth2Client()
		clientID := fakeClient.ID

		mockManager.On("GetOAuth2Client", testutils.ContextMatcher, clientID).Return(fakeClient, nil)

		request := &oauthsvc.GetOAuth2ClientRequest{
			OAuth2ClientID: clientID,
		}

		response, err := service.GetOAuth2Client(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.ResponseDetails)
		assert.NotNil(t, response.Result)
		assert.Equal(t, fakeClient.ID, response.Result.ID)

		mock.AssertExpectationsForObjects(t, mockManager)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service, mockManager := buildTestService(t)

		clientID := "nonexistent-client"

		mockManager.On("GetOAuth2Client", testutils.ContextMatcher, clientID).Return((*oauth.OAuth2Client)(nil), errors.New("manager error"))

		request := &oauthsvc.GetOAuth2ClientRequest{
			OAuth2ClientID: clientID,
		}

		response, err := service.GetOAuth2Client(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockManager)
	})
}

func TestServiceImpl_GetOAuth2Clients(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service, mockManager := buildTestService(t)

		fakeClients := oauthfakes.BuildFakeOAuth2ClientsList()
		page := uint16(1)
		pageSize := uint8(20)
		filter := &filtering.QueryFilter{
			Page:     &page,
			PageSize: &pageSize,
		}

		mockManager.On("GetOAuth2Clients", testutils.ContextMatcher, mock.AnythingOfType("*filtering.QueryFilter")).Return(fakeClients, nil)

		grpcPageSize := uint32(*filter.PageSize)
		request := &oauthsvc.GetOAuth2ClientsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &grpcPageSize,
			},
		}

		response, err := service.GetOAuth2Clients(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.ResponseDetails)
		assert.Len(t, response.Results, len(fakeClients.Data))

		mock.AssertExpectationsForObjects(t, mockManager)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service, mockManager := buildTestService(t)

		page := uint16(1)
		pageSize := uint8(20)
		filter := &filtering.QueryFilter{
			Page:     &page,
			PageSize: &pageSize,
		}

		mockManager.On("GetOAuth2Clients", testutils.ContextMatcher, mock.AnythingOfType("*filtering.QueryFilter")).Return((*filtering.QueryFilteredResult[oauth.OAuth2Client])(nil), errors.New("manager error"))

		grpcPageSize := uint32(*filter.PageSize)
		request := &oauthsvc.GetOAuth2ClientsRequest{
			Filter: &grpcfiltering.QueryFilter{
				PageSize: &grpcPageSize,
			},
		}

		response, err := service.GetOAuth2Clients(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockManager)
	})
}

func TestServiceImpl_ArchiveOAuth2Client(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service, mockManager := buildTestService(t)

		clientID := "test-client-id"

		mockManager.On("ArchiveOAuth2Client", testutils.ContextMatcher, clientID).Return(nil)

		request := &oauthsvc.ArchiveOAuth2ClientRequest{
			OAuth2ClientID: clientID,
		}

		response, err := service.ArchiveOAuth2Client(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.ResponseDetails)

		mock.AssertExpectationsForObjects(t, mockManager)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service, mockManager := buildTestService(t)

		clientID := "nonexistent-client"

		mockManager.On("ArchiveOAuth2Client", testutils.ContextMatcher, clientID).Return(errors.New("manager error"))

		request := &oauthsvc.ArchiveOAuth2ClientRequest{
			OAuth2ClientID: clientID,
		}

		response, err := service.ArchiveOAuth2Client(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockManager)
	})
}
