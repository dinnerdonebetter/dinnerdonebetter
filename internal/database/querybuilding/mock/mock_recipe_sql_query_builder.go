package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.RecipeSQLQueryBuilder = (*RecipeSQLQueryBuilder)(nil)

// RecipeSQLQueryBuilder is a mocked types.RecipeSQLQueryBuilder for testing.
type RecipeSQLQueryBuilder struct {
	mock.Mock
}

// BuildRecipeExistsQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildRecipeExistsQuery(ctx context.Context, recipeID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildGetRecipeQuery(ctx context.Context, recipeID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllRecipesCountQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildGetAllRecipesCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfRecipesQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildGetBatchOfRecipesQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipesQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildGetRecipesQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipesWithIDsQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildGetRecipesWithIDsQuery(ctx context.Context, accountID uint64, limit uint8, ids []uint64, restrictToAccount bool) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, accountID, limit, ids, restrictToAccount)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateRecipeQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildCreateRecipeQuery(ctx context.Context, input *types.RecipeCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForRecipeQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildGetAuditLogEntriesForRecipeQuery(ctx context.Context, recipeID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateRecipeQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildUpdateRecipeQuery(ctx context.Context, input *types.Recipe) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveRecipeQuery implements our interface.
func (m *RecipeSQLQueryBuilder) BuildArchiveRecipeQuery(ctx context.Context, recipeID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
