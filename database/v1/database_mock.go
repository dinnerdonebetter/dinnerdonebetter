package database

import (
	"context"

	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"

	"github.com/stretchr/testify/mock"
)

var _ DataManager = (*MockDatabase)(nil)

// BuildMockDatabase builds a mock database.
func BuildMockDatabase() *MockDatabase {
	return &MockDatabase{
		ValidInstrumentDataManager:               &mockmodels.ValidInstrumentDataManager{},
		ValidIngredientDataManager:               &mockmodels.ValidIngredientDataManager{},
		ValidPreparationDataManager:              &mockmodels.ValidPreparationDataManager{},
		ValidIngredientPreparationDataManager:    &mockmodels.ValidIngredientPreparationDataManager{},
		RequiredPreparationInstrumentDataManager: &mockmodels.RequiredPreparationInstrumentDataManager{},
		RecipeDataManager:                        &mockmodels.RecipeDataManager{},
		RecipeStepDataManager:                    &mockmodels.RecipeStepDataManager{},
		RecipeStepInstrumentDataManager:          &mockmodels.RecipeStepInstrumentDataManager{},
		RecipeStepIngredientDataManager:          &mockmodels.RecipeStepIngredientDataManager{},
		RecipeStepProductDataManager:             &mockmodels.RecipeStepProductDataManager{},
		RecipeIterationDataManager:               &mockmodels.RecipeIterationDataManager{},
		RecipeStepEventDataManager:               &mockmodels.RecipeStepEventDataManager{},
		IterationMediaDataManager:                &mockmodels.IterationMediaDataManager{},
		InvitationDataManager:                    &mockmodels.InvitationDataManager{},
		ReportDataManager:                        &mockmodels.ReportDataManager{},
		UserDataManager:                          &mockmodels.UserDataManager{},
		OAuth2ClientDataManager:                  &mockmodels.OAuth2ClientDataManager{},
		WebhookDataManager:                       &mockmodels.WebhookDataManager{},
	}
}

// MockDatabase is our mock database structure.
type MockDatabase struct {
	mock.Mock

	*mockmodels.ValidInstrumentDataManager
	*mockmodels.ValidIngredientDataManager
	*mockmodels.ValidPreparationDataManager
	*mockmodels.ValidIngredientPreparationDataManager
	*mockmodels.RequiredPreparationInstrumentDataManager
	*mockmodels.RecipeDataManager
	*mockmodels.RecipeStepDataManager
	*mockmodels.RecipeStepInstrumentDataManager
	*mockmodels.RecipeStepIngredientDataManager
	*mockmodels.RecipeStepProductDataManager
	*mockmodels.RecipeIterationDataManager
	*mockmodels.RecipeStepEventDataManager
	*mockmodels.IterationMediaDataManager
	*mockmodels.InvitationDataManager
	*mockmodels.ReportDataManager
	*mockmodels.UserDataManager
	*mockmodels.OAuth2ClientDataManager
	*mockmodels.WebhookDataManager
}

// Migrate satisfies the DataManager interface.
func (m *MockDatabase) Migrate(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// IsReady satisfies the DataManager interface.
func (m *MockDatabase) IsReady(ctx context.Context) (ready bool) {
	return m.Called(ctx).Bool(0)
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
