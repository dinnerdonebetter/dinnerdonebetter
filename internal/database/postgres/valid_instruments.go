package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.ValidInstrumentDataManager = (*SQLQuerier)(nil)
)

// ValidInstrumentExists fetches whether a valid instrument exists from the database.
func (q *SQLQuerier) ValidInstrumentExists(ctx context.Context, validInstrumentID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	result, err := q.generatedQuerier.ValidInstrumentExists(ctx, validInstrumentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, observability.PrepareError(err, logger, span, "performing valid instrument existence check")
	}

	return result, nil
}

// GetValidInstrument fetches a valid instrument from the database.
func (q *SQLQuerier) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	result, err := q.generatedQuerier.GetValidInstrument(ctx, validInstrumentID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid instrument")
	}

	instrument := &types.ValidInstrument{
		Description: result.Description,
		IconPath:    result.IconPath,
		ID:          result.ID,
		Name:        result.Name,
		PluralName:  result.PluralName,
		CreatedOn:   uint64(result.CreatedOn),
	}

	if result.LastUpdatedOn.Valid {
		t := uint64(result.LastUpdatedOn.Int64)
		instrument.LastUpdatedOn = &t
	}

	if result.ArchivedOn.Valid {
		t := uint64(result.ArchivedOn.Int64)
		instrument.ArchivedOn = &t
	}

	return instrument, nil
}

// GetRandomValidInstrument fetches a valid instrument from the database.
func (q *SQLQuerier) GetRandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	result, err := q.generatedQuerier.GetRandomValidInstrument(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid instrument")
	}

	instrument := &types.ValidInstrument{
		Description: result.Description,
		IconPath:    result.IconPath,
		ID:          result.ID,
		Name:        result.Name,
		PluralName:  result.PluralName,
		CreatedOn:   uint64(result.CreatedOn),
	}

	if result.LastUpdatedOn.Valid {
		t := uint64(result.LastUpdatedOn.Int64)
		instrument.LastUpdatedOn = &t
	}

	if result.ArchivedOn.Valid {
		t := uint64(result.ArchivedOn.Int64)
		instrument.ArchivedOn = &t
	}

	return instrument, nil
}

// SearchForValidInstruments fetches a valid instrument from the database.
func (q *SQLQuerier) SearchForValidInstruments(ctx context.Context, query string) ([]*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidInstrumentIDToSpan(span, query)

	results, err := q.generatedQuerier.SearchForValidInstruments(ctx, wrapQueryForILIKE(query))
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredients search query")
	}

	validInstruments := []*types.ValidInstrument{}

	for _, result := range results {
		instrument := &types.ValidInstrument{
			Description: result.Description,
			IconPath:    result.IconPath,
			ID:          result.ID,
			Name:        result.Name,
			PluralName:  result.PluralName,
			CreatedOn:   uint64(result.CreatedOn),
		}

		if result.LastUpdatedOn.Valid {
			t := uint64(result.LastUpdatedOn.Int64)
			instrument.LastUpdatedOn = &t
		}

		if result.ArchivedOn.Valid {
			t := uint64(result.ArchivedOn.Int64)
			instrument.ArchivedOn = &t
		}

		validInstruments = append(validInstruments, instrument)
	}

	return validInstruments, nil
}

// GetTotalValidInstrumentCount fetches the count of valid instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalValidInstrumentCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.generatedQuerier.GetTotalValidInstrumentCount(ctx)
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid instruments")
	}

	return uint64(count), nil
}

// GetValidInstruments fetches a list of valid instruments from the database that meet a particular filter.
func (q *SQLQuerier) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.ValidInstrumentList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidInstrumentList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	} else {
		filter = types.DefaultQueryFilter()
	}

	args := &generated.GetValidInstrumentsParams{
		CreatedAfter:  nullInt64ForUint64Field(filter.CreatedAfter),
		CreatedBefore: nullInt64ForUint64Field(filter.CreatedBefore),
		UpdatedAfter:  nullInt64ForUint64Field(filter.UpdatedAfter),
		UpdatedBefore: nullInt64ForUint64Field(filter.UpdatedBefore),
		Limit:         nullInt32ForUint8Field(filter.Limit),
	}

	results, err := q.generatedQuerier.GetValidInstruments(ctx, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid instruments list retrieval query")
	}

	for _, result := range results {
		instrument := &types.ValidInstrument{
			Description: result.Description,
			IconPath:    result.IconPath,
			ID:          result.ID,
			Name:        result.Name,
			PluralName:  result.PluralName,
			CreatedOn:   uint64(result.CreatedOn),
		}

		if result.LastUpdatedOn.Valid {
			t := uint64(result.LastUpdatedOn.Int64)
			instrument.LastUpdatedOn = &t
		}

		if result.ArchivedOn.Valid {
			t := uint64(result.ArchivedOn.Int64)
			instrument.ArchivedOn = &t
		}

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)

		x.ValidInstruments = append(x.ValidInstruments, instrument)
	}

	return x, nil
}

// CreateValidInstrument creates a valid instrument in the database.
func (q *SQLQuerier) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentDatabaseCreationInput) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, input.ID)

	args := &generated.CreateValidInstrumentParams{
		ID:          input.ID,
		Name:        input.Name,
		PluralName:  input.PluralName,
		Description: input.Description,
		IconPath:    input.IconPath,
	}

	if err := q.generatedQuerier.CreateValidInstrument(ctx, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating valid instrument")
	}

	x := &types.ValidInstrument{
		Description: input.Description,
		IconPath:    input.IconPath,
		ID:          input.ID,
		Name:        input.Name,
		PluralName:  input.PluralName,
		CreatedOn:   q.currentTime(),
	}

	tracing.AttachValidInstrumentIDToSpan(span, x.ID)
	logger.Info("valid instrument created")

	return x, nil
}

// UpdateValidInstrument updates a particular valid instrument.
func (q *SQLQuerier) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, updated.ID)
	tracing.AttachValidInstrumentIDToSpan(span, updated.ID)

	args := &generated.UpdateValidInstrumentParams{
		ID:          updated.ID,
		Name:        updated.Name,
		PluralName:  updated.PluralName,
		Description: updated.Description,
		IconPath:    updated.IconPath,
	}

	if err := q.generatedQuerier.UpdateValidInstrument(ctx, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid instrument")
	}

	logger.Info("valid instrument updated")

	return nil
}

// ArchiveValidInstrument archives a valid instrument from the database by its ID.
func (q *SQLQuerier) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	if err := q.generatedQuerier.ArchiveValidInstrument(ctx, validInstrumentID); err != nil {
		return observability.PrepareError(err, logger, span, "archiving valid instrument")
	}

	logger.Info("valid instrument archived")

	return nil
}
