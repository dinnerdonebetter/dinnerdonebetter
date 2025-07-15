package grpc

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var _ identitysvc.IdentityServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		identitysvc.UnimplementedIdentityServiceServer
		tracer             tracing.Tracer
		logger             logging.Logger
		identityRepository identity.Repository
	}
)

func NewService(
	identityRepository identity.Repository,
) identitysvc.IdentityServiceServer {
	return &ServiceImpl{
		identityRepository: identityRepository,
	}
}
