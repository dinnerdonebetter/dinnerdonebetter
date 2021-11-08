package database

import (
	"context"
	"database/sql"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/mock"

	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
)

var _ DataManager = (*MockDatabase)(nil)

// BuildMockDatabase builds a mock database.
func BuildMockDatabase() *MockDatabase {
	return &MockDatabase{
		HouseholdDataManager:                  &mocktypes.HouseholdDataManager{},
		HouseholdInvitationDataManager:        &mocktypes.HouseholdInvitationDataManager{},
		HouseholdUserMembershipDataManager:    &mocktypes.HouseholdUserMembershipDataManager{},
		ValidInstrumentDataManager:            &mocktypes.ValidInstrumentDataManager{},
		ValidIngredientDataManager:            &mocktypes.ValidIngredientDataManager{},
		ValidPreparationDataManager:           &mocktypes.ValidPreparationDataManager{},
		ValidIngredientPreparationDataManager: &mocktypes.ValidIngredientPreparationDataManager{},
		RecipeDataManager:                     &mocktypes.RecipeDataManager{},
		RecipeStepDataManager:                 &mocktypes.RecipeStepDataManager{},
		RecipeStepInstrumentDataManager:       &mocktypes.RecipeStepInstrumentDataManager{},
		RecipeStepIngredientDataManager:       &mocktypes.RecipeStepIngredientDataManager{},
		RecipeStepProductDataManager:          &mocktypes.RecipeStepProductDataManager{},
		MealPlanDataManager:                   &mocktypes.MealPlanDataManager{},
		MealPlanOptionDataManager:             &mocktypes.MealPlanOptionDataManager{},
		MealPlanOptionVoteDataManager:         &mocktypes.MealPlanOptionVoteDataManager{},
		UserDataManager:                       &mocktypes.UserDataManager{},
		AdminUserDataManager:                  &mocktypes.AdminUserDataManager{},
		APIClientDataManager:                  &mocktypes.APIClientDataManager{},
		WebhookDataManager:                    &mocktypes.WebhookDataManager{},
	}
}

// MockDatabase is our mock database structure. Note, when using this in tests, you must directly access the type name of all the implicit fields.
// So `mockDB.On("GetUserByUsername"...)` is destined to fail, whereas `mockDB.UserDataManager.On("GetUserByUsername"...)` would do what you want it to do.
type MockDatabase struct {
	*mocktypes.AdminUserDataManager
	*mocktypes.HouseholdUserMembershipDataManager
	*mocktypes.ValidInstrumentDataManager
	*mocktypes.ValidIngredientDataManager
	*mocktypes.ValidPreparationDataManager
	*mocktypes.ValidIngredientPreparationDataManager
	*mocktypes.RecipeDataManager
	*mocktypes.RecipeStepDataManager
	*mocktypes.RecipeStepInstrumentDataManager
	*mocktypes.RecipeStepIngredientDataManager
	*mocktypes.RecipeStepProductDataManager
	*mocktypes.MealPlanDataManager
	*mocktypes.MealPlanOptionDataManager
	*mocktypes.MealPlanOptionVoteDataManager
	*mocktypes.UserDataManager
	*mocktypes.APIClientDataManager
	*mocktypes.WebhookDataManager
	*mocktypes.HouseholdDataManager
	*mocktypes.HouseholdInvitationDataManager
	mock.Mock
}

// ProvideSessionStore satisfies the DataManager interface.
func (m *MockDatabase) ProvideSessionStore() scs.Store {
	return m.Called().Get(0).(scs.Store)
}

// Migrate satisfies the DataManager interface.
func (m *MockDatabase) Migrate(ctx context.Context, maxAttempts uint8) error {
	return m.Called(ctx, maxAttempts).Error(0)
}

// IsReady satisfies the DataManager interface.
func (m *MockDatabase) IsReady(ctx context.Context, maxAttempts uint8) (ready bool) {
	return m.Called(ctx, maxAttempts).Bool(0)
}

// BeginTx satisfies the DataManager interface.
func (m *MockDatabase) BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

var _ ResultIterator = (*MockResultIterator)(nil)

// MockResultIterator is our mock sql.Rows structure.
type MockResultIterator struct {
	mock.Mock
}

// Scan satisfies the ResultIterator interface.
func (m *MockResultIterator) Scan(dest ...interface{}) error {
	return m.Called(dest...).Error(0)
}

// Next satisfies the ResultIterator interface.
func (m *MockResultIterator) Next() bool {
	return m.Called().Bool(0)
}

// Err satisfies the ResultIterator interface.
func (m *MockResultIterator) Err() error {
	return m.Called().Error(0)
}

// Close satisfies the ResultIterator interface.
func (m *MockResultIterator) Close() error {
	return m.Called().Error(0)
}

// MockSQLResult mocks a sql.Result.
type MockSQLResult struct {
	mock.Mock
}

// LastInsertId implements our interface.
func (m *MockSQLResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// RowsAffected implements our interface.
func (m *MockSQLResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

var _ SQLQueryExecutor = (*MockQueryExecutor)(nil)

// MockQueryExecutor mocks a sql.Tx|DB.
type MockQueryExecutor struct {
	mock.Mock
}

// ExecContext is a mock function.
func (m *MockQueryExecutor) ExecContext(ctx context.Context, query string, queryArgs ...interface{}) (sql.Result, error) {
	args := m.Called(ctx, query, queryArgs)
	return args.Get(0).(sql.Result), args.Error(1)
}

// QueryContext is a mock function.
func (m *MockQueryExecutor) QueryContext(ctx context.Context, query string, queryArgs ...interface{}) (*sql.Rows, error) {
	args := m.Called(ctx, query, queryArgs)
	return args.Get(0).(*sql.Rows), args.Error(1)
}

// QueryRowContext is a mock function.
func (m *MockQueryExecutor) QueryRowContext(ctx context.Context, query string, queryArgs ...interface{}) *sql.Row {
	args := m.Called(ctx, query, queryArgs)
	return args.Get(0).(*sql.Row)
}
