package recipepreptasks

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
	// RecipePrepTaskIDURIParamKey is a standard string that we'll use to refer to recipe prep task IDs with.
	RecipePrepTaskIDURIParamKey = "recipePrepTaskID"
)

// CreateHandler is our recipe prep task creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
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

	// read parsed input struct from request body.
	providedInput := new(types.RecipePrepTaskCreationRequestInput)
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

	input := converters.ConvertRecipePrepTaskCreationRequestInputToRecipePrepTaskDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	for i := range input.TaskSteps {
		input.TaskSteps[i].ID = identifiers.New()
		input.TaskSteps[i].BelongsToRecipePrepTask = input.ID
	}

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	input.BelongsToRecipe = recipeID
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, input.ID)

	recipePrepTask, err := s.recipePrepTaskDataManager.CreateRecipePrepTask(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating recipe prep task")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:      types.RecipePrepTaskCreatedCustomerEventType,
		RecipePrepTask: recipePrepTask,
		HouseholdID:    sessionCtxData.ActiveHouseholdID,
		UserID:         sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
	}

	for _, step := range recipePrepTask.TaskSteps {
		dcm = &types.DataChangeMessage{
			EventType:          types.RecipePrepTaskStepCreatedCustomerEventType,
			RecipePrepTask:     recipePrepTask,
			RecipePrepTaskStep: step,
			HouseholdID:        sessionCtxData.ActiveHouseholdID,
			UserID:             sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
		}
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, recipePrepTask, http.StatusAccepted)
}

// ReadHandler returns a GET handler that returns a recipe prep task.
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe prep task ID.
	recipePrepTaskID := s.recipePrepTaskIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)

	// fetch recipe prep task from database.
	x, err := s.recipePrepTaskDataManager.GetRecipePrepTask(ctx, recipeID, recipePrepTaskID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe prep task")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ListHandler is our list route.
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	recipePrepTasks, err := s.recipePrepTaskDataManager.GetRecipePrepTasksForRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipePrepTasks = []*types.RecipePrepTask{}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe prep tasks")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipePrepTasks)
}

// UpdateHandler returns a handler that updates a recipe prep task.
func (s *service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
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

	// check for parsed input attached to session context data.
	input := new(types.RecipePrepTaskUpdateRequestInput)
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe prep task ID.
	recipePrepTaskID := s.recipePrepTaskIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)

	// fetch recipe prep task from database.
	recipePrepTask, err := s.recipePrepTaskDataManager.GetRecipePrepTask(ctx, recipeID, recipePrepTaskID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe prep task for update")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// update the recipe prep task.
	recipePrepTask.Update(input)

	if err = s.recipePrepTaskDataManager.UpdateRecipePrepTask(ctx, recipePrepTask); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating recipe prep task")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:      types.RecipePrepTaskUpdatedCustomerEventType,
		RecipePrepTask: recipePrepTask,
		HouseholdID:    sessionCtxData.ActiveHouseholdID,
		UserID:         sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipePrepTask)
}

// ArchiveHandler returns a handler that archives a recipe prep task.
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe prep task ID.
	recipePrepTaskID := s.recipePrepTaskIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)

	exists, existenceCheckErr := s.recipePrepTaskDataManager.RecipePrepTaskExists(ctx, recipeID, recipePrepTaskID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking recipe prep task existence")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	}

	if err = s.recipePrepTaskDataManager.ArchiveRecipePrepTask(ctx, recipeID, recipePrepTaskID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe prep task")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipePrepTaskArchivedCustomerEventType,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
