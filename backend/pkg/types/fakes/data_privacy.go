package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func BuildFakeUserDataCollectionResponse() *types.UserDataCollectionResponse {
	return &types.UserDataCollectionResponse{
		ReportID: BuildFakeID(),
	}
}

func BuildFakeUserDataCollection() *types.UserDataCollection {
	user := BuildFakeUser()
	household := BuildFakeHousehold()

	recipeRatings := BuildFakeRecipeRatingsList().Data
	recipes := BuildFakeRecipesList().Data
	meals := BuildFakeMealsList().Data
	receivedHouseholdInvitations := BuildFakeHouseholdInvitationsList().Data
	userIngredientPreferences := BuildFakeUserIngredientPreferencesList().Data
	sentHouseholdInvitations := BuildFakeHouseholdInvitationsList().Data
	serviceSettingConfigurations := BuildFakeServiceSettingConfigurationsList().Data
	auditLogEntries := BuildFakeAuditLogEntriesList().Data
	mealPlans := BuildFakeMealPlansList().Data
	webhooks := BuildFakeWebhooksList().Data
	householdInstrumentOwnerships := BuildFakeHouseholdInstrumentOwnershipsList().Data

	return &types.UserDataCollection{
		ReportID:                         BuildFakeID(),
		User:                             *user,
		RecipeRatings:                    pointer.DereferenceSlice(recipeRatings),
		Recipes:                          pointer.DereferenceSlice(recipes),
		Meals:                            pointer.DereferenceSlice(meals),
		ReceivedInvites:                  pointer.DereferenceSlice(receivedHouseholdInvitations),
		UserIngredientPreferences:        pointer.DereferenceSlice(userIngredientPreferences),
		SentInvites:                      pointer.DereferenceSlice(sentHouseholdInvitations),
		UserServiceSettingConfigurations: pointer.DereferenceSlice(serviceSettingConfigurations),
		UserAuditLogEntries:              pointer.DereferenceSlice(auditLogEntries),
		Households:                       []types.Household{*household},
		Webhooks: map[string][]types.Webhook{
			household.ID: pointer.DereferenceSlice(webhooks),
		},
		ServiceSettingConfigurations: map[string][]types.ServiceSettingConfiguration{
			household.ID: pointer.DereferenceSlice(serviceSettingConfigurations),
		},
		HouseholdInstrumentOwnerships: map[string][]types.HouseholdInstrumentOwnership{
			household.ID: pointer.DereferenceSlice(householdInstrumentOwnerships),
		},
		AuditLogEntries: map[string][]types.AuditLogEntry{
			household.ID: pointer.DereferenceSlice(auditLogEntries),
		},
		MealPlans: map[string][]types.MealPlan{
			household.ID: pointer.DereferenceSlice(mealPlans),
		},
	}
}
