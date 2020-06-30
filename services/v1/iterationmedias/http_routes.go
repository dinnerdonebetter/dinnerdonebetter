package iterationmedias

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to iteration media IDs with.
	URIParamKey = "iterationMediaID"
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

		// determine recipe iteration ID.
		recipeIterationID := s.recipeIterationIDFetcher(req)
		tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)
		logger = logger.WithValue("recipe_iteration_id", recipeIterationID)

		// fetch iteration medias from database.
		iterationMedias, err := s.iterationMediaDataManager.GetIterationMedias(ctx, recipeID, recipeIterationID, filter)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list.
			iterationMedias = &models.IterationMediaList{
				IterationMedia: []models.IterationMedia{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching iteration medias")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, iterationMedias); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our iteration media creation route.
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check request context for parsed input struct.
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.IterationMediaCreationInput)
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

		// determine recipe iteration ID.
		recipeIterationID := s.recipeIterationIDFetcher(req)
		logger = logger.WithValue("recipe_iteration_id", recipeIterationID)
		tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)

		input.BelongsToRecipeIteration = recipeIterationID

		recipeIterationExists, err := s.recipeIterationDataManager.RecipeIterationExists(ctx, recipeID, recipeIterationID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking recipe iteration existence")
			res.WriteHeader(http.StatusInternalServerError)
			return
		} else if !recipeIterationExists {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		// create iteration media in database.
		x, err := s.iterationMediaDataManager.CreateIterationMedia(ctx, input)
		if err != nil {
			logger.Error(err, "error creating iteration media")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tracing.AttachIterationMediaIDToSpan(span, x.ID)
		logger = logger.WithValue("iteration_media_id", x.ID)

		// notify relevant parties.
		s.iterationMediaCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if an iteration media exists, 404 otherwise.
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

		// determine recipe iteration ID.
		recipeIterationID := s.recipeIterationIDFetcher(req)
		tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)
		logger = logger.WithValue("recipe_iteration_id", recipeIterationID)

		// determine iteration media ID.
		iterationMediaID := s.iterationMediaIDFetcher(req)
		tracing.AttachIterationMediaIDToSpan(span, iterationMediaID)
		logger = logger.WithValue("iteration_media_id", iterationMediaID)

		// fetch iteration media from database.
		exists, err := s.iterationMediaDataManager.IterationMediaExists(ctx, recipeID, recipeIterationID, iterationMediaID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking iteration media existence in database")
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

// ReadHandler returns a GET handler that returns an iteration media.
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

		// determine recipe iteration ID.
		recipeIterationID := s.recipeIterationIDFetcher(req)
		tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)
		logger = logger.WithValue("recipe_iteration_id", recipeIterationID)

		// determine iteration media ID.
		iterationMediaID := s.iterationMediaIDFetcher(req)
		tracing.AttachIterationMediaIDToSpan(span, iterationMediaID)
		logger = logger.WithValue("iteration_media_id", iterationMediaID)

		// fetch iteration media from database.
		x, err := s.iterationMediaDataManager.GetIterationMedia(ctx, recipeID, recipeIterationID, iterationMediaID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching iteration media from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates an iteration media.
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check for parsed input attached to request context.
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.IterationMediaUpdateInput)
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

		// determine recipe iteration ID.
		recipeIterationID := s.recipeIterationIDFetcher(req)
		logger = logger.WithValue("recipe_iteration_id", recipeIterationID)
		tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)

		input.BelongsToRecipeIteration = recipeIterationID

		// determine iteration media ID.
		iterationMediaID := s.iterationMediaIDFetcher(req)
		logger = logger.WithValue("iteration_media_id", iterationMediaID)
		tracing.AttachIterationMediaIDToSpan(span, iterationMediaID)

		// fetch iteration media from database.
		x, err := s.iterationMediaDataManager.GetIterationMedia(ctx, recipeID, recipeIterationID, iterationMediaID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting iteration media")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure.
		x.Update(input)

		// update iteration media in database.
		if err = s.iterationMediaDataManager.UpdateIterationMedia(ctx, x); err != nil {
			logger.Error(err, "error encountered updating iteration media")
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

// ArchiveHandler returns a handler that archives an iteration media.
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

		// determine recipe iteration ID.
		recipeIterationID := s.recipeIterationIDFetcher(req)
		logger = logger.WithValue("recipe_iteration_id", recipeIterationID)
		tracing.AttachRecipeIterationIDToSpan(span, recipeIterationID)

		recipeIterationExists, err := s.recipeIterationDataManager.RecipeIterationExists(ctx, recipeID, recipeIterationID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking recipe iteration existence")
			res.WriteHeader(http.StatusInternalServerError)
			return
		} else if !recipeIterationExists {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		// determine iteration media ID.
		iterationMediaID := s.iterationMediaIDFetcher(req)
		logger = logger.WithValue("iteration_media_id", iterationMediaID)
		tracing.AttachIterationMediaIDToSpan(span, iterationMediaID)

		// archive the iteration media in the database.
		err = s.iterationMediaDataManager.ArchiveIterationMedia(ctx, recipeIterationID, iterationMediaID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting iteration media")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties.
		s.iterationMediaCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.IterationMedia{ID: iterationMediaID},
			Topics:    []string{topicName},
		})

		// encode our response and peace.
		res.WriteHeader(http.StatusNoContent)
	}
}
