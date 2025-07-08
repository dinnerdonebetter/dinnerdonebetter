package mocks

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

// ValidInstrumentExists is a mock function.
func (m *Repository) ValidInstrumentExists(ctx context.Context, validInstrumentID string) (bool, error) {
	returnValues := m.Called(ctx, validInstrumentID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidInstrument is a mock function.
func (m *Repository) GetValidInstrument(ctx context.Context, validInstrumentID string) (*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx, validInstrumentID)
	return returnValues.Get(0).(*recipeenums.ValidInstrument), returnValues.Error(1)
}

// GetRandomValidInstrument is a mock function.
func (m *Repository) GetRandomValidInstrument(ctx context.Context) (*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(*recipeenums.ValidInstrument), returnValues.Error(1)
}

// SearchForValidInstruments is a mock function.
func (m *Repository) SearchForValidInstruments(ctx context.Context, query string) ([]*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx, query)
	return returnValues.Get(0).([]*recipeenums.ValidInstrument), returnValues.Error(1)
}

// GetValidInstruments is a mock function.
func (m *Repository) GetValidInstruments(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidInstrument], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidInstrument]), returnValues.Error(1)
}

// CreateValidInstrument is a mock function.
func (m *Repository) CreateValidInstrument(ctx context.Context, input *recipeenums.ValidInstrumentDatabaseCreationInput) (*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidInstrument), returnValues.Error(1)
}

// UpdateValidInstrument is a mock function.
func (m *Repository) UpdateValidInstrument(ctx context.Context, updated *recipeenums.ValidInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidInstrument is a mock function.
func (m *Repository) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	return m.Called(ctx, validInstrumentID).Error(0)
}

// MarkValidInstrumentAsIndexed is a mock function.
func (m *Repository) MarkValidInstrumentAsIndexed(ctx context.Context, validInstrumentID string) error {
	return m.Called(ctx, validInstrumentID).Error(0)
}

// GetValidInstrumentIDsThatNeedSearchIndexing is a mock function.
func (m *Repository) GetValidInstrumentIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidInstrumentsWithIDs is a mock function.
func (m *Repository) GetValidInstrumentsWithIDs(ctx context.Context, ids []string) ([]*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*recipeenums.ValidInstrument), returnValues.Error(1)
}

// ValidIngredientExists is a mock function.
func (m *Repository) ValidIngredientExists(ctx context.Context, validIngredientID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredient is a mock function.
func (m *Repository) GetValidIngredient(ctx context.Context, validIngredientID string) (*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx, validIngredientID)
	return returnValues.Get(0).(*recipeenums.ValidIngredient), returnValues.Error(1)
}

// GetRandomValidIngredient is a mock function.
func (m *Repository) GetRandomValidIngredient(ctx context.Context) (*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(*recipeenums.ValidIngredient), returnValues.Error(1)
}

// SearchForValidIngredients is a mock function.
func (m *Repository) SearchForValidIngredients(ctx context.Context, query string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredient], error) {
	returnValues := m.Called(ctx, query, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredient]), returnValues.Error(1)
}

// SearchForValidIngredientsForPreparation is a mock function.
func (m *Repository) SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredient], error) {
	returnValues := m.Called(ctx, preparationID, query, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredient]), returnValues.Error(1)
}

// GetValidIngredients is a mock function.
func (m *Repository) GetValidIngredients(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredient], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredient]), returnValues.Error(1)
}

// CreateValidIngredient is a mock function.
func (m *Repository) CreateValidIngredient(ctx context.Context, input *recipeenums.ValidIngredientDatabaseCreationInput) (*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidIngredient), returnValues.Error(1)
}

// UpdateValidIngredient is a mock function.
func (m *Repository) UpdateValidIngredient(ctx context.Context, updated *recipeenums.ValidIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredient is a mock function.
func (m *Repository) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	return m.Called(ctx, validIngredientID).Error(0)
}

// MarkValidIngredientAsIndexed is a mock function.
func (m *Repository) MarkValidIngredientAsIndexed(ctx context.Context, validIngredientID string) error {
	return m.Called(ctx, validIngredientID).Error(0)
}

