package validingredients

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to valid ingredient IDs with.
	URIParamKey = "validIngredientID"
)

// ListHandler is our list route.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// ensure query filter.
	filter := models.ExtractQueryFilter(req)

	// fetch valid ingredients from database.
	validIngredients, err := s.validIngredientDataManager.GetValidIngredients(ctx, filter)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		validIngredients = &models.ValidIngredientList{
			ValidIngredients: []models.ValidIngredient{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching valid ingredients")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, validIngredients); err != nil {
		logger.Error(err, "encoding response")
	}
}

// SearchHandler is our search route.
func (s *Service) SearchHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "SearchHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// we only parse the filter here because it will contain the limit
	filter := models.ExtractQueryFilter(req)
	query := req.URL.Query().Get(models.SearchQueryKey)
	logger = logger.WithValue("search_query", query)

	relevantIDs, searchErr := s.search.Search(ctx, query)
	if searchErr != nil {
		logger.Error(searchErr, "error encountered executing search query")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// fetch valid ingredients from database.
	validIngredients, err := s.validIngredientDataManager.GetValidIngredientsWithIDs(ctx, filter.Limit, relevantIDs)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		validIngredients = []models.ValidIngredient{}
	} else if err != nil {
		logger.Error(err, "error encountered fetching valid ingredients")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, validIngredients); err != nil {
		logger.Error(err, "encoding response")
	}
}

// CreateHandler is our valid ingredient creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check request context for parsed input struct.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*models.ValidIngredientCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// create valid ingredient in database.
	x, err := s.validIngredientDataManager.CreateValidIngredient(ctx, input)
	if err != nil {
		logger.Error(err, "error creating valid ingredient")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracing.AttachValidIngredientIDToSpan(span, x.ID)
	logger = logger.WithValue("valid_ingredient_id", x.ID)

	// notify relevant parties.
	s.validIngredientCounter.Increment(ctx)
	s.reporter.Report(newsman.Event{
		Data:      x,
		Topics:    []string{topicName},
		EventType: string(models.Create),
	})
	if searchIndexErr := s.search.Index(ctx, x.ID, x); searchIndexErr != nil {
		logger.Error(searchIndexErr, "adding valid ingredient to search index")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusCreated)
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ExistenceHandler returns a HEAD handler that returns 200 if a valid ingredient exists, 404 otherwise.
func (s *Service) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	logger = logger.WithValue("valid_ingredient_id", validIngredientID)

	// fetch valid ingredient from database.
	exists, err := s.validIngredientDataManager.ValidIngredientExists(ctx, validIngredientID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking valid ingredient existence in database")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if exists {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

// ReadHandler returns a GET handler that returns a valid ingredient.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	logger = logger.WithValue("valid_ingredient_id", validIngredientID)

	// fetch valid ingredient from database.
	x, err := s.validIngredientDataManager.GetValidIngredient(ctx, validIngredientID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching valid ingredient from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// UpdateHandler returns a handler that updates a valid ingredient.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check for parsed input attached to request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*models.ValidIngredientUpdateInput)
	if !ok {
		logger.Info("no input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	logger = logger.WithValue("valid_ingredient_id", validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	// fetch valid ingredient from database.
	x, err := s.validIngredientDataManager.GetValidIngredient(ctx, validIngredientID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting valid ingredient")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the data structure.
	x.Update(input)

	// update valid ingredient in database.
	if err = s.validIngredientDataManager.UpdateValidIngredient(ctx, x); err != nil {
		logger.Error(err, "error encountered updating valid ingredient")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.reporter.Report(newsman.Event{
		Data:      x,
		Topics:    []string{topicName},
		EventType: string(models.Update),
	})
	if searchIndexErr := s.search.Index(ctx, x.ID, x); searchIndexErr != nil {
		logger.Error(searchIndexErr, "updating valid ingredient in search index")
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ArchiveHandler returns a handler that archives a valid ingredient.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	logger = logger.WithValue("valid_ingredient_id", validIngredientID)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

	// archive the valid ingredient in the database.
	err = s.validIngredientDataManager.ArchiveValidIngredient(ctx, validIngredientID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting valid ingredient")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.validIngredientCounter.Decrement(ctx)
	s.reporter.Report(newsman.Event{
		EventType: string(models.Archive),
		Data:      &models.ValidIngredient{ID: validIngredientID},
		Topics:    []string{topicName},
	})
	if indexDeleteErr := s.search.Delete(ctx, validIngredientID); indexDeleteErr != nil {
		logger.Error(indexDeleteErr, "error removing valid ingredient from search index")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
