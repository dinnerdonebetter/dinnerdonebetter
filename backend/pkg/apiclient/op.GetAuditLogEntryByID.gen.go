// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

func (c *Client) GetAuditLogEntryByID(
	ctx context.Context,
	auditLogEntryID string,
	reqMods ...RequestModifier,
) (*AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if auditLogEntryID == "" {
		return nil, buildInvalidIDError("auditLogEntry")
	}
	logger = logger.WithValue(keys.AuditLogEntryIDKey, auditLogEntryID)
	tracing.AttachToSpan(span, keys.AuditLogEntryIDKey, auditLogEntryID)

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/audit_log_entries/%s", auditLogEntryID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a AuditLogEntry")
	}

	for _, mod := range reqMods {
		mod(req)
	}

	var apiResponse *APIResponse[*AuditLogEntry]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading AuditLogEntry response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}
