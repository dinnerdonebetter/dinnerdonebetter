package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleInvitationID uint64 = 123
)

func TestBuildInvitationCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildInvitationCreationEventEntry(&types.Invitation{}, exampleAccountID))
}

func TestBuildInvitationUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildInvitationUpdateEventEntry(exampleUserID, exampleInvitationID, exampleAccountID, nil))
}

func TestBuildInvitationArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildInvitationArchiveEventEntry(exampleUserID, exampleInvitationID, exampleAccountID))
}
