package frontend

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	invitationIDURLParamKey = "invitation"
)

func (s *service) fetchInvitation(ctx context.Context, req *http.Request) (invitation *types.Invitation, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	// determine invitation ID.
	invitationID := s.invitationIDFetcher(req)
	tracing.AttachInvitationIDToSpan(span, invitationID)
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)

	invitation, err = s.dataStore.GetInvitation(ctx, invitationID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching invitation data")
	}

	return invitation, nil
}

//go:embed templates/partials/generated/creators/invitation_creator.gotpl
var invitationCreatorTemplate string

func (s *service) buildInvitationCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)
		tracing.AttachRequestToSpan(span, req)

		sessionCtxData, err := s.sessionContextDataFetcher(req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
			http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
			return
		}

		invitation := &types.Invitation{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(invitationCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New Invitation",
				ContentData: invitation,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", invitationCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, invitation, res)
		}
	}
}

const (
	invitationCodeFormKey     = "code"
	invitationConsumedFormKey = "consumed"

	invitationCreationInputCodeFormKey     = invitationCodeFormKey
	invitationCreationInputConsumedFormKey = invitationConsumedFormKey

	invitationUpdateInputCodeFormKey     = invitationCodeFormKey
	invitationUpdateInputConsumedFormKey = invitationConsumedFormKey
)

// parseFormEncodedInvitationCreationInput checks a request for an InvitationCreationInput.
func (s *service) parseFormEncodedInvitationCreationInput(ctx context.Context, req *http.Request, sessionCtxData *types.SessionContextData) (creationInput *types.InvitationCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing invitation creation input")
		return nil
	}

	creationInput = &types.InvitationCreationInput{
		Code:             form.Get(invitationCreationInputCodeFormKey),
		Consumed:         s.stringToBool(form, invitationCreationInputConsumedFormKey),
		BelongsToAccount: sessionCtxData.ActiveAccountID,
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid invitation creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleInvitationCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("invitation creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for invitation creation route")

	creationInput := s.parseFormEncodedInvitationCreationInput(ctx, req, sessionCtxData)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing invitation creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debug("invitation creation input parsed successfully")

	if _, err = s.dataStore.CreateInvitation(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing invitation to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("invitation created")

	htmxRedirectTo(res, "/invitations")
	res.WriteHeader(http.StatusCreated)
}

//go:embed templates/partials/generated/editors/invitation_editor.gotpl
var invitationEditorTemplate string

func (s *service) buildInvitationEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)
		tracing.AttachRequestToSpan(span, req)

		sessionCtxData, err := s.sessionContextDataFetcher(req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
			http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
			return
		}

		invitation, err := s.fetchInvitation(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching invitation from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.Invitation) string {
				return fmt.Sprintf("Invitation #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(invitationEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("Invitation #%d", invitation.ID),
				ContentData: invitation,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", invitationEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, invitation, res)
		}
	}
}

func (s *service) fetchInvitations(ctx context.Context, req *http.Request) (invitations *types.InvitationList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilter(req)
	tracing.AttachQueryFilterToSpan(span, filter)

	invitations, err = s.dataStore.GetInvitations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching invitation data")
	}

	return invitations, nil
}

//go:embed templates/partials/generated/tables/invitations_table.gotpl
var invitationsTableTemplate string

func (s *service) buildInvitationsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)
		tracing.AttachRequestToSpan(span, req)

		sessionCtxData, err := s.sessionContextDataFetcher(req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
			http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
			return
		}

		invitations, err := s.fetchInvitations(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching invitations from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.Invitation) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/invitations/%d", x.ID))
			},
			"pushURL": func(x *types.Invitation) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/invitations/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(invitationsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "Invitations",
				ContentData: invitations,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", invitationsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, invitations, res)
		}
	}
}

// parseFormEncodedInvitationUpdateInput checks a request for an InvitationUpdateInput.
func (s *service) parseFormEncodedInvitationUpdateInput(ctx context.Context, req *http.Request, sessionCtxData *types.SessionContextData) (updateInput *types.InvitationUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing invitation creation input")
		return nil
	}

	updateInput = &types.InvitationUpdateInput{
		Code:             form.Get(invitationUpdateInputCodeFormKey),
		Consumed:         s.stringToBool(form, invitationUpdateInputConsumedFormKey),
		BelongsToAccount: sessionCtxData.ActiveAccountID,
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid invitation creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleInvitationUpdateRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	updateInput := s.parseFormEncodedInvitationUpdateInput(ctx, req, sessionCtxData)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	invitation, err := s.fetchInvitation(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching invitation from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := invitation.Update(updateInput)

	if err = s.dataStore.UpdateInvitation(ctx, invitation, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching invitation from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.Invitation) string {
			return fmt.Sprintf("Invitation #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", invitationEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, invitation, res)
}

func (s *service) handleInvitationArchiveRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	invitationID := s.invitationIDFetcher(req)
	tracing.AttachInvitationIDToSpan(span, invitationID)
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)

	if err = s.dataStore.ArchiveInvitation(ctx, invitationID, sessionCtxData.ActiveAccountID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving invitations in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	invitations, err := s.fetchInvitations(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching invitations from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.Invitation) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/invitations/%d", x.ID))
		},
		"pushURL": func(x *types.Invitation) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/invitations/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", invitationsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, invitations, res)
}
