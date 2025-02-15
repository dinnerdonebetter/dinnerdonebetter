package mailgun

import (
	"context"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/lib/email"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/mailgun/mailgun-go/v4"
)

const (
	name = "mailgun_emailer"
)

var (
	_ email.Emailer = (*Emailer)(nil)

	// ErrNilConfig indicates a nil config was provided.
	ErrNilConfig = errors.New("mailgun config is nil")
	// ErrEmptyDomain indicates an empty domain was provided.
	ErrEmptyDomain = errors.New("empty domain")
	// ErrEmptyPrivateAPIKey indicates an empty API token was provided.
	ErrEmptyPrivateAPIKey = errors.New("empty Mailgun API token")
	// ErrNilHTTPClient indicates a nil HTTP client was provided.
	ErrNilHTTPClient = errors.New("nil HTTP client")
)

type (
	// Emailer uses Mailgun to send email.
	Emailer struct {
		logger         logging.Logger
		tracer         tracing.Tracer
		client         mailgun.Mailgun
		circuitBreaker circuitbreaking.CircuitBreaker
	}
)

// NewMailgunEmailer returns a new Mailgun-backed Emailer.
func NewMailgunEmailer(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, client *http.Client, circuitBreaker circuitbreaking.CircuitBreaker) (*Emailer, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	if cfg.Domain == "" {
		return nil, ErrEmptyDomain
	}

	if cfg.PrivateAPIKey == "" {
		return nil, ErrEmptyPrivateAPIKey
	}

	if client == nil {
		return nil, ErrNilHTTPClient
	}

	mg := mailgun.NewMailgun(cfg.Domain, cfg.PrivateAPIKey)
	mg.SetClient(client)

	e := &Emailer{
		logger:         logging.EnsureLogger(logger).WithName(name),
		tracer:         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(name)),
		client:         mg,
		circuitBreaker: circuitBreaker,
	}

	return e, nil
}

// SendEmail sends an email.
func (e *Emailer) SendEmail(ctx context.Context, details *email.OutboundEmailMessage) error {
	ctx, span := e.tracer.StartSpan(ctx)
	defer span.End()

	logger := e.logger.WithValue("email.subject", details.Subject).WithValue("email.to_address", details.ToAddress)
	tracing.AttachToSpan(span, "to_email", details.ToAddress)

	if e.circuitBreaker.CannotProceed() {
		return internalerrors.ErrCircuitBroken
	}

	msg := mailgun.NewMessage(details.FromName, details.Subject, details.HTMLContent, details.ToAddress)
	if _, _, err := e.client.Send(ctx, msg); err != nil {
		e.circuitBreaker.Failed()
		return observability.PrepareAndLogError(err, logger, span, "sending email")
	}

	e.circuitBreaker.Succeeded()
	return nil
}
