package webhooks

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

const (
	// WebhookIDURIParamKey is a standard string that we'll use to refer to webhook IDs with.
	WebhookIDURIParamKey = "webhookID"
)

var (
	_ types.WebhookDataService = (*service)(nil)
)

// CreateWebhookHandler is our webhook creation route.
func (s *service) CreateWebhookHandler(res http.ResponseWriter, req *http.Request) {
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

	providedInput := new(types.WebhookCreationRequestInput)
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

	input := converters.ConvertWebhookCreationRequestInputToWebhookDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	tracing.AttachToSpan(span, keys.WebhookIDKey, input.ID)
	input.BelongsToHousehold = sessionCtxData.ActiveHouseholdID

	webhook, err := s.webhookDataManager.CreateWebhook(ctx, input)
	logger.Debug("database call executed")
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating webhook in database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.WebhookCreatedCustomerEventType,
		Webhook:     webhook,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, webhook, http.StatusCreated)
}

// ListWebhooksHandler is our list route.
func (s *service) ListWebhooksHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	filter := types.ExtractQueryFilterFromRequest(req)
	logger := filter.AttachToLogger(s.logger)

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// find the webhooks.
	webhooks, err := s.webhookDataManager.GetWebhooks(ctx, sessionCtxData.ActiveHouseholdID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		webhooks = &types.QueryFilteredResult[types.Webhook]{
			Data: []*types.Webhook{},
		}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching webhooks")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode the response.
	s.encoderDecoder.RespondWithData(ctx, res, webhooks)
}

// ReadWebhookHandler returns a GET handler that returns an webhook.
func (s *service) ReadWebhookHandler(res http.ResponseWriter, req *http.Request) {
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
	webhookID := s.webhookIDFetcher(req)
	tracing.AttachToSpan(span, keys.WebhookIDKey, webhookID)
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)

	tracing.AttachToSpan(span, keys.HouseholdIDKey, sessionCtxData.ActiveHouseholdID)
	logger = logger.WithValue(keys.HouseholdIDKey, sessionCtxData.ActiveHouseholdID)

	// fetch the webhook from the database.
	webhook, err := s.webhookDataManager.GetWebhook(ctx, webhookID, sessionCtxData.ActiveHouseholdID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("No rows found in webhook database")
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching webhook from database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode the response.
	s.encoderDecoder.RespondWithData(ctx, res, webhook)
}

// ArchiveWebhookHandler returns a handler that archives an webhook.
func (s *service) ArchiveWebhookHandler(res http.ResponseWriter, req *http.Request) {
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
	webhookID := s.webhookIDFetcher(req)
	tracing.AttachToSpan(span, keys.WebhookIDKey, webhookID)
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)

	exists, webhookExistenceCheckErr := s.webhookDataManager.WebhookExists(ctx, webhookID, householdID)
	if webhookExistenceCheckErr != nil && !errors.Is(webhookExistenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		observability.AcknowledgeError(webhookExistenceCheckErr, logger, span, "checking item existence")
		return
	} else if !exists || errors.Is(webhookExistenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.webhookDataManager.ArchiveWebhook(ctx, webhookID, householdID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving webhook in database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.WebhookArchivedCustomerEventType,
		HouseholdID: householdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// let everybody go home.
	res.WriteHeader(http.StatusNoContent)
}
