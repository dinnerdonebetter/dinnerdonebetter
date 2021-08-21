package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ReportDataManager = (*ReportDataManager)(nil)

// ReportDataManager is a mocked types.ReportDataManager for testing.
type ReportDataManager struct {
	mock.Mock
}

// ReportExists is a mock function.
func (m *ReportDataManager) ReportExists(ctx context.Context, reportID uint64) (bool, error) {
	args := m.Called(ctx, reportID)
	return args.Bool(0), args.Error(1)
}

// GetReport is a mock function.
func (m *ReportDataManager) GetReport(ctx context.Context, reportID uint64) (*types.Report, error) {
	args := m.Called(ctx, reportID)
	return args.Get(0).(*types.Report), args.Error(1)
}

// GetAllReportsCount is a mock function.
func (m *ReportDataManager) GetAllReportsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllReports is a mock function.
func (m *ReportDataManager) GetAllReports(ctx context.Context, results chan []*types.Report, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetReports is a mock function.
func (m *ReportDataManager) GetReports(ctx context.Context, filter *types.QueryFilter) (*types.ReportList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ReportList), args.Error(1)
}

// GetReportsWithIDs is a mock function.
func (m *ReportDataManager) GetReportsWithIDs(ctx context.Context, householdID uint64, limit uint8, ids []uint64) ([]*types.Report, error) {
	args := m.Called(ctx, householdID, limit, ids)
	return args.Get(0).([]*types.Report), args.Error(1)
}

// CreateReport is a mock function.
func (m *ReportDataManager) CreateReport(ctx context.Context, input *types.ReportCreationInput, createdByUser uint64) (*types.Report, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.Report), args.Error(1)
}

// UpdateReport is a mock function.
func (m *ReportDataManager) UpdateReport(ctx context.Context, updated *types.Report, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveReport is a mock function.
func (m *ReportDataManager) ArchiveReport(ctx context.Context, reportID, householdID, archivedBy uint64) error {
	return m.Called(ctx, reportID, householdID, archivedBy).Error(0)
}

// GetAuditLogEntriesForReport is a mock function.
func (m *ReportDataManager) GetAuditLogEntriesForReport(ctx context.Context, reportID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, reportID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
