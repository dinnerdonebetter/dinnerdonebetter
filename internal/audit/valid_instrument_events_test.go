package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleValidInstrumentID uint64 = 123
)

func TestBuildValidInstrumentCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidInstrumentCreationEventEntry(&types.ValidInstrument{}, exampleAccountID))
}

func TestBuildValidInstrumentUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidInstrumentUpdateEventEntry(exampleUserID, exampleValidInstrumentID, nil))
}

func TestBuildValidInstrumentArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildValidInstrumentArchiveEventEntry(exampleUserID, exampleValidInstrumentID))
}
