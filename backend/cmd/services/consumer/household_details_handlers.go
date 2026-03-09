package main

import (
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/components"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	webappauth "github.com/dinnerdonebetter/backend/internal/webapp/auth"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (s *ConsumerFrontendServer) HouseholdDetailsPage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		return page("Household Details",
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
		return page("Household Details",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.H2(ghtml.Class("text-xl font-semibold"), g.Text("Household Details")),
				ghtml.P(
					ghtml.Class("text-gray-600"),
					g.Text("Create or join a household to edit household details."),
				),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	account := activeRes.Result
	selfRes, err := c.GetSelf(ctx, &authsvc.GetSelfRequest{})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting self")
		return page("Household Details",
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
	for _, m := range account.GetMembers() {
		if m.BelongsToUser != nil && m.BelongsToUser.Id == currentUserID {
			if m.AccountRole == authorization.AccountAdminRoleName {
				isAdmin = true
			}
			break
		}
	}

	formErrors := (*components.HouseholdDetailsFormErrors)(nil)
	if req.URL.Query().Get("error") == errorParamInvalidName {
		formErrors = &components.HouseholdDetailsFormErrors{Name: "Household name is required."}
	}

	flashMsg := ""
	if req.URL.Query().Get("updated") == "1" {
		flashMsg = "Household details updated."
	}
	if e := req.URL.Query().Get("error"); e != "" && flashMsg == "" {
		flashMsg = householdDetailsErrorForParam(e)
	}

	return page("Household Details",
		ghtml.Div(
			ghtml.Class("space-y-6"),
			ghtml.H2(ghtml.Class("text-xl font-semibold"), g.Text("Household Details")),
			ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline text-sm"), g.Text("Back to Account Settings")),
			g.If(flashMsg != "", buildHouseholdDetailsFlashMessage(flashMsg)),
			s.componentRenderer.HouseholdDetailsContent(account, isAdmin, formErrors),
		),
	), nil
}

func buildHouseholdDetailsFlashMessage(flashMsg string) g.Node {
	cls := "p-3 rounded-md text-sm bg-red-50 text-red-800"
	if strings.HasPrefix(flashMsg, "Household") {
		cls = "p-3 rounded-md text-sm bg-green-50 text-green-800"
	}
	return ghtml.Div(ghtml.Class(cls), g.Text(flashMsg))
}

func householdDetailsErrorForParam(e string) string {
	switch e {
	case errorParamInvalid, errorParamInvalidName:
		return errorMsgInvalidInput
	case errorParamUpdateFailed:
		return "Failed to save household details. Please try again."
	case errorParamServer:
		return errorMsgServer
	default:
		return errorMsgSomethingWrong
	}
}

func (s *ConsumerFrontendServer) UpdateHouseholdDetailsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if req.Method != http.MethodPost {
		http.Redirect(res, req, "/account/household-details", http.StatusFound)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Redirect(res, req, "/account/household-details?error="+errorParamInvalid, http.StatusFound)
		return
	}

	name := strings.TrimSpace(req.FormValue("name"))
	if name == "" {
		http.Redirect(res, req, "/account/household-details?error="+errorParamInvalidName, http.StatusFound)
		return
	}

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		http.Redirect(res, req, "/account/household-details?error="+errorParamServer, http.StatusFound)
		return
	}

	activeRes, err := c.GetActiveAccount(ctx, &authsvc.GetActiveAccountRequest{})
	if err != nil || activeRes == nil || activeRes.Result == nil {
		observability.AcknowledgeError(err, logger, span, "getting active account")
		http.Redirect(res, req, "/account/household-details?error="+errorParamServer, http.StatusFound)
		return
	}

	accountID := activeRes.Result.Id

	contactPhone := strings.TrimSpace(req.FormValue("contact_phone"))
	addressLine1 := strings.TrimSpace(req.FormValue("address_line_1"))
	addressLine2 := strings.TrimSpace(req.FormValue("address_line_2"))
	city := strings.TrimSpace(req.FormValue("city"))
	state := strings.TrimSpace(req.FormValue("state"))
	zipCode := strings.TrimSpace(req.FormValue("zip_code"))
	country := strings.TrimSpace(req.FormValue("country"))

	input := &identitysvc.AccountUpdateRequestInput{
		Name:         &name,
		ContactPhone: &contactPhone,
		AddressLine1: &addressLine1,
		AddressLine2: &addressLine2,
		City:         &city,
		State:        &state,
		ZipCode:      &zipCode,
		Country:      &country,
	}

	_, err = c.UpdateAccount(ctx, &identitysvc.UpdateAccountRequest{
		AccountId: accountID,
		Input:     input,
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "updating household details")
		http.Redirect(res, req, "/account/household-details?error="+errorParamUpdateFailed, http.StatusFound)
		return
	}

	http.Redirect(res, req, "/account/household-details?updated=1", http.StatusFound)
}
