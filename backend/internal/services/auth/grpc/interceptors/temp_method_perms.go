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
		"/mealplanning.MealPlanningService/UpdateValidIngredient": {
			authorization.UpdateValidIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/ArchiveValidIngredient": {
			authorization.ArchiveValidIngredientsPermission,
		},
		"/mealplanning.MealPlanningService/GetRandomValidIngredient": {
			authorization.ReadValidIngredientsPermission,
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
