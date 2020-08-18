package validingredientpreparations

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to valid ingredient preparation IDs with.
	URIParamKey = "validIngredientPreparationID"
)

// ListHandler is our list route.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// ensure query filter.
	filter := models.ExtractQueryFilter(req)

	// fetch valid ingredient preparations from database.
	validIngredientPreparations, err := s.validIngredientPreparationDataManager.GetValidIngredientPreparations(ctx, filter)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		validIngredientPreparations = &models.ValidIngredientPreparationList{
			ValidIngredientPreparations: []models.ValidIngredientPreparation{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching valid ingredient preparations")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, validIngredientPreparations); err != nil {
		logger.Error(err, "encoding response")
	}
}

// CreateHandler is our valid ingredient preparation creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check request context for parsed input struct.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*models.ValidIngredientPreparationCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// create valid ingredient preparation in database.
	x, err := s.validIngredientPreparationDataManager.CreateValidIngredientPreparation(ctx, input)
	if err != nil {
		logger.Error(err, "error creating valid ingredient preparation")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracing.AttachValidIngredientPreparationIDToSpan(span, x.ID)
	logger = logger.WithValue("valid_ingredient_preparation_id", x.ID)

	// notify relevant parties.
	s.validIngredientPreparationCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if a valid ingredient preparation exists, 404 otherwise.
func (s *Service) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid ingredient preparation ID.
	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)
	logger = logger.WithValue("valid_ingredient_preparation_id", validIngredientPreparationID)

	// fetch valid ingredient preparation from database.
	exists, err := s.validIngredientPreparationDataManager.ValidIngredientPreparationExists(ctx, validIngredientPreparationID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking valid ingredient preparation existence in database")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if exists {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

// ReadHandler returns a GET handler that returns a valid ingredient preparation.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid ingredient preparation ID.
	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)
	logger = logger.WithValue("valid_ingredient_preparation_id", validIngredientPreparationID)

	// fetch valid ingredient preparation from database.
	x, err := s.validIngredientPreparationDataManager.GetValidIngredientPreparation(ctx, validIngredientPreparationID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching valid ingredient preparation from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// UpdateHandler returns a handler that updates a valid ingredient preparation.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check for parsed input attached to request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*models.ValidIngredientPreparationUpdateInput)
	if !ok {
		logger.Info("no input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine valid ingredient preparation ID.
	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	logger = logger.WithValue("valid_ingredient_preparation_id", validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	// fetch valid ingredient preparation from database.
	x, err := s.validIngredientPreparationDataManager.GetValidIngredientPreparation(ctx, validIngredientPreparationID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting valid ingredient preparation")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the data structure.
	x.Update(input)

	// update valid ingredient preparation in database.
	if err = s.validIngredientPreparationDataManager.UpdateValidIngredientPreparation(ctx, x); err != nil {
		logger.Error(err, "error encountered updating valid ingredient preparation")
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

// ArchiveHandler returns a handler that archives a valid ingredient preparation.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid ingredient preparation ID.
	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	logger = logger.WithValue("valid_ingredient_preparation_id", validIngredientPreparationID)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)

	// archive the valid ingredient preparation in the database.
	err = s.validIngredientPreparationDataManager.ArchiveValidIngredientPreparation(ctx, validIngredientPreparationID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting valid ingredient preparation")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.validIngredientPreparationCounter.Decrement(ctx)
	s.reporter.Report(newsman.Event{
		EventType: string(models.Archive),
		Data:      &models.ValidIngredientPreparation{ID: validIngredientPreparationID},
		Topics:    []string{topicName},
	})

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
