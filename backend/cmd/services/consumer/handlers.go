package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/components"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	webappauth "github.com/dinnerdonebetter/backend/internal/webapp/auth"
	"github.com/dinnerdonebetter/backend/pkg/client"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	appStoreURL     = "https://apps.apple.com/app/dinner-done-better/id0000000000" // TODO: replace with actual App Store URL
	acceptInviteMsg = "To accept this invitation, open the Dinner Done Better app. If you don't have the app, download it from the App Store."
)

func (s *ConsumerFrontendServer) HomePage(_ http.ResponseWriter, _ *http.Request) (g.Node, error) {
	return page("Home",
		ghtml.Div(
			ghtml.Class("text-center space-y-4"),
			ghtml.H2(
				ghtml.Class("text-2xl font-bold"),
				g.Text("Welcome to Dinner Done Better"),
			),
			ghtml.P(
				ghtml.Class("text-gray-600"),
				g.Text("Manage your meal plans and grocery lists."),
			),
		),
	), nil
}

func (s *ConsumerFrontendServer) AccountSettingsPage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		return page("Account Settings",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(
					ghtml.Class("text-red-600"),
					g.Text("Unable to load account settings. Please try again."),
				),
			),
		), nil
	}

	selfRes, err := c.GetSelf(ctx, &authsvc.GetSelfRequest{})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting self")
		return page("Account Settings",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(
					ghtml.Class("text-red-600"),
					g.Text("Unable to load account settings. Please try again."),
				),
			),
		), nil
	}

	user := selfRes.GetResult()
	userID := ""
	if user != nil {
		userID = user.Id
	}

	hasAccount := false
	if userID != "" {
		var pageSize uint32 = 1
		accountsRes, getErr := c.GetAccountsForUser(ctx, &identitysvc.GetAccountsForUserRequest{
			UserId: userID,
			Filter: &filtering.QueryFilter{
				MaxResponseSize: &pageSize,
			},
		})
		if getErr == nil && len(accountsRes.GetResults()) > 0 {
			hasAccount = true
		}
	}

	return page("Account Settings",
		ghtml.Div(
			ghtml.Class("space-y-6"),
			ghtml.H2(
				ghtml.Class("text-xl font-semibold"),
				g.Text("Account Settings"),
			),
			ghtml.P(
				ghtml.Class("text-gray-600"),
				g.Text("Manage your account preferences here."),
			),
			s.componentRenderer.AccountLinks(&components.AccountLinksProps{HasAccount: hasAccount}),
			s.componentRenderer.PasskeySection(),
		),
	), nil
}

// AcceptInvitationPage shows Option A: "Open in app" or redirect to App Store.
// For web users who land on /accept_invitation (e.g. from a shared link when not on iOS).
func (s *ConsumerFrontendServer) AcceptInvitationPage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	// Option A: Show page explaining to open in app, or redirect to App Store if on mobile web
	userAgent := req.Header.Get("User-Agent")
	if isIOSWebBrowser(userAgent) {
		// On iOS Safari/web - suggest opening in app
		return page("Accept Invitation",
			ghtml.Div(
				ghtml.Class("space-y-6 text-center"),
				ghtml.H2(
					ghtml.Class("text-xl font-semibold"),
					g.Text("Accept Invitation"),
				),
				ghtml.P(
					ghtml.Class("text-gray-600"),
					g.Text(acceptInviteMsg),
				),
				ghtml.A(
					ghtml.Href(appStoreURL),
					ghtml.Class("inline-block px-6 py-3 rounded-md bg-blue-600 text-white font-medium hover:bg-blue-700"),
					g.Text("Open in App Store"),
				),
			),
		), nil
	}

	// On desktop or non-iOS - show generic message
	return page("Accept Invitation",
		ghtml.Div(
			ghtml.Class("space-y-6 text-center"),
			ghtml.H2(
				ghtml.Class("text-xl font-semibold"),
				g.Text("Accept Invitation"),
			),
			ghtml.P(
				ghtml.Class("text-gray-600"),
				g.Text("This invitation link is for the Dinner Done Better mobile app. Open this link on your iPhone or iPad to accept the invitation."),
			),
			ghtml.A(
				ghtml.Href(appStoreURL),
				ghtml.Class("inline-block px-6 py-3 rounded-md bg-blue-600 text-white font-medium hover:bg-blue-700"),
				g.Text("Get the App"),
			),
		),
	), nil
}

