package components

import (
	"strings"

	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// ConfigurableSetting represents a service setting with its current value for display.
type ConfigurableSetting struct {
	Setting      *settingssvc.ServiceSetting
	Config       *settingssvc.ServiceSettingConfiguration
	CurrentValue string
}

// PreferencesContent renders the preferences (service settings) form.
func (r *ComponentRenderer) PreferencesContent(settings []ConfigurableSetting, errorMsg string) g.Node {
	if len(settings) == 0 && errorMsg == "" {
		return emptyPreferencesState()
	}

	if errorMsg != "" {
		return ghtml.Div(
			ghtml.Class("space-y-4"),
			ghtml.Div(
				ghtml.Class("p-3 rounded-md text-sm bg-red-50 text-red-800"),
				g.Text(errorMsg),
			),
			ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
		)
	}

	var rows []g.Node
	for i := range settings {
		rows = append(rows, settingRow(settings[i]))
	}

	return ghtml.Div(
		ghtml.Class("space-y-6"),
		ghtml.Div(
			ghtml.Class("space-y-4"),
			g.Group(rows),
		),
	)
}

func emptyPreferencesState() g.Node {
	return ghtml.Div(
		ghtml.Class("space-y-4 text-center py-8"),
		ghtml.P(
			ghtml.Class("text-gray-500"),
			g.Text("No preferences to configure."),
		),
		ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
	)
}

func settingRow(item ConfigurableSetting) g.Node {
	if len(item.Setting.GetEnumeration()) == 0 {
		return g.El("")
	}

	humanName := humanReadableSettingName(item.Setting.GetName())
	desc := item.Setting.GetDescription()
	if desc != "" {
		humanName = humanName + " — " + desc
	}

	formID := "pref-" + item.Setting.GetId()
	action := "/account/preferences/update"

	return ghtml.Div(
		ghtml.Class("p-4 rounded-lg border border-gray-200 bg-white"),
		ghtml.Form(
			ghtml.Method("POST"),
			ghtml.Action(action),
			ghtml.ID(formID),
			ghtml.Class("space-y-3"),
			ghtml.Div(
				ghtml.Class("flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3"),
				ghtml.Label(
					ghtml.For(formID+"-value"),
					ghtml.Class("text-sm font-medium text-gray-700"),
					g.Text(humanName),
				),
				ghtml.Div(
					ghtml.Class("flex items-center gap-2"),
					ghtml.Select(
						ghtml.Name("value"),
						ghtml.ID(formID+"-value"),
						ghtml.Class("rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500"),
						g.Attr("onchange", "this.form.submit()"),
						enumerationOptions(item.Setting.GetEnumeration(), item.CurrentValue),
					),
				),
			),
			ghtml.Input(
				ghtml.Type("hidden"),
				ghtml.Name("setting_id"),
				ghtml.Value(item.Setting.GetId()),
			),
			g.If(item.Config != nil && item.Config.GetId() != "",
				ghtml.Input(
					ghtml.Type("hidden"),
					ghtml.Name("config_id"),
					ghtml.Value(item.Config.GetId()),
				),
			),
			g.If(desc != "",
				ghtml.P(ghtml.Class("text-xs text-gray-500"), g.Text(desc)),
			),
		),
	)
}

func enumerationOptions(enum []string, current string) g.Node {
	var opts []g.Node
	for _, v := range enum {
		sel := ""
		if v == current {
			sel = "selected"
		}
		opts = append(opts, ghtml.Option(
			ghtml.Value(v),
			g.If(sel != "", ghtml.Selected()),
			g.Text(humanReadableOption(v)),
		))
	}
	return g.Group(opts)
}

func humanReadableSettingName(name string) string {
	return strings.ReplaceAll(name, "_", " ")
}

func humanReadableOption(v string) string {
	if v == "" {
		return v
	}
	return strings.ToUpper(v[:1]) + strings.ToLower(v[1:])
}
