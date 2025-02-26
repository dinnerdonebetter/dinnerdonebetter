package managers

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/events"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
)

/*

TODO:
- [x] all loggers are instantiated from spans
- [x] no more references to `GetUnfinalizedMealPlansWithExpiredVotingPeriods`
- [x] all returned errors have description strings
- [x] all relevant input params are accounted for in logs
- [x] all relevant input params are accounted for in traces
- [x] all pointer inputs have nil checks
- [x] all query filters are defaulted when nil
- [x] all CUD functions fire a data change event
- [x] list routes
- [x] read routes
- [x] search routes
- [x] create routes
- [x] update routes
- [x] archive routes
- [ ] unit tests lmfao

*/

type (
	ValidEnumerationsManager interface {
		SearchValidIngredientGroups(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, error)
		ListValidIngredientGroups(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, string, error)
		CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupCreationRequestInput) (*types.ValidIngredientGroup, error)
		ReadValidIngredientGroup(ctx context.Context, validIngredientGroupID string) (*types.ValidIngredientGroup, error)
		UpdateValidIngredientGroup(ctx context.Context, validIngredientGroupID string, input *types.ValidIngredientGroupUpdateRequestInput) error
		ArchiveValidIngredientGroup(ctx context.Context, validIngredientGroupID string) error

		ListValidIngredientMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, string, error)
		CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitCreationRequestInput) (*types.ValidIngredientMeasurementUnit, error)
		ReadValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error)
		UpdateValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string, input *types.ValidIngredientMeasurementUnitUpdateRequestInput) error
		ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error
		SearchValidIngredientMeasurementUnitsByIngredient(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, error)
		SearchValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, validMeasurementUnitID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, error)

		ListValidIngredientPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, string, error)
		CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationRequestInput) (*types.ValidIngredientPreparation, error)
		ReadValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error)
		UpdateValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string, input *types.ValidIngredientPreparationUpdateRequestInput) error
		ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error
		SearchValidIngredientPreparationsByIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, error)
		SearchValidIngredientPreparationsByPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, error)

		SearchValidIngredients(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error)
		ListValidIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredient, string, error)
		CreateValidIngredient(ctx context.Context, input *types.ValidIngredientCreationRequestInput) (*types.ValidIngredient, error)
		ReadValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error)
		RandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error)
		UpdateValidIngredient(ctx context.Context, validIngredientID string, input *types.ValidIngredientUpdateRequestInput) error
		ArchiveValidIngredient(ctx context.Context, validIngredientID string) error
		SearchValidIngredientsByPreparationAndIngredientName(ctx context.Context, preparationID, query string, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error)

		ListValidIngredientStateIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, string, error)
		CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientCreationRequestInput) (*types.ValidIngredientStateIngredient, error)
		ReadValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error)
		UpdateValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string, input *types.ValidIngredientStateIngredientUpdateRequestInput) error
		ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error
		SearchValidIngredientStateIngredientsByIngredient(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, error)
		SearchValidIngredientStateIngredientsByIngredientState(ctx context.Context, validIngredientStateID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, error)

		SearchValidIngredientStates(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, error)
		ListValidIngredientStates(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, string, error)
		CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateCreationRequestInput) (*types.ValidIngredientState, error)
		ReadValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error)
		UpdateValidIngredientState(ctx context.Context, validIngredientStateID string, input *types.ValidIngredientStateUpdateRequestInput) error
		ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error

		SearchValidMeasurementUnits(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error)
		SearchValidMeasurementUnitsByIngredientID(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error)
		ListValidMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, string, error)
		CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitCreationRequestInput) (*types.ValidMeasurementUnit, error)
		ReadValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error)
		UpdateValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string, input *types.ValidMeasurementUnitUpdateRequestInput) error
		ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error

		SearchValidInstruments(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidInstrument, error)
		ListValidInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidInstrument, string, error)
		CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentCreationRequestInput) (*types.ValidInstrument, error)
		ReadValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error)
		RandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error)
		UpdateValidInstrument(ctx context.Context, validInstrumentID string, input *types.ValidInstrumentUpdateRequestInput) error
		ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error

		ValidMeasurementUnitConversionsFromMeasurementUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error)
		ValidMeasurementUnitConversionsToMeasurementUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error)
		CreateValidMeasurementUnitConversion(ctx context.Context, input *types.ValidMeasurementUnitConversionCreationRequestInput) (*types.ValidMeasurementUnitConversion, error)
		ReadValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) (*types.ValidMeasurementUnitConversion, error)
		UpdateValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string, input *types.ValidMeasurementUnitConversionUpdateRequestInput) error
		ArchiveValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) error

		ListValidPreparationInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, string, error)
		CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationRequestInput) (*types.ValidPreparationInstrument, error)
		ReadValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error)
		UpdateValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string, input *types.ValidPreparationInstrumentUpdateRequestInput) error
		ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error
		SearchValidPreparationInstrumentsByPreparation(ctx context.Context, validPreparationID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, error)
		SearchValidPreparationInstrumentsByInstrument(ctx context.Context, validInstrumentID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, error)

		SearchValidPreparations(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidPreparation, error)
		ListValidPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparation, string, error)
		CreateValidPreparation(ctx context.Context, input *types.ValidPreparationCreationRequestInput) (*types.ValidPreparation, error)
		ReadValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error)
		RandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error)
		UpdateValidPreparation(ctx context.Context, validPreparationID string, input *types.ValidPreparationUpdateRequestInput) error
		ArchiveValidPreparation(ctx context.Context, validPreparationID string) error

		ListValidPreparationVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, string, error)
		CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselCreationRequestInput) (*types.ValidPreparationVessel, error)
		ReadValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error)
		UpdateValidPreparationVessel(ctx context.Context, validPreparationVesselID string, input *types.ValidPreparationVesselUpdateRequestInput) error
		ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error
		SearchValidPreparationVesselsByPreparation(ctx context.Context, validPreparationID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, error)
		SearchValidPreparationVesselsByVessel(ctx context.Context, validVesselID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, error)

		SearchValidVessels(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidVessel, error)
		ListValidVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidVessel, string, error)
		CreateValidVessel(ctx context.Context, input *types.ValidVesselCreationRequestInput) (*types.ValidVessel, error)
		ReadValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error)
		RandomValidVessel(ctx context.Context) (*types.ValidVessel, error)
		UpdateValidVessel(ctx context.Context, validVesselID string, input *types.ValidVesselUpdateRequestInput) error
		ArchiveValidVessel(ctx context.Context, validVesselID string) error
	}

	validEnumerationManager struct {
		logger                           logging.Logger
		tracer                           tracing.Tracer
		db                               types.ValidEnumerationDataManager
		dataChangesPublisher             messagequeue.Publisher
		validIngredientStatesSearchIndex textsearch.IndexSearcher[eatingindexing.ValidIngredientStateSearchSubset]
		validInstrumentSearchIndex       textsearch.IndexSearcher[eatingindexing.ValidInstrumentSearchSubset]
		validMeasurementUnitSearchIndex  textsearch.IndexSearcher[eatingindexing.ValidMeasurementUnitSearchSubset]
		validIngredientSearchIndex       textsearch.IndexSearcher[eatingindexing.ValidIngredientSearchSubset]
		validPreparationsSearchIndex     textsearch.IndexSearcher[eatingindexing.ValidPreparationSearchSubset]
		validVesselsSearchIndex          textsearch.IndexSearcher[eatingindexing.ValidVesselSearchSubset]
	}
)

