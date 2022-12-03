package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/prixfixeco/backend/internal/authorization"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	// householdsUserMembershipTableName is what the households membership table calls itself.
	householdsUserMembershipTableName = "household_user_memberships"
)

var (
	_ types.HouseholdUserMembershipDataManager = (*Querier)(nil)

	// householdsUserMembershipTableColumns are the columns for the household user memberships table.
	householdsUserMembershipTableColumns = []string{
		"household_user_memberships.id",
		"household_user_memberships.belongs_to_user",
		"household_user_memberships.belongs_to_household",
		"household_user_memberships.household_role",
		"household_user_memberships.default_household",
		"household_user_memberships.created_at",
		"household_user_memberships.last_updated_at",
		"household_user_memberships.archived_at",
	}
)

// scanHouseholdUserMembership takes a database Scanner (i.e. *sql.Row) and scans the result into a householdUserMembership struct.
func (q *Querier) scanHouseholdUserMembership(ctx context.Context, scan database.Scanner) (x *types.HouseholdUserMembership, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.HouseholdUserMembership{}

	targetVars := []any{
		&x.ID,
		&x.BelongsToUser,
		&x.BelongsToHousehold,
		&x.HouseholdRole,
		&x.DefaultHousehold,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, observability.PrepareError(err, span, "scanning household user memberships")
	}

	return x, nil
}

// scanHouseholdUserMemberships takes some database rows and turns them into a slice of memberships.
func (q *Querier) scanHouseholdUserMemberships(ctx context.Context, rows database.ResultIterator) (defaultHousehold string, householdRolesMap map[string]string, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	householdRolesMap = map[string]string{}

	for rows.Next() {
		x, scanErr := q.scanHouseholdUserMembership(ctx, rows)
		if scanErr != nil {
			return "", nil, scanErr
		}

		if x.DefaultHousehold && defaultHousehold == "" {
			defaultHousehold = x.BelongsToHousehold
		}

		householdRolesMap[x.BelongsToHousehold] = x.HouseholdRole
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return "", nil, observability.PrepareError(err, span, "handling rows")
	}

	return defaultHousehold, householdRolesMap, nil
}

//go:embed queries/household_user_memberships/get_for_user.sql
var getHouseholdMembershipsForUserQuery string

// BuildSessionContextDataForUser queries the database for the memberships of a user and constructs a SessionContextData struct from the results.
func (q *Querier) BuildSessionContextDataForUser(ctx context.Context, userID string) (*types.SessionContextData, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	user, err := q.GetUser(ctx, userID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching user from database")
	}

	getHouseholdMembershipsArgs := []any{userID}

	membershipRows, membershipReadErr := q.getRows(ctx, q.db, "household memberships for user", getHouseholdMembershipsForUserQuery, getHouseholdMembershipsArgs)
	if membershipReadErr != nil {
		return nil, observability.PrepareError(membershipReadErr, span, "fetching user's memberships from database")
	}

	defaultHouseholdID, householdRolesMap, membershipsScanErr := q.scanHouseholdUserMemberships(ctx, membershipRows)
	if membershipsScanErr != nil {
		return nil, observability.PrepareError(membershipsScanErr, span, "scanning user's memberships from database")
	}

	actualHouseholdRolesMap := map[string]authorization.HouseholdRolePermissionsChecker{}
	for householdID, roles := range householdRolesMap {
		actualHouseholdRolesMap[householdID] = authorization.NewHouseholdRolePermissionChecker(roles)
	}

	sessionCtxData := &types.SessionContextData{
		Requester: types.RequesterInfo{
			Username:                 user.Username,
			EmailAddress:             user.EmailAddress,
			UserID:                   user.ID,
			AccountStatus:            user.AccountStatus,
			AccountStatusExplanation: user.AccountStatusExplanation,
			ServicePermissions:       authorization.NewServiceRolePermissionChecker(user.ServiceRole),
		},
		HouseholdPermissions: actualHouseholdRolesMap,
		ActiveHouseholdID:    defaultHouseholdID,
	}

	logger.Debug("fetched session context data for user")

	return sessionCtxData, nil
}

//go:embed queries/household_user_memberships/get_default_household_id_for_user.sql
var getDefaultHouseholdIDForUserQuery string

// GetDefaultHouseholdIDForUser retrieves the default household ID for a user.
func (q *Querier) GetDefaultHouseholdIDForUser(ctx context.Context, userID string) (string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return "", ErrInvalidIDProvided
	}

	args := []any{userID, true}

	var id string
	if err := q.getOneRow(ctx, q.db, "default household ID query", getDefaultHouseholdIDForUserQuery, args).Scan(&id); err != nil {
		return "", observability.PrepareError(err, span, "executing default household ID query")
	}

	return id, nil
}

//go:embed queries/household_user_memberships/mark_as_user_default.sql
var markHouseholdAsUserDefaultQuery string

// markHouseholdAsUserDefault does a thing.
func (q *Querier) markHouseholdAsUserDefault(ctx context.Context, querier database.SQLQueryExecutor, userID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" || householdID == "" {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValues(map[string]any{
		keys.UserIDKey:      userID,
		keys.HouseholdIDKey: householdID,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []any{
		userID,
		householdID,
		userID,
	}

	// create the household.
	if err := q.performWriteQuery(ctx, querier, "user default household assignment", markHouseholdAsUserDefaultQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "assigning user default household")
	}

	logger.Debug("household marked as default")

	return nil
}

