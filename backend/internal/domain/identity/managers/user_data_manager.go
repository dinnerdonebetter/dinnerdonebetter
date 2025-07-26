package managers

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

/*

TODO:
- all string and pointered params are checked for emptiness and return errors
- all methods write a data change message
- all values are observed in both logs and analytics messages

*/

const (
	o11yName = "identity_data_manager"
)

type (
	IdentityDataManager struct {
		tracer               tracing.Tracer
		logger               logging.Logger
		dataChangesPublisher messagequeue.Publisher
		identityRepo         identity.Repository
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
	cfg *Config,
) (*IdentityDataManager, error) {
	publisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up data changes publisher: %w", err)
	}

	return &IdentityDataManager{
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:               logging.EnsureLogger(logger).WithName(o11yName),
		identityRepo:         identityRepo,
		dataChangesPublisher: publisher,
	}, nil
}

func (i *IdentityDataManager) AcceptAccountInvitation(ctx context.Context, accountInvitationID string, input *identity.AccountInvitationUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

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

func (i *IdentityDataManager) ArchiveAccount(ctx context.Context, accountID, ownerID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.ArchiveAccount(ctx, accountID, ownerID)

	return nil
}

func (i *IdentityDataManager) ArchiveUserMembership(ctx context.Context, userID, accountID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.RemoveUserFromAccount(ctx, userID, accountID)

	return nil
}

func (i *IdentityDataManager) ArchiveUser(ctx context.Context, userID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.ArchiveUser(ctx, userID)

	return nil
}

func (i *IdentityDataManager) CancelAccountInvitation(ctx context.Context, accountInvitationID, note string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.CancelAccountInvitation(ctx, accountInvitationID, note)

	return nil
}

func (i *IdentityDataManager) CreateAccount(ctx context.Context, input *identity.AccountCreationRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.CreateAccount(ctx, converters.ConvertAccountCreationInputToAccountDatabaseCreationInput(input))

	return nil
}

func (i *IdentityDataManager) CreateAccountInvitation(ctx context.Context, input *identity.AccountInvitationCreationRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.CreateAccountInvitation(ctx, converters.ConvertAccountInvitationCreationInputToAccountInvitationDatabaseCreationInput(input))

	return nil
}

func (i *IdentityDataManager) CreateUser(ctx context.Context, input *identity.UserRegistrationInput) (*identity.User, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	// TODO: port

	return nil, nil
}

func (i *IdentityDataManager) GetAccount(ctx context.Context, accountID string) (*identity.Account, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.GetAccount(ctx, accountID)

	return nil, nil
}

func (i *IdentityDataManager) GetAccountInvitation(ctx context.Context, accountID, accountInvitationID string) (*identity.AccountInvitation, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.GetAccountInvitationByAccountAndID(ctx, accountID, accountInvitationID)

	return nil, nil
}

func (i *IdentityDataManager) GetAccounts(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.Account, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.GetAccounts(ctx, userID, filter)

	return nil, "", nil
}

func (i *IdentityDataManager) GetReceivedAccountInvitations(ctx context.Context, userID string, filter *filtering.QueryFilter) (*identity.Account, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.GetPendingAccountInvitationsForUser(ctx, userID, filter)

	return nil, nil
}

func (i *IdentityDataManager) GetSentAccountInvitations(ctx context.Context, userID string, filter *filtering.QueryFilter) (*identity.Account, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.GetPendingAccountInvitationsFromUser(ctx, userID, filter)

	return nil, nil
}

func (i *IdentityDataManager) GetUser(ctx context.Context, userID string) (*identity.User, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.GetUser(ctx, userID)

	return nil, nil
}

func (i *IdentityDataManager) RejectAccountInvitation(ctx context.Context) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) GetUsers(ctx context.Context, filter *filtering.QueryFilter) ([]*identity.User, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	i.identityRepo.GetUsers(ctx, filter)

	return nil, "", nil
}

func (i *IdentityDataManager) SearchForUsers(ctx context.Context, query string, filter *filtering.QueryFilter) (*identity.Account, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil, nil
}

func (i *IdentityDataManager) SetDefaultAccount(ctx context.Context, accountID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) TransferAccountOwnership(ctx context.Context) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) UpdateAccount(ctx context.Context, input *identity.AccountUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) UpdateAccountMemberPermissions(ctx context.Context, input *identity.ModifyUserPermissionsInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) UpdateUserDetails(ctx context.Context, input *identity.UserDetailsUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) UpdateUserEmailAddress(ctx context.Context, newEmail string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) UpdateUserUsername(ctx context.Context, newUsername string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) UploadUserAvatar(ctx context.Context, newAvatar []byte) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) AdminUpdateUserStatus(ctx context.Context, input identity.UserAccountStatusUpdateInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}
