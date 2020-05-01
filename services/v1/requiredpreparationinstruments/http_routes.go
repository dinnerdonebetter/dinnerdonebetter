package requiredpreparationinstruments

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to required preparation instrument IDs with.
	URIParamKey = "requiredPreparationInstrumentID"
)

// ListHandler is our list route.
func (s *Service) ListHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// ensure query filter.
		filter := models.ExtractQueryFilter(req)

		// determine valid preparation ID.
		validPreparationID := s.validPreparationIDFetcher(req)
		tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
		logger = logger.WithValue("valid_preparation_id", validPreparationID)

		// fetch required preparation instruments from database.
		requiredPreparationInstruments, err := s.requiredPreparationInstrumentDataManager.GetRequiredPreparationInstruments(ctx, validPreparationID, filter)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list.
			requiredPreparationInstruments = &models.RequiredPreparationInstrumentList{
				RequiredPreparationInstruments: []models.RequiredPreparationInstrument{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching required preparation instruments")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, requiredPreparationInstruments); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our required preparation instrument creation route.
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check request context for parsed input struct.
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.RequiredPreparationInstrumentCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine valid preparation ID.
		validPreparationID := s.validPreparationIDFetcher(req)
		logger = logger.WithValue("valid_preparation_id", validPreparationID)
		tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

		input.BelongsToValidPreparation = validPreparationID

		validPreparationExists, err := s.validPreparationDataManager.ValidPreparationExists(ctx, validPreparationID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking valid preparation existence")
			res.WriteHeader(http.StatusInternalServerError)
			return
		} else if !validPreparationExists {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		// create required preparation instrument in database.
		x, err := s.requiredPreparationInstrumentDataManager.CreateRequiredPreparationInstrument(ctx, input)
		if err != nil {
			logger.Error(err, "error creating required preparation instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tracing.AttachRequiredPreparationInstrumentIDToSpan(span, x.ID)
		logger = logger.WithValue("required_preparation_instrument_id", x.ID)

		// notify relevant parties.
		s.requiredPreparationInstrumentCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if a required preparation instrument exists, 404 otherwise.
func (s *Service) ExistenceHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// determine valid preparation ID.
		validPreparationID := s.validPreparationIDFetcher(req)
		tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
		logger = logger.WithValue("valid_preparation_id", validPreparationID)

		// determine required preparation instrument ID.
		requiredPreparationInstrumentID := s.requiredPreparationInstrumentIDFetcher(req)
		tracing.AttachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)
		logger = logger.WithValue("required_preparation_instrument_id", requiredPreparationInstrumentID)

		// fetch required preparation instrument from database.
		exists, err := s.requiredPreparationInstrumentDataManager.RequiredPreparationInstrumentExists(ctx, validPreparationID, requiredPreparationInstrumentID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking required preparation instrument existence in database")
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

// ReadHandler returns a GET handler that returns a required preparation instrument.
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// determine valid preparation ID.
		validPreparationID := s.validPreparationIDFetcher(req)
		tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
		logger = logger.WithValue("valid_preparation_id", validPreparationID)

		// determine required preparation instrument ID.
		requiredPreparationInstrumentID := s.requiredPreparationInstrumentIDFetcher(req)
		tracing.AttachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)
		logger = logger.WithValue("required_preparation_instrument_id", requiredPreparationInstrumentID)

		// fetch required preparation instrument from database.
		x, err := s.requiredPreparationInstrumentDataManager.GetRequiredPreparationInstrument(ctx, validPreparationID, requiredPreparationInstrumentID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching required preparation instrument from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a required preparation instrument.
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check for parsed input attached to request context.
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.RequiredPreparationInstrumentUpdateInput)
		if !ok {
			logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine valid preparation ID.
		validPreparationID := s.validPreparationIDFetcher(req)
		logger = logger.WithValue("valid_preparation_id", validPreparationID)
		tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

		input.BelongsToValidPreparation = validPreparationID

		// determine required preparation instrument ID.
		requiredPreparationInstrumentID := s.requiredPreparationInstrumentIDFetcher(req)
		logger = logger.WithValue("required_preparation_instrument_id", requiredPreparationInstrumentID)
		tracing.AttachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)

		// fetch required preparation instrument from database.
		x, err := s.requiredPreparationInstrumentDataManager.GetRequiredPreparationInstrument(ctx, validPreparationID, requiredPreparationInstrumentID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting required preparation instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure.
		x.Update(input)

		// update required preparation instrument in database.
		if err = s.requiredPreparationInstrumentDataManager.UpdateRequiredPreparationInstrument(ctx, x); err != nil {
			logger.Error(err, "error encountered updating required preparation instrument")
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

// ArchiveHandler returns a handler that archives a required preparation instrument.
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var err error
		ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// determine valid preparation ID.
		validPreparationID := s.validPreparationIDFetcher(req)
		logger = logger.WithValue("valid_preparation_id", validPreparationID)
		tracing.AttachValidPreparationIDToSpan(span, validPreparationID)

		validPreparationExists, err := s.validPreparationDataManager.ValidPreparationExists(ctx, validPreparationID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking valid preparation existence")
			res.WriteHeader(http.StatusInternalServerError)
			return
		} else if !validPreparationExists {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		// determine required preparation instrument ID.
		requiredPreparationInstrumentID := s.requiredPreparationInstrumentIDFetcher(req)
		logger = logger.WithValue("required_preparation_instrument_id", requiredPreparationInstrumentID)
		tracing.AttachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)

		// archive the required preparation instrument in the database.
		err = s.requiredPreparationInstrumentDataManager.ArchiveRequiredPreparationInstrument(ctx, validPreparationID, requiredPreparationInstrumentID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting required preparation instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties.
		s.requiredPreparationInstrumentCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.RequiredPreparationInstrument{ID: requiredPreparationInstrumentID},
			Topics:    []string{topicName},
		})

		// encode our response and peace.
		res.WriteHeader(http.StatusNoContent)
	}
}
