package notifications

import "context"

// PushNotificationSender sends push notifications to device tokens.
// Implementations route by platform: APNs for ios, FCM for android.
type PushNotificationSender interface {
	// SendPush sends a push notification to a single device token.
	// platform is "ios" or "android"; implementations filter by platform.
	SendPush(ctx context.Context, platform, token, title, body string) error
}
