package datachangemessagehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/keys"
	domainnotifications "github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
)

// MobileNotificationsEventHandler handles mobile notification requests from the mobile_notifications queue.
// It decodes the request, validates it, and routes to the appropriate type-specific handler.
func (a *AsyncDataChangeMessageHandler) MobileNotificationsEventHandler(ctx context.Context, rawMsg []byte) error {
	ctx, span := a.tracer.StartSpan(ctx)
	defer span.End()

	var req notifications.MobileNotificationRequest
	if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&req); err != nil {
		return fmt.Errorf("decoding mobile notification request: %w", err)
	}
	if req.Title == "" || req.Body == "" {
		return fmt.Errorf("title and body are required")
	}
	if req.RequestType == "" {
		return fmt.Errorf("request type is required")
	}

	switch req.RequestType {
	case notifications.MobileNotificationRequestTypeMealPlanTask:
		return a.handleMealPlanTaskNotification(ctx, &req)
	default:
		return fmt.Errorf("unknown request type: %q", req.RequestType)
	}
}

// handleMealPlanTaskNotification processes a meal plan task reminder notification.
// It performs idempotency checks, sends push notifications to recipients, and marks the task as notified.
func (a *AsyncDataChangeMessageHandler) handleMealPlanTaskNotification(ctx context.Context, req *notifications.MobileNotificationRequest) error {
	ctx, span := a.tracer.StartSpan(ctx)
	defer span.End()

	mealPlanTaskID := ""
	if req.Context != nil {
		mealPlanTaskID = req.Context[notifications.MealPlanTaskIDContextKey]
	}
	if mealPlanTaskID == "" {
		return fmt.Errorf("meal plan task notification requires mealPlanTaskID in context")
	}

	logger := a.logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)

	sent, err := a.mealPlanRepo.MealPlanTaskNotificationHasBeenSent(ctx, mealPlanTaskID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "checking if notification already sent")
	}
	if sent {
		return nil // idempotent skip
	}

	if len(req.RecipientUserIDs) == 0 {
		return a.mealPlanRepo.MarkMealPlanTaskNotificationSent(ctx, mealPlanTaskID)
	}

	deviceTokens, err := a.collectDeviceTokensForUsers(ctx, req.RecipientUserIDs)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "collecting device tokens")
	}
	if len(deviceTokens) == 0 {
		return a.mealPlanRepo.MarkMealPlanTaskNotificationSent(ctx, mealPlanTaskID)
	}

	atLeastOneSent := false
	for _, t := range deviceTokens {
		if sendErr := a.pushNotificationSender.SendPush(ctx, t.Platform, t.DeviceToken, req.Title, req.Body); sendErr != nil {
			logger.WithValue("user_device_token_id", t.ID).WithValue("error", sendErr).Error("sending push notification to device", sendErr)
			if strings.Contains(sendErr.Error(), "BadDeviceToken") {
				if archiveErr := a.notificationsRepo.ArchiveUserDeviceToken(ctx, t.BelongsToUser, t.ID); archiveErr != nil {
					logger.WithValue("user_device_token_id", t.ID).Error("archiving invalid device token", archiveErr)
				} else {
					logger.WithValue("user_device_token_id", t.ID).Info("archived invalid device token (BadDeviceToken from APNs)")
				}
			}
			// Continue to other tokens; don't fail entire batch
		} else {
			atLeastOneSent = true
		}
	}

	if !atLeastOneSent {
		return nil // Don't mark as notified if every attempt failed; scheduler will retry
	}
	return a.mealPlanRepo.MarkMealPlanTaskNotificationSent(ctx, mealPlanTaskID)
}

func (a *AsyncDataChangeMessageHandler) collectDeviceTokensForUsers(ctx context.Context, userIDs []string) ([]*domainnotifications.UserDeviceToken, error) {
	var tokens []*domainnotifications.UserDeviceToken
	filter := filtering.DefaultQueryFilter()
	for _, userID := range userIDs {
		result, err := a.notificationsRepo.GetUserDeviceTokens(ctx, userID, filter, nil)
		if err != nil {
			return nil, fmt.Errorf("getting device tokens for user %s: %w", userID, err)
		}
		for _, t := range result.Data {
			if t != nil && t.DeviceToken != "" {
				tokens = append(tokens, t)
			}
		}
	}
	return tokens, nil
}
