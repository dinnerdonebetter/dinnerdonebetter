package requiredpreparationinstruments

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to required preparation instrument IDs with
	URIParamKey = "requiredPreparationInstrumentID"
)

func attachRequiredPreparationInstrumentIDToSpan(span *trace.Span, requiredPreparationInstrumentID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("required_preparation_instrument_id", strconv.FormatUint(requiredPreparationInstrumentID, 10)))
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

		// fetch required preparation instruments from database
		requiredPreparationInstruments, err := s.requiredPreparationInstrumentDatabase.GetRequiredPreparationInstruments(ctx, qf, userID)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			requiredPreparationInstruments = &models.RequiredPreparationInstrumentList{
				RequiredPreparationInstruments: []models.RequiredPreparationInstrument{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching required preparation instruments")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, requiredPreparationInstruments); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our required preparation instrument creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.RequiredPreparationInstrumentCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)
		input.BelongsTo = userID

		// create required preparation instrument in database
		x, err := s.requiredPreparationInstrumentDatabase.CreateRequiredPreparationInstrument(ctx, input)
		if err != nil {
			logger.Error(err, "error creating required preparation instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.requiredPreparationInstrumentCounter.Increment(ctx)
		attachRequiredPreparationInstrumentIDToSpan(span, x.ID)
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

// ReadHandler returns a GET handler that returns a required preparation instrument
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		requiredPreparationInstrumentID := s.requiredPreparationInstrumentIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":                            userID,
			"required_preparation_instrument_id": requiredPreparationInstrumentID,
		})
		attachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)
		attachUserIDToSpan(span, userID)

		// fetch required preparation instrument from database
		x, err := s.requiredPreparationInstrumentDatabase.GetRequiredPreparationInstrument(ctx, requiredPreparationInstrumentID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching required preparation instrument from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a required preparation instrument
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.RequiredPreparationInstrumentUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		requiredPreparationInstrumentID := s.requiredPreparationInstrumentIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":                            userID,
			"required_preparation_instrument_id": requiredPreparationInstrumentID,
		})
		attachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)
		attachUserIDToSpan(span, userID)

		// fetch required preparation instrument from database
		x, err := s.requiredPreparationInstrumentDatabase.GetRequiredPreparationInstrument(ctx, requiredPreparationInstrumentID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting required preparation instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update required preparation instrument in database
		if err = s.requiredPreparationInstrumentDatabase.UpdateRequiredPreparationInstrument(ctx, x); err != nil {
			logger.Error(err, "error encountered updating required preparation instrument")
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

// ArchiveHandler returns a handler that archives a required preparation instrument
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		requiredPreparationInstrumentID := s.requiredPreparationInstrumentIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"required_preparation_instrument_id": requiredPreparationInstrumentID,
			"user_id":                            userID,
		})
		attachRequiredPreparationInstrumentIDToSpan(span, requiredPreparationInstrumentID)
		attachUserIDToSpan(span, userID)

		// archive the required preparation instrument in the database
		err := s.requiredPreparationInstrumentDatabase.ArchiveRequiredPreparationInstrument(ctx, requiredPreparationInstrumentID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting required preparation instrument")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.requiredPreparationInstrumentCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.RequiredPreparationInstrument{ID: requiredPreparationInstrumentID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
