package audit_test

import (
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuildUserAddedToHouseholdEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildUserAddedToHouseholdEventEntry(exampleAdminUserID, &types.AddUserToHouseholdInput{Reason: t.Name()}))
}

func TestBuildUserRemovedFromHouseholdEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildUserRemovedFromHouseholdEventEntry(exampleAdminUserID, exampleUserID, exampleHouseholdID, "blah blah"))
}

func TestBuildUserMarkedHouseholdAsDefaultEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildUserMarkedHouseholdAsDefaultEventEntry(exampleAdminUserID, exampleUserID, exampleHouseholdID))
}

func TestBuildModifyUserPermissionsEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildModifyUserPermissionsEventEntry(exampleUserID, exampleHouseholdID, exampleAdminUserID, []string{t.Name()}, t.Name()))
}

func TestBuildTransferHouseholdOwnershipEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildTransferHouseholdOwnershipEventEntry(exampleHouseholdID, exampleAdminUserID, fakes.BuildFakeTransferHouseholdOwnershipInput()))
}
