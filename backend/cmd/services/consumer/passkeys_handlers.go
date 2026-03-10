package main

import (
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/components"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	webappauth "github.com/dinnerdonebetter/backend/internal/webapp/auth"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (s *ConsumerFrontendServer) PasskeysPage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		return page("Passkeys",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(ghtml.Class("text-red-600"), g.Text("Unable to load. Please try again.")),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	listRes, err := c.ListPasskeys(ctx, &authsvc.ListPasskeysRequest{})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "listing passkeys")
		return page("Passkeys",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(ghtml.Class("text-red-600"), g.Text("Unable to load passkeys. Please try again.")),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	credentials := listRes.GetResults()
	queryParams := req.URL.Query()
	flashMsg := ""
	if queryParams.Get("deleted") == "1" {
		flashMsg = "Passkey removed successfully."
	}
	if e := queryParams.Get("error"); e != "" && flashMsg == "" {
		flashMsg = passkeyErrorForParam(e)
	}

	return page("Passkeys",
		ghtml.Div(
			ghtml.Class("space-y-6"),
			ghtml.H2(ghtml.Class("text-xl font-semibold"), g.Text("Passkeys")),
			ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline text-sm"), g.Text("Back to Account Settings")),
			g.If(flashMsg != "", buildFlashMessage(flashMsg)),
			s.componentRenderer.PasskeysContent(credentials),
			ghtml.Script(ghtml.Type("text/javascript"), g.Raw(components.PasskeyRegistrationScript)),
		),
	), nil
}

func (s *ConsumerFrontendServer) DeletePasskeyHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if req.Method != http.MethodPost {
		http.Redirect(res, req, "/account/passkeys", http.StatusFound)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Redirect(res, req, "/account/passkeys?error=invalid", http.StatusFound)
		return
	}

	credentialID := strings.TrimSpace(req.FormValue("credential_id"))
	if credentialID == "" {
		http.Redirect(res, req, "/account/passkeys?error=invalid", http.StatusFound)
		return
	}

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		http.Redirect(res, req, "/account/passkeys?error=server", http.StatusFound)
		return
	}

	_, err = c.ArchivePasskey(ctx, &authsvc.ArchivePasskeyRequest{
		CredentialId: credentialID,
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving passkey")
		http.Redirect(res, req, "/account/passkeys?error=delete_failed", http.StatusFound)
		return
	}

	http.Redirect(res, req, "/account/passkeys?deleted=1", http.StatusFound)
}

func passkeyErrorForParam(e string) string {
	switch e {
	case "invalid":
		return "Invalid request."
	case "delete_failed":
		return "Failed to remove passkey. Please try again."
	case "server":
		return "Something went wrong. Please try again."
	default:
		return "Something went wrong. Please try again."
	}
}
