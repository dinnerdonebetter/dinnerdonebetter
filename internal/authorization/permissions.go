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
	// ReadHouseholdAuditLogEntriesPermission is a service admin permission.
	ReadHouseholdAuditLogEntriesPermission Permission = "read.audit_log_entries.household"
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

	// UpdateHouseholdPermission is an household admin permission.
	UpdateHouseholdPermission Permission = "update.household"
	// ArchiveHouseholdPermission is an household admin permission.
	ArchiveHouseholdPermission Permission = "archive.household"
	// AddMemberHouseholdPermission is an household admin permission.
	AddMemberHouseholdPermission Permission = "household.add.member"
	// ModifyMemberPermissionsForHouseholdPermission is an household admin permission.
	ModifyMemberPermissionsForHouseholdPermission Permission = "household.membership.modify"
	// RemoveMemberHouseholdPermission is an household admin permission.
	RemoveMemberHouseholdPermission Permission = "remove_member.household"
	// TransferHouseholdPermission is an household admin permission.
	TransferHouseholdPermission Permission = "transfer.household"
	// CreateWebhooksPermission is an household admin permission.
	CreateWebhooksPermission Permission = "create.webhooks"
	// ReadWebhooksPermission is an household admin permission.
	ReadWebhooksPermission Permission = "read.webhooks"
	// UpdateWebhooksPermission is an household admin permission.
	UpdateWebhooksPermission Permission = "update.webhooks"
	// ArchiveWebhooksPermission is an household admin permission.
	ArchiveWebhooksPermission Permission = "archive.webhooks"
	// CreateAPIClientsPermission is an household admin permission.
	CreateAPIClientsPermission Permission = "create.api_clients"
	// ReadAPIClientsPermission is an household admin permission.
	ReadAPIClientsPermission Permission = "read.api_clients"
	// ArchiveAPIClientsPermission is an household admin permission.
	ArchiveAPIClientsPermission Permission = "archive.api_clients"
	// ReadWebhooksAuditLogEntriesPermission is an household admin permission.
	ReadWebhooksAuditLogEntriesPermission Permission = "read.audit_log_entries.webhooks"
	// ReadValidInstrumentsAuditLogEntriesPermission is an household admin permission.
	ReadValidInstrumentsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_instruments"
	// ReadValidPreparationsAuditLogEntriesPermission is an household admin permission.
	ReadValidPreparationsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_preparations"
	// ReadValidIngredientsAuditLogEntriesPermission is an household admin permission.
	ReadValidIngredientsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_ingredients"
	// ReadValidIngredientPreparationsAuditLogEntriesPermission is an household admin permission.
	ReadValidIngredientPreparationsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_ingredient_preparations"
	// ReadValidPreparationInstrumentsAuditLogEntriesPermission is an household admin permission.
	ReadValidPreparationInstrumentsAuditLogEntriesPermission Permission = "read.audit_log_entries.valid_preparation_instruments"
	// ReadRecipesAuditLogEntriesPermission is an household admin permission.
	ReadRecipesAuditLogEntriesPermission Permission = "read.audit_log_entries.recipes"
	// ReadRecipeStepsAuditLogEntriesPermission is an household admin permission.
	ReadRecipeStepsAuditLogEntriesPermission Permission = "read.audit_log_entries.recipe_steps"
	// ReadRecipeStepIngredientsAuditLogEntriesPermission is an household admin permission.
	ReadRecipeStepIngredientsAuditLogEntriesPermission Permission = "read.audit_log_entries.recipe_step_ingredients"
	// ReadRecipeStepProductsAuditLogEntriesPermission is an household admin permission.
	ReadRecipeStepProductsAuditLogEntriesPermission Permission = "read.audit_log_entries.recipe_step_products"
	// ReadInvitationsAuditLogEntriesPermission is an household admin permission.
	ReadInvitationsAuditLogEntriesPermission Permission = "read.audit_log_entries.invitations"
	// ReadReportsAuditLogEntriesPermission is an household admin permission.
	ReadReportsAuditLogEntriesPermission Permission = "read.audit_log_entries.reports"

	// CreateValidInstrumentsPermission is an household user permission.
	CreateValidInstrumentsPermission Permission = "create.valid_instruments"
	// ReadValidInstrumentsPermission is an household user permission.
	ReadValidInstrumentsPermission Permission = "read.valid_instruments"
	// SearchValidInstrumentsPermission is an household user permission.
	SearchValidInstrumentsPermission Permission = "search.valid_instruments"
	// UpdateValidInstrumentsPermission is an household user permission.
	UpdateValidInstrumentsPermission Permission = "update.valid_instruments"
	// ArchiveValidInstrumentsPermission is an household user permission.
	ArchiveValidInstrumentsPermission Permission = "archive.valid_instruments"

	// CreateValidPreparationsPermission is an household user permission.
	CreateValidPreparationsPermission Permission = "create.valid_preparations"
	// ReadValidPreparationsPermission is an household user permission.
	ReadValidPreparationsPermission Permission = "read.valid_preparations"
	// SearchValidPreparationsPermission is an household user permission.
	SearchValidPreparationsPermission Permission = "search.valid_preparations"
	// UpdateValidPreparationsPermission is an household user permission.
	UpdateValidPreparationsPermission Permission = "update.valid_preparations"
	// ArchiveValidPreparationsPermission is an household user permission.
	ArchiveValidPreparationsPermission Permission = "archive.valid_preparations"

	// CreateValidIngredientsPermission is an household user permission.
	CreateValidIngredientsPermission Permission = "create.valid_ingredients"
	// ReadValidIngredientsPermission is an household user permission.
	ReadValidIngredientsPermission Permission = "read.valid_ingredients"
	// SearchValidIngredientsPermission is an household user permission.
	SearchValidIngredientsPermission Permission = "search.valid_ingredients"
	// UpdateValidIngredientsPermission is an household user permission.
	UpdateValidIngredientsPermission Permission = "update.valid_ingredients"
	// ArchiveValidIngredientsPermission is an household user permission.
	ArchiveValidIngredientsPermission Permission = "archive.valid_ingredients"

	// CreateValidIngredientPreparationsPermission is an household user permission.
	CreateValidIngredientPreparationsPermission Permission = "create.valid_ingredient_preparations"
	// ReadValidIngredientPreparationsPermission is an household user permission.
	ReadValidIngredientPreparationsPermission Permission = "read.valid_ingredient_preparations"
	// SearchValidIngredientPreparationsPermission is an household user permission.
	SearchValidIngredientPreparationsPermission Permission = "search.valid_ingredient_preparations"
	// UpdateValidIngredientPreparationsPermission is an household user permission.
	UpdateValidIngredientPreparationsPermission Permission = "update.valid_ingredient_preparations"
	// ArchiveValidIngredientPreparationsPermission is an household user permission.
	ArchiveValidIngredientPreparationsPermission Permission = "archive.valid_ingredient_preparations"

	// CreateValidPreparationInstrumentsPermission is an household user permission.
	CreateValidPreparationInstrumentsPermission Permission = "create.valid_preparation_instruments"
	// ReadValidPreparationInstrumentsPermission is an household user permission.
	ReadValidPreparationInstrumentsPermission Permission = "read.valid_preparation_instruments"
	// SearchValidPreparationInstrumentsPermission is an household user permission.
	SearchValidPreparationInstrumentsPermission Permission = "search.valid_preparation_instruments"
	// UpdateValidPreparationInstrumentsPermission is an household user permission.
	UpdateValidPreparationInstrumentsPermission Permission = "update.valid_preparation_instruments"
	// ArchiveValidPreparationInstrumentsPermission is an household user permission.
	ArchiveValidPreparationInstrumentsPermission Permission = "archive.valid_preparation_instruments"

	// CreateRecipesPermission is an household user permission.
	CreateRecipesPermission Permission = "create.recipes"
	// ReadRecipesPermission is an household user permission.
	ReadRecipesPermission Permission = "read.recipes"
	// SearchRecipesPermission is an household user permission.
	SearchRecipesPermission Permission = "search.recipes"
	// UpdateRecipesPermission is an household user permission.
	UpdateRecipesPermission Permission = "update.recipes"
	// ArchiveRecipesPermission is an household user permission.
	ArchiveRecipesPermission Permission = "archive.recipes"

	// CreateRecipeStepsPermission is an household user permission.
	CreateRecipeStepsPermission Permission = "create.recipe_steps"
	// ReadRecipeStepsPermission is an household user permission.
	ReadRecipeStepsPermission Permission = "read.recipe_steps"
	// SearchRecipeStepsPermission is an household user permission.
	SearchRecipeStepsPermission Permission = "search.recipe_steps"
	// UpdateRecipeStepsPermission is an household user permission.
	UpdateRecipeStepsPermission Permission = "update.recipe_steps"
	// ArchiveRecipeStepsPermission is an household user permission.
	ArchiveRecipeStepsPermission Permission = "archive.recipe_steps"

	// CreateRecipeStepIngredientsPermission is an household user permission.
	CreateRecipeStepIngredientsPermission Permission = "create.recipe_step_ingredients"
	// ReadRecipeStepIngredientsPermission is an household user permission.
	ReadRecipeStepIngredientsPermission Permission = "read.recipe_step_ingredients"
	// SearchRecipeStepIngredientsPermission is an household user permission.
	SearchRecipeStepIngredientsPermission Permission = "search.recipe_step_ingredients"
	// UpdateRecipeStepIngredientsPermission is an household user permission.
	UpdateRecipeStepIngredientsPermission Permission = "update.recipe_step_ingredients"
	// ArchiveRecipeStepIngredientsPermission is an household user permission.
	ArchiveRecipeStepIngredientsPermission Permission = "archive.recipe_step_ingredients"

	// CreateRecipeStepProductsPermission is an household user permission.
	CreateRecipeStepProductsPermission Permission = "create.recipe_step_products"
	// ReadRecipeStepProductsPermission is an household user permission.
	ReadRecipeStepProductsPermission Permission = "read.recipe_step_products"
	// SearchRecipeStepProductsPermission is an household user permission.
	SearchRecipeStepProductsPermission Permission = "search.recipe_step_products"
	// UpdateRecipeStepProductsPermission is an household user permission.
	UpdateRecipeStepProductsPermission Permission = "update.recipe_step_products"
	// ArchiveRecipeStepProductsPermission is an household user permission.
	ArchiveRecipeStepProductsPermission Permission = "archive.recipe_step_products"

	// CreateInvitationsPermission is an household user permission.
	CreateInvitationsPermission Permission = "create.invitations"
	// ReadInvitationsPermission is an household user permission.
	ReadInvitationsPermission Permission = "read.invitations"
	// SearchInvitationsPermission is an household user permission.
	SearchInvitationsPermission Permission = "search.invitations"
	// UpdateInvitationsPermission is an household user permission.
	UpdateInvitationsPermission Permission = "update.invitations"
	// ArchiveInvitationsPermission is an household user permission.
	ArchiveInvitationsPermission Permission = "archive.invitations"

	// CreateReportsPermission is an household user permission.
	CreateReportsPermission Permission = "create.reports"
	// ReadReportsPermission is an household user permission.
	ReadReportsPermission Permission = "read.reports"
	// SearchReportsPermission is an household user permission.
	SearchReportsPermission Permission = "search.reports"
	// UpdateReportsPermission is an household user permission.
	UpdateReportsPermission Permission = "update.reports"
	// ArchiveReportsPermission is an household user permission.
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
		ReadHouseholdAuditLogEntriesPermission.ID(): ReadHouseholdAuditLogEntriesPermission,
		ReadAPIClientAuditLogEntriesPermission.ID(): ReadAPIClientAuditLogEntriesPermission,
		ReadUserAuditLogEntriesPermission.ID():      ReadUserAuditLogEntriesPermission,
		ReadWebhookAuditLogEntriesPermission.ID():   ReadWebhookAuditLogEntriesPermission,
		UpdateUserStatusPermission.ID():             UpdateUserStatusPermission,
		ReadUserPermission.ID():                     ReadUserPermission,
		SearchUserPermission.ID():                   SearchUserPermission,
	}

	// household admin permissions.
	householdAdminPermissions = map[string]gorbac.Permission{
		UpdateHouseholdPermission.ID():                                UpdateHouseholdPermission,
		ArchiveHouseholdPermission.ID():                               ArchiveHouseholdPermission,
		AddMemberHouseholdPermission.ID():                             AddMemberHouseholdPermission,
		ModifyMemberPermissionsForHouseholdPermission.ID():            ModifyMemberPermissionsForHouseholdPermission,
		RemoveMemberHouseholdPermission.ID():                          RemoveMemberHouseholdPermission,
		TransferHouseholdPermission.ID():                              TransferHouseholdPermission,
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

	// household member permissions.
	householdMemberPermissions = map[string]gorbac.Permission{
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

	// assign household admin permissions.
	for _, perm := range householdAdminPermissions {
		must(householdAdmin.Assign(perm))
	}

	// assign household member permissions.
	for _, perm := range householdMemberPermissions {
		must(householdMember.Assign(perm))
	}
}
