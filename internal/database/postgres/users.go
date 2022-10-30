package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/identifiers"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.UserDataManager = (*Querier)(nil)

	// usersTableColumns are the columns for the users table.
	usersTableColumns = []string{
		"users.id",
		"users.username",
		"users.email_address",
		"users.avatar_src",
		"users.hashed_password",
		"users.requires_password_change",
		"users.password_last_changed_at",
		"users.two_factor_secret",
		"users.two_factor_secret_verified_at",
		"users.service_roles",
		"users.user_account_status",
		"users.user_account_status_explanation",
		"users.birth_day",
		"users.birth_month",
		"users.created_at",
		"users.last_updated_at",
		"users.archived_at",
	}
)

const (
	serviceRolesSeparator = commaSeparator
)

// scanUser provides a consistent way to scan something like a *sql.Row into a Requester struct.
func (q *Querier) scanUser(ctx context.Context, scan database.Scanner, includeCounts bool) (user *types.User, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	user = &types.User{
		ServiceRoles: []string{},
	}

	var (
		rawRoles                  string
		passwordLastChangedAt     sql.NullTime
		twoFactorSecretVerifiedAt sql.NullTime
	)

	targetVars := []interface{}{
		&user.ID,
		&user.Username,
		&user.EmailAddress,
		&user.AvatarSrc,
		&user.HashedPassword,
		&user.RequiresPasswordChange,
		&passwordLastChangedAt,
		&user.TwoFactorSecret,
		&twoFactorSecretVerifiedAt,
		&rawRoles,
		&user.AccountStatus,
		&user.AccountStatusExplanation,
		&user.BirthDay,
		&user.BirthMonth,
		&user.CreatedAt,
		&user.LastUpdatedAt,
		&user.ArchivedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "scanning user")
	}

	if roles := strings.Split(rawRoles, serviceRolesSeparator); len(roles) > 0 {
		user.ServiceRoles = roles
	}

	if passwordLastChangedAt.Valid {
		user.PasswordLastChangedAt = &passwordLastChangedAt.Time
	}

	if twoFactorSecretVerifiedAt.Valid {
		user.TwoFactorSecretVerifiedAt = &twoFactorSecretVerifiedAt.Time
	}

	return user, filteredCount, totalCount, nil
}

// scanUsers takes database rows and loads them into a slice of Requester structs.
func (q *Querier) scanUsers(ctx context.Context, rows database.ResultIterator, includeCounts bool) (users []*types.User, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		user, fc, tc, scanErr := q.scanUser(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, observability.PrepareError(scanErr, span, "scanning user result")
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
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return users, filteredCount, totalCount, nil
}

//go:embed queries/users/get_by_id.sql
var getUserByIDQuery string

//go:embed queries/users/get_with_verified_two_factor.sql
var getUserWithVerified2FAQuery string

// getUser fetches a user.
func (q *Querier) getUser(ctx context.Context, userID string, withVerifiedTOTPSecret bool) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

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
		return nil, observability.PrepareError(err, span, "scanning user")
	}

	return u, nil
}

//go:embed queries/users/exists_with_status.sql
var userHasStatusQuery string

// UserHasStatus fetches whether a user has a particular status.
func (q *Querier) UserHasStatus(ctx context.Context, userID string, statuses ...string) (banned bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return false, ErrInvalidIDProvided
	}

	if len(statuses) == 0 {
		return true, nil
	}

	tracing.AttachUserIDToSpan(span, userID)

	args := []interface{}{userID}
	for _, status := range statuses {
		args = append(args, status)
	}

	result, err := q.performBooleanQuery(ctx, q.db, userHasStatusQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, span, "performing user status check")
	}

	return result, nil
}

// GetUser fetches a user.
func (q *Querier) GetUser(ctx context.Context, userID string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)

	return q.getUser(ctx, userID, false)
}

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified 2FA secret.
func (q *Querier) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)

	return q.getUser(ctx, userID, false)
}

//go:embed queries/users/get_by_username.sql
var getUserByUsernameQuery string

// GetUserByUsername fetches a user by their username.
func (q *Querier) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if username == "" {
		return nil, ErrEmptyInputProvided
	}

	tracing.AttachUsernameToSpan(span, username)

	args := []interface{}{username}

	row := q.getOneRow(ctx, q.db, "user", getUserByUsernameQuery, args)

	u, _, _, err := q.scanUser(ctx, row, false)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, span, "scanning user")
	}

	return u, nil
}