func isIOSWebBrowser(userAgent string) bool {
	return strings.Contains(strings.ToLower(userAgent), "iphone") ||
		strings.Contains(strings.ToLower(userAgent), "ipad")
}

// VerifyEmailAddressPage handles email verification links (e.g., from signup emails).
// Token is passed via query param ?t=TOKEN. No authentication required.
func (s *ConsumerFrontendServer) VerifyEmailAddressPage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	token := req.URL.Query().Get("t")
	if token == "" {
		return page("Verify Email",
			verifyEmailContent("This verification link is invalid. Please check your email for the correct link or sign in to request a new one."),
		), nil
	}

	var unauthedClient client.Client
	var err error
	if s.developingLocally {
		unauthedClient, err = client.BuildUnauthenticatedGRPCClient(s.config.APIServiceConnection.GRPCAPIServerURL)
	} else {
		unauthedClient, err = client.BuildTLSGRPCClient(s.config.APIServiceConnection.GRPCAPIServerURL)
	}
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "building gRPC client")
		return page("Verify Email",
			verifyEmailContent("Unable to verify your email at this time. Please try again later or sign in to request a new verification email."),
		), nil
	}

	_, err = unauthedClient.VerifyEmailAddress(ctx, &authsvc.VerifyEmailAddressRequest{Token: token})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "verifying email address")
		return page("Verify Email",
			verifyEmailContent("This verification link has expired or is invalid. Please sign in to request a new verification email."),
		), nil
	}

	return page("Verify Email",
		verifyEmailContent("Your email has been verified. You can now sign in."),
	), nil
}

func verifyEmailContent(message string) g.Node {
	return ghtml.Div(
		ghtml.Class("space-y-6 text-center"),
		ghtml.H2(
			ghtml.Class("text-xl font-semibold"),
			g.Text("Verify Email"),
		),
		ghtml.P(
			ghtml.Class("text-gray-600"),
			g.Text(message),
		),
		ghtml.A(
			ghtml.Href("/login"),
			ghtml.Class("inline-block px-6 py-3 rounded-md bg-blue-600 text-white font-medium hover:bg-blue-700"),
			g.Text("Sign In"),
		),
	)
}

type aasaStruct struct {
	AppLinks       aasaAppLinks       `json:"applinks"`
	WebCredentials aasaWebCredentials `json:"webcredentials"`
}

type aasaWebCredentials struct {
	Apps []string `json:"apps"`
}

type aasaAppLinks struct {
	Apps    []string     `json:"apps"`
	Details []aasaDetail `json:"details"`
}

type aasaDetail struct {
	AppID string   `json:"appID"`
	Paths []string `json:"paths"`
}

// AppleAppSiteAssociationHandler serves the AASA file for iOS Universal Links.
func (s *ConsumerFrontendServer) AppleAppSiteAssociationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	cfg := s.config.AppleAppSiteAssociation
	appID := fmt.Sprintf("%s.%s", cfg.TeamID, cfg.BundleID)

	aasa := aasaStruct{
		AppLinks: aasaAppLinks{
			Apps: []string{},
			Details: []aasaDetail{
				{AppID: appID, Paths: []string{"/accept_invitation", "/accept_invitation/*"}},
			},
		},
		WebCredentials: aasaWebCredentials{Apps: []string{appID}},
	}

	s.encoder.EncodeResponseWithStatus(ctx, res, aasa, http.StatusOK)
}