var (
	_ ValidEnumerationsManager = (*validEnumerationManager)(nil)
)

func NewValidEnumerationsManager(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	db database.DataManager,
	cfg *msgconfig.QueuesConfig,
	publisherProvider messagequeue.PublisherProvider,
	searchConfig *textsearchcfg.Config,
	metricsProvider metrics.Provider,
) (ValidEnumerationsManager, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("failed to provide publisher for data changes topic: %w", err)
	}

	validIngredientStatesSearchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.ValidIngredientStateSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeValidIngredientStates)
	if err != nil {
		return nil, fmt.Errorf("failed to provide search index for %s index", eatingindexing.IndexTypeValidIngredientStates)
	}

	validInstrumentSearchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.ValidInstrumentSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeValidInstruments)
	if err != nil {
		return nil, fmt.Errorf("failed to provide search index for %s index", eatingindexing.IndexTypeValidInstruments)
	}

	validMeasurementUnitSearchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.ValidMeasurementUnitSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeValidMeasurementUnits)
	if err != nil {
		return nil, fmt.Errorf("failed to provide search index for %s index", eatingindexing.IndexTypeValidMeasurementUnits)
	}

	validIngredientSearchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.ValidIngredientSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeValidIngredients)
	if err != nil {
		return nil, fmt.Errorf("failed to provide search index for %s index", eatingindexing.IndexTypeValidIngredients)
	}

	validPreparationsSearchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.ValidPreparationSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeValidPreparations)
	if err != nil {
		return nil, fmt.Errorf("failed to provide search index for %s index", eatingindexing.IndexTypeValidPreparations)
	}

	validVesselsSearchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.ValidVesselSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeValidVessels)
	if err != nil {
		return nil, fmt.Errorf("failed to provide search index for %s index", eatingindexing.IndexTypeValidVessels)
	}

	m := &validEnumerationManager{
		db:                               db,
		tracer:                           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:                           logging.EnsureLogger(logger).WithName(o11yName),
		dataChangesPublisher:             dataChangesPublisher,
		validIngredientStatesSearchIndex: validIngredientStatesSearchIndex,
		validInstrumentSearchIndex:       validInstrumentSearchIndex,
		validMeasurementUnitSearchIndex:  validMeasurementUnitSearchIndex,
		validIngredientSearchIndex:       validIngredientSearchIndex,
		validPreparationsSearchIndex:     validPreparationsSearchIndex,
		validVesselsSearchIndex:          validVesselsSearchIndex,
	}

	return m, nil
}

