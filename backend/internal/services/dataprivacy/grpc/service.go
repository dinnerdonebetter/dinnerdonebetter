package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/dataprivacy"
	dataprivacysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "data_privacy_service"
)

var _ dataprivacysvc.DataPrivacyServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		dataprivacysvc.UnimplementedDataPrivacyServiceServer
		tracer                tracing.Tracer
		logger                logging.Logger
		dataPrivacyRepository dataprivacy.Repository
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	dataPrivacyRepository dataprivacy.Repository,
) dataprivacysvc.DataPrivacyServiceServer {
	return &serviceImpl{
		logger:                logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		dataPrivacyRepository: dataPrivacyRepository,
	}
}

func (s *serviceImpl) AggregateUserDataReport(ctx context.Context, request *dataprivacysvc.AggregateUserDataReportRequest) (*dataprivacysvc.AggregateUserDataReportResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serviceImpl) DestroyAllUserData(ctx context.Context, request *dataprivacysvc.DestroyAllUserDataRequest) (*dataprivacysvc.DestroyAllUserDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serviceImpl) FetchUserDataReport(ctx context.Context, request *dataprivacysvc.FetchUserDataReportRequest) (*dataprivacysvc.FetchUserDataReportResponse, error) {
	//TODO implement me
	panic("implement me")
}
