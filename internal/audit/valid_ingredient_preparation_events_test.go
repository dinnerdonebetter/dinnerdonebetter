package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleValidIngredientPreparationID uint64 = 123
)

func TestBuildValidIngredientPreparationCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidIngredientPreparationCreationEventEntry(&types.ValidIngredientPreparation{}, exampleAccountID))
}

func TestBuildValidIngredientPreparationUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidIngredientPreparationUpdateEventEntry(exampleUserID, exampleValidIngredientPreparationID, nil))
}

func TestBuildValidIngredientPreparationArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidIngredientPreparationArchiveEventEntry(exampleUserID, exampleValidIngredientPreparationID))
}