// SearchValidIngredientGroups implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredientGroups(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query).WithValue(keys.UseDatabaseKey, useDatabase)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.UseDatabaseKey, useDatabase)

	results, err := m.db.SearchForValidIngredientGroups(ctx, query, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid ingredient groups failed")
	}

	return results, nil
}

// ListValidIngredientGroups implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidIngredientGroups(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientGroups(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing valid ingredient groups")
	}

	return results.Data, "", nil
}

// CreateValidIngredientGroup implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupCreationRequestInput) (*types.ValidIngredientGroup, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientGroupCreationRequestInputToValidIngredientGroupDatabaseCreationInput(input)
	created, err := m.db.CreateValidIngredientGroup(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient group")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientGroupCreated, map[string]any{
		keys.ValidIngredientGroupIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidIngredientGroup implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidIngredientGroup(ctx context.Context, validIngredientGroupID string) (*types.ValidIngredientGroup, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	result, err := m.db.GetValidIngredientGroup(ctx, validIngredientGroupID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient group")
	}

	return result, nil
}

// UpdateValidIngredientGroup implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient group")
	}

	existingValidIngredientGroup.Update(input)
	if err = m.db.UpdateValidIngredientGroup(ctx, existingValidIngredientGroup); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient group")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientGroupUpdated, map[string]any{
		keys.ValidIngredientGroupIDKey: existingValidIngredientGroup.ID,
	}))

	return nil
}

// ArchiveValidIngredientGroup implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidIngredientGroup(ctx context.Context, validIngredientGroupID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, validIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	if err := m.db.ArchiveValidIngredientGroup(ctx, validIngredientGroupID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient group")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientGroupArchived, map[string]any{
		keys.ValidIngredientGroupIDKey: validIngredientGroupID,
	}))

	return nil
}

// ListValidIngredientMeasurementUnits implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidIngredientMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient measurement units")
	}

	return results.Data, "", nil
}

// CreateValidIngredientMeasurementUnit implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitCreationRequestInput) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientMeasurementUnitCreationRequestInputToValidIngredientMeasurementUnitDatabaseCreationInput(input)
	created, err := m.db.CreateValidIngredientMeasurementUnit(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient measurement unit")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientMeasurementUnitCreated, map[string]any{
		keys.ValidIngredientMeasurementUnitIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidIngredientMeasurementUnit implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	result, err := m.db.GetValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient measurement unit")
	}

	return result, nil
}

// UpdateValidIngredientMeasurementUnit implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient measurement unit")
	}

	existingValidIngredientMeasurementUnit.Update(input)
	if err = m.db.UpdateValidIngredientMeasurementUnit(ctx, existingValidIngredientMeasurementUnit); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient measurement unit")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientMeasurementUnitUpdated, map[string]any{
		keys.ValidIngredientMeasurementUnitIDKey: existingValidIngredientMeasurementUnit.ID,
	}))

	return nil
}

// ArchiveValidIngredientMeasurementUnit implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	if err := m.db.ArchiveValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient measurement unit")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientMeasurementUnitArchived, map[string]any{
		keys.ValidIngredientMeasurementUnitIDKey: validIngredientMeasurementUnitID,
	}))

	return nil
}

// SearchValidIngredientMeasurementUnitsByIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredientMeasurementUnitsByIngredient(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	results, err := m.db.GetValidIngredientMeasurementUnitsForIngredient(ctx, validIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient measurement units for ingredient")
	}

	return results.Data, nil
}

// SearchValidIngredientMeasurementUnitsByMeasurementUnit implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, validMeasurementUnitID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	results, err := m.db.GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx, validMeasurementUnitID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient measurement units for measurement unit")
	}

	return results.Data, nil
}

// ListValidIngredientPreparations implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidIngredientPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientPreparations(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient preparations")
	}

	return results.Data, "", nil
}

// CreateValidIngredientPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationRequestInput) (*types.ValidIngredientPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientPreparationCreationRequestInputToValidIngredientPreparationDatabaseCreationInput(input)
	created, err := m.db.CreateValidIngredientPreparation(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient preparation")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientPreparationCreated, map[string]any{
		keys.ValidIngredientPreparationIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidIngredientPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	result, err := m.db.GetValidIngredientPreparation(ctx, validIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient preparation")
	}

	return result, nil
}

