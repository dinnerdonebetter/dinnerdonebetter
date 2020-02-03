package recipestepproducts

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to recipe step product IDs with
	URIParamKey = "recipeStepProductID"
)

func attachRecipeStepProductIDToSpan(span *trace.Span, recipeStepProductID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("recipe_step_product_id", strconv.FormatUint(recipeStepProductID, 10)))
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

		// fetch recipe step products from database
		recipeStepProducts, err := s.recipeStepProductDatabase.GetRecipeStepProducts(ctx, qf, userID)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			recipeStepProducts = &models.RecipeStepProductList{
				RecipeStepProducts: []models.RecipeStepProduct{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching recipe step products")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, recipeStepProducts); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our recipe step product creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.RecipeStepProductCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)
		input.BelongsTo = userID

		// create recipe step product in database
		x, err := s.recipeStepProductDatabase.CreateRecipeStepProduct(ctx, input)
		if err != nil {
			logger.Error(err, "error creating recipe step product")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.recipeStepProductCounter.Increment(ctx)
		attachRecipeStepProductIDToSpan(span, x.ID)
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

// ReadHandler returns a GET handler that returns a recipe step product
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeStepProductID := s.recipeStepProductIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":                userID,
			"recipe_step_product_id": recipeStepProductID,
		})
		attachRecipeStepProductIDToSpan(span, recipeStepProductID)
		attachUserIDToSpan(span, userID)

		// fetch recipe step product from database
		x, err := s.recipeStepProductDatabase.GetRecipeStepProduct(ctx, recipeStepProductID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching recipe step product from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a recipe step product
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.RecipeStepProductUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeStepProductID := s.recipeStepProductIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":                userID,
			"recipe_step_product_id": recipeStepProductID,
		})
		attachRecipeStepProductIDToSpan(span, recipeStepProductID)
		attachUserIDToSpan(span, userID)

		// fetch recipe step product from database
		x, err := s.recipeStepProductDatabase.GetRecipeStepProduct(ctx, recipeStepProductID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting recipe step product")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update recipe step product in database
		if err = s.recipeStepProductDatabase.UpdateRecipeStepProduct(ctx, x); err != nil {
			logger.Error(err, "error encountered updating recipe step product")
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

// ArchiveHandler returns a handler that archives a recipe step product
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		recipeStepProductID := s.recipeStepProductIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"recipe_step_product_id": recipeStepProductID,
			"user_id":                userID,
		})
		attachRecipeStepProductIDToSpan(span, recipeStepProductID)
		attachUserIDToSpan(span, userID)

		// archive the recipe step product in the database
		err := s.recipeStepProductDatabase.ArchiveRecipeStepProduct(ctx, recipeStepProductID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting recipe step product")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.recipeStepProductCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.RecipeStepProduct{ID: recipeStepProductID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
