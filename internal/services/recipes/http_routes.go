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

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	providedInput := new(types.RecipeCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	input, err := converters.ConvertRecipeCreationRequestInputToRecipeDatabaseCreationInput(providedInput)
	if err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	input.CreatedByUser = sessionCtxData.Requester.UserID
	tracing.AttachRecipeIDToSpan(span, input.ID)

	recipe, err := s.recipeDataManager.CreateRecipe(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
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

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ListHandler is our list route.
func (s *service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy)

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

	recipes, err := s.recipeDataManager.GetRecipes(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipes = &types.QueryFilteredResult[types.Recipe]{Data: []*types.Recipe{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipes")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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
	tracing.AttachSearchQueryToSpan(span, query)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy).
		WithValue(keys.SearchQueryKey, query).
		WithValue("using_database", useDB)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
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
			s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// check for parsed input attached to session context data.
	input := new(types.RecipeUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		logger.Error(err, "error encountered decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.Error(err, "provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// fetch recipe from database.
	recipe, err := s.recipeDataManager.GetRecipeByIDAndUser(ctx, recipeID, sessionCtxData.Requester.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the recipe.
	recipe.Update(input)

	if err = s.recipeDataManager.UpdateRecipe(ctx, recipe); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	exists, existenceCheckErr := s.recipeDataManager.RecipeExists(ctx, recipeID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking recipe existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.recipeDataManager.ArchiveRecipe(ctx, recipeID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
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

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
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

	// we deliberately call this with fake data because
	stepInputs, err := s.recipeAnalyzer.GenerateMealPlanTasksForRecipe(ctx, "", x)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating DAG for recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
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
			s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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
