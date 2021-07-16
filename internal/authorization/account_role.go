package authorization

import (
	"encoding/gob"

	"gopkg.in/mikespook/gorbac.v2"
)

type (
	// AccountRole describes a role a user has for an Account context.
	AccountRole role

	// AccountRolePermissionsChecker checks permissions for one or more account Roles.
	AccountRolePermissionsChecker interface {
		HasPermission(Permission) bool

		CanUpdateAccounts() bool
		CanDeleteAccounts() bool
		CanAddMemberToAccounts() bool
		CanRemoveMemberFromAccounts() bool
		CanTransferAccountToNewOwner() bool
		CanCreateWebhooks() bool
		CanSeeWebhooks() bool
		CanUpdateWebhooks() bool
		CanArchiveWebhooks() bool
		CanCreateAPIClients() bool
		CanSeeAPIClients() bool
		CanDeleteAPIClients() bool
		CanSeeAuditLogEntriesForWebhooks() bool
		CanSeeAuditLogEntriesForValidInstruments() bool
		CanSeeAuditLogEntriesForValidPreparations() bool
		CanSeeAuditLogEntriesForValidIngredients() bool
		CanSeeAuditLogEntriesForValidIngredientPreparations() bool
		CanSeeAuditLogEntriesForValidPreparationInstruments() bool
		CanSeeAuditLogEntriesForRecipes() bool
		CanSeeAuditLogEntriesForRecipeSteps() bool
		CanSeeAuditLogEntriesForRecipeStepIngredients() bool
		CanSeeAuditLogEntriesForRecipeStepProducts() bool
		CanSeeAuditLogEntriesForInvitations() bool
		CanSeeAuditLogEntriesForReports() bool
	}
)

const (
	// AccountMemberRole is a role for a plain account participant.
	AccountMemberRole AccountRole = iota
	// AccountAdminRole is a role for someone who can manipulate the specifics of an account.
	AccountAdminRole AccountRole = iota

	accountAdminRoleName  = "account_admin"
	accountMemberRoleName = "account_member"
)

var (
	accountAdmin  = gorbac.NewStdRole(accountAdminRoleName)
	accountMember = gorbac.NewStdRole(accountMemberRoleName)
)

type accountRoleCollection struct {
	Roles []string
}

func init() {
	gob.Register(accountRoleCollection{})
}

// NewAccountRolePermissionChecker returns a new checker for a set of Roles.
func NewAccountRolePermissionChecker(roles ...string) AccountRolePermissionsChecker {
	return &accountRoleCollection{
		Roles: roles,
	}
}

func (r AccountRole) String() string {
	switch r {
	case AccountMemberRole:
		return accountMemberRoleName
	case AccountAdminRole:
		return accountAdminRoleName
	default:
		return ""
	}
}

// HasPermission returns whether a user can do something or not.
func (r accountRoleCollection) HasPermission(p Permission) bool {
	return hasPermission(p, r.Roles...)
}

// CanUpdateAccounts returns whether a user can update accounts or not.
func (r accountRoleCollection) CanUpdateAccounts() bool {
	return hasPermission(UpdateAccountPermission, r.Roles...)
}

// CanDeleteAccounts returns whether a user can delete accounts or not.
func (r accountRoleCollection) CanDeleteAccounts() bool {
	return hasPermission(ArchiveAccountPermission, r.Roles...)
}

// CanAddMemberToAccounts returns whether a user can add members to accounts or not.
func (r accountRoleCollection) CanAddMemberToAccounts() bool {
	return hasPermission(AddMemberAccountPermission, r.Roles...)
}

// CanRemoveMemberFromAccounts returns whether a user can remove members from accounts or not.
func (r accountRoleCollection) CanRemoveMemberFromAccounts() bool {
	return hasPermission(RemoveMemberAccountPermission, r.Roles...)
}

// CanTransferAccountToNewOwner returns whether a user can transfer an account to a new owner or not.
func (r accountRoleCollection) CanTransferAccountToNewOwner() bool {
	return hasPermission(TransferAccountPermission, r.Roles...)
}

// CanCreateWebhooks returns whether a user can create webhooks or not.
func (r accountRoleCollection) CanCreateWebhooks() bool {
	return hasPermission(CreateWebhooksPermission, r.Roles...)
}

// CanSeeWebhooks returns whether a user can view webhooks or not.
func (r accountRoleCollection) CanSeeWebhooks() bool {
	return hasPermission(ReadWebhooksPermission, r.Roles...)
}

// CanUpdateWebhooks returns whether a user can update webhooks or not.
func (r accountRoleCollection) CanUpdateWebhooks() bool {
	return hasPermission(UpdateWebhooksPermission, r.Roles...)
}

// CanArchiveWebhooks returns whether a user can delete webhooks or not.
func (r accountRoleCollection) CanArchiveWebhooks() bool {
	return hasPermission(ArchiveWebhooksPermission, r.Roles...)
}

// CanCreateAPIClients returns whether a user can create API clients or not.
func (r accountRoleCollection) CanCreateAPIClients() bool {
	return hasPermission(CreateAPIClientsPermission, r.Roles...)
}

// CanSeeAPIClients returns whether a user can view API clients or not.
func (r accountRoleCollection) CanSeeAPIClients() bool {
	return hasPermission(ReadAPIClientsPermission, r.Roles...)
}

// CanDeleteAPIClients returns whether a user can delete API clients or not.
func (r accountRoleCollection) CanDeleteAPIClients() bool {
	return hasPermission(ArchiveAPIClientsPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForWebhooks returns whether a user can view webhook audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForWebhooks() bool {
	return hasPermission(ReadWebhooksAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidInstruments returns whether a user can view valid instrument audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForValidInstruments() bool {
	return hasPermission(ReadValidInstrumentsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidPreparations returns whether a user can view valid preparation audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForValidPreparations() bool {
	return hasPermission(ReadValidPreparationsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidIngredients returns whether a user can view valid ingredient audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForValidIngredients() bool {
	return hasPermission(ReadValidIngredientsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidIngredientPreparations returns whether a user can view valid ingredient preparation audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForValidIngredientPreparations() bool {
	return hasPermission(ReadValidIngredientPreparationsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidPreparationInstruments returns whether a user can view valid preparation instrument audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForValidPreparationInstruments() bool {
	return hasPermission(ReadValidPreparationInstrumentsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForRecipes returns whether a user can view recipe audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForRecipes() bool {
	return hasPermission(ReadRecipesAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForRecipeSteps returns whether a user can view recipe step audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForRecipeSteps() bool {
	return hasPermission(ReadRecipeStepsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForRecipeStepIngredients returns whether a user can view recipe step ingredient audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForRecipeStepIngredients() bool {
	return hasPermission(ReadRecipeStepIngredientsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForRecipeStepProducts returns whether a user can view recipe step product audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForRecipeStepProducts() bool {
	return hasPermission(ReadRecipeStepProductsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForInvitations returns whether a user can view invitation audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForInvitations() bool {
	return hasPermission(ReadInvitationsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForReports returns whether a user can view report audit log entries or not.
func (r accountRoleCollection) CanSeeAuditLogEntriesForReports() bool {
	return hasPermission(ReadReportsAuditLogEntriesPermission, r.Roles...)
}
