package oauth2

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/go-oauth2/oauth2/v4"
)

type oauth2TokenStoreImpl struct {
	tracer      tracing.Tracer
	logger      logging.Logger
	dataManager database.DataManager
}

func (s *oauth2TokenStoreImpl) Create(ctx context.Context, info oauth2.TokenInfo) error {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (s *oauth2TokenStoreImpl) RemoveByCode(ctx context.Context, code string) error {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (s *oauth2TokenStoreImpl) RemoveByAccess(ctx context.Context, access string) error {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (s *oauth2TokenStoreImpl) RemoveByRefresh(ctx context.Context, refresh string) error {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (s *oauth2TokenStoreImpl) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, nil
}

func (s *oauth2TokenStoreImpl) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, nil
}

func (s *oauth2TokenStoreImpl) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, nil
}
