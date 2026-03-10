package components

import (
	"github.com/dinnerdonebetter/backend/cmd/services/consumer/design"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// AccountLinksProps holds data for rendering account navigation links.
type AccountLinksProps struct {
	HasAccount bool // whether the user has at least one account (household)
}

// AccountLinks renders card-style navigation links for the account page.
func (r *ComponentRenderer) AccountLinks(props *AccountLinksProps) g.Node {
	if props == nil {
		props = &AccountLinksProps{}
	}
	palette := &design.StandardPalette
	linkClass := "block p-4 rounded-lg border border-gray-200 bg-white hover:bg-gray-50 transition-colors " +
		design.TextColor(palette.Text)

	links := []g.Node{
		accountLink("/account/household-members", "Household Members", "Members and invitations", "person.2", linkClass),
	}
	if props.HasAccount {
		links = append(links, accountLink("/account/household-details", "Household Details", "Edit household details", "house", linkClass))
	}
	links = append(links,
		accountLink("/account/passkeys", "Passkeys", "Add, view, and remove passkeys", "key", linkClass),
		accountLink("/account/preferences", "Preferences", "Configure preferences", "gearshape", linkClass),
		accountLink("/account/profile", "Profile", "Photo, name, and account details", "person.crop.circle", linkClass),
	)

	return ghtml.Div(
		ghtml.Class("space-y-4"),
		g.Group(links),
	)
}

func accountLink(href, title, subtitle, _, linkClass string) g.Node {
	return ghtml.A(
		ghtml.Href(href),
		ghtml.Class(linkClass),
		ghtml.Div(
			ghtml.Class("font-medium"),
			g.Text(title),
		),
		ghtml.Div(
			ghtml.Class("text-sm text-gray-500 mt-0.5"),
			g.Text(subtitle),
		),
	)
}
