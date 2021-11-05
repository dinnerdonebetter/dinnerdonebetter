package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// householdInvitationsTableName is what the household invitations table calls itself.
	householdInvitationsTableName = "household_invitations"
)

var (
	_ types.HouseholdInvitationDataManager = (*SQLQuerier)(nil)

	householdInvitationsTableColumns = []string{
		"household_invitations.id",
		"household_invitations.destination_household",
		"household_invitations.to_user",
		"household_invitations.from_user",
		"household_invitations.status",
		"household_invitations.created_on",
		"household_invitations.last_updated_on",
		"household_invitations.archived_on",
	}
)

// scanHousehold takes a database Scanner (i.e. *sql.Row) and scans the result into a household struct.
func (q *SQLQuerier) scanHouseholdInvitation(ctx context.Context, scan database.Scanner, includeCounts bool) (householdInvitation *types.HouseholdInvitation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	householdInvitation = &types.HouseholdInvitation{}

	targetVars := []interface{}{
		&householdInvitation.ID,
		&householdInvitation.DestinationHousehold,
		&householdInvitation.ToUser,
		&householdInvitation.FromUser,
		&householdInvitation.Status,
		&householdInvitation.CreatedOn,
		&householdInvitation.LastUpdatedOn,
		&householdInvitation.ArchivedOn,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "fetching memberships from database")
	}

	return householdInvitation, filteredCount, totalCount, nil
}

// scanHouseholdInvitations takes some database rows and turns them into a slice of household invitations.
func (q *SQLQuerier) scanHouseholdInvitations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (householdInvitations []*types.HouseholdInvitation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	householdInvitations = []*types.HouseholdInvitation{}

	var currentHousehold *types.HouseholdInvitation
	for rows.Next() {
		household, fc, tc, scanErr := q.scanHouseholdInvitation(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if currentHousehold == nil {
			currentHousehold = household
		}

		if currentHousehold.ID != household.ID {
			householdInvitations = append(householdInvitations, currentHousehold)
			currentHousehold = household
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}
	}

	if currentHousehold != nil {
		householdInvitations = append(householdInvitations, currentHousehold)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return householdInvitations, filteredCount, totalCount, nil
}

const getHouseholdInvitationQuery = `
	SELECT
		household_invitations.id,
		household_invitations.destination_household,
		household_invitations.to_user,
		household_invitations.from_user,
		household_invitations.status,
		household_invitations.created_on,
		household_invitations.last_updated_on,
		household_invitations.archived_on,
	FROM household_invitations
	WHERE household_invitations.archived_on IS NULL
	AND household_invitations.id = $1
`

// GetHouseholdInvitation fetches a household from the database.
func (q *SQLQuerier) GetHouseholdInvitation(ctx context.Context, householdInvitationID, userID string) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdInvitationID == "" || userID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachHouseholdIDToSpan(span, householdInvitationID)
	tracing.AttachUserIDToSpan(span, userID)

	logger := q.logger.WithValues(map[string]interface{}{
		keys.HouseholdIDKey: householdInvitationID,
		keys.UserIDKey:      userID,
	})

	args := []interface{}{
		userID,
		householdInvitationID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "household", getHouseholdInvitationQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing household invitations list retrieval query")
	}

	householdInvitations, _, _, err := q.scanHouseholdInvitations(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	var household *types.HouseholdInvitation
	if len(householdInvitations) > 0 {
		household = householdInvitations[0]
	}

	if household == nil {
		return nil, sql.ErrNoRows
	}

	return household, nil
}

const getAllHouseholdInvitationsCountQuery = `
	SELECT COUNT(household_invitationsid) FROM household_invitations WHERE household_invitationsarchived_on IS NULL
`

// GetAllHouseholdInvitationsCount fetches the count of household invitations from the database that meet a particular filter.
func (q *SQLQuerier) GetAllHouseholdInvitationsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, getAllHouseholdInvitationsCountQuery, "fetching count of all household invitations")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of household invitations")
	}

	return count, nil
}

