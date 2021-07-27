package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.ValidInstrumentSQLQueryBuilder = (*ValidInstrumentSQLQueryBuilder)(nil)

// ValidInstrumentSQLQueryBuilder is a mocked types.ValidInstrumentSQLQueryBuilder for testing.
type ValidInstrumentSQLQueryBuilder struct {
	mock.Mock
}

// BuildValidInstrumentExistsQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildValidInstrumentExistsQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validInstrumentID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidInstrumentIDForNameQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildGetValidInstrumentIDForNameQuery(ctx context.Context, validInstrumentName string) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validInstrumentName)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildSearchForValidInstrumentByNameQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildSearchForValidInstrumentByNameQuery(ctx context.Context, name string) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, name)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidInstrumentQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildGetValidInstrumentQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validInstrumentID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllValidInstrumentsCountQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildGetAllValidInstrumentsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfValidInstrumentsQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildGetBatchOfValidInstrumentsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidInstrumentsQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildGetValidInstrumentsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidInstrumentsWithIDsQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildGetValidInstrumentsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, limit, ids)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateValidInstrumentQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildCreateValidInstrumentQuery(ctx context.Context, input *types.ValidInstrumentCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForValidInstrumentQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildGetAuditLogEntriesForValidInstrumentQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validInstrumentID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateValidInstrumentQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildUpdateValidInstrumentQuery(ctx context.Context, input *types.ValidInstrument) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveValidInstrumentQuery implements our interface.
func (m *ValidInstrumentSQLQueryBuilder) BuildArchiveValidInstrumentQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validInstrumentID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
