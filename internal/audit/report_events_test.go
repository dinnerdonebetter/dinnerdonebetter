package audit_test

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/assert"
)

const (
	exampleReportID uint64 = 123
)

func TestBuildReportCreationEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildReportCreationEventEntry(&types.Report{}, exampleAccountID))
}

func TestBuildReportUpdateEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildReportUpdateEventEntry(exampleUserID, exampleReportID, exampleAccountID, nil))
}

func TestBuildReportArchiveEventEntry(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, audit.BuildReportArchiveEventEntry(exampleUserID, exampleReportID, exampleAccountID))
}
