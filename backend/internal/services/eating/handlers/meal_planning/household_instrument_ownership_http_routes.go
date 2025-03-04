package mealplanning

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	servertiming "github.com/mitchellh/go-server-timing"
)

// CreateHouseholdInstrumentOwnershipHandler is our household instrument ownership creation route.
func (s *service) CreateHouseholdInstrumentOwnershipHandler(res http.ResponseWriter, req *http.Request) {
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
	providedInput := new(types.HouseholdInstrumentOwnershipCreationRequestInput)
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

	input := converters.ConvertHouseholdInstrumentOwnershipCreationRequestInputToHouseholdInstrumentOwnershipDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	input.BelongsToHousehold = sessionCtxData.ActiveHouseholdID

	tracing.AttachToSpan(span, keys.InstrumentOwnershipIDKey, input.ID)

	createTimer := timing.NewMetric("database").WithDesc("create").Start()
	householdInstrumentOwnership, err := s.mealPlanningDataManager.CreateHouseholdInstrumentOwnership(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating household instrument ownership")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	createTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:                    types.HouseholdInstrumentOwnershipCreatedServiceEventType,
		HouseholdInstrumentOwnership: householdInstrumentOwnership,
		UserID:                       sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.HouseholdInstrumentOwnership]{
		Details: responseDetails,
		Data:    householdInstrumentOwnership,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// ReadHouseholdInstrumentOwnershipHandler returns a GET handler that returns a household instrument ownership.
func (s *service) ReadHouseholdInstrumentOwnershipHandler(res http.ResponseWriter, req *http.Request) {
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

	householdID := sessionCtxData.ActiveHouseholdID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine household instrument ownership ID.
	householdInstrumentOwnershipID := s.householdInstrumentOwnershipIDFetcher(req)
	tracing.AttachToSpan(span, keys.InstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	logger = logger.WithValue(keys.InstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	// fetch household instrument ownership from database.
	x, err := s.mealPlanningDataManager.GetHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID, householdID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving household instrument ownership")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[*types.HouseholdInstrumentOwnership]{
		Details: responseDetails,
		Data:    x,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// ListHouseholdInstrumentOwnershipHandler is our list route.
func (s *service) ListHouseholdInstrumentOwnershipHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	filter := filtering.ExtractQueryFilterFromRequest(req)
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

	householdID := sessionCtxData.ActiveHouseholdID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	householdInstrumentOwnerships, err := s.mealPlanningDataManager.GetHouseholdInstrumentOwnerships(ctx, householdID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		householdInstrumentOwnerships = &filtering.QueryFilteredResult[types.HouseholdInstrumentOwnership]{Data: []*types.HouseholdInstrumentOwnership{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving household instrument ownerships")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.HouseholdInstrumentOwnership]{
		Details:    responseDetails,
		Data:       householdInstrumentOwnerships.Data,
		Pagination: &householdInstrumentOwnerships.Pagination,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// UpdateHouseholdInstrumentOwnershipHandler returns a handler that updates a household instrument ownership.
func (s *service) UpdateHouseholdInstrumentOwnershipHandler(res http.ResponseWriter, req *http.Request) {
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

	householdID := sessionCtxData.ActiveHouseholdID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// check for parsed input attached to session context data.
	input := new(types.HouseholdInstrumentOwnershipUpdateRequestInput)
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

	// determine household instrument ownership ID.
	householdInstrumentOwnershipID := s.householdInstrumentOwnershipIDFetcher(req)
	tracing.AttachToSpan(span, keys.InstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	logger = logger.WithValue(keys.InstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	// fetch household instrument ownership from database.
	householdInstrumentOwnership, err := s.mealPlanningDataManager.GetHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID, householdID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving household instrument ownership for update")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// update the household instrument ownership.
	householdInstrumentOwnership.Update(input)

	updateTimer := timing.NewMetric("database").WithDesc("update").Start()
	if err = s.mealPlanningDataManager.UpdateHouseholdInstrumentOwnership(ctx, householdInstrumentOwnership); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating household instrument ownership")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	updateTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:                    types.HouseholdInstrumentOwnershipUpdatedServiceEventType,
		HouseholdInstrumentOwnership: householdInstrumentOwnership,
		UserID:                       sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.HouseholdInstrumentOwnership]{
		Details: responseDetails,
		Data:    householdInstrumentOwnership,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// ArchiveHouseholdInstrumentOwnershipHandler returns a handler that archives a household instrument ownership.
func (s *service) ArchiveHouseholdInstrumentOwnershipHandler(res http.ResponseWriter, req *http.Request) {
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

	householdID := sessionCtxData.ActiveHouseholdID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine household instrument ownership ID.
	householdInstrumentOwnershipID := s.householdInstrumentOwnershipIDFetcher(req)
	tracing.AttachToSpan(span, keys.InstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	logger = logger.WithValue(keys.InstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	existenceTimer := timing.NewMetric("database").WithDesc("existence check").Start()
	exists, err := s.mealPlanningDataManager.HouseholdInstrumentOwnershipExists(ctx, householdInstrumentOwnershipID, householdID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking household instrument ownership existence")
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
	if err = s.mealPlanningDataManager.ArchiveHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID, householdID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving household instrument ownership")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	archiveTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType: types.HouseholdInstrumentOwnershipArchivedServiceEventType,
		UserID:    sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.HouseholdInstrumentOwnership]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}
