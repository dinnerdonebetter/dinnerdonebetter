package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.RecipeStepIngredientSQLQueryBuilder = (*RecipeStepIngredientSQLQueryBuilder)(nil)

// RecipeStepIngredientSQLQueryBuilder is a mocked types.RecipeStepIngredientSQLQueryBuilder for testing.
type RecipeStepIngredientSQLQueryBuilder struct {
	mock.Mock
}

// BuildRecipeStepIngredientExistsQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildRecipeStepIngredientExistsQuery(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeStepIngredientQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildGetRecipeStepIngredientQuery(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllRecipeStepIngredientsCountQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildGetAllRecipeStepIngredientsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfRecipeStepIngredientsQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildGetBatchOfRecipeStepIngredientsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeStepIngredientsQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildGetRecipeStepIngredientsQuery(ctx context.Context, recipeID, recipeStepID uint64, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeID, recipeStepID, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetRecipeStepIngredientsWithIDsQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildGetRecipeStepIngredientsWithIDsQuery(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeStepID, limit, ids)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateRecipeStepIngredientQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildCreateRecipeStepIngredientQuery(ctx context.Context, input *types.RecipeStepIngredientCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForRecipeStepIngredientQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildGetAuditLogEntriesForRecipeStepIngredientQuery(ctx context.Context, recipeStepIngredientID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeStepIngredientID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateRecipeStepIngredientQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildUpdateRecipeStepIngredientQuery(ctx context.Context, input *types.RecipeStepIngredient) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveRecipeStepIngredientQuery implements our interface.
func (m *RecipeStepIngredientSQLQueryBuilder) BuildArchiveRecipeStepIngredientQuery(ctx context.Context, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, recipeStepID, recipeStepIngredientID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
