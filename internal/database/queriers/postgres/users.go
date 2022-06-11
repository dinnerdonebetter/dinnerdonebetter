package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.UserDataManager = (*SQLQuerier)(nil)

	// usersTableColumns are the columns for the users table.
	usersTableColumns = []string{
		"users.id",
		"users.username",
		"users.email_address",
		"users.avatar_src",
		"users.hashed_password",
		"users.requires_password_change",
		"users.password_last_changed_on",
		"users.two_factor_secret",
		"users.two_factor_secret_verified_on",
		"users.service_roles",
		"users.reputation",
		"users.reputation_explanation",
		"users.birth_day",
		"users.birth_month",
		"users.created_on",
		"users.last_updated_on",
		"users.archived_on",
	}
)

const (
	serviceRolesSeparator = commaSeparator
)

// scanUser provides a consistent way to scan something like a *sql.Row into a Requester struct.
func (q *SQLQuerier) scanUser(ctx context.Context, scan database.Scanner, includeCounts bool) (user *types.User, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)
	user = &types.User{
		ServiceRoles: []string{},
	}
	var rawRoles string

	targetVars := []interface{}{
		&user.ID,
		&user.Username,
		&user.EmailAddress,
		&user.AvatarSrc,
		&user.HashedPassword,
		&user.RequiresPasswordChange,
		&user.PasswordLastChangedOn,
		&user.TwoFactorSecret,
		&user.TwoFactorSecretVerifiedOn,
		&rawRoles,
		&user.ServiceHouseholdStatus,
		&user.ReputationExplanation,
		&user.BirthDay,
		&user.BirthMonth,
		&user.CreatedOn,
		&user.LastUpdatedOn,
		&user.ArchivedOn,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "scanning user")
	}

	if roles := strings.Split(rawRoles, serviceRolesSeparator); len(roles) > 0 {
		user.ServiceRoles = roles
	}

	return user, filteredCount, totalCount, nil
}

// scanUsers takes database rows and loads them into a slice of Requester structs.
func (q *SQLQuerier) scanUsers(ctx context.Context, rows database.ResultIterator, includeCounts bool) (users []*types.User, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		user, fc, tc, scanErr := q.scanUser(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, observability.PrepareError(scanErr, logger, span, "scanning user result")
		}

		if includeCounts && filteredCount == 0 {
			filteredCount = fc
		}

		if includeCounts && totalCount == 0 {
			totalCount = tc
		}

		users = append(users, user)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return users, filteredCount, totalCount, nil
}

const getUserByIDQuery = `
	SELECT
		users.id,
		users.username,
		users.email_address,
		users.avatar_src,
		users.hashed_password,
		users.requires_password_change,
		users.password_last_changed_on,
		users.two_factor_secret,
		users.two_factor_secret_verified_on,
		users.service_roles,
		users.reputation,
		users.reputation_explanation,
		users.birth_day,
		users.birth_month,
		users.created_on,
		users.last_updated_on,
		users.archived_on
	FROM users
	WHERE users.archived_on IS NULL
	AND users.id = $1
`

const getUserWithVerified2FAQuery = getUserByIDQuery + `
	AND users.two_factor_secret_verified_on IS NOT NULL
`

// getUser fetches a user.
func (q *SQLQuerier) getUser(ctx context.Context, userID string, withVerifiedTOTPSecret bool) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	var query string
	args := []interface{}{userID}

	if withVerifiedTOTPSecret {
		query = getUserWithVerified2FAQuery
	} else {
		query = getUserByIDQuery
	}

	row := q.getOneRow(ctx, q.db, "user", query, args)

	u, _, _, err := q.scanUser(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning user")
	}

	return u, nil
}

const userHasStatusQuery = `
	SELECT EXISTS ( SELECT users.id FROM users WHERE users.archived_on IS NULL AND users.id = $1 AND (users.reputation = $2 OR users.reputation = $3) )
`

// UserHasStatus fetches whether an user has a particular status.
func (q *SQLQuerier) UserHasStatus(ctx context.Context, userID string, statuses ...string) (banned bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return false, ErrInvalidIDProvided
	}

	if len(statuses) == 0 {
		return true, nil
	}

	logger := q.logger.WithValue(keys.UserIDKey, userID).WithValue("statuses", statuses)
	tracing.AttachUserIDToSpan(span, userID)

	args := []interface{}{userID}
	for _, status := range statuses {
		args = append(args, status)
	}

	result, err := q.performBooleanQuery(ctx, q.db, userHasStatusQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing user status check")
	}

	return result, nil
}

// GetUser fetches a user.
func (q *SQLQuerier) GetUser(ctx context.Context, userID string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)

	return q.getUser(ctx, userID, false)
}

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified 2FA secret.
func (q *SQLQuerier) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)

	return q.getUser(ctx, userID, false)
}

