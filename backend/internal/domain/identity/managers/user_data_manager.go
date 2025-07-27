package managers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/services/identity/indexing"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const (
	o11yName = "identity_data_manager"

	totpIssuer             = "DinnerDoneBetter"
	base64ImagePrefix      = "data:image/jpeg;base64,"
	minimumPasswordEntropy = 60
	totpSecretSize         = 64
)

var (
	userAvatarBase64Encoding = base64.RawURLEncoding

	// ErrInvalidIDProvided indicates a required ID was passed in empty.
	ErrInvalidIDProvided = errors.New("required ID was empty")

	// ErrNilInputProvided indicates that a required parameter was nil
	ErrNilInputProvided = errors.New("nil input provided")

	// ErrEmptyInputProvided indicates that required input was empty
	ErrEmptyInputProvided = errors.New("empty input provided")
)

type (
	manager struct {
		tracer               tracing.Tracer
		logger               logging.Logger
		dataChangesPublisher messagequeue.Publisher
		identityRepo         identity.Repository
		secretGenerator      random.Generator
		authenticator        authentication.Hasher
		userSearchIndex      indexing.UserTextSearcher
	}

	Config struct {
		DataChangesTopicName string `json:"dataChangesTopicName" env:"DATA_CHANGES_TOPIC_NAME"`
	}
)

func NewIdentityDataManager(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	publisherProvider messagequeue.PublisherProvider,
	identityRepo identity.Repository,
	secretGenerator random.Generator,
	authenticator authentication.Hasher,
	userSearchIndex indexing.UserTextSearcher,
	cfg *Config,
) (IdentityDataManager, error) {
	publisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up data changes publisher: %w", err)
	}

	return &manager{
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:               logging.EnsureLogger(logger).WithName(o11yName),
		identityRepo:         identityRepo,
		dataChangesPublisher: publisher,
		secretGenerator:      secretGenerator,
		authenticator:        authenticator,
		userSearchIndex:      userSearchIndex,
	}, nil
}

func (i *manager) AcceptAccountInvitation(ctx context.Context, accountInvitationID string, input *identity.AccountInvitationUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if accountInvitationID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "invalid input attached to request")
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInvitationIDKey: accountInvitationID,
	}, span, i.logger)

	if err := i.identityRepo.AcceptAccountInvitation(ctx, accountInvitationID, input.Token, input.Note); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "accepting account invitation")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountInvitationCanceledServiceEventType, map[string]any{
		keys.AccountInvitationIDKey: accountInvitationID,
	}))

	return nil
}

func (i *manager) RejectAccountInvitation(ctx context.Context, accountInvitationID string, input *identity.AccountInvitationUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if accountInvitationID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInvitationIDKey: accountInvitationID,
	}, span, i.logger)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "invalid input attached to request")
	}

	// note, this is where you would call input.ValidateWithContext, if that currently had any effect.

	invitation, err := i.identityRepo.GetAccountInvitationByTokenAndID(ctx, input.Token, accountInvitationID)
	if errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareError(err, span, "account invitation not found")
	} else if err != nil {
		return observability.PrepareError(err, span, "retrieving invitation")
	}

	if err = i.identityRepo.RejectAccountInvitation(ctx, invitation.ID, input.Note); err != nil {
		return observability.PrepareError(err, span, "rejecting invitation")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountInvitationRejectedServiceEventType, map[string]any{
		keys.AccountInvitationIDKey: accountInvitationID,
	}))

	return nil
}

func (i *manager) CancelAccountInvitation(ctx context.Context, accountInvitationID, note string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if accountInvitationID == "" {
		return ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInvitationIDKey: accountInvitationID,
	}, span, i.logger)

	if err := i.identityRepo.CancelAccountInvitation(ctx, accountInvitationID, note); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "canceling account invitation")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountInvitationCanceledServiceEventType, map[string]any{
		keys.AccountInvitationIDKey: accountInvitationID,
	}))

	return nil
}

