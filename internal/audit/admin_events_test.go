package audit_test

import (
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/audit"

	"github.com/stretchr/testify/assert"
)

func TestBuildUserBanEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildUserBanEventEntry(exampleUserID, exampleUserID, "reason"))
}

func TestBuildHouseholdTerminationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildHouseholdTerminationEventEntry(exampleUserID, exampleUserID, "reason"))
}
