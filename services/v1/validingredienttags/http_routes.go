package validingredienttags

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to valid ingredient tag IDs with.
	URIParamKey = "validIngredientTagID"
)

// ListHandler is our list route.
func (s *Service) ListHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// ensure query filter.
		filter := models.ExtractQueryFilter(req)

		// fetch valid ingredient tags from database.
		validIngredientTags, err := s.validIngredientTagDataManager.GetValidIngredientTags(ctx, filter)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list.
			validIngredientTags = &models.ValidIngredientTagList{
				ValidIngredientTags: []models.ValidIngredientTag{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching valid ingredient tags")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, validIngredientTags); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our valid ingredient tag creation route.
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check request context for parsed input struct.
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.ValidIngredientTagCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// create valid ingredient tag in database.
		x, err := s.validIngredientTagDataManager.CreateValidIngredientTag(ctx, input)
		if err != nil {
			logger.Error(err, "error creating valid ingredient tag")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tracing.AttachValidIngredientTagIDToSpan(span, x.ID)
		logger = logger.WithValue("valid_ingredient_tag_id", x.ID)

		// notify relevant parties.
		s.validIngredientTagCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if a valid ingredient tag exists, 404 otherwise.
func (s *Service) ExistenceHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// determine valid ingredient tag ID.
		validIngredientTagID := s.validIngredientTagIDFetcher(req)
		tracing.AttachValidIngredientTagIDToSpan(span, validIngredientTagID)
		logger = logger.WithValue("valid_ingredient_tag_id", validIngredientTagID)

		// fetch valid ingredient tag from database.
		exists, err := s.validIngredientTagDataManager.ValidIngredientTagExists(ctx, validIngredientTagID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking valid ingredient tag existence in database")
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

// ReadHandler returns a GET handler that returns a valid ingredient tag.
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// determine valid ingredient tag ID.
		validIngredientTagID := s.validIngredientTagIDFetcher(req)
		tracing.AttachValidIngredientTagIDToSpan(span, validIngredientTagID)
		logger = logger.WithValue("valid_ingredient_tag_id", validIngredientTagID)

		// fetch valid ingredient tag from database.
		x, err := s.validIngredientTagDataManager.GetValidIngredientTag(ctx, validIngredientTagID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching valid ingredient tag from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates a valid ingredient tag.
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check for parsed input attached to request context.
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.ValidIngredientTagUpdateInput)
		if !ok {
			logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine valid ingredient tag ID.
		validIngredientTagID := s.validIngredientTagIDFetcher(req)
		logger = logger.WithValue("valid_ingredient_tag_id", validIngredientTagID)
		tracing.AttachValidIngredientTagIDToSpan(span, validIngredientTagID)

		// fetch valid ingredient tag from database.
		x, err := s.validIngredientTagDataManager.GetValidIngredientTag(ctx, validIngredientTagID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting valid ingredient tag")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure.
		x.Update(input)

		// update valid ingredient tag in database.
		if err = s.validIngredientTagDataManager.UpdateValidIngredientTag(ctx, x); err != nil {
			logger.Error(err, "error encountered updating valid ingredient tag")
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

// ArchiveHandler returns a handler that archives a valid ingredient tag.
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var err error
		ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// determine valid ingredient tag ID.
		validIngredientTagID := s.validIngredientTagIDFetcher(req)
		logger = logger.WithValue("valid_ingredient_tag_id", validIngredientTagID)
		tracing.AttachValidIngredientTagIDToSpan(span, validIngredientTagID)

		// archive the valid ingredient tag in the database.
		err = s.validIngredientTagDataManager.ArchiveValidIngredientTag(ctx, validIngredientTagID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting valid ingredient tag")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties.
		s.validIngredientTagCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.ValidIngredientTag{ID: validIngredientTagID},
			Topics:    []string{topicName},
		})

		// encode our response and peace.
		res.WriteHeader(http.StatusNoContent)
	}
}
