package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	_ types.UserDataManager = (*Querier)(nil)
)

// GetUser fetches a user.
func (q *Querier) GetUser(ctx context.Context, userID string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)

	result, err := q.generatedQuerier.GetUserByID(ctx, q.db, userID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting user with verified two factor")
	}

	u := &types.User{
		CreatedAt:                  result.CreatedAt,
		PasswordLastChangedAt:      timePointerFromNullTime(result.PasswordLastChangedAt),
		LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
		LastAcceptedTermsOfService: timePointerFromNullTime(result.LastAcceptedTermsOfService),
		LastAcceptedPrivacyPolicy:  timePointerFromNullTime(result.LastAcceptedPrivacyPolicy),
		TwoFactorSecretVerifiedAt:  timePointerFromNullTime(result.TwoFactorSecretVerifiedAt),
		Birthday:                   timePointerFromNullTime(result.Birthday),
		ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
		AccountStatusExplanation:   result.UserAccountStatusExplanation,
		TwoFactorSecret:            result.TwoFactorSecret,
		HashedPassword:             result.HashedPassword,
		ID:                         result.ID,
		AccountStatus:              result.UserAccountStatus,
		Username:                   result.Username,
		FirstName:                  result.FirstName,
		LastName:                   result.LastName,
		EmailAddress:               result.EmailAddress,
		EmailAddressVerifiedAt:     timePointerFromNullTime(result.EmailAddressVerifiedAt),
		AvatarSrc:                  stringPointerFromNullString(result.AvatarSrc),
		ServiceRole:                result.ServiceRole,
		RequiresPasswordChange:     result.RequiresPasswordChange,
	}

	return u, nil
}

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified 2FA secret.
func (q *Querier) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)

	result, err := q.generatedQuerier.GetUserWithUnverifiedTwoFactor(ctx, q.db, userID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting user with verified two factor")
	}

	u := &types.User{
		CreatedAt:                  result.CreatedAt,
		PasswordLastChangedAt:      timePointerFromNullTime(result.PasswordLastChangedAt),
		LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
		LastAcceptedTermsOfService: timePointerFromNullTime(result.LastAcceptedTermsOfService),
		LastAcceptedPrivacyPolicy:  timePointerFromNullTime(result.LastAcceptedPrivacyPolicy),
		TwoFactorSecretVerifiedAt:  timePointerFromNullTime(result.TwoFactorSecretVerifiedAt),
		Birthday:                   timePointerFromNullTime(result.Birthday),
		ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
		AccountStatusExplanation:   result.UserAccountStatusExplanation,
		TwoFactorSecret:            result.TwoFactorSecret,
		HashedPassword:             result.HashedPassword,
		ID:                         result.ID,
		AccountStatus:              result.UserAccountStatus,
		Username:                   result.Username,
		FirstName:                  result.FirstName,
		LastName:                   result.LastName,
		EmailAddress:               result.EmailAddress,
		EmailAddressVerifiedAt:     timePointerFromNullTime(result.EmailAddressVerifiedAt),
		AvatarSrc:                  stringPointerFromNullString(result.AvatarSrc),
		ServiceRole:                result.ServiceRole,
		RequiresPasswordChange:     result.RequiresPasswordChange,
	}

	return u, nil
}

// GetUserByUsername fetches a user by their username.
func (q *Querier) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if username == "" {
		return nil, ErrEmptyInputProvided
	}
	tracing.AttachUsernameToSpan(span, username)

	result, err := q.generatedQuerier.GetUserByUsername(ctx, q.db, username)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting user by username")
	}

	u := &types.User{
		CreatedAt:                  result.CreatedAt,
		PasswordLastChangedAt:      timePointerFromNullTime(result.PasswordLastChangedAt),
		LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
		LastAcceptedTermsOfService: timePointerFromNullTime(result.LastAcceptedTermsOfService),
		LastAcceptedPrivacyPolicy:  timePointerFromNullTime(result.LastAcceptedPrivacyPolicy),
		TwoFactorSecretVerifiedAt:  timePointerFromNullTime(result.TwoFactorSecretVerifiedAt),
		Birthday:                   timePointerFromNullTime(result.Birthday),
		ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
		AccountStatusExplanation:   result.UserAccountStatusExplanation,
		TwoFactorSecret:            result.TwoFactorSecret,
		HashedPassword:             result.HashedPassword,
		ID:                         result.ID,
		AccountStatus:              result.UserAccountStatus,
		Username:                   result.Username,
		FirstName:                  result.FirstName,
		LastName:                   result.LastName,
		EmailAddress:               result.EmailAddress,
		EmailAddressVerifiedAt:     timePointerFromNullTime(result.EmailAddressVerifiedAt),
		AvatarSrc:                  stringPointerFromNullString(result.AvatarSrc),
		ServiceRole:                result.ServiceRole,
		RequiresPasswordChange:     result.RequiresPasswordChange,
	}

	return u, nil
}

