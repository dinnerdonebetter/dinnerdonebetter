package recipemanagement

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	"github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	servertiming "github.com/mitchellh/go-server-timing"
)

// CreateRecipeHandler is our recipe creation route.
func (s *service) CreateRecipeHandler(res http.ResponseWriter, req *http.Request) {
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
	providedInput := new(types.RecipeCreationRequestInput)
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

	input, err := converters.ConvertRecipeCreationRequestInputToRecipeDatabaseCreationInput(providedInput)
	if err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	input.CreatedByUser = sessionCtxData.Requester.UserID
	tracing.AttachToSpan(span, keys.RecipeIDKey, input.ID)

	createTimer := timing.NewMetric("database").WithDesc("create").Start()
	recipe, err := s.recipeManagementDataManager.CreateRecipe(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating recipe")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	createTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeCreatedServiceEventType,
		Recipe:      recipe,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.Recipe]{
		Details: responseDetails,
		Data:    recipe,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// ReadRecipeHandler returns a GET handler that returns a recipe.
func (s *service) ReadRecipeHandler(res http.ResponseWriter, req *http.Request) {
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

	// fetch recipe from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	x, err := s.recipeManagementDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	logger.Info("recipe retrieved")

	responseValue := &types.APIResponse[*types.Recipe]{
		Details: responseDetails,
		Data:    x,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// ListRecipesHandler is our list route.
func (s *service) ListRecipesHandler(res http.ResponseWriter, req *http.Request) {
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

	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	recipes, err := s.recipeManagementDataManager.GetRecipes(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipes = &filtering.QueryFilteredResult[types.Recipe]{Data: []*types.Recipe{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipes")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	responseValue := &types.APIResponse[[]*types.Recipe]{
		Details:    responseDetails,
		Data:       recipes.Data,
		Pagination: &recipes.Pagination,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// SearchRecipesHandler is our list route.
func (s *service) SearchRecipesHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	query := req.URL.Query().Get(textsearch.QueryKeySearch)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	logger = logger.WithValue(keys.SearchQueryKey, query)

	filter := filtering.ExtractQueryFilterFromRequest(req)
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	useDB := !s.cfg.UseSearchService || strings.TrimSpace(strings.ToLower(req.URL.Query().Get(filtering.QueryKeySearchWithDatabase))) == "true"
	logger = logger.WithValue("using_database", useDB)

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

	recipes := &filtering.QueryFilteredResult[types.Recipe]{
		Pagination: filter.ToPagination(),
	}

	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	if useDB {
		recipes, err = s.recipeManagementDataManager.SearchForRecipes(ctx, query, filter)
	} else {
		var recipeSubsets []*indexing.RecipeSearchSubset
		recipeSubsets, err = s.searchIndex.Search(ctx, query)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "external search for recipes")
			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}

		ids := []string{}
		for _, recipeSubset := range recipeSubsets {
			ids = append(ids, recipeSubset.ID)
		}

		recipes.Data, err = s.recipeManagementDataManager.GetRecipesWithIDs(ctx, ids)
	}

	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipes = &filtering.QueryFilteredResult[types.Recipe]{
			Pagination: filter.ToPagination(),
			Data:       []*types.Recipe{},
		}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for recipes")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	responseValue := &types.APIResponse[[]*types.Recipe]{
		Details:    responseDetails,
		Data:       recipes.Data,
		Pagination: &recipes.Pagination,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// UpdateRecipeHandler returns a handler that updates a recipe.
func (s *service) UpdateRecipeHandler(res http.ResponseWriter, req *http.Request) {
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
	input := new(types.RecipeUpdateRequestInput)
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

	// fetch recipe from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	recipe, err := s.recipeManagementDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe for update")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	if recipe.CreatedByUser != sessionCtxData.Requester.UserID {
		errRes := types.NewAPIErrorResponse("user is not creator", types.ErrUserIsNotAuthorized, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	// update the recipe.
	recipe.Update(input)

	updateTimer := timing.NewMetric("database").WithDesc("update").Start()
	if err = s.recipeManagementDataManager.UpdateRecipe(ctx, recipe); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating recipe")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	updateTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeUpdatedServiceEventType,
		Recipe:      recipe,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.Recipe]{
		Details: responseDetails,
		Data:    recipe,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// ArchiveRecipeHandler returns a handler that archives a recipe.
func (s *service) ArchiveRecipeHandler(res http.ResponseWriter, req *http.Request) {
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

	existenceTimer := timing.NewMetric("database").WithDesc("existence check").Start()
	exists, err := s.recipeManagementDataManager.RecipeExists(ctx, recipeID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking recipe existence")
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
	if err = s.recipeManagementDataManager.ArchiveRecipe(ctx, recipeID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	archiveTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeArchivedServiceEventType,
		RecipeID:    recipeID,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.Recipe]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// RecipeEstimatedPrepStepsHandler is a handler that returns expected prep steps for a given recipe.
func (s *service) RecipeEstimatedPrepStepsHandler(res http.ResponseWriter, req *http.Request) {
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

	// fetch recipe from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	_, err = s.recipeManagementDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	/* TODO:
	stepInputs, err := s.recipeAnalyzer.GenerateMealPlanTasksForRecipe(ctx, "", x)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating DAG for recipe")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	*/

	responseEvents := []*types.MealPlanTaskDatabaseCreationEstimate{}
	/* TODO:
	for _, input := range stepInputs {
		responseEvents = append(responseEvents, &types.MealPlanTaskDatabaseCreationEstimate{
			CreationExplanation: input.CreationExplanation,
		})
	}
	*/

	responseValue := &types.APIResponse[[]*types.MealPlanTaskDatabaseCreationEstimate]{
		Details: responseDetails,
		Data:    responseEvents,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// RecipeImageUploadHandler updates a user's avatar.
func (s *service) RecipeImageUploadHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	logger.Info("about to start processing image")

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
	logger.Info("processed images, saving files")

	var created []*types.RecipeMedia
	for i, img := range images {
		internalPath := fmt.Sprintf("%s/%d_%s", recipeID, time.Now().Unix(), img.Filename)
		logger = logger.WithValue("internal_path", internalPath).WithValue("file_size", len(img.Data))

		if err = s.uploadManager.SaveFile(ctx, internalPath, img.Data); err != nil {
			observability.AcknowledgeError(err, logger, span, "saving provided image")
			errRes := types.NewAPIErrorResponse("saving image", types.ErrMisbehavingDependency, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}

		logger.Info("image uploaded to file store, saving info in database")

		input := &types.RecipeMediaDatabaseCreationInput{
			ID:                  identifiers.New(),
			BelongsToRecipe:     &recipeID,
			BelongsToRecipeStep: nil,
			MimeType:            img.ContentType,
			InternalPath:        internalPath,
			ExternalPath:        fmt.Sprintf("%s/%s", s.cfg.PublicMediaURLPrefix, internalPath),
		}

		var createdMedia *types.RecipeMedia
		createTimer := timing.NewMetric("database").WithDesc(fmt.Sprintf("create #%d", i)).Start()
		createdMedia, err = s.recipeManagementDataManager.CreateRecipeMedia(ctx, input)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "saving recipe media record")
			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}
		createTimer.Stop()
		created = append(created, createdMedia)

		dcm := &types.DataChangeMessage{
			EventType:     types.RecipeMediaCreatedServiceEventType,
			RecipeID:      recipeID,
			RecipeMediaID: input.ID,
			HouseholdID:   sessionCtxData.ActiveHouseholdID,
			UserID:        sessionCtxData.Requester.UserID,
		}

		s.dataChangesPublisher.PublishAsync(ctx, dcm)

		logger.Info("image info saved in database")
	}

	responseValue := &types.APIResponse[[]*types.RecipeMedia]{
		Details: responseDetails,
		Data:    created,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// RecipeMermaidHandler returns a GET handler that returns a recipe in Mermaid format.
func (s *service) RecipeMermaidHandler(res http.ResponseWriter, req *http.Request) {
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

	// fetch recipe from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	_, err = s.recipeManagementDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	// TODO: graphDefinition := s.recipeAnalyzer.RenderMermaidDiagramForRecipe(ctx, x)

	responseValue := &types.APIResponse[string]{
		Details: responseDetails,
		Data:    "",
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// CloneRecipeHandler returns a POST handler that returns a cloned recipe.
func (s *service) CloneRecipeHandler(res http.ResponseWriter, req *http.Request) {
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

	// fetch recipe from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	x, err := s.recipeManagementDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	createTimer := timing.NewMetric("database").WithDesc("create").Start()
	created, err := s.recipeManagementDataManager.CreateRecipe(ctx, cloneRecipe(x, sessionCtxData.Requester.UserID))
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "cloning recipe")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	createTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeClonedServiceEventType,
		Recipe:      created,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.Recipe]{
		Details: responseDetails,
		Data:    created,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

func cloneRecipe(x *types.Recipe, userID string) *types.RecipeDatabaseCreationInput {
	ingredientProductIndicies := map[string]int{}
	instrumentProductIndicies := map[string]int{}
	vesselProductIndicies := map[string]int{}
	for _, step := range x.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.RecipeStepProductID != nil {
				ingredientProductIndicies[ingredient.ID] = x.FindStepIndexByID(x.FindStepForRecipeStepProductID(*ingredient.RecipeStepProductID).ID)
			}
		}

		for _, instrument := range step.Instruments {
			if instrument.RecipeStepProductID != nil {
				instrumentProductIndicies[instrument.ID] = x.FindStepIndexByID(x.FindStepForRecipeStepProductID(*instrument.RecipeStepProductID).ID)
			}
		}

		for _, vessel := range step.Vessels {
			if vessel.RecipeStepProductID != nil {
				vesselProductIndicies[vessel.ID] = x.FindStepIndexByID(x.FindStepForRecipeStepProductID(*vessel.RecipeStepProductID).ID)
			}
		}
	}

	// clone recipe.
	cloneInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(x)
	cloneInput.CreatedByUser = userID
	// TODO: cloneInput.ClonedFromRecipeID = &x.ID

	cloneInput.ID = identifiers.New()
	for i := range cloneInput.Steps {
		newRecipeStepID := identifiers.New()
		cloneInput.Steps[i].ID = newRecipeStepID
		for j := range cloneInput.Steps[i].Ingredients {
			if index, ok := ingredientProductIndicies[x.Steps[i].Ingredients[j].ID]; ok {
				cloneInput.Steps[i].Ingredients[j].ProductOfRecipeStepIndex = pointer.To(uint64(index))
			}
			cloneInput.Steps[i].Ingredients[j].ID = identifiers.New()
			cloneInput.Steps[i].Ingredients[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].Instruments {
			if index, ok := instrumentProductIndicies[x.Steps[i].Instruments[j].ID]; ok {
				cloneInput.Steps[i].Instruments[j].ProductOfRecipeStepIndex = pointer.To(uint64(index))
			}
			cloneInput.Steps[i].Instruments[j].ID = identifiers.New()
			cloneInput.Steps[i].Instruments[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].Vessels {
			if index, ok := vesselProductIndicies[x.Steps[i].Vessels[j].ID]; ok {
				cloneInput.Steps[i].Vessels[j].ProductOfRecipeStepIndex = pointer.To(uint64(index))
			}
			cloneInput.Steps[i].Vessels[j].ID = identifiers.New()
			cloneInput.Steps[i].Vessels[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].Products {
			cloneInput.Steps[i].Products[j].ID = identifiers.New()
			cloneInput.Steps[i].Products[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].CompletionConditions {
			newCompletionConditionID := identifiers.New()
			cloneInput.Steps[i].CompletionConditions[j].ID = newCompletionConditionID
			cloneInput.Steps[i].CompletionConditions[j].BelongsToRecipeStep = newRecipeStepID
			for k := range cloneInput.Steps[i].CompletionConditions[j].Ingredients {
				cloneInput.Steps[i].CompletionConditions[j].Ingredients[k].ID = identifiers.New()
				cloneInput.Steps[i].CompletionConditions[j].Ingredients[k].BelongsToRecipeStepCompletionCondition = newCompletionConditionID
			}
		}
	}

	// TODO: handle media here eventually

	for i := range cloneInput.PrepTasks {
		newPrepTaskID := identifiers.New()
		cloneInput.PrepTasks[i].ID = newPrepTaskID
		for j := range cloneInput.PrepTasks[i].TaskSteps {
			cloneInput.PrepTasks[i].TaskSteps[j].ID = identifiers.New()
			cloneInput.PrepTasks[i].TaskSteps[j].BelongsToRecipePrepTask = newPrepTaskID
		}
	}

	return cloneInput
}
