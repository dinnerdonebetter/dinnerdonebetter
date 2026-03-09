package components

import (
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// ProfileFormErrors holds validation errors for the profile form.
type ProfileFormErrors struct {
	Username    string
	FirstName   string
	LastName    string
	Password    string
	DetailsForm string
}

// ProfileContent renders the profile update form.
func (r *ComponentRenderer) ProfileContent(
	user *identitysvc.User,
	formErrors *ProfileFormErrors,
) g.Node {
	if formErrors == nil {
		formErrors = &ProfileFormErrors{}
	}

	username := ""
	firstName := ""
	lastName := ""
	birthdayStr := ""
	if user != nil {
		username = user.GetUsername()
		firstName = user.GetFirstName()
		lastName = user.GetLastName()
		if user.GetBirthday() != nil {
			t := user.GetBirthday().AsTime()
			birthdayStr = t.Format("2006-01-02")
		}
	}

	return ghtml.Div(
		ghtml.Class("space-y-6"),
		// Username section
		ghtml.Div(
			ghtml.Class("p-4 rounded-lg border border-gray-200 bg-white space-y-4"),
			ghtml.H3(ghtml.Class("text-sm font-medium text-gray-700"), g.Text("Username")),
			ghtml.Form(
				ghtml.Method("POST"),
				ghtml.Action("/account/profile/update-username"),
				ghtml.Class("space-y-3"),
				profileFormField("username", "Username", username, true, formErrors.Username),
				ghtml.Button(
					ghtml.Type("submit"),
					ghtml.Class("px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm font-medium"),
					g.Text("Update Username"),
				),
			),
		),
		// User details section (first name, last name, birthday)
		ghtml.Div(
			ghtml.Class("p-4 rounded-lg border border-gray-200 bg-white space-y-4"),
			ghtml.H3(ghtml.Class("text-sm font-medium text-gray-700"), g.Text("Profile Details")),
			ghtml.P(ghtml.Class("text-sm text-gray-500"), g.Text("Updating your name or birthday requires your password for security.")),
			ghtml.Form(
				ghtml.Method("POST"),
				ghtml.Action("/account/profile/update-details"),
				ghtml.Class("space-y-4"),
				profileFormField("first_name", "First Name", firstName, true, formErrors.FirstName),
				profileFormField("last_name", "Last Name", lastName, false, formErrors.LastName),
				profileDateField("birthday", "Birthday", birthdayStr),
				profileFormField("current_password", "Current Password", "", true, formErrors.Password),
				profileFormField("totp_token", "Authenticator Code (if enabled)", "", false, ""),
				g.If(formErrors.DetailsForm != "", ghtml.Div(ghtml.Class("text-sm text-red-600"), g.Text(formErrors.DetailsForm))),
				ghtml.Button(
					ghtml.Type("submit"),
					ghtml.Class("px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm font-medium"),
					g.Text("Update Details"),
				),
			),
		),
		// Profile photo note
		ghtml.Div(
			ghtml.Class("p-4 rounded-lg border border-gray-200 bg-gray-50"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-600"),
				g.Text("Profile photo can be updated in the Dinner Done Better app."),
			),
		),
	)
}

func profileFormField(id, label, value string, required bool, errorMsg string) g.Node {
	inputType := inputTypeText
	if id == "current_password" {
		inputType = "password"
	}
	if id == "totp_token" {
		inputType = inputTypeText
	}

	inputAttrs := []g.Node{
		ghtml.Type(inputType),
		ghtml.ID(id),
		ghtml.Name(id),
		ghtml.Class("block w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500"),
	}
	if value != "" && inputType != "password" {
		inputAttrs = append(inputAttrs, ghtml.Value(value))
	}
	if required {
		inputAttrs = append(inputAttrs, ghtml.Required())
	}
	if id == "totp_token" {
		inputAttrs = append(inputAttrs, ghtml.Placeholder("Optional"))
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

func profileDateField(id, label, value string) g.Node {
	return ghtml.Div(
		ghtml.Class("space-y-1"),
		ghtml.Label(
			ghtml.For(id),
			ghtml.Class("block text-sm font-medium text-gray-700"),
			g.Text(label),
		),
		ghtml.Input(
			ghtml.Type("date"),
			ghtml.ID(id),
			ghtml.Name(id),
			ghtml.Value(value),
			ghtml.Class("block w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500"),
		),
	)
}
