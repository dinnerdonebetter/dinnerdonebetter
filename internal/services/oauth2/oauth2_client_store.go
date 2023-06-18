package oauth2

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/go-oauth2/oauth2/v4"
)

type oauth2ClientStoreImpl struct {
	domain string
	tracer tracing.Tracer
	logger logging.Logger
}

func (s *oauth2ClientStoreImpl) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, nil
}
