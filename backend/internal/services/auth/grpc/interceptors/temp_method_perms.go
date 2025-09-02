package interceptors

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
)

var (
	noPerms = []authorization.Permission{}

	// TODO: ensure this map doesn't end up with configs for methods that don't exist

	methodPermissions = map[string][]authorization.Permission{
		"/mealplanning.MealPlanningService/CreateValidIngredient": {
			authorization.CreateValidIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredient": {
			authorization.ReadValidIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredients": {
			authorization.ReadValidIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/SearchForValidIngredients": {
			authorization.ReadValidIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/UpdateValidIngredient": {
			authorization.UpdateValidIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveValidIngredient": {
			authorization.ArchiveValidIngredientsPermission,
		},
		"/settings.SettingsService/CreateServiceSetting": {
			authorization.CreateServiceSettingsPermission,
		},
		"/settings.SettingsService/GetServiceSetting": {
			authorization.ReadServiceSettingsPermission,
		},
		"/settings.SettingsService/GetServiceSettings": {
			authorization.ReadServiceSettingsPermission,
		},
		"/settings.SettingsService/SearchForServiceSettings": {
			authorization.ReadServiceSettingsPermission,
		},
		"/settings.SettingsService/ArchiveServiceSetting": {
			authorization.ArchiveServiceSettingsPermission,
		},
		"/settings.SettingsService/CreateServiceSettingConfiguration": {
			authorization.CreateServiceSettingConfigurationsPermission,
		},
		"/settings.SettingsService/GetServiceSettingConfigurationByName": {
			authorization.ReadServiceSettingConfigurationsPermission,
		},
		"/settings.SettingsService/GetServiceSettingConfigurationsForAccount": {
			authorization.ReadServiceSettingConfigurationsPermission,
		},
		"/settings.SettingsService/GetServiceSettingConfigurationsForUser": {
			authorization.ReadServiceSettingConfigurationsPermission,
		},
		"/settings.SettingsService/ArchiveServiceSettingConfiguration": {
			authorization.ArchiveServiceSettingConfigurationsPermission,
		},
		"/oauth.OAuthService/CreateOAuth2Client": {
			authorization.CreateOAuth2ClientsPermission,
		},
		"/oauth.OAuthService/GetOAuth2Client": {
			authorization.ReadOAuth2ClientsPermission,
		},
		"/oauth.OAuthService/GetOAuth2Clients": {
			authorization.ReadOAuth2ClientsPermission,
		},
		"/oauth.OAuthService/ArchiveOAuth2Client": {
			authorization.ArchiveOAuth2ClientsPermission,
		},
		"/mealplanning.MealPlanningService/GetRandomValidIngredient": {
			authorization.ReadValidIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidIngredientGroup": {
			authorization.CreateValidIngredientGroupsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientGroup": {
			authorization.ReadValidIngredientGroupsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientGroups": {
			authorization.ReadValidIngredientGroupsPermission,
		},
		"/mealplanning.MealPlanningService/SearchForValidIngredientGroups": {
			authorization.ReadValidIngredientGroupsPermission,
		},
		"/mealplanning.MealPlanningService/UpdateValidIngredientGroup": {
			authorization.UpdateValidIngredientGroupsPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveValidIngredientGroup": {
			authorization.ArchiveValidIngredientGroupsPermission,
		},
		"/mealplanning.MealPlanningService/GetRandomValidIngredientGroup": {
			authorization.ReadValidIngredientGroupsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidIngredientState": {
			authorization.CreateValidIngredientStatesPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientState": {
			authorization.ReadValidIngredientStatesPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientStates": {
			authorization.ReadValidIngredientStatesPermission,
		},
		"/mealplanning.MealPlanningService/SearchForValidIngredientStates": {
			authorization.ReadValidIngredientStatesPermission,
		},
		"/mealplanning.MealPlanningService/UpdateValidIngredientState": {
			authorization.UpdateValidIngredientStatesPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveValidIngredientState": {
			authorization.ArchiveValidIngredientStatesPermission,
		},
		"/mealplanning.MealPlanningService/GetRandomValidIngredientState": {
			authorization.ReadValidIngredientStatesPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidIngredientStateIngredient": {
			authorization.CreateValidIngredientStateIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientStateIngredient": {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientStateIngredients": {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/SearchForValidIngredientStateIngredients": {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/UpdateValidIngredientStateIngredient": {
			authorization.UpdateValidIngredientStateIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveValidIngredientStateIngredient": {
			authorization.ArchiveValidIngredientStateIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientStateIngredientsByIngredient": {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientStateIngredientsByIngredientState": {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidPreparation": {
			authorization.CreateValidPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparation": {
			authorization.ReadValidPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparations": {
			authorization.ReadValidPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/SearchForValidPreparations": {
			authorization.ReadValidPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/UpdateValidPreparation": {
			authorization.UpdateValidPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveValidPreparation": {
			authorization.ArchiveValidPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/GetRandomValidPreparation": {
			authorization.ReadValidPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidMeasurementUnit": {
			authorization.CreateValidMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidMeasurementUnit": {
			authorization.ReadValidMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidMeasurementUnits": {
			authorization.ReadValidMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/SearchForValidMeasurementUnits": {
			authorization.ReadValidMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/UpdateValidMeasurementUnit": {
			authorization.UpdateValidMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveValidMeasurementUnit": {
			authorization.ArchiveValidMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/GetRandomValidMeasurementUnit": {
			authorization.ReadValidMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidVessel": {
			authorization.CreateValidVesselsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidVessel": {
			authorization.ReadValidVesselsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidVessels": {
			authorization.ReadValidVesselsPermission,
		},
		"/mealplanning.MealPlanningService/SearchForValidVessels": {
			authorization.ReadValidVesselsPermission,
		},
		"/mealplanning.MealPlanningService/UpdateValidVessel": {
			authorization.UpdateValidVesselsPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveValidVessel": {
			authorization.ArchiveValidVesselsPermission,
		},
		"/mealplanning.MealPlanningService/GetRandomValidVessel": {
			authorization.ReadValidVesselsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidInstrument": {
			authorization.CreateValidInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidInstrument": {
			authorization.ReadValidInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidInstruments": {
			authorization.ReadValidInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/SearchForValidInstruments": {
			authorization.ReadValidInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/UpdateValidInstrument": {
			authorization.UpdateValidInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveValidInstrument": {
			authorization.ArchiveValidInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/GetRandomValidInstrument": {
			authorization.ReadValidInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparationVessel": {
			authorization.ReadValidPreparationVesselsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidPreparationVessel": {
			authorization.CreateValidPreparationVesselsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparationVessels": {
			authorization.ReadValidPreparationVesselsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparationVesselsByVessel": {
			authorization.ReadValidPreparationVesselsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparationVesselsByPreparation": {
			authorization.ReadValidPreparationVesselsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientPreparation": {
			authorization.ReadValidIngredientPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidIngredientPreparation": {
			authorization.CreateValidIngredientPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientPreparations": {
			authorization.ReadValidIngredientPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientPreparationsByPreparation": {
			authorization.ReadValidIngredientPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientPreparationsByIngredient": {
			authorization.ReadValidIngredientPreparationsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientMeasurementUnit": {
			authorization.ReadValidIngredientMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidIngredientMeasurementUnit": {
			authorization.CreateValidIngredientMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientMeasurementUnits": {
			authorization.ReadValidIngredientMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientMeasurementUnitsByMeasurementUnit": {
			authorization.ReadValidIngredientMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidIngredientMeasurementUnitsByIngredient": {
			authorization.ReadValidIngredientMeasurementUnitsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparationInstrument": {
			authorization.ReadValidPreparationInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidPreparationInstrument": {
			authorization.CreateValidPreparationInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparationInstruments": {
			authorization.ReadValidPreparationInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparationInstrumentsByInstrument": {
			authorization.ReadValidPreparationInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidPreparationInstrumentsByPreparation": {
			authorization.ReadValidPreparationInstrumentsPermission,
		},
		"/mealplanning.MealPlanningService/CreateValidMeasurementUnitConversion": {
			authorization.CreateValidMeasurementUnitConversionsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidMeasurementUnitConversion": {
			authorization.ReadValidMeasurementUnitConversionsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidMeasurementUnitConversionsFromUnit": {
			authorization.ReadValidMeasurementUnitConversionsPermission,
		},
		"/mealplanning.MealPlanningService/GetValidMeasurementUnitConversionsToUnit": {
			authorization.ReadValidMeasurementUnitConversionsPermission,
		},
		"/mealplanning.MealPlanningService/CreateUserIngredientPreference": {
			authorization.CreateUserIngredientPreferencesPermission,
		},
		"/mealplanning.MealPlanningService/GetUserIngredientPreference": {
			authorization.ReadUserIngredientPreferencesPermission,
		},
		"/mealplanning.MealPlanningService/GetUserIngredientPreferences": {
			authorization.ReadUserIngredientPreferencesPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveUserIngredientPreference": {
			authorization.ReadUserIngredientPreferencesPermission,
		},
		"/identity.IdentityService/AdminUpdateUserStatus": {
			authorization.UpdateUserStatusPermission,
		},
		"/webhooks.WebhooksService/GetWebhook": {
			authorization.ReadWebhooksPermission,
		},
		"/webhooks.WebhooksService/GetWebhooks": {
			authorization.ReadWebhooksPermission,
		},
		"/webhooks.WebhooksService/CreateWebhook": {
			authorization.CreateWebhooksPermission,
		},
		"/webhooks.WebhooksService/ArchiveWebhook": {
			authorization.ArchiveWebhooksPermission,
		},
		"/webhooks.WebhooksService/AddWebhookTriggerEvent": {
			authorization.CreateWebhookTriggerEventsPermission,
		},
		"/webhooks.WebhooksService/ArchiveWebhookTriggerEvent": {
			authorization.ArchiveWebhookTriggerEventsPermission,
		},
		"/identity.IdentityService/UpdateAccount": {
			authorization.UpdateAccountPermission,
		},
		"/identity.IdentityService/ArchiveAccount": {
			authorization.ArchiveAccountPermission,
		},
		"/identity.IdentityService/CreateAccountInvitation": {
			authorization.InviteUserToAccountPermission,
		},
		"/identity.IdentityService/CancelAccountInvitation": {
			authorization.InviteUserToAccountPermission,
		},
		"/identity.IdentityService/TransferAccountOwnership": {
			authorization.TransferAccountPermission,
		},
		"/identity.IdentityService/UpdateAccountMemberPermissions": {
			authorization.ModifyMemberPermissionsForAccountPermission,
		},
		"/identity.IdentityService/ArchiveUserMembership": {
			authorization.RemoveMemberAccountPermission,
		},
		"/identity.IdentityService/GetUser": {
			authorization.ReadUserPermission,
		},
		"/identity.IdentityService/SearchForUsers": {
			authorization.ReadUserPermission,
		},
		"/identity.IdentityService/ArchiveUser": {
			authorization.ArchiveUserPermission,
		},
		"/auth.AuthService/CheckPermissions":                      noPerms,
		"/identity.IdentityService/RejectAccountInvitation":       noPerms,
		"/identity.IdentityService/AcceptAccountInvitation":       noPerms,
		"/identity.IdentityService/GetReceivedAccountInvitations": noPerms,
		"/identity.IdentityService/GetSentAccountInvitations":     noPerms,
		"/identity.IdentityService/SetDefaultAccount":             noPerms,
		"/identity.IdentityService/CreateAccount":                 noPerms,
		"/identity.IdentityService/GetAccount":                    noPerms,
		"/identity.IdentityService/GetAccounts":                   noPerms,
		"/auth.AuthService/GetAuthStatus":                         noPerms,
		"/auth.AuthService/GetActiveAccount":                      noPerms,
		"/auth.AuthService/UpdatePassword":                        noPerms,
		"/auth.AuthService/RefreshTOTPSecret":                     noPerms,
		"/auth.AuthService/VerifyTOTPSecret":                      noPerms,
		"/auth.AuthService/RequestPasswordResetToken":             noPerms,
		"/auth.AuthService/RedeemPasswordResetToken":              noPerms,
	}
)
