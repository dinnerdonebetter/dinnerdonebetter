package postgres

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.HouseholdInvitationDataManager = (*SQLQuerier)(nil)

	// householdInvitationsTableColumns are the columns for the household invitations table.
	householdInvitationsTableColumns = []string{
		"household_invitations.id",
		"household_invitations.destination_household",
		"household_invitations.to_email",
		"household_invitations.to_user",
		"household_invitations.from_user",
		"household_invitations.status",
		"household_invitations.note",
		"household_invitations.token",
		"household_invitations.created_on",
		"household_invitations.last_updated_on",
		"household_invitations.archived_on",
	}
)

// scanHouseholdInvitation is a consistent way to turn a *sql.Row into a webhook struct.
func (q *SQLQuerier) scanHouseholdInvitation(ctx context.Context, scan database.Scanner, includeCounts bool) (householdInvitation *types.HouseholdInvitation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)
	householdInvitation = &types.HouseholdInvitation{}

	targetVars := []interface{}{
		&householdInvitation.ID,
		&householdInvitation.DestinationHousehold,
		&householdInvitation.ToEmail,
		&householdInvitation.ToUser,
		&householdInvitation.FromUser,
		&householdInvitation.Status,
		&householdInvitation.Note,
		&householdInvitation.Token,
		&householdInvitation.CreatedOn,
		&householdInvitation.LastUpdatedOn,
		&householdInvitation.ArchivedOn,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "scanning householdInvitation")
	}

	return householdInvitation, filteredCount, totalCount, nil
}

// scanHouseholdInvitations provides a consistent way to turn sql rows into a slice of household_invitations.
func (q *SQLQuerier) scanHouseholdInvitations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (householdInvitations []*types.HouseholdInvitation, filteredCount, totalCount uint64, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		householdInvitation, fc, tc, scanErr := q.scanHouseholdInvitation(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		householdInvitations = append(householdInvitations, householdInvitation)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "fetching webhook from database")
	}

	if err = rows.Close(); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "fetching webhook from database")
	}

	return householdInvitations, filteredCount, totalCount, nil
}

const householdInvitationExistenceQuery = "SELECT EXISTS ( SELECT household_invitations.id FROM household_invitations WHERE household_invitations.archived_on IS NULL AND household_invitations.id = $1 AND household_invitations.destination_household = $2 )"

// HouseholdInvitationExists fetches whether a household invitation exists from the database.
func (q *SQLQuerier) HouseholdInvitationExists(ctx context.Context, householdID, householdInvitationID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	if householdInvitationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)

	args := []interface{}{
		householdInvitationID,
		householdID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, householdInvitationExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing webhook existence check")
	}

	return result, nil
}

const getHouseholdInvitationQuery = `
SELECT
	household_invitations.id,
	household_invitations.destination_household,
	household_invitations.to_email,
	household_invitations.to_user,
	household_invitations.from_user,
	household_invitations.status,
	household_invitations.note,
	household_invitations.token,
	household_invitations.created_on,
	household_invitations.last_updated_on,
	household_invitations.archived_on
FROM household_invitations 
WHERE household_invitations.archived_on IS NULL
AND household_invitations.id = $1
`

// GetHouseholdInvitation fetches a webhook from the database.
func (q *SQLQuerier) GetHouseholdInvitation(ctx context.Context, householdID, householdInvitationID string) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	if householdInvitationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)

	args := []interface{}{
		householdInvitationID,
	}

	row := q.getOneRow(ctx, q.db, "webhook", getHouseholdInvitationQuery, args)

	webhook, _, _, err := q.scanHouseholdInvitation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning webhook")
	}

	return webhook, nil
}

const getAllHouseholdInvitationsCountQuery = `
	SELECT COUNT(household_invitations.id) FROM household_invitations WHERE household_invitations.archived_on IS NULL
`

// GetAllHouseholdInvitationsCount fetches the count of household invitations from the database that meet a particular filter.
func (q *SQLQuerier) GetAllHouseholdInvitationsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, getAllHouseholdInvitationsCountQuery, "fetching count of household invitations")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of household invitations")
	}

	return count, nil
}

