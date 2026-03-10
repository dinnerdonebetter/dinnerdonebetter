package components

import (
	"fmt"

	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"

	"google.golang.org/protobuf/types/known/timestamppb"
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// PasskeysContent renders the passkeys management page content.
func (r *ComponentRenderer) PasskeysContent(credentials []*authsvc.PasskeyCredential) g.Node {
	nodes := []g.Node{
		addPasskeySection(),
	}

	if len(credentials) > 0 {
		nodes = append(nodes, passkeysListSection(credentials))
	} else {
		nodes = append(nodes, ghtml.Div(
			ghtml.Class("p-4 rounded-lg border border-gray-200 bg-white"),
			ghtml.P(
				ghtml.Class("text-gray-500 text-sm"),
				g.Text("No passkeys yet. Add one to sign in quickly without a password."),
			),
		))
	}

	return ghtml.Div(
		ghtml.Class("space-y-6"),
		g.Group(nodes),
	)
}

func addPasskeySection() g.Node {
	return ghtml.Div(
		ghtml.Class("p-4 rounded-lg border border-gray-200 bg-white"),
		ghtml.Div(
			ghtml.Class("font-medium"),
			g.Text("Add passkey"),
		),
		ghtml.Div(
			ghtml.Class("text-sm text-gray-500 mt-0.5"),
			g.Text("Add a passkey to sign in quickly without a password."),
		),
		ghtml.Button(
			ghtml.Type("button"),
			ghtml.ID("add-passkey-btn"),
			ghtml.Class("mt-3 inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"),
			g.Text("Add passkey"),
		),
	)
}

func passkeysListSection(credentials []*authsvc.PasskeyCredential) g.Node {
	var cards []g.Node
	for _, c := range credentials {
		cards = append(cards, passkeyCard(c))
	}
	return ghtml.Div(
		ghtml.Class("space-y-3"),
		ghtml.H3(ghtml.Class("text-lg font-medium"), g.Text("Your passkeys")),
		ghtml.Div(
			ghtml.Class("space-y-3"),
			g.Group(cards),
		),
	)
}

func passkeyCard(c *authsvc.PasskeyCredential) g.Node {
	name := c.GetFriendlyName()
	if name == "" {
		name = "Passkey"
	}

	createdStr := formatPasskeyTimestamp(c.GetCreatedAt())
	lastUsedStr := ""
	if lu := c.GetLastUsedAt(); lu != nil {
		lastUsedStr = formatPasskeyTimestamp(lu)
	}

	details := createdStr
	if lastUsedStr != "" {
		details = fmt.Sprintf("%s · Last used %s", createdStr, lastUsedStr)
	}

	return ghtml.Div(
		ghtml.Class("p-4 rounded-lg border border-gray-200 bg-white flex items-center justify-between gap-4"),
		ghtml.Div(
			ghtml.Class("flex-1 min-w-0"),
			ghtml.Div(ghtml.Class("font-medium"), g.Text(name)),
			ghtml.Div(
				ghtml.Class("text-sm text-gray-500 mt-0.5"),
				g.Text(details),
			),
		),
		ghtml.Form(
			ghtml.Method("post"),
			ghtml.Action("/account/passkeys/delete"),
			ghtml.Class("flex-shrink-0"),
			ghtml.Input(ghtml.Type("hidden"), ghtml.Name("credential_id"), ghtml.Value(c.GetId())),
			ghtml.Button(
				ghtml.Type("submit"),
				ghtml.Class("text-sm px-3 py-1.5 text-red-600 hover:text-red-700 hover:bg-red-50 rounded-md border border-red-200 hover:border-red-300"),
				g.Text("Remove"),
			),
		),
	)
}

func formatPasskeyTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	t := ts.AsTime()
	return t.Format("Jan 2, 2006")
}
