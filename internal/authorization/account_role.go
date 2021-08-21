package authorization

import (
	"encoding/gob"

	"gopkg.in/mikespook/gorbac.v2"
)

type (
	// HouseholdRole describes a role a user has for an Household context.
	HouseholdRole role

	// HouseholdRolePermissionsChecker checks permissions for one or more household Roles.
	HouseholdRolePermissionsChecker interface {
		HasPermission(Permission) bool

		CanUpdateHouseholds() bool
		CanDeleteHouseholds() bool
		CanAddMemberToHouseholds() bool
		CanRemoveMemberFromHouseholds() bool
		CanTransferHouseholdToNewOwner() bool
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
	// HouseholdMemberRole is a role for a plain household participant.
	HouseholdMemberRole HouseholdRole = iota
	// HouseholdAdminRole is a role for someone who can manipulate the specifics of an household.
	HouseholdAdminRole HouseholdRole = iota

	householdAdminRoleName  = "household_admin"
	householdMemberRoleName = "household_member"
)

var (
	householdAdmin  = gorbac.NewStdRole(householdAdminRoleName)
	householdMember = gorbac.NewStdRole(householdMemberRoleName)
)

type householdRoleCollection struct {
	Roles []string
}

func init() {
	gob.Register(householdRoleCollection{})
}

// NewHouseholdRolePermissionChecker returns a new checker for a set of Roles.
func NewHouseholdRolePermissionChecker(roles ...string) HouseholdRolePermissionsChecker {
	return &householdRoleCollection{
		Roles: roles,
	}
}

func (r HouseholdRole) String() string {
	switch r {
	case HouseholdMemberRole:
		return householdMemberRoleName
	case HouseholdAdminRole:
		return householdAdminRoleName
	default:
		return ""
	}
}

// HasPermission returns whether a user can do something or not.
func (r householdRoleCollection) HasPermission(p Permission) bool {
	return hasPermission(p, r.Roles...)
}

// CanUpdateHouseholds returns whether a user can update households or not.
func (r householdRoleCollection) CanUpdateHouseholds() bool {
	return hasPermission(UpdateHouseholdPermission, r.Roles...)
}

// CanDeleteHouseholds returns whether a user can delete households or not.
func (r householdRoleCollection) CanDeleteHouseholds() bool {
	return hasPermission(ArchiveHouseholdPermission, r.Roles...)
}

// CanAddMemberToHouseholds returns whether a user can add members to households or not.
func (r householdRoleCollection) CanAddMemberToHouseholds() bool {
	return hasPermission(AddMemberHouseholdPermission, r.Roles...)
}

// CanRemoveMemberFromHouseholds returns whether a user can remove members from households or not.
func (r householdRoleCollection) CanRemoveMemberFromHouseholds() bool {
	return hasPermission(RemoveMemberHouseholdPermission, r.Roles...)
}

// CanTransferHouseholdToNewOwner returns whether a user can transfer an household to a new owner or not.
func (r householdRoleCollection) CanTransferHouseholdToNewOwner() bool {
	return hasPermission(TransferHouseholdPermission, r.Roles...)
}

// CanCreateWebhooks returns whether a user can create webhooks or not.
func (r householdRoleCollection) CanCreateWebhooks() bool {
	return hasPermission(CreateWebhooksPermission, r.Roles...)
}

// CanSeeWebhooks returns whether a user can view webhooks or not.
func (r householdRoleCollection) CanSeeWebhooks() bool {
	return hasPermission(ReadWebhooksPermission, r.Roles...)
}

// CanUpdateWebhooks returns whether a user can update webhooks or not.
func (r householdRoleCollection) CanUpdateWebhooks() bool {
	return hasPermission(UpdateWebhooksPermission, r.Roles...)
}

// CanArchiveWebhooks returns whether a user can delete webhooks or not.
func (r householdRoleCollection) CanArchiveWebhooks() bool {
	return hasPermission(ArchiveWebhooksPermission, r.Roles...)
}

// CanCreateAPIClients returns whether a user can create API clients or not.
func (r householdRoleCollection) CanCreateAPIClients() bool {
	return hasPermission(CreateAPIClientsPermission, r.Roles...)
}

// CanSeeAPIClients returns whether a user can view API clients or not.
func (r householdRoleCollection) CanSeeAPIClients() bool {
	return hasPermission(ReadAPIClientsPermission, r.Roles...)
}

// CanDeleteAPIClients returns whether a user can delete API clients or not.
func (r householdRoleCollection) CanDeleteAPIClients() bool {
	return hasPermission(ArchiveAPIClientsPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForWebhooks returns whether a user can view webhook audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForWebhooks() bool {
	return hasPermission(ReadWebhooksAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidInstruments returns whether a user can view valid instrument audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForValidInstruments() bool {
	return hasPermission(ReadValidInstrumentsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidPreparations returns whether a user can view valid preparation audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForValidPreparations() bool {
	return hasPermission(ReadValidPreparationsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidIngredients returns whether a user can view valid ingredient audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForValidIngredients() bool {
	return hasPermission(ReadValidIngredientsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidIngredientPreparations returns whether a user can view valid ingredient preparation audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForValidIngredientPreparations() bool {
	return hasPermission(ReadValidIngredientPreparationsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForValidPreparationInstruments returns whether a user can view valid preparation instrument audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForValidPreparationInstruments() bool {
	return hasPermission(ReadValidPreparationInstrumentsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForRecipes returns whether a user can view recipe audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForRecipes() bool {
	return hasPermission(ReadRecipesAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForRecipeSteps returns whether a user can view recipe step audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForRecipeSteps() bool {
	return hasPermission(ReadRecipeStepsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForRecipeStepIngredients returns whether a user can view recipe step ingredient audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForRecipeStepIngredients() bool {
	return hasPermission(ReadRecipeStepIngredientsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForRecipeStepProducts returns whether a user can view recipe step product audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForRecipeStepProducts() bool {
	return hasPermission(ReadRecipeStepProductsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForInvitations returns whether a user can view invitation audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForInvitations() bool {
	return hasPermission(ReadInvitationsAuditLogEntriesPermission, r.Roles...)
}

// CanSeeAuditLogEntriesForReports returns whether a user can view report audit log entries or not.
func (r householdRoleCollection) CanSeeAuditLogEntriesForReports() bool {
	return hasPermission(ReadReportsAuditLogEntriesPermission, r.Roles...)
}