const getUserByUsernameQuery = `
	SELECT
		users.id,
		users.username,
		users.email_address,
		users.avatar_src,
		users.hashed_password,
		users.requires_password_change,
		users.password_last_changed_on,
		users.two_factor_secret,
		users.two_factor_secret_verified_on,
		users.service_roles,
		users.reputation,
		users.reputation_explanation,
		users.birth_day,
		users.birth_month,
		users.created_on,
		users.last_updated_on,
		users.archived_on
	FROM users
	WHERE users.archived_on IS NULL
	AND users.username = $1
`

// GetUserByUsername fetches a user by their username.
func (q *SQLQuerier) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if username == "" {
		return nil, ErrEmptyInputProvided
	}

	tracing.AttachUsernameToSpan(span, username)
	logger := q.logger.WithValue(keys.UsernameKey, username)

	args := []interface{}{username}

	row := q.getOneRow(ctx, q.db, "user", getUserByUsernameQuery, args)

	u, _, _, err := q.scanUser(ctx, row, false)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, logger, span, "scanning user")
	}

	return u, nil
}

const getAdminUserByUsernameQuery = `
	SELECT
		users.id,
		users.username,
		users.email_address,
		users.avatar_src,
		users.hashed_password,
		users.requires_password_change,
		users.password_last_changed_on,
		users.two_factor_secret,
		users.two_factor_secret_verified_on,
		users.service_roles,
		users.reputation,
		users.reputation_explanation,
		users.birth_day,
		users.birth_month,
		users.created_on,
		users.last_updated_on,
		users.archived_on
	FROM users
	WHERE users.archived_on IS NULL
	AND users.service_roles ILIKE '%service_admin%'
	AND users.username = $1
	AND users.two_factor_secret_verified_on IS NOT NULL
`

// GetAdminUserByUsername fetches a user by their username.
func (q *SQLQuerier) GetAdminUserByUsername(ctx context.Context, username string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if username == "" {
		return nil, ErrEmptyInputProvided
	}

	tracing.AttachUsernameToSpan(span, username)
	logger := q.logger.WithValue(keys.UsernameKey, username)

	args := []interface{}{username}

	row := q.getOneRow(ctx, q.db, "admin user fetch", getAdminUserByUsernameQuery, args)

	u, _, _, err := q.scanUser(ctx, row, false)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, logger, span, "scanning user")
	}

	return u, nil
}

const getUserIDByEmailQuery = `
	SELECT
		users.id
	FROM users
	WHERE users.archived_on IS NULL
	AND users.email_address = $1
	AND users.two_factor_secret_verified_on IS NOT NULL
`

// GetUserIDByEmail fetches a user by their email.
func (q *SQLQuerier) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if email == "" {
		return "", ErrEmptyInputProvided
	}

	tracing.AttachEmailAddressToSpan(span, email)
	logger := q.logger.WithValue(keys.UserEmailAddressKey, email)

	args := []interface{}{email}
	row := q.getOneRow(ctx, q.db, "user", getUserIDByEmailQuery, args)

	var userID string
	if err := row.Scan(&userID); err != nil {
		return "", observability.PrepareError(err, logger, span, "scanning user ID")
	}

	return userID, nil
}

const searchForUserByUsernameQuery = `SELECT 
	users.id,
	users.username,
	users.email_address,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_on,
	users.two_factor_secret,
	users.two_factor_secret_verified_on,
	users.service_roles,
	users.reputation,
	users.reputation_explanation,
	users.birth_day,
	users.birth_month,
	users.created_on,
	users.last_updated_on,
	users.archived_on 
FROM users
WHERE users.username ILIKE $1
AND users.archived_on IS NULL
AND users.two_factor_secret_verified_on IS NOT NULL
`

// SearchForUsersByUsername fetches a list of users whose usernames begin with a given query.
func (q *SQLQuerier) SearchForUsersByUsername(ctx context.Context, usernameQuery string) ([]*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if usernameQuery == "" {
		return []*types.User{}, ErrEmptyInputProvided
	}

	tracing.AttachSearchQueryToSpan(span, usernameQuery)
	logger := q.logger.WithValue(keys.SearchQueryKey, usernameQuery)

	args := []interface{}{
		wrapQueryForILIKE(usernameQuery),
	}

	rows, err := q.performReadQuery(ctx, q.db, "user search by username", searchForUserByUsernameQuery, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, logger, span, "querying database for users")
	}

	u, _, _, err := q.scanUsers(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning user")
	}

	return u, nil
}

const getAllUsersCountQuery = `SELECT COUNT(users.id) FROM users WHERE users.archived_on IS NULL`