//go:embed queries/users/get_admin_by_username.sql
var getAdminUserByUsernameQuery string

// GetAdminUserByUsername fetches a user by their username.
func (q *Querier) GetAdminUserByUsername(ctx context.Context, username string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if username == "" {
		return nil, ErrEmptyInputProvided
	}

	tracing.AttachUsernameToSpan(span, username)

	args := []interface{}{username}

	row := q.getOneRow(ctx, q.db, "admin user fetch", getAdminUserByUsernameQuery, args)

	u, _, _, err := q.scanUser(ctx, row, false)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, span, "scanning user")
	}

	return u, nil
}

//go:embed queries/users/get_by_email.sql
var getUserByEmailQuery string

// GetUserByEmail fetches a user by their email.
func (q *Querier) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if email == "" {
		return nil, ErrEmptyInputProvided
	}

	tracing.AttachEmailAddressToSpan(span, email)

	args := []interface{}{email}
	row := q.getOneRow(ctx, q.db, "user", getUserByEmailQuery, args)

	u, _, _, err := q.scanUser(ctx, row, false)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, span, "scanning user")
	}

	return u, nil
}

//go:embed queries/users/search_by_username.sql
var searchForUserByUsernameQuery string

// SearchForUsersByUsername fetches a list of users whose usernames begin with a given query.
func (q *Querier) SearchForUsersByUsername(ctx context.Context, usernameQuery string) ([]*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if usernameQuery == "" {
		return []*types.User{}, ErrEmptyInputProvided
	}

	tracing.AttachSearchQueryToSpan(span, usernameQuery)

	args := []interface{}{
		wrapQueryForILIKE(usernameQuery),
	}

	rows, err := q.getRows(ctx, q.db, "user search by username", searchForUserByUsernameQuery, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, span, "querying database for users")
	}

	u, _, _, err := q.scanUsers(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning user")
	}

	return u, nil
}

// GetUsers fetches a list of users from the database that meet a particular filter.
func (q *Querier) GetUsers(ctx context.Context, filter *types.QueryFilter) (x *types.UserList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.UserList{}

	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "users", nil, nil, nil, "", usersTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "users", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning user")
	}

	if x.Users, x.FilteredCount, x.TotalCount, err = q.scanUsers(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, span, "loading response from database")
	}

	return x, nil
}

//go:embed queries/users/create.sql
var userCreationQuery string

//go:embed queries/users/create_household_memberships_for_new_user.sql
var createHouseholdMembershipForNewUserQuery string

// CreateUser creates a user.
func (q *Querier) CreateUser(ctx context.Context, input *types.UserDatabaseCreationInput) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachUsernameToSpan(span, input.Username)
	logger := q.logger.WithValues(map[string]interface{}{
		keys.UsernameKey:                 input.Username,
		keys.UserEmailAddressKey:         input.EmailAddress,
		keys.HouseholdInvitationTokenKey: input.InvitationToken,
		"destination_household":          input.DestinationHouseholdID,
	})

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
		CreatedAt:       q.currentTime(),
	}
	logger = logger.WithValue(keys.UserIDKey, user.ID)
	tracing.AttachUserIDToSpan(span, user.ID)

	// begin user creation transaction
	tx, beginTransactionErr := q.db.BeginTx(ctx, nil)
	if beginTransactionErr != nil {
		return nil, observability.PrepareError(beginTransactionErr, span, "beginning transaction")
	}

	if writeErr := q.performWriteQuery(ctx, tx, "user creation", userCreationQuery, userCreationArgs); writeErr != nil {
		q.rollbackTransaction(ctx, tx)

		var e *pq.Error
		if errors.As(writeErr, &e) {
			if e.Code == postgresDuplicateEntryErrorCode {
				return nil, database.ErrUserAlreadyExists
			}
		}

		return nil, observability.PrepareError(writeErr, span, "creating user")
	}

	hasValidInvite := input.InvitationToken != "" && input.DestinationHouseholdID != ""

	if err := q.createHouseholdForUser(ctx, tx, hasValidInvite, user.ID); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating household for new user")
	}

	logger.Debug("household created")

	if hasValidInvite {
		if err := q.acceptInvitationForUser(ctx, tx, input); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "accepting household invitation")
		}
		logger.Debug("accepted invitation and joined household for user")
	}

	if err := q.attachInvitationsToUser(ctx, tx, user.EmailAddress, user.ID); err != nil {
		q.rollbackTransaction(ctx, tx)
		logger = logger.WithValue("email_address", user.EmailAddress).WithValue("user_id", user.ID)
		return nil, observability.PrepareAndLogError(err, logger, span, "attaching existing invitations to new user")
	}

	if err := tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Debug("user and household created")

	return user, nil
}

