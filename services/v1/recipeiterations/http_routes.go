package recipeiterations

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to recipe iteration IDs with.
	URIParamKey = "recipeIterationID"
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

	// fetch recipe iterations from database.
	recipeIterations, err := s.recipeIterationDataManager.GetRecipeIterations(ctx, recipeID, filter)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		recipeIterations = &models.RecipeIterationList{
			RecipeIterations: []models.RecipeIteration{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching recipe iterations")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, recipeIterations); err != nil {
		logger.Error(err, "encoding response")
	}
}

// CreateHandler is our recipe iteration creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check request context for parsed input struct.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*models.RecipeIterationCreationInput)
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

	input.BelongsToRecipe = recipeID

	recipeExists, err := s.recipeDataManager.RecipeExists(ctx, recipeID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking recipe existence")
		res.WriteHeader(http.StatusInternalServerError)
		return
	} else if !recipeExists {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// create recipe iteration in database.
	x, err := s.recipeIterationDataManager.CreateRecipeIteration(ctx, input)
	if err != nil {
		logger.Error(err, "error creating recipe iteration")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracing.AttachRecipeIterationIDToSpan(span, x.ID)
	logger = logger.WithValue("recipe_iteration_id", x.ID)

	// notify relevant parties.
	s.recipeIterationCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if a recipe iteration exists, 404 otherwise.
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

	// determine recipe iteration ID.
	recipeIterationID := s.recipeIterationIDFetcher(req)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)
	logger = logger.WithValue("recipe_iteration_id", recipeIterationID)

	// fetch recipe iteration from database.
	exists, err := s.recipeIterationDataManager.RecipeIterationExists(ctx, recipeID, recipeIterationID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking recipe iteration existence in database")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if exists {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

// ReadHandler returns a GET handler that returns a recipe iteration.
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

	// determine recipe iteration ID.
	recipeIterationID := s.recipeIterationIDFetcher(req)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)
	logger = logger.WithValue("recipe_iteration_id", recipeIterationID)

	// fetch recipe iteration from database.
	x, err := s.recipeIterationDataManager.GetRecipeIteration(ctx, recipeID, recipeIterationID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching recipe iteration from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// UpdateHandler returns a handler that updates a recipe iteration.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check for parsed input attached to request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*models.RecipeIterationUpdateInput)
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

	input.BelongsToRecipe = recipeID

	// determine recipe iteration ID.
	recipeIterationID := s.recipeIterationIDFetcher(req)
	logger = logger.WithValue("recipe_iteration_id", recipeIterationID)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)

	// fetch recipe iteration from database.
	x, err := s.recipeIterationDataManager.GetRecipeIteration(ctx, recipeID, recipeIterationID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting recipe iteration")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the data structure.
	x.Update(input)

	// update recipe iteration in database.
	if err = s.recipeIterationDataManager.UpdateRecipeIteration(ctx, x); err != nil {
		logger.Error(err, "error encountered updating recipe iteration")
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

// ArchiveHandler returns a handler that archives a recipe iteration.
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

	// determine recipe iteration ID.
	recipeIterationID := s.recipeIterationIDFetcher(req)
	logger = logger.WithValue("recipe_iteration_id", recipeIterationID)
	tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)

	// archive the recipe iteration in the database.
	err = s.recipeIterationDataManager.ArchiveRecipeIteration(ctx, recipeID, recipeIterationID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting recipe iteration")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.recipeIterationCounter.Decrement(ctx)
	s.reporter.Report(newsman.Event{
		EventType: string(models.Archive),
		Data:      &models.RecipeIteration{ID: recipeIterationID},
		Topics:    []string{topicName},
	})

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