// GetValidIngredientIDsThatNeedSearchIndexing is a mock function.
func (m *Repository) GetValidIngredientIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidIngredientsWithIDs is a mock function.
func (m *Repository) GetValidIngredientsWithIDs(ctx context.Context, ids []string) ([]*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*recipeenums.ValidIngredient), returnValues.Error(1)
}

// ValidIngredientGroupExists is a mock method.
func (m *Repository) ValidIngredientGroupExists(ctx context.Context, validIngredientID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredientGroup is a mock method.
func (m *Repository) GetValidIngredientGroup(ctx context.Context, validIngredientID string) (*recipeenums.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Get(0).(*recipeenums.ValidIngredientGroup), returnValues.Error(1)
}

// GetValidIngredientGroups is a mock method.
func (m *Repository) GetValidIngredientGroups(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientGroup], error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientGroup]), returnValues.Error(1)
}

// SearchForValidIngredientGroups is a mock method.
func (m *Repository) SearchForValidIngredientGroups(ctx context.Context, query string, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, query, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientGroup), returnValues.Error(1)
}

// CreateValidIngredientGroup is a mock method.
func (m *Repository) CreateValidIngredientGroup(ctx context.Context, input *recipeenums.ValidIngredientGroupDatabaseCreationInput) (*recipeenums.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientGroup), returnValues.Error(1)
}

// UpdateValidIngredientGroup is a mock method.
func (m *Repository) UpdateValidIngredientGroup(ctx context.Context, updated *recipeenums.ValidIngredientGroup) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientGroup is a mock method.
func (m *Repository) ArchiveValidIngredientGroup(ctx context.Context, validIngredientID string) error {
	return m.Called(ctx, validIngredientID).Error(0)
}

// ValidPreparationExists is a mock function.
func (m *Repository) ValidPreparationExists(ctx context.Context, validPreparationID string) (bool, error) {
	returnValues := m.Called(ctx, validPreparationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidPreparation is a mock function.
func (m *Repository) GetValidPreparation(ctx context.Context, validPreparationID string) (*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx, validPreparationID)
	return returnValues.Get(0).(*recipeenums.ValidPreparation), returnValues.Error(1)
}

// GetRandomValidPreparation is a mock function.
func (m *Repository) GetRandomValidPreparation(ctx context.Context) (*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(*recipeenums.ValidPreparation), returnValues.Error(1)
}

// SearchForValidPreparations is a mock function.
func (m *Repository) SearchForValidPreparations(ctx context.Context, query string) ([]*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx, query)
	return returnValues.Get(0).([]*recipeenums.ValidPreparation), returnValues.Error(1)
}

// GetValidPreparations is a mock function.
func (m *Repository) GetValidPreparations(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidPreparation], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidPreparation]), returnValues.Error(1)
}

// CreateValidPreparation is a mock function.
func (m *Repository) CreateValidPreparation(ctx context.Context, input *recipeenums.ValidPreparationDatabaseCreationInput) (*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidPreparation), returnValues.Error(1)
}

// UpdateValidPreparation is a mock function.
func (m *Repository) UpdateValidPreparation(ctx context.Context, updated *recipeenums.ValidPreparation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparation is a mock function.
func (m *Repository) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}

// MarkValidPreparationAsIndexed is a mock function.
func (m *Repository) MarkValidPreparationAsIndexed(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}

// GetValidPreparationIDsThatNeedSearchIndexing is a mock function.
func (m *Repository) GetValidPreparationIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidPreparationsWithIDs is a mock function.
func (m *Repository) GetValidPreparationsWithIDs(ctx context.Context, ids []string) ([]*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*recipeenums.ValidPreparation), returnValues.Error(1)
}

// ValidIngredientPreparationExists is a mock function.
func (m *Repository) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientPreparationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredientPreparation is a mock function.
func (m *Repository) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*recipeenums.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, validIngredientPreparationID)
	return returnValues.Get(0).(*recipeenums.ValidIngredientPreparation), returnValues.Error(1)
}

// GetValidIngredientPreparations is a mock function.
func (m *Repository) GetValidIngredientPreparations(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientPreparation], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientPreparation]), returnValues.Error(1)
}

// GetValidIngredientPreparationsForIngredient is a mock function.
func (m *Repository) GetValidIngredientPreparationsForIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientPreparation], error) {
	returnValues := m.Called(ctx, ingredientID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientPreparation]), returnValues.Error(1)
}

