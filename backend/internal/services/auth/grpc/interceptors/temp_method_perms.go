package interceptors

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/authorization"
)

var (
	noPerms = []authorization.Permission{}

	// TODO: ensure this map doesn't end up with configs for methods that don't exist.

	methodPermissions = map[string][]authorization.Permission{
		mealPlanningPerm("CreateValidIngredient"): {
			authorization.CreateValidIngredientsPermission,
		},
		mealPlanningPerm("GetValidIngredient"): {
			authorization.ReadValidIngredientsPermission,
		},
		mealPlanningPerm("GetValidIngredients"): {
			authorization.ReadValidIngredientsPermission,
		},
		mealPlanningPerm("SearchForValidIngredients"): {
			authorization.ReadValidIngredientsPermission,
		},
		mealPlanningPerm("UpdateValidIngredient"): {
			authorization.UpdateValidIngredientsPermission,
		},
		mealPlanningPerm("ArchiveValidIngredient"): {
			authorization.ArchiveValidIngredientsPermission,
		},
		settingsServicePerm("CreateServiceSetting"): {
			authorization.CreateServiceSettingsPermission,
		},
		settingsServicePerm("GetServiceSetting"): {
			authorization.ReadServiceSettingsPermission,
		},
		settingsServicePerm("GetServiceSettings"): {
			authorization.ReadServiceSettingsPermission,
		},
		settingsServicePerm("SearchForServiceSettings"): {
			authorization.ReadServiceSettingsPermission,
		},
		settingsServicePerm("ArchiveServiceSetting"): {
			authorization.ArchiveServiceSettingsPermission,
		},
		settingsServicePerm("CreateServiceSettingConfiguration"): {
			authorization.CreateServiceSettingConfigurationsPermission,
		},
		settingsServicePerm("GetServiceSettingConfigurationByName"): {
			authorization.ReadServiceSettingConfigurationsPermission,
		},
		settingsServicePerm("GetServiceSettingConfigurationsForAccount"): {
			authorization.ReadServiceSettingConfigurationsPermission,
		},
		settingsServicePerm("GetServiceSettingConfigurationsForUser"): {
			authorization.ReadServiceSettingConfigurationsPermission,
		},
		settingsServicePerm("ArchiveServiceSettingConfiguration"): {
			authorization.ArchiveServiceSettingConfigurationsPermission,
		},
		oauthServicePerm("CreateOAuth2Client"): {
			authorization.CreateOAuth2ClientsPermission,
		},
		oauthServicePerm("GetOAuth2Client"): {
			authorization.ReadOAuth2ClientsPermission,
		},
		oauthServicePerm("GetOAuth2Clients"): {
			authorization.ReadOAuth2ClientsPermission,
		},
		oauthServicePerm("ArchiveOAuth2Client"): {
			authorization.ArchiveOAuth2ClientsPermission,
		},
		mealPlanningPerm("GetRandomValidIngredient"): {
			authorization.ReadValidIngredientsPermission,
		},
		mealPlanningPerm("CreateValidIngredientGroup"): {
			authorization.CreateValidIngredientGroupsPermission,
		},
		mealPlanningPerm("GetValidIngredientGroup"): {
			authorization.ReadValidIngredientGroupsPermission,
		},
		mealPlanningPerm("GetValidIngredientGroups"): {
			authorization.ReadValidIngredientGroupsPermission,
		},
		mealPlanningPerm("SearchForValidIngredientGroups"): {
			authorization.ReadValidIngredientGroupsPermission,
		},
		mealPlanningPerm("UpdateValidIngredientGroup"): {
			authorization.UpdateValidIngredientGroupsPermission,
		},
		mealPlanningPerm("ArchiveValidIngredientGroup"): {
			authorization.ArchiveValidIngredientGroupsPermission,
		},
		mealPlanningPerm("GetRandomValidIngredientGroup"): {
			authorization.ReadValidIngredientGroupsPermission,
		},
		mealPlanningPerm("CreateValidIngredientState"): {
			authorization.CreateValidIngredientStatesPermission,
		},
		mealPlanningPerm("GetValidIngredientState"): {
			authorization.ReadValidIngredientStatesPermission,
		},
		mealPlanningPerm("GetValidIngredientStates"): {
			authorization.ReadValidIngredientStatesPermission,
		},
		mealPlanningPerm("SearchForValidIngredientStates"): {
			authorization.ReadValidIngredientStatesPermission,
		},
		mealPlanningPerm("UpdateValidIngredientState"): {
			authorization.UpdateValidIngredientStatesPermission,
		},
		mealPlanningPerm("ArchiveValidIngredientState"): {
			authorization.ArchiveValidIngredientStatesPermission,
		},
		mealPlanningPerm("GetRandomValidIngredientState"): {
			authorization.ReadValidIngredientStatesPermission,
		},
		mealPlanningPerm("CreateValidIngredientStateIngredient"): {
			authorization.CreateValidIngredientStateIngredientsPermission,
		},
		mealPlanningPerm("GetValidIngredientStateIngredient"): {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		mealPlanningPerm("GetValidIngredientStateIngredients"): {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		mealPlanningPerm("SearchForValidIngredientStateIngredients"): {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		mealPlanningPerm("UpdateValidIngredientStateIngredient"): {
			authorization.UpdateValidIngredientStateIngredientsPermission,
		},
		mealPlanningPerm("ArchiveValidIngredientStateIngredient"): {
			authorization.ArchiveValidIngredientStateIngredientsPermission,
		},
		mealPlanningPerm("GetValidIngredientStateIngredientsByIngredient"): {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		mealPlanningPerm("GetValidIngredientStateIngredientsByIngredientState"): {
			authorization.ReadValidIngredientStateIngredientsPermission,
		},
		mealPlanningPerm("CreateValidPreparation"): {
			authorization.CreateValidPreparationsPermission,
		},
		mealPlanningPerm("GetValidPreparation"): {
			authorization.ReadValidPreparationsPermission,
		},
		mealPlanningPerm("GetValidPreparations"): {
			authorization.ReadValidPreparationsPermission,
		},
		mealPlanningPerm("SearchForValidPreparations"): {
			authorization.ReadValidPreparationsPermission,
		},
		mealPlanningPerm("UpdateValidPreparation"): {
			authorization.UpdateValidPreparationsPermission,
		},
		mealPlanningPerm("ArchiveValidPreparation"): {
			authorization.ArchiveValidPreparationsPermission,
		},
		mealPlanningPerm("GetRandomValidPreparation"): {
			authorization.ReadValidPreparationsPermission,
		},
		mealPlanningPerm("CreateValidMeasurementUnit"): {
			authorization.CreateValidMeasurementUnitsPermission,
		},
		mealPlanningPerm("GetValidMeasurementUnit"): {
			authorization.ReadValidMeasurementUnitsPermission,
		},
		mealPlanningPerm("GetValidMeasurementUnits"): {
			authorization.ReadValidMeasurementUnitsPermission,
		},
		mealPlanningPerm("SearchForValidMeasurementUnits"): {
			authorization.ReadValidMeasurementUnitsPermission,
		},
		mealPlanningPerm("UpdateValidMeasurementUnit"): {
			authorization.UpdateValidMeasurementUnitsPermission,
		},
		mealPlanningPerm("ArchiveValidMeasurementUnit"): {
			authorization.ArchiveValidMeasurementUnitsPermission,
		},
		mealPlanningPerm("GetRandomValidMeasurementUnit"): {
			authorization.ReadValidMeasurementUnitsPermission,
		},
		mealPlanningPerm("CreateValidVessel"): {
			authorization.CreateValidVesselsPermission,
		},
		mealPlanningPerm("GetValidVessel"): {
			authorization.ReadValidVesselsPermission,
		},
		mealPlanningPerm("GetValidVessels"): {
			authorization.ReadValidVesselsPermission,
		},
		mealPlanningPerm("SearchForValidVessels"): {
			authorization.ReadValidVesselsPermission,
		},
		mealPlanningPerm("UpdateValidVessel"): {
			authorization.UpdateValidVesselsPermission,
		},
		mealPlanningPerm("ArchiveValidVessel"): {
			authorization.ArchiveValidVesselsPermission,
		},
		mealPlanningPerm("GetRandomValidVessel"): {
			authorization.ReadValidVesselsPermission,
		},
		mealPlanningPerm("CreateValidInstrument"): {
			authorization.CreateValidInstrumentsPermission,
		},
		mealPlanningPerm("GetValidInstrument"): {
			authorization.ReadValidInstrumentsPermission,
		},
		mealPlanningPerm("GetValidInstruments"): {
			authorization.ReadValidInstrumentsPermission,
		},
		mealPlanningPerm("SearchForValidInstruments"): {
			authorization.ReadValidInstrumentsPermission,
		},
		mealPlanningPerm("UpdateValidInstrument"): {
			authorization.UpdateValidInstrumentsPermission,
		},
		mealPlanningPerm("ArchiveValidInstrument"): {
			authorization.ArchiveValidInstrumentsPermission,
		},
		mealPlanningPerm("GetRandomValidInstrument"): {
			authorization.ReadValidInstrumentsPermission,
		},
		mealPlanningPerm("GetValidPreparationVessel"): {
			authorization.ReadValidPreparationVesselsPermission,
		},
		mealPlanningPerm("CreateValidPreparationVessel"): {
			authorization.CreateValidPreparationVesselsPermission,
		},
		mealPlanningPerm("GetValidPreparationVessels"): {
			authorization.ReadValidPreparationVesselsPermission,
		},
		mealPlanningPerm("GetValidPreparationVesselsByVessel"): {
			authorization.ReadValidPreparationVesselsPermission,
		},
		mealPlanningPerm("GetValidPreparationVesselsByPreparation"): {
			authorization.ReadValidPreparationVesselsPermission,
		},
		mealPlanningPerm("GetValidIngredientPreparation"): {
			authorization.ReadValidIngredientPreparationsPermission,
		},
		mealPlanningPerm("CreateValidIngredientPreparation"): {
			authorization.CreateValidIngredientPreparationsPermission,
		},
		mealPlanningPerm("GetValidIngredientPreparations"): {
			authorization.ReadValidIngredientPreparationsPermission,
		},
		mealPlanningPerm("GetValidIngredientPreparationsByPreparation"): {
			authorization.ReadValidIngredientPreparationsPermission,
		},
		mealPlanningPerm("GetValidIngredientPreparationsByIngredient"): {
			authorization.ReadValidIngredientPreparationsPermission,
		},
		mealPlanningPerm("GetValidIngredientMeasurementUnit"): {
			authorization.ReadValidIngredientMeasurementUnitsPermission,
		},
		mealPlanningPerm("CreateValidIngredientMeasurementUnit"): {
			authorization.CreateValidIngredientMeasurementUnitsPermission,
		},
		mealPlanningPerm("GetValidIngredientMeasurementUnits"): {
			authorization.ReadValidIngredientMeasurementUnitsPermission,
		},
		mealPlanningPerm("GetValidIngredientMeasurementUnitsByMeasurementUnit"): {
			authorization.ReadValidIngredientMeasurementUnitsPermission,
		},
		mealPlanningPerm("GetValidIngredientMeasurementUnitsByIngredient"): {
			authorization.ReadValidIngredientMeasurementUnitsPermission,
		},
		mealPlanningPerm("GetValidPreparationInstrument"): {
			authorization.ReadValidPreparationInstrumentsPermission,
		},
		mealPlanningPerm("CreateValidPreparationInstrument"): {
			authorization.CreateValidPreparationInstrumentsPermission,
		},
		mealPlanningPerm("GetValidPreparationInstruments"): {
			authorization.ReadValidPreparationInstrumentsPermission,
		},
		mealPlanningPerm("GetValidPreparationInstrumentsByInstrument"): {
			authorization.ReadValidPreparationInstrumentsPermission,
		},
		mealPlanningPerm("GetValidPreparationInstrumentsByPreparation"): {
			authorization.ReadValidPreparationInstrumentsPermission,
		},
		mealPlanningPerm("CreateValidMeasurementUnitConversion"): {
			authorization.CreateValidMeasurementUnitConversionsPermission,
		},
		mealPlanningPerm("GetValidMeasurementUnitConversion"): {
			authorization.ReadValidMeasurementUnitConversionsPermission,
		},
		mealPlanningPerm("GetValidMeasurementUnitConversionsFromUnit"): {
			authorization.ReadValidMeasurementUnitConversionsPermission,
		},
		mealPlanningPerm("GetValidMeasurementUnitConversionsToUnit"): {
			authorization.ReadValidMeasurementUnitConversionsPermission,
		},
		mealPlanningPerm("CreateUserIngredientPreference"): {
			authorization.CreateUserIngredientPreferencesPermission,
		},
		mealPlanningPerm("GetUserIngredientPreference"): {
			authorization.ReadUserIngredientPreferencesPermission,
		},
		mealPlanningPerm("GetUserIngredientPreferences"): {
			authorization.ReadUserIngredientPreferencesPermission,
		},
		mealPlanningPerm("ArchiveUserIngredientPreference"): {
			authorization.ReadUserIngredientPreferencesPermission,
		},
		mealPlanningPerm("CreateAccountInstrumentOwnership"): {
			authorization.CreateAccountInstrumentOwnershipsPermission,
		},
		mealPlanningPerm("GetAccountInstrumentOwnership"): {
			authorization.ReadAccountInstrumentOwnershipsPermission,
		},
		mealPlanningPerm("GetAccountInstrumentOwnerships"): {
			authorization.ReadAccountInstrumentOwnershipsPermission,
		},
		mealPlanningPerm("ArchiveAccountInstrumentOwnership"): {
			authorization.ReadAccountInstrumentOwnershipsPermission,
		},
		mealPlanningPerm("CreateRecipe"): {
			authorization.CreateRecipesPermission,
		},
		mealPlanningPerm("GetRecipe"): {
			authorization.ReadRecipesPermission,
		},
		mealPlanningPerm("UpdateRecipe"): {
			authorization.UpdateRecipesPermission,
		},
		mealPlanningPerm("ArchiveRecipe"): {
			authorization.ArchiveRecipesPermission,
		},
		mealPlanningPerm("SearchForRecipes"): {
			authorization.ReadRecipesPermission,
		},
		mealPlanningPerm("CreateMeal"): {
			authorization.CreateMealsPermission,
		},
		mealPlanningPerm("GetMeal"): {
			authorization.ReadMealsPermission,
		},
		mealPlanningPerm("UpdateMeal"): {
			authorization.UpdateMealsPermission,
		},
		mealPlanningPerm("ArchiveMeal"): {
			authorization.ArchiveMealsPermission,
		},
		mealPlanningPerm("SearchForMeals"): {
			authorization.ReadMealsPermission,
		},
		mealPlanningPerm("CreateMealPlan"): {
			authorization.CreateMealPlansPermission,
		},
		mealPlanningPerm("GetMealPlan"): {
			authorization.ReadMealPlansPermission,
		},
		mealPlanningPerm("SearchForMealPlans"): {
			authorization.ReadMealPlansPermission,
		},
		mealPlanningPerm("UpdateMealPlan"): {
			authorization.UpdateMealPlansPermission,
		},
		mealPlanningPerm("ArchiveMealPlan"): {
			authorization.ArchiveMealPlansPermission,
		},
		mealPlanningPerm("CreateMealPlanOption"): {
			authorization.CreateMealPlanOptionsPermission,
		},
		mealPlanningPerm("GetMealPlanOption"): {
			authorization.ReadMealPlanOptionsPermission,
		},
		mealPlanningPerm("GetMealPlanOptions"): {
			authorization.ReadMealPlanOptionsPermission,
		},
		mealPlanningPerm("UpdateMealPlanOption"): {
			authorization.UpdateMealPlanOptionsPermission,
		},
		mealPlanningPerm("ArchiveMealPlanOption"): {
			authorization.ArchiveMealPlanOptionsPermission,
		},
		mealPlanningPerm("CreateMealPlanEvent"): {
			authorization.CreateMealPlanEventsPermission,
		},
		mealPlanningPerm("GetMealPlanEvent"): {
			authorization.ReadMealPlanEventsPermission,
		},
		mealPlanningPerm("GetMealPlanEvents"): {
			authorization.ReadMealPlanEventsPermission,
		},
		mealPlanningPerm("UpdateMealPlanEvent"): {
			authorization.UpdateMealPlanEventsPermission,
		},
		mealPlanningPerm("ArchiveMealPlanEvent"): {
			authorization.ArchiveMealPlanEventsPermission,
		},
		mealPlanningPerm("CreateMealPlanTask"): {
			authorization.CreateMealPlanTasksPermission,
		},
		mealPlanningPerm("GetMealPlanTask"): {
			authorization.ReadMealPlanTasksPermission,
		},
		mealPlanningPerm("GetMealPlanTasks"): {
			authorization.ReadMealPlanTasksPermission,
		},
		mealPlanningPerm("UpdateMealPlanTask"): {
			authorization.UpdateMealPlanTasksPermission,
		},
		mealPlanningPerm("CreateMealPlanEvent"): {
			authorization.CreateMealPlanEventsPermission,
		},
		mealPlanningPerm("GetMealPlanEvent"): {
			authorization.ReadMealPlanEventsPermission,
		},
		mealPlanningPerm("GetMealPlanEvents"): {
			authorization.ReadMealPlanEventsPermission,
		},
		mealPlanningPerm("UpdateMealPlanEvent"): {
			authorization.UpdateMealPlanEventsPermission,
		},
		mealPlanningPerm("ArchiveMealPlanEvent"): {
			authorization.ArchiveMealPlanEventsPermission,
		},
		mealPlanningPerm("CreateMealPlanGroceryListItem"): {
			authorization.CreateMealPlanGroceryListItemsPermission,
		},
		mealPlanningPerm("GetMealPlanGroceryListItem"): {
			authorization.ReadMealPlanGroceryListItemsPermission,
		},
		mealPlanningPerm("GetMealPlanGroceryListItemsForMealPlan"): {
			authorization.ReadMealPlanGroceryListItemsPermission,
		},
		mealPlanningPerm("GetMealPlanGroceryListItems"): {
			authorization.ReadMealPlanGroceryListItemsPermission,
		},
		mealPlanningPerm("UpdateMealPlanGroceryListItem"): {
			authorization.UpdateMealPlanGroceryListItemsPermission,
		},
		mealPlanningPerm("ArchiveMealPlanGroceryListItem"): {
			authorization.ArchiveMealPlanGroceryListItemsPermission,
		},
		mealPlanningPerm("CloneRecipe"): {
			authorization.ReadRecipesPermission, // TODO: this should be its own perm
		},
		mealPlanningPerm("SearchForMeals"): {
			authorization.ReadMealsPermission,
		},
		mealPlanningPerm("GetMeals"): {
			authorization.ReadMealsPermission,
		},
		mealPlanningPerm("GetMealPlansForAccount"): {
			authorization.ReadMealPlansPermission,
		},
		mealPlanningPerm("GetMeal"): {
			authorization.ReadMealsPermission,
		},
		mealPlanningPerm("GetMealPlanTasks"): {
			authorization.ReadMealPlanTasksPermission,
		},
		mealPlanningPerm("CreateRecipeStep"): {
			authorization.CreateRecipeStepsPermission,
		},
		mealPlanningPerm("GetRecipeSteps"): {
			authorization.ReadRecipeStepsPermission,
		},
		mealPlanningPerm("GetRecipeStep"): {
			authorization.ReadRecipeStepsPermission,
		},
		mealPlanningPerm("UpdateRecipeStep"): {
			authorization.UpdateRecipeStepsPermission,
		},
		mealPlanningPerm("ArchiveRecipeStep"): {
			authorization.ArchiveRecipeStepsPermission,
		},
		mealPlanningPerm("CreateRecipeStepVessel"): {
			authorization.CreateRecipeStepVesselsPermission,
		},
		mealPlanningPerm("GetRecipeStepVessels"): {
			authorization.ReadRecipeStepVesselsPermission,
		},
		mealPlanningPerm("GetRecipeStepVessel"): {
			authorization.ReadRecipeStepVesselsPermission,
		},
		mealPlanningPerm("UpdateRecipeStepVessel"): {
			authorization.UpdateRecipeStepVesselsPermission,
		},
		mealPlanningPerm("ArchiveRecipeStepVessel"): {
			authorization.ArchiveRecipeStepVesselsPermission,
		},
		mealPlanningPerm("CreateRecipeStepProduct"): {
			authorization.CreateRecipeStepProductsPermission,
		},
		mealPlanningPerm("GetRecipeStepProducts"): {
			authorization.ReadRecipeStepProductsPermission,
		},
		mealPlanningPerm("GetRecipeStepProduct"): {
			authorization.ReadRecipeStepProductsPermission,
		},
		mealPlanningPerm("UpdateRecipeStepProduct"): {
			authorization.UpdateRecipeStepProductsPermission,
		},
		mealPlanningPerm("ArchiveRecipeStepProduct"): {
			authorization.ArchiveRecipeStepProductsPermission,
		},
		mealPlanningPerm("CreateRecipePrepTask"): {
			authorization.CreateRecipePrepTasksPermission,
		},
		mealPlanningPerm("GetRecipePrepTasks"): {
			authorization.ReadRecipePrepTasksPermission,
		},
		mealPlanningPerm("GetRecipePrepTask"): {
			authorization.ReadRecipePrepTasksPermission,
		},
		mealPlanningPerm("UpdateRecipePrepTask"): {
			authorization.UpdateRecipePrepTasksPermission,
		},
		mealPlanningPerm("ArchiveRecipePrepTask"): {
			authorization.ArchiveRecipePrepTasksPermission,
		},
		mealPlanningPerm("CreateRecipeStepInstrument"): {
			authorization.CreateRecipeStepInstrumentsPermission,
		},
		mealPlanningPerm("GetRecipeStepInstruments"): {
			authorization.ReadRecipeStepInstrumentsPermission,
		},
		mealPlanningPerm("GetRecipeStepInstrument"): {
			authorization.ReadRecipeStepInstrumentsPermission,
		},
		mealPlanningPerm("UpdateRecipeStepInstrument"): {
			authorization.UpdateRecipeStepInstrumentsPermission,
		},
		mealPlanningPerm("ArchiveRecipeStepInstrument"): {
			authorization.ArchiveRecipeStepInstrumentsPermission,
		},
		mealPlanningPerm("CreateRecipeStepIngredient"): {
			authorization.CreateRecipeStepIngredientsPermission,
		},
		mealPlanningPerm("GetRecipeStepIngredients"): {
			authorization.ReadRecipeStepIngredientsPermission,
		},
		mealPlanningPerm("GetRecipeStepIngredient"): {
			authorization.ReadRecipeStepIngredientsPermission,
		},
		mealPlanningPerm("UpdateRecipeStepIngredient"): {
			authorization.UpdateRecipeStepIngredientsPermission,
		},
		mealPlanningPerm("ArchiveRecipeStepIngredient"): {
			authorization.ArchiveRecipeStepIngredientsPermission,
		},
		mealPlanningPerm("CreateRecipeStepCompletionCondition"): {
			authorization.CreateRecipeStepCompletionConditionsPermission,
		},
		mealPlanningPerm("GetRecipeStepCompletionConditions"): {
			authorization.ReadRecipeStepCompletionConditionsPermission,
		},
		mealPlanningPerm("GetRecipeStepCompletionCondition"): {
			authorization.ReadRecipeStepCompletionConditionsPermission,
		},
		mealPlanningPerm("UpdateRecipeStepCompletionCondition"): {
			authorization.UpdateRecipeStepCompletionConditionsPermission,
		},
		mealPlanningPerm("ArchiveRecipeStepCompletionCondition"): {
			authorization.ArchiveRecipeStepCompletionConditionsPermission,
		},
		mealPlanningPerm("CreateRecipeRating"): {
			authorization.CreateRecipeRatingsPermission,
		},
		mealPlanningPerm("GetRecipeRatings"): {
			authorization.ReadRecipeRatingsPermission,
		},
		mealPlanningPerm("GetRecipeRating"): {
			authorization.ReadRecipeRatingsPermission,
		},
		mealPlanningPerm("UpdateRecipeRating"): {
			authorization.UpdateRecipeRatingsPermission,
		},
		mealPlanningPerm("ArchiveRecipeRating"): {
			authorization.ArchiveRecipeRatingsPermission,
		},
		mealPlanningPerm("GetRecipeRatingsForRecipe"): {
			authorization.ReadRecipeRatingsPermission,
		},
		mealPlanningPerm("CreateMealPlanOptionVote"): {
			authorization.CreateMealPlanOptionVotesPermission,
		},
		mealPlanningPerm("UpdateMealPlanOptionVote"): {
			authorization.UpdateMealPlanOptionVotesPermission,
		},
		mealPlanningPerm("GetMealPlanOptionVote"): {
			authorization.ReadMealPlanOptionVotesPermission,
		},
		mealPlanningPerm("GetMealPlanOptionVotes"): {
			authorization.ReadMealPlanOptionVotesPermission,
		},
		mealPlanningPerm("ArchiveMealPlanOptionVote"): {
			authorization.ArchiveMealPlanOptionVotesPermission,
		},
		mealPlanningPerm("RunFinalizeMealPlanWorker"): {
			authorization.UpdateMealPlansPermission, // TODO: this should be its own perm
		},
		userNotifsServicePerm("CreateUserNotification"): {
			authorization.CreateUserNotificationsPermission,
		},
		userNotifsServicePerm("GetUserNotification"): {
			authorization.ReadUserNotificationsPermission,
		},
		userNotifsServicePerm("GetUserNotifications"): {
			authorization.ReadUserNotificationsPermission,
		},
		userNotifsServicePerm("UpdateUserNotification"): {
			authorization.UpdateUserNotificationsPermission,
		},
		identityServicePerm("AdminUpdateUserStatus"): {
			authorization.UpdateUserStatusPermission,
		},
		webhooksServicePerm("GetWebhook"): {
			authorization.ReadWebhooksPermission,
		},
		webhooksServicePerm("GetWebhooks"): {
			authorization.ReadWebhooksPermission,
		},
		webhooksServicePerm("CreateWebhook"): {
			authorization.CreateWebhooksPermission,
		},
		webhooksServicePerm("ArchiveWebhook"): {
			authorization.ArchiveWebhooksPermission,
		},
		webhooksServicePerm("AddWebhookTriggerEvent"): {
			authorization.CreateWebhookTriggerEventsPermission,
		},
		webhooksServicePerm("ArchiveWebhookTriggerEvent"): {
			authorization.ArchiveWebhookTriggerEventsPermission,
		},
		identityServicePerm("UpdateAccount"): {
			authorization.UpdateAccountPermission,
		},
		identityServicePerm("ArchiveAccount"): {
			authorization.ArchiveAccountPermission,
		},
		identityServicePerm("CreateAccountInvitation"): {
			authorization.InviteUserToAccountPermission,
		},
		identityServicePerm("CancelAccountInvitation"): {
			authorization.InviteUserToAccountPermission,
		},
		identityServicePerm("TransferAccountOwnership"): {
			authorization.TransferAccountPermission,
		},
		identityServicePerm("UpdateAccountMemberPermissions"): {
			authorization.ModifyMemberPermissionsForAccountPermission,
		},
		identityServicePerm("ArchiveUserMembership"): {
			authorization.RemoveMemberAccountPermission,
		},
		identityServicePerm("GetUser"): {
			authorization.ReadUserPermission,
		},
		identityServicePerm("SearchForUsers"): {
			authorization.ReadUserPermission,
		},
		identityServicePerm("ArchiveUser"): {
			authorization.ArchiveUserPermission,
		},
		identityServicePerm("RejectAccountInvitation"):       noPerms,
		identityServicePerm("AcceptAccountInvitation"):       noPerms,
		identityServicePerm("GetReceivedAccountInvitations"): noPerms,
		identityServicePerm("GetSentAccountInvitations"):     noPerms,
		identityServicePerm("SetDefaultAccount"):             noPerms,
		identityServicePerm("CreateAccount"):                 noPerms,
		identityServicePerm("GetAccount"):                    noPerms,
		identityServicePerm("GetAccounts"):                   noPerms,
		authServicePerm("CheckPermissions"):                  noPerms,
		authServicePerm("GetAuthStatus"):                     noPerms,
		authServicePerm("GetActiveAccount"):                  noPerms,
		authServicePerm("UpdatePassword"):                    noPerms,
		authServicePerm("RefreshTOTPSecret"):                 noPerms,
		authServicePerm("VerifyTOTPSecret"):                  noPerms,
		authServicePerm("RequestPasswordResetToken"):         noPerms,
		authPerm("RedeemPasswordResetToken"):                 noPerms,
	}
)

func permString(collectionName, serviceName, methodName string) string {
	return fmt.Sprintf("/%s.%s/%s", collectionName, serviceName, methodName)
}

func oauthServicePerm(method string) string {
	return permString("oauth", "OAuthService", method)
}

func settingsServicePerm(method string) string {
	return permString("settings", "SettingsService", method)
}

func userNotifsServicePerm(method string) string {
	return permString("notifications", "UserNotificationsService", method)
}

func webhooksServicePerm(method string) string {
	return permString("webhooks", "WebhooksService", method)
}

func identityServicePerm(method string) string {
	return permString("identity", "IdentityService", method)
}

func authServicePerm(method string) string {
	return permString("auth", "AuthService", method)
}

func mealPlanningPerm(method string) string {
	return permString("mealplanning", "MealPlanningService", method)
}

func authPerm(method string) string {
	return permString("auth", "AuthService", method)
}
