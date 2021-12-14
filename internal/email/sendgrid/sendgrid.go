package sendgrid

import (
	"context"
	"errors"
	"net/http"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	name = "sendgrid_emailer"
)

var _ email.Emailer = (*Emailer)(nil)

var (
	// ErrEmptyAPIToken indicates an empty API token was provided.
	ErrEmptyAPIToken = errors.New("empty API token")
	// ErrNilHTTPClient indicates a nil HTTP client was provided.
	ErrNilHTTPClient = errors.New("nil HTTP client")
)

type (
	// Emailer uses SendGrid to send email.
	Emailer struct {
		logger logging.Logger
		tracer tracing.Tracer
		client *sendgrid.Client
	}
)

// NewSendGridEmailer returns a new SendGrid-backed Emailer.
func NewSendGridEmailer(apiToken string, logger logging.Logger, tracerProvider trace.TracerProvider, client *http.Client) (*Emailer, error) {
	if apiToken == "" {
		return nil, ErrEmptyAPIToken
	}

	if client == nil {
		return nil, ErrNilHTTPClient
	}

	sendgrid.DefaultClient = &rest.Client{HTTPClient: client}
	c := sendgrid.NewSendClient(apiToken)

	e := &Emailer{
		logger: logging.EnsureLogger(logger).WithName(name),
		tracer: tracing.NewTracer(tracerProvider.Tracer(name)),
		client: c,
	}

	return e, nil
}

// ErrSendgridAPIIssue indicates an error occurred in SendGrid.
var ErrSendgridAPIIssue = errors.New("making SendGrid request")

// SendEmail sends an email.
func (e *Emailer) SendEmail(ctx context.Context, details *email.OutboundMessageDetails) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	logger := e.logger.WithValue("to_email", details.ToAddress)

	to := mail.NewEmail(details.ToName, details.ToAddress)
	from := mail.NewEmail(details.FromName, details.FromAddress)
	message := mail.NewSingleEmail(from, details.Subject, to, "", details.HTMLContent)
	res, err := e.client.SendWithContext(ctx, message)
	if err != nil {
		return observability.PrepareError(err, logger, span, "sending email")
	}

	if res.StatusCode != http.StatusOK {
		return observability.PrepareError(ErrSendgridAPIIssue, logger, span, "sending email yielded a %d response", res.StatusCode)
	}

	return nil
}
