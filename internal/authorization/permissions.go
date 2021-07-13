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
	// ReadAllAuditLogEntriesPermission is a service admin permission.
	ReadAllAuditLogEntriesPermission Permission = "read.audit_log_entries.all"
	// ReadAccountAuditLogEntriesPermission is a service admin permission.
	ReadAccountAuditLogEntriesPermission Permission = "read.audit_log_entries.account"
	// ReadAPIClientAuditLogEntriesPermission is a service admin permission.
	ReadAPIClientAuditLogEntriesPermission Permission = "read.audit_log_entries.api_client"
	// ReadUserAuditLogEntriesPermission is a service admin permission.
	ReadUserAuditLogEntriesPermission Permission = "read.audit_log_entries.user"
	// ReadWebhookAuditLogEntriesPermission is a service admin permission.
	ReadWebhookAuditLogEntriesPermission Permission = "read.audit_log_entries.webhook"
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
	// ReadWebhooksAuditLogEntriesPermission is an account admin permission.
	ReadWebhooksAuditLogEntriesPermission Permission = "read.audit_log_entries.webhooks"
	// ReadValidInstrumentsAuditLogEntriesPermission is an account admin permission.
	ReadValidInstrumentsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_instruments"
	// ReadValidPreparationsAuditLogEntriesPermission is an account admin permission.
	ReadValidPreparationsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_preparations"
	// ReadValidIngredientsAuditLogEntriesPermission is an account admin permission.
	ReadValidIngredientsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_ingredients"
	// ReadValidIngredientPreparationsAuditLogEntriesPermission is an account admin permission.
	ReadValidIngredientPreparationsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_ingredient_preparations"
	// ReadValidPreparationInstrumentsAuditLogEntriesPermission is an account admin permission.
	ReadValidPreparationInstrumentsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_preparation_instruments"
	// ReadRecipesAuditLogEntriesPermission is an account admin permission.
	ReadRecipesAuditLogEntriesPermission Permission = "read.audit_log_entries.recipes"
	// ReadRecipeStepsAuditLogEntriesPermission is an account admin permission.
	ReadRecipeStepsAuditLogEntriesPermission Permission = "read.audit_log_entries.recipe_steps"
	// ReadRecipeStepIngredientsAuditLogEntriesPermission is an account admin permission.
	ReadRecipeStepIngredientsAuditLogEntriesPermission Permission = "read.audit_log_entries.recipe_step_ingredients"
	// ReadRecipeStepProductsAuditLogEntriesPermission is an account admin permission.
	ReadRecipeStepProductsAuditLogEntriesPermission Permission = "read.audit_log_entries.recipe_step_products"
	// ReadInvitationsAuditLogEntriesPermission is an account admin permission.
	ReadInvitationsAuditLogEntriesPermission Permission = "read.audit_log_entries.invitations"
	// ReadReportsAuditLogEntriesPermission is an account admin permission.
	ReadReportsAuditLogEntriesPermission Permission = "read.audit_log_entries.reports"

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

	// CreateInvitationsPermission is an account user permission.
	CreateInvitationsPermission Permission = "create.invitations"
	// ReadInvitationsPermission is an account user permission.
	ReadInvitationsPermission Permission = "read.invitations"
	// SearchInvitationsPermission is an account user permission.
	SearchInvitationsPermission Permission = "search.invitations"
	// UpdateInvitationsPermission is an account user permission.
	UpdateInvitationsPermission Permission = "update.invitations"
	// ArchiveInvitationsPermission is an account user permission.
	ArchiveInvitationsPermission Permission = "archive.invitations"

	// CreateReportsPermission is an account user permission.
	CreateReportsPermission Permission = "create.reports"
	// ReadReportsPermission is an account user permission.
	ReadReportsPermission Permission = "read.reports"
	// SearchReportsPermission is an account user permission.
	SearchReportsPermission Permission = "search.reports"
	// UpdateReportsPermission is an account user permission.
	UpdateReportsPermission Permission = "update.reports"
	// ArchiveReportsPermission is an account user permission.
	ArchiveReportsPermission Permission = "archive.reports"
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
		CycleCookieSecretPermission.ID():            CycleCookieSecretPermission,
		ReadAllAuditLogEntriesPermission.ID():       ReadAllAuditLogEntriesPermission,
		ReadAccountAuditLogEntriesPermission.ID():   ReadAccountAuditLogEntriesPermission,
		ReadAPIClientAuditLogEntriesPermission.ID(): ReadAPIClientAuditLogEntriesPermission,
		ReadUserAuditLogEntriesPermission.ID():      ReadUserAuditLogEntriesPermission,
		ReadWebhookAuditLogEntriesPermission.ID():   ReadWebhookAuditLogEntriesPermission,
		UpdateUserStatusPermission.ID():             UpdateUserStatusPermission,
		ReadUserPermission.ID():                     ReadUserPermission,
		SearchUserPermission.ID():                   SearchUserPermission,
	}

	// account admin permissions.
	accountAdminPermissions = map[string]gorbac.Permission{
		UpdateAccountPermission.ID():                                  UpdateAccountPermission,
		ArchiveAccountPermission.ID():                                 ArchiveAccountPermission,
		AddMemberAccountPermission.ID():                               AddMemberAccountPermission,
		ModifyMemberPermissionsForAccountPermission.ID():              ModifyMemberPermissionsForAccountPermission,
		RemoveMemberAccountPermission.ID():                            RemoveMemberAccountPermission,
		TransferAccountPermission.ID():                                TransferAccountPermission,
		CreateWebhooksPermission.ID():                                 CreateWebhooksPermission,
		ReadWebhooksPermission.ID():                                   ReadWebhooksPermission,
		UpdateWebhooksPermission.ID():                                 UpdateWebhooksPermission,
		ArchiveWebhooksPermission.ID():                                ArchiveWebhooksPermission,
		CreateAPIClientsPermission.ID():                               CreateAPIClientsPermission,
		ReadAPIClientsPermission.ID():                                 ReadAPIClientsPermission,
		ArchiveAPIClientsPermission.ID():                              ArchiveAPIClientsPermission,
		ReadWebhooksAuditLogEntriesPermission.ID():                    ReadWebhooksAuditLogEntriesPermission,
		ReadValidInstrumentsAuditLogEntriesPermission.ID():            ReadValidInstrumentsAuditLogEntriesPermission,
		ReadValidPreparationsAuditLogEntriesPermission.ID():           ReadValidPreparationsAuditLogEntriesPermission,
		ReadValidIngredientsAuditLogEntriesPermission.ID():            ReadValidIngredientsAuditLogEntriesPermission,
		ReadValidIngredientPreparationsAuditLogEntriesPermission.ID(): ReadValidIngredientPreparationsAuditLogEntriesPermission,
		ReadValidPreparationInstrumentsAuditLogEntriesPermission.ID(): ReadValidPreparationInstrumentsAuditLogEntriesPermission,
		ReadRecipesAuditLogEntriesPermission.ID():                     ReadRecipesAuditLogEntriesPermission,
		ReadRecipeStepsAuditLogEntriesPermission.ID():                 ReadRecipeStepsAuditLogEntriesPermission,
		ReadRecipeStepIngredientsAuditLogEntriesPermission.ID():       ReadRecipeStepIngredientsAuditLogEntriesPermission,
		ReadRecipeStepProductsAuditLogEntriesPermission.ID():          ReadRecipeStepProductsAuditLogEntriesPermission,
		ReadInvitationsAuditLogEntriesPermission.ID():                 ReadInvitationsAuditLogEntriesPermission,
		ReadReportsAuditLogEntriesPermission.ID():                     ReadReportsAuditLogEntriesPermission,
	}

	// account member permissions.
	accountMemberPermissions = map[string]gorbac.Permission{
		CreateValidInstrumentsPermission.ID():  CreateValidInstrumentsPermission,
		ReadValidInstrumentsPermission.ID():    ReadValidInstrumentsPermission,
		SearchValidInstrumentsPermission.ID():  SearchValidInstrumentsPermission,
		UpdateValidInstrumentsPermission.ID():  UpdateValidInstrumentsPermission,
		ArchiveValidInstrumentsPermission.ID(): ArchiveValidInstrumentsPermission,

		CreateValidPreparationsPermission.ID():  CreateValidPreparationsPermission,
		ReadValidPreparationsPermission.ID():    ReadValidPreparationsPermission,
		SearchValidPreparationsPermission.ID():  SearchValidPreparationsPermission,
		UpdateValidPreparationsPermission.ID():  UpdateValidPreparationsPermission,
		ArchiveValidPreparationsPermission.ID(): ArchiveValidPreparationsPermission,

		CreateValidIngredientsPermission.ID():  CreateValidIngredientsPermission,
		ReadValidIngredientsPermission.ID():    ReadValidIngredientsPermission,
		SearchValidIngredientsPermission.ID():  SearchValidIngredientsPermission,
		UpdateValidIngredientsPermission.ID():  UpdateValidIngredientsPermission,
		ArchiveValidIngredientsPermission.ID(): ArchiveValidIngredientsPermission,

		CreateValidIngredientPreparationsPermission.ID():  CreateValidIngredientPreparationsPermission,
		ReadValidIngredientPreparationsPermission.ID():    ReadValidIngredientPreparationsPermission,
		SearchValidIngredientPreparationsPermission.ID():  SearchValidIngredientPreparationsPermission,
		UpdateValidIngredientPreparationsPermission.ID():  UpdateValidIngredientPreparationsPermission,
		ArchiveValidIngredientPreparationsPermission.ID(): ArchiveValidIngredientPreparationsPermission,

		CreateValidPreparationInstrumentsPermission.ID():  CreateValidPreparationInstrumentsPermission,
		ReadValidPreparationInstrumentsPermission.ID():    ReadValidPreparationInstrumentsPermission,
		SearchValidPreparationInstrumentsPermission.ID():  SearchValidPreparationInstrumentsPermission,
		UpdateValidPreparationInstrumentsPermission.ID():  UpdateValidPreparationInstrumentsPermission,
		ArchiveValidPreparationInstrumentsPermission.ID(): ArchiveValidPreparationInstrumentsPermission,

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

		CreateInvitationsPermission.ID():  CreateInvitationsPermission,
		ReadInvitationsPermission.ID():    ReadInvitationsPermission,
		SearchInvitationsPermission.ID():  SearchInvitationsPermission,
		UpdateInvitationsPermission.ID():  UpdateInvitationsPermission,
		ArchiveInvitationsPermission.ID(): ArchiveInvitationsPermission,

		CreateReportsPermission.ID():  CreateReportsPermission,
		ReadReportsPermission.ID():    ReadReportsPermission,
		SearchReportsPermission.ID():  SearchReportsPermission,
		UpdateReportsPermission.ID():  UpdateReportsPermission,
		ArchiveReportsPermission.ID(): ArchiveReportsPermission,
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
