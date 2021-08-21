package database

import (
	"context"
	"database/sql"

	"gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	mockquerybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding/mock"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"

	"github.com/stretchr/testify/mock"
)

var _ DataManager = (*MockDatabase)(nil)

// BuildMockDatabase builds a mock database.
func BuildMockDatabase() *MockDatabase {
	return &MockDatabase{
		AuditLogEntryDataManager:              &mocktypes.AuditLogEntryDataManager{},
		HouseholdDataManager:                  &mocktypes.HouseholdDataManager{},
		HouseholdUserMembershipDataManager:    &mocktypes.HouseholdUserMembershipDataManager{},
		ValidInstrumentDataManager:            &mocktypes.ValidInstrumentDataManager{},
		ValidPreparationDataManager:           &mocktypes.ValidPreparationDataManager{},
		ValidIngredientDataManager:            &mocktypes.ValidIngredientDataManager{},
		ValidIngredientPreparationDataManager: &mocktypes.ValidIngredientPreparationDataManager{},
		ValidPreparationInstrumentDataManager: &mocktypes.ValidPreparationInstrumentDataManager{},
		RecipeDataManager:                     &mocktypes.RecipeDataManager{},
		RecipeStepDataManager:                 &mocktypes.RecipeStepDataManager{},
		RecipeStepIngredientDataManager:       &mocktypes.RecipeStepIngredientDataManager{},
		RecipeStepProductDataManager:          &mocktypes.RecipeStepProductDataManager{},
		InvitationDataManager:                 &mocktypes.InvitationDataManager{},
		ReportDataManager:                     &mocktypes.ReportDataManager{},
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
	*mocktypes.AuditLogEntryDataManager
	*mocktypes.HouseholdUserMembershipDataManager
	*mocktypes.ValidInstrumentDataManager
	*mocktypes.ValidPreparationDataManager
	*mocktypes.ValidIngredientDataManager
	*mocktypes.ValidIngredientPreparationDataManager
	*mocktypes.ValidPreparationInstrumentDataManager
	*mocktypes.RecipeDataManager
	*mocktypes.RecipeStepDataManager
	*mocktypes.RecipeStepIngredientDataManager
	*mocktypes.RecipeStepProductDataManager
	*mocktypes.InvitationDataManager
	*mocktypes.ReportDataManager
	*mocktypes.UserDataManager
	*mocktypes.APIClientDataManager
	*mocktypes.WebhookDataManager
	*mocktypes.HouseholdDataManager
	mock.Mock
}

// Migrate satisfies the DataManager interface.
func (m *MockDatabase) Migrate(ctx context.Context, maxAttempts uint8, ucc *types.TestUserCreationConfig) error {
	return m.Called(ctx, maxAttempts, ucc).Error(0)
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

var _ querybuilding.SQLQueryBuilder = (*MockSQLQueryBuilder)(nil)

// BuildMockSQLQueryBuilder builds a MockSQLQueryBuilder.
func BuildMockSQLQueryBuilder() *MockSQLQueryBuilder {
	return &MockSQLQueryBuilder{
		HouseholdSQLQueryBuilder:                  &mockquerybuilding.HouseholdSQLQueryBuilder{},
		HouseholdUserMembershipSQLQueryBuilder:    &mockquerybuilding.HouseholdUserMembershipSQLQueryBuilder{},
		AuditLogEntrySQLQueryBuilder:              &mockquerybuilding.AuditLogEntrySQLQueryBuilder{},
		ValidInstrumentSQLQueryBuilder:            &mockquerybuilding.ValidInstrumentSQLQueryBuilder{},
		ValidPreparationSQLQueryBuilder:           &mockquerybuilding.ValidPreparationSQLQueryBuilder{},
		ValidIngredientSQLQueryBuilder:            &mockquerybuilding.ValidIngredientSQLQueryBuilder{},
		ValidIngredientPreparationSQLQueryBuilder: &mockquerybuilding.ValidIngredientPreparationSQLQueryBuilder{},
		ValidPreparationInstrumentSQLQueryBuilder: &mockquerybuilding.ValidPreparationInstrumentSQLQueryBuilder{},
		RecipeSQLQueryBuilder:                     &mockquerybuilding.RecipeSQLQueryBuilder{},
		RecipeStepSQLQueryBuilder:                 &mockquerybuilding.RecipeStepSQLQueryBuilder{},
		RecipeStepIngredientSQLQueryBuilder:       &mockquerybuilding.RecipeStepIngredientSQLQueryBuilder{},
		RecipeStepProductSQLQueryBuilder:          &mockquerybuilding.RecipeStepProductSQLQueryBuilder{},
		InvitationSQLQueryBuilder:                 &mockquerybuilding.InvitationSQLQueryBuilder{},
		ReportSQLQueryBuilder:                     &mockquerybuilding.ReportSQLQueryBuilder{},
		APIClientSQLQueryBuilder:                  &mockquerybuilding.APIClientSQLQueryBuilder{},
		UserSQLQueryBuilder:                       &mockquerybuilding.UserSQLQueryBuilder{},
		WebhookSQLQueryBuilder:                    &mockquerybuilding.WebhookSQLQueryBuilder{},
	}
}

// MockSQLQueryBuilder is our mock database structure.
type MockSQLQueryBuilder struct {
	*mockquerybuilding.UserSQLQueryBuilder
	*mockquerybuilding.HouseholdSQLQueryBuilder
	*mockquerybuilding.HouseholdUserMembershipSQLQueryBuilder
	*mockquerybuilding.AuditLogEntrySQLQueryBuilder
	*mockquerybuilding.ValidInstrumentSQLQueryBuilder
	*mockquerybuilding.ValidPreparationSQLQueryBuilder
	*mockquerybuilding.ValidIngredientSQLQueryBuilder
	*mockquerybuilding.ValidIngredientPreparationSQLQueryBuilder
	*mockquerybuilding.ValidPreparationInstrumentSQLQueryBuilder
	*mockquerybuilding.RecipeSQLQueryBuilder
	*mockquerybuilding.RecipeStepSQLQueryBuilder
	*mockquerybuilding.RecipeStepIngredientSQLQueryBuilder
	*mockquerybuilding.RecipeStepProductSQLQueryBuilder
	*mockquerybuilding.InvitationSQLQueryBuilder
	*mockquerybuilding.ReportSQLQueryBuilder
	*mockquerybuilding.APIClientSQLQueryBuilder
	*mockquerybuilding.WebhookSQLQueryBuilder
	mock.Mock
}

// BuildMigrationFunc implements our interface.
func (m *MockSQLQueryBuilder) BuildMigrationFunc(db *sql.DB) func() {
	args := m.Called(db)

	return args.Get(0).(func())
}

// BuildTestUserCreationQuery implements our interface.
func (m *MockSQLQueryBuilder) BuildTestUserCreationQuery(ctx context.Context, testUserConfig *types.TestUserCreationConfig) (query string, args []interface{}) {
	returnValues := m.Called(ctx, testUserConfig)

	return returnValues.Get(0).(string), returnValues.Get(1).([]interface{})
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
