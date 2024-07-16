package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.AuditLogEntryDataManager = (*AuditLogEntryDataManagerMock)(nil)

// AuditLogEntryDataManagerMock is a mocked types.AuditLogEntryDataManager for testing.
type AuditLogEntryDataManagerMock struct {
	mock.Mock
}

func (m *AuditLogEntryDataManagerMock) GetAuditLogEntry(ctx context.Context, auditLogID string) (*types.AuditLogEntry, error) {
	args := m.Called(ctx, auditLogID)
	return args.Get(0).(*types.AuditLogEntry), args.Error(1)
}

func (m *AuditLogEntryDataManagerMock) GetAuditLogEntriesForUser(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.AuditLogEntry]), args.Error(1)
}

func (m *AuditLogEntryDataManagerMock) GetAuditLogEntriesForUserAndResourceType(ctx context.Context, userID string, resourceTypes []string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	args := m.Called(ctx, userID, resourceTypes, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.AuditLogEntry]), args.Error(1)
}

func (m *AuditLogEntryDataManagerMock) GetAuditLogEntriesForHousehold(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	args := m.Called(ctx, householdID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.AuditLogEntry]), args.Error(1)
}

func (m *AuditLogEntryDataManagerMock) GetAuditLogEntriesForHouseholdAndResourceType(ctx context.Context, householdID string, resourceTypes []string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	args := m.Called(ctx, householdID, resourceTypes, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.AuditLogEntry]), args.Error(1)
}

// CreateAuditLogEntry is a mock function.
func (m *AuditLogEntryDataManagerMock) CreateAuditLogEntry(ctx context.Context, input *types.AuditLogEntryDatabaseCreationInput) (*types.AuditLogEntry, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.AuditLogEntry), args.Error(1)
}
