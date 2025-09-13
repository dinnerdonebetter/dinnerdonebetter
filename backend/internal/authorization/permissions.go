package authorization

import (
	"github.com/mikespook/gorbac/v2"
)

var _ gorbac.Permission = (*Permission)(nil)

type (
	// Permission is a simple string alias.
	Permission string
)

const (
	// PublishArbitraryQueueMessagesPermission is a service admin permission.
	PublishArbitraryQueueMessagesPermission Permission = "admin.publish_queue_messages"
	// UpdateUserStatusPermission is a service admin permission.
	UpdateUserStatusPermission Permission = "update.user_status"
	// ImpersonateUserPermission is a service admin permission.
	ImpersonateUserPermission Permission = "imitate.user"
	// ReadUserPermission is a service admin permission.
	ReadUserPermission Permission = "read.user"
	// SearchUserPermission is a service admin permission.
	SearchUserPermission Permission = "search.user"
	// ArchiveUserPermission is a service admin permission.
	ArchiveUserPermission Permission = "archive.user"

	// UpdateAccountPermission is an account admin permission.
	UpdateAccountPermission Permission = "update.account"
	// ArchiveAccountPermission is an account admin permission.
	ArchiveAccountPermission Permission = "archive.account"
	// InviteUserToAccountPermission is an account admin permission.
	InviteUserToAccountPermission Permission = "account.add.member"
	// ModifyMemberPermissionsForAccountPermission is an account admin permission.
	ModifyMemberPermissionsForAccountPermission Permission = "account.membership.modify"
	// PublishArbitraryQueueMessagePermission is a service admin permission.
	PublishArbitraryQueueMessagePermission Permission = "queues.publish.message"
	// RemoveMemberAccountPermission is an account admin permission.
	RemoveMemberAccountPermission Permission = "remove_member.account"
	// TransferAccountPermission is an account admin permission.
	TransferAccountPermission Permission = "transfer.account"

	// CreateWebhooksPermission is an account admin permission.
	CreateWebhooksPermission Permission = "create.webhooks"
	// ReadWebhooksPermission is an account admin permission.
	ReadWebhooksPermission Permission = "read.webhooks"
	// UpdateWebhooksPermission is an account admin permission.
	UpdateWebhooksPermission Permission = "update.webhooks"
	// ArchiveWebhooksPermission is an account admin permission.
	ArchiveWebhooksPermission Permission = "archive.webhooks"
	// CreateWebhookTriggerEventsPermission is an account admin permission.
	CreateWebhookTriggerEventsPermission Permission = "create.webhook_trigger_events"
	// ArchiveWebhookTriggerEventsPermission is an account admin permission.
	ArchiveWebhookTriggerEventsPermission Permission = "archive.webhook_trigger_events"

	// ReadAuditLogEntriesPermission is a service permission.
	ReadAuditLogEntriesPermission Permission = "read.audit_log_entries"

	// CreateValidInstrumentsPermission is a permission.
	CreateValidInstrumentsPermission Permission = "create.valid_instruments"
	// ReadValidInstrumentsPermission is a permission.
	ReadValidInstrumentsPermission Permission = "read.valid_instruments"
	// SearchValidInstrumentsPermission is a permission.
	SearchValidInstrumentsPermission Permission = "search.valid_instruments"
	// UpdateValidInstrumentsPermission is a permission.
	UpdateValidInstrumentsPermission Permission = "update.valid_instruments"
	// ArchiveValidInstrumentsPermission is a permission.
	ArchiveValidInstrumentsPermission Permission = "archive.valid_instruments"

	// CreateValidVesselsPermission is a permission.
	CreateValidVesselsPermission Permission = "create.valid_vessels"
	// ReadValidVesselsPermission is a permission.
	ReadValidVesselsPermission Permission = "read.valid_vessels"
	// SearchValidVesselsPermission is a permission.
	SearchValidVesselsPermission Permission = "search.valid_vessels"
	// UpdateValidVesselsPermission is a permission.
	UpdateValidVesselsPermission Permission = "update.valid_vessels"
	// ArchiveValidVesselsPermission is a permission.
	ArchiveValidVesselsPermission Permission = "archive.valid_vessels"

	// CreateValidIngredientsPermission is a permission.
	CreateValidIngredientsPermission Permission = "create.valid_ingredients"
	// ReadValidIngredientsPermission is a permission.
	ReadValidIngredientsPermission Permission = "read.valid_ingredients"
	// SearchValidIngredientsPermission is a permission.
	SearchValidIngredientsPermission Permission = "search.valid_ingredients"
	// UpdateValidIngredientsPermission is a permission.
	UpdateValidIngredientsPermission Permission = "update.valid_ingredients"
	// ArchiveValidIngredientsPermission is a permission.
	ArchiveValidIngredientsPermission Permission = "archive.valid_ingredients"

	// CreateValidIngredientGroupsPermission is a permission.
	CreateValidIngredientGroupsPermission Permission = "create.valid_ingredient_groups"
	// ReadValidIngredientGroupsPermission is a permission.
	ReadValidIngredientGroupsPermission Permission = "read.valid_ingredient_groups"
	// SearchValidIngredientGroupsPermission is a permission.
	SearchValidIngredientGroupsPermission Permission = "search.valid_ingredient_groups"
	// UpdateValidIngredientGroupsPermission is a permission.
	UpdateValidIngredientGroupsPermission Permission = "update.valid_ingredient_groups"
	// ArchiveValidIngredientGroupsPermission is a permission.
	ArchiveValidIngredientGroupsPermission Permission = "archive.valid_ingredient_groups"

	// CreateValidPreparationsPermission is a permission.
	CreateValidPreparationsPermission Permission = "create.valid_preparations"
	// ReadValidPreparationsPermission is a permission.
	ReadValidPreparationsPermission Permission = "read.valid_preparations"
	// SearchValidPreparationsPermission is a permission.
	SearchValidPreparationsPermission Permission = "search.valid_preparations"
	// UpdateValidPreparationsPermission is a permission.
	UpdateValidPreparationsPermission Permission = "update.valid_preparations"
	// ArchiveValidPreparationsPermission is a permission.
	ArchiveValidPreparationsPermission Permission = "archive.valid_preparations"

	// CreateValidMeasurementUnitsPermission is a permission.
	CreateValidMeasurementUnitsPermission Permission = "create.measurement_units"
	// ReadValidMeasurementUnitsPermission is a permission.
	ReadValidMeasurementUnitsPermission Permission = "read.measurement_units"
	// SearchValidMeasurementUnitsPermission is a permission.
	SearchValidMeasurementUnitsPermission Permission = "search.measurement_units"
	// UpdateValidMeasurementUnitsPermission is a permission.
	UpdateValidMeasurementUnitsPermission Permission = "update.measurement_units"
	// ArchiveValidMeasurementUnitsPermission is a permission.
	ArchiveValidMeasurementUnitsPermission Permission = "archive.measurement_units"

	// CreateValidIngredientStatesPermission is a permission.
	CreateValidIngredientStatesPermission Permission = "create.valid_ingredient_states"
	// ReadValidIngredientStatesPermission is a permission.
	ReadValidIngredientStatesPermission Permission = "read.valid_ingredient_states"
	// UpdateValidIngredientStatesPermission is a permission.
	UpdateValidIngredientStatesPermission Permission = "update.valid_ingredient_states"
	// ArchiveValidIngredientStatesPermission is a permission.
	ArchiveValidIngredientStatesPermission Permission = "archive.valid_ingredient_states"

	// CreateValidMeasurementUnitConversionsPermission is a permission.
	CreateValidMeasurementUnitConversionsPermission Permission = "create.measurement_conversions"
	// ReadValidMeasurementUnitConversionsPermission is a permission.
	ReadValidMeasurementUnitConversionsPermission Permission = "read.measurement_conversions"
	// UpdateValidMeasurementUnitConversionsPermission is a permission.
	UpdateValidMeasurementUnitConversionsPermission Permission = "update.measurement_conversions"
	// ArchiveValidMeasurementUnitConversionsPermission is a permission.
	ArchiveValidMeasurementUnitConversionsPermission Permission = "archive.measurement_conversions"

	// CreateValidIngredientPreparationsPermission is a permission.
	CreateValidIngredientPreparationsPermission Permission = "create.valid_ingredient_preparations"
	// ReadValidIngredientPreparationsPermission is a permission.
	ReadValidIngredientPreparationsPermission Permission = "read.valid_ingredient_preparations"
	// SearchValidIngredientPreparationsPermission is a permission.
	SearchValidIngredientPreparationsPermission Permission = "search.valid_ingredient_preparations"
	// UpdateValidIngredientPreparationsPermission is a permission.
	UpdateValidIngredientPreparationsPermission Permission = "update.valid_ingredient_preparations"
	// ArchiveValidIngredientPreparationsPermission is a permission.
	ArchiveValidIngredientPreparationsPermission Permission = "archive.valid_ingredient_preparations"

	// CreateValidIngredientStateIngredientsPermission is a permission.
	CreateValidIngredientStateIngredientsPermission Permission = "create.valid_ingredient_state_ingredients"
	// ReadValidIngredientStateIngredientsPermission is a permission.
	ReadValidIngredientStateIngredientsPermission Permission = "read.valid_ingredient_state_ingredients"
	// SearchValidIngredientStateIngredientsPermission is a permission.
	SearchValidIngredientStateIngredientsPermission Permission = "search.valid_ingredient_state_ingredients"
	// UpdateValidIngredientStateIngredientsPermission is a permission.
	UpdateValidIngredientStateIngredientsPermission Permission = "update.valid_ingredient_state_ingredients"
	// ArchiveValidIngredientStateIngredientsPermission is a permission.
	ArchiveValidIngredientStateIngredientsPermission Permission = "archive.valid_ingredient_state_ingredients"

	// CreateValidPreparationInstrumentsPermission is a permission.
	CreateValidPreparationInstrumentsPermission Permission = "create.valid_preparation_instruments"
	// ReadValidPreparationInstrumentsPermission is a permission.
	ReadValidPreparationInstrumentsPermission Permission = "read.valid_preparation_instruments"
	// SearchValidPreparationInstrumentsPermission is a permission.
	SearchValidPreparationInstrumentsPermission Permission = "search.valid_preparation_instruments"
	// UpdateValidPreparationInstrumentsPermission is a permission.
	UpdateValidPreparationInstrumentsPermission Permission = "update.valid_preparation_instruments"
	// ArchiveValidPreparationInstrumentsPermission is a permission.
	ArchiveValidPreparationInstrumentsPermission Permission = "archive.valid_preparation_instruments"

	// CreateValidPreparationVesselsPermission is a permission.
	CreateValidPreparationVesselsPermission Permission = "create.valid_preparation_vessels"
	// ReadValidPreparationVesselsPermission is a permission.
	ReadValidPreparationVesselsPermission Permission = "read.valid_preparation_vessels"
	// SearchValidPreparationVesselsPermission is a permission.
	SearchValidPreparationVesselsPermission Permission = "search.valid_preparation_vessels"
	// UpdateValidPreparationVesselsPermission is a permission.
	UpdateValidPreparationVesselsPermission Permission = "update.valid_preparation_vessels"
	// ArchiveValidPreparationVesselsPermission is a permission.
	ArchiveValidPreparationVesselsPermission Permission = "archive.valid_preparation_vessels"

	// CreateValidIngredientMeasurementUnitsPermission is a permission.
	CreateValidIngredientMeasurementUnitsPermission Permission = "create.valid_ingredient_measurement_units"
	// ReadValidIngredientMeasurementUnitsPermission is a permission.
	ReadValidIngredientMeasurementUnitsPermission Permission = "read.valid_ingredient_measurement_units"
	// SearchValidIngredientMeasurementUnitsPermission is a permission.
	SearchValidIngredientMeasurementUnitsPermission Permission = "search.valid_ingredient_measurement_units"
	// UpdateValidIngredientMeasurementUnitsPermission is a permission.
	UpdateValidIngredientMeasurementUnitsPermission Permission = "update.valid_ingredient_measurement_units"
	// ArchiveValidIngredientMeasurementUnitsPermission is a permission.
	ArchiveValidIngredientMeasurementUnitsPermission Permission = "archive.valid_ingredient_measurement_units"

	// CreateMealsPermission is a permission.
	CreateMealsPermission Permission = "create.meals"
	// ReadMealsPermission is a permission.
	ReadMealsPermission Permission = "read.meals"
	// UpdateMealsPermission is a permission.
	UpdateMealsPermission Permission = "update.meals"
	// ArchiveMealsPermission is a permission.
	ArchiveMealsPermission Permission = "archive.meals"

	// TODO: clone.recipes permission

	// CreateRecipesPermission is a permission.
	CreateRecipesPermission Permission = "create.recipes"
	// ReadRecipesPermission is a permission.
	ReadRecipesPermission Permission = "read.recipes"
	// SearchRecipesPermission is a permission.
	SearchRecipesPermission Permission = "search.recipes"
	// UpdateRecipesPermission is a permission.
	UpdateRecipesPermission Permission = "update.recipes"
	// ArchiveRecipesPermission is a permission.
	ArchiveRecipesPermission Permission = "archive.recipes"

	// CreateRecipePrepTasksPermission is a permission.
	CreateRecipePrepTasksPermission Permission = "create.recipe_prep_tasks"
	// ReadRecipePrepTasksPermission is a permission.
	ReadRecipePrepTasksPermission Permission = "read.recipe_prep_tasks"
	// UpdateRecipePrepTasksPermission is a permission.
	UpdateRecipePrepTasksPermission Permission = "update.recipe_prep_tasks"
	// ArchiveRecipePrepTasksPermission is a permission.
	ArchiveRecipePrepTasksPermission Permission = "archive.recipe_prep_tasks"

	// CreateRecipeStepsPermission is a permission.
	CreateRecipeStepsPermission Permission = "create.recipe_steps"
	// ReadRecipeStepsPermission is a permission.
	ReadRecipeStepsPermission Permission = "read.recipe_steps"
	// SearchRecipeStepsPermission is a permission.
	SearchRecipeStepsPermission Permission = "search.recipe_steps"
	// UpdateRecipeStepsPermission is a permission.
	UpdateRecipeStepsPermission Permission = "update.recipe_steps"
	// ArchiveRecipeStepsPermission is a permission.
	ArchiveRecipeStepsPermission Permission = "archive.recipe_steps"

	// CreateRecipeStepInstrumentsPermission is a permission.
	CreateRecipeStepInstrumentsPermission Permission = "create.recipe_step_instruments"
	// ReadRecipeStepInstrumentsPermission is a permission.
	ReadRecipeStepInstrumentsPermission Permission = "read.recipe_step_instruments"
	// SearchRecipeStepInstrumentsPermission is a permission.
	SearchRecipeStepInstrumentsPermission Permission = "search.recipe_step_instruments"
	// UpdateRecipeStepInstrumentsPermission is a permission.
	UpdateRecipeStepInstrumentsPermission Permission = "update.recipe_step_instruments"
	// ArchiveRecipeStepInstrumentsPermission is a permission.
	ArchiveRecipeStepInstrumentsPermission Permission = "archive.recipe_step_instruments"

	// CreateRecipeStepVesselsPermission is a permission.
	CreateRecipeStepVesselsPermission Permission = "create.recipe_step_vessels"
	// ReadRecipeStepVesselsPermission is a permission.
	ReadRecipeStepVesselsPermission Permission = "read.recipe_step_vessels"
	// SearchRecipeStepVesselsPermission is a permission.
	SearchRecipeStepVesselsPermission Permission = "search.recipe_step_vessels"
	// UpdateRecipeStepVesselsPermission is a permission.
	UpdateRecipeStepVesselsPermission Permission = "update.recipe_step_vessels"
	// ArchiveRecipeStepVesselsPermission is a permission.
	ArchiveRecipeStepVesselsPermission Permission = "archive.recipe_step_vessels"

	// CreateRecipeStepIngredientsPermission is a permission.
	CreateRecipeStepIngredientsPermission Permission = "create.recipe_step_ingredients"
	// ReadRecipeStepIngredientsPermission is a permission.
	ReadRecipeStepIngredientsPermission Permission = "read.recipe_step_ingredients"
	// SearchRecipeStepIngredientsPermission is a permission.
	SearchRecipeStepIngredientsPermission Permission = "search.recipe_step_ingredients"
	// UpdateRecipeStepIngredientsPermission is a permission.
	UpdateRecipeStepIngredientsPermission Permission = "update.recipe_step_ingredients"
	// ArchiveRecipeStepIngredientsPermission is a permission.
	ArchiveRecipeStepIngredientsPermission Permission = "archive.recipe_step_ingredients"

	// CreateRecipeStepCompletionConditionsPermission is a permission.
	CreateRecipeStepCompletionConditionsPermission Permission = "create.recipe_step_completion_conditions"
	// ReadRecipeStepCompletionConditionsPermission is a permission.
	ReadRecipeStepCompletionConditionsPermission Permission = "read.recipe_step_completion_conditions"
	// SearchRecipeStepCompletionConditionsPermission is a permission.
	SearchRecipeStepCompletionConditionsPermission Permission = "search.recipe_step_completion_conditions"
	// UpdateRecipeStepCompletionConditionsPermission is a permission.
	UpdateRecipeStepCompletionConditionsPermission Permission = "update.recipe_step_completion_conditions"
	// ArchiveRecipeStepCompletionConditionsPermission is a permission.
	ArchiveRecipeStepCompletionConditionsPermission Permission = "archive.recipe_step_completion_conditions"

	// CreateRecipeStepProductsPermission is a permission.
	CreateRecipeStepProductsPermission Permission = "create.recipe_step_products"
	// ReadRecipeStepProductsPermission is a permission.
	ReadRecipeStepProductsPermission Permission = "read.recipe_step_products"
	// SearchRecipeStepProductsPermission is a permission.
	SearchRecipeStepProductsPermission Permission = "search.recipe_step_products"
	// UpdateRecipeStepProductsPermission is a permission.
	UpdateRecipeStepProductsPermission Permission = "update.recipe_step_products"
	// ArchiveRecipeStepProductsPermission is a permission.
	ArchiveRecipeStepProductsPermission Permission = "archive.recipe_step_products"

	// CreateMealPlansPermission is a permission.
	CreateMealPlansPermission Permission = "create.meal_plans"
	// ReadMealPlansPermission is a permission.
	ReadMealPlansPermission Permission = "read.meal_plans"
	// SearchMealPlansPermission is a permission.
	SearchMealPlansPermission Permission = "search.meal_plans"
	// UpdateMealPlansPermission is a permission.
	UpdateMealPlansPermission Permission = "update.meal_plans"
	// ArchiveMealPlansPermission is a permission.
	ArchiveMealPlansPermission Permission = "archive.meal_plans"

	// CreateMealPlanEventsPermission is a permission.
	CreateMealPlanEventsPermission Permission = "create.meal_plan_events"
	// ReadMealPlanEventsPermission is a permission.
	ReadMealPlanEventsPermission Permission = "read.meal_plan_events"
	// UpdateMealPlanEventsPermission is a permission.
	UpdateMealPlanEventsPermission Permission = "update.meal_plan_events"
	// ArchiveMealPlanEventsPermission is a permission.
	ArchiveMealPlanEventsPermission Permission = "archive.meal_plan_events"

	// CreateMealPlanOptionsPermission is a permission.
	CreateMealPlanOptionsPermission Permission = "create.meal_plan_options"
	// ReadMealPlanOptionsPermission is a permission.
	ReadMealPlanOptionsPermission Permission = "read.meal_plan_options"
	// SearchMealPlanOptionsPermission is a permission.
	SearchMealPlanOptionsPermission Permission = "search.meal_plan_options"
	// UpdateMealPlanOptionsPermission is a permission.
	UpdateMealPlanOptionsPermission Permission = "update.meal_plan_options"
	// ArchiveMealPlanOptionsPermission is a permission.
	ArchiveMealPlanOptionsPermission Permission = "archive.meal_plan_options"

	// CreateMealPlanGroceryListItemsPermission is a permission.
	CreateMealPlanGroceryListItemsPermission Permission = "create.meal_plan_grocery_list_items"
	// ReadMealPlanGroceryListItemsPermission is a permission.
	ReadMealPlanGroceryListItemsPermission Permission = "read.meal_plan_grocery_list_items"
	// UpdateMealPlanGroceryListItemsPermission is a permission.
	UpdateMealPlanGroceryListItemsPermission Permission = "update.meal_plan_grocery_list_items"
	// ArchiveMealPlanGroceryListItemsPermission is a permission.
	ArchiveMealPlanGroceryListItemsPermission Permission = "archive.meal_plan_grocery_list_items"

	// CreateMealPlanOptionVotesPermission is a permission.
	CreateMealPlanOptionVotesPermission Permission = "create.meal_plan_option_votes"
	// ReadMealPlanOptionVotesPermission is a permission.
	ReadMealPlanOptionVotesPermission Permission = "read.meal_plan_option_votes"
	// SearchMealPlanOptionVotesPermission is a permission.
	SearchMealPlanOptionVotesPermission Permission = "search.meal_plan_option_votes"
	// UpdateMealPlanOptionVotesPermission is a permission.
	UpdateMealPlanOptionVotesPermission Permission = "update.meal_plan_option_votes"
	// ArchiveMealPlanOptionVotesPermission is a permission.
	ArchiveMealPlanOptionVotesPermission Permission = "archive.meal_plan_option_votes"

	// ReadMealPlanTasksPermission is a permission.
	ReadMealPlanTasksPermission Permission = "read.meal_plan_tasks"
	// CreateMealPlanTasksPermission is a permission.
	CreateMealPlanTasksPermission Permission = "create.meal_plan_tasks"
	// UpdateMealPlanTasksPermission is a permission.
	UpdateMealPlanTasksPermission Permission = "update.meal_plan_tasks"

	// CreateServiceSettingsPermission is an admin user permission.
	CreateServiceSettingsPermission Permission = "create.service_settings"
	// ReadServiceSettingsPermission is an admin user permission.
	ReadServiceSettingsPermission Permission = "read.service_settings"
	// SearchServiceSettingsPermission is an admin user permission.
	SearchServiceSettingsPermission Permission = "search.service_settings"
	// ArchiveServiceSettingsPermission is an admin user permission.
	ArchiveServiceSettingsPermission Permission = "archive.service_settings"

	// CreateServiceSettingConfigurationsPermission is an admin user permission.
	CreateServiceSettingConfigurationsPermission Permission = "create.service_setting_configurations"
	// ReadServiceSettingConfigurationsPermission is an admin user permission.
	ReadServiceSettingConfigurationsPermission Permission = "read.service_setting_configurations"
	// UpdateServiceSettingConfigurationsPermission is an admin user permission.
	UpdateServiceSettingConfigurationsPermission Permission = "update.service_setting_configurations"
	// ArchiveServiceSettingConfigurationsPermission is an admin user permission.
	ArchiveServiceSettingConfigurationsPermission Permission = "archive.service_setting_configurations"

	// CreateUserIngredientPreferencesPermission is a permission.
	CreateUserIngredientPreferencesPermission Permission = "create.user_ingredient_preferences"
	// ReadUserIngredientPreferencesPermission is a permission.
	ReadUserIngredientPreferencesPermission Permission = "read.user_ingredient_preferences"
	// UpdateUserIngredientPreferencesPermission is a permission.
	UpdateUserIngredientPreferencesPermission Permission = "update.user_ingredient_preferences"
	// ArchiveUserIngredientPreferencesPermission is a permission.
	ArchiveUserIngredientPreferencesPermission Permission = "archive.user_ingredient_preferences"

	// CreateAccountInstrumentOwnershipsPermission is a permission.
	CreateAccountInstrumentOwnershipsPermission Permission = "create.account_instrument_ownerships"
	// ReadAccountInstrumentOwnershipsPermission is a permission.
	ReadAccountInstrumentOwnershipsPermission Permission = "read.account_instrument_ownerships"
	// UpdateAccountInstrumentOwnershipsPermission is a permission.
	UpdateAccountInstrumentOwnershipsPermission Permission = "update.account_instrument_ownerships"
	// ArchiveAccountInstrumentOwnershipsPermission is a permission.
	ArchiveAccountInstrumentOwnershipsPermission Permission = "archive.account_instrument_ownerships"

	// CreateRecipeRatingsPermission is a permission.
	CreateRecipeRatingsPermission Permission = "create.recipe_ratings"
	// ReadRecipeRatingsPermission is a permission.
	ReadRecipeRatingsPermission Permission = "read.recipe_ratings"
	// UpdateRecipeRatingsPermission is a permission.
	UpdateRecipeRatingsPermission Permission = "update.recipe_ratings"
	// ArchiveRecipeRatingsPermission is a permission.
	ArchiveRecipeRatingsPermission Permission = "archive.recipe_ratings"

	// CreateOAuth2ClientsPermission is an account admin permission.
	CreateOAuth2ClientsPermission Permission = "create.oauth2_clients"
	// ReadOAuth2ClientsPermission is an account admin permission.
	ReadOAuth2ClientsPermission Permission = "read.oauth2_clients"
	// ArchiveOAuth2ClientsPermission is an account admin permission.
	ArchiveOAuth2ClientsPermission Permission = "archive.oauth2_clients"

	// CreateUserNotificationsPermission is an admin user permission.
	CreateUserNotificationsPermission Permission = "create.user_notifications"
	// ReadUserNotificationsPermission is a permission.
	ReadUserNotificationsPermission Permission = "read.user_notifications"
	// UpdateUserNotificationsPermission is a permission.
	UpdateUserNotificationsPermission Permission = "update.user_notifications"
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
		PublishArbitraryQueueMessagesPermission,
		UpdateUserStatusPermission,
		ReadUserPermission,
		SearchUserPermission,
		ArchiveUserPermission,
		CreateOAuth2ClientsPermission,
		ArchiveOAuth2ClientsPermission,
		ArchiveServiceSettingsPermission,
		CreateRecipesPermission,
		CreateUserNotificationsPermission,
		ImpersonateUserPermission,
		PublishArbitraryQueueMessagePermission,
		// only admins can arbitrarily create these via the API, this is exclusively for integration test purposes.
		CreateServiceSettingsPermission,
		CreateMealPlanTasksPermission,
		CreateMealPlanGroceryListItemsPermission,
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
		CreateWebhookTriggerEventsPermission,
		ArchiveWebhookTriggerEventsPermission,
	}

	// AccountMemberPermissions is every account member permission.
	AccountMemberPermissions = []gorbac.Permission{
		ReadWebhooksPermission,
		ReadAuditLogEntriesPermission,
		ReadOAuth2ClientsPermission,
		ReadServiceSettingsPermission,
		SearchServiceSettingsPermission,
		CreateMealsPermission,
		ReadMealsPermission,
		UpdateMealsPermission,
		ArchiveMealsPermission,
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
		UpdateRecipeRatingsPermission,
		ArchiveRecipeRatingsPermission,
		ReadUserNotificationsPermission,
		UpdateUserNotificationsPermission,
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
