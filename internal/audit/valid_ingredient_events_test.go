package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleValidIngredientID uint64 = 123
)

func TestBuildValidIngredientCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidIngredientCreationEventEntry(&types.ValidIngredient{}, exampleHouseholdID))
}

func TestBuildValidIngredientUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidIngredientUpdateEventEntry(exampleUserID, exampleValidIngredientID, nil))
}

func TestBuildValidIngredientArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidIngredientArchiveEventEntry(exampleUserID, exampleValidIngredientID))
}
