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

	// UpdateAccountPermission is an account admin permission.
	UpdateAccountPermission Permission = "update.account"
	// ArchiveAccountPermission is an account admin permission.
	ArchiveAccountPermission Permission = "archive.account"
	// AddMemberAccountPermission is an account admin permission.
	AddMemberAccountPermission Permission = "account.add.member"
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
	// CreateAPIClientsPermission is an account admin permission.
	CreateAPIClientsPermission Permission = "create.api_clients"
	// ReadAPIClientsPermission is an account admin permission.
	ReadAPIClientsPermission Permission = "read.api_clients"
	// ArchiveAPIClientsPermission is an account admin permission.
	ArchiveAPIClientsPermission Permission = "archive.api_clients"

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

	// account admin permissions.
	accountAdminPermissions = map[string]gorbac.Permission{
		UpdateAccountPermission.ID():                     UpdateAccountPermission,
		ArchiveAccountPermission.ID():                    ArchiveAccountPermission,
		AddMemberAccountPermission.ID():                  AddMemberAccountPermission,
		ModifyMemberPermissionsForAccountPermission.ID(): ModifyMemberPermissionsForAccountPermission,
		RemoveMemberAccountPermission.ID():               RemoveMemberAccountPermission,
		TransferAccountPermission.ID():                   TransferAccountPermission,
		CreateWebhooksPermission.ID():                    CreateWebhooksPermission,
		ReadWebhooksPermission.ID():                      ReadWebhooksPermission,
		UpdateWebhooksPermission.ID():                    UpdateWebhooksPermission,
		ArchiveWebhooksPermission.ID():                   ArchiveWebhooksPermission,
		CreateAPIClientsPermission.ID():                  CreateAPIClientsPermission,
		ReadAPIClientsPermission.ID():                    ReadAPIClientsPermission,
		ArchiveAPIClientsPermission.ID():                 ArchiveAPIClientsPermission,
	}

	// account member permissions.
	accountMemberPermissions = map[string]gorbac.Permission{
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

	// assign account admin permissions.
	for _, perm := range accountAdminPermissions {
		must(accountAdmin.Assign(perm))
	}

	// assign account member permissions.
	for _, perm := range accountMemberPermissions {
		must(accountMember.Assign(perm))
	}
}
