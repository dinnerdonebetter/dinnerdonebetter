package householdinvitations

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/prixfixeco/api_server/internal/email"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// HouseholdInvitationIDURIParamKey is a standard string that we'll use to refer to household invitation IDs with.
	HouseholdInvitationIDURIParamKey = "householdInvitationID"
)

var _ types.HouseholdInvitationDataService = (*service)(nil)

// InviteMemberHandler is our household creation route.
func (s *service) InviteMemberHandler(res http.ResponseWriter, req *http.Request) {
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

	requester := sessionCtxData.Requester.UserID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = logger.WithValue(keys.RequesterIDKey, requester)

	householdID := s.householdIDFetcher(req)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdInvitationCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}
	providedInput.ID = ksuid.New().String()
	providedInput.DestinationHousehold = householdID
	providedInput.FromUser = requester

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	input := types.HouseholdInvitationDatabaseCreationInputFromHouseholdInvitationCreationInput(providedInput)

	token, err := s.secretGenerator.GenerateBase64EncodedString(ctx, 64)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating invitation token")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}
	input.Token = token

	userID, err := s.userDataManager.GetUserIDByEmail(ctx, input.ToEmail)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "fetching user ID by email")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if userID != "" {
		input.ToUser = &userID
	}

	householdInvitation, err := s.householdInvitationDataManager.CreateHouseholdInvitation(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating household invitation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.HouseholdInvitationDataType,
			EventType:                 types.HouseholdInvitationCreatedCustomerEventType,
			HouseholdInvitation:       householdInvitation,
			AttributableToUserID:      requester,
			HouseholdID:               householdID,
			AttributableToHouseholdID: householdID,
		}
		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	if s.emailer != nil {
		msg, emailGenerationErr := email.BuildInviteMemberEmail(householdInvitation)
		if emailGenerationErr != nil {
			observability.AcknowledgeError(emailGenerationErr, logger, span, "building email message")
		} else {
			if err = s.emailer.SendEmail(ctx, msg); err != nil {
				observability.AcknowledgeError(err, logger, span, "sending email notice")
			}
		}
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, householdInvitation, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a household invitation.
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

	// determine relevant household invitation ID.
	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	tracing.AttachHouseholdIDToSpan(span, sessionCtxData.ActiveHouseholdID)
	logger = logger.WithValue(keys.HouseholdIDKey, sessionCtxData.ActiveHouseholdID)

	// fetch the household invitation from the database.
	householdInvitation, err := s.householdInvitationDataManager.GetHouseholdInvitationByHouseholdAndID(ctx, sessionCtxData.ActiveHouseholdID, householdInvitationID)
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

func (s *service) InboundInvitesHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilter(req)
	filter.AttachToLogger(logger)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := s.householdIDFetcher(req)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	invitations, err := s.householdInvitationDataManager.GetPendingHouseholdInvitationsForUser(ctx, userID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching outbound invites")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	s.encoderDecoder.RespondWithData(ctx, res, invitations)
}

func (s *service) OutboundInvitesHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilter(req)
	filter.AttachToLogger(logger)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := s.householdIDFetcher(req)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	invitations, err := s.householdInvitationDataManager.GetPendingHouseholdInvitationsFromUser(ctx, userID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching outbound invites")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	s.encoderDecoder.RespondWithData(ctx, res, invitations)
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

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdInvitationUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	// note, this is where you would call providedInput.ValidateWithContext, if that currently had any effect.

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	invitation, err := s.householdInvitationDataManager.GetHouseholdInvitationByTokenAndID(ctx, providedInput.Token, householdInvitationID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving invitation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if err = s.householdInvitationDataManager.AcceptHouseholdInvitation(ctx, invitation.ID, providedInput.Token, providedInput.Note); err != nil {
		observability.AcknowledgeError(err, logger, span, "accepting invitation")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "error accepting invitation", http.StatusInternalServerError)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.HouseholdInvitationDataType,
			EventType:                 types.HouseholdInvitationAcceptedCustomerEventType,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			AttributableToHouseholdID: invitation.DestinationHousehold,
			HouseholdID:               invitation.DestinationHousehold,
			HouseholdInvitationID:     householdInvitationID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	res.WriteHeader(http.StatusAccepted)
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

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdInvitationUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	// note, this is where you would call providedInput.ValidateWithContext, if that currently had any effect.

	userID := sessionCtxData.Requester.UserID
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	invitation, err := s.householdInvitationDataManager.GetHouseholdInvitationByTokenAndID(ctx, providedInput.Token, householdInvitationID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving invitation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if err = s.householdInvitationDataManager.CancelHouseholdInvitation(ctx, invitation.ID, providedInput.Token, providedInput.Note); err != nil {
		observability.AcknowledgeError(err, logger, span, "cancelling invitation")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "error cancelling invitation", http.StatusInternalServerError)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.HouseholdInvitationDataType,
			EventType:                 types.HouseholdInvitationCanceledCustomerEventType,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			HouseholdID:               invitation.DestinationHousehold,
			AttributableToHouseholdID: invitation.DestinationHousehold,
			HouseholdInvitationID:     householdInvitationID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	res.WriteHeader(http.StatusAccepted)
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

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdInvitationUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	// note, this is where you would call providedInput.ValidateWithContext, if that currently had any effect.

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachHouseholdInvitationIDToSpan(span, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	invitation, err := s.householdInvitationDataManager.GetHouseholdInvitationByTokenAndID(ctx, providedInput.Token, householdInvitationID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving invitation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if err = s.householdInvitationDataManager.RejectHouseholdInvitation(ctx, invitation.ID, providedInput.Token, providedInput.Note); err != nil {
		observability.AcknowledgeError(err, logger, span, "rejecting invitation")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "error rejecting invitation", http.StatusInternalServerError)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.HouseholdInvitationDataType,
			EventType:                 types.HouseholdInvitationRejectedCustomerEventType,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			HouseholdID:               invitation.DestinationHousehold,
			AttributableToHouseholdID: invitation.DestinationHousehold,
			HouseholdInvitationID:     householdInvitationID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	res.WriteHeader(http.StatusAccepted)
}
