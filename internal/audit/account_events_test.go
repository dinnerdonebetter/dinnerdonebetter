package audit_test

import (
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleHouseholdID uint64 = 123
)

func TestBuildHouseholdCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildHouseholdCreationEventEntry(&types.Household{}, exampleUserID))
}

func TestBuildHouseholdUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildHouseholdUpdateEventEntry(exampleUserID, exampleHouseholdID, exampleUserID, nil))
}

func TestBuildHouseholdArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildHouseholdArchiveEventEntry(exampleUserID, exampleHouseholdID, exampleUserID))
}
