package householdinvitations

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
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

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	userID := sessionCtxData.Requester.UserID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = logger.WithValue(keys.RequesterIDKey, userID)

	householdID := s.householdIDFetcher(req)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdInvitationCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	providedInput.ToEmail = strings.TrimSpace(strings.ToLower(providedInput.ToEmail))

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if providedInput.ExpiresAt == nil {
		providedInput.ExpiresAt = pointers.Pointer(time.Now().Add((time.Hour * 24) * 7))
	}

	input := converters.ConvertHouseholdInvitationCreationInputToHouseholdInvitationDatabaseCreationInput(providedInput)

	input.ID = identifiers.New()
	input.DestinationHouseholdID = householdID
	input.FromUser = userID

	token, err := s.secretGenerator.GenerateBase64EncodedString(ctx, 64)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating invitation token")
		errRes := types.NewAPIErrorResponse("generating invitation token", types.ErrSecretGeneration, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	input.Token = token

	user, err := s.userDataManager.GetUserByEmail(ctx, input.ToEmail)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "fetching user ID by email")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	if user != nil {
		input.ToUser = &user.ID
	}

	householdInvitation, err := s.householdInvitationDataManager.CreateHouseholdInvitation(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating household invitation")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:           types.HouseholdInvitationCreatedCustomerEventType,
		HouseholdInvitation: householdInvitation,
		HouseholdID:         householdID,
		UserID:              userID,
	}
	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, householdInvitation, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a household invitation.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine relevant household invitation ID.
	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	tracing.AttachToSpan(span, keys.HouseholdIDKey, sessionCtxData.ActiveHouseholdID)
	logger = logger.WithValue(keys.HouseholdIDKey, sessionCtxData.ActiveHouseholdID)

	// fetch the household invitation from the database.
	householdInvitation, err := s.householdInvitationDataManager.GetHouseholdInvitationByHouseholdAndID(ctx, sessionCtxData.ActiveHouseholdID, householdInvitationID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("No rows found in household invitation database")
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching household invitation from database")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
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

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	filter := types.ExtractQueryFilterFromRequest(req)
	filter.AttachToLogger(logger)

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := s.householdIDFetcher(req)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	invitations, err := s.householdInvitationDataManager.GetPendingHouseholdInvitationsForUser(ctx, userID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching outbound invites")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	s.encoderDecoder.RespondWithData(ctx, res, invitations)
}

func (s *service) OutboundInvitesHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	filter := types.ExtractQueryFilterFromRequest(req)
	filter.AttachToLogger(logger)

	logger.Debug("fetching outbound invites for household")

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	userID := sessionCtxData.Requester.UserID
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdID := s.householdIDFetcher(req)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	invitations, err := s.householdInvitationDataManager.GetPendingHouseholdInvitationsFromUser(ctx, userID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching outbound invites")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	logger.Debug("responding with outbound invites for household")

	s.encoderDecoder.RespondWithData(ctx, res, invitations)
}

func (s *service) AcceptInviteHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdInvitationUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	// note, this is where you would call providedInput.ValidateWithContext, if that currently had any effect.

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	invitation, err := s.householdInvitationDataManager.GetHouseholdInvitationByTokenAndID(ctx, providedInput.Token, householdInvitationID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving invitation")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	if err = s.householdInvitationDataManager.AcceptHouseholdInvitation(ctx, invitation.ID, providedInput.Token, providedInput.Note); err != nil {
		observability.AcknowledgeError(err, logger, span, "accepting invitation")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:             types.HouseholdInvitationAcceptedCustomerEventType,
		HouseholdID:           invitation.DestinationHousehold.ID,
		HouseholdInvitationID: householdInvitationID,
		UserID:                sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

func (s *service) CancelInviteHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdInvitationUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	// note, this is where you would call providedInput.ValidateWithContext, if that currently had any effect.

	userID := sessionCtxData.Requester.UserID
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	invitation, err := s.householdInvitationDataManager.GetHouseholdInvitationByTokenAndID(ctx, providedInput.Token, householdInvitationID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving invitation")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	if err = s.householdInvitationDataManager.CancelHouseholdInvitation(ctx, invitation.ID, providedInput.Note); err != nil {
		observability.AcknowledgeError(err, logger, span, "cancelling invitation")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:             types.HouseholdInvitationCanceledCustomerEventType,
		HouseholdID:           invitation.DestinationHousehold.ID,
		HouseholdInvitationID: householdInvitationID,
		UserID:                sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

func (s *service) RejectInviteHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine relevant user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdInvitationUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	// note, this is where you would call providedInput.ValidateWithContext, if that currently had any effect.

	userID := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.UserIDKey, userID)

	householdInvitationID := s.householdInvitationIDFetcher(req)
	tracing.AttachToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	invitation, err := s.householdInvitationDataManager.GetHouseholdInvitationByTokenAndID(ctx, providedInput.Token, householdInvitationID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving invitation")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	if err = s.householdInvitationDataManager.RejectHouseholdInvitation(ctx, invitation.ID, providedInput.Note); err != nil {
		observability.AcknowledgeError(err, logger, span, "rejecting invitation")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:             types.HouseholdInvitationRejectedCustomerEventType,
		HouseholdID:           invitation.DestinationHousehold.ID,
		HouseholdInvitationID: householdInvitationID,
		UserID:                sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}
