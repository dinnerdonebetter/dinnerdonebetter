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
	_ types.HouseholdUserMembershipDataManager = (*SQLQuerier)(nil)
)

const (
	householdMemberRolesSeparator = ","
)

// scanHouseholdUserMembership takes a database Scanner (i.e. *sql.Row) and scans the result into an HouseholdUserMembership struct.
func (q *SQLQuerier) scanHouseholdUserMembership(ctx context.Context, scan database.Scanner) (x *types.HouseholdUserMembership, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.HouseholdUserMembership{}

	var (
		rawHouseholdRoles string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.BelongsToUser,
		&x.BelongsToHousehold,
		&rawHouseholdRoles,
		&x.DefaultHousehold,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, observability.PrepareError(err, q.logger, span, "scanning household user memberships")
	}

	if roles := strings.Split(rawHouseholdRoles, householdMemberRolesSeparator); len(roles) > 0 {
		x.HouseholdRoles = roles
	} else {
		x.HouseholdRoles = []string{}
	}

	return x, nil
}

// scanHouseholdUserMemberships takes some database rows and turns them into a slice of memberships.
func (q *SQLQuerier) scanHouseholdUserMemberships(ctx context.Context, rows database.ResultIterator) (defaultHousehold uint64, householdRolesMap map[uint64][]string, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	householdRolesMap = map[uint64][]string{}
	logger := q.logger

	for rows.Next() {
		x, scanErr := q.scanHouseholdUserMembership(ctx, rows)
		if scanErr != nil {
			return 0, nil, scanErr
		}

		if x.DefaultHousehold && defaultHousehold == 0 {
			defaultHousehold = x.BelongsToHousehold
		}

		householdRolesMap[x.BelongsToHousehold] = x.HouseholdRoles
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return 0, nil, observability.PrepareError(err, logger, span, "handling rows")
	}

	return defaultHousehold, householdRolesMap, nil
}

// BuildSessionContextDataForUser does .
func (q *SQLQuerier) BuildSessionContextDataForUser(ctx context.Context, userID uint64) (*types.SessionContextData, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == 0 {
		return nil, ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	user, err := q.GetUser(ctx, userID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching user from database")
	}

	getHouseholdMembershipsQuery, getHouseholdMembershipsArgs := q.sqlQueryBuilder.BuildGetHouseholdMembershipsForUserQuery(ctx, userID)

	membershipRows, err := q.performReadQuery(ctx, q.db, "household memberships for user", getHouseholdMembershipsQuery, getHouseholdMembershipsArgs...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching user's memberships from database")
	}

	defaultHouseholdID, householdRolesMap, err := q.scanHouseholdUserMemberships(ctx, membershipRows)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning user's memberships from database")
	}

	actualHouseholdRolesMap := map[uint64]authorization.HouseholdRolePermissionsChecker{}
	for householdID, roles := range householdRolesMap {
		actualHouseholdRolesMap[householdID] = authorization.NewHouseholdRolePermissionChecker(roles...)
	}

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
			UserID:                user.ID,
			Reputation:            user.ServiceHouseholdStatus,
			ReputationExplanation: user.ReputationExplanation,
			ServicePermissions:    authorization.NewServiceRolePermissionChecker(user.ServiceRoles...),
		},
		HouseholdPermissions: actualHouseholdRolesMap,
		ActiveHouseholdID:    defaultHouseholdID,
	}

	return sessionCtxData, nil
}

// GetDefaultHouseholdIDForUser does .
func (q *SQLQuerier) GetDefaultHouseholdIDForUser(ctx context.Context, userID uint64) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == 0 {
		return 0, ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	query, args := q.sqlQueryBuilder.BuildGetDefaultHouseholdIDForUserQuery(ctx, userID)

	var id uint64
	if err := q.getOneRow(ctx, q.db, "default household ID query", query, args...).Scan(&id); err != nil {
		return 0, observability.PrepareError(err, logger, span, "executing id query")
	}

	return id, nil
}

// MarkHouseholdAsUserDefault does a thing.
func (q *SQLQuerier) MarkHouseholdAsUserDefault(ctx context.Context, userID, householdID, changedByUser uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == 0 || householdID == 0 || changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.UserIDKey:      userID,
		keys.HouseholdIDKey: householdID,
		keys.RequesterIDKey: changedByUser,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildMarkHouseholdAsUserDefaultQuery(ctx, userID, householdID)

	// create the household.
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "user default household assignment", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "assigning user default household")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildUserMarkedHouseholdAsDefaultEventEntry(userID, householdID, changedByUser)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "household not found for user")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("household marked as default")

	return nil
}

