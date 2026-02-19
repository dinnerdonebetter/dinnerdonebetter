package manager

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/uploadedmedia"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "uploaded_media_data_manager"
)

var (
	_ uploadedmedia.Repository = (*uploadedMediaManager)(nil)
	_ UploadedMediaManager     = (*uploadedMediaManager)(nil)
)

type uploadedMediaManager struct {
	tracer tracing.Tracer
	logger logging.Logger
	repo   uploadedmedia.Repository
}

// NewUploadedMediaDataManager returns a new UploadedMediaManager that wraps the uploaded media repository.
func NewUploadedMediaDataManager(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	repo uploadedmedia.Repository,
) UploadedMediaManager {
	return &uploadedMediaManager{
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger: logging.EnsureLogger(logger).WithName(o11yName),
		repo:   repo,
	}
}

func (m *uploadedMediaManager) GetUploadedMedia(ctx context.Context, uploadedMediaID string) (*uploadedmedia.UploadedMedia, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()
	return m.repo.GetUploadedMedia(ctx, uploadedMediaID)
}

func (m *uploadedMediaManager) GetUploadedMediaWithIDs(ctx context.Context, ids []string) ([]*uploadedmedia.UploadedMedia, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()
	return m.repo.GetUploadedMediaWithIDs(ctx, ids)
}

func (m *uploadedMediaManager) GetUploadedMediaForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[uploadedmedia.UploadedMedia], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()
	return m.repo.GetUploadedMediaForUser(ctx, userID, filter)
}

func (m *uploadedMediaManager) CreateUploadedMedia(ctx context.Context, input *uploadedmedia.UploadedMediaDatabaseCreationInput) (*uploadedmedia.UploadedMedia, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.UserIDKey, input.CreatedByUser)

	created, err := m.repo.CreateUploadedMedia(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating uploaded media")
	}

	return created, nil
}

func (m *uploadedMediaManager) UpdateUploadedMedia(ctx context.Context, uploadedMedia *uploadedmedia.UploadedMedia) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()
	return m.repo.UpdateUploadedMedia(ctx, uploadedMedia)
}

func (m *uploadedMediaManager) ArchiveUploadedMedia(ctx context.Context, uploadedMediaID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()
	return m.repo.ArchiveUploadedMedia(ctx, uploadedMediaID)
}
