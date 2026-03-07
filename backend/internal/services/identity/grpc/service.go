package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	uploadedmediamanager "github.com/dinnerdonebetter/backend/internal/domain/uploadedmedia/manager"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads"
)

const (
	o11yName = "identity_service"
)

var _ identitysvc.IdentityServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		identitysvc.UnimplementedIdentityServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(ctx context.Context) (*sessions.ContextData, error)
		identityDataManager       manager.IdentityDataManager
		uploadedMediaManager      uploadedmediamanager.UploadedMediaManager
		uploadManager             uploads.UploadManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	sessionContextDataFetcher func(ctx context.Context) (*sessions.ContextData, error),
	identityDataManager manager.IdentityDataManager,
	uploadedMediaManager uploadedmediamanager.UploadedMediaManager,
	uploadManager uploads.UploadManager,
) identitysvc.IdentityServiceServer {
	return &serviceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		sessionContextDataFetcher: sessionContextDataFetcher,
		identityDataManager:       identityDataManager,
		uploadedMediaManager:      uploadedMediaManager,
		uploadManager:             uploadManager,
	}
}

func (s *serviceImpl) buildResponseDetails(ctx context.Context, span tracing.Span) *types.ResponseDetails {
	out := &types.ResponseDetails{}
	if span != nil {
		out.TraceId = span.SpanContext().TraceID().String()
	}

	if sessionContextData, err := s.sessionContextDataFetcher(ctx); err == nil && sessionContextData != nil {
		out.CurrentAccountId = sessionContextData.GetActiveAccountID()
	}

	return out
}
