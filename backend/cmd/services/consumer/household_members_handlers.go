package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/components"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	webappauth "github.com/dinnerdonebetter/backend/internal/webapp/auth"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (s *ConsumerFrontendServer) HouseholdMembersPage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		return page("Household Members",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(ghtml.Class("text-red-600"), g.Text("Unable to load. Please try again.")),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	activeRes, err := c.GetActiveAccount(ctx, &authsvc.GetActiveAccountRequest{})
	if err != nil || activeRes == nil || activeRes.Result == nil {
		observability.AcknowledgeError(err, logger, span, "getting active account")
		return page("Household Members",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.H2(ghtml.Class("text-xl font-semibold"), g.Text("Household Members")),
				ghtml.P(
					ghtml.Class("text-gray-600"),
					g.Text("No household found. Create an account to get started."),
				),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	account := activeRes.Result
	selfRes, err := c.GetSelf(ctx, &authsvc.GetSelfRequest{})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting self")
		return page("Household Members",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(ghtml.Class("text-red-600"), g.Text("Unable to load. Please try again.")),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	currentUserID := ""
	if selfRes != nil && selfRes.Result != nil {
		currentUserID = selfRes.Result.Id
	}

	isAdmin := false
	for _, m := range account.Members {
		if m.BelongsToUser != nil && m.BelongsToUser.Id == currentUserID {
			if m.AccountRole == authorization.AccountAdminRoleName {
				isAdmin = true
			}
			break
		}
	}

	var invitations []*identitysvc.AccountInvitation
	invRes, err := c.GetSentAccountInvitations(ctx, &identitysvc.GetSentAccountInvitationsRequest{
		Filter: &filtering.QueryFilter{MaxResponseSize: new(uint32(50))},
	})
	if err == nil && invRes != nil {
		for _, inv := range invRes.Results {
			if inv.DestinationAccount != nil && inv.DestinationAccount.Id == account.Id {
				invitations = append(invitations, inv)
			}
		}
	}

	baseURL := buildBaseURL(req)

	queryParams := req.URL.Query()
	formErrors := (*components.HouseholdMembersFormErrors)(nil)
	if queryParams.Get("error") == "invalid_email" {
		formErrors = &components.HouseholdMembersFormErrors{Email: "Please enter a valid email address."}
	}

	flashMsg := ""
	if queryParams.Get("invited") == "1" {
		flashMsg = "Invitation sent successfully."
	}
	if e := queryParams.Get("error"); e != "" && flashMsg == "" {
		flashMsg = errorMessageForParam(e)
	}

	return page("Household Members",
		ghtml.Div(
			ghtml.Class("space-y-6"),
			ghtml.H2(ghtml.Class("text-xl font-semibold"), g.Text("Household Members")),
			ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline text-sm"), g.Text("Back to Account Settings")),
			g.If(flashMsg != "", buildFlashMessage(flashMsg)),
			s.componentRenderer.HouseholdMembersContent(account, invitations, currentUserID, isAdmin, baseURL, formErrors),
			ghtml.Script(ghtml.Type("text/javascript"), g.Raw(copyInviteLinkScript)),
		),
	), nil
}

func buildFlashMessage(flashMsg string) g.Node {
	cls := "p-3 rounded-md text-sm bg-red-50 text-red-800"
	if strings.HasPrefix(flashMsg, "Invitation") {
		cls = "p-3 rounded-md text-sm bg-green-50 text-green-800"
	}
	return ghtml.Div(ghtml.Class(cls), g.Text(flashMsg))
}

func errorMessageForParam(e string) string {
	switch e {
	case errorParamInvalid, "invalid_email":
		return errorMsgInvalidInput
	case "invalid_role":
		return "Invalid role selected."
	case "invitation_failed":
		return "Failed to send invitation. Please try again."
	case "cancel_failed":
		return "Failed to cancel invitation."
	case "role_update_failed":
		return "Failed to update member role."
	case errorParamServer:
		return errorMsgServer
	default:
		return errorMsgSomethingWrong
	}
}

func (s *ConsumerFrontendServer) SendInvitationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if req.Method != http.MethodPost {
		http.Redirect(res, req, "/account/household-members", http.StatusFound)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Redirect(res, req, "/account/household-members?error=invalid", http.StatusFound)
		return
	}

	email := strings.TrimSpace(req.FormValue("email"))
	name := strings.TrimSpace(req.FormValue("name"))
	note := strings.TrimSpace(req.FormValue("note"))

	if email == "" || !emailRegex.MatchString(email) {
		http.Redirect(res, req, "/account/household-members?error=invalid_email", http.StatusFound)
		return
	}

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		http.Redirect(res, req, "/account/household-members?error=server", http.StatusFound)
		return
	}

	_, err = c.CreateAccountInvitation(ctx, &identitysvc.CreateAccountInvitationRequest{
		Input: &identitysvc.AccountInvitationCreationRequestInput{
			ToEmail: email,
			ToName:  name,
			Note:    note,
		},
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating invitation")
		http.Redirect(res, req, "/account/household-members?error=invitation_failed", http.StatusFound)
		return
	}

	http.Redirect(res, req, "/account/household-members?invited=1", http.StatusFound)
}

func (s *ConsumerFrontendServer) CancelInvitationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if req.Method != http.MethodPost {
		http.Redirect(res, req, "/account/household-members", http.StatusFound)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Redirect(res, req, "/account/household-members?error=invalid", http.StatusFound)
		return
	}

	invitationID := strings.TrimSpace(req.FormValue("invitation_id"))
	if invitationID == "" {
		http.Redirect(res, req, "/account/household-members?error=invalid", http.StatusFound)
		return
	}

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		http.Redirect(res, req, "/account/household-members?error=server", http.StatusFound)
		return
	}

	_, err = c.CancelAccountInvitation(ctx, &identitysvc.CancelAccountInvitationRequest{
		AccountInvitationId: invitationID,
		Input:               &identitysvc.AccountInvitationUpdateRequestInput{},
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "cancelling invitation")
		http.Redirect(res, req, "/account/household-members?error=cancel_failed", http.StatusFound)
		return
	}

	http.Redirect(res, req, "/account/household-members", http.StatusFound)
}

