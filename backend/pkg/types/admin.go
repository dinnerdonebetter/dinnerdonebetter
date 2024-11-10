package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	dataChangesQueueName              = "data_changes"
	outboundEmailsQueueName           = "outbound_emails"
	searchIndexRequestsQueueName      = "search_index_requests"
	userDataAggregatorQueueName       = "user_data_aggregator"
	webhookExecutionRequestsQueueName = "webhook_execution_requests"
)

var (
	ValidQueueNames = []string{
		dataChangesQueueName,
		outboundEmailsQueueName,
		searchIndexRequestsQueueName,
		userDataAggregatorQueueName,
		webhookExecutionRequestsQueueName,
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
			dataChangesQueueName,
			outboundEmailsQueueName,
			searchIndexRequestsQueueName,
			userDataAggregatorQueueName,
			webhookExecutionRequestsQueueName,
		)),
		validation.Field(&i.Body, validation.Required),
	)
}
