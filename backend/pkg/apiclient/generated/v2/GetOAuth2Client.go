// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient




import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
)


func (c *Client) GetOAuth2Client(
	ctx context.Context,
oauth2ClientID string,
) ( *types.OAuth2Client, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if oauth2ClientID == "" {
		return nil, buildInvalidIDError("oauth2Client")
	} 
	logger = logger.WithValue(keys.OAuth2ClientIDKey, oauth2ClientID)
	tracing.AttachToSpan(span, keys.OAuth2ClientIDKey, oauth2ClientID)

 

	u := c.BuildURL(ctx, nil, fmt.Sprintf("/api/v1/oauth2_clients/%s" , oauth2ClientID ))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building request to fetch a OAuth2Client")
	}

	var apiResponse *types.APIResponse[  *types.OAuth2Client]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading OAuth2Client response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}


	return apiResponse.Data, nil
}