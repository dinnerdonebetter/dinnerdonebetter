package sendgrid

import (
	"context"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/circuitbreaking"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	name = "sendgrid_emailer"
)

var (
	_ email.Emailer = (*Emailer)(nil)
	// ErrNilConfig indicates a nil config was provided.
	ErrNilConfig = errors.New("SendGrid config is nil")
	// ErrEmptyAPIToken indicates an empty API token was provided.
	ErrEmptyAPIToken = errors.New("empty Sendgrid API token")
	// ErrNilHTTPClient indicates a nil HTTP client was provided.
	ErrNilHTTPClient = errors.New("nil HTTP client")
)

type (
	// Config configures SendGrid to send email.
	Config struct {
		APIToken string `json:"apiToken" toml:"api_token,omitempty"`
	}

	// Emailer uses SendGrid to send email.
	Emailer struct {
		logger         logging.Logger
		tracer         tracing.Tracer
		circuitBreaker circuitbreaking.CircuitBreaker
		client         *sendgrid.Client
		config         Config
	}
)

// NewSendGridEmailer returns a new SendGrid-backed Emailer.
func NewSendGridEmailer(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, client *http.Client, circuitBreaker circuitbreaking.CircuitBreaker) (*Emailer, error) {
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
		logger:         logging.EnsureLogger(logger).WithName(name),
		tracer:         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(name)),
		client:         sendgrid.NewSendClient(cfg.APIToken),
		config:         *cfg,
		circuitBreaker: circuitBreaker,
	}

	return e, nil
}

// ErrSendgridAPIResponse indicates an error occurred in SendGrid.
var ErrSendgridAPIResponse = errors.New("sendgrid request error")

// SendEmail sends an email.
func (e *Emailer) SendEmail(ctx context.Context, details *email.OutboundEmailMessage) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachToSpan(span, "to_email", details.ToAddress)

	if e.circuitBreaker.CannotProceed() {
		return types.ErrCircuitBroken
	}

	to := mail.NewEmail(details.ToName, details.ToAddress)
	from := mail.NewEmail(details.FromName, details.FromAddress)
	message := mail.NewSingleEmail(from, details.Subject, to, "", details.HTMLContent)

	res, err := e.client.SendWithContext(ctx, message)
	if err != nil {
		return observability.PrepareError(err, span, "sending email")
	}

	// Fun fact: if your account is limited and not able to send an email, there is
	// no distinguishing feature of the response to let you know. Thanks, SendGrid!
	if res.StatusCode != http.StatusAccepted {
		e.logger.WithValue("sendgrid_api_token", e.config.APIToken).Info("sending email yielded an invalid response")
		tracing.AttachToSpan(span, e.config.APIToken, "sendgrid_api_token")
		e.circuitBreaker.Failed()
		return observability.PrepareError(ErrSendgridAPIResponse, span, "sending email yielded a %d response", res.StatusCode)
	}

	e.circuitBreaker.Succeeded()
	return nil
}

func (e *Emailer) preparePersonalization(to *mail.Email, data map[string]any) *mail.Personalization {
	p := mail.NewPersonalization()
	p.AddTos(to)

	for k, v := range data {
		p.SetDynamicTemplateData(k, v)
	}

	return p
}

// sendDynamicTemplateEmail sends an email.
func (e *Emailer) sendDynamicTemplateEmail(ctx context.Context, to, from *mail.Email, templateID string, data map[string]any, request rest.Request) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachToSpan(span, "to_email", to.Address)

	m := mail.NewV3Mail()
	m.SetFrom(from).SetTemplateID(templateID).AddPersonalizations(e.preparePersonalization(to, data))

	request.Body = mail.GetRequestBody(m)

	res, err := sendgrid.MakeRequestWithContext(ctx, request)
	if err != nil {
		return observability.PrepareError(err, span, "sending dynamic email")
	}

	if res.StatusCode != http.StatusAccepted {
		return observability.PrepareError(ErrSendgridAPIResponse, span, "sending dynamic email yielded a %d response", res.StatusCode)
	}

	return nil
}
