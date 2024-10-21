package workers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"

	servertiming "github.com/mitchellh/go-server-timing"
)

const (
	// ReportIDURIParamKey is a standard string that we'll use to refer to report IDs with.
	ReportIDURIParamKey = "userDataAggregationReportID"
)

// DataDeletionHandler deletes a user, which consequently deletes all data in the system.
func (s *service) DataDeletionHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	logger.Info("meal plan finalization worker invoked")

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
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

	destroyUserDataTimer := timing.NewMetric("database").WithDesc("destroy user data").Start()
	if err = s.dataPrivacyDataManager.DeleteUser(ctx, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "deleting user")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	destroyUserDataTimer.Stop()

	if err = s.dataChangesPublisher.Publish(ctx, &types.DataChangeMessage{
		Context:     nil,
		EventType:   types.UserDataDestroyedServiceEventType,
		UserID:      sessionCtxData.Requester.UserID,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
	}); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data")
	}

	responseValue := &types.APIResponse[types.DataDeletionResponse]{
		Data:    types.DataDeletionResponse{Successful: true},
		Details: responseDetails,
	}

	logger.Info("user data deleted")

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusAccepted)
}

// UserDataAggregationRequestHandler sends a message to the queue for data aggregation.
func (s *service) UserDataAggregationRequestHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	logger.Info("meal plan finalization worker invoked")

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
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

	reportID := identifiers.New()
	tracing.AttachToSpan(span, keys.UserDataAggregationReportIDKey, reportID)
	logger = logger.WithValue(keys.UserDataAggregationReportIDKey, reportID)

	if err = s.userDataAggregationPublisher.Publish(ctx, &types.UserDataAggregationRequest{
		UserID:   sessionCtxData.Requester.UserID,
		ReportID: reportID,
	}); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data aggregation request")
	}

	if err = s.dataChangesPublisher.Publish(ctx, &types.DataChangeMessage{
		Context: map[string]any{
			keys.UserDataAggregationReportIDKey: reportID,
		},
		EventType:   types.UserDataAggregationRequestServiceEventType,
		UserID:      sessionCtxData.Requester.UserID,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
	}); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	responseValue := &types.APIResponse[types.UserDataCollectionResponse]{
		Data:    types.UserDataCollectionResponse{ReportID: reportID},
		Details: responseDetails,
	}

	logger.Info("user data deleted")

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusAccepted)
}

// ReadUserDataAggregationReportHandler sends a message to the queue for data aggregation.
func (s *service) ReadUserDataAggregationReportHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	logger.Info("meal plan finalization worker invoked")

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
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

	// determine report ID.
	reportID := s.reportIDFetcher(req)
	tracing.AttachToSpan(span, keys.UserDataAggregationReportIDKey, reportID)
	logger = logger.WithValue(keys.UserDataAggregationReportIDKey, reportID)

	fileContents, err := s.uploader.ReadFile(ctx, fmt.Sprintf("%s.json", reportID))
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving report")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	var collection types.UserDataCollection
	if err = json.Unmarshal(fileContents, &collection); err != nil {
		observability.AcknowledgeError(err, logger, span, "unmarshalling report")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	logger.Info("read report from storage")

	responseValue := &types.APIResponse[types.UserDataCollection]{
		Data:    collection,
		Details: responseDetails,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}
