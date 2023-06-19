package oauth2

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/go-oauth2/oauth2/v4"
)

type oauth2ClientStoreImpl struct {
	tracer tracing.Tracer
	logger logging.Logger
	domain string
}

func (s *oauth2ClientStoreImpl) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.OAuth2ClientIDKey, id)
	logger.Debug("getting oauth2 client by ID")

	return &oauth2ClientInfoImpl{}, nil
}