// GetValidIngredientPreparationsForPreparation is a mock function.
func (m *Repository) GetValidIngredientPreparationsForPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientPreparation], error) {
	returnValues := m.Called(ctx, preparationID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientPreparation]), returnValues.Error(1)
}

// GetValidIngredientPreparationsForIngredientNameQuery is a mock function.
func (m *Repository) GetValidIngredientPreparationsForIngredientNameQuery(ctx context.Context, preparationID, query string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientPreparation], error) {
	returnValues := m.Called(ctx, preparationID, query, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientPreparation]), returnValues.Error(1)
}

// CreateValidIngredientPreparation is a mock function.
func (m *Repository) CreateValidIngredientPreparation(ctx context.Context, input *recipeenums.ValidIngredientPreparationDatabaseCreationInput) (*recipeenums.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidIngredientPreparation), returnValues.Error(1)
}

// UpdateValidIngredientPreparation is a mock function.
func (m *Repository) UpdateValidIngredientPreparation(ctx context.Context, updated *recipeenums.ValidIngredientPreparation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientPreparation is a mock function.
func (m *Repository) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	return m.Called(ctx, validIngredientPreparationID).Error(0)
}

// ValidMeasurementUnitsForIngredientID is a mock function.
func (m *Repository) ValidMeasurementUnitsForIngredientID(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidMeasurementUnit], error) {
	returnValues := m.Called(ctx, validIngredientID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidMeasurementUnit]), returnValues.Error(1)
}

// ValidMeasurementUnitExists is a mock function.
func (m *Repository) ValidMeasurementUnitExists(ctx context.Context, validMeasurementUnitID string) (bool, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidMeasurementUnit is a mock function.
func (m *Repository) GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*recipeenums.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)
	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnit), returnValues.Error(1)
}

// SearchForValidMeasurementUnitsByName is a mock function.
func (m *Repository) SearchForValidMeasurementUnits(ctx context.Context, query string) ([]*recipeenums.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, query)
	return returnValues.Get(0).([]*recipeenums.ValidMeasurementUnit), returnValues.Error(1)
}

// GetValidMeasurementUnits is a mock function.
func (m *Repository) GetValidMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidMeasurementUnit], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidMeasurementUnit]), returnValues.Error(1)
}

// CreateValidMeasurementUnit is a mock function.
func (m *Repository) CreateValidMeasurementUnit(ctx context.Context, input *recipeenums.ValidMeasurementUnitDatabaseCreationInput) (*recipeenums.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnit), returnValues.Error(1)
}

// UpdateValidMeasurementUnit is a mock function.
func (m *Repository) UpdateValidMeasurementUnit(ctx context.Context, updated *recipeenums.ValidMeasurementUnit) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidMeasurementUnit is a mock function.
func (m *Repository) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	return m.Called(ctx, validMeasurementUnitID).Error(0)
}

// MarkValidMeasurementUnitAsIndexed is a mock function.
func (m *Repository) MarkValidMeasurementUnitAsIndexed(ctx context.Context, validMeasurementUnitID string) error {
	return m.Called(ctx, validMeasurementUnitID).Error(0)
}

// GetValidMeasurementUnitIDsThatNeedSearchIndexing is a mock function.
func (m *Repository) GetValidMeasurementUnitIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidMeasurementUnitsWithIDs is a mock function.
func (m *Repository) GetValidMeasurementUnitsWithIDs(ctx context.Context, ids []string) ([]*recipeenums.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*recipeenums.ValidMeasurementUnit), returnValues.Error(1)
}

// ValidPreparationInstrumentExists is a mock function.
func (m *Repository) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID string) (bool, error) {
	returnValues := m.Called(ctx, validPreparationInstrumentID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidPreparationInstrument is a mock function.
func (m *Repository) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*recipeenums.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validPreparationInstrumentID)
	return returnValues.Get(0).(*recipeenums.ValidPreparationInstrument), returnValues.Error(1)
}

// GetValidPreparationInstruments is a mock function.
func (m *Repository) GetValidPreparationInstruments(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidPreparationInstrument], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidPreparationInstrument]), returnValues.Error(1)
}

