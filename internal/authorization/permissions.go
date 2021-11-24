package authorization

import (
	"gopkg.in/mikespook/gorbac.v2"
)

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

	// CreateMealsPermission is a household user permission.
	CreateMealsPermission Permission = "create.mealss"
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
	}

	// household admin permissions.
	householdAdminPermissions = map[string]gorbac.Permission{
		UpdateHouseholdPermission.ID():                     UpdateHouseholdPermission,
		ArchiveHouseholdPermission.ID():                    ArchiveHouseholdPermission,
		InviteUserToHouseholdPermission.ID():               InviteUserToHouseholdPermission,
		ModifyMemberPermissionsForHouseholdPermission.ID(): ModifyMemberPermissionsForHouseholdPermission,
		RemoveMemberHouseholdPermission.ID():               RemoveMemberHouseholdPermission,
		TransferHouseholdPermission.ID():                   TransferHouseholdPermission,
		CreateWebhooksPermission.ID():                      CreateWebhooksPermission,
		UpdateWebhooksPermission.ID():                      UpdateWebhooksPermission,
		ArchiveWebhooksPermission.ID():                     ArchiveWebhooksPermission,

		CreateValidInstrumentsPermission.ID():  CreateValidInstrumentsPermission,
		ReadValidInstrumentsPermission.ID():    ReadValidInstrumentsPermission,
		SearchValidInstrumentsPermission.ID():  SearchValidInstrumentsPermission,
		UpdateValidInstrumentsPermission.ID():  UpdateValidInstrumentsPermission,
		ArchiveValidInstrumentsPermission.ID(): ArchiveValidInstrumentsPermission,

		CreateValidIngredientsPermission.ID():  CreateValidIngredientsPermission,
		ReadValidIngredientsPermission.ID():    ReadValidIngredientsPermission,
		SearchValidIngredientsPermission.ID():  SearchValidIngredientsPermission,
		UpdateValidIngredientsPermission.ID():  UpdateValidIngredientsPermission,
		ArchiveValidIngredientsPermission.ID(): ArchiveValidIngredientsPermission,

		CreateValidPreparationsPermission.ID():  CreateValidPreparationsPermission,
		ReadValidPreparationsPermission.ID():    ReadValidPreparationsPermission,
		SearchValidPreparationsPermission.ID():  SearchValidPreparationsPermission,
		UpdateValidPreparationsPermission.ID():  UpdateValidPreparationsPermission,
		ArchiveValidPreparationsPermission.ID(): ArchiveValidPreparationsPermission,

		CreateValidIngredientPreparationsPermission.ID():  CreateValidIngredientPreparationsPermission,
		ReadValidIngredientPreparationsPermission.ID():    ReadValidIngredientPreparationsPermission,
		SearchValidIngredientPreparationsPermission.ID():  SearchValidIngredientPreparationsPermission,
		UpdateValidIngredientPreparationsPermission.ID():  UpdateValidIngredientPreparationsPermission,
		ArchiveValidIngredientPreparationsPermission.ID(): ArchiveValidIngredientPreparationsPermission,

		CreateMealPlansPermission.ID():  CreateMealPlansPermission,
		ReadMealPlansPermission.ID():    ReadMealPlansPermission,
		SearchMealPlansPermission.ID():  SearchMealPlansPermission,
		UpdateMealPlansPermission.ID():  UpdateMealPlansPermission,
		ArchiveMealPlansPermission.ID(): ArchiveMealPlansPermission,

		CreateMealPlanOptionsPermission.ID():  CreateMealPlanOptionsPermission,
		ReadMealPlanOptionsPermission.ID():    ReadMealPlanOptionsPermission,
		SearchMealPlanOptionsPermission.ID():  SearchMealPlanOptionsPermission,
		UpdateMealPlanOptionsPermission.ID():  UpdateMealPlanOptionsPermission,
		ArchiveMealPlanOptionsPermission.ID(): ArchiveMealPlanOptionsPermission,
	}

	// household member permissions.
	householdMemberPermissions = map[string]gorbac.Permission{
		ReadWebhooksPermission.ID():      ReadWebhooksPermission,
		CreateAPIClientsPermission.ID():  CreateAPIClientsPermission,
		ReadAPIClientsPermission.ID():    ReadAPIClientsPermission,
		ArchiveAPIClientsPermission.ID(): ArchiveAPIClientsPermission,

		CreateMealsPermission.ID():  CreateMealsPermission,
		ReadMealsPermission.ID():    ReadMealsPermission,
		UpdateMealsPermission.ID():  UpdateMealsPermission,
		ArchiveMealsPermission.ID(): ArchiveMealsPermission,

		CreateRecipesPermission.ID():  CreateRecipesPermission,
		ReadRecipesPermission.ID():    ReadRecipesPermission,
		SearchRecipesPermission.ID():  SearchRecipesPermission,
		UpdateRecipesPermission.ID():  UpdateRecipesPermission,
		ArchiveRecipesPermission.ID(): ArchiveRecipesPermission,

		CreateRecipeStepsPermission.ID():  CreateRecipeStepsPermission,
		ReadRecipeStepsPermission.ID():    ReadRecipeStepsPermission,
		SearchRecipeStepsPermission.ID():  SearchRecipeStepsPermission,
		UpdateRecipeStepsPermission.ID():  UpdateRecipeStepsPermission,
		ArchiveRecipeStepsPermission.ID(): ArchiveRecipeStepsPermission,

		CreateRecipeStepInstrumentsPermission.ID():  CreateRecipeStepInstrumentsPermission,
		ReadRecipeStepInstrumentsPermission.ID():    ReadRecipeStepInstrumentsPermission,
		SearchRecipeStepInstrumentsPermission.ID():  SearchRecipeStepInstrumentsPermission,
		UpdateRecipeStepInstrumentsPermission.ID():  UpdateRecipeStepInstrumentsPermission,
		ArchiveRecipeStepInstrumentsPermission.ID(): ArchiveRecipeStepInstrumentsPermission,

		CreateRecipeStepIngredientsPermission.ID():  CreateRecipeStepIngredientsPermission,
		ReadRecipeStepIngredientsPermission.ID():    ReadRecipeStepIngredientsPermission,
		SearchRecipeStepIngredientsPermission.ID():  SearchRecipeStepIngredientsPermission,
		UpdateRecipeStepIngredientsPermission.ID():  UpdateRecipeStepIngredientsPermission,
		ArchiveRecipeStepIngredientsPermission.ID(): ArchiveRecipeStepIngredientsPermission,

		CreateRecipeStepProductsPermission.ID():  CreateRecipeStepProductsPermission,
		ReadRecipeStepProductsPermission.ID():    ReadRecipeStepProductsPermission,
		SearchRecipeStepProductsPermission.ID():  SearchRecipeStepProductsPermission,
		UpdateRecipeStepProductsPermission.ID():  UpdateRecipeStepProductsPermission,
		ArchiveRecipeStepProductsPermission.ID(): ArchiveRecipeStepProductsPermission,

		CreateMealPlanOptionVotesPermission.ID():  CreateMealPlanOptionVotesPermission,
		ReadMealPlanOptionVotesPermission.ID():    ReadMealPlanOptionVotesPermission,
		SearchMealPlanOptionVotesPermission.ID():  SearchMealPlanOptionVotesPermission,
		UpdateMealPlanOptionVotesPermission.ID():  UpdateMealPlanOptionVotesPermission,
		ArchiveMealPlanOptionVotesPermission.ID(): ArchiveMealPlanOptionVotesPermission,
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
