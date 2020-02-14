package preparations

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to preparation IDs with
	URIParamKey = "preparationID"
)

func attachPreparationIDToSpan(span *trace.Span, preparationID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("preparation_id", strconv.FormatUint(preparationID, 10)))
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

		// fetch preparations from database
		preparations, err := s.preparationDatabase.GetPreparations(ctx, qf)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			preparations = &models.PreparationList{
				Preparations: []models.Preparation{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching preparations")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, preparations); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our preparation creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.PreparationCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)

		// create preparation in database
		x, err := s.preparationDatabase.CreatePreparation(ctx, input)
		if err != nil {
			logger.Error(err, "error creating preparation")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.preparationCounter.Increment(ctx)
		attachPreparationIDToSpan(span, x.ID)
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

// ReadHandler returns a GET handler that returns a preparation
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		preparationID := s.preparationIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":        userID,
			"preparation_id": preparationID,
		})
		attachPreparationIDToSpan(span, preparationID)
		attachUserIDToSpan(span, userID)

		// fetch preparation from database
		x, err := s.preparationDatabase.GetPreparation(ctx, preparationID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching preparation from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a preparation
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.PreparationUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		preparationID := s.preparationIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":        userID,
			"preparation_id": preparationID,
		})
		attachPreparationIDToSpan(span, preparationID)
		attachUserIDToSpan(span, userID)

		// fetch preparation from database
		x, err := s.preparationDatabase.GetPreparation(ctx, preparationID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting preparation")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update preparation in database
		if err = s.preparationDatabase.UpdatePreparation(ctx, x); err != nil {
			logger.Error(err, "error encountered updating preparation")
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

// ArchiveHandler returns a handler that archives a preparation
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		preparationID := s.preparationIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"preparation_id": preparationID,
			"user_id":        userID,
		})
		attachPreparationIDToSpan(span, preparationID)
		attachUserIDToSpan(span, userID)

		// archive the preparation in the database
		err := s.preparationDatabase.ArchivePreparation(ctx, preparationID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting preparation")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.preparationCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.Preparation{ID: preparationID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
