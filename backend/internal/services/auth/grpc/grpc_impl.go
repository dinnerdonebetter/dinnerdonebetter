package grpc

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var _ authsvc.AuthServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		authsvc.UnimplementedAuthServiceServer
		tracer             tracing.Tracer
		logger             logging.Logger
		identityRepository identity.Repository
	}
)

func NewService(
	identityRepository identity.Repository,
) authsvc.AuthServiceServer {
	return &ServiceImpl{
		identityRepository: identityRepository,
	}
}
