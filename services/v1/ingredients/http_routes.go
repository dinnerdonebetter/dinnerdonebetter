package ingredients

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to ingredient IDs with
	URIParamKey = "ingredientID"
)

func attachIngredientIDToSpan(span *trace.Span, ingredientID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("ingredient_id", strconv.FormatUint(ingredientID, 10)))
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

		// fetch ingredients from database
		ingredients, err := s.ingredientDatabase.GetIngredients(ctx, qf, userID)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			ingredients = &models.IngredientList{
				Ingredients: []models.Ingredient{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching ingredients")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, ingredients); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our ingredient creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.IngredientCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)

		// create ingredient in database
		x, err := s.ingredientDatabase.CreateIngredient(ctx, input)
		if err != nil {
			logger.Error(err, "error creating ingredient")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.ingredientCounter.Increment(ctx)
		attachIngredientIDToSpan(span, x.ID)
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

// ReadHandler returns a GET handler that returns an ingredient
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		ingredientID := s.ingredientIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":       userID,
			"ingredient_id": ingredientID,
		})
		attachIngredientIDToSpan(span, ingredientID)
		attachUserIDToSpan(span, userID)

		// fetch ingredient from database
		x, err := s.ingredientDatabase.GetIngredient(ctx, ingredientID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching ingredient from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates an ingredient
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.IngredientUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		ingredientID := s.ingredientIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":       userID,
			"ingredient_id": ingredientID,
		})
		attachIngredientIDToSpan(span, ingredientID)
		attachUserIDToSpan(span, userID)

		// fetch ingredient from database
		x, err := s.ingredientDatabase.GetIngredient(ctx, ingredientID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting ingredient")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update ingredient in database
		if err = s.ingredientDatabase.UpdateIngredient(ctx, x); err != nil {
			logger.Error(err, "error encountered updating ingredient")
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

// ArchiveHandler returns a handler that archives an ingredient
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		ingredientID := s.ingredientIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"ingredient_id": ingredientID,
			"user_id":       userID,
		})
		attachIngredientIDToSpan(span, ingredientID)
		attachUserIDToSpan(span, userID)

		// archive the ingredient in the database
		err := s.ingredientDatabase.ArchiveIngredient(ctx, ingredientID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting ingredient")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.ingredientCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.Ingredient{ID: ingredientID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
