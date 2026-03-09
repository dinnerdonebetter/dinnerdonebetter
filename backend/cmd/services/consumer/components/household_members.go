package components

import (
	"strings"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// HouseholdMembersFormErrors holds validation errors for the invitation form.
type HouseholdMembersFormErrors struct {
	Email string
}

// HouseholdMembersContent renders the main content for the household members page.
func (r *ComponentRenderer) HouseholdMembersContent(
	account *identitysvc.Account,
	invitations []*identitysvc.AccountInvitation,
	currentUserID string,
	isAdmin bool,
	baseURL string,
	formErrors *HouseholdMembersFormErrors,
) g.Node {
	if formErrors == nil {
		formErrors = &HouseholdMembersFormErrors{}
	}

	nodes := []g.Node{
		membersSection(account, currentUserID, isAdmin),
	}

	if isAdmin {
		nodes = append(nodes, sendInvitationSection(formErrors))
	}

	if len(invitations) > 0 {
		nodes = append(nodes, pendingInvitationsSection(invitations, isAdmin, baseURL))
	}

	return ghtml.Div(
		ghtml.Class("space-y-8"),
		g.Group(nodes),
	)
}

func membersSection(account *identitysvc.Account, currentUserID string, isAdmin bool) g.Node {
	content := ghtml.P(
		ghtml.Class("text-gray-500 text-sm"),
		g.Text("No members yet. Invite someone to join your household."),
	)

	if len(account.GetMembers()) > 0 {
		var memberCards []g.Node
		for _, m := range account.GetMembers() {
			memberCards = append(memberCards, memberCard(m, currentUserID, isAdmin))
		}
		content = ghtml.Div(ghtml.Class("space-y-3"), g.Group(memberCards))
	}

	return ghtml.Div(
		ghtml.Class("space-y-3"),
		ghtml.H3(ghtml.Class("text-lg font-medium"), g.Text("Household Members")),
		content,
	)
}

func memberCard(m *identitysvc.AccountUserMembershipWithUser, currentUserID string, isAdmin bool) g.Node {
	displayName := memberDisplayName(m)
	isYou := m.BelongsToUser != nil && m.BelongsToUser.Id == currentUserID

	roleLabel := "Member"
	if m.AccountRole == authorization.AccountAdminRoleName {
		roleLabel = "Admin"
	}

	roleNode := ghtml.Span(
		ghtml.Class("text-sm text-gray-500 px-2 py-0.5 rounded bg-gray-100"),
		g.Text(roleLabel),
	)

	if isAdmin && !isYou && m.BelongsToUser != nil && m.BelongsToUser.Id != "" {
		currentRole := m.AccountRole
		roleNode = ghtml.Form(
			ghtml.Method("post"),
			ghtml.Action("/account/household-members/update-role"),
			ghtml.Class("inline-flex items-center gap-2 flex-wrap"),
			ghtml.Select(
				ghtml.Name("new_role"),
				ghtml.Class("text-sm border border-gray-300 rounded px-2 py-1"),
				ghtml.Option(g.If(currentRole == authorization.AccountMemberRoleName, ghtml.Selected()), ghtml.Value(authorization.AccountMemberRoleName), g.Text("Member")),
				ghtml.Option(g.If(currentRole == authorization.AccountAdminRoleName, ghtml.Selected()), ghtml.Value(authorization.AccountAdminRoleName), g.Text("Admin")),
			),
			ghtml.Input(ghtml.Type("hidden"), ghtml.Name("user_id"), ghtml.Value(m.BelongsToUser.Id)),
			ghtml.Input(ghtml.Type(inputTypeText), ghtml.Name("reason"), ghtml.Placeholder("Reason (required)"), ghtml.Class("text-sm border border-gray-300 rounded px-2 py-1 min-w-[120px]")),
			ghtml.Button(ghtml.Type("submit"), ghtml.Class("text-sm px-2 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"), g.Text("Update")),
		)
	}

	return ghtml.Div(
		ghtml.Class("p-4 rounded-lg border border-gray-200 bg-white flex items-center justify-between gap-4"),
		ghtml.Div(
			ghtml.Class("flex-1 min-w-0"),
			ghtml.Div(ghtml.Class("font-medium"), g.Text(displayName)),
			g.If(isYou, ghtml.Span(ghtml.Class("text-xs text-gray-500"), g.Text("(You)"))),
		),
		roleNode,
	)
}

func memberDisplayName(m *identitysvc.AccountUserMembershipWithUser) string {
	if m.BelongsToUser == nil {
		return "Unknown User"
	}
	u := m.BelongsToUser
	if u.FirstName != "" {
		if u.LastName != "" {
			return u.FirstName + " " + u.LastName
		}
		return u.FirstName
	}
	if u.Username != "" {
		return u.Username
	}
	return "Unknown User"
}

func sendInvitationSection(errors *HouseholdMembersFormErrors) g.Node {
	return ghtml.Div(
		ghtml.Class("space-y-3"),
		ghtml.H3(ghtml.Class("text-lg font-medium"), g.Text("Add Someone to Your Household")),
		ghtml.P(ghtml.Class("text-sm text-gray-500"), g.Text("Send an invitation by email. They can join once they have an account.")),
		ghtml.Form(
			ghtml.Method("post"),
			ghtml.Action("/account/household-members/send-invitation"),
			ghtml.Class("space-y-3"),
			ghtml.Div(
				ghtml.Class("space-y-1"),
				ghtml.Label(ghtml.For("email"), ghtml.Class("block text-sm font-medium"), g.Text("Email Address")),
				ghtml.Input(
					ghtml.Type("email"),
					ghtml.ID("email"),
					ghtml.Name("email"),
					ghtml.Required(),
					ghtml.Class("block w-full rounded-md border border-gray-300 px-3 py-2 text-sm"),
				),
				g.If(errors.Email != "", ghtml.Span(ghtml.Class("text-sm text-red-600"), g.Text(errors.Email))),
			),
			ghtml.Div(
				ghtml.Class("space-y-1"),
				ghtml.Label(ghtml.For("name"), ghtml.Class("block text-sm font-medium"), g.Text("Name (Optional)")),
				ghtml.Input(
					ghtml.Type(inputTypeText),
					ghtml.ID("name"),
					ghtml.Name("name"),
					ghtml.Class("block w-full rounded-md border border-gray-300 px-3 py-2 text-sm"),
				),
			),
			ghtml.Div(
				ghtml.Class("space-y-1"),
				ghtml.Label(ghtml.For("note"), ghtml.Class("block text-sm font-medium"), g.Text("Note (Optional)")),
				ghtml.Textarea(
					ghtml.ID("note"),
					ghtml.Name("note"),
					ghtml.Class("block w-full rounded-md border border-gray-300 px-3 py-2 text-sm"),
				),
			),
			ghtml.Button(
				ghtml.Type("submit"),
				ghtml.Class("px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm font-medium"),
				g.Text("Send Invitation"),
			),
		),
	)
}

func pendingInvitationsSection(invitations []*identitysvc.AccountInvitation, isAdmin bool, baseURL string) g.Node {
	var cards []g.Node
	for _, inv := range invitations {
		cards = append(cards, invitationCard(inv, isAdmin, baseURL))
	}

	return ghtml.Div(
		ghtml.Class("space-y-3"),
		ghtml.H3(ghtml.Class("text-lg font-medium"), g.Text("Invitations")),
		ghtml.P(ghtml.Class("text-sm text-gray-500"), g.Text("Invitations you've sent for this household and their status.")),
		ghtml.Div(ghtml.Class("space-y-3"), g.Group(cards)),
	)
}

func invitationCard(inv *identitysvc.AccountInvitation, isAdmin bool, baseURL string) g.Node {
	status := inv.Status
	if status == "" {
		status = "pending"
	}

	inviteURL := baseURL + "/accept_invitation?i=" + inv.Id + "&t=" + inv.Token

	nodes := []g.Node{
		ghtml.Div(
			ghtml.Class("flex-1 min-w-0"),
			ghtml.Div(ghtml.Class("font-medium"), g.Text(inv.ToEmail)),
			g.If(inv.ToName != "", ghtml.Div(ghtml.Class("text-sm text-gray-500"), g.Text(inv.ToName))),
			ghtml.Span(
				ghtml.Class("inline-block mt-1 text-xs px-2 py-0.5 rounded bg-gray-100"),
				g.Text(status),
			),
		),
	}

	if strings.EqualFold(status, "pending") {
		actionNodes := []g.Node{
			copyLinkButton(inviteURL),
		}
		if isAdmin {
			actionNodes = append(actionNodes, cancelInvitationButton(inv.Id))
		}
		nodes = append(nodes, ghtml.Div(ghtml.Class("flex gap-2"), g.Group(actionNodes)))
	}

	return ghtml.Div(
		ghtml.Class("p-4 rounded-lg border border-gray-200 bg-white flex items-center justify-between gap-4"),
		g.Group(nodes),
	)
}

func copyLinkButton(url string) g.Node {
	return ghtml.Button(
		ghtml.Type("button"),
		ghtml.Class("text-sm px-2 py-1 border border-gray-300 rounded hover:bg-gray-50 copy-invite-link"),
		g.Attr("data-url", url),
		g.Text("Copy Link"),
	)
}

func cancelInvitationButton(invitationID string) g.Node {
	return ghtml.Form(
		ghtml.Method("post"),
		ghtml.Action("/account/household-members/cancel-invitation"),
		ghtml.Class("inline"),
		ghtml.Input(ghtml.Type("hidden"), ghtml.Name("invitation_id"), ghtml.Value(invitationID)),
		ghtml.Button(
			ghtml.Type("submit"),
			ghtml.Class("text-sm px-2 py-1 text-red-600 border border-red-300 rounded hover:bg-red-50"),
			g.Text("Cancel"),
		),
	)
}
