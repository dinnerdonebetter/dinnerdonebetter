// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IWebhookTriggerEvent {
  id: string;
  triggerEvent: string;
  archivedAt?: string;
  belongsToWebhook: string;
  createdAt: string;
}

export class WebhookTriggerEvent implements IWebhookTriggerEvent {
  id: string;
  triggerEvent: string;
  archivedAt?: string;
  belongsToWebhook: string;
  createdAt: string;
  constructor(input: Partial<WebhookTriggerEvent> = {}) {
    this.id = input.id = '';
    this.triggerEvent = input.triggerEvent = '';
    this.archivedAt = input.archivedAt;
    this.belongsToWebhook = input.belongsToWebhook = '';
    this.createdAt = input.createdAt = '';
  }
}
