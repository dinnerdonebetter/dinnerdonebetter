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
	householdsTableName = "households"
)

type (
	// Household represents a household.
	Household struct {
		_ struct{}

		CreatedAt                  time.Time  `db:"created_at"                    goqu:"skipinsert"`
		LastUpdatedAt              *time.Time `db:"last_updated_at"               goqu:"skipinsert"`
		ArchivedAt                 *time.Time `db:"archived_at"                   goqu:"skipinsert"`
		ID                         string     `db:"id"                            goqu:"skipupdate"`
		SubscriptionPlanID         *string    `db:"subscription_plan_id"`
		ContactPhone               string     `db:"contact_phone"`
		BillingStatus              string     `db:"billing_status"`
		AddressLine1               string     `db:"address_line_1"`
		AddressLine2               string     `db:"address_line_2"`
		City                       string     `db:"city"`
		State                      string     `db:"state"`
		ZipCode                    string     `db:"zip_code"`
		Country                    string     `db:"country"`
		Latitude                   *float64   `db:"latitude"`
		Longitude                  *float64   `db:"longitude"`
		PaymentProcessorCustomerID string     `db:"payment_processor_customer_id"`
		BelongsToUser              string     `db:"belongs_to_user"`
		Name                       string     `db:"name"`
		// Members                    []*HouseholdUserMembershipWithUser `json:"members"`
	}
)

// CreateHousehold gets a household from the database.
func (c *DatabaseClient) CreateHousehold(ctx context.Context, input *Household) (*types.Household, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Insert(householdsTableName).Rows(
		input,
	)

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "creating household")
	}

	var output types.Household
	if err := copier.Copy(&output, input); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetHousehold gets a household from the database.
func (c *DatabaseClient) GetHousehold(ctx context.Context, householdID string) (*types.Household, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	x := &Household{}
	q := c.xdb.From(householdsTableName).Where(goqu.Ex{
		idColumn:         householdID,
		archivedAtColumn: nil,
	})

	found, err := q.ScanStructContext(ctx, x)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, sql.ErrNoRows
	}

	var output types.Household
	if err = copier.Copy(&output, x); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetHouseholds gets a household from the database.
func (c *DatabaseClient) GetHouseholds(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Household], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	q := c.xdb.From(householdsTableName).Where(goqu.Ex{
		archivedAtColumn: nil,
	})
	q = queryFilterToGoqu(q, filter)

	var x []Household
	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	output := &types.QueryFilteredResult[types.Household]{
		Data:       []*types.Household{},
		Pagination: filter.ToPagination(),
	}
	for _, y := range x {
		var z types.Household
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output.Data = append(output.Data, &z)
	}

	return output, nil
}

// UpdateHousehold gets a household from the database.
func (c *DatabaseClient) UpdateHousehold(ctx context.Context, input *types.Household) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	var updateInput Household
	if err := copier.Copy(&updateInput, input); err != nil {
		return observability.PrepareError(err, span, "copying input to output")
	}

	q := c.xdb.Update(householdsTableName).Set(
		updateInput,
	).Set(goqu.Ex{lastUpdatedAtColumn: goqu.L("NOW()")})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "updating household")
	}

	return nil
}

// ArchiveHousehold gets a household from the database.
func (c *DatabaseClient) ArchiveHousehold(ctx context.Context, householdID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Update(householdsTableName).
		Set(goqu.Record{archivedAtColumn: goqu.L("NOW()")}).
		Where(goqu.Ex{idColumn: householdID})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "archiving household")
	}

	return nil
}
