package oauth2clients

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

var _ types.OAuth2ClientDataService = (*service)(nil)

const (
	// OAuth2ClientIDURIParamKey is used for referring to OAuth2 client IDs in router params.
	OAuth2ClientIDURIParamKey = "oauth2ClientID"

	clientIDSize     = 16
	clientSecretSize = 16
)

// ListHandler is a handler that returns a list of OAuth2 clients.
func (s *service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req)
	logger = filter.AttachToLogger(logger)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	// determine user.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// fetch OAuth2 clients.
	oauth2Clients, err := s.oauth2ClientDataManager.GetOAuth2Clients(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// just return an empty list if there are no results.
		oauth2Clients = &types.QueryFilteredResult[types.OAuth2Client]{
			Data: []*types.OAuth2Client{},
		}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching OAuth2 clients from database")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// encode response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, oauth2Clients)
}

// CreateHandler is our OAuth2 client creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// check session context data for user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// fetch creation input from session context data.
	input := new(types.OAuth2ClientCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	// keep relevant data in mind.
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	dbInput := converters.ConvertOAuth2ClientCreationRequestInputToOAuth2ClientDatabaseCreationInput(input)

	// set some data.
	if dbInput.ClientID, err = s.secretGenerator.GenerateHexEncodedString(ctx, clientIDSize); err != nil {
		observability.AcknowledgeError(err, logger, span, "generating client id")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if dbInput.ClientSecret, err = s.secretGenerator.GenerateHexEncodedString(ctx, clientSecretSize); err != nil {
		observability.AcknowledgeError(err, logger, span, "generating client secret")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dbInput.ID = identifiers.New()

	// create the client.
	client, err := s.oauth2ClientDataManager.CreateOAuth2Client(ctx, dbInput)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating OAuth2 client")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	tracing.AttachToSpan(span, keys.OAuth2ClientClientIDKey, client.ID)

	dcm := &types.DataChangeMessage{
		EventType:      types.OAuth2ClientCreatedCustomerEventType,
		OAuth2ClientID: client.ID,
		HouseholdID:    sessionCtxData.ActiveHouseholdID,
		UserID:         sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	resObj := &types.OAuth2ClientCreationResponse{
		ID:           client.ID,
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, resObj, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns an OAuth2 client.
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

	// determine OAuth2 client ID.
	oauth2ClientID := s.urlClientIDExtractor(req)
	tracing.AttachToSpan(span, keys.OAuth2ClientClientIDKey, oauth2ClientID)
	logger = logger.WithValue(keys.OAuth2ClientIDKey, oauth2ClientID)

	// fetch OAuth2 client from database.
	x, err := s.oauth2ClientDataManager.GetOAuth2ClientByDatabaseID(ctx, oauth2ClientID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching OAuth2 client from database")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ArchiveHandler returns a handler that archives an OAuth2 client.
func (s *service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine OAuth2 client ID.
	oauth2ClientID := s.urlClientIDExtractor(req)
	logger = logger.WithValue(keys.OAuth2ClientIDKey, oauth2ClientID)
	tracing.AttachToSpan(span, keys.OAuth2ClientClientIDKey, oauth2ClientID)

	// archive the OAuth2 client in the database.
	err = s.oauth2ClientDataManager.ArchiveOAuth2Client(ctx, oauth2ClientID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving OAuth2 client")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:      types.OAuth2ClientArchivedCustomerEventType,
		OAuth2ClientID: oauth2ClientID,
		HouseholdID:    sessionCtxData.ActiveHouseholdID,
		UserID:         sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