// UpdateValidIngredientPreparation implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient preparation")
	}

	existingValidIngredientPreparation.Update(input)
	if err = m.db.UpdateValidIngredientPreparation(ctx, existingValidIngredientPreparation); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientPreparationUpdated, map[string]any{
		keys.ValidIngredientPreparationIDKey: existingValidIngredientPreparation.ID,
	}))

	return nil
}

// ArchiveValidIngredientPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	if err := m.db.ArchiveValidIngredientPreparation(ctx, validIngredientPreparationID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient preparation")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientPreparationArchived, map[string]any{
		keys.ValidIngredientPreparationIDKey: validIngredientPreparationID,
	}))

	return nil
}

// SearchValidIngredientPreparationsByIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredientPreparationsByIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, ingredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, ingredientID)

	results, err := m.db.GetValidIngredientPreparationsForIngredient(ctx, ingredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient preparations for ingredient")
	}

	return results.Data, nil
}

// SearchValidIngredientPreparationsByPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredientPreparationsByPreparation(ctx context.Context, validPreparationID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	results, err := m.db.GetValidIngredientPreparationsForPreparation(ctx, validPreparationID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient preparations for ingredient")
	}

	return results.Data, nil
}

// SearchValidIngredients implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredients(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query).WithValue(keys.UseDatabaseKey, useDatabase)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.UseDatabaseKey, useDatabase)

	var (
		results []*types.ValidIngredient
	)
	if useDatabase {
		rawResults, err := m.db.SearchForValidIngredients(ctx, query, filter)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching database for valid ingredients")
		}

		results = rawResults.Data
	} else {
		var validIngredientSubsets []*eatingindexing.ValidIngredientSearchSubset
		validIngredientSubsets, err := m.validIngredientSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching valid ingredient search index for valid ingredients")
		}

		ids := []string{}
		for _, validIngredientSubset := range validIngredientSubsets {
			ids = append(ids, validIngredientSubset.ID)
		}

		results, err = m.db.GetValidIngredientsWithIDs(ctx, ids)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredients from database")
		}
	}

	return results, nil
}

// ListValidIngredients implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredient, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredients(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing valid ingredients")
	}

	return results.Data, "", nil
}

// CreateValidIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientCreationRequestInput) (*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientCreationRequestInputToValidIngredientDatabaseCreationInput(input)
	created, err := m.db.CreateValidIngredient(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientCreated, map[string]any{
		keys.ValidIngredientIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	result, err := m.db.GetValidIngredient(ctx, validIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient")
	}

	return result, nil
}

// RandomValidIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) RandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	result, err := m.db.GetRandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching random valid ingredient")
	}

	return result, nil
}

// UpdateValidIngredient implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient")
	}

	existingValidIngredient.Update(input)
	if err = m.db.UpdateValidIngredient(ctx, existingValidIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientUpdated, map[string]any{
		keys.ValidIngredientIDKey: existingValidIngredient.ID,
	}))

	return nil
}

// ArchiveValidIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	if err := m.db.ArchiveValidIngredient(ctx, validIngredientID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientArchived, map[string]any{
		keys.ValidIngredientIDKey: validIngredientID,
	}))

	return nil
}

// SearchValidIngredientsByPreparationAndIngredientName implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredientsByPreparationAndIngredientName(ctx context.Context, validPreparationID, query string, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query).WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	validIngredients, err := m.db.SearchForValidIngredientsForPreparation(ctx, validPreparationID, query, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid ingredient preparations")
	}

	return validIngredients.Data, nil
}

// ListValidIngredientStateIngredients implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidIngredientStateIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientStateIngredients(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing valid ingredient state ingredients")
	}

	return results.Data, "", nil
}

// CreateValidIngredientStateIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientCreationRequestInput) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientStateIngredientCreationRequestInputToValidIngredientStateIngredientDatabaseCreationInput(input)
	created, err := m.db.CreateValidIngredientStateIngredient(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient state ingredient")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientStateIngredientCreated, map[string]any{
		keys.ValidIngredientStateIngredientIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidIngredientStateIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	result, err := m.db.GetValidIngredientStateIngredient(ctx, validIngredientStateIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient state ingredient")
	}

	return result, nil
}

// UpdateValidIngredientStateIngredient implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient state ingredient")
	}

	existingValidIngredientStateIngredient.Update(input)
	if err = m.db.UpdateValidIngredientStateIngredient(ctx, existingValidIngredientStateIngredient); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state ingredient")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientStateIngredientUpdated, map[string]any{
		keys.ValidIngredientStateIngredientIDKey: existingValidIngredientStateIngredient.ID,
	}))

	return nil
}

// ArchiveValidIngredientStateIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	if err := m.db.ArchiveValidIngredientStateIngredient(ctx, validIngredientStateIngredientID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient state ingredient")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientStateIngredientArchived, map[string]any{
		keys.ValidIngredientStateIngredientIDKey: validIngredientStateIngredientID,
	}))

	return nil
}

// SearchValidIngredientStateIngredientsByIngredient implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredientStateIngredientsByIngredient(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	results, err := m.db.GetValidIngredientStateIngredientsForIngredient(ctx, validIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient state ingredients for ingredient")
	}

	return results.Data, nil
}

// SearchValidIngredientStateIngredientsByIngredientState implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredientStateIngredientsByIngredientState(ctx context.Context, validIngredientStateID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	results, err := m.db.GetValidIngredientStateIngredientsForIngredientState(ctx, validIngredientStateID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient state ingredients for ingredient state")
	}

	return results.Data, nil
}

// SearchValidIngredientStates implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidIngredientStates(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query).WithValue(keys.UseDatabaseKey, useDatabase)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.UseDatabaseKey, useDatabase)

	var (
		results []*types.ValidIngredientState
		err     error
	)
	if useDatabase {
		results, err = m.db.SearchForValidIngredientStates(ctx, query)
	} else {
		var validIngredientStateSubsets []*eatingindexing.ValidIngredientStateSearchSubset
		validIngredientStateSubsets, err = m.validIngredientStatesSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching valid ingredient states")
		}

		ids := []string{}
		for _, validIngredientStateSubset := range validIngredientStateSubsets {
			ids = append(ids, validIngredientStateSubset.ID)
		}

		results, err = m.db.GetValidIngredientStatesWithIDs(ctx, ids)
	}

	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient states")
	}

	return results, nil
}

// ListValidIngredientStates implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidIngredientStates(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidIngredientStates(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing valid ingredient states")
	}

	return results.Data, "", nil
}

// CreateValidIngredientState implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateCreationRequestInput) (*types.ValidIngredientState, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidIngredientStateCreationRequestInputToValidIngredientStateDatabaseCreationInput(input)
	created, err := m.db.CreateValidIngredientState(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient state")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientStateCreated, map[string]any{
		keys.ValidIngredientStateIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidIngredientState implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	result, err := m.db.GetValidIngredientState(ctx, validIngredientStateID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient state")
	}

	return result, nil
}

// UpdateValidIngredientState implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient state")
	}

	existingValidIngredientState.Update(input)
	if err = m.db.UpdateValidIngredientState(ctx, existingValidIngredientState); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientStateUpdated, map[string]any{
		keys.ValidIngredientStateIDKey: existingValidIngredientState.ID,
	}))

	return nil
}

// ArchiveValidIngredientState implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	if err := m.db.ArchiveValidIngredientState(ctx, validIngredientStateID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient state")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidIngredientStateArchived, map[string]any{
		keys.ValidIngredientStateIDKey: validIngredientStateID,
	}))

	return nil
}

// SearchValidMeasurementUnits implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidMeasurementUnits(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query).WithValue(keys.UseDatabaseKey, useDatabase)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.UseDatabaseKey, useDatabase)

	var (
		results []*types.ValidMeasurementUnit
		err     error
	)
	if useDatabase {
		results, err = m.db.SearchForValidMeasurementUnits(ctx, query)
	} else {
		var validMeasurementUnitSubsets []*eatingindexing.ValidMeasurementUnitSearchSubset
		validMeasurementUnitSubsets, err = m.validMeasurementUnitSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid measurement units")
		}

		ids := []string{}
		for _, validMeasurementUnitSubset := range validMeasurementUnitSubsets {
			ids = append(ids, validMeasurementUnitSubset.ID)
		}

		results, err = m.db.GetValidMeasurementUnitsWithIDs(ctx, ids)
	}

	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid measurement units")
	}

	return results, nil
}

// SearchValidMeasurementUnitsByIngredientID implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidMeasurementUnitsByIngredientID(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	validMeasurementUnits, err := m.db.ValidMeasurementUnitsForIngredientID(ctx, validIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching valid measurement units for ingredient")
	}

	return validMeasurementUnits.Data, nil
}

// ListValidMeasurementUnits implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching valid measurement units")
	}

	return results.Data, "", nil
}

// CreateValidMeasurementUnit implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitCreationRequestInput) (*types.ValidMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidMeasurementUnitCreationRequestInputToValidMeasurementUnitDatabaseCreationInput(input)
	created, err := m.db.CreateValidMeasurementUnit(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid measurement unit")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidMeasurementUnitCreated, map[string]any{
		keys.ValidMeasurementUnitIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidMeasurementUnit implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	result, err := m.db.GetValidMeasurementUnit(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid measurement unit")
	}

	return result, nil
}

