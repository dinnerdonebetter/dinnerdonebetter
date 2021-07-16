package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleRecipeStepID uint64 = 123
)

func TestBuildRecipeStepCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeStepCreationEventEntry(&types.RecipeStep{}, exampleAccountID))
}

func TestBuildRecipeStepUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeStepUpdateEventEntry(exampleUserID, exampleRecipeStepID, nil))
}

func TestBuildRecipeStepArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeStepArchiveEventEntry(exampleUserID, exampleRecipeStepID))
}
