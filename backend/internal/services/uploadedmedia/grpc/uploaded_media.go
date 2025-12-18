package grpc

import (
	"context"
	"errors"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	uploadedmediasvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/uploaded_media"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/services/uploadedmedia/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) CreateUploadedMedia(ctx context.Context, request *uploadedmediasvc.CreateUploadedMediaRequest) (*uploadedmediasvc.CreateUploadedMediaResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.Requester.UserID)

	input := converters.ConvertGRPCUploadedMediaCreationRequestInputToUploadedMediaDatabaseCreationInput(request.Input, sessionContextData.Requester.UserID)
	if err = input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "failed to validate uploaded media creation request")
	}

	created, err := s.uploadedMediaRepository.CreateUploadedMedia(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create uploaded media")
	}

	x := &uploadedmediasvc.CreateUploadedMediaResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId:          span.SpanContext().TraceID().String(),
			CurrentAccountId: sessionContextData.ActiveAccountID,
		},
		Created: converters.ConvertUploadedMediaToGRPCUploadedMedia(created),
	}

	return x, nil
}

func (s *serviceImpl) GetUploadedMedia(ctx context.Context, request *uploadedmediasvc.GetUploadedMediaRequest) (*uploadedmediasvc.GetUploadedMediaResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.UploadedMediaIDKey, request.UploadedMediaId)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.Requester.UserID)

	uploadedMedia, err := s.uploadedMediaRepository.GetUploadedMedia(ctx, request.UploadedMediaId)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch uploaded media")
	}

	// Verify the uploaded media belongs to the user
	if uploadedMedia.CreatedByUser != sessionContextData.Requester.UserID {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("permission denied"), logger, span, codes.PermissionDenied, "uploaded media does not belong to user")
	}

	x := &uploadedmediasvc.GetUploadedMediaResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId:          span.SpanContext().TraceID().String(),
			CurrentAccountId: sessionContextData.ActiveAccountID,
		},
		Result: converters.ConvertUploadedMediaToGRPCUploadedMedia(uploadedMedia),
	}

	return x, nil
}

func (s *serviceImpl) GetUploadedMediaWithIDs(ctx context.Context, request *uploadedmediasvc.GetUploadedMediaWithIDsRequest) (*uploadedmediasvc.GetUploadedMediaWithIDsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.Requester.UserID)

	if len(request.Ids) == 0 {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("no IDs provided"), logger, span, codes.InvalidArgument, "no IDs provided")
	}

	uploadedMediaList, err := s.uploadedMediaRepository.GetUploadedMediaWithIDs(ctx, request.Ids)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch uploaded media")
	}

	x := &uploadedmediasvc.GetUploadedMediaWithIDsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId:          span.SpanContext().TraceID().String(),
			CurrentAccountId: sessionContextData.ActiveAccountID,
		},
	}

	for _, uploadedMedia := range uploadedMediaList {
		// Only return media that belongs to the user
		if uploadedMedia.CreatedByUser == sessionContextData.Requester.UserID {
			x.Results = append(x.Results, converters.ConvertUploadedMediaToGRPCUploadedMedia(uploadedMedia))
		}
	}

	return x, nil
}

func (s *serviceImpl) GetUploadedMediaForUser(ctx context.Context, request *uploadedmediasvc.GetUploadedMediaForUserRequest) (*uploadedmediasvc.GetUploadedMediaForUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.UserIDKey, request.UserId)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	// Verify the user is requesting their own media
	if request.UserId != sessionContextData.Requester.UserID {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("permission denied"), logger, span, codes.PermissionDenied, "cannot access other user's media")
	}

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	uploadedMediaList, err := s.uploadedMediaRepository.GetUploadedMediaForUser(ctx, request.UserId, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch uploaded media for user")
	}

	x := &uploadedmediasvc.GetUploadedMediaForUserResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId:          span.SpanContext().TraceID().String(),
			CurrentAccountId: sessionContextData.ActiveAccountID,
		},
		Pagination: grpcconverters.ConvertPaginationToGRPCPagination(uploadedMediaList.Pagination, filter),
	}

	for _, uploadedMedia := range uploadedMediaList.Data {
		x.Results = append(x.Results, converters.ConvertUploadedMediaToGRPCUploadedMedia(uploadedMedia))
	}

	return x, nil
}

func (s *serviceImpl) UpdateUploadedMedia(ctx context.Context, request *uploadedmediasvc.UpdateUploadedMediaRequest) (*uploadedmediasvc.UpdateUploadedMediaResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.UploadedMediaIDKey, request.UploadedMediaId)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.Requester.UserID)

	// Fetch the existing uploaded media
	uploadedMedia, err := s.uploadedMediaRepository.GetUploadedMedia(ctx, request.UploadedMediaId)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch uploaded media")
	}

	// Verify the uploaded media belongs to the user
	if uploadedMedia.CreatedByUser != sessionContextData.Requester.UserID {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("permission denied"), logger, span, codes.PermissionDenied, "uploaded media does not belong to user")
	}

	// Apply updates
	updateInput := converters.ConvertGRPCUploadedMediaUpdateRequestInputToUploadedMediaUpdateRequestInput(request.Input)
	if err = updateInput.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "failed to validate uploaded media update request")
	}

	uploadedMedia.Update(updateInput)

	if err = s.uploadedMediaRepository.UpdateUploadedMedia(ctx, uploadedMedia); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update uploaded media")
	}

	x := &uploadedmediasvc.UpdateUploadedMediaResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId:          span.SpanContext().TraceID().String(),
			CurrentAccountId: sessionContextData.ActiveAccountID,
		},
		Updated: converters.ConvertUploadedMediaToGRPCUploadedMedia(uploadedMedia),
	}

	return x, nil
}

func (s *serviceImpl) ArchiveUploadedMedia(ctx context.Context, request *uploadedmediasvc.ArchiveUploadedMediaRequest) (*uploadedmediasvc.ArchiveUploadedMediaResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.UploadedMediaIDKey, request.UploadedMediaId)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.Requester.UserID)

	// Fetch the existing uploaded media to verify ownership
	uploadedMedia, err := s.uploadedMediaRepository.GetUploadedMedia(ctx, request.UploadedMediaId)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch uploaded media")
	}

	// Verify the uploaded media belongs to the user
	if uploadedMedia.CreatedByUser != sessionContextData.Requester.UserID {
		return nil, observability.PrepareAndLogGRPCStatus(errors.New("permission denied"), logger, span, codes.PermissionDenied, "uploaded media does not belong to user")
	}

	if err = s.uploadedMediaRepository.ArchiveUploadedMedia(ctx, request.UploadedMediaId); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive uploaded media")
	}

	x := &uploadedmediasvc.ArchiveUploadedMediaResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceId:          span.SpanContext().TraceID().String(),
			CurrentAccountId: sessionContextData.ActiveAccountID,
		},
	}

	return x, nil
}
