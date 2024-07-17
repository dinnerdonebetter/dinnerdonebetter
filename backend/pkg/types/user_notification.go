package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// UserNotificationCreatedCustomerEventType indicates a user notification was created.
	UserNotificationCreatedCustomerEventType ServiceEventType = "user_notification_created"
	// UserNotificationUpdatedCustomerEventType indicates a user notification was updated.
	UserNotificationUpdatedCustomerEventType ServiceEventType = "user_notification_updated"
	// UserNotificationArchivedCustomerEventType indicates a user notification was archived.
	UserNotificationArchivedCustomerEventType ServiceEventType = "user_notification_archived"

	// UserNotificationStatusTypeUnread represents the user notification status type for unread.
	UserNotificationStatusTypeUnread = "unread"
	// UserNotificationStatusTypeRead represents the user notification status type for read.
	UserNotificationStatusTypeRead = "read"
	// UserNotificationStatusTypeDismissed represents the user notification status type for dismissed.
	UserNotificationStatusTypeDismissed = "dismissed"
)

func init() {
	gob.Register(new(UserNotification))
	gob.Register(new(UserNotificationCreationRequestInput))
	gob.Register(new(UserNotificationUpdateRequestInput))
}

type (
	// UserNotification represents a user notification.
	UserNotification struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time  `json:"createdAt"`
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		ID            string     `json:"id"`
		Content       string     `json:"content"`
		Status        string     `json:"status"`
		BelongsToUser string     `json:"belongsToUser"`
	}

	// UserNotificationCreationRequestInput represents what a user could set as input for creating user notifications.
	UserNotificationCreationRequestInput struct {
		_ struct{} `json:"-"`

		Content       string `json:"content"`
		Status        string `json:"status"`
		BelongsToUser string `json:"belongsToUser"`
	}

	// UserNotificationDatabaseCreationInput represents what a user could set as input for creating user notifications.
	UserNotificationDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID            string
		Content       string
		BelongsToUser string
	}

	// UserNotificationUpdateRequestInput represents what a user could set as input for updating user notifications.
	UserNotificationUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Status *string `json:"status,omitempty"`
	}

	// UserNotificationDataManager describes a structure capable of storing user notifications permanently.
	UserNotificationDataManager interface {
		UserNotificationExists(ctx context.Context, userID, userNotificationID string) (bool, error)
		GetUserNotification(ctx context.Context, userID, userNotificationID string) (*UserNotification, error)
		GetUserNotifications(ctx context.Context, userID string, filter *QueryFilter) (*QueryFilteredResult[UserNotification], error)
		CreateUserNotification(ctx context.Context, input *UserNotificationDatabaseCreationInput) (*UserNotification, error)
		UpdateUserNotification(ctx context.Context, updated *UserNotification) error
	}

	// UserNotificationDataService describes a structure capable of serving traffic related to user notifications.
	UserNotificationDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an UserNotificationUpdateRequestInput with a user notification.
func (x *UserNotification) Update(input *UserNotificationUpdateRequestInput) {
	if input.Status != nil && *input.Status != x.Status {
		x.Status = *input.Status
	}
}

var _ validation.ValidatableWithContext = (*UserNotificationCreationRequestInput)(nil)

// ValidateWithContext validates a UserNotificationCreationRequestInput.
func (x *UserNotificationCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Content, validation.Required),
		validation.Field(&x.Status, validation.In(
			UserNotificationStatusTypeUnread,
			UserNotificationStatusTypeRead,
			UserNotificationStatusTypeDismissed,
		)),
	)
}

var _ validation.ValidatableWithContext = (*UserNotificationDatabaseCreationInput)(nil)

// ValidateWithContext validates a UserNotificationDatabaseCreationInput.
func (x *UserNotificationDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Content, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*UserNotificationUpdateRequestInput)(nil)

// ValidateWithContext validates a UserNotificationUpdateRequestInput.
func (x *UserNotificationUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Status, validation.Required, validation.In(
			UserNotificationStatusTypeUnread,
			UserNotificationStatusTypeRead,
			UserNotificationStatusTypeDismissed,
		)),
	)
}
