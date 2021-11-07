package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// householdsUserMembershipTableName is what the households membership table calls itself.
	householdsUserMembershipTableName = "household_user_memberships"
)

var (
	_ types.HouseholdUserMembershipDataManager = (*SQLQuerier)(nil)

	// householdsUserMembershipTableColumns are the columns for the household user memberships table.
	householdsUserMembershipTableColumns = []string{
		"household_user_memberships.id",
		"household_user_memberships.belongs_to_user",
		"household_user_memberships.belongs_to_household",
		"household_user_memberships.household_roles",
		"household_user_memberships.default_household",
		"household_user_memberships.created_on",
		"household_user_memberships.last_updated_on",
		"household_user_memberships.archived_on",
	}
)

const (
	householdMemberRolesSeparator = commaSeparator
)

// scanHouseholdUserMembership takes a database Scanner (i.e. *sql.Row) and scans the result into a householdUserMembership struct.
func (q *SQLQuerier) scanHouseholdUserMembership(ctx context.Context, scan database.Scanner) (x *types.HouseholdUserMembership, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.HouseholdUserMembership{
		HouseholdRoles: []string{},
	}

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
	}

	return x, nil
}

// scanHouseholdUserMemberships takes some database rows and turns them into a slice of memberships.
func (q *SQLQuerier) scanHouseholdUserMemberships(ctx context.Context, rows database.ResultIterator) (defaultHousehold string, householdRolesMap map[string][]string, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	householdRolesMap = map[string][]string{}
	logger := q.logger

	for rows.Next() {
		x, scanErr := q.scanHouseholdUserMembership(ctx, rows)
		if scanErr != nil {
			return "", nil, scanErr
		}

		if x.DefaultHousehold && defaultHousehold == "" {
			defaultHousehold = x.BelongsToHousehold
		}

		householdRolesMap[x.BelongsToHousehold] = x.HouseholdRoles
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return "", nil, observability.PrepareError(err, logger, span, "handling rows")
	}

	return defaultHousehold, householdRolesMap, nil
}

const getHouseholdMembershipsForUserQuery = `
	SELECT
		household_user_memberships.id,
		household_user_memberships.belongs_to_user,
		household_user_memberships.belongs_to_household,
		household_user_memberships.household_roles,
		household_user_memberships.default_household,
		household_user_memberships.created_on,
		household_user_memberships.last_updated_on,
		household_user_memberships.archived_on
	FROM household_user_memberships
	JOIN households ON households.id = household_user_memberships.belongs_to_household
	WHERE household_user_memberships.archived_on IS NULL
	AND household_user_memberships.belongs_to_user = $1
`

// BuildSessionContextDataForUser queries the database for the memberships of a user and constructs a SessionContextData struct from the results.
func (q *SQLQuerier) BuildSessionContextDataForUser(ctx context.Context, userID string) (*types.SessionContextData, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	user, err := q.GetUser(ctx, userID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching user from database")
	}

	getHouseholdMembershipsArgs := []interface{}{userID}

	membershipRows, err := q.performReadQuery(ctx, q.db, "household memberships for user", getHouseholdMembershipsForUserQuery, getHouseholdMembershipsArgs)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching user's memberships from database")
	}

	defaultHouseholdID, householdRolesMap, err := q.scanHouseholdUserMemberships(ctx, membershipRows)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning user's memberships from database")
	}

	actualHouseholdRolesMap := map[string]authorization.HouseholdRolePermissionsChecker{}
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

const getDefaultHouseholdIDForUserQuery = `
	SELECT households.id
	FROM households 
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id 
	WHERE household_user_memberships.belongs_to_user = $1
	AND household_user_memberships.default_household = $2
`

// GetDefaultHouseholdIDForUser retrieves the default household ID for a user.
func (q *SQLQuerier) GetDefaultHouseholdIDForUser(ctx context.Context, userID string) (string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return "", ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	args := []interface{}{userID, true}

	var id string
	if err := q.getOneRow(ctx, q.db, "default household ID query", getDefaultHouseholdIDForUserQuery, args).Scan(&id); err != nil {
		return "", observability.PrepareError(err, logger, span, "executing default household ID query")
	}

	return id, nil
}

const markHouseholdAsUserDefaultQuery = `
	UPDATE household_user_memberships
	SET default_household = (belongs_to_user = $1 AND belongs_to_household = $2)
	WHERE archived_on IS NULL
	AND belongs_to_user = $3
`

// MarkHouseholdAsUserDefault does a thing.
func (q *SQLQuerier) MarkHouseholdAsUserDefault(ctx context.Context, userID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" || householdID == "" {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.UserIDKey:      userID,
		keys.HouseholdIDKey: householdID,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		userID,
		householdID,
		userID,
	}

	// create the household.
	if err := q.performWriteQuery(ctx, q.db, "user default household assignment", markHouseholdAsUserDefaultQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "assigning user default household")
	}

	logger.Info("household marked as default")

	return nil
}

