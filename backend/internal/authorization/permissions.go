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

	// UpdateAccountPermission is a account admin permission.
	UpdateAccountPermission Permission = "update.account"
	// ArchiveAccountPermission is a account admin permission.
	ArchiveAccountPermission Permission = "archive.account"
	// InviteUserToAccountPermission is a account admin permission.
	InviteUserToAccountPermission Permission = "account.add.member"
	// ModifyMemberPermissionsForAccountPermission is a account admin permission.
	ModifyMemberPermissionsForAccountPermission Permission = "account.membership.modify"
	// RemoveMemberAccountPermission is a account admin permission.
	RemoveMemberAccountPermission Permission = "remove_member.account"
	// TransferAccountPermission is a account admin permission.
	TransferAccountPermission Permission = "transfer.account"

	// CreateWebhooksPermission is a account admin permission.
	CreateWebhooksPermission Permission = "create.webhooks"
	// ReadWebhooksPermission is a account admin permission.
	ReadWebhooksPermission Permission = "read.webhooks"
	// UpdateWebhooksPermission is a account admin permission.
	UpdateWebhooksPermission Permission = "update.webhooks"
	// ArchiveWebhooksPermission is a account admin permission.
	ArchiveWebhooksPermission Permission = "archive.webhooks"
	// CreateWebhookTriggerEventsPermission is a account admin permission.
	CreateWebhookTriggerEventsPermission Permission = "create.webhook_trigger_events"
	// ArchiveWebhookTriggerEventsPermission is a account admin permission.
	ArchiveWebhookTriggerEventsPermission Permission = "archive.webhook_trigger_events"

	// ReadAuditLogEntriesPermission is a service permission.
	ReadAuditLogEntriesPermission Permission = "read.audit_log_entries"

	// CreateValidInstrumentsPermission is a account user permission.
	CreateValidInstrumentsPermission Permission = "create.valid_instruments"
	// ReadValidInstrumentsPermission is a account user permission.
	ReadValidInstrumentsPermission Permission = "read.valid_instruments"
	// SearchValidInstrumentsPermission is a account user permission.
	SearchValidInstrumentsPermission Permission = "search.valid_instruments"
	// UpdateValidInstrumentsPermission is a account user permission.
	UpdateValidInstrumentsPermission Permission = "update.valid_instruments"
	// ArchiveValidInstrumentsPermission is a account user permission.
	ArchiveValidInstrumentsPermission Permission = "archive.valid_instruments"

	// CreateValidVesselsPermission is a account user permission.
	CreateValidVesselsPermission Permission = "create.valid_vessels"
	// ReadValidVesselsPermission is a account user permission.
	ReadValidVesselsPermission Permission = "read.valid_vessels"
	// SearchValidVesselsPermission is a account user permission.
	SearchValidVesselsPermission Permission = "search.valid_vessels"
	// UpdateValidVesselsPermission is a account user permission.
	UpdateValidVesselsPermission Permission = "update.valid_vessels"
	// ArchiveValidVesselsPermission is a account user permission.
	ArchiveValidVesselsPermission Permission = "archive.valid_vessels"

	// CreateValidIngredientsPermission is a account user permission.
	CreateValidIngredientsPermission Permission = "create.valid_ingredients"
	// ReadValidIngredientsPermission is a account user permission.
	ReadValidIngredientsPermission Permission = "read.valid_ingredients"
	// SearchValidIngredientsPermission is a account user permission.
	SearchValidIngredientsPermission Permission = "search.valid_ingredients"
	// UpdateValidIngredientsPermission is a account user permission.
	UpdateValidIngredientsPermission Permission = "update.valid_ingredients"
	// ArchiveValidIngredientsPermission is a account user permission.
	ArchiveValidIngredientsPermission Permission = "archive.valid_ingredients"

	// CreateValidIngredientGroupsPermission is a account user permission.
	CreateValidIngredientGroupsPermission Permission = "create.valid_ingredient_groups"
	// ReadValidIngredientGroupsPermission is a account user permission.
	ReadValidIngredientGroupsPermission Permission = "read.valid_ingredient_groups"
	// SearchValidIngredientGroupsPermission is a account user permission.
	SearchValidIngredientGroupsPermission Permission = "search.valid_ingredient_groups"
	// UpdateValidIngredientGroupsPermission is a account user permission.
	UpdateValidIngredientGroupsPermission Permission = "update.valid_ingredient_groups"
	// ArchiveValidIngredientGroupsPermission is a account user permission.
	ArchiveValidIngredientGroupsPermission Permission = "archive.valid_ingredient_groups"

	// CreateValidPreparationsPermission is a account user permission.
	CreateValidPreparationsPermission Permission = "create.valid_preparations"
	// ReadValidPreparationsPermission is a account user permission.
	ReadValidPreparationsPermission Permission = "read.valid_preparations"
	// SearchValidPreparationsPermission is a account user permission.
	SearchValidPreparationsPermission Permission = "search.valid_preparations"
	// UpdateValidPreparationsPermission is a account user permission.
	UpdateValidPreparationsPermission Permission = "update.valid_preparations"
	// ArchiveValidPreparationsPermission is a account user permission.
	ArchiveValidPreparationsPermission Permission = "archive.valid_preparations"

	// CreateValidMeasurementUnitsPermission is a account user permission.
	CreateValidMeasurementUnitsPermission Permission = "create.measurement_units"
	// ReadValidMeasurementUnitsPermission is a account user permission.
	ReadValidMeasurementUnitsPermission Permission = "read.measurement_units"
	// SearchValidMeasurementUnitsPermission is a account user permission.
	SearchValidMeasurementUnitsPermission Permission = "search.measurement_units"
	// UpdateValidMeasurementUnitsPermission is a account user permission.
	UpdateValidMeasurementUnitsPermission Permission = "update.measurement_units"
	// ArchiveValidMeasurementUnitsPermission is a account user permission.
	ArchiveValidMeasurementUnitsPermission Permission = "archive.measurement_units"

	// CreateValidIngredientStatesPermission is a account user permission.
	CreateValidIngredientStatesPermission Permission = "create.valid_ingredient_states"
	// ReadValidIngredientStatesPermission is a account user permission.
	ReadValidIngredientStatesPermission Permission = "read.valid_ingredient_states"
	// UpdateValidIngredientStatesPermission is a account user permission.
	UpdateValidIngredientStatesPermission Permission = "update.valid_ingredient_states"
	// ArchiveValidIngredientStatesPermission is a account user permission.
	ArchiveValidIngredientStatesPermission Permission = "archive.valid_ingredient_states"

	// CreateValidMeasurementUnitConversionsPermission is a account user permission.
	CreateValidMeasurementUnitConversionsPermission Permission = "create.measurement_conversions"
	// ReadValidMeasurementUnitConversionsPermission is a account user permission.
	ReadValidMeasurementUnitConversionsPermission Permission = "read.measurement_conversions"
	// UpdateValidMeasurementUnitConversionsPermission is a account user permission.
	UpdateValidMeasurementUnitConversionsPermission Permission = "update.measurement_conversions"
	// ArchiveValidMeasurementUnitConversionsPermission is a account user permission.
	ArchiveValidMeasurementUnitConversionsPermission Permission = "archive.measurement_conversions"

	// CreateValidIngredientPreparationsPermission is a account user permission.
	CreateValidIngredientPreparationsPermission Permission = "create.valid_ingredient_preparations"
	// ReadValidIngredientPreparationsPermission is a account user permission.
	ReadValidIngredientPreparationsPermission Permission = "read.valid_ingredient_preparations"
	// SearchValidIngredientPreparationsPermission is a account user permission.
	SearchValidIngredientPreparationsPermission Permission = "search.valid_ingredient_preparations"
	// UpdateValidIngredientPreparationsPermission is a account user permission.
	UpdateValidIngredientPreparationsPermission Permission = "update.valid_ingredient_preparations"
	// ArchiveValidIngredientPreparationsPermission is a account user permission.
	ArchiveValidIngredientPreparationsPermission Permission = "archive.valid_ingredient_preparations"

	// CreateValidIngredientStateIngredientsPermission is a account user permission.
	CreateValidIngredientStateIngredientsPermission Permission = "create.valid_ingredient_state_ingredients"
	// ReadValidIngredientStateIngredientsPermission is a account user permission.
	ReadValidIngredientStateIngredientsPermission Permission = "read.valid_ingredient_state_ingredients"
	// SearchValidIngredientStateIngredientsPermission is a account user permission.
	SearchValidIngredientStateIngredientsPermission Permission = "search.valid_ingredient_state_ingredients"
	// UpdateValidIngredientStateIngredientsPermission is a account user permission.
	UpdateValidIngredientStateIngredientsPermission Permission = "update.valid_ingredient_state_ingredients"
	// ArchiveValidIngredientStateIngredientsPermission is a account user permission.
	ArchiveValidIngredientStateIngredientsPermission Permission = "archive.valid_ingredient_state_ingredients"

	// CreateValidPreparationInstrumentsPermission is a account user permission.
	CreateValidPreparationInstrumentsPermission Permission = "create.valid_preparation_instruments"
	// ReadValidPreparationInstrumentsPermission is a account user permission.
	ReadValidPreparationInstrumentsPermission Permission = "read.valid_preparation_instruments"
	// SearchValidPreparationInstrumentsPermission is a account user permission.
	SearchValidPreparationInstrumentsPermission Permission = "search.valid_preparation_instruments"
	// UpdateValidPreparationInstrumentsPermission is a account user permission.
	UpdateValidPreparationInstrumentsPermission Permission = "update.valid_preparation_instruments"
	// ArchiveValidPreparationInstrumentsPermission is a account user permission.
	ArchiveValidPreparationInstrumentsPermission Permission = "archive.valid_preparation_instruments"

	// CreateValidPreparationVesselsPermission is a account user permission.
	CreateValidPreparationVesselsPermission Permission = "create.valid_preparation_vessels"
	// ReadValidPreparationVesselsPermission is a account user permission.
	ReadValidPreparationVesselsPermission Permission = "read.valid_preparation_vessels"
	// SearchValidPreparationVesselsPermission is a account user permission.
	SearchValidPreparationVesselsPermission Permission = "search.valid_preparation_vessels"
	// UpdateValidPreparationVesselsPermission is a account user permission.
	UpdateValidPreparationVesselsPermission Permission = "update.valid_preparation_vessels"
	// ArchiveValidPreparationVesselsPermission is a account user permission.
	ArchiveValidPreparationVesselsPermission Permission = "archive.valid_preparation_vessels"

	// CreateValidIngredientMeasurementUnitsPermission is a account user permission.
	CreateValidIngredientMeasurementUnitsPermission Permission = "create.valid_ingredient_measurement_units"
	// ReadValidIngredientMeasurementUnitsPermission is a account user permission.
	ReadValidIngredientMeasurementUnitsPermission Permission = "read.valid_ingredient_measurement_units"
	// SearchValidIngredientMeasurementUnitsPermission is a account user permission.
	SearchValidIngredientMeasurementUnitsPermission Permission = "search.valid_ingredient_measurement_units"
	// UpdateValidIngredientMeasurementUnitsPermission is a account user permission.
	UpdateValidIngredientMeasurementUnitsPermission Permission = "update.valid_ingredient_measurement_units"
	// ArchiveValidIngredientMeasurementUnitsPermission is a account user permission.
	ArchiveValidIngredientMeasurementUnitsPermission Permission = "archive.valid_ingredient_measurement_units"

	// CreateMealsPermission is a account user permission.
	CreateMealsPermission Permission = "create.meals"
	// ReadMealsPermission is a account user permission.
	ReadMealsPermission Permission = "read.meals"
	// UpdateMealsPermission is a account user permission.
	UpdateMealsPermission Permission = "update.meals"
	// ArchiveMealsPermission is a account user permission.
	ArchiveMealsPermission Permission = "archive.meals"

	// CreateRecipesPermission is a account user permission.
	CreateRecipesPermission Permission = "create.recipes"
	// ReadRecipesPermission is a account user permission.
	ReadRecipesPermission Permission = "read.recipes"
	// SearchRecipesPermission is a account user permission.
	SearchRecipesPermission Permission = "search.recipes"
	// UpdateRecipesPermission is a account user permission.
	UpdateRecipesPermission Permission = "update.recipes"
	// ArchiveRecipesPermission is a account user permission.
	ArchiveRecipesPermission Permission = "archive.recipes"

	// CreateRecipePrepTasksPermission is a account user permission.
	CreateRecipePrepTasksPermission Permission = "create.recipe_prep_tasks"
	// ReadRecipePrepTasksPermission is a account user permission.
	ReadRecipePrepTasksPermission Permission = "read.recipe_prep_tasks"
	// UpdateRecipePrepTasksPermission is a account user permission.
	UpdateRecipePrepTasksPermission Permission = "update.recipe_prep_tasks"
	// ArchiveRecipePrepTasksPermission is a account user permission.
	ArchiveRecipePrepTasksPermission Permission = "archive.recipe_prep_tasks"

	// CreateRecipeStepsPermission is a account user permission.
	CreateRecipeStepsPermission Permission = "create.recipe_steps"
	// ReadRecipeStepsPermission is a account user permission.
	ReadRecipeStepsPermission Permission = "read.recipe_steps"
	// SearchRecipeStepsPermission is a account user permission.
	SearchRecipeStepsPermission Permission = "search.recipe_steps"
	// UpdateRecipeStepsPermission is a account user permission.
	UpdateRecipeStepsPermission Permission = "update.recipe_steps"
	// ArchiveRecipeStepsPermission is a account user permission.
	ArchiveRecipeStepsPermission Permission = "archive.recipe_steps"

	// CreateRecipeStepInstrumentsPermission is a account user permission.
	CreateRecipeStepInstrumentsPermission Permission = "create.recipe_step_instruments"
	// ReadRecipeStepInstrumentsPermission is a account user permission.
	ReadRecipeStepInstrumentsPermission Permission = "read.recipe_step_instruments"
	// SearchRecipeStepInstrumentsPermission is a account user permission.
	SearchRecipeStepInstrumentsPermission Permission = "search.recipe_step_instruments"
	// UpdateRecipeStepInstrumentsPermission is a account user permission.
	UpdateRecipeStepInstrumentsPermission Permission = "update.recipe_step_instruments"
	// ArchiveRecipeStepInstrumentsPermission is a account user permission.
	ArchiveRecipeStepInstrumentsPermission Permission = "archive.recipe_step_instruments"

	// CreateRecipeStepVesselsPermission is a account user permission.
	CreateRecipeStepVesselsPermission Permission = "create.recipe_step_vessels"
	// ReadRecipeStepVesselsPermission is a account user permission.
	ReadRecipeStepVesselsPermission Permission = "read.recipe_step_vessels"
	// SearchRecipeStepVesselsPermission is a account user permission.
	SearchRecipeStepVesselsPermission Permission = "search.recipe_step_vessels"
	// UpdateRecipeStepVesselsPermission is a account user permission.
	UpdateRecipeStepVesselsPermission Permission = "update.recipe_step_vessels"
	// ArchiveRecipeStepVesselsPermission is a account user permission.
	ArchiveRecipeStepVesselsPermission Permission = "archive.recipe_step_vessels"

	// CreateRecipeStepIngredientsPermission is a account user permission.
	CreateRecipeStepIngredientsPermission Permission = "create.recipe_step_ingredients"
	// ReadRecipeStepIngredientsPermission is a account user permission.
	ReadRecipeStepIngredientsPermission Permission = "read.recipe_step_ingredients"
	// SearchRecipeStepIngredientsPermission is a account user permission.
	SearchRecipeStepIngredientsPermission Permission = "search.recipe_step_ingredients"
	// UpdateRecipeStepIngredientsPermission is a account user permission.
	UpdateRecipeStepIngredientsPermission Permission = "update.recipe_step_ingredients"
	// ArchiveRecipeStepIngredientsPermission is a account user permission.
	ArchiveRecipeStepIngredientsPermission Permission = "archive.recipe_step_ingredients"

	// CreateRecipeStepCompletionConditionsPermission is a account user permission.
	CreateRecipeStepCompletionConditionsPermission Permission = "create.recipe_step_completion_conditions"
	// ReadRecipeStepCompletionConditionsPermission is a account user permission.
	ReadRecipeStepCompletionConditionsPermission Permission = "read.recipe_step_completion_conditions"
	// SearchRecipeStepCompletionConditionsPermission is a account user permission.
	SearchRecipeStepCompletionConditionsPermission Permission = "search.recipe_step_completion_conditions"
	// UpdateRecipeStepCompletionConditionsPermission is a account user permission.
	UpdateRecipeStepCompletionConditionsPermission Permission = "update.recipe_step_completion_conditions"
	// ArchiveRecipeStepCompletionConditionsPermission is a account user permission.
	ArchiveRecipeStepCompletionConditionsPermission Permission = "archive.recipe_step_completion_conditions"

	// CreateRecipeStepProductsPermission is a account user permission.
	CreateRecipeStepProductsPermission Permission = "create.recipe_step_products"
	// ReadRecipeStepProductsPermission is a account user permission.
	ReadRecipeStepProductsPermission Permission = "read.recipe_step_products"
	// SearchRecipeStepProductsPermission is a account user permission.
	SearchRecipeStepProductsPermission Permission = "search.recipe_step_products"
	// UpdateRecipeStepProductsPermission is a account user permission.
	UpdateRecipeStepProductsPermission Permission = "update.recipe_step_products"
	// ArchiveRecipeStepProductsPermission is a account user permission.
	ArchiveRecipeStepProductsPermission Permission = "archive.recipe_step_products"

	// CreateMealPlansPermission is a account user permission.
	CreateMealPlansPermission Permission = "create.meal_plans"
	// ReadMealPlansPermission is a account user permission.
	ReadMealPlansPermission Permission = "read.meal_plans"
	// SearchMealPlansPermission is a account user permission.
	SearchMealPlansPermission Permission = "search.meal_plans"
	// UpdateMealPlansPermission is a account user permission.
	UpdateMealPlansPermission Permission = "update.meal_plans"
	// ArchiveMealPlansPermission is a account user permission.
	ArchiveMealPlansPermission Permission = "archive.meal_plans"

	// CreateMealPlanEventsPermission is a account user permission.
	CreateMealPlanEventsPermission Permission = "create.meal_plan_events"
	// ReadMealPlanEventsPermission is a account user permission.
	ReadMealPlanEventsPermission Permission = "read.meal_plan_events"
	// UpdateMealPlanEventsPermission is a account user permission.
	UpdateMealPlanEventsPermission Permission = "update.meal_plan_events"
	// ArchiveMealPlanEventsPermission is a account user permission.
	ArchiveMealPlanEventsPermission Permission = "archive.meal_plan_events"

	// CreateMealPlanOptionsPermission is a account user permission.
	CreateMealPlanOptionsPermission Permission = "create.meal_plan_options"
	// ReadMealPlanOptionsPermission is a account user permission.
	ReadMealPlanOptionsPermission Permission = "read.meal_plan_options"
	// SearchMealPlanOptionsPermission is a account user permission.
	SearchMealPlanOptionsPermission Permission = "search.meal_plan_options"
	// UpdateMealPlanOptionsPermission is a account user permission.
	UpdateMealPlanOptionsPermission Permission = "update.meal_plan_options"
	// ArchiveMealPlanOptionsPermission is a account user permission.
	ArchiveMealPlanOptionsPermission Permission = "archive.meal_plan_options"

	// CreateMealPlanGroceryListItemsPermission is a account user permission.
	CreateMealPlanGroceryListItemsPermission Permission = "create.meal_plan_grocery_list_items"
	// ReadMealPlanGroceryListItemsPermission is a account user permission.
	ReadMealPlanGroceryListItemsPermission Permission = "read.meal_plan_grocery_list_items"
	// UpdateMealPlanGroceryListItemsPermission is a account user permission.
	UpdateMealPlanGroceryListItemsPermission Permission = "update.meal_plan_grocery_list_items"
	// ArchiveMealPlanGroceryListItemsPermission is a account user permission.
	ArchiveMealPlanGroceryListItemsPermission Permission = "archive.meal_plan_grocery_list_items"

	// CreateMealPlanOptionVotesPermission is a account user permission.
	CreateMealPlanOptionVotesPermission Permission = "create.meal_plan_option_votes"
	// ReadMealPlanOptionVotesPermission is a account user permission.
	ReadMealPlanOptionVotesPermission Permission = "read.meal_plan_option_votes"
	// SearchMealPlanOptionVotesPermission is a account user permission.
	SearchMealPlanOptionVotesPermission Permission = "search.meal_plan_option_votes"
	// UpdateMealPlanOptionVotesPermission is a account user permission.
	UpdateMealPlanOptionVotesPermission Permission = "update.meal_plan_option_votes"
	// ArchiveMealPlanOptionVotesPermission is a account user permission.
	ArchiveMealPlanOptionVotesPermission Permission = "archive.meal_plan_option_votes"

	// ReadMealPlanTasksPermission is a account user permission.
	ReadMealPlanTasksPermission Permission = "read.meal_plan_tasks"
	// CreateMealPlanTasksPermission is a account user permission.
	CreateMealPlanTasksPermission Permission = "create.meal_plan_tasks"
	// UpdateMealPlanTasksPermission is a account user permission.
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

	// CreateUserIngredientPreferencesPermission is a account user permission.
	CreateUserIngredientPreferencesPermission Permission = "create.user_ingredient_preferences"
	// ReadUserIngredientPreferencesPermission is a account user permission.
	ReadUserIngredientPreferencesPermission Permission = "read.user_ingredient_preferences"
	// UpdateUserIngredientPreferencesPermission is a account user permission.
	UpdateUserIngredientPreferencesPermission Permission = "update.user_ingredient_preferences"
	// ArchiveUserIngredientPreferencesPermission is a account user permission.
	ArchiveUserIngredientPreferencesPermission Permission = "archive.user_ingredient_preferences"

	// CreateAccountInstrumentOwnershipsPermission is a account user permission.
	CreateAccountInstrumentOwnershipsPermission Permission = "create.account_instrument_ownerships"
	// ReadAccountInstrumentOwnershipsPermission is a account user permission.
	ReadAccountInstrumentOwnershipsPermission Permission = "read.account_instrument_ownerships"
	// UpdateAccountInstrumentOwnershipsPermission is a account user permission.
	UpdateAccountInstrumentOwnershipsPermission Permission = "update.account_instrument_ownerships"
	// ArchiveAccountInstrumentOwnershipsPermission is a account user permission.
	ArchiveAccountInstrumentOwnershipsPermission Permission = "archive.account_instrument_ownerships"

	// CreateRecipeRatingsPermission is a account user permission.
	CreateRecipeRatingsPermission Permission = "create.recipe_ratings"
	// ReadRecipeRatingsPermission is a account user permission.
	ReadRecipeRatingsPermission Permission = "read.recipe_ratings"
	// UpdateRecipeRatingsPermission is a account user permission.
	UpdateRecipeRatingsPermission Permission = "update.recipe_ratings"
	// ArchiveRecipeRatingsPermission is a account user permission.
	ArchiveRecipeRatingsPermission Permission = "archive.recipe_ratings"

	// CreateOAuth2ClientsPermission is a account admin permission.
	CreateOAuth2ClientsPermission Permission = "create.oauth2_clients"
	// ReadOAuth2ClientsPermission is a account admin permission.
	ReadOAuth2ClientsPermission Permission = "read.oauth2_clients"
	// ArchiveOAuth2ClientsPermission is a account admin permission.
	ArchiveOAuth2ClientsPermission Permission = "archive.oauth2_clients"

	// CreateUserNotificationsPermission is an admin user permission.
	CreateUserNotificationsPermission Permission = "create.user_notifications"
	// ReadUserNotificationsPermission is a account user permission.
	ReadUserNotificationsPermission Permission = "read.user_notifications"
	// UpdateUserNotificationsPermission is a account user permission.
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
		CreateOAuth2ClientsPermission,
		ArchiveOAuth2ClientsPermission,
		ArchiveServiceSettingsPermission,
		CreateRecipesPermission,
		CreateUserNotificationsPermission,
		ImpersonateUserPermission,
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
