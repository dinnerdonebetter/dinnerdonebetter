package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.RecipeStepSQLQueryBuilder = (*RecipeStepSQLQueryBuilder)(nil)

// RecipeStepSQLQueryBuilder is a mocked types.RecipeStepSQLQueryBuilder for testing.
type RecipeStepSQLQueryBuilder struct {
	mock.Mock
}

// BuildRecipeStepExistsQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildRecipeStepExistsQuery(ctx context.Context, recipeID, recipeStepID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, recipeStepID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeStepQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildGetRecipeStepQuery(ctx context.Context, recipeID, recipeStepID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, recipeStepID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllRecipeStepsCountQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildGetAllRecipeStepsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfRecipeStepsQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildGetBatchOfRecipeStepsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeStepsQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildGetRecipeStepsQuery(ctx context.Context, recipeID uint64, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeStepsWithIDsQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildGetRecipeStepsWithIDsQuery(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, limit, ids)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateRecipeStepQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildCreateRecipeStepQuery(ctx context.Context, input *types.RecipeStepCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForRecipeStepQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildGetAuditLogEntriesForRecipeStepQuery(ctx context.Context, recipeStepID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeStepID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateRecipeStepQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildUpdateRecipeStepQuery(ctx context.Context, input *types.RecipeStep) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveRecipeStepQuery implements our interface.
func (m *RecipeStepSQLQueryBuilder) BuildArchiveRecipeStepQuery(ctx context.Context, recipeID, recipeStepID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, recipeStepID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
