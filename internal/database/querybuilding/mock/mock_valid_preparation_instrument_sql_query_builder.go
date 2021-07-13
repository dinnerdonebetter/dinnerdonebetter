package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.ValidPreparationInstrumentSQLQueryBuilder = (*ValidPreparationInstrumentSQLQueryBuilder)(nil)

// ValidPreparationInstrumentSQLQueryBuilder is a mocked types.ValidPreparationInstrumentSQLQueryBuilder for testing.
type ValidPreparationInstrumentSQLQueryBuilder struct {
	mock.Mock
}

// BuildValidPreparationInstrumentExistsQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildValidPreparationInstrumentExistsQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validPreparationInstrumentID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidPreparationInstrumentQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildGetValidPreparationInstrumentQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validPreparationInstrumentID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllValidPreparationInstrumentsCountQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildGetAllValidPreparationInstrumentsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfValidPreparationInstrumentsQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildGetBatchOfValidPreparationInstrumentsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidPreparationInstrumentsQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildGetValidPreparationInstrumentsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidPreparationInstrumentsWithIDsQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildGetValidPreparationInstrumentsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, limit, ids)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateValidPreparationInstrumentQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildCreateValidPreparationInstrumentQuery(ctx context.Context, input *types.ValidPreparationInstrumentCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForValidPreparationInstrumentQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildGetAuditLogEntriesForValidPreparationInstrumentQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validPreparationInstrumentID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateValidPreparationInstrumentQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildUpdateValidPreparationInstrumentQuery(ctx context.Context, input *types.ValidPreparationInstrument) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveValidPreparationInstrumentQuery implements our interface.
func (m *ValidPreparationInstrumentSQLQueryBuilder) BuildArchiveValidPreparationInstrumentQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validPreparationInstrumentID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
