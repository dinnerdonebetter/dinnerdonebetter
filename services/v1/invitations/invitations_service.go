package invitations

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// CreateMiddlewareCtxKey is a string alias we can use for referring to invitation input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "invitation_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to invitation update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "invitation_update_input"

	counterName        metrics.CounterName = "invitations"
	counterDescription string              = "the number of invitations managed by the invitations service"
	topicName          string              = "invitations"
	serviceName        string              = "invitations_service"
)

var (
	_ models.InvitationDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list invitations
	Service struct {
		logger                logging.Logger
		invitationDataManager models.InvitationDataManager
		invitationIDFetcher   InvitationIDFetcher
		userIDFetcher         UserIDFetcher
		invitationCounter     metrics.UnitCounter
		encoderDecoder        encoding.EncoderDecoder
		reporter              newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// InvitationIDFetcher is a function that fetches invitation IDs.
	InvitationIDFetcher func(*http.Request) uint64
)

// ProvideInvitationsService builds a new InvitationsService.
func ProvideInvitationsService(
	logger logging.Logger,
	invitationDataManager models.InvitationDataManager,
	invitationIDFetcher InvitationIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	invitationCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	invitationCounter, err := invitationCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                logger.WithName(serviceName),
		invitationIDFetcher:   invitationIDFetcher,
		userIDFetcher:         userIDFetcher,
		invitationDataManager: invitationDataManager,
		encoderDecoder:        encoder,
		invitationCounter:     invitationCounter,
		reporter:              reporter,
	}

	return svc, nil
}