// UpdateValidMeasurementUnit implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid measurement unit")
	}

	existingValidMeasurementUnit.Update(input)
	if err = m.db.UpdateValidMeasurementUnit(ctx, existingValidMeasurementUnit); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement unit")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidMeasurementUnitUpdated, map[string]any{
		keys.ValidMeasurementUnitIDKey: existingValidMeasurementUnit.ID,
	}))

	return nil
}

// ArchiveValidMeasurementUnit implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	if err := m.db.ArchiveValidMeasurementUnit(ctx, validMeasurementUnitID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement unit")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidMeasurementUnitArchived, map[string]any{
		keys.ValidMeasurementUnitIDKey: validMeasurementUnitID,
	}))

	return nil
}

// SearchValidInstruments implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidInstruments(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query).WithValue(keys.UseDatabaseKey, useDatabase)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.UseDatabaseKey, useDatabase)

	var (
		results []*types.ValidInstrument
		err     error
	)
	if useDatabase {
		results, err = m.db.SearchForValidInstruments(ctx, query)
	} else {
		var validInstrumentSubsets []*eatingindexing.ValidInstrumentSearchSubset
		validInstrumentSubsets, err = m.validInstrumentSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid instruments")
		}

		ids := []string{}
		for _, validInstrumentSubset := range validInstrumentSubsets {
			ids = append(ids, validInstrumentSubset.ID)
		}

		results, err = m.db.GetValidInstrumentsWithIDs(ctx, ids)
	}

	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid instruments")
	}

	return results, nil
}

// ListValidInstruments implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidInstrument, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidInstruments(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing valid instruments")
	}

	return results.Data, "", nil
}

// CreateValidInstrument implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentCreationRequestInput) (*types.ValidInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidInstrumentCreationRequestInputToValidInstrumentDatabaseCreationInput(input)
	created, err := m.db.CreateValidInstrument(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid instrument")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidInstrumentCreated, map[string]any{
		keys.ValidInstrumentIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidInstrument implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

	result, err := m.db.GetValidInstrument(ctx, validInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid instrument")
	}

	return result, nil
}

// RandomValidInstrument implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) RandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	result, err := m.db.GetRandomValidInstrument(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching random valid instrument")
	}

	return result, nil
}

// UpdateValidInstrument implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid instrument")
	}

	existingValidInstrument.Update(input)
	if err = m.db.UpdateValidInstrument(ctx, existingValidInstrument); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid instrument")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidInstrumentUpdated, map[string]any{
		keys.ValidInstrumentIDKey: existingValidInstrument.ID,
	}))

	return nil
}

// ArchiveValidInstrument implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	if err := m.db.ArchiveValidInstrument(ctx, validInstrumentID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid instrument")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidInstrumentArchived, map[string]any{
		keys.ValidInstrumentIDKey: validInstrumentID,
	}))

	return nil
}

// ValidMeasurementUnitConversionsFromMeasurementUnit implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ValidMeasurementUnitConversionsFromMeasurementUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	results, err := m.db.GetValidMeasurementUnitConversionsFromUnit(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid measurement unit conversions from unit")
	}

	return results, nil
}

// ValidMeasurementUnitConversionsToMeasurementUnit implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ValidMeasurementUnitConversionsToMeasurementUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	results, err := m.db.GetValidMeasurementUnitConversionsToUnit(ctx, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid measurement unit conversions to unit")
	}

	return results, nil
}

// CreateValidMeasurementUnitConversion implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidMeasurementUnitConversion(ctx context.Context, input *types.ValidMeasurementUnitConversionCreationRequestInput) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidMeasurementUnitConversionCreationRequestInputToValidMeasurementUnitConversionDatabaseCreationInput(input)
	created, err := m.db.CreateValidMeasurementUnitConversion(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid measurement unit conversion")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidMeasurementUnitConversionCreated, map[string]any{
		keys.ValidMeasurementUnitConversionIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidMeasurementUnitConversion implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)

	result, err := m.db.GetValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid measurement unit conversion")
	}

	return result, nil
}

// UpdateValidMeasurementUnitConversion implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid measurement unit conversion")
	}

	existingValidMeasurementUnitConversion.Update(input)
	if err = m.db.UpdateValidMeasurementUnitConversion(ctx, existingValidMeasurementUnitConversion); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement unit conversion")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidMeasurementUnitConversionUpdated, map[string]any{
		keys.ValidMeasurementUnitConversionIDKey: existingValidMeasurementUnitConversion.ID,
	}))

	return nil
}

