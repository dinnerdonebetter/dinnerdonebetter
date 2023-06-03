package types

import (
	"context"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// UserFeedbackCreatedCustomerEventType indicates a user feedback was created.
	UserFeedbackCreatedCustomerEventType CustomerEventType = "user_feedback_created"
	// UserFeedbackArchivedCustomerEventType indicates a user feedback was archived.
	UserFeedbackArchivedCustomerEventType CustomerEventType = "user_feedback_archived"
)

type (
	// UserFeedback represents a user feedback listener, an endpoint to send an HTTP request to upon an event.
	UserFeedback struct {
		_         struct{}
		CreatedAt time.Time      `json:"createdAt"`
		Context   map[string]any `json:"context"`
		Prompt    string         `json:"prompt"`
		Feedback  string         `json:"feedback"`
		ByUser    string         `json:"byUser"`
		ID        string         `json:"id"`
		Rating    float32        `json:"rating"`
	}

	// UserFeedbackCreationRequestInput represents what a User could set as input for creating a user feedback.
	UserFeedbackCreationRequestInput struct {
		_        struct{}
		Context  map[string]any `json:"context"`
		Prompt   string         `json:"prompt"`
		Feedback string         `json:"feedback"`
		Rating   float32        `json:"rating"`
	}

	// UserFeedbackDatabaseCreationInput is used for creating a user feedback.
	UserFeedbackDatabaseCreationInput struct {
		_        struct{}
		Context  map[string]any
		ID       string
		Prompt   string
		Feedback string
		ByUser   string
		Rating   float32
	}

	// UserFeedbackDataManager describes a structure capable of storing user feedback.
	UserFeedbackDataManager interface {
		GetUserFeedbacks(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[UserFeedback], error)
		CreateUserFeedback(ctx context.Context, input *UserFeedbackDatabaseCreationInput) (*UserFeedback, error)
	}

	// UserFeedbackDataService describes a structure capable of serving traffic related to user feedback.
	UserFeedbackDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*UserFeedbackCreationRequestInput)(nil)

// ValidateWithContext validates a UserFeedbackCreationRequestInput.
func (w *UserFeedbackCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, w,
		validation.Field(&w.Prompt, validation.Required),
		validation.Field(&w.Feedback, validation.Required),
		validation.Field(&w.Rating, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*UserFeedbackDatabaseCreationInput)(nil)

// ValidateWithContext validates a UserFeedbackDatabaseCreationInput.
func (w *UserFeedbackDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, w,
		validation.Field(&w.ID, validation.Required),
		validation.Field(&w.Prompt, validation.Required),
		validation.Field(&w.Feedback, validation.Required),
		validation.Field(&w.Rating, validation.Required),
		validation.Field(&w.ByUser, validation.Required),
	)
}
