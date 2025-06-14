// GENERATED CODE, DO NOT EDIT MANUALLY

import { WebhookTriggerEvent } from './WebhookTriggerEvent.gen';

export interface IWebhook {
  archivedAt: string;
  belongsToAccount: string;
  contentType: string;
  createdAt: string;
  events: WebhookTriggerEvent[];
  id: string;
  lastUpdatedAt: string;
  method: string;
  name: string;
  url: string;
}

export class Webhook implements IWebhook {
  archivedAt: string;
  belongsToAccount: string;
  contentType: string;
  createdAt: string;
  events: WebhookTriggerEvent[];
  id: string;
  lastUpdatedAt: string;
  method: string;
  name: string;
  url: string;
  constructor(input: Partial<Webhook> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToAccount = input.belongsToAccount || '';
    this.contentType = input.contentType || '';
    this.createdAt = input.createdAt || '';
    this.events = input.events || [];
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.method = input.method || '';
    this.name = input.name || '';
    this.url = input.url || '';
  }
}
