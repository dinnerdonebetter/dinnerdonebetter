package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleRecipeID uint64 = 123
)

func TestBuildRecipeCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeCreationEventEntry(&types.Recipe{}, exampleAccountID))
}

func TestBuildRecipeUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeUpdateEventEntry(exampleUserID, exampleRecipeID, exampleAccountID, nil))
}

func TestBuildRecipeArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeArchiveEventEntry(exampleUserID, exampleRecipeID, exampleAccountID))
}
