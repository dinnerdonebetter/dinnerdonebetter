package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.ValidIngredientPreparationSQLQueryBuilder = (*ValidIngredientPreparationSQLQueryBuilder)(nil)

// ValidIngredientPreparationSQLQueryBuilder is a mocked types.ValidIngredientPreparationSQLQueryBuilder for testing.
type ValidIngredientPreparationSQLQueryBuilder struct {
	mock.Mock
}

// BuildValidIngredientPreparationExistsQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildValidIngredientPreparationExistsQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validIngredientPreparationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidIngredientPreparationQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildGetValidIngredientPreparationQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validIngredientPreparationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllValidIngredientPreparationsCountQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildGetAllValidIngredientPreparationsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfValidIngredientPreparationsQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildGetBatchOfValidIngredientPreparationsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidIngredientPreparationsQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildGetValidIngredientPreparationsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetValidIngredientPreparationsWithIDsQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildGetValidIngredientPreparationsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, limit, ids)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateValidIngredientPreparationQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildCreateValidIngredientPreparationQuery(ctx context.Context, input *types.ValidIngredientPreparationCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForValidIngredientPreparationQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildGetAuditLogEntriesForValidIngredientPreparationQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validIngredientPreparationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateValidIngredientPreparationQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildUpdateValidIngredientPreparationQuery(ctx context.Context, input *types.ValidIngredientPreparation) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveValidIngredientPreparationQuery implements our interface.
func (m *ValidIngredientPreparationSQLQueryBuilder) BuildArchiveValidIngredientPreparationQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, validIngredientPreparationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
