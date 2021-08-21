package querier

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

var (
	_ types.HouseholdDataManager = (*SQLQuerier)(nil)
)

// scanHousehold takes a database Scanner (i.e. *sql.Row) and scans the result into an Household struct.
func (q *SQLQuerier) scanHousehold(ctx context.Context, scan database.Scanner, includeCounts bool) (household *types.Household, membership *types.HouseholdUserMembership, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	household = &types.Household{Members: []*types.HouseholdUserMembership{}}
	membership = &types.HouseholdUserMembership{}

	var (
		rawRoles string
	)

	targetVars := []interface{}{
		&household.ID,
		&household.ExternalID,
		&household.Name,
		&household.BillingStatus,
		&household.ContactEmail,
		&household.ContactPhone,
		&household.PaymentProcessorCustomerID,
		&household.SubscriptionPlanID,
		&household.CreatedOn,
		&household.LastUpdatedOn,
		&household.ArchivedOn,
		&household.BelongsToUser,
		&membership.ID,
		&membership.BelongsToUser,
		&membership.BelongsToHousehold,
		&rawRoles,
		&membership.DefaultHousehold,
		&membership.CreatedOn,
		&membership.LastUpdatedOn,
		&membership.ArchivedOn,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, nil, 0, 0, observability.PrepareError(err, logger, span, "fetching memberships from database")
	}

	membership.HouseholdRoles = strings.Split(rawRoles, householdMemberRolesSeparator)

	return household, membership, filteredCount, totalCount, nil
}

// scanHouseholds takes some database rows and turns them into a slice of households.
func (q *SQLQuerier) scanHouseholds(ctx context.Context, rows database.ResultIterator, includeCounts bool) (households []*types.Household, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)
	households = []*types.Household{}

	var currentHousehold *types.Household
	for rows.Next() {
		household, membership, fc, tc, scanErr := q.scanHousehold(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if currentHousehold == nil {
			currentHousehold = household
		}

		if currentHousehold.ID != household.ID {
			households = append(households, currentHousehold)
			currentHousehold = household
		}

		currentHousehold.Members = append(currentHousehold.Members, membership)

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
		households = append(households, currentHousehold)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return households, filteredCount, totalCount, nil
}

// GetHousehold fetches an household from the database.
func (q *SQLQuerier) GetHousehold(ctx context.Context, householdID, userID uint64) (*types.Household, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == 0 || userID == 0 {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachUserIDToSpan(span, userID)

	logger := q.logger.WithValues(map[string]interface{}{
		keys.HouseholdIDKey: householdID,
		keys.UserIDKey:      userID,
	})

	query, args := q.sqlQueryBuilder.BuildGetHouseholdQuery(ctx, householdID, userID)
	rows, err := q.performReadQuery(ctx, q.db, "household", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing households list retrieval query")
	}

	households, _, _, err := q.scanHouseholds(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	var household *types.Household
	if len(households) > 0 {
		household = households[0]
	}

	if household == nil {
		return nil, sql.ErrNoRows
	}

	return household, nil
}

// GetAllHouseholdsCount fetches the count of households from the database that meet a particular filter.
func (q *SQLQuerier) GetAllHouseholdsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger

	count, err := q.performCountQuery(ctx, q.db, q.sqlQueryBuilder.BuildGetAllHouseholdsCountQuery(ctx), "fetching count of all households")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of households")
	}

	return count, nil
}

// GetAllHouseholds fetches a list of all households in the database.
func (q *SQLQuerier) GetAllHouseholds(ctx context.Context, results chan []*types.Household, batchSize uint16) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if results == nil {
		return ErrNilInputProvided
	}

	if batchSize == 0 {
		batchSize = defaultBatchSize
	}

	logger := q.logger.WithValue("batch_size", batchSize)

	count, err := q.GetAllHouseholdsCount(ctx)
	if err != nil {
		return observability.PrepareError(err, logger, span, "fetching count of households")
	}

	for beginID := uint64(1); beginID <= count; beginID += uint64(batchSize) {
		endID := beginID + uint64(batchSize)
		go func(begin, end uint64) {
			query, args := q.sqlQueryBuilder.BuildGetBatchOfHouseholdsQuery(ctx, begin, end)
			logger = logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, queryErr := q.db.Query(query, args...)
			if errors.Is(queryErr, sql.ErrNoRows) {
				return
			} else if queryErr != nil {
				observability.AcknowledgeError(queryErr, logger, span, "querying for database rows")
				return
			}

			households, _, _, scanErr := q.scanHouseholds(ctx, rows, false)
			if scanErr != nil {
				observability.AcknowledgeError(scanErr, logger, span, "scanning database rows")
				return
			}

			results <- households
		}(beginID, endID)
	}

	return nil
}

