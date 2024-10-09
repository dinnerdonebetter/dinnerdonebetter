// GENERATED CODE, DO NOT EDIT MANUALLY

import { WebhookTriggerEvent } from './WebhookTriggerEvent';

export interface IWebhook {
  archivedAt?: string;
  belongsToHousehold: string;
  createdAt: string;
  events: WebhookTriggerEvent;
  method: string;
  name: string;
  contentType: string;
  id: string;
  lastUpdatedAt?: string;
  url: string;
}

export class Webhook implements IWebhook {
  archivedAt?: string;
  belongsToHousehold: string;
  createdAt: string;
  events: WebhookTriggerEvent;
  method: string;
  name: string;
  contentType: string;
  id: string;
  lastUpdatedAt?: string;
  url: string;
  constructor(input: Partial<Webhook> = {}) {
    this.archivedAt = input.archivedAt;
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.createdAt = input.createdAt = '';
    this.events = input.events = new WebhookTriggerEvent();
    this.method = input.method = '';
    this.name = input.name = '';
    this.contentType = input.contentType = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.url = input.url = '';
  }
}
