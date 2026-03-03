package datachangemessagehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/keys"
	domainnotifications "github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
)

// MobileNotificationsEventHandler handles meal plan task notification requests from the mobile_notifications queue.
func (a *AsyncDataChangeMessageHandler) MobileNotificationsEventHandler(ctx context.Context, rawMsg []byte) error {
	ctx, span := a.tracer.StartSpan(ctx)
	defer span.End()

	// TODO: genericize this

	var req notifications.MealPlanTaskNotificationRequest
	if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&req); err != nil {
		return fmt.Errorf("decoding meal plan task notification request: %w", err)
	}
	if req.MealPlanTaskID == "" {
		return fmt.Errorf("meal plan task ID is required")
	}

	logger := a.logger.WithValue(keys.MealPlanTaskIDKey, req.MealPlanTaskID)

	sent, err := a.mealPlanRepo.MealPlanTaskNotificationHasBeenSent(ctx, req.MealPlanTaskID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "checking if notification already sent")
	}
	if sent {
		return nil // idempotent skip
	}

	task, err := a.mealPlanRepo.GetMealPlanTask(ctx, req.MealPlanTaskID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching meal plan task")
	}
	if task == nil {
		return fmt.Errorf("meal plan task %s not found", req.MealPlanTaskID)
	}

	recipientUserIDs, err := a.resolveNotificationRecipients(ctx, task)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "resolving notification recipients")
	}
	if len(recipientUserIDs) == 0 {
		logger.WithValue("mealPlanTaskID", req.MealPlanTaskID).Info("no recipients for meal plan task notification, marking as sent")
		return a.mealPlanRepo.MarkMealPlanTaskNotificationSent(ctx, req.MealPlanTaskID)
	}

	deviceTokens, err := a.collectDeviceTokensForUsers(ctx, recipientUserIDs)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "collecting device tokens")
	}
	if len(deviceTokens) == 0 {
		logger.Info("no device tokens for recipients, marking as sent")
		return a.mealPlanRepo.MarkMealPlanTaskNotificationSent(ctx, req.MealPlanTaskID)
	}

	title, body := buildNotificationContent(task)
	for _, t := range deviceTokens {
		if err = a.pushNotificationSender.SendPush(ctx, t.Platform, t.DeviceToken, title, body); err != nil {
			logger.WithValue("user_device_token_id", t.ID).WithValue("error", err).Error("sending push notification to device", err)
			if strings.Contains(err.Error(), "BadDeviceToken") {
				if archiveErr := a.notificationsRepo.ArchiveUserDeviceToken(ctx, t.BelongsToUser, t.ID); archiveErr != nil {
					logger.WithValue("user_device_token_id", t.ID).Error("archiving invalid device token", archiveErr)
				} else {
					logger.WithValue("user_device_token_id", t.ID).Info("archived invalid device token (BadDeviceToken from APNs)")
				}
			}
			// Continue to other tokens; don't fail entire batch
		}
	}

	return a.mealPlanRepo.MarkMealPlanTaskNotificationSent(ctx, req.MealPlanTaskID)
}

func (a *AsyncDataChangeMessageHandler) resolveNotificationRecipients(ctx context.Context, task *mealplanning.MealPlanTask) ([]string, error) {
	if task.AssignedToUser != nil && *task.AssignedToUser != "" {
		return []string{*task.AssignedToUser}, nil
	}

	accountID, err := a.mealPlanRepo.GetMealPlanTaskAccountID(ctx, task.ID)
	if err != nil {
		return nil, fmt.Errorf("getting account ID for task: %w", err)
	}
	if accountID == "" {
		return nil, fmt.Errorf("task has no account")
	}

	usersResult, err := a.identityRepo.GetUsersForAccount(ctx, accountID, filtering.DefaultQueryFilter())
	if err != nil {
		return nil, fmt.Errorf("getting users for account: %w", err)
	}
	userIDs := make([]string, 0, len(usersResult.Data))
	for _, u := range usersResult.Data {
		if u != nil && u.ID != "" {
			userIDs = append(userIDs, u.ID)
		}
	}
	return userIDs, nil
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

func buildNotificationContent(task *mealplanning.MealPlanTask) (title, body string) {
	recipeName := task.RecipePrepTask.Name
	if recipeName == "" {
		recipeName = "A task"
	}
	return "Meal plan task", fmt.Sprintf("%s needs your attention", recipeName)
}
