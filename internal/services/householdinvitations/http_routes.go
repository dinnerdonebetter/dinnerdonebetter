package householdinvitations

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// HouseholdInvitationIDURIParamKey is a standard string that we'll use to refer to webhook IDs with.
	HouseholdInvitationIDURIParamKey = "householdInvitationID"
)

// CreateHandler is our webhook creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	providedInput := new(types.HouseholdInvitationCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	input := types.HouseholdInvitationDatabaseCreationInputFromHouseholdInvitationCreationInput(providedInput)
	input.ID = ksuid.New().String()
	tracing.AttachHouseholdInvitationIDToSpan(span, input.ID)

	preWrite := &types.PreWriteMessage{
		DataType:                  types.HouseholdInvitationDataType,
		HouseholdInvitation:       input,
		AttributableToUserID:      sessionCtxData.Requester.UserID,
		AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
	}
	if err = s.preWritesPublisher.Publish(ctx, preWrite); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing household invitation write message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	pwr := types.PreWriteResponse{ID: input.ID}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, pwr, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns an webhook.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine relevant webhook ID.
	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	tracing.AttachHouseholdIDToSpan(span, sessionCtxData.ActiveHouseholdID)
	logger = logger.WithValue(keys.HouseholdIDKey, sessionCtxData.ActiveHouseholdID)

	// fetch the webhook from the database.
	householdInvitation, err := s.householdInvitationDataManager.GetHouseholdInvitation(ctx, householdInvitationID, sessionCtxData.ActiveHouseholdID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("No rows found in household invitation database")
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching household invitation from database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode the response.
	s.encoderDecoder.RespondWithData(ctx, res, householdInvitation)
}

// ArchiveHandler returns a handler that archives an webhook.
func (s *service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := sessionCtxData.ActiveHouseholdID
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	// determine relevant webhook ID.
	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	exists, householdInvitationExistenceCheckErr := s.householdInvitationDataManager.HouseholdInvitationExists(ctx, householdID, householdInvitationID)
	if householdInvitationExistenceCheckErr != nil && !errors.Is(householdInvitationExistenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		observability.AcknowledgeError(householdInvitationExistenceCheckErr, logger, span, "checking item existence")
		return
	} else if !exists || errors.Is(householdInvitationExistenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	pam := &types.PreArchiveMessage{
		DataType:                  types.HouseholdInvitationDataType,
		HouseholdInvitationID:     householdInvitationID,
		AttributableToUserID:      sessionCtxData.Requester.UserID,
		AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
	}
	if err = s.preArchivesPublisher.Publish(ctx, pam); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing household invitation archive message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// let everybody go home.
	res.WriteHeader(http.StatusNoContent)
}

func (s *service) InboundInvitesHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := sessionCtxData.ActiveHouseholdID
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
}

func (s *service) OutboundInvitesHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := sessionCtxData.ActiveHouseholdID
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
}

func (s *service) CancelInviteHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := sessionCtxData.ActiveHouseholdID
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
}

func (s *service) AcceptInviteHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := sessionCtxData.ActiveHouseholdID
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
}

func (s *service) RejectInviteHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := sessionCtxData.ActiveHouseholdID
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
}

func (s *service) LeaveHouseholdHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := sessionCtxData.ActiveHouseholdID
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
}
