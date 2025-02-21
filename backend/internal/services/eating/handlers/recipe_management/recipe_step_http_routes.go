package recipemanagement

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	servertiming "github.com/mitchellh/go-server-timing"
)

// CreateRecipeStepHandler is our recipe step creation route.
func (s *service) CreateRecipeStepHandler(res http.ResponseWriter, req *http.Request) {
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
	providedInput := new(types.RecipeStepCreationRequestInput)
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

	input := converters.ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	for i := range input.Ingredients {
		input.Ingredients[i].ID = identifiers.New()
	}

	for j := range input.Products {
		input.Products[j].ID = identifiers.New()
	}

	for j := range input.Instruments {
		input.Instruments[j].ID = identifiers.New()
	}

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	input.BelongsToRecipe = recipeID
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, input.ID)

	createTimer := timing.NewMetric("database").WithDesc("create").Start()
	recipeStep, err := s.recipeManagementDataManager.CreateRecipeStep(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating recipe step")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	createTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeStepCreatedServiceEventType,
		RecipeStep:  recipeStep,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.RecipeStep]{
		Details: responseDetails,
		Data:    recipeStep,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// ReadRecipeStepHandler returns a GET handler that returns a recipe step.
func (s *service) ReadRecipeStepHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	// fetch recipe step from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	x, err := s.recipeManagementDataManager.GetRecipeStep(ctx, recipeID, recipeStepID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe step")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	responseValue := &types.APIResponse[*types.RecipeStep]{
		Details: responseDetails,
		Data:    x,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// ListRecipeStepsHandler is our list route.
func (s *service) ListRecipeStepsHandler(res http.ResponseWriter, req *http.Request) {
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

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	recipeSteps, err := s.recipeManagementDataManager.GetRecipeSteps(ctx, recipeID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipeSteps = &filtering.QueryFilteredResult[types.RecipeStep]{Data: []*types.RecipeStep{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe steps")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	responseValue := &types.APIResponse[[]*types.RecipeStep]{
		Details:    responseDetails,
		Data:       recipeSteps.Data,
		Pagination: &recipeSteps.Pagination,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// UpdateRecipeStepHandler returns a handler that updates a recipe step.
func (s *service) UpdateRecipeStepHandler(res http.ResponseWriter, req *http.Request) {
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
	input := new(types.RecipeStepUpdateRequestInput)
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	// fetch recipe step from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	recipeStep, err := s.recipeManagementDataManager.GetRecipeStep(ctx, recipeID, recipeStepID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe step for update")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	// update the recipe step.
	recipeStep.Update(input)

	updateTimer := timing.NewMetric("database").WithDesc("update").Start()
	if err = s.recipeManagementDataManager.UpdateRecipeStep(ctx, recipeStep); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating recipe step")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	updateTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeStepUpdatedServiceEventType,
		RecipeStep:  recipeStep,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.RecipeStep]{
		Details: responseDetails,
		Data:    recipeStep,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// ArchiveRecipeStepHandler returns a handler that archives a recipe step.
func (s *service) ArchiveRecipeStepHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	existenceTimer := timing.NewMetric("database").WithDesc("existence check").Start()
	exists, err := s.recipeManagementDataManager.RecipeStepExists(ctx, recipeID, recipeStepID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking recipe step existence")
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
	if err = s.recipeManagementDataManager.ArchiveRecipeStep(ctx, recipeID, recipeStepID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe step")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	archiveTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeStepArchivedServiceEventType,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.RecipeStep]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// RecipeStepImageUploadHandler updates a user's avatar.
func (s *service) RecipeStepImageUploadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	logger.Info("ImageUploadRecipeStepHandler invoked")

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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	images, err := s.imageUploadProcessor.ProcessFiles(ctx, req, "upload")
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "processing provided image")
		errorResponse := &types.APIResponse[any]{
			Details: types.ResponseDetails{
				TraceID: span.SpanContext().TraceID().String(),
			},
			Error: &types.APIError{
				Message: "invalid input attached to request",
			},
		}

		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errorResponse, http.StatusBadRequest)
		return
	}

	logger = logger.WithValue("image_qty", len(images))

	createdImages := []*types.RecipeMedia{}
	for i, img := range images {
		internalPath := fmt.Sprintf("%s/steps/%s/%d_%s", recipeID, recipeStepID, time.Now().Unix(), img.Filename)
		logger = logger.WithValue("internal_path", internalPath).WithValue("file_size", len(img.Data))

		if err = s.uploadManager.SaveFile(ctx, internalPath, img.Data); err != nil {
			observability.AcknowledgeError(err, logger, span, "saving provided image")

			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}

		input := &types.RecipeMediaDatabaseCreationInput{
			ID:                  identifiers.New(),
			BelongsToRecipe:     &recipeID,
			BelongsToRecipeStep: &recipeStepID,
			MimeType:            img.ContentType,
			InternalPath:        internalPath,
			ExternalPath:        fmt.Sprintf("%s/%s", s.cfg.PublicMediaURLPrefix, internalPath),
		}

		var created *types.RecipeMedia
		createTimer := timing.NewMetric("database").WithDesc(fmt.Sprintf("create #%d", i)).Start()
		created, err = s.recipeManagementDataManager.CreateRecipeMedia(ctx, input)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "saving recipe media record")
			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}
		createTimer.Stop()
		createdImages = append(createdImages, created)
	}

	logger.Info("media data saved in database")

	responseValue := &types.APIResponse[[]*types.RecipeMedia]{
		Details: responseDetails,
		Data:    createdImages,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}
