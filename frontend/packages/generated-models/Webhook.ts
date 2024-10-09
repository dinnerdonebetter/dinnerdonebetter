// GENERATED CODE, DO NOT EDIT MANUALLY

import { WebhookTriggerEvent } from './WebhookTriggerEvent';

export interface IWebhook {
  archivedAt?: string;
  events: WebhookTriggerEvent;
  id: string;
  name: string;
  url: string;
  belongsToHousehold: string;
  contentType: string;
  createdAt: string;
  lastUpdatedAt?: string;
  method: string;
}

export class Webhook implements IWebhook {
  archivedAt?: string;
  events: WebhookTriggerEvent;
  id: string;
  name: string;
  url: string;
  belongsToHousehold: string;
  contentType: string;
  createdAt: string;
  lastUpdatedAt?: string;
  method: string;
  constructor(input: Partial<Webhook> = {}) {
    this.archivedAt = input.archivedAt;
    this.events = input.events = new WebhookTriggerEvent();
    this.id = input.id = '';
    this.name = input.name = '';
    this.url = input.url = '';
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.contentType = input.contentType = '';
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.method = input.method = '';
  }
}
