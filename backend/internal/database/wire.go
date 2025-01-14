package database

import (
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/google/wire"
)

var (
	// DBProviders represents what we provide to dependency injectors.
	DBProviders = wire.NewSet(
		ProvideMealPlanTaskDataManager,
		ProvideHouseholdDataManager,
		ProvideHouseholdInvitationDataManager,
		ProvideHouseholdUserMembershipDataManager,
		ProvideMealDataManager,
		ProvideMealPlanDataManager,
		ProvideMealPlanEventDataManager,
		ProvideMealPlanOptionDataManager,
		ProvideMealPlanOptionVoteDataManager,
		ProvideUserDataManager,
		ProvideAdminUserDataManager,
		ProvidePasswordResetTokenDataManager,
		ProvideWebhookDataManager,
		ProvideMealPlanGroceryListItemDataManager,
		ProvideServiceSettingDataManager,
		ProvideServiceSettingConfigurationDataManager,
		ProvideUserIngredientPreferenceDataManager,
		ProvideHouseholdInstrumentOwnershipDataManager,
		ProvideOAuth2ClientDataManager,
		ProvideOAuth2ClientTokenDataManager,
		ProvideUserNotificationDataManager,
		ProvideAuditLogEntryDataManager,
		ProvideDataPrivacyDataManager,
		ProvideValidEnumerationDataManager,
		ProvideRecipeManagementDataManager,
	)
)

// ProvideMealPlanTaskDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanTaskDataManager(db DataManager) types.MealPlanTaskDataManager {
	return db
}

// ProvideHouseholdDataManager is an arbitrary function for dependency injection's sake.
func ProvideHouseholdDataManager(db DataManager) types.HouseholdDataManager {
	return db
}

// ProvideHouseholdInvitationDataManager is an arbitrary function for dependency injection's sake.
func ProvideHouseholdInvitationDataManager(db DataManager) types.HouseholdInvitationDataManager {
	return db
}

// ProvideHouseholdUserMembershipDataManager is an arbitrary function for dependency injection's sake.
func ProvideHouseholdUserMembershipDataManager(db DataManager) types.HouseholdUserMembershipDataManager {
	return db
}

// ProvideRecipeDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeDataManager(db DataManager) types.RecipeDataManager {
	return db
}

// ProvideRecipeStepDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepDataManager(db DataManager) types.RecipeStepDataManager {
	return db
}

// ProvideRecipeStepInstrumentDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepInstrumentDataManager(db DataManager) types.RecipeStepInstrumentDataManager {
	return db
}

// ProvideRecipeStepProductDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepProductDataManager(db DataManager) types.RecipeStepProductDataManager {
	return db
}

// ProvideRecipeStepIngredientDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepIngredientDataManager(db DataManager) types.RecipeStepIngredientDataManager {
	return db
}

// ProvideMealDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealDataManager(db DataManager) types.MealDataManager {
	return db
}

// ProvideMealPlanDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanDataManager(db DataManager) types.MealPlanDataManager {
	return db
}

// ProvideMealPlanEventDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanEventDataManager(db DataManager) types.MealPlanEventDataManager {
	return db
}

// ProvideMealPlanOptionDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanOptionDataManager(db DataManager) types.MealPlanOptionDataManager {
	return db
}

// ProvideMealPlanOptionVoteDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanOptionVoteDataManager(db DataManager) types.MealPlanOptionVoteDataManager {
	return db
}

// ProvideUserDataManager is an arbitrary function for dependency injection's sake.
func ProvideUserDataManager(db DataManager) types.UserDataManager {
	return db
}

// ProvideAdminUserDataManager is an arbitrary function for dependency injection's sake.
func ProvideAdminUserDataManager(db DataManager) types.AdminUserDataManager {
	return db
}

// ProvidePasswordResetTokenDataManager is an arbitrary function for dependency injection's sake.
func ProvidePasswordResetTokenDataManager(db DataManager) types.PasswordResetTokenDataManager {
	return db
}

// ProvideWebhookDataManager is an arbitrary function for dependency injection's sake.
func ProvideWebhookDataManager(db DataManager) types.WebhookDataManager {
	return db
}

// ProvideRecipePrepTaskDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipePrepTaskDataManager(db DataManager) types.RecipePrepTaskDataManager {
	return db
}

// ProvideMealPlanGroceryListItemDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanGroceryListItemDataManager(db DataManager) types.MealPlanGroceryListItemDataManager {
	return db
}

// ProvideRecipeMediaDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeMediaDataManager(db DataManager) types.RecipeMediaDataManager {
	return db
}

// ProvideRecipeStepCompletionConditionDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepCompletionConditionDataManager(db DataManager) types.RecipeStepCompletionConditionDataManager {
	return db
}

// ProvideRecipeStepVesselDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepVesselDataManager(db DataManager) types.RecipeStepVesselDataManager {
	return db
}

// ProvideServiceSettingDataManager is an arbitrary function for dependency injection's sake.
func ProvideServiceSettingDataManager(db DataManager) types.ServiceSettingDataManager {
	return db
}

// ProvideServiceSettingConfigurationDataManager is an arbitrary function for dependency injection's sake.
func ProvideServiceSettingConfigurationDataManager(db DataManager) types.ServiceSettingConfigurationDataManager {
	return db
}

// ProvideUserIngredientPreferenceDataManager is an arbitrary function for dependency injection's sake.
func ProvideUserIngredientPreferenceDataManager(db DataManager) types.UserIngredientPreferenceDataManager {
	return db
}

// ProvideRecipeRatingDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeRatingDataManager(db DataManager) types.RecipeRatingDataManager {
	return db
}

// ProvideHouseholdInstrumentOwnershipDataManager is an arbitrary function for dependency injection's sake.
func ProvideHouseholdInstrumentOwnershipDataManager(db DataManager) types.HouseholdInstrumentOwnershipDataManager {
	return db
}

// ProvideOAuth2ClientDataManager is an arbitrary function for dependency injection's sake.
func ProvideOAuth2ClientDataManager(db DataManager) types.OAuth2ClientDataManager {
	return db
}

// ProvideOAuth2ClientTokenDataManager is an arbitrary function for dependency injection's sake.
func ProvideOAuth2ClientTokenDataManager(db DataManager) types.OAuth2ClientTokenDataManager {
	return db
}

// ProvideUserNotificationDataManager is an arbitrary function for dependency injection's sake.
func ProvideUserNotificationDataManager(db DataManager) types.UserNotificationDataManager {
	return db
}

// ProvideAuditLogEntryDataManager is an arbitrary function for dependency injection's sake.
func ProvideAuditLogEntryDataManager(db DataManager) types.AuditLogEntryDataManager {
	return db
}

// ProvideDataPrivacyDataManager is an arbitrary function for dependency injection's sake.
func ProvideDataPrivacyDataManager(db DataManager) types.DataPrivacyDataManager {
	return db
}

// ProvideValidEnumerationDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidEnumerationDataManager(db DataManager) types.ValidEnumerationDataManager {
	return db
}

// ProvideRecipeManagementDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeManagementDataManager(db DataManager) types.RecipeManagementDataManager {
	return db
}