// GetValidPreparationInstrumentsForPreparation is a mock function.
func (m *Repository) GetValidPreparationInstrumentsForPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidPreparationInstrument], error) {
	returnValues := m.Called(ctx, preparationID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidPreparationInstrument]), returnValues.Error(1)
}

// GetValidPreparationInstrumentsForInstrument is a mock function.
func (m *Repository) GetValidPreparationInstrumentsForInstrument(ctx context.Context, instrumentID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidPreparationInstrument], error) {
	returnValues := m.Called(ctx, instrumentID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidPreparationInstrument]), returnValues.Error(1)
}

// CreateValidPreparationInstrument is a mock function.
func (m *Repository) CreateValidPreparationInstrument(ctx context.Context, input *recipeenums.ValidPreparationInstrumentDatabaseCreationInput) (*recipeenums.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidPreparationInstrument), returnValues.Error(1)
}

// UpdateValidPreparationInstrument is a mock function.
func (m *Repository) UpdateValidPreparationInstrument(ctx context.Context, updated *recipeenums.ValidPreparationInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparationInstrument is a mock function.
func (m *Repository) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	return m.Called(ctx, validPreparationInstrumentID).Error(0)
}

// ValidIngredientMeasurementUnitExists is a mock function.
func (m *Repository) ValidIngredientMeasurementUnitExists(ctx context.Context, validIngredientMeasurementUnitID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientMeasurementUnitID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredientMeasurementUnit is a mock function.
func (m *Repository) GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*recipeenums.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, validIngredientMeasurementUnitID)
	return returnValues.Get(0).(*recipeenums.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

// GetValidIngredientMeasurementUnits is a mock function.
func (m *Repository) GetValidIngredientMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientMeasurementUnit], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientMeasurementUnit]), returnValues.Error(1)
}

// GetValidIngredientMeasurementUnitsForIngredient is a mock function.
func (m *Repository) GetValidIngredientMeasurementUnitsForIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientMeasurementUnit], error) {
	returnValues := m.Called(ctx, ingredientID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientMeasurementUnit]), returnValues.Error(1)
}

// GetValidIngredientMeasurementUnitsForMeasurementUnit is a mock function.
func (m *Repository) GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx context.Context, measurementUnitID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientMeasurementUnit], error) {
	returnValues := m.Called(ctx, measurementUnitID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientMeasurementUnit]), returnValues.Error(1)
}

// CreateValidIngredientMeasurementUnit is a mock function.
func (m *Repository) CreateValidIngredientMeasurementUnit(ctx context.Context, input *recipeenums.ValidIngredientMeasurementUnitDatabaseCreationInput) (*recipeenums.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

// UpdateValidIngredientMeasurementUnit is a mock function.
func (m *Repository) UpdateValidIngredientMeasurementUnit(ctx context.Context, updated *recipeenums.ValidIngredientMeasurementUnit) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientMeasurementUnit is a mock function.
func (m *Repository) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	return m.Called(ctx, validIngredientMeasurementUnitID).Error(0)
}

// GetValidMeasurementUnitConversionsFromUnit is a mock function.
func (m *Repository) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*recipeenums.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*recipeenums.ValidMeasurementUnitConversion), returnValues.Error(1)
}

// GetValidMeasurementUnitConversionsToUnit is a mock function.
func (m *Repository) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, validMeasurementUnitID string) ([]*recipeenums.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*recipeenums.ValidMeasurementUnitConversion), returnValues.Error(1)
}

// ValidMeasurementUnitConversionExists is a mock function.
func (m *Repository) ValidMeasurementUnitConversionExists(ctx context.Context, validPreparationID string) (bool, error) {
	returnValues := m.Called(ctx, validPreparationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidMeasurementUnitConversion is a mock function.
func (m *Repository) GetValidMeasurementUnitConversion(ctx context.Context, validPreparationID string) (*recipeenums.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validPreparationID)
	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnitConversion), returnValues.Error(1)
}

// CreateValidMeasurementUnitConversion is a mock function.
func (m *Repository) CreateValidMeasurementUnitConversion(ctx context.Context, input *recipeenums.ValidMeasurementUnitConversionDatabaseCreationInput) (*recipeenums.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnitConversion), returnValues.Error(1)
}