// GetAdminUserByUsername fetches a user by their username.
func (q *Querier) GetAdminUserByUsername(ctx context.Context, username string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if username == "" {
		return nil, ErrEmptyInputProvided
	}
	tracing.AttachUsernameToSpan(span, username)

	result, err := q.generatedQuerier.GetAdminUserByUsername(ctx, q.db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, observability.PrepareError(err, span, "getting admin user by username")
	}

	u := &types.User{
		CreatedAt:                  result.CreatedAt,
		PasswordLastChangedAt:      timePointerFromNullTime(result.PasswordLastChangedAt),
		LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
		LastAcceptedTermsOfService: timePointerFromNullTime(result.LastAcceptedTermsOfService),
		LastAcceptedPrivacyPolicy:  timePointerFromNullTime(result.LastAcceptedPrivacyPolicy),
		TwoFactorSecretVerifiedAt:  timePointerFromNullTime(result.TwoFactorSecretVerifiedAt),
		Birthday:                   timePointerFromNullTime(result.Birthday),
		ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
		AccountStatusExplanation:   result.UserAccountStatusExplanation,
		TwoFactorSecret:            result.TwoFactorSecret,
		HashedPassword:             result.HashedPassword,
		ID:                         result.ID,
		AccountStatus:              result.UserAccountStatus,
		Username:                   result.Username,
		FirstName:                  result.FirstName,
		LastName:                   result.LastName,
		EmailAddress:               result.EmailAddress,
		EmailAddressVerifiedAt:     timePointerFromNullTime(result.EmailAddressVerifiedAt),
		AvatarSrc:                  stringPointerFromNullString(result.AvatarSrc),
		ServiceRole:                result.ServiceRole,
		RequiresPasswordChange:     result.RequiresPasswordChange,
	}

	return u, nil
}

// GetUserByEmail fetches a user by their email.
func (q *Querier) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if email == "" {
		return nil, ErrEmptyInputProvided
	}
	tracing.AttachEmailAddressToSpan(span, email)

	result, err := q.generatedQuerier.GetUserByEmail(ctx, q.db, email)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting user by email")
	}

	u := &types.User{
		CreatedAt:                  result.CreatedAt,
		PasswordLastChangedAt:      timePointerFromNullTime(result.PasswordLastChangedAt),
		LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
		LastAcceptedTermsOfService: timePointerFromNullTime(result.LastAcceptedTermsOfService),
		LastAcceptedPrivacyPolicy:  timePointerFromNullTime(result.LastAcceptedPrivacyPolicy),
		TwoFactorSecretVerifiedAt:  timePointerFromNullTime(result.TwoFactorSecretVerifiedAt),
		Birthday:                   timePointerFromNullTime(result.Birthday),
		ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
		AccountStatusExplanation:   result.UserAccountStatusExplanation,
		TwoFactorSecret:            result.TwoFactorSecret,
		HashedPassword:             result.HashedPassword,
		ID:                         result.ID,
		AccountStatus:              result.UserAccountStatus,
		Username:                   result.Username,
		FirstName:                  result.FirstName,
		LastName:                   result.LastName,
		EmailAddress:               result.EmailAddress,
		EmailAddressVerifiedAt:     timePointerFromNullTime(result.EmailAddressVerifiedAt),
		AvatarSrc:                  stringPointerFromNullString(result.AvatarSrc),
		ServiceRole:                result.ServiceRole,
		RequiresPasswordChange:     result.RequiresPasswordChange,
	}

	return u, nil
}

