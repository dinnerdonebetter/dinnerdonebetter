package authorization

import (
	"gopkg.in/mikespook/gorbac.v2"
)

var _ gorbac.Permission = (*Permission)(nil)

type (
	// Permission is a simple string alias.
	Permission string
)

const (
	// CycleCookieSecretPermission is a service admin permission.
	CycleCookieSecretPermission Permission = "update.cookie_secret"
	// UpdateUserStatusPermission is a service admin permission.
	UpdateUserStatusPermission Permission = "update.user_status"
	// ReadUserPermission is a service admin permission.
	ReadUserPermission Permission = "read.user"
	// SearchUserPermission is a service admin permission.
	SearchUserPermission Permission = "search.user"

	// UpdateHouseholdPermission is a household admin permission.
	UpdateHouseholdPermission Permission = "update.household"
	// ArchiveHouseholdPermission is a household admin permission.
	ArchiveHouseholdPermission Permission = "archive.household"
	// InviteUserToHouseholdPermission is a household admin permission.
	InviteUserToHouseholdPermission Permission = "household.add.member"
	// ModifyMemberPermissionsForHouseholdPermission is a household admin permission.
	ModifyMemberPermissionsForHouseholdPermission Permission = "household.membership.modify"
	// RemoveMemberHouseholdPermission is a household admin permission.
	RemoveMemberHouseholdPermission Permission = "remove_member.household"
	// TransferHouseholdPermission is a household admin permission.
	TransferHouseholdPermission Permission = "transfer.household"
	// CreateWebhooksPermission is a household admin permission.
	CreateWebhooksPermission Permission = "create.webhooks"
	// ReadWebhooksPermission is a household admin permission.
	ReadWebhooksPermission Permission = "read.webhooks"
	// UpdateWebhooksPermission is a household admin permission.
	UpdateWebhooksPermission Permission = "update.webhooks"
	// ArchiveWebhooksPermission is a household admin permission.
	ArchiveWebhooksPermission Permission = "archive.webhooks"
	// CreateAPIClientsPermission is a household admin permission.
	CreateAPIClientsPermission Permission = "create.api_clients"
	// ReadAPIClientsPermission is a household admin permission.
	ReadAPIClientsPermission Permission = "read.api_clients"
	// ArchiveAPIClientsPermission is a household admin permission.
	ArchiveAPIClientsPermission Permission = "archive.api_clients"

	// CreateValidInstrumentsPermission is a household user permission.
	CreateValidInstrumentsPermission Permission = "create.valid_instruments"
	// ReadValidInstrumentsPermission is a household user permission.
	ReadValidInstrumentsPermission Permission = "read.valid_instruments"
	// SearchValidInstrumentsPermission is a household user permission.
	SearchValidInstrumentsPermission Permission = "search.valid_instruments"
	// UpdateValidInstrumentsPermission is a household user permission.
	UpdateValidInstrumentsPermission Permission = "update.valid_instruments"
	// ArchiveValidInstrumentsPermission is a household user permission.
	ArchiveValidInstrumentsPermission Permission = "archive.valid_instruments"

	// CreateValidIngredientsPermission is a household user permission.
	CreateValidIngredientsPermission Permission = "create.valid_ingredients"
	// ReadValidIngredientsPermission is a household user permission.
	ReadValidIngredientsPermission Permission = "read.valid_ingredients"
	// SearchValidIngredientsPermission is a household user permission.
	SearchValidIngredientsPermission Permission = "search.valid_ingredients"
	// UpdateValidIngredientsPermission is a household user permission.
	UpdateValidIngredientsPermission Permission = "update.valid_ingredients"
	// ArchiveValidIngredientsPermission is a household user permission.
	ArchiveValidIngredientsPermission Permission = "archive.valid_ingredients"

	// CreateValidIngredientGroupsPermission is a household user permission.
	CreateValidIngredientGroupsPermission Permission = "create.valid_ingredient_groups"
	// ReadValidIngredientGroupsPermission is a household user permission.
	ReadValidIngredientGroupsPermission Permission = "read.valid_ingredient_groups"
	// SearchValidIngredientGroupsPermission is a household user permission.
	SearchValidIngredientGroupsPermission Permission = "search.valid_ingredient_groups"
	// UpdateValidIngredientGroupsPermission is a household user permission.
	UpdateValidIngredientGroupsPermission Permission = "update.valid_ingredient_groups"
	// ArchiveValidIngredientGroupsPermission is a household user permission.
	ArchiveValidIngredientGroupsPermission Permission = "archive.valid_ingredient_groups"

	// CreateValidPreparationsPermission is a household user permission.
	CreateValidPreparationsPermission Permission = "create.valid_preparations"
	// ReadValidPreparationsPermission is a household user permission.
	ReadValidPreparationsPermission Permission = "read.valid_preparations"
	// SearchValidPreparationsPermission is a household user permission.
	SearchValidPreparationsPermission Permission = "search.valid_preparations"
	// UpdateValidPreparationsPermission is a household user permission.
	UpdateValidPreparationsPermission Permission = "update.valid_preparations"
	// ArchiveValidPreparationsPermission is a household user permission.
	ArchiveValidPreparationsPermission Permission = "archive.valid_preparations"

	// CreateValidMeasurementUnitsPermission is a household user permission.
	CreateValidMeasurementUnitsPermission Permission = "create.measurement_units"
	// ReadValidMeasurementUnitsPermission is a household user permission.
	ReadValidMeasurementUnitsPermission Permission = "read.measurement_units"
	// SearchValidMeasurementUnitsPermission is a household user permission.
	SearchValidMeasurementUnitsPermission Permission = "search.measurement_units"
	// UpdateValidMeasurementUnitsPermission is a household user permission.
	UpdateValidMeasurementUnitsPermission Permission = "update.measurement_units"
	// ArchiveValidMeasurementUnitsPermission is a household user permission.
	ArchiveValidMeasurementUnitsPermission Permission = "archive.measurement_units"

	// CreateValidIngredientStatesPermission is a household user permission.
	CreateValidIngredientStatesPermission Permission = "create.valid_ingredient_states"
	// ReadValidIngredientStatesPermission is a household user permission.
	ReadValidIngredientStatesPermission Permission = "read.valid_ingredient_states"
	// UpdateValidIngredientStatesPermission is a household user permission.
	UpdateValidIngredientStatesPermission Permission = "update.valid_ingredient_states"
	// ArchiveValidIngredientStatesPermission is a household user permission.
	ArchiveValidIngredientStatesPermission Permission = "archive.valid_ingredient_states"

	// CreateValidMeasurementConversionsPermission is a household user permission.
	CreateValidMeasurementConversionsPermission Permission = "create.measurement_conversions"
	// ReadValidMeasurementConversionsPermission is a household user permission.
	ReadValidMeasurementConversionsPermission Permission = "read.measurement_conversions"
	// UpdateValidMeasurementConversionsPermission is a household user permission.
	UpdateValidMeasurementConversionsPermission Permission = "update.measurement_conversions"
	// ArchiveValidMeasurementConversionsPermission is a household user permission.
	ArchiveValidMeasurementConversionsPermission Permission = "archive.measurement_conversions"

	// CreateValidIngredientPreparationsPermission is a household user permission.
	CreateValidIngredientPreparationsPermission Permission = "create.valid_ingredient_preparations"
	// ReadValidIngredientPreparationsPermission is a household user permission.
	ReadValidIngredientPreparationsPermission Permission = "read.valid_ingredient_preparations"
	// SearchValidIngredientPreparationsPermission is a household user permission.
	SearchValidIngredientPreparationsPermission Permission = "search.valid_ingredient_preparations"
	// UpdateValidIngredientPreparationsPermission is a household user permission.
	UpdateValidIngredientPreparationsPermission Permission = "update.valid_ingredient_preparations"
	// ArchiveValidIngredientPreparationsPermission is a household user permission.
	ArchiveValidIngredientPreparationsPermission Permission = "archive.valid_ingredient_preparations"

	// CreateValidIngredientStateIngredientsPermission is a household user permission.
	CreateValidIngredientStateIngredientsPermission Permission = "create.valid_ingredient_state_ingredients"
	// ReadValidIngredientStateIngredientsPermission is a household user permission.
	ReadValidIngredientStateIngredientsPermission Permission = "read.valid_ingredient_state_ingredients"
	// SearchValidIngredientStateIngredientsPermission is a household user permission.
	SearchValidIngredientStateIngredientsPermission Permission = "search.valid_ingredient_state_ingredients"
	// UpdateValidIngredientStateIngredientsPermission is a household user permission.
	UpdateValidIngredientStateIngredientsPermission Permission = "update.valid_ingredient_state_ingredients"
	// ArchiveValidIngredientStateIngredientsPermission is a household user permission.
	ArchiveValidIngredientStateIngredientsPermission Permission = "archive.valid_ingredient_state_ingredients"

	// CreateValidPreparationInstrumentsPermission is a household user permission.
	CreateValidPreparationInstrumentsPermission Permission = "create.valid_preparation_instruments"
	// ReadValidPreparationInstrumentsPermission is a household user permission.
	ReadValidPreparationInstrumentsPermission Permission = "read.valid_preparation_instruments"
	// SearchValidPreparationInstrumentsPermission is a household user permission.
	SearchValidPreparationInstrumentsPermission Permission = "search.valid_preparation_instruments"
	// UpdateValidPreparationInstrumentsPermission is a household user permission.
	UpdateValidPreparationInstrumentsPermission Permission = "update.valid_preparation_instruments"
	// ArchiveValidPreparationInstrumentsPermission is a household user permission.
	ArchiveValidPreparationInstrumentsPermission Permission = "archive.valid_preparation_instruments"

	// CreateValidIngredientMeasurementUnitsPermission is a household user permission.
	CreateValidIngredientMeasurementUnitsPermission Permission = "create.valid_ingredient_measurement_units"
	// ReadValidIngredientMeasurementUnitsPermission is a household user permission.
	ReadValidIngredientMeasurementUnitsPermission Permission = "read.valid_ingredient_measurement_units"
	// SearchValidIngredientMeasurementUnitsPermission is a household user permission.
	SearchValidIngredientMeasurementUnitsPermission Permission = "search.valid_ingredient_measurement_units"
	// UpdateValidIngredientMeasurementUnitsPermission is a household user permission.
	UpdateValidIngredientMeasurementUnitsPermission Permission = "update.valid_ingredient_measurement_units"
	// ArchiveValidIngredientMeasurementUnitsPermission is a household user permission.
	ArchiveValidIngredientMeasurementUnitsPermission Permission = "archive.valid_ingredient_measurement_units"

	// CreateMealsPermission is a household user permission.
	CreateMealsPermission Permission = "create.meals"
	// ReadMealsPermission is a household user permission.
	ReadMealsPermission Permission = "read.meals"
	// UpdateMealsPermission is a household user permission.
	UpdateMealsPermission Permission = "update.meals"
	// ArchiveMealsPermission is a household user permission.
	ArchiveMealsPermission Permission = "archive.meals"

	// CreateRecipesPermission is a household user permission.
	CreateRecipesPermission Permission = "create.recipes"
	// ReadRecipesPermission is a household user permission.
	ReadRecipesPermission Permission = "read.recipes"
	// SearchRecipesPermission is a household user permission.
	SearchRecipesPermission Permission = "search.recipes"
	// UpdateRecipesPermission is a household user permission.
	UpdateRecipesPermission Permission = "update.recipes"
	// ArchiveRecipesPermission is a household user permission.
	ArchiveRecipesPermission Permission = "archive.recipes"

	// CreateRecipePrepTasksPermission is a household user permission.
	CreateRecipePrepTasksPermission Permission = "create.recipe_prep_tasks"
	// ReadRecipePrepTasksPermission is a household user permission.
	ReadRecipePrepTasksPermission Permission = "read.recipe_prep_tasks"
	// UpdateRecipePrepTasksPermission is a household user permission.
	UpdateRecipePrepTasksPermission Permission = "update.recipe_prep_tasks"
	// ArchiveRecipePrepTasksPermission is a household user permission.
	ArchiveRecipePrepTasksPermission Permission = "archive.recipe_prep_tasks"

	// CreateRecipeStepsPermission is a household user permission.
	CreateRecipeStepsPermission Permission = "create.recipe_steps"
	// ReadRecipeStepsPermission is a household user permission.
	ReadRecipeStepsPermission Permission = "read.recipe_steps"
	// SearchRecipeStepsPermission is a household user permission.
	SearchRecipeStepsPermission Permission = "search.recipe_steps"
	// UpdateRecipeStepsPermission is a household user permission.
	UpdateRecipeStepsPermission Permission = "update.recipe_steps"
	// ArchiveRecipeStepsPermission is a household user permission.
	ArchiveRecipeStepsPermission Permission = "archive.recipe_steps"

	// CreateRecipeStepInstrumentsPermission is a household user permission.
	CreateRecipeStepInstrumentsPermission Permission = "create.recipe_step_instruments"
	// ReadRecipeStepInstrumentsPermission is a household user permission.
	ReadRecipeStepInstrumentsPermission Permission = "read.recipe_step_instruments"
	// SearchRecipeStepInstrumentsPermission is a household user permission.
	SearchRecipeStepInstrumentsPermission Permission = "search.recipe_step_instruments"
	// UpdateRecipeStepInstrumentsPermission is a household user permission.
	UpdateRecipeStepInstrumentsPermission Permission = "update.recipe_step_instruments"
	// ArchiveRecipeStepInstrumentsPermission is a household user permission.
	ArchiveRecipeStepInstrumentsPermission Permission = "archive.recipe_step_instruments"

	// CreateRecipeStepVesselsPermission is a household user permission.
	CreateRecipeStepVesselsPermission Permission = "create.recipe_step_vessels"
	// ReadRecipeStepVesselsPermission is a household user permission.
	ReadRecipeStepVesselsPermission Permission = "read.recipe_step_vessels"
	// SearchRecipeStepVesselsPermission is a household user permission.
	SearchRecipeStepVesselsPermission Permission = "search.recipe_step_vessels"
	// UpdateRecipeStepVesselsPermission is a household user permission.
	UpdateRecipeStepVesselsPermission Permission = "update.recipe_step_vessels"
	// ArchiveRecipeStepVesselsPermission is a household user permission.
	ArchiveRecipeStepVesselsPermission Permission = "archive.recipe_step_vessels"

	// CreateRecipeStepIngredientsPermission is a household user permission.
	CreateRecipeStepIngredientsPermission Permission = "create.recipe_step_ingredients"
	// ReadRecipeStepIngredientsPermission is a household user permission.
	ReadRecipeStepIngredientsPermission Permission = "read.recipe_step_ingredients"
	// SearchRecipeStepIngredientsPermission is a household user permission.
	SearchRecipeStepIngredientsPermission Permission = "search.recipe_step_ingredients"
	// UpdateRecipeStepIngredientsPermission is a household user permission.
	UpdateRecipeStepIngredientsPermission Permission = "update.recipe_step_ingredients"
	// ArchiveRecipeStepIngredientsPermission is a household user permission.
	ArchiveRecipeStepIngredientsPermission Permission = "archive.recipe_step_ingredients"

	// CreateRecipeStepCompletionConditionsPermission is a household user permission.
	CreateRecipeStepCompletionConditionsPermission Permission = "create.recipe_step_completion_conditions"
	// ReadRecipeStepCompletionConditionsPermission is a household user permission.
	ReadRecipeStepCompletionConditionsPermission Permission = "read.recipe_step_completion_conditions"
	// SearchRecipeStepCompletionConditionsPermission is a household user permission.
	SearchRecipeStepCompletionConditionsPermission Permission = "search.recipe_step_completion_conditions"
	// UpdateRecipeStepCompletionConditionsPermission is a household user permission.
	UpdateRecipeStepCompletionConditionsPermission Permission = "update.recipe_step_completion_conditions"
	// ArchiveRecipeStepCompletionConditionsPermission is a household user permission.
	ArchiveRecipeStepCompletionConditionsPermission Permission = "archive.recipe_step_completion_conditions"

	// CreateRecipeStepProductsPermission is a household user permission.
	CreateRecipeStepProductsPermission Permission = "create.recipe_step_products"
	// ReadRecipeStepProductsPermission is a household user permission.
	ReadRecipeStepProductsPermission Permission = "read.recipe_step_products"
	// SearchRecipeStepProductsPermission is a household user permission.
	SearchRecipeStepProductsPermission Permission = "search.recipe_step_products"
	// UpdateRecipeStepProductsPermission is a household user permission.
	UpdateRecipeStepProductsPermission Permission = "update.recipe_step_products"
	// ArchiveRecipeStepProductsPermission is a household user permission.
	ArchiveRecipeStepProductsPermission Permission = "archive.recipe_step_products"

	// CreateMealPlansPermission is a household user permission.
	CreateMealPlansPermission Permission = "create.meal_plans"
	// ReadMealPlansPermission is a household user permission.
	ReadMealPlansPermission Permission = "read.meal_plans"
	// SearchMealPlansPermission is a household user permission.
	SearchMealPlansPermission Permission = "search.meal_plans"
	// UpdateMealPlansPermission is a household user permission.
	UpdateMealPlansPermission Permission = "update.meal_plans"
	// ArchiveMealPlansPermission is a household user permission.
	ArchiveMealPlansPermission Permission = "archive.meal_plans"

	// CreateMealPlanEventsPermission is a household user permission.
	CreateMealPlanEventsPermission Permission = "create.meal_plan_events"
	// ReadMealPlanEventsPermission is a household user permission.
	ReadMealPlanEventsPermission Permission = "read.meal_plan_events"
	// UpdateMealPlanEventsPermission is a household user permission.
	UpdateMealPlanEventsPermission Permission = "update.meal_plan_events"
	// ArchiveMealPlanEventsPermission is a household user permission.
	ArchiveMealPlanEventsPermission Permission = "archive.meal_plan_events"

	// CreateMealPlanOptionsPermission is a household user permission.
	CreateMealPlanOptionsPermission Permission = "create.meal_plan_options"
	// ReadMealPlanOptionsPermission is a household user permission.
	ReadMealPlanOptionsPermission Permission = "read.meal_plan_options"
	// SearchMealPlanOptionsPermission is a household user permission.
	SearchMealPlanOptionsPermission Permission = "search.meal_plan_options"
	// UpdateMealPlanOptionsPermission is a household user permission.
	UpdateMealPlanOptionsPermission Permission = "update.meal_plan_options"
	// ArchiveMealPlanOptionsPermission is a household user permission.
	ArchiveMealPlanOptionsPermission Permission = "archive.meal_plan_options"

	// CreateMealPlanGroceryListItemsPermission is a household user permission.
	CreateMealPlanGroceryListItemsPermission Permission = "create.meal_plan_grocery_list_items"
	// ReadMealPlanGroceryListItemsPermission is a household user permission.
	ReadMealPlanGroceryListItemsPermission Permission = "read.meal_plan_grocery_list_items"
	// UpdateMealPlanGroceryListItemsPermission is a household user permission.
	UpdateMealPlanGroceryListItemsPermission Permission = "update.meal_plan_grocery_list_items"
	// ArchiveMealPlanGroceryListItemsPermission is a household user permission.
	ArchiveMealPlanGroceryListItemsPermission Permission = "archive.meal_plan_grocery_list_items"

	// CreateMealPlanOptionVotesPermission is a household user permission.
	CreateMealPlanOptionVotesPermission Permission = "create.meal_plan_option_votes"
	// ReadMealPlanOptionVotesPermission is a household user permission.
	ReadMealPlanOptionVotesPermission Permission = "read.meal_plan_option_votes"
	// SearchMealPlanOptionVotesPermission is a household user permission.
	SearchMealPlanOptionVotesPermission Permission = "search.meal_plan_option_votes"
	// UpdateMealPlanOptionVotesPermission is a household user permission.
	UpdateMealPlanOptionVotesPermission Permission = "update.meal_plan_option_votes"
	// ArchiveMealPlanOptionVotesPermission is a household user permission.
	ArchiveMealPlanOptionVotesPermission Permission = "archive.meal_plan_option_votes"

	// ReadMealPlanTasksPermission is a household user permission.
	ReadMealPlanTasksPermission Permission = "read.meal_plan_tasks"
	// CreateMealPlanTasksPermission is a household user permission.
	CreateMealPlanTasksPermission Permission = "create.meal_plan_tasks"
	// UpdateMealPlanTasksPermission is a household user permission.
	UpdateMealPlanTasksPermission Permission = "update.meal_plan_tasks"

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

	// CreateUserIngredientPreferencesPermission is a household user permission.
	CreateUserIngredientPreferencesPermission Permission = "create.user_ingredient_preferences"
	// ReadUserIngredientPreferencesPermission is a household user permission.
	ReadUserIngredientPreferencesPermission Permission = "read.user_ingredient_preferences"
	// UpdateUserIngredientPreferencesPermission is a household user permission.
	UpdateUserIngredientPreferencesPermission Permission = "update.user_ingredient_preferences"
	// ArchiveUserIngredientPreferencesPermission is a household user permission.
	ArchiveUserIngredientPreferencesPermission Permission = "archive.user_ingredient_preferences"

	// CreateHouseholdInstrumentOwnershipsPermission is a household user permission.
	CreateHouseholdInstrumentOwnershipsPermission Permission = "create.household_instrument_ownerships"
	// ReadHouseholdInstrumentOwnershipsPermission is a household user permission.
	ReadHouseholdInstrumentOwnershipsPermission Permission = "read.household_instrument_ownerships"
	// UpdateHouseholdInstrumentOwnershipsPermission is a household user permission.
	UpdateHouseholdInstrumentOwnershipsPermission Permission = "update.household_instrument_ownerships"
	// ArchiveHouseholdInstrumentOwnershipsPermission is a household user permission.
	ArchiveHouseholdInstrumentOwnershipsPermission Permission = "archive.household_instrument_ownerships"

	// CreateRecipeRatingsPermission is a household user permission.
	CreateRecipeRatingsPermission Permission = "create.recipe_ratings"
	// ReadRecipeRatingsPermission is a household user permission.
	ReadRecipeRatingsPermission Permission = "read.recipe_ratings"
	// UpdateRecipeRatingsPermission is a household user permission.
	UpdateRecipeRatingsPermission Permission = "update.recipe_ratings"
	// ArchiveRecipeRatingsPermission is a household user permission.
	ArchiveRecipeRatingsPermission Permission = "archive.recipe_ratings"

	// CreateOAuth2ClientsPermission is a household admin permission.
	CreateOAuth2ClientsPermission Permission = "create.oauth2_clients"
	// ReadOAuth2ClientsPermission is a household admin permission.
	ReadOAuth2ClientsPermission Permission = "read.oauth2_clients"
	// ArchiveOAuth2ClientsPermission is a household admin permission.
	ArchiveOAuth2ClientsPermission Permission = "archive.oauth2_clients"
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
	// service admin permissions.
	serviceAdminPermissions = map[string]gorbac.Permission{
		CycleCookieSecretPermission.ID(): CycleCookieSecretPermission,
		UpdateUserStatusPermission.ID():  UpdateUserStatusPermission,
		ReadUserPermission.ID():          ReadUserPermission,
		SearchUserPermission.ID():        SearchUserPermission,

		CreateOAuth2ClientsPermission.ID():  CreateOAuth2ClientsPermission,
		ArchiveOAuth2ClientsPermission.ID(): ArchiveOAuth2ClientsPermission,

		ArchiveServiceSettingsPermission.ID(): ArchiveServiceSettingsPermission,

		CreateRecipesPermission.ID(): CreateRecipesPermission,

		CreateValidInstrumentsPermission.ID():  CreateValidInstrumentsPermission,
		UpdateValidInstrumentsPermission.ID():  UpdateValidInstrumentsPermission,
		ArchiveValidInstrumentsPermission.ID(): ArchiveValidInstrumentsPermission,

		CreateValidIngredientsPermission.ID():  CreateValidIngredientsPermission,
		UpdateValidIngredientsPermission.ID():  UpdateValidIngredientsPermission,
		ArchiveValidIngredientsPermission.ID(): ArchiveValidIngredientsPermission,

		CreateValidIngredientGroupsPermission.ID():  CreateValidIngredientGroupsPermission,
		UpdateValidIngredientGroupsPermission.ID():  UpdateValidIngredientGroupsPermission,
		ArchiveValidIngredientGroupsPermission.ID(): ArchiveValidIngredientGroupsPermission,

		CreateValidPreparationsPermission.ID():  CreateValidPreparationsPermission,
		UpdateValidPreparationsPermission.ID():  UpdateValidPreparationsPermission,
		ArchiveValidPreparationsPermission.ID(): ArchiveValidPreparationsPermission,

		CreateValidMeasurementUnitsPermission.ID():  CreateValidMeasurementUnitsPermission,
		UpdateValidMeasurementUnitsPermission.ID():  UpdateValidMeasurementUnitsPermission,
		ArchiveValidMeasurementUnitsPermission.ID(): ArchiveValidMeasurementUnitsPermission,

		CreateValidMeasurementConversionsPermission.ID():  CreateValidMeasurementConversionsPermission,
		UpdateValidMeasurementConversionsPermission.ID():  UpdateValidMeasurementConversionsPermission,
		ArchiveValidMeasurementConversionsPermission.ID(): ArchiveValidMeasurementConversionsPermission,

		CreateValidIngredientPreparationsPermission.ID():  CreateValidIngredientPreparationsPermission,
		UpdateValidIngredientPreparationsPermission.ID():  UpdateValidIngredientPreparationsPermission,
		ArchiveValidIngredientPreparationsPermission.ID(): ArchiveValidIngredientPreparationsPermission,

		CreateValidIngredientStateIngredientsPermission.ID():  CreateValidIngredientStateIngredientsPermission,
		UpdateValidIngredientStateIngredientsPermission.ID():  UpdateValidIngredientStateIngredientsPermission,
		ArchiveValidIngredientStateIngredientsPermission.ID(): ArchiveValidIngredientStateIngredientsPermission,

		CreateValidPreparationInstrumentsPermission.ID():  CreateValidPreparationInstrumentsPermission,
		UpdateValidPreparationInstrumentsPermission.ID():  UpdateValidPreparationInstrumentsPermission,
		ArchiveValidPreparationInstrumentsPermission.ID(): ArchiveValidPreparationInstrumentsPermission,

		CreateValidIngredientMeasurementUnitsPermission.ID():  CreateValidIngredientMeasurementUnitsPermission,
		UpdateValidIngredientMeasurementUnitsPermission.ID():  UpdateValidIngredientMeasurementUnitsPermission,
		ArchiveValidIngredientMeasurementUnitsPermission.ID(): ArchiveValidIngredientMeasurementUnitsPermission,

		CreateValidIngredientStatesPermission.ID():  CreateValidIngredientStatesPermission,
		UpdateValidIngredientStatesPermission.ID():  UpdateValidIngredientStatesPermission,
		ArchiveValidIngredientStatesPermission.ID(): ArchiveValidIngredientStatesPermission,

		// only admins can arbitrarily create these via the API, and this is basically for integration test purposes.
		CreateMealPlanTasksPermission.ID():            CreateMealPlanTasksPermission,
		CreateMealPlanGroceryListItemsPermission.ID(): CreateMealPlanGroceryListItemsPermission,
	}

	// household admin permissions.
	householdAdminPermissions = map[string]gorbac.Permission{
		UpdateHouseholdPermission.ID():   UpdateHouseholdPermission,
		ArchiveHouseholdPermission.ID():  ArchiveHouseholdPermission,
		TransferHouseholdPermission.ID(): TransferHouseholdPermission,

		InviteUserToHouseholdPermission.ID():               InviteUserToHouseholdPermission,
		ModifyMemberPermissionsForHouseholdPermission.ID(): ModifyMemberPermissionsForHouseholdPermission,
		RemoveMemberHouseholdPermission.ID():               RemoveMemberHouseholdPermission,

		CreateWebhooksPermission.ID():  CreateWebhooksPermission,
		UpdateWebhooksPermission.ID():  UpdateWebhooksPermission,
		ArchiveWebhooksPermission.ID(): ArchiveWebhooksPermission,

		CreateMealPlansPermission.ID():  CreateMealPlansPermission,
		UpdateMealPlansPermission.ID():  UpdateMealPlansPermission,
		ArchiveMealPlansPermission.ID(): ArchiveMealPlansPermission,

		CreateMealPlanEventsPermission.ID():  CreateMealPlanEventsPermission,
		UpdateMealPlanEventsPermission.ID():  UpdateMealPlanEventsPermission,
		ArchiveMealPlanEventsPermission.ID(): ArchiveMealPlanEventsPermission,

		CreateMealPlanOptionsPermission.ID():  CreateMealPlanOptionsPermission,
		UpdateMealPlanOptionsPermission.ID():  UpdateMealPlanOptionsPermission,
		ArchiveMealPlanOptionsPermission.ID(): ArchiveMealPlanOptionsPermission,

		CreateHouseholdInstrumentOwnershipsPermission.ID():  CreateHouseholdInstrumentOwnershipsPermission,
		UpdateHouseholdInstrumentOwnershipsPermission.ID():  UpdateHouseholdInstrumentOwnershipsPermission,
		ArchiveHouseholdInstrumentOwnershipsPermission.ID(): ArchiveHouseholdInstrumentOwnershipsPermission,
	}

	// household member permissions.
	householdMemberPermissions = map[string]gorbac.Permission{
		ReadWebhooksPermission.ID(): ReadWebhooksPermission,

		CreateAPIClientsPermission.ID():  CreateAPIClientsPermission,
		ReadAPIClientsPermission.ID():    ReadAPIClientsPermission,
		ArchiveAPIClientsPermission.ID(): ArchiveAPIClientsPermission,
		ReadOAuth2ClientsPermission.ID(): ReadOAuth2ClientsPermission,

		ReadServiceSettingsPermission.ID():   ReadServiceSettingsPermission,
		SearchServiceSettingsPermission.ID(): SearchServiceSettingsPermission,

		CreateMealsPermission.ID():  CreateMealsPermission,
		ReadMealsPermission.ID():    ReadMealsPermission,
		UpdateMealsPermission.ID():  UpdateMealsPermission,
		ArchiveMealsPermission.ID(): ArchiveMealsPermission,

		ReadRecipesPermission.ID():    ReadRecipesPermission,
		SearchRecipesPermission.ID():  SearchRecipesPermission,
		UpdateRecipesPermission.ID():  UpdateRecipesPermission,
		ArchiveRecipesPermission.ID(): ArchiveRecipesPermission,

		CreateRecipeStepsPermission.ID():  CreateRecipeStepsPermission,
		ReadRecipeStepsPermission.ID():    ReadRecipeStepsPermission,
		SearchRecipeStepsPermission.ID():  SearchRecipeStepsPermission,
		UpdateRecipeStepsPermission.ID():  UpdateRecipeStepsPermission,
		ArchiveRecipeStepsPermission.ID(): ArchiveRecipeStepsPermission,

		CreateRecipePrepTasksPermission.ID():  CreateRecipePrepTasksPermission,
		ReadRecipePrepTasksPermission.ID():    ReadRecipePrepTasksPermission,
		UpdateRecipePrepTasksPermission.ID():  UpdateRecipePrepTasksPermission,
		ArchiveRecipePrepTasksPermission.ID(): ArchiveRecipePrepTasksPermission,

		CreateRecipeStepInstrumentsPermission.ID():  CreateRecipeStepInstrumentsPermission,
		ReadRecipeStepInstrumentsPermission.ID():    ReadRecipeStepInstrumentsPermission,
		SearchRecipeStepInstrumentsPermission.ID():  SearchRecipeStepInstrumentsPermission,
		UpdateRecipeStepInstrumentsPermission.ID():  UpdateRecipeStepInstrumentsPermission,
		ArchiveRecipeStepInstrumentsPermission.ID(): ArchiveRecipeStepInstrumentsPermission,

		CreateRecipeStepVesselsPermission.ID():  CreateRecipeStepVesselsPermission,
		ReadRecipeStepVesselsPermission.ID():    ReadRecipeStepVesselsPermission,
		SearchRecipeStepVesselsPermission.ID():  SearchRecipeStepVesselsPermission,
		UpdateRecipeStepVesselsPermission.ID():  UpdateRecipeStepVesselsPermission,
		ArchiveRecipeStepVesselsPermission.ID(): ArchiveRecipeStepVesselsPermission,

		CreateRecipeStepIngredientsPermission.ID():  CreateRecipeStepIngredientsPermission,
		ReadRecipeStepIngredientsPermission.ID():    ReadRecipeStepIngredientsPermission,
		SearchRecipeStepIngredientsPermission.ID():  SearchRecipeStepIngredientsPermission,
		UpdateRecipeStepIngredientsPermission.ID():  UpdateRecipeStepIngredientsPermission,
		ArchiveRecipeStepIngredientsPermission.ID(): ArchiveRecipeStepIngredientsPermission,

		CreateRecipeStepCompletionConditionsPermission.ID():  CreateRecipeStepCompletionConditionsPermission,
		ReadRecipeStepCompletionConditionsPermission.ID():    ReadRecipeStepCompletionConditionsPermission,
		SearchRecipeStepCompletionConditionsPermission.ID():  SearchRecipeStepCompletionConditionsPermission,
		UpdateRecipeStepCompletionConditionsPermission.ID():  UpdateRecipeStepCompletionConditionsPermission,
		ArchiveRecipeStepCompletionConditionsPermission.ID(): ArchiveRecipeStepCompletionConditionsPermission,

		CreateRecipeStepProductsPermission.ID():  CreateRecipeStepProductsPermission,
		ReadRecipeStepProductsPermission.ID():    ReadRecipeStepProductsPermission,
		SearchRecipeStepProductsPermission.ID():  SearchRecipeStepProductsPermission,
		UpdateRecipeStepProductsPermission.ID():  UpdateRecipeStepProductsPermission,
		ArchiveRecipeStepProductsPermission.ID(): ArchiveRecipeStepProductsPermission,

		ReadValidInstrumentsPermission.ID():   ReadValidInstrumentsPermission,
		SearchValidInstrumentsPermission.ID(): SearchValidInstrumentsPermission,

		ReadValidIngredientsPermission.ID():   ReadValidIngredientsPermission,
		SearchValidIngredientsPermission.ID(): SearchValidIngredientsPermission,

		ReadValidIngredientGroupsPermission.ID():   ReadValidIngredientGroupsPermission,
		SearchValidIngredientGroupsPermission.ID(): SearchValidIngredientGroupsPermission,

		ReadValidPreparationsPermission.ID():   ReadValidPreparationsPermission,
		SearchValidPreparationsPermission.ID(): SearchValidPreparationsPermission,

		ReadValidMeasurementUnitsPermission.ID():   ReadValidMeasurementUnitsPermission,
		SearchValidMeasurementUnitsPermission.ID(): SearchValidMeasurementUnitsPermission,

		ReadValidMeasurementConversionsPermission.ID(): ReadValidMeasurementConversionsPermission,

		ReadValidIngredientPreparationsPermission.ID():   ReadValidIngredientPreparationsPermission,
		SearchValidIngredientPreparationsPermission.ID(): SearchValidIngredientPreparationsPermission,

		ReadValidIngredientStateIngredientsPermission.ID():   ReadValidIngredientStateIngredientsPermission,
		SearchValidIngredientStateIngredientsPermission.ID(): SearchValidIngredientStateIngredientsPermission,

		ReadValidPreparationInstrumentsPermission.ID():   ReadValidPreparationInstrumentsPermission,
		SearchValidPreparationInstrumentsPermission.ID(): SearchValidPreparationInstrumentsPermission,

		ReadValidIngredientMeasurementUnitsPermission.ID():   ReadValidIngredientMeasurementUnitsPermission,
		SearchValidIngredientMeasurementUnitsPermission.ID(): SearchValidIngredientMeasurementUnitsPermission,

		ReadMealPlansPermission.ID():   ReadMealPlansPermission,
		SearchMealPlansPermission.ID(): SearchMealPlansPermission,

		ReadMealPlanEventsPermission.ID(): ReadMealPlanEventsPermission,

		ReadMealPlanOptionsPermission.ID():   ReadMealPlanOptionsPermission,
		SearchMealPlanOptionsPermission.ID(): SearchMealPlanOptionsPermission,

		ReadValidIngredientStatesPermission.ID(): ReadValidIngredientStatesPermission,

		ReadMealPlanGroceryListItemsPermission.ID():    ReadMealPlanGroceryListItemsPermission,
		UpdateMealPlanGroceryListItemsPermission.ID():  UpdateMealPlanGroceryListItemsPermission,
		ArchiveMealPlanGroceryListItemsPermission.ID(): ArchiveMealPlanGroceryListItemsPermission,

		CreateMealPlanOptionVotesPermission.ID():  CreateMealPlanOptionVotesPermission,
		ReadMealPlanOptionVotesPermission.ID():    ReadMealPlanOptionVotesPermission,
		SearchMealPlanOptionVotesPermission.ID():  SearchMealPlanOptionVotesPermission,
		UpdateMealPlanOptionVotesPermission.ID():  UpdateMealPlanOptionVotesPermission,
		ArchiveMealPlanOptionVotesPermission.ID(): ArchiveMealPlanOptionVotesPermission,

		CreateServiceSettingConfigurationsPermission.ID():  CreateServiceSettingConfigurationsPermission,
		ReadServiceSettingConfigurationsPermission.ID():    ReadServiceSettingConfigurationsPermission,
		UpdateServiceSettingConfigurationsPermission.ID():  UpdateServiceSettingConfigurationsPermission,
		ArchiveServiceSettingConfigurationsPermission.ID(): ArchiveServiceSettingConfigurationsPermission,

		ReadMealPlanTasksPermission.ID():   ReadMealPlanTasksPermission,
		UpdateMealPlanTasksPermission.ID(): UpdateMealPlanTasksPermission,

		CreateUserIngredientPreferencesPermission.ID():  CreateUserIngredientPreferencesPermission,
		ReadUserIngredientPreferencesPermission.ID():    ReadUserIngredientPreferencesPermission,
		UpdateUserIngredientPreferencesPermission.ID():  UpdateUserIngredientPreferencesPermission,
		ArchiveUserIngredientPreferencesPermission.ID(): ArchiveUserIngredientPreferencesPermission,

		ReadHouseholdInstrumentOwnershipsPermission.ID(): ReadHouseholdInstrumentOwnershipsPermission,

		CreateRecipeRatingsPermission.ID():  CreateRecipeRatingsPermission,
		ReadRecipeRatingsPermission.ID():    ReadRecipeRatingsPermission,
		UpdateRecipeRatingsPermission.ID():  UpdateRecipeRatingsPermission,
		ArchiveRecipeRatingsPermission.ID(): ArchiveRecipeRatingsPermission,
	}
)

func init() {
	// assign service admin permissions.
	for _, perm := range serviceAdminPermissions {
		must(serviceAdmin.Assign(perm))
	}

	// assign household admin permissions.
	for _, perm := range householdAdminPermissions {
		must(householdAdmin.Assign(perm))
	}

	// assign household member permissions.
	for _, perm := range householdMemberPermissions {
		must(householdMember.Assign(perm))
	}
}
