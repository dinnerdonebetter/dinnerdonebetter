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
	usersTableName = "users"
)

type (
	// User represents a user.
	User struct {
		_ struct{}

		CreatedAt                 time.Time  `db:"created_at"                      goqu:"skipinsert"`
		LastUpdatedAt             *time.Time `db:"last_updated_at"                 goqu:"skipinsert"`
		ArchivedAt                *time.Time `db:"archived_at"                     goqu:"skipinsert"`
		PasswordLastChangedAt     *time.Time `db:"password_last_changed_at"`
		LastAcceptedTOS           *time.Time `db:"last_accepted_terms_of_service"`
		LastAcceptedPrivacyPolicy *time.Time `db:"last_accepted_privacy_policy"`
		TwoFactorSecretVerifiedAt *time.Time `db:"two_factor_secret_verified_at"`
		AvatarSrc                 *string    `db:"avatar_src"`
		Birthday                  *time.Time `db:"birthday"`
		AccountStatusExplanation  string     `db:"user_account_status_explanation"`
		TwoFactorSecret           string     `db:"two_factor_secret"`
		HashedPassword            string     `db:"hashed_password"`
		ID                        string     `db:"id"                              goqu:"skipupdate"`
		AccountStatus             string     `db:"user_account_status"`
		Username                  string     `db:"username"`
		FirstName                 string     `db:"first_name"`
		LastName                  string     `db:"last_name"`
		EmailAddress              string     `db:"email_address"`
		EmailAddressVerifiedAt    *time.Time `db:"email_address_verified_at"`
		ServiceRole               string     `db:"service_role"`
		RequiresPasswordChange    bool       `db:"requires_password_change"`
	}
)

// CreateUser gets a user from the database.
func (c *DatabaseClient) CreateUser(ctx context.Context, input *User) (*types.User, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Insert(usersTableName).Rows(
		input,
	)

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "creating user")
	}

	var output types.User
	if err := copier.Copy(&output, input); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetUser gets a user from the database.
func (c *DatabaseClient) GetUser(ctx context.Context, userID string) (*types.User, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	x := &User{}
	q := c.xdb.From(usersTableName).Where(goqu.Ex{
		idColumn:         userID,
		archivedAtColumn: nil,
	})

	found, err := q.ScanStructContext(ctx, x)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, sql.ErrNoRows
	}

	var output types.User
	if err = copier.Copy(&output, x); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetUsers gets a user from the database.
func (c *DatabaseClient) GetUsers(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.User], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	q := c.xdb.From(usersTableName).Where(goqu.Ex{
		archivedAtColumn: nil,
	})
	q = queryFilterToGoqu(q, filter)

	var x []User
	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	output := &types.QueryFilteredResult[types.User]{
		Data:       []*types.User{},
		Pagination: filter.ToPagination(),
	}
	for _, y := range x {
		var z types.User
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output.Data = append(output.Data, &z)
	}

	return output, nil
}

// SearchForUsers gets a user from the database.
func (c *DatabaseClient) SearchForUsers(ctx context.Context, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.User], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	q := c.xdb.From(usersTableName).Where(goqu.Ex{
		archivedAtColumn: nil,
		"username":       goqu.Op{"like": "%" + query + "%"},
	})
	q = queryFilterToGoqu(q, filter)

	var x []User
	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	output := &types.QueryFilteredResult[types.User]{
		Data:       []*types.User{},
		Pagination: filter.ToPagination(),
	}
	for _, y := range x {
		var z types.User
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output.Data = append(output.Data, &z)
	}

	return output, nil
}

// GetUsersWithIDs gets a user from the database.
func (c *DatabaseClient) GetUsersWithIDs(ctx context.Context, ids []string) ([]*types.User, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	x := []*User{}

	q := c.xdb.From(usersTableName).Where(goqu.Ex{
		idColumn:         ids,
		archivedAtColumn: nil,
	})

	if err := q.ScanStructsContext(ctx, &x); err != nil {
		return nil, err
	}

	var output []*types.User
	for _, y := range x {
		var z types.User
		if err := copier.Copy(&z, y); err != nil {
			return nil, observability.PrepareError(err, span, "copying input to output")
		}

		output = append(output, &z)
	}

	return output, nil
}

// ArchiveUser gets a user from the database.
func (c *DatabaseClient) ArchiveUser(ctx context.Context, userID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Update(usersTableName).
		Set(goqu.Record{archivedAtColumn: goqu.L("NOW()")}).
		Where(goqu.Ex{idColumn: userID})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "archiving user")
	}

	return nil
}
