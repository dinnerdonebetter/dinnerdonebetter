package validinstruments

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to valid instrument IDs with.
	URIParamKey = "validInstrumentID"
)

// ListHandler is our list route.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// ensure query filter.
	filter := models.ExtractQueryFilter(req)

	// fetch valid instruments from database.
	validInstruments, err := s.validInstrumentDataManager.GetValidInstruments(ctx, filter)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		validInstruments = &models.ValidInstrumentList{
			ValidInstruments: []models.ValidInstrument{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching valid instruments")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, validInstruments); err != nil {
		logger.Error(err, "encoding response")
	}
}

// CreateHandler is our valid instrument creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check request context for parsed input struct.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*models.ValidInstrumentCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// create valid instrument in database.
	x, err := s.validInstrumentDataManager.CreateValidInstrument(ctx, input)
	if err != nil {
		logger.Error(err, "error creating valid instrument")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracing.AttachValidInstrumentIDToSpan(span, x.ID)
	logger = logger.WithValue("valid_instrument_id", x.ID)

	// notify relevant parties.
	s.validInstrumentCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if a valid instrument exists, 404 otherwise.
func (s *Service) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid instrument ID.
	validInstrumentID := s.validInstrumentIDFetcher(req)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)
	logger = logger.WithValue("valid_instrument_id", validInstrumentID)

	// fetch valid instrument from database.
	exists, err := s.validInstrumentDataManager.ValidInstrumentExists(ctx, validInstrumentID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking valid instrument existence in database")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if exists {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

// ReadHandler returns a GET handler that returns a valid instrument.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid instrument ID.
	validInstrumentID := s.validInstrumentIDFetcher(req)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)
	logger = logger.WithValue("valid_instrument_id", validInstrumentID)

	// fetch valid instrument from database.
	x, err := s.validInstrumentDataManager.GetValidInstrument(ctx, validInstrumentID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching valid instrument from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// UpdateHandler returns a handler that updates a valid instrument.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check for parsed input attached to request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*models.ValidInstrumentUpdateInput)
	if !ok {
		logger.Info("no input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine valid instrument ID.
	validInstrumentID := s.validInstrumentIDFetcher(req)
	logger = logger.WithValue("valid_instrument_id", validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	// fetch valid instrument from database.
	x, err := s.validInstrumentDataManager.GetValidInstrument(ctx, validInstrumentID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting valid instrument")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the data structure.
	x.Update(input)

	// update valid instrument in database.
	if err = s.validInstrumentDataManager.UpdateValidInstrument(ctx, x); err != nil {
		logger.Error(err, "error encountered updating valid instrument")
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

// ArchiveHandler returns a handler that archives a valid instrument.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine valid instrument ID.
	validInstrumentID := s.validInstrumentIDFetcher(req)
	logger = logger.WithValue("valid_instrument_id", validInstrumentID)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)

	// archive the valid instrument in the database.
	err = s.validInstrumentDataManager.ArchiveValidInstrument(ctx, validInstrumentID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting valid instrument")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.validInstrumentCounter.Decrement(ctx)
	s.reporter.Report(newsman.Event{
		EventType: string(models.Archive),
		Data:      &models.ValidInstrument{ID: validInstrumentID},
		Topics:    []string{topicName},
	})

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