// ArchiveValidMeasurementUnitConversion implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)

	if err := m.db.ArchiveValidMeasurementUnitConversion(ctx, validMeasurementUnitConversionID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement unit conversion")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidMeasurementUnitConversionArchived, map[string]any{
		keys.ValidMeasurementUnitConversionIDKey: validMeasurementUnitConversionID,
	}))

	return nil
}

// ListValidPreparationInstruments implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidPreparationInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidPreparationInstruments(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing valid preparation instruments")
	}

	return results.Data, "", nil
}

// CreateValidPreparationInstrument implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationRequestInput) (*types.ValidPreparationInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidPreparationInstrumentCreationRequestInputToValidPreparationInstrumentDatabaseCreationInput(input)
	created, err := m.db.CreateValidPreparationInstrument(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid preparation instrument")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidPreparationInstrumentCreated, map[string]any{
		keys.ValidPreparationInstrumentIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidPreparationInstrument implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	result, err := m.db.GetValidPreparationInstrument(ctx, validPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid preparation instrument")
	}

	return result, nil
}

// UpdateValidPreparationInstrument implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid preparation instrument")
	}

	existingValidPreparationInstrument.Update(input)
	if err = m.db.UpdateValidPreparationInstrument(ctx, existingValidPreparationInstrument); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation instrument")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidPreparationInstrumentUpdated, map[string]any{
		keys.ValidPreparationInstrumentIDKey: existingValidPreparationInstrument.ID,
	}))

	return nil
}

// ArchiveValidPreparationInstrument implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	if err := m.db.ArchiveValidPreparationInstrument(ctx, validPreparationInstrumentID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid preparation instrument")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidPreparationInstrumentArchived, map[string]any{
		keys.ValidPreparationInstrumentIDKey: validPreparationInstrumentID,
	}))

	return nil
}

// SearchValidPreparationInstrumentsByPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidPreparationInstrumentsByPreparation(ctx context.Context, validPreparationID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	validPreparationInstruments, err := m.db.GetValidPreparationInstrumentsForPreparation(ctx, validPreparationID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching valid preparation instruments by preparation")
	}

	return validPreparationInstruments.Data, nil
}

// SearchValidPreparationInstrumentsByInstrument implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidPreparationInstrumentsByInstrument(ctx context.Context, validInstrumentID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	validPreparationInstruments, err := m.db.GetValidPreparationInstrumentsForInstrument(ctx, validInstrumentID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching valid preparation instruments by instrument")
	}

	return validPreparationInstruments.Data, nil
}

// SearchValidPreparations implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidPreparations(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query).WithValue(keys.UseDatabaseKey, useDatabase)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.UseDatabaseKey, useDatabase)

	var (
		results []*types.ValidPreparation
		err     error
	)
	if useDatabase {
		results, err = m.db.SearchForValidPreparations(ctx, query)
	} else {
		var validPreparationSubsets []*eatingindexing.ValidPreparationSearchSubset
		validPreparationSubsets, err = m.validPreparationsSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching valid preparations")
		}

		ids := []string{}
		for _, validPreparationSubset := range validPreparationSubsets {
			ids = append(ids, validPreparationSubset.ID)
		}

		results, err = m.db.GetValidPreparationsWithIDs(ctx, ids)
	}

	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching valid preparations")
	}

	return results, nil
}

// ListValidPreparations implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparation, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidPreparations(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing valid preparations")
	}

	return results.Data, "", nil
}

// CreateValidPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationCreationRequestInput) (*types.ValidPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidPreparationCreationRequestInputToValidPreparationDatabaseCreationInput(input)
	created, err := m.db.CreateValidPreparation(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid preparation")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidPreparationCreated, map[string]any{
		keys.ValidPreparationIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	result, err := m.db.GetValidPreparation(ctx, validPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	return result, nil
}

// RandomValidPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) RandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	result, err := m.db.GetRandomValidPreparation(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching random valid preparation")
	}

	return result, nil
}

// UpdateValidPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) UpdateValidPreparation(ctx context.Context, validPreparationID string, input *types.ValidPreparationUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidPreparation, err := m.db.GetValidPreparation(ctx, validPreparationID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching valid preparation")
	}

	existingValidPreparation.Update(input)
	if err = m.db.UpdateValidPreparation(ctx, existingValidPreparation); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidPreparationUpdated, map[string]any{keys.ValidPreparationIDKey: existingValidPreparation.ID}))

	return nil
}

// ArchiveValidPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	if err := m.db.ArchiveValidPreparation(ctx, validPreparationID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid preparation")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidPreparationArchived, map[string]any{
		keys.ValidPreparationIDKey: validPreparationID,
	}))

	return nil
}

// ListValidPreparationVessels implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidPreparationVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidPreparationVessels(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing valid preparation vessels")
	}

	return results.Data, "", nil
}

