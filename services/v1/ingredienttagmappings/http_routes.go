package ingredienttagmappings

import (
	"database/sql"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// URIParamKey is a standard string that we'll use to refer to ingredient tag mapping IDs with.
	URIParamKey = "ingredientTagMappingID"
)

// ListHandler is our list route.
func (s *Service) ListHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// ensure query filter.
		filter := models.ExtractQueryFilter(req)

		// determine valid ingredient ID.
		validIngredientID := s.validIngredientIDFetcher(req)
		tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
		logger = logger.WithValue("valid_ingredient_id", validIngredientID)

		// fetch ingredient tag mappings from database.
		ingredientTagMappings, err := s.ingredientTagMappingDataManager.GetIngredientTagMappings(ctx, validIngredientID, filter)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list.
			ingredientTagMappings = &models.IngredientTagMappingList{
				IngredientTagMappings: []models.IngredientTagMapping{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching ingredient tag mappings")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, ingredientTagMappings); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our ingredient tag mapping creation route.
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check request context for parsed input struct.
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.IngredientTagMappingCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine valid ingredient ID.
		validIngredientID := s.validIngredientIDFetcher(req)
		logger = logger.WithValue("valid_ingredient_id", validIngredientID)
		tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

		input.BelongsToValidIngredient = validIngredientID

		validIngredientExists, err := s.validIngredientDataManager.ValidIngredientExists(ctx, validIngredientID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking valid ingredient existence")
			res.WriteHeader(http.StatusInternalServerError)
			return
		} else if !validIngredientExists {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		// create ingredient tag mapping in database.
		x, err := s.ingredientTagMappingDataManager.CreateIngredientTagMapping(ctx, input)
		if err != nil {
			logger.Error(err, "error creating ingredient tag mapping")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tracing.AttachIngredientTagMappingIDToSpan(span, x.ID)
		logger = logger.WithValue("ingredient_tag_mapping_id", x.ID)

		// notify relevant parties.
		s.ingredientTagMappingCounter.Increment(ctx)
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

// ExistenceHandler returns a HEAD handler that returns 200 if an ingredient tag mapping exists, 404 otherwise.
func (s *Service) ExistenceHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// determine valid ingredient ID.
		validIngredientID := s.validIngredientIDFetcher(req)
		tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
		logger = logger.WithValue("valid_ingredient_id", validIngredientID)

		// determine ingredient tag mapping ID.
		ingredientTagMappingID := s.ingredientTagMappingIDFetcher(req)
		tracing.AttachIngredientTagMappingIDToSpan(span, ingredientTagMappingID)
		logger = logger.WithValue("ingredient_tag_mapping_id", ingredientTagMappingID)

		// fetch ingredient tag mapping from database.
		exists, err := s.ingredientTagMappingDataManager.IngredientTagMappingExists(ctx, validIngredientID, ingredientTagMappingID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking ingredient tag mapping existence in database")
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

// ReadHandler returns a GET handler that returns an ingredient tag mapping.
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// determine valid ingredient ID.
		validIngredientID := s.validIngredientIDFetcher(req)
		tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
		logger = logger.WithValue("valid_ingredient_id", validIngredientID)

		// determine ingredient tag mapping ID.
		ingredientTagMappingID := s.ingredientTagMappingIDFetcher(req)
		tracing.AttachIngredientTagMappingIDToSpan(span, ingredientTagMappingID)
		logger = logger.WithValue("ingredient_tag_mapping_id", ingredientTagMappingID)

		// fetch ingredient tag mapping from database.
		x, err := s.ingredientTagMappingDataManager.GetIngredientTagMapping(ctx, validIngredientID, ingredientTagMappingID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching ingredient tag mapping from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace.
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates an ingredient tag mapping.
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// check for parsed input attached to request context.
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.IngredientTagMappingUpdateInput)
		if !ok {
			logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine valid ingredient ID.
		validIngredientID := s.validIngredientIDFetcher(req)
		logger = logger.WithValue("valid_ingredient_id", validIngredientID)
		tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

		input.BelongsToValidIngredient = validIngredientID

		// determine ingredient tag mapping ID.
		ingredientTagMappingID := s.ingredientTagMappingIDFetcher(req)
		logger = logger.WithValue("ingredient_tag_mapping_id", ingredientTagMappingID)
		tracing.AttachIngredientTagMappingIDToSpan(span, ingredientTagMappingID)

		// fetch ingredient tag mapping from database.
		x, err := s.ingredientTagMappingDataManager.GetIngredientTagMapping(ctx, validIngredientID, ingredientTagMappingID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting ingredient tag mapping")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure.
		x.Update(input)

		// update ingredient tag mapping in database.
		if err = s.ingredientTagMappingDataManager.UpdateIngredientTagMapping(ctx, x); err != nil {
			logger.Error(err, "error encountered updating ingredient tag mapping")
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

// ArchiveHandler returns a handler that archives an ingredient tag mapping.
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var err error
		ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		logger := s.logger.WithRequest(req)

		// determine valid ingredient ID.
		validIngredientID := s.validIngredientIDFetcher(req)
		logger = logger.WithValue("valid_ingredient_id", validIngredientID)
		tracing.AttachValidIngredientIDToSpan(span, validIngredientID)

		validIngredientExists, err := s.validIngredientDataManager.ValidIngredientExists(ctx, validIngredientID)
		if err != nil && err != sql.ErrNoRows {
			logger.Error(err, "error checking valid ingredient existence")
			res.WriteHeader(http.StatusInternalServerError)
			return
		} else if !validIngredientExists {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		// determine ingredient tag mapping ID.
		ingredientTagMappingID := s.ingredientTagMappingIDFetcher(req)
		logger = logger.WithValue("ingredient_tag_mapping_id", ingredientTagMappingID)
		tracing.AttachIngredientTagMappingIDToSpan(span, ingredientTagMappingID)

		// archive the ingredient tag mapping in the database.
		err = s.ingredientTagMappingDataManager.ArchiveIngredientTagMapping(ctx, validIngredientID, ingredientTagMappingID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting ingredient tag mapping")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties.
		s.ingredientTagMappingCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.IngredientTagMapping{ID: ingredientTagMappingID},
			Topics:    []string{topicName},
		})

		// encode our response and peace.
		res.WriteHeader(http.StatusNoContent)
	}
}
