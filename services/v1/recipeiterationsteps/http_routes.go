package recipeiterationsteps

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to recipe iteration step IDs with.
	URIParamKey = "recipeIterationStepID"
)

// ListHandler is our list route.
func (s *Service) ListHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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

		// fetch recipe iteration steps from database.
		recipeIterationSteps, err := s.recipeIterationStepDataManager.GetRecipeIterationSteps(ctx, recipeID, filter)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list.
			recipeIterationSteps = &models.RecipeIterationStepList{
				RecipeIterationSteps: []models.RecipeIterationStep{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching recipe iteration steps")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, recipeIterationSteps); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our recipe iteration step creation route.
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check request context for parsed input struct.
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.RecipeIterationStepCreationInput)
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

		// create recipe iteration step in database.
		x, err := s.recipeIterationStepDataManager.CreateRecipeIterationStep(ctx, input)
		if err != nil {
			logger.Error(err, "error creating recipe iteration step")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tracing.AttachRecipeIterationStepIDToSpan(span, x.ID)
		logger = logger.WithValue("recipe_iteration_step_id", x.ID)

		// notify relevant parties.
		s.recipeIterationStepCounter.Increment(ctx)
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
}

// ExistenceHandler returns a HEAD handler that returns 200 if a recipe iteration step exists, 404 otherwise.
func (s *Service) ExistenceHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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

		// determine recipe iteration step ID.
		recipeIterationStepID := s.recipeIterationStepIDFetcher(req)
		tracing.AttachRecipeIterationStepIDToSpan(span, recipeIterationStepID)
		logger = logger.WithValue("recipe_iteration_step_id", recipeIterationStepID)

		// fetch recipe iteration step from database.
		exists, err := s.recipeIterationStepDataManager.RecipeIterationStepExists(ctx, recipeID, recipeIterationStepID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking recipe iteration step existence in database")
			res.WriteHeader(http.StatusNotFound)
			return
		}

		if exists {
			res.WriteHeader(http.StatusOK)
		} else {
			res.WriteHeader(http.StatusNotFound)
		}
	}
}

// ReadHandler returns a GET handler that returns a recipe iteration step.
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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

		// determine recipe iteration step ID.
		recipeIterationStepID := s.recipeIterationStepIDFetcher(req)
		tracing.AttachRecipeIterationStepIDToSpan(span, recipeIterationStepID)
		logger = logger.WithValue("recipe_iteration_step_id", recipeIterationStepID)

		// fetch recipe iteration step from database.
		x, err := s.recipeIterationStepDataManager.GetRecipeIterationStep(ctx, recipeID, recipeIterationStepID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching recipe iteration step from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a recipe iteration step.
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check for parsed input attached to request context.
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.RecipeIterationStepUpdateInput)
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

		// determine recipe iteration step ID.
		recipeIterationStepID := s.recipeIterationStepIDFetcher(req)
		logger = logger.WithValue("recipe_iteration_step_id", recipeIterationStepID)
		tracing.AttachRecipeIterationStepIDToSpan(span, recipeIterationStepID)

		// fetch recipe iteration step from database.
		x, err := s.recipeIterationStepDataManager.GetRecipeIterationStep(ctx, recipeID, recipeIterationStepID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting recipe iteration step")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure.
		x.Update(input)

		// update recipe iteration step in database.
		if err = s.recipeIterationStepDataManager.UpdateRecipeIterationStep(ctx, x); err != nil {
			logger.Error(err, "error encountered updating recipe iteration step")
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
}

// ArchiveHandler returns a handler that archives a recipe iteration step.
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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

		// determine recipe iteration step ID.
		recipeIterationStepID := s.recipeIterationStepIDFetcher(req)
		logger = logger.WithValue("recipe_iteration_step_id", recipeIterationStepID)
		tracing.AttachRecipeIterationStepIDToSpan(span, recipeIterationStepID)

		// archive the recipe iteration step in the database.
		err = s.recipeIterationStepDataManager.ArchiveRecipeIterationStep(ctx, recipeID, recipeIterationStepID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting recipe iteration step")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties.
		s.recipeIterationStepCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.RecipeIterationStep{ID: recipeIterationStepID},
			Topics:    []string{topicName},
		})

		// encode our response and peace.
		res.WriteHeader(http.StatusNoContent)
	}
}
