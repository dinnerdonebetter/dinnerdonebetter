package validenumerations

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

// CreateValidIngredientMeasurementUnitHandler is our valid ingredient measurement unit creation route.
func (s *service) CreateValidIngredientMeasurementUnitHandler(res http.ResponseWriter, req *http.Request) {
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
	providedInput := new(types.ValidIngredientMeasurementUnitCreationRequestInput)
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

	input := converters.ConvertValidIngredientMeasurementUnitCreationRequestInputToValidIngredientMeasurementUnitDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, input.ID)

	createTimer := timing.NewMetric("database").WithDesc("create").Start()
	validIngredientMeasurementUnit, err := s.validEnumerationDataManager.CreateValidIngredientMeasurementUnit(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid ingredient measurement unit")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	createTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:                      types.ValidIngredientMeasurementUnitCreatedServiceEventType,
		ValidIngredientMeasurementUnit: validIngredientMeasurementUnit,
		UserID:                         sessionCtxData.Requester.UserID,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.ValidIngredientMeasurementUnit]{
		Details: responseDetails,
		Data:    validIngredientMeasurementUnit,
	}

	// encode the response.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// ReadValidIngredientMeasurementUnitHandler returns a GET handler that returns a valid ingredient measurement unit.
func (s *service) ReadValidIngredientMeasurementUnitHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine valid ingredient measurement unit ID.
	validIngredientMeasurementUnitID := s.validIngredientMeasurementUnitIDFetcher(req)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	// fetch valid ingredient measurement unit from database.
	x, err := s.validEnumerationDataManager.GetValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient measurement unit")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[*types.ValidIngredientMeasurementUnit]{
		Details: responseDetails,
		Data:    x,
	}

	// encode the response.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ListValidIngredientMeasurementUnitsHandler is our list route.
func (s *service) ListValidIngredientMeasurementUnitsHandler(res http.ResponseWriter, req *http.Request) {
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
	tracing.AttachQueryFilterToSpan(span, filter)

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

	validIngredientMeasurementUnits, err := s.validEnumerationDataManager.GetValidIngredientMeasurementUnits(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validIngredientMeasurementUnits = &types.QueryFilteredResult[types.ValidIngredientMeasurementUnit]{Data: []*types.ValidIngredientMeasurementUnit{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient measurement units")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.ValidIngredientMeasurementUnit]{
		Details:    responseDetails,
		Pagination: &validIngredientMeasurementUnits.Pagination,
		Data:       validIngredientMeasurementUnits.Data,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// UpdateValidIngredientMeasurementUnitHandler returns a handler that updates a valid ingredient measurement unit.
func (s *service) UpdateValidIngredientMeasurementUnitHandler(res http.ResponseWriter, req *http.Request) {
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
	input := new(types.ValidIngredientMeasurementUnitUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		logger.Error("error encountered decoding request body", err)
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.Error("provided input was invalid", err)
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	// determine valid ingredient measurement unit ID.
	validIngredientMeasurementUnitID := s.validIngredientMeasurementUnitIDFetcher(req)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	// fetch valid ingredient measurement unit from database.
	validIngredientMeasurementUnit, err := s.validEnumerationDataManager.GetValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient measurement unit for update")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// update the valid ingredient measurement unit.
	validIngredientMeasurementUnit.Update(input)

	updateTimer := timing.NewMetric("database").WithDesc("update").Start()
	if err = s.validEnumerationDataManager.UpdateValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnit); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid ingredient measurement unit")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	updateTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:                      types.ValidIngredientMeasurementUnitUpdatedServiceEventType,
		ValidIngredientMeasurementUnit: validIngredientMeasurementUnit,
		UserID:                         sessionCtxData.Requester.UserID,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.ValidIngredientMeasurementUnit]{
		Details: responseDetails,
		Data:    validIngredientMeasurementUnit,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ArchiveValidIngredientMeasurementUnitHandler returns a handler that archives a valid ingredient measurement unit.
func (s *service) ArchiveValidIngredientMeasurementUnitHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine valid ingredient measurement unit ID.
	validIngredientMeasurementUnitID := s.validIngredientMeasurementUnitIDFetcher(req)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	existenceTimer := timing.NewMetric("database").WithDesc("existence check").Start()
	exists, err := s.validEnumerationDataManager.ValidIngredientMeasurementUnitExists(ctx, validIngredientMeasurementUnitID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking valid ingredient measurement unit existence")
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
	if err = s.validEnumerationDataManager.ArchiveValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnitID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid ingredient measurement unit")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	archiveTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType: types.ValidIngredientMeasurementUnitArchivedServiceEventType,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.ValidIngredientMeasurementUnit]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// SearchValidIngredientMeasurementUnitsByIngredientHandler is our valid ingredient measurement unit search route.
func (s *service) SearchValidIngredientMeasurementUnitsByIngredientHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachQueryFilterToSpan(span, filter)
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

	validIngredientMeasurementUnits, err := s.validEnumerationDataManager.GetValidIngredientMeasurementUnitsForIngredient(ctx, validIngredientID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid ingredient measurement units")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.ValidIngredientMeasurementUnit]{
		Details:    responseDetails,
		Pagination: &validIngredientMeasurementUnits.Pagination,
		Data:       validIngredientMeasurementUnits.Data,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// SearchValidIngredientMeasurementUnitsByMeasurementUnitHandler is our valid ingredient measurement unit search route.
func (s *service) SearchValidIngredientMeasurementUnitsByMeasurementUnitHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	validMeasurementUnitID := s.validMeasurementUnitIDFetcher(req)
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

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

	validIngredientMeasurementUnits, err := s.validEnumerationDataManager.GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx, validMeasurementUnitID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid ingredient measurement units")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.ValidIngredientMeasurementUnit]{
		Details:    responseDetails,
		Pagination: &validIngredientMeasurementUnits.Pagination,
		Data:       validIngredientMeasurementUnits.Data,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}
