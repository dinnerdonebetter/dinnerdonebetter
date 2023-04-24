package database

import (
	"context"
	"database/sql"

	mocktypes "github.com/prixfixeco/backend/pkg/types/mock"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/mock"
)

var _ DataManager = (*MockDatabase)(nil)

// NewMockDatabase builds a mock database.
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		HouseholdDataManager:                      &mocktypes.HouseholdDataManager{},
		HouseholdInvitationDataManager:            &mocktypes.HouseholdInvitationDataManager{},
		HouseholdUserMembershipDataManager:        &mocktypes.HouseholdUserMembershipDataManager{},
		ValidInstrumentDataManager:                &mocktypes.ValidInstrumentDataManager{},
		ValidIngredientDataManager:                &mocktypes.ValidIngredientDataManager{},
		ValidPreparationDataManager:               &mocktypes.ValidPreparationDataManager{},
		ValidIngredientPreparationDataManager:     &mocktypes.ValidIngredientPreparationDataManager{},
		MealDataManager:                           &mocktypes.MealDataManager{},
		RecipeDataManager:                         &mocktypes.RecipeDataManager{},
		RecipeStepDataManager:                     &mocktypes.RecipeStepDataManager{},
		RecipeStepProductDataManager:              &mocktypes.RecipeStepProductDataManager{},
		RecipeStepInstrumentDataManager:           &mocktypes.RecipeStepInstrumentDataManager{},
		RecipeStepIngredientDataManager:           &mocktypes.RecipeStepIngredientDataManager{},
		MealPlanDataManager:                       &mocktypes.MealPlanDataManager{},
		MealPlanOptionDataManager:                 &mocktypes.MealPlanOptionDataManager{},
		MealPlanOptionVoteDataManager:             &mocktypes.MealPlanOptionVoteDataManager{},
		UserDataManager:                           &mocktypes.UserDataManager{},
		AdminUserDataManager:                      &mocktypes.AdminUserDataManager{},
		APIClientDataManager:                      &mocktypes.APIClientDataManager{},
		PasswordResetTokenDataManager:             &mocktypes.PasswordResetTokenDataManager{},
		WebhookDataManager:                        &mocktypes.WebhookDataManager{},
		ValidMeasurementUnitDataManager:           &mocktypes.ValidMeasurementUnitDataManager{},
		ValidPreparationInstrumentDataManager:     &mocktypes.ValidPreparationInstrumentDataManager{},
		ValidIngredientMeasurementUnitDataManager: &mocktypes.ValidIngredientMeasurementUnitDataManager{},
		MealPlanEventDataManager:                  &mocktypes.MealPlanEventDataManager{},
		MealPlanTaskDataManager:                   &mocktypes.MealPlanTaskDataManager{},
		RecipePrepTaskDataManager:                 &mocktypes.RecipePrepTaskDataManager{},
		MealPlanGroceryListItemDataManager:        &mocktypes.MealPlanGroceryListItemDataManager{},
		ValidMeasurementConversionDataManager:     &mocktypes.ValidMeasurementConversionDataManager{},
		RecipeMediaDataManager:                    &mocktypes.RecipeMediaDataManager{},
		ValidIngredientStateDataManager:           &mocktypes.ValidIngredientStateDataManager{},
		ValidIngredientStateIngredientDataManager: &mocktypes.ValidIngredientStateIngredientDataManager{},
		RecipeStepCompletionConditionDataManager:  &mocktypes.RecipeStepCompletionConditionDataManager{},
		RecipeStepVesselDataManager:               &mocktypes.RecipeStepVesselDataManager{},
		ServiceSettingDataManager:                 &mocktypes.ServiceSettingDataManager{},
	}
}

// MockDatabase is our mock database structure. Note, when using this in tests, you must directly access the type name of all the implicit fields.
// So `mockDB.On("GetUserByUsername"...)` is destined to fail, whereas `mockDB.UserDataManager.On("GetUserByUsername"...)` would do what you want it to do.
type MockDatabase struct {
	*mocktypes.ValidIngredientStateDataManager
	*mocktypes.AdminUserDataManager
	*mocktypes.HouseholdUserMembershipDataManager
	*mocktypes.ValidInstrumentDataManager
	*mocktypes.ValidIngredientDataManager
	*mocktypes.ValidPreparationDataManager
	*mocktypes.ValidIngredientPreparationDataManager
	*mocktypes.MealDataManager
	*mocktypes.RecipeDataManager
	*mocktypes.RecipeStepDataManager
	*mocktypes.RecipeStepProductDataManager
	*mocktypes.RecipeStepInstrumentDataManager
	*mocktypes.RecipeStepIngredientDataManager
	*mocktypes.MealPlanDataManager
	*mocktypes.MealPlanOptionDataManager
	*mocktypes.MealPlanOptionVoteDataManager
	*mocktypes.UserDataManager
	*mocktypes.APIClientDataManager
	*mocktypes.PasswordResetTokenDataManager
	*mocktypes.WebhookDataManager
	*mocktypes.HouseholdDataManager
	*mocktypes.HouseholdInvitationDataManager
	*mocktypes.ValidMeasurementUnitDataManager
	*mocktypes.ValidPreparationInstrumentDataManager
	*mocktypes.ValidIngredientMeasurementUnitDataManager
	*mocktypes.MealPlanEventDataManager
	*mocktypes.MealPlanTaskDataManager
	*mocktypes.RecipePrepTaskDataManager
	*mocktypes.MealPlanGroceryListItemDataManager
	*mocktypes.ValidMeasurementConversionDataManager
	*mocktypes.RecipeMediaDataManager
	*mocktypes.RecipeStepCompletionConditionDataManager
	*mocktypes.ValidIngredientStateIngredientDataManager
	*mocktypes.RecipeStepVesselDataManager
	*mocktypes.ServiceSettingDataManager

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

// Close satisfies the DataManager interface.
func (m *MockDatabase) Close() {
	m.Called()
}

// DB satisfies the DataManager interface.
func (m *MockDatabase) DB() *sql.DB {
	return m.Called().Get(0).(*sql.DB)
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
func (m *MockResultIterator) Scan(dest ...any) error {
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
func (m *MockQueryExecutor) ExecContext(ctx context.Context, query string, queryArgs ...any) (sql.Result, error) {
	args := m.Called(ctx, query, queryArgs)
	return args.Get(0).(sql.Result), args.Error(1)
}

// PrepareContext is a mock function.
func (m *MockQueryExecutor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*sql.Stmt), args.Error(1)
}

// QueryContext is a mock function.
func (m *MockQueryExecutor) QueryContext(ctx context.Context, query string, queryArgs ...any) (*sql.Rows, error) {
	args := m.Called(ctx, query, queryArgs)
	return args.Get(0).(*sql.Rows), args.Error(1)
}

// QueryRowContext is a mock function.
func (m *MockQueryExecutor) QueryRowContext(ctx context.Context, query string, queryArgs ...any) *sql.Row {
	args := m.Called(ctx, query, queryArgs)
	return args.Get(0).(*sql.Row)
}
