package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleRecipeStepProductID uint64 = 123
)

func TestBuildRecipeStepProductCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeStepProductCreationEventEntry(&types.RecipeStepProduct{}, exampleHouseholdID))
}

func TestBuildRecipeStepProductUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeStepProductUpdateEventEntry(exampleUserID, exampleRecipeStepProductID, nil))
}

func TestBuildRecipeStepProductArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeStepProductArchiveEventEntry(exampleUserID, exampleRecipeStepProductID))
}