// SearchForUsersByUsername fetches a list of users whose usernames begin with a given query.
func (q *Querier) SearchForUsersByUsername(ctx context.Context, usernameQuery string) ([]*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if usernameQuery == "" {
		return []*types.User{}, ErrEmptyInputProvided
	}
	tracing.AttachSearchQueryToSpan(span, usernameQuery)

	results, err := q.generatedQuerier.SearchUsersByUsername(ctx, q.db, usernameQuery)
	if err != nil {
		return nil, observability.PrepareError(err, span, "querying database for users")
	}

	users := make([]*types.User, len(results))
	for i, result := range results {
		users[i] = &types.User{
			CreatedAt:                  result.CreatedAt,
			PasswordLastChangedAt:      timePointerFromNullTime(result.PasswordLastChangedAt),
			LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
			LastAcceptedTermsOfService: timePointerFromNullTime(result.LastAcceptedTermsOfService),
			LastAcceptedPrivacyPolicy:  timePointerFromNullTime(result.LastAcceptedPrivacyPolicy),
			TwoFactorSecretVerifiedAt:  timePointerFromNullTime(result.TwoFactorSecretVerifiedAt),
			Birthday:                   timePointerFromNullTime(result.Birthday),
			ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
			AccountStatusExplanation:   result.UserAccountStatusExplanation,
			TwoFactorSecret:            result.TwoFactorSecret,
			HashedPassword:             result.HashedPassword,
			ID:                         result.ID,
			AccountStatus:              result.UserAccountStatus,
			Username:                   result.Username,
			FirstName:                  result.FirstName,
			LastName:                   result.LastName,
			EmailAddress:               result.EmailAddress,
			EmailAddressVerifiedAt:     timePointerFromNullTime(result.EmailAddressVerifiedAt),
			AvatarSrc:                  stringPointerFromNullString(result.AvatarSrc),
			ServiceRole:                result.ServiceRole,
			RequiresPasswordChange:     result.RequiresPasswordChange,
		}
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users, nil
}

// GetUsers fetches a list of users from the database that meet a particular filter.
func (q *Querier) GetUsers(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.User], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	filter.AttachToLogger(logger)

	x = &types.QueryFilteredResult[types.User]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetUsers(ctx, q.db, &generated.GetUsersParams{
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning user")
	}

	for _, result := range results {
		u := &types.User{
			CreatedAt:                  result.CreatedAt,
			PasswordLastChangedAt:      timePointerFromNullTime(result.PasswordLastChangedAt),
			LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
			LastAcceptedTermsOfService: timePointerFromNullTime(result.LastAcceptedTermsOfService),
			LastAcceptedPrivacyPolicy:  timePointerFromNullTime(result.LastAcceptedPrivacyPolicy),
			TwoFactorSecretVerifiedAt:  timePointerFromNullTime(result.TwoFactorSecretVerifiedAt),
			Birthday:                   timePointerFromNullTime(result.Birthday),
			ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
			AccountStatusExplanation:   result.UserAccountStatusExplanation,
			TwoFactorSecret:            result.TwoFactorSecret,
			HashedPassword:             result.HashedPassword,
			ID:                         result.ID,
			AccountStatus:              result.UserAccountStatus,
			Username:                   result.Username,
			FirstName:                  result.FirstName,
			LastName:                   result.LastName,
			EmailAddress:               result.EmailAddress,
			EmailAddressVerifiedAt:     timePointerFromNullTime(result.EmailAddressVerifiedAt),
			AvatarSrc:                  stringPointerFromNullString(result.AvatarSrc),
			ServiceRole:                result.ServiceRole,
			RequiresPasswordChange:     result.RequiresPasswordChange,
		}

		x.Data = append(x.Data, u)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetUserIDsThatNeedSearchIndexing fetches a list of valid vessels from the database that meet a particular filter.
func (q *Querier) GetUserIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetUserIDsNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing users list retrieval query")
	}

	return results, nil
}

// MarkUserAsIndexed updates a particular user's last_indexed_at value.
func (q *Querier) MarkUserAsIndexed(ctx context.Context, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	if _, err := q.generatedQuerier.UpdateUserLastIndexedAt(ctx, q.db, userID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking user as indexed")
	}

	logger.Info("user marked as indexed")

	return nil
}

// CreateUser creates a user. TODO: this should return a household as well.
func (q *Querier) CreateUser(ctx context.Context, input *types.UserDatabaseCreationInput) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachUsernameToSpan(span, input.Username)
	logger := q.logger.WithValues(map[string]any{
		keys.UsernameKey:                 input.Username,
		keys.UserEmailAddressKey:         input.EmailAddress,
		keys.HouseholdInvitationTokenKey: input.InvitationToken,
		"destination_household":          input.DestinationHouseholdID,
	})

	// begin user creation transaction
	tx, beginTransactionErr := q.db.BeginTx(ctx, nil)
	if beginTransactionErr != nil {
		return nil, observability.PrepareError(beginTransactionErr, span, "beginning transaction")
	}

	token, err := q.secretGenerator.GenerateBase64EncodedString(ctx, 32)
	if err != nil {
		return nil, observability.PrepareError(err, span, "generating email verification token")
	}

	if writeErr := q.generatedQuerier.CreateUser(ctx, tx, &generated.CreateUserParams{
		ID:                            input.ID,
		FirstName:                     input.FirstName,
		LastName:                      input.LastName,
		Username:                      input.Username,
		EmailAddress:                  input.EmailAddress,
		HashedPassword:                input.HashedPassword,
		TwoFactorSecret:               input.TwoFactorSecret,
		AvatarSrc:                     nullStringFromStringPointer(input.AvatarSrc),
		UserAccountStatus:             string(types.UnverifiedHouseholdStatus),
		Birthday:                      nullTimeFromTimePointer(input.Birthday),
		ServiceRole:                   authorization.ServiceUserRole.String(),
		EmailAddressVerificationToken: nullStringFromString(token),
	}); writeErr != nil {
		q.rollbackTransaction(ctx, tx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresDuplicateEntryErrorCode {
				return nil, database.ErrUserAlreadyExists
			}
		}

		return nil, observability.PrepareError(writeErr, span, "creating user")
	}

	hasValidInvite := input.InvitationToken != "" && input.DestinationHouseholdID != ""

	user := &types.User{
		ID:              input.ID,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Username:        input.Username,
		EmailAddress:    input.EmailAddress,
		HashedPassword:  input.HashedPassword,
		TwoFactorSecret: input.TwoFactorSecret,
		AccountStatus:   string(types.UnverifiedHouseholdStatus),
		Birthday:        input.Birthday,
		ServiceRole:     authorization.ServiceUserRole.String(),
		CreatedAt:       q.currentTime(),
	}
	logger = logger.WithValue(keys.UserIDKey, user.ID)
	tracing.AttachUserIDToSpan(span, user.ID)

	if strings.TrimSpace(input.HouseholdName) == "" {
		input.HouseholdName = fmt.Sprintf("%s's cool household", input.Username)
	}

	if err = q.createHouseholdForUser(ctx, tx, hasValidInvite, input.HouseholdName, user.ID); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating household for new user")
	}

	logger.Debug("household created")

	if hasValidInvite {
		if err = q.acceptInvitationForUser(ctx, tx, input); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "accepting household invitation")
		}
		logger.Debug("accepted invitation and joined household for user")
	}

	if err = q.attachInvitationsToUser(ctx, tx, user.EmailAddress, user.ID); err != nil {
		q.rollbackTransaction(ctx, tx)
		logger = logger.WithValue("email_address", user.EmailAddress).WithValue("user_id", user.ID)
		return nil, observability.PrepareAndLogError(err, logger, span, "attaching existing invitations to new user")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Debug("user and household created")

	return user, nil
}

