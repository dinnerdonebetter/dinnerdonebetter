package recipes

import (
	"database/sql"
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	// RecipeIDURIParamKey is a standard string that we'll use to refer to recipe IDs with.
	RecipeIDURIParamKey = "recipeID"
)

// CreateHandler is our recipe creation route.
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

	recipe, err := s.recipeDataManager.CreateRecipe(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating recipe")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeCreatedCustomerEventType,
		Recipe:      recipe,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, recipe, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a recipe.
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

	// fetch recipe from database.
	x, err := s.recipeDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
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

	recipes, err := s.recipeDataManager.GetRecipes(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipes = &types.QueryFilteredResult[types.Recipe]{Data: []*types.Recipe{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipes")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipes)
}

// SearchHandler is our list route.
func (s *service) SearchHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	tracing.AttachRequestToSpan(span, req)
	useDB := !s.cfg.UseSearchService || strings.TrimSpace(strings.ToLower(req.URL.Query().Get("useDB"))) == "true"

	query := req.URL.Query().Get(types.SearchQueryKey)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	logger := s.logger.WithRequest(req).
		WithValue(keys.SearchQueryKey, query).
		WithValue("using_database", useDB)
	logger = filter.AttachToLogger(logger)

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

	recipes := &types.QueryFilteredResult[types.Recipe]{
		Pagination: filter.ToPagination(),
	}
	if useDB {
		recipes, err = s.recipeDataManager.SearchForRecipes(ctx, query, filter)
	} else {
		var recipeSubsets []*types.RecipeSearchSubset
		recipeSubsets, err = s.searchIndex.Search(ctx, query)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "external search for recipes")
			errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}

		ids := []string{}
		for _, recipeSubset := range recipeSubsets {
			ids = append(ids, recipeSubset.ID)
		}

		recipes.Data, err = s.recipeDataManager.GetRecipesWithIDs(ctx, ids)
	}

	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipes = &types.QueryFilteredResult[types.Recipe]{
			Pagination: filter.ToPagination(),
			Data:       []*types.Recipe{},
		}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for recipes")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipes)
}

// UpdateHandler returns a handler that updates a recipe.
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
	input := new(types.RecipeUpdateRequestInput)
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

	// fetch recipe from database.
	recipe, err := s.recipeDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe for update")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	if recipe.CreatedByUser != sessionCtxData.Requester.UserID {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthorized", http.StatusUnauthorized)
		return
	}

	// update the recipe.
	recipe.Update(input)

	if err = s.recipeDataManager.UpdateRecipe(ctx, recipe); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating recipe")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeUpdatedCustomerEventType,
		Recipe:      recipe,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipe)
}

