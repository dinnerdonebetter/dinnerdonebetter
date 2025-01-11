package apiclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
)

const (
	codeKey = "code"
)

var (
	ErrCodeNotReturned = errors.New("code not returned")
)

type ClientOption func(*Client) error

// SetOptions sets a new ClientOption on the client.
func (c *Client) SetOptions(opts ...ClientOption) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}

	return nil
}

// UsingURL sets the url on the client.
func UsingURL(u string) func(*Client) error {
	return func(c *Client) error {
		parsed, err := url.Parse(u)
		if err != nil {
			return err
		}

		c.url = parsed

		return nil
	}
}

// UsingTracingProvider sets the tracing provider on the client.
func UsingTracingProvider(tracerProvider tracing.TracerProvider) ClientOption {
	return func(c *Client) error {
		c.tracer = tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(clientName))

		return nil
	}
}

// UsingJSON sets the content type on the client.
func UsingJSON() func(*Client) error {
	return func(c *Client) error {
		c.encoder = encoding.ProvideClientEncoder(c.logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		return nil
	}
}

// UsingXML sets the content type on the client.
func UsingXML() func(*Client) error {
	return func(c *Client) error {
		c.encoder = encoding.ProvideClientEncoder(c.logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeXML)

		return nil
	}
}

// UsingLogger sets the logger on the client.
func UsingLogger(logger logging.Logger) func(*Client) error {
	return func(c *Client) error {
		c.logger = logging.EnsureLogger(logger)

		return nil
	}
}

// UsingDebug sets the debug value on the client.
func UsingDebug(debug bool) func(*Client) error {
	return func(c *Client) error {
		c.debug = debug
		return nil
	}
}

// UsingTimeout sets the debug value on the client.
func UsingTimeout(timeout time.Duration) func(*Client) error {
	return func(c *Client) error {
		if timeout == 0 {
			timeout = defaultTimeout
		}

		c.authedClient.Timeout = timeout
		c.unauthenticatedClient.Timeout = timeout

		return nil
	}
}

// UsingOAuth2 sets the client to use OAuth2.
func UsingOAuth2(ctx context.Context, clientID, clientSecret string, scopes []string, token string) func(*Client) error {
	return func(c *Client) error {
		state, err := random.GenerateBase64EncodedString(ctx, 32)
		if err != nil {
			return fmt.Errorf("generating state: %w", err)
		}

		oauth2Config := oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			RedirectURL:  c.url.String(),
			Endpoint: oauth2.Endpoint{
				AuthStyle: oauth2.AuthStyleInParams,
				AuthURL:   c.URL().String() + "/oauth2/authorize",
				TokenURL:  c.URL().String() + "/oauth2/token",
			},
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			oauth2Config.AuthCodeURL(
				state,
				oauth2.SetAuthURLParam("code_challenge_method", "plain"),
			),
			http.NoBody,
		)
		if err != nil {
			return fmt.Errorf("failed to get oauth2 code: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		client := otelhttp.DefaultClient
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}

		res, err := otelhttp.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to get oauth2 code: %w", err)
		}
		defer func() {
			if closeErr := res.Body.Close(); closeErr != nil {
				c.logger.Error("failed to close oauth2 response body", err)
			}
		}()

		rl, err := res.Location()
		if err != nil {
			return err
		}

		code := rl.Query().Get(codeKey)
		if code == "" {
			return ErrCodeNotReturned
		}

		oauth2Token, err := oauth2Config.Exchange(ctx, code)
		if err != nil {
			return err
		}

		c.authedClient.Transport = &oauth2.Transport{
			Source: oauth2.ReuseTokenSource(oauth2Token, oauth2.StaticTokenSource(oauth2Token)),
			Base:   otelhttp.DefaultClient.Transport,
		}

		c.authedClient = buildRetryingClient(c.authedClient, c.logger, c.tracer)

		return nil
	}
}

// ImpersonatingUser sets the impersonatedUserID value on the client.
func ImpersonatingUser(userID string) func(*Client) error {
	return func(c *Client) error {
		c.impersonatedUserID = userID

		return nil
	}
}

// ImpersonatingHousehold sets the impersonatedHouseholdID value on the client.
func ImpersonatingHousehold(householdID string) func(*Client) error {
	return func(c *Client) error {
		c.impersonatedHouseholdID = householdID

		return nil
	}
}

// WithoutImpersonating clears the impersonatedUserID and impersonatedHouseholdID value on the client.
func WithoutImpersonating() func(*Client) error {
	return func(c *Client) error {
		c.impersonatedUserID = ""
		c.impersonatedHouseholdID = ""

		return nil
	}
}
