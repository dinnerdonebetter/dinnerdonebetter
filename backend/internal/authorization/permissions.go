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

	// UpdateAccountPermission is an account admin permission.
	UpdateAccountPermission Permission = "update.account"
	// ArchiveAccountPermission is an account admin permission.
	ArchiveAccountPermission Permission = "archive.account"
	// InviteUserToAccountPermission is an account admin permission.
	InviteUserToAccountPermission Permission = "account.add.member"
	// ModifyMemberPermissionsForAccountPermission is an account admin permission.
	ModifyMemberPermissionsForAccountPermission Permission = "account.membership.modify"
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

	// CreateValidInstrumentsPermission is an account user permission.
	CreateValidInstrumentsPermission Permission = "create.valid_instruments"
	// ReadValidInstrumentsPermission is an account user permission.
	ReadValidInstrumentsPermission Permission = "read.valid_instruments"
	// SearchValidInstrumentsPermission is an account user permission.
	SearchValidInstrumentsPermission Permission = "search.valid_instruments"
	// UpdateValidInstrumentsPermission is an account user permission.
	UpdateValidInstrumentsPermission Permission = "update.valid_instruments"
	// ArchiveValidInstrumentsPermission is an account user permission.
	ArchiveValidInstrumentsPermission Permission = "archive.valid_instruments"

	// CreateValidVesselsPermission is an account user permission.
	CreateValidVesselsPermission Permission = "create.valid_vessels"
	// ReadValidVesselsPermission is an account user permission.
	ReadValidVesselsPermission Permission = "read.valid_vessels"
	// SearchValidVesselsPermission is an account user permission.
	SearchValidVesselsPermission Permission = "search.valid_vessels"
	// UpdateValidVesselsPermission is an account user permission.
	UpdateValidVesselsPermission Permission = "update.valid_vessels"
	// ArchiveValidVesselsPermission is an account user permission.
	ArchiveValidVesselsPermission Permission = "archive.valid_vessels"

	// CreateValidIngredientsPermission is an account user permission.
	CreateValidIngredientsPermission Permission = "create.valid_ingredients"
	// ReadValidIngredientsPermission is an account user permission.
	ReadValidIngredientsPermission Permission = "read.valid_ingredients"
	// SearchValidIngredientsPermission is an account user permission.
	SearchValidIngredientsPermission Permission = "search.valid_ingredients"
	// UpdateValidIngredientsPermission is an account user permission.
	UpdateValidIngredientsPermission Permission = "update.valid_ingredients"
	// ArchiveValidIngredientsPermission is an account user permission.
	ArchiveValidIngredientsPermission Permission = "archive.valid_ingredients"

	// CreateValidIngredientGroupsPermission is an account user permission.
	CreateValidIngredientGroupsPermission Permission = "create.valid_ingredient_groups"
	// ReadValidIngredientGroupsPermission is an account user permission.
	ReadValidIngredientGroupsPermission Permission = "read.valid_ingredient_groups"
	// SearchValidIngredientGroupsPermission is an account user permission.
	SearchValidIngredientGroupsPermission Permission = "search.valid_ingredient_groups"
	// UpdateValidIngredientGroupsPermission is an account user permission.
	UpdateValidIngredientGroupsPermission Permission = "update.valid_ingredient_groups"
	// ArchiveValidIngredientGroupsPermission is an account user permission.
	ArchiveValidIngredientGroupsPermission Permission = "archive.valid_ingredient_groups"

	// CreateValidPreparationsPermission is an account user permission.
	CreateValidPreparationsPermission Permission = "create.valid_preparations"
	// ReadValidPreparationsPermission is an account user permission.
	ReadValidPreparationsPermission Permission = "read.valid_preparations"
	// SearchValidPreparationsPermission is an account user permission.
	SearchValidPreparationsPermission Permission = "search.valid_preparations"
	// UpdateValidPreparationsPermission is an account user permission.
	UpdateValidPreparationsPermission Permission = "update.valid_preparations"
	// ArchiveValidPreparationsPermission is an account user permission.
	ArchiveValidPreparationsPermission Permission = "archive.valid_preparations"

	// CreateValidMeasurementUnitsPermission is an account user permission.
	CreateValidMeasurementUnitsPermission Permission = "create.measurement_units"
	// ReadValidMeasurementUnitsPermission is an account user permission.
	ReadValidMeasurementUnitsPermission Permission = "read.measurement_units"
	// SearchValidMeasurementUnitsPermission is an account user permission.
	SearchValidMeasurementUnitsPermission Permission = "search.measurement_units"
	// UpdateValidMeasurementUnitsPermission is an account user permission.
	UpdateValidMeasurementUnitsPermission Permission = "update.measurement_units"
	// ArchiveValidMeasurementUnitsPermission is an account user permission.
	ArchiveValidMeasurementUnitsPermission Permission = "archive.measurement_units"

	// CreateValidIngredientStatesPermission is an account user permission.
	CreateValidIngredientStatesPermission Permission = "create.valid_ingredient_states"
	// ReadValidIngredientStatesPermission is an account user permission.
	ReadValidIngredientStatesPermission Permission = "read.valid_ingredient_states"
	// UpdateValidIngredientStatesPermission is an account user permission.
	UpdateValidIngredientStatesPermission Permission = "update.valid_ingredient_states"
	// ArchiveValidIngredientStatesPermission is an account user permission.
	ArchiveValidIngredientStatesPermission Permission = "archive.valid_ingredient_states"

	// CreateValidMeasurementUnitConversionsPermission is an account user permission.
	CreateValidMeasurementUnitConversionsPermission Permission = "create.measurement_conversions"
	// ReadValidMeasurementUnitConversionsPermission is an account user permission.
	ReadValidMeasurementUnitConversionsPermission Permission = "read.measurement_conversions"
	// UpdateValidMeasurementUnitConversionsPermission is an account user permission.
	UpdateValidMeasurementUnitConversionsPermission Permission = "update.measurement_conversions"
	// ArchiveValidMeasurementUnitConversionsPermission is an account user permission.
	ArchiveValidMeasurementUnitConversionsPermission Permission = "archive.measurement_conversions"

	// CreateValidIngredientPreparationsPermission is an account user permission.
	CreateValidIngredientPreparationsPermission Permission = "create.valid_ingredient_preparations"
	// ReadValidIngredientPreparationsPermission is an account user permission.
	ReadValidIngredientPreparationsPermission Permission = "read.valid_ingredient_preparations"
	// SearchValidIngredientPreparationsPermission is an account user permission.
	SearchValidIngredientPreparationsPermission Permission = "search.valid_ingredient_preparations"
	// UpdateValidIngredientPreparationsPermission is an account user permission.
	UpdateValidIngredientPreparationsPermission Permission = "update.valid_ingredient_preparations"
	// ArchiveValidIngredientPreparationsPermission is an account user permission.
	ArchiveValidIngredientPreparationsPermission Permission = "archive.valid_ingredient_preparations"

	// CreateValidIngredientStateIngredientsPermission is an account user permission.
	CreateValidIngredientStateIngredientsPermission Permission = "create.valid_ingredient_state_ingredients"
	// ReadValidIngredientStateIngredientsPermission is an account user permission.
	ReadValidIngredientStateIngredientsPermission Permission = "read.valid_ingredient_state_ingredients"
	// SearchValidIngredientStateIngredientsPermission is an account user permission.
	SearchValidIngredientStateIngredientsPermission Permission = "search.valid_ingredient_state_ingredients"
	// UpdateValidIngredientStateIngredientsPermission is an account user permission.
	UpdateValidIngredientStateIngredientsPermission Permission = "update.valid_ingredient_state_ingredients"
	// ArchiveValidIngredientStateIngredientsPermission is an account user permission.
	ArchiveValidIngredientStateIngredientsPermission Permission = "archive.valid_ingredient_state_ingredients"

	// CreateValidPreparationInstrumentsPermission is an account user permission.
	CreateValidPreparationInstrumentsPermission Permission = "create.valid_preparation_instruments"
	// ReadValidPreparationInstrumentsPermission is an account user permission.
	ReadValidPreparationInstrumentsPermission Permission = "read.valid_preparation_instruments"
	// SearchValidPreparationInstrumentsPermission is an account user permission.
	SearchValidPreparationInstrumentsPermission Permission = "search.valid_preparation_instruments"
	// UpdateValidPreparationInstrumentsPermission is an account user permission.
	UpdateValidPreparationInstrumentsPermission Permission = "update.valid_preparation_instruments"
	// ArchiveValidPreparationInstrumentsPermission is an account user permission.
	ArchiveValidPreparationInstrumentsPermission Permission = "archive.valid_preparation_instruments"

	// CreateValidPreparationVesselsPermission is an account user permission.
	CreateValidPreparationVesselsPermission Permission = "create.valid_preparation_vessels"
	// ReadValidPreparationVesselsPermission is an account user permission.
	ReadValidPreparationVesselsPermission Permission = "read.valid_preparation_vessels"
	// SearchValidPreparationVesselsPermission is an account user permission.
	SearchValidPreparationVesselsPermission Permission = "search.valid_preparation_vessels"
	// UpdateValidPreparationVesselsPermission is an account user permission.
	UpdateValidPreparationVesselsPermission Permission = "update.valid_preparation_vessels"
	// ArchiveValidPreparationVesselsPermission is an account user permission.
	ArchiveValidPreparationVesselsPermission Permission = "archive.valid_preparation_vessels"

	// CreateValidIngredientMeasurementUnitsPermission is an account user permission.
	CreateValidIngredientMeasurementUnitsPermission Permission = "create.valid_ingredient_measurement_units"
	// ReadValidIngredientMeasurementUnitsPermission is an account user permission.
	ReadValidIngredientMeasurementUnitsPermission Permission = "read.valid_ingredient_measurement_units"
	// SearchValidIngredientMeasurementUnitsPermission is an account user permission.
	SearchValidIngredientMeasurementUnitsPermission Permission = "search.valid_ingredient_measurement_units"
	// UpdateValidIngredientMeasurementUnitsPermission is an account user permission.
	UpdateValidIngredientMeasurementUnitsPermission Permission = "update.valid_ingredient_measurement_units"
	// ArchiveValidIngredientMeasurementUnitsPermission is an account user permission.
	ArchiveValidIngredientMeasurementUnitsPermission Permission = "archive.valid_ingredient_measurement_units"

	// CreateMealsPermission is an account user permission.
	CreateMealsPermission Permission = "create.meals"
	// ReadMealsPermission is an account user permission.
	ReadMealsPermission Permission = "read.meals"
	// UpdateMealsPermission is an account user permission.
	UpdateMealsPermission Permission = "update.meals"
	// ArchiveMealsPermission is an account user permission.
	ArchiveMealsPermission Permission = "archive.meals"

	// CreateRecipesPermission is an account user permission.
	CreateRecipesPermission Permission = "create.recipes"
	// ReadRecipesPermission is an account user permission.
	ReadRecipesPermission Permission = "read.recipes"
	// SearchRecipesPermission is an account user permission.
	SearchRecipesPermission Permission = "search.recipes"
	// UpdateRecipesPermission is an account user permission.
	UpdateRecipesPermission Permission = "update.recipes"
	// ArchiveRecipesPermission is an account user permission.
	ArchiveRecipesPermission Permission = "archive.recipes"

	// CreateRecipePrepTasksPermission is an account user permission.
	CreateRecipePrepTasksPermission Permission = "create.recipe_prep_tasks"
	// ReadRecipePrepTasksPermission is an account user permission.
	ReadRecipePrepTasksPermission Permission = "read.recipe_prep_tasks"
	// UpdateRecipePrepTasksPermission is an account user permission.
	UpdateRecipePrepTasksPermission Permission = "update.recipe_prep_tasks"
	// ArchiveRecipePrepTasksPermission is an account user permission.
	ArchiveRecipePrepTasksPermission Permission = "archive.recipe_prep_tasks"

	// CreateRecipeStepsPermission is an account user permission.
	CreateRecipeStepsPermission Permission = "create.recipe_steps"
	// ReadRecipeStepsPermission is an account user permission.
	ReadRecipeStepsPermission Permission = "read.recipe_steps"
	// SearchRecipeStepsPermission is an account user permission.
	SearchRecipeStepsPermission Permission = "search.recipe_steps"
	// UpdateRecipeStepsPermission is an account user permission.
	UpdateRecipeStepsPermission Permission = "update.recipe_steps"
	// ArchiveRecipeStepsPermission is an account user permission.
	ArchiveRecipeStepsPermission Permission = "archive.recipe_steps"

	// CreateRecipeStepInstrumentsPermission is an account user permission.
	CreateRecipeStepInstrumentsPermission Permission = "create.recipe_step_instruments"
	// ReadRecipeStepInstrumentsPermission is an account user permission.
	ReadRecipeStepInstrumentsPermission Permission = "read.recipe_step_instruments"
	// SearchRecipeStepInstrumentsPermission is an account user permission.
	SearchRecipeStepInstrumentsPermission Permission = "search.recipe_step_instruments"
	// UpdateRecipeStepInstrumentsPermission is an account user permission.
	UpdateRecipeStepInstrumentsPermission Permission = "update.recipe_step_instruments"
	// ArchiveRecipeStepInstrumentsPermission is an account user permission.
	ArchiveRecipeStepInstrumentsPermission Permission = "archive.recipe_step_instruments"

	// CreateRecipeStepVesselsPermission is an account user permission.
	CreateRecipeStepVesselsPermission Permission = "create.recipe_step_vessels"
	// ReadRecipeStepVesselsPermission is an account user permission.
	ReadRecipeStepVesselsPermission Permission = "read.recipe_step_vessels"
	// SearchRecipeStepVesselsPermission is an account user permission.
	SearchRecipeStepVesselsPermission Permission = "search.recipe_step_vessels"
	// UpdateRecipeStepVesselsPermission is an account user permission.
	UpdateRecipeStepVesselsPermission Permission = "update.recipe_step_vessels"
	// ArchiveRecipeStepVesselsPermission is an account user permission.
	ArchiveRecipeStepVesselsPermission Permission = "archive.recipe_step_vessels"

	// CreateRecipeStepIngredientsPermission is an account user permission.
	CreateRecipeStepIngredientsPermission Permission = "create.recipe_step_ingredients"
	// ReadRecipeStepIngredientsPermission is an account user permission.
	ReadRecipeStepIngredientsPermission Permission = "read.recipe_step_ingredients"
	// SearchRecipeStepIngredientsPermission is an account user permission.
	SearchRecipeStepIngredientsPermission Permission = "search.recipe_step_ingredients"
	// UpdateRecipeStepIngredientsPermission is an account user permission.
	UpdateRecipeStepIngredientsPermission Permission = "update.recipe_step_ingredients"
	// ArchiveRecipeStepIngredientsPermission is an account user permission.
	ArchiveRecipeStepIngredientsPermission Permission = "archive.recipe_step_ingredients"

	// CreateRecipeStepCompletionConditionsPermission is an account user permission.
	CreateRecipeStepCompletionConditionsPermission Permission = "create.recipe_step_completion_conditions"
	// ReadRecipeStepCompletionConditionsPermission is an account user permission.
	ReadRecipeStepCompletionConditionsPermission Permission = "read.recipe_step_completion_conditions"
	// SearchRecipeStepCompletionConditionsPermission is an account user permission.
	SearchRecipeStepCompletionConditionsPermission Permission = "search.recipe_step_completion_conditions"
	// UpdateRecipeStepCompletionConditionsPermission is an account user permission.
	UpdateRecipeStepCompletionConditionsPermission Permission = "update.recipe_step_completion_conditions"
	// ArchiveRecipeStepCompletionConditionsPermission is an account user permission.
	ArchiveRecipeStepCompletionConditionsPermission Permission = "archive.recipe_step_completion_conditions"

	// CreateRecipeStepProductsPermission is an account user permission.
	CreateRecipeStepProductsPermission Permission = "create.recipe_step_products"
	// ReadRecipeStepProductsPermission is an account user permission.
	ReadRecipeStepProductsPermission Permission = "read.recipe_step_products"
	// SearchRecipeStepProductsPermission is an account user permission.
	SearchRecipeStepProductsPermission Permission = "search.recipe_step_products"
	// UpdateRecipeStepProductsPermission is an account user permission.
	UpdateRecipeStepProductsPermission Permission = "update.recipe_step_products"
	// ArchiveRecipeStepProductsPermission is an account user permission.
	ArchiveRecipeStepProductsPermission Permission = "archive.recipe_step_products"

	// CreateMealPlansPermission is an account user permission.
	CreateMealPlansPermission Permission = "create.meal_plans"
	// ReadMealPlansPermission is an account user permission.
	ReadMealPlansPermission Permission = "read.meal_plans"
	// SearchMealPlansPermission is an account user permission.
	SearchMealPlansPermission Permission = "search.meal_plans"
	// UpdateMealPlansPermission is an account user permission.
	UpdateMealPlansPermission Permission = "update.meal_plans"
	// ArchiveMealPlansPermission is an account user permission.
	ArchiveMealPlansPermission Permission = "archive.meal_plans"

	// CreateMealPlanEventsPermission is an account user permission.
	CreateMealPlanEventsPermission Permission = "create.meal_plan_events"
	// ReadMealPlanEventsPermission is an account user permission.
	ReadMealPlanEventsPermission Permission = "read.meal_plan_events"
	// UpdateMealPlanEventsPermission is an account user permission.
	UpdateMealPlanEventsPermission Permission = "update.meal_plan_events"
	// ArchiveMealPlanEventsPermission is an account user permission.
	ArchiveMealPlanEventsPermission Permission = "archive.meal_plan_events"

	// CreateMealPlanOptionsPermission is an account user permission.
	CreateMealPlanOptionsPermission Permission = "create.meal_plan_options"
	// ReadMealPlanOptionsPermission is an account user permission.
	ReadMealPlanOptionsPermission Permission = "read.meal_plan_options"
	// SearchMealPlanOptionsPermission is an account user permission.
	SearchMealPlanOptionsPermission Permission = "search.meal_plan_options"
	// UpdateMealPlanOptionsPermission is an account user permission.
	UpdateMealPlanOptionsPermission Permission = "update.meal_plan_options"
	// ArchiveMealPlanOptionsPermission is an account user permission.
	ArchiveMealPlanOptionsPermission Permission = "archive.meal_plan_options"

	// CreateMealPlanGroceryListItemsPermission is an account user permission.
	CreateMealPlanGroceryListItemsPermission Permission = "create.meal_plan_grocery_list_items"
	// ReadMealPlanGroceryListItemsPermission is an account user permission.
	ReadMealPlanGroceryListItemsPermission Permission = "read.meal_plan_grocery_list_items"
	// UpdateMealPlanGroceryListItemsPermission is an account user permission.
	UpdateMealPlanGroceryListItemsPermission Permission = "update.meal_plan_grocery_list_items"
	// ArchiveMealPlanGroceryListItemsPermission is an account user permission.
	ArchiveMealPlanGroceryListItemsPermission Permission = "archive.meal_plan_grocery_list_items"

	// CreateMealPlanOptionVotesPermission is an account user permission.
	CreateMealPlanOptionVotesPermission Permission = "create.meal_plan_option_votes"
	// ReadMealPlanOptionVotesPermission is an account user permission.
	ReadMealPlanOptionVotesPermission Permission = "read.meal_plan_option_votes"
	// SearchMealPlanOptionVotesPermission is an account user permission.
	SearchMealPlanOptionVotesPermission Permission = "search.meal_plan_option_votes"
	// UpdateMealPlanOptionVotesPermission is an account user permission.
	UpdateMealPlanOptionVotesPermission Permission = "update.meal_plan_option_votes"
	// ArchiveMealPlanOptionVotesPermission is an account user permission.
	ArchiveMealPlanOptionVotesPermission Permission = "archive.meal_plan_option_votes"

	// ReadMealPlanTasksPermission is an account user permission.
	ReadMealPlanTasksPermission Permission = "read.meal_plan_tasks"
	// CreateMealPlanTasksPermission is an account user permission.
	CreateMealPlanTasksPermission Permission = "create.meal_plan_tasks"
	// UpdateMealPlanTasksPermission is an account user permission.
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

	// CreateUserIngredientPreferencesPermission is an account user permission.
	CreateUserIngredientPreferencesPermission Permission = "create.user_ingredient_preferences"
	// ReadUserIngredientPreferencesPermission is an account user permission.
	ReadUserIngredientPreferencesPermission Permission = "read.user_ingredient_preferences"
	// UpdateUserIngredientPreferencesPermission is an account user permission.
	UpdateUserIngredientPreferencesPermission Permission = "update.user_ingredient_preferences"
	// ArchiveUserIngredientPreferencesPermission is an account user permission.
	ArchiveUserIngredientPreferencesPermission Permission = "archive.user_ingredient_preferences"

	// CreateAccountInstrumentOwnershipsPermission is an account user permission.
	CreateAccountInstrumentOwnershipsPermission Permission = "create.account_instrument_ownerships"
	// ReadAccountInstrumentOwnershipsPermission is an account user permission.
	ReadAccountInstrumentOwnershipsPermission Permission = "read.account_instrument_ownerships"
	// UpdateAccountInstrumentOwnershipsPermission is an account user permission.
	UpdateAccountInstrumentOwnershipsPermission Permission = "update.account_instrument_ownerships"
	// ArchiveAccountInstrumentOwnershipsPermission is an account user permission.
	ArchiveAccountInstrumentOwnershipsPermission Permission = "archive.account_instrument_ownerships"

	// CreateRecipeRatingsPermission is an account user permission.
	CreateRecipeRatingsPermission Permission = "create.recipe_ratings"
	// ReadRecipeRatingsPermission is an account user permission.
	ReadRecipeRatingsPermission Permission = "read.recipe_ratings"
	// UpdateRecipeRatingsPermission is an account user permission.
	UpdateRecipeRatingsPermission Permission = "update.recipe_ratings"
	// ArchiveRecipeRatingsPermission is an account user permission.
	ArchiveRecipeRatingsPermission Permission = "archive.recipe_ratings"

	// CreateOAuth2ClientsPermission is an account admin permission.
	CreateOAuth2ClientsPermission Permission = "create.oauth2_clients"
	// ReadOAuth2ClientsPermission is an account admin permission.
	ReadOAuth2ClientsPermission Permission = "read.oauth2_clients"
	// ArchiveOAuth2ClientsPermission is an account admin permission.
	ArchiveOAuth2ClientsPermission Permission = "archive.oauth2_clients"

	// CreateUserNotificationsPermission is an admin user permission.
	CreateUserNotificationsPermission Permission = "create.user_notifications"
	// ReadUserNotificationsPermission is an account user permission.
	ReadUserNotificationsPermission Permission = "read.user_notifications"
	// UpdateUserNotificationsPermission is an account user permission.
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
