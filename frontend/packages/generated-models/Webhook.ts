// GENERATED CODE, DO NOT EDIT MANUALLY

import { WebhookTriggerEvent } from './WebhookTriggerEvent';

export interface IWebhook {
  belongsToHousehold: string;
  events: WebhookTriggerEvent;
  name: string;
  url: string;
  archivedAt?: string;
  contentType: string;
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  method: string;
}

export class Webhook implements IWebhook {
  belongsToHousehold: string;
  events: WebhookTriggerEvent;
  name: string;
  url: string;
  archivedAt?: string;
  contentType: string;
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  method: string;
  constructor(input: Partial<Webhook> = {}) {
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.events = input.events = new WebhookTriggerEvent();
    this.name = input.name = '';
    this.url = input.url = '';
    this.archivedAt = input.archivedAt;
    this.contentType = input.contentType = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.method = input.method = '';
  }
}
