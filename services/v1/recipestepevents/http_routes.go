package recipestepevents

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to recipe step event IDs with.
	URIParamKey = "recipeStepEventID"
)

// ListHandler is our list route.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// ensure query filter.
	filter := models.ExtractQueryFilter(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue("recipe_id", recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue("recipe_step_id", recipeStepID)

	// fetch recipe step events from database.
	recipeStepEvents, err := s.recipeStepEventDataManager.GetRecipeStepEvents(ctx, recipeID, recipeStepID, filter)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		recipeStepEvents = &models.RecipeStepEventList{
			RecipeStepEvents: []models.RecipeStepEvent{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching recipe step events")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, recipeStepEvents); err != nil {
		logger.Error(err, "encoding response")
	}
}

// CreateHandler is our recipe step event creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check request context for parsed input struct.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*models.RecipeStepEventCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine user ID.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	logger = logger.WithValue("recipe_id", recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	recipeExists, err := s.recipeDataManager.RecipeExists(ctx, recipeID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking recipe existence")
		res.WriteHeader(http.StatusInternalServerError)
		return
	} else if !recipeExists {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	logger = logger.WithValue("recipe_step_id", recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	input.BelongsToRecipeStep = recipeStepID

	recipeStepExists, err := s.recipeStepDataManager.RecipeStepExists(ctx, recipeID, recipeStepID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking recipe step existence")
		res.WriteHeader(http.StatusInternalServerError)
		return
	} else if !recipeStepExists {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// create recipe step event in database.
	x, err := s.recipeStepEventDataManager.CreateRecipeStepEvent(ctx, input)
	if err != nil {
		logger.Error(err, "error creating recipe step event")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracing.AttachRecipeStepEventIDToSpan(span, x.ID)
	logger = logger.WithValue("recipe_step_event_id", x.ID)

	// notify relevant parties.
	s.recipeStepEventCounter.Increment(ctx)
	s.reporter.Report(newsman.Event{
		Data:      x,
		Topics:    []string{topicName},
		EventType: string(models.Create),
	})

	// encode our response and peace.
	res.WriteHeader(http.StatusCreated)
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ExistenceHandler returns a HEAD handler that returns 200 if a recipe step event exists, 404 otherwise.
func (s *Service) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue("recipe_id", recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue("recipe_step_id", recipeStepID)

	// determine recipe step event ID.
	recipeStepEventID := s.recipeStepEventIDFetcher(req)
	tracing.AttachRecipeStepEventIDToSpan(span, recipeStepEventID)
	logger = logger.WithValue("recipe_step_event_id", recipeStepEventID)

	// fetch recipe step event from database.
	exists, err := s.recipeStepEventDataManager.RecipeStepEventExists(ctx, recipeID, recipeStepID, recipeStepEventID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking recipe step event existence in database")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if exists {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

// ReadHandler returns a GET handler that returns a recipe step event.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue("recipe_id", recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue("recipe_step_id", recipeStepID)

	// determine recipe step event ID.
	recipeStepEventID := s.recipeStepEventIDFetcher(req)
	tracing.AttachRecipeStepEventIDToSpan(span, recipeStepEventID)
	logger = logger.WithValue("recipe_step_event_id", recipeStepEventID)

	// fetch recipe step event from database.
	x, err := s.recipeStepEventDataManager.GetRecipeStepEvent(ctx, recipeID, recipeStepID, recipeStepEventID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching recipe step event from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// UpdateHandler returns a handler that updates a recipe step event.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check for parsed input attached to request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*models.RecipeStepEventUpdateInput)
	if !ok {
		logger.Info("no input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine user ID.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	logger = logger.WithValue("recipe_id", recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	logger = logger.WithValue("recipe_step_id", recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	input.BelongsToRecipeStep = recipeStepID

	// determine recipe step event ID.
	recipeStepEventID := s.recipeStepEventIDFetcher(req)
	logger = logger.WithValue("recipe_step_event_id", recipeStepEventID)
	tracing.AttachRecipeStepEventIDToSpan(span, recipeStepEventID)

	// fetch recipe step event from database.
	x, err := s.recipeStepEventDataManager.GetRecipeStepEvent(ctx, recipeID, recipeStepID, recipeStepEventID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting recipe step event")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the data structure.
	x.Update(input)

	// update recipe step event in database.
	if err = s.recipeStepEventDataManager.UpdateRecipeStepEvent(ctx, x); err != nil {
		logger.Error(err, "error encountered updating recipe step event")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.reporter.Report(newsman.Event{
		Data:      x,
		Topics:    []string{topicName},
		EventType: string(models.Update),
	})

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ArchiveHandler returns a handler that archives a recipe step event.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	logger = logger.WithValue("recipe_id", recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	recipeExists, err := s.recipeDataManager.RecipeExists(ctx, recipeID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking recipe existence")
		res.WriteHeader(http.StatusInternalServerError)
		return
	} else if !recipeExists {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	logger = logger.WithValue("recipe_step_id", recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	recipeStepExists, err := s.recipeStepDataManager.RecipeStepExists(ctx, recipeID, recipeStepID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking recipe step existence")
		res.WriteHeader(http.StatusInternalServerError)
		return
	} else if !recipeStepExists {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// determine recipe step event ID.
	recipeStepEventID := s.recipeStepEventIDFetcher(req)
	logger = logger.WithValue("recipe_step_event_id", recipeStepEventID)
	tracing.AttachRecipeStepEventIDToSpan(span, recipeStepEventID)

	// archive the recipe step event in the database.
	err = s.recipeStepEventDataManager.ArchiveRecipeStepEvent(ctx, recipeStepID, recipeStepEventID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting recipe step event")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.recipeStepEventCounter.Decrement(ctx)
	s.reporter.Report(newsman.Event{
		EventType: string(models.Archive),
		Data:      &models.RecipeStepEvent{ID: recipeStepEventID},
		Topics:    []string{topicName},
	})

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
