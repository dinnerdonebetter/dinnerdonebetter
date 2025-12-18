package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/uploadedmedia"
	uploadedmediafakes "github.com/dinnerdonebetter/backend/internal/domain/uploadedmedia/fakes"
	uploadedmediamock "github.com/dinnerdonebetter/backend/internal/domain/uploadedmedia/mock"
	grpcfiltering "github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	uploadedmediasvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/uploaded_media"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"
	"github.com/dinnerdonebetter/backend/internal/services/uploadedmedia/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func buildTestService(t *testing.T) (*serviceImpl, *uploadedmediamock.Repository) {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())
	uploadedMediaRepo := &uploadedmediamock.Repository{}

	service := &serviceImpl{
		tracer: tracer,
		logger: logger,
		sessionContextDataFetcher: func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: "test-account-id",
				Requester: sessions.RequesterInfo{
					UserID: "test-user-id",
				},
			}, nil
		},
		uploadedMediaRepository: uploadedMediaRepo,
	}

	return service, uploadedMediaRepo
}

func buildTestServiceWithSessionError(t *testing.T) *serviceImpl {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())

	service := &serviceImpl{
		tracer: tracer,
		logger: logger,
		sessionContextDataFetcher: func(ctx context.Context) (*sessions.ContextData, error) {
			return nil, errors.New("session error")
		},
		uploadedMediaRepository: &uploadedmediamock.Repository{},
	}

	return service
}

func TestServiceImpl_CreateUploadedMedia(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeInput := uploadedmediafakes.BuildFakeUploadedMediaCreationRequestInput()

		mockRepo.On("CreateUploadedMedia", testutils.ContextMatcher, mock.AnythingOfType("*uploadedmedia.UploadedMediaDatabaseCreationInput")).Return(fakeUploadedMedia, nil)

		request := &uploadedmediasvc.CreateUploadedMediaRequest{
			Input: converters.ConvertUploadedMediaCreationRequestInputToGRPCUploadedMediaCreationRequestInput(fakeInput),
		}

		response, err := service.CreateUploadedMedia(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.Created)
		assert.NotNil(t, response.ResponseDetails)
		assert.Equal(t, fakeUploadedMedia.ID, response.Created.Id)
		assert.Equal(t, fakeUploadedMedia.StoragePath, response.Created.StoragePath)

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("session context error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service := buildTestServiceWithSessionError(t)

		request := &uploadedmediasvc.CreateUploadedMediaRequest{
			Input: &uploadedmediasvc.UploadedMediaCreationRequestInput{},
		}

		response, err := service.CreateUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeInput := uploadedmediafakes.BuildFakeUploadedMediaCreationRequestInput()

		mockRepo.On("CreateUploadedMedia", testutils.ContextMatcher, mock.AnythingOfType("*uploadedmedia.UploadedMediaDatabaseCreationInput")).Return(nil, errors.New("repository error"))

		request := &uploadedmediasvc.CreateUploadedMediaRequest{
			Input: converters.ConvertUploadedMediaCreationRequestInputToGRPCUploadedMediaCreationRequestInput(fakeInput),
		}

		response, err := service.CreateUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})
}

func TestServiceImpl_GetUploadedMedia(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia.CreatedByUser = "test-user-id"

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(fakeUploadedMedia, nil)

		request := &uploadedmediasvc.GetUploadedMediaRequest{
			UploadedMediaId: fakeUploadedMedia.ID,
		}

		response, err := service.GetUploadedMedia(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.Result)
		assert.Equal(t, fakeUploadedMedia.ID, response.Result.Id)

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("session context error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service := buildTestServiceWithSessionError(t)

		request := &uploadedmediasvc.GetUploadedMediaRequest{
			UploadedMediaId: "some-id",
		}

		response, err := service.GetUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("permission denied - different user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia.CreatedByUser = "different-user-id"

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(fakeUploadedMedia, nil)

		request := &uploadedmediasvc.GetUploadedMediaRequest{
			UploadedMediaId: fakeUploadedMedia.ID,
		}

		response, err := service.GetUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.PermissionDenied, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, "some-id").Return(nil, errors.New("repository error"))

		request := &uploadedmediasvc.GetUploadedMediaRequest{
			UploadedMediaId: "some-id",
		}

		response, err := service.GetUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})
}

