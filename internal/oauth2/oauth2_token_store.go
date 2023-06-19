package oauth2

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/go-oauth2/oauth2/v4"
)

type oauth2TokenStoreImpl struct {
	tracer tracing.Tracer
	logger logging.Logger
	// dataManager database.DataManager
}

func (s *oauth2TokenStoreImpl) Create(ctx context.Context, info oauth2.TokenInfo) error {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("info", info)
	logger.Debug("Create invoked")

	return nil
}

func (s *oauth2TokenStoreImpl) RemoveByCode(ctx context.Context, code string) error {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("code", code)
	logger.Debug("RemoveByCode invoked")

	return nil
}

func (s *oauth2TokenStoreImpl) RemoveByAccess(ctx context.Context, access string) error {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("access", access)
	logger.Debug("RemoveByAccess invoked")

	return nil
}

func (s *oauth2TokenStoreImpl) RemoveByRefresh(ctx context.Context, refresh string) error {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("refresh", refresh)
	logger.Debug("RemoveByRefresh invoked")

	return nil
}

func (s *oauth2TokenStoreImpl) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("code", code)
	logger.Debug("GetByCode invoked")

	return nil, nil
}

func (s *oauth2TokenStoreImpl) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("access", access)
	logger.Debug("GetByAccess invoked")

	return nil, nil
}

func (s *oauth2TokenStoreImpl) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue("refresh", refresh)
	logger.Debug("GetByRefresh invoked")

	return nil, nil
}