// GetAllUsersCount fetches a count of users from the database that meet a particular filter.
func (q *SQLQuerier) GetAllUsersCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getAllUsersCountQuery, "fetching count of users")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of users")
	}

	return count, nil
}

// GetUsers fetches a list of users from the database that meet a particular filter.
func (q *SQLQuerier) GetUsers(ctx context.Context, filter *types.QueryFilter) (x *types.UserList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.UserList{}

	tracing.AttachQueryFilterToSpan(span, filter)
	logger := filter.AttachToLogger(q.logger)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "users", nil, nil, nil, "", usersTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "users", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning user")
	}

	if x.Users, x.FilteredCount, x.TotalCount, err = q.scanUsers(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "loading response from database")
	}

	return x, nil
}

const userCreationQuery = `
	INSERT INTO users (id,username,email_address,hashed_password,two_factor_secret,avatar_src,reputation,birth_day,birth_month,service_roles) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
`

const createHouseholdMembershipForNewUserQuery = `
	INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,default_household,household_roles)
	VALUES ($1,$2,$3,$4,$5)
`

// CreateUser creates a user.
func (q *SQLQuerier) CreateUser(ctx context.Context, input *types.UserDatabaseCreationInput) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachUsernameToSpan(span, input.Username)
	logger := q.logger.WithValue(keys.UsernameKey, input.Username)

	userCreationArgs := []interface{}{
		input.ID,
		input.Username,
		input.EmailAddress,
		input.HashedPassword,
		input.TwoFactorSecret,
		input.AvatarSrc,
		types.UnverifiedHouseholdStatus,
		input.BirthDay,
		input.BirthMonth,
		authorization.ServiceUserRole.String(),
	}

	user := &types.User{
		ID:              input.ID,
		Username:        input.Username,
		EmailAddress:    input.EmailAddress,
		HashedPassword:  input.HashedPassword,
		TwoFactorSecret: input.TwoFactorSecret,
		BirthMonth:      input.BirthMonth,
		BirthDay:        input.BirthDay,
		ServiceRoles:    []string{authorization.ServiceUserRole.String()},
		CreatedOn:       q.currentTime(),
	}
	logger = logger.WithValue(keys.UserIDKey, user.ID)
	tracing.AttachUserIDToSpan(span, user.ID)

	// begin user creation transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	if writeErr := q.performWriteQuery(ctx, tx, "user creation", userCreationQuery, userCreationArgs); writeErr != nil {
		q.rollbackTransaction(ctx, tx)

		var e *pq.Error
		if errors.As(err, &e) {
			if e.Code == "23505" {
				observability.AcknowledgeError(database.ErrUserAlreadyExists, logger, span, "creating user")
				return nil, database.ErrUserAlreadyExists
			}
		}

		return nil, observability.PrepareError(writeErr, logger, span, "creating user")
	}

	hasValidInvite := input.InvitationToken != "" && input.DestinationHousehold != ""

	// standard registration: we need to create the household
	householdID := ksuid.New().String()
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	householdCreationInput := types.HouseholdCreationInputForNewUser(user)
	householdCreationInput.ID = householdID

	householdCreationArgs := []interface{}{
		householdCreationInput.ID,
		householdCreationInput.Name,
		types.UnpaidHouseholdBillingStatus,
		householdCreationInput.ContactEmail,
		householdCreationInput.ContactPhone,
		householdCreationInput.BelongsToUser,
	}

	if writeErr := q.performWriteQuery(ctx, tx, "household creation", householdCreationQuery, householdCreationArgs); writeErr != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(writeErr, logger, span, "create household")
	}

	createHouseholdMembershipForNewUserArgs := []interface{}{
		ksuid.New().String(),
		user.ID,
		householdID,
		!hasValidInvite,
		authorization.HouseholdAdminRole.String(),
	}

	if err = q.performWriteQuery(ctx, tx, "household user membership creation", createHouseholdMembershipForNewUserQuery, createHouseholdMembershipForNewUserArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "writing household user membership")
	}

	if hasValidInvite {
		invitation, tokenCheckErr := q.GetHouseholdInvitationByEmailAndToken(ctx, input.EmailAddress, input.InvitationToken)
		if tokenCheckErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(tokenCheckErr, logger, span, "creating user")
		}

		createHouseholdMembershipForNewUserArgs = []interface{}{
			ksuid.New().String(),
			user.ID,
			input.DestinationHousehold,
			true,
			authorization.HouseholdMemberRole.String(),
		}

		if err = q.performWriteQuery(ctx, tx, "household user membership creation", createHouseholdMembershipForNewUserQuery, createHouseholdMembershipForNewUserArgs); err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(err, logger, span, "writing destination household membership")
		}

		if err = q.setInvitationStatus(ctx, tx, invitation.DestinationHousehold, invitation.ID, "", types.AcceptedHouseholdInvitationStatus); err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareError(err, logger, span, "accepting household invitation")
		}
	}

	if err = q.attachInvitationsToUser(ctx, tx, user.EmailAddress, user.ID); err != nil {
		q.rollbackTransaction(ctx, tx)
		logger = logger.WithValue("email_address", user.EmailAddress).WithValue("user_id", user.ID)
		return nil, observability.PrepareError(err, logger, span, "writing household user membership")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Debug("user and household created")

	return user, nil
}

