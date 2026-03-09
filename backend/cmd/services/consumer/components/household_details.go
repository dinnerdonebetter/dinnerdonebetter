package components

import (
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const inputTypeText = "text"

// HouseholdDetailsFormErrors holds validation errors for the household details form.
type HouseholdDetailsFormErrors struct {
	Name string
}

// HouseholdDetailsContent renders the form for editing household details.
func (r *ComponentRenderer) HouseholdDetailsContent(
	account *identitysvc.Account,
	isAdmin bool,
	formErrors *HouseholdDetailsFormErrors,
) g.Node {
	if formErrors == nil {
		formErrors = &HouseholdDetailsFormErrors{}
	}

	if !isAdmin {
		return ghtml.Div(
			ghtml.Class("space-y-4"),
			ghtml.P(
				ghtml.Class("text-gray-600"),
				g.Text("Only household admins can edit household details."),
			),
		)
	}

	return ghtml.Div(
		ghtml.Class("space-y-6"),
		ghtml.Form(
			ghtml.Method("POST"),
			ghtml.Action("/account/household-details/update"),
			ghtml.Class("space-y-6"),
			householdDetailsFormFields(account, formErrors),
			ghtml.Button(
				ghtml.Type("submit"),
				ghtml.Class("px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm font-medium"),
				g.Text("Save Changes"),
			),
		),
	)
}

func householdDetailsFormFields(account *identitysvc.Account, errors *HouseholdDetailsFormErrors) g.Node {
	name := ""
	contactPhone := ""
	addressLine1 := ""
	addressLine2 := ""
	city := ""
	state := ""
	zipCode := ""
	country := ""
	if account != nil {
		name = account.GetName()
		contactPhone = account.GetContactPhone()
		addressLine1 = account.GetAddressLine1()
		addressLine2 = account.GetAddressLine2()
		city = account.GetCity()
		state = account.GetState()
		zipCode = account.GetZipCode()
		country = account.GetCountry()
	}

	return ghtml.Div(
		ghtml.Class("space-y-4"),
		// Household name
		formField("name", "Household Name", name, true, errors.Name),
		// Contact phone
		formField("contact_phone", "Contact Phone", contactPhone, false, ""),
		// Address section
		ghtml.Div(
			ghtml.Class("space-y-4 p-4 rounded-lg border border-gray-200 bg-white"),
			ghtml.H4(ghtml.Class("text-sm font-medium text-gray-700"), g.Text("Address")),
			formField("address_line_1", "Address Line 1", addressLine1, false, ""),
			formField("address_line_2", "Address Line 2", addressLine2, false, ""),
			formField("city", "City", city, false, ""),
			formField("state", "State / Province", state, false, ""),
			formField("zip_code", "ZIP / Postal Code", zipCode, false, ""),
			formField("country", "Country", country, false, ""),
		),
	)
}

func formField(id, label, value string, required bool, errorMsg string) g.Node {
	inputType := inputTypeText
	if id == "contact_phone" {
		inputType = "tel"
	}

	inputAttrs := []g.Node{
		ghtml.Type(inputType),
		ghtml.ID(id),
		ghtml.Name(id),
		ghtml.Value(value),
		ghtml.Class("block w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500"),
	}
	if required {
		inputAttrs = append(inputAttrs, ghtml.Required())
	}

	nodes := []g.Node{
		ghtml.Label(
			ghtml.For(id),
			ghtml.Class("block text-sm font-medium text-gray-700"),
			g.Text(label),
		),
		ghtml.Input(inputAttrs...),
	}
	if errorMsg != "" {
		nodes = append(nodes, ghtml.Span(ghtml.Class("text-sm text-red-600"), g.Text(errorMsg)))
	}

	return ghtml.Div(ghtml.Class("space-y-1"), g.Group(nodes))
}
