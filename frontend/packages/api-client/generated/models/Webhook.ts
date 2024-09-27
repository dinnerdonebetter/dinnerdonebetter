/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { WebhookTriggerEvent } from './WebhookTriggerEvent';
export type Webhook = {
  archivedAt?: string;
  belongsToHousehold?: string;
  contentType?: string;
  createdAt?: string;
  events?: Array<WebhookTriggerEvent>;
  id?: string;
  lastUpdatedAt?: string;
  method?: string;
  name?: string;
  url?: string;
};
