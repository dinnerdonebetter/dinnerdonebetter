package recipesteppreparations

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to recipe step preparation IDs with.
	URIParamKey = "recipeStepPreparationID"
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

		// determine recipe step ID.
		recipeStepID := s.recipeStepIDFetcher(req)
		tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
		logger = logger.WithValue("recipe_step_id", recipeStepID)

		// fetch recipe step preparations from database.
		recipeStepPreparations, err := s.recipeStepPreparationDataManager.GetRecipeStepPreparations(ctx, recipeID, recipeStepID, filter)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list.
			recipeStepPreparations = &models.RecipeStepPreparationList{
				RecipeStepPreparations: []models.RecipeStepPreparation{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching recipe step preparations")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, recipeStepPreparations); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our recipe step preparation creation route.
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check request context for parsed input struct.
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.RecipeStepPreparationCreationInput)
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

		// create recipe step preparation in database.
		x, err := s.recipeStepPreparationDataManager.CreateRecipeStepPreparation(ctx, input)
		if err != nil {
			logger.Error(err, "error creating recipe step preparation")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tracing.AttachRecipeStepPreparationIDToSpan(span, x.ID)
		logger = logger.WithValue("recipe_step_preparation_id", x.ID)

		// notify relevant parties.
		s.recipeStepPreparationCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if a recipe step preparation exists, 404 otherwise.
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

		// determine recipe step ID.
		recipeStepID := s.recipeStepIDFetcher(req)
		tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
		logger = logger.WithValue("recipe_step_id", recipeStepID)

		// determine recipe step preparation ID.
		recipeStepPreparationID := s.recipeStepPreparationIDFetcher(req)
		tracing.AttachRecipeStepPreparationIDToSpan(span, recipeStepPreparationID)
		logger = logger.WithValue("recipe_step_preparation_id", recipeStepPreparationID)

		// fetch recipe step preparation from database.
		exists, err := s.recipeStepPreparationDataManager.RecipeStepPreparationExists(ctx, recipeID, recipeStepID, recipeStepPreparationID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking recipe step preparation existence in database")
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

// ReadHandler returns a GET handler that returns a recipe step preparation.
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

		// determine recipe step ID.
		recipeStepID := s.recipeStepIDFetcher(req)
		tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
		logger = logger.WithValue("recipe_step_id", recipeStepID)

		// determine recipe step preparation ID.
		recipeStepPreparationID := s.recipeStepPreparationIDFetcher(req)
		tracing.AttachRecipeStepPreparationIDToSpan(span, recipeStepPreparationID)
		logger = logger.WithValue("recipe_step_preparation_id", recipeStepPreparationID)

		// fetch recipe step preparation from database.
		x, err := s.recipeStepPreparationDataManager.GetRecipeStepPreparation(ctx, recipeID, recipeStepID, recipeStepPreparationID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching recipe step preparation from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a recipe step preparation.
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check for parsed input attached to request context.
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.RecipeStepPreparationUpdateInput)
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

		// determine recipe step preparation ID.
		recipeStepPreparationID := s.recipeStepPreparationIDFetcher(req)
		logger = logger.WithValue("recipe_step_preparation_id", recipeStepPreparationID)
		tracing.AttachRecipeStepPreparationIDToSpan(span, recipeStepPreparationID)

		// fetch recipe step preparation from database.
		x, err := s.recipeStepPreparationDataManager.GetRecipeStepPreparation(ctx, recipeID, recipeStepID, recipeStepPreparationID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting recipe step preparation")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure.
		x.Update(input)

		// update recipe step preparation in database.
		if err = s.recipeStepPreparationDataManager.UpdateRecipeStepPreparation(ctx, x); err != nil {
			logger.Error(err, "error encountered updating recipe step preparation")
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

// ArchiveHandler returns a handler that archives a recipe step preparation.
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

		// determine recipe step preparation ID.
		recipeStepPreparationID := s.recipeStepPreparationIDFetcher(req)
		logger = logger.WithValue("recipe_step_preparation_id", recipeStepPreparationID)
		tracing.AttachRecipeStepPreparationIDToSpan(span, recipeStepPreparationID)

		// archive the recipe step preparation in the database.
		err = s.recipeStepPreparationDataManager.ArchiveRecipeStepPreparation(ctx, recipeStepID, recipeStepPreparationID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting recipe step preparation")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties.
		s.recipeStepPreparationCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.RecipeStepPreparation{ID: recipeStepPreparationID},
			Topics:    []string{topicName},
		})

		// encode our response and peace.
		res.WriteHeader(http.StatusNoContent)
	}
}