// UpdateValidMeasurementUnitConversion is a mock function.
func (m *Repository) UpdateValidMeasurementUnitConversion(ctx context.Context, updated *recipeenums.ValidMeasurementUnitConversion) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidMeasurementUnitConversion is a mock function.
func (m *Repository) ArchiveValidMeasurementUnitConversion(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}

// ValidIngredientStateExists is a mock function.
func (m *Repository) ValidIngredientStateExists(ctx context.Context, validIngredientStateID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientStateID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredientState is a mock function.
func (m *Repository) GetValidIngredientState(ctx context.Context, validIngredientStateID string) (*recipeenums.ValidIngredientState, error) {
	returnValues := m.Called(ctx, validIngredientStateID)
	return returnValues.Get(0).(*recipeenums.ValidIngredientState), returnValues.Error(1)
}

// SearchForValidIngredientStates is a mock function.
func (m *Repository) SearchForValidIngredientStates(ctx context.Context, query string) ([]*recipeenums.ValidIngredientState, error) {
	returnValues := m.Called(ctx, query)
	return returnValues.Get(0).([]*recipeenums.ValidIngredientState), returnValues.Error(1)
}

// GetValidIngredientStates is a mock function.
func (m *Repository) GetValidIngredientStates(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientState], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientState]), returnValues.Error(1)
}

// CreateValidIngredientState is a mock function.
func (m *Repository) CreateValidIngredientState(ctx context.Context, input *recipeenums.ValidIngredientStateDatabaseCreationInput) (*recipeenums.ValidIngredientState, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidIngredientState), returnValues.Error(1)
}

// UpdateValidIngredientState is a mock function.
func (m *Repository) UpdateValidIngredientState(ctx context.Context, updated *recipeenums.ValidIngredientState) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientState is a mock function.
func (m *Repository) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	return m.Called(ctx, validIngredientStateID).Error(0)
}

// MarkValidIngredientStateAsIndexed is a mock function.
func (m *Repository) MarkValidIngredientStateAsIndexed(ctx context.Context, validIngredientStateID string) error {
	return m.Called(ctx, validIngredientStateID).Error(0)
}

// GetValidIngredientStatesWithIDs is a mock function.
func (m *Repository) GetValidIngredientStatesWithIDs(ctx context.Context, ids []string) ([]*recipeenums.ValidIngredientState, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*recipeenums.ValidIngredientState), returnValues.Error(1)
}

// ValidIngredientStateIngredientExists is a mock function.
func (m *Repository) ValidIngredientStateIngredientExists(ctx context.Context, validIngredientStateIngredientID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientStateIngredientID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredientStateIngredient is a mock function.
func (m *Repository) GetValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*recipeenums.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, validIngredientStateIngredientID)
	return returnValues.Get(0).(*recipeenums.ValidIngredientStateIngredient), returnValues.Error(1)
}

// GetValidIngredientStateIngredients is a mock function.
func (m *Repository) GetValidIngredientStateIngredients(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientStateIngredient], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientStateIngredient]), returnValues.Error(1)
}

// GetValidIngredientStateIngredientsForIngredient is a mock function.
func (m *Repository) GetValidIngredientStateIngredientsForIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientStateIngredient], error) {
	returnValues := m.Called(ctx, ingredientID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientStateIngredient]), returnValues.Error(1)
}

// GetValidIngredientStateIngredientsForIngredientState is a mock function.
func (m *Repository) GetValidIngredientStateIngredientsForIngredientState(ctx context.Context, ingredientStateID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidIngredientStateIngredient], error) {
	returnValues := m.Called(ctx, ingredientStateID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidIngredientStateIngredient]), returnValues.Error(1)
}

// CreateValidIngredientStateIngredient is a mock function.
func (m *Repository) CreateValidIngredientStateIngredient(ctx context.Context, input *recipeenums.ValidIngredientStateIngredientDatabaseCreationInput) (*recipeenums.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidIngredientStateIngredient), returnValues.Error(1)
}

