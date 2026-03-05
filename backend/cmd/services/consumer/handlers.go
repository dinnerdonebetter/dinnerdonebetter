package main

import (
	"fmt"
	"net/http"
	"strings"

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

func (s *ConsumerFrontendServer) AccountSettingsPage(_ http.ResponseWriter, _ *http.Request) (g.Node, error) {
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
