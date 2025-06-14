package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func BuildFakeUserDataCollectionResponse() *types.UserDataCollectionResponse {
	return &types.UserDataCollectionResponse{
		ReportID: BuildFakeID(),
	}
}

func BuildFakeUserDataCollection() *types.UserDataCollection {
	user := BuildFakeUser()
	account := BuildFakeAccount()

	recipeRatings := BuildFakeRecipeRatingsList().Data
	recipes := BuildFakeRecipesList().Data
	meals := BuildFakeMealsList().Data
	receivedAccountInvitations := BuildFakeAccountInvitationsList().Data
	userIngredientPreferences := BuildFakeUserIngredientPreferencesList().Data
	sentAccountInvitations := BuildFakeAccountInvitationsList().Data
	serviceSettingConfigurations := BuildFakeServiceSettingConfigurationsList().Data
	auditLogEntries := BuildFakeAuditLogEntriesList().Data
	mealPlans := BuildFakeMealPlansList().Data
	webhooks := BuildFakeWebhooksList().Data
	accountInstrumentOwnerships := BuildFakeAccountInstrumentOwnershipsList().Data

	return &types.UserDataCollection{
		ReportID: BuildFakeID(),
		User:     *user,
		Core: types.CoreUserDataCollection{
			ReceivedInvites:                  pointer.DereferenceSlice(receivedAccountInvitations),
			SentInvites:                      pointer.DereferenceSlice(sentAccountInvitations),
			UserServiceSettingConfigurations: pointer.DereferenceSlice(serviceSettingConfigurations),
			UserAuditLogEntries:              pointer.DereferenceSlice(auditLogEntries),
			Accounts:                         []types.Account{*account},
			Webhooks: map[string][]types.Webhook{
				account.ID: pointer.DereferenceSlice(webhooks),
			},
			ServiceSettingConfigurations: map[string][]types.ServiceSettingConfiguration{
				account.ID: pointer.DereferenceSlice(serviceSettingConfigurations),
			},
			AuditLogEntries: map[string][]types.AuditLogEntry{
				account.ID: pointer.DereferenceSlice(auditLogEntries),
			},
		},
		Eating: types.EatingUserDataCollection{
			RecipeRatings:             pointer.DereferenceSlice(recipeRatings),
			Recipes:                   pointer.DereferenceSlice(recipes),
			Meals:                     pointer.DereferenceSlice(meals),
			UserIngredientPreferences: pointer.DereferenceSlice(userIngredientPreferences),
			AccountInstrumentOwnerships: map[string][]types.AccountInstrumentOwnership{
				account.ID: pointer.DereferenceSlice(accountInstrumentOwnerships),
			},
			MealPlans: map[string][]types.MealPlan{
				account.ID: pointer.DereferenceSlice(mealPlans),
			},
		},
	}
}