// buildGetHouseholdInvitationsQuery builds a SQL query selecting household invitations that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (q *SQLQuerier) buildGetHouseholdInvitationsQuery(ctx context.Context, userID string, forAdmin bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	var includeArchived bool
	if filter != nil {
		includeArchived = filter.IncludeArchived
	}

	filteredCountQuery, filteredCountQueryArgs := q.buildFilteredCountQuery(ctx, householdInvitationsTableName, nil, nil, userOwnershipColumn, userID, forAdmin, includeArchived, filter)
	totalCountQuery, totalCountQueryArgs := q.buildTotalCountQuery(ctx, householdInvitationsTableName, nil, nil, userOwnershipColumn, userID, forAdmin, includeArchived)

	builder := q.sqlBuilder.Select(append(
		householdInvitationsTableColumns,
		fmt.Sprintf("(%s) as total_count", totalCountQuery),
		fmt.Sprintf("(%s) as filtered_count", filteredCountQuery),
	)...).
		From(householdInvitationsTableName)

	if !forAdmin {
		where := squirrel.Eq{
			"household_invitations.archived_on": nil,
		}

		if userID != "" {
			where["household_invitations.belongs_to_user"] = userID
		}

		builder = builder.Where(where)
	}

	builder = builder.GroupBy(fmt.Sprintf(
		"%s.%s",
		householdInvitationsTableName,
		"id",
	))

	if filter != nil {
		builder = applyFilterToQueryBuilder(filter, householdInvitationsTableName, builder)
	}

	query, selectArgs := q.buildQuery(span, builder)

	return query, append(append(filteredCountQueryArgs, totalCountQueryArgs...), selectArgs...)
}

// GetHouseholdInvitations fetches a list of household invitations from the database that meet a particular filter.
func (q *SQLQuerier) GetHouseholdInvitations(ctx context.Context, userID string, filter *types.QueryFilter) (x *types.HouseholdInvitationList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := filter.AttachToLogger(q.logger).WithValue(keys.UserIDKey, userID)
	tracing.AttachQueryFilterToSpan(span, filter)
	tracing.AttachUserIDToSpan(span, userID)

	x = &types.HouseholdInvitationList{}
	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildGetHouseholdInvitationsQuery(ctx, userID, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "household invitations", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing household invitations list retrieval query")
	}

	if x.HouseholdInvitations, x.FilteredCount, x.TotalCount, err = q.scanHouseholdInvitations(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning household invitations from database")
	}

	return x, nil
}

const householdInvitationCreationQuery = `
	INSERT INTO household_invitations (id,name,billing_status,contact_email,contact_phone,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6)
`

// CreateHouseholdInvitation creates a household in the database.
func (q *SQLQuerier) CreateHouseholdInvitation(ctx context.Context, input *types.HouseholdInvitationCreationInput) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, input.BelongsToUser)

	// begin household creation transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	householdCreationArgs := []interface{}{
		input.ID,
		input.BelongsToUser,
	}

	// create the household.
	if writeErr := q.performWriteQuery(ctx, tx, "household creation", householdInvitationCreationQuery, householdCreationArgs); writeErr != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(writeErr, logger, span, "creating household")
	}

	household := &types.HouseholdInvitation{
		ID:        input.ID,
		CreatedOn: q.currentTime(),
	}

	addInput := &types.AddUserToHouseholdInput{
		ID:             ksuid.New().String(),
		UserID:         input.BelongsToUser,
		HouseholdID:    household.ID,
		HouseholdRoles: []string{authorization.HouseholdAdminRole.String()},
	}

	addUserToHouseholdArgs := []interface{}{
		addInput.ID,
		addInput.UserID,
		addInput.HouseholdID,
		strings.Join(addInput.HouseholdRoles, householdMemberRolesSeparator),
	}

	if err = q.performWriteQuery(ctx, tx, "household user membership creation", addUserToHouseholdDuringCreationQuery, addUserToHouseholdArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "performing household membership creation query")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachHouseholdIDToSpan(span, household.ID)
	logger.Info("household created")

	return household, nil
}

const archiveHouseholdInvitationQuery = `
	UPDATE household_invitations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2
`

// ArchiveHouseholdInvitation archives a household from the database by its ID.
func (q *SQLQuerier) ArchiveHouseholdInvitation(ctx context.Context, householdInvitationID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdInvitationID == "" || userID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdInvitationID)

	logger := q.logger.WithValues(map[string]interface{}{
		keys.HouseholdIDKey: householdInvitationID,
		keys.UserIDKey:      userID,
	})

	args := []interface{}{
		userID,
		householdInvitationID,
	}

	if err := q.performWriteQuery(ctx, q.db, "household archive", archiveHouseholdInvitationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "archiving household")
	}

	logger.Info("household archived")

	return nil
}
