// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestClient_GetAuditLogEntryByID(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/audit_log_entries/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		auditLogEntryID := fakes.BuildFakeID()

		data := fakes.BuildFakeAuditLogEntry()
		expected := &types.APIResponse[*types.AuditLogEntry]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, auditLogEntryID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetAuditLogEntryByID(ctx, auditLogEntryID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected.Data, actual)
	})

	T.Run("with invalid auditLogEntry ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetAuditLogEntryByID(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		auditLogEntryID := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetAuditLogEntryByID(ctx, auditLogEntryID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		auditLogEntryID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, auditLogEntryID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetAuditLogEntryByID(ctx, auditLogEntryID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
