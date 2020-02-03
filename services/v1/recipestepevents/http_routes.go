package recipestepevents

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to recipe step event IDs with
	URIParamKey = "recipeStepEventID"
)

func attachRecipeStepEventIDToSpan(span *trace.Span, recipeStepEventID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_step_event_id", strconv.FormatUint(recipeStepEventID, 10)))
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

		// fetch recipe step events from database
		recipeStepEvents, err := s.recipeStepEventDatabase.GetRecipeStepEvents(ctx, qf, userID)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			recipeStepEvents = &models.RecipeStepEventList{
				RecipeStepEvents: []models.RecipeStepEvent{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching recipe step events")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, recipeStepEvents); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our recipe step event creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.RecipeStepEventCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)
		input.BelongsTo = userID

		// create recipe step event in database
		x, err := s.recipeStepEventDatabase.CreateRecipeStepEvent(ctx, input)
		if err != nil {
			logger.Error(err, "error creating recipe step event")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.recipeStepEventCounter.Increment(ctx)
		attachRecipeStepEventIDToSpan(span, x.ID)
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

// ReadHandler returns a GET handler that returns a recipe step event
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeStepEventID := s.recipeStepEventIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":              userID,
			"recipe_step_event_id": recipeStepEventID,
		})
		attachRecipeStepEventIDToSpan(span, recipeStepEventID)
		attachUserIDToSpan(span, userID)

		// fetch recipe step event from database
		x, err := s.recipeStepEventDatabase.GetRecipeStepEvent(ctx, recipeStepEventID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching recipe step event from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a recipe step event
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.RecipeStepEventUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeStepEventID := s.recipeStepEventIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":              userID,
			"recipe_step_event_id": recipeStepEventID,
		})
		attachRecipeStepEventIDToSpan(span, recipeStepEventID)
		attachUserIDToSpan(span, userID)

		// fetch recipe step event from database
		x, err := s.recipeStepEventDatabase.GetRecipeStepEvent(ctx, recipeStepEventID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting recipe step event")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update recipe step event in database
		if err = s.recipeStepEventDatabase.UpdateRecipeStepEvent(ctx, x); err != nil {
			logger.Error(err, "error encountered updating recipe step event")
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

// ArchiveHandler returns a handler that archives a recipe step event
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeStepEventID := s.recipeStepEventIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"recipe_step_event_id": recipeStepEventID,
			"user_id":              userID,
		})
		attachRecipeStepEventIDToSpan(span, recipeStepEventID)
		attachUserIDToSpan(span, userID)

		// archive the recipe step event in the database
		err := s.recipeStepEventDatabase.ArchiveRecipeStepEvent(ctx, recipeStepEventID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting recipe step event")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.recipeStepEventCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.RecipeStepEvent{ID: recipeStepEventID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
