package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleValidPreparationInstrumentID uint64 = 123
)

func TestBuildValidPreparationInstrumentCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidPreparationInstrumentCreationEventEntry(&types.ValidPreparationInstrument{}, exampleAccountID))
}

func TestBuildValidPreparationInstrumentUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidPreparationInstrumentUpdateEventEntry(exampleUserID, exampleValidPreparationInstrumentID, nil))
}

func TestBuildValidPreparationInstrumentArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidPreparationInstrumentArchiveEventEntry(exampleUserID, exampleValidPreparationInstrumentID))
}
