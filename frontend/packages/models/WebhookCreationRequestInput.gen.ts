// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IWebhookCreationRequestInput {
  contentType: string;
  events: string[];
  method: string;
  name: string;
  url: string;
}

export class WebhookCreationRequestInput implements IWebhookCreationRequestInput {
  contentType: string;
  events: string[];
  method: string;
  name: string;
  url: string;
  constructor(input: Partial<WebhookCreationRequestInput> = {}) {
    this.contentType = input.contentType || '';
    this.events = input.events || [];
    this.method = input.method || '';
    this.name = input.name || '';
    this.url = input.url || '';
  }
}
