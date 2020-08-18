package validpreparations

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to valid preparation IDs with.
	URIParamKey = "validPreparationID"
)

// ListHandler is our list route.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// ensure query filter.
	filter := models.ExtractQueryFilter(req)

	// fetch valid preparations from database.
	validPreparations, err := s.validPreparationDataManager.GetValidPreparations(ctx, filter)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		validPreparations = &models.ValidPreparationList{
			ValidPreparations: []models.ValidPreparation{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching valid preparations")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, validPreparations); err != nil {
		logger.Error(err, "encoding response")
	}
}

// CreateHandler is our valid preparation creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check request context for parsed input struct.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*models.ValidPreparationCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// create valid preparation in database.
	x, err := s.validPreparationDataManager.CreateValidPreparation(ctx, input)
	if err != nil {
		logger.Error(err, "error creating valid preparation")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracing.AttachValidPreparationIDToSpan(span, x.ID)
	logger = logger.WithValue("valid_preparation_id", x.ID)

	// notify relevant parties.
	s.validPreparationCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if a valid preparation exists, 404 otherwise.
func (s *Service) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
	logger = logger.WithValue("valid_preparation_id", validPreparationID)

	// fetch valid preparation from database.
	exists, err := s.validPreparationDataManager.ValidPreparationExists(ctx, validPreparationID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking valid preparation existence in database")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if exists {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

// ReadHandler returns a GET handler that returns a valid preparation.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
	logger = logger.WithValue("valid_preparation_id", validPreparationID)

	// fetch valid preparation from database.
	x, err := s.validPreparationDataManager.GetValidPreparation(ctx, validPreparationID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching valid preparation from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// UpdateHandler returns a handler that updates a valid preparation.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check for parsed input attached to request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*models.ValidPreparationUpdateInput)
	if !ok {
		logger.Info("no input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	logger = logger.WithValue("valid_preparation_id", validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	// fetch valid preparation from database.
	x, err := s.validPreparationDataManager.GetValidPreparation(ctx, validPreparationID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting valid preparation")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the data structure.
	x.Update(input)

	// update valid preparation in database.
	if err = s.validPreparationDataManager.UpdateValidPreparation(ctx, x); err != nil {
		logger.Error(err, "error encountered updating valid preparation")
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

// ArchiveHandler returns a handler that archives a valid preparation.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	logger = logger.WithValue("valid_preparation_id", validPreparationID)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

	// archive the valid preparation in the database.
	err = s.validPreparationDataManager.ArchiveValidPreparation(ctx, validPreparationID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting valid preparation")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.validPreparationCounter.Decrement(ctx)
	s.reporter.Report(newsman.Event{
		EventType: string(models.Archive),
		Data:      &models.ValidPreparation{ID: validPreparationID},
		Topics:    []string{topicName},
	})

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
