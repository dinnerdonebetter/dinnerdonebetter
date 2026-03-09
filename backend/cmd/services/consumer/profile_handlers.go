package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/components"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	webappauth "github.com/dinnerdonebetter/backend/internal/webapp/auth"

	"google.golang.org/protobuf/types/known/timestamppb"
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (s *ConsumerFrontendServer) ProfilePage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		return page("Profile",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(ghtml.Class("text-red-600"), g.Text("Unable to load. Please try again.")),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	selfRes, err := c.GetSelf(ctx, &authsvc.GetSelfRequest{})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting self")
		return page("Profile",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(ghtml.Class("text-red-600"), g.Text("Unable to load profile. Please try again.")),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	user := selfRes.GetResult()
	formErrors := (*components.ProfileFormErrors)(nil)
	if e := req.URL.Query().Get("error"); e != "" {
		formErrors = profileErrorToFormErrors(e)
	}

	flashMsg := ""
	if req.URL.Query().Get("updated") == "1" {
		flashMsg = "Profile updated."
	}

	return page("Profile",
		ghtml.Div(
			ghtml.Class("space-y-6"),
			ghtml.H2(ghtml.Class("text-xl font-semibold"), g.Text("Profile")),
			ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline text-sm"), g.Text("Back to Account Settings")),
			g.If(flashMsg != "", ghtml.Div(ghtml.Class("p-3 rounded-md text-sm bg-green-50 text-green-800"), g.Text(flashMsg))),
			s.componentRenderer.ProfileContent(user, formErrors),
		),
	), nil
}

func profileErrorToFormErrors(e string) *components.ProfileFormErrors {
	errs := &components.ProfileFormErrors{}
	switch e {
	case "invalid_username":
		errs.Username = "Username is required."
	case "invalid_first_name":
		errs.FirstName = "First name is required."
	case "invalid_password":
		errs.Password = "Password is required to update details."
	case errorParamInvalid, "invalid_input":
		errs.DetailsForm = errorMsgInvalidInput
	case errorParamUpdateFailed:
		errs.DetailsForm = "Failed to save. Please try again."
	case errorParamServer:
		errs.DetailsForm = errorMsgServer
	default:
		errs.DetailsForm = errorMsgSomethingWrong
	}
	return errs
}

func (s *ConsumerFrontendServer) UpdateProfileUsernameHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if req.Method != http.MethodPost {
		http.Redirect(res, req, "/account/profile", http.StatusFound)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Redirect(res, req, "/account/profile?error=invalid_input", http.StatusFound)
		return
	}

	username := strings.TrimSpace(req.FormValue("username"))
	if username == "" {
		http.Redirect(res, req, "/account/profile?error=invalid_username", http.StatusFound)
		return
	}

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		http.Redirect(res, req, "/account/profile?error=server", http.StatusFound)
		return
	}

	_, err = c.UpdateUserUsername(ctx, &identitysvc.UpdateUserUsernameRequest{
		NewUsername: username,
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "updating username")
		http.Redirect(res, req, "/account/profile?error=update_failed", http.StatusFound)
		return
	}

	http.Redirect(res, req, "/account/profile?updated=1", http.StatusFound)
}

func (s *ConsumerFrontendServer) UpdateProfileDetailsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if req.Method != http.MethodPost {
		http.Redirect(res, req, "/account/profile", http.StatusFound)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Redirect(res, req, "/account/profile?error=invalid_input", http.StatusFound)
		return
	}

	firstName := strings.TrimSpace(req.FormValue("first_name"))
	lastName := strings.TrimSpace(req.FormValue("last_name"))
	currentPassword := strings.TrimSpace(req.FormValue("current_password"))
	totpToken := strings.TrimSpace(req.FormValue("totp_token"))

	if firstName == "" {
		http.Redirect(res, req, "/account/profile?error=invalid_first_name", http.StatusFound)
		return
	}
	if currentPassword == "" {
		http.Redirect(res, req, "/account/profile?error=invalid_password", http.StatusFound)
		return
	}

	var birthday *timestamppb.Timestamp
	if b := strings.TrimSpace(req.FormValue("birthday")); b != "" {
		t, err := time.Parse("2006-01-02", b)
		if err == nil {
			birthday = timestamppb.New(t)
		}
	}

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		http.Redirect(res, req, "/account/profile?error=server", http.StatusFound)
		return
	}

	_, err = c.UpdateUserDetails(ctx, &identitysvc.UpdateUserDetailsRequest{
		Input: &identitysvc.UserDetailsUpdateRequestInput{
			FirstName:       firstName,
			LastName:        lastName,
			Birthday:        birthday,
			CurrentPassword: currentPassword,
			TotpToken:       totpToken,
		},
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user details")
		http.Redirect(res, req, "/account/profile?error=update_failed", http.StatusFound)
		return
	}

	http.Redirect(res, req, "/account/profile?updated=1", http.StatusFound)
}