func (q *Querier) createHouseholdForUser(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, hasValidInvite bool, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	// standard registration: we need to create the household
	householdID := identifiers.New()
	tracing.AttachHouseholdIDToSpan(span, householdID)

	householdCreationInput := &types.HouseholdCreationRequestInput{
		ID:            householdID,
		Name:          fmt.Sprintf("%s_default", userID),
		BelongsToUser: userID,
		TimeZone:      types.DefaultHouseholdTimeZone,
	}

	householdCreationArgs := []interface{}{
		householdCreationInput.ID,
		householdCreationInput.Name,
		types.UnpaidHouseholdBillingStatus,
		householdCreationInput.ContactEmail,
		householdCreationInput.ContactPhone,
		householdCreationInput.TimeZone,
		householdCreationInput.BelongsToUser,
	}

	if writeErr := q.performWriteQuery(ctx, querier, "household creation", householdCreationQuery, householdCreationArgs); writeErr != nil {
		q.rollbackTransaction(ctx, querier)
		return observability.PrepareError(writeErr, span, "create household")
	}

	createHouseholdMembershipForNewUserArgs := []interface{}{
		identifiers.New(),
		userID,
		householdID,
		!hasValidInvite,
		authorization.HouseholdAdminRole.String(),
	}

	if err := q.performWriteQuery(ctx, querier, "household user membership creation", createHouseholdMembershipForNewUserQuery, createHouseholdMembershipForNewUserArgs); err != nil {
		q.rollbackTransaction(ctx, querier)
		return observability.PrepareError(err, span, "writing household user membership")
	}

	return nil
}

//go:embed queries/users/update.sql
var updateUserQuery string

// UpdateUser receives a complete Requester struct and updates its record in the database.
// NOTE: this function uses the ID provided in the input to make its query.
func (q *Querier) UpdateUser(ctx context.Context, updated *types.User) error {
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
		updated.TwoFactorSecretVerifiedAt,
		updated.BirthDay,
		updated.BirthMonth,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "user update", updateUserQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user")
	}

	logger.Info("user updated")

	return nil
}

/* #nosec G101 */
//go:embed queries/users/update_password.sql
var updateUserPasswordQuery string

// UpdateUserPassword updates a user's passwords hash in the database.
func (q *Querier) UpdateUserPassword(ctx context.Context, userID, newHash string) error {
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
		return observability.PrepareAndLogError(err, logger, span, "updating user password")
	}

	logger.Info("user password updated")

	return nil
}

/* #nosec G101 */
//go:embed queries/users/update_two_factor_secret.sql
var updateUserTwoFactorSecretQuery string

// UpdateUserTwoFactorSecret marks a user's two factor secret as validated.
func (q *Querier) UpdateUserTwoFactorSecret(ctx context.Context, userID, newSecret string) error {
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
		return observability.PrepareAndLogError(err, logger, span, "updating user 2FA secret")
	}
	logger.Info("user two factor secret updated")

	return nil
}

/* #nosec G101 */
//go:embed queries/users/mark_two_factor_secret_as_verified.sql
var markUserTwoFactorSecretAsVerified string

// MarkUserTwoFactorSecretAsVerified marks a user's two factor secret as validated.
func (q *Querier) MarkUserTwoFactorSecretAsVerified(ctx context.Context, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	args := []interface{}{
		types.GoodStandingUserAccountStatus,
		userID,
	}

	if err := q.performWriteQuery(ctx, q.db, "user two factor secret verification", markUserTwoFactorSecretAsVerified, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "writing verified two factor status to database")
	}

	logger.Info("user two factor secret verified")

	return nil
}

//go:embed queries/users/archive.sql
var archiveUserQuery string

//go:embed queries/users/archive_memberships.sql
var archiveMembershipsQuery string

// ArchiveUser archives a user.
func (q *Querier) ArchiveUser(ctx context.Context, userID string) error {
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
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	archiveUserArgs := []interface{}{userID}

	if err = q.performWriteQuery(ctx, tx, "user archive", archiveUserQuery, archiveUserArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving user")
	}

	archiveMembershipsArgs := []interface{}{userID}

	if err = q.performWriteQuery(ctx, tx, "user memberships archive", archiveMembershipsQuery, archiveMembershipsArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving user household memberships")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("user archived")

	return nil
}
