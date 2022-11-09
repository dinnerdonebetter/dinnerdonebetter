package mocksearch

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/internal/search"
)

var _ search.IndexManager = (*IndexManager)(nil)

// IndexManager is a mock IndexManager.
type IndexManager struct {
	mock.Mock
}

// Index implements our interface.
func (m *IndexManager) Index(ctx context.Context, id string, value any) error {
	args := m.Called(ctx, id, value)
	return args.Error(0)
}

// Search implements our interface.
func (m *IndexManager) Search(ctx context.Context, byField, query, householdID string) (ids []string, err error) {
	args := m.Called(ctx, byField, query, householdID)
	return args.Get(0).([]string), args.Error(1)
}

// Delete implements our interface.
func (m *IndexManager) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
