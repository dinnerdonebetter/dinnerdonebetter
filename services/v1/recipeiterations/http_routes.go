package recipeiterations

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to recipe iteration IDs with
	URIParamKey = "recipeIterationID"
)

func attachRecipeIterationIDToSpan(span *trace.Span, recipeIterationID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_iteration_id", strconv.FormatUint(recipeIterationID, 10)))
	}
}

func attachUserIDToSpan(span *trace.Span, userID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("user_id", strconv.FormatUint(userID, 10)))
	}
}

// ListHandler is our list route
func (s *Service) ListHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ListHandler")
		defer span.End()

		// ensure query filter
		qf := models.ExtractQueryFilter(req)

		// determine user ID
		userID := s.userIDFetcher(req)
		logger := s.logger.WithValue("user_id", userID)
		attachUserIDToSpan(span, userID)

		// fetch recipe iterations from database
		recipeIterations, err := s.recipeIterationDatabase.GetRecipeIterations(ctx, qf, userID)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			recipeIterations = &models.RecipeIterationList{
				RecipeIterations: []models.RecipeIteration{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching recipe iterations")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, recipeIterations); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our recipe iteration creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.RecipeIterationCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)
		input.BelongsTo = userID

		// create recipe iteration in database
		x, err := s.recipeIterationDatabase.CreateRecipeIteration(ctx, input)
		if err != nil {
			logger.Error(err, "error creating recipe iteration")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.recipeIterationCounter.Increment(ctx)
		attachRecipeIterationIDToSpan(span, x.ID)
		s.reporter.Report(newsman.Event{
			Data:      x,
			Topics:    []string{topicName},
			EventType: string(models.Create),
		})

		// encode our response and peace
		res.WriteHeader(http.StatusCreated)
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// ReadHandler returns a GET handler that returns a recipe iteration
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeIterationID := s.recipeIterationIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":             userID,
			"recipe_iteration_id": recipeIterationID,
		})
		attachRecipeIterationIDToSpan(span, recipeIterationID)
		attachUserIDToSpan(span, userID)

		// fetch recipe iteration from database
		x, err := s.recipeIterationDatabase.GetRecipeIteration(ctx, recipeIterationID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching recipe iteration from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a recipe iteration
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.RecipeIterationUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeIterationID := s.recipeIterationIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":             userID,
			"recipe_iteration_id": recipeIterationID,
		})
		attachRecipeIterationIDToSpan(span, recipeIterationID)
		attachUserIDToSpan(span, userID)

		// fetch recipe iteration from database
		x, err := s.recipeIterationDatabase.GetRecipeIteration(ctx, recipeIterationID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting recipe iteration")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update recipe iteration in database
		if err = s.recipeIterationDatabase.UpdateRecipeIteration(ctx, x); err != nil {
			logger.Error(err, "error encountered updating recipe iteration")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.reporter.Report(newsman.Event{
			Data:      x,
			Topics:    []string{topicName},
			EventType: string(models.Update),
		})

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// ArchiveHandler returns a handler that archives a recipe iteration
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeIterationID := s.recipeIterationIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"recipe_iteration_id": recipeIterationID,
			"user_id":             userID,
		})
		attachRecipeIterationIDToSpan(span, recipeIterationID)
		attachUserIDToSpan(span, userID)

		// archive the recipe iteration in the database
		err := s.recipeIterationDatabase.ArchiveRecipeIteration(ctx, recipeIterationID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting recipe iteration")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.recipeIterationCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.RecipeIteration{ID: recipeIterationID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
