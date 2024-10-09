// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IWebhookCreationRequestInput {
   name: string;
 url: string;
 contentType: string;
 events: string;
 method: string;

}

export class WebhookCreationRequestInput implements IWebhookCreationRequestInput {
   name: string;
 url: string;
 contentType: string;
 events: string;
 method: string;
constructor(input: Partial<WebhookCreationRequestInput> = {}) {
	 this.name = input.name = '';
 this.url = input.url = '';
 this.contentType = input.contentType = '';
 this.events = input.events = '';
 this.method = input.method = '';
}
}