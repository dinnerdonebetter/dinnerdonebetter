package sendgrid

import (
	"context"
	"errors"
	"net/http"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	name = "sendgrid_emailer"
)

var (
	_ email.Emailer = (*Emailer)(nil)
	// ErrNilConfig indicates an empty API token was provided.
	ErrNilConfig = errors.New("SendGrid config is nil")
	// ErrEmptyAPIToken indicates an empty API token was provided.
	ErrEmptyAPIToken = errors.New("empty API token")
	// ErrNilHTTPClient indicates a nil HTTP client was provided.
	ErrNilHTTPClient = errors.New("nil HTTP client")
)

type (
	// Config configures SendGrid to send email.
	Config struct {
		APIToken  string `json:"apiToken" mapstructure:"api_token" toml:"api_token,omitempty"`
		WebAppURL string `json:"webAppURL" mapstructure:"web_app_url" toml:"web_app_url,omitempty"`
	}

	// Emailer uses SendGrid to send email.
	Emailer struct {
		logger logging.Logger
		tracer tracing.Tracer
		client *sendgrid.Client
		config Config
	}
)

// NewSendGridEmailer returns a new SendGrid-backed Emailer.
func NewSendGridEmailer(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, client *http.Client) (*Emailer, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	if cfg.APIToken == "" {
		return nil, ErrEmptyAPIToken
	}

	if client == nil {
		return nil, ErrNilHTTPClient
	}

	// this line causes data races when the unit tests in this package are run in parallel.
	// that sucks, but I also basically can't do anything about it because of how SendGrid's dogshit client works.
	sendgrid.DefaultClient = &rest.Client{HTTPClient: client}

	e := &Emailer{
		logger: logging.EnsureLogger(logger).WithName(name),
		tracer: tracing.NewTracer(tracerProvider.Tracer(name)),
		client: sendgrid.NewSendClient(cfg.APIToken),
		config: *cfg,
	}

	return e, nil
}

// ErrSendgridAPIIssue indicates an error occurred in SendGrid.
var ErrSendgridAPIIssue = errors.New("making SendGrid request")

// SendEmail sends an email.
func (e *Emailer) SendEmail(ctx context.Context, details *email.OutboundEmailMessage) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachStringToSpan(span, "to_email", details.ToAddress)

	to := mail.NewEmail(details.ToName, details.ToAddress)
	from := mail.NewEmail(details.FromName, details.FromAddress)
	message := mail.NewSingleEmail(from, details.Subject, to, "", details.HTMLContent)

	res, err := e.client.SendWithContext(ctx, message)
	if err != nil {
		return observability.PrepareError(err, span, "sending email")
	}

	if res.StatusCode != http.StatusAccepted {
		return observability.PrepareError(ErrSendgridAPIIssue, span, "sending email yielded a %d response", res.StatusCode)
	}

	return nil
}

func (e *Emailer) preparePersonalization(to *mail.Email, data map[string]interface{}) *mail.Personalization {
	p := mail.NewPersonalization()
	p.AddTos(to)

	for k, v := range data {
		p.SetDynamicTemplateData(k, v)
	}

	return p
}

// sendDynamicTemplateEmail sends an email.
func (e *Emailer) sendDynamicTemplateEmail(ctx context.Context, to, from *mail.Email, templateID string, data map[string]interface{}, request rest.Request) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachStringToSpan(span, "to_email", to.Address)

	m := mail.NewV3Mail()
	m.SetFrom(from).SetTemplateID(templateID).AddPersonalizations(e.preparePersonalization(to, data))

	request.Body = mail.GetRequestBody(m)

	res, err := sendgrid.MakeRequestWithContext(ctx, request)
	if err != nil {
		return observability.PrepareError(err, span, "sending dynamic email")
	}

	if res.StatusCode != http.StatusAccepted {
		return observability.PrepareError(ErrSendgridAPIIssue, span, "sending dynamic email yielded a %d response", res.StatusCode)
	}

	return nil
}
