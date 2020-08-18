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
	invitationsTableCodeColumn     = "code"
	invitationsTableConsumedColumn = "consumed"
	invitationsUserOwnershipColumn = "belongs_to_user"
)

var (
	invitationsTableColumns = []string{
		fmt.Sprintf("%s.%s", invitationsTableName, idColumn),
		fmt.Sprintf("%s.%s", invitationsTableName, invitationsTableCodeColumn),
		fmt.Sprintf("%s.%s", invitationsTableName, invitationsTableConsumedColumn),
		fmt.Sprintf("%s.%s", invitationsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", invitationsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", invitationsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", invitationsTableName, invitationsUserOwnershipColumn),
	}
)

// scanInvitation takes a database Scanner (i.e. *sql.Row) and scans the result into an Invitation struct
func (p *Postgres) scanInvitation(scan database.Scanner) (*models.Invitation, error) {
	x := &models.Invitation{}

	targetVars := []interface{}{
		&x.ID,
		&x.Code,
		&x.Consumed,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanInvitations takes a logger and some database rows and turns them into a slice of invitations.
func (p *Postgres) scanInvitations(rows database.ResultIterator) ([]models.Invitation, error) {
	var (
		list []models.Invitation
	)

	for rows.Next() {
		x, err := p.scanInvitation(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		p.logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

//
func (p *Postgres) buildInvitationExistsQuery(invitationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", invitationsTableName, idColumn)).
		Prefix(existencePrefix).
		From(invitationsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", invitationsTableName, idColumn): invitationID,
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
			fmt.Sprintf("%s.%s", invitationsTableName, idColumn): invitationID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetInvitation fetches an invitation from the database.
func (p *Postgres) GetInvitation(ctx context.Context, invitationID uint64) (*models.Invitation, error) {
	query, args := p.buildGetInvitationQuery(invitationID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanInvitation(row)
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
				fmt.Sprintf("%s.%s", invitationsTableName, archivedOnColumn): nil,
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

// buildGetBatchOfInvitationsQuery returns a query that fetches every invitation in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfInvitationsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(invitationsTableColumns...).
		From(invitationsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", invitationsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", invitationsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllInvitations fetches every invitation from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllInvitations(ctx context.Context, resultChannel chan []models.Invitation) error {
	count, err := p.GetAllInvitationsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfInvitationsQuery(begin, end)
			logger := p.logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, err := p.db.Query(query, args...)
			if err == sql.ErrNoRows {
				return
			} else if err != nil {
				logger.Error(err, "querying for database rows")
				return
			}

			invitations, err := p.scanInvitations(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- invitations
		}(beginID, endID)
	}

	return nil
}

// buildGetInvitationsQuery builds a SQL query selecting invitations that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetInvitationsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(invitationsTableColumns...).
		From(invitationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", invitationsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", invitationsTableName, idColumn))

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

	invitations, err := p.scanInvitations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.InvitationList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Invitations: invitations,
	}

	return list, nil
}

// buildGetInvitationsWithIDsQuery builds a SQL query selecting invitations
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetInvitationsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(invitationsTableColumns...).
		From(invitationsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(invitationsTableColumns...).
		FromSelect(subqueryBuilder, invitationsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", invitationsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetInvitationsWithIDs fetches a list of invitations from the database that exist within a given set of IDs.
func (p *Postgres) GetInvitationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.Invitation, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetInvitationsWithIDsQuery(limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for invitations")
	}

	invitations, err := p.scanInvitations(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return invitations, nil
}

// buildCreateInvitationQuery takes an invitation and returns a creation query for that invitation and the relevant arguments.
func (p *Postgres) buildCreateInvitationQuery(input *models.Invitation) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(invitationsTableName).
		Columns(
			invitationsTableCodeColumn,
			invitationsTableConsumedColumn,
			invitationsUserOwnershipColumn,
		).
		Values(
			input.Code,
			input.Consumed,
			input.BelongsToUser,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
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
		Set(invitationsTableCodeColumn, input.Code).
		Set(invitationsTableConsumedColumn, input.Consumed).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                       input.ID,
			invitationsUserOwnershipColumn: input.BelongsToUser,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateInvitation updates a particular invitation. Note that UpdateInvitation expects the provided input to have a valid ID.
func (p *Postgres) UpdateInvitation(ctx context.Context, input *models.Invitation) error {
	query, args := p.buildUpdateInvitationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveInvitationQuery returns a SQL query which marks a given invitation belonging to a given user as archived.
func (p *Postgres) buildArchiveInvitationQuery(invitationID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(invitationsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                       invitationID,
			archivedOnColumn:               nil,
			invitationsUserOwnershipColumn: userID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
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
