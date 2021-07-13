package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.ReportSQLQueryBuilder = (*ReportSQLQueryBuilder)(nil)

// ReportSQLQueryBuilder is a mocked types.ReportSQLQueryBuilder for testing.
type ReportSQLQueryBuilder struct {
	mock.Mock
}

// BuildReportExistsQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildReportExistsQuery(ctx context.Context, reportID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, reportID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetReportQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildGetReportQuery(ctx context.Context, reportID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, reportID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllReportsCountQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildGetAllReportsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfReportsQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildGetBatchOfReportsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetReportsQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildGetReportsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetReportsWithIDsQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildGetReportsWithIDsQuery(ctx context.Context, accountID uint64, limit uint8, ids []uint64, restrictToAccount bool) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, accountID, limit, ids, restrictToAccount)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateReportQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildCreateReportQuery(ctx context.Context, input *types.ReportCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForReportQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildGetAuditLogEntriesForReportQuery(ctx context.Context, reportID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, reportID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateReportQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildUpdateReportQuery(ctx context.Context, input *types.Report) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveReportQuery implements our interface.
func (m *ReportSQLQueryBuilder) BuildArchiveReportQuery(ctx context.Context, reportID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, reportID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