func (i *manager) ArchiveAccount(ctx context.Context, accountID, ownerID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if accountID == "" || ownerID == "" {
		return ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: accountID,
		keys.UserIDKey:    ownerID,
	}, span, i.logger)

	if err := i.identityRepo.ArchiveAccount(ctx, accountID, ownerID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "ArchiveAccount")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountArchivedServiceEventType, map[string]any{
		keys.AccountIDKey: accountID,
		keys.UserIDKey:    ownerID,
	}))

	return nil
}

func (i *manager) ArchiveUserMembership(ctx context.Context, userID, accountID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if accountID == "" || userID == "" {
		return ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: accountID,
		keys.UserIDKey:    userID,
	}, span, i.logger)

	if err := i.identityRepo.RemoveUserFromAccount(ctx, userID, accountID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "RemoveUserFromAccount")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountMemberRemovedServiceEventType, map[string]any{
		keys.AccountIDKey: accountID,
		keys.UserIDKey:    userID,
	}))

	return nil
}

func (i *manager) ArchiveUser(ctx context.Context, userID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	if err := i.identityRepo.ArchiveUser(ctx, userID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "ArchiveUser")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.UserArchivedServiceEventType, map[string]any{
		keys.UserIDKey: userID,
	}))

	return nil
}

func (i *manager) CreateAccount(ctx context.Context, input *identity.AccountCreationRequestInput) (*identity.Account, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "invalid input attached to request")
	}

	logger := i.logger.WithSpan(span)

	created, err := i.identityRepo.CreateAccount(ctx, converters.ConvertAccountCreationInputToAccountDatabaseCreationInput(input))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating Account")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountCreatedServiceEventType, map[string]any{
		keys.AccountIDKey: created.ID,
	}))

	return created, nil
}

func (i *manager) CreateAccountInvitation(ctx context.Context, userID string, input *identity.AccountInvitationCreationRequestInput) (*identity.AccountInvitation, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "invalid input attached to request")
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	created, err := i.identityRepo.CreateAccountInvitation(ctx, converters.ConvertAccountInvitationCreationInputToAccountInvitationDatabaseCreationInput(input))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating AccountInvitation")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountInvitationCreatedServiceEventType, map[string]any{
		keys.AccountInvitationIDKey: created.ID,
		keys.UserIDKey:              userID,
	}))

	return created, nil
}

