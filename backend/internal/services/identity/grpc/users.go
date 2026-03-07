package grpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	identitykeys "github.com/dinnerdonebetter/backend/internal/domain/identity/keys"
	"github.com/dinnerdonebetter/backend/internal/domain/uploadedmedia"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	uploadedmediasvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/uploaded_media"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	platformkeys "github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/services/identity/grpc/converters"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const maxAvatarUploadSize = 5 * 1024 * 1024 // 5 MB for avatars

func (s *serviceImpl) CreateUser(ctx context.Context, request *identitysvc.CreateUserRequest) (*identitysvc.CreateUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	input := converters.ConvertGRPCUserRegistrationInputToUserRegistrationInput(request.Input)

	created, err := s.identityDataManager.CreateUser(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to create user")
	}

	x := &identitysvc.CreateUserResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Created:         converters.ConvertUserCreationResponseToGRPCUserCreationResponse(created),
	}

	return x, nil
}

func (s *serviceImpl) ArchiveUser(ctx context.Context, request *identitysvc.ArchiveUserRequest) (*identitysvc.ArchiveUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	if err := s.identityDataManager.ArchiveUser(ctx, request.UserId); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to archive user")
	}

	x := &identitysvc.ArchiveUserResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) GetUser(ctx context.Context, request *identitysvc.GetUserRequest) (*identitysvc.GetUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		identitykeys.UserIDKey: request.UserId,
	}, span, s.logger)

	user, err := s.identityDataManager.GetUser(ctx, request.UserId)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch users from database")
	}

	x := &identitysvc.GetUserResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Result:          converters.ConvertUserToGRPCUser(user),
	}

	return x, nil
}

func (s *serviceImpl) GetUsers(ctx context.Context, request *identitysvc.GetUsersRequest) (*identitysvc.GetUsersResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	users, err := s.identityDataManager.GetUsers(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch users from database")
	}

	x := &identitysvc.GetUsersResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Pagination:      grpcconverters.ConvertPaginationToGRPCPagination(users.Pagination, filter),
	}

	for _, user := range users.Data {
		x.Results = append(x.Results, converters.ConvertUserToGRPCUser(user))
	}

	return x, nil
}

func (s *serviceImpl) GetUsersForAccount(ctx context.Context, request *identitysvc.GetUsersForAccountRequest) (*identitysvc.GetUsersForAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	users, err := s.identityDataManager.GetUsersForAccount(ctx, request.AccountId, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch users from database")
	}

	x := &identitysvc.GetUsersForAccountResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Pagination:      grpcconverters.ConvertPaginationToGRPCPagination(users.Pagination, filter),
	}

	for _, user := range users.Data {
		x.Results = append(x.Results, converters.ConvertUserToGRPCUser(user))
	}

	return x, nil
}

func (s *serviceImpl) SearchForUsers(ctx context.Context, request *identitysvc.SearchForUsersRequest) (*identitysvc.SearchForUsersResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		platformkeys.SearchQueryKey: request.Query,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	users, err := s.identityDataManager.SearchForUsers(ctx, request.Query, request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to search for users")
	}

	x := &identitysvc.SearchForUsersResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Pagination:      grpcconverters.ConvertPaginationToGRPCPagination(users.Pagination, filter),
	}

	for _, user := range users.Data {
		x.Results = append(x.Results, converters.ConvertUserToGRPCUser(user))
	}

	return x, nil
}

