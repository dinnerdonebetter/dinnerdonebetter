package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	invitationsTableName = "invitations"
)

var (
	invitationsTableColumns = []string{
		"id",
		"code",
		"consumed",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanInvitation takes a database Scanner (i.e. *sql.Row) and scans the result into an Invitation struct
func scanInvitation(scan database.Scanner) (*models.Invitation, error) {
	x := &models.Invitation{}

	if err := scan.Scan(
		&x.ID,
		&x.Code,
		&x.Consumed,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanInvitations takes a logger and some database rows and turns them into a slice of invitations
func scanInvitations(logger logging.Logger, rows *sql.Rows) ([]models.Invitation, error) {
	var list []models.Invitation

	for rows.Next() {
		x, err := scanInvitation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildGetInvitationQuery constructs a SQL query for fetching an invitation with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetInvitationQuery(invitationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(invitationsTableColumns...).
		From(invitationsTableName).
		Where(squirrel.Eq{
			"id":         invitationID,
			"belongs_to": userID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetInvitation fetches an invitation from the postgres database
func (p *Postgres) GetInvitation(ctx context.Context, invitationID, userID uint64) (*models.Invitation, error) {
	query, args := p.buildGetInvitationQuery(invitationID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanInvitation(row)
}

// buildGetInvitationCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of invitations belonging to a given user that meet a given query
func (p *Postgres) buildGetInvitationCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(invitationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetInvitationCount will fetch the count of invitations from the database that meet a particular filter and belong to a particular user.
func (p *Postgres) GetInvitationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := p.buildGetInvitationCountQuery(filter, userID)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allInvitationsCountQueryBuilder sync.Once
	allInvitationsCountQuery        string
)

// buildGetAllInvitationsCountQuery returns a query that fetches the total number of invitations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllInvitationsCountQuery() string {
	allInvitationsCountQueryBuilder.Do(func() {
		var err error
		allInvitationsCountQuery, _, err = p.sqlBuilder.
			Select(CountQuery).
			From(invitationsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allInvitationsCountQuery
}

// GetAllInvitationsCount will fetch the count of invitations from the database
func (p *Postgres) GetAllInvitationsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllInvitationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetInvitationsQuery builds a SQL query selecting invitations that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetInvitationsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(invitationsTableColumns...).
		From(invitationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetInvitations fetches a list of invitations from the database that meet a particular filter
func (p *Postgres) GetInvitations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.InvitationList, error) {
	query, args := p.buildGetInvitationsQuery(filter, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for invitations")
	}

	list, err := scanInvitations(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetInvitationCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching invitation count: %w", err)
	}

	x := &models.InvitationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Invitations: list,
	}

	return x, nil
}

// GetAllInvitationsForUser fetches every invitation belonging to a user
func (p *Postgres) GetAllInvitationsForUser(ctx context.Context, userID uint64) ([]models.Invitation, error) {
	query, args := p.buildGetInvitationsQuery(nil, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching invitations for user")
	}

	list, err := scanInvitations(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateInvitationQuery takes an invitation and returns a creation query for that invitation and the relevant arguments.
func (p *Postgres) buildCreateInvitationQuery(input *models.Invitation) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(invitationsTableName).
		Columns(
			"code",
			"consumed",
			"belongs_to",
		).
		Values(
			input.Code,
			input.Consumed,
			input.BelongsTo,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateInvitation creates an invitation in the database
func (p *Postgres) CreateInvitation(ctx context.Context, input *models.InvitationCreationInput) (*models.Invitation, error) {
	x := &models.Invitation{
		Code:      input.Code,
		Consumed:  input.Consumed,
		BelongsTo: input.BelongsTo,
	}

	query, args := p.buildCreateInvitationQuery(x)

	// create the invitation
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing invitation creation query: %w", err)
	}

	return x, nil
}

// buildUpdateInvitationQuery takes an invitation and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdateInvitationQuery(input *models.Invitation) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(invitationsTableName).
		Set("code", input.Code).
		Set("consumed", input.Consumed).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateInvitation updates a particular invitation. Note that UpdateInvitation expects the provided input to have a valid ID.
func (p *Postgres) UpdateInvitation(ctx context.Context, input *models.Invitation) error {
	query, args := p.buildUpdateInvitationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveInvitationQuery returns a SQL query which marks a given invitation belonging to a given user as archived.
func (p *Postgres) buildArchiveInvitationQuery(invitationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(invitationsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          invitationID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveInvitation marks an invitation as archived in the database
func (p *Postgres) ArchiveInvitation(ctx context.Context, invitationID, userID uint64) error {
	query, args := p.buildArchiveInvitationQuery(invitationID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