func (i *manager) CreateUser(ctx context.Context, input *identity.UserRegistrationInput) (*identity.UserCreationResponse, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UsernameKey: input.Username,
	}, span, i.logger)

	input.Username = strings.TrimSpace(input.Username)
	tracing.AttachToSpan(span, keys.UsernameKey, input.Username)
	input.EmailAddress = strings.TrimSpace(strings.ToLower(input.EmailAddress))
	tracing.AttachToSpan(span, keys.UserEmailAddressKey, input.EmailAddress)
	input.Password = strings.TrimSpace(input.Password)

	logger = logger.WithValues(map[string]any{
		keys.UsernameKey:               input.Username,
		keys.UserEmailAddressKey:       input.EmailAddress,
		keys.AccountInvitationIDKey:    input.InvitationID,
		keys.AccountInvitationTokenKey: input.InvitationToken,
	})

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided dbInput was invalid")
		return nil, observability.PrepareError(err, span, "invalid user creation dbInput provided")
	}

	// ensure the password is not garbage-tier
	if err := passwordvalidator.Validate(strings.TrimSpace(input.Password), minimumPasswordEntropy); err != nil {
		logger.WithValue("password_validation_error", err).Debug("weak password provided to user creation route")
		return nil, observability.PrepareAndLogError(err, logger, span, "weak password provided for user creation")
	}

	var invitation *identity.AccountInvitation
	if input.InvitationID != "" && input.InvitationToken != "" {
		invite, err := i.identityRepo.GetAccountInvitationByTokenAndID(ctx, input.InvitationToken, input.InvitationID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, observability.PrepareAndLogError(err, logger, span, "no account invitation found")
		} else if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "getting account invitation by token and ID")
		}

		invitation = invite
		logger.Debug("retrieved account invitation")
	}

	logger.Debug("completed invitation check")

	// hash the password
	hp, err := i.authenticator.HashPassword(ctx, strings.TrimSpace(input.Password))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "hashing user creation password")
	}

	// generate a two factor secret.
	tfs, err := i.secretGenerator.GenerateBase32EncodedString(ctx, totpSecretSize)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "generating two factor secret")
	}

	dbInput := &identity.UserDatabaseCreationInput{
		ID:              identifiers.New(),
		Username:        input.Username,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		EmailAddress:    input.EmailAddress,
		HashedPassword:  hp,
		TwoFactorSecret: tfs,
		InvitationToken: input.InvitationToken,
		Birthday:        input.Birthday,
		AccountName:     input.AccountName,
	}

	if invitation != nil {
		logger.Debug("supplementing user creation dbInput with invitation data")
		dbInput.DestinationAccountID = invitation.DestinationAccount.ID
		dbInput.InvitationToken = invitation.Token
	}

	// create the user.
	user, err := i.identityRepo.CreateUser(ctx, dbInput)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating user")
		if errors.Is(err, database.ErrUserAlreadyExists) {
			return nil, observability.PrepareAndLogError(err, logger, span, "user already exists")
		}
		return nil, observability.PrepareAndLogError(err, logger, span, "creating user in database")
	}

	logger.Debug("user created")

	defaultAccountID, err := i.identityRepo.GetDefaultAccountIDForUser(ctx, user.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching default account ID for user")
	}

	emailVerificationToken, err := i.identityRepo.GetEmailAddressVerificationTokenForUser(ctx, user.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching email verification token for user")
	}

	// notify the relevant parties.
	tracing.AttachToSpan(span, keys.UserIDKey, user.ID)

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.UserSignedUpServiceEventType, map[string]any{
		keys.AccountIDKey:                  defaultAccountID,
		keys.UserIDKey:                     user.ID,
		keys.UserEmailVerificationTokenKey: emailVerificationToken,
	}))

	/* TODO: this should be done in a downstream handler

	if err = i.analyticsReporter.AddUser(ctx, user.ID, map[string]any{
		"username":        user.Username,
		"default_account": defaultAccountID,
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
	}); err != nil {
		observability.AcknowledgeError(err, logger, span, "identifying user for analytics")
	}

	if err = s.featureFlagManager.Identify(ctx, user); err != nil {
		observability.AcknowledgeError(err, logger, span, "identifying user in feature flag manager")
	}
	*/

	// UserCreationResponse is a struct we can use to notify the user of their two factor secret, but ideally just this once and then never again.
	ucr := &identity.UserCreationResponse{
		CreatedUserID:   user.ID,
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		EmailAddress:    user.EmailAddress,
		CreatedAt:       user.CreatedAt,
		TwoFactorSecret: user.TwoFactorSecret,
		Birthday:        user.Birthday,
		TwoFactorQRCode: i.buildQRCode(ctx, user.Username, user.TwoFactorSecret),
	}

	return ucr, nil
}

func (i *manager) GetAccount(ctx context.Context, accountID string) (*identity.Account, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if accountID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: accountID,
	}, span, i.logger)

	account, err := i.identityRepo.GetAccount(ctx, accountID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching account")
	}

	return account, nil
}

func (i *manager) GetAccountInvitation(ctx context.Context, accountID, accountInvitationID string) (*identity.AccountInvitation, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if accountID == "" || accountInvitationID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey:           accountID,
		keys.AccountInvitationIDKey: accountInvitationID,
	}, span, i.logger)

	invitation, err := i.identityRepo.GetAccountInvitationByAccountAndID(ctx, accountID, accountInvitationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting invitation")
	}

	return invitation, nil
}

func (i *manager) GetAccounts(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.Account, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, "", ErrInvalidIDProvided
	}

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	accounts, err := i.identityRepo.GetAccounts(ctx, userID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching accounts")
	}

	return accounts.Data, "TODO", nil
}

