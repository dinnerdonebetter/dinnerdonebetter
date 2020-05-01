package recipetags

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to recipe tag IDs with.
	URIParamKey = "recipeTagID"
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

		// fetch recipe tags from database.
		recipeTags, err := s.recipeTagDataManager.GetRecipeTags(ctx, recipeID, filter)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list.
			recipeTags = &models.RecipeTagList{
				RecipeTags: []models.RecipeTag{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching recipe tags")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, recipeTags); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our recipe tag creation route.
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check request context for parsed input struct.
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.RecipeTagCreationInput)
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

		// create recipe tag in database.
		x, err := s.recipeTagDataManager.CreateRecipeTag(ctx, input)
		if err != nil {
			logger.Error(err, "error creating recipe tag")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tracing.AttachRecipeTagIDToSpan(span, x.ID)
		logger = logger.WithValue("recipe_tag_id", x.ID)

		// notify relevant parties.
		s.recipeTagCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if a recipe tag exists, 404 otherwise.
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

		// determine recipe tag ID.
		recipeTagID := s.recipeTagIDFetcher(req)
		tracing.AttachRecipeTagIDToSpan(span, recipeTagID)
		logger = logger.WithValue("recipe_tag_id", recipeTagID)

		// fetch recipe tag from database.
		exists, err := s.recipeTagDataManager.RecipeTagExists(ctx, recipeID, recipeTagID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking recipe tag existence in database")
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

// ReadHandler returns a GET handler that returns a recipe tag.
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

		// determine recipe tag ID.
		recipeTagID := s.recipeTagIDFetcher(req)
		tracing.AttachRecipeTagIDToSpan(span, recipeTagID)
		logger = logger.WithValue("recipe_tag_id", recipeTagID)

		// fetch recipe tag from database.
		x, err := s.recipeTagDataManager.GetRecipeTag(ctx, recipeID, recipeTagID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching recipe tag from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a recipe tag.
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check for parsed input attached to request context.
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.RecipeTagUpdateInput)
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

		// determine recipe tag ID.
		recipeTagID := s.recipeTagIDFetcher(req)
		logger = logger.WithValue("recipe_tag_id", recipeTagID)
		tracing.AttachRecipeTagIDToSpan(span, recipeTagID)

		// fetch recipe tag from database.
		x, err := s.recipeTagDataManager.GetRecipeTag(ctx, recipeID, recipeTagID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting recipe tag")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure.
		x.Update(input)

		// update recipe tag in database.
		if err = s.recipeTagDataManager.UpdateRecipeTag(ctx, x); err != nil {
			logger.Error(err, "error encountered updating recipe tag")
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

// ArchiveHandler returns a handler that archives a recipe tag.
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

		// determine recipe tag ID.
		recipeTagID := s.recipeTagIDFetcher(req)
		logger = logger.WithValue("recipe_tag_id", recipeTagID)
		tracing.AttachRecipeTagIDToSpan(span, recipeTagID)

		// archive the recipe tag in the database.
		err = s.recipeTagDataManager.ArchiveRecipeTag(ctx, recipeID, recipeTagID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting recipe tag")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties.
		s.recipeTagCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.RecipeTag{ID: recipeTagID},
			Topics:    []string{topicName},
		})

		// encode our response and peace.
		res.WriteHeader(http.StatusNoContent)
	}
}
