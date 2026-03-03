package notifications

import "context"

// PushNotificationSender sends push notifications to device tokens.
// APNs/FCM integration is out of scope for the initial implementation;
// callers may use a noop implementation until a real provider is wired.
type PushNotificationSender interface {
	// SendPush sends a push notification to the given device tokens.
	// The title and body are the notification content.
	SendPush(ctx context.Context, deviceTokens []string, title, body string) error
}