func (s *serviceImpl) UpdateUserDetails(ctx context.Context, request *identitysvc.UpdateUserDetailsRequest) (*identitysvc.UpdateUserDetailsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCUserDetailsUpdateRequestInputToUserDetailsUpdateRequestInput(request.Input)

	if err = s.identityDataManager.UpdateUserDetails(ctx, sessionContextData.GetUserID(), input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update user details")
	}

	x := &identitysvc.UpdateUserDetailsResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) UpdateUserEmailAddress(ctx context.Context, request *identitysvc.UpdateUserEmailAddressRequest) (*identitysvc.UpdateUserEmailAddressResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if err = s.identityDataManager.UpdateUserEmailAddress(ctx, sessionContextData.GetUserID(), request.NewEmailAddress); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update user email address")
	}

	x := &identitysvc.UpdateUserEmailAddressResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) UpdateUserUsername(ctx context.Context, request *identitysvc.UpdateUserUsernameRequest) (*identitysvc.UpdateUserUsernameResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if err = s.identityDataManager.UpdateUserUsername(ctx, sessionContextData.GetUserID(), request.NewUsername); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update user username")
	}

	x := &identitysvc.UpdateUserUsernameResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) UploadUserAvatar(stream grpc.ClientStreamingServer[uploadedmediasvc.UploadRequest, identitysvc.UploadUserAvatarResponse]) error {
	ctx := stream.Context()
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}
	userID := sessionContextData.GetUserID()
	logger = logger.WithValue(identitykeys.UserIDKey, userID)

	firstReq, err := stream.Recv()
	if err != nil {
		return observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "failed to receive metadata")
	}

	metadata := firstReq.GetMetadata()
	if metadata == nil {
		return observability.PrepareAndLogGRPCStatus(
			errors.New("first message must contain metadata"),
			logger, span, codes.InvalidArgument, "first message must contain metadata",
		)
	}

	if metadata.ObjectName == "" {
		return observability.PrepareAndLogGRPCStatus(
			errors.New("object_name is required"),
			logger, span, codes.InvalidArgument, "object_name is required",
		)
	}

	if metadata.ContentType == "" {
		return observability.PrepareAndLogGRPCStatus(
			errors.New("content_type is required"),
			logger, span, codes.InvalidArgument, "content_type is required",
		)
	}

	mimeType := metadata.ContentType
	if !uploadedmedia.IsValidMimeType(mimeType) {
		return observability.PrepareAndLogGRPCStatus(
			fmt.Errorf("unsupported content type: %s", mimeType),
			logger, span, codes.InvalidArgument, "unsupported content type",
		)
	}

	var fileData bytes.Buffer
	totalSize := int64(0)

	for {
		req, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return observability.PrepareAndLogGRPCStatus(recvErr, logger, span, codes.Internal, "failed to receive chunk")
		}

		chunk := req.GetChunk()
		if chunk == nil {
			continue
		}

		chunkSize := int64(len(chunk))
		if totalSize+chunkSize > maxAvatarUploadSize {
			return observability.PrepareAndLogGRPCStatus(
				fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxAvatarUploadSize),
				logger, span, codes.InvalidArgument, "file too large",
			)
		}

		if _, err = fileData.Write(chunk); err != nil {
			return observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to write chunk")
		}

		totalSize += chunkSize
	}

	if totalSize == 0 {
		return observability.PrepareAndLogGRPCStatus(
			errors.New("no file data received"),
			logger, span, codes.InvalidArgument, "no file data received",
		)
	}

	fileID := identifiers.New()
	storagePath := filepath.Join(userID, fileID, metadata.ObjectName)

	if err = s.uploadManager.SaveFile(ctx, storagePath, fileData.Bytes()); err != nil {
		return observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to save file")
	}

	uploadedMediaInput := &uploadedmedia.UploadedMediaDatabaseCreationInput{
		ID:            fileID,
		StoragePath:   storagePath,
		MimeType:      mimeType,
		CreatedByUser: userID,
	}

	if err = uploadedMediaInput.ValidateWithContext(ctx); err != nil {
		return observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "failed to validate uploaded media")
	}

	created, err := s.uploadedMediaManager.CreateUploadedMedia(ctx, uploadedMediaInput)
	if err != nil {
		return observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create uploaded media record")
	}

	if err = s.identityDataManager.SetUserAvatar(ctx, userID, created.ID); err != nil {
		return observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to set user avatar")
	}

	response := &identitysvc.UploadUserAvatarResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	if err = stream.SendAndClose(response); err != nil {
		return observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to send response")
	}

	logger.Info("avatar uploaded successfully")
	return nil
}
