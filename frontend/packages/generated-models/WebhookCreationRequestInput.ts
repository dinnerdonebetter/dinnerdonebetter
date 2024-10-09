// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IWebhookCreationRequestInput {
  method: string;
  name: string;
  url: string;
  contentType: string;
  events: string;
}

export class WebhookCreationRequestInput implements IWebhookCreationRequestInput {
  method: string;
  name: string;
  url: string;
  contentType: string;
  events: string;
  constructor(input: Partial<WebhookCreationRequestInput> = {}) {
    this.method = input.method = '';
    this.name = input.name = '';
    this.url = input.url = '';
    this.contentType = input.contentType = '';
    this.events = input.events = '';
  }
}
