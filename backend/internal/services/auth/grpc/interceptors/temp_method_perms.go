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
		"/identity.IdentityService/AdminUpdateUserStatus": {
			authorization.UpdateUserStatusPermission,
		},
		"/webhooks.WebhooksService/GetWebhook": {
			authorization.ReadWebhooksPermission,
		},
		"/webhooks.WebhooksService/CreateWebhook": {
			authorization.CreateWebhooksPermission,
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
