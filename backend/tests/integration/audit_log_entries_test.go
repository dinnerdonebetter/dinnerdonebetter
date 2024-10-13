package integration

import (
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/stretchr/testify/assert"
)

func (s *TestSuite) TestAuditLogEntries_Listing() {
	s.runTest("should be able to be read in a list", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.userClient.GetAuditLogEntriesForUser(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)

			assert.Equal(t, 4, len(actual.Data))

			actual, err = testClients.userClient.GetAuditLogEntriesForHousehold(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)

			assert.Equal(t, 2, len(actual.Data))
		}
	})
}

func (s *TestSuite) TestWebhooks_Retrieving_Returns404ForNonexistentAuditLogEntry() {
	s.runTest("should error when archiving a non-existent auditLogEntry", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.userClient.GetAuditLogEntryByID(ctx, nonexistentID)
			assert.Nil(t, actual)
			assert.Error(t, err)
		}
	})
}
