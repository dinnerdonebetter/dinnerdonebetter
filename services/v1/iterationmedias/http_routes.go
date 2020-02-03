package iterationmedias

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to iteration media IDs with
	URIParamKey = "iterationMediaID"
)

func attachIterationMediaIDToSpan(span *trace.Span, iterationMediaID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("iteration_media_id", strconv.FormatUint(iterationMediaID, 10)))
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

		// fetch iteration medias from database
		iterationMedias, err := s.iterationMediaDatabase.GetIterationMedias(ctx, qf, userID)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			iterationMedias = &models.IterationMediaList{
				IterationMedias: []models.IterationMedia{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching iteration medias")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, iterationMedias); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our iteration media creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.IterationMediaCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)
		input.BelongsTo = userID

		// create iteration media in database
		x, err := s.iterationMediaDatabase.CreateIterationMedia(ctx, input)
		if err != nil {
			logger.Error(err, "error creating iteration media")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.iterationMediaCounter.Increment(ctx)
		attachIterationMediaIDToSpan(span, x.ID)
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

// ReadHandler returns a GET handler that returns an iteration media
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		iterationMediaID := s.iterationMediaIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":            userID,
			"iteration_media_id": iterationMediaID,
		})
		attachIterationMediaIDToSpan(span, iterationMediaID)
		attachUserIDToSpan(span, userID)

		// fetch iteration media from database
		x, err := s.iterationMediaDatabase.GetIterationMedia(ctx, iterationMediaID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching iteration media from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates an iteration media
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.IterationMediaUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		iterationMediaID := s.iterationMediaIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":            userID,
			"iteration_media_id": iterationMediaID,
		})
		attachIterationMediaIDToSpan(span, iterationMediaID)
		attachUserIDToSpan(span, userID)

		// fetch iteration media from database
		x, err := s.iterationMediaDatabase.GetIterationMedia(ctx, iterationMediaID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting iteration media")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update iteration media in database
		if err = s.iterationMediaDatabase.UpdateIterationMedia(ctx, x); err != nil {
			logger.Error(err, "error encountered updating iteration media")
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

// ArchiveHandler returns a handler that archives an iteration media
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		iterationMediaID := s.iterationMediaIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"iteration_media_id": iterationMediaID,
			"user_id":            userID,
		})
		attachIterationMediaIDToSpan(span, iterationMediaID)
		attachUserIDToSpan(span, userID)

		// archive the iteration media in the database
		err := s.iterationMediaDatabase.ArchiveIterationMedia(ctx, iterationMediaID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting iteration media")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.iterationMediaCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.IterationMedia{ID: iterationMediaID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