func (q *Querier) createHouseholdForUser(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, hasValidInvite bool, householdName, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	// standard registration: we need to create the household
	householdID := identifiers.New()
	tracing.AttachHouseholdIDToSpan(span, householdID)

	hn := householdName
	if householdName == "" {
		hn = fmt.Sprintf("%s_default", userID)
	}

	householdCreationInput := &types.HouseholdDatabaseCreationInput{
		ID:            householdID,
		Name:          hn,
		BelongsToUser: userID,
	}

	// create the household.
	if err := q.generatedQuerier.CreateHousehold(ctx, querier, &generated.CreateHouseholdParams{
		City:          householdCreationInput.City,
		Name:          householdCreationInput.Name,
		BillingStatus: string(types.UnpaidHouseholdBillingStatus),
		ContactPhone:  householdCreationInput.ContactPhone,
		AddressLine1:  householdCreationInput.AddressLine1,
		AddressLine2:  householdCreationInput.AddressLine2,
		ID:            householdCreationInput.ID,
		State:         householdCreationInput.State,
		ZipCode:       householdCreationInput.ZipCode,
		Country:       householdCreationInput.Country,
		BelongsToUser: householdCreationInput.BelongsToUser,
		Latitude:      nullStringFromFloat64Pointer(householdCreationInput.Latitude),
		Longitude:     nullStringFromFloat64Pointer(householdCreationInput.Longitude),
	}); err != nil {
		q.rollbackTransaction(ctx, querier)
		return observability.PrepareError(err, span, "creating household")
	}

	if err := q.generatedQuerier.CreateHouseholdUserMembershipForNewUser(ctx, querier, &generated.CreateHouseholdUserMembershipForNewUserParams{
		ID:                 identifiers.New(),
		BelongsToUser:      userID,
		BelongsToHousehold: householdID,
		HouseholdRole:      authorization.HouseholdAdminRole.String(),
		DefaultHousehold:   !hasValidInvite,
	}); err != nil {
		q.rollbackTransaction(ctx, querier)
		return observability.PrepareError(err, span, "writing household user membership")
	}

	return nil
}

