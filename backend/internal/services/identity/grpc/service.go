package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/identity/managers"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
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
		identityDataManager       managers.IdentityDataManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	sessionContextDataFetcher func(ctx context.Context) (*sessions.ContextData, error),
	identityDataManager managers.IdentityDataManager,
) identitysvc.IdentityServiceServer {
	return &serviceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		sessionContextDataFetcher: sessionContextDataFetcher,
		identityDataManager:       identityDataManager,
	}
}

func (s *serviceImpl) buildResponseDetails(ctx context.Context, span tracing.Span) *types.ResponseDetails {
	out := &types.ResponseDetails{}
	if span != nil {
		out.TraceID = span.SpanContext().TraceID().String()
	}

	if sessionContextData, err := s.sessionContextDataFetcher(ctx); err == nil && sessionContextData != nil {
		out.CurrentAccountID = sessionContextData.GetActiveAccountID()
	}

	return out
}
