package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.ValidIngredientSQLQueryBuilder = (*ValidIngredientSQLQueryBuilder)(nil)

// ValidIngredientSQLQueryBuilder is a mocked types.ValidIngredientSQLQueryBuilder for testing.
type ValidIngredientSQLQueryBuilder struct {
	mock.Mock
}

// BuildValidIngredientExistsQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildValidIngredientExistsQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validIngredientID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidIngredientIDForNameQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildGetValidIngredientIDForNameQuery(ctx context.Context, validIngredientName string) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validIngredientName)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildSearchForValidIngredientByNameQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildSearchForValidIngredientByNameQuery(ctx context.Context, name string) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, name)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidIngredientQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildGetValidIngredientQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validIngredientID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllValidIngredientsCountQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildGetAllValidIngredientsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfValidIngredientsQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildGetBatchOfValidIngredientsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidIngredientsQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildGetValidIngredientsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidIngredientsWithIDsQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildGetValidIngredientsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, limit, ids)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateValidIngredientQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildCreateValidIngredientQuery(ctx context.Context, input *types.ValidIngredientCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForValidIngredientQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildGetAuditLogEntriesForValidIngredientQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validIngredientID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateValidIngredientQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildUpdateValidIngredientQuery(ctx context.Context, input *types.ValidIngredient) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveValidIngredientQuery implements our interface.
func (m *ValidIngredientSQLQueryBuilder) BuildArchiveValidIngredientQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validIngredientID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
