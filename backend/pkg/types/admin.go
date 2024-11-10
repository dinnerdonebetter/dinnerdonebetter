package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	QueueTypeDataChanges              = "data_changes"
	QueueTypeOutboundEmails           = "outbound_emails"
	QueueTypeSearchIndexRequests      = "search_index_requests"
	QueueTypeUserDataAggregator       = "user_data_aggregator"
	QueueTypeWebhookExecutionRequests = "webhook_execution_requests"
)

var (
	ValidQueueNames = []string{
		QueueTypeDataChanges,
		QueueTypeOutboundEmails,
		QueueTypeSearchIndexRequests,
		QueueTypeUserDataAggregator,
		QueueTypeWebhookExecutionRequests,
	}
)

type (
	// AdminDataService describes a structure capable of serving traffic related to users.
	AdminDataService interface {
		UserAccountStatusChangeHandler(http.ResponseWriter, *http.Request)
		WriteArbitraryQueueMessageHandler(http.ResponseWriter, *http.Request)
	}

	// UserAccountStatusUpdateInput represents what an admin User could provide as input for changing statuses.
	UserAccountStatusUpdateInput struct {
		_ struct{} `json:"-"`

		NewStatus    string `json:"newStatus"`
		Reason       string `json:"reason"`
		TargetUserID string `json:"targetUserID"`
	}

	ArbitraryQueueMessageRequestInput struct {
		_ struct{} `json:"-"`

		QueueName string `json:"queueName"`
		Body      string `json:"body"`
	}

	ArbitraryQueueMessageResponse struct {
		_ struct{} `json:"-"`

		Success bool `json:"success"`
	}

	// FrontendService serves static frontend files.
	FrontendService interface {
		StaticDir(ctx context.Context, staticFilesDirectory string) (http.HandlerFunc, error)
	}
)

var _ validation.ValidatableWithContext = (*UserAccountStatusUpdateInput)(nil)

// ValidateWithContext ensures our struct is validatable.
func (i *UserAccountStatusUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.NewStatus, validation.Required),
		validation.Field(&i.Reason, validation.Required),
		validation.Field(&i.TargetUserID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ArbitraryQueueMessageRequestInput)(nil)

// ValidateWithContext ensures our struct is validatable.
func (i *ArbitraryQueueMessageRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.QueueName, validation.In(
			QueueTypeDataChanges,
			QueueTypeOutboundEmails,
			QueueTypeSearchIndexRequests,
			QueueTypeUserDataAggregator,
			QueueTypeWebhookExecutionRequests,
		)),
		validation.Field(&i.Body, validation.Required),
	)
}
