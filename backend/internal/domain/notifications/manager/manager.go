package manager

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	notificationkeys "github.com/dinnerdonebetter/backend/internal/domain/notifications/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "notifications_data_manager"
)

// notificationsRepo avoids wire cycles: manager takes this interface and produces notifications.Repository.
type notificationsRepo interface {
	notifications.Repository
}

var (
	_ notifications.Repository = (*notificationsManager)(nil)
	_ NotificationsDataManager = (*notificationsManager)(nil)
)

type notificationsManager struct {
	tracer               tracing.Tracer
	logger               logging.Logger
	repo                 notificationsRepo
	dataChangesPublisher messagequeue.Publisher
}

// NewNotificationsDataManager returns a new NotificationsDataManager implementing notifications.Repository.
func NewNotificationsDataManager(
	ctx context.Context,
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	repo notificationsRepo,
	cfg *msgconfig.QueuesConfig,
	publisherProvider messagequeue.PublisherProvider,
) (NotificationsDataManager, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(ctx, cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("failed to provide publisher for data changes topic: %w", err)
	}

	return &notificationsManager{
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:               logging.EnsureLogger(logger).WithName(o11yName),
		repo:                 repo,
		dataChangesPublisher: dataChangesPublisher,
	}, nil
}

func (m *notificationsManager) UserNotificationExists(ctx context.Context, userID, userNotificationID string) (bool, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.UserNotificationExists(ctx, userID, userNotificationID)
}

func (m *notificationsManager) GetUserNotification(ctx context.Context, userID, userNotificationID string) (*notifications.UserNotification, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetUserNotification(ctx, userID, userNotificationID)
}

func (m *notificationsManager) GetUserNotifications(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[notifications.UserNotification], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetUserNotifications(ctx, userID, filter)
}

func (m *notificationsManager) CreateUserNotification(ctx context.Context, input *notifications.UserNotificationDatabaseCreationInput) (*notifications.UserNotification, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}
	logger := m.logger.WithSpan(span).WithValue(notificationkeys.UserNotificationIDKey, input.ID)
	tracing.AttachToSpan(span, notificationkeys.UserNotificationIDKey, input.ID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating user notification creation input")
	}

	created, err := m.repo.CreateUserNotification(ctx, input)
	if err != nil {
		return nil, err
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, notifications.UserNotificationCreatedServiceEventType, map[string]any{
		notificationkeys.UserNotificationIDKey: created.ID,
	}))

	return created, nil
}

func (m *notificationsManager) UpdateUserNotification(ctx context.Context, updated *notifications.UserNotification) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return internalerrors.ErrNilInputParameter
	}
	logger := m.logger.WithSpan(span).WithValue(notificationkeys.UserNotificationIDKey, updated.ID)
	tracing.AttachToSpan(span, notificationkeys.UserNotificationIDKey, updated.ID)

	if err := m.repo.UpdateUserNotification(ctx, updated); err != nil {
		return err
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, notifications.UserNotificationUpdatedServiceEventType, map[string]any{
		notificationkeys.UserNotificationIDKey: updated.ID,
	}))

	return nil
}
