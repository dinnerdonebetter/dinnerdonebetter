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
	validPreparationsTableName = "valid_preparations"
)

type (
	// ValidPreparation represents a valid preparation.
	ValidPreparation struct {
		_ struct{}

		CreatedAt                   time.Time  `db:"created_at"                                   goqu:"skipinsert"`
		LastUpdatedAt               *time.Time `db:"last_updated_at"                              goqu:"skipinsert"`
		ArchivedAt                  *time.Time `db:"archived_at"                                  goqu:"skipinsert"`
		MaximumInstrumentCount      *int32     `db:"maximum_instrument_count"`
		MaximumIngredientCount      *int32     `db:"maximum_ingredient_count"`
		MaximumVesselCount          *int32     `db:"maximum_vessel_count"`
		IconPath                    string     `db:"icon_path"`
		PastTense                   string     `db:"past_tense"`
		ID                          string     `db:"id"`
		Name                        string     `db:"name"`
		Description                 string     `db:"description"`
		Slug                        string     `db:"slug"`
		MinimumIngredientCount      int32      `db:"minimum_ingredient_count"`
		MinimumInstrumentCount      int32      `db:"minimum_instrument_count"`
		MinimumVesselCount          int32      `db:"minimum_vessel_count"`
		RestrictToIngredients       bool       `db:"restrict_to_ingredients"`
		TemperatureRequired         bool       `db:"temperature_required"`
		TimeEstimateRequired        bool       `db:"time_estimate_required"`
		ConditionExpressionRequired bool       `db:"condition_expression_required"`
		ConsumesVessel              bool       `db:"consumes_vessel"`
		OnlyForVessels              bool       `db:"only_for_vessels"`
		YieldsNothing               bool       `db:"yields_nothing"`
	}
)

// CreateValidPreparation gets a valid preparation from the database.
func (c *DatabaseClient) CreateValidPreparation(ctx context.Context, input *ValidPreparation) (*types.ValidPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Insert(validPreparationsTableName).Rows(
		input,
	)

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "creating valid preparation")
	}

	var output types.ValidPreparation
	if err := copier.Copy(&output, input); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetValidPreparation gets a valid preparation from the database.
func (c *DatabaseClient) GetValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	x := &ValidPreparation{}
	q := c.xdb.From(validPreparationsTableName).Where(goqu.Ex{
		idColumn:         validPreparationID,
		archivedAtColumn: nil,
	})

	found, err := q.ScanStructContext(ctx, x)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, sql.ErrNoRows
	}

	var output types.ValidPreparation
	if err = copier.Copy(&output, x); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetValidPreparations gets a valid preparation from the database.
func (c *DatabaseClient) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	q := c.xdb.From(validPreparationsTableName).Where(goqu.Ex{
		archivedAtColumn: nil,
	})
	q = queryFilterToGoqu(q, filter)

	var x []ValidPreparation
	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	output := &types.QueryFilteredResult[types.ValidPreparation]{
		Data:       []*types.ValidPreparation{},
		Pagination: filter.ToPagination(),
	}
	for _, y := range x {
		var z types.ValidPreparation
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output.Data = append(output.Data, &z)
	}

	return output, nil
}

// SearchForValidPreparations gets a valid preparation from the database.
func (c *DatabaseClient) SearchForValidPreparations(ctx context.Context, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	q := c.xdb.From(validPreparationsTableName).Where(goqu.Ex{
		archivedAtColumn: nil,
		"name":           goqu.Op{"like": "%" + query + "%"},
	})
	q = queryFilterToGoqu(q, filter)

	var x []ValidPreparation
	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	output := &types.QueryFilteredResult[types.ValidPreparation]{
		Data:       []*types.ValidPreparation{},
		Pagination: filter.ToPagination(),
	}
	for _, y := range x {
		var z types.ValidPreparation
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output.Data = append(output.Data, &z)
	}

	return output, nil
}

// GetValidPreparationsWithIDs gets a valid preparation from the database.
func (c *DatabaseClient) GetValidPreparationsWithIDs(ctx context.Context, ids []string) ([]*types.ValidPreparation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	x := []*ValidPreparation{}

	q := c.xdb.From(validPreparationsTableName).Where(goqu.Ex{
		idColumn:         ids,
		archivedAtColumn: nil,
	})

	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	var output []*types.ValidPreparation
	for _, y := range x {
		var z types.ValidPreparation
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output = append(output, &z)
	}

	return output, nil
}

// UpdateValidPreparation gets a valid preparation from the database.
func (c *DatabaseClient) UpdateValidPreparation(ctx context.Context, input *types.ValidPreparation) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	var updateInput ValidPreparation
	if err := copier.Copy(&updateInput, input); err != nil {
		return observability.PrepareError(err, span, "copying input to output")
	}

	q := c.xdb.Update(validPreparationsTableName).Set(
		updateInput,
	).Set(goqu.Ex{lastUpdatedAtColumn: goqu.L("NOW()")})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "updating valid preparation")
	}

	return nil
}

// ArchiveValidPreparation gets a valid preparation from the database.
func (c *DatabaseClient) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Update(validPreparationsTableName).
		Set(goqu.Record{archivedAtColumn: goqu.L("NOW()")}).
		Where(goqu.Ex{idColumn: validPreparationID})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "archiving valid preparation")
	}

	return nil
}
