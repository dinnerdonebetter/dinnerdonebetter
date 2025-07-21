package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	configurationsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/configuration"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

func ConvertGPRCServiceSettingCreationRequestInputToServiceSettingDatabaseCreationInput(input *configurationsvc.ServiceSettingCreationRequestInput) *settings.ServiceSettingDatabaseCreationInput {
	return &settings.ServiceSettingDatabaseCreationInput{
		DefaultValue: &input.DefaultValue,
		Name:         input.Name,
		Type:         input.Type,
		Description:  input.Description,
		Enumeration:  input.Enumeration,
		AdminsOnly:   input.AdminsOnly,
	}
}

func ConvertServiceSettingToGRPCServiceSetting(input *settings.ServiceSetting) *configurationsvc.ServiceSetting {
	x := &configurationsvc.ServiceSetting{
		CreatedAt:     grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		ArchivedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		LastUpdatedAt: grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		Name:          input.Name,
		ID:            input.ID,
		Type:          input.Type,
		Description:   input.Description,
		Enumeration:   input.Enumeration,
		AdminsOnly:    input.AdminsOnly,
	}

	if input.DefaultValue != nil {
		x.DefaultValue = *input.DefaultValue
	}

	return x
}

func ConvertGRPCServiceSettingConfigurationCreationRequestInputToServiceSettingConfigurationDatabaseCreationInput(input *configurationsvc.ServiceSettingConfigurationCreationRequestInput) *settings.ServiceSettingConfigurationDatabaseCreationInput {
	return &settings.ServiceSettingConfigurationDatabaseCreationInput{
		ID:               identifiers.New(),
		Value:            input.Value,
		Notes:            input.Notes,
		ServiceSettingID: input.ServiceSettingID,
		BelongsToUser:    input.BelongsToUser,
		BelongsToAccount: input.BelongsToAccount,
	}
}

func ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(input *settings.ServiceSettingConfiguration) *configurationsvc.ServiceSettingConfiguration {
	return &configurationsvc.ServiceSettingConfiguration{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(input.CreatedAt),
		LastUpdatedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(input.LastUpdatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(input.ArchivedAt),
		ID:               input.ID,
		Value:            input.Value,
		Notes:            input.Notes,
		BelongsToUser:    input.BelongsToUser,
		BelongsToAccount: input.BelongsToAccount,
		ServiceSetting:   ConvertServiceSettingToGRPCServiceSetting(&input.ServiceSetting),
	}
}

func ConvertGRPCServiceSettingConfigurationUpdateRequestInputToServiceSettingConfigurationUpdateRequestInputTo(input *configurationsvc.ServiceSettingConfigurationUpdateRequestInput) *settings.ServiceSettingConfigurationUpdateRequestInput {
	return &settings.ServiceSettingConfigurationUpdateRequestInput{
		Value:            input.Value,
		Notes:            input.Notes,
		ServiceSettingID: input.ServiceSettingID,
		BelongsToUser:    input.BelongsToUser,
		BelongsToAccount: input.BelongsToAccount,
	}
}

func ConvertWebhookToGRPCWebhook(webhook *webhooks.Webhook) *configurationsvc.Webhook {
	converted := &configurationsvc.Webhook{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(webhook.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(webhook.ArchivedAt),
		LastUpdatedAt:    grpcconverters.ConvertTimePointerToPBTimestamp(webhook.LastUpdatedAt),
		Name:             webhook.Name,
		URL:              webhook.URL,
		Method:           webhook.Method,
		ID:               webhook.ID,
		BelongsToAccount: webhook.BelongsToAccount,
		ContentType:      webhook.ContentType,
	}

	for _, event := range webhook.Events {
		converted.Events = append(converted.Events, ConvertWebhookTriggerEventToGRPCWebhookTriggerEvent(event))
	}

	return converted
}

func ConvertWebhookTriggerEventToGRPCWebhookTriggerEvent(z *webhooks.WebhookTriggerEvent) *configurationsvc.WebhookTriggerEvent {
	return &configurationsvc.WebhookTriggerEvent{
		CreatedAt:        grpcconverters.ConvertTimeToPBTimestamp(z.CreatedAt),
		ArchivedAt:       grpcconverters.ConvertTimePointerToPBTimestamp(z.ArchivedAt),
		ID:               z.ID,
		BelongsToWebhook: z.BelongsToWebhook,
		TriggerEvent:     z.TriggerEvent,
	}
}

func ConvertGRPCWebhookCreationRequestInputToWebhookDatabaseCreationInput(input *configurationsvc.WebhookCreationRequestInput, accountID string) *webhooks.WebhookDatabaseCreationInput {
	x := &webhooks.WebhookDatabaseCreationInput{
		ID:               identifiers.New(),
		Name:             input.Name,
		ContentType:      input.ContentType,
		URL:              input.URL,
		Method:           input.Method,
		BelongsToAccount: accountID,
	}

	for _, event := range input.Events {
		x.Events = append(x.Events, &webhooks.WebhookTriggerEventDatabaseCreationInput{
			ID:               identifiers.New(),
			BelongsToWebhook: x.ID,
			TriggerEvent:     event,
		})
	}

	return x
}

func ConvertGRPCWebhookTriggerEventDatabaseCreationInputToWebhookTriggerEventDatabaseCreationInput(input *configurationsvc.WebhookTriggerEventCreationRequestInput) *webhooks.WebhookTriggerEventDatabaseCreationInput {
	return &webhooks.WebhookTriggerEventDatabaseCreationInput{
		ID:               identifiers.New(),
		BelongsToWebhook: input.BelongsToWebhook,
		TriggerEvent:     input.TriggerEvent,
	}
}
