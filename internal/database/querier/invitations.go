package querier

import (
	"context"
	"database/sql"
	"errors"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	database "gitlab.com/prixfixe/prixfixe/internal/database"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var (
	_ types.InvitationDataManager = (*SQLQuerier)(nil)
)

// scanInvitation takes a database Scanner (i.e. *sql.Row) and scans the result into an invitation struct.
func (q *SQLQuerier) scanInvitation(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.Invitation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.Invitation{}

	targetVars := []interface{}{
		&x.ID,
		&x.ExternalID,
		&x.Code,
		&x.Consumed,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToAccount,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanInvitations takes some database rows and turns them into a slice of invitations.
func (q *SQLQuerier) scanInvitations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (invitations []*types.Invitation, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	invitations = []*types.Invitation{}

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanInvitation(ctx, rows, includeCounts)
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

		invitations = append(invitations, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return invitations, filteredCount, totalCount, nil
}

// InvitationExists fetches whether an invitation exists from the database.
func (q *SQLQuerier) InvitationExists(ctx context.Context, invitationID uint64) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if invitationID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	query, args := q.sqlQueryBuilder.BuildInvitationExistsQuery(ctx, invitationID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing invitation existence check")
	}

	return result, nil
}

// GetInvitation fetches an invitation from the database.
func (q *SQLQuerier) GetInvitation(ctx context.Context, invitationID uint64) (*types.Invitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if invitationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	query, args := q.sqlQueryBuilder.BuildGetInvitationQuery(ctx, invitationID)
	row := q.getOneRow(ctx, q.db, "invitation", query, args...)

	invitation, _, _, err := q.scanInvitation(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning invitation")
	}

	return invitation, nil
}

// GetAllInvitationsCount fetches the count of invitations from the database that meet a particular filter.
func (q *SQLQuerier) GetAllInvitationsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllInvitationsCountQuery(ctx), "fetching count of invitations")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of invitations")
	}

	return count, nil
}

// GetAllInvitations fetches a list of all invitations in the database.
func (q *SQLQuerier) GetAllInvitations(ctx context.Context, results chan []*types.Invitation, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllInvitationsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of invitations")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfInvitationsQuery(ctx, begin, end)
			logger = logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, queryErr := q.db.Query(query, args...)
			if errors.Is(queryErr, sql.ErrNoRows) {
				return
			} else if queryErr != nil {
				logger.Error(queryErr, "querying for database rows")
				return
			}

			invitations, _, _, scanErr := q.scanInvitations(ctx, rows, false)
			if scanErr != nil {
				logger.Error(scanErr, "scanning database rows")
				return
			}

			results <- invitations
		}(beginID, endID)
	}

	return nil
}

// GetInvitations fetches a list of invitations from the database that meet a particular filter.
func (q *SQLQuerier) GetInvitations(ctx context.Context, filter *types.QueryFilter) (x *types.InvitationList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	x = &types.InvitationList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetInvitationsQuery(ctx, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "invitations", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing invitations list retrieval query")
	}

	if x.Invitations, x.FilteredCount, x.TotalCount, err = q.scanInvitations(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning invitations")
	}

	return x, nil
}

// GetInvitationsWithIDs fetches invitations from the database within a given set of IDs.
func (q *SQLQuerier) GetInvitationsWithIDs(ctx context.Context, accountID uint64, limit uint8, ids []uint64) ([]*types.Invitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if accountID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachAccountIDToSpan(span, accountID)

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.sqlQueryBuilder.BuildGetInvitationsWithIDsQuery(ctx, accountID, limit, ids, false)

	rows, err := q.performReadQuery(ctx, q.db, "invitations with IDs", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching invitations from database")
	}

	invitations, _, _, err := q.scanInvitations(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning invitations")
	}

	return invitations, nil
}

// CreateInvitation creates an invitation in the database.
func (q *SQLQuerier) CreateInvitation(ctx context.Context, input *types.InvitationCreationInput, createdByUser uint64) (*types.Invitation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if createdByUser == 0 {
		return nil, ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.RequesterIDKey, createdByUser)
	tracing.AttachRequestingUserIDToSpan(span, createdByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildCreateInvitationQuery(ctx, input)

	// create the invitation.
	id, err := q.performWriteQuery(ctx, tx, false, "invitation creation", query, args)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating invitation")
	}

	x := &types.Invitation{
		ID:               id,
		Code:             input.Code,
		Consumed:         input.Consumed,
		BelongsToAccount: input.BelongsToAccount,
		CreatedOn:        q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildInvitationCreationEventEntry(x, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing invitation creation audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachInvitationIDToSpan(span, x.ID)
	logger.Info("invitation created")

	return x, nil
}

// UpdateInvitation updates a particular invitation. Note that UpdateInvitation expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateInvitation(ctx context.Context, updated *types.Invitation, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.InvitationIDKey, updated.ID)
	tracing.AttachInvitationIDToSpan(span, updated.ID)
	tracing.AttachAccountIDToSpan(span, updated.BelongsToAccount)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateInvitationQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "invitation update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating invitation")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildInvitationUpdateEventEntry(changedByUser, updated.ID, updated.BelongsToAccount, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing invitation update audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("invitation updated")

	return nil
}

// ArchiveInvitation archives an invitation from the database by its ID.
func (q *SQLQuerier) ArchiveInvitation(ctx context.Context, invitationID, accountID, archivedBy uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if invitationID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	if accountID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachAccountIDToSpan(span, accountID)

	if archivedBy == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RequesterIDKey, archivedBy)
	tracing.AttachUserIDToSpan(span, archivedBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveInvitationQuery(ctx, invitationID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "invitation archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating invitation")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildInvitationArchiveEventEntry(archivedBy, accountID, invitationID)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing invitation archive audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("invitation archived")

	return nil
}

// GetAuditLogEntriesForInvitation fetches a list of audit log entries from the database that relate to a given invitation.
func (q *SQLQuerier) GetAuditLogEntriesForInvitation(ctx context.Context, invitationID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	if invitationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForInvitationQuery(ctx, invitationID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for invitation", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
