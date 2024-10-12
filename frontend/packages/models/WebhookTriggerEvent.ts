// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IWebhookTriggerEvent {
   archivedAt: string;
 belongsToWebhook: string;
 createdAt: string;
 id: string;
 triggerEvent: string;

}

export class WebhookTriggerEvent implements IWebhookTriggerEvent {
   archivedAt: string;
 belongsToWebhook: string;
 createdAt: string;
 id: string;
 triggerEvent: string;
constructor(input: Partial<WebhookTriggerEvent> = {}) {
	 this.archivedAt = input.archivedAt || '';
 this.belongsToWebhook = input.belongsToWebhook || '';
 this.createdAt = input.createdAt || '';
 this.id = input.id || '';
 this.triggerEvent = input.triggerEvent || '';
}
}