// UserIsMemberOfHousehold does a thing.
func (q *SQLQuerier) UserIsMemberOfHousehold(ctx context.Context, userID, householdID uint64) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == 0 || householdID == 0 {
		return false, ErrInvalidIDProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.UserIDKey:      userID,
		keys.HouseholdIDKey: householdID,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	query, args := q.sqlQueryBuilder.BuildUserIsMemberOfHouseholdQuery(ctx, userID, householdID)

	result, err := q.performBooleanQuery(ctx, q.db, query, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing user membership check query")
	}

	return result, nil
}

// ModifyUserPermissions does a thing.
func (q *SQLQuerier) ModifyUserPermissions(ctx context.Context, userID, householdID, changedByUser uint64, input *types.ModifyUserPermissionsInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == 0 || userID == 0 || changedByUser == 0 {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.HouseholdIDKey: householdID,
		keys.UserIDKey:      userID,
		keys.RequesterIDKey: changedByUser,
		"new_roles":         input.NewRoles,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachRequestingUserIDToSpan(span, changedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildModifyUserPermissionsQuery(ctx, userID, householdID, input.NewRoles)

	// modify the membership.
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "user household permissions modification", query, args); err != nil && !errors.Is(err, sql.ErrNoRows) {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "modifying user household permissions")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildModifyUserPermissionsEventEntry(userID, householdID, changedByUser, input.NewRoles, input.Reason)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing user household membership permission modification audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("user permissions modified")

	return nil
}

// TransferHouseholdOwnership does a thing.
func (q *SQLQuerier) TransferHouseholdOwnership(ctx context.Context, householdID, transferredBy uint64, input *types.HouseholdOwnershipTransferInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == 0 || transferredBy == 0 {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.HouseholdIDKey: householdID,
		keys.RequesterIDKey: transferredBy,
		"current_owner":     input.CurrentOwner,
		"new_owner":         input.NewOwner,
	})

	tracing.AttachUserIDToSpan(span, input.NewOwner)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachRequestingUserIDToSpan(span, transferredBy)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	transferHouseholdOwnershipQuery, transferHouseholdOwnershipArgs := q.sqlQueryBuilder.BuildTransferHouseholdOwnershipQuery(ctx, input.CurrentOwner, input.NewOwner, householdID)

	// create the membership.
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "user ownership transfer", transferHouseholdOwnershipQuery, transferHouseholdOwnershipArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "transferring household to new owner")
	}

	transferHouseholdMembershipQuery, transferHouseholdMembershipArgs := q.sqlQueryBuilder.BuildTransferHouseholdMembershipsQuery(ctx, input.CurrentOwner, input.NewOwner, householdID)

	// create the membership.
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "user memberships transfer", transferHouseholdMembershipQuery, transferHouseholdMembershipArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "transferring household memberships")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildTransferHouseholdOwnershipEventEntry(householdID, transferredBy, input)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing household ownership transfer audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("household transferred to new owner")

	return nil
}

// AddUserToHousehold does a thing.
func (q *SQLQuerier) AddUserToHousehold(ctx context.Context, input *types.AddUserToHouseholdInput, addedByUser uint64) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if addedByUser == 0 {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.RequesterIDKey: addedByUser,
		keys.UserIDKey:      input.UserID,
		keys.HouseholdIDKey: input.HouseholdID,
	})

	tracing.AttachUserIDToSpan(span, input.UserID)
	tracing.AttachHouseholdIDToSpan(span, input.HouseholdID)
	tracing.AttachRequestingUserIDToSpan(span, addedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildAddUserToHouseholdQuery(ctx, input)

	// create the membership.
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "user household membership creation", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "creating user household membership")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildUserAddedToHouseholdEventEntry(addedByUser, input)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing user added to household audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("user added to household")

	return nil
}

// RemoveUserFromHousehold removes a user's membership to an household.
func (q *SQLQuerier) RemoveUserFromHousehold(ctx context.Context, userID, householdID, removedByUser uint64, reason string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == 0 || householdID == 0 || removedByUser == 0 {
		return ErrInvalidIDProvided
	}

	if reason == "" {
		return ErrEmptyInputProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.UserIDKey:      userID,
		keys.HouseholdIDKey: householdID,
		keys.ReasonKey:      reason,
		keys.RequesterIDKey: removedByUser,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachRequestingUserIDToSpan(span, removedByUser)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	query, args := q.sqlQueryBuilder.BuildRemoveUserFromHouseholdQuery(ctx, userID, householdID)

	// create the membership.
	if err = q.performWriteQueryIgnoringReturn(ctx, tx, "user membership removal", query, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "removing user from household")
	}

	if err = q.createAuditLogEntryInTransaction(ctx, tx, audit.BuildUserRemovedFromHouseholdEventEntry(userID, householdID, removedByUser, reason)); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "writing remove user from household audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("user removed from household")

	return nil
}
