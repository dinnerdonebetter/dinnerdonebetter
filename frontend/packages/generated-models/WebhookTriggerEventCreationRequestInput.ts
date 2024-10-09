// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IWebhookTriggerEventCreationRequestInput {
   triggerEvent: string;
 belongsToWebhook: string;

}

export class WebhookTriggerEventCreationRequestInput implements IWebhookTriggerEventCreationRequestInput {
   triggerEvent: string;
 belongsToWebhook: string;
constructor(input: Partial<WebhookTriggerEventCreationRequestInput> = {}) {
	 this.triggerEvent = input.triggerEvent = '';
 this.belongsToWebhook = input.belongsToWebhook = '';
}
}