const userIsMemberOfHouseholdQuery = `
	SELECT EXISTS ( 
		SELECT household_user_memberships.id 
		FROM household_user_memberships 
		WHERE household_user_memberships.archived_on IS NULL 
		AND household_user_memberships.belongs_to_household = $1 
		AND household_user_memberships.belongs_to_user = $2 
	)
`

// UserIsMemberOfHousehold does a thing.
func (q *SQLQuerier) UserIsMemberOfHousehold(ctx context.Context, userID, householdID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" || householdID == "" {
		return false, ErrInvalidIDProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.UserIDKey:      userID,
		keys.HouseholdIDKey: householdID,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		householdID,
		userID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, userIsMemberOfHouseholdQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing user membership check query")
	}

	return result, nil
}

const modifyUserPermissionsQuery = `
	UPDATE household_user_memberships SET household_roles = $1 WHERE belongs_to_household = $2 AND belongs_to_user = $3
`

// ModifyUserPermissions does a thing.
func (q *SQLQuerier) ModifyUserPermissions(ctx context.Context, householdID, userID string, input *types.ModifyUserPermissionsInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" || userID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.HouseholdIDKey: householdID,
		keys.UserIDKey:      userID,
		"new_roles":         input.NewRoles,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		strings.Join(input.NewRoles, householdMemberRolesSeparator),
		householdID,
		userID,
	}

	// modify the membership.
	if err := q.performWriteQuery(ctx, q.db, "user household permissions modification", modifyUserPermissionsQuery, args); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareError(err, logger, span, "modifying user household permissions")
	}

	logger.Info("user permissions modified")

	return nil
}

const transferHouseholdOwnershipQuery = `
	UPDATE households SET belongs_to_user = $1 WHERE archived_on IS NULL AND belongs_to_user = $2 AND id = $3
`

const transferHouseholdMembershipQuery = `
	UPDATE household_user_memberships SET belongs_to_user = $1 WHERE archived_on IS NULL AND belongs_to_household = $2 AND belongs_to_user = $3
`

// TransferHouseholdOwnership does a thing.
func (q *SQLQuerier) TransferHouseholdOwnership(ctx context.Context, householdID string, input *types.HouseholdOwnershipTransferInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.HouseholdIDKey: householdID,
		"current_owner":     input.CurrentOwner,
		"new_owner":         input.NewOwner,
	})

	tracing.AttachUserIDToSpan(span, input.NewOwner)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// begin household transfer transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	transferHouseholdOwnershipArgs := []interface{}{
		input.NewOwner,
		input.CurrentOwner,
		householdID,
	}

	// create the membership.
	if err = q.performWriteQuery(ctx, tx, "user ownership transfer", transferHouseholdOwnershipQuery, transferHouseholdOwnershipArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "transferring household to new owner")
	}

	transferHouseholdMembershipArgs := []interface{}{
		input.NewOwner,
		householdID,
		input.CurrentOwner,
	}

	// create the membership.
	if err = q.performWriteQuery(ctx, tx, "user memberships transfer", transferHouseholdMembershipQuery, transferHouseholdMembershipArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "transferring household memberships")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("household transferred to new owner")

	return nil
}

const addUserToHouseholdQuery = `
	INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,household_roles)
	VALUES ($1,$2,$3,$4)
`

// AddUserToHousehold does a thing.
func (q *SQLQuerier) AddUserToHousehold(ctx context.Context, input *types.HouseholdUserMembershipDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.UserIDKey:      input.UserID,
		keys.HouseholdIDKey: input.HouseholdID,
	})

	tracing.AttachUserIDToSpan(span, input.UserID)
	tracing.AttachHouseholdIDToSpan(span, input.HouseholdID)

	addUserToHouseholdArgs := []interface{}{
		input.ID,
		input.UserID,
		input.HouseholdID,
		strings.Join(input.HouseholdRoles, householdMemberRolesSeparator),
	}

	// create the membership.
	if err := q.performWriteQuery(ctx, q.db, "user household membership creation", addUserToHouseholdQuery, addUserToHouseholdArgs); err != nil {
		return observability.PrepareError(err, logger, span, "performing user household membership creation query")
	}

	logger.Info("user added to household")

	return nil
}

const removeUserFromHouseholdQuery = `
	DELETE FROM household_user_memberships 
	WHERE household_user_memberships.archived_on IS NULL
	AND household_user_memberships.belongs_to_household = $1 
	AND household_user_memberships.belongs_to_user = $2
`

// RemoveUserFromHousehold removes a user's membership to a household.
func (q *SQLQuerier) RemoveUserFromHousehold(ctx context.Context, userID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" || householdID == "" {
		return ErrInvalidIDProvided
	}

	logger := q.logger.WithValues(map[string]interface{}{
		keys.UserIDKey:      userID,
		keys.HouseholdIDKey: householdID,
	})

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	args := []interface{}{
		householdID,
		userID,
	}

	// create the membership.
	if err := q.performWriteQuery(ctx, q.db, "user membership removal", removeUserFromHouseholdQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "removing user from household")
	}

	logger.Info("user removed from household")

	return nil
}
