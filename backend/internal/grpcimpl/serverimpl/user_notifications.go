package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) CreateUserNotification(ctx context.Context, request *messages.CreateUserNotificationRequest) (*messages.CreateUserNotificationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUserNotification(ctx context.Context, request *messages.GetUserNotificationRequest) (*messages.GetUserNotificationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUserNotifications(ctx context.Context, request *messages.GetUserNotificationsRequest) (*messages.GetUserNotificationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