func TestServiceImpl_GetUploadedMediaWithIDs(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia1 := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia1.CreatedByUser = "test-user-id"
		fakeUploadedMedia2 := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia2.CreatedByUser = "test-user-id"

		fakeUploadedMediaList := []*uploadedmedia.UploadedMedia{
			fakeUploadedMedia1,
			fakeUploadedMedia2,
		}

		ids := []string{fakeUploadedMedia1.ID, fakeUploadedMedia2.ID}

		mockRepo.On("GetUploadedMediaWithIDs", testutils.ContextMatcher, ids).Return(fakeUploadedMediaList, nil)

		request := &uploadedmediasvc.GetUploadedMediaWithIDsRequest{
			Ids: ids,
		}

		response, err := service.GetUploadedMediaWithIDs(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Results, 2)

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("filters out other users' media", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia1 := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia1.CreatedByUser = "test-user-id"
		fakeUploadedMedia2 := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia2.CreatedByUser = "other-user-id"

		fakeUploadedMediaList := []*uploadedmedia.UploadedMedia{
			fakeUploadedMedia1,
			fakeUploadedMedia2,
		}

		ids := []string{fakeUploadedMedia1.ID, fakeUploadedMedia2.ID}

		mockRepo.On("GetUploadedMediaWithIDs", testutils.ContextMatcher, ids).Return(fakeUploadedMediaList, nil)

		request := &uploadedmediasvc.GetUploadedMediaWithIDsRequest{
			Ids: ids,
		}

		response, err := service.GetUploadedMediaWithIDs(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Results, 1)
		assert.Equal(t, fakeUploadedMedia1.ID, response.Results[0].Id)

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("session context error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service := buildTestServiceWithSessionError(t)

		request := &uploadedmediasvc.GetUploadedMediaWithIDsRequest{
			Ids: []string{"id1", "id2"},
		}

		response, err := service.GetUploadedMediaWithIDs(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("no IDs provided", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, _ := buildTestService(t)

		request := &uploadedmediasvc.GetUploadedMediaWithIDsRequest{
			Ids: []string{},
		}

		response, err := service.GetUploadedMediaWithIDs(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		ids := []string{"id1", "id2"}

		mockRepo.On("GetUploadedMediaWithIDs", testutils.ContextMatcher, ids).Return(nil, errors.New("repository error"))

		request := &uploadedmediasvc.GetUploadedMediaWithIDsRequest{
			Ids: ids,
		}

		response, err := service.GetUploadedMediaWithIDs(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})
}

func TestServiceImpl_GetUploadedMediaForUser(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMediaList := &filtering.QueryFilteredResult[uploadedmedia.UploadedMedia]{
			Data: []*uploadedmedia.UploadedMedia{
				uploadedmediafakes.BuildFakeUploadedMedia(),
				uploadedmediafakes.BuildFakeUploadedMedia(),
			},
			Pagination: filtering.Pagination{
				TotalCount:    2,
				FilteredCount: 2,
			},
		}

		mockRepo.On("GetUploadedMediaForUser", testutils.ContextMatcher, "test-user-id", mock.AnythingOfType("*filtering.QueryFilter")).Return(fakeUploadedMediaList, nil)

		request := &uploadedmediasvc.GetUploadedMediaForUserRequest{
			UserId: "test-user-id",
			Filter: &grpcfiltering.QueryFilter{},
		}

		response, err := service.GetUploadedMediaForUser(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Results, 2)

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("session context error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service := buildTestServiceWithSessionError(t)

		request := &uploadedmediasvc.GetUploadedMediaForUserRequest{
			UserId: "test-user-id",
			Filter: &grpcfiltering.QueryFilter{},
		}

		response, err := service.GetUploadedMediaForUser(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("permission denied - different user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, _ := buildTestService(t)

		request := &uploadedmediasvc.GetUploadedMediaForUserRequest{
			UserId: "different-user-id",
			Filter: &grpcfiltering.QueryFilter{},
		}

		response, err := service.GetUploadedMediaForUser(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.PermissionDenied, status.Code(err))
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		mockRepo.On("GetUploadedMediaForUser", testutils.ContextMatcher, "test-user-id", mock.AnythingOfType("*filtering.QueryFilter")).Return(nil, errors.New("repository error"))

		request := &uploadedmediasvc.GetUploadedMediaForUserRequest{
			UserId: "test-user-id",
			Filter: &grpcfiltering.QueryFilter{},
		}

		response, err := service.GetUploadedMediaForUser(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})
}

func TestServiceImpl_UpdateUploadedMedia(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia.CreatedByUser = "test-user-id"

		newStoragePath := "updated/path.jpg"

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(fakeUploadedMedia, nil)
		mockRepo.On("UpdateUploadedMedia", testutils.ContextMatcher, mock.AnythingOfType("*uploadedmedia.UploadedMedia")).Return(nil)

		request := &uploadedmediasvc.UpdateUploadedMediaRequest{
			UploadedMediaId: fakeUploadedMedia.ID,
			Input: &uploadedmediasvc.UploadedMediaUpdateRequestInput{
				StoragePath: &newStoragePath,
			},
		}

		response, err := service.UpdateUploadedMedia(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.Updated)

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("session context error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service := buildTestServiceWithSessionError(t)

		request := &uploadedmediasvc.UpdateUploadedMediaRequest{
			UploadedMediaId: "some-id",
			Input:           &uploadedmediasvc.UploadedMediaUpdateRequestInput{},
		}

		response, err := service.UpdateUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("permission denied - different user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia.CreatedByUser = "different-user-id"

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(fakeUploadedMedia, nil)

		request := &uploadedmediasvc.UpdateUploadedMediaRequest{
			UploadedMediaId: fakeUploadedMedia.ID,
			Input:           &uploadedmediasvc.UploadedMediaUpdateRequestInput{},
		}

		response, err := service.UpdateUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.PermissionDenied, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("repository error on get", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, "some-id").Return(nil, errors.New("repository error"))

		request := &uploadedmediasvc.UpdateUploadedMediaRequest{
			UploadedMediaId: "some-id",
			Input:           &uploadedmediasvc.UploadedMediaUpdateRequestInput{},
		}

		response, err := service.UpdateUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("repository error on update", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia.CreatedByUser = "test-user-id"

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(fakeUploadedMedia, nil)
		mockRepo.On("UpdateUploadedMedia", testutils.ContextMatcher, mock.AnythingOfType("*uploadedmedia.UploadedMedia")).Return(errors.New("repository error"))

		request := &uploadedmediasvc.UpdateUploadedMediaRequest{
			UploadedMediaId: fakeUploadedMedia.ID,
			Input:           &uploadedmediasvc.UploadedMediaUpdateRequestInput{},
		}

		response, err := service.UpdateUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})
}

func TestServiceImpl_ArchiveUploadedMedia(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia.CreatedByUser = "test-user-id"

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(fakeUploadedMedia, nil)
		mockRepo.On("ArchiveUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(nil)

		request := &uploadedmediasvc.ArchiveUploadedMediaRequest{
			UploadedMediaId: fakeUploadedMedia.ID,
		}

		response, err := service.ArchiveUploadedMedia(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("session context error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service := buildTestServiceWithSessionError(t)

		request := &uploadedmediasvc.ArchiveUploadedMediaRequest{
			UploadedMediaId: "some-id",
		}

		response, err := service.ArchiveUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("permission denied - different user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia.CreatedByUser = "different-user-id"

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(fakeUploadedMedia, nil)

		request := &uploadedmediasvc.ArchiveUploadedMediaRequest{
			UploadedMediaId: fakeUploadedMedia.ID,
		}

		response, err := service.ArchiveUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.PermissionDenied, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("repository error on get", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, "some-id").Return(nil, errors.New("repository error"))

		request := &uploadedmediasvc.ArchiveUploadedMediaRequest{
			UploadedMediaId: "some-id",
		}

		response, err := service.ArchiveUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})

	t.Run("repository error on archive", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		service, mockRepo := buildTestService(t)

		fakeUploadedMedia := uploadedmediafakes.BuildFakeUploadedMedia()
		fakeUploadedMedia.CreatedByUser = "test-user-id"

		mockRepo.On("GetUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(fakeUploadedMedia, nil)
		mockRepo.On("ArchiveUploadedMedia", testutils.ContextMatcher, fakeUploadedMedia.ID).Return(errors.New("repository error"))

		request := &uploadedmediasvc.ArchiveUploadedMediaRequest{
			UploadedMediaId: fakeUploadedMedia.ID,
		}

		response, err := service.ArchiveUploadedMedia(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))

		mock.AssertExpectationsForObjects(t, mockRepo)
	})
}
