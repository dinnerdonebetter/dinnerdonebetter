package managers

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
)

/*

TODO:
- [ ] all returned errors have description strings
- [x] all loggers are instantiated from spans
- [ ] all relevant input params are accounted for in logs
- [ ] all relevant input params are accounted for in traces
- [ ] no more references to `GetUnfinalizedMealPlansWithExpiredVotingPeriods`
- [x] all pointer inputs have nil checks
- [ ] all query filters are defaulted when nil
- [ ] all CUD functions fire a data change event
- [ ] unit tests lmfao

*/

var (
	errUnimplemented = errors.New("not implemented")

	_ ValidEnumerationsManager = (*validEnumerationManager)(nil)
)

type (
	validEnumerationManager struct {
		logger logging.Logger
		tracer tracing.Tracer
		db     database.DataManager
	}
)

func (m *validEnumerationManager) SearchValidIngredientGroups(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	results, err := m.db.SearchForValidIngredientGroups(ctx, query, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid ingredient groups failed")
	}

	return results, nil
}

func (m *validEnumerationManager) ListValidIngredientGroups(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientGroups(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupCreationRequestInput) (*types.ValidIngredientGroup, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientGroupCreationRequestInputToValidIngredientGroupDatabaseCreationInput(input)
	result, err := m.db.CreateValidIngredientGroup(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidIngredientGroup(ctx context.Context, validIngredientGroupID string) (*types.ValidIngredientGroup, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidIngredientGroup(ctx context.Context, validIngredientGroupID string, input *types.ValidIngredientGroupUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidIngredientGroup, err := m.db.GetValidIngredientGroup(ctx, validIngredientGroupID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidIngredientGroup.Update(input)

	if err = m.db.UpdateValidIngredientGroup(ctx, existingValidIngredientGroup); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidIngredientGroup(ctx context.Context, validIngredientGroupID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidIngredientMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitCreationRequestInput) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientMeasurementUnitCreationRequestInputToValidIngredientMeasurementUnitDatabaseCreationInput(input)
	result, err := m.db.CreateValidIngredientMeasurementUnit(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string, input *types.ValidIngredientMeasurementUnitUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidIngredientMeasurementUnit, err := m.db.GetValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidIngredientMeasurementUnit.Update(input)

	if err = m.db.UpdateValidIngredientMeasurementUnit(ctx, existingValidIngredientMeasurementUnit); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidIngredientMeasurementUnitsByIngredient(ctx context.Context, query string) ([]*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, query string) ([]*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidIngredientPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientPreparations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationRequestInput) (*types.ValidIngredientPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientPreparationCreationRequestInputToValidIngredientPreparationDatabaseCreationInput(input)
	result, err := m.db.CreateValidIngredientPreparation(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string, input *types.ValidIngredientPreparationUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidIngredientPreparation, err := m.db.GetValidIngredientPreparation(ctx, validIngredientPreparationID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidIngredientPreparation.Update(input)

	if err = m.db.UpdateValidIngredientPreparation(ctx, existingValidIngredientPreparation); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidIngredientPreparationsByIngredient(ctx context.Context, query string) ([]*types.ValidIngredientPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidIngredientPreparationsByPreparation(ctx context.Context, query string) ([]*types.ValidIngredientPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidIngredients(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredients(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientCreationRequestInput) (*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientCreationRequestInputToValidIngredientDatabaseCreationInput(input)
	result, err := m.db.CreateValidIngredient(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, validIngredientID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) RandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidIngredient(ctx context.Context, validIngredientID string, input *types.ValidIngredientUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidIngredient, err := m.db.GetValidIngredient(ctx, validIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidIngredient.Update(input)

	if err = m.db.UpdateValidIngredient(ctx, existingValidIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, validIngredientID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidIngredientsByPreparationAndIngredientName(ctx context.Context, query string) ([]*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidIngredientStateIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientStateIngredients(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientCreationRequestInput) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientStateIngredientCreationRequestInputToValidIngredientStateIngredientDatabaseCreationInput(input)
	result, err := m.db.CreateValidIngredientStateIngredient(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string, input *types.ValidIngredientStateIngredientUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidIngredientStateIngredient, err := m.db.GetValidIngredientStateIngredient(ctx, validIngredientStateIngredientID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidIngredientStateIngredient.Update(input)

	if err = m.db.UpdateValidIngredientStateIngredient(ctx, existingValidIngredientStateIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidIngredientStateIngredientsByIngredient(ctx context.Context, query string) ([]*types.ValidIngredientStateIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidIngredientStateIngredientsByIngredientState(ctx context.Context, query string) ([]*types.ValidIngredientStateIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidIngredientStates(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidIngredientStates(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientStates(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateCreationRequestInput) (*types.ValidIngredientState, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientStateCreationRequestInputToValidIngredientStateDatabaseCreationInput(input)
	result, err := m.db.CreateValidIngredientState(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidIngredientState(ctx context.Context, validIngredientStateID string, input *types.ValidIngredientStateUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidIngredientState, err := m.db.GetValidIngredientState(ctx, validIngredientStateID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidIngredientState.Update(input)

	if err = m.db.UpdateValidIngredientState(ctx, existingValidIngredientState); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidMeasurementUnits(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidMeasurementUnitsByIngredientID(ctx context.Context, query string) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitCreationRequestInput) (*types.ValidMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidMeasurementUnitCreationRequestInputToValidMeasurementUnitDatabaseCreationInput(input)
	result, err := m.db.CreateValidMeasurementUnit(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string, input *types.ValidMeasurementUnitUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidMeasurementUnit, err := m.db.GetValidMeasurementUnit(ctx, validMeasurementUnitID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidMeasurementUnit.Update(input)

	if err = m.db.UpdateValidMeasurementUnit(ctx, existingValidMeasurementUnit); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidInstruments(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidInstruments(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentCreationRequestInput) (*types.ValidInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidInstrumentCreationRequestInputToValidInstrumentDatabaseCreationInput(input)
	result, err := m.db.CreateValidInstrument(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) RandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidInstrument(ctx context.Context, validInstrumentID string, input *types.ValidInstrumentUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidInstrument, err := m.db.GetValidInstrument(ctx, validInstrumentID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidInstrument.Update(input)

	if err = m.db.UpdateValidInstrument(ctx, existingValidInstrument); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ValidMeasurementUnitConversionsFromMeasurementUnit(ctx context.Context) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ValidMeasurementUnitConversionsToMeasurementUnit(ctx context.Context) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) CreateValidMeasurementUnitConversion(ctx context.Context, input *types.ValidMeasurementUnitConversionCreationRequestInput) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidMeasurementUnitConversionCreationRequestInputToValidMeasurementUnitConversionDatabaseCreationInput(input)
	result, err := m.db.CreateValidMeasurementUnitConversion(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string, input *types.ValidMeasurementUnitConversionUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidMeasurementUnitConversion, err := m.db.GetValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidMeasurementUnitConversion.Update(input)
	if err = m.db.UpdateValidMeasurementUnitConversion(ctx, existingValidMeasurementUnitConversion); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidPreparationInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidPreparationInstruments(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationRequestInput) (*types.ValidPreparationInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidPreparationInstrumentCreationRequestInputToValidPreparationInstrumentDatabaseCreationInput(input)
	result, err := m.db.CreateValidPreparationInstrument(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string, input *types.ValidPreparationInstrumentUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidPreparationInstrument, err := m.db.GetValidPreparationInstrument(ctx, validPreparationInstrumentID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidPreparationInstrument.Update(input)

	if err = m.db.UpdateValidPreparationInstrument(ctx, existingValidPreparationInstrument); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidPreparationInstrumentsByPreparation(ctx context.Context, query string) ([]*types.ValidPreparationInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidPreparationInstrumentsByInstrument(ctx context.Context, query string) ([]*types.ValidPreparationInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidPreparations(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidPreparations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationCreationRequestInput) (*types.ValidPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidPreparationCreationRequestInputToValidPreparationDatabaseCreationInput(input)
	result, err := m.db.CreateValidPreparation(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, validPreparationID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) RandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidPreparation(ctx context.Context, validPreparationID string, input *types.ValidPreparationUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidPreparation, err := m.db.GetValidPreparation(ctx, validPreparationID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidPreparation.Update(input)

	if err = m.db.UpdateValidPreparation(ctx, existingValidPreparation); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, validPreparationID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidPreparationVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidPreparationVessels(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselCreationRequestInput) (*types.ValidPreparationVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidPreparationVesselCreationRequestInputToValidPreparationVesselDatabaseCreationInput(input)
	result, err := m.db.CreateValidPreparationVessel(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidPreparationVessel(ctx context.Context, validPreparationVesselID string, input *types.ValidPreparationVesselUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidPreparationVessel, err := m.db.GetValidPreparationVessel(ctx, validPreparationVesselID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidPreparationVessel.Update(input)

	if err = m.db.UpdateValidPreparationVessel(ctx, existingValidPreparationVessel); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidPreparationVesselsByPreparation(ctx context.Context, query string) ([]*types.ValidPreparationVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidPreparationVesselsByVessel(ctx context.Context, query string) ([]*types.ValidPreparationVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) SearchValidVessels(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ListValidVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidVessels(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return results.Data, nil
}

func (m *validEnumerationManager) CreateValidVessel(ctx context.Context, input *types.ValidVesselCreationRequestInput) (*types.ValidVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidVesselCreationRequestInputToValidVesselDatabaseCreationInput(input)
	result, err := m.db.CreateValidVessel(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "")
	}

	return result, nil
}

func (m *validEnumerationManager) ReadValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, validVesselID)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) RandomValidVessel(ctx context.Context) (*types.ValidVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) UpdateValidVessel(ctx context.Context, validVesselID string, input *types.ValidVesselUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidVessel, err := m.db.GetValidVessel(ctx, validVesselID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving        ")
	}

	existingValidVessel.Update(input)

	if err = m.db.UpdateValidVessel(ctx, existingValidVessel); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating          ")
	}

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}

func (m *validEnumerationManager) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, validVesselID)

	return observability.PrepareAndLogError(errUnimplemented, logger, span, "")
}
