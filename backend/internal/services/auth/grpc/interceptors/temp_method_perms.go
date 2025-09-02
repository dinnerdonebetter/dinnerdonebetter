package interceptors

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
)

var (
	noPerms = []authorization.Permission{}

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
		"/mealplanning.MealPlanningService/GetRandomValidIngredient": {
			authorization.ReadValidIngredientsPermission,
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
