package main

import (
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/components"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	webappauth "github.com/dinnerdonebetter/backend/internal/webapp/auth"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (s *ConsumerFrontendServer) PreferencesPage(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		return page("Preferences",
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
		return page("Preferences",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.H2(ghtml.Class("text-xl font-semibold"), g.Text("Preferences")),
				ghtml.P(
					ghtml.Class("text-gray-600"),
					g.Text("Create or join a household to configure your preferences."),
				),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	accountID := activeRes.Result.Id
	filter := &filtering.QueryFilter{MaxResponseSize: new(uint32(100))}

	settingsRes, err := c.GetServiceSettings(ctx, &settingssvc.GetServiceSettingsRequest{Filter: filter})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting service settings")
		return page("Preferences",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(ghtml.Class("text-red-600"), g.Text("Unable to load preferences. Please try again.")),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	configsRes, err := c.GetServiceSettingConfigurationsForUser(ctx, &settingssvc.GetServiceSettingConfigurationsForUserRequest{Filter: filter})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting service setting configurations")
		return page("Preferences",
			ghtml.Div(
				ghtml.Class("space-y-6"),
				ghtml.P(ghtml.Class("text-red-600"), g.Text("Unable to load preferences. Please try again.")),
				ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline"), g.Text("Back to Account Settings")),
			),
		), nil
	}

	configurableSettings := mergeAndFilterSettings(
		settingsRes.GetResults(),
		configsRes.GetResults(),
		accountID,
	)

	flashMsg := ""
	if req.URL.Query().Get("updated") == "1" {
		flashMsg = "Preferences updated."
	}
	errorMsg := ""
	if e := req.URL.Query().Get("error"); e != "" {
		errorMsg = preferenceErrorForParam(e)
	}

	return page("Preferences",
		ghtml.Div(
			ghtml.Class("space-y-6"),
			ghtml.H2(ghtml.Class("text-xl font-semibold"), g.Text("Preferences")),
			ghtml.A(ghtml.Href("/account/settings"), ghtml.Class("text-blue-600 hover:underline text-sm"), g.Text("Back to Account Settings")),
			g.If(flashMsg != "", ghtml.Div(ghtml.Class("p-3 rounded-md text-sm bg-green-50 text-green-800"), g.Text(flashMsg))),
			s.componentRenderer.PreferencesContent(configurableSettings, errorMsg),
		),
	), nil
}

func (s *ConsumerFrontendServer) UpdatePreferenceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if req.Method != http.MethodPost {
		http.Redirect(res, req, "/account/preferences", http.StatusFound)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Redirect(res, req, "/account/preferences?error="+errorParamInvalid, http.StatusFound)
		return
	}

	settingID := strings.TrimSpace(req.FormValue("setting_id"))
	configID := strings.TrimSpace(req.FormValue("config_id"))
	value := strings.TrimSpace(req.FormValue("value"))

	if settingID == "" || value == "" {
		http.Redirect(res, req, "/account/preferences?error="+errorParamInvalid, http.StatusFound)
		return
	}

	c, err := webappauth.ClientFromContext(req.Context())
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting API client")
		http.Redirect(res, req, "/account/preferences?error="+errorParamServer, http.StatusFound)
		return
	}

	if configID != "" {
		_, err = c.UpdateServiceSettingConfiguration(ctx, &settingssvc.UpdateServiceSettingConfigurationRequest{
			ServiceSettingConfigurationId: configID,
			Input: &settingssvc.ServiceSettingConfigurationUpdateRequestInput{
				Value: &value,
			},
		})
	} else {
		_, err = c.CreateServiceSettingConfiguration(ctx, &settingssvc.CreateServiceSettingConfigurationRequest{
			Input: &settingssvc.ServiceSettingConfigurationCreationRequestInput{
				ServiceSettingId: settingID,
				Value:            value,
				Notes:            "",
			},
		})
	}

	if err != nil {
		observability.AcknowledgeError(err, logger, span, "updating preference")
		http.Redirect(res, req, "/account/preferences?error="+errorParamUpdateFailed, http.StatusFound)
		return
	}

	http.Redirect(res, req, "/account/preferences?updated=1", http.StatusFound)
}

func mergeAndFilterSettings(
	settings []*settingssvc.ServiceSetting,
	configs []*settingssvc.ServiceSettingConfiguration,
	accountID string,
) []components.ConfigurableSetting {
	configsForAccount := make(map[string]*settingssvc.ServiceSettingConfiguration)
	for _, cfg := range configs {
		if cfg.GetBelongsToAccount() == accountID {
			if ss := cfg.GetServiceSetting(); ss != nil {
				configsForAccount[ss.GetId()] = cfg
			}
		}
	}

	var result []components.ConfigurableSetting
	for _, setting := range settings {
		if setting.GetType() != "user" || setting.GetAdminsOnly() {
			continue
		}
		if len(setting.GetEnumeration()) == 0 {
			continue
		}

		config := configsForAccount[setting.GetId()]
		currentValue := ""
		switch {
		case config != nil && config.GetValue() != "":
			currentValue = config.GetValue()
		case setting.GetDefaultValue() != "":
			currentValue = setting.GetDefaultValue()
		case len(setting.GetEnumeration()) > 0:
			currentValue = setting.GetEnumeration()[0]
		}
		if currentValue == "" && len(setting.GetEnumeration()) > 0 {
			currentValue = setting.GetEnumeration()[0]
		}

		result = append(result, components.ConfigurableSetting{
			Setting:      setting,
			Config:       config,
			CurrentValue: currentValue,
		})
	}
	return result
}

func preferenceErrorForParam(e string) string {
	switch e {
	case errorParamInvalid:
		return "Invalid input. Please try again."
	case errorParamUpdateFailed:
		return "Failed to save preference. Please try again."
	case errorParamServer:
		return errorMsgServer
	default:
		return errorMsgSomethingWrong
	}
}
