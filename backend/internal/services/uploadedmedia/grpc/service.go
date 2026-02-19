package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	uploadedmediamanager "github.com/dinnerdonebetter/backend/internal/domain/uploadedmedia/manager"
	uploadedmediasvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/uploaded_media"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads"
)

const (
	o11yName = "uploaded_media_service"
)

var _ uploadedmediasvc.UploadedMediaServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		uploadedmediasvc.UnimplementedUploadedMediaServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
		uploadedMediaManager      uploadedmediamanager.UploadedMediaManager
		uploadManager             uploads.UploadManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	uploadedMediaManager uploadedmediamanager.UploadedMediaManager,
	uploadManager uploads.UploadManager,
) uploadedmediasvc.UploadedMediaServiceServer {
	return &serviceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
		uploadedMediaManager:      uploadedMediaManager,
		uploadManager:             uploadManager,
	}
}