// MarkHouseholdAsUserDefault does a thing.
func (q *Querier) MarkHouseholdAsUserDefault(ctx context.Context, userID, householdID string) error {
	return q.markHouseholdAsUserDefault(ctx, q.db, userID, householdID)
}

//go:embed queries/household_user_memberships/user_is_member.sql
var userIsMemberOfHouseholdQuery string

// UserIsMemberOfHousehold does a thing.
func (q *Querier) UserIsMemberOfHousehold(ctx context.Context, userID, householdID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" || householdID == "" {
		return false, ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []any{
		householdID,
		userID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, userIsMemberOfHouseholdQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, span, "performing user membership check query")
	}

	return result, nil
}

//go:embed queries/household_user_memberships/modify_user_permissions.sql
var modifyUserPermissionsQuery string

// ModifyUserPermissions does a thing.
func (q *Querier) ModifyUserPermissions(ctx context.Context, householdID, userID string, input *types.ModifyUserPermissionsInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" || userID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]any{
		keys.HouseholdIDKey: householdID,
		keys.UserIDKey:      userID,
		"new_roles":         input.NewRole,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []any{
		input.NewRole,
		householdID,
		userID,
	}

	// modify the membership.
	if err := q.performWriteQuery(ctx, q.db, "user household permissions modification", modifyUserPermissionsQuery, args); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareAndLogError(err, logger, span, "modifying user household permissions")
	}

	logger.Info("user permissions modified")

	return nil
}

//go:embed queries/household_user_memberships/transfer_ownership.sql
var transferHouseholdOwnershipQuery string

//go:embed queries/household_user_memberships/transfer_membership.sql
var transferHouseholdMembershipQuery string

// TransferHouseholdOwnership does a thing.
func (q *Querier) TransferHouseholdOwnership(ctx context.Context, householdID string, input *types.HouseholdOwnershipTransferInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]any{
		keys.HouseholdIDKey: householdID,
		"current_owner":     input.CurrentOwner,
		"new_owner":         input.NewOwner,
	})

	tracing.AttachUserIDToSpan(span, input.NewOwner)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// begin household transfer transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	transferHouseholdOwnershipArgs := []any{
		input.NewOwner,
		input.CurrentOwner,
		householdID,
	}

	// create the membership.
	if err = q.performWriteQuery(ctx, tx, "user ownership transfer", transferHouseholdOwnershipQuery, transferHouseholdOwnershipArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "transferring household to new owner")
	}

	transferHouseholdMembershipArgs := []any{
		input.NewOwner,
		householdID,
		input.CurrentOwner,
	}

	// create the membership.
	if err = q.performWriteQuery(ctx, tx, "user memberships transfer", transferHouseholdMembershipQuery, transferHouseholdMembershipArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "transferring household memberships")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("household transferred to new owner")

	return nil
}

//go:embed queries/household_user_memberships/add_user_to_household.sql
var addUserToHouseholdQuery string

// addUserToHousehold does a thing.
func (q *Querier) addUserToHousehold(ctx context.Context, querier database.SQLQueryExecutor, input *types.HouseholdUserMembershipDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]any{
		keys.UserIDKey:      input.UserID,
		keys.HouseholdIDKey: input.HouseholdID,
	})

	tracing.AttachUserIDToSpan(span, input.UserID)
	tracing.AttachHouseholdIDToSpan(span, input.HouseholdID)

	addUserToHouseholdArgs := []any{
		input.ID,
		input.UserID,
		input.HouseholdID,
		input.HouseholdRole,
	}

	// create the membership.
	if err := q.performWriteQuery(ctx, querier, "user household membership creation", addUserToHouseholdQuery, addUserToHouseholdArgs); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "performing user household membership creation query")
	}

	logger.Info("user added to household")

	return nil
}

//go:embed queries/household_user_memberships/remove_user_from_household.sql
var removeUserFromHouseholdQuery string

// RemoveUserFromHousehold removes a user's membership to a household.
func (q *Querier) RemoveUserFromHousehold(ctx context.Context, userID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachHouseholdIDToSpan(span, householdID)

	logger := q.logger.WithValues(map[string]any{
		keys.UserIDKey:      userID,
		keys.HouseholdIDKey: householdID,
	})

	logger.Info("creating")

	tx, createTransactionErr := q.db.BeginTx(ctx, nil)
	if createTransactionErr != nil {
		return observability.PrepareError(createTransactionErr, span, "beginning transaction")
	}

	logger.Info("created transaction")

	args := []any{
		householdID,
		userID,
	}

	// remove the membership.
	if err := q.performWriteQuery(ctx, tx, "user membership removal", removeUserFromHouseholdQuery, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "removing user from household")
	}

	logger.Info("user removed from household")

	remainingHouseholds, fetchRemainingHouseholdsErr := q.getHouseholds(ctx, tx, userID, false, nil)
	if fetchRemainingHouseholdsErr != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(fetchRemainingHouseholdsErr, span, "fetching remaining households")
	}

	logger = logger.WithValue("count", len(remainingHouseholds.Data))
	logger.Info("remaining households fetched")

	if len(remainingHouseholds.Data) == 0 {
		if err := q.createHouseholdForUser(ctx, tx, false, "", userID); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "creating household for new user")
		}
		return nil
	}

	household := remainingHouseholds.Data[0]

	logger = logger.WithValue(keys.HouseholdIDKey, household.ID)
	logger.Info("about to mark household as default")

	if err := q.markHouseholdAsUserDefault(ctx, tx, userID, household.ID); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "marking household as default")
	}

	logger.Info("marked household as default, committing transaction")

	if err := tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return nil
}