func (s *ConsumerFrontendServer) UpdateMemberRoleHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if req.Method != http.MethodPost {
		http.Redirect(res, req, "/account/household-members", http.StatusFound)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Redirect(res, req, "/account/household-members?error=invalid", http.StatusFound)
		return
	}

	userID := strings.TrimSpace(req.FormValue("user_id"))
	newRole := strings.TrimSpace(req.FormValue("new_role"))
	reason := strings.TrimSpace(req.FormValue("reason"))

	if userID == "" || newRole == "" || reason == "" {
		http.Redirect(res, req, "/account/household-members?error=invalid", http.StatusFound)
		return
	}

	if newRole != authorization.AccountAdminRoleName && newRole != authorization.AccountMemberRoleName {
		http.Redirect(res, req, "/account/household-members?error=invalid_role", http.StatusFound)
		return
	}

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		http.Redirect(res, req, "/account/household-members?error=server", http.StatusFound)
		return
	}

	_, err = c.UpdateAccountMemberPermissions(ctx, &identitysvc.UpdateAccountMemberPermissionsRequest{
		UserId: userID,
		Input: &identitysvc.ModifyUserPermissionsInput{
			NewRole: newRole,
			Reason:  reason,
		},
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "updating member role")
		http.Redirect(res, req, "/account/household-members?error=role_update_failed", http.StatusFound)
		return
	}

	http.Redirect(res, req, "/account/household-members", http.StatusFound)
}

func buildBaseURL(req *http.Request) string {
	scheme := "https"
	if req.TLS == nil && (req.Header.Get("X-Forwarded-Proto") == "http" || req.Host == "localhost" || strings.HasPrefix(req.Host, "localhost:")) {
		scheme = "http"
	}
	if proto := req.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	return fmt.Sprintf("%s://%s", scheme, req.Host)
}

const copyInviteLinkScript = `
document.querySelectorAll('.copy-invite-link').forEach(function(btn) {
	btn.addEventListener('click', function() {
		var url = this.getAttribute('data-url');
		if (url && navigator.clipboard && navigator.clipboard.writeText) {
			navigator.clipboard.writeText(url).then(function() { alert('Link copied to clipboard'); });
		}
	});
});
`
