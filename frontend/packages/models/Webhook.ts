// GENERATED CODE, DO NOT EDIT MANUALLY

import { WebhookTriggerEvent } from './WebhookTriggerEvent';

export interface IWebhook {
  createdAt: string;
  events: WebhookTriggerEvent[];
  url: string;
  contentType: string;
  belongsToHousehold: string;
  id: string;
  lastUpdatedAt: string;
  method: string;
  name: string;
  archivedAt: string;
}

export class Webhook implements IWebhook {
  createdAt: string;
  events: WebhookTriggerEvent[];
  url: string;
  contentType: string;
  belongsToHousehold: string;
  id: string;
  lastUpdatedAt: string;
  method: string;
  name: string;
  archivedAt: string;
  constructor(input: Partial<Webhook> = {}) {
    this.createdAt = input.createdAt || '';
    this.events = input.events || [];
    this.url = input.url || '';
    this.contentType = input.contentType || '';
    this.belongsToHousehold = input.belongsToHousehold || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.method = input.method || '';
    this.name = input.name || '';
    this.archivedAt = input.archivedAt || '';
  }
}
