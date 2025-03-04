package admin

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/email"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	"github.com/dinnerdonebetter/backend/pkg/types"

	servertiming "github.com/mitchellh/go-server-timing"
)

var _ types.AdminDataService = (*service)(nil)

// UserAccountStatusChangeHandler changes a user's status.
func (s *service) UserAccountStatusChangeHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// read parsed input struct from request body.
	decodeTimer := timing.NewMetric("decoding")
	decodeTimer.Start()
	input := new(types.UserAccountStatusUpdateInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}
	decodeTimer.Stop()

	validationTimer := timing.NewMetric("validation")
	validationTimer.Start()
	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}
	validationTimer.Stop()

	logger = logger.WithValue("new_status", input.NewStatus)

	if !sessionCtxData.Requester.ServicePermissions.CanUpdateUserAccountStatuses() {
		// this should never happen in production
		errorResponse := &types.APIResponse[any]{
			Details: types.ResponseDetails{
				TraceID: span.SpanContext().TraceID().String(),
			},
			Error: &types.APIError{
				Message: "inadequate permissions for route",
			},
		}
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errorResponse, http.StatusForbidden)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue("ban_giver", requester).WithValue("status_change_recipient", input.TargetUserID)

	if err = s.userDB.UpdateUserAccountStatus(ctx, input.TargetUserID, input); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
			return
		} else {
			observability.AcknowledgeError(err, logger, span, "retrieving session context data")
			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		}
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

// WriteArbitraryQueueMessageHandler publishes an arbitrary message to a given queue.
func (s *service) WriteArbitraryQueueMessageHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)
	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// read parsed input struct from request body.
	decodeTimer := timing.NewMetric("decoding request body")
	decodeTimer.Start()
	input := new(types.ArbitraryQueueMessageRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}
	decodeTimer.Stop()

	validationTimer := timing.NewMetric("validation")
	validationTimer.Start()
	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}
	validationTimer.Stop()

	var (
		dest      any
		topicName string
	)

	decodeTimer = timing.NewMetric("decoding message body")
	decodeTimer.Start()
	switch input.QueueName {
	case "data_changes":
		topicName = s.queuesConfig.DataChangesTopicName
		dest = &types.DataChangeMessage{}
		if err := s.encoderDecoder.DecodeBytes(ctx, []byte(input.Body), dest); err != nil {
			observability.AcknowledgeError(err, logger, span, "decoding message queue body")
			errRes := types.NewAPIErrorResponse("decoding message queue body", types.ErrDecodingRequestInput, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
			return
		}
		if x, ok := dest.(*types.DataChangeMessage); ok {
			x.RequestID = identifiers.New()
		}
	case "outbound_emails":
		topicName = s.queuesConfig.OutboundEmailsTopicName
		dest = &email.OutboundEmailMessage{}
		if err := s.encoderDecoder.DecodeBytes(ctx, []byte(input.Body), dest); err != nil {
			observability.AcknowledgeError(err, logger, span, "decoding message queue body")
			errRes := types.NewAPIErrorResponse("decoding message queue body", types.ErrDecodingRequestInput, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
			return
		}
	case "search_index_requests":
		topicName = s.queuesConfig.SearchIndexRequestsTopicName
		dest = &textsearch.IndexRequest{}
		if err := s.encoderDecoder.DecodeBytes(ctx, []byte(input.Body), dest); err != nil {
			observability.AcknowledgeError(err, logger, span, "decoding message queue body")
			errRes := types.NewAPIErrorResponse("decoding message queue body", types.ErrDecodingRequestInput, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
			return
		}
		if x, ok := dest.(*textsearch.IndexRequest); ok {
			x.RequestID = identifiers.New()
		}
	case "user_data_aggregator":
		topicName = s.queuesConfig.UserDataAggregationTopicName
		dest = &types.UserDataAggregationRequest{}
		if err := s.encoderDecoder.DecodeBytes(ctx, []byte(input.Body), dest); err != nil {
			observability.AcknowledgeError(err, logger, span, "decoding message queue body")
			errRes := types.NewAPIErrorResponse("decoding message queue body", types.ErrDecodingRequestInput, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
			return
		}
		if x, ok := dest.(*types.UserDataAggregationRequest); ok {
			x.RequestID = identifiers.New()
		}
	case "webhook_execution_requests":
		topicName = s.queuesConfig.WebhookExecutionRequestsTopicName
		dest = &types.WebhookExecutionRequest{}
		if err := s.encoderDecoder.DecodeBytes(ctx, []byte(input.Body), dest); err != nil {
			observability.AcknowledgeError(err, logger, span, "decoding message queue body")
			errRes := types.NewAPIErrorResponse("decoding message queue body", types.ErrDecodingRequestInput, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
			return
		}
		if x, ok := dest.(*types.WebhookExecutionRequest); ok {
			x.RequestID = identifiers.New()
		}
	}
	decodeTimer.Stop()

	publisherTimer := timing.NewMetric("instantiating publisher")
	publisherTimer.Start()
	publisher, err := s.publisherProvider.ProvidePublisher(topicName)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "instantiating message queue publisher")
		errRes := types.NewAPIErrorResponse("instantiating message queue publisher", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}
	publisherTimer.Stop()

	logger.WithValue("topic_name", topicName).WithValue("body", input.Body).Info("writing message to queue")

	publishTimer := timing.NewMetric("instantiating publisher")
	publishTimer.Start()
	if err = publisher.Publish(ctx, s.encoderDecoder.MustEncodeJSON(ctx, dest)); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing message to queue")
		errRes := types.NewAPIErrorResponse("publishing message to queue", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}
	publishTimer.Stop()

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, &types.ArbitraryQueueMessageResponse{Success: true}, http.StatusOK)
}
