package recipestepinstruments

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to recipe step instrument IDs with
	URIParamKey = "recipeStepInstrumentID"
)

func attachRecipeStepInstrumentIDToSpan(span *trace.Span, recipeStepInstrumentID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_step_instrument_id", strconv.FormatUint(recipeStepInstrumentID, 10)))
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

		// fetch recipe step instruments from database
		recipeStepInstruments, err := s.recipeStepInstrumentDatabase.GetRecipeStepInstruments(ctx, qf, userID)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			recipeStepInstruments = &models.RecipeStepInstrumentList{
				RecipeStepInstruments: []models.RecipeStepInstrument{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching recipe step instruments")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, recipeStepInstruments); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our recipe step instrument creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.RecipeStepInstrumentCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)
		input.BelongsTo = userID

		// create recipe step instrument in database
		x, err := s.recipeStepInstrumentDatabase.CreateRecipeStepInstrument(ctx, input)
		if err != nil {
			logger.Error(err, "error creating recipe step instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.recipeStepInstrumentCounter.Increment(ctx)
		attachRecipeStepInstrumentIDToSpan(span, x.ID)
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

// ReadHandler returns a GET handler that returns a recipe step instrument
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeStepInstrumentID := s.recipeStepInstrumentIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":                   userID,
			"recipe_step_instrument_id": recipeStepInstrumentID,
		})
		attachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)
		attachUserIDToSpan(span, userID)

		// fetch recipe step instrument from database
		x, err := s.recipeStepInstrumentDatabase.GetRecipeStepInstrument(ctx, recipeStepInstrumentID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching recipe step instrument from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a recipe step instrument
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.RecipeStepInstrumentUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeStepInstrumentID := s.recipeStepInstrumentIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":                   userID,
			"recipe_step_instrument_id": recipeStepInstrumentID,
		})
		attachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)
		attachUserIDToSpan(span, userID)

		// fetch recipe step instrument from database
		x, err := s.recipeStepInstrumentDatabase.GetRecipeStepInstrument(ctx, recipeStepInstrumentID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting recipe step instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update recipe step instrument in database
		if err = s.recipeStepInstrumentDatabase.UpdateRecipeStepInstrument(ctx, x); err != nil {
			logger.Error(err, "error encountered updating recipe step instrument")
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

// ArchiveHandler returns a handler that archives a recipe step instrument
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeStepInstrumentID := s.recipeStepInstrumentIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"recipe_step_instrument_id": recipeStepInstrumentID,
			"user_id":                   userID,
		})
		attachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)
		attachUserIDToSpan(span, userID)

		// archive the recipe step instrument in the database
		err := s.recipeStepInstrumentDatabase.ArchiveRecipeStepInstrument(ctx, recipeStepInstrumentID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting recipe step instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.recipeStepInstrumentCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.RecipeStepInstrument{ID: recipeStepInstrumentID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
