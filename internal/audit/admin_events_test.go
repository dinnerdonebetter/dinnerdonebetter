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

func TestBuildAccountTerminationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildAccountTerminationEventEntry(exampleUserID, exampleUserID, "reason"))
}