// UpdateValidIngredientStateIngredient is a mock function.
func (m *Repository) UpdateValidIngredientStateIngredient(ctx context.Context, updated *recipeenums.ValidIngredientStateIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientStateIngredient is a mock function.
func (m *Repository) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	return m.Called(ctx, validIngredientStateIngredientID).Error(0)
}

// GetValidIngredientStateIDsThatNeedSearchIndexing is a mock function.
func (m *Repository) GetValidIngredientStateIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// ValidVesselExists is a mock function.
func (m *Repository) ValidVesselExists(ctx context.Context, validVesselID string) (bool, error) {
	returnValues := m.Called(ctx, validVesselID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidVessel is a mock function.
func (m *Repository) GetValidVessel(ctx context.Context, validVesselID string) (*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx, validVesselID)
	return returnValues.Get(0).(*recipeenums.ValidVessel), returnValues.Error(1)
}

// GetRandomValidVessel is a mock function.
func (m *Repository) GetRandomValidVessel(ctx context.Context) (*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(*recipeenums.ValidVessel), returnValues.Error(1)
}

// SearchForValidVessels is a mock function.
func (m *Repository) SearchForValidVessels(ctx context.Context, query string) ([]*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx, query)
	return returnValues.Get(0).([]*recipeenums.ValidVessel), returnValues.Error(1)
}

// GetValidVessels is a mock function.
func (m *Repository) GetValidVessels(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidVessel], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidVessel]), returnValues.Error(1)
}

// CreateValidVessel is a mock function.
func (m *Repository) CreateValidVessel(ctx context.Context, input *recipeenums.ValidVesselDatabaseCreationInput) (*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidVessel), returnValues.Error(1)
}

// UpdateValidVessel is a mock function.
func (m *Repository) UpdateValidVessel(ctx context.Context, updated *recipeenums.ValidVessel) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidVessel is a mock function.
func (m *Repository) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	return m.Called(ctx, validVesselID).Error(0)
}

// MarkValidVesselAsIndexed is a mock function.
func (m *Repository) MarkValidVesselAsIndexed(ctx context.Context, validVesselID string) error {
	return m.Called(ctx, validVesselID).Error(0)
}

// GetValidVesselIDsThatNeedSearchIndexing is a mock function.
func (m *Repository) GetValidVesselIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidVesselsWithIDs is a mock function.
func (m *Repository) GetValidVesselsWithIDs(ctx context.Context, ids []string) ([]*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*recipeenums.ValidVessel), returnValues.Error(1)
}

// ValidPreparationVesselExists is a mock function.
func (m *Repository) ValidPreparationVesselExists(ctx context.Context, validPreparationVesselID string) (bool, error) {
	returnValues := m.Called(ctx, validPreparationVesselID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidPreparationVessel is a mock function.
func (m *Repository) GetValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*recipeenums.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validPreparationVesselID)
	return returnValues.Get(0).(*recipeenums.ValidPreparationVessel), returnValues.Error(1)
}

// GetValidPreparationVessels is a mock function.
func (m *Repository) GetValidPreparationVessels(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidPreparationVessel], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidPreparationVessel]), returnValues.Error(1)
}

// GetValidPreparationVesselsForPreparation is a mock function.
func (m *Repository) GetValidPreparationVesselsForPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidPreparationVessel], error) {
	returnValues := m.Called(ctx, preparationID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidPreparationVessel]), returnValues.Error(1)
}

// GetValidPreparationVesselsForVessel is a mock function.
func (m *Repository) GetValidPreparationVesselsForVessel(ctx context.Context, vesselID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[recipeenums.ValidPreparationVessel], error) {
	returnValues := m.Called(ctx, vesselID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[recipeenums.ValidPreparationVessel]), returnValues.Error(1)
}

// CreateValidPreparationVessel is a mock function.
func (m *Repository) CreateValidPreparationVessel(ctx context.Context, input *recipeenums.ValidPreparationVesselDatabaseCreationInput) (*recipeenums.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*recipeenums.ValidPreparationVessel), returnValues.Error(1)
}

// UpdateValidPreparationVessel is a mock function.
func (m *Repository) UpdateValidPreparationVessel(ctx context.Context, updated *recipeenums.ValidPreparationVessel) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparationVessel is a mock function.
func (m *Repository) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	return m.Called(ctx, validPreparationVesselID).Error(0)
}