const updateUserQuery = `UPDATE users SET 
	username = $1,
	hashed_password = $2,
	avatar_src = $3,
	two_factor_secret = $4,
	two_factor_secret_verified_on = $5,
	birth_day = $6,
	birth_month = $7,
	last_updated_on = extract(epoch FROM NOW()) 
WHERE archived_on IS NULL 
AND id = $8
`

// UpdateUser receives a complete Requester struct and updates its record in the database.
// NOTE: this function uses the ID provided in the input to make its query.
func (q *SQLQuerier) UpdateUser(ctx context.Context, updated *types.User) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	tracing.AttachUsernameToSpan(span, updated.Username)
	logger := q.logger.WithValue(keys.UsernameKey, updated.Username)

	args := []interface{}{
		updated.Username,
		updated.HashedPassword,
		updated.AvatarSrc,
		updated.TwoFactorSecret,
		updated.TwoFactorSecretVerifiedOn,
		updated.BirthDay,
		updated.BirthMonth,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "user update", updateUserQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating user")
	}

	logger.Info("user updated")

	return nil
}

/* #nosec G101 */
const updateUserPasswordQuery = `UPDATE users SET
	hashed_password = $1,
	requires_password_change = $2,
	password_last_changed_on = extract(epoch FROM NOW()),
	last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND id = $3
`

// UpdateUserPassword updates a user's passwords hash in the database.
func (q *SQLQuerier) UpdateUserPassword(ctx context.Context, userID, newHash string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	if newHash == "" {
		return ErrEmptyInputProvided
	}

	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	args := []interface{}{
		newHash,
		false,
		userID,
	}

	if err := q.performWriteQuery(ctx, q.db, "user passwords update", updateUserPasswordQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating user password")
	}

	logger.Info("user password updated")

	return nil
}

/* #nosec G101 */
const updateUserTwoFactorSecretQuery = `UPDATE users SET 
	two_factor_secret_verified_on = $1,
	two_factor_secret = $2
WHERE archived_on IS NULL
AND id = $3
`

// UpdateUserTwoFactorSecret marks a user's two factor secret as validated.
func (q *SQLQuerier) UpdateUserTwoFactorSecret(ctx context.Context, userID, newSecret string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	if newSecret == "" {
		return ErrEmptyInputProvided
	}

	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	args := []interface{}{
		nil,
		newSecret,
		userID,
	}

	if err := q.performWriteQuery(ctx, q.db, "user 2FA secret update", updateUserTwoFactorSecretQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating user 2FA secret")
	}
	logger.Info("user two factor secret updated")

	return nil
}

/* #nosec G101 */
const markUserTwoFactorSecretAsVerified = `UPDATE users SET
	two_factor_secret_verified_on = extract(epoch FROM NOW()),
	reputation = $1
WHERE archived_on IS NULL
AND id = $2
`

// MarkUserTwoFactorSecretAsVerified marks a user's two factor secret as validated.
func (q *SQLQuerier) MarkUserTwoFactorSecretAsVerified(ctx context.Context, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	args := []interface{}{
		types.GoodStandingHouseholdStatus,
		userID,
	}

	if err := q.performWriteQuery(ctx, q.db, "user two factor secret verification", markUserTwoFactorSecretAsVerified, args); err != nil {
		return observability.PrepareError(err, logger, span, "writing verified two factor status to database")
	}

	logger.Info("user two factor secret verified")

	return nil
}

const archiveUserQuery = `UPDATE users SET
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND id = $1
`

const archiveMembershipsQuery = `UPDATE household_user_memberships SET
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND belongs_to_user = $1
`

// ArchiveUser archives a user.
func (q *SQLQuerier) ArchiveUser(ctx context.Context, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	// begin archive user transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareError(err, logger, span, "beginning transaction")
	}

	archiveUserArgs := []interface{}{userID}

	if err = q.performWriteQuery(ctx, tx, "user archive", archiveUserQuery, archiveUserArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "archiving user")
	}

	archiveMembershipsArgs := []interface{}{userID}

	if err = q.performWriteQuery(ctx, tx, "user memberships archive", archiveMembershipsQuery, archiveMembershipsArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareError(err, logger, span, "archiving user household memberships")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareError(err, logger, span, "committing transaction")
	}

	logger.Info("user archived")

	return nil
}