// ArchiveHandler returns a handler that archives a recipe.
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

	exists, existenceCheckErr := s.recipeDataManager.RecipeExists(ctx, recipeID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking recipe existence")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.recipeDataManager.ArchiveRecipe(ctx, recipeID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.RecipeArchivedCustomerEventType,
		RecipeID:    recipeID,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// DAGHandler is a handler that returns a DAG image.
func (s *service) DAGHandler(res http.ResponseWriter, req *http.Request) {
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

	// fetch recipe from database.
	x, err := s.recipeDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dag, err := s.recipeAnalyzer.GenerateDAGDiagramForRecipe(ctx, x)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating DAG for recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	res.Header().Set("Content-type", "image/png")
	if err = png.Encode(res, dag); err != nil {
		observability.AcknowledgeError(err, logger, span, "encoding response")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}
}

// EstimatedPrepStepsHandler is a handler that returns expected prep steps for a given recipe.
func (s *service) EstimatedPrepStepsHandler(res http.ResponseWriter, req *http.Request) {
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

	// fetch recipe from database.
	x, err := s.recipeDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// we deliberately call this with fake data because
	stepInputs, err := s.recipeAnalyzer.GenerateMealPlanTasksForRecipe(ctx, "", x)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating DAG for recipe")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseEvents := []*types.MealPlanTaskDatabaseCreationEstimate{}
	for _, input := range stepInputs {
		responseEvents = append(responseEvents, &types.MealPlanTaskDatabaseCreationEstimate{
			CreationExplanation: input.CreationExplanation,
		})
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseEvents)
}

// ImageUploadHandler updates a user's avatar.
func (s *service) ImageUploadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

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

	logger.Info("about to start processing image")

	images, err := s.imageUploadProcessor.ProcessFiles(ctx, req, "upload")
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "processing provided image")
		s.encoderDecoder.EncodeInvalidInputResponse(ctx, res)
		return
	}

	logger = logger.WithValue("image_qty", len(images))
	logger.Info("processed images, saving files")

	for _, img := range images {
		internalPath := fmt.Sprintf("%s/%d_%s", recipeID, time.Now().Unix(), img.Filename)
		logger = logger.WithValue("internal_path", internalPath).WithValue("file_size", len(img.Data))

		if err = s.uploadManager.SaveFile(ctx, internalPath, img.Data); err != nil {
			observability.AcknowledgeError(err, logger, span, "saving provided image")
			s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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

		if _, dbWriteErr := s.recipeMediaDataManager.CreateRecipeMedia(ctx, input); dbWriteErr != nil {
			observability.AcknowledgeError(dbWriteErr, logger, span, "saving recipe media record")
			errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}

		dcm := &types.DataChangeMessage{
			EventType:     types.RecipeMediaCreatedCustomerEventType,
			RecipeID:      recipeID,
			RecipeMediaID: input.ID,
			HouseholdID:   sessionCtxData.ActiveHouseholdID,
			UserID:        sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}

		logger.Info("image info saved in database")
	}

	res.WriteHeader(http.StatusCreated)
}

// MermaidHandler returns a GET handler that returns a recipe in Mermaid format.
func (s *service) MermaidHandler(res http.ResponseWriter, req *http.Request) {
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

	// fetch recipe from database.
	x, err := s.recipeDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	graphDefinition := s.recipeAnalyzer.RenderMermaidDiagramForRecipe(ctx, x)

	res.Header().Set("Content-Type", "text/mermaid")

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, graphDefinition)
}

// CloneHandler returns a POST handler that returns a cloned recipe.
func (s *service) CloneHandler(res http.ResponseWriter, req *http.Request) {
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

	// fetch recipe from database.
	x, err := s.recipeDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

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
	cloneInput.CreatedByUser = sessionCtxData.Requester.UserID
	// TODO: cloneInput.ClonedFromRecipeID = &x.ID
	cloneInput.ID = identifiers.New()
	for i := range cloneInput.Steps {
		newRecipeStepID := identifiers.New()
		cloneInput.Steps[i].ID = newRecipeStepID
		for j := range cloneInput.Steps[i].Ingredients {
			if index, ok := ingredientProductIndicies[x.Steps[i].Ingredients[j].ID]; ok {
				cloneInput.Steps[i].Ingredients[j].ProductOfRecipeStepIndex = pointers.Pointer(uint64(index))
			}
			cloneInput.Steps[i].Ingredients[j].ID = identifiers.New()
			cloneInput.Steps[i].Ingredients[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].Instruments {
			if index, ok := instrumentProductIndicies[x.Steps[i].Instruments[j].ID]; ok {
				cloneInput.Steps[i].Instruments[j].ProductOfRecipeStepIndex = pointers.Pointer(uint64(index))
			}
			cloneInput.Steps[i].Instruments[j].ID = identifiers.New()
			cloneInput.Steps[i].Instruments[j].BelongsToRecipeStep = newRecipeStepID
		}
		for j := range cloneInput.Steps[i].Vessels {
			if index, ok := vesselProductIndicies[x.Steps[i].Vessels[j].ID]; ok {
				cloneInput.Steps[i].Vessels[j].ProductOfRecipeStepIndex = pointers.Pointer(uint64(index))
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

	created, err := s.recipeDataManager.CreateRecipe(ctx, cloneInput)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "cloning recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, created)
}
