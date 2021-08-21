package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// UserAddedToHouseholdEvent events indicate a user created a membership.
	UserAddedToHouseholdEvent = "user_added_to_household"
	// UserHouseholdPermissionsModifiedEvent events indicate a user created a membership.
	UserHouseholdPermissionsModifiedEvent = "user_household_permissions_modified"
	// UserRemovedFromHouseholdEvent events indicate a user deleted a membership.
	UserRemovedFromHouseholdEvent = "user_removed_from_household"
	// HouseholdMarkedAsDefaultEvent events indicate a user deleted a membership.
	HouseholdMarkedAsDefaultEvent = "household_marked_as_default"
	// HouseholdTransferredEvent events indicate a user deleted a membership.
	HouseholdTransferredEvent = "household_transferred"
)

// BuildUserAddedToHouseholdEventEntry builds an entry creation input for when a membership is created.
func BuildUserAddedToHouseholdEventEntry(addedBy uint64, input *types.AddUserToHouseholdInput) *types.AuditLogEntryCreationInput {
	contextMap := map[string]interface{}{
		ActorAssignmentKey:     addedBy,
		HouseholdAssignmentKey: input.HouseholdID,
		UserAssignmentKey:      input.UserID,
	}

	if input.Reason != "" {
		contextMap[ReasonKey] = input.Reason
	}

	return &types.AuditLogEntryCreationInput{
		EventType: UserAddedToHouseholdEvent,
		Context:   contextMap,
	}
}

// BuildUserRemovedFromHouseholdEventEntry builds an entry creation input for when a membership is archived.
func BuildUserRemovedFromHouseholdEventEntry(removedBy, removed, householdID uint64, reason string) *types.AuditLogEntryCreationInput {
	contextMap := map[string]interface{}{
		ActorAssignmentKey:     removedBy,
		HouseholdAssignmentKey: householdID,
		UserAssignmentKey:      removed,
	}

	if reason != "" {
		contextMap[ReasonKey] = reason
	}

	return &types.AuditLogEntryCreationInput{
		EventType: UserRemovedFromHouseholdEvent,
		Context:   contextMap,
	}
}

// BuildUserMarkedHouseholdAsDefaultEventEntry builds an entry creation input for when a membership is created.
func BuildUserMarkedHouseholdAsDefaultEventEntry(performedBy, userID, householdID uint64) *types.AuditLogEntryCreationInput {
	contextMap := map[string]interface{}{
		ActorAssignmentKey:     performedBy,
		UserAssignmentKey:      userID,
		HouseholdAssignmentKey: householdID,
	}

	return &types.AuditLogEntryCreationInput{
		EventType: HouseholdMarkedAsDefaultEvent,
		Context:   contextMap,
	}
}

// BuildModifyUserPermissionsEventEntry builds an entry creation input for when a membership is created.
func BuildModifyUserPermissionsEventEntry(userID, householdID, modifiedBy uint64, newRoles []string, reason string) *types.AuditLogEntryCreationInput {
	contextMap := map[string]interface{}{
		ActorAssignmentKey:     modifiedBy,
		HouseholdAssignmentKey: householdID,
		UserAssignmentKey:      userID,
		HouseholdRolesKey:      newRoles,
	}

	if reason != "" {
		contextMap[ReasonKey] = reason
	}

	return &types.AuditLogEntryCreationInput{
		EventType: UserHouseholdPermissionsModifiedEvent,
		Context:   contextMap,
	}
}

// BuildTransferHouseholdOwnershipEventEntry builds an entry creation input for when a membership is created.
func BuildTransferHouseholdOwnershipEventEntry(householdID, changedBy uint64, input *types.HouseholdOwnershipTransferInput) *types.AuditLogEntryCreationInput {
	contextMap := map[string]interface{}{
		ActorAssignmentKey:     changedBy,
		"old_owner":            input.CurrentOwner,
		"new_owner":            input.NewOwner,
		HouseholdAssignmentKey: householdID,
	}

	if input.Reason != "" {
		contextMap[ReasonKey] = input.Reason
	}

	return &types.AuditLogEntryCreationInput{
		EventType: HouseholdTransferredEvent,
		Context:   contextMap,
	}
}
