// GENERATED CODE, DO NOT EDIT MANUALLY

 import { WebhookTriggerEvent } from './WebhookTriggerEvent';


export interface IWebhook {
   url: string;
 contentType: string;
 createdAt: string;
 lastUpdatedAt?: string;
 name: string;
 method: string;
 archivedAt?: string;
 belongsToHousehold: string;
 events: WebhookTriggerEvent;
 id: string;

}

export class Webhook implements IWebhook {
   url: string;
 contentType: string;
 createdAt: string;
 lastUpdatedAt?: string;
 name: string;
 method: string;
 archivedAt?: string;
 belongsToHousehold: string;
 events: WebhookTriggerEvent;
 id: string;
constructor(input: Partial<Webhook> = {}) {
	 this.url = input.url = '';
 this.contentType = input.contentType = '';
 this.createdAt = input.createdAt = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.name = input.name = '';
 this.method = input.method = '';
 this.archivedAt = input.archivedAt;
 this.belongsToHousehold = input.belongsToHousehold = '';
 this.events = input.events = new WebhookTriggerEvent();
 this.id = input.id = '';
}
}