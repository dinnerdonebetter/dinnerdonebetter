package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.ValidPreparationSQLQueryBuilder = (*ValidPreparationSQLQueryBuilder)(nil)

// ValidPreparationSQLQueryBuilder is a mocked types.ValidPreparationSQLQueryBuilder for testing.
type ValidPreparationSQLQueryBuilder struct {
	mock.Mock
}

// BuildValidPreparationExistsQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildValidPreparationExistsQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validPreparationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidPreparationQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildGetValidPreparationQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validPreparationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllValidPreparationsCountQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildGetAllValidPreparationsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfValidPreparationsQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildGetBatchOfValidPreparationsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidPreparationsQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildGetValidPreparationsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidPreparationsWithIDsQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildGetValidPreparationsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, limit, ids)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateValidPreparationQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildCreateValidPreparationQuery(ctx context.Context, input *types.ValidPreparationCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForValidPreparationQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildGetAuditLogEntriesForValidPreparationQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validPreparationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateValidPreparationQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildUpdateValidPreparationQuery(ctx context.Context, input *types.ValidPreparation) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveValidPreparationQuery implements our interface.
func (m *ValidPreparationSQLQueryBuilder) BuildArchiveValidPreparationQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validPreparationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
