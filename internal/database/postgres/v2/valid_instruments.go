package v2

import (
	"context"
	"database/sql"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/doug-martin/goqu/v9"
	"github.com/jinzhu/copier"
)

const (
	validInstrumentsTableName = "valid_instruments"
)

type (
	// ValidInstrument represents a valid instrument.
	ValidInstrument struct {
		_ struct{}

		CreatedAt                      time.Time  `db:"created_at"`
		LastUpdatedAt                  *time.Time `db:"last_updated_at"`
		ArchivedAt                     *time.Time `db:"archived_at"`
		IconPath                       string     `db:"icon_path"`
		ID                             string     `db:"id"                                goqu:"skipupdate"`
		Name                           string     `db:"name"`
		PluralName                     string     `db:"plural_name"`
		Description                    string     `db:"description"`
		Slug                           string     `db:"slug"`
		DisplayInSummaryLists          bool       `db:"display_in_summary_lists"`
		IncludeInGeneratedInstructions bool       `db:"include_in_generated_instructions"`
		UsableForStorage               bool       `db:"usable_for_storage"`
	}
)

// CreateValidInstrument gets a valid instrument from the database.
func (c *DatabaseClient) CreateValidInstrument(ctx context.Context, input *ValidInstrument) (*types.ValidInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Insert(validInstrumentsTableName).Rows(
		input,
	)

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "creating valid instrument")
	}

	var output types.ValidInstrument
	if err := copier.Copy(&output, input); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetValidInstrument gets a valid instrument from the database.
func (c *DatabaseClient) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	x := &ValidInstrument{}
	q := c.xdb.From(validInstrumentsTableName).Where(goqu.Ex{
		idColumn:         validInstrumentID,
		archivedAtColumn: nil,
	})

	found, err := q.ScanStructContext(ctx, x)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, sql.ErrNoRows
	}

	var output types.ValidInstrument
	if err = copier.Copy(&output, x); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetValidInstruments gets a valid instrument from the database.
func (c *DatabaseClient) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidInstrument], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	q := c.xdb.From(validInstrumentsTableName).Where(goqu.Ex{
		archivedAtColumn: nil,
	})
	q = queryFilterToGoqu(q, filter)

	var x []ValidInstrument
	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	output := &types.QueryFilteredResult[types.ValidInstrument]{
		Data:       []*types.ValidInstrument{},
		Pagination: filter.ToPagination(),
	}
	for _, y := range x {
		var z types.ValidInstrument
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output.Data = append(output.Data, &z)
	}

	return output, nil
}

// SearchForValidInstruments gets a valid instrument from the database.
func (c *DatabaseClient) SearchForValidInstruments(ctx context.Context, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidInstrument], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	q := c.xdb.From(validInstrumentsTableName).Where(goqu.Ex{
		archivedAtColumn: nil,
		"name":           goqu.Op{"like": "%" + query + "%"},
	})
	q = queryFilterToGoqu(q, filter)

	var x []ValidInstrument
	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	output := &types.QueryFilteredResult[types.ValidInstrument]{
		Data:       []*types.ValidInstrument{},
		Pagination: filter.ToPagination(),
	}
	for _, y := range x {
		var z types.ValidInstrument
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output.Data = append(output.Data, &z)
	}

	return output, nil
}

// GetValidInstrumentsWithIDs gets a valid instrument from the database.
func (c *DatabaseClient) GetValidInstrumentsWithIDs(ctx context.Context, ids []string) ([]*types.ValidInstrument, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	x := []*ValidInstrument{}

	q := c.xdb.From(validInstrumentsTableName).Where(goqu.Ex{
		idColumn:         ids,
		archivedAtColumn: nil,
	})

	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	var output []*types.ValidInstrument
	for _, y := range x {
		var z types.ValidInstrument
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output = append(output, &z)
	}

	return output, nil
}

// UpdateValidInstrument gets a valid instrument from the database.
func (c *DatabaseClient) UpdateValidInstrument(ctx context.Context, input *types.ValidInstrument) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	var updateInput ValidInstrument
	if err := copier.Copy(&updateInput, input); err != nil {
		return observability.PrepareError(err, span, "copying input to output")
	}

	q := c.xdb.Update(validInstrumentsTableName).Set(
		updateInput,
	).Set(goqu.Ex{lastUpdatedAtColumn: goqu.L("NOW()")})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "updating valid instrument")
	}

	return nil
}

// ArchiveValidInstrument gets a valid instrument from the database.
func (c *DatabaseClient) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Update(validInstrumentsTableName).
		Set(goqu.Record{archivedAtColumn: goqu.L("NOW()")}).
		Where(goqu.Ex{idColumn: validInstrumentID})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "archiving valid instrument")
	}

	return nil
}
