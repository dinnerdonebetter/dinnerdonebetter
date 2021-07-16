package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleRecipeStepIngredientID uint64 = 123
)

func TestBuildRecipeStepIngredientCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeStepIngredientCreationEventEntry(&types.RecipeStepIngredient{}, exampleAccountID))
}

func TestBuildRecipeStepIngredientUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeStepIngredientUpdateEventEntry(exampleUserID, exampleRecipeStepIngredientID, nil))
}

func TestBuildRecipeStepIngredientArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildRecipeStepIngredientArchiveEventEntry(exampleUserID, exampleRecipeStepIngredientID))
}
