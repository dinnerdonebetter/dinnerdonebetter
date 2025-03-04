package gprc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
)

func (s *serviceImpl) LoginForToken(ctx context.Context, request *messages.LoginForTokenRequest) (*messages.LoginForTokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *serviceImpl) RequestPasswordResetToken(ctx context.Context, request *messages.RequestPasswordResetTokenRequest) (*messages.RequestPasswordResetTokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *serviceImpl) ExchangeToken(ctx context.Context, request *messages.ExchangeTokenRequest) (*messages.ExchangeTokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *serviceImpl) AdminLoginForToken(ctx context.Context, request *messages.AdminLoginForTokenRequest) (*messages.AdminLoginForTokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *serviceImpl) CheckPermissions(ctx context.Context, request *messages.CheckPermissionsRequest) (*messages.CheckPermissionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