func (i *manager) GetReceivedAccountInvitations(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.AccountInvitation, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, "", ErrInvalidIDProvided
	}

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	invites, err := i.identityRepo.GetPendingAccountInvitationsForUser(ctx, userID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching invites")
	}

	return invites.Data, "TODO", nil
}

func (i *manager) GetSentAccountInvitations(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.AccountInvitation, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, "", ErrInvalidIDProvided
	}

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	invites, err := i.identityRepo.GetPendingAccountInvitationsFromUser(ctx, userID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching invites")
	}

	return invites.Data, "TODO", nil
}

func (i *manager) GetUser(ctx context.Context, userID string) (*identity.User, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	user, err := i.identityRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting user")
	}

	return user, nil
}

func (i *manager) GetUsers(ctx context.Context, filter *filtering.QueryFilter) ([]*identity.User, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := i.logger.WithSpan(span)

	users, err := i.identityRepo.GetUsers(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching users")
	}

	return users.Data, "", nil
}

func (i *manager) SearchForUsers(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*identity.User, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if query == "" {
		return nil, "", errors.New("query cannot be empty")
	}

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UseDatabaseKey: useDatabase,
	}, span, i.logger)

	if useDatabase {
		users, err := i.identityRepo.SearchForUsersByUsername(ctx, query, filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, "", observability.PrepareError(err, span, "no users found")
			}
			return nil, "", observability.PrepareAndLogError(err, logger, span, "searching for users")
		}

		return users, "TODO", nil
	} else {
		uss, err := i.userSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, "", observability.PrepareAndLogError(err, logger, span, "searching for users")
		}

		userIDs := []string{}
		for _, us := range uss {
			userIDs = append(userIDs, us.ID)
		}

		users, err := i.identityRepo.GetUsersWithIDs(ctx, userIDs)
		if err != nil {
			return nil, "", observability.PrepareAndLogError(err, logger, span, "searching for users")
		}

		return users, "TODO", nil
	}
}

func (i *manager) SetDefaultAccount(ctx context.Context, userID, accountID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" || accountID == "" {
		return ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey:    userID,
		keys.AccountIDKey: accountID,
	}, span, i.logger)

	// mark household as default in database.
	if err := i.identityRepo.MarkAccountAsUserDefault(ctx, userID, accountID); err != nil {
		return observability.PrepareError(err, span, "marking default account as user")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountSetAsDefaultServiceEventType, map[string]any{
		keys.AccountIDKey: accountID,
		keys.UserIDKey:    userID,
	}))

	return nil
}

func (i *manager) TransferAccountOwnership(ctx context.Context, accountID string, input *identity.AccountOwnershipTransferInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if accountID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: accountID,
	}, span, i.logger)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "")
	}

	// transfer ownership of household in database.
	if err := i.identityRepo.TransferAccountOwnership(ctx, accountID, input); err != nil {
		return observability.PrepareError(err, span, "")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountOwnershipTransferredServiceEventType, map[string]any{
		keys.AccountIDKey: accountID,
	}))

	return nil
}

func (i *manager) UpdateAccount(ctx context.Context, accountID string, input *identity.AccountUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if accountID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: accountID,
	}, span, i.logger)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "validating account update")
	}

	// determine account ID.
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	// fetch account from database.
	account, err := i.identityRepo.GetAccount(ctx, accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareError(err, span, "no account found")
	} else if err != nil {
		return observability.PrepareError(err, span, "fetching account")
	}

	// update the data structure.
	account.Update(input)

	// update account in database.
	if err = i.identityRepo.UpdateAccount(ctx, account); err != nil {
		return observability.PrepareError(err, span, "updating account")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountUpdatedServiceEventType, map[string]any{
		keys.AccountIDKey: accountID,
	}))

	return nil
}

