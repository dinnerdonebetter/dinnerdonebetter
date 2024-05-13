package sendgrid

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	circuit "github.com/rubyist/circuitbreaker"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	name = "sendgrid_emailer"
)

var (
	_ email.Emailer = (*Emailer)(nil)

	// ErrNilConfig indicates an empty API token was provided.
	ErrNilConfig = errors.New("SendGrid config is nil")
	// ErrNilHTTPClient indicates a nil HTTP client was provided.
	ErrNilHTTPClient = errors.New("nil HTTP client")
)

type (
	// Config configures SendGrid to send email.
	Config struct {
		APIToken                             string  `json:"apiToken"                             toml:"api_token,omitempty"`
		CircuitBreakerFailureRate            float64 `json:"circuitBreakerFailureRate"            toml:"circuit_breaker_failure_rate,omitempty"`
		CircuitBreakerFailureSampleThreshold int64   `json:"circuitBreakerFailureSampleThreshold" toml:"circuit_breaker_failure_sample_threshold,omitempty"`
	}

	// Emailer uses SendGrid to send email.
	Emailer struct {
		logger         logging.Logger
		tracer         tracing.Tracer
		client         *sendgrid.Client
		circuitBreaker *circuit.Breaker
		config         Config
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	cfg.EnsureDefaults()

	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.APIToken, validation.Required),
		validation.Field(&cfg.CircuitBreakerFailureRate, validation.Required),
		validation.Field(&cfg.CircuitBreakerFailureSampleThreshold, validation.Required),
	)
}

func (cfg *Config) EnsureDefaults() {
	if cfg.CircuitBreakerFailureRate == 0 {
		cfg.CircuitBreakerFailureRate = .5
	}

	if cfg.CircuitBreakerFailureSampleThreshold == 0 {
		cfg.CircuitBreakerFailureSampleThreshold = 10
	}
}

// NewSendGridEmailer returns a new SendGrid-backed Emailer.
func NewSendGridEmailer(ctx context.Context, cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, client *http.Client) (*Emailer, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	if err := cfg.ValidateWithContext(ctx); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if client == nil {
		return nil, ErrNilHTTPClient
	}

	// this line causes data races when the unit tests in this package are run in parallel.
	// that sucks, but I also basically can't do anything about it because of how SendGrid's dogshit client works.
	sendgrid.DefaultClient = &rest.Client{HTTPClient: client}

	e := &Emailer{
		logger: logging.EnsureLogger(logger).WithName(name),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(name)),
		client: sendgrid.NewSendClient(cfg.APIToken),
		config: *cfg,
		circuitBreaker: circuit.NewBreakerWithOptions(&circuit.Options{
			ShouldTrip: func(cb *circuit.Breaker) bool {
				samples := cb.Failures() + cb.Successes()
				return samples >= cfg.CircuitBreakerFailureSampleThreshold && cb.ErrorRate() >= cfg.CircuitBreakerFailureRate
			},
		}),
	}

	return e, nil
}

// ErrSendgridAPIIssue indicates an error occurred in SendGrid.
var ErrSendgridAPIIssue = errors.New("making SendGrid request")

// SendEmail sends an email.
func (e *Emailer) SendEmail(ctx context.Context, details *email.OutboundEmailMessage) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	logger := e.logger.WithValue("to_email", details.ToAddress).WithValue("subject", details.Subject)
	tracing.AttachToSpan(span, "to_email", details.ToAddress)
	tracing.AttachToSpan(span, e.config.APIToken, "sendgrid_api_token")

	to := mail.NewEmail(details.ToName, details.ToAddress)
	from := mail.NewEmail(details.FromName, details.FromAddress)
	message := mail.NewSingleEmail(from, details.Subject, to, "", details.HTMLContent)

	res, err := e.client.SendWithContext(ctx, message)
	if err != nil {
		e.circuitBreaker.Fail()
		return observability.PrepareError(err, span, "sending email")
	}

	// Fun fact: if your account is limited and not able to send an email, there is localhost:9000/accept_invitation?i=cgk80ta23akg00eo8pf0&t=DHriGX_38me7f7NY7yjHjrh47jU9RIbloCAQU4SJhRwDoHSI23dw0kaD2CBHepRnnpbmHxcPOohCC5t8foniGA
	// no distinguishing feature of the response to let you know. Thanks, SendGrid!
	if res.StatusCode != http.StatusAccepted {
		e.circuitBreaker.Fail()
		return observability.PrepareAndLogError(ErrSendgridAPIIssue, logger, span, "sending email yielded a %d response", res.StatusCode)
	}

	e.circuitBreaker.Success()
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
		return observability.PrepareError(ErrSendgridAPIIssue, span, "sending dynamic email yielded a %d response", res.StatusCode)
	}

	return nil
}
