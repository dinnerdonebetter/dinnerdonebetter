package mocksearch

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/search"
)

var _ search.IndexManagerProvider = (*IndexManagerProvider)(nil)

// IndexManagerProvider is a mock IndexManagerProvider.
type IndexManagerProvider struct {
	mock.Mock
}

// ProvideIndexManager implements our interface.
func (m *IndexManagerProvider) ProvideIndexManager(ctx context.Context, logger logging.Logger, name search.IndexName, fields ...string) (search.IndexManager, error) {
	args := m.Called(ctx, logger, name, fields)
	return args.Get(0).(search.IndexManager), args.Error(1)
}
