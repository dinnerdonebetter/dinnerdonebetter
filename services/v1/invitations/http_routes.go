package invitations

import (
	"database/sql"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

const (
	// URIParamKey is a standard string that we'll use to refer to invitation IDs with
	URIParamKey = "invitationID"
)

func attachInvitationIDToSpan(span *trace.Span, invitationID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("invitation_id", strconv.FormatUint(invitationID, 10)))
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

		// fetch invitations from database
		invitations, err := s.invitationDatabase.GetInvitations(ctx, qf, userID)
		if err == sql.ErrNoRows {
			// in the event no rows exist return an empty list
			invitations = &models.InvitationList{
				Invitations: []models.Invitation{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching invitations")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, invitations); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our invitation creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()

		// determine user ID
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)

		// check request context for parsed input struct
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.InvitationCreationInput)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		logger = logger.WithValue("input", input)
		input.BelongsTo = userID

		// create invitation in database
		x, err := s.invitationDatabase.CreateInvitation(ctx, input)
		if err != nil {
			logger.Error(err, "error creating invitation")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.invitationCounter.Increment(ctx)
		attachInvitationIDToSpan(span, x.ID)
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

// ReadHandler returns a GET handler that returns an invitation
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		invitationID := s.invitationIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":       userID,
			"invitation_id": invitationID,
		})
		attachInvitationIDToSpan(span, invitationID)
		attachUserIDToSpan(span, userID)

		// fetch invitation from database
		x, err := s.invitationDatabase.GetInvitation(ctx, invitationID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching invitation from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// encode our response and peace
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates an invitation
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()

		// check for parsed input attached to request context
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.InvitationUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// determine relevant information
		userID := s.userIDFetcher(req)
		invitationID := s.invitationIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id":       userID,
			"invitation_id": invitationID,
		})
		attachInvitationIDToSpan(span, invitationID)
		attachUserIDToSpan(span, userID)

		// fetch invitation from database
		x, err := s.invitationDatabase.GetInvitation(ctx, invitationID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting invitation")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update the data structure
		x.Update(input)

		// update invitation in database
		if err = s.invitationDatabase.UpdateInvitation(ctx, x); err != nil {
			logger.Error(err, "error encountered updating invitation")
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

// ArchiveHandler returns a handler that archives an invitation
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()

		// determine relevant information
		userID := s.userIDFetcher(req)
		invitationID := s.invitationIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"invitation_id": invitationID,
			"user_id":       userID,
		})
		attachInvitationIDToSpan(span, invitationID)
		attachUserIDToSpan(span, userID)

		// archive the invitation in the database
		err := s.invitationDatabase.ArchiveInvitation(ctx, invitationID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting invitation")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// notify relevant parties
		s.invitationCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data:      &models.Invitation{ID: invitationID},
			Topics:    []string{topicName},
		})

		// encode our response and peace
		res.WriteHeader(http.StatusNoContent)
	}
}
