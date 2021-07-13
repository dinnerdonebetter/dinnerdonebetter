package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.RecipeStepProductSQLQueryBuilder = (*RecipeStepProductSQLQueryBuilder)(nil)

// RecipeStepProductSQLQueryBuilder is a mocked types.RecipeStepProductSQLQueryBuilder for testing.
type RecipeStepProductSQLQueryBuilder struct {
	mock.Mock
}

// BuildRecipeStepProductExistsQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildRecipeStepProductExistsQuery(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeStepProductQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildGetRecipeStepProductQuery(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllRecipeStepProductsCountQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildGetAllRecipeStepProductsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfRecipeStepProductsQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildGetBatchOfRecipeStepProductsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeStepProductsQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildGetRecipeStepProductsQuery(ctx context.Context, recipeID, recipeStepID uint64, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, recipeStepID, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeStepProductsWithIDsQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildGetRecipeStepProductsWithIDsQuery(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeStepID, limit, ids)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateRecipeStepProductQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildCreateRecipeStepProductQuery(ctx context.Context, input *types.RecipeStepProductCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForRecipeStepProductQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildGetAuditLogEntriesForRecipeStepProductQuery(ctx context.Context, recipeStepProductID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeStepProductID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateRecipeStepProductQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildUpdateRecipeStepProductQuery(ctx context.Context, input *types.RecipeStepProduct) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveRecipeStepProductQuery implements our interface.
func (m *RecipeStepProductSQLQueryBuilder) BuildArchiveRecipeStepProductQuery(ctx context.Context, recipeStepID, recipeStepProductID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeStepID, recipeStepProductID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
