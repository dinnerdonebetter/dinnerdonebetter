package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.ReportDataManager = (*ReportDataManager)(nil)

// ReportDataManager is a mocked models.ReportDataManager for testing
type ReportDataManager struct {
	mock.Mock
}

// GetReport is a mock function
func (m *ReportDataManager) GetReport(ctx context.Context, reportID, userID uint64) (*models.Report, error) {
	args := m.Called(ctx, reportID, userID)
	return args.Get(0).(*models.Report), args.Error(1)
}

// GetReportCount is a mock function
func (m *ReportDataManager) GetReportCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllReportsCount is a mock function
func (m *ReportDataManager) GetAllReportsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetReports is a mock function
func (m *ReportDataManager) GetReports(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.ReportList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.ReportList), args.Error(1)
}

// GetAllReportsForUser is a mock function
func (m *ReportDataManager) GetAllReportsForUser(ctx context.Context, userID uint64) ([]models.Report, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Report), args.Error(1)
}

// CreateReport is a mock function
func (m *ReportDataManager) CreateReport(ctx context.Context, input *models.ReportCreationInput) (*models.Report, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.Report), args.Error(1)
}

// UpdateReport is a mock function
func (m *ReportDataManager) UpdateReport(ctx context.Context, updated *models.Report) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveReport is a mock function
func (m *ReportDataManager) ArchiveReport(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
