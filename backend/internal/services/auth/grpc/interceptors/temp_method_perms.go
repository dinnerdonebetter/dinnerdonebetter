package interceptors

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
)

var (
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
		"/auth.AuthService/GetAuthStatus":             {},
		"/auth.AuthService/UpdatePassword":            {},
		"/auth.AuthService/RefreshTOTPSecret":         {},
		"/auth.AuthService/VerifyTOTPSecret":          {},
		"/auth.AuthService/RequestPasswordResetToken": {},
		"/auth.AuthService/RedeemPasswordResetToken":  {},
	}
)
