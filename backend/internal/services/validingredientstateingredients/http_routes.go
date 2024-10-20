package validingredientstateingredients

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	servertiming "github.com/mitchellh/go-server-timing"
)

const (
	// ValidIngredientStateIngredientIDURIParamKey is a standard string that we'll use to refer to valid ingredient state ingredient IDs with.
	ValidIngredientStateIngredientIDURIParamKey = "validIngredientStateIngredientID"
	// ValidIngredientStateIDURIParamKey is a standard string that we'll use to refer to valid preparation IDs with.
	ValidIngredientStateIDURIParamKey = "validIngredientStateID"
	// ValidIngredientIDURIParamKey is a standard string that we'll use to refer to valid ingredient IDs with.
	ValidIngredientIDURIParamKey = "validIngredientID"
)

// CreateHandler is our valid ingredient state ingredient creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

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

	// read parsed input struct from request body.
	providedInput := new(types.ValidIngredientStateIngredientCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	input := converters.ConvertValidIngredientStateIngredientCreationRequestInputToValidIngredientStateIngredientDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, input.ID)

	createTimer := timing.NewMetric("database").WithDesc("create").Start()
	validIngredientStateIngredient, err := s.validIngredientStateIngredientDataManager.CreateValidIngredientStateIngredient(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid ingredient state ingredient")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	createTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:                      types.ValidIngredientStateIngredientCreatedServiceEventType,
		ValidIngredientStateIngredient: validIngredientStateIngredient,
		UserID:                         sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
	}

	responseValue := &types.APIResponse[*types.ValidIngredientStateIngredient]{
		Details: responseDetails,
		Data:    validIngredientStateIngredient,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid ingredient state ingredient.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

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

	// determine valid ingredient state ingredient ID.
	validIngredientStateIngredientID := s.validIngredientStateIngredientIDFetcher(req)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	// fetch valid ingredient state ingredient from database.
	x, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredient(ctx, validIngredientStateIngredientID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient state ingredient")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[*types.ValidIngredientStateIngredient]{
		Details: responseDetails,
		Data:    x,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ListHandler is our list route.
func (s *service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req).WithSpan(span)
	logger = filter.AttachToLogger(logger)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

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

	validIngredientStateIngredients, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredients(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validIngredientStateIngredients = &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{Data: []*types.ValidIngredientStateIngredient{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient state ingredients")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.ValidIngredientStateIngredient]{
		Details:    responseDetails,
		Data:       validIngredientStateIngredients.Data,
		Pagination: &validIngredientStateIngredients.Pagination,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// UpdateHandler returns a handler that updates a valid ingredient state ingredient.
func (s *service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

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

	// check for parsed input attached to session context data.
	input := new(types.ValidIngredientStateIngredientUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		logger.Error(err, "error encountered decoding request body")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.Error(err, "provided input was invalid")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	// determine valid ingredient state ingredient ID.
	validIngredientStateIngredientID := s.validIngredientStateIngredientIDFetcher(req)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	// fetch valid ingredient state ingredient from database.
	validIngredientStateIngredient, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredient(ctx, validIngredientStateIngredientID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient state ingredient for update")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// update the valid ingredient state ingredient.
	validIngredientStateIngredient.Update(input)

	updateTimer := timing.NewMetric("database").WithDesc("update").Start()
	if err = s.validIngredientStateIngredientDataManager.UpdateValidIngredientStateIngredient(ctx, validIngredientStateIngredient); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid ingredient state ingredient")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	updateTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:                      types.ValidIngredientStateIngredientUpdatedServiceEventType,
		ValidIngredientStateIngredient: validIngredientStateIngredient,
		UserID:                         sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	responseValue := &types.APIResponse[*types.ValidIngredientStateIngredient]{
		Details: responseDetails,
		Data:    validIngredientStateIngredient,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ArchiveHandler returns a handler that archives a valid ingredient state ingredient.
func (s *service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

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

	// determine valid ingredient state ingredient ID.
	validIngredientStateIngredientID := s.validIngredientStateIngredientIDFetcher(req)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	existenceTimer := timing.NewMetric("database").WithDesc("existence check").Start()
	exists, err := s.validIngredientStateIngredientDataManager.ValidIngredientStateIngredientExists(ctx, validIngredientStateIngredientID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking valid ingredient state ingredient existence")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	} else if !exists || errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	}
	existenceTimer.Stop()

	archiveTimer := timing.NewMetric("database").WithDesc("archive").Start()
	if err = s.validIngredientStateIngredientDataManager.ArchiveValidIngredientStateIngredient(ctx, validIngredientStateIngredientID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid ingredient state ingredient")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	archiveTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType: types.ValidIngredientStateIngredientArchivedServiceEventType,
		UserID:    sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	responseValue := &types.APIResponse[*types.ValidIngredientStateIngredient]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// SearchByIngredientHandler is our valid ingredient state ingredient search route.
func (s *service) SearchByIngredientHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)
	logger = filter.AttachToLogger(logger)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	validIngredientID := s.validIngredientIDFetcher(req)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)

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

	validIngredientStateIngredients, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredientsForIngredient(ctx, validIngredientID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid ingredient state ingredients")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.ValidIngredientStateIngredient]{
		Details:    responseDetails,
		Data:       validIngredientStateIngredients.Data,
		Pagination: &validIngredientStateIngredients.Pagination,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// SearchByIngredientStateHandler is our valid ingredient state ingredient search route.
func (s *service) SearchByIngredientStateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)
	logger = filter.AttachToLogger(logger)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	validIngredientStateID := s.validIngredientStateIDFetcher(req)
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)

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

	validIngredientStateIngredients, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredientsForIngredientState(ctx, validIngredientStateID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid ingredient state ingredients")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.ValidIngredientStateIngredient]{
		Details:    responseDetails,
		Data:       validIngredientStateIngredients.Data,
		Pagination: &validIngredientStateIngredients.Pagination,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}
