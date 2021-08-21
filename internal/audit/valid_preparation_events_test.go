package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleValidPreparationID uint64 = 123
)

func TestBuildValidPreparationCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidPreparationCreationEventEntry(&types.ValidPreparation{}, exampleHouseholdID))
}

func TestBuildValidPreparationUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidPreparationUpdateEventEntry(exampleUserID, exampleValidPreparationID, nil))
}

func TestBuildValidPreparationArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidPreparationArchiveEventEntry(exampleUserID, exampleValidPreparationID))
}
