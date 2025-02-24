package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidInstrumentDataManager = (*ValidInstrumentDataManagerMock)(nil)

// ValidInstrumentDataManagerMock is a mocked types.ValidInstrumentDataManager for testing.
type ValidInstrumentDataManagerMock struct {
	mock.Mock
}

// ValidInstrumentExists is a mock function.
func (m *ValidInstrumentDataManagerMock) ValidInstrumentExists(ctx context.Context, validInstrumentID string) (bool, error) {
	returnValues := m.Called(ctx, validInstrumentID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	returnValues := m.Called(ctx, validInstrumentID)
	return returnValues.Get(0).(*types.ValidInstrument), returnValues.Error(1)
}

// GetRandomValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) GetRandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(*types.ValidInstrument), returnValues.Error(1)
}

// SearchForValidInstruments is a mock function.
func (m *ValidInstrumentDataManagerMock) SearchForValidInstruments(ctx context.Context, query string) ([]*types.ValidInstrument, error) {
	returnValues := m.Called(ctx, query)
	return returnValues.Get(0).([]*types.ValidInstrument), returnValues.Error(1)
}

// GetValidInstruments is a mock function.
func (m *ValidInstrumentDataManagerMock) GetValidInstruments(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidInstrument], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidInstrument]), returnValues.Error(1)
}

// CreateValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentDatabaseCreationInput) (*types.ValidInstrument, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ValidInstrument), returnValues.Error(1)
}

// UpdateValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	return m.Called(ctx, validInstrumentID).Error(0)
}

// MarkValidInstrumentAsIndexed is a mock function.
func (m *ValidInstrumentDataManagerMock) MarkValidInstrumentAsIndexed(ctx context.Context, validInstrumentID string) error {
	return m.Called(ctx, validInstrumentID).Error(0)
}

// GetValidInstrumentIDsThatNeedSearchIndexing is a mock function.
func (m *ValidInstrumentDataManagerMock) GetValidInstrumentIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidInstrumentsWithIDs is a mock function.
func (m *ValidInstrumentDataManagerMock) GetValidInstrumentsWithIDs(ctx context.Context, ids []string) ([]*types.ValidInstrument, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*types.ValidInstrument), returnValues.Error(1)
}
