package instruments

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to instrument IDs with
	URIParamKey = "instrumentID"
)

func attachInstrumentIDToSpan(span *trace.Span, instrumentID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("instrument_id", strconv.FormatUint(instrumentID, 10)))
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

		// fetch instruments from database
		instruments, err := s.instrumentDatabase.GetInstruments(ctx, qf, userID)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			instruments = &models.InstrumentList{
				Instruments: []models.Instrument{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching instruments")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, instruments); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our instrument creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.InstrumentCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)
		input.BelongsTo = userID

		// create instrument in database
		x, err := s.instrumentDatabase.CreateInstrument(ctx, input)
		if err != nil {
			logger.Error(err, "error creating instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.instrumentCounter.Increment(ctx)
		attachInstrumentIDToSpan(span, x.ID)
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

// ReadHandler returns a GET handler that returns an instrument
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		instrumentID := s.instrumentIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":       userID,
			"instrument_id": instrumentID,
		})
		attachInstrumentIDToSpan(span, instrumentID)
		attachUserIDToSpan(span, userID)

		// fetch instrument from database
		x, err := s.instrumentDatabase.GetInstrument(ctx, instrumentID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching instrument from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates an instrument
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.InstrumentUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		instrumentID := s.instrumentIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":       userID,
			"instrument_id": instrumentID,
		})
		attachInstrumentIDToSpan(span, instrumentID)
		attachUserIDToSpan(span, userID)

		// fetch instrument from database
		x, err := s.instrumentDatabase.GetInstrument(ctx, instrumentID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update instrument in database
		if err = s.instrumentDatabase.UpdateInstrument(ctx, x); err != nil {
			logger.Error(err, "error encountered updating instrument")
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

// ArchiveHandler returns a handler that archives an instrument
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		instrumentID := s.instrumentIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"instrument_id": instrumentID,
			"user_id":       userID,
		})
		attachInstrumentIDToSpan(span, instrumentID)
		attachUserIDToSpan(span, userID)

		// archive the instrument in the database
		err := s.instrumentDatabase.ArchiveInstrument(ctx, instrumentID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.instrumentCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.Instrument{ID: instrumentID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
