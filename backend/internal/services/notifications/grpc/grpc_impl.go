package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
	grpctypes "github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"google.golang.org/grpc/codes"
)

const (
	o11yName = "notifications_service"
)

var _ notificationssvc.UserNotificationsServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		notificationssvc.UnimplementedUserNotificationsServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (sessions.ContextData, error)
		notificationsRepository   notifications.Repository
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	notificationsRepository notifications.Repository,
) notificationssvc.UserNotificationsServiceServer {
	return &serviceImpl{
		logger:                  logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		notificationsRepository: notificationsRepository,
	}
}

func ConvertUserNotificationToGRPCUserNotification(notification *notifications.UserNotification) *notificationssvc.UserNotification {
	return &notificationssvc.UserNotification{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(notification.CreatedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(notification.LastUpdatedAt),
		ID:            notification.ID,
		Content:       notification.Content,
		Status:        notification.Status,
		BelongsToUser: notification.BelongsToUser,
	}
}

func (s *serviceImpl) GetUserNotification(ctx context.Context, request *notificationssvc.GetUserNotificationRequest) (*notificationssvc.GetUserNotificationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.UserNotificationIDKey, request.UserNotificationID)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "unable to determine authentication")
	}

	notification, err := s.notificationsRepository.GetUserNotification(ctx, sessionContextData.GetUserID(), request.UserNotificationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "")
	}

	x := &notificationssvc.GetUserNotificationResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: ConvertUserNotificationToGRPCUserNotification(notification),
	}

	return x, nil
}

func (s *serviceImpl) GetUserNotifications(ctx context.Context, request *notificationssvc.GetUserNotificationsRequest) (*notificationssvc.GetUserNotificationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "unable to determine authentication")
	}

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	notifs, err := s.notificationsRepository.GetUserNotifications(ctx, sessionContextData.GetUserID(), filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching user notifs")
	}

	x := &notificationssvc.GetUserNotificationsResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}
	for _, notif := range notifs.Data {
		x.Results = append(x.Results, ConvertUserNotificationToGRPCUserNotification(notif))
	}

	return x, nil
}

func (s *serviceImpl) UpdateUserNotification(ctx context.Context, request *notificationssvc.UpdateUserNotificationRequest) (*notificationssvc.UpdateUserNotificationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "unable to determine authentication")
	}

	logger = logger.WithValue(keys.UserNotificationIDKey, request.UserNotificationID)

	existing, err := s.notificationsRepository.GetUserNotification(ctx, sessionContextData.GetUserID(), request.UserNotificationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching existing notification")
	}

	existing.Update(&notifications.UserNotificationUpdateRequestInput{Status: &request.Input.Status})
	if err = s.notificationsRepository.UpdateUserNotification(ctx, existing); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating existing notification")
	}

	x := &notificationssvc.UpdateUserNotificationResponse{
		ResponseDetails: &grpctypes.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: ConvertUserNotificationToGRPCUserNotification(existing),
	}

	return x, nil
}