// GetHouseholdInvitations fetches a list of household invitations from the database that meet a particular filter.
func (q *SQLQuerier) GetHouseholdInvitations(ctx context.Context, householdInvitationID string, filter *types.QueryFilter) (*types.HouseholdInvitationList, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdInvitationID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdIDToSpan(span, householdInvitationID)
	tracing.AttachQueryFilterToSpan(span, filter)

	x := &types.HouseholdInvitationList{}
	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "household_invitations", nil, nil, nil, "belongs_to_household", householdInvitationsTableColumns, householdInvitationID, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "household invitations", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching webhook from database")
	}

	if x.HouseholdInvitations, x.FilteredCount, x.TotalCount, err = q.scanHouseholdInvitations(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning database response")
	}

	return x, nil
}

const createHouseholdInvitationQuery = `
	INSERT INTO household_invitations (id,from_user,to_user,note,to_email,token,destination_household) VALUES ($1,$2,$3,$4,$5,$6,$7)
`

// CreateHouseholdInvitation creates a webhook in a database.
func (q *SQLQuerier) CreateHouseholdInvitation(ctx context.Context, input *types.HouseholdInvitationDatabaseCreationInput) (*types.HouseholdInvitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.HouseholdInvitationIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.FromUser,
		input.ToUser,
		input.Note,
		input.ToEmail,
		input.Token,
		input.DestinationHousehold,
	}

	if err := q.performWriteQuery(ctx, q.db, "webhook creation", createHouseholdInvitationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing webhook creation query")
	}

	x := &types.HouseholdInvitation{
		ID:        input.ID,
		CreatedOn: q.currentTime(),
	}

	tracing.AttachHouseholdInvitationIDToSpan(span, x.ID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, x.ID)

	logger.Info("webhook created")

	return x, nil
}

func (q *SQLQuerier) GetSentPendingHouseholdInvitations(ctx context.Context, userID string, filter *types.QueryFilter) ([]*types.HouseholdInvitation, error) {
	return nil, nil
}

func (q *SQLQuerier) GetReceivedPendingHouseholdInvitations(ctx context.Context, userID string, filter *types.QueryFilter) ([]*types.HouseholdInvitation, error) {
	return nil, nil
}

const cancelHouseholdInvitationQuery = `
UPDATE household_invitations SET
	status = 'cancelled',
	last_updated_on = extract(epoch FROM NOW()), 
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL 
AND belongs_to_household = $1
AND id = $2
`

func (q *SQLQuerier) CancelHouseholdInvitation(ctx context.Context, householdInvitationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)

	args := []interface{}{householdInvitationID}

	if err := q.performWriteQuery(ctx, q.db, "household invitation cancel", cancelHouseholdInvitationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "cancelling household invitation")
	}

	logger.Info("household invitation cancelled")

	return nil
}

const acceptHouseholdInvitationQuery = `
UPDATE household_invitations SET
	status = 'accepted',
	last_updated_on = extract(epoch FROM NOW()), 
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL 
AND belongs_to_household = $1
AND id = $2
`

func (q *SQLQuerier) AcceptHouseholdInvitation(ctx context.Context, householdInvitationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)

	args := []interface{}{householdInvitationID}

	if err := q.performWriteQuery(ctx, q.db, "household invitation accept", acceptHouseholdInvitationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "accepting household invitation")
	}

	logger.Info("household invitation accepted")

	return nil
}

const rejectHouseholdInvitationQuery = `
UPDATE household_invitations SET
	status = 'rejected',
	last_updated_on = extract(epoch FROM NOW()), 
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL 
AND belongs_to_household = $1
AND id = $2
`

func (q *SQLQuerier) RejectHouseholdInvitation(ctx context.Context, householdInvitationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)

	args := []interface{}{householdInvitationID}

	if err := q.performWriteQuery(ctx, q.db, "household invitation reject", rejectHouseholdInvitationQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "rejecting household invitation")
	}

	logger.Info("household invitation rejected")

	return nil
}