// CreateValidPreparationVessel implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselCreationRequestInput) (*types.ValidPreparationVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidPreparationVesselCreationRequestInputToValidPreparationVesselDatabaseCreationInput(input)
	created, err := m.db.CreateValidPreparationVessel(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid preparation vessel")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidPreparationVesselCreated, map[string]any{
		keys.ValidPreparationVesselIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidPreparationVessel implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	result, err := m.db.GetValidPreparationVessel(ctx, validPreparationVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid preparation vessel")
	}

	return result, nil
}

// UpdateValidPreparationVessel implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) UpdateValidPreparationVessel(ctx context.Context, validPreparationVesselID string, input *types.ValidPreparationVesselUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	existingValidPreparationVessel, err := m.db.GetValidPreparationVessel(ctx, validPreparationVesselID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching valid preparation vessel")
	}

	existingValidPreparationVessel.Update(input)
	if err = m.db.UpdateValidPreparationVessel(ctx, existingValidPreparationVessel); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation vessel")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidPreparationVesselUpdated, map[string]any{
		keys.ValidPreparationVesselIDKey: existingValidPreparationVessel.ID,
	}))

	return nil
}

// ArchiveValidPreparationVessel implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	if err := m.db.ArchiveValidPreparationVessel(ctx, validPreparationVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid preparation vessel")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidPreparationVesselArchived, map[string]any{
		keys.ValidPreparationVesselIDKey: validPreparationVesselID,
	}))

	return nil
}

// SearchValidPreparationVesselsByPreparation implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidPreparationVesselsByPreparation(ctx context.Context, validPreparationID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	validPreparationVessels, err := m.db.GetValidPreparationVesselsForPreparation(ctx, validPreparationID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching valid preparation vessels by preparation")
	}

	return validPreparationVessels.Data, nil
}

// SearchValidPreparationVesselsByVessel implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidPreparationVesselsByVessel(ctx context.Context, validVesselID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	validPreparationVessels, err := m.db.GetValidPreparationVesselsForVessel(ctx, validVesselID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching valid preparation vessels by vessel")
	}

	return validPreparationVessels.Data, nil
}

// SearchValidVessels implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) SearchValidVessels(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query).WithValue(keys.UseDatabaseKey, useDatabase)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.UseDatabaseKey, useDatabase)

	var (
		validVessels []*types.ValidVessel
		err          error
	)
	if useDatabase {
		validVessels, err = m.db.SearchForValidVessels(ctx, query)
	} else {
		var validVesselSubsets []*eatingindexing.ValidVesselSearchSubset
		validVesselSubsets, err = m.validVesselsSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching index for valid vessels")
		}

		ids := []string{}
		for _, validVesselSubset := range validVesselSubsets {
			ids = append(ids, validVesselSubset.ID)
		}

		validVessels, err = m.db.GetValidVesselsWithIDs(ctx, ids)
	}
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "searching valid vessels")
	}

	return validVessels, nil
}

// ListValidVessels implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ListValidVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidVessel, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	logger := m.logger.WithSpan(span)

	results, err := m.db.GetValidVessels(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "listing valid vessels")
	}

	return results.Data, "", nil
}

// CreateValidVessel implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) CreateValidVessel(ctx context.Context, input *types.ValidVesselCreationRequestInput) (*types.ValidVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertValidVesselCreationRequestInputToValidVesselDatabaseCreationInput(input)
	created, err := m.db.CreateValidVessel(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid vessel")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidVesselCreated, map[string]any{
		keys.ValidVesselIDKey: created.ID,
	}))

	return created, nil
}

// ReadValidVessel implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ReadValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	result, err := m.db.GetValidVessel(ctx, validVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid vessel")
	}

	return result, nil
}

// RandomValidVessel implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) RandomValidVessel(ctx context.Context) (*types.ValidVessel, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span)

	result, err := m.db.GetRandomValidVessel(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching random valid vessel")
	}

	return result, nil
}

// UpdateValidVessel implements the ValidEnumerationsManager interface.
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
		return observability.PrepareAndLogError(err, logger, span, "fetching valid vessel")
	}

	existingValidVessel.Update(input)
	if err = m.db.UpdateValidVessel(ctx, existingValidVessel); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid vessel")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidVesselUpdated, map[string]any{
		keys.ValidVesselIDKey: existingValidVessel.ID,
	}))

	return nil
}

// ArchiveValidVessel implements the ValidEnumerationsManager interface.
func (m *validEnumerationManager) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	if err := m.db.ArchiveValidVessel(ctx, validVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid vessel")
	}

	m.dataChangesPublisher.PublishAsync(ctx, buildDataChangeMessageFromContext(ctx, logger, events.ValidVesselArchived, map[string]any{
		keys.ValidVesselIDKey: validVesselID,
	}))

	return nil
}