// UpdateUserUsername updates a user's username.
func (q *Querier) UpdateUserUsername(ctx context.Context, userID, newUsername string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	if newUsername == "" {
		return ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.UsernameKey, newUsername)
	tracing.AttachUsernameToSpan(span, newUsername)

	if _, err := q.generatedQuerier.UpdateUserUsername(ctx, q.db, &generated.UpdateUserUsernameParams{
		Username: newUsername,
		ID:       userID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating username")
	}

	logger.Info("username updated")

	return nil
}

// UpdateUserEmailAddress updates a user's username.
func (q *Querier) UpdateUserEmailAddress(ctx context.Context, userID, newEmailAddress string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(keys.UserEmailAddressKey, newEmailAddress).WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	if newEmailAddress == "" {
		return ErrEmptyInputProvided
	}
	tracing.AttachEmailAddressToSpan(span, newEmailAddress)

	if _, err := q.generatedQuerier.UpdateUserEmailAddress(ctx, q.db, &generated.UpdateUserEmailAddressParams{
		EmailAddress: newEmailAddress,
		ID:           userID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user email address")
	}

	logger.Info("user email address updated")

	return nil
}

// UpdateUserDetails updates a user's username.
func (q *Querier) UpdateUserDetails(ctx context.Context, userID string, input *types.UserDetailsDatabaseUpdateInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return ErrEmptyInputProvided
	}

	if userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	if _, err := q.generatedQuerier.UpdateUserDetails(ctx, q.db, &generated.UpdateUserDetailsParams{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Birthday:  nullTimeFromTime(input.Birthday),
		ID:        userID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user details")
	}

	logger.Info("user details updated")

	return nil
}

// UpdateUserAvatar updates a user's avatar source.
func (q *Querier) UpdateUserAvatar(ctx context.Context, userID, newAvatarSrc string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if newAvatarSrc == "" {
		return ErrEmptyInputProvided
	}

	if userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	if _, err := q.generatedQuerier.UpdateUserAvatarSrc(ctx, q.db, &generated.UpdateUserAvatarSrcParams{
		AvatarSrc: nullStringFromString(newAvatarSrc),
		ID:        userID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user avatar")
	}

	logger.Info("user avatar updated")

	return nil
}

// UpdateUserPassword updates a user's passwords hash in the database.
func (q *Querier) UpdateUserPassword(ctx context.Context, userID, newHash string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if newHash == "" {
		return ErrEmptyInputProvided
	}

	if userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	if _, err := q.generatedQuerier.UpdateUserPassword(ctx, q.db, &generated.UpdateUserPasswordParams{
		HashedPassword: newHash,
		ID:             userID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user password")
	}

	logger.Info("user password updated")

	return nil
}

// UpdateUserTwoFactorSecret marks a user's two factor secret as validated.
func (q *Querier) UpdateUserTwoFactorSecret(ctx context.Context, userID, newSecret string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if newSecret == "" {
		return ErrEmptyInputProvided
	}

	if userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	if _, err := q.generatedQuerier.UpdateUserTwoFactorSecret(ctx, q.db, &generated.UpdateUserTwoFactorSecretParams{
		TwoFactorSecret: newSecret,
		ID:              userID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user 2FA secret")
	}
	logger.Info("user two factor secret updated")

	return nil
}

// MarkUserTwoFactorSecretAsVerified marks a user's two factor secret as validated.
func (q *Querier) MarkUserTwoFactorSecretAsVerified(ctx context.Context, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	if err := q.generatedQuerier.MarkTwoFactorSecretAsVerified(ctx, q.db, userID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "writing verified two factor status to database")
	}

	logger.Info("user two factor secret verified")

	return nil
}

// MarkUserTwoFactorSecretAsUnverified marks a user's two factor secret as unverified.
func (q *Querier) MarkUserTwoFactorSecretAsUnverified(ctx context.Context, userID, newSecret string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if newSecret == "" {
		return ErrEmptyInputProvided
	}

	if userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	if err := q.generatedQuerier.MarkTwoFactorSecretAsUnverified(ctx, q.db, &generated.MarkTwoFactorSecretAsUnverifiedParams{
		TwoFactorSecret: newSecret,
		ID:              userID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "writing verified two factor status to database")
	}

	logger.Info("user two factor secret unverified")

	return nil
}

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

	if _, err = q.generatedQuerier.ArchiveUser(ctx, tx, userID); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving user")
	}

	if _, err = q.generatedQuerier.ArchiveUserMemberships(ctx, tx, userID); err != nil {
		q.rollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving user household memberships")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("user archived")

	return nil
}

func (q *Querier) GetEmailAddressVerificationTokenForUser(ctx context.Context, userID string) (string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return "", ErrInvalidIDProvided
	}
	tracing.AttachStringToSpan(span, keys.UserIDKey, userID)

	result, err := q.generatedQuerier.GetEmailVerificationTokenByUserID(ctx, q.db, userID)
	if err != nil {
		return "", observability.PrepareError(err, span, "getting user by email address verification token")
	}

	return result.String, nil
}

func (q *Querier) GetUserByEmailAddressVerificationToken(ctx context.Context, token string) (*types.User, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if token == "" {
		return nil, ErrEmptyInputProvided
	}

	result, err := q.generatedQuerier.GetUserByEmailAddressVerificationToken(ctx, q.db, nullStringFromString(token))
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting user by email address verification token")
	}

	u := &types.User{
		CreatedAt:                  result.CreatedAt,
		PasswordLastChangedAt:      timePointerFromNullTime(result.PasswordLastChangedAt),
		LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
		LastAcceptedTermsOfService: timePointerFromNullTime(result.LastAcceptedTermsOfService),
		LastAcceptedPrivacyPolicy:  timePointerFromNullTime(result.LastAcceptedPrivacyPolicy),
		TwoFactorSecretVerifiedAt:  timePointerFromNullTime(result.TwoFactorSecretVerifiedAt),
		AvatarSrc:                  stringPointerFromNullString(result.AvatarSrc),
		Birthday:                   timePointerFromNullTime(result.Birthday),
		ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
		AccountStatusExplanation:   result.UserAccountStatusExplanation,
		TwoFactorSecret:            result.TwoFactorSecret,
		HashedPassword:             result.HashedPassword,
		ID:                         result.ID,
		AccountStatus:              result.UserAccountStatus,
		Username:                   result.Username,
		FirstName:                  result.FirstName,
		LastName:                   result.LastName,
		EmailAddress:               result.EmailAddress,
		EmailAddressVerifiedAt:     timePointerFromNullTime(result.EmailAddressVerifiedAt),
		ServiceRole:                result.ServiceRole,
		RequiresPasswordChange:     result.RequiresPasswordChange,
	}

	return u, nil
}

func (q *Querier) MarkUserEmailAddressAsVerified(ctx context.Context, userID, token string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)

	if token == "" {
		return ErrEmptyInputProvided
	}

	if err := q.generatedQuerier.MarkEmailAddressAsVerified(ctx, q.db, &generated.MarkEmailAddressAsVerifiedParams{
		ID:                            userID,
		EmailAddressVerificationToken: nullStringFromString(token),
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}

		return observability.PrepareAndLogError(err, logger, span, "writing verified email address status to database")
	}

	if err := q.UpdateUserAccountStatus(ctx, userID, &types.UserAccountStatusUpdateInput{
		NewStatus:    string(types.GoodStandingUserAccountStatus),
		Reason:       "verified email address",
		TargetUserID: userID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user account status")
	}

	return nil
}