func (i *manager) UpdateAccountMemberPermissions(ctx context.Context, userID, accountID string, input *identity.ModifyUserPermissionsInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" || accountID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "invalid input attached to request")
	}

	// create account in database.
	if err := i.identityRepo.ModifyUserPermissions(ctx, accountID, userID, input); err != nil {
		return observability.PrepareError(err, span, "modifying user permissions")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountMembershipPermissionsUpdatedServiceEventType, map[string]any{
		keys.AccountIDKey: accountID,
	}))

	return nil
}

func (i *manager) UpdateUserDetails(ctx context.Context, userID string, input *identity.UserDetailsUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "invalid input attached to request")
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	if err := i.identityRepo.UpdateUserDetails(ctx, userID, converters.ConvertUserDetailsUpdateRequestInputToUserDetailsUpdateInput(input)); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user details")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.UserDetailsChangedEventType, map[string]any{
		keys.UserIDKey: userID,
	}))

	return nil
}

func (i *manager) UpdateUserEmailAddress(ctx context.Context, userID, newEmail string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	if err := i.identityRepo.UpdateUserEmailAddress(ctx, userID, newEmail); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user email address")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.EmailAddressChangedEventType, map[string]any{
		keys.UserIDKey: userID,
	}))

	return nil
}

func (i *manager) UpdateUserUsername(ctx context.Context, userID, newUsername string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	if newUsername == "" {
		return ErrEmptyInputProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	if err := i.identityRepo.UpdateUserUsername(ctx, userID, newUsername); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user username")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.UsernameChangedEventType, map[string]any{
		keys.UserIDKey: userID,
	}))

	return nil
}

func (i *manager) UploadUserAvatar(ctx context.Context, userID, base64EncodedImageData string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, i.logger)

	data, err := userAvatarBase64Encoding.DecodeString(base64EncodedImageData)
	if err != nil {
		return observability.PrepareError(err, span, "decoding base64 encoded image data")
	}

	logger = observability.ObserveValues(map[string]any{
		"data.length": len(data),
	}, span, logger)

	if err = i.identityRepo.UpdateUserAvatar(ctx, userID, base64EncodedImageData); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user avatar")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.UserAvatarChangedEventType, map[string]any{
		keys.UserIDKey: userID,
	}))

	return nil
}

func (i *manager) AdminUpdateUserStatus(ctx context.Context, input *identity.UserAccountStatusUpdateInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return ErrNilInputProvided
	}

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: input.TargetUserID,
		keys.ReasonKey: input.Reason,
	}, span, i.logger)

	if err := i.identityRepo.UpdateUserAccountStatus(ctx, input.TargetUserID, input); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user account status")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.UserStatusChangedServiceEventType, map[string]any{
		keys.UserIDKey: input.TargetUserID,
	}))

	return nil
}

// buildQRCode builds a QR code for a given username and secret.
func (i *manager) buildQRCode(ctx context.Context, username, twoFactorSecret string) string {
	_, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.UsernameKey: username,
	}, span, i.logger)

	// "otpauth://totp/{{ .Issuer }}:{{ .EnsureUsername }}?secret={{ .Secret }}&issuer={{ .Issuer }}",
	otpString := fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		totpIssuer,
		username,
		twoFactorSecret,
		totpIssuer,
	)

	// encode two factor secret as authenticator-friendly QR code
	qrCode, err := qr.Encode(otpString, qr.L, qr.Auto)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "encoding OTP string")
		return ""
	}

	// scale the QR code so that it's not a PNG for ants.
	qrCode, err = barcode.Scale(qrCode, 256, 256)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "scaling QR code")
		return ""
	}

	// encode the QR code to PNG.
	var b bytes.Buffer
	if err = png.Encode(&b, qrCode); err != nil {
		observability.AcknowledgeError(err, logger, span, "encoding QR code to PNG")
		return ""
	}

	// base64 encode the image for easy HTML use.
	return fmt.Sprintf("%s%s", base64ImagePrefix, base64.StdEncoding.EncodeToString(b.Bytes()))
}
