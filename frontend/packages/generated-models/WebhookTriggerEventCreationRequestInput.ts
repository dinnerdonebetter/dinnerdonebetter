// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IWebhookTriggerEventCreationRequestInput {
  belongsToWebhook: string;
  triggerEvent: string;
}

export class WebhookTriggerEventCreationRequestInput implements IWebhookTriggerEventCreationRequestInput {
  belongsToWebhook: string;
  triggerEvent: string;
  constructor(input: Partial<WebhookTriggerEventCreationRequestInput> = {}) {
    this.belongsToWebhook = input.belongsToWebhook = '';
    this.triggerEvent = input.triggerEvent = '';
  }
}
