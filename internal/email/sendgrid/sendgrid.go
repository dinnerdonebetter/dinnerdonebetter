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
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	name = "sendgrid_emailer"
)

var (
	_ email.Emailer = (*Emailer)(nil)
	// ErrEmptyAPIToken indicates an empty API token was provided.
	ErrEmptyAPIToken = errors.New("empty API token")
	// ErrNilHTTPClient indicates a nil HTTP client was provided.
	ErrNilHTTPClient = errors.New("nil HTTP client")
)

type (
	Config struct {
		APIToken                            string `json:"apiToken" mapstructure:"api_token" toml:"api_token,omitempty"`
		WebAppURL                           string `json:"webAppURL" mapstructure:"web_app_url" toml:"web_app_url,omitempty"`
		HouseholdInviteOutboundEmailAddress string `json:"householdInviteOutboundEmailAddress" mapstructure:"household_invitation_outbound_email_address" toml:"household_invitation_outbound_email_address,omitempty"`
		HouseholdInviteTemplateID           string `json:"householdInviteTemplateID" mapstructure:"household_invite_template_id" toml:"household_invite_template_id,omitempty"`
	}

	// Emailer uses SendGrid to send email.
	Emailer struct {
		logger logging.Logger
		tracer tracing.Tracer
		config Config
		client *sendgrid.Client
	}
)

// NewSendGridEmailer returns a new SendGrid-backed Emailer.
func NewSendGridEmailer(cfg Config, logger logging.Logger, tracerProvider tracing.TracerProvider, client *http.Client) (*Emailer, error) {
	if cfg.APIToken == "" {
		return nil, ErrEmptyAPIToken
	}

	if client == nil {
		return nil, ErrNilHTTPClient
	}

	sendgrid.DefaultClient = &rest.Client{HTTPClient: client}
	c := sendgrid.NewSendClient(cfg.APIToken)

	e := &Emailer{
		logger: logging.EnsureLogger(logger).WithName(name),
		tracer: tracing.NewTracer(tracerProvider.Tracer(name)),
		client: c,
		config: cfg,
	}

	return e, nil
}

// ErrSendgridAPIIssue indicates an error occurred in SendGrid.
var ErrSendgridAPIIssue = errors.New("making SendGrid request")

// SendEmail sends an email.
func (e *Emailer) SendEmail(ctx context.Context, details *email.OutboundEmailMessage) error {
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

	if res.StatusCode != http.StatusAccepted {
		return observability.PrepareError(ErrSendgridAPIIssue, logger, span, "sending email yielded a %d response", res.StatusCode)
	}

	return nil
}

func preparePersonalization(to *mail.Email, data map[string]string) *mail.Personalization {
	p := mail.NewPersonalization()
	p.AddTos(to)

	for k, v := range data {
		p.SetDynamicTemplateData(k, v)
	}

	return p
}

// SendHouseholdInvitationEmail sends an email.
func (e *Emailer) SendHouseholdInvitationEmail(ctx context.Context, householdInvitation *types.HouseholdInvitation) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	to := mail.NewEmail("PrixFixe Invitee", householdInvitation.ToEmail)
	from := mail.NewEmail("Example Test", e.config.HouseholdInviteOutboundEmailAddress)

	templateData := map[string]string{
		"webAppURL":    e.config.WebAppURL,
		"invitationID": householdInvitation.ID,
		"token":        householdInvitation.Token,
	}

	request := sendgrid.GetRequest(e.config.APIToken, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = http.MethodPost

	return e.sendDynamicTemplateEmail(ctx, to, from, e.config.HouseholdInviteTemplateID, templateData, request)
}

// sendDynamicTemplateEmail sends an email.
func (e *Emailer) sendDynamicTemplateEmail(ctx context.Context, to, from *mail.Email, templateID string, data map[string]string, request rest.Request) error {
	_, span := e.tracer.StartSpan(ctx)
	defer span.End()

	logger := e.logger.WithValue("to_email", to.Address)

	m := mail.NewV3Mail()
	m.SetFrom(from).SetTemplateID(templateID).AddPersonalizations(preparePersonalization(to, data))

	request.Body = mail.GetRequestBody(m)

	res, err := sendgrid.MakeRequestWithContext(ctx, request)
	if err != nil {
		return observability.PrepareError(err, logger, span, "sending dynamic email")
	}

	if res.StatusCode != http.StatusAccepted {
		return observability.PrepareError(ErrSendgridAPIIssue, logger, span, "sending dynamic email yielded a %d response", res.StatusCode)
	}

	return nil
}
