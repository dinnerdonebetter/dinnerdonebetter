package authorization

import (
	"github.com/mikespook/gorbac/v2"
)

var _ gorbac.Permission = (*Permission)(nil)

type (
	// Permission is a simple string alias.
	Permission string
)

// ID implements the gorbac Permission interface.
func (p Permission) ID() string {
	return string(p)
}

// Match implements the gorbac Permission interface.
func (p Permission) Match(perm gorbac.Permission) bool {
	return p.ID() == perm.ID()
}

var (
	// ServiceAdminPermissions is every service admin permission.
	ServiceAdminPermissions = []gorbac.Permission{
		ReadUserDataPermission,
		UpdateUserStatusPermission,
		ReadUserPermission,
		SearchUserPermission,
		ArchiveUserPermission,
		CreateOAuth2ClientsPermission,
		ArchiveOAuth2ClientsPermission,
		ArchiveServiceSettingsPermission,
		CreateUserNotificationsPermission,
		ImpersonateUserPermission,
		PublishArbitraryQueueMessagePermission,
		UpdateRecipesStatusPermission,
		// only admins can arbitrarily create these via the API, this is exclusively for integration test purposes.
		CreateServiceSettingsPermission,
		CreateMealPlanTasksPermission,
		CreateMealPlanGroceryListItemsPermission,
		CreateWaitlistsPermission,
		UpdateWaitlistsPermission,
		ArchiveWaitlistsPermission,
		CreateProductsPermission,
		ReadProductsPermission,
		UpdateProductsPermission,
		ArchiveProductsPermission,
		CreateSubscriptionsPermission,
		ReadSubscriptionsPermission,
		UpdateSubscriptionsPermission,
		ArchiveSubscriptionsPermission,
	}

	// ServiceDataAdminPermissions is every service admin permission.
	ServiceDataAdminPermissions = []gorbac.Permission{
		CreateValidInstrumentsPermission,
		UpdateValidInstrumentsPermission,
		ArchiveValidInstrumentsPermission,
		CreateValidVesselsPermission,
		UpdateValidVesselsPermission,
		ArchiveValidVesselsPermission,
		CreateValidIngredientsPermission,
		UpdateValidIngredientsPermission,
		ArchiveValidIngredientsPermission,
		CreateValidIngredientGroupsPermission,
		UpdateValidIngredientGroupsPermission,
		ArchiveValidIngredientGroupsPermission,
		CreateValidPreparationsPermission,
		UpdateValidPreparationsPermission,
		ArchiveValidPreparationsPermission,
		CreateValidMeasurementUnitsPermission,
		UpdateValidMeasurementUnitsPermission,
		ArchiveValidMeasurementUnitsPermission,
		CreateValidMeasurementUnitConversionsPermission,
		UpdateValidMeasurementUnitConversionsPermission,
		ArchiveValidMeasurementUnitConversionsPermission,
		CreateValidIngredientPreparationsPermission,
		UpdateValidIngredientPreparationsPermission,
		ArchiveValidIngredientPreparationsPermission,
		CreateValidPrepTaskConfigsPermission,
		UpdateValidPrepTaskConfigsPermission,
		ArchiveValidPrepTaskConfigsPermission,
		CreateValidIngredientStateIngredientsPermission,
		UpdateValidIngredientStateIngredientsPermission,
		ArchiveValidIngredientStateIngredientsPermission,
		CreateValidPreparationInstrumentsPermission,
		UpdateValidPreparationInstrumentsPermission,
		ArchiveValidPreparationInstrumentsPermission,
		CreateValidPreparationVesselsPermission,
		UpdateValidPreparationVesselsPermission,
		ArchiveValidPreparationVesselsPermission,
		CreateValidIngredientMeasurementUnitsPermission,
		UpdateValidIngredientMeasurementUnitsPermission,
		ArchiveValidIngredientMeasurementUnitsPermission,
		CreateValidIngredientStatesPermission,
		UpdateValidIngredientStatesPermission,
		ArchiveValidIngredientStatesPermission,
	}

	// AccountAdminPermissions is every account admin permission.
	AccountAdminPermissions = []gorbac.Permission{
		UpdateAccountPermission,
		ArchiveAccountPermission,
		TransferAccountPermission,
		InviteUserToAccountPermission,
		ModifyMemberPermissionsForAccountPermission,
		RemoveMemberAccountPermission,
		CreateWebhooksPermission,
		UpdateWebhooksPermission,
		ArchiveWebhooksPermission,
		CreateIssueReportsPermission,
		UpdateIssueReportsPermission,
		ArchiveIssueReportsPermission,
		CreateMealPlansPermission,
		UpdateMealPlansPermission,
		ArchiveMealPlansPermission,
		CreateMealPlanEventsPermission,
		UpdateMealPlanEventsPermission,
		ArchiveMealPlanEventsPermission,
		CreateMealPlanOptionsPermission,
		UpdateMealPlanOptionsPermission,
		ArchiveMealPlanOptionsPermission,
		CreateAccountInstrumentOwnershipsPermission,
		UpdateAccountInstrumentOwnershipsPermission,
		ArchiveAccountInstrumentOwnershipsPermission,
		CreateWebhookTriggerConfigsPermission,
		ArchiveWebhookTriggerConfigsPermission,
		CreateWebhookTriggerEventsPermission,
		ReadWebhookTriggerEventsPermission,
		UpdateWebhookTriggerEventsPermission,
		ArchiveWebhookTriggerEventsPermission,
		CreateMealListsPermission,
		ReadMealListsPermission,
		UpdateMealListsPermission,
		ArchiveMealListsPermission,
		CreateRecipeListsPermission,
		ReadRecipeListsPermission,
		UpdateRecipeListsPermission,
		ArchiveRecipeListsPermission,
		CreateCheckoutSessionPermission,
		CancelSubscriptionPermission,
		ReadPurchasesPermission,
		ReadPaymentHistoryPermission,
		ReadSubscriptionsPermission,
	}

	// AccountMemberPermissions is every account member permission.
	AccountMemberPermissions = []gorbac.Permission{
		ReportAnalyticsEventsPermission,
		ReadWebhooksPermission,
		ReadIssueReportsPermission,
		ReadAuditLogEntriesPermission,
		ReadOAuth2ClientsPermission,
		ReadServiceSettingsPermission,
		SearchServiceSettingsPermission,
		CreateUploadedMediaPermission,
		ReadUploadedMediaPermission,
		UpdateUploadedMediaPermission,
		ArchiveUploadedMediaPermission,
		CreateMealsPermission,
		ReadMealsPermission,
		UpdateMealsPermission,
		ArchiveMealsPermission,
		CreateRecipesPermission,
		ReadRecipesPermission,
		SearchRecipesPermission,
		UpdateRecipesPermission,
		ArchiveRecipesPermission,
		CreateRecipeStepsPermission,
		ReadRecipeStepsPermission,
		SearchRecipeStepsPermission,
		UpdateRecipeStepsPermission,
		ArchiveRecipeStepsPermission,
		CreateRecipePrepTasksPermission,
		ReadRecipePrepTasksPermission,
		UpdateRecipePrepTasksPermission,
		ArchiveRecipePrepTasksPermission,
		CreateRecipeStepInstrumentsPermission,
		ReadRecipeStepInstrumentsPermission,
		SearchRecipeStepInstrumentsPermission,
		UpdateRecipeStepInstrumentsPermission,
		ArchiveRecipeStepInstrumentsPermission,
		CreateRecipeStepVesselsPermission,
		ReadRecipeStepVesselsPermission,
		SearchRecipeStepVesselsPermission,
		UpdateRecipeStepVesselsPermission,
		ArchiveRecipeStepVesselsPermission,
		CreateRecipeStepIngredientsPermission,
		ReadRecipeStepIngredientsPermission,
		SearchRecipeStepIngredientsPermission,
		UpdateRecipeStepIngredientsPermission,
		ArchiveRecipeStepIngredientsPermission,
		CreateRecipeStepCompletionConditionsPermission,
		ReadRecipeStepCompletionConditionsPermission,
		SearchRecipeStepCompletionConditionsPermission,
		UpdateRecipeStepCompletionConditionsPermission,
		ArchiveRecipeStepCompletionConditionsPermission,
		CreateRecipeStepProductsPermission,
		ReadRecipeStepProductsPermission,
		SearchRecipeStepProductsPermission,
		UpdateRecipeStepProductsPermission,
		ArchiveRecipeStepProductsPermission,
		ReadValidInstrumentsPermission,
		SearchValidInstrumentsPermission,
		ReadValidVesselsPermission,
		SearchValidVesselsPermission,
		ReadValidIngredientsPermission,
		SearchValidIngredientsPermission,
		ReadValidIngredientGroupsPermission,
		SearchValidIngredientGroupsPermission,
		ReadValidPreparationsPermission,
		SearchValidPreparationsPermission,
		ReadValidMeasurementUnitsPermission,
		SearchValidMeasurementUnitsPermission,
		ReadValidMeasurementUnitConversionsPermission,
		ReadValidIngredientPreparationsPermission,
		SearchValidIngredientPreparationsPermission,
		ReadValidIngredientStateIngredientsPermission,
		SearchValidIngredientStateIngredientsPermission,
		ReadValidPreparationInstrumentsPermission,
		SearchValidPreparationInstrumentsPermission,
		ReadValidPreparationVesselsPermission,
		SearchValidPreparationVesselsPermission,
		ReadValidIngredientMeasurementUnitsPermission,
		SearchValidIngredientMeasurementUnitsPermission,
		ReadMealPlansPermission,
		SearchMealPlansPermission,
		ReadMealPlanEventsPermission,
		ReadMealPlanOptionsPermission,
		SearchMealPlanOptionsPermission,
		ReadValidIngredientStatesPermission,
		ReadMealPlanGroceryListItemsPermission,
		UpdateMealPlanGroceryListItemsPermission,
		ArchiveMealPlanGroceryListItemsPermission,
		CreateMealPlanOptionVotesPermission,
		ReadMealPlanOptionVotesPermission,
		SearchMealPlanOptionVotesPermission,
		UpdateMealPlanOptionVotesPermission,
		ArchiveMealPlanOptionVotesPermission,
		CreateMealPlanRecipeOptionSelectionsPermission,
		ReadMealPlanRecipeOptionSelectionsPermission,
		UpdateMealPlanRecipeOptionSelectionsPermission,
		ArchiveMealPlanRecipeOptionSelectionsPermission,
		CreateServiceSettingConfigurationsPermission,
		ReadServiceSettingConfigurationsPermission,
		UpdateServiceSettingConfigurationsPermission,
		ArchiveServiceSettingConfigurationsPermission,
		ReadMealPlanTasksPermission,
		UpdateMealPlanTasksPermission,
		CreateUserIngredientPreferencesPermission,
		ReadUserIngredientPreferencesPermission,
		UpdateUserIngredientPreferencesPermission,
		ArchiveUserIngredientPreferencesPermission,
		ReadAccountInstrumentOwnershipsPermission,
		CreateRecipeRatingsPermission,
		ReadRecipeRatingsPermission,
		CreateCommentsPermission,
		ReadCommentsPermission,
		UpdateCommentsPermission,
		ArchiveCommentsPermission,
		UpdateRecipeRatingsPermission,
		ArchiveRecipeRatingsPermission,
		ReadUserNotificationsPermission,
		UpdateUserNotificationsPermission,
		CreateUserDeviceTokensPermission,
		ReadUserDeviceTokensPermission,
		ArchiveUserDeviceTokensPermission,
		CreateWaitlistSignupsPermission,
		UpdateWaitlistSignupsPermission,
		ArchiveWaitlistSignupsPermission,
		ReadWaitlistsPermission,
		ReadWaitlistSignupsPermission,
		CreateWaitlistSignupsPermission,
		UpdateWaitlistSignupsPermission,
		ArchiveWaitlistSignupsPermission,
		ReadValidPrepTaskConfigsPermission,
		CreateCheckoutSessionPermission,
		CancelSubscriptionPermission,
		ReadPurchasesPermission,
		ReadPaymentHistoryPermission,
		ReadSubscriptionsPermission,
	}
)

func init() {
	// assign service admin permissions.
	for _, perm := range ServiceAdminPermissions {
		must(serviceAdmin.Assign(perm))
	}

	// assign service data admin permissions.
	for _, perm := range ServiceDataAdminPermissions {
		must(serviceDataAdmin.Assign(perm))
		must(serviceAdmin.Assign(perm)) // these aren't separate things yet
	}

	// assign account admin permissions.
	for _, perm := range AccountAdminPermissions {
		must(accountAdmin.Assign(perm))
	}

	// assign account member permissions.
	for _, perm := range AccountMemberPermissions {
		must(accountMember.Assign(perm))
	}
}
