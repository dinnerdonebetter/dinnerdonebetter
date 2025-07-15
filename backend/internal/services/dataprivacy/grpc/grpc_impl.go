package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/dataprivacy"
	dataprivacysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var _ dataprivacysvc.DataPrivacyServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		dataprivacysvc.UnimplementedDataPrivacyServiceServer
		tracer                tracing.Tracer
		logger                logging.Logger
		dataPrivacyRepository dataprivacy.Repository
	}
)

func NewService(
	dataPrivacyRepository dataprivacy.Repository,
) dataprivacysvc.DataPrivacyServiceServer {
	return &ServiceImpl{
		dataPrivacyRepository: dataPrivacyRepository,
	}
}

func (s *ServiceImpl) AggregateUserDataReport(ctx context.Context, request *dataprivacysvc.AggregateUserDataReportRequest) (*dataprivacysvc.AggregateUserDataReportResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) DestroyAllUserData(ctx context.Context, request *dataprivacysvc.DestroyAllUserDataRequest) (*dataprivacysvc.DestroyAllUserDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) FetchUserDataReport(ctx context.Context, request *dataprivacysvc.FetchUserDataReportRequest) (*dataprivacysvc.FetchUserDataReportResponse, error) {
	//TODO implement me
	panic("implement me")
}
