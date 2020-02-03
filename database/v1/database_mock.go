package database

import (
	"context"

	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"

	"github.com/stretchr/testify/mock"
)

var _ Database = (*MockDatabase)(nil)

// BuildMockDatabase builds a mock database
func BuildMockDatabase() *MockDatabase {
	return &MockDatabase{
		InstrumentDataManager:                    &mockmodels.InstrumentDataManager{},
		IngredientDataManager:                    &mockmodels.IngredientDataManager{},
		PreparationDataManager:                   &mockmodels.PreparationDataManager{},
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

// MockDatabase is our mock database structure
type MockDatabase struct {
	mock.Mock

	*mockmodels.InstrumentDataManager
	*mockmodels.IngredientDataManager
	*mockmodels.PreparationDataManager
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

// Migrate satisfies the database.Database interface
func (m *MockDatabase) Migrate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// IsReady satisfies the database.Database interface
func (m *MockDatabase) IsReady(ctx context.Context) (ready bool) {
	args := m.Called(ctx)
	return args.Bool(0)
}
