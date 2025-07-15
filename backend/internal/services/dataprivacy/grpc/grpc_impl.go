package grpc

import (
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