// GetHouseholds fetches a list of households from the database that meet a particular filter.
func (q *SQLQuerier) GetHouseholds(ctx context.Context, userID uint64, filter *types.QueryFilter) (x *types.HouseholdList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == 0 {
		return nil, ErrInvalidIDProvided
	}

	logger := filter.AttachToLogger(q.logger).WithValue(keys.UserIDKey, userID)
	tracing.AttachQueryFilterToSpan(span, filter)
	tracing.AttachUserIDToSpan(span, userID)

	x = &types.HouseholdList{}
	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetHouseholdsQuery(ctx, userID, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "households", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing households list retrieval query")
	}

	if x.Households, x.FilteredCount, x.TotalCount, err = q.scanHouseholds(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning households from database")
	}

	return x, nil
}

// GetHouseholdsForAdmin fetches a list of households from the database that meet a particular filter for all users.
func (q *SQLQuerier) GetHouseholdsForAdmin(ctx context.Context, filter *types.QueryFilter) (x *types.HouseholdList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(q.logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.HouseholdList{}
	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.sqlQueryBuilder.BuildGetHouseholdsQuery(ctx, 0, true, filter)

	rows, err := q.performReadQuery(ctx, q.db, "households for admin", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for households")
	}

	if x.Households, x.FilteredCount, x.TotalCount, err = q.scanHouseholds(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning households")
	}

	return x, nil
}

// CreateHousehold creates an household in the database.
func (q *SQLQuerier) CreateHousehold(ctx context.Context, input *types.HouseholdCreationInput, createdByUser uint64) (*types.Household, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if createdByUser == 0 {
		return nil, ErrInvalidIDProvided
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RequesterIDKey, createdByUser).WithValue(keys.UserIDKey, input.BelongsToUser)
	tracing.AttachRequestingUserIDToSpan(span, createdByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	householdCreationQuery, householdCreationArgs := q.sqlQueryBuilder.BuildHouseholdCreationQuery(ctx, input)

	// create the household.
	id, err := q.performWriteQuery(ctx, tx, false, "household creation", householdCreationQuery, householdCreationArgs)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating household")
	}

	logger = logger.WithValue(keys.HouseholdIDKey, id)

	household := &types.Household{
		ID:            id,
		Name:          input.Name,
		BelongsToUser: input.BelongsToUser,
		BillingStatus: types.UnpaidHouseholdBillingStatus,
		ContactEmail:  input.ContactEmail,
		ContactPhone:  input.ContactPhone,
		CreatedOn:     q.currentTime(),
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildHouseholdCreationEventEntry(household, createdByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing household creation audit log event entry")
	}

	addInput := &types.AddUserToHouseholdInput{
		UserID:         input.BelongsToUser,
		HouseholdID:    household.ID,
		Reason:         "household creation",
		HouseholdRoles: []string{authorization.HouseholdAdminRole.String()},
	}

	addUserToHouseholdQuery, addUserToHouseholdArgs := q.sqlQueryBuilder.BuildAddUserToHouseholdQuery(ctx, addInput)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "household user membership creation", addUserToHouseholdQuery, addUserToHouseholdArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "creating household membership")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildUserAddedToHouseholdEventEntry(createdByUser, addInput)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing household membership creation audit log event entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachHouseholdIDToSpan(span, household.ID)
	logger.Info("household created")

	return household, nil
}

// UpdateHousehold updates a particular household. Note that UpdateHousehold expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateHousehold(ctx context.Context, updated *types.Household, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.HouseholdIDKey, updated.ID)
	tracing.AttachHouseholdIDToSpan(span, updated.ID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)
	tracing.AttachChangeSummarySpan(span, "household", changes)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildUpdateHouseholdQuery(ctx, updated)
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "household update", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "updating household")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildHouseholdUpdateEventEntry(updated.BelongsToUser, updated.ID, changedByUser, changes)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing household update audit log event entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("household updated")

	return nil
}

// ArchiveHousehold archives an household from the database by its ID.
func (q *SQLQuerier) ArchiveHousehold(ctx context.Context, householdID, userID, archivedByUser uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == 0 || userID == 0 || archivedByUser == 0 {
		return ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	logger := q.logger.WithValues(map[string]interface{}{
		keys.RequesterIDKey: archivedByUser,
		keys.HouseholdIDKey: householdID,
		keys.UserIDKey:      userID,
	})

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildArchiveHouseholdQuery(ctx, householdID, userID)

	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "household archive", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "archiving household")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildHouseholdArchiveEventEntry(userID, householdID, archivedByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing household archive audit log event entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("household archived")

	return nil
}

// GetAuditLogEntriesForHousehold fetches a list of audit log entries from the database that relate to a given household.
func (q *SQLQuerier) GetAuditLogEntriesForHousehold(ctx context.Context, householdID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == 0 {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachHouseholdIDToSpan(span, householdID)
	logger := q.logger.WithValue(keys.HouseholdIDKey, householdID)

	query, args := q.sqlQueryBuilder.BuildGetAuditLogEntriesForHouseholdQuery(ctx, householdID)

	rows, err := q.performReadQuery(ctx, q.db, "audit log entries for household", query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "querying database for audit log entries")
	}

	auditLogEntries, _, err := q.scanAuditLogEntries(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning audit log entries")
	}

	return auditLogEntries, nil
}
