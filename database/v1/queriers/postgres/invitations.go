package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/Masterminds/squirrel"
)

const (
	invitationsTableName           = "invitations"
	invitationsUserOwnershipColumn = "belongs_to_user"
)

var (
	invitationsTableColumns = []string{
		fmt.Sprintf("%s.%s", invitationsTableName, "id"),
		fmt.Sprintf("%s.%s", invitationsTableName, "code"),
		fmt.Sprintf("%s.%s", invitationsTableName, "consumed"),
		fmt.Sprintf("%s.%s", invitationsTableName, "created_on"),
		fmt.Sprintf("%s.%s", invitationsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", invitationsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", invitationsTableName, invitationsUserOwnershipColumn),
	}
)

// scanInvitation takes a database Scanner (i.e. *sql.Row) and scans the result into an Invitation struct
func (p *Postgres) scanInvitation(scan database.Scanner, includeCount bool) (*models.Invitation, uint64, error) {
	x := &models.Invitation{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Code,
		&x.Consumed,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}

// scanInvitations takes a logger and some database rows and turns them into a slice of invitations.
func (p *Postgres) scanInvitations(rows database.ResultIterator) ([]models.Invitation, uint64, error) {
	var (
		list  []models.Invitation
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanInvitation(rows, true)
		if err != nil {
			return nil, 0, err
		}

		if count == 0 {
			count = c
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		p.logger.Error(closeErr, "closing database rows")
	}

	return list, count, nil
}

//
func (p *Postgres) buildInvitationExistsQuery(invitationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", invitationsTableName)).
		Prefix(existencePrefix).
		From(invitationsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", invitationsTableName): invitationID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// InvitationExists queries the database to see if a given invitation belonging to a given user exists.
func (p *Postgres) InvitationExists(ctx context.Context, invitationID uint64) (exists bool, err error) {
	query, args := p.buildInvitationExistsQuery(invitationID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

//
func (p *Postgres) buildGetInvitationQuery(invitationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(invitationsTableColumns...).
		From(invitationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", invitationsTableName): invitationID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetInvitation fetches an invitation from the database.
func (p *Postgres) GetInvitation(ctx context.Context, invitationID uint64) (*models.Invitation, error) {
	query, args := p.buildGetInvitationQuery(invitationID)
	row := p.db.QueryRowContext(ctx, query, args...)

	invitation, _, err := p.scanInvitation(row, false)
	return invitation, err
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
			Select(fmt.Sprintf(countQuery, invitationsTableName)).
			From(invitationsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", invitationsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allInvitationsCountQuery
}

// GetAllInvitationsCount will fetch the count of invitations from the database.
func (p *Postgres) GetAllInvitationsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllInvitationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetInvitationsQuery builds a SQL query selecting invitations that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetInvitationsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(invitationsTableColumns, fmt.Sprintf(countQuery, invitationsTableName))...).
		From(invitationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", invitationsTableName): nil,
		}).
		GroupBy(fmt.Sprintf("%s.id", invitationsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, invitationsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetInvitations fetches a list of invitations from the database that meet a particular filter.
func (p *Postgres) GetInvitations(ctx context.Context, filter *models.QueryFilter) (*models.InvitationList, error) {
	query, args := p.buildGetInvitationsQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for invitations")
	}

	invitations, count, err := p.scanInvitations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.InvitationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Invitations: invitations,
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
			invitationsUserOwnershipColumn,
		).
		Values(
			input.Code,
			input.Consumed,
			input.BelongsToUser,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateInvitation creates an invitation in the database.
func (p *Postgres) CreateInvitation(ctx context.Context, input *models.InvitationCreationInput) (*models.Invitation, error) {
	x := &models.Invitation{
		Code:          input.Code,
		Consumed:      input.Consumed,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := p.buildCreateInvitationQuery(x)

	// create the invitation.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing invitation creation query: %w", err)
	}

	return x, nil
}

// buildUpdateInvitationQuery takes an invitation and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateInvitationQuery(input *models.Invitation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(invitationsTableName).
		Set("code", input.Code).
		Set("consumed", input.Consumed).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                           input.ID,
			invitationsUserOwnershipColumn: input.BelongsToUser,
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
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                           invitationID,
			"archived_on":                  nil,
			invitationsUserOwnershipColumn: userID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveInvitation marks an invitation as archived in the database.
func (p *Postgres) ArchiveInvitation(ctx context.Context, invitationID, userID uint64) error {
	query, args := p.buildArchiveInvitationQuery(invitationID, userID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
